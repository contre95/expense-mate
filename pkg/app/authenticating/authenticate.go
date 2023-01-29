package authenticating

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/user"
)

type LoginResp struct {
	UserID string
	Alias  string
}

type LoginReq struct {
	Username string
	Password string
}

type UserAuthenticator struct {
	logger app.Logger
	hasher app.Hasher
	users  user.Users
}

func NewUserAuthenticator(l app.Logger, h app.Hasher, u user.Users) *UserAuthenticator {
	return &UserAuthenticator{l, h, u}
}

func (a *UserAuthenticator) Authenticate(req LoginReq) (*LoginResp, error) {
	a.logger.Info("Attampting to authenticate user %s", req.Username)
	u, err := a.users.Get(req.Username)
	if err != nil {
		if err.Error() == user.UserNotFoundErr {
			a.logger.Warn("There was an attempt to authenticate with unexistent user %s", req.Username)
		}
		return nil, errors.New("Error trying to retrieve user.")
	}
	if u.Password == req.Password {
		a.logger.Info("User %s has been authenticated", req.Username)
		return &LoginResp{
			UserID: u.ID.String(),
			Alias:  u.Alias,
		}, nil
	}
	a.logger.Warn("User %s not authorized to login.", req.Username)
	return nil, errors.New("Authentication error")
}
