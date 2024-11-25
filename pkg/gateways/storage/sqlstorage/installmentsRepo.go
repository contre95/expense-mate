package sqlstorage

import (
	"database/sql"
	"expenses-app/pkg/domain/expense"
	"strings"
	"time"

	"github.com/google/uuid"
)

//	func (i *InstallmentsStorage) All() ([]expense.Installent, error) {
//		query := `SELECT i.id, i.repeat_every, i.category_id,
//	                     GROUP_CONCAT(ie.expense_id) as expense_ids,
//	                     GROUP_CONCAT(iu.user_id) as user_ids
//	              FROM installments i
//	              LEFT JOIN installment_expenses ie ON i.id = ie.installment_id
//	              LEFT JOIN installment_users iu ON i.id = iu.installment_id
//	              GROUP BY i.id`
//		rows, err := i.db.Query(query)
//		if err != nil {
//			return nil, err
//		}
//
//		defer rows.Close()
//		var installments []expense.Installent
//		for rows.Next() {
//			var installment expense.Installent
//			var expenseIDs sql.NullString
//			var userIDs sql.NullString
//			var repeatEvery int64
//
//			if err := rows.Scan(&installment.ID, &repeatEvery, &installment.CategoryID, &expenseIDs, &userIDs); err != nil {
//				return nil, err
//			}
//
//			installment.RepeatEvery = time.Duration(repeatEvery) * time.Second
//			installment.ExpensesID = []uuid.UUID{}
//			installment.UsersID = []uuid.UUID{}
//
//			// Parse Expense IDs
//			if expenseIDs.Valid && expenseIDs.String != "" {
//				expenseIDStrs := strings.Split(expenseIDs.String, ",")
//				for _, expenseIDStr := range expenseIDStrs {
//					expenseID, err := uuid.Parse(expenseIDStr)
//					if err != nil {
//						return nil, err
//					}
//					installment.ExpensesID = append(installment.ExpensesID, expenseID)
//				}
//			}
//
//			// Parse User IDs
//			if userIDs.Valid && userIDs.String != "" {
//				userIDStrs := strings.Split(userIDs.String, ",")
//				for _, userIDStr := range userIDStrs {
//					userID, err := uuid.Parse(userIDStr)
//					if err != nil {
//						return nil, err
//					}
//					installment.UsersID = append(installment.UsersID, userID)
//				}
//			}
//
//			installments = append(installments, installment)
//		}
//		if err := rows.Err(); err != nil {
//			return nil, err
//		}
//		return installments, nil
//	}

func (i *InstallmentsStorage) All() ([]expense.Installent, error) {
	query := `SELECT i.id, i.repeat_every, i.category_id, 
                     i.start_date, i.end_date, i.amount, i.description, i.product, i.shop,
                     GROUP_CONCAT(ie.expense_id) as expense_ids, 
                     GROUP_CONCAT(iu.user_id) as user_ids 
              FROM installments i 
              LEFT JOIN installment_expenses ie ON i.id = ie.installment_id 
              LEFT JOIN installment_users iu ON i.id = iu.installment_id 
              GROUP BY i.id`
	rows, err := i.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var installments []expense.Installent
	for rows.Next() {
		var installment expense.Installent
		var expenseIDs sql.NullString
		var userIDs sql.NullString
		var repeatEvery int64

		if err := rows.Scan(&installment.ID, &repeatEvery, &installment.CategoryID,
			&installment.StartDate, &installment.EndDate, &installment.Amount,
			&installment.Description, &installment.Product, &installment.Shop,
			&expenseIDs, &userIDs); err != nil {
			return nil, err
		}

		installment.RepeatEvery = time.Duration(repeatEvery) * time.Second
		installment.ExpensesID = []uuid.UUID{}
		installment.UsersID = []uuid.UUID{}

		// Parse Expense IDs
		if expenseIDs.Valid && expenseIDs.String != "" {
			expenseIDStrs := strings.Split(expenseIDs.String, ",")
			for _, expenseIDStr := range expenseIDStrs {
				expenseID, err := uuid.Parse(expenseIDStr)
				if err != nil {
					return nil, err
				}
				installment.ExpensesID = append(installment.ExpensesID, expenseID)
			}
		}

		// Parse User IDs
		if userIDs.Valid && userIDs.String != "" {
			userIDStrs := strings.Split(userIDs.String, ",")
			for _, userIDStr := range userIDStrs {
				userID, err := uuid.Parse(userIDStr)
				if err != nil {
					return nil, err
				}
				installment.UsersID = append(installment.UsersID, userID)
			}
		}

		installments = append(installments, installment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return installments, nil
}

func (i *InstallmentsStorage) Add(installment expense.Installent) error {
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}

	query := "INSERT INTO installments (id, repeat_every, category_id) VALUES (?, ?, ?)"
	_, err = tx.Exec(query, installment.ID, installment.RepeatEvery.Seconds(), installment.CategoryID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, expenseID := range installment.ExpensesID {
		expenseQuery := "INSERT INTO installment_expenses (installment_id, expense_id) VALUES (?, ?)"
		_, err := tx.Exec(expenseQuery, installment.ID, expenseID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, userID := range installment.UsersID {
		userQuery := "INSERT INTO installment_users (installment_id, user_id) VALUES (?, ?)"
		_, err := tx.Exec(userQuery, installment.ID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (i *InstallmentsStorage) Delete(id string) error {
	query := "DELETE FROM installments WHERE id = ?"
	_, err := i.db.Exec(query, id)
	return err
}
