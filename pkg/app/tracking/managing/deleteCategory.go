package categories

import (
	"errors"
	"expenses/pkg/app"
	"expenses/pkg/domain/expense"
	"time"
)

type DeleteResponse struct {
	DeletedDate time.Time
	Softdelete  bool
}

type DeleteRequest struct {
	ID expense.CategoryID
}

type DeleteUseCase struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewDeleteUseCase(l app.Logger, e expense.Expenses) *DeleteUseCase {
	return &DeleteUseCase{l, e}
}

func (s *DeleteUseCase) Delete(req DeleteRequest) (*DeleteResponse, error) {
	err := s.expenses.DeleteCategory(req.ID)
	if err != nil {
		s.logger.Err("Error updating client", err)
		return nil, errors.New("Could not Delete client information.")
	}
	resp := &DeleteResponse{
		DeletedDate: time.Now(),
		Softdelete:  false,
	}
	s.logger.Info("Category %s deleted", req.ID)
	return resp, nil

}
