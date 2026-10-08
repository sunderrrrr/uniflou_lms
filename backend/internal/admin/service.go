package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/auth"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

var errLastAdmin = errors.New("cannot remove the last administrator")

type Service struct {
	repo *Repository
}

func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepository(db)} }

func (s *Service) Users(ctx context.Context, p pagination.Params) ([]User, int64, error) {
	return s.repo.ListUsers(ctx, p)
}

func (s *Service) UpdateUserRole(ctx context.Context, actorID, userID uuid.UUID, role string) error {
	if actorID == userID {
		return apperrors.ErrForbidden
	}
	role = auth.NormalizeRole(role)
	if role == "" {
		return apperrors.ErrValidation
	}
	err := s.repo.UpdateUserRole(ctx, actorID, userID, role)
	return mapAdminError(err)
}

func (s *Service) DeleteUser(ctx context.Context, actorID, userID uuid.UUID) error {
	if actorID == userID {
		return apperrors.ErrForbidden
	}
	return mapAdminError(s.repo.DeleteUser(ctx, actorID, userID))
}

func (s *Service) Subjects(ctx context.Context) ([]Subject, error) { return s.repo.Subjects(ctx) }

func (s *Service) CreateSubject(ctx context.Context, input SubjectInput) (*Subject, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Code = strings.ToLower(strings.TrimSpace(input.Code))
	if input.Name == "" || len(input.Name) > 100 || input.Code == "" || len(input.Code) > 20 {
		return nil, apperrors.ErrValidation
	}
	for _, char := range input.Code {
		if !(char >= 'a' && char <= 'z' || char >= '0' && char <= '9' || char == '_' || char == '-') {
			return nil, apperrors.ErrValidation
		}
	}
	return s.repo.CreateSubject(ctx, input)
}

func (s *Service) CreateTopic(ctx context.Context, input TopicInput) (*Topic, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.SubjectID == uuid.Nil || input.Name == "" || len(input.Name) > 255 {
		return nil, apperrors.ErrValidation
	}
	exists, err := s.repo.SubjectExists(ctx, input.SubjectID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrValidation
	}
	return s.repo.CreateTopic(ctx, input)
}

func (s *Service) Tasks(ctx context.Context, p pagination.Params) ([]Task, int64, error) {
	return s.repo.ListTasks(ctx, p)
}

func (s *Service) CreateTask(ctx context.Context, input TaskInput) (*Task, error) {
	if err := s.validateTask(ctx, &input); err != nil {
		return nil, err
	}
	item := &Task{
		ID: uuid.New(), SubjectID: input.SubjectID, TopicID: input.TopicID, TaskNumber: input.TaskNumber,
		Type: input.Type, Question: strings.TrimSpace(input.Question), Options: input.Options,
		CorrectAnswer: strings.TrimSpace(input.CorrectAnswer), Explanation: strings.TrimSpace(input.Explanation),
		Source: strings.TrimSpace(input.Source), CreatedAt: time.Now(),
	}
	if err := s.repo.SaveTask(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateTask(ctx context.Context, id uuid.UUID, input TaskInput) (*Task, error) {
	if err := s.validateTask(ctx, &input); err != nil {
		return nil, err
	}
	if err := mapAdminError(s.repo.UpdateTask(ctx, id, input)); err != nil {
		return nil, err
	}
	return s.repo.GetTask(ctx, id)
}

func (s *Service) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return mapAdminError(s.repo.DeleteTask(ctx, id))
}

func (s *Service) validateTask(ctx context.Context, input *TaskInput) error {
	input.Question = strings.TrimSpace(input.Question)
	input.CorrectAnswer = strings.TrimSpace(input.CorrectAnswer)
	if input.SubjectID == uuid.Nil || input.TopicID == uuid.Nil || input.TaskNumber < 1 || input.Question == "" || input.CorrectAnswer == "" || len(input.CorrectAnswer) > 255 {
		return apperrors.ErrValidation
	}
	if input.Type != "single_choice" && input.Type != "short_answer" {
		return apperrors.ErrValidation
	}
	validTopic, err := s.repo.TopicMatchesSubject(ctx, input.TopicID, input.SubjectID)
	if err != nil {
		return err
	}
	if !validTopic {
		return apperrors.ErrValidation
	}
	if len(input.Options) > 0 && !json.Valid(input.Options) {
		return apperrors.ErrValidation
	}
	if input.Type == "single_choice" {
		var options []string
		if len(input.Options) == 0 || json.Unmarshal(input.Options, &options) != nil || len(options) < 2 {
			return apperrors.ErrValidation
		}
		found := false
		for _, option := range options {
			if strings.EqualFold(strings.TrimSpace(option), input.CorrectAnswer) {
				found = true
				break
			}
		}
		if !found {
			return apperrors.ErrValidation
		}
	}
	return nil
}

func mapAdminError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return apperrors.ErrNotFound
	case errors.Is(err, errLastAdmin):
		return apperrors.ErrConflict
	default:
		return err
	}
}
