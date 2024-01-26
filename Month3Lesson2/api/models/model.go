package models


type GetListRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PrimaryKey struct {
	ID string `json:"id"`
}