package expense

import (
	"time"

	"github.com/google/uuid"
)

// NewExpense acts a Factory Method for new Expenses enforcinf invariants for the Expense entity
func NewExpense(amount float64, product string, shop string, date time.Time, category Category) (*Expense, error) {
	newExpense := Expense{
		Amount:   amount,
		Product:  product,
		Shop:     shop,
		Date:     date,
		Category: category, // I'm asking for a valid category when I create an user
	}

	newExpense.ID = uuid.New()
	return newExpense.Validate()
}

// NewCategory creates a new category and validates the field. Still don't know if this can be created separately or always under an Expense ?
func NewCategory(name string) (*Category, error) {
	newCategory := Category{
		ID:   CategoryID(uuid.New()),
		Name: name,
	}
	return newCategory.Validate()
}

func NewRule(p string, uids []uuid.UUID, cid CategoryID) (*Rule, error) {
	// TODO: Check if it is a valid regex pattern
	rule := &Rule{uuid.New().String(), p, cid, uids}
	return rule.Validate()
}
