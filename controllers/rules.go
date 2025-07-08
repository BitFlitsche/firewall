// controllers/rules.go
package controllers

import (
	"firewall/models"
	"firewall/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateIPAddress fügt eine neue IP-Adresse hinzu
// @Summary      Neue IP-Adresse anlegen
// @Description  Legt eine neue IP-Adresse mit Status an
// @Tags         ip
// @Accept       json
// @Produce      json
// @Param        ip  body      models.IP  true  "IP-Daten"
// @Success      200 {object}  models.IP
// @Failure      400 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /ip [post]
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

// GetIPAddresses listet alle IP-Adressen
// @Summary      IP-Adressen auflisten
// @Description  Gibt alle gespeicherten IP-Adressen zurück
// @Tags         ip
// @Produce      json
// @Success      200 {array}   models.IP
// @Router       /ip [get]
func GetIPAddresses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ips []models.IP
		db.Find(&ips)
		c.JSON(http.StatusOK, ips)
	}
}

// CreateEmail fügt eine neue E-Mail hinzu
// @Summary      Neue E-Mail anlegen
// @Description  Legt eine neue E-Mail-Adresse mit Status an
// @Tags         email
// @Accept       json
// @Produce      json
// @Param        email  body      models.Email  true  "E-Mail-Daten"
// @Success      200 {object}  models.Email
// @Failure      400 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /email [post]
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

// GetEmails listet alle E-Mails
// @Summary      E-Mails auflisten
// @Description  Gibt alle gespeicherten E-Mail-Adressen zurück
// @Tags         email
// @Produce      json
// @Success      200 {array}   models.Email
// @Router       /email [get]
func GetEmails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var emails []models.Email
		db.Find(&emails)
		c.JSON(http.StatusOK, emails)
	}
}

// CreateUserAgent fügt einen neuen User-Agent hinzu
// @Summary      Neuen User-Agent anlegen
// @Description  Legt einen neuen User-Agent mit Status an
// @Tags         useragent
// @Accept       json
// @Produce      json
// @Param        useragent  body      models.UserAgent  true  "User-Agent-Daten"
// @Success      200 {object}  models.UserAgent
// @Failure      400 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /useragent [post]
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

// GetUserAgents listet alle User-Agents
// @Summary      User-Agents auflisten
// @Description  Gibt alle gespeicherten User-Agents zurück
// @Tags         useragent
// @Produce      json
// @Success      200 {array}   models.UserAgent
// @Router       /useragent [get]
func GetUserAgents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userAgents []models.UserAgent
		db.Find(&userAgents)
		c.JSON(http.StatusOK, userAgents)
	}
}

// CreateCountry fügt ein neues Land hinzu
// @Summary      Neues Land anlegen
// @Description  Legt einen neuen Ländercode mit Status an
// @Tags         country
// @Accept       json
// @Produce      json
// @Param        country  body      models.Country  true  "Länder-Daten"
// @Success      200 {object}  models.Country
// @Failure      400 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /country [post]
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

// GetCountries listet alle Länder
// @Summary      Länder auflisten
// @Description  Gibt alle gespeicherten Länder zurück
// @Tags         country
// @Produce      json
// @Success      200 {array}   models.Country
// @Router       /country [get]
func GetCountries(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var countries []models.Country
		db.Find(&countries)
		c.JSON(http.StatusOK, countries)
	}
}
