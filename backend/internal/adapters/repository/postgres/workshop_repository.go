package postgres

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"gorm.io/gorm"
)

// WorkshopRepository is a PostgreSQL implementation of the WorkshopRepository interface
type WorkshopRepository struct {
	db *gorm.DB
}

// NewWorkshopRepository creates a new WorkshopRepository
func NewWorkshopRepository(db *gorm.DB) ports.WorkshopRepository {
	return &WorkshopRepository{db: db}
}

// Create implements WorkshopRepository.Create
func (r *WorkshopRepository) CreateWorkshop(ctx context.Context, workshop *domain.Workshop) error {
	// TODO: Implementar
	return nil
}

// GetByID implements WorkshopRepository.GetByID
func (r *WorkshopRepository) GetWorkshopByID(ctx context.Context, id string) (*domain.Workshop, error) {
	// TODO: Implementar
	return nil, nil
}

// Update implements WorkshopRepository.Update
func (r *WorkshopRepository) UpdateWorkshop(ctx context.Context, workshop *domain.Workshop) error {
	// TODO: Implementar
	return nil
}

// Delete implements WorkshopRepository.Delete
func (r *WorkshopRepository) DeleteWorkshop(ctx context.Context, id string) error {
	// TODO: Implementar
	return nil
}

// List implements WorkshopRepository.List
func (r *WorkshopRepository) ListWorkshops(ctx context.Context, limit, offset int) ([]*domain.Workshop, int64, error) {
	// TODO: Implementar
	return nil, 0, nil
}

// SearchByName implements WorkshopRepository.SearchByName
func (r *WorkshopRepository) SearchByName(ctx context.Context, name string, limit int) ([]*domain.Workshop, error) {
	// TODO: Implementar
	return nil, nil
}
