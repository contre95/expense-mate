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
	err := newExpense.validate()
	if err != nil {
		return nil, err
	}
	return &newExpense, nil
}

func NewCategory(name string) (*Category, error) {
	newCategory := Category{
		ID:   CategoryID(strings.ReplaceAll(strings.ToLower(name), " ", "-")),
		Name: CategoryName(name),
	}
	if err := newCategory.validate(); err != nil {
		return nil, newCategory.validate()
	}
	return &newCategory, nil
}

// NewCategory creates a new category for an expense (Factory Method)
func (e *Expense) Categorize(name string) {
}
