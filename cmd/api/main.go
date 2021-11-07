package main

import (
	"expenses/pkg/app/tracking/managing"
	"expenses/pkg/gateways/logger"
	"expenses/pkg/gateways/storage/sql"
	"expenses/pkg/presenters/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	db, _ := gorm.Open(sqlite.Open("db/ims.db"), &gorm.Config{})
	storage := sql.NewStorage(db)
	storage.Migrate()

	healthLogger := logger.NewSTDLogger("HEALTH", logger.GREEN2)
	managingLogger := logger.NewSTDLogger("Managing", logger.VIOLET)

	createCategory := managing.NewCreateCategoryUseCase(managingLogger, storage)
	deleteCategory := managing.NewDeleteCategoryUseCase(managingLogger, storage)

	http.MapRoutes(fiberApp, &healthService, &categoriesService)
	fiberApp.Listen(":3000")
}
