package main

import (
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
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
	// SQL Storage
	db, _ := gorm.Open(sqlite.Open("db/ims.db"), &gorm.Config{})
	storage := sql.NewStorage(db)
	storage.Migrate()

	// Application Healthching
	healthLogger1 := logger.NewSTDLogger("HEALTH", logger.GREEN2)
	healthLogger2 := logger.NewSTDLogger("HEALTH", logger.GREEN)

	healthService := health.NewService(healthLogger1, healthLogger2)

	// Application Managing
	managingLogger := logger.NewSTDLogger("Managing", logger.VIOLET)
	createCategory := managing.NewCreateCategoryUseCase(managingLogger, storage)
	deleteCategory := managing.NewDeleteCategoryUseCase(managingLogger, storage)
	managingService := managing.NewService(*createCategory, *deleteCategory)

	// API
	fiberApp := fiber.New()
	http.MapRoutes(fiberApp, &healthService, &managingService)
	fiberApp.Listen(":3000")
}
