package managing

import (
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
type UserCreator struct {
	logger app.Logger
	users  user.Users
	hashed app.Hasher
}

func NewUserCreator(l app.Logger, h app.Hasher, u user.Users) *UserCreator {
	return &UserCreator{l, u, h}
}

func (s *UserCreator) Create(req CreateUserReq) (*CreateUserResp, error) {
	encryptedPass := req.Password
	user, _ := user.NewUser(req.Username, encryptedPass, req.Alias)
	err := s.users.Add(*user)
	if err != nil {
		s.logger.Info("Could not create user %s", req.Username)
		return nil, err
	}
	return nil, nil
}
