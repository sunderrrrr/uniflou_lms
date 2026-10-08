package ai

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;index" json:"user_id"`
	ContextType string     `json:"context_type"`
	ContextID   *uuid.UUID `gorm:"type:uuid" json:"context_id,omitempty"`
	Prompt      string     `json:"prompt"`
	Response    string     `json:"response"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (Report) TableName() string { return "ai_reports" }
