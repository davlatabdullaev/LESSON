package postgres

import (
	"context"
	"fmt"
	"test/config"
	"test/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Store struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config) (storage.IStorage, error) {
	poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB))
	if err != nil {
		fmt.Println("error while parsing config", err.Error())
		return nil, err
	}

	poolConfig.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		fmt.Println("error while connecting to db", err.Error())
		return nil, err
	}

	return Store{
		Pool: pool,
	}, nil

}

func (s Store) Close() {
	s.Pool.Close()
}

func (s Store) User() storage.IUserStorage {
	newUser := NewUserRepo(s.Pool)
	return newUser
}

func (s Store) Basket() storage.IBasketStorage {
	newBasket := NewBasketRepo(s.Pool)
	return newBasket
}

func (s Store) Category() storage.ICategoryStorage {
	newCategory := NewCategoryRepo(s.Pool)
	return newCategory
}

func (s Store) Product() storage.IProductStorage {
	newProduct := NewProductRepo(s.Pool)
	return newProduct
}

func (s Store) BasketProduct() storage.IBasketProductStorage {
	newBasketProduct := NewBasketProductRepo(s.Pool)
	return newBasketProduct
}
