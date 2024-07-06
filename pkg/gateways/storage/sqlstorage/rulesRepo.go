package sqlstorage

import "expenses-app/pkg/domain/expense"

func (r *RulesStorage) All() ([]expense.Rule, error) {
	query := "SELECT id, pattern, category_id FROM rules"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []expense.Rule
	for rows.Next() {
		var rule expense.Rule
		if err := rows.Scan(&rule.ID, &rule.Pattern, &rule.CategoryID); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rules, nil
}

func (r *RulesStorage) Add(rule expense.Rule) error {
	query := "INSERT INTO rules (id, pattern, category_id) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, rule.ID, rule.Pattern, rule.CategoryID)
	return err
}

func (r *RulesStorage) Delete(id string) error {
	deleteQuery := "DELETE FROM rules WHERE id = ?"
	_, err := r.db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}
	return nil
}
