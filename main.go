package main

import (
	//install "FoxxoOS/installation"
	s "FoxxoOS/main_server"

	"github.com/gofiber/fiber/v2"
)

func main() {
	//install.Config()

	app := fiber.New(fiber.Config{
		AppName: "Foxxo OS",
	})

	s.MainServer(app)
}
