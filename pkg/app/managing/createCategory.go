package managing

import (
	"expenses/pkg/app"
	"expenses/pkg/domain/expense"
)

type CreateCategoryResp struct {
	ID  string
	Msg string
}

type CreateCategoryReq struct {
	Name string
}

// The createCategory use case creates a category for a expense
type CreateCategoryUseCase struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewCreateCategoryUseCase(l app.Logger, e expense.Expenses) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{l, e}
}

// Create use cases function creates a new category
func (u *CreateCategoryUseCase) Create(req CreateCategoryReq) (*CreateCategoryResp, error) {
	category := expense.NewCategory(req.Name)
	err := u.expenses.SaveCategory(category)
	if err != nil {
		u.logger.Err("Could not create category: %v", err)
		return nil, err
	}
	resp := &CreateCategoryResp{ID: string(category.ID), Msg: "Category created"}
	u.logger.Info("Category %s, created", category.Name)
	return resp, nil
}
