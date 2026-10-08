package plan

import (
	"context"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListPlans(ctx context.Context, studentID uuid.UUID) ([]StudentPlan, error) {
	var items []StudentPlan
	err := r.db.WithContext(ctx).Preload("University").Where("student_id = ?", studentID).Find(&items).Error
	return items, err
}

func (r *Repository) ListUniversities(ctx context.Context) ([]University, error) {
	var items []University
	err := r.db.WithContext(ctx).Order("name asc").Find(&items).Error
	return items, err
}

func (r *Repository) AddPlan(ctx context.Context, studentID, universityID uuid.UUID, priority string) (*StudentPlan, error) {
	sp := StudentPlan{
		ID: uuid.New(), StudentID: studentID, UniversityID: universityID, Priority: priority, CreatedAt: time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(&sp).Error; err != nil {
		return nil, err
	}
	return &sp, nil
}

func (r *Repository) ReplacePlan(ctx context.Context, studentID uuid.UUID, universities []PlanSelection, deadlines []DeadlineInput) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("student_id = ?", studentID).Delete(&StudentPlan{}).Error; err != nil {
			return err
		}
		if err := tx.Where("student_id = ?", studentID).Delete(&Deadline{}).Error; err != nil {
			return err
		}
		for _, selection := range universities {
			item := StudentPlan{ID: uuid.New(), StudentID: studentID, UniversityID: selection.UniversityID, Priority: selection.Priority, CreatedAt: time.Now()}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		for _, deadline := range deadlines {
			item := Deadline{ID: uuid.New(), StudentID: studentID, Title: deadline.Title, Date: deadline.Date, Type: deadline.Type, CreatedAt: time.Now()}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) CountPlans(ctx context.Context, studentID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&StudentPlan{}).Where("student_id = ?", studentID).Count(&count).Error
	return count, err
}

func (r *Repository) HasPlan(ctx context.Context, studentID, universityID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&StudentPlan{}).Where("student_id = ? AND university_id = ?", studentID, universityID).Count(&count).Error
	return count > 0, err
}

func (r *Repository) UniversityExists(ctx context.Context, universityID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&University{}).Where("id = ?", universityID).Count(&count).Error
	return count == 1, err
}

func (r *Repository) DeletePlan(ctx context.Context, studentID, planID uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ? AND student_id = ?", planID, studentID).Delete(&StudentPlan{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *Repository) ListDeadlines(ctx context.Context, studentID uuid.UUID) ([]Deadline, error) {
	var items []Deadline
	err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Order("date asc").Find(&items).Error
	return items, err
}

func (r *Repository) AddDeadline(ctx context.Context, studentID uuid.UUID, input DeadlineInput) (*Deadline, error) {
	item := &Deadline{ID: uuid.New(), StudentID: studentID, Title: input.Title, Date: input.Date, Type: input.Type, CreatedAt: time.Now()}
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r *Repository) DeleteDeadline(ctx context.Context, studentID, deadlineID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("student_id = ? AND id = ?", studentID, deadlineID).Delete(&Deadline{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *Repository) NextDeadline(ctx context.Context, studentID uuid.UUID) (*Deadline, error) {
	var d Deadline
	err := r.db.WithContext(ctx).Where("student_id = ? AND date >= CURRENT_DATE", studentID).Order("date asc").First(&d).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &d, err
}
