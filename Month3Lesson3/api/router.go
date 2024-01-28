package api

import (
	"test/api/handler"
	"test/storage"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	_"test/api/docs"
)

// New ...
// @title              Swagger Example API
// @version            1.0
// @description        This is a sample server celler server
func New(store storage.IStorage) *gin.Engine {
	h := handler.New(store)

	r := gin.New()

	r.POST("/user", h.CreateUser)
	r.GET("/user", h.GetUserList)
	r.GET("/user/:id", h.GetUser)
	r.PUT("/user/:id", h.UpdateUser)
	r.DELETE("user/:id", h.DeleteUser)
	r.PATCH("user/:id", h.UpdateUserPassword)

	r.POST("/category", h.CreateCategory)
	r.GET("/category/:id", h.GetCategory)
	r.GET("/category", h.GetCategoriesList)
	r.PUT("/category/:id", h.UpdateCategory)
	r.DELETE("/category/:id", h.DeleteCategory)

	r.POST("/product", h.CreateProduct)
	r.GET("/product/:id", h.GetProduct)
	r.GET("/product", h.GetProductList)
	r.PUT("/product/:id", h.UpdateProduct)
	r.DELETE("/product/:id", h.DeleteProduct)

	r.POST("/basket", h.CreateBasket)
	r.GET("/basket/:id", h.GetBasket)
	r.GET("/basket", h.GetBasketList)
	r.PUT("basket/:id", h.UpdateBasket)
	r.DELETE("basket/:id", h.DeleteBasket)

	r.POST("/basketProduct", h.CreateBasketProduct)
	r.GET("/basketProduct/:id", h.GetBasketProduct)
	r.GET("/basketProduct", h.GetBasketProductList)
	r.PUT("/basketProduct/:id", h.UpdateBasketProduct)
	r.DELETE("/basketProduct/:id", h.DeleteBasketProduct)
	
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
	return r
}
