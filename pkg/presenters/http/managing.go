package http

import (
	"expenses-app/pkg/app/managing"
	"fmt"
	"net/http"

	fiber "github.com/gofiber/fiber/v2"
)

// createCategory handler handles the POST HTTP request for creating an categories
func createCategory(cc managing.CategoryCreator) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bodyJSON := categoryJSON{}
		err := c.BodyParser(&bodyJSON)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"cat_id":  nil,
				"err":     fmt.Sprintf("%v", err),
				"msg":     "Could not parse JSON body",
			})
		}
		req := managing.CreateCategoryReq{
			Name: bodyJSON.Name,
		}
		resp, err := cc.Create(req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"msg":     fmt.Sprintf("Could not create category: %s", req.Name),
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

// createUsers handler handles the POST HTTP request for creating an user
func createUsers(uc managing.UsersCreator) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bodyJSON := usersJSON{}
		err := c.BodyParser(&bodyJSON)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"err":     fmt.Sprintf("%v", err),
				"msg":     "Could not parse JSON body",
			})
		}
		//id := c.Params("id")
		req := managing.CreateUserReq{
			Username: bodyJSON.Name,
			Password: bodyJSON.Pass,
			Alias:    bodyJSON.Name,
		}
		resp, err := uc.Create(req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"msg":     fmt.Sprintf("Could not create user: %s", req.Username),
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
