package ui

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// Home hanlder reders the homescreen
func Home(c *fiber.Ctx) error {
	slog.Info("HOME")
	// render index template
	c.Append("Hx-Trigger", "expensesTable")
	return c.Render("main", fiber.Map{
		"ExpensesTrigger": "revealed",
	})
}

// Home hanlder reders the homescreen
func Empty() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("")
	}
}
