package installation

import (
	"FoxxoOS/files"
	"FoxxoOS/util"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Installation() {
	var time time.Time

	fmt.Println("Startng installation...\n\n")

	fmt.Println("Partitioning...")
	util.StartTime(&time)
	parts := Partitioning()
	util.EndTime(time, "Partitioning")

	fmt.Println("Formatting...")
	util.StartTime(&time)
	Formating(parts)
	util.EndTime(time, "Formatting")

	fmt.Println("Mounting...")
	util.StartTime(&time)
	Mounting(parts)
	util.EndTime(time, "Mounting")

	fmt.Println("Building nix files...")
	util.StartTime(&time)
	Config()
	util.EndTime(time, "Building nix files")

	fmt.Println("Installation...")
	util.StartTime(&time)

	command := exec.Command("bash", "-c", "sudo nixos-install --no-root-passwd")
	command.Stderr = os.Stderr
	err := command.Run()

	util.ErrorCheck(err)
	util.EndTime(time, "Installation")

	fmt.Println("Configuring...")
	util.StartTime(&time)
	Chroot()
	util.EndTime(time, "Configuring")

	fmt.Println("Umounting...")
	UMounting()
	util.EndTime(time, "Umounting")

	Restart()
}

type Partitions struct {
	Disk string
	Root string
	Swap string
	Boot string
}

func partAuto(parts *Partitions, diskInfo map[string]string) {
	_, err := os.Stat("/sys/firmware/efi/efivars")
	fmt.Println(diskInfo)

	rootStart := "512M"

	if err == nil {
		util.Partitioning(
			diskInfo["disk"],
			"mklabel",
			[]string{"gpt"},
			[]string{},
		)
	} else {
		util.Partitioning(
			diskInfo["disk"],
			"mklabel",
			[]string{"msdos"},
			[]string{},
		)
	}
	parts.Disk = diskInfo["disk"]

	partitionRoot := util.Partitioning(
		diskInfo["disk"],
		"mkpart",
		[]string{"primary"},
		[]string{rootStart, "-4G"},
		1,
	)
	parts.Root = partitionRoot

	partitionSwap := util.Partitioning(
		diskInfo["disk"],
		"mkpart",
		[]string{"primary", "linux-swap"},
		[]string{"-4G", "100%"},
		2,
	)
	parts.Swap = partitionSwap

	if err == nil {
		partitionBoot := util.Partitioning(
			diskInfo["disk"],
			"mkpart",
			[]string{"ESP", "fat32"},
			[]string{"1M", rootStart},
			3,
		)
		parts.Boot = partitionBoot
	}
}

func partManual(parts *Partitions, diskInfo map[string]string) {
	_, err := os.Stat("/sys/firmware/efi/efivars")
	if err == nil {
		parts.Boot = diskInfo["boot"]
	}

	parts.Root = diskInfo["root"]
	parts.Swap = diskInfo["swap"]
}

func Formating(parts Partitions) {
	util.FormatFS("fs.btrfs", parts.Root)
	util.FormatFS("swap", parts.Swap)

	_, err := os.Stat("/sys/firmware/efi/efivars")
	if err == nil {
		util.FormatFS("fs.fat -F32", parts.Boot)
	}
}

func Mounting(parts Partitions) {
	util.Mount(parts.Root, "/mnt")

	_, err := os.Stat("/sys/firmware/efi/efivars")
	if err == nil {
		util.SudoExec("mkdir /mnt/boot")
		util.SudoExec("mkdir /mnt/boot/efi")

		util.Mount(parts.Boot, "/mnt/boot/efi")
	}

	command := fmt.Sprintf("swapon %v", parts.Swap)
	cmd := exec.Command("sudo " + command)
	cmd.Run()
}

func UMounting() {
	_, err := os.Stat("/sys/firmware/efi/efivars")
	if err == nil {
		util.UMount("/mnt/boot/efi")
	}

	util.UMount("/mnt")
}

func Partitioning() Partitions {
	file, err := os.ReadFile(files.FilesJSON[2])
	util.ErrorCheck(err)

	var JSON map[string]map[string]string
	json.Unmarshal(file, &JSON)

	diskInfo := JSON["disk"]
	parts := Partitions{}

	switch diskInfo["type"] {
	case "auto":
		partAuto(&parts, diskInfo)
	case "manual":
		partManual(&parts, diskInfo)
	}

	return parts
}

func Config() {
	fileSave, err := os.ReadFile(files.FilesJSON[2])
	util.ErrorCheck(err)

	fileNIX, err := os.ReadFile(files.FilesNIX[0])
	util.ErrorCheck(err)

	var JSON map[string]interface{}
	err = json.Unmarshal(fileSave, &JSON)

	util.ErrorCheck(err)

	util.ReplaceFile(&fileNIX, "$keyboard", JSON["keyboard"])
	util.ReplaceFile(&fileNIX, "$locales", JSON["lang"])
	util.ReplaceFile(&fileNIX, "$timezone", JSON["timezone"])
	util.ReplaceFile(&fileNIX, "$hostname", JSON["hostname"])
	util.ReplaceFile(&fileNIX, "$printing", util.StringInSlice("printing", JSON["drivers"]))
	util.ReplaceFile(&fileNIX, "$touchpad", util.StringInSlice("touchpad", JSON["drivers"]))
	util.ReplaceFile(&fileNIX, "$wifi", util.StringInSlice("wifi", JSON["drivers"]))
	util.ReplaceFile(&fileNIX, "$user.name", util.GetString(JSON["user"], "name"))
	util.ReplaceFile(&fileNIX, "$desktop", JSON["desktop"])

	bootEfi := `boot.loader = {
		efi = {
		  canTouchEfiVariables = true;
		  efiSysMountPoint = "/boot/efi";
		};
		grub = {
		   efiSupport = true;
		   device = "nodev";
		};
	  };
	  `
	bootBIOS := fmt.Sprintf("boot.loader.grub.enable = true;\n  boot.loader.grub.version = 2;\n  boot.loader.grub.device = \"%v\";", util.GetString(JSON["disk"], "disk"))

	_, err = os.Stat("/sys/firmware/efi")
	if err == nil {
		util.ReplaceFile(&fileNIX, "$boot", bootEfi)
	} else {
		util.ReplaceFile(&fileNIX, "$boot", bootBIOS)
	}

	if util.StringInSlice("nvidia", JSON["drivers"]) {
		util.ReplaceFile(&fileNIX, "$nvidia", "services.xserver.videoDrivers = [ \"nvidia\" ];")
	} else {
		util.ReplaceFile(&fileNIX, "$nvidia", "")
	}

	if util.StringInSlice("bluetooth", JSON["drivers"]) {
		util.ReplaceFile(&fileNIX, "$bluetooth", "hardware.bluetooth.enable = true;")
	} else {
		util.ReplaceFile(&fileNIX, "$bluetooth", "")
	}

	if util.StringInSlice("blueman", JSON["drivers"]) {
		util.ReplaceFile(&fileNIX, "$blueman", "services.blueman.enable = true;")
	} else {
		util.ReplaceFile(&fileNIX, "$blueman", "")
	}

	if util.StringInSlice("scanner_hp", JSON["drivers"]) {
		util.ReplaceFile(&fileNIX, "$scanner.hp", "hardware.sane.extraBackends = [ pkgs.hplipWithPlugin ];")
	} else {
		util.ReplaceFile(&fileNIX, "$scanner.hp", "")
	}

	if util.StringInSlice("scanner_airscan", JSON["drivers"]) {
		util.ReplaceFile(&fileNIX, "$scanner.airscan", "hardware.sane.extraBackends = [ pkgs.sane-airscan ];")
	} else {
		util.ReplaceFile(&fileNIX, "$scanner.airscan", "")
	}

	if util.StringInSlice("scanner_epson", JSON["drivers"]) {
		util.ReplaceFile(&fileNIX, "$scanner.epson", "hardware.sane.extraBackends = [ pkgs.epkowa ]; \n hardware.sane.extraBackends = [ pkgs.utsushi ]; \n services.udev.packages = [ pkgs.utsushi ];")
	} else {
		util.ReplaceFile(&fileNIX, "$scanner.epson", "")
	}

	if util.StringInSlice("scanner_brother", JSON["drivers"]) {
		util.ReplaceFile(&fileNIX, "$scanner.brother", `imports = [ 
    	<nixpkgs/nixos/modules/services/hardware/sane_extra_backends/brscan4.nix>
    	./hardware-configuration.nix
	];
	hardware.sane.brscan4.enable = true;
		`)
	} else {
		util.ReplaceFile(&fileNIX, "$scanner.brother", "")
	}

	if util.StringInSlice("scanner_gimp", JSON["drivers"]) {
		util.ReplaceFile(&fileNIX, "$scanner.gimp", `nixpkgs.config.packageOverrides = pkgs: {
			xsaneGimp = pkgs.xsane.override { gimpSupport = true; }; 
		};`)
	} else {
		util.ReplaceFile(&fileNIX, "$scanner.gimp", "")
	}

	if util.StringInSlice("scanner", JSON["drivers"]) {
		util.ReplaceFile(&fileNIX, "$scanner", "hardware.sane.enable = true;")
	} else {
		util.ReplaceFile(&fileNIX, "$scanner", "")
	}

	util.ReplaceFile(&fileNIX, "$pkg.webbrowser", util.Stringing(JSON["webbrowser"], "\n  "))
	util.ReplaceFile(&fileNIX, "$pkg.programming", util.Stringing(JSON["programming"], "\n  "))
	util.ReplaceFile(&fileNIX, "$pkg.gaming", util.Stringing(JSON["gaming"], "\n  "))
	util.ReplaceFile(&fileNIX, "$pkg.utils", util.Stringing(JSON["utils"], "\n  "))
	util.ReplaceFile(&fileNIX, "$pkg.mediagrap", util.Stringing(JSON["mediagrap"], "\n  "))
	util.ReplaceFile(&fileNIX, "$pkg.office", util.Stringing(JSON["office"], "\n  "))

	util.SaveFile("nix/configuration.nix", fileNIX)

	util.SudoExec("nixos-generate-config --root /mnt")

	util.SudoExec("cp %v %v", "./nix/configuration.nix", "/mnt/etc/nixos/configuration.nix")
}

func Chroot() {
	file, err := os.ReadFile(files.FilesJSON[2])
	util.ErrorCheck(err)

	var JSON map[string]map[string]string
	json.Unmarshal(file, &JSON)

	userInfo := JSON["user"]

	util.Chroot("echo -e \"%v\n%v\" | passwd %v", userInfo["password"], userInfo["password"], userInfo["name"])
	util.Chroot("echo -e \"%v\n%v\" | passwd %v", userInfo["password"], userInfo["password"], "root")
}

func Restart() {
	util.Clean()
	fmt.Println("Restart in 20 seconds! \n Click CTRL+C to stop it")

	for i := 19; i >= -1; i-- {
		time.Sleep(1 * time.Second)
		util.Clean()
		fmt.Printf("Restart in %v seconds! \n Click CTRL+C to stop it", i)
	}

	util.SudoExec("reboot --no-wall")
}
