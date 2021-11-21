package user

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string
	Password string
	Alias    string
}

// Users is the repository for users
type Users interface {
	Add(u User) error
	Get(uname string) (*User, error)
}
