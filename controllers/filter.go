package controllers

import (
	"context"
	"errors"
	"firewall/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// FilterRequest defines the structure for the incoming JSON request
type FilterRequest struct {
	IP        string `json:"ip"`
	Email     string `json:"email"`
	UserAgent string `json:"user_agent"`
	Country   string `json:"country"`
}

// FilterRequestHandler handles filtering of IP, email, user agents, and countries
func FilterRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input FilterRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
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
