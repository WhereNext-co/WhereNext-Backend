package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/:value", func(c *fiber.Ctx) error {
		return c.SendString("Value is: " + c.Params("value"))
	})

	app.Get("/api/*", func(c *fiber.Ctx) error {
		return c.SendString("API path is: " + c.Params("*"))
	})

	app.Listen(":3000")
}
