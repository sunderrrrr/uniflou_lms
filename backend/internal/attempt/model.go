package attempt

import (
	"time"

	"github.com/google/uuid"
)

type TaskAttempt struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID uuid.UUID `gorm:"type:uuid;index" json:"student_id"`
	TaskID    uuid.UUID `gorm:"type:uuid;index" json:"task_id"`
	Answer    string    `json:"answer"`
	IsCorrect bool      `json:"is_correct"`
	TimeSpent int       `json:"time_spent"`
	CreatedAt time.Time `json:"created_at"`
}

func (TaskAttempt) TableName() string { return "task_attempts" }
