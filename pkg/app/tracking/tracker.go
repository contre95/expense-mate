package tracking

// Service just holds all the managing use cases
type Service struct {
	ExpenseCreator   ExpenseCreator
	ExpenseUpdater   ExpenseUpdater
	ExpenseDeleter   ExpenseDeleter
	ExpenseCataloger ExpenseCataloger
}

// NewService is the interctor for all Managing Use cases
func NewService(ec ExpenseCreator, eu ExpenseUpdater, ed ExpenseDeleter, eca ExpenseCataloger) Service {
	return Service{ec, eu, ed, eca}
}
