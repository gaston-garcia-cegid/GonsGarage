package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
)

type PostgresServiceJobRepository struct {
	db *gorm.DB
}

func NewPostgresServiceJobRepository(db *gorm.DB) ports.ServiceJobRepository {
	return &PostgresServiceJobRepository{db: db}
}

func (r *PostgresServiceJobRepository) Create(ctx context.Context, job *domain.ServiceJob) error {
	if err := r.db.WithContext(ctx).Create(job).Error; err != nil {
		return fmt.Errorf("create service job: %w", err)
	}
	return nil
}

func (r *PostgresServiceJobRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ServiceJob, error) {
	var j domain.ServiceJob
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&j).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrServiceJobNotFound
		}
		return nil, err
	}
	return &j, nil
}

func (r *PostgresServiceJobRepository) Update(ctx context.Context, job *domain.ServiceJob) error {
	if err := r.db.WithContext(ctx).Save(job).Error; err != nil {
		return fmt.Errorf("update service job: %w", err)
	}
	return nil
}

func (r *PostgresServiceJobRepository) ListByCarID(ctx context.Context, carID uuid.UUID) ([]*domain.ServiceJob, error) {
	var rows []*domain.ServiceJob
	if err := r.db.WithContext(ctx).Where("car_id = ?", carID).Order("opened_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []*domain.ServiceJob{}
	}
	return rows, nil
}

// ListByOpenedOn returns service jobs with OpenedAt in [day 00:00:00 UTC, next calendar day 00:00:00 UTC).
// The `day` argument is normalized: only year, month, and day in UTC are used.
func (r *PostgresServiceJobRepository) ListByOpenedOn(ctx context.Context, day time.Time) ([]*domain.ServiceJob, error) {
	y, m, d := day.UTC().Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	var rows []*domain.ServiceJob
	if err := r.db.WithContext(ctx).Where("opened_at >= ? AND opened_at < ? AND deleted_at IS NULL", start, end).Order("opened_at asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []*domain.ServiceJob{}
	}
	return rows, nil
}

func (r *PostgresServiceJobRepository) SaveReception(ctx context.Context, rec *domain.ServiceJobReception) error {
	var m domain.ServiceJobReception
	tx := r.db.WithContext(ctx)
	err := tx.Where("service_job_id = ?", rec.ServiceJobID).First(&m).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(rec).Error; err != nil {
			return fmt.Errorf("create reception: %w", err)
		}
		return nil
	}
	m.OdometerKM = rec.OdometerKM
	m.OilLevel = rec.OilLevel
	m.CoolantLevel = rec.CoolantLevel
	m.TiresNote = rec.TiresNote
	m.GeneralNotes = rec.GeneralNotes
	m.RecordedByUserID = rec.RecordedByUserID
	m.RecordedAt = rec.RecordedAt
	if rec.SchemaVersion > 0 {
		m.SchemaVersion = rec.SchemaVersion
	}
	if err := tx.Save(&m).Error; err != nil {
		return fmt.Errorf("update reception: %w", err)
	}
	return nil
}

func (r *PostgresServiceJobRepository) GetReception(ctx context.Context, serviceJobID uuid.UUID) (*domain.ServiceJobReception, error) {
	var out domain.ServiceJobReception
	if err := r.db.WithContext(ctx).Where("service_job_id = ?", serviceJobID).First(&out).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func (r *PostgresServiceJobRepository) SaveHandover(ctx context.Context, h *domain.ServiceJobHandover) error {
	var m domain.ServiceJobHandover
	tx := r.db.WithContext(ctx)
	err := tx.Where("service_job_id = ?", h.ServiceJobID).First(&m).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(h).Error; err != nil {
			return fmt.Errorf("create handover: %w", err)
		}
		return nil
	}
	m.OdometerKM = h.OdometerKM
	m.TiresNote = h.TiresNote
	m.GeneralNotes = h.GeneralNotes
	m.RecordedByUserID = h.RecordedByUserID
	m.RecordedAt = h.RecordedAt
	if h.SchemaVersion > 0 {
		m.SchemaVersion = h.SchemaVersion
	}
	if err := tx.Save(&m).Error; err != nil {
		return fmt.Errorf("update handover: %w", err)
	}
	return nil
}

func (r *PostgresServiceJobRepository) GetHandover(ctx context.Context, serviceJobID uuid.UUID) (*domain.ServiceJobHandover, error) {
	var out domain.ServiceJobHandover
	if err := r.db.WithContext(ctx).Where("service_job_id = ?", serviceJobID).First(&out).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}
