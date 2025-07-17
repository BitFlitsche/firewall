package controllers

import (
	"context"
	"errors"
	"firewall/services"
	"net/http"
	"time"

	"firewall/models"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FilterRequest defines the structure for the incoming JSON request
type FilterRequest struct {
	IP        string `json:"ip"`
	Email     string `json:"email"`
	UserAgent string `json:"user_agent"`
	Country   string `json:"country"`
	Content   string `json:"content"`  // optional
	Username  string `json:"username"` // optional
}

// Helper: Charset-Erkennung (sehr einfach, z.B. ASCII, Latin, Cyrillic, etc.)
func detectCharset(s string) string {
	if s == "" {
		return "ASCII"
	}
	ascii := true
	latin := true
	cyrillic := true
	for _, r := range s {
		if r > 127 {
			ascii = false
		}
		if !(r >= 0x0020 && r <= 0x007E) && !(r >= 0x00A0 && r <= 0x00FF) {
			latin = false
		}
		if !(r >= 0x0400 && r <= 0x04FF) {
			cyrillic = false
		}
	}
	if ascii {
		return "ASCII"
	}
	if latin {
		return "Latin"
	}
	if cyrillic {
		return "Cyrillic"
	}
	if utf8.ValidString(s) {
		return "UTF-8"
	}
	return "Other"
}

// FilterRequestHandler prüft jetzt auch Charset-Regeln
func FilterRequestHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input FilterRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Generate a cache key based on the filter input
		cache := services.GetCache()
		cacheKey := "filter:" + input.IP + ":" + input.Email + ":" + input.UserAgent + ":" + input.Country + ":" + input.Username

		// Try to get from cache first
		if cached, exists := cache.Get(cacheKey); exists {
			c.JSON(http.StatusOK, cached)
			return
		}

		// Lade alle Charset-Regeln
		var charsetRules []models.CharsetRule
		db.Find(&charsetRules)

		// Prüfe Email, UserAgent, Content, Username auf Charset-Regeln
		fields := map[string]string{
			"email":      input.Email,
			"user_agent": input.UserAgent,
			"content":    input.Content,
			"username":   input.Username,
		}
		for field, value := range fields {
			if value == "" {
				continue
			}
			cs := detectCharset(value)
			for _, rule := range charsetRules {
				if rule.Charset == cs {
					if rule.Status == "denied" {
						cache.Set(cacheKey, gin.H{"result": "denied", "reason": "charset denied", "field": field, field: value}, 5*time.Minute)
						c.JSON(200, gin.H{"result": "denied", "reason": "charset denied", "field": field, field: value})
						return
					}
					if rule.Status == "whitelisted" {
						cache.Set(cacheKey, gin.H{"result": "whitelisted", "reason": "charset whitelisted", "field": field, field: value}, 5*time.Minute)
						c.JSON(200, gin.H{"result": "whitelisted", "reason": "charset whitelisted", "field": field, field: value})
						return
					}
				}
			}
		}

		// Username-Filter: Prüfe gegen UsernameRule-Liste
		// Note: Username filtering is now handled by Elasticsearch regex filtering
		// This allows for both exact matches and regex patterns

		// Timeout for the entire operation (e.g., 5 seconds)
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// Call the service to evaluate filters
		finalResult, err := services.EvaluateFilters(ctx, input.IP, input.Email, input.UserAgent, input.Country, input.Username)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timed out"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			return
		}

		// Cache the result for 5 minutes
		cache.Set(cacheKey, finalResult, 5*time.Minute)

		c.JSON(http.StatusOK, finalResult)
	}
}
