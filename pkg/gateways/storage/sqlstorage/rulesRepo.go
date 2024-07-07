package sqlstorage

import (
	"database/sql"
	"expenses-app/pkg/domain/expense"
	"strings"

	"github.com/google/uuid"
)

func (r *RulesStorage) All() ([]expense.Rule, error) {
	query := `SELECT r.id, r.pattern, r.category_id, 
                     GROUP_CONCAT(ru.user_id) as user_ids 
              FROM rules r 
              LEFT JOIN rule_users ru ON r.id = ru.rule_id 
              GROUP BY r.id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []expense.Rule
	for rows.Next() {
		var rule expense.Rule
		var userIDs sql.NullString // Use sql.NullString to handle NULL values

		if err := rows.Scan(&rule.ID, &rule.Pattern, &rule.CategoryID, &userIDs); err != nil {
			return nil, err
		}

		// Initialize UsersID slice
		rule.UsersID = []uuid.UUID{}

		// Split user IDs into a slice and convert to uuid.UUID
		if userIDs.Valid && userIDs.String != "" {
			userIDStrs := strings.Split(userIDs.String, ",")
			for _, userIDStr := range userIDStrs {
				userID, err := uuid.Parse(userIDStr)
				if err != nil {
					return nil, err
				}
				rule.UsersID = append(rule.UsersID, userID)
			}
		}

		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rules, nil
}
func (r *RulesStorage) Add(rule expense.Rule) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := "INSERT INTO rules (id, pattern, category_id) VALUES (?, ?, ?)"
	_, err = tx.Exec(query, rule.ID, rule.Pattern, rule.CategoryID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, userID := range rule.UsersID {
		userQuery := "INSERT INTO rule_users (rule_id, user_id) VALUES (?, ?)"
		_, err := tx.Exec(userQuery, rule.ID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *RulesStorage) Delete(id string) error {
	query := "DELETE FROM rules WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
