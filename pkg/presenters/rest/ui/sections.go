package ui

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LoadSettingsSection() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			return c.Render("main", fiber.Map{
				"SettingsTrigger": "revealed",
			})
		}
		return c.Render("sections/settings/index", fiber.Map{})
	}
}

func LoadExpensesSection() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			return c.Render("main", fiber.Map{
				"ExpensesTrigger": "revealed",
			})
		}
		return c.Render("sections/expenses/index", fiber.Map{})
	}
}

func LoadImporterSection() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			return c.Render("main", fiber.Map{
				"ImporterTrigger": "revealed",
			})
		}
		return c.Render("sections/importers/index", fiber.Map{})
	}
}
