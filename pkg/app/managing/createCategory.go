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
	newCategory, createErr := expense.NewCategory(req.Name)
	if createErr != nil {
		s.logger.Debug("Invalid category %s: %v", req, createErr)
		return nil, createErr
	}
	exist, err := s.expenses.CategoryExists(req.Name)
	if err != nil {
		s.logger.Err("Could not check category existance", req, createErr)
		return nil, err
	}
	if !exist {
		err := s.expenses.AddCategory(*newCategory)
		if err != nil && err.Error() != expense.CategoryAlreadyExists {
			s.logger.Err("Could not save category: %v", err)
			return nil, err
		}
		resp := &CreateCategoryResp{ID: string(newCategory.ID), Msg: "Category created"}
		s.logger.Info("Category %s, created", newCategory.Name)
		return resp, nil
	}
	s.logger.Warn("Tried to add existing category", req, createErr)
	return nil, errors.New(expense.CategoryAlreadyExists)
}
