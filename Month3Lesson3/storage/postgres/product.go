package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"test/api/models"
	"test/storage"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

	query := `select id, name, price, original_price, quantity, category_id from products where id = $1`

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
			&product.OriginalPrice,
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

func (p *productRepo) Search(customerProductIDs map[string]int) (map[string]int, map[string]int, error) {

	var (
		selectedProducts = models.SellRequest{
			Products: map[string]int{},
		}
		products      = make([]string, len(customerProductIDs))
		productPrices = make(map[string]int, 0)
	)

	for key := range customerProductIDs {
		products = append(products, key)
	}

	query := `  
	select id, quantity, price, original_price from products where id::varchar = ANY($1)
	`
	rows, err := p.db.Query(query, pq.Array(products))
	if err != nil {
		fmt.Println("Error while getting products by product ids", err.Error())
		return nil, nil, err
	}

	for rows.Next() {
		var (
			quantity, price, originalPrice int
			productID                      string
		)
		if err = rows.Scan(
			&productID,
			&quantity,
			&price,
			&originalPrice,
		); err != nil {
			fmt.Println("error while scanning rows one by one", err.Error())
			return nil, nil, err
		}

		if customerProductIDs[productID] <= quantity {
			selectedProducts.Products[productID] = price
			productPrices[productID] = originalPrice
		}

	}

	return selectedProducts.Products, productPrices, nil
}

func (p *productRepo) TakeProduct(products map[string]int) error {
	query := `
	update products set quantity = quantity - $1 where id = $2
	`
	for productID, quantity := range products {
		if _, err := p.db.Exec(query, quantity, productID); err != nil {
			fmt.Println("error while updating product quantity", err.Error())
			return err
		}
	}

	return nil

}
