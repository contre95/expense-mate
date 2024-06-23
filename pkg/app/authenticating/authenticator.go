package authenticating

type Service struct {
	UserAuthenticator UserAuthenticator
}

func NewService(u UserAuthenticator) Service {
	return Service{u}
}
