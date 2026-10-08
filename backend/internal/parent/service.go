package parent

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/plan"
	"github.com/uniflou/uniflou-lms/backend/internal/progress"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/tutor"
	"gorm.io/gorm"
)

type Service struct {
	repo     *Repository
	progress *progress.Repository
	plan     *plan.Repository
	tutor    *tutor.Repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		repo:     NewRepository(db),
		progress: progress.NewRepository(db),
		plan:     plan.NewRepository(db),
		tutor:    tutor.NewRepository(db),
	}
}

func (s *Service) Children(ctx context.Context, parentID uuid.UUID) (any, error) {
	return s.repo.Children(ctx, parentID)
}

func (s *Service) ChildProgress(ctx context.Context, parentID, studentID uuid.UUID) (map[string]any, error) {
	if !s.repo.LinkActive(ctx, parentID, studentID) {
		return nil, apperrors.ErrForbidden
	}
	items, err := s.progress.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}
	total, _, err := s.progress.WeeklyAttemptStats(ctx, studentID)
	if err != nil {
		return nil, err
	}
	summary := "Пока нет данных о решённых заданиях."
	if total > 0 {
		summary = "За эту неделю решено заданий: " + fmt.Sprint(total) + "."
	}
	return map[string]any{
		"goal_percent": s.progress.GoalPercent(ctx, studentID, nil),
		"summary":      summary,
		"topics":       items,
	}, nil
}

func (s *Service) ChildReports(ctx context.Context, parentID, studentID uuid.UUID) (map[string]any, error) {
	if !s.repo.LinkActive(ctx, parentID, studentID) {
		return nil, apperrors.ErrForbidden
	}
	deadlines, err := s.plan.ListDeadlines(ctx, studentID)
	if err != nil {
		return nil, err
	}
	comments, err := s.tutor.LatestComments(ctx, studentID, 3)
	if err != nil {
		return nil, err
	}
	total, correct, err := s.progress.WeeklyAttemptStats(ctx, studentID)
	if err != nil {
		return nil, err
	}
	report := []string{}
	if total > 0 {
		report = append(report, fmt.Sprintf("За неделю решено заданий: %d", total))
		report = append(report, fmt.Sprintf("Верных ответов: %d из %d", correct, total))
	}
	return map[string]any{
		"weekly_report": report,
		"deadlines": deadlines,
		"tutor_notes": comments,
	}, nil
}
