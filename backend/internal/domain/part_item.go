package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Part UoM values (MVP; extend in code when adding units).
const (
	PartUOMUnit  = "unit"
	PartUOMLiter = "liter"
)

// PartItem is a spare-parts catalog row with on-hand quantity (parts-inventory spec).
type PartItem struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Reference       string     `json:"reference" gorm:"column:reference;type:varchar(120);not null"`
	Brand           string     `json:"brand" gorm:"column:brand;type:varchar(120);not null"`
	Name            string     `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Barcode         string     `json:"barcode,omitempty" gorm:"column:barcode;type:varchar(64)"`
	Quantity        float64    `json:"quantity" gorm:"column:quantity;not null"`
	UOM             string     `json:"uom" gorm:"column:uom;type:varchar(16);not null"`
	MinimumQuantity *float64   `json:"minimumQuantity,omitempty" gorm:"column:minimum_quantity"`
	CreatedAt       time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty" gorm:"column:deleted_at;index"`
}

func (PartItem) TableName() string {
	return "part_items"
}

// Validate checks required fields and UoM (repository does not enforce business duplicates).
func (p *PartItem) Validate() error {
	if p == nil {
		return errors.New("part item is nil")
	}
	if strings.TrimSpace(p.Reference) == "" {
		return errors.New("reference is required")
	}
	if strings.TrimSpace(p.Brand) == "" {
		return errors.New("brand is required")
	}
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("name is required")
	}
	if p.Quantity < 0 {
		return errors.New("quantity must be non-negative")
	}
	switch strings.TrimSpace(p.UOM) {
	case PartUOMUnit, PartUOMLiter:
	default:
		return errors.New("invalid uom")
	}
	if p.MinimumQuantity != nil && *p.MinimumQuantity < 0 {
		return errors.New("minimum quantity must be non-negative")
	}
	p.Reference = strings.TrimSpace(p.Reference)
	p.Brand = strings.TrimSpace(p.Brand)
	p.Name = strings.TrimSpace(p.Name)
	p.Barcode = strings.TrimSpace(p.Barcode)
	p.UOM = strings.TrimSpace(p.UOM)
	return nil
}
