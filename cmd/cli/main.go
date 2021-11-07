package main

import (
	"codelamp-ims/pkg/gateways/logger"
	"codelamp-ims/pkg/gateways/storage/sql"
	"spain-gastos/pkg/domain/categories"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("db/ims.db"), &gorm.Config{})
	storage := sql.NewStorage(db)
	storage.Migrate()

	categoriesLogger := logger.NewSTDLogger("CATEGORIES", logger.VIOLET)

	add := categories.NewAddUseCase(categoriesLogger, storage)
	del := categories.NewDeleteUseCase(categoriesLogger, storage)

	categoriesService := categories.NewService(add, del)

}
