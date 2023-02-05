package telegram

import "expenses-app/pkg/app/querying"

func getCategories(cg querying.CategoryGetter) ([]string, error) {
	resp, err := cg.Get()
	if err != nil {
		return nil, err
	}
	categories := []string{}
	for _, name := range resp.Categories {
		categories = append(categories, name)
	}
	return categories, nil
}
