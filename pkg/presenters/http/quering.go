package http

import (
	"expenses-app/pkg/app/querying"
	"fmt"
	"net/http"
	"strconv"
	"time"

	fiber "github.com/gofiber/fiber/v2"
)

type Category struct {
	ID   string
	Name string
}

const DEFAULT_DAYS_FROM_PARAM = "90"
const DEFAULT_DAYS_TO_PARAM = "0" // Now
const DEFAULT_PSIZE_PARAM = "50"
const DEFAULT_PNUM_PARAM = "0"

func getExpenses(eq querying.ExpenseQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		daysFrom, err := strconv.Atoi(c.Query("days_from", DEFAULT_DAYS_FROM_PARAM))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"err":     fmt.Sprintf("%v", err),
				"dat":     nil,
			})
		}
		daysTo, err := strconv.Atoi(c.Query("days_to", DEFAULT_DAYS_TO_PARAM))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"err":     fmt.Sprintf("%v", err),
				"dat":     nil,
			})
		}
		pageNum, err := strconv.Atoi(c.Query("page_num", DEFAULT_PNUM_PARAM))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"err":     fmt.Sprintf("%v", err),
				"dat":     nil,
			})
		}
		pageSize, err := strconv.Atoi(c.Query("page_size", DEFAULT_PSIZE_PARAM))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"err":     fmt.Sprintf("%v", err),
				"dat":     nil,
			})
		}
		req := querying.ExpenseQuerierReq{
			From:     time.Now().Add(-1 * time.Hour * 24 * time.Duration(daysFrom)),
			To:       time.Now().Add(-1 * time.Hour * 24 * time.Duration(daysTo)),
			Page:     uint(pageNum),
			PageSize: uint(pageSize),
		}
		fmt.Println(req)
		resp, err := eq.Query(req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"err":     fmt.Sprintf("%v", err),
				"dat":     nil,
			})
		}
		return c.Status(http.StatusAccepted).JSON(&fiber.Map{
			"success": true,
			"err":     nil,
			"data":    resp,
		})

	}
}

func getCategories(cg querying.CategoryQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		resp, err := cg.Query()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"err":     fmt.Sprintf("%v", err),
				"dat":     nil,
			})
		}
		categories := []Category{}
		for id, name := range resp.Categories {
			categories = append(categories, Category{id, name})
		}
		return c.Status(http.StatusAccepted).JSON(&fiber.Map{
			"success": true,
			"err":     nil,
			"data":    categories,
		})
	}
}
