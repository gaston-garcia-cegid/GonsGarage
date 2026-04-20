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

type postgresBillingDocumentRepository struct {
	db *gorm.DB
}

// NewPostgresBillingDocumentRepository returns ports.BillingDocumentRepository.
func NewPostgresBillingDocumentRepository(db *gorm.DB) ports.BillingDocumentRepository {
	return &postgresBillingDocumentRepository{db: db}
}

func (r *postgresBillingDocumentRepository) Create(ctx context.Context, doc *domain.BillingDocument) error {
	if doc == nil {
		return fmt.Errorf("billing document is nil")
	}
	if err := r.db.WithContext(ctx).Create(doc).Error; err != nil {
		return fmt.Errorf("failed to create billing document: %w", err)
	}
	return nil
}

func (r *postgresBillingDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.BillingDocument, error) {
	var doc domain.BillingDocument
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBillingDocumentNotFound
		}
		return nil, fmt.Errorf("failed to get billing document: %w", err)
	}
	return &doc, nil
}

func (r *postgresBillingDocumentRepository) Update(ctx context.Context, doc *domain.BillingDocument) error {
	if doc == nil {
		return fmt.Errorf("billing document is nil")
	}
	res := r.db.WithContext(ctx).Model(&domain.BillingDocument{}).
		Where("id = ? AND deleted_at IS NULL", doc.ID).
		Updates(map[string]interface{}{
			"kind":         string(doc.Kind),
			"title":        doc.Title,
			"amount":       doc.Amount,
			"customer_id":  doc.CustomerID,
			"reference":    doc.Reference,
			"notes":        doc.Notes,
			"updated_at":   time.Now().UTC(),
		})
	if res.Error != nil {
		return fmt.Errorf("failed to update billing document: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrBillingDocumentNotFound
	}
	return nil
}

func (r *postgresBillingDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	res := r.db.WithContext(queryCtx).Model(&domain.BillingDocument{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now().UTC())
	if res.Error != nil {
		return fmt.Errorf("failed to delete billing document: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrBillingDocumentNotFound
	}
	return nil
}

func (r *postgresBillingDocumentRepository) List(ctx context.Context, limit, offset int) ([]*domain.BillingDocument, int64, error) {
	limit, offset = clampRepoList(limit, offset)
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.BillingDocument{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count billing documents: %w", err)
	}
	var rows []domain.BillingDocument
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Model(&domain.BillingDocument{}).
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list billing documents: %w", err)
	}
	out := make([]*domain.BillingDocument, 0, len(rows))
	for i := range rows {
		row := rows[i]
		out = append(out, &row)
	}
	return out, total, nil
}

var _ ports.BillingDocumentRepository = (*postgresBillingDocumentRepository)(nil)
