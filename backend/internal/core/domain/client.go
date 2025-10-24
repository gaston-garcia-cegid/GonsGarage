package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Client represents a client user in the system (following Agent.md domain modeling)
type Client struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	FirstName string     `json:"firstName" gorm:"column:first_name;not null"` // ✅ camelCase per Agent.md
	LastName  string     `json:"lastName" gorm:"column:last_name;not null"`   // ✅ camelCase per Agent.md
	Phone     string     `json:"phone"`
	Address   string     `json:"address"`
	City      string     `json:"city"`
	State     string     `json:"state"`
	ZipCode   string     `json:"zipCode" gorm:"column:zip_code"`                     // ✅ camelCase per Agent.md
	IsActive  bool       `json:"isActive" gorm:"column:is_active;default:true"`      // ✅ camelCase per Agent.md
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`  // ✅ camelCase per Agent.md
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`  // ✅ camelCase per Agent.md
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"column:deleted_at;index"` // ✅ camelCase per Agent.md

	// Relationships
	Cars    []Car    `json:"cars,omitempty" gorm:"foreignKey:OwnerID;references:ID"`
	Repairs []Repair `json:"repairs,omitempty" gorm:"many2many:client_repairs;"`
}

// TableName specifies the table name for GORM
func (Client) TableName() string {
	return "clients"
}

// Validate validates the client data according to business rules (following Agent.md validation)
func (c *Client) Validate() error {
	if c.Email == "" {
		return ErrInvalidClientData
	}

	// Email format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(c.Email) {
		return ErrInvalidClientData
	}

	if c.FirstName == "" {
		return ErrInvalidClientData
	}

	if c.LastName == "" {
		return ErrInvalidClientData
	}

	// Trim whitespace
	c.Email = strings.TrimSpace(strings.ToLower(c.Email))
	c.FirstName = strings.TrimSpace(c.FirstName)
	c.LastName = strings.TrimSpace(c.LastName)
	c.Phone = strings.TrimSpace(c.Phone)
	c.Address = strings.TrimSpace(c.Address)
	c.City = strings.TrimSpace(c.City)
	c.State = strings.TrimSpace(c.State)
	c.ZipCode = strings.TrimSpace(c.ZipCode)

	return nil
}

// GetFullName returns the client's full name
func (c *Client) GetFullName() string {
	return strings.TrimSpace(c.FirstName + " " + c.LastName)
}

// IsDeleted checks if the client is soft deleted
func (c *Client) IsDeleted() bool {
	return c.DeletedAt != nil
}

// Client-related domain errors (following Agent.md error handling)
// ErrClientNotFound is defined in errors.go
var (
	ErrClientAlreadyExists = errors.New("client already exists")
)
