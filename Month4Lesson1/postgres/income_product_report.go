
func (i *incomeProductRepo) IncomeProductList(ctx context.Context, request models.IncomeProductReportRequest) (models.IncomeProductReportList, error) {
	var (
		page                      = request.Page
		offset                    = (page - 1) * request.Limit
		pagination, query, filter string
		overallPrice              int
		incomeProducts            []models.IncomeProductReport
	)
	pagination = ` limit $1 offset $2`

	if request.From != "" {
		filter += fmt.Sprintf(` and i.created_at::text <= '%s'`, request.From)
	}

	if request.To != "" {
		filter += fmt.Sprintf(` and i.created_at::text >= '%s'`, request.To)
	}

	if request.BranchID != "" {
		filter += fmt.Sprintf(` and i.branch_id = '%s'`, request.BranchID)
	}

	query = `select 
	p.name, 
	i.quantity, 
	i.price, 
	i.quantity*i.price
	 from income_products as i 
                inner join products as p ON i.product_id = p.id where i.deleted_at = 0 ` + filter + pagination

	rows, err := i.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		i.log.Error("error while income product report", logger.Error(err))
		return models.IncomeProductReportList{}, err
	}
	for rows.Next() {
		incomeProduct := models.IncomeProductReport{}
		if err := rows.Scan(
			&incomeProduct.ProductName,
			&incomeProduct.Quantity,
			&incomeProduct.Price,
			&incomeProduct.TotalPrice,
		); err != nil {
			i.log.Error("error is while scanning all from income products", logger.Error(err))
			return models.IncomeProductReportList{}, err
		}
		overallPrice += incomeProduct.TotalPrice
		incomeProducts = append(incomeProducts, incomeProduct)
	}
	return models.IncomeProductReportList{
		IncomeProducts: incomeProducts,
		OverallPrice:   overallPrice,
	}, err
}