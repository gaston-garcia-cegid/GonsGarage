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

type PostgresCarRepository struct {
	db *gorm.DB
}

func NewPostgresCarRepository(db *gorm.DB) ports.CarRepository {
	return &PostgresCarRepository{db: db}
}

// CarModel represents the database table structure
type CarModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClientID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	Make         string     `gorm:"not null"`
	Model        string     `gorm:"not null"`
	Year         int        `gorm:"not null"`
	LicensePlate string     `gorm:"column:license_plate;not null;uniqueIndex"`
	VIN          string     `gorm:"column:vin"`
	Color        string     `gorm:"not null"`
	Mileage      int        `gorm:"not null;default:0"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index"`
}

func (CarModel) TableName() string {
	return "cars"
}

func (r *PostgresCarRepository) Create(ctx context.Context, car *domain.Car) error {
	dbCar := &CarModel{
		ID:           car.ID,
		ClientID:     car.OwnerID,
		Make:         car.Make,
		Model:        car.Model,
		Year:         car.Year,
		LicensePlate: car.LicensePlate,
		VIN:          car.VIN,
		Color:        car.Color,
		CreatedAt:    car.CreatedAt,
		UpdatedAt:    car.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(dbCar).Error; err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}

	return nil
}

func (r *PostgresCarRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	var dbCar CarModel

	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&dbCar).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCarNotFound
		}
		return nil, fmt.Errorf("failed to get car by ID: %w", err)
	}

	return r.toDomainCar(&dbCar), nil
}

func (r *PostgresCarRepository) GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*domain.Car, error) {
	var dbCars []CarModel

	err := r.db.WithContext(ctx).Where("client_id = ? AND deleted_at IS NULL", clientID).
		Order("created_at DESC").Find(&dbCars).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get cars by client ID: %w", err)
	}

	cars := make([]*domain.Car, len(dbCars))
	for i, dbCar := range dbCars {
		cars[i] = r.toDomainCar(&dbCar)
	}

	return cars, nil
}

func (r *PostgresCarRepository) List(ctx context.Context, limit, offset int) ([]*domain.Car, error) {
	var dbCars []CarModel

	query := r.db.WithContext(ctx).Where("deleted_at IS NULL").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&dbCars).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list cars: %w", err)
	}

	cars := make([]*domain.Car, len(dbCars))
	for i, dbCar := range dbCars {
		cars[i] = r.toDomainCar(&dbCar)
	}

	return cars, nil
}

func (r *PostgresCarRepository) Update(ctx context.Context, car *domain.Car) error {
	dbCar := &CarModel{
		ID:           car.ID,
		ClientID:     car.OwnerID,
		Make:         car.Make,
		Model:        car.Model,
		Year:         car.Year,
		LicensePlate: car.LicensePlate,
		VIN:          car.VIN,
		Color:        car.Color,
		UpdatedAt:    car.UpdatedAt,
	}

	result := r.db.WithContext(ctx).Model(dbCar).Where("id = ?", car.ID).Updates(dbCar)
	if result.Error != nil {
		return fmt.Errorf("failed to update car: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrCarNotFound
	}

	return nil
}

func (r *PostgresCarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&CarModel{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete car: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrCarNotFound
	}

	return nil
}

func (r *PostgresCarRepository) GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error) {
	var dbCar CarModel

	err := r.db.WithContext(ctx).Where("license_plate = ? AND deleted_at IS NULL", licensePlate).First(&dbCar).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCarNotFound
		}
		return nil, fmt.Errorf("failed to get car by license plate: %w", err)
	}

	return r.toDomainCar(&dbCar), nil
}

func (r *PostgresCarRepository) toDomainCar(dbCar *CarModel) *domain.Car {
	return &domain.Car{
		ID:           dbCar.ID,
		OwnerID:      dbCar.ClientID,
		Make:         dbCar.Make,
		Model:        dbCar.Model,
		Year:         dbCar.Year,
		LicensePlate: dbCar.LicensePlate,
		VIN:          dbCar.VIN,
		Color:        dbCar.Color,
		CreatedAt:    dbCar.CreatedAt,
		UpdatedAt:    dbCar.UpdatedAt,
		DeletedAt:    dbCar.DeletedAt,
	}
}
