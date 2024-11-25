package expense

import (
	"time"
)

type Installent struct {
	ID          string        `validate:"required"`
	RepeatEvery time.Duration `validate:"required"`
	ExpensesID  []ExpenseID   `validate:"required"`
	CategoryID  CategoryID    `validate:"required"`
	UsersID     []UserID
	StartDate   time.Time `validate:"required"`
	EndDate     time.Time `validate:"required"`
	Amount      float64   `validate:"required"`
	Description string
	Product     string
	Shop        string
}

// Installments is the Installments repository.
type Installments interface {
	All() ([]Installent, error)
	Add(Installent) error
	Delete(id string) error
}

func (r *Installent) IsOver() bool {
	// TODO: Implement this function
	return false
}
