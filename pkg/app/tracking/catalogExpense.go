package tracking

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
)

type CatalogExpenseResp struct {
	CategoryID string
	RuleID     string
	Matched    bool
}

type CatalogExpenseReq struct {
	Product string
	Shop    string
}

// ExpenseCataloger use case creates a category for a expense
type ExpenseCataloger struct {
	logger app.Logger
	rules  expense.Rules
}

func NewExpenseCataloger(l app.Logger, r expense.Rules) *ExpenseCataloger {
	return &ExpenseCataloger{l, r}
}

// Create use cases function creates a new expense
func (s *ExpenseCataloger) Catalog(req CatalogExpenseReq) *CatalogExpenseResp {
	rules, err := s.rules.All()
	if err != nil {
		return nil
	}
	for _, rule := range rules {
		// I'm not using req.Product for now
		s.logger.Debug("Matching %s against %s", req.Shop, rule.Pattern)
		if rule.Matches(req.Shop) {
			s.logger.Info("Shop:%s matched against Pattern: %s", req.Shop, rule.Pattern)
			return &CatalogExpenseResp{
				CategoryID: rule.CategoryID.String(),
				RuleID:     rule.ID,
				Matched:    true,
			}
		}
	}
	return &CatalogExpenseResp{Matched: false}
}
