package querying

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"fmt"
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
	e, err := s.expenses.Get(expense.ID(id))
	if err != nil {
		s.logger.Err("Could not get expense from storage: %v", err)
		return nil, expense.ErrNotFound
	}
	resp := ExpenseQuerierResp{
		Expenses: []ExpensesBasics{
			{
				Amount:     e.Price.Amount,
				Category:   string(e.Category.Name),
				CategoryID: id,
				Date:       e.Date,
				ID:         string(e.ID),
				People:     e.People,
				Product:    e.Product,
				Shop:       e.Place.Shop,
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
			Amount:     exp.Price.Amount,
			Category:   string(exp.Category.Name),
			CategoryID: string(exp.Category.ID),
			Date:       exp.Date,
			ID:         string(exp.ID),
			People:     exp.People,
			Product:    exp.Product,
			Shop:       exp.Place.Shop,
		}
		resp.Expenses = append(resp.Expenses, expBasic)
	}
	fmt.Println(resp)
	resp.PageSize = uint(len(resp.Expenses))
	return &resp, nil
}
