package expense

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator"
)

// ID is the unique identifier for the domain objects of type Expense
type ID uint32

// Expense is the aggregate root for other entities such as Category
type Expense struct {
	ID      ID        `validate:"required,min=3,max=32"`
	Product string    `validate:"required,min=3,max=32"`
	Shop    string    `validate:"required,min=3,max=32"`
	Date    time.Time `validate:"required,min=3,max=32"`
	City    string    `validate:"required,min=3,max=32"`
	Town    string    `validate:"required,min=3,max=32"`

	Category Category
}

// CategoryID is the unique identifier for the domain object of type Category
type CategoryID string

// Category is an entity that is supposed to be accessed only from the Expense aggregate
type Category struct {
	ID   CategoryID
	Name string
}

// Expenses is the repository for all the command actions for Expense
type Expenses interface {
	// Add is used to add a new Expense to the system
	Add(e Expense) error
	// Delete is used to remove a Expense from the system
	Delete(id ID) error
	// Add is used to save a new category for future expenses
	SaveCategory(c Category) error
	// Add is used to save a new category for future expenses
	DeleteCategory(id CategoryID) error
}

func (c *Expense) validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}
		return errors.New(fmt.Sprintf("Invalid category data: %v", err))
	}
	return nil

}
