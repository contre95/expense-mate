package main

import (
	"database/sql"
	"expenses-app/pkg/app/analyzing"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/gateways/ai"
	"expenses-app/pkg/gateways/logger"
	"expenses-app/pkg/gateways/storage/jsonstorage"
	"expenses-app/pkg/gateways/storage/sqlstorage"
	"expenses-app/pkg/presenters/rest"
	"expenses-app/pkg/presenters/rest/ui"
	"expenses-app/pkg/presenters/telegram"
	"os"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Infrastructure / Gateways

	//Loggers
	initLogger := logger.NewSTDLogger("INIT", logger.VIOLET)
	healthLogger := logger.NewSTDLogger("HEALTH", logger.GREEN2)
	analyzingLogger := logger.NewSTDLogger("Analyzing", logger.RED)
	managerLogger := logger.NewSTDLogger("Managing", logger.CYAN)
	querierLogger := logger.NewSTDLogger("Querying", logger.YELLOW2)
	trackerLogger := logger.NewSTDLogger("Tracker", logger.CYAN)
	telegramLogger := logger.NewSTDLogger("TELEGRAM", logger.BLUE)
	commanderLogger := logger.NewSTDLogger("TELEGRAM COMMANDER", logger.BLUE2)

	//AI
	txtModel := os.Getenv("TEXT_MODEL")
	visioModel := os.Getenv("VISION_MODEL")
	ollamaURL := os.Getenv("OLLAMA_ENDPOINT")
	guesser, err := ai.NewGuesser(txtModel, visioModel, ollamaURL)
	if err != nil {
		initLogger.Err("failed to create guesser: %v", err)
	}
	// SQL storage
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
			if os.Getenv("LOAD_SAMPLE_DATA") == "true" {
				_, err = db.Exec(sqlstorage.MySQLInserts)
				if err != nil {
					initLogger.Err("Error creating sqlite tables: %v", err)
					return
				}
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
		if os.Getenv("LOAD_SAMPLE_DATA") == "true" {
			_, err = db.Exec(sqlstorage.SQLiteInserts)
			if err != nil {
				initLogger.Err("Error creating sqlite tables: %v", err)
				return
			}
		}
		initLogger.Info("SQLte storage initialized on %s", path)
	case "":
		initLogger.Err("No storage set. Please set STORAGE_ENGINE variabel")
	}
	expensesStorage := sqlstorage.NewExpensesStorage(db)
	ruleStorage := sqlstorage.NewRulesStorage(db)

	// JSON Storage
	path := os.Getenv("JSON_STORAGE_PATH")
	if path == "" {
		initLogger.Err("No storage set. Please set STORAGE_ENGINE variabel")
		return
	}
	initLogger.Info("LOAD_SAMPLE_DATA=%s", os.Getenv("LOAD_SAMPLE_DATA"))
	if os.Getenv("LOAD_SAMPLE_DATA") == "true" {
		err := jsonstorage.CreateFileIfNotExists(path, jsonstorage.SampleUsers)
		if err != nil {
			initLogger.Err("Couldn't create sample file: %v", err)
			return
		}
	}
	userStorage := jsonstorage.NewStorage(path)

	// Healthching
	healthChecker := health.NewService(healthLogger)

	// Querying
	getCategories := querying.NewCategoryQuerier(querierLogger, expensesStorage)
	getExpenses := querying.NewExpenseQuerier(querierLogger, expensesStorage, userStorage)
	querier := querying.NewService(*getCategories, *getExpenses)

	// Managing
	telegramCommandsSends := make(chan string)
	telegramCommandsReceived := make(chan string)
	commandTelegram := managing.NewTelegramCommander(commanderLogger, telegramCommandsSends, telegramCommandsReceived)
	createCategory := managing.NewCategoryCreator(managerLogger, expensesStorage)
	deleteCategory := managing.NewCategoryDeleter(managerLogger, expensesStorage, ruleStorage)
	updateCategory := managing.NewCategoryUpdater(managerLogger, expensesStorage)
	ruleManager := managing.NewRuleManager(managerLogger, ruleStorage, expensesStorage, userStorage)
	userManager := managing.NewUserManager(managerLogger, userStorage, expensesStorage)
	manager := managing.NewService(*deleteCategory, *createCategory, *updateCategory, *commandTelegram, *ruleManager, *userManager)

	// Analyzing
	sumerizePerCategory := analyzing.NewSummarizer(analyzingLogger, expensesStorage)
	analyzer := analyzing.NewService(*sumerizePerCategory)

	// Tracking
	createExpense := tracking.NewExpenseCreator(trackerLogger, expensesStorage)
	updateExpense := tracking.NewExpenseUpdater(trackerLogger, expensesStorage)
	deleteExpense := tracking.NewExpenseDeleter(trackerLogger, expensesStorage)
	catalogExpense := tracking.NewRuleApplier(trackerLogger, ruleStorage)
	tracker := tracking.NewService(*createExpense, *updateExpense, *deleteExpense, *catalogExpense)

	// Telegram Bot
	token := os.Getenv("TELEGRAM_APITOKEN")
	if token != "" {
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			initLogger.Err("Couldn't init Telegram bot, please check the token: %v", err)
			return
		}
		tgbotapi.SetLogger(telegramLogger)
		tgbot := telegram.Bot{
			API: bot,
		}
		ctx := telegram.BotContext{
			BotAPI:       bot,
			Health:       &healthChecker,
			Tracking:     &tracker,
			Querying:     &querier,
			Managing:     &manager,
			Analyzing:    &analyzer,
			AI:           guesser,
			AllowedUsers: &tgbot.AllowedUsers,
			Mu:           &sync.Mutex{},
		}

		// Run the bot with the updated context
		go tgbot.Run(bot, telegramCommandsSends, telegramCommandsReceived, ctx)
		// go tgbot.Run(bot, telegramCommandsSends, telegramCommandsReceived, &healthChecker, &tracker, &querier, &manager, &analyzer, guesser)
		// go tgbot.Run(bot, telegramCommandsSends, telegramCommandsReceived, &healthChecker, &tracker, &querier, &manager, &analyzer, guesser)
	}

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
	rest.MapRoutes(fiberApp, &healthChecker, &manager, &tracker, &querier, &analyzer)
	if err := rest.Run(fiberApp, 3535); err != nil {
		initLogger.Err("Error instanciating mysql: %v", err)
		return
	}

}
