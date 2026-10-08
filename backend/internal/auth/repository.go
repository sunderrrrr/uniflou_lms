package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *Repository) EnsureStudentProfile(ctx context.Context, userID uuid.UUID) error {
	type row struct {
		ID uuid.UUID `gorm:"type:uuid;primaryKey"`
	}
	var existing row
	err := r.db.WithContext(ctx).Table("student_profiles").Where("user_id = ?", userID).First(&existing).Error
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	target := 80
	if err := r.db.WithContext(ctx).Exec(
		`INSERT INTO student_profiles (id, user_id, target_score, subjects, created_at, updated_at) VALUES (gen_random_uuid(), ?, ?, '[]', NOW(), NOW())`,
		userID, target,
	).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO deadlines (id, student_id, title, date, type, created_at) VALUES (gen_random_uuid(), ?, 'Пробный ЕГЭ', CURRENT_DATE + INTERVAL '30 days', 'exam', NOW())`,
		userID,
	).Error
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) SetUserRole(ctx context.Context, userID uuid.UUID, role string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(map[string]any{
		"role": role, "email_verified": true, "updated_at": time.Now(),
	}).Error
}

func (r *Repository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	rt := RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(&rt).Error
}

func (r *Repository) FindRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string) error {
	var rt RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND token_hash = ? AND expires_at > ?", userID, tokenHash, time.Now()).
		First(&rt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.ErrUnauthorized
		}
		return err
	}
	return nil
}

func (r *Repository) DeleteRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}

func (r *Repository) SavePasswordReset(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_ = r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&PasswordResetToken{}).Error
	pr := PasswordResetToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(&pr).Error
}

func (r *Repository) ConsumePasswordReset(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var token PasswordResetToken
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).First(&token).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return apperrors.ErrNotFound
			}
			return err
		}
		if err := tx.Delete(&token).Error; err != nil {
			return err
		}
		userID = token.UserID
		return nil
	})
	return userID, err
}

func (r *Repository) UpdatePassword(ctx context.Context, userID uuid.UUID, hash string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&User{}).Where("id = ?", userID).Update("password_hash", hash)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return apperrors.ErrNotFound
		}
		return tx.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
	})
}
