package main

// basket go

package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type basketRepo struct {
	pool *pgxpool.Pool
}

func NewBasketRepo(pool *pgxpool.Pool) storage.IBasketRepo {
	return &basketRepo{
		pool: pool,
	}
}

func (b *basketRepo) Create(ctx context.Context, basket models.CreateBasket) (string, error) {

	id := uuid.New()

	query := `insert into basket (id, sale_id, product_id, quantity, price) values ($1, $2, $3, $4, $5)`

	_, err := b.pool.Exec(ctx, query,
		id,
		basket.SaleID,
		basket.ProductID,
		basket.Quantity,
		basket.Price,
	)
	if err != nil {
		log.Println("error while inserting basket", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (b *basketRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Basket, error) {

	var updatedAt = sql.NullTime{}

	basket := models.Basket{}

	row := b.pool.QueryRow(ctx, `select 
	id,
    sale_id,
    product_id, 
	quantity, 
	price,
    created_at, 
	updated_at
	from basket where deleted_at is null and id = $1`, id.ID)

	err := row.Scan(
		&basket.ID,
		&basket.SaleID,
		&basket.ProductID,
		&basket.Quantity,
		&basket.Price,
		&basket.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting basket", err.Error())
		return models.Basket{}, err
	}

	if updatedAt.Valid {
		basket.UpdatedAt = updatedAt.Time
	}

	return basket, nil
}

func (b *basketRepo) GetList(ctx context.Context, request models.GetBasketsListRequest) (models.BasketsResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		baskets           = []models.Basket{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from basket where deleted_at is null `

	if search != "" {
		countQuery += fmt.Sprintf(` and product_id = '%s' or sale_id = '%s'`, search, search)
	}
	if err := b.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.BasketsResponse{}, err
	}

	query = `select 
	id, 
	sale_id, 
	product_id, 
	quantity, 
	price, 
	created_at, 
	updated_at
	from basket where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` and product_id = '%s' or sale_id = '%s'`, search, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := b.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting basket", err.Error())
		return models.BasketsResponse{}, err
	}

	for rows.Next() {
		basket := models.Basket{}
		if err = rows.Scan(
			&basket.ID,
			&basket.SaleID,
			&basket.ProductID,
			&basket.Quantity,
			&basket.Price,
			&basket.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning basket data", err.Error())
			return models.BasketsResponse{}, err
		}

		if updatedAt.Valid {
			basket.UpdatedAt = updatedAt.Time
		}

		baskets = append(baskets, basket)

	}

	return models.BasketsResponse{
		Baskets: baskets,
		Count:   count,
	}, nil
}

func (b *basketRepo) Update(ctx context.Context, request models.UpdateBasket) (string, error) {

	query := `update basket
   set sale_id = $1,
    product_id = $2, 
	quantity = $3,
	price = $4,
	updated_at = $5 
    where id = $6 
   `

	_, err := b.pool.Exec(ctx, query,
		request.SaleID,
		request.ProductID,
		request.Quantity,
		request.Price,
		time.Now(),
		request.ID)
	if err != nil {
		log.Println("error while updating basket data...", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (b *basketRepo) Delete(ctx context.Context, id string) error {

	query := `
	  update basket
	  set deleted_at = $1
	  where id = $2
	`

	_, err := b.pool.Exec(ctx,
		query,
		time.Now(),
		id)
	if err != nil {
		log.Println("error while deleting basket by id", err.Error())
		return err
	}
	return nil
}


func (b *basketRepo) UpdateBasketQuantity(ctx context.Context, request models.UpdateBasketQuantity) (string, error) {

	query := `update basket
    set 
    quantity = quantity + $1,
    updated_at = $2 
    where id = $3 
   `

	_, err := b.pool.Exec(ctx, query,
		request.Quantity,
		time.Now(),
		request.ID)
	if err != nil {
		log.Println("error while updating basket quantity...", err.Error())
		return "", err
	}

	return request.ID, nil
}

// branch go

package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type branchRepo struct {
	pool *pgxpool.Pool
}

func NewBranchRepo(pool *pgxpool.Pool) storage.IBranchRepo {
	return &branchRepo{
		pool: pool,
	}
}

func (b *branchRepo) Create(ctx context.Context, branch models.CreateBranch) (string, error) {

	id := uuid.New()

	query := `insert into branch (
		id, 
		name, 
		address) 
		values ($1, $2, $3)`

	_, err := b.pool.Exec(ctx, query,
		id,
		branch.Name,
		branch.Address,
	)
	if err != nil {
		log.Println("error while inserting branch", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (b *branchRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Branch, error) {

	var updatedAt = sql.NullTime{}

	branch := models.Branch{}

	query := `select
	 id, 
	 name, 
	 address, 
	 created_at, 
	 updated_at
	 from branch where deleted_at is null and id = $1`

	row := b.pool.QueryRow(ctx, query, id.ID)

	err := row.Scan(
		&branch.ID,
		&branch.Name,
		&branch.Address,
		&branch.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting branch", err.Error())
		return models.Branch{}, err
	}

	if updatedAt.Valid {
		branch.UpdatedAt = updatedAt.Time
	}

	return branch, nil

}

func (b *branchRepo) GetList(ctx context.Context, request models.GetListRequest) (models.BranchsResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		branchs           = []models.Branch{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from branch where deleted_at is null`

	if search != "" {
		countQuery += fmt.Sprintf(` and name ilike '%%%s%%' or address ilike '%%%s%%'`, search, search)
	}
	if err := b.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		log.Println("error is while selecting count", err.Error())
		return models.BranchsResponse{}, err
	}

	query = `select 
	id, 
	name, 
	address, 
	created_at, 
	updated_at 
	from branch where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` and name ilike '%%%s%%' or address ilike '%%%s%%'`, search, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := b.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		log.Println("error is while selecting branch", err.Error())
		return models.BranchsResponse{}, err
	}

	for rows.Next() {
		branch := models.Branch{}
		if err = rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.Address,
			&branch.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning branch data", err.Error())
			return models.BranchsResponse{}, err
		}

		if updatedAt.Valid {
			branch.UpdatedAt = updatedAt.Time
		}

		branchs = append(branchs, branch)

	}

	return models.BranchsResponse{
		Branchs: branchs,
		Count:   count,
	}, nil
}

func (b *branchRepo) Update(ctx context.Context, request models.UpdateBranch) (string, error) {

	query := `update branch
   set name = $1,
    address = $2, 
	updated_at = $3
   where id = $4  
   `

	_, err := b.pool.Exec(ctx, query,
		request.Name,
		request.Address,
		time.Now(),
		request.ID)
	if err != nil {
		log.Println("error while updating branch data...", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (b *branchRepo) Delete(ctx context.Context, id string) error {

	query := `
	update branch
	 set deleted_at = $1
	  where id = $2
	`

	_, err := b.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting branch by id", err.Error())
		return err
	}
	return nil
}

// category go

package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoryRepo struct {
	pool *pgxpool.Pool
}

func NewCategoryRepo(pool *pgxpool.Pool) storage.ICategoryRepo {
	return &categoryRepo{
		pool: pool,
	}
}

func (c *categoryRepo) Create(ctx context.Context, category models.CreateCategory) (string, error) {

	id := uuid.New()

	query := `insert into category (id, name, parent_id) values ($1, $2, $3)`

	_, err := c.pool.Exec(ctx, query,
		id,
		category.Name,
		category.ParentID,
	)
	if err != nil {
		log.Println("error while inserting category", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (c *categoryRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Category, error) {

	var updatedAt = sql.NullTime{}

	category := models.Category{}

	row := c.pool.QueryRow(ctx, `select 
	id, 
	name, 
	parent_id, 
	created_at, 
	updated_at from category where deleted_at is null and id = $1`, id.ID)

	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.ParentID,
		&category.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting category", err.Error())
		return models.Category{}, err
	}

	if updatedAt.Valid {
		category.UpdatedAt = updatedAt.Time
	}

	return category, nil
}

func (c *categoryRepo) GetList(ctx context.Context, request models.GetListRequest) (models.CategoriesResponse, error) {
	var (
		updatedAt         = sql.NullTime{}
		categories        = []models.Category{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from category where deleted_at is null`

	if search != "" {
		countQuery += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}
	if err := c.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.CategoriesResponse{}, err
	}

	query = `select id, name, parent_id, created_at, updated_at from category where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := c.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting category", err.Error())
		return models.CategoriesResponse{}, err
	}

	for rows.Next() {
		category := models.Category{}
		var parentID sql.NullString
		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&parentID,
			&category.CreatedAt,
			&updatedAt); err != nil {
			fmt.Println("error is while scanning category data", err.Error())
			return models.CategoriesResponse{}, err
		}

		if parentID.Valid {
			category.ParentID = parentID.String
		} else {
			category.ParentID = ""
		}

		if updatedAt.Valid {
			category.UpdatedAt = updatedAt.Time
		}

		categories = append(categories, category)

	}

	return models.CategoriesResponse{
		Categories: categories,
		Count:      count,
	}, nil
}

func (c *categoryRepo) Update(ctx context.Context, request models.UpdateCategory) (string, error) {

	query := `update category
   set name = $1, parent_id = $2, updated_at = $3 
   where id = $4  
   `

	_, err := c.pool.Exec(ctx, query,
		request.Name,
		request.ParentID,
		time.Now(),
		request.ID)
	if err != nil {
		log.Println("error while updating category data...", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (c *categoryRepo) Delete(ctx context.Context, id string) error {

	query := `
	update category
	 set deleted_at = $1
	  where id = $2
	`

	_, err := c.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting category by id", err.Error())
		return err
	}
	return nil
}

//product go

package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) storage.IProductRepo {
	return &productRepo{
		pool: pool,
	}
}

func (p *productRepo) Create(ctx context.Context, product models.CreateProduct) (string, error) {

	id := uuid.New()

	query := `insert into product (id, name, price, category_id) values ($1, $2, $3, $4)`

	_, err := p.pool.Exec(ctx, query,
		id,
		product.Name,
		product.Price,
		product.CategoryID,
	)
	if err != nil {
		log.Println("error while inserting product", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (p *productRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Product, error) {

	var updatedAt = sql.NullTime{}

	product := models.Product{}

	row := p.pool.QueryRow(ctx, `select
	 id, 
	 name, 
	 price, 
	 barcode, 
	 category_id, 
	 created_at, 
	 updated_at  from product where deleted_at is null and id = $1`, id.ID)

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Barcode,
		&product.CategoryID,
		&product.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting product", err.Error())
		return models.Product{}, err
	}

	if updatedAt.Valid {
		product.UpdatedAt = updatedAt.Time
	}

	return product, nil
}

func (p *productRepo) GetList(ctx context.Context, request models.ProductGetListRequest) (models.ProductsResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		products          = []models.Product{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from product where deleted_at is null `

	if search != "" {
		countQuery += fmt.Sprintf(`and name ilike '%%%s%%' or barcode ilike '%%%s%%'`, search, search)
	}
	if err := p.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.ProductsResponse{}, err
	}

	query = `select 
	id, 
	name, 
	price, 
	barcode, 
	category_id, 
	created_at, 
	updated_at from product where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` and name ilike '%%%s%%' or barcode ilike '%%%s%%'`, search, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := p.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting product", err.Error())
		return models.ProductsResponse{}, err
	}

	for rows.Next() {
		product := models.Product{}
		if err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Barcode,
			&product.CategoryID,
			&product.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning product data", err.Error())
			return models.ProductsResponse{}, err
		}

		if updatedAt.Valid {
			product.UpdatedAt = updatedAt.Time
		}

		products = append(products, product)

	}

	return models.ProductsResponse{
		Products: products,
		Count:    count,
	}, nil
}

func (p *productRepo) Update(ctx context.Context, request models.UpdateProduct) (string, error) {

	query := `update product
   set 
    name = $1,
    price = $2, 
	category_id = $3, 
	updated_at = $4
   where id = $5  
   `

	_, err := p.pool.Exec(ctx, query,
		request.Name,
		request.Price,
		request.CategoryID,
		time.Now(),
		request.ID)
	if err != nil {
		log.Println("error while updating product data...", err.Error())
		return "", err
	}
	return request.ID, nil
}

func (p *productRepo) Delete(ctx context.Context, id string) error {

	query := `
	update product
	 set deleted_at = $1
	  where id = $2
	`

	_, err := p.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting product by id", err.Error())
		return err
	}

	return nil
}

// sale go

package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type saleRepo struct {
	pool *pgxpool.Pool
}

func NewSaleRepo(pool *pgxpool.Pool) storage.ISaleRepo {
	return &saleRepo{
		pool: pool,
	}
}

func (s *saleRepo) Create(ctx context.Context, sale models.CreateSale) (string, error) {

	id := uuid.New()

	query := `insert into sale (id, branch_id, shop_assistent_id, cashier_id, payment_type, price, status, client_name) values ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := s.pool.Exec(ctx, query,
		id,
		sale.BranchID,
		sale.ShopAssistantID,
		sale.CashierID,
		sale.PaymentType,
		sale.Price,
		sale.Status,
		sale.ClientName,
	)
	if err != nil {
		log.Println("error while inserting sale", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (s *saleRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Sale, error) {

	var updatedAt = sql.NullTime{}

	sale := models.Sale{}

	row := s.pool.QueryRow(ctx, `select id, branch_id, shop_assistent_id, cashier_id, payment_type, price, status, client_name, created_at, updated_at  from sale where deleted_at is null and id = $1`, id.ID)

	err := row.Scan(
		&sale.ID,
		&sale.BranchID,
		&sale.ShopAssistantID,
		&sale.CashierID,
		&sale.PaymentType,
		&sale.Price,
		&sale.Status,
		&sale.ClientName,
		&sale.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting sale", err.Error())
		return models.Sale{}, err
	}

	if updatedAt.Valid {
		sale.UpdatedAt = updatedAt.Time
	}

	return sale, nil
}

func (s *saleRepo) GetList(ctx context.Context, request models.GetListRequest) (models.SalesResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		sales             = []models.Sale{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from sale where deleted_at is null`

	if search != "" {
		countQuery += fmt.Sprintf(` and status ilike '%%%s%%' or payment_type ilike '%%%s%%'`, search, search)
	}
	if err := s.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.SalesResponse{}, err
	}

	query = `select 
	id, 
	branch_id, 
	shop_assistent_id, 
	cashier_id, 
	payment_type, 
	price, 
	status, 
	client_name, 
	created_at, 
	updated_at from sale where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` and status ilike '%%%s%%' or payment_type ilike '%%%s%%'`, search, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := s.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting product", err.Error())
		return models.SalesResponse{}, err
	}

	for rows.Next() {
		sale := models.Sale{}
		if err = rows.Scan(
			&sale.ID,
			&sale.BranchID,
			&sale.ShopAssistantID,
			&sale.CashierID,
			&sale.PaymentType,
			&sale.Price,
			&sale.Status,
			&sale.ClientName,
			&sale.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning sale data", err.Error())
			return models.SalesResponse{}, err
		}

		if updatedAt.Valid {
			sale.UpdatedAt = updatedAt.Time
		}

		sales = append(sales, sale)

	}

	return models.SalesResponse{
		Sales: sales,
		Count: count,
	}, nil
}

func (s *saleRepo) Update(ctx context.Context, request models.UpdateSale) (string, error) {

	query := `update sale set 
	branch_id = $1, 
	shop_assistent_id = $2,
	cashier_id = $3, 
	payment_type = $4, 
	price = $5, 
	status = $6, 
	client_name = $7, 
	updated_at = $8 
	 where id = $9`

	_, err := s.pool.Exec(ctx, query,
		request.BranchID,
		request.ShopAssistantID,
		request.CashierID,
		request.PaymentType,
		request.Price,
		request.Status,
		request.ClientName,
		time.Now(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating sale data...", err.Error())
		return "", err
	}
	return request.ID, nil
}

func (s *saleRepo) Delete(ctx context.Context, id string) error {

	query := `update sale 
	set deleted_at = $1 
	where id = $2`

	_, err := s.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting sale by id", err.Error())
		return err
	}

	return nil
}

func (s *saleRepo) UpdateSalePrice(ctx context.Context, request models.SaleRequest) (string, error) {

	query := `update sale set 
  price = $1,
  status = $2,
  updated_at = $3 
  where id = $4 `

	if rowsAffected, err := s.pool.Exec(ctx, query, request.TotalPrice, request.Status, time.Now(), request.ID); err != nil {
		if r := rowsAffected.RowsAffected(); r == 0 {
			log.Println("error in rows affected ", err.Error())
			return "", err
		}
		log.Println("error while updating sale price and status...", err.Error())
		return "", err
	}
	return request.ID, nil

}

