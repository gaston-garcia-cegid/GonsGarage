package car

import (
	"context"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/google/uuid"
)

// carService implements the CarService interface
type carService struct {
	carRepo  ports.CarRepository
	userRepo ports.UserRepository
	logger   ports.Logger
}

// NewCarService creates a new car service instance
func NewCarServiceTest(
	carRepo ports.CarRepository,
	userRepo ports.UserRepository,
) ports.CarService {
	return &carService{
		carRepo:  carRepo,
		userRepo: userRepo,
	}
}

// CreateCar creates a new car with proper business rules and authorization
func (s *carService) CreateCar(ctx context.Context, car *domain.Car, requestingUserID uuid.UUID) (*domain.Car, error) {
	// Get requesting user to validate permissions
	requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		s.logger.Error("failed to get requesting user", "userID", requestingUserID, "error", err)
		return nil, fmt.Errorf("failed to validate user: %w", err)
	}

	// Apply business rules based on user role
	if err := s.validateCarCreationPermissions(requestingUser, car); err != nil {
		return nil, err
	}

	// Validate car owner exists and is a client
	if err := s.validateCarOwner(ctx, car.OwnerID); err != nil {
		return nil, err
	}

	// Check for duplicate license plate
	if err := s.checkDuplicateLicensePlate(ctx, car.LicensePlate); err != nil {
		return nil, err
	}

	// Validate car domain rules
	if err := car.Validate(); err != nil {
		return nil, fmt.Errorf("car validation failed: %w", err)
	}

	// Set car metadata
	car.ID = uuid.New()
	car.CreatedAt = time.Now()
	car.UpdatedAt = time.Now()

	// Persist car
	if err := s.carRepo.Create(ctx, car); err != nil {
		s.logger.Error("failed to create car", "carID", car.ID, "error", err)
		return nil, fmt.Errorf("failed to create car: %w", err)
	}

	s.logger.Info("car created successfully", "carID", car.ID, "ownerID", car.OwnerID)
	return car, nil
}

// GetCar retrieves a car with authorization checks
func (s *carService) GetCar(ctx context.Context, carID uuid.UUID, requestingUserID uuid.UUID) (*domain.Car, error) {
	// Get requesting user
	requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requesting user: %w", err)
	}

	// Get car
	car, err := s.carRepo.GetByID(ctx, carID)
	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	if car == nil {
		return nil, domain.ErrCarNotFound
	}

	// Check authorization
	if err := s.validateCarAccess(requestingUser, car); err != nil {
		return nil, err
	}

	return car, nil
}

// GetCarsByOwner retrieves cars for a specific owner with authorization
func (s *carService) GetCarsByOwner(ctx context.Context, ownerID uuid.UUID, requestingUserID uuid.UUID) ([]*domain.Car, error) {
	// Get requesting user
	requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requesting user: %w", err)
	}

	// Validate access permissions
	if err := s.validateOwnerAccess(requestingUser, ownerID, requestingUserID); err != nil {
		return nil, err
	}

	// Get cars
	cars, err := s.carRepo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		s.logger.Error("failed to get cars by owner", "ownerID", ownerID, "error", err)
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}

	return cars, nil
}

// UpdateCar updates an existing car with authorization
func (s *carService) UpdateCar(ctx context.Context, car *domain.Car, requestingUserID uuid.UUID) (*domain.Car, error) {
	// Get requesting user
	requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requesting user: %w", err)
	}

	// Get existing car
	existingCar, err := s.carRepo.GetByID(ctx, car.ID)
	if err != nil || existingCar == nil {
		return nil, domain.ErrCarNotFound
	}

	// Check authorization
	if err := s.validateCarAccess(requestingUser, existingCar); err != nil {
		return nil, err
	}

	// Validate updated data
	if err := car.Validate(); err != nil {
		return nil, fmt.Errorf("car validation failed: %w", err)
	}

	// Preserve immutable fields
	car.ID = existingCar.ID
	car.OwnerID = existingCar.OwnerID
	car.CreatedAt = existingCar.CreatedAt
	car.UpdatedAt = time.Now()

	// Update car
	if err := s.carRepo.Update(ctx, car); err != nil {
		s.logger.Error("failed to update car", "carID", car.ID, "error", err)
		return nil, fmt.Errorf("failed to update car: %w", err)
	}

	s.logger.Info("car updated successfully", "carID", car.ID)
	return car, nil
}

// DeleteCar removes a car with authorization
func (s *carService) DeleteCar(ctx context.Context, carID uuid.UUID, requestingUserID uuid.UUID) error {
	// Get requesting user
	requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return fmt.Errorf("failed to get requesting user: %w", err)
	}

	// Get car
	car, err := s.carRepo.GetByID(ctx, carID)
	if err != nil || car == nil {
		return domain.ErrCarNotFound
	}

	// Check authorization
	if err := s.validateCarAccess(requestingUser, car); err != nil {
		return err
	}

	// Delete car
	if err := s.carRepo.Delete(ctx, carID); err != nil {
		s.logger.Error("failed to delete car", "carID", carID, "error", err)
		return fmt.Errorf("failed to delete car: %w", err)
	}

	s.logger.Info("car deleted successfully", "carID", carID)
	return nil
}

// GetCarWithRepairs retrieves a car with its repair history
func (s *carService) GetCarWithRepairs(ctx context.Context, carID uuid.UUID, requestingUserID uuid.UUID) (*domain.Car, error) {
	// Get requesting user
	requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requesting user: %w", err)
	}

	// Get car with repairs
	car, err := s.carRepo.GetWithRepairs(ctx, carID)
	if err != nil || car == nil {
		return nil, domain.ErrCarNotFound
	}

	// Check authorization
	if err := s.validateCarAccess(requestingUser, car); err != nil {
		return nil, err
	}

	return car, nil
}

// Helper methods for validation (following SRP)

func (s *carService) validateCarCreationPermissions(user *domain.User, car *domain.Car) error {
	if user.IsClient() {
		// Clients can only create cars for themselves
		car.OwnerID = user.ID
		return nil
	}

	if !user.CanManageUsers() {
		return domain.ErrUnauthorizedAccess
	}

	return nil
}

func (s *carService) validateCarOwner(ctx context.Context, ownerID uuid.UUID) error {
	owner, err := s.userRepo.GetByID(ctx, ownerID)
	if err != nil {
		return fmt.Errorf("car owner not found: %w", err)
	}

	if !owner.IsClient() {
		return fmt.Errorf("car owner must be a client")
	}

	return nil
}

func (s *carService) checkDuplicateLicensePlate(ctx context.Context, licensePlate string) error {
	existingCar, err := s.carRepo.GetByLicensePlate(ctx, licensePlate)
	if err != nil {
		return nil // Error means not found, which is OK
	}

	if existingCar != nil {
		return domain.ErrCarAlreadyExists
	}

	return nil
}

func (s *carService) validateCarAccess(user *domain.User, car *domain.Car) error {
	if user.IsClient() && !car.IsOwnedBy(user.ID) {
		return domain.ErrUnauthorizedAccess
	}

	if !user.IsClient() && !user.CanManageUsers() {
		return domain.ErrUnauthorizedAccess
	}

	return nil
}

func (s *carService) validateOwnerAccess(user *domain.User, ownerID uuid.UUID, requestingUserID uuid.UUID) error {
	if user.IsClient() && ownerID != requestingUserID {
		return domain.ErrUnauthorizedAccess
	}

	if !user.IsClient() && !user.CanManageUsers() {
		return domain.ErrUnauthorizedAccess
	}

	return nil
}
