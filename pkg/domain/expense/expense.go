package expense

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator"
)

// ID is the unique identifier for the domain objects of type Expense
type ID string

type Place struct {
	City string `validate:"min=2,max=64"`
	Town string `validate:"min=2,max=64"`
	Shop string `validate:"min=2,max=64"`
}
type Price struct {
	Currency string  `validate:"required"`
	Amount   float64 `validate:"required"`
}

// Expense is the aggregate root for other entities such as Category
type Expense struct {
	ID      ID `validate:"required"`
	Price   Price
	Place   Place
	Product string    `validate:"required,min=3"`
	Date    time.Time `validate:"required"`
	People  string    `validate:"required,min=3"`

	User     string
	Category Category
}

const CategoryNotFoundErr = "Category not found"
const CategoryAlreadyExists = "Category already exists"

// CategoryID is the unique identifier for the domain object of type Category
type CategoryID string

// CategoryName is type for the Name of a category. This value should be unique among all categories
type CategoryName string

// Category is an entity that is supposed to be accessed only from the Expense aggregate
type Category struct {
	ID   CategoryID   `validate:"required,min=3"`
	Name CategoryName `validate:"required,min=3,excludesall=!-@#"`
}

// Expenses is the repository for all the command actions for Expense
type Expenses interface {
	// Get is used to retrieve all expenses from a certain time range with the ability to "paginate" using limit and offset. By passing limit = 0 paginating will be dismissed.
	GetFromTimeRange(from, to time.Time, limit, offset uint) ([]Expense, error)
	// Count how many expenses are in a time range for a given category
	//Count(from, to *time.Time, categories []Category) (uint, error)
	// Add is used to add a new Expense to the system
	Add(e Expense) error
	// Delete is used to remove a Expense from the system
	Delete(id ID) error
	// Add is used to save a new category for future expenses
	GetCategories() ([]Category, error)
	// Creates a new category returns expense.CategoryAlreadyExistsErr if category is duplicated.
	AddCategory(c Category) error
	// Validates if a category exists
	CategoryExists(id CategoryID) (bool, error)
}

// Validate
func (c *Expense) Validate() (*Expense, error) {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}
		return nil, errors.New(fmt.Sprintf("Invalid expense data: %v", err))
	}
	return c, nil
}

func (c *Category) Validate() (*Category, error) {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}
		return nil, errors.New(fmt.Sprintf("Invalid category data: %v", err))
	}
	return c, nil
}
