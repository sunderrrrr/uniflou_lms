package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/auth"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

func (s *Service) Me(ctx context.Context, userID uuid.UUID) (*auth.User, *StudentProfile, error) {
	u, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	var profile *StudentProfile
	if u.Role == auth.RoleStudent {
		profile, _ = s.repo.GetOrCreateProfile(ctx, userID)
	}
	return u, profile, nil
}

func (s *Service) LinkStudent(ctx context.Context, parentID uuid.UUID, studentEmail string) error {
	studentID, err := s.repo.FindStudentByEmail(ctx, studentEmail)
	if err != nil {
		return err
	}
	return s.repo.LinkParentStudent(ctx, parentID, studentID)
}

func (s *Service) UpdateMe(ctx context.Context, userID uuid.UUID, grade, targetScore *int) (*StudentProfile, error) {
	updates := map[string]any{}
	if grade != nil {
		updates["grade"] = *grade
	}
	if targetScore != nil {
		updates["target_score"] = *targetScore
	}
	if len(updates) == 0 {
		return s.repo.GetOrCreateProfile(ctx, userID)
	}
	return s.repo.UpdateProfile(ctx, userID, updates)
}
