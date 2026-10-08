package progress

import (
	"context"

	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/plan"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/user"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	user *user.Repository
	plan *plan.Repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewRepository(db), user: user.NewRepository(db), plan: plan.NewRepository(db)}
}

func (s *Service) Get(ctx context.Context, studentID uuid.UUID) (map[string]any, error) {
	profile, _ := s.user.GetOrCreateProfile(ctx, studentID)
	var target *int
	if profile != nil {
		target = profile.TargetScore
	}
	items, err := s.repo.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"goal_percent": s.repo.GoalPercent(ctx, studentID, target),
		"items":        items,
	}, nil
}

func (s *Service) Subjects(ctx context.Context, studentID uuid.UUID) ([]Progress, error) {
	return s.repo.ListByStudent(ctx, studentID)
}

func (s *Service) Topics(ctx context.Context, studentID uuid.UUID) ([]Progress, error) {
	return s.repo.ListByStudent(ctx, studentID)
}

func (s *Service) Weak(ctx context.Context, studentID uuid.UUID) ([]Progress, error) {
	return s.repo.WeakTopics(ctx, studentID, 5)
}

func (s *Service) CreateMockResult(ctx context.Context, studentID, subjectID uuid.UUID, title string, totalTasks, score int) (*MockExamResult, error) {
	if title == "" || len(title) > 255 || totalTasks < 1 || score < 0 || score > 100 || subjectID == uuid.Nil {
		return nil, apperrors.ErrValidation
	}
	exists, err := s.repo.SubjectExists(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrValidation
	}
	return s.repo.SaveMockResult(ctx, studentID, subjectID, title, totalTasks, score)
}

func (s *Service) MockResults(ctx context.Context, studentID uuid.UUID) ([]MockExamResult, error) {
	return s.repo.ListMockResults(ctx, studentID)
}

func (s *Service) Dashboard(ctx context.Context, studentID uuid.UUID) (*DashboardSummary, error) {
	profile, _ := s.user.GetOrCreateProfile(ctx, studentID)
	var target *int
	if profile != nil {
		target = profile.TargetScore
	}
	weak, _ := s.repo.WeakTopics(ctx, studentID, 3)
	taskIDs := make([]uuid.UUID, 0, len(weak))
	for _, w := range weak {
		taskIDs = append(taskIDs, w.TopicID)
	}
	nextDeadline, _ := s.plan.NextDeadline(ctx, studentID)
	var deadline any
	if nextDeadline != nil {
		deadline = map[string]string{"title": nextDeadline.Title, "date": nextDeadline.Date.Format("2006-01-02")}
	}
	return &DashboardSummary{
		GoalPercent:  min(100, s.repo.GoalPercent(ctx, studentID, target)),
		TodayTasks:   taskIDs,
		WeeklyScores: s.repo.WeeklyScores(ctx, studentID),
		NextDeadline: deadline,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
