package main

import (
	"database/sql"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/gateways/logger"
	"expenses-app/pkg/gateways/storage/sqlstorage"
	"expenses-app/pkg/presenters/rest"
	"expenses-app/pkg/presenters/rest/ui"
	"expenses-app/pkg/presenters/telegram"
	"os"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Infrastructure / Gateways

	//Loggers
	initLogger := logger.NewSTDLogger("INIT", logger.VIOLET)
	healthLogger := logger.NewSTDLogger("HEALTH", logger.GREEN2)
	// authLogger := logger.NewSTDLogger("Authenticator", logger.YELLOW)
	managerLogger := logger.NewSTDLogger("Managing", logger.CYAN)
	// importerLogger := logger.NewSTDLogger("Importing", logger.BEIGE)
	querierLogger := logger.NewSTDLogger("Querying", logger.YELLOW2)
	trackerLogger := logger.NewSTDLogger("Tracker", logger.CYAN)
	telergamLogger := logger.NewSTDLogger("TELEGRAM", logger.BLUE)

	// SQL Storage
	mysqlUser := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASS")
	mysqlUrl := "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DB") + "?parseTime=true"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := sql.Open("mysql", mysqlUser+mysqlUrl)
	defer db.Close()
	if err != nil {
		initLogger.Err("%v", err)
		return
	}
	sqlStorage := sqlstorage.NewStorage(db)
	initLogger.Info("SQL storage initializer on %s", mysqlUrl)

	// Authenticating
	// authenticator := authenticating.NewService()

	// Healthching
	healthChecker := health.NewService(healthLogger)

	// Querying
	getCategories := querying.NewCategoryQuerier(querierLogger, sqlStorage)
	getExpenses := querying.NewExpenseQuerier(querierLogger, sqlStorage)
	querier := querying.NewService(*getCategories, *getExpenses)

	// Importing
	// importExpenses := importing.NewExpenseImporter(importerLogger, sqlStorage)

	// Managing
	createCategory := managing.NewCategoryCreator(managerLogger, sqlStorage)
	deleteCategory := managing.NewCategoryDeleter(managerLogger, sqlStorage)
	manager := managing.NewService(*deleteCategory, *createCategory)
	// Tracking
	createExpense := tracking.NewExpenseCreator(trackerLogger, sqlStorage)
	updateExpense := tracking.NewExpenseUpdater(trackerLogger, sqlStorage)
	tracker := tracking.NewService(*createExpense, *updateExpense)

	// API
	engine := html.New("./views", ".html")
	engine.AddFunc("nameToColor", ui.NameToColor)
	engine.AddFunc("unescape", ui.Unescape)
	engine.Debug(true)

	fiberApp := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ALLOWLIST"),
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// rest.MapRoutes(fiberApp, &healthChecker, &tracker, &querier)
	rest.MapRoutes(fiberApp, &healthChecker, &manager, &tracker, &querier)
	rest.Run(fiberApp, 8080)

	// Telegram Bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		initLogger.Err("%v", err)
		return
	}

	botConfig := telegram.BotConfig{
		AllowedUsers: []string{"contre", "anouxx"},
		People:       []string{"Anoux", "Contre", "Anoux / Contre"},
		PeopleUsers:  map[string]string{"contre": "Contre", "anouxx": "Anoux"},
		AuthUsers:    []int64{527377846, 751504879},
	}
	tgbotapi.SetLogger(telergamLogger)
	initLogger.Info("Telegram %s running.", bot.Self.FirstName)
	telegram.Run(bot, botConfig, &healthChecker, &tracker, &querier)
}
