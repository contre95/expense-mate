package rest

import (
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/presenters/rest/api"
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
	// UI
	fi.Static("/assets", "./public/assets")
	fi.Get("/", ui.Home)
	fi.Get("/expenses", ui.LoadExpensesSection())
	fi.Post("/expenses", ui.CreateExpense(t.ExpenseCreator))
	fi.Get("/expenses/filter", ui.LoadExpenseFilter(q.CategoryQuerier))
	fi.Get("/expenses/table", ui.LoadExpensesTable(q.ExpenseQuerier))
	fi.Get("/expenses/add", ui.LoadExpensesAddRow(q.CategoryQuerier))
	fi.Get("/expenses/:id/edit", ui.LoadExpenseEditRow(q.ExpenseQuerier, q.CategoryQuerier))
	fi.Get("/expenses/:id/row", ui.LoadExpenseRow(q.ExpenseQuerier, q.CategoryQuerier))
	fi.Put("/expenses/:id", ui.EditExpense(q.ExpenseQuerier, t.ExpenseUpdater))
	fi.Get("/empty", ui.Empty())
	fi.Delete("/expenses/:id", ui.DeleteExpense(t.ExpenseDeleter))
	fi.Get("/importers", ui.LoadImporterSection())
	fi.Get("/importers/n26", ui.LoadN26Importer())
	fi.Get("/importers/table", ui.LoadImportersTable(q.ExpenseQuerier, q.CategoryQuerier))
	fi.Post("/importers/n26", ui.ImportN26CSV(t.ExpenseCreator))
	fi.Get("/importers/revolut", ui.LoadRevolutImporter())

	// Restricted endpoints below
	fi.Use(jwtware.New(jwtware.Config{SigningKey: []byte(os.Getenv("JWT_SECRET_SEED"))}))

	// fi.Post("/login", login(*a))
	fi.Get("/api/ping", api.Ping(*he))
	fi.Get("/api/expenses", api.GetExpenses(q.ExpenseQuerier))
	fi.Get("/restricted", restricted)
	// fi.Post("/users", createUsers(m.UserCreator))
	fi.Get("/api/expenses/categories", api.GetCategories(q.CategoryQuerier))
	fi.Post("/api/expenses/categories", api.CreateCategory(m.CategoryCreator))
	// fi.Post("/importers/:id", api.ImportExpenses(i.ImportExpenses))
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
