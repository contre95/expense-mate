package managing

// Service just holds all the managing use cases
type Service struct {
	CategoryCreator CategoryCreator
	CategoryDeleter CategoryDeleter
	UserCreator     UserCreator
}

// NewService is the interctor for all Managing Use cases
func NewService(cc CategoryCreator, cd CategoryDeleter, uc UserCreator) Service {
	return Service{cc, cd, uc}
}
