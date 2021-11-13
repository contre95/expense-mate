package http

import (
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/importing"
	"expenses-app/pkg/app/managing"

	"github.com/gofiber/fiber/v2"
)

// MapRoutes is where http REST routes are mapped to functions
func MapRoutes(fi *fiber.App, he *health.Service, m *managing.Service, i *importing.Service) {
	fi.Get("/ping", ping(*he)) // /api/v1/ping
	api := fi.Group("/api")    // /api
	v1 := api.Group("/v1")     // /api/v1
	v1.Post("/impoters/:id", importExpenses(i.ImportExpenses))
	v1.Delete("/categories/:id", deleteCategory(m.DeleteCategory))
	//v1.Get("/categories", listClients(*&c.))
}

func ping(h health.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"ping": h.Ping(),
		})
	}
}
