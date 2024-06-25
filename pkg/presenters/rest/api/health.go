package api

import (
	"expenses-app/pkg/app/health"

	"github.com/gofiber/fiber/v2"
)

func Ping(h health.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"ping": h.Ping(),
		})
	}
}
