package postgres

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"gorm.io/gorm"
)

// AccountingRepository is a PostgreSQL implementation of the AccountingRepository interface
type AccountingRepository struct {
	db *gorm.DB
}

// UpdateAccountingEntry implements ports.AccountingRepository.
func (r *AccountingRepository) UpdateAccountingEntry(ctx context.Context, entry *domain.AccountingEntry) error {
	panic("unimplemented")
}

// NewAccountingRepository creates a new AccountingRepository
func NewAccountingRepository(db *gorm.DB) ports.AccountingRepository {
	return &AccountingRepository{db: db}
}

// CreateAccountingEntry implements AccountingRepository.CreateAccountingEntry
func (r *AccountingRepository) CreateAccountingEntry(ctx context.Context, entry *domain.AccountingEntry) error {
	// TODO: Implementar
	return nil
}

// Create implements AccountingRepository.Create
func (r *AccountingRepository) Create(ctx context.Context, entry *domain.AccountingEntry) error {
	// TODO: Implementar
	return nil
}

// GetByID implements AccountingRepository.GetByID
func (r *AccountingRepository) GetByID(ctx context.Context, id string) (*domain.AccountingEntry, error) {
	// TODO: Implementar
	return nil, nil
}

// Update implements AccountingRepository.Update
func (r *AccountingRepository) Update(ctx context.Context, entry *domain.AccountingEntry) error {
	// TODO: Implementar
	return nil
}

// Delete implements AccountingRepository.Delete
func (r *AccountingRepository) Delete(ctx context.Context, id string) error {
	// TODO: Implementar
	return nil
}

// List implements AccountingRepository.List
func (r *AccountingRepository) List(ctx context.Context, limit, offset int) ([]*domain.AccountingEntry, int64, error) {
	// TODO: Implementar
	return nil, 0, nil
}

// SearchByName implements AccountingRepository.SearchByName
func (r *AccountingRepository) SearchByName(ctx context.Context, name string, limit int) ([]*domain.AccountingEntry, error) {
	// TODO: Implementar
	return nil, nil
}

// DeleteAccountingEntry implements AccountingRepository.DeleteAccountingEntry
func (r *AccountingRepository) DeleteAccountingEntry(ctx context.Context, id string) error {
	// TODO: Implementar
	return nil
}

// GetAccountingEntryByID implements AccountingRepository.GetAccountingEntryByID
func (r *AccountingRepository) GetAccountingEntryByID(ctx context.Context, id string) (*domain.AccountingEntry, error) {
	// TODO: Implementar
	return nil, nil
}

// ListAccountingEntries implements AccountingRepository.ListAccountingEntries
func (r *AccountingRepository) ListAccountingEntries(ctx context.Context, limit, offset int) ([]*domain.AccountingEntry, int64, error) {
	// TODO: Implementar
	return nil, 0, nil
}
