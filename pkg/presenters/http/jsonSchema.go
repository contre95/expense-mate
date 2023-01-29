package http

type categoriesJSON struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type expenseImporterJSON struct {
	BypassWrongExpenses bool `json:"bypass_wrong_expenses"`
}

type usersJSON struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}
