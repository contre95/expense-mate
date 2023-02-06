package querying

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"
)

type ExpensesBasics struct {
	ID       string
	Price    float64
	Product  string
	People   string
	Category string
}

type ExpenseQuerierResp struct {
	Expenses []ExpensesBasics
	Page     uint
	PageSize uint
}

type ExpenseQuerierReq struct {
	From     time.Time
	To       time.Time
	Page     uint
	PageSize uint
}

type ExpenseQuerier struct {
	logger   app.Logger
	expenses expense.Expenses // Expenses Repository
}

func NewExpenseQuerier(l app.Logger, e expense.Expenses) *ExpenseQuerier {
	return &ExpenseQuerier{l, e}
}

func (s *ExpenseQuerier) Query(req ExpenseQuerierReq) (*ExpenseQuerierResp, error) {
	s.logger.Info("Getting all expenses")
	expenses, err := s.expenses.GetFromTimeRange(req.From, req.To, req.PageSize, req.Page*req.PageSize)
	if err != nil {
		s.logger.Err("Could not get expenses from storage: %v", err)
		return nil, err
	}
	resp := ExpenseQuerierResp{
		Expenses: []ExpensesBasics{},
		Page:     req.Page,
	}
	for _, exp := range expenses {
		expBasic := ExpensesBasics{
			ID:       string(exp.ID),
			Price:    exp.Price.Amount,
			Product:  exp.Product,
			People:   exp.People,
			Category: string(exp.Category.Name),
		}
		resp.Expenses = append(resp.Expenses, expBasic)
	}
	resp.PageSize = uint(len(resp.Expenses))
	return &resp, nil
}
