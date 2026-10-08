package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/pkg/hash"
	"gorm.io/gorm"
)

type Mailer interface {
	SendResetPassword(to, token string) error
}

type noopMailer struct{}

func (noopMailer) SendResetPassword(to, token string) error { return nil }

type Service struct {
	repo   *Repository
	jwt    *JWTManager
	log    zerolog.Logger
	mailer Mailer
}

func NewService(db *gorm.DB, jwt *JWTManager, log zerolog.Logger, mailer Mailer) *Service {
	if mailer == nil {
		mailer = noopMailer{}
	}
	return &Service{repo: NewRepository(db), jwt: jwt, log: log, mailer: mailer}
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *Service) Register(ctx context.Context, email, password, role string) (*TokenPair, *User, error) {
	normalizedRole := NormalizeRole(role)
	if normalizedRole == "" {
		return nil, nil, fmt.Errorf("%w: допустимые роли: student, parent", apperrors.ErrValidation)
	}
	if normalizedRole == RoleAdmin || normalizedRole == RoleTutor {
		return nil, nil, fmt.Errorf("%w: создание наставника или администратора доступно только через внутренний админ-процесс", apperrors.ErrForbidden)
	}
	if err := ValidatePassword(password); err != nil {
		return nil, nil, fmt.Errorf("%w: пароль должен быть не короче 8 символов и содержать буквы и цифры", apperrors.ErrValidation)
	}
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return nil, nil, apperrors.ErrConflict
	} else if err != apperrors.ErrNotFound {
		return nil, nil, err
	}

	pwHash, err := hash.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	user := &User{
		ID:            uuid.New(),
		Email:         email,
		PasswordHash:  pwHash,
		Role:          normalizedRole,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, nil, err
	}
	if user.Role == RoleStudent {
		_ = s.repo.EnsureStudentProfile(ctx, user.ID)
	}

	pair, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	s.log.Info().Str("user_id", user.ID.String()).Str("event", "register").Msg("auth")
	return pair, user, nil
}

func (s *Service) CreateTutor(ctx context.Context, email, password string) (*User, error) {
	if err := ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("%w: пароль наставника должен содержать буквы и цифры", apperrors.ErrValidation)
	}
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return nil, apperrors.ErrConflict
	} else if err != apperrors.ErrNotFound {
		return nil, err
	}
	pwHash, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}
	tutor := &User{ID: uuid.New(), Email: email, PasswordHash: pwHash, Role: RoleTutor, EmailVerified: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.repo.CreateUser(ctx, tutor); err != nil {
		return nil, err
	}
	return tutor, nil
}

func (s *Service) EnsureBootstrapAdmin(ctx context.Context, email, password string) error {
	if email == "" || password == "" {
		return nil
	}
	if err := ValidatePassword(password); err != nil {
		return fmt.Errorf("%w: пароль администратора должен содержать буквы и цифры", apperrors.ErrValidation)
	}
	if existing, err := s.repo.FindByEmail(ctx, email); err == nil {
		if existing.Role == RoleAdmin {
			return nil
		}
		if !hash.CheckPassword(existing.PasswordHash, password) {
			return fmt.Errorf("%w: существующий аккаунт не совпадает с bootstrap-паролем", apperrors.ErrConflict)
		}
		return s.repo.SetUserRole(ctx, existing.ID, RoleAdmin)
	} else if err != apperrors.ErrNotFound {
		return err
	}

	pwHash, err := hash.HashPassword(password)
	if err != nil {
		return err
	}
	admin := &User{
		ID:            uuid.New(),
		Email:         email,
		PasswordHash:  pwHash,
		Role:          RoleAdmin,
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return s.repo.CreateUser(ctx, admin)
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPair, *User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, apperrors.ErrInvalidCredentials
	}
	if !hash.CheckPassword(user.PasswordHash, password) {
		return nil, nil, apperrors.ErrInvalidCredentials
	}
	pair, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	s.log.Info().Str("user_id", user.ID.String()).Str("event", "login").Msg("auth")
	return pair, user, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	userID, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}
	tokenHash := hashToken(refreshToken)
	if err := s.repo.FindRefreshToken(ctx, userID, tokenHash); err != nil {
		return nil, err
	}
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, user)
}

func (s *Service) Logout(ctx context.Context, userID uuid.UUID) error {
	s.log.Info().Str("user_id", userID.String()).Str("event", "logout").Msg("auth")
	return s.repo.DeleteRefreshTokens(ctx, userID)
}

func (s *Service) ResetPasswordRequest(ctx context.Context, email string) error {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		// не раскрываем наличие email
		return nil
	}
	token := newResetToken()
	if err := s.repo.SavePasswordReset(ctx, user.ID, hashToken(token), time.Now().Add(time.Hour)); err != nil {
		return err
	}
	if err := s.mailer.SendResetPassword(user.Email, token); err != nil {
		s.log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("password reset email failed")
		return err
	}
	return nil
}

func (s *Service) ResetPasswordConfirm(ctx context.Context, token, newPassword string) error {
	if err := ValidatePassword(newPassword); err != nil {
		return fmt.Errorf("%w: слабый пароль", apperrors.ErrValidation)
	}
	userID, err := s.repo.ConsumePasswordReset(ctx, hashToken(token))
	if err != nil {
		return err
	}
	pwHash, err := hash.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePassword(ctx, userID, pwHash); err != nil {
		return err
	}
	return s.repo.DeleteRefreshTokens(ctx, userID)
}

func (s *Service) issueTokens(ctx context.Context, user *User) (*TokenPair, error) {
	pair, err := s.jwt.GeneratePair(user)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveRefreshToken(ctx, user.ID, hashToken(pair.RefreshToken), time.Now().Add(s.jwt.RefreshTTL())); err != nil {
		return nil, err
	}
	return pair, nil
}
