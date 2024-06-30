package ui

import (
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LoadCategoriesTable(cq querying.CategoryQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		resp, err := cq.Query()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not load categories",
			})
		}
		return c.Render("sections/settings/categories", fiber.Map{
			"Categories": resp.Categories,
		})
	}
}

func CreateCategory(cc managing.CategoryCreator) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		categoryName := c.FormValue("category_name")
		req := managing.CreateCategoryReq{
			Name: categoryName,
		}
		resp, err := cc.Create(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Msg": fmt.Sprintf("Could create category: %v", err),
			})
		}
		c.Append("Hx-Trigger", "reloadCategoriesTable")
		return c.Render("alerts/toastOk", fiber.Map{
			"Msg": fmt.Sprintf("Category %s created.", resp.ID),
		})
	}
}

func EditCategory(cc managing.CategoryUpdater) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Append("Hx-Trigger", "reloadCategoriesTable")
		id := c.Params("id")
		newCategoryName := c.FormValue("category_name")
		req := managing.UpdateCategoryReq{
			ID:      id,
			NewName: newCategoryName,
		}
		resp, err := cc.Update(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Msg": fmt.Sprintf("Could update category: %v", err),
			})
		}
		return c.Render("alerts/toastOk", fiber.Map{
			"Msg": fmt.Sprintf("Category %s updated.", resp.ID),
		})
	}
}

func DeleteCategory(cd managing.CategoryDeleter) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := managing.DeleteCategoryReq{
			ID: c.Params("id"),
		}
		resp, err := cd.Delete(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Msg": err,
			})
		}
		c.Append("Hx-Trigger", "reloadCategoriesTable")
		return c.Render("alerts/toastOk", fiber.Map{
			"Msg": fmt.Sprintf("Category %s deleted.", resp.ID),
		})
	}
}
