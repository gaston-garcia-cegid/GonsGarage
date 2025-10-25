package repair

import (
	"context"
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
	// Get user to check permissions
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Only employees can create repairs
	if !user.IsEmployee() {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Validate car exists
	car, err := uc.carRepo.GetByID(ctx, repair.CarID)
	if err != nil {
		return nil, fmt.Errorf("car not found: %w", err)
	}

	// Check if repair already exists for the same car with same description and start date
	existingRepairs, err := uc.repairRepo.GetByCarID(ctx, repair.CarID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing repairs: %w", err)
	}
	for _, r := range existingRepairs {
		if r.Description == repair.Description && r.StartedAt == repair.StartedAt {
			return nil, fmt.Errorf("repair already exists for the same car with same description and start date")
		}
	}

	if car.OwnerID != userID {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Set metadata
	repair.ID = uuid.New()
	repair.TechnicianID = userID
	repair.CreatedAt = time.Now()
	repair.UpdatedAt = time.Now()

	if repair.Status == "" {
		repair.Status = domain.RepairStatusPending
	}

	// Validate repair data
	if err := uc.validateRepair(repair); err != nil {
		return nil, err
	}

	// Create repair
	if err := uc.repairRepo.Create(ctx, repair); err != nil {
		return nil, fmt.Errorf("failed to create repair: %w", err)
	}

	return repair, nil
}

func (uc *RepairService) GetRepair(ctx context.Context, repairID uuid.UUID, userID uuid.UUID) (*domain.Repair, error) {
	// Get user to check permissions
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	repair, err := uc.repairRepo.GetByID(ctx, repairID)
	if err != nil {
		return nil, fmt.Errorf("failed to get repair: %w", err)
	}

	// Check permissions: clients can only see repairs for their cars
	if user.IsClient() {
		car, err := uc.carRepo.GetByID(ctx, repair.CarID)
		if err != nil {
			return nil, fmt.Errorf("failed to get car: %w", err)
		}
		if car.OwnerID != userID {
			return nil, domain.ErrUnauthorizedAccess
		}
	}

	return repair, nil
}

func (uc *RepairService) GetRepairsByCarID(ctx context.Context, carID uuid.UUID, userID uuid.UUID) ([]*domain.Repair, error) {
	// Get user to check permissions
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if car exists and user has access
	car, err := uc.carRepo.GetByID(ctx, carID)
	if err != nil {
		return nil, fmt.Errorf("car not found: %w", err)
	}

	// Check permissions: clients can only see repairs for their cars
	if user.IsClient() && car.OwnerID != userID {
		return nil, domain.ErrUnauthorizedAccess
	}

	return uc.repairRepo.GetByCarID(ctx, carID)
}

func (uc *RepairService) UpdateRepair(ctx context.Context, repair *domain.Repair, userID uuid.UUID) (*domain.Repair, error) {
	// Get user to check permissions
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Only employees can update repairs
	if !user.IsEmployee() {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Get existing repair
	existingRepair, err := uc.repairRepo.GetByID(ctx, repair.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing repair: %w", err)
	}

	// Validate repair data
	if err := uc.validateRepair(repair); err != nil {
		return nil, err
	}

	// Update metadata
	repair.UpdatedAt = time.Now()
	repair.CreatedAt = existingRepair.CreatedAt // Preserve original creation time

	// If marking as completed, set end date
	if repair.Status == domain.RepairStatusCompleted && repair.CompletedAt == nil {
		now := time.Now()
		repair.CompletedAt = &now
	}

	// Update repair
	if err := uc.repairRepo.Update(ctx, repair); err != nil {
		return nil, fmt.Errorf("failed to update repair: %w", err)
	}

	return repair, nil
}

func (uc *RepairService) validateRepair(repair *domain.Repair) error {
	if repair.Description == "" {
		return fmt.Errorf("description is required")
	}
	if !domain.ValidateRepairStatus(repair.Status) {
		return fmt.Errorf("invalid repair status")
	}
	if repair.Cost < 0 {
		return fmt.Errorf("cost cannot be negative")
	}

	return nil
}
