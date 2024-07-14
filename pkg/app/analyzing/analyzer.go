package analyzing

type Service struct {
	ExpenseAnalyzer ExpenseAnalyzer
}

func NewService(ea ExpenseAnalyzer) Service {
	return Service{ea}
}
