package analyzing

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"
)

// New code for summarizing expenses per category and per date

type Summary struct {
	Category string
	Month    time.Month
	Total    float64
	Count    uint
}

type ExpenseSummaryResp struct {
	Summaries map[string]Summary
	Month     time.Month
	Total     float64
}

type ExpenseSummaryReq struct {
	TimeRange [2]time.Time
}

type ExpenseAnalyzer struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewSummarizer(l app.Logger, e expense.Expenses) *ExpenseAnalyzer {
	return &ExpenseAnalyzer{l, e}
}

func (a *ExpenseAnalyzer) Summarize(req ExpenseSummaryReq) (*ExpenseSummaryResp, error) {
	var err error
	expenses, err := a.expenses.Filter(nil, []string{}, 0, 0, "", "", req.TimeRange[0], req.TimeRange[1], 0, 0)
	if err != nil {
		a.logger.Err("Could not get expenses from storage: %v", err)
		return nil, err
	}
	summaries := make(map[string]Summary)
	total := 0.0
	for _, e := range expenses {
		total += e.Amount
		if _, exists := summaries[e.Category.ID.String()]; !exists {
			summaries[e.Category.ID.String()] = Summary{
				Category: e.Category.Name,
				Month:    e.Date.Month(),
				Total:    e.Amount,
				Count:    1,
			}
		}
		summary := summaries[e.Category.ID.String()]
		summary.Count += 1
		summary.Total += e.Amount
		summaries[e.Category.ID.String()] = summary

	}
	resp := ExpenseSummaryResp{
		Summaries: summaries,
		Month:     req.TimeRange[0].Month(),
		Total:     total,
	}
	return &resp, nil
}
