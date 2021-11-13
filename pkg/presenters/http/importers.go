package http

import (
	"expenses-app/pkg/app/importing"
	"fmt"
	"net/http"

	fiber "github.com/gofiber/fiber/v2"
)

func importExpenses(i importing.ImportExpenses) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bodyJSON := expenseImporterJSON{}
		err := c.BodyParser(&bodyJSON)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"cat_id":  nil,
				"err":     fmt.Sprintf("%v", err),
				"msg":     "Could not parse JSON body",
			})
		}
		id := c.Params("id")
		req := importing.ImportExpensesReq{
			ByPassWrongExpenses: false,
			ImporterID:          id,
		}
		resp, err := i.Import(req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"msg":     fmt.Sprintf("Could not delete category: %s", req.ImporterID),
				"err":     fmt.Sprintf("%v", err),
			})
		}
		return c.Status(http.StatusAccepted).JSON(&fiber.Map{
			"success": true,
			"msg":     resp,
			"err":     nil,
		})
	}
}
