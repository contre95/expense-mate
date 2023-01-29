package sql

import (
	"errors"
	"expenses-app/pkg/domain/expense"

	"gorm.io/gorm"
)

// Add is used to add a new Expense to the system
func (sql *SQLStorage) Add(e expense.Expense) error {
	var category Category
	res := sql.db.First(&category, e.Category.ID)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		category = Category{
			ID:   string(e.Category.ID),
			Name: string(e.Category.Name),
		}
		sql.db.Create(&category)
	}
	result := sql.db.Create(&Expense{
		ID:       uint64(e.ID),
		Price:    e.Price.Amount,
		Currency: e.Price.Currency,
		Product:  e.Product,
		Shop:     e.Place.Shop,
		City:     e.Place.Town,
		Date:     e.Date,
		People:   e.People,
		Category: category,
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

func parseCategories(categories []Category) []expense.Category {
	domainCategories := []expense.Category{}
	for _, c := range categories {
		domainCat := expense.Category{
			ID:   expense.CategoryID(c.ID),
			Name: expense.CategoryName(c.Name),
		}
		domainCategories = append(domainCategories, domainCat)
	}
	return domainCategories
}

// GetCategories is used to retrieve all categories
func (sql *SQLStorage) CategoryExist() (bool, error) {
	return false, nil
}

// GetCategories is used to retrieve all categories
func (sql *SQLStorage) GetCategories() ([]expense.Category, error) {
	var categories []Category
	result := sql.db.Raw("select * from categories").Scan(&categories)
	//result := sql.db.Delete(&Category{ID: string(id)})
	if result.Error != nil {
		return nil, result.Error
	}
	return parseCategories(categories), nil
}

// DeleteCategory deletes a category from the database using Gorm ORM
func (sql *SQLStorage) DeleteCategory(id expense.CategoryID) error {
	result := sql.db.Raw("delete from categories where id = ?", id)
	//result := sql.db.Delete(&Category{ID: string(id)})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
