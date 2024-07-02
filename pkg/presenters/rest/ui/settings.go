package ui

import (
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LoadCategoriesConfig(cq querying.CategoryQuerier) func(*fiber.Ctx) error {
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

func LoadRulesConfig(cq querying.CategoryQuerier, rm managing.RuleManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		categoriesResp, err := cq.Query()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not load categories",
			})
		}
		rulesResp, err := rm.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not load categories",
			})
		}
		return c.Render("sections/settings/rules", fiber.Map{
			"Categories": categoriesResp.Categories,
			"Rules":      rulesResp.Rules,
		})
	}
}

func CreateRule(rm managing.RuleManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Append("Hx-Trigger", "reloadRulesConfig")
		pattern := c.FormValue("rule_pattern")
		categoryID := c.FormValue("category_id")
		req := managing.CreateRuleReq{
			Pattern:    pattern,
			CategoryID: categoryID,
		}
		err := rm.Create(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not create rule.",
			})
		}
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Error",
			"Msg":   "Rule created for " + categoryID,
		})
	}
}

func DeleteRule(rm managing.RuleManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Append("Hx-Trigger", "reloadRulesConfig")
		req := managing.DeleteRuleReq{
			ID: c.Params("id"),
		}
		err := rm.Delete(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not delete rule.",
			})
		}
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Error",
			"Msg":   "Rule deleted",
		})
	}
}

func CreateCategory(cc managing.CategoryCreator) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Append("Hx-Trigger", "reloadCategoriesConfig")
		c.Append("Hx-Trigger", "reloadRulesConfig")
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
		return c.Render("alerts/toastOk", fiber.Map{
			"Msg": fmt.Sprintf("Category %s created.", resp.ID),
		})
	}
}

func LoadTelegramConfig() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("sections/settings/telegram", fiber.Map{})
	}
}

func LoadTelegramStatus(h health.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("sections/settings/telegramStatus", fiber.Map{
			"Status": h.CheckBotHealth(),
		})
	}
}

func SendTelegramCommandOutput(tc managing.TelegramCommander) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cmd := c.FormValue("command")
		fmt.Println("Command to send: ", cmd)
		resp, err := tc.Command(cmd)
		fmt.Println("Response received: ", resp.Msg)
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString(resp.Msg)
	}
}

func SendTelegramCommand(tc managing.TelegramCommander) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cmd := c.FormValue("command")
		fmt.Println("Command to send: ", cmd)
		resp, err := tc.Command(cmd)
		fmt.Println("Response received: ", resp.Msg)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Msg": "Failed to send command",
			})
		}
		return c.Render("alerts/toastOk", fiber.Map{
			"Msg": "Command sent",
		})
	}
}

func EditCategory(cc managing.CategoryUpdater) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Append("Hx-Trigger", "reloadCategoriesConfig")
		c.Append("Hx-Trigger", "reloadRulesConfig")
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
		c.Append("Hx-Trigger", "reloadCategoriesConfig")
		c.Append("Hx-Trigger", "reloadRulesConfig")
		req := managing.DeleteCategoryReq{
			ID: c.Params("id"),
		}
		resp, err := cd.Delete(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Msg": err,
			})
		}
		return c.Render("alerts/toastOk", fiber.Map{
			"Msg": fmt.Sprintf("Category %s deleted.", resp.ID),
		})
	}
}
