package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/auth"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/response"
)

func Auth(jwt *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwt.ParseAccess(token)
		if err != nil {
			response.Fail(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func OptionalUserID(c *gin.Context) uuid.UUID {
	v, ok := c.Get("userID")
	if !ok {
		return uuid.Nil
	}
	return v.(uuid.UUID)
}
