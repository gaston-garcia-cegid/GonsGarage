package domain

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	Address     string
	City        string
	State       string
	ZipCode     string
	DateOfBirth *time.Time
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Cars        []Car
}
