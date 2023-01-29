package user

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string `validate:"min=3,max=64"`
	Password string `validate:"min=8,max=64"`
	Alias    string `validate:"max=64"`
}

const UserNotFoundErr = "User not found"
const UserAlreadyExists = "User already exist"

// Users is the repository for users
type Users interface {
	// Adds an new user to the database
	Add(u User) error
	// Retrieves an user from the database, returns user.UserNotFoundErr if user doesn't exist
	Get(uname string) (*User, error)
}

func (u *User) validate() error {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}
		return errors.New(fmt.Sprintf("Invalid category data: %v", err))
	}
	return nil
}
