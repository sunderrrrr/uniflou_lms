package middleware

import (
	"github.com/gin-gonic/gin"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/response"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role := c.GetString("role")
		if _, ok := allowed[role]; !ok {
			response.Fail(c, apperrors.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
