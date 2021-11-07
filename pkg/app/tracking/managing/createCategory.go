package managing

import (
	"expenses/pkg/app"
	"expenses/pkg/domain/expense"
)

type CreateResponse struct {
	ID  expense.CategoryID
	Msg string
}

type CreateRequest struct {
	Name string
}

// The createCategory use case creates a category for a expense
type CreateUseCase struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewCategoryCreator(l app.Logger, e expense.Expenses) *CreateUseCase {
	return &CreateUseCase{l, e}
}

// Create use cases function creates a new category
func (u *CreateUseCase) Create(req CreateRequest) (*CreateResponse, error) {
	category := expense.NewCategory(req.Name)
	err := u.expenses.SaveCategory(category)
	if err != nil {
		u.logger.Err("Could not create category: %v", err)
		return nil, err
	}
	resp := &CreateResponse{ID: category.ID, Msg: "Category created"}
	u.logger.Info("Category %s, created", category.Name)
	return resp, nil
}
