package car

import (
	"context"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/google/uuid"
)

type CarUseCase struct {
	carRepo   ports.CarRepository
	userRepo  ports.UserRepository
	cacheRepo ports.CacheRepository
}

func NewCarUseCase(carRepo ports.CarRepository, userRepo ports.UserRepository, cacheRepo ports.CacheRepository) *CarUseCase {
	return &CarUseCase{
		carRepo:   carRepo,
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
	}
}

func (uc *CarUseCase) CreateCar(ctx context.Context, car *domain.Car, userID uuid.UUID) (*domain.Car, error) {
	// Get user to check permissions
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check permissions: only client can create their own car, or admin/manager
	if user.IsClient() && car.OwnerID != userID {
		return nil, domain.ErrUnauthorizedAccess
	}

	if !user.IsClient() && !user.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Validate client exists
	client, err := uc.userRepo.GetByID(ctx, car.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	if !client.IsClient() {
		return nil, fmt.Errorf("user is not a client")
	}

	// Check if license plate already exists
	existingCar, err := uc.carRepo.GetByLicensePlate(ctx, car.LicensePlate)
	if err == nil && existingCar != nil {
		return nil, fmt.Errorf("car with license plate %s already exists", car.LicensePlate)
	}

	// Set metadata
	car.ID = uuid.New()
	car.CreatedAt = time.Now()
	car.UpdatedAt = time.Now()

	// Validate car data
	if err := uc.validateCar(car); err != nil {
		return nil, err
	}

	// Create car
	if err := uc.carRepo.Create(ctx, car); err != nil {
		return nil, fmt.Errorf("failed to create car: %w", err)
	}

	return car, nil
}

func (uc *CarUseCase) GetCar(ctx context.Context, carID uuid.UUID, userID uuid.UUID) (*domain.Car, error) {
	// Get user to check permissions
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	car, err := uc.carRepo.GetByID(ctx, carID)
	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	// Check permissions: clients can only see their own cars
	if user.IsClient() && car.OwnerID != userID {
		return nil, domain.ErrUnauthorizedAccess
	}

	return car, nil
}

func (uc *CarUseCase) ListCars(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Car, error) {
	// Get user to check permissions
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// If client, only return their cars
	if user.IsClient() {
		return uc.carRepo.GetByClientID(ctx, userID)
	}

	// Employees, managers, and admins can see all cars
	return uc.carRepo.List(ctx, limit, offset)
}

func (uc *CarUseCase) UpdateCar(ctx context.Context, car *domain.Car, userID uuid.UUID) (*domain.Car, error) {
	// Get user to check permissions
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get existing car
	existingCar, err := uc.carRepo.GetByID(ctx, car.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing car: %w", err)
	}

	// Check permissions
	if user.IsClient() && existingCar.OwnerID != userID {
		return nil, domain.ErrUnauthorizedAccess
	}

	if !user.IsClient() && !user.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Validate car data
	if err := uc.validateCar(car); err != nil {
		return nil, err
	}

	// Update metadata
	car.UpdatedAt = time.Now()
	car.CreatedAt = existingCar.CreatedAt // Preserve original creation time

	// Update car
	if err := uc.carRepo.Update(ctx, car); err != nil {
		return nil, fmt.Errorf("failed to update car: %w", err)
	}

	return car, nil
}

func (uc *CarUseCase) DeleteCar(ctx context.Context, carID uuid.UUID, userID uuid.UUID) error {
	// Get user to check permissions
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Get car to check ownership
	car, err := uc.carRepo.GetByID(ctx, carID)
	if err != nil {
		return fmt.Errorf("failed to get car: %w", err)
	}

	// Check permissions: only admin/manager can delete cars
	if !user.CanManageUsers() {
		return domain.ErrUnauthorizedAccess
	}

	if car == nil {
		return domain.ErrCarNotFound
	}

	// Delete car (soft delete)
	if err := uc.carRepo.Delete(ctx, carID); err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	return nil
}

func (uc *CarUseCase) validateCar(car *domain.Car) error {
	if car.Make == "" {
		return fmt.Errorf("make is required")
	}
	if car.Model == "" {
		return fmt.Errorf("model is required")
	}
	if car.Year < 1900 || car.Year > time.Now().Year()+1 {
		return fmt.Errorf("invalid year")
	}
	if car.LicensePlate == "" {
		return fmt.Errorf("license plate is required")
	}
	if car.Color == "" {
		return fmt.Errorf("color is required")
	}

	return nil
}
