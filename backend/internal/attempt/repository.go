package attempt

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

func (r *Repository) Create(ctx context.Context, studentID, taskID uuid.UUID, answer string, ok bool, timeSpent int) (*TaskAttempt, error) {
	a := &TaskAttempt{
		ID:        uuid.New(),
		StudentID: studentID,
		TaskID:    taskID,
		Answer:    answer,
		IsCorrect: ok,
		TimeSpent: timeSpent,
		CreatedAt: time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (r *Repository) GetByID(ctx context.Context, id, studentID uuid.UUID) (*TaskAttempt, error) {
	var a TaskAttempt
	if err := r.db.WithContext(ctx).Where("id = ? AND student_id = ?", id, studentID).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) ListByStudent(ctx context.Context, studentID uuid.UUID, p pagination.Params) ([]TaskAttempt, int64, error) {
	q := r.db.WithContext(ctx).Model(&TaskAttempt{}).Where("student_id = ?", studentID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []TaskAttempt
	if err := q.Order("created_at desc").Offset(p.Offset()).Limit(p.Limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
