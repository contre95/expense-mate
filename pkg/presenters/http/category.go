package http

type addCategoriesJSON struct {
	Names []string `json:"categories"`
}

type addCategoryJSON struct {
	Name string `json:"name"`
}
