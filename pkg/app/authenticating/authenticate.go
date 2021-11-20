package authenticating

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/user"
)

type LoginResp struct {
	Authenticated bool
	UserID        string
}

type LoginReq struct {
	Username string
	Password string
}

type Authenticator struct {
	logger app.Logger
	users  user.Users
}

func NewAuthenticator(l app.Logger, u user.Users) *Authenticator {
	return &Authenticator{l, u}
}

func (auth *Authenticator) Authenticate(req LoginReq) (*LoginResp, error) {
	auth.logger.Info("Attampting to authenticate user %s", req.Username)
	user, err := auth.users.Get(req.Username)
	if err != nil {
		return nil, errors.New("Failed retrieving users data")
	}
	resp := &LoginResp{
		Authenticated: user.Password == req.Password,
		UserID:        user.ID.String(),
	}
	return resp, nil

}
