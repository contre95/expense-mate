package main

import (
	"database/sql"
	"expenses-app/pkg/app/analyzing"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/config"
	"expenses-app/pkg/gateways/logger"
	"expenses-app/pkg/gateways/ollama"
	"expenses-app/pkg/gateways/storage/jsonstorage"
	"expenses-app/pkg/gateways/storage/sqlstorage"
	"expenses-app/pkg/presenters/rest"
	"expenses-app/pkg/presenters/rest/ui"
	"expenses-app/pkg/presenters/telegram"
	"os"
	"os/signal"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Initialize configuration
	cfg := config.Load()
	initLogger := logger.NewSTDLogger("INIT", logger.VIOLET)

	// Initialize SQLite database
	db, err := sql.Open("sqlite3", cfg.SQLitePath())
	if err != nil {
		initLogger.Err("Error initializing SQLite: %v", err)
		return
	}
	defer db.Close()

	// Create database tables
	if _, err = db.Exec(sqlstorage.SQLiteTables); err != nil {
		initLogger.Err("Error creating SQLite tables: %v", err)
		return
	}

	// Load sample data if enabled
	if cfg.LoadSampleData() {
		if _, err = db.Exec(sqlstorage.SQLiteInserts); err != nil {
			initLogger.Err("Error loading sample data: %v", err)
			return
		}
		initLogger.Info("Loaded sample data into SQLite database")
	}

	// Initialize storage
	expensesStorage := sqlstorage.NewExpensesStorage(db)
	ruleStorage := sqlstorage.NewRulesStorage(db)

	// Initialize JSON storage
	userStorage := jsonstorage.NewStorage(cfg.JSONStoragePath())
	if cfg.LoadSampleData() {
		if err := jsonstorage.CreateFileIfNotExists(cfg.JSONStoragePath(), jsonstorage.SampleUsers); err != nil {
			initLogger.Err("Couldn't create sample user file: %v", err)
			return
		}
	}

	// Initialize services
	healthChecker := health.NewService(logger.NewSTDLogger("HEALTH", logger.GREEN2))

	// Querying
	querier := querying.NewService(
		*querying.NewCategoryQuerier(logger.NewSTDLogger("Querying", logger.YELLOW2), expensesStorage),
		*querying.NewExpenseQuerier(logger.NewSTDLogger("Querying", logger.YELLOW2), expensesStorage, userStorage),
	)

	// Managing
	telegramCommandsSends := make(chan string)
	telegramCommandsReceived := make(chan string)
	manager := managing.NewService(
		*managing.NewCategoryDeleter(logger.NewSTDLogger("Managing", logger.CYAN), expensesStorage, ruleStorage),
		*managing.NewCategoryCreator(logger.NewSTDLogger("Managing", logger.CYAN), expensesStorage),
		*managing.NewCategoryUpdater(logger.NewSTDLogger("Managing", logger.CYAN), expensesStorage),
		*managing.NewTelegramCommander(logger.NewSTDLogger("TELEGRAM COMMANDER", logger.BLUE2), telegramCommandsSends, telegramCommandsReceived),
		*managing.NewRuleManager(logger.NewSTDLogger("Managing", logger.CYAN), ruleStorage, expensesStorage, userStorage),
		*managing.NewUserManager(logger.NewSTDLogger("Managing", logger.CYAN), userStorage, expensesStorage),
	)

	// Analyzing
	analyzer := analyzing.NewService(
		*analyzing.NewSummarizer(logger.NewSTDLogger("Analyzing", logger.RED), expensesStorage),
	)

	// Tracking
	tracker := tracking.NewService(
		*tracking.NewExpenseCreator(logger.NewSTDLogger("Tracker", logger.CYAN), expensesStorage),
		*tracking.NewExpenseUpdater(logger.NewSTDLogger("Tracker", logger.CYAN), expensesStorage),
		*tracking.NewExpenseDeleter(logger.NewSTDLogger("Tracker", logger.CYAN), expensesStorage),
		*tracking.NewRuleApplier(logger.NewSTDLogger("Tracker", logger.CYAN), ruleStorage),
	)

	// Initialize Ollama if enabled
	var guesser *ollama.OllamaAPI
	if cfg.OllamaEnabled() {
		endpoint, textModel, visionModel, timeout := cfg.OllamaConfig()
		guesser, err = ollama.NewOllamaAPI(textModel, visionModel, endpoint, timeout)
		if err != nil {
			initLogger.Err("Failed to initialize Ollama: %v", err)
		}
	}

	// Initialize Telegram Bot if enabled
	if cfg.TelegramEnabled() {
		bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken())
		if err != nil {
			initLogger.Err("Couldn't initialize Telegram bot: %v", err)
			return
		}

		tgbot := telegram.Bot{
			API:    bot,
			Config: cfg,
		}
		ctx := telegram.BotConfig{
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

		tgbotapi.SetLogger(logger.NewSTDLogger("TELEGRAM", logger.BLUE))
		go tgbot.Run(bot, telegramCommandsSends, telegramCommandsReceived, ctx)
	}

	// Initialize web server
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

	go func() {
		if err := rest.Run(fiberApp, 3535); err != nil {
			initLogger.Err("Error starting web server: %v", err)
		}
	}()

	// Shutdown handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	_ = fiberApp.Shutdown()
}
