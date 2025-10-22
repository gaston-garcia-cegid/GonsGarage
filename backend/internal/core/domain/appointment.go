package domain

import (
	"time"

	"github.com/google/uuid"
)

type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
)

type Appointment struct {
	ID          uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID  uuid.UUID         `json:"customer_id" gorm:"type:uuid;column:customer_id;not null;index"`
	CarID       uuid.UUID         `json:"car_id" gorm:"type:uuid;column:car_id;not null;index"`
	ServiceType string            `json:"service_type" gorm:"column:service_type;not null"`
	Status      AppointmentStatus `json:"status" gorm:"not null;default:'scheduled'"`
	ScheduledAt time.Time         `json:"scheduled_at" gorm:"column:scheduled_at;not null"`
	Notes       string            `json:"notes"`
	CreatedAt   time.Time         `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   *time.Time        `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships - these will be ignored by GORM for auto-migration
	Customer User `json:"customer,omitempty" gorm:"-"`
	Car      Car  `json:"car,omitempty" gorm:"-"`
}

// TableName specifies the table name for GORM
func (Appointment) TableName() string {
	return "appointments"
}

// ValidateAppointmentStatus checks if appointment status is valid
func ValidateAppointmentStatus(status AppointmentStatus) bool {
	switch status {
	case AppointmentStatusScheduled, AppointmentStatusConfirmed, AppointmentStatusCancelled, AppointmentStatusCompleted:
		return true
	default:
		return false
	}
}
