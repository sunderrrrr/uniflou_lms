package tutor

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TutorID   uuid.UUID `gorm:"type:uuid" json:"tutor_id"`
	StudentID uuid.UUID `gorm:"type:uuid" json:"student_id"`
	Body      string    `json:"body"`
	NeedsCall bool      `json:"needs_call"`
	CreatedAt time.Time `json:"created_at"`
}

func (Comment) TableName() string { return "tutor_comments" }

type TutorRequest struct {
	ID             uuid.UUID `json:"id"`
	StudentID      uuid.UUID `json:"student_id"`
	StudentEmail   string    `json:"student_email"`
	RequesterEmail string    `json:"requester_email"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
}

type storedTutorRequest struct {
	ID        uuid.UUID       `gorm:"column:id"`
	Payload   json.RawMessage `gorm:"column:payload"`
	CreatedAt time.Time       `gorm:"column:sent_at"`
}
