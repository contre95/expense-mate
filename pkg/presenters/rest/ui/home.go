package ui

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

const DEFAULT_DAYS_FROM_PARAM = "190"
const DEFAULT_DAYS_TO_PARAM = "0" // Now
const DEFAULT_PSIZE_PARAM = "4"
const DEFAULT_PNUM_PARAM = "0"

// Home hanlder reders the homescreen
func Home(c *fiber.Ctx) error {
	slog.Info("HOME")
	// render index template
	c.Append("Hx-Trigger", "expensesTable")
	return c.Render("main", fiber.Map{
		"DashboardTrigger": "revealed",
	})
}

// Home hanlder reders the homescreen
func Empty() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("")
	}
}
