package ui

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LoadN26Importer() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("sections/importers/n26", fiber.Map{})
	}
}

func LoadRevolutImporter() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("sections/importers/revolut", fiber.Map{})
	}
}

func LoadImporterSection() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			// c.Append("hx-trigger", "newPair")  // Not working :(
			return c.Render("main", fiber.Map{
				"ImporterTrigger": "revealed",
			})
		}
		return c.Render("sections/importers/index", fiber.Map{})
	}
}
