package api

import (
	"test/api/handler"
	"test/storage"

	"github.com/gin-gonic/gin"
)


func New(store storage.IStorage) *gin.Engine {
    h := handler.New(store)

	r := gin.New()

	r.POST("/user", h.CreateUser)
	r.GET("/users", h.GetUserList)
	r.GET("/user/:id", h.GetUser)
	r.PUT("/user/:id", h.UpdateUser)
	r.DELETE("user/:id", h.DeleteUser)
	r.PATCH("user/:id", h.UpdateUserPassword)

	return r
}