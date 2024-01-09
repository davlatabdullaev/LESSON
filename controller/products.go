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

func (c Controller) Products(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.InsertProduct(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		fmt.Println("values: ", values, time.Now())
		_, ok := values["id"]
		if ok {
			c.GetProductByID(w, r)
		} else {
			c.GetProductList(w, r)
		}
	case http.MethodPut:
		c.UpdateProductByID(w, r)
	case http.MethodDelete:
		c.DeleteProductByID(w, r)
	}
}

func (c Controller) InsertProduct(w http.ResponseWriter, r *http.Request) {
	var product structs.Products

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Fatal("error while reading data from client", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.Store.ProductsStorage.InsertProducts(product)
	if err != nil {
		log.Fatal("error while inserting product inside controller err: ", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handResponse(w, http.StatusCreated, nil)
}

func (c Controller) GetProductByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if id == "" {
		handResponse(w, http.StatusBadRequest, errors.New("missing product ID"))
		return
	}

	uid := uuid.MustParse(id)
	product, err := c.Store.ProductsStorage.GetByIDProduct(uid)
	if err != nil {
		log.Fatal("error while getting product by id ", err.Error())
		handResponse(w, http.StatusInternalServerError, err)
		return
	}

	handResponse(w, http.StatusOK, product)
}

func (c Controller) GetProductList(w http.ResponseWriter, r *http.Request) {
	products, err := c.Store.ProductsStorage.GetListProduct()
	if err != nil {
		log.Fatal("error while get list of products ", err.Error())
		handResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handResponse(w, http.StatusOK, products)
}

func (c Controller) UpdateProductByID(w http.ResponseWriter, r *http.Request) {
	var product structs.Products

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Fatal("error while reading data from client", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.Store.ProductsStorage.UpdateProducts(product)
	if err != nil {
		log.Fatal("error while updating product by ID inside controller err: ", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handResponse(w, http.StatusOK, "Product updated successfully")
}

func (c Controller) DeleteProductByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if id == "" {
		handResponse(w, http.StatusBadRequest, errors.New("missing product ID"))
		return
	}

	uid := uuid.MustParse(id)
	err := c.Store.ProductsStorage.DeleteProducts(uid)
	if err != nil {
		log.Fatal("error while deleting product by ID inside controller err: ", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handResponse(w, http.StatusOK, "Product deleted successfully")
}
