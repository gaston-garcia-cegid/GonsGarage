package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
)

// ClientModel represents the database model for clients (following Agent.md naming)
type ClientModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FirstName   string     `gorm:"column:first_name;not null"`
	LastName    string     `gorm:"column:last_name;not null"`
	Email       string     `gorm:"column:email;uniqueIndex;not null"`
	Phone       string     `gorm:"column:phone"`
	Address     string     `gorm:"column:address"`
	City        string     `gorm:"column:city"`
	State       string     `gorm:"column:state"`
	ZipCode     string     `gorm:"column:zip_code"`
	DateOfBirth *time.Time `gorm:"column:date_of_birth"`
	IsActive    bool       `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index"`

	// Relationships
	Cars []CarModel `gorm:"foreignKey:OwnerID;references:ID"`
}

// TableName specifies the table name for GORM (following Agent.md database conventions)
func (ClientModel) TableName() string {
	return "clients"
}

// toDomain converts ClientModel to domain.Client (following Agent.md conversion patterns)
func (c *ClientModel) toDomain() *domain.Client {
	client := &domain.Client{
		ID:          c.ID,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Email:       c.Email,
		Phone:       c.Phone,
		Address:     c.Address,
		City:        c.City,
		State:       c.State,
		ZipCode:     c.ZipCode,
		DateOfBirth: c.DateOfBirth,
		IsActive:    c.IsActive,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		DeletedAt:   c.DeletedAt,
	}

	// Convert associated cars
	if len(c.Cars) > 0 {
		client.Cars = make([]domain.Car, len(c.Cars))
		for i, car := range c.Cars {
			if carDomain := car.toDomain(); carDomain != nil {
				client.Cars[i] = *carDomain
			}
		}
	}

	return client
}

// fromDomain converts domain.Client to ClientModel (following Agent.md conversion patterns)
func (c *ClientModel) fromDomain(client *domain.Client) {
	c.ID = client.ID
	c.FirstName = client.FirstName
	c.LastName = client.LastName
	c.Email = client.Email
	c.Phone = client.Phone
	c.Address = client.Address
	c.City = client.City
	c.State = client.State
	c.ZipCode = client.ZipCode
	c.DateOfBirth = client.DateOfBirth
	c.IsActive = client.IsActive
	c.CreatedAt = client.CreatedAt
	c.UpdatedAt = client.UpdatedAt
	c.DeletedAt = client.DeletedAt
}

// postgresClientRepository implements the ClientRepository interface (following Agent.md Clean Architecture)
type postgresClientRepository struct {
	db *gorm.DB
}

// NewPostgresClientRepository creates a new PostgreSQL client repository (following Agent.md dependency injection)
func NewPostgresClientRepository(db *gorm.DB) ports.ClientRepository {
	return &postgresClientRepository{db: db}
}

// Create creates a new client in the database (following Agent.md CRUD patterns)
func (r *postgresClientRepository) Create(ctx context.Context, client *domain.Client) error {
	var model ClientModel
	model.fromDomain(client)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Update the domain object with any auto-generated fields
	*client = *model.toDomain()
	return nil
}

// GetByID retrieves a client by ID (following Agent.md query patterns)
func (r *postgresClientRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	var model ClientModel

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound
		}
		return nil, fmt.Errorf("failed to get client by ID: %w", err)
	}

	return model.toDomain(), nil
}

// GetByEmail retrieves a client by email address (following Agent.md query patterns)
func (r *postgresClientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	var model ClientModel

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound
		}
		return nil, fmt.Errorf("failed to get client by email: %w", err)
	}

	return model.toDomain(), nil
}

// GetWithCars retrieves a client with their associated cars (following Agent.md relationship patterns)
func (r *postgresClientRepository) GetWithCars(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	var model ClientModel

	err := r.db.WithContext(ctx).
		Preload("Cars", "deleted_at IS NULL"). // Only load non-deleted cars
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound
		}
		return nil, fmt.Errorf("failed to get client with cars: %w", err)
	}

	return model.toDomain(), nil
}

