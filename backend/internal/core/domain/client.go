package domain

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	// Add other relevant fields here
}
