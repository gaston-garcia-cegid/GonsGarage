package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/google/uuid"
)

// ClientUseCase handles client-related business logic (following Agent.md Clean Architecture)
type ClientUseCase struct {
	clientRepo ports.ClientRepository
	userRepo   ports.UserRepository
	carRepo    ports.CarRepository
	repairRepo ports.RepairRepository
}

// NewClientUseCase creates a new client use case instance (following Agent.md dependency injection)
func NewClientUseCase(
	clientRepo ports.ClientRepository,
	userRepo ports.UserRepository,
	carRepo ports.CarRepository,
	repairRepo ports.RepairRepository,
) ports.ClientUseCase {
	return &ClientUseCase{
		clientRepo: clientRepo,
		userRepo:   userRepo,
		carRepo:    carRepo,
		repairRepo: repairRepo,
	}
}

// CreateClient creates a new client - ✅ Fixed: using *domain.Client parameter and return type
func (uc *ClientUseCase) CreateClient(ctx context.Context, client *domain.Client) (*domain.Client, error) {
	// Check if email already exists
	existingClient, err := uc.clientRepo.GetByEmail(ctx, client.Email)
	if err == nil && existingClient != nil {
		return nil, domain.ErrClientAlreadyExists
	}

	// Validate client data (following Agent.md validation rules)
	if err := client.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client data: %w", err)
	}

	// Set metadata
	client.ID = uuid.New()
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()
	client.IsActive = true

	// Create the client
	if err := uc.clientRepo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

// GetClient retrieves a client by ID with permission checks (following Agent.md authorization)
func (uc *ClientUseCase) GetClient(ctx context.Context, clientID uuid.UUID, requestingUserID uuid.UUID) (*domain.Client, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get the client
	client, err := uc.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	if client == nil {
		return nil, domain.ErrClientNotFound
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

// GetClientProfile allows clients to get their own profile (following Agent.md self-service patterns)
func (uc *ClientUseCase) GetClientProfile(ctx context.Context, clientID uuid.UUID) (*domain.Client, error) {
	client, err := uc.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	if client == nil {
		return nil, domain.ErrClientNotFound
	}

	return client, nil
}

// UpdateClient updates an existing client (following Agent.md business logic)
func (uc *ClientUseCase) UpdateClient(ctx context.Context, client *domain.Client, requestingUserID uuid.UUID) (*domain.Client, error) {
	// Get the requesting user
	requestingUser, err := uc.userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get the existing client
	existingClient, err := uc.clientRepo.GetByID(ctx, client.ID)
	if err != nil || existingClient == nil {
		return nil, domain.ErrClientNotFound
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
		return nil, fmt.Errorf("invalid client data: %w", err)
	}

	// Preserve some fields (following Agent.md data integrity rules)
	client.CreatedAt = existingClient.CreatedAt
	client.UpdatedAt = time.Now()

	// Update the client
	updatedClient, err := uc.clientRepo.Update(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	return updatedClient, nil
}

// UpdateClientProfile allows clients to update their own profile
func (uc *ClientUseCase) UpdateClientProfile(ctx context.Context, clientID uuid.UUID, client *domain.Client) (*domain.Client, error) {
	// Get the existing client
	existingClient, err := uc.clientRepo.GetByID(ctx, clientID)
	if err != nil || existingClient == nil {
		return nil, domain.ErrClientNotFound
	}

	// Validate client data
	if err := client.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client data: %w", err)
	}

	// Preserve critical fields (clients can't change certain fields themselves)
	client.ID = clientID
	client.CreatedAt = existingClient.CreatedAt
	client.UpdatedAt = time.Now()
	client.IsActive = existingClient.IsActive // Clients can't activate/deactivate themselves

	// Update the client
	updatedClient, err := uc.clientRepo.Update(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	return updatedClient, nil
}

// DeleteClient deletes a client (soft delete)
func (uc *ClientUseCase) DeleteClient(ctx context.Context, clientID uuid.UUID) error {
	// Get the client
	client, err := uc.clientRepo.GetByID(ctx, clientID)
	if err != nil || client == nil {
		return domain.ErrClientNotFound
	}

	// Note: Only admins/managers should be able to call this method
	// The authorization should be handled at the handler/service layer

	// Delete the client (soft delete)
	if err := uc.clientRepo.Delete(ctx, clientID); err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	return nil
}

// ListClients retrieves all clients (admin/manager only) - updated to match interface signature
func (uc *ClientUseCase) ListClients(ctx context.Context) ([]*domain.Client, error) {
	// NOTE: Authorization should be handled at the handler/service layer since we no longer have requestingUserID here.

	// Get clients
	clients, err := uc.clientRepo.List(ctx)
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

	// Verify the target client exists
	client, err := uc.clientRepo.GetByID(ctx, clientID)
	if err != nil || client == nil {
		return nil, domain.ErrClientNotFound
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

	// Verify the target client exists
	client, err := uc.clientRepo.GetByID(ctx, clientID)
	if err != nil || client == nil {
		return nil, domain.ErrClientNotFound
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
