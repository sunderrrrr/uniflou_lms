package admin

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

type Subject struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Code   string    `json:"code"`
	Topics []Topic   `json:"topics"`
}

type Topic struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Task struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	SubjectID     uuid.UUID       `gorm:"type:uuid" json:"subject_id"`
	TopicID       uuid.UUID       `gorm:"type:uuid" json:"topic_id"`
	TaskNumber    int             `json:"task_number"`
	Type          string          `json:"type"`
	Question      string          `json:"question"`
	Options       json.RawMessage `gorm:"type:jsonb" json:"options,omitempty"`
	CorrectAnswer string          `json:"correct_answer"`
	Explanation   string          `json:"explanation"`
	Source        string          `json:"source"`
	CreatedAt     time.Time       `json:"created_at"`
}

func (Task) TableName() string { return "tasks" }

type TaskInput struct {
	SubjectID     uuid.UUID       `json:"subject_id" binding:"required"`
	TopicID       uuid.UUID       `json:"topic_id" binding:"required"`
	TaskNumber    int             `json:"task_number" binding:"required,min=1"`
	Type          string          `json:"type" binding:"required"`
	Question      string          `json:"question" binding:"required"`
	Options       json.RawMessage `json:"options"`
	CorrectAnswer string          `json:"correct_answer" binding:"required"`
	Explanation   string          `json:"explanation"`
	Source        string          `json:"source"`
}

type SubjectInput struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type TopicInput struct {
	SubjectID uuid.UUID `json:"subject_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
}
