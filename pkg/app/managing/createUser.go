package managing

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/user"
	"fmt"
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
	var err error
	var exists bool
	exists, err = s.users.Exists(req.Username)
	if err != nil {
		s.logger.Err("Could not check if user %s exists: %v", req.Username, err)
		return nil, err
	}
	s.logger.Info("Attemp to create user %s", req.Username)
	if !exists {
		newUser, _ := user.NewUser(req.Username, req.Password, req.Alias)
		s.logger.Info("Creating user %s uuid %s", newUser.Username, newUser.ID)
		err := s.users.Add(*newUser)
		if err != nil {
			s.logger.Err("Could not create user %s: %v", req.Username, err)
			return nil, err
		}
		s.logger.Info("User %s with ID: %s created succesfully.", newUser.Username, newUser.ID.String())
		return &CreateUserResp{
			UUID: newUser.ID.String(),
		}, nil
	} else {
		existingUser, err := s.users.Get(req.Username)
		if err != nil {
			s.logger.Err("Could not retrieve existing user %s: %v", req.Username, err)
			return nil, err
		}
		s.logger.Warn("User %s already exist with uuid %s", existingUser.Username, existingUser.ID)
		return &CreateUserResp{
			UUID: existingUser.ID.String(),
		}, errors.New(fmt.Sprintf("User with name %s already exists with uuid %s", req.Username, existingUser.ID.String()))

	}
}
