package tracking

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"time"
)

// Same as createExpense.go, but for updating expenses

type UpdateExpenseResp struct {
	ExpenseID string
}

type UpdateExpenseReq struct {
	Amount     float64
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
	oldExpense, getErr := s.expenses.Get(expense.ID(req.ExpenseID))
	fmt.Println(oldExpense)
	switch {
	case errors.Is(getErr, expense.ErrNotFound):
		s.logger.Debug("Expense %s not found in storage: %v", req.ExpenseID, getErr)
	case getErr != nil:
		s.logger.Debug("Failed to update expense %s: %v", req.ExpenseID, getErr)
		return nil, getErr
	}
	// Update the values
	oldExpense.Amount = req.Amount
	oldExpense.Shop = req.Shop
	oldExpense.Product = req.Product
	oldExpense.Date = req.Date
	newCategory, err := s.expenses.GetCategory(expense.CategoryID(req.CategoryID))
	switch {
	case errors.Is(err, expense.ErrNotFound):
		s.logger.Err("The category you are trying to reach doesn't exists", req, err)
		return nil, err
	case err != nil:
		s.logger.Err("Could not get category.", req, err)
		return nil, err
	}
	oldExpense.Category = *newCategory
	oldExpense, err = oldExpense.Validate()
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
