package http

import (
	"expenses-app/pkg/app/managing"
	"fmt"
	"log"
	"net/http"

	fiber "github.com/gofiber/fiber/v2"
)

func createCategories(s managing.CreateCategory) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		catJSON := addCategoriesJSON{}
		err := c.BodyParser(&catJSON)
		if err != nil || len(catJSON.Names) == 0 {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"cat_id":  nil,
				"err":     fmt.Sprintf("%v", err),
				"msg":     "",
			})
		}
		failed_categories := map[string]string{}
		for _, name := range catJSON.Names {
			req := &managing.CreateCategoryReq{
				Name: name,
			}
			_, err := s.Create(*req)
			if err != nil {
				failed_categories[req.Name] = err.Error()
			}
		}
		if len(failed_categories) > 0 {
			return c.Status(http.StatusMultiStatus).JSON(&fiber.Map{
				"success": false,
				"msg":     "Some categories could not be created",
				"err":     failed_categories,
			})
		}
		return c.Status(http.StatusAccepted).JSON(&fiber.Map{
			"success": true,
			"err":     "",
			"msg":     fmt.Sprintf("All the categories were created"),
		})
	}
}

func createCategory(s managing.CreateCategory) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		catJSON := addCategoryJSON{}
		err := c.BodyParser(&catJSON)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"cat_id":  nil,
				"err":     fmt.Sprintf("%v", err),
				"msg":     "",
			})
		}
		log.Println(catJSON)
		req := &managing.CreateCategoryReq{
			Name: catJSON.Name,
		}
		resp, err := s.Create(*req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"cat_id":  nil,
				"err":     fmt.Sprintf("%v", err),
				"msg":     "",
			})
		}
		return c.Status(http.StatusAccepted).JSON(&fiber.Map{
			"success": true,
			"cat_id":  resp.ID,
			"err":     nil,
			"msg":     fmt.Sprintf("Category %s, was created with id %s", catJSON.Name, resp.ID),
		})
	}
}

func deleteCategory(s managing.DeleteCategory) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := &managing.DeleteCategoryReq{ID: id}
		resp, err := s.Delete(*req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"msg":     fmt.Sprintf("Could not delete category: %s", req.ID),
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

//func listClients(s clients.List) func(*fiber.Ctx) error {
//return func(c *fiber.Ctx) error {
//page, err1 := strconv.Atoi(c.Query("page", "0"))
//pageSize, err2 := strconv.Atoi(c.Query("pageSize", "10"))
//if err1 != nil || err2 != nil {
//return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
//"success": false,
//"clients": nil,
//"err":     fmt.Sprintf("Wrong page input"),
//})
//}
//listRequest := &clients.ListRequest{
//Page:     uint(page),
//PageSize: uint(pageSize),
//}
//resp, err := s.List(*listRequest)
//if err != nil {
//return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
//"success": false,
//"clients": nil,
//"err":     fmt.Sprintf("%v", err),
//})
//}
//return c.Status(http.StatusAccepted).JSON(&fiber.Map{
//"success": true,
//"clients": resp,
//"err":     nil,
//})
//}
//}

//func sampleHanlder(s clients.Service) func(*fiber.Ctx) error {
//return func(c *fiber.Ctx) error {
//return nil
//}
//}
