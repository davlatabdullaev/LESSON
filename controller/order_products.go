package controller

import (
	"basa/structs"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c Controller) OrderProduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.CreateOrderProduct(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		fmt.Println("values: ", values, time.Now())
		_, ok := values["id"]
		if ok {
			c.GetOrderProductByID(w, r)
		} else {
			c.GetOrderProductsList(w, r)
		}
	case http.MethodPut:
		c.UpdateOrderProductByID(w, r)
	case http.MethodDelete:
		c.DeleteOrderProductByID(w,r)
	}
}

// CREATE

func (c Controller) CreateOrderProduct(w http.ResponseWriter, r *http.Request) {
	orderProduct := structs.OrderProducts{}

	if err := json.NewDecoder(r.Body).Decode(&orderProduct); err != nil {
		log.Fatal("error while reading data from client", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.Store.OrderProductsStorage.InsertOrderProduct(orderProduct)
	if err != nil {
		fmt.Println("error while inserting driver inside controller err: ", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handResponse(w, http.StatusCreated, nil)
}

// READ BY ID

func (c Controller) GetOrderProductByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	fmt.Println("values : ", values)
	id := values["id"][0]
	uid := uuid.MustParse(id)
	driver, err := c.Store.OrderProductsStorage.GetByIDOrderProduct(uid)
	if err != nil {
		fmt.Println("error while getting driver by id ", err.Error())
		handResponse(w, http.StatusInternalServerError, err)
		return
	}

	handResponse(w, http.StatusOK, driver)

}

// READ ALL

func (c Controller) GetOrderProductsList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.GetOrderProducts(w, r)
	default:
		handResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (c Controller) GetOrderProducts(w http.ResponseWriter, r *http.Request) {
	orderProducts, err := c.Store.OrderProductsStorage.GetListOrderProducts()
	if err != nil {
		fmt.Println("error while get list order products ", err.Error())
		handResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handResponse(w, http.StatusOK, orderProducts)
}

// UPDATE

func (c Controller) UpdateOrderProductByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if id == "" {
		handResponse(w, http.StatusBadRequest, errors.New("missing driver ID"))
		return
	}

	orderProduct := structs.OrderProducts{}
	if err := json.NewDecoder(r.Body).Decode(&orderProduct); err != nil {
		log.Fatal("error while reading data from client", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.Store.OrderProductsStorage.UpdateOrderProducts(orderProduct)
	if err != nil {
		fmt.Println("error while updating driver by ID inside controller err: ", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handResponse(w, http.StatusOK, "Driver updated successfully")

}

// DELETE

func (c Controller) DeleteOrderProductByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if id == "" {
		handResponse(w, http.StatusBadRequest, errors.New("missing driver ID"))
		return
	}

	uid := uuid.MustParse(id)
	err := c.Store.OrderProductsStorage.DeleteOrderProducts(uid)
	if err != nil {
		fmt.Println("error while deleting driver by ID inside controller err: ", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handResponse(w, http.StatusOK, "Driver deleted successfully")

}
