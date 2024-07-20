package rest

import (
	"expenses-app/pkg/app/analyzing"
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
func MapRoutes(fi *fiber.App, he *health.Service, m *managing.Service, t *tracking.Service, q *querying.Service, a *analyzing.Service) {
	// Home and others
	fi.Static("/assets", "./public/assets")
	fi.Get("/empty", ui.Empty())
	fi.Get("/", ui.Home)
	// Expenses
	fi.Get("/expenses", ui.LoadExpensesSection())
	fi.Post("/expenses", ui.CreateExpense(t.ExpenseCreator))
	fi.Delete("/expenses/:id", ui.DeleteExpense(t.ExpenseDeleter))
	fi.Get("/expenses/table", ui.LoadExpensesTable(q.ExpenseQuerier))
	fi.Put("/expenses/:id", ui.EditExpense(q.ExpenseQuerier, t.ExpenseUpdater))
	fi.Get("/expenses/add", ui.LoadAddExpensesRow(q.CategoryQuerier, m.UserManager))
	fi.Get("/expenses/filter", ui.LoadExpenseFilter(q.CategoryQuerier, m.UserManager))
	fi.Get("/expenses/:id/row", ui.LoadExpenseRow(q.ExpenseQuerier, q.CategoryQuerier))
	fi.Get("/expenses/:id/edit", ui.LoadExpenseEditRow(q.ExpenseQuerier, q.CategoryQuerier, m.UserManager))
	// Importers
	fi.Get("/importers/n26", ui.LoadN26Importer(m.UserManager))
	fi.Get("/importers", ui.LoadImporterSection())
	fi.Get("/importers/revolut", ui.LoadRevolutImporter())
	fi.Post("/importers/n26", ui.ImportN26CSV(t.ExpenseCreator, t.RuleApplier))
	fi.Get("/importers/table", ui.LoadImportersTable(q.ExpenseQuerier, q.CategoryQuerier, m.UserManager))
	fi.Get("/export/csv", ui.ExportCSV(q.ExpenseQuerier, q.CategoryQuerier))
	fi.Get("/export/json", ui.ExportJSON(q.ExpenseQuerier, q.CategoryQuerier))
	// Settings
	fi.Get("/settings", ui.LoadSettingsSection())
	// Users
	fi.Post("/settings/users", ui.CreateUser(m.UserManager))
	fi.Get("/settings/users", ui.LoadUsersConfig(m.UserManager))
	fi.Delete("/settings/users/:id", ui.DeleteUser(m.UserManager))
	// Categores
	fi.Post("/settings/categories", ui.CreateCategory(m.CategoryCreator))
	fi.Put("/settings/categories/:id", ui.EditCategory(m.CategoryUpdater))
	fi.Get("/settings/categories", ui.LoadCategoriesConfig(q.CategoryQuerier))
	fi.Delete("/settings/categories/:id", ui.DeleteCategory(m.CategoryDeleter))
	// Rules
	fi.Post("/settings/rules/", ui.CreateRule(m.RuleManager))
	fi.Delete("/settings/rules/:id", ui.DeleteRule(m.RuleManager))
	fi.Get("/settings/rules", ui.LoadRulesConfig(q.CategoryQuerier, m.RuleManager, m.UserManager))
	// Telegram
	fi.Get("/settings/telegram", ui.LoadTelegramConfig())
	// fi.Get("/settings/telegram/status", ui.LoadTelegramStatus(*he))
	fi.Post("/telegram/command", ui.SendTelegramCommand(m.TelegramCommander))
	fi.Get("/telegram/users", ui.GetTelegramUsers(m.TelegramCommander))
	fi.Get("/telegram/status", ui.GetTelegramStatus(m.TelegramCommander))
	fi.Get("/dashboard/categories/summary", ui.LoadCategorySummaryTable(a.ExpenseAnalyzer))
	fi.Get("/dashboard/table/mini", ui.LoadExpensesMiniTable(q.ExpenseQuerier))
	fi.Get("/dashboard", ui.LoadDashboardSection())

	fi.Get("/api/health/app", api.Ping(*he))
	// fi.Get("/api/health/bot", api.BotPing(*he))
	// Restricted endpoints below
	fi.Use(jwtware.New(jwtware.Config{SigningKey: []byte(os.Getenv("JWT_SECRET_SEED"))}))

	// fi.Post("/login", login(*a))
	// fi.Post("/users", createUsers(m.UserCreator))
	// fi.Post("/importers/:id", api.ImportExpenses(i.ImportExpenses))
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
