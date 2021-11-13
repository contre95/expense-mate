package main

import (
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/importing"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/gateways/importers"
	"expenses-app/pkg/gateways/logger"
	"expenses-app/pkg/gateways/storage/sql"
	"expenses-app/pkg/presenters/http"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Starting")
	// Infrastructure / Gateways

	// SQL Storage
	db, _ := gorm.Open(sqlite.Open("db/ims.db"), &gorm.Config{})
	storage := sql.NewStorage(db)
	storage.Migrate()

	// Importers
	exampleImporter := importers.NewExampleImporter("example data")
	importers := map[string]importing.Importer{
		"example": exampleImporter,
	}

	// Healthching
	healthLogger1 := logger.NewSTDLogger("HEALTH", logger.GREEN2)
	healthLogger2 := logger.NewSTDLogger("HEALTH", logger.GREEN)
	healthChecker := health.NewService(healthLogger1, healthLogger2)

	// Managing
	managerLogger := logger.NewSTDLogger("Managing", logger.VIOLET)
	createCategory := managing.NewCreateCategory(managerLogger, storage)
	deleteCategory := managing.NewDeleteCategory(managerLogger, storage)
	manager := managing.NewService(*createCategory, *deleteCategory)

	// Importing
	importerLogger := logger.NewSTDLogger("Importing", logger.VIOLET)
	importExpenses := importing.NewExpenseImporter(importerLogger, importers, storage)
	importer := importing.NewService(*importExpenses)

	// API
	fiberApp := fiber.New()
	http.MapRoutes(fiberApp, &healthChecker, &manager, &importer)
	fiberApp.Listen(":3000")
}
