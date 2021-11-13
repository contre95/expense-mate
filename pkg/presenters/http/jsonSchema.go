package http

type addCategoriesJSON struct {
	Names []string `json:"categories"`
}

type addCategoryJSON struct {
	Name string `json:"name"`
}

type expenseImporterJSON struct {
	ByPassWrongExpenses bool `json:"by_pass_wrong_expenses"`
}
