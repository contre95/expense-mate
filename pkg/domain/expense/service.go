package expense

import (
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
func NewExpense(product, shop, city, town string, date time.Time, category Category) (*Expense, error) {
	expenseID := expenseID(hash(shop + town + string(category.ID) + date.String()))
	newExpense := Expense{
		ID:       expenseID,
		Product:  product,
		Shop:     shop,
		Date:     date,
		City:     city,
		Town:     town,
		Category: category,
	}
	err := newExpense.validate()
	if err != nil {
		return nil, err
	}
	return &newExpense, nil
}

// NewCategory creates a new category for expenses (Factory Method)
func NewCategory(name string) Category {
	return Category{
		ID:   CategoryID(strings.ReplaceAll(strings.ToLower(name), " ", "-")),
		Name: name,
	}
}
