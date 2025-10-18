package services

import (
	"context"

	"gonsgarage/backend/internal/core/domain"
)

// NotificationService defines notification operations
type NotificationService interface {
	// Send SMS notification
	SendSMS(ctx context.Context, notification *domain.SMSNotification) error

	// Send WhatsApp message
	SendWhatsApp(ctx context.Context, notification *domain.WhatsAppNotification) error

	// Send email notification
	SendEmail(ctx context.Context, notification *domain.EmailNotification) error

	// Schedule notification for later
	ScheduleNotification(ctx context.Context, notification *domain.ScheduledNotification) error

	// Cancel scheduled notification
	CancelScheduledNotification(ctx context.Context, notificationID string) error

	// Get notification delivery status
	GetDeliveryStatus(ctx context.Context, notificationID string) (*domain.NotificationStatus, error)

	// Send bulk notifications
	SendBulkNotifications(ctx context.Context, notifications []*domain.BulkNotification) error

	// Get notification templates
	GetTemplate(ctx context.Context, templateID string) (*domain.NotificationTemplate, error)
}
