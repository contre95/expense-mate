package sqlstorage

import (
	"database/sql"
	"errors"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

const SQL_DATE_FORMAT = "2006-01-02 15:04:05"

func (sqls *ExpensesStorage) Add(e expense.Expense) error {
	q := "INSERT INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES (?,?,?,?,?,?);"
	stmt, err := sqls.db.Prepare(q)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.ID, e.Amount, e.Product, e.Shop, e.Date, e.Category.ID)
	stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (sqls *ExpensesStorage) Update(e expense.Expense) error {
	q := "UPDATE expenses SET amount=?, product=?, shop=?, expend_date=?, category_id=? WHERE id=?"
	stmt, err := sqls.db.Prepare(q)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.Amount, e.Product, e.Shop, e.Date, e.Category.ID, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (sqls *ExpensesStorage) Delete(id expense.ExpenseID) error {
	stmt, err := sqls.db.Prepare("DELETE FROM expenses WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil

}

// Get retrieves an expense from the db. It returns a valid expense.Expense
// func (sqls *ExpensesStorage) Get(id expense.ExpenseID) (*expense.Expense, error) {
// 	q := "SELECT id, amount, product, shop, expend_date, category_id FROM expenses where id=?"
// 	var catID expense.CategoryID
// 	var e expense.Expense
// 	err := sqls.db.QueryRow(q, id).Scan(&e.ID, &e.Amount, &e.Product, &e.Shop, &e.Date, &catID)
// 	switch {
// 	case errors.Is(err, sql.ErrNoRows):
// 		return nil, expense.ErrNotFound
// 	case err != nil:
// 		return nil, err
// 	}
// 	category, err := sqls.GetCategory(catID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	e.Category = *category
// 	return &e, nil
// }

func (sqls *ExpensesStorage) Get(id expense.ExpenseID) (*expense.Expense, error) {
	q := "SELECT id, amount, product, shop, expend_date, category_id FROM expenses WHERE id=?"
	var catID expense.CategoryID
	var e expense.Expense
	err := sqls.db.QueryRow(q, id).Scan(&e.ID, &e.Amount, &e.Product, &e.Shop, &e.Date, &catID)
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
	// Retrieve associated user IDs
	userQuery := "SELECT user_id FROM expense_users WHERE expense_id=?"
	rows, err := sqls.db.Query(userQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID // What is this magic and why am I scanning directly into a uuid.UUID ?
		// Maybe adding a vlidation here should be the best thing to do. Like.. what if someone gets to store "not-uuid" into the DB ? :O
		if err := rows.Scan(&userID); err != nil { // Read the Scan docu. Try to change uuid.UUID to uint64 and see what happens.
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	e.UserIDS = userIDs
	fmt.Println(e)
	return &e, nil
}

func (sqls *ExpensesStorage) UpdateCategory(c expense.Category) error {
	stmt, err := sqls.db.Prepare("UPDATE categories SET name=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c.Name, c.ID)
	if err != nil {
		return err
	}
	return nil
}

func (sqls *ExpensesStorage) AddCategory(c expense.Category) error {
	stmt, err := sqls.db.Prepare("INSERT INTO categories (id, name) VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c.ID, c.Name)
	if err != nil {
		return err
	}
	return nil
}

func (sqls *ExpensesStorage) GetCategory(id expense.CategoryID) (*expense.Category, error) {
	q := "SELECT id, name FROM categories where id=?"
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

// CountWithFilter
func (sqls *ExpensesStorage) CountWithFilter(users_ids, categories_ids []string, minAmount, maxAmount uint, shop, product string, from, to time.Time) (uint, error) {
	var conditions []string
	query := "SELECT COUNT(DISTINCT e.id) FROM expenses e JOIN categories c ON e.category_id = c.id"
	if !(len(users_ids) == 0 || (len(users_ids) == 1 && len(users_ids[0]) == 0)) {
		query += " JOIN expense_users eu ON e.id = eu.expense_id"
	}
	if !from.IsZero() {
		conditions = append(conditions, fmt.Sprintf("e.expend_date >= '%s'", from.Format("2006-01-02")))
	}
	if !to.IsZero() {
		conditions = append(conditions, fmt.Sprintf("e.expend_date <= '%s'", to.Format("2006-01-02")))
	}
	if minAmount > 0 {
		conditions = append(conditions, fmt.Sprintf("e.amount >= %.2f", float64(minAmount)))
	}
	if maxAmount > 0 {
		conditions = append(conditions, fmt.Sprintf("e.amount <= %.2f", float64(maxAmount)))
	}
	if shop != "" {
		conditions = append(conditions, fmt.Sprintf("e.shop LIKE '%%%s%%'", shop))
	}
	if product != "" {
		conditions = append(conditions, fmt.Sprintf("e.product LIKE '%%%s%%'", product))
	}
	if !(len(categories_ids) == 0 || (len(categories_ids) == 1 && len(categories_ids[0]) == 0)) { // This mean they are sending []string{""}.
		categoryConditions := make([]string, len(categories_ids))
		for i, cat := range categories_ids {
			categoryConditions[i] = fmt.Sprintf("c.id = '%s'", cat)
		}
		conditions = append(conditions, "("+strings.Join(categoryConditions, " OR ")+")")
	}
	if len(users_ids) > 0 {
		userConditions := make([]string, len(users_ids))
		for i, uid := range users_ids {
			userConditions[i] = fmt.Sprintf("eu.user_id = '%s'", uid)
		}
		conditions = append(conditions, "("+strings.Join(userConditions, " OR ")+")")
	}
	if len(conditions) > 0 {
		whereClause := strings.Join(conditions, " AND ")
		query += " WHERE " + whereClause
	}
	var count uint
	err := sqls.db.QueryRow(query).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, expense.ErrNotFound
	} else if err != nil {
		return 0, err
	}
	return count, nil
}

// // Filter retrieves expenses from the db based on the given filters. It skips the filter parameters with zero value
//
//	func (sqls *ExpensesStorage) Filter(categories []string, minAmount, maxAmount uint, shop, product string, from time.Time, to time.Time, limit, offset uint) ([]expense.Expense, error) {
//		var conditions []string
//		query := "SELECT e.id, e.amount, e.product, e.shop, e.expend_date, c.id, c.name FROM expenses e JOIN categories c ON e.category_id = c.id"
//		if !from.IsZero() {
//			conditions = append(conditions, fmt.Sprintf("expend_date >= '%s'", from.Format("2006-01-02")))
//		}
//		if !to.IsZero() {
//			conditions = append(conditions, fmt.Sprintf("expend_date <= '%s'", to.Format("2006-01-02")))
//		}
//		if minAmount > 0 {
//			conditions = append(conditions, fmt.Sprintf("amount >= %.2f", float64(minAmount)))
//		}
//		if maxAmount > 0 {
//			conditions = append(conditions, fmt.Sprintf("amount <= %.2f", float64(maxAmount)))
//		}
//		if shop != "" {
//			conditions = append(conditions, fmt.Sprintf("shop LIKE '%%%s%%'", shop))
//		}
//		if product != "" {
//			conditions = append(conditions, fmt.Sprintf("product LIKE '%%%s%%'", product))
//		}
//		if !(len(categories) == 1 && len(categories[0]) == 0) { // This mean they are sending []string{""}.
//			categoryConditions := make([]string, len(categories))
//			for i, cat := range categories {
//				categoryConditions[i] = fmt.Sprintf("c.id ='%s'", cat)
//			}
//			conditions = append(conditions, "("+strings.Join(categoryConditions, " OR ")+")")
//		}
//		if len(conditions) > 0 {
//			whereClause := " " + strings.Join(conditions, " AND ")
//			query += " WHERE " + whereClause
//		}
//		query += fmt.Sprintf(" ORDER BY e.expend_date")
//		rows, err := sqls.db.Query(query)
//		if limit > 0 {
//			query += fmt.Sprintf(" DESC LIMIT ? OFFSET ?")
//			rows, err = sqls.db.Query(query, limit, offset)
//		}
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, expense.ErrNotFound
//		} else if err != nil {
//			return nil, err
//		}
//		defer rows.Close()
//		var expenses []expense.Expense
//		for rows.Next() {
//			var e expense.Expense
//			var cat expense.Category
//			err := rows.Scan(&e.ID, &e.Amount, &e.Product, &e.Shop, &e.Date, &cat.ID, &cat.Name)
//			if err != nil {
//				return nil, err
//			}
//			e.Category = cat
//			if _, err = e.Validate(); err != nil {
//				return nil, err
//			}
//			expenses = append(expenses, e)
//		}
//		return expenses, nil
//	}

func (sqls *ExpensesStorage) Filter(user_ids, categories_ids []string, minAmount, maxAmount uint, shop, product string, from time.Time, to time.Time, limit, offset uint) ([]expense.Expense, error) {
	var conditions []string
	query := `
		SELECT e.id, e.amount, e.product, e.shop, e.expend_date, c.id, c.name, eu.user_id
		FROM expenses e
		JOIN categories c ON e.category_id = c.id
		LEFT JOIN expense_users eu ON e.id = eu.expense_id
	`
	if !from.IsZero() {
		conditions = append(conditions, fmt.Sprintf("expend_date >= '%s'", from.Format("2006-01-02")))
	}
	if !to.IsZero() {
		conditions = append(conditions, fmt.Sprintf("expend_date <= '%s'", to.Format("2006-01-02")))
	}
	if minAmount > 0 {
		conditions = append(conditions, fmt.Sprintf("amount >= %.2f", float64(minAmount)))
	}
	if maxAmount > 0 {
		conditions = append(conditions, fmt.Sprintf("amount <= %.2f", float64(maxAmount)))
	}
	if shop != "" {
		conditions = append(conditions, fmt.Sprintf("shop LIKE '%%%s%%'", shop))
	}
	if product != "" {
		conditions = append(conditions, fmt.Sprintf("product LIKE '%%%s%%'", product))
	}
	if !(len(categories_ids) == 1 && len(categories_ids[0]) == 0) { // This means they are sending []string{""}.
		categoryConditions := make([]string, len(categories_ids))
		for i, cat := range categories_ids {
			categoryConditions[i] = fmt.Sprintf("c.id ='%s'", cat)
		}
		conditions = append(conditions, "("+strings.Join(categoryConditions, " OR ")+")")
	}
	if len(user_ids) > 0 { // This means they are sending []string{""}.
		userConditions := make([]string, len(user_ids))
		for i, uid := range user_ids {
			userConditions[i] = fmt.Sprintf("eu.user_id ='%s'", uid)
		}
		conditions = append(conditions, "("+strings.Join(userConditions, " OR ")+")")
	}
	if len(conditions) > 0 {
		whereClause := " " + strings.Join(conditions, " AND ")
		query += " WHERE " + whereClause
	}
	query += " ORDER BY e.expend_date"
	if limit > 0 {
		query += fmt.Sprintf(" DESC LIMIT %d OFFSET %d", limit, offset)
	}
	rows, err := sqls.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, expense.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenseMap := make(map[string]*expense.Expense)
	for rows.Next() {
		var e expense.Expense
		var cat expense.Category
		var userID sql.NullString
		err := rows.Scan(&e.ID, &e.Amount, &e.Product, &e.Shop, &e.Date, &cat.ID, &cat.Name, &userID)
		if err != nil {
			return nil, err
		}
		e.Category = cat
		if existingExpense, exists := expenseMap[e.ID.String()]; exists {
			if userID.Valid {
				uid, err := uuid.Parse(userID.String)
				if err != nil {
					return nil, err
				}
				existingExpense.UserIDS = append(existingExpense.UserIDS, uid)
			}
		} else {
			if userID.Valid {
				uid, err := uuid.Parse(userID.String)
				if err != nil {
					return nil, err
				}
				e.UserIDS = []expense.UserID{uid}
			} else {
				e.UserIDS = []expense.UserID{}
			}
			expenseMap[e.ID.String()] = &e
		}
	}
	var expenses []expense.Expense
	for _, e := range expenseMap {
		if _, err = e.Validate(); err != nil {
			return nil, err
		}
		expenses = append(expenses, *e)
	}
	return expenses, nil
}

func (sqls *ExpensesStorage) All(limit, offset uint) ([]expense.Expense, error) {
	query := `
		SELECT e.id, e.amount, e.product, e.shop,  e.expend_date, c.id, c.name FROM expenses e JOIN categories c ON e.category_id = c.id
		ORDER BY e.expend_date DESC
		LIMIT ? OFFSET ?`
	rows, err := sqls.db.Query(query, limit, offset)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, expense.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()
	var expenses []expense.Expense
	for rows.Next() {
		var e expense.Expense
		var cat expense.Category
		err := rows.Scan(
			&e.ID, &e.Amount, &e.Product, &e.Shop, &e.Date, &cat.ID, &cat.Name)
		if err != nil {
			return nil, err
		}
		e.Category = cat
		if _, err = e.Validate(); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (sqls *ExpensesStorage) CategoryExists(id expense.CategoryID) (bool, error) {
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

func (sqls *ExpensesStorage) GetCategories() ([]expense.Category, error) {
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

func (sqls *ExpensesStorage) DeleteCategory(id expense.CategoryID) error {
	_, err := sqls.db.Exec(fmt.Sprintf("delete from categories where id='%s'", id))
	if err != nil {
		return err
	}
	return nil
}
