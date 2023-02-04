package importers

import (
	"expenses-app/pkg/app/importing"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

type SheetsImporter struct {
	srv           *sheets.Service // DIP for better testeability. I'm not testing this though
	spreadsheetId string
	pageRange     string
	rangeLenght   int
}

func NewSheetsImporter(srv *sheets.Service, sheetID, pageRange string, rangeLenght int) *SheetsImporter {
	return &SheetsImporter{srv, sheetID, pageRange, rangeLenght}
}

// GetAllCategories() ([]string, error)
func (si *SheetsImporter) GetImportedExpenses() ([]importing.ImportedExpense, error) {
	//readRange := "Expenses!A2:E"
	resp, err := si.srv.Spreadsheets.Values.Get(si.spreadsheetId, si.pageRange).Do()
	if err != nil {
		return nil, err
	}
	importedExpenses := []importing.ImportedExpense{}
	if len(resp.Values) == 0 {
		log.Println("No data found")
	} else {
		//TODO: fix this values
		for _, row := range resp.Values[1:] {
			if len(row) < si.rangeLenght {
				log.Printf("Skipping row %v of lenthg %d", row, len(row))
				continue
			}
			date, dateErr := time.Parse("1/2/2006", row[7].(string))
			if dateErr != nil {
				fmt.Printf("Couldn't parase date from product %s\n", row[2])
				return nil, dateErr
			}
			price, priceErr := strconv.ParseFloat(row[0].(string), 32)
			if priceErr != nil {
				fmt.Printf("Couldn't parase price from product %s\n", row[2])
				return nil, priceErr
			}
			e := importing.ImportedExpense{
				Amount:   price,
				Currency: row[1].(string),
				Product:  row[2].(string),
				Shop:     strings.TrimSpace(row[4].(string)),
				Date:     date,
				City:     row[5].(string),
				Town:     row[6].(string),
				People:   row[8].(string),
				Category: row[3].(string),
			}
			importedExpenses = append(importedExpenses, e)
		}
	}
	return importedExpenses, nil
}
