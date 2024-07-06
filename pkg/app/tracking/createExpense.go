package tracking

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreateExpenseResp struct {
	ID  string
	Msg string
}

type CreateExpenseReq struct {
	Product    string
	Amount     float64
	Shop       string
	Date       time.Time
	UserIDS    []string
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
	idC, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, expense.ErrInvalidID
	}
	category, err := s.expenses.GetCategory(idC)
	if err != nil {
		return nil, errors.New("Couldn't fetch the category %s" + req.CategoryID)
	}
	newExpense, createErr := expense.NewExpense(req.Amount, req.Product, req.Shop, req.Date, *category)
	for _, sid := range req.UserIDS {
		pid, err := uuid.Parse(sid)
		if err != nil {
			s.logger.Err("Failed to parse UUID %s", err.Error())
			return nil, errors.New("Failed to parse UUID %s" + sid)
		}
		newExpense.UserIDS = append(newExpense.UserIDS, pid)
	}
	fmt.Println(newExpense.UserIDS)
	if createErr != nil {
		s.logger.Debug("Failed to validate expense %s: %v", req, createErr)
		return nil, createErr
	}
	err = s.expenses.Add(*newExpense)
	if err != nil {
		s.logger.Err("Could add expense: %s", err)
		return nil, err
	}
	s.logger.Info("Expense %s created: %s", newExpense.ID, newExpense)
	return &CreateExpenseResp{
		ID:  newExpense.ID.String(),
		Msg: fmt.Sprintf("Expense %s created: %v", newExpense.ID, newExpense),
	}, nil
}
