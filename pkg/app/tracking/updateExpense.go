package tracking

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"
)

// Same as createExpense.go, but for updating expenses

type UpdateExpenseResp struct {
	ExpenseID string
}

type UpdateExpenseReq struct {
	ExpenseID  string
	Amount     float64
	Product    string
	Price      float64
	City       string
	Currency   string
	Shop       string
	Date       time.Time
	People     string
	CategoryID string
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
	if getErr != nil {
		s.logger.Debug("Failed to update expense %s: %v", req, getErr)
		return nil, getErr
	}
	// Update the values
	oldExpense.Price.Amount = req.Price
	oldExpense.Price.Currency = req.Currency
	oldExpense.Place.City = req.City
	oldExpense.Place.Shop = req.Shop
	oldExpense.Product = req.Product
	oldExpense.Date = req.Date
	oldExpense.People = req.People
	newCategory, err := s.expenses.GetCategory(expense.CategoryID(req.CategoryID))
	switch {
	case errors.Is(err, expense.ErrNotFound):
		s.logger.Err("The category you are trying to reach doesn't exists", req, err)
		return nil, err
	case err != nil:
		s.logger.Err("Could not update category.", req, err)
		return nil, err
	}
	oldExpense.Category = *newCategory
	s.expenses.Update(*oldExpense)
	return &UpdateExpenseResp{ExpenseID: req.ExpenseID}, nil
}
