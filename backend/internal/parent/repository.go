package parent

import (
	"context"

	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/auth"
	"github.com/uniflou/uniflou-lms/backend/internal/user"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Children(ctx context.Context, parentID uuid.UUID) ([]auth.User, error) {
	var links []user.ParentLink
	if err := r.db.WithContext(ctx).Where("parent_id = ? AND status = ?", parentID, "active").Find(&links).Error; err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return []auth.User{}, nil
	}
	ids := make([]uuid.UUID, len(links))
	for i, l := range links {
		ids[i] = l.StudentID
	}
	var users []auth.User
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	return users, err
}

func (r *Repository) LinkActive(ctx context.Context, parentID, studentID uuid.UUID) bool {
	var n int64
	r.db.WithContext(ctx).Model(&user.ParentLink{}).
		Where("parent_id = ? AND student_id = ? AND status = ?", parentID, studentID, "active").
		Count(&n)
	return n > 0
}
