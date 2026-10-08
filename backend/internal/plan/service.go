package plan

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

func (s *Service) GetPlan(ctx context.Context, studentID uuid.UUID) (map[string]any, error) {
	plans, err := s.repo.ListPlans(ctx, studentID)
	if err != nil {
		return nil, err
	}
	deadlines, err := s.repo.ListDeadlines(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return map[string]any{"plans": plans, "deadlines": deadlines}, nil
}

func (s *Service) Universities(ctx context.Context) ([]University, error) {
	return s.repo.ListUniversities(ctx)
}

func (s *Service) AddUniversity(ctx context.Context, studentID, universityID uuid.UUID, priority string) (*StudentPlan, error) {
	if priority == "" {
		priority = "realistic"
	}
	if !validPriority(priority) || universityID == uuid.Nil {
		return nil, apperrors.ErrValidation
	}
	count, err := s.repo.CountPlans(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, apperrors.ErrConflict
	}
	exists, err := s.repo.HasPlan(ctx, studentID, universityID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.ErrConflict
	}
	valid, err := s.repo.UniversityExists(ctx, universityID)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, apperrors.ErrValidation
	}
	return s.repo.AddPlan(ctx, studentID, universityID, priority)
}

func (s *Service) RemoveUniversity(ctx context.Context, studentID, planID uuid.UUID) error {
	return s.repo.DeletePlan(ctx, studentID, planID)
}

func (s *Service) Deadlines(ctx context.Context, studentID uuid.UUID) ([]Deadline, error) {
	return s.repo.ListDeadlines(ctx, studentID)
}

func (s *Service) SavePlan(ctx context.Context, studentID uuid.UUID, selections []PlanSelection, deadlines []DeadlineInput) error {
	if len(selections) > 5 {
		return apperrors.ErrValidation
	}
	seen := make(map[uuid.UUID]struct{}, len(selections))
	for i := range selections {
		if selections[i].UniversityID == uuid.Nil || !validPriority(selections[i].Priority) {
			return apperrors.ErrValidation
		}
		if _, ok := seen[selections[i].UniversityID]; ok {
			return apperrors.ErrValidation
		}
		seen[selections[i].UniversityID] = struct{}{}
		valid, err := s.repo.UniversityExists(ctx, selections[i].UniversityID)
		if err != nil {
			return err
		}
		if !valid {
			return apperrors.ErrValidation
		}
	}
	for i := range deadlines {
		deadlines[i].Title = strings.TrimSpace(deadlines[i].Title)
		if deadlines[i].Title == "" || len(deadlines[i].Title) > 255 || deadlines[i].Date.IsZero() || len(deadlines[i].Type) > 50 {
			return apperrors.ErrValidation
		}
		deadlines[i].Date = time.Date(deadlines[i].Date.Year(), deadlines[i].Date.Month(), deadlines[i].Date.Day(), 0, 0, 0, 0, time.UTC)
	}
	return s.repo.ReplacePlan(ctx, studentID, selections, deadlines)
}

func (s *Service) AddDeadline(ctx context.Context, studentID uuid.UUID, title, date, deadlineType string) (*Deadline, error) {
	title = strings.TrimSpace(title)
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil || title == "" || len(title) > 255 || len(deadlineType) > 50 {
		return nil, apperrors.ErrValidation
	}
	return s.repo.AddDeadline(ctx, studentID, DeadlineInput{Title: title, Date: parsedDate, Type: deadlineType})
}

func (s *Service) RemoveDeadline(ctx context.Context, studentID, deadlineID uuid.UUID) error {
	return s.repo.DeleteDeadline(ctx, studentID, deadlineID)
}

func validPriority(priority string) bool {
	switch priority {
	case "dream", "realistic", "safety":
		return true
	default:
		return false
	}
}
