package task

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/attempt"
	"github.com/uniflou/uniflou-lms/backend/internal/progress"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

type Service struct {
	repo     *Repository
	attempts *attempt.Repository
	progress *progress.Repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		repo:     NewRepository(db),
		attempts: attempt.NewRepository(db),
		progress: progress.NewRepository(db),
	}
}

func (s *Service) List(ctx context.Context, topicID *uuid.UUID, p pagination.Params) ([]TaskPublic, int64, error) {
	items, total, err := s.repo.List(ctx, topicID, p)
	if err != nil {
		return nil, 0, err
	}
	out := make([]TaskPublic, len(items))
	for i := range items {
		out[i] = ToPublic(&items[i])
	}
	return out, total, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*TaskPublic, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	p := ToPublic(t)
	return &p, nil
}

func (s *Service) Next(ctx context.Context, studentID uuid.UUID, mode, subjectIDStr, topicIDStr string, limit int) (*TaskPublic, error) {
	var subjectID *uuid.UUID
	if subjectIDStr != "" {
		id, err := uuid.Parse(subjectIDStr)
		if err == nil {
			subjectID = &id
		}
	}
	var topicID *uuid.UUID
	if topicIDStr != "" {
		id, err := uuid.Parse(topicIDStr)
		if err == nil {
			topicID = &id
		}
	}
	if mode == "" {
		mode = "weak"
	}
	items, err := s.repo.NextForStudent(ctx, studentID, mode, subjectID, topicID, limit)
	if err != nil {
		return nil, err
	}
	p := ToPublic(&items[0])
	return &p, nil
}

func (s *Service) NextBatch(ctx context.Context, studentID uuid.UUID, mode, subjectIDStr, topicIDStr string, limit int) ([]TaskPublic, error) {
	var subjectID *uuid.UUID
	if subjectIDStr != "" {
		id, err := uuid.Parse(subjectIDStr)
		if err == nil {
			subjectID = &id
		}
	}
	var topicID *uuid.UUID
	if topicIDStr != "" {
		id, err := uuid.Parse(topicIDStr)
		if err == nil {
			topicID = &id
		}
	}
	if limit <= 0 {
		limit = 5
	}
	if limit > 30 {
		limit = 30
	}
	items, err := s.repo.NextForStudent(ctx, studentID, mode, subjectID, topicID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]TaskPublic, 0, len(items))
	for i := range items {
		p := ToPublic(&items[i])
		out = append(out, p)
	}
	return out, nil
}

func (s *Service) TrainerCatalog(ctx context.Context) ([]TrainerSubject, error) {
	return s.repo.TrainerCatalog(ctx)
}

func (s *Service) TrainerTopics(ctx context.Context) ([]TrainerTopic, error) {
	return s.repo.TrainerTopics(ctx)
}

func (s *Service) MockSession(ctx context.Context, subjectID uuid.UUID) ([]TaskPublic, error) {
	items, err := s.repo.MockSession(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	out := make([]TaskPublic, 0, len(items))
	for i := range items {
		out = append(out, ToPublic(&items[i]))
	}
	return out, nil
}

func (s *Service) TrainingBatch(ctx context.Context, subjectID uuid.UUID, taskNumber, limit int) ([]TaskPublic, error) {
	items, err := s.repo.TrainingBatch(ctx, subjectID, taskNumber, limit)
	if err != nil {
		return nil, err
	}
	out := make([]TaskPublic, 0, len(items))
	for i := range items {
		out = append(out, ToPublic(&items[i]))
	}
	return out, nil
}

type AttemptResult struct {
	Attempt     attempt.TaskAttempt `json:"attempt"`
	IsCorrect   bool                `json:"is_correct"`
	Explanation string              `json:"explanation"`
}

func (s *Service) SubmitAttempt(ctx context.Context, studentID, taskID uuid.UUID, answer string, timeSpent int) (*AttemptResult, error) {
	t, err := s.repo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	ok := normalize(answer) == normalize(t.CorrectAnswer)
	a, err := s.attempts.Create(ctx, studentID, taskID, answer, ok, timeSpent)
	if err != nil {
		return nil, err
	}
	if err := s.progress.UpsertFromAttempt(ctx, studentID, t.SubjectID, t.TopicID, ok); err != nil {
		return nil, err
	}
	return &AttemptResult{Attempt: *a, IsCorrect: ok, Explanation: t.Explanation}, nil
}

func (s *Service) History(ctx context.Context, studentID uuid.UUID, p pagination.Params) ([]attempt.TaskAttempt, int64, error) {
	return s.attempts.ListByStudent(ctx, studentID, p)
}

func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
