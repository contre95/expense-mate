package managing

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"time"
)

type DeleteCategoryResp struct {
	ID string
}

type DeleteCategoryReq struct {
	ID string
}

type CategoryDeleter struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewCategoryDeleter(l app.Logger, e expense.Expenses) *CategoryDeleter {
	return &CategoryDeleter{l, e}
}

func (s *CategoryDeleter) Delete(req DeleteCategoryReq) (*DeleteCategoryResp, error) {
	i, err := s.expenses.CountWithFilter([]string{req.ID}, 0, 0, "", "", time.Time{}, time.Time{})
	if err != nil {
		s.logger.Err(fmt.Sprintf("Could count expenses for category %s", req.ID), err)
		return nil, errors.New("Could count expenses for category %s.")
	}
	s.logger.Debug("Amount of expenses associated with %s : %d", req.ID, i)
	if i != 0 {
		return nil, errors.New(fmt.Sprintf("Could not delete category. %d expenses are still associated, please delete them.", i))
	}
	err = s.expenses.DeleteCategory(expense.CategoryID(req.ID))
	if err != nil {
		s.logger.Err("Error deleting category", err)
		return nil, err
	}
	resp := &DeleteCategoryResp{
		ID: req.ID,
	}
	s.logger.Info("Category %s deleted", req.ID)
	return resp, nil
}
