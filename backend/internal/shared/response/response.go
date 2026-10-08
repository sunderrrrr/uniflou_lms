package response

import (
	"errors"
	"net/http"

	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

func Fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{Error: ErrorBody{Code: "NOT_FOUND", Message: "Ресурс не найден"}})
	case errors.Is(err, apperrors.ErrUnauthorized), errors.Is(err, apperrors.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: ErrorBody{Code: "UNAUTHORIZED", Message: "Неверные учётные данные"}})
	case errors.Is(err, apperrors.ErrForbidden):
		c.JSON(http.StatusForbidden, ErrorResponse{Error: ErrorBody{Code: "FORBIDDEN", Message: "Доступ запрещён"}})
	case errors.Is(err, apperrors.ErrConflict):
		c.JSON(http.StatusConflict, ErrorResponse{Error: ErrorBody{Code: "CONFLICT", Message: errMsg(err, "Email уже зарегистрирован")}})
	case errors.Is(err, apperrors.ErrValidation):
		msg := err.Error()
		if msg == apperrors.ErrValidation.Error() {
			msg = "Ошибка валидации"
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrorBody{Code: "VALIDATION", Message: msg}})
	case errors.Is(err, apperrors.ErrRateLimit):
		c.JSON(http.StatusTooManyRequests, ErrorResponse{Error: ErrorBody{Code: "RATE_LIMIT", Message: "Превышен лимит запросов"}})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: ErrorBody{Code: "INTERNAL", Message: "Внутренняя ошибка сервера"}})
	}
}

func errMsg(err error, fallback string) string {
	if err != nil && err.Error() != "" && err.Error() != fallback {
		return err.Error()
	}
	return fallback
}

type Page struct {
	Items any `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
}
