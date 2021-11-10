package sql

import (
	"expenses-app/pkg/domain/expense"
)

// Add is used to add a new Expense to the system
func (sql *SQLStorage) Add(e expense.Expense) error {
	result := sql.db.Create(&Expense{
		ID:         string(e.ID),
		Product:    e.Product,
		Shop:       e.Shop,
		City:       e.Town,
		Date:       e.Date,
		CategoryID: string(e.Category.ID),
		//Category: Category{
		//ID:   string(e.Category.ID),
		//Name: e.Category.Name,
		//},
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete is used to remove a Expense from the system
func (sql *SQLStorage) Delete(id expense.ID) error {
	panic("not implemented") // TODO: Implement
}

// SaveCategory stores a category insto a sql database using Gorm ORM
func (sql *SQLStorage) SaveCategory(c expense.Category) error {
	var category Category
	// Filter for "unscoped" rows (i.e already soft-deleted) due to unique constraints at DB level
	result := sql.db.Unscoped().FirstOrCreate(&category, &Category{ID: string(c.ID), Name: c.Name})
	sql.db.Model(&category).Update("deleted_at", 0) // Updated deleted at, I'm I supposed to do this manually
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteCategory deletes a category from the database using Gorm ORM
func (sql *SQLStorage) DeleteCategory(id expense.CategoryID) error {
	result := sql.db.Delete(&Category{ID: string(id)})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
