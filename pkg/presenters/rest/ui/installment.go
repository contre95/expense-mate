package ui

import (
	"expenses-app/pkg/app/managing"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func LoadInsatllmentsTable(rm managing.InstallmentManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		resp, err := rm.List()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(resp)
		return nil
	}
}

func CreateInstallment(im managing.InstallmentManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Append("Hx-Trigger", "reloadInstallmentsConfig")
		amount, err := strconv.ParseFloat(c.FormValue("amount"), 64)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Invalid amount provided.",
			})
		}
		description := c.FormValue("description")
		// date := c.FormValue("date")
		categoryID := c.FormValue("category_id")
		users := slices.DeleteFunc(strings.Split(c.FormValue("users"), ","), func(s string) bool { return s == "" })

		req := managing.CreateInstallmentReq{
			Amount:      amount,
			Description: description,
			// Start: ,
			// End:         date,
			CategoryID: categoryID,
			UsersID:    users,
		}
		err = im.Create(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not create installment.",
			})
		}
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Success",
			"Msg":   "Installment created for " + categoryID,
		})
	}
}

func DeleteInstallment(im managing.InstallmentManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Append("Hx-Trigger", "reloadInstallmentsConfig")
		req := managing.DeleteInstallmentReq{
			ID: c.Params("id"),
		}
		err := im.Delete(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not delete installment.",
			})
		}
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Success",
			"Msg":   "Installment deleted",
		})
	}
}
