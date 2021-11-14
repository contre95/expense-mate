package importers

import (
	"context"
	"expenses-app/pkg/app/importing"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsImporter struct {
	srv           *sheets.Service // DIP for better testeability. I'm not testing this though
	spreadsheetId string
	pageRange     string
}

func NewSheetsImporter(srv *sheets.Service, sheetID, pageRange string) *SheetsImporter {
	return &SheetsImporter{srv, sheetID, pageRange}
}

func NewSheetService(credPath string) (*sheets.Service, error) {
	ctx := context.Background()
	fmt.Println(credPath)
	srv, err := sheets.NewService(ctx, option.WithServiceAccountFile(credPath))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return nil, err
	}
	return srv, nil
}

//GetAllCategories() ([]string, error)
func (si *SheetsImporter) GetImportedExpenses() ([]importing.ImportedExpense, error) {
	//readRange := "Expenses!A2:E"
	resp, err := si.srv.Spreadsheets.Values.Get(si.spreadsheetId, si.pageRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	importedExpenses := []importing.ImportedExpense{}
	if len(resp.Values) == 0 {
		log.Println("No data found")
	} else {
		//TODO: fix this values
		for _, row := range resp.Values {
			e := importing.ImportedExpense{
				Amount:   10,
				Currency: row[1].(string),
				Product:  row[2].(string),
				Shop:     row[4].(string),
				Date:     time.Now(),
				City:     row[5].(string),
				Town:     row[6].(string),
				Category: row[3].(string),
			}
			importedExpenses = append(importedExpenses, e)
		}
	}
	return importedExpenses, nil
}
