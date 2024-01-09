package controller

import (
	"basa/structs"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
)

func (c Controller) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.InsertUser(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		fmt.Println("values: ", values, time.Now())
		_, ok := values["id"]
		if ok {
			c.GetUserByID(w, r)
		} else {
			c.GetUsersList(w, r)
		}
	case http.MethodPut:
		c.UpdateUser(w, r)
	case http.MethodDelete:
		c.DeleteUser(w, r)
	}
}

func (c Controller) InsertUser(w http.ResponseWriter, r *http.Request) {
	var user structs.Users

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal("error while reading data from client", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	phoneRegex := regexp.MustCompile(`^\+\d{12}$`)
	if !phoneRegex.MatchString(user.Phone) {
		handResponse(w, http.StatusBadRequest, errors.New("phone number is not correct"))
		return
	}

	id, err := c.Store.UsersStorage.InsertUser(user)
	if err != nil {
		log.Fatal("error while creating data inside controller err: ", err.Error())
		handResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handResponse(w, http.StatusCreated, id)
}

func (c Controller) GetUserByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if id == "" {
		handResponse(w, http.StatusBadRequest, errors.New("missing user ID"))
		return
	}

	uid := uuid.MustParse(id)
	user, err := c.Store.UsersStorage.GetByIDUser(uid)
	if err != nil {
		log.Fatal("error while getting user by id ", err.Error())
		handResponse(w, http.StatusInternalServerError, err)
		return
	}

	handResponse(w, http.StatusOK, user)
}

func (c Controller) GetUsersList(w http.ResponseWriter, r *http.Request) {
	users, err := c.Store.UsersStorage.GetListUser()
	if err != nil {
		log.Fatal("error while get list of users ", err.Error())
		handResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handResponse(w, http.StatusOK, users)
}

func (c Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user structs.Users

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal("error while reading data from client", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	c.Store.UsersStorage.UpdateUser(user)

	handResponse(w, http.StatusOK, "User updated successfully")
}

func (c Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if id == "" {
		handResponse(w, http.StatusBadRequest, errors.New("missing user ID"))
		return
	}

	uid := uuid.MustParse(id)
	c.Store.UsersStorage.DeleteUser(uid)

	handResponse(w, http.StatusOK, "User deleted successfully")
}

func getUserInfoForUpdate(w http.ResponseWriter, r *http.Request) structs.Users {
	var user structs.Users

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal("error while reading data from client", err.Error())
		handResponse(w, http.StatusBadRequest, err.Error())
		return user
	}

	return user
}
