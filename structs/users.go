package structs

import "github.com/google/uuid"

type Users struct {
	ID         uuid.UUID
	First_Name string
	Last_Name  string
	Email      string
	Phone      string
}
