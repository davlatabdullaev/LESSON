package api

import (
	"mini_market/api/handler"
	"net/http"
)

func New(h handler.Handler) {
	http.HandleFunc("/staff", h.Staff)
}
