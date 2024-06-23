package ui

import (
	"expenses-app/pkg/app/querying"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const DEFAULT_DAYS_FROM_PARAM = "30"
const DEFAULT_DAYS_TO_PARAM = "0" // Now
const DEFAULT_PSIZE_PARAM = "50"
const DEFAULT_PNUM_PARAM = "0"

// Home hanlder reders the homescreen
func ExpensesTable(eq querying.ExpenseQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		daysFrom, err := strconv.Atoi(c.Query("days_from", DEFAULT_DAYS_FROM_PARAM))
		if err != nil {
			// return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			// 	"success": false,
			// 	"err":     fmt.Sprintf("%v", err),
			// 	"dat":     nil,
			// })
		}
		daysTo, err := strconv.Atoi(c.Query("days_to", DEFAULT_DAYS_TO_PARAM))
		if err != nil {
			// return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			// 	"success": false,
			// 	"err":     fmt.Sprintf("%v", err),
			// 	"dat":     nil,
			// })
		}
		pageNum, err := strconv.Atoi(c.Query("page_num", DEFAULT_PNUM_PARAM))
		if err != nil {
			// return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			// 	"success": false,
			// 	"err":     fmt.Sprintf("%v", err),
			// 	"dat":     nil,
			// })
		}
		pageSize, err := strconv.Atoi(c.Query("page_size", DEFAULT_PSIZE_PARAM))
		if err != nil {
			// return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			// 	"success": false,
			// 	"err":     fmt.Sprintf("%v", err),
			// 	"dat":     nil,
			// })
		}
		req := querying.ExpenseQuerierReq{
			From:     time.Now().Add(-1 * time.Hour * 24 * time.Duration(daysFrom)),
			To:       time.Now().Add(-1 * time.Hour * 24 * time.Duration(daysTo)),
			Page:     uint(pageNum),
			PageSize: uint(pageSize),
		}
		resp, err := eq.Query(req)
		if err != nil {
			panic("Implement error UI")
		}
		return c.Render("expensesTable", fiber.Map{
			"Expenses": resp.Expenses,
		})
	}
}
