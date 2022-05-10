package main

import (
	"expenses-app/pkg/app/authenticating"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/importing"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/gateways/hasher"
	"expenses-app/pkg/gateways/importers"
	"expenses-app/pkg/gateways/logger"
	"expenses-app/pkg/gateways/storage/json"
	"expenses-app/pkg/gateways/storage/sql"
	"expenses-app/pkg/presenters/http"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Starting..")
	// Infrastructure / Gateways

	// JSON Storage
	jsonStorage := json.NewStorage(os.Getenv("JSON_STORAGE_PATH"))
	// SQL Storage
	dsn := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASS") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DB")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlStorage := sql.NewStorage(db)

	// Importers
	exampleImporter := importers.NewExampleImporter("example data")
	rangeLength, _ := strconv.Atoi(os.Getenv("SHEETS_IMPORTER_RAGENLEN"))
	sheetsImporter := importers.NewSheetsImporter(nil, os.Getenv("SHEETS_IMPORTER_ID"), os.Getenv("SHEETS_IMPORTER_SA_PATH"), os.Getenv("SHEETS_IMPORTER_PAGERANGE"), rangeLength)

	importers := map[string]importing.Importer{
		"example": exampleImporter,
		"sheets":  sheetsImporter,
	}

	//Loggers
	healthLogger := logger.NewSTDLogger("HEALTH", logger.GREEN2)
	managerLogger := logger.NewSTDLogger("Managing", logger.VIOLET)
	importerLogger := logger.NewSTDLogger("Importing", logger.BEIGE)
	querierLogger := logger.NewSTDLogger("Querying", logger.YELLOW2)

	// Hashers
	passHasher := hasher.NewPasswordHasher()

	// Healthching
	healthChecker := health.NewService(healthLogger)

	// Querying
	getCategories := querying.NewCategoryGetter(querierLogger, sqlStorage)
	querier := querying.NewService(*getCategories)

	// Managing
	createUser := managing.NewUserCreator(managerLogger, passHasher, jsonStorage)
	manager := managing.NewService(*createUser)

	// Importing
	importExpenses := importing.NewExpenseImporter(importerLogger, importers, sqlStorage)
	importer := importing.NewService(*importExpenses)

	// Authenticating
	authLogger := logger.NewSTDLogger("Authenticator", logger.RED2)
	authenticateUser := authenticating.NewUserAuthenticator(authLogger, passHasher, jsonStorage)
	authenticator := authenticating.NewAuthenticator(*authenticateUser)

	// API
	fiberApp := fiber.New()
	http.MapRoutes(fiberApp, &healthChecker, &manager, &importer, &authenticator, &querier)
	fiberApp.Listen(":8080")
}
