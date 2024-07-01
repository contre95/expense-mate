package expense

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/go-playground/validator"
)

func (c *Expense) Validate() (*Expense, error) {
	validate := validator.New()
	validate.RegisterValidation("alphanum_space", spaceOrLetter)
	err := validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}
		fmt.Println(err)
		return nil, errors.New(fmt.Sprintf("Invalid expense data. Please review the fields."))
	}
	return c, nil
}

func (c *Category) Validate() (*Category, error) {
	validate := validator.New()
	validate.RegisterValidation("alphanum_space", spaceOrLetter)
	err := validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}
		fmt.Println(err)
		return nil, ErrInvalidEntity
	}
	return c, nil
}

// spaceOrLetter checks that a string contains only ASCII characters and spaces.
func spaceOrLetter(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		// only letters, spaces, numbers and this: &
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && !unicode.IsSpace(char) && char != 38 {
			return false
		}
	}
	return true
}
