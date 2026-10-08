package task

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	SubjectID     uuid.UUID       `gorm:"type:uuid" json:"subject_id"`
	TopicID       uuid.UUID       `gorm:"type:uuid" json:"topic_id"`
	TaskNumber    int             `json:"task_number"`
	Type          string          `json:"type"`
	Question      string          `json:"question"`
	Options       json.RawMessage `json:"options,omitempty" gorm:"type:jsonb"`
	CorrectAnswer string          `json:"-"`
	Explanation   string          `json:"explanation,omitempty"`
	Source        string          `json:"source,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

func (Task) TableName() string { return "tasks" }

type TaskPublic struct {
	ID         uuid.UUID       `json:"id"`
	SubjectID  uuid.UUID       `json:"subject_id"`
	TopicID    uuid.UUID       `json:"topic_id"`
	TaskNumber int             `json:"task_number"`
	Type       string          `json:"type"`
	Question   string          `json:"question"`
	Options    json.RawMessage `json:"options,omitempty"`
}

func ToPublic(t *Task) TaskPublic {
	return TaskPublic{
		ID: t.ID, SubjectID: t.SubjectID, TopicID: t.TopicID,
		TaskNumber: t.TaskNumber, Type: t.Type, Question: t.Question, Options: t.Options,
	}
}

type TrainerSubject struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name"`
	Code        string              `json:"code"`
	TaskNumbers []TrainerTaskNumber `json:"task_numbers"`
}

type TrainerTopic struct {
	ID          uuid.UUID `json:"id"`
	SubjectID   uuid.UUID `json:"subject_id"`
	SubjectName string    `json:"subject_name"`
	Name        string    `json:"name"`
	Count       int       `json:"count"`
}

type TrainerTaskNumber struct {
	Number int `json:"number"`
	Count  int `json:"count"`
}
