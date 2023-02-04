package sqlstorage

import (
	"expenses-app/pkg/domain/expense"
	"fmt"
	"log"
)

// Add is used to add a new Expense to the system
func (sql *SQLStorage) Add(e expense.Expense) error {
	stmt, err := sql.db.Prepare("INSTER INTO expenses VALUES(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.ID, e.Price.Amount, e.Product, e.Price.Currency, e.Place.Shop, e.Place.City, e.People, e.Date, e.Category, nil, nil, e.Category.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete is used to remove a Expense from the system
func (sql *SQLStorage) Delete(id expense.ID) error {
	//result := sql.db.Raw("delete from expenses where id = ?", id)
	_, err := sql.db.Exec(fmt.Sprintf("delete from expenses where id = %d", id))
	if err != nil {
		return err
	}
	return nil
}

func (sql *SQLStorage) AddCategory(c expense.Category) error {
	stmt, err := sql.db.Prepare("INSTER INTO expenses VALUES(?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c.ID, c.Name)
	if err != nil {
		return err
	}
	return nil

}

// GetCategories is used to retrieve all categories
func (sql *SQLStorage) CategoryExists(name string) (bool, error) {
	return false, nil
}

// GetCategories is used to retrieve all categories
func (sql *SQLStorage) GetCategories() ([]expense.Category, error) {
	rows, err := sql.db.Query("select * from categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []expense.Category
	for rows.Next() {
		var category expense.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		if _, err := category.Validate(); err != nil {
			log.Printf("Error retrieving category: %v", err)
			continue
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// DeleteCategory deletes a category from the database using Gorm ORM
func (sql *SQLStorage) DeleteCategory(id expense.CategoryID) error {
	_, err := sql.db.Exec(fmt.Sprintf("delete from categories where id = %s", id))
	if err != nil {
		return err
	}
	return nil
}
