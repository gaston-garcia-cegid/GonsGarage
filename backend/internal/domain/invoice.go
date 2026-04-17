package domain

import (
	"time"

	"github.com/google/uuid"
)

// Invoice represents a customer invoice (client may read/update own rows — see invoice service).
type Invoice struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	Amount      float64
	Status      string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
