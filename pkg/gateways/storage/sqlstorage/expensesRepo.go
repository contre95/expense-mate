package sqlstorage

import (
	"database/sql"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"log"
)

// Add is used to add a new Expense to the system
func (sqls *SQLStorage) Add(e expense.Expense) error {
	//q := "INSTER INTO `expenses` (price, product, currency, shop, city, people, `date`, created_at, updated_at, category_id)  VALUES (?,?,?,?,?,?,?,?,?,?);"
	exist, err := sqls.CategoryExists(e.Category.ID)
	if err != nil {
		return err
	}
	if !exist {
		err := sqls.AddCategory(e.Category)
		if err != nil {
			return err
		}
	}
	q := "INSERT INTO `expenses` (id, price, product, currency, shop, city, people, expend_date, category_id) VALUES (?,?,?,?,?,?,?,?,?);"
	stmt, err := sqls.db.Prepare(q)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.ID, e.Price.Amount, e.Product, e.Price.Currency, e.Place.Shop, e.Place.City, e.People, e.Date, e.Category.ID)
	stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

// Delete is used to remove a Expense from the system
func (sqls *SQLStorage) Delete(id expense.ID) error {
	_, err := sqls.db.Exec(fmt.Sprintf("delete from expenses where id = %d", id))
	if err != nil {
		return err
	}
	return nil
}

func (sqls *SQLStorage) AddCategory(c expense.Category) error {
	stmt, err := sqls.db.Prepare("INSERT INTO `categories` (id, name) VALUES (?,?)")
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
func (sqls *SQLStorage) CategoryExists(id expense.CategoryID) (bool, error) {
	q := "SELECT id FROM categories where id=?"
	var cat_id string
	err := sqls.db.QueryRow(q, id).Scan(&cat_id)
	switch err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	}
	return true, err
}

// GetCategories is used to retrieve all categories
func (sqls *SQLStorage) GetCategories() ([]expense.Category, error) {
	rows, err := sqls.db.Query("select * from categories")
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
func (sqls *SQLStorage) DeleteCategory(id expense.CategoryID) error {
	_, err := sqls.db.Exec(fmt.Sprintf("delete from categories where id = %s", id))
	if err != nil {
		return err
	}
	return nil
}
