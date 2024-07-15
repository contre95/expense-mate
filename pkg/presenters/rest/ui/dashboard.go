package ui

import (
	"expenses-app/pkg/app/analyzing"
	"expenses-app/pkg/app/querying"
	"fmt"
	"slices"
	"strconv"
	"strings"
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

func LoadExpensesMiniTable(eq querying.ExpenseQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		pageNum, err := strconv.Atoi(c.Query("page_num", DEFAULT_PNUM_PARAM))
		if err != nil {
			panic("Atoi parse error")
		}
		pageSize, err := strconv.Atoi(c.Query("page_size", "10000"))
		if err != nil {
			panic("Atoi parse error")
		}
		fromDate, err := time.Parse("2006-01-02", c.Query("from-date", time.Time{}.Format("2006-01-02")))
		if err != nil {
			panic("Date parse error")
		}
		toDate, err := time.Parse("2006-01-02", c.Query("to-date", time.Time{}.Format("2006-01-02")))
		if err != nil {
			panic("Date parse error")
		}
		categories := slices.DeleteFunc(strings.Split(c.Query("categories"), ","), func(s string) bool { return s == "" })
		req := querying.ExpenseQuerierReq{
			Page:        uint(pageNum),
			MaxPageSize: uint(pageSize),
			ExpenseFilter: querying.ExpenseQuerierFilter{
				ByCategoryID: categories,
				ByTime:       [2]time.Time{fromDate, toDate},
			},
		}
		fmt.Println(req.ExpenseFilter.ByUsers)
		resp, err := eq.Query(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Table error",
				"Msg":   "Error loading expenses table.",
			})
		}
		return c.Render("sections/dashboard/miniTableView", fiber.Map{
			"Expenses":      resp.Expenses,
			"ExpensesCount": resp.ExpensesCount,
		})
	}
}
