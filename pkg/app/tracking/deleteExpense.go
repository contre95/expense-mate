package tracking

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"fmt"
)

type DeleteExpenseResp struct {
	SuccessfulDeletes []string
	FailedDeletes     map[string]error
}

type DeleteExpenseReq struct {
	IDS []string
}

type ExpenseDeleter struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewExpenseDeleter(l app.Logger, e expense.Expenses) *ExpenseDeleter {
	return &ExpenseDeleter{l, e}
}

func (s *ExpenseDeleter) Delete(req DeleteExpenseReq) (*DeleteExpenseResp, error) {
	resp := &DeleteExpenseResp{}
	for _, id := range req.IDS {
		s.logger.Debug("Attempting to delte Expense ", string(id))
		err := s.expenses.Delete(expense.ID(id))
		if err != nil {
			s.logger.Err("Error updating client", err)
			resp.FailedDeletes[id] = err
		} else {
			s.logger.Debug("Expense delted successfully", err)
			resp.SuccessfulDeletes = append(resp.SuccessfulDeletes, id)
		}
	}
	if len(resp.FailedDeletes) > 0 {
		return nil, errors.New(fmt.Sprintf("Could not delete any of %d the expenses.", len(resp.FailedDeletes)))
	}
	s.logger.Debug("Deleted: %d, Failed: %d", len(resp.SuccessfulDeletes), len(resp.FailedDeletes))
	return resp, nil
}
