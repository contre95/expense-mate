package sql

import (
	"gorm.io/gorm"
)

type SQLStorage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *SQLStorage {
	return &SQLStorage{db}
}

func (sql *SQLStorage) Migrate() {
	sql.db.AutoMigrate(&Category{})
	sql.db.AutoMigrate(&Expense{})
	//sql.db.Model(&Expense{}).AddForeignKey()

}

// Paginations for GORM: https://gorm.io/docs/scopes.html#Pagination
func (sql *SQLStorage) paginate(pageSize, page uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}
