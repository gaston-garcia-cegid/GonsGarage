package domain

import (
	"time"

	"github.com/google/uuid"
)

// Car represents a vehicle entity in the garage management system
type Car struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Make         string     `json:"make" gorm:"not null"`
	Model        string     `json:"model" gorm:"not null"`
	Year         int        `json:"year" gorm:"not null"`
	LicensePlate string     `json:"licensePlate" gorm:"column:license_plate;uniqueIndex;not null"`
	VIN          string     `json:"vin" gorm:"column:vin;uniqueIndex"`
	Color        string     `json:"color" gorm:"not null"`
	Mileage      int        `json:"mileage" gorm:"default:0"`
	OwnerID      uuid.UUID  `json:"ownerID" gorm:"type:uuid;column:owner_id;not null;index"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty" gorm:"column:deleted_at;index"`

	// Relationships
	Owner   User     `json:"owner,omitempty" gorm:"foreignKey:OwnerID;references:ID"`
	Repairs []Repair `json:"repairs,omitempty" gorm:"foreignKey:CarID;references:ID"`
}

// TableName specifies the database table name
func (Car) TableName() string {
	return "cars"
}

// Validate validates the car entity according to business rules
func (c *Car) Validate() error {
	if c.Make == "" {
		return ErrInvalidCarData // ✅ Use domain error
	}
	if c.Model == "" {
		return ErrInvalidCarData
	}
	if c.Year < minCarYear || c.Year > maxCarYear {
		return ErrInvalidCarData
	}
	if c.LicensePlate == "" {
		return ErrInvalidCarData
	}
	if c.Color == "" {
		return ErrInvalidCarData
	}
	if c.Mileage < 0 {
		return ErrInvalidCarData
	}
	return nil
}

// IsOwnedBy checks if the car belongs to the specified user
func (c *Car) IsOwnedBy(userID uuid.UUID) bool {
	return c.OwnerID == userID
}

// Business constants
const (
	minCarYear = 1900
	maxCarYear = 2030 // Allow future model years
)
