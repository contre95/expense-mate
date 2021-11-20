package http

import (
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/importing"
	"expenses-app/pkg/app/managing"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v4"
)

// MapRoutes is where http REST routes are mapped to functions
func MapRoutes(fi *fiber.App, he *health.Service, m *managing.Service, i *importing.Service) {
	log.Println("asdasd")
	api := fi.Group("/api")    // /api
	v1 := api.Group("/v1")     // /api/v1
	fi.Get("/ping", ping(*he)) // /api/v1/ping
	v1.Post("/importers/:id", importExpenses(i.ImportExpenses))
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

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
