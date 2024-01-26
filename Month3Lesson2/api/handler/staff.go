package handler


import (
	"encoding/json"
	"errors"
	"mini_market/api/models"
	"net/http"
)

func (h Handler) Staff(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateStaff(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetStaffList(w)
		} else {
			h.GetStaffByID(w, r)
		}
	case http.MethodPut:
		values := r.URL.Query()
		if _, ok := values["route"]; ok {
		} else {
			h.UpdateStaff(w, r)
		}
	case http.MethodDelete:
		h.DeleteStaff(w, r)
	}
}

func (h Handler) CreateStaff(w http.ResponseWriter, r *http.Request) {
	createStaff := models.CreateStaff{}

	if err := json.NewDecoder(r.Body).Decode(&createStaff); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Staff().Create(createStaff)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	car, err := h.storage.Staff().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, car)

}

func (h Handler) GetStaffByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	car, err := h.storage.Staff().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, car)

}

func (h Handler) GetStaffList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		err         error
	)

	response, err := h.storage.Staff().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})

	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, response)

}

func (h Handler) UpdateStaff(w http.ResponseWriter, r *http.Request) {
	updateCar := models.UpdateStaff{}

	if err := json.NewDecoder(r.Body).Decode(&updateCar); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Staff().Update(updateCar)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	car, err := h.storage.Staff().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, car)

}

func (h Handler) DeleteStaff(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Staff().Delete(models.DeleteStaff{
		ID: id,
	}); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
