package managing

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
)

type CreateRuleReq struct {
	Pattern    string
	CategoryID string
}

type DeleteRuleReq struct {
	ID string
}

type RulesBasic struct {
	Pattern    string
	CategoryID string
}

type ListRulesResp struct {
	Rules map[string]RulesBasic
}

type RuleManager struct {
	logger app.Logger
	rules  expense.Rules
}

func NewRuleManager(l app.Logger, r expense.Rules) *RuleManager {
	return &RuleManager{l, r}
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
		rulesMap[string(r.ID)] = RulesBasic{
			Pattern:    r.Pattern,
			CategoryID: string(r.CategoryID),
		}
	}
	return &ListRulesResp{Rules: rulesMap}, nil
}

func (s *RuleManager) Create(req CreateRuleReq) error {
	newRule, createErr := expense.NewRule(req.Pattern, expense.CategoryID(req.CategoryID))
	if createErr != nil {
		return createErr
	}
	err := s.rules.Add(*newRule)
	if err != nil {
		s.logger.Err("Failed to create rule:", err)
		return err
	}
	s.logger.Info("Rule created successfully")
	return nil
}
