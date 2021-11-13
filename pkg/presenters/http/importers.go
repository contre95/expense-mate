package http

import (
	"expenses-app/pkg/app/importing"
	"fmt"
	"log"
	"net/http"

	fiber "github.com/gofiber/fiber/v2"
)

func importExpenses(i importing.ImportExpenses) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bodyJSON := expenseImporterJSON{}
		log.Println(bodyJSON)
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
			BypassWrongExpenses: bodyJSON.BypassWrongExpenses,
			ImporterID:          id,
		}
		log.Println(req)
		log.Println(bodyJSON)
		resp, err := i.Import(req)
		log.Println("Ahora si")
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"msg":     fmt.Sprintf("Could import expense from importer: %s", req.ImporterID),
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
