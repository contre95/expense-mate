package ui

import (
	"encoding/csv"
	"expenses-app/pkg/app/tracking"
	"fmt"
	"io"
	"log/slog"
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
		fmt.Println("Importing N26 CSV")
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
			// Exclude all incomes (this includes internal transfers between spaces)
			if line[3] == "Income" {
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
			}
			price, err := strconv.ParseFloat(line[5], 64)
			if err != nil {
				slog.Debug("Excluiding internal transfer", "Internal transfer", line[1])
				failedImports = append(failedImports, lineNumber)
			}
			req := tracking.CreateExpenseReq{
				Product:  line[3],
				Price:    price,
				Currency: "Euro",
				Shop:     line[1],
				City:     "City",
				Date:     date,
				People:   "People",
				Category: "unknown",
			}
			fmt.Println(req)
			resp, err := t.Create(req)
			if err != nil {
				failedImports = append(failedImports, lineNumber)
				slog.Error("Can't create expense:", err)
			}
			fmt.Println(resp)
		}
		fmt.Println("Imported")
		c.Append("HX-Trigger", "uncategorizeTable")
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   fmt.Sprintf("%d Expenses imported successfully. %d failed", int(lineNumber)-len(failedImports), len(failedImports)),
		})
	}
}
func LoadUncotegorizedExpensesTable() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			// c.Append("hx-trigger", "newPair")  // Not working :(
			return c.Render("main", fiber.Map{
				"ImporterTrigger": "revealed",
			})
		}
		return c.Render("sections/importers/index", fiber.Map{})
	}
}

func LoadImporterSection() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			// c.Append("hx-trigger", "newPair")  // Not working :(
			return c.Render("main", fiber.Map{
				"ImporterTrigger": "revealed",
			})
		}
		return c.Render("sections/importers/index", fiber.Map{})
	}
}
