package querying

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
)

type GetCategoriesResp struct {
	Categories map[string]string
}

type GetCategoriesReq struct {
	CategoriesIDs map[string]bool
}

type CategoryGetter struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewCategoryGetter(l app.Logger, e expense.Expenses) *CategoryGetter {
	return &CategoryGetter{l, e}
}

func (s *CategoryGetter) Get(req GetCategoriesReq) (*GetCategoriesResp, error) {
	categories, err := s.expenses.GetCategories()
	if err != nil {
		s.logger.Err("Could not get categories from storage")
		return nil, err
	}
	var resp GetCategoriesResp
	for _, c := range categories {
		if len(req.CategoriesIDs) == 0 || req.CategoriesIDs[string(c.ID)] {
			resp.Categories[string(c.ID)] = string(c.Name)
		}
	}
	return &resp, nil

}
