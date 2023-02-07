package tracking

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"time"
)

type CreateExpenseResp struct {
	ID  string
	Msg string
}

type CreateExpenseReq struct {
	Product  string
	Price    float64
	Currency string
	Place    string
	City     string
	Date     time.Time
	People   string
	Category string
}

// ExpenseCreator use case creates a category for a expense
type ExpenseCreator struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewExpenseCreator(l app.Logger, e expense.Expenses) *ExpenseCreator {
	return &ExpenseCreator{l, e}
}

// Create use cases function creates a new category
func (s *ExpenseCreator) Create(req CreateExpenseReq) (*CreateExpenseResp, error) {
	price := expense.Price{
		Currency: req.Currency,
		Amount:   req.Price,
	}
	place := expense.Place{
		City: req.Place,
		Shop: req.City,
	}
	newExpense, createErr := expense.NewExpense(price, req.Product, req.People, place, req.Date, req.Category)
	if createErr != nil {
		s.logger.Debug("Failed to validate expense %s: %v", req, createErr)
		return nil, createErr
	}
	err := s.expenses.Add(*newExpense)
	if err != nil {
		s.logger.Err("Could add expense: %s", err)
		return nil, err
	}
	s.logger.Info("Expense %s created: %s", newExpense.ID, newExpense)
	return &CreateExpenseResp{
		ID:  string(newExpense.ID),
		Msg: fmt.Sprintf("Expense %s created: %v", newExpense.ID, newExpense),
	}, nil
}
