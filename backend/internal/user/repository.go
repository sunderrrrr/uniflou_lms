package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/auth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (*auth.User, error) {
	var u auth.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetOrCreateProfile(ctx context.Context, userID uuid.UUID) (*StudentProfile, error) {
	var p StudentProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&p).Error
	if err == nil {
		return &p, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	p = StudentProfile{ID: uuid.New(), UserID: userID}
	if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) LinkParentStudent(ctx context.Context, parentID, studentID uuid.UUID) error {
	link := ParentLink{
		ID: uuid.New(), ParentID: parentID, StudentID: studentID, Status: "active", CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "parent_id"}, {Name: "student_id"}},
		DoUpdates: clause.Assignments(map[string]any{"status": "active"}),
	}).Create(&link).Error
}

func (r *Repository) FindStudentByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	var u auth.User
	if err := r.db.WithContext(ctx).Where("email = ? AND role = ?", email, auth.RoleStudent).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, apperrors.ErrNotFound
		}
		return uuid.Nil, err
	}
	return u.ID, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]any) (*StudentProfile, error) {
	p, err := r.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(p).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetOrCreateProfile(ctx, userID)
}
