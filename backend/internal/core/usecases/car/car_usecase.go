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
	logger    ports.Logger
	cacheRepo ports.CacheRepository
}

func NewCarUseCase(
	carRepo ports.CarRepository,
	userRepo ports.UserRepository,
	cacheRepo ports.CacheRepository) *CarUseCase {
	return &CarUseCase{
		carRepo:   carRepo,
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
	}
}

// CreateCar creates a new car for a client
func (uc *CarUseCase) CreateCar(ctx context.Context, car *domain.Car, requestingUserID uuid.UUID) (*domain.Car, error) {
	// Get the requesting user to check permissions
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		uc.logger.Error("failed to get requesting user", "user_id", requestingUserID, "error", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Validate permissions: clients can only create cars for themselves
	// Admins and managers can create cars for any client
	if requestingUser.IsClient() {
		car.OwnerID = requestingUserID // Force owner to be the requesting client
	} else if !requestingUser.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Validate that the owner exists and is a client
	owner, err := uc.userRepo.GetByID(ctx, car.OwnerID)
	if err != nil {
		uc.logger.Error("failed to get car owner", "owner_id", car.OwnerID, "error", err)
		return nil, fmt.Errorf("car owner not found: %w", err)
	}

	if !owner.IsClient() {
		return nil, fmt.Errorf("car owner must be a client")
	}

	// Check if license plate already exists
	existingCar, err := uc.carRepo.GetByLicensePlate(ctx, car.LicensePlate)
	if err == nil && existingCar != nil {
		return nil, domain.ErrCarAlreadyExists
	}

	// Validate car data
	if err := car.Validate(); err != nil {
		return nil, fmt.Errorf("invalid car data: %w", err)
	}

	// Set metadata
	car.ID = uuid.New()
	car.CreatedAt = time.Now()
	car.UpdatedAt = time.Now()

	// Create the car
	if err := uc.carRepo.Create(ctx, car); err != nil {
		uc.logger.Error("failed to create car", "car", car, "error", err)
		return nil, fmt.Errorf("failed to create car: %w", err)
	}

	uc.logger.Info("car created successfully", "car_id", car.ID, "owner_id", car.OwnerID)
	return car, nil
}

// GetCar retrieves a car by ID with permission checks
func (uc *CarUseCase) GetCar(ctx context.Context, carID uuid.UUID, requestingUserID uuid.UUID) (*domain.Car, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get the car
	car, err := uc.carRepo.GetByID(ctx, carID)
	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	if car == nil {
		return nil, domain.ErrCarNotFound
	}

	// Check permissions: clients can only see their own cars
	if requestingUser.IsClient() && !car.IsOwnedBy(requestingUserID) {
		return nil, domain.ErrUnauthorizedAccess
	}

	return car, nil
}

// GetCarsByOwner retrieves all cars owned by a specific user
func (uc *CarUseCase) GetCarsByOwner(ctx context.Context, ownerID uuid.UUID, requestingUserID uuid.UUID) ([]*domain.Car, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check permissions: clients can only see their own cars
	if requestingUser.IsClient() && ownerID != requestingUserID {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Admins and managers can see any user's cars
	if !requestingUser.IsClient() && !requestingUser.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}

	cars, err := uc.carRepo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		uc.logger.Error("failed to get cars by owner", "owner_id", ownerID, "error", err)
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}

	return cars, nil
}

// UpdateCar updates an existing car
func (uc *CarUseCase) UpdateCar(ctx context.Context, car *domain.Car, requestingUserID uuid.UUID) (*domain.Car, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get the existing car
	existingCar, err := uc.carRepo.GetByID(ctx, car.ID)
	if err != nil || existingCar == nil {
		return nil, domain.ErrCarNotFound
	}

	// Check permissions: clients can only update their own cars
	if requestingUser.IsClient() && !existingCar.IsOwnedBy(requestingUserID) {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Admins and managers can update any car
	if !requestingUser.IsClient() && !requestingUser.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Validate car data
	if err := car.Validate(); err != nil {
		return nil, fmt.Errorf("invalid car data: %w", err)
	}

	// Preserve some fields
	car.ID = existingCar.ID
	car.OwnerID = existingCar.OwnerID
	car.CreatedAt = existingCar.CreatedAt
	car.UpdatedAt = time.Now()

	// Update the car
	if err := uc.carRepo.Update(ctx, car); err != nil {
		uc.logger.Error("failed to update car", "car_id", car.ID, "error", err)
		return nil, fmt.Errorf("failed to update car: %w", err)
	}

	uc.logger.Info("car updated successfully", "car_id", car.ID)
	return car, nil
}

// DeleteCar deletes a car (soft delete)
func (uc *CarUseCase) DeleteCar(ctx context.Context, carID uuid.UUID, requestingUserID uuid.UUID) error {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Get the car
	car, err := uc.carRepo.GetByID(ctx, carID)
	if err != nil || car == nil {
		return domain.ErrCarNotFound
	}

	// Check permissions: clients can only delete their own cars
	if requestingUser.IsClient() && !car.IsOwnedBy(requestingUserID) {
		return domain.ErrUnauthorizedAccess
	}

	// Admins and managers can delete any car
	if !requestingUser.IsClient() && !requestingUser.CanManageUsers() {
		return domain.ErrUnauthorizedAccess
	}

	// Delete the car
	if err := uc.carRepo.Delete(ctx, carID); err != nil {
		uc.logger.Error("failed to delete car", "car_id", carID, "error", err)
		return fmt.Errorf("failed to delete car: %w", err)
	}

	uc.logger.Info("car deleted successfully", "car_id", carID)
	return nil
}

// GetCarWithRepairs retrieves a car with its repair history
func (uc *CarUseCase) GetCarWithRepairs(ctx context.Context, carID uuid.UUID, requestingUserID uuid.UUID) (*domain.Car, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get the car with repairs
	car, err := uc.carRepo.GetWithRepairs(ctx, carID)
	if err != nil || car == nil {
		return nil, domain.ErrCarNotFound
	}

	// Check permissions: clients can only see their own cars
	if requestingUser.IsClient() && !car.IsOwnedBy(requestingUserID) {
		return nil, domain.ErrUnauthorizedAccess
	}

	return car, nil
}
