package importing

// Service just holds all the importing use cases
type Service struct {
	ImportExpenses ImportExpenses
}

// NewService returns a new improting.Service
func NewService(ie ImportExpenses) Service {
	return Service{ie}
}
