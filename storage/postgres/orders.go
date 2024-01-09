package postgres

import (
	"basa/structs"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ordersRepo struct {
	DB *sql.DB
}

func NewOrdersRepo(db *sql.DB) ordersRepo {
	return ordersRepo{
		DB: db,
	}
}

// INSERT ORDER

func (o ordersRepo) InsertOrders(s structs.Orders) error {
	id := uuid.New()
	createdAt := time.Now()
	if _, err := o.DB.Exec(`insert into orders values ($1, $2, $3, $4)`, id, s.Amount, s.UserID, createdAt); err != nil {
		return err
	}
	return nil
}

// GET ORDER BY ID

func (o ordersRepo) GetByIDOrder(id uuid.UUID) (structs.Orders, error) {
	order := structs.Orders{}

	if err := o.DB.QueryRow(`select id, amount, user_id, created_at from orders where id = $1`, id).Scan(
		&order.ID,
		&order.Amount,
		&order.UserID,
		&order.CreatedAt,
	); err != nil {
		return structs.Orders{}, err
	}
	return order, nil
}

// GET ORDERS LIST

func (o ordersRepo) GetListOrder() ([]structs.Orders, error) {
	Orders := []structs.Orders{}

	rows, err := o.DB.Query(`select * from orders`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		Order := structs.Orders{}

		rows.Scan(&Order.ID, &Order.Amount, &Order.UserID, &Order.CreatedAt)

		Orders = append(Orders, Order)

	}
	return Orders, nil
}

// UPDATE ORDER

func (o ordersRepo) UpdateOrders(order structs.Orders) error {

	_, err := o.DB.Exec(`update orders set amount = $1, user_id = $2 where id = $3`, order.Amount, order.UserID, order.ID)
	if err != nil {
		return err
	}
	return nil
}

// DELETE ORDER

func (o ordersRepo) DeleteOrders(id uuid.UUID) error {
	_, err := o.DB.Exec(`delete from users where id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
