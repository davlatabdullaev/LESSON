package postgres

import (
	"context"
	"fmt"
	"log"
	"test/api/models"
	"test/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type basketRepo struct {
	pool *pgxpool.Pool
}

func NewBasketRepo(pool *pgxpool.Pool) storage.IBasketStorage {
	return &basketRepo{
		pool: pool,
	}
}

func (b *basketRepo) Create(createBasket models.CreateBasket) (string, error) {

	uid := uuid.New()

	query := `insert into baskets values ($1, $2, $3)`

	if _, err := b.pool.Exec(context.Background(), query, uid, createBasket.CustomerID, createBasket.TotalSum); err != nil {
		log.Println("error while inserting data", err.Error())
		return "", err
	}

	return uid.String(), nil
}

func (b *basketRepo) GetByID(pKey models.PrimaryKey) (models.Basket, error) {

	basket := models.Basket{}

	query := `select id, customer_id, total_sum from baskets where id = $1`

	err := b.pool.QueryRow(context.Background(), query, pKey.ID).Scan(
		&basket.ID,
		&basket.CustomerID,
		&basket.TotalSum,
	)
	if err != nil {
		log.Println("error while scanning basket", err.Error())
		return models.Basket{}, err
	}

	return basket, nil
}

func (b *basketRepo) GetList(request models.GetListRequest) (models.BasketResponse, error) {

	var (
		baskets           = []models.Basket{}
		count             = 0
		countQuery, query string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `
	SELECT count(1) from baskets
	`
	if search != "" {
		countQuery += fmt.Sprintf(` and (total_sum ilike '%%%s%%')`, search)
	}

	if err := b.pool.QueryRow(context.Background(), countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of baskets", err.Error())
		return models.BasketResponse{}, err
	}

	query = `
	SELECT id, customer_id, total_sum from baskets
	`
	if search != "" {
		query += fmt.Sprintf(` and (total_sum ilike '%%%s%%')`, search)
	}

	query += `LIMIT $1 OFFSET $2`

	rows, err := b.pool.Query(context.Background(), query, request.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.BasketResponse{}, err
	}

	for rows.Next() {
		basket := models.Basket{}

		if err = rows.Scan(
			&basket.ID,
			&basket.CustomerID,
			&basket.TotalSum,
		); err != nil {
			log.Println("error while scanning row", err.Error())
			return models.BasketResponse{}, err
		}

		baskets = append(baskets, basket)

	}

	return models.BasketResponse{
		Baskets: baskets,
		Count:   count,
	}, nil
}

func (b *basketRepo) Update(request models.UpdateBasket) (string, error) {

	query := `update baskets set customer_id = $1, total_sum = $2 where id = $3`

	if _, err := b.pool.Exec(context.Background(), query, request.CustomerID, request.TotalSum, request.ID); err != nil {
		log.Println("error while updating basket data", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (b *basketRepo) Delete(request models.PrimaryKey) error {

	query := `delete from 
	baskets
	 where id = $1`

	if _, err := b.pool.Exec(context.Background(), query, request.ID); err != nil {
		log.Println("error while deleting basket by id", err.Error())
		return err
	}

	return nil
}
