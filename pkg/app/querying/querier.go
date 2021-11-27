package querying

// Service just holds all the querying use cases
type Service struct {
	CategoryGetter CategoryGetter
}

func NewService(cg CategoryGetter) Service {
	return Service{cg}
}
