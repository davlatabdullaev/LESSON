package structs

import (
	"time"

	"github.com/google/uuid"
)


type Orders struct {
ID           uuid.UUID
Amount       string
UserID       uuid.UUID 
CreatedAt    time.Time
}
