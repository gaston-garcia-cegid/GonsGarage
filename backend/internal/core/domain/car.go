package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Car entity
type Car struct {
	ID           uuid.UUID  `json:"id"`
	ClientID     uuid.UUID  `json:"client_id"`
	Make         string     `json:"make"`
	Model        string     `json:"model"`
	Year         int        `json:"year"`
	LicensePlate string     `json:"license_plate"`
	VIN          string     `json:"vin"`
	Color        string     `json:"color"`
	Mileage      int        `json:"mileage"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Client  *User    `json:"client,omitempty"`
	Repairs []Repair `json:"repairs,omitempty"`
}

// Repair entity
type Repair struct {
	ID          uuid.UUID  `json:"id"`
	CarID       uuid.UUID  `json:"car_id"`
	EmployeeID  uuid.UUID  `json:"employee_id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Cost        float64    `json:"cost"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Car      *Car  `json:"car,omitempty"`
	Employee *User `json:"employee,omitempty"`
}

// Appointment entity
type Appointment struct {
	ID            uuid.UUID  `json:"id"`
	CarID         uuid.UUID  `json:"car_id"`
	ClientID      uuid.UUID  `json:"client_id"`
	ScheduledDate time.Time  `json:"scheduled_date"`
	Status        string     `json:"status"`
	Description   string     `json:"description"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Car    *Car  `json:"car,omitempty"`
	Client *User `json:"client,omitempty"`
}

// Repair statuses
const (
	RepairStatusPending    = "pending"
	RepairStatusInProgress = "in_progress"
	RepairStatusCompleted  = "completed"
	RepairStatusCancelled  = "cancelled"
)

// Appointment statuses
const (
	AppointmentStatusScheduled = "scheduled"
	AppointmentStatusConfirmed = "confirmed"
	AppointmentStatusCancelled = "cancelled"
	AppointmentStatusCompleted = "completed"
)

// Car validation errors
var (
	ErrCarNotFound         = errors.New("car not found")
	ErrRepairNotFound      = errors.New("repair not found")
	ErrAppointmentNotFound = errors.New("appointment not found")
	ErrInvalidCarData      = errors.New("invalid car data")
	ErrUnauthorizedAccess  = errors.New("unauthorized access")
)

// ValidateRepairStatus checks if repair status is valid
func ValidateRepairStatus(status string) bool {
	switch status {
	case RepairStatusPending, RepairStatusInProgress, RepairStatusCompleted, RepairStatusCancelled:
		return true
	default:
		return false
	}
}

// ValidateAppointmentStatus checks if appointment status is valid
func ValidateAppointmentStatus(status string) bool {
	switch status {
	case AppointmentStatusScheduled, AppointmentStatusConfirmed, AppointmentStatusCancelled, AppointmentStatusCompleted:
		return true
	default:
		return false
	}
}
