package sql

import (
	"expenses/pkg/domain/expense"

	"time"

	"gorm.io/plugin/soft_delete"
)

type Category struct {
	ID        string `gorm:"index:idx_name,uniqueIndex:udx_name,primaryKey"`
	Name      string `gorm:"uniqueIndex:udx_name"`
	CreatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:udx_name"`
}

func (sql *SQLStorage) SaveCategory(c expense.Category) error {
	var category Category
	result := sql.db.Unscoped().FirstOrCreate(&category, &Category{ID: string(c.ID), Name: c.Name}) // Filter for "unscoped" rows (i.e already soft-deleted) due to unique constraints at DB level
	sql.db.Model(&category).Update("deleted_at", 0)                                                 // Updated deleted at, I'm I supposed to do this manually
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//func (sql *SQLStorage) CategoryExist(id categories.CategoryID) bool {
//var category Category
//var result *gorm.DB
//result = sql.db.First(&category, "id = ?", id)
//return !errors.Is(result.Error, gorm.ErrRecordNotFound)
//}

func (sql *SQLStorage) DeleteCategory(id categories.CategoryID) error {
	result := sql.db.Delete(&Category{ID: id})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
