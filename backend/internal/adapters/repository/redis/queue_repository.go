package redis

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports/repositories"
)

// QueueRepository is a Redis implementation of the QueueRepository interface
type QueueRepository struct {
	// client *redis.Client
}

// NewQueueRepository creates a new QueueRepository
func NewQueueRepository() repositories.QueueRepository {
	return &QueueRepository{}
}

// Enqueue implements QueueRepository.Enqueue
func (r *QueueRepository) Enqueue(ctx context.Context, job *domain.Job) error {
	// TODO: Implementar
	return nil
}

// Dequeue implements QueueRepository.Dequeue
func (r *QueueRepository) Dequeue(ctx context.Context) (*domain.Job, error) {
	// TODO: Implementar
	return nil, nil
}

// Acknowledge implements QueueRepository.Acknowledge
func (r *QueueRepository) Acknowledge(ctx context.Context, id string) error {
	// TODO: Implementar
	return nil
}

// Reject implements QueueRepository.Reject
func (r *QueueRepository) Reject(ctx context.Context, id string) error {
	// TODO: Implementar
	return nil
}

// List implements QueueRepository.List
func (r *QueueRepository) List(ctx context.Context, limit, offset int) ([]*domain.Job, int64, error) {
	// TODO: Implementar
	return nil, 0, nil
}

// SearchByName implements QueueRepository.SearchByName
func (r *QueueRepository) SearchByName(ctx context.Context, name string, limit int) ([]*domain.Job, error) {
	// TODO: Implementar
	return nil, nil
}
