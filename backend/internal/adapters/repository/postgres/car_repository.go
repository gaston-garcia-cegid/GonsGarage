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

// postgresCarRepository implements CarRepository using PostgreSQL
type postgresCarRepository struct {
	db *gorm.DB
}

// NewPostgresCarRepository creates a new PostgreSQL car repository
func NewPostgresCarRepository(db *gorm.DB) ports.CarRepository {
	return &postgresCarRepository{db: db}
}

// CarModel represents the database table structure
type CarModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Make         string     `gorm:"not null"`
	Model        string     `gorm:"not null"`
	Year         int        `gorm:"not null"`
	LicensePlate string     `gorm:"column:license_plate;not null;uniqueIndex"` // ✅ Fixed column name
	VIN          string     `gorm:"column:vin"`
	Color        string     `gorm:"not null"`
	Mileage      int        `gorm:"not null;default:0"`
	OwnerID      uuid.UUID  `gorm:"type:uuid;column:owner_id;not null;index"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index"`

	// Relationships
	Owner   UserModel     `gorm:"foreignKey:OwnerID;references:ID"`
	Repairs []RepairModel `gorm:"foreignKey:CarID;references:ID"`
}

// TableName specifies the database table name
func (CarModel) TableName() string {
	return "cars"
}

// Create stores a new car in the database
func (r *postgresCarRepository) Create(ctx context.Context, car *domain.Car) error {
	dbCar := r.toCarModel(car)

	if err := r.db.WithContext(ctx).Create(dbCar).Error; err != nil {
		return fmt.Errorf("failed to create car in database: %w", err)
	}

	return nil
}

// GetByID retrieves a car by its unique identifier
func (r *postgresCarRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	var dbCar CarModel

	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Owner").
		First(&dbCar).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCarNotFound
		}
		return nil, fmt.Errorf("failed to get car by ID: %w", err)
	}

	return r.toDomainCar(&dbCar), nil
}

// GetByOwnerID retrieves all cars owned by a specific user
func (r *postgresCarRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Car, error) {
	var dbCars []CarModel

	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND deleted_at IS NULL", ownerID).
		Preload("Owner").
		Order("created_at DESC").
		Find(&dbCars).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get cars by owner ID: %w", err)
	}

	cars := make([]*domain.Car, len(dbCars))
	for i, dbCar := range dbCars {
		cars[i] = r.toDomainCar(&dbCar)
	}

	return cars, nil
}

// GetByLicensePlate retrieves a car by its license plate
func (r *postgresCarRepository) GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error) {
	var dbCar CarModel

	err := r.db.WithContext(ctx).
		Where("license_plate = ? AND deleted_at IS NULL", licensePlate).
		Preload("Owner").
		First(&dbCar).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not finding a license plate is not an error
		}
		return nil, fmt.Errorf("failed to get car by license plate: %w", err)
	}

	return r.toDomainCar(&dbCar), nil
}

// List retrieves cars with pagination
func (r *postgresCarRepository) List(ctx context.Context, limit, offset int) ([]*domain.Car, error) {
	var dbCars []CarModel

	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Preload("Owner").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&dbCars).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list cars: %w", err)
	}

	cars := make([]*domain.Car, len(dbCars))
	for i, dbCar := range dbCars {
		cars[i] = r.toDomainCar(&dbCar)
	}

	return cars, nil
}

// Update modifies an existing car
func (r *postgresCarRepository) Update(ctx context.Context, car *domain.Car) error {
	dbCar := r.toCarModel(car)

	result := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", car.ID).
		Updates(dbCar)

	if result.Error != nil {
		return fmt.Errorf("failed to update car: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrCarNotFound
	}

	return nil
}

// Delete removes a car using soft delete
func (r *postgresCarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("failed to delete car: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrCarNotFound
	}

	return nil
}

// GetWithRepairs retrieves a car with its complete repair history
func (r *postgresCarRepository) GetWithRepairs(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	var dbCar CarModel

	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Owner").
		Preload("Repairs", "deleted_at IS NULL").
		First(&dbCar).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCarNotFound
		}
		return nil, fmt.Errorf("failed to get car with repairs: %w", err)
	}

	return r.toDomainCarWithRepairs(&dbCar), nil
}

// Conversion methods following Clean Architecture principles

func (r *postgresCarRepository) toCarModel(car *domain.Car) *CarModel {
	return &CarModel{
		ID:           car.ID,
		Make:         car.Make,
		Model:        car.Model,
		Year:         car.Year,
		LicensePlate: car.LicensePlate,
		VIN:          car.VIN,
		Color:        car.Color,
		Mileage:      car.Mileage,
		OwnerID:      car.OwnerID,
		CreatedAt:    car.CreatedAt,
		UpdatedAt:    car.UpdatedAt,
		DeletedAt:    car.DeletedAt,
	}
}

func (r *postgresCarRepository) toDomainCar(dbCar *CarModel) *domain.Car {
	car := &domain.Car{
		ID:           dbCar.ID,
		Make:         dbCar.Make,
		Model:        dbCar.Model,
		Year:         dbCar.Year,
		LicensePlate: dbCar.LicensePlate,
		VIN:          dbCar.VIN,
		Color:        dbCar.Color,
		Mileage:      dbCar.Mileage,
		OwnerID:      dbCar.OwnerID,
		CreatedAt:    dbCar.CreatedAt,
		UpdatedAt:    dbCar.UpdatedAt,
		DeletedAt:    dbCar.DeletedAt,
	}

	// Convert owner if preloaded
	if dbCar.Owner.ID != uuid.Nil {
		car.Owner = r.toDomainUser(&dbCar.Owner)
	}

	return car
}

func (r *postgresCarRepository) toDomainCarWithRepairs(dbCar *CarModel) *domain.Car {
	car := r.toDomainCar(dbCar)

	// Convert repairs if preloaded
	if len(dbCar.Repairs) > 0 {
		car.Repairs = make([]domain.Repair, len(dbCar.Repairs))
		for i, dbRepair := range dbCar.Repairs {
			car.Repairs[i] = r.toDomainRepair(&dbRepair)
		}
	}

	return car
}

// Helper conversion methods (these would be shared or in a separate file)
func (r *postgresCarRepository) toDomainUser(dbUser *UserModel) domain.User {
	return domain.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Role:      dbUser.Role,
		IsActive:  dbUser.IsActive,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		DeletedAt: dbUser.DeletedAt,
	}
}

func (r *postgresCarRepository) toDomainRepair(dbRepair *RepairModel) domain.Repair {
	return domain.Repair{
		ID:           dbRepair.ID,
		CarID:        dbRepair.CarID,
		TechnicianID: dbRepair.TechnicianID,
		Description:  dbRepair.Description,
		Status:       domain.RepairStatus(dbRepair.Status),
		Cost:         dbRepair.Cost,
		StartedAt:    &dbRepair.StartedAt,
		CompletedAt:  dbRepair.CompletedAt,
		CreatedAt:    dbRepair.CreatedAt,
		UpdatedAt:    dbRepair.UpdatedAt,
		DeletedAt:    dbRepair.DeletedAt,
	}
}
