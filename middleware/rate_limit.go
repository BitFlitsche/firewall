// Package middleware middleware/rate_limit.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

var limiter = rate.NewLimiter(1, 5) // 1 request per second with a burst of 5

// RateLimitMiddleware limits the rate of requests
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
