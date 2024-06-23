package rest

type categoryJSON struct {
	Name string `json:"name"`
}

type expenseImporterJSON struct {
	BypassWrongExpenses bool `json:"bypass_wrong_expenses"`
}

type usersJSON struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}
