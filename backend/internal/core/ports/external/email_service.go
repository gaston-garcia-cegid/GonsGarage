package external

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
)

// EmailService defines the interface for the email service
type EmailService interface {
	SendWorkshopCreatedEmail(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopUpdatedEmail(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopDeletedEmail(ctx context.Context, workshop *domain.Workshop) error
}
