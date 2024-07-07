package querying

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"

	"github.com/google/uuid"
)

const NoUserID = "00000000-0000-0000-0000-000000000000"

type ExpensesBasics struct {
	Amount   float64
	Category struct {
		Name string
		ID   string
	}
	Users map[string]struct {
		DisplayName      string
		TelegramUsername string
	}
	Date    time.Time
	ID      string
	Product string
	Shop    string
}

type ExpenseQuerierResp struct {
	Expenses      []ExpensesBasics
	Page          uint
	PageSize      uint
	ExpensesCount uint
}

type ExpenseQuerierFilter struct {
	ByCategoryID []string
	ByUsers      []string
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
	expenses expense.Expenses
	users    expense.Users
}

func NewExpenseQuerier(l app.Logger, e expense.Expenses, u expense.Users) *ExpenseQuerier {
	return &ExpenseQuerier{l, e, u}
}

func (s *ExpenseQuerier) GetByID(id string) (*ExpenseQuerierResp, error) {
	s.logger.Info("Getting expense " + id)
	idE, err := uuid.Parse(id)
	if err != nil {
		s.logger.Err("Failed parsing id: %v", err)
		return nil, expense.ErrInvalidID
	}
	e, err := s.expenses.Get(idE)
	if err != nil {
		s.logger.Err("Could not get expense %s from storage: %v", idE, err)
		return nil, expense.ErrNotFound
	}
	expBasic := ExpensesBasics{
		Amount: e.Amount,
		Category: struct {
			Name string
			ID   string
		}{e.Category.Name, e.Category.ID.String()},
		Users: map[string]struct {
			DisplayName      string
			TelegramUsername string
		}{},
		Date:    e.Date,
		ID:      e.ID.String(),
		Product: e.Product,
		Shop:    e.Shop,
	}
	expenseMap := map[string]ExpensesBasics{e.ID.String(): expBasic}
	resp := ExpenseQuerierResp{
		Expenses:      []ExpensesBasics{expBasic},
		Page:          1,
		PageSize:      1,
		ExpensesCount: 1,
	}
	users, err := s.users.All()
	if err != nil {
		return nil, err
	}
	for _, uid := range e.UserIDS {
		for _, u := range users {
			if u.ID == uid {
				user := struct {
					DisplayName      string
					TelegramUsername string
				}{u.DisplayName, u.TelegramUsername}
				expenseMap[e.ID.String()].Users[uid.String()] = user
			}
		}
	}
	resp.Expenses = append(resp.Expenses, expenseMap[e.ID.String()])
	return &resp, nil
}

func (s *ExpenseQuerier) Query(req ExpenseQuerierReq) (*ExpenseQuerierResp, error) {
	var expenses []expense.Expense
	var err error
	s.logger.Info("Getting all expenses")
	totalExpenses, err := s.expenses.CountWithFilter(req.ExpenseFilter.ByUsers, req.ExpenseFilter.ByCategoryID, req.ExpenseFilter.ByAmount[0], req.ExpenseFilter.ByAmount[1], req.ExpenseFilter.ByShop, req.ExpenseFilter.ByProduct, req.ExpenseFilter.ByTime[0], req.ExpenseFilter.ByTime[1])
	if err != nil {
		s.logger.Err("Could count expenses storage: %v", err)
		return nil, err
	}
	s.logger.Debug("Total Filtered expenses", totalExpenses)
	expenses, err = s.expenses.Filter(req.ExpenseFilter.ByUsers, req.ExpenseFilter.ByCategoryID, req.ExpenseFilter.ByAmount[0], req.ExpenseFilter.ByAmount[1], req.ExpenseFilter.ByShop, req.ExpenseFilter.ByProduct, req.ExpenseFilter.ByTime[0], req.ExpenseFilter.ByTime[1], req.MaxPageSize, req.Page*req.MaxPageSize)
	if err != nil {
		s.logger.Err("Could not get expenses from storage: %v", err)
		return nil, err
	}
	users, err := s.users.All()
	if err != nil {
		return nil, err
	}
	resp := ExpenseQuerierResp{
		Expenses:      []ExpensesBasics{},
		Page:          req.Page,
		ExpensesCount: totalExpenses,
		PageSize:      uint(len(expenses)),
	}
	expenseMap := map[string]ExpensesBasics{}
	for _, e := range expenses {
		expBasic := ExpensesBasics{
			Amount: e.Amount,
			Category: struct {
				Name string
				ID   string
			}{e.Category.Name, e.Category.ID.String()},
			Users: map[string]struct {
				DisplayName      string
				TelegramUsername string
			}{},
			Date:    e.Date,
			ID:      e.ID.String(),
			Product: e.Product,
			Shop:    e.Shop,
		}
		expenseMap[e.ID.String()] = expBasic
		for _, u := range users {
			for _, uid := range e.UserIDS {
				if u.ID == uid {
					user := struct {
						DisplayName      string
						TelegramUsername string
					}{u.DisplayName, u.TelegramUsername}
					expenseMap[e.ID.String()].Users[uid.String()] = user
				}
			}
		}
		resp.Expenses = append(resp.Expenses, expenseMap[e.ID.String()])
	}
	resp.PageSize = uint(len(resp.Expenses))
	return &resp, nil
}
