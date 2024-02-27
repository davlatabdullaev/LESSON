
func (i incomeProductService) IncomeProductList(ctx context.Context, request models.IncomeProductReportRequest) (models.IncomeProductReportList, error) {
	incomeProductList, err := i.storage.IncomeProduct().IncomeProductList(ctx, request)
	if err != nil {
		i.log.Error("error in service layer income product report", logger.Error(err))
		return models.IncomeProductReportList{}, err
	}
	return incomeProductList, err
}