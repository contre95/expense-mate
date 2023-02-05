package managing

// Service just holds all the managing use cases
type Service struct {
	UserCreator     UsersCreator
	CategoryCreator CategoryCreator
}

// NewService is the interctor for all Managing Use cases
func NewService(uc UsersCreator, cc CategoryCreator) Service {
	return Service{uc, cc}
}
