package tracking

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
)

type ApplyRuleResp struct {
	CategoryID string
	UsersID    []string
	RuleID     string
	Matched    bool
}

type ApplyRuleReq struct {
	Product string
	Shop    string
}

// RuleApplier use case creates a category for a expense
type RuleApplier struct {
	logger app.Logger
	rules  expense.Rules
}

func NewRuleApplier(l app.Logger, r expense.Rules) *RuleApplier {
	return &RuleApplier{l, r}
}

func (s *RuleApplier) Apply(req ApplyRuleReq) *ApplyRuleResp {
	rules, err := s.rules.All()
	if err != nil {
		return nil
	}
	for _, rule := range rules {
		// I'm not using req.Product for now
		s.logger.Debug("Matching %s against %s", req.Shop, rule.Pattern)
		if rule.Matches(req.Shop) {
			s.logger.Info("Shop:%s matched against Pattern: %s", req.Shop, rule.Pattern)
			uids := []string{}
			for _, ruid := range rule.UsersID {
				uids = append(uids, ruid.String())
			}
			return &ApplyRuleResp{
				CategoryID: rule.CategoryID.String(),
				UsersID:    uids,
				RuleID:     rule.ID,
				Matched:    true,
			}
		}
	}
	return &ApplyRuleResp{Matched: false}
}
