package handler

import (
	"errors"
	"net/http"
	"strconv"
	"test/api/models"
	"test/pkg/check"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ShowAccount godoc
// @Router       /user [POST]
// @Summary      Creates a new user
// @Description  create a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user1 body models.CreateUser  true  "user2"
// @Success      201  {object}  models.User
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateUser(c *gin.Context) {
	createUser := models.CreateUser{}

	if err := c.ShouldBindJSON(&createUser); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
		return
	}

	pKey, err := h.storage.User().Create(createUser)
	if err != nil {
		handleResponse(c, "error while creating user ", http.StatusInternalServerError, err)
		return
	}

	user, err := h.storage.User().GetByID(models.PrimaryKey{
		ID: pKey,
	})

	if err != nil {
		handleResponse(c, "error while getting user by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "data created succesfully", http.StatusCreated, user)

}

// GetUser godoc
// @Router       /user/{id} [GET]
// @Summary      Gets user
// @Description  Get user by id
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path string true  "user"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetUser(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "id is not uuid", http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.storage.User().GetByID(models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while getting user by id", http.StatusInternalServerError, err)
		return
	}
	handleResponse(c, "", http.StatusOK, user)
}

// GetUserList godoc
// @Router       /user [GET]
// @Summary      Get user list
// @Description  Get user list
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.UsersResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetUserList(c *gin.Context) {
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

	response, err := h.storage.User().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting users", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

// UpdateUser godoc
// @Router       /user/{id} [PUT]
// @Summary      Update user
// @Description  Update user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path string  true  "user_id"
// @Param        user body models.UpdateUser true "user"
// @Success      201  {object}  models.User
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateUser(c *gin.Context) {
	updateUser := models.UpdateUser{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateUser.ID = uid

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	pKey, err := h.storage.User().Update(updateUser)
	if err != nil {
		handleResponse(c, "error while updating user", http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.storage.User().GetByID(models.PrimaryKey{
		ID: pKey,
	})

	if err != nil {
		handleResponse(c, "error while getting user by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, user)

}

// DeleteUser godoc
// @Router       /user/{id} [DELETE]
// @Summary      Delete user
// @Description  Delete user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path string  true  "user_id"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteUser(c *gin.Context) {
	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err = h.storage.User().Delete(models.PrimaryKey{
		ID: id.String(),
	}); err != nil {
		handleResponse(c, "error while deleting user by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}

// UpdateUserpassword godoc
// @Router       /user/{id} [PATCH]
// @Summary      update user
// @Description  update user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Param        user body models.UpdateUserPassword true "user"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateUserPassword(c *gin.Context) {
	updateUserPassword := models.UpdateUserPassword{}

	if err := c.ShouldBindJSON(&updateUserPassword); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err)
		return
	}

	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		handleResponse(c, "error while parsing uuid", http.StatusBadRequest, err.Error())
		return
	}

	updateUserPassword.ID = uid.String()

	oldPassword, err := h.storage.User().GetPassword(updateUserPassword.ID)
	if err != nil {
		handleResponse(c, "error while getting password by id", http.StatusInternalServerError, err.Error())
		return
	}

	if oldPassword != updateUserPassword.OldPassword {
		handleResponse(c, "old password is not correct", http.StatusBadRequest, "old password is not correct")
		return
	}

	if err = check.ValidatePassword(updateUserPassword.NewPassword); err != nil {
		handleResponse(c, "new password is weak", http.StatusBadRequest, err.Error())
		return
	}

	if err = h.storage.User().UpdatePassword(updateUserPassword); err != nil {
		handleResponse(c, "error while updating user password by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "password succesfully updated")

}
