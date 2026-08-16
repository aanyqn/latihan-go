package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	var port string
	app.Get("/", func(c *fiber.Ctx) error {
		port = "8080"
		return c.SendString("Server is running")
	})

	log.Fatal(app.Listen(port))
}