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
func NewExpense(price Price, product string, place Place, date time.Time, category string) (*Expense, error) {
	newExpense := Expense{
		Product: product,
		Price:   price,
		Place:   place,
		Date:    date,
		Category: Category{
			ID:   CategoryID(strings.ReplaceAll(strings.ToLower(category), " ", "-")),
			Name: CategoryName(category),
		},
	}
	newExpense.ID = ID(hash(place.Shop + place.Town + string(newExpense.Category.ID) + date.String()))
	fmt.Println(newExpense.ID)
	err := newExpense.validate()
	if err != nil {
		return nil, err
	}
	return &newExpense, nil
}

// NewCategory creates a new category for expenses (Factory Method)
func (e *Expense) NewCategory(name string) Category {
	return Category{
		ID:   CategoryID(strings.ReplaceAll(strings.ToLower(name), " ", "-")),
		Name: CategoryName(name),
	}
}
