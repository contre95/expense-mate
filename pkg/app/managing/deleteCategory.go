package managing

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"
)

type DeleteCategoryResp struct {
	DeletedDate time.Time
	Softdelete  bool
}

type DeleteCategoryReq struct {
	ID string
}

type CategoryDeleter struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewCategoryDeleter(l app.Logger, e expense.Expenses) *CategoryDeleter {
	return &CategoryDeleter{l, e}
}

func (s *CategoryDeleter) Delete(req DeleteCategoryReq) (*DeleteCategoryResp, error) {
	panic("Implement me ?")
	//err := s.expenses.DeleteCategory(expense.CategoryID(req.ID))
	//if err != nil {
	//s.logger.Err("Error updating client", err)
	//return nil, errors.New("Could not Delete client information.")
	//}
	//resp := &DeleteCategoryResp{
	//DeletedDate: time.Now(),
	//Softdelete:  false,
	//}
	//s.logger.Info("Category %s deleted", req.ID)
	//return resp, nil
}
