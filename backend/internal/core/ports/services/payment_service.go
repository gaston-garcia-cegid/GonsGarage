package services

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
)

// PaymentService defines the interface for the payment service
type PaymentService interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	UpdatePayment(ctx context.Context, payment *domain.Payment) error
	DeletePayment(ctx context.Context, id string) error
	GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error)
	ListPayments(ctx context.Context, limit, offset int) ([]*domain.Payment, int64, error)
}
