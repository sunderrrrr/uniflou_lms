package tutor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/auth"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Students(ctx context.Context) ([]auth.User, error) {
	var users []auth.User
	err := r.db.WithContext(ctx).Where("role = ?", auth.RoleStudent).Find(&users).Error
	return users, err
}

func (r *Repository) GetStudent(ctx context.Context, id uuid.UUID) (*auth.User, error) {
	var u auth.User
	if err := r.db.WithContext(ctx).Where("id = ? AND role = ?", id, auth.RoleStudent).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) AddComment(ctx context.Context, tutorID, studentID uuid.UUID, body string, needsCall bool) (*Comment, error) {
	c := Comment{
		ID: uuid.New(), TutorID: tutorID, StudentID: studentID,
		Body: body, NeedsCall: needsCall, CreatedAt: time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) LatestComments(ctx context.Context, studentID uuid.UUID, limit int) ([]Comment, error) {
	var items []Comment
	err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Order("created_at desc").Limit(limit).Find(&items).Error
	return items, err
}

func (r *Repository) ParentLinkActive(ctx context.Context, parentID, studentID uuid.UUID) bool {
	var count int64
	r.db.WithContext(ctx).Table("parent_links").
		Where("parent_id = ? AND student_id = ? AND status = ?", parentID, studentID, "active").
		Count(&count)
	return count > 0
}

func (r *Repository) CreateRequest(ctx context.Context, request TutorRequest) error {
	var recipients []auth.User
	if err := r.db.WithContext(ctx).Where("role IN ?", []string{auth.RoleTutor, auth.RoleAdmin}).Find(&recipients).Error; err != nil {
		return err
	}
	if len(recipients) == 0 {
		return apperrors.ErrNotFound
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}
	rows := make([]map[string]any, 0, len(recipients))
	for _, recipient := range recipients {
		rows = append(rows, map[string]any{
			"id": uuid.New(), "user_id": recipient.ID, "type": "tutor_request",
			"payload": string(payload), "sent_at": request.CreatedAt,
		})
	}
	return r.db.WithContext(ctx).Table("notifications").Create(&rows).Error
}

func (r *Repository) Requests(ctx context.Context, tutorID uuid.UUID) ([]TutorRequest, error) {
	var rows []storedTutorRequest
	err := r.db.WithContext(ctx).Table("notifications").
		Select("id, payload, sent_at").
		Where("user_id = ? AND type = ? AND read_at IS NULL", tutorID, "tutor_request").
		Order("sent_at DESC").Limit(100).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	requests := make([]TutorRequest, 0, len(rows))
	for _, row := range rows {
		var request TutorRequest
		if err := json.Unmarshal(row.Payload, &request); err != nil {
			return nil, err
		}
		request.ID = row.ID
		if !row.CreatedAt.IsZero() {
			request.CreatedAt = row.CreatedAt
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func (r *Repository) MarkRequestRead(ctx context.Context, tutorID, requestID uuid.UUID) error {
	result := r.db.WithContext(ctx).Table("notifications").
		Where("id = ? AND user_id = ? AND type = ?", requestID, tutorID, "tutor_request").
		Update("read_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
