package rest

import (
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/presenters/rest/ui"
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
// func MapRoutes(fi *fiber.App, he *health.Service, m *managing.Service, i *importing.Service, t *tracking.Service, q *querying.Service) {
func MapRoutes(fi *fiber.App, he *health.Service, m *managing.Service, t *tracking.Service, q *querying.Service) {
	// Unrestricted
	fi.Static("/assets", "./public/assets")
	fi.Get("/", ui.Home)
	fi.Get("/expenses", ui.ExpenseSection(q.ExpenseQuerier))
	fi.Get("/expenses/:id/edit", ui.ExpenseRowEdit(q.ExpenseQuerier, q.CategoryQuerier))
	// fi.Get("/expenses/:id/edit", func(c *fiber.Ctx) error {
	// 	fmt.Println(c.Params("id"))
	// 	return nil
	// })
	fi.Get("/importer", ui.Importer())
	fi.Get("/ping", ping(*he))
	// fi.Post("/login", login(*a))
	fi.Get("/api/expenses", getExpenses(q.ExpenseQuerier))

	// Restricted
	fi.Use(jwtware.New(jwtware.Config{SigningKey: []byte(os.Getenv("JWT_SECRET_SEED"))}))
	fi.Get("/restricted", restricted)
	// fi.Post("/users", createUsers(m.UserCreator))
	fi.Get("/api/expenses/categories", getCategories(q.CategoryQuerier))
	fi.Post("/api/expenses/categories", createCategory(m.CategoryCreator))
	// fi.Post("/importers/:id", importExpenses(i.ImportExpenses))
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
