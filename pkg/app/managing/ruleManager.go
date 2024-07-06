package managing

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"

	"github.com/google/uuid"
)

type CreateRuleReq struct {
	Pattern    string
	CategoryID string
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
}

type ListRulesResp struct {
	Rules map[string]RulesBasic
}

type RuleManager struct {
	logger   app.Logger
	rules    expense.Rules
	expenses expense.Expenses
}

func NewRuleManager(l app.Logger, r expense.Rules, e expense.Expenses) *RuleManager {
	return &RuleManager{l, r, e}
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
		s.logger.Err("Failed to list rules:", err)
		return nil, err
	}
	s.logger.Info("Rules listed successfully")
	rulesMap := map[string]RulesBasic{}
	for _, r := range rules {
		category, err := s.expenses.GetCategory(r.CategoryID)
		if err != nil {
			s.logger.Err("Failed to retrievecategory:", err)
			return nil, err
		}
		rulesMap[string(r.ID)] = RulesBasic{
			Pattern: r.Pattern,
			Category: struct {
				ID   string
				Name string
			}{category.ID.String(), category.Name},
		}
	}
	return &ListRulesResp{Rules: rulesMap}, nil
}

func (s *RuleManager) Create(req CreateRuleReq) error {
	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return expense.ErrInvalidID
	}
	newRule, createErr := expense.NewRule(req.Pattern, catID)
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
