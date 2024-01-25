package models

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateCategory struct {
	Name string `json:"name"`
}

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
	Count      int        `json:"count"`
}
