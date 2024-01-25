package handler

import (
	"net/http"
	"strconv"
	"test/api/models"

	"github.com/gin-gonic/gin"
)

func (h Handler) CreateBasketProduct(c *gin.Context) {
	basketProduct := models.CreateBasketProduct{}

	if err := c.ShouldBindJSON(&basketProduct); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.BasketProduct().Create(basketProduct)
	if err != nil {
		handleResponse(c, "error is while creating basket product", http.StatusInternalServerError, err)
		return
	}

	createdBasketProduct, err := h.storage.BasketProduct().GetByID(models.PrimaryKey{ID: id})
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, createdBasketProduct)
}

func (h Handler) GetBasketProduct(c *gin.Context) {
	uid := c.Param("id")

	basketProduct, err := h.storage.BasketProduct().GetByID(models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, basketProduct)
}

func (h Handler) GetBasketProductList(c *gin.Context) {
	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error is while converting page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error is while converting page", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	basketProducts, err := h.storage.BasketProduct().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	handleResponse(c, "", http.StatusOK, basketProducts)
}

func (h Handler) UpdateBasketProduct(c *gin.Context) {
	basketProduct := models.UpdateBasketProduct{}
	uid := c.Param("id")

	if err := c.ShouldBindJSON(&basketProduct); err != nil {
		handleResponse(c, "error is while reading from body", http.StatusBadRequest, err.Error())
		return
	}

	basketProduct.ID = uid
	id, err := h.storage.BasketProduct().Update(basketProduct)
	if err != nil {
		handleResponse(c, "error is while updating basket", http.StatusInternalServerError, err.Error())
		return
	}

	updatedBasketProduct, err := h.storage.BasketProduct().GetByID(models.PrimaryKey{ID: id})
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedBasketProduct)
}

func (h Handler) DeleteBasketProduct(c *gin.Context) {
	uid := c.Param("id")

	if err := h.storage.BasketProduct().Delete(models.PrimaryKey{ID: uid}); err != nil {
		handleResponse(c, "error is while deleting", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "basket product deleted!")
}
