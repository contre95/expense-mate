package managing

// Service just holds all the managing use cases
type Service struct {
	UserCreator UsersCreator
}

// NewService is the interctor for all Managing Use cases
func NewService(uc UsersCreator) Service {
	return Service{uc}
}
