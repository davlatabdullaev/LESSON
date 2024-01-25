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

func (h Handler) GetUser(c *gin.Context) {
	var err error

	uid := c.Param("id")

	user, err := h.storage.User().GetByID(models.PrimaryKey{
		ID: uid,
	})
	if err != nil {
		handleResponse(c, "error while getting user by id", http.StatusInternalServerError, err)
		return
	}
	handleResponse(c, "", http.StatusOK, user)
}

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
