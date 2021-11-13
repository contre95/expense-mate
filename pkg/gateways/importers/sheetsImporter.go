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

	if len(resp.Values) == 0 {
		log.Println("No data found")
	} else {
		for _, row := range resp.Values {
			fmt.Println(row)
		}
	}
	return []importing.ImportedExpense{
		{
			Amount:   1.0,
			Currency: "euro",
			Product:  "Wine",
			Shop:     "Mercadona",
			Date:     time.Now(),
			City:     "Barcelona",
			Town:     "Spain",
			Category: "Alimentos",
		}}, nil
}
