package managing

// Service just holds all the managing use cases
type Service struct {
	CreateCategory CreateCategory
	DeleteCategory DeleteCategory
}

// NewService returns a new manging.Service
func NewService(c CreateCategory, d DeleteCategory) Service {
	return Service{c, d}
}
