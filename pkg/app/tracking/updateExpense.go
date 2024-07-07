package tracking

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"

	"github.com/google/uuid"
)

// Same as createExpense.go, but for updating expenses

type UpdateExpenseResp struct {
	ExpenseID string
}

type UpdateExpenseReq struct {
	Amount     float64
	UsersID    []string
	CategoryID string
	Date       time.Time
	ExpenseID  string
	Product    string
	Shop       string
}

type ExpenseUpdater struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewExpenseUpdater(l app.Logger, e expense.Expenses) *ExpenseUpdater {
	return &ExpenseUpdater{l, e}
}

// Update use case function retrieves an expense from the db and updates it with the new values
func (s *ExpenseUpdater) Update(req UpdateExpenseReq) (*UpdateExpenseResp, error) {
	// Get the oldExpense from the db
	pidE, err := uuid.Parse(req.ExpenseID)
	if err != nil {
		s.logger.Err("Failed parse expense ID: %s", err.Error())
		return nil, expense.ErrInvalidID
	}
	oldExpense, getErr := s.expenses.Get(pidE)
	switch {
	case errors.Is(getErr, expense.ErrNotFound):
		s.logger.Err("Expense %s not found in storage: %v", req.ExpenseID, getErr)
		return nil, getErr
	case getErr != nil:
		s.logger.Err("Failed to update expense %s: %v", req.ExpenseID, getErr)
		return nil, getErr
	}
	s.logger.Debug("Old expense to update: %s", oldExpense)
	oldExpense.Amount = req.Amount
	oldExpense.Shop = req.Shop
	oldExpense.Product = req.Product
	oldExpense.Date = req.Date
	oldExpense.UsersID = []uuid.UUID{}
	for _, sid := range req.UsersID {
		pid, err := uuid.Parse(sid)
		if err != nil {
			s.logger.Err("Failed to parse UUID %s", err.Error())
			return nil, errors.New("Failed to parse UUID %s" + sid)
		}
		oldExpense.UsersID = append(oldExpense.UsersID, pid)
	}
	pidC, err := uuid.Parse(req.CategoryID)
	if err != nil {
		s.logger.Err("Failed to parse UUID %s", err.Error())
		return nil, expense.ErrInvalidID
	}
	newCategory, err := s.expenses.GetCategory(pidC)
	switch {
	case errors.Is(err, expense.ErrNotFound):
		s.logger.Err("The category you are trying to reach doesn't exists", req, err)
		return nil, err
	case err != nil:
		s.logger.Err("Failed to retrieve category from storage.", req, err)
		return nil, err
	}
	oldExpense.Category = *newCategory
	oldExpense, err = oldExpense.Validate()
	s.logger.Debug("Old expense updated: %s", oldExpense)
	if err != nil {
		s.logger.Err("Could not update expense", req, err)
		return nil, err
	}
	updateErr := s.expenses.Update(*oldExpense)
	if updateErr != nil {
		s.logger.Err("Could not update expense", req, err)
		return nil, err
	}
	return &UpdateExpenseResp{ExpenseID: req.ExpenseID}, nil
}
