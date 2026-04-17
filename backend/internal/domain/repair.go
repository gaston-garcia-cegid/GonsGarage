package domain

import (
	"time"

	"github.com/google/uuid"
)

type RepairStatus string

const (
	RepairStatusPending    RepairStatus = "pending"
	RepairStatusInProgress RepairStatus = "in_progress"
	RepairStatusCompleted  RepairStatus = "completed"
	RepairStatusCancelled  RepairStatus = "cancelled"
)

type Repair struct {
	ID           uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CarID        uuid.UUID    `json:"car_id" gorm:"type:uuid;column:car_id;not null;index"`
	TechnicianID uuid.UUID    `json:"technician_id" gorm:"type:uuid;column:technician_id;not null;index"`
	Description  string       `json:"description" gorm:"not null"`
	Status       RepairStatus `json:"status" gorm:"not null;default:'pending'"`
	Cost         float64      `json:"cost" gorm:"type:decimal(10,2);default:0"`
	StartedAt    *time.Time   `json:"started_at,omitempty" gorm:"column:started_at"`
	CompletedAt  *time.Time   `json:"completed_at,omitempty" gorm:"column:completed_at"`
	CreatedAt    time.Time    `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time    `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    *time.Time   `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships - these will be ignored by GORM for auto-migration
	Car        Car  `json:"car,omitempty" gorm:"-"`
	Technician User `json:"technician,omitempty" gorm:"-"`
}

// TableName specifies the table name for GORM
func (Repair) TableName() string {
	return "repairs"
}

// ValidateRepairStatus checks if repair status is valid
func ValidateRepairStatus(status RepairStatus) bool {
	switch status {
	case RepairStatusPending, RepairStatusInProgress, RepairStatusCompleted, RepairStatusCancelled:
		return true
	default:
		return false
	}
}
