package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Car entity
type Car struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Make         string     `json:"make" gorm:"not null"`
	Model        string     `json:"model" gorm:"not null"`
	Year         int        `json:"year" gorm:"not null"`
	LicensePlate string     `json:"license_plate" gorm:"column:license_plate;uniqueIndex;not null"`
	VIN          string     `json:"vin" gorm:"column:vin;uniqueIndex"`
	Color        string     `json:"color"`
	OwnerID      uuid.UUID  `json:"owner_id" gorm:"type:uuid;column:owner_id;not null;index"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships - these will be ignored by GORM for auto-migration
	Owner   User     `json:"owner,omitempty" gorm:"-"`
	Repairs []Repair `json:"repairs,omitempty" gorm:"-"`
}

// TableName specifies the table name for GORM
func (Car) TableName() string {
	return "cars"
}

// Car validation errors
var (
	ErrCarNotFound         = errors.New("car not found")
	ErrRepairNotFound      = errors.New("repair not found")
	ErrAppointmentNotFound = errors.New("appointment not found")
	ErrInvalidCarData      = errors.New("invalid car data")
	ErrUnauthorizedAccess  = errors.New("unauthorized access")
)
