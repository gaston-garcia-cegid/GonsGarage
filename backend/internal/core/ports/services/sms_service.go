package services

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
)

// SMSService defines the interface for the SMS service
type SMSService interface {
	SendWorkshopCreatedSMS(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopUpdatedSMS(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopDeletedSMS(ctx context.Context, workshop *domain.Workshop) error
}
