package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"test/api/models"
	"test/storage"

	"github.com/google/uuid"
)

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) storage.IProductStorage {
	return &productRepo{
		db: db,
	}
}

func (p *productRepo) Create(createProduct models.CreateProduct) (string, error) {
	uid := uuid.New()

	query := `insert into products values ($1, $2, $3, $4, $5, $6)`

	if _, err := p.db.Exec(query, uid, createProduct.Name, createProduct.Price, createProduct.OriginalPrice, createProduct.Quantity, createProduct.CategoryID); err != nil {
		log.Println("error while inserting data", err.Error())
		return "", err
	}

	return uid.String(), nil
}
func (p *productRepo) GetByID(pKey models.PrimaryKey) (models.Product, error) {

	product := models.Product{}

	query := `select id, name, price, original_price, quantity, catefory_id from products where id = $1`

	err := p.db.QueryRow(query, pKey.ID).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.OriginalPrice,
		&product.Quantity,
		&product.CategoryID,
	)
	if err != nil {
		log.Println("error while scanning product", err.Error())
		return models.Product{}, err
	}

	return product, nil

}

func (p *productRepo) GetList(request models.GetListRequest) (models.ProductsResponse, error) {

	var (
		products          = []models.Product{}
		count             = 0
		countQuery, query string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `
	SELECT count(1) from products
	`
	if search != "" {
		countQuery += fmt.Sprintf(` and (name ilike '%%%s%%')`, search)
	}

	if err := p.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of baskets", err.Error())
		return models.ProductsResponse{}, err
	}

	query = `
	SELECT id, name, price, original_price, quantity, category_id from products
	`
	if search != "" {
		query += fmt.Sprintf(` and (name ilike '%%%s%%')`, search)
	}

	query += `LIMIT $1 OFFSET $2`

	rows, err := p.db.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.ProductsResponse{}, err
	}

	for rows.Next() {
		product := models.Product{}

		if err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Quantity,
			&product.CategoryID,
		); err != nil {
			log.Println("error while scanning row", err.Error())
			return models.ProductsResponse{}, err
		}

		products = append(products, product)

	}

	return models.ProductsResponse{
		Products: products,
		Count:    count,
	}, nil
}

func (p *productRepo) Update(request models.UpdateProduct) (string, error) {

	query := `update products set name = $1, price = $2, original_price = $3, quantity = $4, category_id = $5 where id = $6`

	if _, err := p.db.Exec(query, request.Name, request.Price, request.OriginalPrice, request.Quantity, request.CategoryID, request.ID); err != nil {
		log.Println("error while updating product data", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (p *productRepo) Delete(request models.PrimaryKey) error {

	query := `delete from 
	products
	 where id = $1`

	if _, err := p.db.Exec(query, request.ID); err != nil {
		log.Println("error while deleting product by id", err.Error())
		return err
	}

	return nil
}
