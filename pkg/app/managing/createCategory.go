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
	newCategory, err := expense.NewCategory(req.Name)
	if errors.Is(err, expense.ErrInvalidEntity) {
		s.logger.Debug("Invalid category %s: %v", req, err)
		return nil, err
	}
	if errors.Is(s.expenses.AddCategory(*newCategory), expense.ErrAlreadyExists) {
		s.logger.Debug("Couldn't create Category %s: %v", req.Name, expense.ErrAlreadyExists)
		return nil, expense.ErrAlreadyExists
	}
	s.logger.Info("Category %s, created", newCategory.Name)
	resp := &CreateCategoryResp{ID: string(newCategory.ID), Msg: "Category created"}
	return resp, nil
}
