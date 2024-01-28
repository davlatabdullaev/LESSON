package handler

import (
	"errors"
	"net/http"
	"strconv"
	"test/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ShowCategory godoc
// @Router       /category [POST]
// @Summary      Creates a new category
// @Description  create a new category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        category1 body models.CreateCategory  true  "category2"
// @Success      201  {object}  models.Category
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateCategory(c *gin.Context) {
	createCategory := models.CreateCategory{}

	if err := c.ShouldBindJSON(&createCategory); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
		return
	}

	pKey, err := h.storage.Category().Create(createCategory)
	if err != nil {
		handleResponse(c, "error while creating category ", http.StatusInternalServerError, err)
		return
	}

	user, err := h.storage.Category().GetByID(models.PrimaryKey{
		ID: pKey,
	})

	if err != nil {
		handleResponse(c, "error while getting c category by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "data created succesfully", http.StatusCreated, user)

}

// GetCategory godoc
// @Router       /category/{id} [GET]
// @Summary      Gets category
// @Description  Get category by id
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id path string true  "category"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetCategory(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "id is not uuid", http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.storage.Category().GetByID(models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while getting user by id", http.StatusInternalServerError, err)
		return
	}
	handleResponse(c, "", http.StatusOK, user)
}

// GetUserCategory godoc
// @Router       /category [GET]
// @Summary      Get category list
// @Description  Get category list
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.CategoriesResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetCategoriesList(c *gin.Context) {
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

	response, err := h.storage.Category().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting categories", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

// UpdateCategory godoc
// @Router       /category/{id} [PUT]
// @Summary      Update category
// @Description  Update category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id path string  true  "category_id"
// @Param        user body models.Category true "user"
// @Success      201  {object}  models.Category
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateCategory(c *gin.Context) {
	updateCategory := models.Category{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateCategory.ID = uid

	if err := c.ShouldBindJSON(&updateCategory); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	pKey, err := h.storage.Category().Update(updateCategory)
	if err != nil {
		handleResponse(c, "error while updating category", http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.storage.Category().GetByID(models.PrimaryKey{
		ID: pKey,
	})

	if err != nil {
		handleResponse(c, "error while getting category by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, user)

}

// DeleteCategory godoc
// @Router       /category/{id} [DELETE]
// @Summary      Delete category
// @Description  Delete category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id path string  true  "category_id"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteCategory(c *gin.Context) {
	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err = h.storage.Category().Delete(models.PrimaryKey{
		ID: id.String(),
	}); err != nil {
		handleResponse(c, "error while deleting  by category id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
