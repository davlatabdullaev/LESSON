package postgres

import (
	"basa/config"
	"database/sql"
	"fmt"
	_"github.com/lib/pq"
)



type Store struct{
	DB                      *sql.DB
	UsersStorage            usersRepo
    OrdersStorage           ordersRepo
	ProductsStorage         productsRepo
	OrderProductsStorage    orderProductsRepo
}

func New(cfg config.Config) (Store, error) {
     url  := fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s sslmode=disable", cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	db, err := sql.Open("postgres", url)
	if err!=nil{
		return Store{}, err
	}

	usersRepo := NewUsersRepo(db)
    ordersRepo := NewOrdersRepo(db)
	productsRepo := NewProductsRepo(db)
	orderProductsRepo := NewOrderProductRepo(db)


	return Store{
		DB: db,
		UsersStorage: usersRepo,
		OrdersStorage: ordersRepo,
		ProductsStorage: productsRepo,
		OrderProductsStorage: orderProductsRepo,
	}, nil
}