package progress

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) UpsertFromAttempt(ctx context.Context, studentID, subjectID, topicID uuid.UUID, correct bool) error {
	var stats struct {
		Total   int64
		Correct int64
	}
	if err := r.db.WithContext(ctx).Table("task_attempts").
		Select("COUNT(*) AS total, COUNT(*) FILTER (WHERE task_attempts.is_correct) AS correct").
		Joins("JOIN tasks ON tasks.id = task_attempts.task_id").
		Where("task_attempts.student_id = ? AND tasks.topic_id = ?", studentID, topicID).
		Scan(&stats).Error; err != nil {
		return err
	}
	score := 0
	if stats.Total > 0 {
		score = int(stats.Correct * 100 / stats.Total)
	}
	p := Progress{
		ID:        uuid.New(),
		StudentID: studentID,
		SubjectID: subjectID,
		TopicID:   topicID,
		Score:     score,
		UpdatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "student_id"}, {Name: "topic_id"}},
		DoUpdates: clause.Assignments(map[string]any{"score": score, "updated_at": time.Now()}),
	}).Create(&p).Error
}

func (r *Repository) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]Progress, error) {
	var items []Progress
	err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Find(&items).Error
	return items, err
}

func (r *Repository) WeakTopics(ctx context.Context, studentID uuid.UUID, limit int) ([]Progress, error) {
	var items []Progress
	err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Order("score asc").Limit(limit).Find(&items).Error
	return items, err
}

func (r *Repository) SaveMockResult(ctx context.Context, studentID, subjectID uuid.UUID, title string, totalTasks, score int) (*MockExamResult, error) {
	item := &MockExamResult{ID: uuid.New(), StudentID: studentID, SubjectID: subjectID, Title: title, TotalTasks: totalTasks, Score: score, CreatedAt: time.Now()}
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r *Repository) ListMockResults(ctx context.Context, studentID uuid.UUID) ([]MockExamResult, error) {
	var items []MockExamResult
	err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *Repository) SubjectExists(ctx context.Context, subjectID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("subjects").Where("id = ?", subjectID).Count(&count).Error
	return count == 1, err
}

func (r *Repository) GoalPercent(ctx context.Context, studentID uuid.UUID, target *int) int {
	var avg float64
	r.db.WithContext(ctx).Model(&Progress{}).Where("student_id = ?", studentID).Select("COALESCE(AVG(score),0)").Scan(&avg)
	if target != nil && *target > 0 {
		return int(avg * 100 / float64(*target))
	}
	return int(avg)
}

func (r *Repository) WeeklyScores(ctx context.Context, studentID uuid.UUID) []WeeklyPoint {
	type weeklyResult struct {
		Week  time.Time
		Score float64
	}
	var results []weeklyResult
	start := time.Now().AddDate(0, 0, -56)
	if err := r.db.WithContext(ctx).Table("task_attempts").
		Select("date_trunc('week', created_at) AS week, AVG(CASE WHEN is_correct THEN 100.0 ELSE 0.0 END) AS score").
		Where("student_id = ? AND created_at >= ?", studentID, start).
		Group("date_trunc('week', created_at)").Order("week asc").Scan(&results).Error; err != nil {
		return []WeeklyPoint{}
	}
	points := make([]WeeklyPoint, 0, len(results))
	for _, result := range results {
		points = append(points, WeeklyPoint{Week: result.Week.Format("02.01"), Score: int(result.Score)})
	}
	return points
}

func (r *Repository) WeeklyAttemptStats(ctx context.Context, studentID uuid.UUID) (total int64, correct int64, err error) {
	now := time.Now()
	weekStart := now.Truncate(24*time.Hour).AddDate(0, 0, -int(now.Weekday())+1)
	if now.Weekday() == time.Sunday {
		weekStart = weekStart.AddDate(0, 0, -7)
	}
	query := r.db.WithContext(ctx).Table("task_attempts").Where("student_id = ? AND created_at >= ?", studentID, weekStart)
	if err := query.Count(&total).Error; err != nil {
		return 0, 0, err
	}
	err = query.Where("is_correct = TRUE").Count(&correct).Error
	return total, correct, err
}
