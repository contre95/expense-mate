package http

import (
	"expenses-app/pkg/app/authenticating"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/importing"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func Run(fi *fiber.App, port int) {
	fi.Listen(":" + strconv.Itoa(port))
}

// MapRoutes is where http REST routes are mapped to functions
func MapRoutes(fi *fiber.App, he *health.Service, m *managing.Service, i *importing.Service, a *authenticating.Service, q *querying.Service) {
	// Unrestricted
	fi.Get("/ping", ping(*he))
	fi.Post("/login", login(*a))
	fi.Use(jwtware.New(jwtware.Config{SigningKey: []byte(os.Getenv("JWT_SECRET_SEED"))}))
	// Restricted
	fi.Get("/restricted", restricted)
	fi.Post("/users", createUsers(m.UserCreator))
	fi.Get("/expenses/categories", getCategories(q.CategoryGetter))
	fi.Post("/expenses/categories", createCategory(m.CategoryCreator))
	fi.Post("/importers/:id", importExpenses(i.ImportExpenses))
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
