package managing

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"fmt"
)

type UpdateCategoryResp struct {
	ID string
}

type UpdateCategoryReq struct {
	ID      string
	NewName string
}

type CategoryUpdater struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewCategoryUpdater(l app.Logger, e expense.Expenses) *CategoryUpdater {
	return &CategoryUpdater{l, e}
}

func (s *CategoryUpdater) Update(req UpdateCategoryReq) (*UpdateCategoryResp, error) {
	s.logger.Debug("Updating category %s", req.ID, "name", req.NewName)
	cat, err := s.expenses.GetCategory(expense.CategoryID(req.ID))
	if errors.Is(err, expense.ErrNotFound) {
		s.logger.Err("Category not found %s", req.ID)
		return nil, fmt.Errorf("Category %s not found", req.ID)
	}
	if err != nil {
		s.logger.Err("Could not find category %s: %v", req.ID, err)
		return nil, errors.New("Could not update category information.")
	}
	s.logger.Debug("Updating category id:%s new_name:%s", req.ID, req.NewName)
	cat.Name = expense.CategoryName(cat.Name)
	err = s.expenses.UpdateCategory(*cat)
	if err != nil {
		s.logger.Err(fmt.Sprintf("Error updating category %s", req.ID), err)
		return nil, err
	}
	resp := &UpdateCategoryResp{
		ID: req.ID,
	}
	s.logger.Info("Category %s updated", req.ID)
	return resp, nil
}
