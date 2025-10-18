package external

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
)

// NotificationService defines the interface for the notification service
type NotificationService interface {
	SendWorkshopCreatedNotification(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopUpdatedNotification(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopDeletedNotification(ctx context.Context, workshop *domain.Workshop) error
}
