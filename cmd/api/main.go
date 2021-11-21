package main

import (
	"expenses-app/pkg/app/authenticating"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/importing"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/gateways/importers"
	"expenses-app/pkg/gateways/logger"
	"expenses-app/pkg/gateways/storage/json"
	"expenses-app/pkg/gateways/storage/sql"
	"expenses-app/pkg/presenters/http"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Starting")
	// Infrastructure / Gateways

	// JSON Storage
	jsonStorage := json.NewStorage()
	// SQL Storage
	dsn := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASS") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DB")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlStorage := sql.NewStorage(db)

	// Importers
	exampleImporter := importers.NewExampleImporter("example data")

	srv, _ := importers.NewSheetService(os.Getenv("SHEETS_IMPORTER_SA_PATH"))
	sheetsImporter := importers.NewSheetsImporter(srv, os.Getenv("SHEETS_IMPORTER_ID"), os.Getenv("SHEETS_IMPORTER_PAGERANGE"))

	importers := map[string]importing.Importer{
		"example": exampleImporter,
		"sheets":  sheetsImporter,
	}

	// Healthching
	healthLogger1 := logger.NewSTDLogger("HEALTH", logger.GREEN2)
	healthLogger2 := logger.NewSTDLogger("HEALTH", logger.GREEN)
	healthChecker := health.NewService(healthLogger1, healthLogger2)

	// Managing
	managerLogger := logger.NewSTDLogger("Managing", logger.VIOLET)
	createCategory := managing.NewCategoryCreator(managerLogger, sqlStorage)
	deleteCategory := managing.NewCategoryDeleter(managerLogger, sqlStorage)
	createUser := managing.NewUserCreator(managerLogger, jsonStorage)
	manager := managing.NewService(*createCategory, *deleteCategory, *createUser)

	// Importing
	importerLogger := logger.NewSTDLogger("Importing", logger.VIOLET)
	importExpenses := importing.NewExpenseImporter(importerLogger, importers, sqlStorage)
	importer := importing.NewService(*importExpenses)

	// Authenticating
	authLogger := logger.NewSTDLogger("Authenticator", logger.RED2)
	authenticateUser := authenticating.NewUserAuthenticator(authLogger, jsonStorage)
	authenticator := authenticating.NewAuthenticator(*authenticateUser)

	// API
	fiberApp := fiber.New()
	http.MapRoutes(fiberApp, &healthChecker, &manager, &importer, &authenticator)
	fiberApp.Listen(":3000")
}
