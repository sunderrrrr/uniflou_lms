package ai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/uniflou/uniflou-lms/backend/internal/attempt"
	"github.com/uniflou/uniflou-lms/backend/internal/plan"
	"github.com/uniflou/uniflou-lms/backend/internal/progress"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/pagination"
	"github.com/uniflou/uniflou-lms/backend/internal/task"
	"github.com/uniflou/uniflou-lms/backend/internal/user"
	"gorm.io/gorm"
)

type Service struct {
	repo      *Repository
	client    *Client
	log       zerolog.Logger
	limitDay  int
	attempts  *attempt.Repository
	tasks     *task.Repository
	progress  *progress.Repository
	plans     *plan.Repository
	users     *user.Repository
	limitMu   sync.Mutex
	userLocks map[uuid.UUID]*sync.Mutex
}

func NewService(db *gorm.DB, client *Client, log zerolog.Logger, limitDay int) *Service {
	return &Service{
		repo:      NewRepository(db),
		client:    client,
		log:       log,
		limitDay:  limitDay,
		attempts:  attempt.NewRepository(db),
		tasks:     task.NewRepository(db),
		progress:  progress.NewRepository(db),
		plans:     plan.NewRepository(db),
		users:     user.NewRepository(db),
		userLocks: make(map[uuid.UUID]*sync.Mutex),
	}
}

func (s *Service) analyze(ctx context.Context, userID uuid.UUID, contextType string, contextID *uuid.UUID, prompt string) (*Report, error) {
	lock := s.lockForUser(userID)
	lock.Lock()
	defer lock.Unlock()
	if cached, err := s.repo.FindCached(ctx, userID, contextType, contextID, prompt); err == nil && cached != nil {
		return cached, nil
	}
	n, err := s.repo.CountToday(ctx, userID)
	if err != nil {
		return nil, err
	}
	if int(n) >= s.limitDay {
		return nil, apperrors.ErrRateLimit
	}
	resp, err := s.client.Complete(ctx, prompt)
	if err != nil {
		return nil, err
	}
	s.log.Info().Str("user_id", userID.String()).Str("context", contextType).Str("prompt", prompt).Str("response", resp).Msg("llm")
	rep := &Report{
		ID: uuid.New(), UserID: userID, ContextType: contextType, ContextID: contextID,
		Prompt: prompt, Response: resp, CreatedAt: time.Now(),
	}
	if err := s.repo.Save(ctx, rep); err != nil {
		return nil, err
	}
	return rep, nil
}

func (s *Service) lockForUser(userID uuid.UUID) *sync.Mutex {
	s.limitMu.Lock()
	defer s.limitMu.Unlock()
	lock := s.userLocks[userID]
	if lock == nil {
		lock = &sync.Mutex{}
		s.userLocks[userID] = lock
	}
	return lock
}

func (s *Service) AnalyzeTask(ctx context.Context, userID, attemptID uuid.UUID) (*Report, error) {
	a, err := s.attempts.GetByID(ctx, attemptID, userID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	t, err := s.tasks.GetByID(ctx, a.TaskID)
	if err != nil {
		return nil, err
	}
	prompt := TaskPrompt(t.Question, a.Answer, a.IsCorrect)
	return s.analyze(ctx, userID, "task", &attemptID, prompt)
}

func (s *Service) AnalyzeTrainingSession(ctx context.Context, userID uuid.UUID, attemptIDs []uuid.UUID, subject, taskNumber string) (*Report, error) {
	if len(attemptIDs) == 0 {
		return nil, apperrors.ErrNotFound
	}
	results := make([]string, 0, len(attemptIDs))
	correct := 0
	totalTime := 0
	for _, attemptID := range attemptIDs {
		a, err := s.attempts.GetByID(ctx, attemptID, userID)
		if err != nil {
			return nil, apperrors.ErrNotFound
		}
		t, err := s.tasks.GetByID(ctx, a.TaskID)
		if err != nil {
			return nil, err
		}
		status := "ошибка"
		if a.IsCorrect {
			status = "верно"
			correct++
		}
		totalTime += a.TimeSpent
		results = append(results, fmt.Sprintf("Задание: %s; ответ: %s; результат: %s; объяснение: %s", t.Question, a.Answer, status, t.Explanation))
	}
	accuracy := correct * 100 / len(attemptIDs)
	prompt := fmt.Sprintf("Ученик прошёл тренировочную сессию по предмету %s, номер задания %s. Всего заданий: %d, верно: %d, точность: %d%%, общее время: %d секунд. Ответы: %v. Проанализируй типичные ошибки, назови правило, которое нужно повторить, и дай конкретный план следующей тренировки. Не выдумывай данные вне отчёта.", subject, taskNumber, len(attemptIDs), correct, accuracy, totalTime, results)
	contextID := attemptIDs[len(attemptIDs)-1]
	return s.analyze(ctx, userID, "trainer_session", &contextID, prompt)
}

func (s *Service) AnalyzeDashboard(ctx context.Context, userID uuid.UUID) (*Report, error) {
	profile, err := s.users.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	goal := 0
	if profile.TargetScore != nil {
		goal = *profile.TargetScore
	}
	weak, err := s.progress.WeakTopics(ctx, userID, 5)
	if err != nil {
		return nil, err
	}
	prompt := DashboardPrompt(goal, fmt.Sprintf("%v", weak))
	return s.analyze(ctx, userID, "dashboard", nil, prompt)
}

func (s *Service) AnalyzePlan(ctx context.Context, userID uuid.UUID) (*Report, error) {
	plans, err := s.plans.ListPlans(ctx, userID)
	if err != nil {
		return nil, err
	}
	deadlines, err := s.plans.ListDeadlines(ctx, userID)
	if err != nil {
		return nil, err
	}
	summary := fmt.Sprintf("Вузы: %v. Дедлайны: %v", plans, deadlines)
	return s.analyze(ctx, userID, "plan", nil, PlanPrompt(summary))
}

func (s *Service) AnalyzeProgress(ctx context.Context, userID uuid.UUID) (*Report, error) {
	items, err := s.progress.ListByStudent(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.analyze(ctx, userID, "progress", nil, ProgressPrompt(fmt.Sprintf("%v", items)))
}

func (s *Service) Reports(ctx context.Context, userID uuid.UUID, p pagination.Params) ([]Report, int64, error) {
	return s.repo.List(ctx, userID, p)
}
