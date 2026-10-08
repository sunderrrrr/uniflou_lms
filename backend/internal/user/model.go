package user

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type StudentProfile struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID       `gorm:"type:uuid;uniqueIndex" json:"user_id"`
	Grade          *int            `json:"grade"`
	TargetScore    *int            `json:"target_score"`
	UniversityName string          `json:"university_name"`
	SpecialtyName  string          `json:"specialty_name"`
	PassScore      int             `json:"pass_score"`
	Subjects       json.RawMessage `json:"subjects" gorm:"type:jsonb"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func (StudentProfile) TableName() string { return "student_profiles" }

type ParentLink struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ParentID  uuid.UUID `gorm:"type:uuid" json:"parent_id"`
	StudentID uuid.UUID `gorm:"type:uuid" json:"student_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (ParentLink) TableName() string { return "parent_links" }
