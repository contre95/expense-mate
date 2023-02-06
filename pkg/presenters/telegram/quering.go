package telegram

import "expenses-app/pkg/app/querying"

func getCategories(cg querying.CategoryQuerier) ([]string, error) {
	resp, err := cg.Query()
	if err != nil {
		return nil, err
	}
	categories := []string{}
	for _, name := range resp.Categories {
		categories = append(categories, name)
	}
	return categories, nil
}
