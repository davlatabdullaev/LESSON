package controller

import (
	"basa/structs"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c Controller) Orders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.InsertOrders(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		fmt.Println("values: ", values, time.Now())
		_, ok := values["id"]
		if ok {
			c.GetOrderById(w, r)
		} else {
			c.GetOrdersList(w, r)
		}
	case http.MethodPut:
		c.UpdateOrder(w, r)
	case http.MethodDelete:
		c.DeleteOrder(w, r)
	}
}

// CREATE
func (c Controller) InsertOrders(w http.ResponseWriter, r *http.Request) {
	order := structs.Orders{}

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		fmt.Println("error while reading data from client", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.Store.OrdersStorage.InsertOrders(order)
	if err != nil {
		fmt.Println("error while inserting order inside controller err: ", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handResponse(w, http.StatusCreated, nil)
}

// READ BY ID
func (c Controller) GetOrderById(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	fmt.Println("values : ", values)
	id := values["id"][0]
	uid := uuid.MustParse(id)
	order, err := c.Store.OrdersStorage.GetByIDOrder(uid)
	if err != nil {
		fmt.Println("error while getting order by id ", err.Error())
		handResponse(w, http.StatusInternalServerError, err)
		return
	}

	handResponse(w, http.StatusOK, order)
}

// READ ALL
func (c Controller) GetOrdersList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.GetOrders(w, r)
	default:
		handResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (c Controller) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := c.Store.OrdersStorage.GetListOrder()
	if err != nil {
		fmt.Println("error while get list of orders ", err.Error())
		handResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handResponse(w, http.StatusOK, orders)
}

// UPDATE
func (c Controller) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if id == "" {
		handResponse(w, http.StatusBadRequest, errors.New("missing order ID"))
		return
	}

	order := structs.Orders{}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		fmt.Println("error while reading data from client", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.Store.OrdersStorage.UpdateOrders(order)
	if err != nil {
		fmt.Println("error while updating order by ID inside controller err: ", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handResponse(w, http.StatusOK, "Order updated successfully")
}

// DELETE
func (c Controller) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if id == "" {
		handResponse(w, http.StatusBadRequest, errors.New("missing order ID"))
		return
	}

	uid := uuid.MustParse(id)
	err := c.Store.OrdersStorage.DeleteOrders(uid)
	if err != nil {
		fmt.Println("error while deleting order by ID inside controller err: ", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handResponse(w, http.StatusOK, "Order deleted successfully")
}
