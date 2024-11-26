package ui

import (
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func LoadInstallmentsAddRow(im managing.InstallmentManager, um managing.UserManager, cq querying.CategoryQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("asdf")
		usersResp, err := um.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not load users.",
			})
		}
		installmentsResp, err := im.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not load installments.",
			})
		}
		rc, err := cq.Query()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}

		return c.Render("sections/installments/rowAdd", fiber.Map{
			"Installments": installmentsResp.Installments,
			"Users":        usersResp.Users,
			"Categories":   rc.Categories,
		})
	}
}

func LoadInstallmentsTable(im managing.InstallmentManager, um managing.UserManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		usersResp, err := um.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not load users.",
			})
		}
		installmentsResp, err := im.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   "Could not load installments.",
			})
		}
		return c.Render("sections/installments/table", fiber.Map{
			"Installments": installmentsResp.Installments,
			"Users":        usersResp.Users,
		})
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
				"Msg":   "Could not delete installment.",
				"Title": "Error",
			})
		}
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Success",
			"Msg":   "Installment deleted",
		})
	}
}
