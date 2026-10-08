package ai

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, report *Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *Repository) CountToday(ctx context.Context, userID uuid.UUID) (int64, error) {
	var n int64
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	err := r.db.WithContext(ctx).Model(&Report{}).
		Where("user_id = ? AND created_at >= ?", userID, start).
		Count(&n).Error
	return n, err
}

func (r *Repository) FindCached(ctx context.Context, userID uuid.UUID, contextType string, contextID *uuid.UUID, prompt string) (*Report, error) {
	q := r.db.WithContext(ctx).Where("user_id = ? AND context_type = ? AND prompt = ?", userID, contextType, prompt)
	if contextID != nil {
		q = q.Where("context_id = ?", *contextID)
	}
	var rep Report
	err := q.Order("created_at desc").First(&rep).Error
	if err != nil {
		return nil, err
	}
	if time.Since(rep.CreatedAt) > 6*time.Hour {
		return nil, gorm.ErrRecordNotFound
	}
	return &rep, nil
}

func (r *Repository) List(ctx context.Context, userID uuid.UUID, p pagination.Params) ([]Report, int64, error) {
	q := r.db.WithContext(ctx).Model(&Report{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []Report
	if err := q.Order("created_at desc").Offset(p.Offset()).Limit(p.Limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
