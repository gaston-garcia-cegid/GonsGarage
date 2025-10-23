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
	clientRepo ports.ClientRepository
	userRepo   ports.UserRepository
	cacheRepo  ports.CacheRepository
}

func NewClientUseCase(clientRepo ports.ClientRepository, userRepo ports.UserRepository, cacheRepo ports.CacheRepository) *ClientUseCase {
	return &ClientUseCase{
		clientRepo: clientRepo,
		userRepo:   userRepo,
		cacheRepo:  cacheRepo,
	}
}
func (uc *ClientUseCase) CreateClient(ctx context.Context, client *domain.Client) error {
	// Check if client with the same email already exists
	existingClient, err := uc.clientRepo.GetByEmail(ctx, client.Email)
	if err == nil && existingClient != nil {
		return fmt.Errorf("client with email %s already exists", client.Email)
	}

	// Create the client
	client.ID = uuid.New()
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()

	return uc.clientRepo.Create(ctx, client)
}
func (uc *ClientUseCase) ListClients(ctx context.Context) ([]*domain.Client, error) {
	return uc.clientRepo.List(ctx)
}
func (uc *ClientUseCase) GetClient(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	return uc.clientRepo.GetByID(ctx, id)
}
func (uc *ClientUseCase) UpdateClient(ctx context.Context, id uuid.UUID, updatedData *domain.Client) (*domain.Client, error) {
	// Get existing client
	client, err := uc.clientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	// Update client fields
	client.FirstName = updatedData.FirstName
	client.LastName = updatedData.LastName
	client.Email = updatedData.Email
	client.UpdatedAt = time.Now()

	return uc.clientRepo.Update(ctx, client)
}
func (uc *ClientUseCase) DeleteClient(ctx context.Context, id uuid.UUID) error {
	return uc.clientRepo.Delete(ctx, id)
}
