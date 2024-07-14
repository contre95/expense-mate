package ui

import (
	"expenses-app/pkg/app/analyzing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoadExpensesTable rendersn the Expenses section
func LoadCategorySummaryTable(ea analyzing.ExpenseAnalyzer) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		now := time.Now().Add(time.Hour * 24)
		startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)
		endOfLastMonth := startOfThisMonth.Add(-time.Nanosecond)
		pastMonth := [2]time.Time{startOfLastMonth, endOfLastMonth}
		thisMonth := [2]time.Time{startOfThisMonth, now}

		// Function to summarize expenses
		// summarizeExpenses := func(timeRange [2]time.Time) ([]analyzing.Summary, error) {
		// 	summaryReq := analyzing.ExpenseSummaryReq{
		// 		TimeRange: timeRange,
		// 	}
		// 	summaryResp, err :=
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	return summaryResp.Summaries, nil
		// }

		pastMonthSummary, err := ea.Summarize(analyzing.ExpenseSummaryReq{
			TimeRange: pastMonth,
		})
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not load category summary.",
			})
		}
		thisMonthSummary, err := ea.Summarize(analyzing.ExpenseSummaryReq{
			TimeRange: thisMonth,
		})
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not load category summary.",
			})
		}

		return c.Render("sections/dashboard/categorySummary", fiber.Map{
			"SummariesThisMonth": thisMonthSummary,
			"SummariesPastMonth": pastMonthSummary,
			"StartOfThisMonth":   startOfThisMonth.Format("2006-01-02"),
			"StartOfLastMonth":   startOfLastMonth.Format("2006-01-02"),
			"EndOfLastMonth":     endOfLastMonth.Format("2006-01-02"),
			"Today":              now.Format("2006-01-02"),
		})
	}
}
