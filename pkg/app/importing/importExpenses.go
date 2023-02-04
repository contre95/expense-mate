package importing

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"time"
)

// ImportExpensesResp is the Response Model por the ImportExpenses use case
type ImportExpensesResp struct {
	Msg               string
	SuccesfullImports int
	FailedImports     int
}

// ImportExpensesReq is the Request Model por the ImportExpenses use case
type ImportExpensesReq struct {
	BypassWrongExpenses bool
	ReImport            bool
	ImporterID          string
}

// ImportedExpense holds the values that should be imported by the Importer
type ImportedExpense struct {
	Amount   float64
	Currency string
	Product  string
	Shop     string
	Date     time.Time
	City     string
	Town     string
	People   string

	Category string
}

// Importer is the main dependency of the ImportExpenses and defines how an importer should behave
type Importer interface {
	//GetAllCategories() ([]string, error)
	GetImportedExpenses() ([]ImportedExpense, error)
}

// The ImportExpenses use case creates a category for a expense
type ImportExpenses struct {
	logger    app.Logger
	importers map[string]Importer
	expenses  expense.Expenses
}

// NewExpenseImporter returns a valid ExpenseImporter use case
func NewExpenseImporter(l app.Logger, i map[string]Importer, e expense.Expenses) *ImportExpenses {
	return &ImportExpenses{l, i, e}
}

func parseExpense(e ImportedExpense) (*expense.Expense, error) {
	price := expense.Price{
		Currency: e.Currency,
		Amount:   e.Amount,
	}
	place := expense.Place{
		City: e.City,
		Town: e.Town,
		Shop: e.Shop,
	}
	return expense.NewExpense(price, e.Product, e.People, place, e.Date, e.Category)
}

// Import imports a all the categories provided by the importer
func (u *ImportExpenses) Import(req ImportExpensesReq) (*ImportExpensesResp, error) {
	importedExpenses, err := u.importers[req.ImporterID].GetImportedExpenses()
	if err != nil {
		u.logger.Err("Could not import expenses: %s", err)
		return nil, errors.New("Could not import expenses from importer" + req.ImporterID)
	}
	failedExpenses := 0
	for _, e := range importedExpenses {
		newExp, err := parseExpense(e)
		if err != nil {
			failedExpenses++
			u.logger.Err("Could not import expense: %s of %f %s: %s", e.Product, e.Amount, e.Currency, err)
			if !req.BypassWrongExpenses {
				fmt.Println(req.BypassWrongExpenses)
				return nil, fmt.Errorf("Failed to import expense: %s of %f %s", e.Product, e.Amount, e.Currency)
			}
		}
		err = u.expenses.Add(*newExp)
		if err != nil {
			failedExpenses++
			u.logger.Err("Failed to save expense %s : %s", newExp.ID, err)
			if !req.BypassWrongExpenses {
				fmt.Println(req.BypassWrongExpenses)
				return nil, fmt.Errorf("Failed to save expense %d : %s", newExp.ID, err)
			}
		}
	}
	var msg string
	if failedExpenses == 0 {
		msg = "All the expenses where imported"
	} else {
		msg = "Some expenses could not be imported"
	}
	return &ImportExpensesResp{
		SuccesfullImports: len(importedExpenses) - failedExpenses,
		FailedImports:     failedExpenses,
		Msg:               msg,
	}, nil
}
