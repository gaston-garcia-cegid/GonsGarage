package redis

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
)

// SessionRepository is a Redis implementation of the SessionRepository interface
type SessionRepository struct {
	// client *redis.Client
}

// NewSessionRepository creates a new SessionRepository
func NewSessionRepository() ports.SessionRepository {
	return &SessionRepository{}
}

// Create implements SessionRepository.Create
func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	// TODO: Implementar
	return nil
}

// GetByID implements SessionRepository.GetByID
func (r *SessionRepository) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	// TODO: Implementar
	return nil, nil
}

// Update implements SessionRepository.Update
func (r *SessionRepository) Update(ctx context.Context, session *domain.Session) error {
	// TODO: Implementar
	return nil
}

// Delete implements SessionRepository.Delete
func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	// TODO: Implementar
	return nil
}

// List implements SessionRepository.List
func (r *SessionRepository) List(ctx context.Context, limit, offset int) ([]*domain.Session, int64, error) {
	// TODO: Implementar
	return nil, 0, nil
}

// SearchByName implements SessionRepository.SearchByName
func (r *SessionRepository) SearchByName(ctx context.Context, name string, limit int) ([]*domain.Session, error) {
	// TODO: Implementar
	return nil, nil
}
