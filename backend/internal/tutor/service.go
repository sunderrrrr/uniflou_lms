package tutor

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/auth"
	"github.com/uniflou/uniflou-lms/backend/internal/progress"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/user"
	"gorm.io/gorm"
)

type Service struct {
	repo     *Repository
	progress *progress.Repository
	users    *user.Repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewRepository(db), progress: progress.NewRepository(db), users: user.NewRepository(db)}
}

func (s *Service) RequestHelp(ctx context.Context, actorID, requestedStudentID uuid.UUID, message string) error {
	actor, err := s.users.GetUser(ctx, actorID)
	if err != nil {
		return err
	}
	studentID := actorID
	switch actor.Role {
	case auth.RoleStudent:
		if requestedStudentID != uuid.Nil && requestedStudentID != actorID {
			return apperrors.ErrForbidden
		}
	case auth.RoleParent:
		if requestedStudentID == uuid.Nil || !s.repo.ParentLinkActive(ctx, actorID, requestedStudentID) {
			return apperrors.ErrForbidden
		}
		studentID = requestedStudentID
	default:
		return apperrors.ErrForbidden
	}
	student, err := s.repo.GetStudent(ctx, studentID)
	if err != nil {
		return apperrors.ErrNotFound
	}
	message = strings.TrimSpace(message)
	if message == "" {
		message = "Просьба связаться по поводу подготовки."
	}
	if len(message) > 1000 {
		return apperrors.ErrValidation
	}
	return s.repo.CreateRequest(ctx, TutorRequest{
		StudentID: student.ID, StudentEmail: student.Email,
		RequesterEmail: actor.Email, Message: message, CreatedAt: time.Now(),
	})
}

func (s *Service) Requests(ctx context.Context, tutorID uuid.UUID) ([]TutorRequest, error) {
	return s.repo.Requests(ctx, tutorID)
}

func (s *Service) MarkRequestRead(ctx context.Context, tutorID, requestID uuid.UUID) error {
	return s.repo.MarkRequestRead(ctx, tutorID, requestID)
}

func (s *Service) Students(ctx context.Context) ([]StudentCard, error) {
	users, err := s.repo.Students(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]StudentCard, len(users))
	for i, u := range users {
		out[i] = StudentCard{
			ID:          u.ID,
			Email:       u.Email,
			GoalPercent: s.progress.GoalPercent(ctx, u.ID, nil),
		}
	}
	return out, nil
}

func (s *Service) SetStudentTarget(ctx context.Context, actorID, studentID uuid.UUID, university, specialty string, passScore int) (*user.StudentProfile, error) {
	actor, err := s.users.GetUser(ctx, actorID)
	if err != nil {
		return nil, err
	}
	if actor.Role != auth.RoleAdmin && actor.Role != auth.RoleTutor {
		return nil, apperrors.ErrForbidden
	}
	if _, err := s.repo.GetStudent(ctx, studentID); err != nil {
		return nil, apperrors.ErrNotFound
	}
	if passScore < 0 || passScore > 100 || len(strings.TrimSpace(university)) > 255 || len(strings.TrimSpace(specialty)) > 255 {
		return nil, apperrors.ErrValidation
	}
	updates := map[string]any{}
	if university = strings.TrimSpace(university); university != "" {
		updates["university_name"] = university
	}
	if specialty = strings.TrimSpace(specialty); specialty != "" {
		updates["specialty_name"] = specialty
	}
	if passScore > 0 {
		updates["target_score"] = passScore
		updates["pass_score"] = passScore
	}
	if len(updates) == 0 {
		return nil, apperrors.ErrValidation
	}
	return s.users.UpdateProfile(ctx, studentID, updates)
}

type StudentCard struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	GoalPercent int       `json:"goal_percent"`
}

func (s *Service) StudentDetail(ctx context.Context, id uuid.UUID) (map[string]any, error) {
	u, err := s.repo.GetStudent(ctx, id)
	if err != nil {
		return nil, err
	}
	profile, err := s.users.GetOrCreateProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	weak, _ := s.progress.WeakTopics(ctx, id, 5)
	comments, _ := s.repo.LatestComments(ctx, id, 5)
	return map[string]any{
		"student":  u,
		"profile":  profile,
		"progress": s.progress.GoalPercent(ctx, id, nil),
		"weak":     weak,
		"comments": comments,
	}, nil
}

func (s *Service) Comment(ctx context.Context, tutorID, studentID uuid.UUID, body string, needsCall bool) (*Comment, error) {
	return s.repo.AddComment(ctx, tutorID, studentID, body, needsCall)
}
