package managing

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"

	"github.com/google/uuid"
)

type CreateRuleReq struct {
	Pattern    string
	CategoryID string
	UsersID    []string
}

type DeleteRuleReq struct {
	ID string
}

type RulesBasic struct {
	Pattern  string
	Category struct {
		ID   string
		Name string
	}
	Users []struct {
		ID          string
		DisplayName string
	}
}

type ListRulesResp struct {
	Rules map[string]RulesBasic
}

type RuleManager struct {
	logger   app.Logger
	rules    expense.Rules
	expenses expense.Expenses
	users    expense.Users
}

func NewRuleManager(l app.Logger, r expense.Rules, e expense.Expenses, u expense.Users) *RuleManager {
	return &RuleManager{l, r, e, u}
}

func (s *RuleManager) Delete(req DeleteRuleReq) error {
	err := s.rules.Delete(req.ID)
	if err != nil {
		s.logger.Err("Failed to delete rule:", err)
		return err
	}
	s.logger.Info("Rule deleted successfully")
	return nil
}

func (s *RuleManager) List() (*ListRulesResp, error) {
	rules, err := s.rules.All()
	if err != nil {
		s.logger.Err("Failed to retrieve rules from storage:", err)
		return nil, err
	}
	s.logger.Info("Rules listed successfully")
	allUsers, err := s.users.All()
	if err != nil {
		s.logger.Err("Failed to get users:", err)
		return nil, err
	}
	rulesMap := map[string]RulesBasic{}
	for _, r := range rules {
		category, err := s.expenses.GetCategory(r.CategoryID)
		if err != nil {
			s.logger.Err("Failed to retrieve category:", err)
			return nil, err
		}
		rulesBasic := RulesBasic{
			Pattern: r.Pattern,
			Category: struct {
				ID   string
				Name string
			}{category.ID.String(), category.Name},
			Users: []struct {
				ID          string
				DisplayName string
			}{},
		}
		for _, au := range allUsers {
			for _, ruid := range r.UsersID {
				if au.ID == ruid {
					rulesBasic.Users = append(rulesBasic.Users, struct {
						ID          string
						DisplayName string
					}{au.ID.String(), au.DisplayName})
				}
			}
		}
		rulesMap[string(r.ID)] = rulesBasic
	}
	return &ListRulesResp{Rules: rulesMap}, nil
}

func (s *RuleManager) Create(req CreateRuleReq) error {
	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return expense.ErrInvalidID
	}
	usersID := []uuid.UUID{}
	for _, sid := range req.UsersID {
		pid, err := uuid.Parse(sid)
		if err != nil {
			s.logger.Err("Failed to parse UUID %s", err.Error())
			return errors.New("Failed to parse UUID %s" + sid)
		}
		usersID = append(usersID, pid)
	}
	newRule, createErr := expense.NewRule(req.Pattern, usersID, catID)
	if createErr != nil {
		return createErr
	}
	err = s.rules.Add(*newRule)
	if err != nil {
		s.logger.Err("Failed to create rule:", err)
		return err
	}
	s.logger.Info("Rule created successfully")
	return nil
}
