package routes

import (
	"encoding/json"
	"fmt"
	"os"

	"FoxxoOS/files"
	"FoxxoOS/util"

	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
)

func User(c *fiber.Ctx) error {
	user := [3]string{
		c.Query("name"),
		c.Query("password"),
		c.Query("hostname"),
	}

	userJSON := fmt.Sprintf("%v", user)

	util.SetOnceSave("user.name", user[0])
	util.SetOnceSave("user.password", user[1])
	util.SetOnceSave("hostname", user[2])

	return c.SendString(userJSON)
}

func Save(c *fiber.Ctx) error {
	saveRead, err := os.ReadFile(files.FilesJSON[2])

	util.ErrorCheck(err)

	saveJSON := string(saveRead)

	if c.Query("done") != "ok" {
		return c.SendString("not ok")
	}

	return c.SendString(saveJSON)
}

func Lang(c *fiber.Ctx) error {
	langRead, err := os.ReadFile(files.FilesJSON[0])

	util.ErrorCheck(err)

	langJSON := string(langRead)
	lang := c.Query("lang")

	value := gjson.Get(langJSON, lang)

	util.SetOnceSave("lang", value.String())

	return c.SendString(value.String())
}

func Keyboard(c *fiber.Ctx) error {
	keyRead, err := os.ReadFile(files.FilesJSON[1])

	util.ErrorCheck(err)

	keyJSON := string(keyRead)
	key := c.Query("keyboard")

	value := gjson.Get(keyJSON, key)

	util.SetOnceSave("keyboard", value.String())

	return c.SendString(value.String())
}

type Time struct {
	Timezone []string
}

func Timezone(c *fiber.Ctx) error {
	timeRead, err := os.ReadFile(files.FilesJSON[3])

	util.ErrorCheck(err)

	var times Time

	err = json.Unmarshal(timeRead, &times)

	util.ErrorCheck(err)

	time := c.Query("time")

	if !util.StringInSlice(time, times.Timezone) {
		return c.SendString("no ok")
	}

	util.SetOnceSave("timezone", time)

	return c.SendString(time)
}

func DE(c *fiber.Ctx) error {
	DERead, err := os.ReadFile(files.FilesJSON[4])

	util.ErrorCheck(err)

	DEJSON := string(DERead)
	DE := c.Query("desktop")

	value := gjson.Get(DEJSON, DE)

	util.SetOnceSave("desktop", value.String())

	return c.SendString(value.String())
}

func Web(c *fiber.Ctx) error {
	read, err := os.ReadFile(files.FilesJSON[5])

	util.ErrorCheck(err)

	JSON := string(read)

	var lMap map[string]string
	json.Unmarshal(read, &lMap)

	array := []string{}
	for key, _ := range lMap {
		array = append(array, c.Query(key))
	}

	list := []string{}
	for i := 0; i < len(array); i++ {
		if array[i] != "" {
			list = append(list, gjson.Get(JSON, array[i]).String())
		}
	}

	util.SetMultiSave("webbrowser", list)

	return c.SendString(fmt.Sprintf("%v", list))
}

func Program(c *fiber.Ctx) error {
	read, err := os.ReadFile(files.FilesJSON[6])

	util.ErrorCheck(err)

	JSON := string(read)

	var lMap map[string]string
	json.Unmarshal(read, &lMap)

	array := []string{}
	for key, _ := range lMap {
		array = append(array, c.Query(key))
	}

	list := []string{}
	for i := 0; i < len(array); i++ {
		if array[i] != "" {
			list = append(list, gjson.Get(JSON, array[i]).String())
		}
	}

	util.SetMultiSave("programming", list)

	return c.SendString(fmt.Sprintf("%v", list))
}

func Office(c *fiber.Ctx) error {
	read, err := os.ReadFile(files.FilesJSON[7])

	util.ErrorCheck(err)

	JSON := string(read)

	var lMap map[string]string
	json.Unmarshal(read, &lMap)

	array := []string{}
	for key, _ := range lMap {
		array = append(array, c.Query(key))
	}

	list := []string{}
	for i := 0; i < len(array); i++ {
		if array[i] != "" {
			list = append(list, gjson.Get(JSON, array[i]).String())
		}
	}

	util.SetMultiSave("office", list)

	return c.SendString(fmt.Sprintf("%v", list))
}

