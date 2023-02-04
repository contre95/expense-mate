package tracking

// Service just holds all the managing use cases
type Service struct {
	ExpenseCreator ExpenseCreator
}

// NewService is the interctor for all Managing Use cases
func NewService(ec ExpenseCreator) Service {
	return Service{ec}
}
