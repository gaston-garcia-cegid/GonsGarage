package postgres

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"gorm.io/gorm"
)

// NotificationRepository is a PostgreSQL implementation of the NotificationRepository interface
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new NotificationRepository
func NewNotificationRepository(db *gorm.DB) ports.NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create implements NotificationRepository.Create
func (r *NotificationRepository) SendWorkshopCreatedNotification(ctx context.Context, workshop *domain.Workshop) error {
	// TODO: Implementar
	return nil
}

// Update implements NotificationRepository.Update
func (r *NotificationRepository) SendWorkshopUpdatedNotification(ctx context.Context, workshop *domain.Workshop) error {
	// TODO: Implementar
	return nil
}

// Delete implements NotificationRepository.Delete
func (r *NotificationRepository) SendWorkshopDeletedNotification(ctx context.Context, workshop *domain.Workshop) error {
	// TODO: Implementar
	return nil
}
