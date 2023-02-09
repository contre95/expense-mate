package main

import (
	"context"
	"database/sql"
	"expenses-app/pkg/app/authenticating"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/importing"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/gateways/hasher"
	"expenses-app/pkg/gateways/importers"
	"expenses-app/pkg/gateways/logger"
	"expenses-app/pkg/gateways/storage/json"
	"expenses-app/pkg/gateways/storage/sqlstorage"
	"expenses-app/pkg/presenters/http"
	"expenses-app/pkg/presenters/telegram"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	// Infrastructure / Gateways

	//Loggers
	initLogger := logger.NewSTDLogger("INIT", logger.VIOLET)
	healthLogger := logger.NewSTDLogger("HEALTH", logger.GREEN2)
	authLogger := logger.NewSTDLogger("Authenticator", logger.YELLOW)
	managerLogger := logger.NewSTDLogger("Managing", logger.CYAN)
	importerLogger := logger.NewSTDLogger("Importing", logger.BEIGE)
	querierLogger := logger.NewSTDLogger("Querying", logger.YELLOW2)
	trackerLogger := logger.NewSTDLogger("Tracker", logger.CYAN)
	telergamLogger := logger.NewSTDLogger("TELEGRAM", logger.BLUE)

	// JSON Storage
	jsonStorage := json.NewStorage(os.Getenv("JSON_STORAGE_PATH"))
	initLogger.Info("Json storage initializer on %s", os.Getenv("JSON_STORAGE_PATH"))

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

	// Example importer
	exampleImporter := importers.NewExampleImporter("example data")
	initLogger.Info("Example importer initialized")
	// Sheets importer
	sheetsRangeLength, _ := strconv.Atoi(os.Getenv("SHEETS_IMPORTER_RAGENLEN"))
	sheetsID := os.Getenv("SHEETS_IMPORTER_ID")
	sheetsPageRange := os.Getenv("SHEETS_IMPORTER_PAGERANGE")
	sheetsPath := os.Getenv("SHEETS_IMPORTER_SA_PATH")
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithServiceAccountFile(sheetsPath))
	if err != nil {
		initLogger.Err("Error intilizing Google sheets: %v", err)
		return
	}
	sheetsImporter := importers.NewSheetsImporter(srv, sheetsID, sheetsPageRange, sheetsRangeLength)
	initLogger.Info("Sheets importer importer initialized for page range %s and range length %s", sheetsPageRange, sheetsRangeLength)

	// Importers
	importers := map[string]importing.Importer{
		"example": exampleImporter,
		"sheets":  sheetsImporter,
	}

	// Hashers
	passHasher := hasher.NewPasswordHasher()

	// Healthching
	healthChecker := health.NewService(healthLogger)

	// Querying
	getCategories := querying.NewCategoryQuerier(querierLogger, sqlStorage)
	getExpenses := querying.NewExpenseQuerier(querierLogger, sqlStorage)
	querier := querying.NewService(*getCategories, *getExpenses)

	// Tracking
	createExpense := tracking.NewExpenseCreator(trackerLogger, sqlStorage)
	tracker := tracking.NewService(*createExpense)

	// Managing
	createUser := managing.NewUserCreator(managerLogger, passHasher, jsonStorage)
	createCategory := managing.NewCategoryCreator(managerLogger, sqlStorage)
	manager := managing.NewService(*createUser, *createCategory)

	// Importing
	importExpenses := importing.NewExpenseImporter(importerLogger, importers, sqlStorage)
	importer := importing.NewService(*importExpenses)

	// Authenticating
	authenticateUser := authenticating.NewUserAuthenticator(authLogger, passHasher, jsonStorage)
	authenticator := authenticating.NewAuthenticator(*authenticateUser)

	// API
	fiberApp := fiber.New()
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ALLOWLIST"),
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	http.MapRoutes(fiberApp, &healthChecker, &manager, &importer, &authenticator, &querier)
	initLogger.Info("Startin fiber server on port  %d", 3000)
	for _, route := range fiberApp.GetRoutes(true) {
		initLogger.Info("Route %s and method %s registered.", route.Path, route.Params)
	}
	go http.Run(fiberApp, 3000)

	// Telegram Bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	tgbotapi.SetLogger(telergamLogger)
	initLogger.Info("Telegram %s running.", bot.Self.FirstName)
	telegram.Run(bot, &healthChecker, &manager, &tracker, &authenticator, &querier)
}
