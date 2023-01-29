package managing

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
)

type CreateCategoryResp struct {
	ID  string
	Msg string
}

type CreateCategoryReq struct {
	Name string
}

// CategoryCreator use case creates a category for an expense
type CategoryCreator struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewCategoryCreator(l app.Logger, e expense.Expenses) *CategoryCreator {
	return &CategoryCreator{l, e}
}

// Create use cases function creates a new category
func (s *CategoryCreator) Create(req CreateCategoryReq) (*CreateCategoryResp, error) {
	panic("Implement me ?")
	//category := expense.NewCategory(req.Name)
	//err := s.expenses.SaveCategory(category)
	//if err != nil {
	//s.logger.Err("Could not create category: %v", err)
	//return nil, err
	//}
	//resp := &CreateCategoryResp{ID: string(category.ID), Msg: "Category created"}
	//s.logger.Info("Category %s, created", category.Name)
	//return resp, nil
}
