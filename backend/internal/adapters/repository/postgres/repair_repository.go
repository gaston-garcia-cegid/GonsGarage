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

type PostgresRepairRepository struct {
	db *gorm.DB
}

func NewPostgresRepairRepository(db *gorm.DB) ports.RepairRepository {
	return &PostgresRepairRepository{db: db}
}

// RepairModel represents the database table structure
type RepairModel struct {
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
	CarID        uuid.UUID  `gorm:"type:uuid;not null;index"`
	TechnicianID uuid.UUID  `gorm:"type:uuid;not null;index"`
	Description  string     `gorm:"not null"`
	Status       string     `gorm:"not null;default:'pending'"`
	Cost         float64    `gorm:"not null;default:0"`
	StartedAt    time.Time  `gorm:"column:start_date"`
	CompletedAt  *time.Time `gorm:"column:completed_at"`
}

func (RepairModel) TableName() string {
	return "repairs"
}

func (r *PostgresRepairRepository) Create(ctx context.Context, repair *domain.Repair) error {
	dbRepair := &RepairModel{
		ID:        repair.ID,
		CreatedAt: repair.CreatedAt,
		UpdatedAt: repair.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(dbRepair).Error; err != nil {
		return fmt.Errorf("failed to create repair: %w", err)
	}

	return nil
}

func (r *PostgresRepairRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Repair, error) {
	var dbRepair RepairModel

	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&dbRepair).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRepairNotFound
		}
		return nil, fmt.Errorf("failed to get repair by ID: %w", err)
	}

	return r.toDomainRepair(&dbRepair), nil
}

func (r *PostgresRepairRepository) GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*domain.Repair, error) {
	var dbRepairs []RepairModel

	err := r.db.WithContext(ctx).Where("client_id = ? AND deleted_at IS NULL", clientID).
		Order("created_at DESC").Find(&dbRepairs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get repairs by client ID: %w", err)
	}

	repairs := make([]*domain.Repair, len(dbRepairs))
	for i, dbRepair := range dbRepairs {
		repairs[i] = r.toDomainRepair(&dbRepair)
	}

	return repairs, nil
}

// GetByCarID implements ports.RepairRepository.
func (r *PostgresRepairRepository) GetByCarID(ctx context.Context, carID uuid.UUID) ([]*domain.Repair, error) {
	var dbRepairs []RepairModel

	err := r.db.WithContext(ctx).Where("car_id = ? AND deleted_at IS NULL", carID).
		Order("created_at DESC").Find(&dbRepairs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get repairs by car ID: %w", err)
	}

	repairs := make([]*domain.Repair, len(dbRepairs))
	for i, dbRepair := range dbRepairs {
		repairs[i] = r.toDomainRepair(&dbRepair)
	}

	return repairs, nil
}

func (r *PostgresRepairRepository) List(ctx context.Context, limit, offset int) ([]*domain.Repair, error) {
	var dbRepairs []RepairModel

	query := r.db.WithContext(ctx).Where("deleted_at IS NULL").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&dbRepairs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list repairs: %w", err)
	}

	repairs := make([]*domain.Repair, len(dbRepairs))
	for i, dbRepair := range dbRepairs {
		repairs[i] = r.toDomainRepair(&dbRepair)
	}

	return repairs, nil
}

func (r *PostgresRepairRepository) Update(ctx context.Context, repair *domain.Repair) error {
	dbRepair := &RepairModel{
		ID:        repair.ID,
		CreatedAt: repair.CreatedAt,
		UpdatedAt: repair.UpdatedAt,
	}

	result := r.db.WithContext(ctx).Model(dbRepair).Where("id = ?", repair.ID).Updates(dbRepair)

	if result.Error != nil {
		return fmt.Errorf("failed to update repair: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrRepairNotFound
	}

	return nil
}

func (r *PostgresRepairRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&RepairModel{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete repair: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrRepairNotFound
	}

	return nil
}

func (r *PostgresRepairRepository) GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Repair, error) {
	var dbRepair RepairModel

	err := r.db.WithContext(ctx).Where("license_plate = ? AND deleted_at IS NULL", licensePlate).First(&dbRepair).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRepairNotFound
		}
		return nil, fmt.Errorf("failed to get repair by license plate: %w", err)
	}

	return r.toDomainRepair(&dbRepair), nil
}

func (r *PostgresRepairRepository) toDomainRepair(dbRepair *RepairModel) *domain.Repair {
	return &domain.Repair{
		ID:        dbRepair.ID,
		CreatedAt: dbRepair.CreatedAt,
		UpdatedAt: dbRepair.UpdatedAt,
	}
}
