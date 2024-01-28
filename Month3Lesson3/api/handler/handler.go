package handler

import (
	"fmt"
	"test/api/models"
	"test/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage storage.IStorage
}

func New(store storage.IStorage) Handler {
	return Handler{
		storage: store,
	}
}

func handleResponse(c *gin.Context, msg string, statusCode int, data interface{}) {

	resp := models.Response{}

	switch code := statusCode; {
	case code < 400:
		resp.Description = "succes"
	case code < 500:
		resp.Description = "bad request"
		fmt.Println("BAD REQUEST: "+msg, "reason: ", data)
	default:
		resp.Description = "internal server error"
		fmt.Println("INTERNAL SERVER ERROR: "+msg, "reason: ", data)

	}

	resp.StatusCode = statusCode
	resp.Data = data

	c.JSON(resp.StatusCode, resp)

}