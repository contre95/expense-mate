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
	// Adds an new user to the database
	Add(u User) error
	// Checks if an user exists, returns false on error
	Exists(uname string) (bool, error)
	// Retrieves an user from the database
	Get(uname string) (*User, error)
}
