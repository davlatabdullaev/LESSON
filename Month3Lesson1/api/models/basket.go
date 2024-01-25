package models

type Basket struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	TotalSum   string `json:"total_sum"`
}

type CreateBasket struct {
	CustomerID string `json:"customer_id"`
	TotalSum   uint   `json:"total_sum"`
}

type UpdateBasket struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	TotalSum   uint   `json:"total_sum"`
}

type BasketResponse struct {
	Baskets []Basket `json:"baskets"`
	Count   int      `json:"count"`
}
