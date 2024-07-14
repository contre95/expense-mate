package ui

import (
	"expenses-app/pkg/app/analyzing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoadExpensesTable rendersn the Expenses section
func LoadCategorySummaryTable(ea analyzing.ExpenseAnalyzer) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := analyzing.ExpenseSummaryReq{
			TimeRange: [2]time.Time{},
		}
		resp, err := ea.Summarize(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not load category summary.",
			})
		}
		return c.Render("sections/dashboard/categorySummary", fiber.Map{
			"Summaries":    resp.Summaries,
			"SummaryCount": resp.SummaryCount,
		})
	}
}

