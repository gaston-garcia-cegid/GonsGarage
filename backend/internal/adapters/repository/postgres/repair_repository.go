package postgres

import (
	"context"
	"errors"
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

// RepairModel refleja la tabla repairs alineada con domain.Repair (migración GORM).
type RepairModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CarID        uuid.UUID  `gorm:"column:car_id;type:uuid;not null;index"`
	TechnicianID uuid.UUID  `gorm:"column:technician_id;type:uuid;not null;index"`
	Description  string     `gorm:"not null"`
	Status       string     `gorm:"not null;default:pending"`
	Cost         float64    `gorm:"type:decimal(10,2);not null;default:0"`
	StartedAt    *time.Time `gorm:"column:started_at"`
	CompletedAt  *time.Time `gorm:"column:completed_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index"`
}

func (RepairModel) TableName() string {
	return "repairs"
}

func (r *PostgresRepairRepository) Create(ctx context.Context, repair *domain.Repair) error {
	m := r.fromDomain(repair)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("failed to create repair: %w", err)
	}
	return nil
}

func (r *PostgresRepairRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Repair, error) {
	var m RepairModel
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRepairNotFound
		}
		return nil, fmt.Errorf("failed to get repair by ID: %w", err)
	}
	return r.toDomain(&m), nil
}

func (r *PostgresRepairRepository) GetByCarID(ctx context.Context, carID uuid.UUID) ([]*domain.Repair, error) {
	var rows []RepairModel
	err := r.db.WithContext(ctx).
		Where("car_id = ? AND deleted_at IS NULL", carID).
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get repairs by car ID: %w", err)
	}
	out := make([]*domain.Repair, 0, len(rows))
	for i := range rows {
		out = append(out, r.toDomain(&rows[i]))
	}
	return out, nil
}

func (r *PostgresRepairRepository) Update(ctx context.Context, repair *domain.Repair) error {
	var m RepairModel
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", repair.ID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrRepairNotFound
		}
		return fmt.Errorf("failed to load repair: %w", err)
	}
	m.CarID = repair.CarID
	m.TechnicianID = repair.TechnicianID
	m.Description = repair.Description
	m.Status = string(repair.Status)
	m.Cost = repair.Cost
	m.StartedAt = repair.StartedAt
	m.CompletedAt = repair.CompletedAt
	m.UpdatedAt = repair.UpdatedAt

	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("failed to update repair: %w", err)
	}
	return nil
}

func (r *PostgresRepairRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Model(&RepairModel{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", time.Now().UTC())
	if res.Error != nil {
		return fmt.Errorf("failed to delete repair: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrRepairNotFound
	}
	return nil
}

func (r *PostgresRepairRepository) fromDomain(d *domain.Repair) RepairModel {
	return RepairModel{
		ID:           d.ID,
		CarID:        d.CarID,
		TechnicianID: d.TechnicianID,
		Description:  d.Description,
		Status:       string(d.Status),
		Cost:         d.Cost,
		StartedAt:    d.StartedAt,
		CompletedAt:  d.CompletedAt,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		DeletedAt:    d.DeletedAt,
	}
}

func (r *PostgresRepairRepository) toDomain(m *RepairModel) *domain.Repair {
	return &domain.Repair{
		ID:           m.ID,
		CarID:        m.CarID,
		TechnicianID: m.TechnicianID,
		Description:  m.Description,
		Status:       domain.RepairStatus(m.Status),
		Cost:         m.Cost,
		StartedAt:    m.StartedAt,
		CompletedAt:  m.CompletedAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    m.DeletedAt,
	}
}
