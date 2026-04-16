package repair

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/google/uuid"
)

type RepairService struct {
	repairRepo ports.RepairRepository
	carRepo    ports.CarRepository
	userRepo   ports.UserRepository
}

func NewRepairService(repairRepo ports.RepairRepository, carRepo ports.CarRepository, userRepo ports.UserRepository) *RepairService {
	return &RepairService{
		repairRepo: repairRepo,
		carRepo:    carRepo,
		userRepo:   userRepo,
	}
}

func (uc *RepairService) CreateRepair(ctx context.Context, repair *domain.Repair, userID uuid.UUID) (*domain.Repair, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	if !user.IsEmployee() {
		return nil, domain.ErrUnauthorizedAccess
	}

	car, err := uc.carRepo.GetByID(ctx, repair.CarID)
	if err != nil {
		if errors.Is(err, domain.ErrCarNotFound) {
			return nil, domain.ErrCarNotFound
		}
		return nil, fmt.Errorf("failed to get car: %w", err)
	}
	if car == nil {
		return nil, domain.ErrCarNotFound
	}

	repair.ID = uuid.New()
	repair.TechnicianID = userID
	repair.CreatedAt = time.Now().UTC()
	repair.UpdatedAt = repair.CreatedAt

	if repair.Status == "" {
		repair.Status = domain.RepairStatusPending
	}

	if err := uc.validateRepair(repair); err != nil {
		return nil, err
	}

	if err := uc.repairRepo.Create(ctx, repair); err != nil {
		return nil, fmt.Errorf("failed to create repair: %w", err)
	}

	return repair, nil
}

func (uc *RepairService) GetRepair(ctx context.Context, repairID uuid.UUID, userID uuid.UUID) (*domain.Repair, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	repair, err := uc.repairRepo.GetByID(ctx, repairID)
	if err != nil {
		if errors.Is(err, domain.ErrRepairNotFound) {
			return nil, domain.ErrRepairNotFound
		}
		return nil, err
	}
	if repair == nil {
		return nil, domain.ErrRepairNotFound
	}

	if user.IsClient() {
		car, err := uc.carRepo.GetByID(ctx, repair.CarID)
		if err != nil {
			return nil, fmt.Errorf("failed to get car: %w", err)
		}
		if car == nil || car.OwnerID != userID {
			return nil, domain.ErrUnauthorizedAccess
		}
	}

	return repair, nil
}

func (uc *RepairService) GetRepairsByCarID(ctx context.Context, carID uuid.UUID, userID uuid.UUID) ([]*domain.Repair, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	car, err := uc.carRepo.GetByID(ctx, carID)
	if err != nil {
		if errors.Is(err, domain.ErrCarNotFound) {
			return nil, domain.ErrCarNotFound
		}
		return nil, fmt.Errorf("failed to get car: %w", err)
	}
	if car == nil {
		return nil, domain.ErrCarNotFound
	}

	if user.IsClient() && car.OwnerID != userID {
		return nil, domain.ErrUnauthorizedAccess
	}

	return uc.repairRepo.GetByCarID(ctx, carID)
}

func (uc *RepairService) UpdateRepair(ctx context.Context, repair *domain.Repair, userID uuid.UUID) (*domain.Repair, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	if !user.IsEmployee() {
		return nil, domain.ErrUnauthorizedAccess
	}

	existing, err := uc.repairRepo.GetByID(ctx, repair.ID)
	if err != nil {
		if errors.Is(err, domain.ErrRepairNotFound) {
			return nil, domain.ErrRepairNotFound
		}
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrRepairNotFound
	}

	merged := *existing
	if repair.Description != "" {
		merged.Description = repair.Description
	}
	if repair.Status != "" {
		merged.Status = repair.Status
	}
	merged.Cost = repair.Cost
	if repair.StartedAt != nil {
		merged.StartedAt = repair.StartedAt
	}
	if repair.CompletedAt != nil {
		merged.CompletedAt = repair.CompletedAt
	}

	if merged.Status == domain.RepairStatusCompleted && merged.CompletedAt == nil {
		now := time.Now().UTC()
		merged.CompletedAt = &now
	}

	merged.UpdatedAt = time.Now().UTC()

	if err := uc.validateRepair(&merged); err != nil {
		return nil, err
	}

	if err := uc.repairRepo.Update(ctx, &merged); err != nil {
		return nil, fmt.Errorf("failed to update repair: %w", err)
	}

	return &merged, nil
}

func (uc *RepairService) validateRepair(repair *domain.Repair) error {
	if repair.Description == "" {
		return domain.ErrInvalidRepairData
	}
	if !domain.ValidateRepairStatus(repair.Status) {
		return domain.ErrInvalidRepairData
	}
	if repair.Cost < 0 {
		return domain.ErrInvalidRepairData
	}
	return nil
}
