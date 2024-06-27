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
	Expenses map[string]ExpensesBasics
	Page     uint
	PageSize uint
}

type ExpenseQuerierFilter struct {
	ByCategoryID []string
	ByShop       string
	ByProduct    string
	ByPrice      [2]uint
	ByTime       [2]time.Time
}

type ExpenseQuerierReq struct {
	Page          uint
	MaxPageSize   uint
	ExpenseFilter ExpenseQuerierFilter
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
		Expenses: map[string]ExpensesBasics{
			string(e.ID): {
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
	var expenses []expense.Expense
	var err error
	s.logger.Info("Getting all expenses")
	// expenses, err = s.expenses.All(req.MaxPageSize, req.Page*req.MaxPageSize)
	expenses, err = s.expenses.Filter(req.ExpenseFilter.ByCategoryID, req.ExpenseFilter.ByPrice[0], req.ExpenseFilter.ByPrice[1], req.ExpenseFilter.ByShop, req.ExpenseFilter.ByProduct, req.ExpenseFilter.ByTime[0], req.ExpenseFilter.ByTime[1], req.MaxPageSize, req.Page*req.MaxPageSize)
	if err != nil {
		s.logger.Err("Could not get expenses from storage: %v", err)
		return nil, err
	}
	resp := ExpenseQuerierResp{
		Expenses: map[string]ExpensesBasics{},
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
		resp.Expenses[string(expBasic.ID)] = expBasic
	}
	resp.PageSize = uint(len(resp.Expenses))
	return &resp, nil
}
