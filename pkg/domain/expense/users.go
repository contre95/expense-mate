package expense

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type UserID = uuid.UUID

type User struct {
	ID               UserID `json:"id"`
	DisplayName      string `json:"display_name"`
	TelegramUsername string `json:"telegram_username"`
}

type Users interface {
	All() ([]User, error)
	Add(User) error
	Delete(UserID) error
}

func NewUser(name, telegram string) (*User, error) {
	u := User{
		ID:               UserID(uuid.New()),
		DisplayName:      name,
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
