package importing

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
	"time"
)

type ImporertResp struct {
	ID  string
	Msg string
}

type ImporertReq struct {
}

// Importer (is a Repository ?) is an Interface that defines a dependency for the importing expenses use cases
// TODO: Está bien que esto viva en el dominio ? No, porque no tiene que
// ver con mi dominio, es una feature más de mi app que me permite
// trackear mi expenses. Interfáz tiene que vivir en este mismo archivo,
// lo más cerca de donde se vaya a usar posible.
type ImportedExpense struct {
	Product string
	Shop    string
	Date    time.Time
	City    string
	Town    string

	Category string
}

type Importer interface {
	//GetAllCategories() ([]string, error)
	GetImpoertedExpenses(categoryName string) ([]ImportedExpense, error)
}

// The createCategory use case creates a category for a expense
type ImportUseCase struct {
	logger   app.Logger
	importer Importer
	expenses expense.Expenses
}

func NewImporterUseCase(l app.Logger, i Importer, e expense.Expenses) *ImportUseCase {
	return &ImportUseCase{l, i, e}
}

// Create use cases function creates a new category
func (u *ImportUseCase) Import(req ImporertReq) (*ImporertResp, error) {

}
