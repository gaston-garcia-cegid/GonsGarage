package invoice

import (
	"context"
	"fmt"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/google/uuid"
)

type InvoiceService struct {
	invoiceRepo ports.InvoiceRepository
	userRepo    ports.UserRepository
}

func NewInvoiceService(invoiceRepo ports.InvoiceRepository, userRepo ports.UserRepository) *InvoiceService {
	return &InvoiceService{invoiceRepo: invoiceRepo, userRepo: userRepo}
}

func (s *InvoiceService) GetInvoice(ctx context.Context, invoiceID uuid.UUID, requestingUserID uuid.UUID) (*domain.Invoice, error) {
	u, err := s.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if u == nil {
		return nil, domain.ErrUserNotFound
	}

	inv, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, domain.ErrInvoiceNotFound
	}

	if u.IsClient() && inv.CustomerID != requestingUserID {
		return nil, domain.ErrUnauthorizedAccess
	}
	if !u.IsClient() && !u.IsEmployee() {
		return nil, domain.ErrUnauthorizedAccess
	}
	return inv, nil
}

// UpdateInvoice merges updates. Clients may only update invoices they own and may only change Notes (RU billing).
func (s *InvoiceService) UpdateInvoice(ctx context.Context, invoice *domain.Invoice, requestingUserID uuid.UUID) (*domain.Invoice, error) {
	if invoice == nil {
		return nil, fmt.Errorf("invoice is required")
	}
	u, err := s.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if u == nil {
		return nil, domain.ErrUserNotFound
	}

	existing, err := s.invoiceRepo.GetByID(ctx, invoice.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrInvoiceNotFound
	}

	if u.IsClient() {
		if existing.CustomerID != requestingUserID {
			return nil, domain.ErrUnauthorizedAccess
		}
		merged := *existing
		merged.Notes = invoice.Notes
		if err := s.invoiceRepo.Update(ctx, &merged); err != nil {
			return nil, err
		}
		return s.invoiceRepo.GetByID(ctx, merged.ID)
	}

	if !u.IsEmployee() {
		return nil, domain.ErrUnauthorizedAccess
	}

	merged := *existing
	merged.Notes = invoice.Notes
	if invoice.Status != "" {
		merged.Status = invoice.Status
	}
	if invoice.Amount != 0 {
		merged.Amount = invoice.Amount
	}
	if err := s.invoiceRepo.Update(ctx, &merged); err != nil {
		return nil, err
	}
	return s.invoiceRepo.GetByID(ctx, merged.ID)
}

func (s *InvoiceService) ListMyInvoices(ctx context.Context, requestingUserID uuid.UUID, limit, offset int) ([]*domain.Invoice, int64, error) {
	u, err := s.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user: %w", err)
	}
	if u == nil {
		return nil, 0, domain.ErrUserNotFound
	}
	if !u.IsClient() {
		return nil, 0, domain.ErrUnauthorizedAccess
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.invoiceRepo.ListByCustomerID(ctx, requestingUserID, limit, offset)
}
