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
	Product    string
	Amount     float64
	Currency   string
	Shop       string
	City       string
	Date       time.Time
	People     string
	CategoryID string
}

// ExpenseCreator use case creates a category for a expense
type ExpenseCreator struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewExpenseCreator(l app.Logger, e expense.Expenses) *ExpenseCreator {
	return &ExpenseCreator{l, e}
}

// Create use cases function creates a new expense
func (s *ExpenseCreator) Create(req CreateExpenseReq) (*CreateExpenseResp, error) {
	newExpense, createErr := expense.NewExpense(req.Amount, req.Product, req.Shop, req.Date, req.CategoryID)
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
