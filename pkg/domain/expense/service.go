package expense

import (
	"fmt"
	"hash/fnv"
	"strings"
	"time"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// NewExpense acts a Factory Method for new Expenses enforcinf invariants for the Expense entity
func NewExpense(price Price, product, people string, place Place, date time.Time, category string) (*Expense, error) {
	newExpense := Expense{
		Price:   price,
		Product: product,
		People:  people,
		Place:   place,
		Date:    date,
	}

	newExpense.ID = ID(hash(place.Shop + fmt.Sprintf("%f", price.Amount) + newExpense.People + newExpense.Product + string(newExpense.Category.ID) + date.String()))
	newExpense.Categorize(category)
	err := newExpense.validate()
	if err != nil {
		return nil, err
	}
	return &newExpense, nil
}

// NewCategory creates a new category for an expense (Factory Method)
func (e *Expense) Categorize(name string) {
	e.Category = Category{
		ID:   CategoryID(strings.ReplaceAll(strings.ToLower(name), " ", "-")),
		Name: CategoryName(name),
	}
}
