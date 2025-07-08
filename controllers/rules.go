// controllers/rules.go
package controllers

import (
	"firewall/models"
	"firewall/services"
	"fmt"
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

// GetIPAddresses listet alle IP-Adressen mit Paginierung, Filterung und Sortierung
// @Summary      IP-Adressen auflisten
// @Description  Gibt paginierte, gefilterte und sortierte IP-Adressen zurück
// @Tags         ip
// @Produce      json
// @Param        page     query     int     false  "Seite (beginnend bei 1)"
// @Param        limit    query     int     false  "Einträge pro Seite"
// @Param        status   query     string  false  "Status-Filter (allowed, denied, whitelisted)"
// @Param        search   query     string  false  "Suche nach IP-Adresse"
// @Param        orderBy  query     string  false  "Sortierfeld (ID, Address, Status)"
// @Param        order    query     string  false  "asc oder desc"
// @Success      200 {object} map[string]interface{}
// @Router       /ips [get]
func GetIPAddresses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ips []models.IP
		var total int64

		// Query-Parameter
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")
		status := c.Query("status")
		search := c.Query("search")
		orderBy := c.DefaultQuery("orderBy", "ID")
		order := c.DefaultQuery("order", "desc")

		// Umwandlung
		pageNum := 1
		limitNum := 10
		fmt.Sscanf(page, "%d", &pageNum)
		fmt.Sscanf(limit, "%d", &limitNum)
		if pageNum < 1 {
			pageNum = 1
		}
		if limitNum < 1 {
			limitNum = 10
		}

		dbQuery := db.Model(&models.IP{})
		if status != "" {
			dbQuery = dbQuery.Where("status = ?", status)
		}
		if search != "" {
			dbQuery = dbQuery.Where("address LIKE ?", "%"+search+"%")
		}

		dbQuery.Count(&total)

		if orderBy != "ID" && orderBy != "Address" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		dbQuery = dbQuery.Order(orderBy + " " + order)
		dbQuery = dbQuery.Offset((pageNum - 1) * limitNum).Limit(limitNum)
		dbQuery.Find(&ips)

		c.JSON(http.StatusOK, gin.H{
			"items": ips,
			"total": total,
		})
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

// GetEmails listet alle E-Mails mit Paginierung, Filterung und Sortierung
// @Summary      E-Mails auflisten
// @Description  Gibt paginierte, gefilterte und sortierte E-Mails zurück
// @Tags         email
// @Produce      json
// @Param        page     query     int     false  "Seite (beginnend bei 1)"
// @Param        limit    query     int     false  "Einträge pro Seite"
// @Param        status   query     string  false  "Status-Filter (allowed, denied, whitelisted)"
// @Param        search   query     string  false  "Suche nach E-Mail-Adresse"
// @Param        orderBy  query     string  false  "Sortierfeld (ID, Address, Status)"
// @Param        order    query     string  false  "asc oder desc"
// @Success      200 {object} map[string]interface{}
// @Router       /emails [get]
func GetEmails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var emails []models.Email
		var total int64

		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")
		status := c.Query("status")
		search := c.Query("search")
		orderBy := c.DefaultQuery("orderBy", "ID")
		order := c.DefaultQuery("order", "desc")

		pageNum := 1
		limitNum := 10
		fmt.Sscanf(page, "%d", &pageNum)
		fmt.Sscanf(limit, "%d", &limitNum)
		if pageNum < 1 {
			pageNum = 1
		}
		if limitNum < 1 {
			limitNum = 10
		}

		dbQuery := db.Model(&models.Email{})
		if status != "" {
			dbQuery = dbQuery.Where("status = ?", status)
		}
		if search != "" {
			dbQuery = dbQuery.Where("address LIKE ?", "%"+search+"%")
		}

		dbQuery.Count(&total)

		if orderBy != "ID" && orderBy != "Address" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		dbQuery = dbQuery.Order(orderBy + " " + order)
		dbQuery = dbQuery.Offset((pageNum - 1) * limitNum).Limit(limitNum)
		dbQuery.Find(&emails)

		c.JSON(http.StatusOK, gin.H{
			"items": emails,
			"total": total,
		})
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

// GetUserAgents listet alle User-Agents mit Paginierung, Filterung und Sortierung
// @Summary      User-Agents auflisten
// @Description  Gibt paginierte, gefilterte und sortierte User-Agents zurück
// @Tags         useragent
// @Produce      json
// @Param        page     query     int     false  "Seite (beginnend bei 1)"
// @Param        limit    query     int     false  "Einträge pro Seite"
// @Param        status   query     string  false  "Status-Filter (allowed, denied, whitelisted)"
// @Param        search   query     string  false  "Suche nach User-Agent"
// @Param        orderBy  query     string  false  "Sortierfeld (ID, UserAgent, Status)"
// @Param        order    query     string  false  "asc oder desc"
// @Success      200 {object} map[string]interface{}
// @Router       /user-agents [get]
func GetUserAgents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userAgents []models.UserAgent
		var total int64

		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")
		status := c.Query("status")
		search := c.Query("search")
		orderBy := c.DefaultQuery("orderBy", "ID")
		order := c.DefaultQuery("order", "desc")

		pageNum := 1
		limitNum := 10
		fmt.Sscanf(page, "%d", &pageNum)
		fmt.Sscanf(limit, "%d", &limitNum)
		if pageNum < 1 {
			pageNum = 1
		}
		if limitNum < 1 {
			limitNum = 10
		}

		dbQuery := db.Model(&models.UserAgent{})
		if status != "" {
			dbQuery = dbQuery.Where("status = ?", status)
		}
		if search != "" {
			dbQuery = dbQuery.Where("user_agent LIKE ?", "%"+search+"%")
		}

		dbQuery.Count(&total)

		if orderBy != "ID" && orderBy != "UserAgent" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		dbQuery = dbQuery.Order(orderBy + " " + order)
		dbQuery = dbQuery.Offset((pageNum - 1) * limitNum).Limit(limitNum)
		dbQuery.Find(&userAgents)

		c.JSON(http.StatusOK, gin.H{
			"items": userAgents,
			"total": total,
		})
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

// GetCountries listet alle Länder mit Paginierung, Filterung und Sortierung
// @Summary      Länder auflisten
// @Description  Gibt paginierte, gefilterte und sortierte Länder zurück
// @Tags         country
// @Produce      json
// @Param        page     query     int     false  "Seite (beginnend bei 1)"
// @Param        limit    query     int     false  "Einträge pro Seite"
// @Param        status   query     string  false  "Status-Filter (allowed, denied, whitelisted)"
// @Param        search   query     string  false  "Suche nach Country Code"
// @Param        orderBy  query     string  false  "Sortierfeld (ID, Code, Status)"
// @Param        order    query     string  false  "asc oder desc"
// @Success      200 {object} map[string]interface{}
// @Router       /countries [get]
func GetCountries(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var countries []models.Country
		var total int64

		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")
		status := c.Query("status")
		search := c.Query("search")
		orderBy := c.DefaultQuery("orderBy", "ID")
		order := c.DefaultQuery("order", "desc")

		pageNum := 1
		limitNum := 10
		fmt.Sscanf(page, "%d", &pageNum)
		fmt.Sscanf(limit, "%d", &limitNum)
		if pageNum < 1 {
			pageNum = 1
		}
		if limitNum < 1 {
			limitNum = 10
		}

		dbQuery := db.Model(&models.Country{})
		if status != "" {
			dbQuery = dbQuery.Where("status = ?", status)
		}
		if search != "" {
			dbQuery = dbQuery.Where("code LIKE ?", "%"+search+"%")
		}

		dbQuery.Count(&total)

		if orderBy != "ID" && orderBy != "Code" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		dbQuery = dbQuery.Order(orderBy + " " + order)
		dbQuery = dbQuery.Offset((pageNum - 1) * limitNum).Limit(limitNum)
		dbQuery.Find(&countries)

		c.JSON(http.StatusOK, gin.H{
			"items": countries,
			"total": total,
		})
	}
}
