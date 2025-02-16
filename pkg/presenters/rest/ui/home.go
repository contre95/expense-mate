package ui

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

const DEFAULT_DAYS_FROM_PARAM = "190"
const DEFAULT_DAYS_TO_PARAM = "0" // Now
const DEFAULT_PSIZE_PARAM = "30"
const DEFAULT_PNUM_PARAM = "0"

// Home hanlder reders the homescreen
func Home(c *fiber.Ctx) error {
	slog.Info("HOME")
	slog.Info(os.Getenv("IMAGE_TAG"))
	// render index template
	c.Append("Hx-Trigger", "expensesTable")
	return c.Render("main", fiber.Map{
		"DashboardTrigger": "revealed",
		"MateVersion":      os.Getenv("IMAGE_TAG"),
	})
}

// Home hanlder reders the homescreen
func Empty() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("")
	}
}
