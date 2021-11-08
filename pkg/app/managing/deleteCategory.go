package managing

import (
	"errors"
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

type DeleteCategoryUseCase struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewDeleteCategoryUseCase(l app.Logger, e expense.Expenses) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{l, e}
}

func (s *DeleteCategoryUseCase) Delete(req DeleteCategoryReq) (*DeleteCategoryResp, error) {
	err := s.expenses.DeleteCategory(expense.CategoryID(req.ID))
	if err != nil {
		s.logger.Err("Error updating client", err)
		return nil, errors.New("Could not Delete client information.")
	}
	resp := &DeleteCategoryResp{
		DeletedDate: time.Now(),
		Softdelete:  false,
	}
	s.logger.Info("Category %s deleted", req.ID)
	return resp, nil

}
