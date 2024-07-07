package tracking

// Service just holds all the managing use cases
type Service struct {
	ExpenseCreator ExpenseCreator
	ExpenseUpdater ExpenseUpdater
	ExpenseDeleter ExpenseDeleter
	RuleApplier    RuleApplier
}

// NewService is the interctor for all Managing Use cases
func NewService(ec ExpenseCreator, eu ExpenseUpdater, ed ExpenseDeleter, ra RuleApplier) Service {
	return Service{ec, eu, ed, ra}
}
