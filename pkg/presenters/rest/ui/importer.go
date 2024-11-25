package ui

import (
	"encoding/csv"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/domain/expense"
	"fmt"
	"io"
	"log/slog"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoadGenericImporter(mu managing.UserManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		respUsers, err := mu.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "User error",
				"Msg":   "Could not load users.",
			})
		}
		return c.Render("sections/importers/generic", fiber.Map{
			"Users": respUsers.Users,
		})
	}
}

func LoadRevolutImporter() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("sections/importers/revolut", fiber.Map{})
	}
}
func ImportGenericCSV(ec tracking.ExpenseCreator, eca tracking.RuleApplier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		csvOrder := strings.Split(c.FormValue("csvOrder"), ",")
		fmt.Println(csvOrder)
		useRules := c.FormValue("useRules") == "checked"
		selectedUsers := slices.DeleteFunc(strings.Split(c.FormValue("users"), ","), func(s string) bool { return s == "" })
		fmt.Println(selectedUsers)
		var matched, skipped, total uint = 0, 0, 0
		failedLines := []uint{}
		file, err := c.FormFile("genericCSV")
		if err != nil || file == nil {
			slog.Error("error", err)
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   fmt.Sprintln("Error reading CSV req:", err),
			})
		}
		uploadedFile, err := file.Open()
		if err != nil {
			slog.Error("error", err)
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   fmt.Sprintln("Error reading CSV req:", err),
			})
		}
		defer uploadedFile.Close()
		csvReader := csv.NewReader(uploadedFile)
		// Iterate over the CSV records
		for {
			line, err := csvReader.Read()
			if err == io.EOF { // Finished processing CSV
				break
			}
			if err != nil || len(line) != 4 {
				return c.Render("alerts/toastErr", fiber.Map{
					"Msg": "Error importing CSV, check Header order.",
				})
			}
			total++
			date, err := time.Parse("2006-01-02", line[slices.Index(csvOrder, "Date")])
			if err != nil {
				slog.Error("error", err)
				failedLines = append(failedLines, total)
				continue
			}
			// Exclude all incomes (i.e amount > 0, this includes internal transfers between spaces)
			amount, err := strconv.ParseFloat(strings.ReplaceAll(line[slices.Index(csvOrder, "Amount")], " ", ""), 64)
			if err != nil {
				slog.Error("error", err)
				failedLines = append(failedLines, total)
				continue
			}
			req := tracking.CreateExpenseReq{
				Product:    line[slices.Index(csvOrder, "Product")],
				Shop:       line[slices.Index(csvOrder, "Shop")],
				Amount:     amount,
				Date:       date,
				UsersID:    selectedUsers,
				CategoryID: expense.UnkownCategoryID,
			}
			if useRules {
				resp := eca.Apply(tracking.ApplyRuleReq{
					Product: req.Product,
					Shop:    req.Shop,
				})
				if resp.Matched {
					matched++
					req.CategoryID = resp.CategoryID
					req.UsersID = resp.UsersID
				}
			}
			_, err = ec.Create(req)
			if err != nil {
				failedLines = append(failedLines, total)
				slog.Error("Can't create expense:", err)
			}
		}
		c.Append("HX-Trigger", "reloadImportTable")
		if int(total)-len(failedLines) == 0 {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Created",
				"Msg":   fmt.Sprintf("Success: %d\nFailed: %d\nSkipped:%d\nAutomatically categorized:%d\nFailed lines:%v", int(total)-len(failedLines), len(failedLines), skipped, matched, failedLines),
			})
		}
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   fmt.Sprintf("Success: %d\nFailed: %d\nSkipped:%d\nAutomatically categorized:%d\nFailed lines:%v", int(total)-len(failedLines), len(failedLines), skipped, matched, failedLines),
		})
	}
}

func LoadImportersTable(eq querying.ExpenseQuerier, cq querying.CategoryQuerier, mu managing.UserManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		pageNum, err := strconv.Atoi(c.Query("page_num", DEFAULT_PNUM_PARAM))
		if err != nil {
			panic("Atoi parse error")
		}
		pageSize, err := strconv.Atoi(c.Query("page_size", DEFAULT_PSIZE_PARAM))
		if err != nil {
			panic("Atoi parse error")
		}
		req := querying.ExpenseQuerierReq{
			Page:        uint(pageNum),
			MaxPageSize: uint(pageSize),
			ExpenseFilter: querying.ExpenseQuerierFilter{
				ByCategoryID: []string{expense.UnkownCategoryID},
			},
		}
		re, err := eq.Query(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "User error",
				"Msg":   err,
			})
		}
		rc, err := cq.Query()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		respUsers, err := mu.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "User error",
				"Msg":   "Could not load users.",
			})
		}

		return c.Render("sections/importers/table", fiber.Map{
			"Expenses":      re.Expenses,
			"NoUserID":      querying.NoUserID,
			"Users":         respUsers.Users,
			"Categories":    rc.Categories,
			"CurrentPage":   req.Page,    // Add this line
			"NextPage":      re.Page + 1, // Add this line
			"PrevPage":      re.Page - 1, // Add this line
			"PageSize":      re.PageSize,
			"ExpensesCount": re.ExpensesCount,
			"TotalPages":    uint(math.Ceil(float64(re.ExpensesCount / req.MaxPageSize))),
		})
	}
}
