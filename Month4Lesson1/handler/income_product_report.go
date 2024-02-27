
// IncomeProductReport godoc
// @Router       /income_product/report [GET]
// @Summary      Get  income product report
// @Description  get income product report
// @Tags         report
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        from query string false "from"
// @Param        to query string false "to"
// @Param        branch_id query string false "branch_id"
// @Success      201  {object}  models.IncomeProductReportList
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) IncomeProductReport(c *gin.Context) {
	var (
		page, limit int
		branchID    string
		from, to    string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, h.log, "error while converting pege", http.StatusBadRequest, err)
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, h.log, "error while converting limit", http.StatusBadRequest, err)
		return
	}

	from = c.Query("from")
	to = c.Query("to")
	branchID = c.Query("branch_id")

	incomeProductReport, err := h.services.IncomeProduct().IncomeProductList(context.Background(), models.IncomeProductReportRequest{
		Page:     page,
		Limit:    limit,
		BranchID: branchID,
		From:     from,
		To:       to,
	})
	if err != nil {
		handleResponse(c, h.log, "error is while getting income product report", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, incomeProductReport)
}