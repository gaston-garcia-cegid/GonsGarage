package repositories

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
)

// NotificationRepository defines the interface for the notification repository
type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *domain.Notification) error
	UpdateNotification(ctx context.Context, notification *domain.Notification) error
	DeleteNotification(ctx context.Context, id string) error
	GetNotificationByID(ctx context.Context, id string) (*domain.Notification, error)
	ListNotifications(ctx context.Context, limit, offset int) ([]*domain.Notification, int64, error)
}
