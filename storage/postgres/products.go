package postgres

import (
	"basa/structs"
	"database/sql"

	"github.com/google/uuid"
)

type productsRepo struct {
	DB *sql.DB
}

func NewProductsRepo(db *sql.DB) productsRepo {
	return productsRepo{
		DB: db,
	}
}

// INSERT PRODUCT

func (p productsRepo) InsertProducts(s structs.Products) error {
	id := uuid.New()
	if _, err := p.DB.Exec(`insert into products values ($1, $2, $3)`, id, s.Price, s.ProductName); err != nil {
		return err
	}
	return nil
}

// GET PRODUCT BY ID

func (p productsRepo) GetByIDProduct(id uuid.UUID) (structs.Products, error) {
	product := structs.Products{}

	if err := p.DB.QueryRow(`select id, price, product_name from products where id = $1`, id).Scan(
		&product.ID,
		&product.Price,
		&product.ProductName,
	); err != nil {
		return structs.Products{}, err
	}
	return product, nil
}

// GET PRODUCTS LIST

func (p productsRepo) GetListProduct() ([]structs.Products, error) {
	Products := []structs.Products{}

	rows, err := p.DB.Query(`select * from products`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		Product := structs.Products{}

		rows.Scan(&Product.ID, &Product.Price, &Product.ProductName)

		Products = append(Products, Product)

	}
	return Products, nil
}

// UPDATE PRODUCT

func (p productsRepo) UpdateProducts(product structs.Products) error {

	_, err := p.DB.Exec(`update products set price = $1, product_name = $2 where id = $3`, product.Price, product.ProductName, product.ID)
	if err != nil {
		return err
	}
	return nil
}

// DELETE PRODUCT

func (p productsRepo) DeleteProducts(id uuid.UUID) error {
	_, err := p.DB.Exec(`delete from products where id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
