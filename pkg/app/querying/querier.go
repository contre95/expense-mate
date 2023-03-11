package querying

// Service just holds all the querying use cases
type Service struct {
	CategoryQuerier CategoryQuerier
	ExpenseQuerier  ExpenseQuerier
}

func NewService(cq CategoryQuerier, eq ExpenseQuerier) Service {
	return Service{cq, eq}
}
