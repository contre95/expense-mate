package expense

import (
	"strings"
	"time"
)

type CategoryID string

// Category is an entity that is supposed to be accessed only from the Expense aggregate
type Category struct {
	ID   CategoryID
	Name string
}

type ExpenseID uint

// Expense is the aggregate root for other entities such as Category
type Expense struct {
	ID      ExpenseID
	Product string
	Shop    string
	Date    time.Time
	City    string

	Category Category
}

// Expenses is the repository for all the command actions for Expense
type Expenses interface {
	// Add is used to add a new Expense to the system
	Add(e Expense) error
	Delete(id ExpenseID) error
	// Add is used to save a new category for future expenses
	SaveCategory(c Category) error
	// Add is used to save a new category for future expenses
	DeleteCategory(id CategoryID) error
}

// NewCategory creates a new category for expenses
func NewCategory(name string) Category {
	return Category{
		ID:   CategoryID(strings.ReplaceAll(strings.ToLower(name), " ", "-")),
		Name: name,
	}
}
