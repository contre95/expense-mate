package user

import "github.com/google/uuid"

func NewUser(username, pass string) (*User, error) {
	newUser := User{
		ID:       uuid.New(),
		Username: username,
		Password: pass,
	}
	return &newUser, nil
}
