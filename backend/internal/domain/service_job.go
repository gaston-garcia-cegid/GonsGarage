package domain

import (
	"time"

	"github.com/google/uuid"
)

// ServiceJobStatus is the workshop visit lifecycle.
type ServiceJobStatus string

const (
	ServiceJobStatusOpen        ServiceJobStatus = "open"
	ServiceJobStatusInProgress  ServiceJobStatus = "in_progress"
	ServiceJobStatusClosed      ServiceJobStatus = "closed"
	ServiceJobStatusCancelled   ServiceJobStatus = "cancelled"
)

// ServiceJob is one vehicle visit in the workshop (taller). Root aggregate for reception/handover.
type ServiceJob struct {
	ID             uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CarID          uuid.UUID        `json:"car_id" gorm:"type:uuid;not null;index"`
	Status         ServiceJobStatus `json:"status" gorm:"not null;default:'open'"`
	OpenedByUserID uuid.UUID        `json:"opened_by_user_id" gorm:"type:uuid;not null;index"`
	OpenedAt       time.Time        `json:"opened_at" gorm:"not null"`
	ClosedAt       *time.Time       `json:"closed_at,omitempty"`
	AppointmentID  *uuid.UUID       `json:"appointment_id,omitempty" gorm:"type:uuid;index"`
	CreatedAt      time.Time        `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time        `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt      *time.Time       `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (ServiceJob) TableName() string { return "service_jobs" }

// ValidateServiceJobStatus returns true for known status values.
func ValidateServiceJobStatus(s ServiceJobStatus) bool {
	switch s {
	case ServiceJobStatusOpen, ServiceJobStatusInProgress, ServiceJobStatusClosed, ServiceJobStatusCancelled:
		return true
	default:
		return false
	}
}

// ServiceJobReception is 1:1 reception checklist for a service job.
type ServiceJobReception struct {
	ServiceJobID     uuid.UUID `json:"service_job_id" gorm:"type:uuid;primaryKey"`
	OdometerKM       int       `json:"odometer_km" gorm:"not null"`
	OilLevel         string    `json:"oil_level,omitempty" gorm:"type:text"`
	CoolantLevel     string    `json:"coolant_level,omitempty" gorm:"type:text"`
	TiresNote        string    `json:"tires_note,omitempty" gorm:"type:text"`
	GeneralNotes     string    `json:"general_notes,omitempty" gorm:"type:text"`
	RecordedByUserID uuid.UUID `json:"recorded_by_user_id" gorm:"type:uuid;not null"`
	RecordedAt       time.Time `json:"recorded_at" gorm:"not null"`
	SchemaVersion    int       `json:"schema_version" gorm:"not null;default:1"`
}

func (ServiceJobReception) TableName() string { return "service_job_receptions" }

// ServiceJobHandover is 1:1 delivery checklist (close visit).
type ServiceJobHandover struct {
	ServiceJobID     uuid.UUID `json:"service_job_id" gorm:"type:uuid;primaryKey"`
	OdometerKM       int       `json:"odometer_km" gorm:"not null"`
	TiresNote        string    `json:"tires_note,omitempty" gorm:"type:text"`
	GeneralNotes     string    `json:"general_notes,omitempty" gorm:"type:text"`
	RecordedByUserID uuid.UUID `json:"recorded_by_user_id" gorm:"type:uuid;not null"`
	RecordedAt       time.Time `json:"recorded_at" gorm:"not null"`
	SchemaVersion    int       `json:"schema_version" gorm:"not null;default:1"`
}

func (ServiceJobHandover) TableName() string { return "service_job_handovers" }
