package services

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
)

// NotificationService defines notification operations
type NotificationService interface {
	SendWorkshopCreatedNotification(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopUpdatedNotification(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopDeletedNotification(ctx context.Context, workshop *domain.Workshop) error

	QueueNotification(ctx context.Context, notification NotificationRequest) error
}

type NotificationRequest struct {
	Type    string `json:"type"` // "sms" or "whatsapp"
	To      string `json:"to"`
	Message string `json:"message"`
}
