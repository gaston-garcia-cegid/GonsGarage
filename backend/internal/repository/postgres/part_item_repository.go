package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
)

const partItemActiveRow = "id = ? AND deleted_at IS NULL"

type postgresPartItemRepository struct {
	db *gorm.DB
}

// NewPostgresPartItemRepository returns a PartItemRepository backed by GORM (PostgreSQL or sqlite tests).
func NewPostgresPartItemRepository(db *gorm.DB) ports.PartItemRepository {
	return &postgresPartItemRepository{db: db}
}

func (r *postgresPartItemRepository) Create(ctx context.Context, p *domain.PartItem) error {
	if p == nil {
		return fmt.Errorf("part item is nil")
	}
	if err := p.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return fmt.Errorf("failed to create part item: %w", err)
	}
	return nil
}

func (r *postgresPartItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PartItem, error) {
	var row domain.PartItem
	err := r.db.WithContext(ctx).Where(partItemActiveRow, id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPartItemNotFound
		}
		return nil, fmt.Errorf("failed to get part item: %w", err)
	}
	return &row, nil
}

func (r *postgresPartItemRepository) GetByBarcode(ctx context.Context, barcode string) (*domain.PartItem, error) {
	b := strings.TrimSpace(barcode)
	if b == "" {
		return nil, domain.ErrPartItemNotFound
	}
	var row domain.PartItem
	err := r.db.WithContext(ctx).Where("barcode = ? AND deleted_at IS NULL", b).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPartItemNotFound
		}
		return nil, fmt.Errorf("failed to get part item by barcode: %w", err)
	}
	return &row, nil
}

func (r *postgresPartItemRepository) Update(ctx context.Context, p *domain.PartItem) error {
	if p == nil {
		return fmt.Errorf("part item is nil")
	}
	if err := p.Validate(); err != nil {
		return err
	}
	res := r.db.WithContext(ctx).Model(&domain.PartItem{}).
		Where(partItemActiveRow, p.ID).
		Updates(map[string]interface{}{
			"reference":          p.Reference,
			"brand":              p.Brand,
			"name":               p.Name,
			"barcode":            p.Barcode,
			"quantity":           p.Quantity,
			"uom":                p.UOM,
			"minimum_quantity":   p.MinimumQuantity,
			"updated_at":         time.Now().UTC(),
		})
	if res.Error != nil {
		return fmt.Errorf("failed to update part item: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrPartItemNotFound
	}
	return nil
}

func (r *postgresPartItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Model(&domain.PartItem{}).
		Where(partItemActiveRow, id).
		Update("deleted_at", time.Now().UTC())
	if res.Error != nil {
		return fmt.Errorf("failed to delete part item: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrPartItemNotFound
	}
	return nil
}

func (r *postgresPartItemRepository) listQuery(ctx context.Context, f ports.PartItemListFilters) *gorm.DB {
	q := r.db.WithContext(ctx).Model(&domain.PartItem{}).Where("deleted_at IS NULL")
	if f.Barcode != nil {
		b := strings.TrimSpace(*f.Barcode)
		if b != "" {
			q = q.Where("barcode = ?", b)
		}
	}
	if f.Search != nil {
		s := strings.TrimSpace(*f.Search)
		if s != "" {
			pat := "%" + s + "%"
			q = q.Where("reference LIKE ? OR brand LIKE ? OR name LIKE ?", pat, pat, pat)
		}
	}
	return q
}

func (r *postgresPartItemRepository) List(ctx context.Context, f ports.PartItemListFilters) ([]*domain.PartItem, int64, error) {
	limit, offset := clampRepoList(f.Limit, f.Offset)
	base := r.listQuery(ctx, f)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count part items: %w", err)
	}
	var rows []domain.PartItem
	q2 := r.listQuery(ctx, f)
	if err := q2.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list part items: %w", err)
	}
	out := make([]*domain.PartItem, 0, len(rows))
	for i := range rows {
		out = append(out, &rows[i])
	}
	return out, total, nil
}

var _ ports.PartItemRepository = (*postgresPartItemRepository)(nil)
