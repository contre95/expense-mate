package ui

import (
	"expenses-app/pkg/app/querying"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const DEFAULT_DAYS_FROM_PARAM = "190"
const DEFAULT_DAYS_TO_PARAM = "0" // Now
const DEFAULT_PSIZE_PARAM = "20"
const DEFAULT_PNUM_PARAM = "0"

func ExpenseEdit(eq querying.ExpenseQuerier, cq querying.CategoryQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func ExpenseRowEdit(eq querying.ExpenseQuerier, cq querying.CategoryQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			// c.Append("hx-trigger", "newPair")  // Not working :(
			return c.Render("main", fiber.Map{
				"ExpensesTrigger": "revealed",
			})
		}
		respCategories, err := cq.Query()
		if err != nil {
			panic("Implement error")
		}
		respExpense, err := eq.GetByID(c.Params("id"))
		if err != nil {
			panic("Implement error")
		}
		fmt.Println(respCategories)
		fmt.Println(respExpense)
		return c.Render("sections/expenses/rowEdit", fiber.Map{
			"Expense":    respExpense.Expenses[0],
			"Categories": respCategories.Categories,
		})

	}
}

// ExpenseSection rendersn the Expenses section
func ExpenseSection(eq querying.ExpenseQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			// c.Append("hx-trigger", "newPair")  // Not working :(
			return c.Render("main", fiber.Map{
				"ExpensesTrigger": "revealed",
			})
		}
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
			From:        time.Now().Add(-1 * time.Hour * 24 * time.Duration(daysFrom)),
			To:          time.Now().Add(-1 * time.Hour * 24 * time.Duration(daysTo)),
			Page:        uint(pageNum),
			MaxPageSize: uint(pageSize),
		}
		resp, err := eq.Query(req)
		if err != nil {
			// panic("Implement error UI")
			log.Panic(err)
		}
		return c.Render("sections/expenses/index", fiber.Map{
			"Expenses":    resp.Expenses,
			"From":        req.From,
			"To":          req.To,
			"CurrentPage": req.Page,      // Add this line
			"NextPage":    resp.Page + 1, // Add this line
			"PrevPage":    resp.Page - 1, // Add this line
			"PageSize":    resp.PageSize,
		})
	}
}
