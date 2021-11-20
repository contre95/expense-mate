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
}

// CreateUser is the use case in charge of creating new users
type UserCreator struct {
	logger app.Logger
	users  user.Users
}

func NewUserCreator(l app.Logger, u user.Users) *UserCreator {
	return &UserCreator{l, u}
}

func (s *UserCreator) Create(req CreateUserReq) (*CreateUserResp, error) {
	user, _ := user.NewUser(req.Username, req.Password)
	s.users.Add(*user)
	return nil, nil
}
