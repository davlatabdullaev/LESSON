package models

type Product struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	OriginalPrice int    `json:"original_price"`
	Quantity      int    `json:"quantity"`
	CategoryID    string `json:"category_id"`
}

type CreateProduct struct {
	Name          string `json:"name"`
	Price         int    `json:"price"`
	OriginalPrice int    `json:"original_price"`
	Quantity      int    `json:"quantity"`
	CategoryID    string `json:"category_id"`
}

type UpdateProduct struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	OriginalPrice int    `json:"original_price"`
	Quantity      int    `json:"quantity"`
	CategoryID    string `json:"category_id"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}


type ProductSell struct {
	Name          string `json:"name"`
	Price         uint   `json:"-"`
	OriginalPrice uint   `json:"-"`
	Quantity      uint   `json:"quantity"`
}

type SellRequest struct {
	Products map[string]int `json:"products"`
	BasketID string         `json:"basket_id"`
}