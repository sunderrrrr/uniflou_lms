package plan

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type University struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string          `json:"name"`
	City         string          `json:"city"`
	TargetScores json.RawMessage `json:"target_scores" gorm:"type:jsonb"`
	CreatedAt    time.Time       `json:"created_at"`
}

func (University) TableName() string { return "universities" }

type StudentPlan struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID    uuid.UUID  `gorm:"type:uuid" json:"student_id"`
	UniversityID uuid.UUID  `gorm:"type:uuid" json:"university_id"`
	Priority     string     `json:"priority"`
	CreatedAt    time.Time  `json:"created_at"`
	University   University `gorm:"foreignKey:UniversityID" json:"university,omitempty"`
}

func (StudentPlan) TableName() string { return "student_plans" }

type Deadline struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID uuid.UUID `gorm:"type:uuid" json:"student_id"`
	Title     string    `json:"title"`
	Date      time.Time `gorm:"type:date" json:"date"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func (Deadline) TableName() string { return "deadlines" }

type PlanSelection struct {
	UniversityID uuid.UUID `json:"university_id"`
	Priority     string    `json:"priority"`
}

type DeadlineInput struct {
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
	Type  string    `json:"type"`
}
