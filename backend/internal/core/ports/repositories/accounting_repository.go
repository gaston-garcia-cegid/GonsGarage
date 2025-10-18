package repositories

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
)

// AccountingRepository defines the interface for the accounting repository
type AccountingRepository interface {
	CreateInvoice(ctx context.Context, invoice *domain.Invoice) error
	UpdateInvoice(ctx context.Context, invoice *domain.Invoice) error
	DeleteInvoice(ctx context.Context, id string) error
	GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error)
	ListInvoices(ctx context.Context, limit, offset int) ([]*domain.Invoice, int64, error)
}
