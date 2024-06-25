package querying

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"
)

type ExpensesBasics struct {
	Amount     float64
	Category   string
	CategoryID string
	Date       time.Time
	ID         string
	People     string
	Product    string
	Shop       string
}

type ExpenseQuerierResp struct {
	Expenses []ExpensesBasics
	Page     uint
	PageSize uint
}

type ExpenseQuerierReq struct {
	From        time.Time
	To          time.Time
	Page        uint
	MaxPageSize uint
}

type ExpenseQuerier struct {
	logger   app.Logger
	expenses expense.Expenses // Expenses Repository
}

func NewExpenseQuerier(l app.Logger, e expense.Expenses) *ExpenseQuerier {
	return &ExpenseQuerier{l, e}
}

func (s *ExpenseQuerier) GetByID(id string) (*ExpenseQuerierResp, error) {
	s.logger.Info("Getting expense " + id)
	expense, err := s.expenses.Get(expense.ID(id))
	if err != nil {
	}
	resp := ExpenseQuerierResp{
		Expenses: []ExpensesBasics{
			{
				ID:       string(expense.ID),
				Date:     expense.Date,
				Amount:   expense.Price.Amount,
				Product:  expense.Product,
				People:   expense.People,
				Category: string(expense.Category.Name),
			},
		},
		Page:     0,
		PageSize: 1,
	}
	return &resp, nil
}

func (s *ExpenseQuerier) Query(req ExpenseQuerierReq) (*ExpenseQuerierResp, error) {
	s.logger.Info("Getting all expenses")
	expenses, err := s.expenses.GetFromTimeRange(req.From, req.To, req.MaxPageSize, req.Page*req.MaxPageSize)
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
			Date:     exp.Date,
			Amount:   exp.Price.Amount,
			Product:  exp.Product,
			People:   exp.People,
			Category: string(exp.Category.Name),
		}
		resp.Expenses = append(resp.Expenses, expBasic)
	}
	resp.PageSize = uint(len(resp.Expenses))
	return &resp, nil
}
