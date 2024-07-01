package ui

import (
	"encoding/csv"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"fmt"
	"io"
	"log/slog"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoadN26Importer() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("sections/importers/n26", fiber.Map{})
	}
}

func LoadRevolutImporter() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("sections/importers/revolut", fiber.Map{})
	}
}

func ImportN26CSV(t tracking.ExpenseCreator) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		includeSpaces := c.FormValue("spacesTransactions") == "checked"
		includeTransfers := c.FormValue("externalTransactions") == "checked"
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
		failedImports := []uint{}
		var lineNumber uint
		for {
			line, err := csvReader.Read()
			if err == io.EOF { // Finished processing CSV
				break
			}
			if err != nil {
				return c.Render("toastErr", fiber.Map{"Title": "Error", "Msg": fmt.Sprintln("Error reading CSV req:", err)})
			}
			lineNumber++
			// Exclude outgoing transactions (Direct Debit not included)
			// Example1: "2024-03-02","Esteban","ES6315636852323267845001","MoneyBeam","MoneyBeam","-25.2","","",""
			// Example2: "2024-01-29","PEDRO GONZALES","ES0355491146272210003281","Income","Sin concepto","7.6","","",""
			if !includeTransfers && line[2] != "" && line[3] == "Outgoing Transfer" { // Extenal transfers conaint IBAN of the recipients in the 3rd col of csv
				// Example:
				slog.Debug("Excluiding income", "Income", line[3])
				continue
			}
			// Exclude transfers between spaces
			// Example1: "2024-02-22","From Main to Holidays","","Outgoing Transfer","2x Round-up","-1.0","","",""
			// Example2: "2024-02-22","From Car to Main","","Income","PAPELERIA Y LLIB LLORE","4.5","","",""
			internalTransferPattern := regexp.MustCompile(`^From \S+ to \S+$`)
			if !includeSpaces && internalTransferPattern.MatchString(line[1]) {
				slog.Debug("Excluiding internal transfer", "Internal transfer", line[1])
				continue
			}
			date, err := time.Parse("2006-01-02", line[0])
			if err != nil {
				slog.Debug("Excluiding internal transfer", "Internal transfer", line[1])
				failedImports = append(failedImports, lineNumber)
				continue
			}
			// Exclude all incomes (i.e amount > 0, this includes internal transfers between spaces)
			amount, err := strconv.ParseFloat(line[5], 64)
			if err != nil {
				slog.Debug("Excluiding internal transfer", "Internal transfer", line[1])
				failedImports = append(failedImports, lineNumber)
				continue
			}
			if amount > 0 {
				slog.Debug("Excluiding income", "Income", line[1])
				continue
			}
			req := tracking.CreateExpenseReq{
				Product:    line[3],
				Amount:     amount * -1,
				Currency:   "Euro",
				Shop:       line[1],
				City:       "City",
				Date:       date,
				People:     "People",
				CategoryID: "unknown",
			}
			_, err = t.Create(req)
			if err != nil {
				failedImports = append(failedImports, lineNumber)
				slog.Error("Can't create expense:", err)
			}
		}
		c.Append("HX-Trigger", "reloadImportTable")
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   fmt.Sprintf("%d Expenses imported successfully. %d failed", int(lineNumber)-len(failedImports), len(failedImports)),
		})
	}
}

func LoadImportersTable(eq querying.ExpenseQuerier, cq querying.CategoryQuerier) func(*fiber.Ctx) error {
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
				ByCategoryID: []string{"unknown"},
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
		return c.Render("sections/importers/table", fiber.Map{
			"Expenses":      re.Expenses,
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
