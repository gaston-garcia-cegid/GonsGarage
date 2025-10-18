package postgres

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports/repositories"
)

// WorkshopRepository is a PostgreSQL implementation of the WorkshopRepository interface
type WorkshopRepository struct {
	// db *gorm.DB
}

// NewWorkshopRepository creates a new WorkshopRepository
func NewWorkshopRepository() repositories.WorkshopRepository {
	return &WorkshopRepository{}
}

// Create implements WorkshopRepository.Create
func (r *WorkshopRepository) Create(ctx context.Context, workshop *domain.Workshop) error {
	// TODO: Implementar
	return nil
}

// GetByID implements WorkshopRepository.GetByID
func (r *WorkshopRepository) GetByID(ctx context.Context, id string) (*domain.Workshop, error) {
	// TODO: Implementar
	return nil, nil
}

// Update implements WorkshopRepository.Update
func (r *WorkshopRepository) Update(ctx context.Context, workshop *domain.Workshop) error {
	// TODO: Implementar
	return nil
}

// Delete implements WorkshopRepository.Delete
func (r *WorkshopRepository) Delete(ctx context.Context, id string) error {
	// TODO: Implementar
	return nil
}

// List implements WorkshopRepository.List
func (r *WorkshopRepository) List(ctx context.Context, limit, offset int) ([]*domain.Workshop, int64, error) {
	// TODO: Implementar
	return nil, 0, nil
}

// SearchByName implements WorkshopRepository.SearchByName
func (r *WorkshopRepository) SearchByName(ctx context.Context, name string, limit int) ([]*domain.Workshop, error) {
	// TODO: Implementar
	return nil, nil
}
