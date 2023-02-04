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
func NewExpense(price Price, product, people string, place Place, date time.Time, cname string) (*Expense, error) {
	newCat, err := NewCategory(cname)
	if err != nil {
		return nil, err
	}
	newExpense := Expense{
		Price:    price,
		Product:  product,
		People:   people,
		Place:    place,
		Date:     date,
		Category: *newCat,
	}

	newExpense.ID = ID(hash(place.Shop + fmt.Sprintf("%f", price.Amount) + newExpense.People + newExpense.Product + string(newExpense.Category.ID) + date.String()))
	return newExpense.Validate()
}

// NewCategory creats a new category and validates the field. Still don't know if this can be created separately or always under an Expense ?
func NewCategory(name string) (*Category, error) {
	newCategory := Category{
		ID:   CategoryID(strings.ReplaceAll(strings.ToLower(name), " ", "-")),
		Name: CategoryName(name),
	}
	return newCategory.Validate()
}