func Gaming(c *fiber.Ctx) error {
	read, err := os.ReadFile(files.FilesJSON[8])

	util.ErrorCheck(err)

	JSON := string(read)

	var lMap map[string]string
	json.Unmarshal(read, &lMap)

	array := []string{}
	for key, _ := range lMap {
		array = append(array, c.Query(key))
	}

	list := []string{}
	for i := 0; i < len(array); i++ {
		if array[i] != "" {
			list = append(list, gjson.Get(JSON, array[i]).String())
		}
	}

	util.SetMultiSave("gaming", list)

	return c.SendString(fmt.Sprintf("%v", list))
}

func Utils(c *fiber.Ctx) error {
	read, err := os.ReadFile(files.FilesJSON[9])

	util.ErrorCheck(err)

	JSON := string(read)

	var lMap map[string]string
	json.Unmarshal(read, &lMap)

	array := []string{}
	for key, _ := range lMap {
		array = append(array, c.Query(key))
	}

	list := []string{}
	for i := 0; i < len(array); i++ {
		if array[i] != "" {
			list = append(list, gjson.Get(JSON, array[i]).String())
		}
	}

	util.SetMultiSave("utils", list)

	return c.SendString(fmt.Sprintf("%v", list))
}

func MediaGrap(c *fiber.Ctx) error {
	read, err := os.ReadFile(files.FilesJSON[10])

	util.ErrorCheck(err)

	JSON := string(read)

	var lMap map[string]string
	json.Unmarshal(read, &lMap)

	array := []string{}
	for key, _ := range lMap {
		array = append(array, c.Query(key))
	}

	list := []string{}
	for i := 0; i < len(array); i++ {
		if array[i] != "" {
			list = append(list, gjson.Get(JSON, array[i]).String())
		}
	}

	util.SetMultiSave("mediagrap", list)

	return c.SendString(fmt.Sprintf("%v", list))
}

func Drivers(c *fiber.Ctx) error {
	read, err := os.ReadFile(files.FilesJSON[11])

	util.ErrorCheck(err)

	JSON := string(read)

	var lMap map[string]string
	json.Unmarshal(read, &lMap)

	array := []string{}
	for key, _ := range lMap {
		array = append(array, c.Query(key))
	}

	list := []string{}
	for i := 0; i < len(array); i++ {
		if array[i] != "" {
			list = append(list, gjson.Get(JSON, array[i]).String())
		}
	}

	util.SetMultiSave("drivers", list)

	return c.SendString(fmt.Sprintf("%v", list))
}

type Disk struct {
	Type string
	Disk string
	Root string
	BootEFI string
	Swap string
}

func Partitions(c *fiber.Ctx) error {
	types := c.Query("type")

	disk := Disk{}
	disk.Type = types

	switch types {
	case "auto":
		disk.Disk = c.Query("disk")

		_, err := os.Stat("/sys/firmware/efi") 
		if err == nil {
			disk.BootEFI = fmt.Sprintf("%v%v", disk.Disk, 3)
		} 

		disk.Swap = fmt.Sprintf("%v%v", disk.Disk, 2)

		disk.Root = fmt.Sprintf("%v%v", disk.Disk, 1)
	case "manually":
		disk.Disk = c.Query("disk")

		_, err := os.Stat("/sys/firmware/efi") 
		if err == nil {
			disk.BootEFI = c.Query("boot")
		}

		disk.Swap = c.Query("swap")

		disk.Root = c.Query("root")
	}

	util.SetOnceSave("disk.type", disk.Type)
	util.SetOnceSave("disk.disk", disk.Disk)
	util.SetOnceSave("disk.swap", disk.Swap)

	_, err := os.Stat("/sys/firmware/efi") 
	if err == nil {
		util.SetOnceSave("disk.boot", disk.BootEFI)
	}

	util.SetOnceSave("disk.root", disk.Root)

	return c.SendString(fmt.Sprintf("%v", disk))
}