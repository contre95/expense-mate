package main

import (
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
	categoriesLogger := logger.NewSTDLogger("CATEGORIES", logger.VIOLET)

	http.MapRoutes(fiberApp, &healthService, &categoriesService)
	fiberApp.Listen(":3000")
}
