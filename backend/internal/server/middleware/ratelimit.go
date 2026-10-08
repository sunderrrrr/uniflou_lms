package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/response"
)

type ipLimiter struct {
	mu    sync.Mutex
	hits  map[string][]time.Time
	limit int
	window time.Duration
}

func newIPLimiter(limit int, window time.Duration) *ipLimiter {
	return &ipLimiter{hits: map[string][]time.Time{}, limit: limit, window: window}
}

func (l *ipLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-l.window)
	var kept []time.Time
	for _, t := range l.hits[key] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= l.limit {
		l.hits[key] = kept
		return false
	}
	l.hits[key] = append(kept, now)
	return true
}

func RateLimitLogin(limit int) gin.HandlerFunc {
	l := newIPLimiter(limit, 15*time.Minute)
	return func(c *gin.Context) {
		if !l.allow(c.ClientIP()) {
			response.Fail(c, apperrors.ErrRateLimit)
			c.Abort()
			return
		}
		c.Next()
	}
}

func RateLimitRegister(limit int) gin.HandlerFunc {
	l := newIPLimiter(limit, time.Hour)
	return func(c *gin.Context) {
		if !l.allow(c.ClientIP()) {
			response.Fail(c, apperrors.ErrRateLimit)
			c.Abort()
			return
		}
		c.Next()
	}
}
