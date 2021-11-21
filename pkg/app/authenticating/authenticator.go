package authenticating

type Service struct {
	UserAuthenticator UserAuthenticator
}

func NewAuthenticator(u UserAuthenticator) Service {
	return Service{u}
}
