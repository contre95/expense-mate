package expense

import (
	"errors"
	"time"
)

// ID is the unique identifier for the domain objects of type Expense
type ID string

// Expense is the aggregate root for other entities such as Category
type Expense struct {
	ID      ID        `validate:"required"`
	Amount  float64   `validate:"required"`
	Product string    `validate:"required,min=3"`
	Shop    string    `validate:"min=2,max=64"`
	Date    time.Time `validate:"required"`

	User     string
	Category Category
}

var (
	ErrNotFound      = errors.New("The resource you are trying to get does not exist")
	ErrAlreadyExists = errors.New("The resource you are trying to create already exists")
	ErrInvalidEntity = errors.New("The entity you are trying create has invalid fields.")
)
var (
	ErrInvalidAmount   = errors.New("The amount of the expense is invalid")
	ErrInvalidProduct  = errors.New("The product name is invalid")
	ErrInvalidShop     = errors.New("The shop name is invalid")
	ErrInvalidDate     = errors.New("The date of the expense is invalid")
	ErrInvalidCategory = errors.New("The category of the expense is invalid")
)

// CategoryID is the unique identifier for the domain object of type Category
type CategoryID string

// CategoryName is type for the Name of a category. This value should be unique among all categories
type CategoryName string

// Category is an entity that is supposed to be accessed only from the Expense aggregate
type Category struct {
	ID   CategoryID   `validate:"required,min=3"`
	Name CategoryName `validate:"required,min=3,alphanum_space"`
}

// Expenses is the repository for all the command actions for Expense
type Expenses interface {
	// All retrieves all Expenses with pagination
	Filter(categories []string, minAmount, maxAmount uint, shop, product string, from time.Time, to time.Time, limit, offset uint) ([]Expense, error)
	// All retrieves all Expenses with pagination
	All(limit, offset uint) ([]Expense, error)
	// Get retrieves an Expense from storage
	CountWithFilter(categories []string, minAmount, maxAmount uint, shop, product string, from time.Time, to time.Time) (uint, error)
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
	// DeleteCategory is used to remove a category
	DeleteCategory(id CategoryID) error
	// UpdateCategory updates a category name
	UpdateCategory(c Category) error
	// Creates a new category returns expense.CategoryAlreadyExistsErr if category is duplicated.
	AddCategory(c Category) error
	// GetCategory retrieves a category by ID
	GetCategory(id CategoryID) (*Category, error)
}
