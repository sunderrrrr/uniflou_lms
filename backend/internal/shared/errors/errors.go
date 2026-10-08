package apperrors

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrConflict       = errors.New("conflict")
	ErrValidation     = errors.New("validation")
	ErrRateLimit      = errors.New("rate limit exceeded")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
