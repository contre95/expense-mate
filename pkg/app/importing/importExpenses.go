package importing

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"
)

type ImporertResp struct {
	ID                string
	Msg               string
	SuccesfullImports int
	FailedImports     int
}

type ImporertReq struct {
	ByPassWrongExpenses bool
}

// IImportedExpense holds the values that should be imported by the Importer
type ImportedExpense struct {
	Amount   float32
	Currency string
	Product  string
	Shop     string
	Date     time.Time
	City     string
	Town     string

	Category string
}

// Importer is the main dependency of the ImportExpensesUseCase and defines how an importer should behave
type Importer interface {
	//GetAllCategories() ([]string, error)
	GetImportedExpenses() ([]ImportedExpense, error)
}

// The createCategory use case creates a category for a expense
type ImportExpensesUseCase struct {
	logger   app.Logger
	importer Importer
	expenses expense.Expenses
}

// Contructor for Import
func NewImporterUseCase(l app.Logger, i Importer, e expense.Expenses) *ImportExpensesUseCase {
	return &ImportExpensesUseCase{l, i, e}
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
	return expense.NewExpense(price, e.Product, place, e.Date, e.Category)
}

// Import imports a all the categories provided by the importer
func (u *ImportExpensesUseCase) Import(req ImporertReq) (*ImporertResp, error) {
	importedExpenses, err := u.importer.GetImportedExpenses()
	if err != nil {
		u.logger.Err("Could not import expenses: %s", err)
		return nil, err
	}
	var expensesToAdd []expense.Expense
	for _, e := range importedExpenses {
		newExp, err := parseExpense(e)
		if err != nil && req.ByPassWrongExpenses {
			u.logger.Err("Could not import expense: %s of %d %s: %s", e.Product, e.Amount, e.Currency, err)
			if !req.ByPassWrongExpenses {
				return nil, err
			}
		} else {
			expensesToAdd = append(expensesToAdd, *newExp)
		}
	}
	for _, exp := range expensesToAdd {
		err := u.expenses.Add(exp)
		if err != nil {
			u.logger.Err("Failed to save expense %s : %s", exp.ID, err)
			return nil, err
		}
	}
	return &ImporertResp{
		SuccesfullImports: len(expensesToAdd),
		FailedImports:     len(importedExpenses) - len(expensesToAdd),
	}, nil
}
