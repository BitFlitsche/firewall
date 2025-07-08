// controllers/rules.go
package controllers

import (
	"firewall/models"
	"firewall/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateIPAddress adds a new IP address to the database
func CreateIPAddress(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ip models.IP

		if err := c.ShouldBindJSON(&ip); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save to MySQL first
		if err := db.Create(&ip).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save IP address"})
			return
		}

		// Publish event for async processing
		services.PublishEvent("ip", "created", ip)

		c.JSON(http.StatusOK, ip)
	}
}

// GetIPAddresses lists all IP addresses
func GetIPAddresses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ips []models.IP
		db.Find(&ips)
		c.JSON(http.StatusOK, ips)
	}
}

// CreateEmail adds a new email address to the database
func CreateEmail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var email models.Email
		if err := c.ShouldBindJSON(&email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save to MySQL first
		if err := db.Create(&email).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save email"})
			return
		}

		// Publish event for async processing
		services.PublishEvent("email", "created", email)

		c.JSON(http.StatusOK, email)
	}
}

// GetEmails lists all emails
func GetEmails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var emails []models.Email
		db.Find(&emails)
		c.JSON(http.StatusOK, emails)
	}
}

// CreateUserAgent adds a new user agent to the database
func CreateUserAgent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userAgent models.UserAgent
		if err := c.ShouldBindJSON(&userAgent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save to MySQL first
		if err := db.Create(&userAgent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user agent"})
			return
		}

		// Publish event for async processing
		services.PublishEvent("user_agent", "created", userAgent)

		c.JSON(http.StatusOK, userAgent)
	}
}

// GetUserAgents lists all user agents
func GetUserAgents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userAgents []models.UserAgent
		db.Find(&userAgents)
		c.JSON(http.StatusOK, userAgents)
	}
}

// CreateCountry adds a new country code to the database
func CreateCountry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var country models.Country
		if err := c.ShouldBindJSON(&country); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save to MySQL first
		if err := db.Create(&country).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save country"})
			return
		}

		// Publish event for async processing
		services.PublishEvent("country", "created", country)

		c.JSON(http.StatusOK, country)
	}
}

// GetCountries lists all country codes
func GetCountries(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var countries []models.Country
		db.Find(&countries)
		c.JSON(http.StatusOK, countries)
	}
}
