package auth

import (
	"unicode"

	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return apperrors.ErrValidation
	}
	var hasLetter, hasDigit bool
	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return apperrors.ErrValidation
	}
	return nil
}

func ValidateRole(role string) error {
	normalized := NormalizeRole(role)
	if normalized == "" {
		return apperrors.ErrValidation
	}
	return nil
}
