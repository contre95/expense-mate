package main

import (
	"database/sql"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/gateways/logger"
	"expenses-app/pkg/gateways/storage/jsonstorage"
	"expenses-app/pkg/gateways/storage/sqlstorage"
	"expenses-app/pkg/presenters/rest"
	"expenses-app/pkg/presenters/rest/ui"
	"expenses-app/pkg/presenters/telegram"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Infrastructure / Gateways

	//Loggers
	initLogger := logger.NewSTDLogger("INIT", logger.VIOLET)
	healthLogger := logger.NewSTDLogger("HEALTH", logger.GREEN2)
	managerLogger := logger.NewSTDLogger("Managing", logger.CYAN)
	querierLogger := logger.NewSTDLogger("Querying", logger.YELLOW2)
	trackerLogger := logger.NewSTDLogger("Tracker", logger.CYAN)
	telegramLogger := logger.NewSTDLogger("TELEGRAM", logger.BLUE)
	commanderLogger := logger.NewSTDLogger("TELEGRAM COMMANDER", logger.BLUE2)

	// SQL storage
	var err error
	var db *sql.DB
	defer db.Close()
	switch os.Getenv("STORAGE_ENGINE") {
	case "mysql":
		mysqlUser := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASS")
		mysqlUrl := "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DB") + "?parseTime=true"
		db, err = sql.Open("mysql", mysqlUser+mysqlUrl)
		defer db.Close()
		if err != nil {
			initLogger.Err("Error instanciating mysql: %v", err)
			return
		}
		statements := strings.Split(sqlstorage.MySQLTables, ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			_, err = db.Exec(stmt)
			if err != nil {
				initLogger.Err("Error creating mysql tables: %v", err)
			}
		}
		initLogger.Info("MySQL storage initialized on %s", mysqlUrl)
	case "sqlite":
		path := os.Getenv("SQLITE_PATH")
		db, err = sql.Open("sqlite3", path)
		if err != nil {
			initLogger.Err("Error instanciating sqlite3: %v", err)
			return
		}
		_, err = db.Exec(sqlstorage.SQLiteTables)
		if err != nil {
			initLogger.Err("Error creating sqlite tables: %v", err)
			return
		}
		initLogger.Info("SQLte storage initialized on %s", path)
	case "":
		initLogger.Err("No storage set. Please set STORAGE_ENGINE variabel")
	}

	expensesStorage := sqlstorage.NewExpensesStorage(db)
	path := os.Getenv("JSON_STORAGE_PATH")
	ruleStorage := sqlstorage.NewRulesStorage(db)
	if path == "" {
		initLogger.Err("No storage set. Please set STORAGE_ENGINE variabel")
		return
	}
	userStorage := jsonstorage.NewStorage(path)

	// Healthching
	var botRunning int32 = 1
	healthChecker := health.NewService(healthLogger, &botRunning)

	// Querying
	getCategories := querying.NewCategoryQuerier(querierLogger, expensesStorage)
	getExpenses := querying.NewExpenseQuerier(querierLogger, expensesStorage, userStorage)
	querier := querying.NewService(*getCategories, *getExpenses)

	// Importing
	// importExpenses := importing.NewExpenseImporter(importerLogger, sqlStorage)

	// Managing
	telegramCommands := make(chan string)
	commandTelegram := managing.NewTelegramCommander(commanderLogger, telegramCommands)
	createCategory := managing.NewCategoryCreator(managerLogger, expensesStorage)
	deleteCategory := managing.NewCategoryDeleter(managerLogger, expensesStorage, ruleStorage)
	updateCategory := managing.NewCategoryUpdater(managerLogger, expensesStorage)
	ruleManager := managing.NewRuleManager(managerLogger, ruleStorage, expensesStorage)
	userManager := managing.NewUserManager(managerLogger, userStorage, expensesStorage)
	manager := managing.NewService(*deleteCategory, *createCategory, *updateCategory, *commandTelegram, *ruleManager, *userManager)
	// Tracking
	createExpense := tracking.NewExpenseCreator(trackerLogger, expensesStorage)
	updateExpense := tracking.NewExpenseUpdater(trackerLogger, expensesStorage)
	deleteExpense := tracking.NewExpenseDeleter(trackerLogger, expensesStorage)
	catalogExpense := tracking.NewExpenseCataloger(trackerLogger, ruleStorage)
	tracker := tracking.NewService(*createExpense, *updateExpense, *deleteExpense, *catalogExpense)

	// Telegram Bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		initLogger.Err("%v", err)
		return
	}
	tgbotapi.SetLogger(telegramLogger)
	allowedUsers := strings.Split(os.Getenv("TELEGRAM_ALLOWED_USERNAMES"), ",")
	go telegram.Run(bot, allowedUsers, telegramCommands, &botRunning, &healthChecker, &tracker, &querier)

	// API
	engine := html.New("./views", ".html")
	engine.AddFunc("nameToColor", ui.NameToColor)
	engine.AddFunc("userInMap", ui.UserInMap)
	engine.AddFunc("unescape", ui.Unescape)
	engine.Debug(true)

	fiberApp := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})
	rest.MapRoutes(fiberApp, &healthChecker, &manager, &tracker, &querier)
	rest.Run(fiberApp, 8080)
}
