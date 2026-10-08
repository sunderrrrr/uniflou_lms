package auth

import (
	"testing"

	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
)

func TestValidatePassword(t *testing.T) {
	if err := ValidatePassword("short"); err != apperrors.ErrValidation {
		t.Fatalf("expected validation error for short password")
	}
	if err := ValidatePassword("validpass1"); err != nil {
		t.Fatalf("expected valid password, got %v", err)
	}
}
