package expense

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type UserID uuid.UUID

type User struct {
	ID               UserID
	Name             string
	TelegramUsername string
}

type Users interface {
	All() ([]User, error)
}

func NewUser(name, telegram string) (*User, error) {
	u := User{
		ID:               UserID(uuid.New()),
		Name:             name,
		TelegramUsername: telegram,
	}
	return u.Validate()
}

func (u *User) Validate() (*User, error) {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}
		return nil, errors.New(fmt.Sprintf("Invalid user data: %v", err))
	}
	return u, nil
}
