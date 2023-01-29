package tracking

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
)

type CreateExpenseResp struct {
	ID  string
	Msg string
}

type CreateExpenseReq struct {
	Name string
}

// ExpenseCreator use case creates a category for a expense
type ExpenseCreator struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewExpenseCreator(l app.Logger, e expense.Expenses) *ExpenseCreator {
	return &ExpenseCreator{l, e}
}

// Create use cases function creates a new category
func (s *ExpenseCreator) Create(req CreateExpenseReq) (*CreateExpenseResp, error) {
	panic("Implement me ?")
	//category := expense.NewCategory(req.Name)
	//err := u.expenses.SaveCategory(category)
	//if err != nil {
	//u.logger.Err("Could not create category: %v", err)
	//return nil, err
	//}
	//resp := &CreateExpenseResp{ID: string(category.ID), Msg: "Category created"}
	//u.logger.Info("Category %s, created", category.Name)
	//return resp, nil
}
