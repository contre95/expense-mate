package http

import (
	"expenses-app/pkg/app/querying"
	"fmt"
	"net/http"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
)

func getCategories(s querying.CategoryGetter) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		catIDs := strings.Split(c.Query("categories_id"), ",")
		fmt.Println(catIDs)
		var req querying.GetCategoriesReq
		for _, val := range catIDs {
			req.CategoriesIDs[val] = true
		}
		resp, err := s.Get(req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success":    false,
				"err":        fmt.Sprintf("%v", err),
				"categories": "",
			})
		}
		return c.Status(http.StatusAccepted).JSON(&fiber.Map{
			"success":    true,
			"err":        nil,
			"categories": resp.Categories,
		})
	}
}
