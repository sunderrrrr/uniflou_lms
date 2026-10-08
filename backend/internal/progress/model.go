package progress

import (
	"time"

	"github.com/google/uuid"
)

type Progress struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID uuid.UUID `gorm:"type:uuid;index" json:"student_id"`
	SubjectID uuid.UUID `gorm:"type:uuid" json:"subject_id"`
	TopicID   uuid.UUID `gorm:"type:uuid" json:"topic_id"`
	Score     int       `json:"score"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Progress) TableName() string { return "progress" }

type Achievement struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID  uuid.UUID `gorm:"type:uuid" json:"student_id"`
	Code       string    `json:"code"`
	UnlockedAt time.Time `json:"unlocked_at"`
}

func (Achievement) TableName() string { return "achievements" }

type WeeklyPoint struct {
	Week  string `json:"week"`
	Score int    `json:"score"`
}

type DashboardSummary struct {
	GoalPercent   int           `json:"goal_percent"`
	TodayTasks    []uuid.UUID   `json:"today_task_ids"`
	WeeklyScores  []WeeklyPoint `json:"weekly_scores"`
	NextDeadline  any           `json:"next_deadline"`
}

type MockExamResult struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID  uuid.UUID `gorm:"type:uuid;index" json:"student_id"`
	SubjectID  uuid.UUID `gorm:"type:uuid" json:"subject_id"`
	Title      string    `json:"title"`
	TotalTasks int       `json:"total_tasks"`
	Score      int       `json:"score"`
	CreatedAt  time.Time `json:"created_at"`
}

func (MockExamResult) TableName() string { return "mock_exam_results" }