// List retrieves clients with pagination (following Agent.md pagination patterns)
func (r *postgresClientRepository) List(ctx context.Context) ([]*domain.Client, error) {
	var models []ClientModel

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}

	// Convert to domain objects
	clients := make([]*domain.Client, len(models))
	for i, model := range models {
		clients[i] = model.toDomain()
	}

	return clients, nil
}

// Search searches clients by name or email (following Agent.md search patterns)
func (r *postgresClientRepository) Search(ctx context.Context, query string, limit int) ([]*domain.Client, error) {
	var models []ClientModel

	searchPattern := "%" + query + "%"

	dbQuery := r.db.WithContext(ctx).
		Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			searchPattern, searchPattern, searchPattern).
		Order("created_at DESC")

	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}

	err := dbQuery.Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to search clients: %w", err)
	}

	// Convert to domain objects
	clients := make([]*domain.Client, len(models))
	for i, model := range models {
		clients[i] = model.toDomain()
	}

	return clients, nil
}

// Update updates an existing client (following Agent.md update patterns)
func (r *postgresClientRepository) Update(ctx context.Context, client *domain.Client) (*domain.Client, error) {
	var model ClientModel
	model.fromDomain(client)

	// Update only non-zero fields and ensure updated_at is set
	model.UpdatedAt = time.Now()

	err := r.db.WithContext(ctx).
		Model(&model).
		Where("id = ?", client.ID).
		Updates(&model).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	// Check if any rows were affected
	if r.db.RowsAffected == 0 {
		return nil, domain.ErrClientNotFound
	}

	// Retrieve the updated client to return
	var updatedModel ClientModel
	err = r.db.WithContext(ctx).
		Where("id = ?", client.ID).
		First(&updatedModel).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated client: %w", err)
	}

	return updatedModel.toDomain(), nil
}

// Delete soft deletes a client (following Agent.md soft delete patterns)
func (r *postgresClientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Model(&ClientModel{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error

	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	// Check if any rows were affected
	if r.db.RowsAffected == 0 {
		return domain.ErrClientNotFound
	}

	return nil
}

// Count returns the total number of active clients (following Agent.md analytics patterns)
func (r *postgresClientRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&ClientModel{}).
		Where("deleted_at IS NULL").
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count clients: %w", err)
	}

	return count, nil
}

// GetActiveClients retrieves only active clients (following Agent.md filtering patterns)
func (r *postgresClientRepository) GetActiveClients(ctx context.Context) ([]*domain.Client, error) {
	var models []ClientModel

	err := r.db.WithContext(ctx).
		Where("is_active = ? AND deleted_at IS NULL", true).
		Order("created_at DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get active clients: %w", err)
	}

	// Convert to domain objects
	clients := make([]*domain.Client, len(models))
	for i, model := range models {
		clients[i] = model.toDomain()
	}

	return clients, nil
}

// DeactivateClient deactivates a client without deleting (following Agent.md state management)
func (r *postgresClientRepository) DeactivateClient(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Model(&ClientModel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to deactivate client: %w", err)
	}

	// Check if any rows were affected
	if r.db.RowsAffected == 0 {
		return domain.ErrClientNotFound
	}

	return nil
}

// ActivateClient reactivates a client (following Agent.md state management)
func (r *postgresClientRepository) ActivateClient(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Model(&ClientModel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"is_active":  true,
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to activate client: %w", err)
	}

	// Check if any rows were affected
	if r.db.RowsAffected == 0 {
		return domain.ErrClientNotFound
	}

	return nil
}

// ExistsByEmail checks if a client with the given email exists (following Agent.md validation patterns)
func (r *postgresClientRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&ClientModel{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// GetClientsByDateRange retrieves clients created within a date range (following Agent.md analytics patterns)
func (r *postgresClientRepository) GetClientsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Client, error) {
	var models []ClientModel

	err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).
		Order("created_at DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get clients by date range: %w", err)
	}

	// Convert to domain objects
	clients := make([]*domain.Client, len(models))
	for i, model := range models {
		clients[i] = model.toDomain()
	}

	return clients, nil
}
