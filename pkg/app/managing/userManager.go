package managing

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreateUserReq struct {
	DisplayName      string
	TelegramUsername string
}

type DeleteUserReq struct {
	ID string
}

type ListUsersResp struct {
	Users map[string]struct {
		DisplayName      string
		TelegramUsername string
	}
}

type UserManager struct {
	logger   app.Logger
	users    expense.Users
	expenses expense.Expenses
}

func NewUserManager(l app.Logger, u expense.Users, e expense.Expenses) *UserManager {
	return &UserManager{l, u, e}
}

func (s *UserManager) Delete(req DeleteUserReq) error {
	// s.expenses.CountWithFilter()
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		return expense.ErrInvalidID
	}
	i, err := s.expenses.CountWithFilter([]string{req.ID}, []string{}, 0, 0, "", "", time.Time{}, time.Time{})
	if err != nil {
		s.logger.Err("Failed count expenses:", err)
		return err
	}
	if i > 0 {
		s.logger.Debug("Failed to delete user: Associated expenses %d", i)
		return errors.New(fmt.Sprintf("Could not delete user, %d expenses are still associated with it", i))
	}
	err = s.users.Delete(expense.UserID(userID))
	if err != nil {
		s.logger.Err("Failed to delete user:", err)
		return err
	}
	s.logger.Info("User deleted successfully")
	return nil
}

func (s *UserManager) List() (*ListUsersResp, error) {
	users, err := s.users.All()
	if err != nil {
		s.logger.Err("Failed to list users:", err)
		return nil, err
	}
	s.logger.Info("Users listed successfully")
	resp := ListUsersResp{
		Users: map[string]struct {
			DisplayName      string
			TelegramUsername string
		}{},
	}
	for _, u := range users {
		resp.Users[u.ID.String()] = struct {
			DisplayName      string
			TelegramUsername string
		}{
			DisplayName:      u.DisplayName,
			TelegramUsername: u.TelegramUsername,
		}
	}
	return &resp, nil
}

func (s *UserManager) Create(req CreateUserReq) error {
	newUser, err := expense.NewUser(req.DisplayName, req.TelegramUsername)
	if err != nil {
		s.logger.Err("Failed to create user:", err)
		return err
	}
	err = s.users.Add(*newUser)
	if err != nil {
		s.logger.Err("Failed to create user:", err)
		return err
	}
	s.logger.Info("User created successfully")
	return nil
}
