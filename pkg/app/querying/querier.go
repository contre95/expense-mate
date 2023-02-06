package querying

// Service just holds all the querying use cases
type Service struct {
	CategoryGetter CategoryQuerier
}

func NewService(cg CategoryQuerier) Service {
	return Service{cg}
}
