package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
)

type postgresReceivedInvoiceRepository struct {
	db *gorm.DB
}

// NewPostgresReceivedInvoiceRepository returns ports.ReceivedInvoiceRepository.
func NewPostgresReceivedInvoiceRepository(db *gorm.DB) ports.ReceivedInvoiceRepository {
	return &postgresReceivedInvoiceRepository{db: db}
}

func (r *postgresReceivedInvoiceRepository) Create(ctx context.Context, inv *domain.ReceivedInvoice) error {
	if inv == nil {
		return fmt.Errorf("received invoice is nil")
	}
	if err := r.db.WithContext(ctx).Create(inv).Error; err != nil {
		return fmt.Errorf("failed to create received invoice: %w", err)
	}
	return nil
}

func (r *postgresReceivedInvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ReceivedInvoice, error) {
	var inv domain.ReceivedInvoice
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrReceivedInvoiceNotFound
		}
		return nil, fmt.Errorf("failed to get received invoice: %w", err)
	}
	return &inv, nil
}

func (r *postgresReceivedInvoiceRepository) Update(ctx context.Context, inv *domain.ReceivedInvoice) error {
	if inv == nil {
		return fmt.Errorf("received invoice is nil")
	}
	res := r.db.WithContext(ctx).Model(&domain.ReceivedInvoice{}).
		Where("id = ? AND deleted_at IS NULL", inv.ID).
		Updates(map[string]interface{}{
			"supplier_id":  inv.SupplierID,
			"vendor_name":  inv.VendorName,
			"category":     inv.Category,
			"amount":       inv.Amount,
			"invoice_date": inv.InvoiceDate,
			"notes":        inv.Notes,
			"updated_at":   time.Now().UTC(),
		})
	if res.Error != nil {
		return fmt.Errorf("failed to update received invoice: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrReceivedInvoiceNotFound
	}
	return nil
}

func (r *postgresReceivedInvoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	res := r.db.WithContext(queryCtx).Model(&domain.ReceivedInvoice{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now().UTC())
	if res.Error != nil {
		return fmt.Errorf("failed to delete received invoice: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrReceivedInvoiceNotFound
	}
	return nil
}

func (r *postgresReceivedInvoiceRepository) List(ctx context.Context, limit, offset int) ([]*domain.ReceivedInvoice, int64, error) {
	limit, offset = clampRepoList(limit, offset)
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.ReceivedInvoice{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count received invoices: %w", err)
	}
	var rows []domain.ReceivedInvoice
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Model(&domain.ReceivedInvoice{}).
		Order("invoice_date DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list received invoices: %w", err)
	}
	out := make([]*domain.ReceivedInvoice, 0, len(rows))
	for i := range rows {
		row := rows[i]
		out = append(out, &row)
	}
	return out, total, nil
}

var _ ports.ReceivedInvoiceRepository = (*postgresReceivedInvoiceRepository)(nil)
