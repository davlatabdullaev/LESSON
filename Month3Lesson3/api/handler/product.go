package handler

import (
	"errors"
	"net/http"
	"strconv"
	"test/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ShowAccount godoc
// @Router       /product [POST]
// @Summary      Creates a new product
// @Description  create a new product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        user1 body models.CreateProduct  true  "user2"
// @Success      201  {object}  models.Product
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateProduct(c *gin.Context) {
	createProduct := models.CreateProduct{}

	if err := c.ShouldBindJSON(&createProduct); err != nil {
		handleResponse(c, "error while reading body from client ", http.StatusBadRequest, err)
		return
	}

	pKey, err := h.storage.Product().Create(createProduct)
	if err != nil {
		handleResponse(c, "error while creating product", http.StatusInternalServerError, err)
		return
	}

	product, err := h.storage.Product().GetByID(models.PrimaryKey{
		ID: pKey,
	})

	if err != nil {
		handleResponse(c, "error while getting product by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "data created succesfully ", http.StatusCreated, product)
}

// GetProduct godoc
// @Router       /product/{id} [GET]
// @Summary      Gets product
// @Description  Get product by id
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id path string true  "product"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetProduct(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "id is not uuid", http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.storage.Product().GetByID(models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get product by id", http.StatusInternalServerError, err)
		return
	}
	handleResponse(c, "", http.StatusOK, product)
}

// GetProductList godoc
// @Router       /product [GET]
// @Summary      Get user list
// @Description  Get user list
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.ProductsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetProductList(c *gin.Context) {
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

	response, err := h.storage.Product().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting product", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, response)
}

// UpdateProduct godoc
// @Router       /product/{id} [PUT]
// @Summary      Update product
// @Description  Update product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id path string  true  "product_id"
// @Param        product body models.UpdateProduct true "product"
// @Success      201  {object}  models.Product
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateProduct(c *gin.Context) {
	updateProduct := models.UpdateProduct{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateProduct.ID = uid

	if err := c.ShouldBindJSON(&updateProduct); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	pKey, err := h.storage.Product().Update(updateProduct)
	if err != nil {
		handleResponse(c, "error while updating product", http.StatusInternalServerError, err.Error())
		return
	}

	product, err := h.storage.Product().GetByID(models.PrimaryKey{
		ID: pKey,
	})

	if err != nil {
		handleResponse(c, "error while getting product by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, product)

}

// DeleteProduct godoc
// @Router       /product/{id} [DELETE]
// @Summary      Delete product
// @Description  Delete product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id path string  true  "product_id"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteProduct(c *gin.Context) {
	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err = h.storage.Product().Delete(models.PrimaryKey{
		ID: id.String(),
	}); err != nil {
		handleResponse(c, "error while deleting product by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}

// StartSellNew godoc
// @Router               /sell-new [POST]
// @Summary              Selling products
// @Description          selling products
// @Tags                 product
// @Accept               json
// @Produce              json
// @Param                sell_request body models.SellRequest false "sell_request"
// @Success              200  {object}  models.Response
// @Failure              400  {object}  models.Response
// @Failure              404  {object}  models.Response
// @Failure              500  {object}  models.Response
func (h Handler) StartSellNew(c *gin.Context) {
	request := models.SellRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	selectedProducts, productPrices, err := h.storage.Product().Search(request.Products)
	if err != nil {
		handleResponse(c, "error while searching products", http.StatusInternalServerError, err.Error())
		return
	}
	basket, err := h.storage.Basket().GetByID(models.PrimaryKey{
		ID: request.BasketID,
	})
	if err != nil {
		handleResponse(c, "error while searching products", http.StatusInternalServerError, err.Error())
		return
	}
	customer, err := h.storage.User().GetByID(models.PrimaryKey{
		ID: basket.CustomerID,
	})
	if err != nil {
		handleResponse(c, "error while getting customer by id", http.StatusInternalServerError, err.Error())
		return
	}
	totalSum, profit := 0, 0
	basketProducts := map[string]int{}

	for productID, price := range selectedProducts {
		customerQuantity := request.Products[productID]
		totalSum += price * customerQuantity

		profit += customerQuantity * (price - productPrices[productID])
		basketProducts[productID] = customerQuantity
	}

	if customer.Cash < uint(totalSum) {
		handleResponse(c, "not enough customer cash", http.StatusBadRequest, err.Error())
		return
	}

	if err = h.storage.Product().TakeProduct(basketProducts); err != nil {
		handleResponse(c, "error while taking products", http.StatusBadRequest, err.Error())
		return
	}

	if err = h.storage.BasketProduct().AddProducts(basket.ID, basketProducts); err != nil {
		handleResponse(c, "error while adding products", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "succesfully finished the purchase", http.StatusOK, productPrices)

}
