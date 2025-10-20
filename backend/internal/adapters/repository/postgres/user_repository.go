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

// PostgresUserRepository implements UserRepository interface using PostgreSQL
type PostgresUserRepository struct {
	db *gorm.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *gorm.DB) ports.UserRepository {
	return &PostgresUserRepository{db: db}
}

// Create implements UserRepository.Create
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	// Convert domain entity to database model
	dbUser := &UserModel{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		IsActive:     user.IsActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Update the domain entity with generated data
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt

	return nil
}

// GetByID implements UserRepository.GetByID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var dbUser UserModel

	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return r.toDomainUser(&dbUser), nil
}

// GetByEmail implements UserRepository.GetByEmail
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser UserModel

	err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.toDomainUser(&dbUser), nil
}

// Update implements UserRepository.Update
func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	dbUser := &UserModel{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		IsActive:     user.IsActive,
		UpdatedAt:    time.Now(),
	}

	result := r.db.WithContext(ctx).Model(dbUser).Where("id = ?", user.ID).Updates(dbUser)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	user.UpdatedAt = dbUser.UpdatedAt
	return nil
}

// Delete implements UserRepository.Delete (soft delete)
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List implements UserRepository.List
func (r *PostgresUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	var dbUsers []UserModel

	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&dbUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.toDomainUser(&dbUser)
	}

	return users, nil
}

// EmailExists implements UserRepository.EmailExists
func (r *PostgresUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UserModel{}).Where("email = ? AND deleted_at IS NULL", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// toDomainUser converts database model to domain entity
func (r *PostgresUserRepository) toDomainUser(dbUser *UserModel) *domain.User {
	return &domain.User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		FirstName:    dbUser.FirstName,
		LastName:     dbUser.LastName,
		Role:         dbUser.Role,
		IsActive:     dbUser.IsActive,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}

// UserModel represents the database table structure
type UserModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string     `gorm:"uniqueIndex;not null"`
	PasswordHash string     `gorm:"not null"`
	FirstName    string     `gorm:"not null"`
	LastName     string     `gorm:"not null"`
	Role         string     `gorm:"not null;default:'employee'"`
	IsActive     bool       `gorm:"default:true"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `gorm:"index"`
}

// TableName specifies the table name for GORM
func (UserModel) TableName() string {
	return "users"
}
