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
type Importer interface {
	GetAllCategories() ([]expense.Category, error) // Está bien pedir que esto me devuelva un objeto de dominio ? Le va a quedar la responsabilidad de inicializar ese objeto de dominio a la infra. Me manejo todo con strings ?
	GetAllExpenses(categoryName string, fromDate time.Time) ([]expense.Expense, error)
}

// The createCategory use case creates a category for a expense
type ImportUseCase struct {
	logger   app.Logger
	importer Importer
	expenses expense.Expenses
}

func NewSheetsImporter(l app.Logger, i Importer, e expense.Expenses) *ImportUseCase {
	return &ImportUseCase{l, i, e}
}

// Create use cases function creates a new category
func (u *ImportUseCase) Import(req ImporertReq) (*ImporertResp, error) {
	return nil, nil
	//categories, err := u.importer.GetAllCategories()
	//if err != nil {
	//return nil, err
	//}
	// En blanco, preguntarle a Lois
}
