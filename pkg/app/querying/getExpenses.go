package querying

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"

	"github.com/google/uuid"
)

type ExpensesBasics struct {
	Amount   float64
	Category struct {
		Name string
		ID   string
	}
	Date    time.Time
	ID      string
	Product string
	Shop    string
}

type ExpenseQuerierResp struct {
	Expenses      map[string]ExpensesBasics
	Page          uint
	PageSize      uint
	ExpensesCount uint
}

type ExpenseQuerierFilter struct {
	ByCategoryID []string
	ByShop       string
	ByProduct    string
	ByAmount     [2]uint
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
	idE, err := uuid.Parse(id)
	if err != nil {
		return nil, expense.ErrInvalidID
	}
	e, err := s.expenses.Get(idE)
	if err != nil {
		s.logger.Err("Could not get expense from storage: %v", err)
		return nil, expense.ErrNotFound
	}
	// Get users as well

	resp := ExpenseQuerierResp{
		Expenses: map[string]ExpensesBasics{
			e.ID.String(): {
				Amount: e.Amount,
				Category: struct {
					Name string
					ID   string
				}{
					e.Category.Name,
					e.Category.ID.String(),
				},
				Date:    e.Date,
				ID:      e.ID.String(),
				Product: e.Product,
				Shop:    e.Shop,
			},
		},
		Page:     1,
		PageSize: 1,
	}
	return &resp, nil
}

func (s *ExpenseQuerier) Query(req ExpenseQuerierReq) (*ExpenseQuerierResp, error) {
	var expenses []expense.Expense
	var err error
	s.logger.Info("Getting all expenses")
	totalExpenses, err := s.expenses.CountWithFilter(req.ExpenseFilter.ByCategoryID, req.ExpenseFilter.ByAmount[0], req.ExpenseFilter.ByAmount[1], req.ExpenseFilter.ByShop, req.ExpenseFilter.ByProduct, req.ExpenseFilter.ByTime[0], req.ExpenseFilter.ByTime[1])
	if err != nil {
		s.logger.Err("Could count expenses storage: %v", err)
		return nil, err
	}
	s.logger.Debug("Total Filtered expenses", totalExpenses)
	expenses, err = s.expenses.Filter(req.ExpenseFilter.ByCategoryID, req.ExpenseFilter.ByAmount[0], req.ExpenseFilter.ByAmount[1], req.ExpenseFilter.ByShop, req.ExpenseFilter.ByProduct, req.ExpenseFilter.ByTime[0], req.ExpenseFilter.ByTime[1], req.MaxPageSize, req.Page*req.MaxPageSize)
	if err != nil {
		s.logger.Err("Could not get expenses from storage: %v", err)
		return nil, err
	}
	resp := ExpenseQuerierResp{
		Expenses:      map[string]ExpensesBasics{},
		Page:          req.Page,
		ExpensesCount: totalExpenses,
		PageSize:      uint(len(expenses)),
	}
	for _, e := range expenses {
		expBasic := ExpensesBasics{
			Amount: e.Amount,
			Category: struct {
				Name string
				ID   string
			}{
				e.Category.Name,
				e.Category.ID.String(),
			},

			Date:    e.Date,
			ID:      e.ID.String(),
			Product: e.Product,
			Shop:    e.Shop,
		}
		resp.Expenses[string(expBasic.ID)] = expBasic
	}
	resp.PageSize = uint(len(resp.Expenses))
	return &resp, nil
}
