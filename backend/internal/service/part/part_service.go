package part

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/google/uuid"
)

// PartService implements ports.PartService (manager/admin inventory; business rules on top of repo).
type PartService struct {
	repo     ports.PartItemRepository
	userRepo ports.UserRepository
}

func NewPartService(repo ports.PartItemRepository, userRepo ports.UserRepository) *PartService {
	return &PartService{repo: repo, userRepo: userRepo}
}

var _ ports.PartService = (*PartService)(nil)

func (s *PartService) requireManagerOrAdmin(ctx context.Context, requestingUserID uuid.UUID) (*domain.User, error) {
	u, err := s.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if u == nil {
		return nil, domain.ErrUserNotFound
	}
	if !u.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}
	return u, nil
}

func (s *PartService) ensureNoDuplicateBarcode(ctx context.Context, barcode string, excludeID uuid.UUID) error {
	b := strings.TrimSpace(barcode)
	if b == "" {
		return nil
	}
	existing, err := s.repo.GetByBarcode(ctx, b)
	if err != nil {
		if errors.Is(err, domain.ErrPartItemNotFound) {
			return nil
		}
		return err
	}
	if existing != nil && existing.ID != excludeID {
		return domain.ErrPartItemDuplicateBarcode
	}
	return nil
}

// Create persists a part item (manager/admin only).
func (s *PartService) Create(ctx context.Context, item *domain.PartItem, requestingUserID uuid.UUID) (*domain.PartItem, error) {
	if _, err := s.requireManagerOrAdmin(ctx, requestingUserID); err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("part item is required")
	}
	if err := item.Validate(); err != nil {
		return nil, err
	}
	if err := s.ensureNoDuplicateBarcode(ctx, item.Barcode, uuid.Nil); err != nil {
		return nil, err
	}
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, item.ID)
}

// Get returns a part item by id (manager/admin only).
func (s *PartService) Get(ctx context.Context, id uuid.UUID, requestingUserID uuid.UUID) (*domain.PartItem, error) {
	if _, err := s.requireManagerOrAdmin(ctx, requestingUserID); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

// List returns part items (manager/admin only).
func (s *PartService) List(ctx context.Context, filters ports.PartItemListFilters, requestingUserID uuid.UUID) ([]*domain.PartItem, int64, error) {
	if _, err := s.requireManagerOrAdmin(ctx, requestingUserID); err != nil {
		return nil, 0, err
	}
	return s.repo.List(ctx, filters)
}

// Update updates a part item (manager/admin only).
func (s *PartService) Update(ctx context.Context, item *domain.PartItem, requestingUserID uuid.UUID) (*domain.PartItem, error) {
	if _, err := s.requireManagerOrAdmin(ctx, requestingUserID); err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("part item is required")
	}
	if err := item.Validate(); err != nil {
		return nil, err
	}
	if err := s.ensureNoDuplicateBarcode(ctx, item.Barcode, item.ID); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, item.ID)
}

// Delete soft-deletes a part item (manager/admin only).
func (s *PartService) Delete(ctx context.Context, id uuid.UUID, requestingUserID uuid.UUID) error {
	if _, err := s.requireManagerOrAdmin(ctx, requestingUserID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
