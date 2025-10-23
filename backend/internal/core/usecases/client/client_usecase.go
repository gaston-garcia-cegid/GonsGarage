package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/google/uuid"
)

type ClientUseCase struct {
	userRepo   ports.UserRepository
	carRepo    ports.CarRepository
	repairRepo ports.RepairRepository
}

func NewClientUseCase(
	userRepo ports.UserRepository,
	carRepo ports.CarRepository,
	repairRepo ports.RepairRepository) *ClientUseCase {
	return &ClientUseCase{
		userRepo:   userRepo,
		carRepo:    carRepo,
		repairRepo: repairRepo,
	}
}
func (uc *ClientUseCase) CreateClient(ctx context.Context, client *domain.User) ([]*domain.User, error) {
	// Validate that the user is being created as a client
	if client.Role != "client" {
		return nil, fmt.Errorf("invalid role: expected 'client', got '%s'", client.Role)
	}

	// Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, client.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Validate client data (following Agent.md validation rules)
	if err := client.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client data: %w", err.(error))
	}

	// Set metadata
	client.ID = uuid.New()
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()
	client.IsActive = true

	// Create the client
	if err := uc.userRepo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return []*domain.User{client}, nil
}
func (uc *ClientUseCase) ListClients(ctx context.Context) ([]*domain.User, error) {
	return uc.userRepo.List(ctx, 0, 0)
}
func (uc *ClientUseCase) GetClient(ctx context.Context, clientID uuid.UUID, requestingUserID uuid.UUID) (*domain.User, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get the client
	client, err := uc.userRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	if client == nil || !client.IsClient() {
		return nil, domain.ErrUserNotFound
	}

	// Check permissions: clients can only see their own profile (Agent.md security rules)
	if requestingUser.IsClient() && clientID != requestingUserID {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Admins and managers can see any client
	if !requestingUser.IsClient() && !requestingUser.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}

	return client, nil
}

func (uc *ClientUseCase) UpdateClient(ctx context.Context, client *domain.User, requestingUserID uuid.UUID) (*domain.User, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get the existing client
	existingClient, err := uc.userRepo.GetByID(ctx, client.ID)
	if err != nil || existingClient == nil || !existingClient.IsClient() {
		return nil, domain.ErrUserNotFound
	}

	// Check permissions
	if requestingUser.IsClient() && client.ID != requestingUserID {
		return nil, domain.ErrUnauthorizedAccess
	}

	if !requestingUser.IsClient() && !requestingUser.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Validate client data
	if err := client.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client data: %w", err.(error))
	}

	// Preserve some fields (following Agent.md data integrity rules)
	client.ID = existingClient.ID
	client.Role = "client" // Ensure role remains client
	client.CreatedAt = existingClient.CreatedAt
	client.UpdatedAt = time.Now()

	// Update the client
	if err := uc.userRepo.Update(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	return client, nil
}

// UpdateClientProfile allows clients to update their own profile
func (uc *ClientUseCase) UpdateClientProfile(ctx context.Context, clientID uuid.UUID, client *domain.User) (*domain.User, error) {
	// Get the existing client
	existingClient, err := uc.userRepo.GetByID(ctx, clientID)
	if err != nil || existingClient == nil || !existingClient.IsClient() {
		return nil, domain.ErrUserNotFound
	}

	// Validate client data
	if err := client.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client data: %w", err.(error))
	}

	// Preserve critical fields (clients can't change their own role, etc.)
	client.ID = clientID
	client.Role = "client"
	client.CreatedAt = existingClient.CreatedAt
	client.UpdatedAt = time.Now()
	client.IsActive = existingClient.IsActive // Clients can't activate/deactivate themselves

	// Update the client
	if err := uc.userRepo.Update(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	return client, nil
}

// DeleteClient deletes a client (soft delete) - ✅ Fixed: using uuid.UUID parameter
func (uc *ClientUseCase) DeleteClient(ctx context.Context, clientID uuid.UUID) error {
	// Get the client
	client, err := uc.userRepo.GetByID(ctx, clientID)
	if err != nil || client == nil || !client.IsClient() {
		return domain.ErrUserNotFound
	}

	// Note: Only admins/managers should be able to call this method
	// The authorization should be handled at the handler/service layer

	// Delete the client (soft delete)
	if err := uc.userRepo.Delete(ctx, clientID); err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	return nil
}

// ListClients retrieves all clients (admin/manager only)
func (uc *ClientUseCase) ListClientsWithAuth(ctx context.Context, requestingUserID uuid.UUID, limit, offset int) ([]*domain.User, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check permissions: only admins and managers can list all clients
	if !requestingUser.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Get clients by role
	clients, err := uc.userRepo.GetByRole(ctx, "client", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}

	return clients, nil
}

// GetClientCars retrieves all cars owned by a specific client
func (uc *ClientUseCase) GetClientCars(ctx context.Context, clientID uuid.UUID, requestingUserID uuid.UUID) ([]*domain.Car, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify the target user is a client
	client, err := uc.userRepo.GetByID(ctx, clientID)
	if err != nil || client == nil || !client.IsClient() {
		return nil, domain.ErrUserNotFound
	}

	// Check permissions: clients can only see their own cars
	if requestingUser.IsClient() && clientID != requestingUserID {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Admins and managers can see any client's cars
	if !requestingUser.IsClient() && !requestingUser.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}

	cars, err := uc.carRepo.GetByOwnerID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client cars: %w", err)
	}

	return cars, nil
}

// GetClientRepairs retrieves all repairs for a specific client's cars
func (uc *ClientUseCase) GetClientRepairs(ctx context.Context, clientID uuid.UUID, requestingUserID uuid.UUID) ([]*domain.Repair, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify the target user is a client
	client, err := uc.userRepo.GetByID(ctx, clientID)
	if err != nil || client == nil || !client.IsClient() {
		return nil, domain.ErrUserNotFound
	}

	// Check permissions: clients can only see their own repairs
	if requestingUser.IsClient() && clientID != requestingUserID {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Admins and managers can see any client's repairs
	if !requestingUser.IsClient() && !requestingUser.CanManageUsers() {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Get all client's cars first
	cars, err := uc.carRepo.GetByOwnerID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client cars: %w", err)
	}

	// Get repairs for all client's cars
	var allRepairs []*domain.Repair
	for _, car := range cars {
		repairs, err := uc.repairRepo.GetByCarID(ctx, car.ID)
		if err != nil {
			continue // Continue with other cars instead of failing completely
		}
		allRepairs = append(allRepairs, repairs...)
	}

	return allRepairs, nil
}
