package tracking

// Service just holds all the managing use cases
type Service struct {
	ExpenseCreator ExpenseCreator
	ExpenseUpdater ExpenseUpdater
}

// NewService is the interctor for all Managing Use cases
func NewService(ec ExpenseCreator, eu ExpenseUpdater) Service {
	return Service{ec, eu}
}
