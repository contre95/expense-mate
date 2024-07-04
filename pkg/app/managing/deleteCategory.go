package managing

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	rules    expense.Rules
}

func NewCategoryDeleter(l app.Logger, e expense.Expenses, r expense.Rules) *CategoryDeleter {
	return &CategoryDeleter{l, e, r}
}

func (s *CategoryDeleter) Delete(req DeleteCategoryReq) (*DeleteCategoryResp, error) {
	catID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, expense.ErrInvalidID
	}
	i, err := s.expenses.CountWithFilter([]string{req.ID}, 0, 0, "", "", time.Time{}, time.Time{})
	if err != nil {
		s.logger.Err(fmt.Sprintf("Could count expenses for category %s", req.ID), err)
		return nil, errors.New("Could count expenses for category %s.")
	}
	s.logger.Debug("Amount of expenses associated with %s : %d", req.ID, i)
	if i != 0 {
		return nil, errors.New(fmt.Sprintf("Could not delete category. %d expenses are still associated, please delete them.", i))
	}
	// Note: Making this verification cause SQLite ON CASCADE DELETE doesn't work :(
	// Also a good things, cause it is independent of the storage engine and dependant only on the repository interface
	rules, err := s.rules.All()
	if err != nil {
		s.logger.Err(fmt.Sprintf("Could count rules for category %s", req.ID), err)
		return nil, errors.New("Could count rules for category %s.")
	}
	for _, r := range rules {
		// if r.CategoryID.String() == catID.String() {
		if r.CategoryID == catID {
			return nil, errors.New("Could not delete category. Rules are still associated, please delete them.")
		}
	}
	err = s.expenses.DeleteCategory(expense.CategoryID(catID))
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
