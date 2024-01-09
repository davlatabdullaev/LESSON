package structs

import "github.com/google/uuid"

type Products struct {
	ID          uuid.UUID
	Price       int
	ProductName string
}
