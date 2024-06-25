package sqlstorage

import (
	"database/sql"
	"errors"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"log"
	"time"
)

const SQL_DATE_FORMAT = "2006-01-02 15:04:05"

func (sqls *SQLStorage) Add(e expense.Expense) error {
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

func (sqls *SQLStorage) Update(e expense.Expense) error {
	q := "UPDATE `expenses` SET price=?, product=?, currency=?, shop=?, city=?, people=?, expend_date=?, category_id=? WHERE id=?"
	fmt.Println("Hola")
	stmt, err := sqls.db.Prepare(q)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.Price.Amount, e.Product, e.Price.Currency, e.Place.Shop, e.Place.City, e.People, e.Date, e.Category.ID, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (sqls *SQLStorage) Delete(id expense.ID) error {
	_, err := sqls.db.Exec(fmt.Sprintf("delete from expenses where id = %s", id))
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves an expense from the db. It returns a valid expense.Expense
func (sqls *SQLStorage) Get(id expense.ID) (*expense.Expense, error) {
	q := "SELECT * FROM expenses where id=?"
	var catID expense.CategoryID
	var e expense.Expense
	err := sqls.db.QueryRow(q, id).Scan(&e.ID, &e.Price.Amount, &e.Product, &e.Price.Currency, &e.Place.Shop, &e.Place.City, &e.People, &e.Date, &catID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, expense.ErrNotFound
	case err != nil:
		return nil, err
	}
	category, err := sqls.GetCategory(catID)
	if err != nil {
		return nil, err
	}
	e.Category = *category
	return &e, nil
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

func (sqls *SQLStorage) GetCategory(id expense.CategoryID) (*expense.Category, error) {
	q := "SELECT * FROM categories where id=?"
	var category expense.Category
	err := sqls.db.QueryRow(q, id).Scan(&category.ID, &category.Name)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, expense.ErrNotFound
	case err != nil:
		return nil, err
	}
	return &category, nil
}

func (sqls *SQLStorage) GetFromTimeRange(from, to time.Time, limit, offset uint) ([]expense.Expense, error) {
	query := `
		SELECT 
			e.id, e.price, e.product, e.currency, e.shop, e.city, 
			e.people, e.expend_date, c.id, c.name 
		FROM 
			expenses e
		JOIN 
			categories c ON e.category_id = c.id
		WHERE 
			e.expend_date >= ? AND e.expend_date <= ?
		ORDER BY 
			e.expend_date DESC 
		LIMIT ? OFFSET ?`

	rows, err := sqls.db.Query(query, from, to, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []expense.Expense
	for rows.Next() {
		var e expense.Expense
		var cat expense.Category

		err := rows.Scan(
			&e.ID, &e.Price.Amount, &e.Product, &e.Price.Currency, &e.Place.Shop, &e.Place.City,
			&e.People, &e.Date, &cat.ID, &cat.Name)
		if err != nil {
			return nil, err
		}

		e.Category = cat
		if _, err = e.Validate(); err != nil {
			return nil, err
		}

		expenses = append(expenses, e)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (sqls *SQLStorage) CategoryExists(id expense.CategoryID) (bool, error) {
	q := "SELECT id FROM categories where id=?"
	var cat_id string
	err := sqls.db.QueryRow(q, id).Scan(&cat_id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, err
}

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

func (sqls *SQLStorage) DeleteCategory(id expense.CategoryID) error {
	_, err := sqls.db.Exec(fmt.Sprintf("delete from categories where id = %s", id))
	if err != nil {
		return err
	}
	return nil
}
