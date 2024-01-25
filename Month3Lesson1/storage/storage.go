package storage

import "test/api/models"

type IStorage interface {
	Close()
	User() IUserStorage
	Basket() IBasketStorage
	Category() ICategoryStorage
	Product() IProductStorage
	BasketProduct() IBasketProductStorage
}

type IUserStorage interface {
	Create(createUser models.CreateUser) (string, error)
	GetByID(models.PrimaryKey) (models.User, error)
	GetList(models.GetListRequest) (models.UsersResponse, error)
	Update(models.UpdateUser) (string, error)
	Delete(models.PrimaryKey) error
	GetPassword(id string) (string, error)
	UpdatePassword(models.UpdateUserPassword) error
}

type IBasketStorage interface {
	Create(createBasket models.CreateBasket) (string, error)
	GetByID(pKey models.PrimaryKey) (models.Basket, error)
	GetList(models.GetListRequest) (models.BasketResponse, error)
	Update(models.UpdateBasket) (string, error)
	Delete(models.PrimaryKey) error
}

type ICategoryStorage interface {
	Create(models.CreateCategory) (string, error)
	GetByID(models.PrimaryKey) (models.Category, error)
	GetList(models.GetListRequest) (models.CategoriesResponse, error)
	Update(models.Category) (string, error)
	Delete(models.PrimaryKey) error
}

type IProductStorage interface {
	Create(models.CreateProduct) (string, error)
	GetByID(models.PrimaryKey) (models.Product, error)
	GetList(models.GetListRequest) (models.ProductsResponse, error)
	Update(models.UpdateProduct) (string, error)
	Delete(models.PrimaryKey) error
}

type IBasketProductStorage interface {
	Create(models.CreateBasketProduct) (string, error)
	GetByID(models.PrimaryKey) (models.BasketProduct, error)
	GetList(models.GetListRequest) (models.BasketProductResponse, error)
	Update(models.UpdateBasketProduct) (string, error)
	Delete(models.PrimaryKey) error
}