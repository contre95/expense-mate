package importers

import (
	"expenses-app/pkg/app/importing"
	"time"
)

type ExampleImporter struct {
	exampleData string
}

func NewExampleImporter(data string) *ExampleImporter {
	return &ExampleImporter{data}
}

//GetAllCategories() ([]string, error)
func (i *ExampleImporter) GetImportedExpenses() ([]importing.ImportedExpense, error) {
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
		},
		{
			Amount:   20,
			Currency: "euro",
			Product:  "Guitarr",
			Shop:     "Local Music Shop",
			Date:     time.Now(),
			City:     "Barcelona",
			Town:     "Spain",
			Category: "Hogar",
		},
		{
			Amount:   1.5,
			Currency: "euro",
			Product:  "Beer",
			Shop:     "Paki's Shop",
			Date:     time.Now().Add(time.Hour * 24 * 5),
			City:     "Barcelona",
			Town:     "Spain",
			Category: "Alimentos",
		},
	}, nil
}
