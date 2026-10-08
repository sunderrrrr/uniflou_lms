package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email         string    `gorm:"uniqueIndex;size:255" json:"email"`
	PasswordHash  string    `gorm:"size:255" json:"-"`
	Role          string    `gorm:"size:20" json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index"`
	TokenHash string    `gorm:"size:255"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

type PasswordResetToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	TokenHash string    `gorm:"size:255"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (PasswordResetToken) TableName() string { return "password_reset_tokens" }

const (
	RoleStudent = "student"
	RoleParent  = "parent"
	RoleTutor   = "tutor"
	RoleAdmin   = "admin"
)

var AllowedRoles = map[string]struct{}{
	RoleStudent: {},
	RoleParent:  {},
	RoleTutor:   {},
	RoleAdmin:   {},
}

func NormalizeRole(role string) string {
	switch role {
	case "", "student":
		return RoleStudent
	case RoleParent, RoleTutor, RoleAdmin:
		return role
	default:
		return ""
	}
}
