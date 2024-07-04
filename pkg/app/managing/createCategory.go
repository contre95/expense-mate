package managing

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
)

type CreateCategoryResp struct {
	ID  string
	Msg string
}

type CreateCategoryReq struct {
	Name string
}

// CategoryCreator use case creates a category for an expense
type CategoryCreator struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewCategoryCreator(l app.Logger, e expense.Expenses) *CategoryCreator {
	return &CategoryCreator{l, e}
}

// Create use cases function creates a new category
func (s *CategoryCreator) Create(req CreateCategoryReq) (*CreateCategoryResp, error) {
	s.logger.Debug("Creating category %s", "name", req.Name)
	newCategory, err := expense.NewCategory(req.Name)
	if err != nil {
		s.logger.Debug("Invalid category %s: %v", req.Name, err)
		return nil, err
	}
	_, err = s.expenses.GetCategory(*&newCategory.ID)
	if err != nil && !errors.Is(err, expense.ErrNotFound) {
		s.logger.Debug("Attempt to create an existing category", req.Name, expense.ErrAlreadyExists)
		return nil, expense.ErrAlreadyExists
	}
	if s.expenses.AddCategory(*newCategory) != nil {
		s.logger.Debug("Couldn't create Category %s: %v", req.Name, expense.ErrAlreadyExists)
		return nil, expense.ErrAlreadyExists
	}
	s.logger.Info("Category %s, created", newCategory.Name)

	resp := &CreateCategoryResp{ID: newCategory.ID.String(), Msg: "Category created"}
	return resp, nil
}
