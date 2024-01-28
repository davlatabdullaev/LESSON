package postgres

import (
	"database/sql"
	"fmt"
	"test/api/models"
	"test/storage"

	"github.com/google/uuid"
)

type basketProductRepo struct {
	db *sql.DB
}

func NewBasketProductRepo(db *sql.DB) storage.IBasketProductStorage {
	return basketProductRepo{db: db}
}

func (b basketProductRepo) Create(product models.CreateBasketProduct) (string, error) {
	id := uuid.New()
	query := `insert into basket_products(id, basket_id, product_id, quantity) 
					values($1, $2, $3, $4)`

	fmt.Println("id", id)
	if _, err := b.db.Exec(query,
		id,
		product.BasketID,
		product.ProductID,
		product.Quantity); err != nil {
		fmt.Println("error is while insert", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (b basketProductRepo) GetByID(key models.PrimaryKey) (models.BasketProduct, error) {
	product := models.BasketProduct{}
	query := `select id, basket_id, product_id, quantity from basket_products where id = $1`

	if err := b.db.QueryRow(query, key.ID).Scan(
		&product.ID,
		&product.BasketID,
		&product.ProductID,
		&product.Quantity); err != nil {
		fmt.Println("error is while selecting by id", err.Error())
		return models.BasketProduct{}, err
	}

	return product, nil
}

func (b basketProductRepo) GetList(request models.GetListRequest) (models.BasketProductResponse, error) {
	var (
		count             = 0
		basketProducts    = []models.BasketProduct{}
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from basket_products `
	if search != "" {
		countQuery += fmt.Sprintf(` where CAST(quantity AS TEXT) ilike '%%%s%%'`, search)
	}

	if err := b.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning count", err.Error())
		return models.BasketProductResponse{}, err
	}

	query = `select id, basket_id, product_id, quantity from basket_products `
	if search != "" {
		query += fmt.Sprintf(` where CAST(quantity AS TEXT) ilike '%%%s%%'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := b.db.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting basket products", err.Error())
		return models.BasketProductResponse{}, err
	}

	for rows.Next() {
		basketProd := models.BasketProduct{}
		if err = rows.Scan(&basketProd.ID, &basketProd.BasketID, &basketProd.ProductID, &basketProd.Quantity); err != nil {
			fmt.Println("error is while scanning basket products", err.Error())
			return models.BasketProductResponse{}, err
		}
		basketProducts = append(basketProducts, basketProd)
	}

	return models.BasketProductResponse{
		BasketProducts: basketProducts,
		Count:          count,
	}, err
}

func (b basketProductRepo) Update(product models.UpdateBasketProduct) (string, error) {
	query := `update basket_products set basket_id = $1, product_id = $2, quantity = $3 where id = $4`
	if _, err := b.db.Exec(query,
		&product.BasketID,
		&product.ProductID,
		&product.Quantity,
		&product.ID); err != nil {
		fmt.Println("error is while updating basket_products", err.Error())
		return "", err
	}

	return product.ID, nil
}

func (b basketProductRepo) Delete(key models.PrimaryKey) error {
	query := `delete from basket_products where id = $1`

	if _, err := b.db.Exec(query, key.ID); err != nil {
		fmt.Println("error is while deleting basket products", err.Error())
		return err
	}
	return nil
}

func (b basketProductRepo) AddProducts(basketID string, products map[string]int) error {
	query := ` 
	insert into basket_products
	(id, basket_id, product_id, quantity)
	values ($1, $2, $3, $4)
	`
	for productID, quantity := range products {
		if _, err := b.db.Exec(query, uuid.New(), basketID, productID, quantity); err != nil {
			fmt.Println("error while adding product to basketProduct table", err.Error())
			return err
		}
	}
	return nil
}