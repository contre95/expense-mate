package user

import (
	"github.com/google/uuid"
)

func NewUser(username, pass, alias string) (*User, error) {
	newUser := User{
		ID:       uuid.New(),
		Username: username,
		Password: pass,
		Alias:    alias,
	}
	return &newUser, newUser.validate()
}
