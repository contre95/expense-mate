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
	Town string
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

var (
	ErrNotFound      = errors.New("The resource you are trying to get does not exist")
	ErrAlreadyExists = errors.New("The resource you are trying to get already exists")
	ErrInvalidEntity = errors.New("The entity you are trying create is not valid")
)

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
	// All retrieves all Expenses with pagination
	Filter(categories []string, minPrice, maxPrice uint, shop, product string, from time.Time, to time.Time, limit, offset uint) ([]Expense, error)
	// All retrieves all Expenses with pagination
	All(limit, offset uint) ([]Expense, error)
	// Get retrieves an Expense from storage
	CountWithFilter(categories []string, minPrice, maxPrice uint, shop, product string, from time.Time, to time.Time) (uint, error)
	// Get retrieves an Expense from storage
	Get(id ID) (*Expense, error)
	// // Add is used to add a new Expense to the system
	Add(e Expense) error
	// Delete is used to remove an Expense
	Delete(id ID) error
	// Update is used to update an Expnese
	Update(Expense) error
	// Add is used to save a new category for future expenses
	GetCategories() ([]Category, error)
	// Creates a new category returns expense.CategoryAlreadyExistsErr if category is duplicated.
	AddCategory(c Category) error
	// Validates if a category exists
	GetCategory(id CategoryID) (*Category, error)
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
