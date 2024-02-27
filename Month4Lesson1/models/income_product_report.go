
type IncomeProductReport struct {
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	TotalPrice  int    `json:"total_price"`
}

type IncomeProductReportList struct {
	IncomeProducts []IncomeProductReport `json:"income_products"`
	OverallPrice   int                   `json:"overall_price"`
}

type IncomeProductReportRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	BranchID string `json:"branch_id"`
	From     string `json:"from"`
	To       string `json:"to"`
}