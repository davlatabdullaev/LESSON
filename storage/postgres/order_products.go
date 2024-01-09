package postgres

import (
	"basa/structs"
	"database/sql"

	"github.com/google/uuid"
)

type orderProductsRepo struct {
	DB *sql.DB
}

func NewOrderProductRepo(db *sql.DB) orderProductsRepo {
	return orderProductsRepo{
		DB: db,
	}
}


// INSERT ORDER PRODUCT

func (o orderProductsRepo) InsertOrderProduct(orderProduct structs.OrderProducts) error {
	id := uuid.New()
	if _, err := o.DB.Exec(`insert into order_products values ($1, $2, $3, $4, $5)`, id, orderProduct.OrderID, orderProduct.ProductID, orderProduct.Quantity, orderProduct.Price ); err != nil {
		return err
	}
	return nil
}

// GET ORDER PRODUCT BY ID

func (o orderProductsRepo) GetByIDOrderProduct(id uuid.UUID) (structs.OrderProducts, error) {
	orderProduct := structs.OrderProducts{}

	if err := o.DB.QueryRow(`select id, order_id, product_id, quantity, price from order_products where id = $1`, id).Scan(
		&orderProduct.Id,
		&orderProduct.OrderID,
		&orderProduct.ProductID,
		&orderProduct.Quantity,
		&orderProduct.Price,
	); err != nil {
		return structs.OrderProducts{}, err
	}
	return orderProduct, nil
}

// GET ORDER PRODUCTS LIST

func (o orderProductsRepo) GetListOrderProducts() ([]structs.OrderProducts, error) {
	OrderProducts := []structs.OrderProducts{}

	rows, err := o.DB.Query(`select * from order_products`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		OrderProduct := structs.OrderProducts{}

		rows.Scan(&OrderProduct.Id, &OrderProduct.OrderID, &OrderProduct.ProductID, &OrderProduct.Quantity, &OrderProduct.Price)

		OrderProducts = append(OrderProducts, OrderProduct)

	}
	return OrderProducts, nil
}

// UPDATE ORDER PRODUCTS

func (o orderProductsRepo) UpdateOrderProducts(OP structs.OrderProducts) error {

	_, err := o.DB.Exec(`update order_products set order_id = $1, product_id = $2, quantity = $3, price = $4 where id = $5`,)
	if err != nil {
		return err
	}
	return nil
}

// DELETE ORDER PRODUCT

func (o orderProductsRepo) DeleteOrderProducts(id uuid.UUID) error {
	_, err := o.DB.Exec(`delete from order_products where id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
