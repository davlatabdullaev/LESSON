package handler

import (
	"bazaar/api/models"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateIncome godoc
// @Router       /income [POST]
// @Summary      Create a new income
// @Description  Create a new income
// @Tags         income
// @Accept       json
// @Produce      json
// @Param        income  body  models.CreateIncome  true  "income data"
// @Success      201  {object}  models.Income
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateIncome(c *gin.Context) {
	createIncome := models.CreateIncome{}

	if err := c.ShouldBindJSON(&createIncome); err != nil {
		handleResponse(c, "error while reading income body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Income().Create(context.Background(), createIncome)
	if err != nil {
		handleResponse(c, "error while creating income", http.StatusInternalServerError, err)
		return
	}

	income, err := h.storage.Income().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get income", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, income)

}

// GetIncomeByID godoc
// @Router       /income/{id} [GET]
// @Summary      Get income by id
// @Description  Get income by id
// @Tags         income
// @Accept       json
// @Produce      json
// @Param        id path string true "income"
// @Success      200  {object}  models.Income
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomeByID(c *gin.Context) {

	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	income, err := h.storage.Income().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get income by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, income)

}

// GetIncomesList godoc
// @Router       /incomes [GET]
// @Summary      Get incomes list
// @Description  Get incomes list
// @Tags         income
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.IncomesResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomesList(c *gin.Context) {

	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error while parsing page ", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	response, err := h.storage.Income().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting income", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

// UpdateIncome godoc
// @Router       /income/{id} [PUT]
// @Summary      Update income by id
// @Description  Update income by id
// @Tags         income
// @Accept       json
// @Produce      json
// @Param        id path string true "income id"
// @Param        income body models.UpdateIncome true "income"
// @Success      200  {object}  models.Income
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateIncome(c *gin.Context) {
	updateIncome := models.UpdateIncome{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateIncome.ID = uid

	if err := c.ShouldBindJSON(&updateIncome); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Income().Update(context.Background(), updateIncome)
	if err != nil {
		handleResponse(c, "error while updating income", http.StatusInternalServerError, err.Error())
		return
	}

	income, err := h.storage.Income().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting income by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, income)

}

// DeleteIncome godoc
// @Router       /income/{id} [DELETE]
// @Summary      Delete Income
// @Description  Delete Income
// @Tags         income
// @Accept       json
// @Produce      json
// @Param        id path string true "income id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteIncome(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Income().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting income by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}

// income product

package handler

import (
	"bazaar/api/models"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateIncomeProduct godoc
// @Router       /income_product [POST]
// @Summary      Create a new income product
// @Description  Create a new income product
// @Tags         income_product
// @Accept       json
// @Produce      json
// @Param        income_product  body  models.CreateIncomeProduct  true  "income product data"
// @Success      201  {object}  models.IncomeProduct
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateIncomeProduct(c *gin.Context) {
	createIncomeProduct := models.CreateIncomeProduct{}

	if err := c.ShouldBindJSON(&createIncomeProduct); err != nil {
		handleResponse(c, "error while reading income product body from client", http.StatusBadRequest, err)
		return
	}

	incomeData, err := h.storage.Income().Get(context.Background(), models.PrimaryKey{
		ID: createIncomeProduct.IncomeID,
	})
	if err != nil {
		handleResponse(c, "error while search income for create income product", http.StatusInternalServerError, err)
		return
	}

	storageDataForProduct, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  100,
		Search: createIncomeProduct.ProductID,
	})
	if err != nil {
		handleResponse(c, "error while search storage for create income product", http.StatusInternalServerError, err)
		return
	}

	storageDataForBranch, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  100,
		Search: incomeData.BranchID,
	})
	if err != nil {
		handleResponse(c, "error while search storage for create income product", http.StatusInternalServerError, err)
		return
	}

	if createIncomeProduct.ProductID == storageDataForProduct.Storages[0].ProductID && incomeData.BranchID == storageDataForBranch.Storages[0].BranchID {

		if err := h.storage.Storage().UpdateCount(context.Background(), models.UpdateCount{
			ID:    storageDataForBranch.Storages[0].ID,
			Count: createIncomeProduct.Count,
		}); err != nil {
			handleResponse(c, "error while update storage count", http.StatusInternalServerError, err)
			return
		}

	} else {
		id, err := h.storage.IncomeProduct().Create(context.Background(), createIncomeProduct)
		if err != nil {
			handleResponse(c, "error while creating income product", http.StatusInternalServerError, err)
			return
		}

		incomeProduct, err := h.storage.IncomeProduct().Get(context.Background(), models.PrimaryKey{
			ID: id,
		})
		if err != nil {
			handleResponse(c, "error while get income product", http.StatusInternalServerError, err)
			return
		}

		handleResponse(c, "", http.StatusCreated, incomeProduct)
	}
}

// GetIncomeProductByID godoc
// @Router       /income_product/{id} [GET]
// @Summary      Get income product by id
// @Description  Get income product by id
// @Tags         income_product
// @Accept       json
// @Produce      json
// @Param        id path string true "income_product"
// @Success      200  {object}  models.IncomeProduct
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomeProductByID(c *gin.Context) {

	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	incomeProduct, err := h.storage.IncomeProduct().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get income product by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, incomeProduct)

}

// GetIncomesList godoc
// @Router       /income_products [GET]
// @Summary      Get income products list
// @Description  Get income products list
// @Tags         income_product
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.IncomeProductsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomeProductsList(c *gin.Context) {

	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error while parsing page ", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	response, err := h.storage.IncomeProduct().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting income product", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

// UpdateIncome godoc
// @Router       /income_product/{id} [PUT]
// @Summary      Update income product by id
// @Description  Update income product by id
// @Tags         income_product
// @Accept       json
// @Produce      json
// @Param        id path string true "income id"
// @Param        income_product body models.UpdateIncomeProduct true "income product"
// @Success      200  {object}  models.IncomeProduct
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateIncomeProduct(c *gin.Context) {
	updateIncomeProduct := models.UpdateIncomeProduct{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateIncomeProduct.ID = uid

	if err := c.ShouldBindJSON(&updateIncomeProduct); err != nil {
		handleResponse(c, "error while reading income products body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.IncomeProduct().Update(context.Background(), updateIncomeProduct)
	if err != nil {
		handleResponse(c, "error while updating income product", http.StatusInternalServerError, err.Error())
		return
	}

	incomeProduct, err := h.storage.IncomeProduct().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting income by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, incomeProduct)

}

// DeleteIncomeProduct godoc
// @Router       /income_product/{id} [DELETE]
// @Summary      Delete Income
// @Description  Delete Income
// @Tags         income_product
// @Accept       json
// @Produce      json
// @Param        id path string true "income product id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteIncomeProduct(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.IncomeProduct().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting income product by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}

