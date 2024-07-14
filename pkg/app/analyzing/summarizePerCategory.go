package analyzing

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"
)

// New code for summarizing expenses per category and per date

type ExpensesSummary struct {
	Category string
	Date     time.Time
	Total    float64
	Count    uint
}

type ExpenseSummaryResp struct {
	Summaries    []ExpensesSummary
	SummaryCount uint
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
	summaryMap := make(map[string]map[time.Time]ExpensesSummary)
	for _, e := range expenses {
		if _, exists := summaryMap[e.Category.ID.String()]; !exists {
			summaryMap[e.Category.ID.String()] = make(map[time.Time]ExpensesSummary)
		}

		dateKey := time.Date(e.Date.Year(), e.Date.Month(), e.Date.Day(), 0, 0, 0, 0, e.Date.Location())
		if summary, exists := summaryMap[e.Category.ID.String()][dateKey]; exists {
			summary.Total += e.Amount
			summary.Count++
			summaryMap[e.Category.ID.String()][dateKey] = summary
		} else {
			summaryMap[e.Category.ID.String()][dateKey] = ExpensesSummary{
				Category: e.Category.Name,
				Date:     dateKey,
				Total:    e.Amount,
				Count:    1,
			}
		}
	}

	var summaries []ExpensesSummary
	for _, dateSummaries := range summaryMap {
		for _, summary := range dateSummaries {
			summaries = append(summaries, summary)
		}
	}
	a.logger.Info("Summarized %d expenses", len(summaries))
	resp := ExpenseSummaryResp{
		Summaries:    summaries,
		SummaryCount: uint(len(summaries)),
	}

	return &resp, nil
}
