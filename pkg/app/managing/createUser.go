package managing

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/user"
)

type CreateUserResp struct {
	UUID string
}

type CreateUserReq struct {
	Username string
	Password string
	Alias    string
}

// CreateUser is the use case in charge of creating new users
type UsersCreator struct {
	logger app.Logger
	users  user.Users
	hashed app.Hasher
}

func NewUserCreator(l app.Logger, h app.Hasher, u user.Users) *UsersCreator {
	return &UsersCreator{l, u, h}
}

// Create tries to create an new user, if it exist it returs an error with the UUID of the existing user in the response
func (s *UsersCreator) Create(req CreateUserReq) (*CreateUserResp, error) {
	s.logger.Info("Attempting to create user %s", req.Username)
	newUser, createErr := user.NewUser(req.Username, req.Password, req.Alias)
	if createErr != nil {
		return nil, createErr
	}
	_, err := s.users.Get(req.Username)
	if err != nil && err.Error() == user.UserNotFoundErr {
		if s.users.Add(*newUser) != nil {
			s.logger.Err("Could not create user %s: %v", req.Username, err)
			return nil, err
		}
		s.logger.Info("User %s with ID: %s created succesfully.", newUser.Username, newUser.ID.String())
		return &CreateUserResp{UUID: newUser.ID.String()}, nil
	}
	if err != nil && err.Error() != user.UserNotFoundErr {
		s.logger.Err("Coldn't validate if user %s exists: %v", req.Username, err)
		return nil, err
	}
	// err is nil
	s.logger.Err("Could not create user %s, already exists", req.Username)
	return nil, errors.New(user.UserAlreadyExists)
}
