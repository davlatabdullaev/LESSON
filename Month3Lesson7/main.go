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

func (s *saleRepo)   UpdateSalePrice(ctx context.Context, request models.SaleRequest) (string, error) {

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

// staff go

package postgres

import (
	"bazaar/api/models"
	"bazaar/pkg/check"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type staffRepo struct {
	pool *pgxpool.Pool
}

func NewStaffRepo(pool *pgxpool.Pool) storage.IStaffRepo {
	return &staffRepo{
		pool: pool,
	}
}

func (s *staffRepo) Create(ctx context.Context, request models.CreateStaff) (string, error) {

	id := uuid.New()

	query := `insert into staff (
		id, 
		branch_id, 
		tarif_id, 
		type_staff, 
		name, 
		balance, 
		birth_date, 
		age, 
		gender, 
		login, 
		password) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := s.pool.Exec(ctx, query,
		id,
		request.BranchID,
		request.TarifID,
		request.TypeStaff,
		request.Name,
		request.Balance,
		request.BirthDate,
		check.CalculateAge(request.BirthDate),
		request.Gender,
		request.Login,
		request.Password,
	)
	if err != nil {
		log.Println("error while inserting staff data", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (s *staffRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Staff, error) {

	var updatedAt = sql.NullTime{}

	staff := models.Staff{}

	query := `select 
	id, 
	branch_id, 
	tarif_id, 
	type_staff, 
	name, 
	birth_date::text, 
	age, 
	gender, 
	login, 
	balance,
	password, 
	created_at, 
	updated_at from staff where deleted_at is null and id = $1`

	row := s.pool.QueryRow(ctx, query, id.ID)

	err := row.Scan(
		&staff.ID,
		&staff.BranchID,
		&staff.TarifID,
		&staff.TypeStaff,
		&staff.Name,
		&staff.BirthDate,
		&staff.Age,
		&staff.Gender,
		&staff.Login,
		&staff.Balance,
		&staff.Password,
		&staff.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting staff data", err.Error())
		return models.Staff{}, err
	}

	if updatedAt.Valid {
		staff.UpdatedAt = updatedAt.Time
	}

	return staff, nil
}

func (s *staffRepo) GetList(ctx context.Context, request models.GetListRequest) (models.StaffsResponse, error) {
	var (
		updatedAt         = sql.NullTime{}
		staffs            = []models.Staff{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from staff where deleted_at is null `

	if search != "" {
		countQuery += fmt.Sprintf(`and name ilike '%%%s%%' or login ilike '%%%s%%' `, search, search)
	}
	if err := s.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting staff count", err.Error())
		return models.StaffsResponse{}, err
	}

	query = `select 
	id, 
	branch_id, 
	tarif_id, 
	type_staff, 
	name, 
	birth_date::text, 
	age, 
	gender, 
	login,
	balance, 
	password, 
	created_at, 
	updated_at from staff where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` and name ilike '%%%s%%' or login ilike '%%%s%%' `, search, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := s.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting staff", err.Error())
		return models.StaffsResponse{}, err
	}

	for rows.Next() {
		staff := models.Staff{}
		if err = rows.Scan(
			&staff.ID,
			&staff.BranchID,
			&staff.TarifID,
			&staff.TypeStaff,
			&staff.Name,
			&staff.BirthDate,
			&staff.Age,
			&staff.Gender,
			&staff.Login,
			&staff.Balance,
			&staff.Password,
			&staff.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning staff data", err.Error())
			return models.StaffsResponse{}, err
		}

		if updatedAt.Valid {
			staff.UpdatedAt = updatedAt.Time
		}

		staffs = append(staffs, staff)

	}

	return models.StaffsResponse{
		Staffs: staffs,
		Count:  count,
	}, nil
}

func (s *staffRepo) Update(ctx context.Context, request models.UpdateStaff) (string, error) {

	query := `update staff
   set 
   branch_id = $1, 
   tarif_id = $2, 
   type_staff = $3,
   name = $4, 
   birth_date = $5, 
   age = $6, 
   gender = $7, 
   login = $8, 
   password = $9,
   balance = $10, 
   updated_at = $11
   where id = $12
   `
	_, err := s.pool.Exec(ctx, query,
		request.BranchID,
		request.TarifID,
		request.TypeStaff,
		request.Name,
		request.BirthDate,
		check.CalculateAge(request.BirthDate),
		request.Gender,
		request.Login,
		request.Password,
		request.Balance,
		time.Now(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating staff data...", err.Error())
		return "", err
	}
	return request.ID, nil
}

func (s *staffRepo) Delete(ctx context.Context, id string) error {

	query := `
	update staff
	 set deleted_at = $1
	  where id = $2
	`

	_, err := s.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting staff by id", err.Error())
		return err
	}
	return nil
}

func (s *staffRepo) UpdateStaffBalance(ctx context.Context, request models.UpdateStaffBalance) error {

	fmt.Println("pg 123: ", request)

	query := `update staff
	 set 
	 balance = balance + $1,
	 updated_at = $2
	 where id = $3
	 `

	_, err := s.pool.Exec(ctx, query, request.Balance, time.Now(), request.ID)

	if err != nil {
		log.Println("error while updating staff balance", err.Error())
		return err
	}

	return nil

}
// storage transaction go

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

type storageTransactionRepo struct {
	pool *pgxpool.Pool
}

func NewStorageTransactionRepo(pool *pgxpool.Pool) storage.IStorageTransactionRepo {
	return &storageTransactionRepo{
		pool: pool,
	}
}

func (s *storageTransactionRepo) Create(ctx context.Context, request models.CreateStorageTransaction) (string, error) {

	id := uuid.New()

	query := `insert into storage_transaction (id, staff_id, product_id, storage_transaction_type, price, quantity) values ($1, $2, $3, $4, $5, $6)`

	_, err := s.pool.Exec(ctx, query,
		id,
		request.StaffID,
		request.ProductID,
		request.StorageTransactionType,
		request.Price,
		request.Quantity,
	)
	if err != nil {
		log.Println("error while inserting storage transaction data", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (s *storageTransactionRepo) Get(ctx context.Context, id models.PrimaryKey) (models.StorageTransaction, error) {

	var updatedAt = sql.NullTime{}

	storageTransaction := models.StorageTransaction{}

	query := `select id,
	 staff_id, 
	 product_id, 
	 storage_transaction_type, 
	 price, 
	 quantity, 
	 created_at, 
	 updated_at 
	 from storage_transaction
	 where deleted_at is null and id = $1`

	row := s.pool.QueryRow(ctx, query, id.ID)

	err := row.Scan(
		&storageTransaction.ID,
		&storageTransaction.StaffID,
		&storageTransaction.ProductID,
		&storageTransaction.StorageTransactionType,
		&storageTransaction.Price,
		&storageTransaction.Quantity,
		&storageTransaction.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting storage transaction data", err.Error())
		return models.StorageTransaction{}, err
	}

	if updatedAt.Valid {
		storageTransaction.UpdatedAt = updatedAt.Time
	}

	return storageTransaction, nil
}

func (s *storageTransactionRepo) GetList(ctx context.Context, request models.GetListRequest) (models.StorageTransactionsResponse, error) {

	var (
		updatedAt           = sql.NullTime{}
		storageTransactions = []models.StorageTransaction{}
		count               = 0
		query, countQuery   string
		page                = request.Page
		offset              = (page - 1) * request.Limit
		search              = request.Search
	)

	countQuery = `select count(1) from storage_transaction where deleted_at is null `

	if search != "" {
		countQuery += fmt.Sprintf(`and storage_transaction_type ilike '%%%s%%'`, search)
	}
	if err := s.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting storage_transaction count", err.Error())
		return models.StorageTransactionsResponse{}, err
	}

	query = `select 
	id, 
	staff_id, 
	product_id, 
	storage_transaction_type, 
	price, 
	quantity, 
	created_at, 
	updated_at
	from storage_transaction 
	where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where storage_transaction_type ilike '%%%s%%'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := s.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting storage transaction", err.Error())
		return models.StorageTransactionsResponse{}, err
	}

	for rows.Next() {
		storageTransaction := models.StorageTransaction{}
		if err = rows.Scan(
			&storageTransaction.ID,
			&storageTransaction.StaffID,
			&storageTransaction.ProductID,
			&storageTransaction.StorageTransactionType,
			&storageTransaction.Price,
			&storageTransaction.Quantity,
			&storageTransaction.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning storage transaction data", err.Error())
			return models.StorageTransactionsResponse{}, err
		}

		if updatedAt.Valid {
			storageTransaction.UpdatedAt = updatedAt.Time
		}

		storageTransactions = append(storageTransactions, storageTransaction)

	}

	return models.StorageTransactionsResponse{
		StorageTransactions: storageTransactions,
		Count:               count,
	}, nil
}

