package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ReceivedInvoice is a supplier/purchase invoice received by the workshop (invoices spec).
type ReceivedInvoice struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	SupplierID  *uuid.UUID `json:"supplierId,omitempty" gorm:"type:uuid;index"`
	VendorName  string     `json:"vendorName" gorm:"type:varchar(200)"` // free text if no supplier row
	Category    string     `json:"category" gorm:"type:varchar(80);not null"`
	Amount      float64    `json:"amount" gorm:"not null"`
	InvoiceDate time.Time  `json:"invoiceDate" gorm:"column:invoice_date;not null"`
	Notes       string     `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" gorm:"column:deleted_at;index"`
}

func (ReceivedInvoice) TableName() string {
	return "received_invoices"
}

func (r *ReceivedInvoice) Validate() error {
	if r == nil {
		return errors.New("received invoice is nil")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if r.InvoiceDate.IsZero() {
		return errors.New("invoice date is required")
	}
	if r.Category == "" {
		return errors.New("category is required")
	}
	return nil
}
