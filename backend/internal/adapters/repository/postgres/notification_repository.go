package postgres

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports/repositories"
)

// NotificationRepository is a PostgreSQL implementation of the NotificationRepository interface
type NotificationRepository struct {
	// db *gorm.DB
}

// NewNotificationRepository creates a new NotificationRepository
func NewNotificationRepository() repositories.NotificationRepository {
	return &NotificationRepository{}
}

// Create implements NotificationRepository.Create
func (r *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	// TODO: Implementar
	return nil
}

// GetByID implements NotificationRepository.GetByID
func (r *NotificationRepository) GetByID(ctx context.Context, id string) (*domain.Notification, error) {
	// TODO: Implementar
	return nil, nil
}

// Update implements NotificationRepository.Update
func (r *NotificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	// TODO: Implementar
	return nil
}

// Delete implements NotificationRepository.Delete
func (r *NotificationRepository) Delete(ctx context.Context, id string) error {
	// TODO: Implementar
	return nil
}

// List implements NotificationRepository.List
func (r *NotificationRepository) List(ctx context.Context, limit, offset int) ([]*domain.Notification, int64, error) {
	// TODO: Implementar
	return nil, 0, nil
}

// SearchByName implements NotificationRepository.SearchByName
func (r *NotificationRepository) SearchByName(ctx context.Context, name string, limit int) ([]*domain.Notification, error) {
	// TODO: Implementar
	return nil, nil
}
