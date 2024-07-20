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
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoadN26Importer(mu managing.UserManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		respUsers, err := mu.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "User error",
				"Msg":   "Could not load users.",
			})
		}
		return c.Render("sections/importers/n26", fiber.Map{
			"Users": respUsers.Users,
		})
	}
}

func LoadRevolutImporter() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("sections/importers/revolut", fiber.Map{})
	}
}
func ImportN26CSV(ec tracking.ExpenseCreator, eca tracking.RuleApplier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		includeSpaces := c.FormValue("spacesTransactions") == "checked"
		includeTransfers := c.FormValue("externalTransactions") == "checked"
		useRules := c.FormValue("useRules") == "checked"
		selectedUsers := slices.DeleteFunc(strings.Split(c.FormValue("users"), ","), func(s string) bool { return s == "" })
		var matched, skipped, total uint = 0, 0, 0
		failedLines := []uint{}
		file, err := c.FormFile("n26csv")
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
		csvReader.Read() // Skip column
		for {
			line, err := csvReader.Read()
			if err == io.EOF { // Finished processing CSV
				break
			}
			if err != nil {
				return c.Render("toastErr", fiber.Map{"Title": "Error", "Msg": fmt.Sprintln("Error reading CSV req:", err)})
			}
			total++
			// Exclude outgoing transactions (Direct Debit not included)
			// Example1: "2024-03-02","Esteban","ES6315636852323267845001","MoneyBeam","MoneyBeam","-25.2","","",""
			// Example2: "2024-01-29","PEDRO GONZALES","ES0355491146272210003281","Income","Sin concepto","7.6","","",""
			if !includeTransfers && line[2] != "" && slices.Contains([]string{"Outgoing Transfer", "MoneyBeam"}, line[3]) { // Extenal transfers conaint IBAN of the recipients in the 3rd col of csv
				// Example:
				slog.Debug("Excluiding income", "Income", line[3])
				skipped++
				continue
			}
			// Exclude transfers between spaces
			// Example1: "2024-02-22","From Main to Holidays","","Outgoing Transfer","2x Round-up","-1.0","","",""
			// Example2: "2024-02-22","From Car to Main","","Income","PAPELERIA Y LLIB LLORE","4.5","","",""
			internalTransferPattern := regexp.MustCompile(`^From \S+ to \S+$`)
			if !includeSpaces && internalTransferPattern.MatchString(line[1]) {
				slog.Debug("Excluiding internal transfer", "Internal transfer", line[1])
				skipped++
				continue
			}
			date, err := time.Parse("2006-01-02", line[0])
			if err != nil {
				slog.Debug("Excluiding internal transfer", "Internal transfer", line[1])
				failedLines = append(failedLines, total)
				continue
			}
			// Exclude all incomes (i.e amount > 0, this includes internal transfers between spaces)
			amount, err := strconv.ParseFloat(line[5], 64)
			if err != nil {
				slog.Debug("Excluiding internal transfer", "Internal transfer", line[1])
				failedLines = append(failedLines, total)
				continue
			}
			if amount > 0 {
				slog.Debug("Excluiding income", "Income", line[1])
				skipped++
				continue
			}
			req := tracking.CreateExpenseReq{
				Product:    line[3],
				Amount:     amount * -1,
				Shop:       line[1],
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
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   fmt.Sprintf("Success: %d\nFailed: %d\nSkipped:%d\nAutomatically categorized:%d\nFailed lines:%v", total, len(failedLines), skipped, matched, failedLines),
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
			panic("Implement error")
		}
		rc, err := cq.Query()
		if err != nil {
			panic("Implement error")
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
