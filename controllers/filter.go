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
						c.JSON(200, gin.H{"result": "denied", "reason": "charset denied", "field": field, "charset": cs})
						return
					}
					if rule.Status == "whitelisted" {
						c.JSON(200, gin.H{"result": "whitelisted", "reason": "charset whitelisted", "field": field, "charset": cs})
						return
					}
				}
			}
		}

		// Username-Filter: Prüfe gegen UsernameRule-Liste
		if input.Username != "" {
			var usernameRule models.UsernameRule
			db.Where("username = ?", input.Username).First(&usernameRule)
			if usernameRule.ID != 0 {
				if usernameRule.Status == "denied" {
					c.JSON(200, gin.H{"result": "denied", "reason": "username denied", "field": "username", "username": input.Username})
					return
				}
				if usernameRule.Status == "whitelisted" {
					c.JSON(200, gin.H{"result": "whitelisted", "reason": "username whitelisted", "field": "username", "username": input.Username})
					return
				}
			}
		}

		// Timeout for the entire operation (e.g., 5 seconds)
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// Call the service to evaluate filters
		finalResult, err := services.EvaluateFilters(ctx, input.IP, input.Email, input.UserAgent, input.Country)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timed out"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, finalResult)
	}
}
