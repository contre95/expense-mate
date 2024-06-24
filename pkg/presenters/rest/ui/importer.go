package ui

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Home hanlder reders the homescreen
func Importer() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			// c.Append("hx-trigger", "newPair")  // Not working :(
			return c.Render("main", fiber.Map{
				"ImporterTrigger": "revealed",
			})
		}
		return c.Render("sections/imports/index", fiber.Map{})
	}
}
