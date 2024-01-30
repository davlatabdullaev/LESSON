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

type categoryRepo struct {
	pool *pgxpool.Pool
}

func NewCategoryRepo(pool *pgxpool.Pool) storage.ICategoryStorage {
	return &categoryRepo{
		pool: pool,
	}
}

func (c *categoryRepo) Create(request models.CreateCategory) (string, error) {

	uid := uuid.New()

	if _, err := c.pool.Exec(context.Background(), `insert into categories values ($1, $2)`,
		uid,
		request.Name,
	); err != nil {
		log.Println("error while inserting data", err.Error())
		return "", err
	}

	return uid.String(), nil

}

func (c *categoryRepo) GetByID(pKey models.PrimaryKey) (models.Category, error) {

	category := models.Category{}

	query := `select id, name from categories where id = $1`

	if err := c.pool.QueryRow(context.Background(), query, pKey.ID).Scan(
		&category.ID,
		&category.Name,
	); err != nil {
		log.Println("error while scanning category", err.Error())
		return models.Category{}, err
	}

	return category, nil
}

func (c *categoryRepo) GetList(request models.GetListRequest) (models.CategoriesResponse, error) {

	var (
		categories        = []models.Category{}
		count             = 0
		countQuery, query string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `
	SELECT count(1) from categories
	`

	if search != "" {
		countQuery += fmt.Sprintf(` and (name ilike '%%%s%%')`, search)

	}

	if err := c.pool.QueryRow(context.Background(), countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of categories", err.Error())
		return models.CategoriesResponse{}, err
	}

	query = `
	SELECT id, name
	FROM categories
	`

	if search != "" {
		query += fmt.Sprintf(` and (name ilike '%%%s%%')`, search)
	}

	query += ` LIMIT $1 OFFSET $2`

	rows, err := c.pool.Query(context.Background(), query, request.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.CategoriesResponse{}, err
	}

	for rows.Next() {
		category := models.Category{}

		if err = rows.Scan(
			&category.ID,
			&category.Name,
		); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.CategoriesResponse{}, err
		}

		categories = append(categories, category)

	}

	return models.CategoriesResponse{
		Categories: categories,
		Count:      count,
	}, nil
}

func (c *categoryRepo) Update(request models.Category) (string, error) {

	query := ` update categories set name = $1 where id = $2`

	if _, err := c.pool.Exec(context.Background(), query, request.Name, request.ID); err != nil {
		fmt.Println("error while updating category data", err.Error())
		return "", err
	}

	return request.ID, nil

}

func (c *categoryRepo) Delete(request models.PrimaryKey) error {

	query := `delete from 
	categories where
	id = $1
	`

	if _, err := c.pool.Exec(context.Background(), query, request.ID); err != nil {
		fmt.Println("error while deleting category by id", err.Error())
		return err
	}

	return nil
}
