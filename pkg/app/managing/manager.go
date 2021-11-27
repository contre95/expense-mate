package managing

// Service just holds all the managing use cases
type Service struct {
	UserCreator UserCreator
}

// NewService is the interctor for all Managing Use cases
func NewService(uc UserCreator) Service {
	return Service{uc}
}
