package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Supplier is a vendor master record (suppliers spec).
type Supplier struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Name          string     `json:"name" gorm:"type:varchar(200);not null"`
	ContactEmail  string     `json:"contactEmail" gorm:"column:contact_email;type:varchar(200)"`
	ContactPhone  string     `json:"contactPhone" gorm:"column:contact_phone;type:varchar(40)"`
	TaxID         string     `json:"taxId,omitempty" gorm:"column:tax_id;type:varchar(64)"`
	Notes         string     `json:"notes" gorm:"type:text"`
	IsActive      bool       `json:"isActive" gorm:"column:is_active;default:true"`
	CreatedAt     time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty" gorm:"column:deleted_at;index"`
}

func (Supplier) TableName() string {
	return "suppliers"
}

func (s *Supplier) Validate() error {
	if s == nil {
		return errors.New("supplier is nil")
	}
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}
