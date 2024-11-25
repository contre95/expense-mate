package managing

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"

	"github.com/google/uuid"
)

type CreateInstallmentReq struct {
	RepeatEvery time.Duration
	CategoryID  string
	ExpensesID  []string
	UsersID     []string
	Amount      float64
	Description string
	Start       time.Time
	End         time.Time
}

type DeleteInstallmentReq struct {
	ID string
}

type InstallmentBasic struct {
	RepeatEvery time.Duration
	Amount      float64
	Description string
	End         time.Time
	Start       time.Time
	Category    struct {
		ID   string
		Name string
	}
	Expenses []struct {
		ID string
	}
	Users []struct {
		ID          string
		DisplayName string
	}
}

type ListInstallmentsResp struct {
	Installments map[string]InstallmentBasic
}

type InstallmentManager struct {
	logger       app.Logger
	installments expense.Installments
	expenses     expense.Expenses
	users        expense.Users
}

func NewInstallmentManager(l app.Logger, i expense.Installments, e expense.Expenses, u expense.Users) *InstallmentManager {
	return &InstallmentManager{l, i, e, u}
}

func (s *InstallmentManager) Delete(req DeleteInstallmentReq) error {
	err := s.installments.Delete(req.ID)
	if err != nil {
		s.logger.Err("Failed to delete installment:", err)
		return err
	}
	s.logger.Info("Installment deleted successfully")
	return nil
}

func (s *InstallmentManager) List() (*ListInstallmentsResp, error) {
	installments, err := s.installments.All()
	if err != nil {
		s.logger.Err("Failed to retrieve installments from storage:", err)
		return nil, err
	}
	s.logger.Info("Installments listed successfully")
	allUsers, err := s.users.All()
	if err != nil {
		s.logger.Err("Failed to get users:", err)
		return nil, err
	}
	installmentsMap := map[string]InstallmentBasic{}
	for _, i := range installments {
		category, err := s.expenses.GetCategory(i.CategoryID)
		if err != nil {
			s.logger.Err("Failed to retrieve category:", err)
			return nil, err
		}
		installmentBasic := InstallmentBasic{
			RepeatEvery: i.RepeatEvery,
			Category: struct {
				ID   string
				Name string
			}{category.ID.String(), category.Name},
			Expenses: []struct {
				ID string
			}{},
			Users: []struct {
				ID          string
				DisplayName string
			}{},
		}
		for _, expID := range i.ExpensesID {
			expense, err := s.expenses.Get(expID)
			if err != nil {
				s.logger.Err("Failed to retrieve expense:", err)
				return nil, err
			}
			installmentBasic.Expenses = append(installmentBasic.Expenses, struct {
				ID string
			}{expense.ID.String()})
		}
		for _, au := range allUsers {
			for _, iuid := range i.UsersID {
				if au.ID == iuid {
					installmentBasic.Users = append(installmentBasic.Users, struct {
						ID          string
						DisplayName string
					}{au.ID.String(), au.DisplayName})
				}
			}
		}
		installmentsMap[string(i.ID)] = installmentBasic
	}
	return &ListInstallmentsResp{Installments: installmentsMap}, nil
}

func (s *InstallmentManager) Create(req CreateInstallmentReq) error {
	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return expense.ErrInvalidID
	}
	expensesID := []uuid.UUID{}
	for _, eid := range req.ExpensesID {
		expID, err := uuid.Parse(eid)
		if err != nil {
			s.logger.Err("Failed to parse UUID %s", err.Error())
			return errors.New("Failed to parse UUID %s" + eid)
		}
		expensesID = append(expensesID, expID)
	}
	usersID := []uuid.UUID{}
	for _, uid := range req.UsersID {
		userID, err := uuid.Parse(uid)
		if err != nil {
			s.logger.Err("Failed to parse UUID %s", err.Error())
			return errors.New("Failed to parse UUID %s" + uid)
		}
		usersID = append(usersID, userID)
	}
	newInstallment := expense.Installent{
		ID:          uuid.New().String(),
		RepeatEvery: req.RepeatEvery,
		ExpensesID:  expensesID,
		CategoryID:  catID,
		UsersID:     usersID,
	}
	err = s.installments.Add(newInstallment)
	if err != nil {
		s.logger.Err("Failed to create installment:", err)
		return err
	}
	s.logger.Info("Installment created successfully")
	return nil
}
