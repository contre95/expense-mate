package expense

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// NewExpense acts a Factory Method for new Expenses enforcinf invariants for the Expense entity
func NewExpense(amount float64, product string, shop string, date time.Time, cname string) (*Expense, error) {
	newCat, err := NewCategory(cname)
	if err != nil {
		return nil, err
	}
	newExpense := Expense{
		Amount:   amount,
		Product:  product,
		Shop:     shop,
		Date:     date,
		Category: *newCat,
	}

	newExpense.ID = ID(uuid.New().String())
	return newExpense.Validate()
}

// NewCategory creates a new category and validates the field. Still don't know if this can be created separately or always under an Expense ?
func NewCategory(name string) (*Category, error) {
	newCategory := Category{
		ID:   CategoryID(strings.ReplaceAll(strings.ToLower(name), " ", "-")),
		Name: CategoryName(name),
	}
	return newCategory.Validate()
}