func (s *storageTransactionRepo) Update(ctx context.Context, request models.UpdateStorageTransaction) (string, error) {

	query := `update storage_transaction
   set 
   staff_id = $1, 
   product_id = $2, 
   storage_transaction_type = $3,
   price = $4, 
   quantity = $5, 
   updated_at = $6
   where id = $7
   `
	_, err := s.pool.Exec(ctx, query,
		request.StaffID,
		request.ProductID,
		request.StorageTransactionType,
		request.Price,
		request.Quantity,
		time.Now(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating storage_transaction data...", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (s *storageTransactionRepo) Delete(ctx context.Context, id string) error {

	query := `
	update storage_transaction
	 set deleted_at = $1
	  where id = $2
	`

	_, err := s.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting storage_transaction by id", err.Error())
		return err
	}

	return nil
}

// storage go

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

type storageRepo struct {
	pool *pgxpool.Pool
}

func NewStorageRepo(pool *pgxpool.Pool) storage.IStorageRepo {
	return &storageRepo{
		pool: pool,
	}
}

func (s *storageRepo) Create(ctx context.Context, storage models.CreateStorage) (string, error) {

	id := uuid.New()

	query := `insert into storage (id, product_id, branch_id, count) values ($1, $2, $3, $4)`

	_, err := s.pool.Exec(ctx, query,
		id,
		storage.ProductID,
		storage.BranchID,
		storage.Count,
	)
	if err != nil {
		log.Println("error while inserting storage", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (s *storageRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Storage, error) {

	var updatedAt = sql.NullTime{}

	storage := models.Storage{}

	row := s.pool.QueryRow(ctx, `select 
	id, 
	product_id, 
	branch_id, 
	count, 
	created_at, 
	updated_at  from storage where deleted_at is null and id = $1`, id.ID)

	err := row.Scan(
		&storage.ID,
		&storage.ProductID,
		&storage.BranchID,
		&storage.Count,
		&storage.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting storage", err.Error())
		return models.Storage{}, err
	}

	if updatedAt.Valid {
		storage.UpdatedAt = updatedAt.Time
	}

	return storage, nil
}

func (s *storageRepo) GetList(ctx context.Context, request models.GetListRequest) (models.StoragesResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		storages          = []models.Storage{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from storage where deleted_at is null`

	if search != "" {
		countQuery += fmt.Sprintf(` and product_id = '%s'`, search)
	}
	if err := s.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.StoragesResponse{}, err
	}

	query = `select 
	id, 
	product_id, 
	branch_id, 
	count, 
	created_at, 
	updated_at from storage where deleted_at is null`

	if search != "" {
		countQuery += fmt.Sprintf(` and product_id = '%s'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := s.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting product", err.Error())
		return models.StoragesResponse{}, err
	}

	for rows.Next() {
		storage := models.Storage{}
		if err = rows.Scan(
			&storage.ID,
			&storage.ProductID,
			&storage.BranchID,
			&storage.Count,
			&storage.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning storage data", err.Error())
			return models.StoragesResponse{}, err
		}

		if updatedAt.Valid {
			storage.UpdatedAt = updatedAt.Time
		}

		storages = append(storages, storage)

	}

	return models.StoragesResponse{
		Storages: storages,
		Count:    count,
	}, nil
}

func (s *storageRepo) Update(ctx context.Context, request models.UpdateStorage) (string, error) {

	query := ` update storage set 
	product_id = $1, 
	branch_id = $2, 
	count = $3, 
	updated_at = $4 
	where id = $5`

	_, err := s.pool.Exec(ctx, query,
		request.ProductID,
		request.BranchID,
		request.Count,
		time.Now(),
		request.ID,
	)

	if err != nil {
		log.Println("error while updating storage data...", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (s *storageRepo) Delete(ctx context.Context, id string) error {

	query := `update storage 
	set deleted_at = $1 
	where id = $2`

	_, err := s.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting storage by id", err.Error())
		return err
	}

	return nil
}

func (s *storageRepo) UpdateCount(ctx context.Context, request models.UpdateCount) error {

	query := `update storage set 
	count = count + $1,
	updated_at = $2
	where id = $3
	`

	_, err := s.pool.Exec(ctx, query, request.Count, time.Now(), request.ID)
	if err != nil {
		log.Println("error while update storage count", err.Error())
		return err
	}

	return nil

}

// tarif go

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

type tarifRepo struct {
	pool *pgxpool.Pool
}

func NewTarifRepo(pool *pgxpool.Pool) storage.ITarifRepo {
	return &tarifRepo{
		pool: pool,
	}
}

func (t *tarifRepo) Create(ctx context.Context, request models.CreateTarif) (string, error) {

	id := uuid.New()

	query := `insert into tarif (id, name, tarif_type, amount_for_cash,
		amount_for_card) 
	values 
	($1, $2, $3, $4, $5)`

	_, err := t.pool.Exec(ctx, query,
		id,
		request.Name,
		request.TarifType,
		request.AmountForCash,
		request.AmountForCard,
	)
	if err != nil {
		log.Println("error while inserting tarif data", err.Error())
		return "", err
	}

	return id.String(), nil

}

func (t *tarifRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Tarif, error) {

	var updatedAt = sql.NullTime{}

	tarif := models.Tarif{}

	query := `select 
	id, 
	name, 
	tarif_type, 
	amount_for_cash,
	amount_for_card, 
	created_at, 
	updated_at from tarif
	 where deleted_at is null and id = $1`

	row := t.pool.QueryRow(ctx, query, id.ID)

	err := row.Scan(
		&tarif.ID,
		&tarif.Name,
		&tarif.TarifType,
		&tarif.AmountForCash,
		&tarif.AmountForCard,
		&tarif.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting tarif data", err.Error())
		return models.Tarif{}, err
	}

	if updatedAt.Valid {
		tarif.UpdatedAt = updatedAt.Time
	}

	return tarif, nil

}

func (t *tarifRepo) GetList(ctx context.Context, request models.GetListRequest) (models.TarifsResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		tarifs            = []models.Tarif{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from tarif where deleted_at is null `

	if search != "" {
		countQuery += fmt.Sprintf(`and name ilike '%%%s%%'`, search)
	}
	if err := t.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting tarif count", err.Error())
		return models.TarifsResponse{}, err
	}

	query = `select 
	id, 
	name, 
	tarif_type, 
	amount_for_cash, 
	amount_for_card,
	created_at, 
	updated_at
	from tarif 
	where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where name ilike '%%%s%%'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := t.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting tarif", err.Error())
		return models.TarifsResponse{}, err
	}

	for rows.Next() {
		tarif := models.Tarif{}
		if err = rows.Scan(
			&tarif.ID,
			&tarif.Name,
			&tarif.TarifType,
			&tarif.AmountForCash,
			&tarif.AmountForCard,
			&tarif.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning tarif data", err.Error())
			return models.TarifsResponse{}, err
		}

		if updatedAt.Valid {
			tarif.UpdatedAt = updatedAt.Time
		}

		tarifs = append(tarifs, tarif)

	}

	return models.TarifsResponse{
		Tarifs: tarifs,
		Count:  count,
	}, nil
}

func (t *tarifRepo) Update(ctx context.Context, request models.UpdateTarif) (string, error) {

	query := `update tarif
   set name = $1, tarif_type = $2, amount_for_cash = $3,
   amount_for_card = $4, updated_at = $5
   where id = $6
   `
	_, err := t.pool.Exec(ctx, query,
		request.Name,
		request.TarifType,
		request.AmountForCash,
		request.AmountForCard,
		time.Now(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating tarif data...", err.Error())
		return "", err
	}
	return request.ID, nil
}

func (t *tarifRepo) Delete(ctx context.Context, id string) error {

	query := `
	update tarif
	 set deleted_at = $1
	  where id = $2
	`

	_, err := t.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting tarif by id", err.Error())
		return err
	}

	return nil
}

// transaction go

package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionRepo struct {
	pool *pgxpool.Pool
}

func NewTransactionRepo(pool *pgxpool.Pool) storage.ITransactionRepo {
	return &transactionRepo{
		pool: pool,
	}
}

func (t *transactionRepo) Create(ctx context.Context, request models.CreateTransactions) (string, error) {

	id := uuid.New()

	query := `insert into transactions (
		id, 
		sale_id, 
		staff_id, 
		transaction_type,
		source_type, 
		amount, 
		description) 
	values 
	($1, $2, $3, $4, $5, $6, $7)`

	_, err := t.pool.Exec(ctx, query,
		id,
		request.SaleID,
		request.StaffID,
		request.TransactionType,
		request.SourceType,
		request.Amount,
		request.Description,
	)
	if err != nil {
		log.Println("error while inserting transaction data", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (t *transactionRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Transactions, error) {

	var updatedAt = sql.NullTime{}

	transaction := models.Transactions{}

	query := `select 
	id, 
	sale_id, 
	staff_id, 
	transaction_type,
	source_type, 
	amount, 
	description, 
	created_at, 
	updated_at 
	from transactions
	 where deleted_at is null and id = $1`

	row := t.pool.QueryRow(ctx, query, id.ID)

	err := row.Scan(
		&transaction.ID,
		&transaction.SaleID,
		&transaction.StaffID,
		&transaction.TransactionType,
		&transaction.SourceType,
		&transaction.Amount,
		&transaction.Description,
		&transaction.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting transaction data", err.Error())
		return models.Transactions{}, err
	}

	if updatedAt.Valid {
		transaction.UpdatedAt = updatedAt.Time
	}

	return transaction, nil
}

func (t *transactionRepo) GetList(ctx context.Context, request models.GetListTransactionsRequest) (models.TransactionsResponse, error) {
	var (
		updatedAt         = sql.NullTime{}
		page              = request.Page
		offset            = (page - 1) * request.Limit
		transactions      = []models.Transactions{}
		fromAmount        = request.FromAmount
		toAmount          = request.ToAmount
		count             = 0
		query, countQuery string
	)

	countQuery = `select count(1) from transactions where deleted_at is null `
	if fromAmount != 0 && toAmount != 0 {
		countQuery += fmt.Sprintf(` and amount between %f and %f `, fromAmount, toAmount)
	} else if fromAmount != 0 && toAmount == 0 {
		countQuery += ` and amount >= ` + strconv.FormatFloat(fromAmount, 'f', 2, 64)
	} else if toAmount != 0 && fromAmount == 0 {
		countQuery += ` and amount <= ` + strconv.FormatFloat(toAmount, 'f', 2, 64)

	}
	if err := t.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning row", err.Error())
		return models.TransactionsResponse{}, err
	}

	query = `select 
	id, 
	sale_id, 
	staff_id, 
	transaction_type, 
	source_type, 
	amount,
    description, 
	created_at, 
	updated_at from transactions where deleted_at is null `

	if fromAmount != 0 && toAmount != 0 {
		query += fmt.Sprintf(` and amount between %f and %f `, fromAmount, toAmount)
	} else if fromAmount != 0 && toAmount == 0 {
		query += ` and amount >= ` + strconv.FormatFloat(fromAmount, 'f', 2, 64)
	} else if toAmount != 0 && fromAmount == 0 {
		query += ` and amount <= ` + strconv.FormatFloat(toAmount, 'f', 2, 64)

	}

	query += ` LIMIT $1 OFFSET $2 `

	rows, err := t.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting all transaction", err.Error())
		return models.TransactionsResponse{}, err
	}

	for rows.Next() {
		transaction := models.Transactions{}
		if err = rows.Scan(
			&transaction.ID,
			&transaction.SaleID,
			&transaction.StaffID,
			&transaction.TransactionType,
			&transaction.SourceType,
			&transaction.Amount,
			&transaction.Description,
			&transaction.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning rows", err.Error())
			return models.TransactionsResponse{}, err
		}

		if updatedAt.Valid {
			transaction.UpdatedAt = updatedAt.Time
		}

		transactions = append(transactions, transaction)
	}
	return models.TransactionsResponse{
		Transactions: transactions,
		Count:        count,
	}, nil
}

func (t *transactionRepo) Update(ctx context.Context, request models.UpdateTransactions) (string, error) {

	query := `update transactions
   set 
   sale_id = $1, 
   staff_id = $2, 
   transaction_type = $3,
   source_type = $4, 
   amount = $5, 
   description = $6, 
   updated_at = $7
   where id = $8
   `
	_, err := t.pool.Exec(ctx, query,
		request.SaleID,
		request.StaffID,
		request.TransactionType,
		request.SourceType,
		request.Amount,
		request.Description,
		time.Now(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating transaction data...", err.Error())
		return "", err
	}
	return request.ID, nil
}

func (t *transactionRepo) Delete(ctx context.Context, id string) error {

	query := `
	update transactions
	 set deleted_at = $1
	  where id = $2
	`

	_, err := t.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting transaction by id", err.Error())
		return err
	}

	return nil
}

func (t *transactionRepo) UpdateStaffBalanceAndCreateTransaction(ctx context.Context, request models.UpdateStaffBalanceAndCreateTransaction) error {

	transaction, err := t.pool.Begin(ctx)

	defer func() {

		if err != nil {
			transaction.Rollback(ctx)
		} else {
			transaction.Commit(ctx)
		}

	}()

	queryForUpdateStaffBalance := `update staff set
	 balance = balance + $1,
    updated_at = $2
	  where id = $3`

	_, err = transaction.Exec(ctx, queryForUpdateStaffBalance, request.UpdateCashierBalance.Amount, time.Now(), request.UpdateCashierBalance.StaffID)
	if err != nil {
		log.Println("error while update staff balance")
		return err
	}

	queryForCreateTransaction := `insert into transactions (
		id, 
		sale_id, 
		staff_id, 
		transaction_type,
		source_type, 
		amount, 
		description) 
	values 
	($1, $2, $3, $4, $5, $6, $7)`

	_, err = transaction.Exec(ctx, queryForCreateTransaction,
		uuid.New().String(),
		request.SaleID,
		request.UpdateCashierBalance.StaffID,
		request.TransactionType,
		request.SourceType,
		request.Amount,
		request.Description,
	)
	if err != nil {
		log.Println("error while creating transaction data", err.Error())
		return err
	}

	if request.UpdateShopAssistantBalance.StaffID != "" {

		queryForUpdateStaffBalance := `update staff set
	 balance = balance + $1,
     updated_at = $2
	 where id = $3`

		_, err = transaction.Exec(ctx, queryForUpdateStaffBalance, request.UpdateShopAssistantBalance.Amount, time.Now(), request.UpdateShopAssistantBalance.StaffID)
		if err != nil {
			log.Println("error while update staff balance")
			return err
		}

		queryForCreateTransaction := `insert into transactions (
	   id, 
	   sale_id, 
	   staff_id, 
	   transaction_type,
	   source_type, 
	   amount, 
	   description) 
   values 
   ($1, $2, $3, $4, $5, $6, $7)`

		_, err = transaction.Exec(ctx, queryForCreateTransaction,
			uuid.New().String(),
			request.SaleID,
			request.UpdateShopAssistantBalance.StaffID,
			request.TransactionType,
			request.SourceType,
			request.Amount,
			request.Description,
		)
		if err != nil {
			log.Println("error while creating transaction data", err.Error())
			return err
		}

	}

	return nil
}
