package managing

type Service struct {
	CreateCategoryUseCase CreateCategoryUseCase
	DeleteCategoryUseCase DeleteCategoryUseCase
}

func NewService(c CreateCategoryUseCase, d DeleteCategoryUseCase) Service {
	return Service{c, d}
}
