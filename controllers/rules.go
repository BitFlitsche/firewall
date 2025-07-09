// controllers/rules.go
package controllers

import (
	"firewall/models"
	"firewall/services"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"gorm.io/gorm"
)

var appStartTime = time.Now()
var requestCount int64
var errorCount int64

// Middleware zum Zählen von Requests/Errors
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestCount++
		c.Next()
		if c.Writer.Status() >= 400 {
			errorCount++
		}
	}
}

// SystemStats liefert System- und App-Metriken
func SystemStatsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uptime := time.Since(appStartTime).Seconds()
		cpuPercent, _ := cpu.Percent(0, false)
		memStats, _ := mem.VirtualMemory()
		diskStats, _ := disk.Usage("/")
		pid := os.Getpid()
		var dbConns int64
		db.Raw("SHOW STATUS WHERE variable_name = 'Threads_connected'").Scan(&dbConns)

		// DB Health
		dbHealth := "ok"
		if err := db.Exec("SELECT 1").Error; err != nil {
			dbHealth = "error"
		}
		// ES Health (Dummy, implementiere falls ES-Client verfügbar)
		esHealth := "unknown"

		c.JSON(200, gin.H{
			"uptime":         uptime,
			"cpu_percent":    cpuPercent,
			"memory_used":    memStats.Used,
			"memory_total":   memStats.Total,
			"memory_percent": memStats.UsedPercent,
			"disk_used":      diskStats.Used,
			"disk_total":     diskStats.Total,
			"disk_percent":   diskStats.UsedPercent,
			"db_health":      dbHealth,
			"db_connections": dbConns,
			"es_health":      esHealth,
			"request_count":  requestCount,
			"error_count":    errorCount,
			"go_routines":    runtime.NumGoroutine(),
			"pid":            pid,
		})
	}
}

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

// Count-Stats für IPs
func GetIPStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var total, allowed, denied, whitelisted int64
		db.Model(&models.IP{}).Count(&total)
		db.Model(&models.IP{}).Where("status = ?", "allowed").Count(&allowed)
		db.Model(&models.IP{}).Where("status = ?", "denied").Count(&denied)
		db.Model(&models.IP{}).Where("status = ?", "whitelisted").Count(&whitelisted)
		c.JSON(http.StatusOK, gin.H{
			"total":       total,
			"allowed":     allowed,
			"denied":      denied,
			"whitelisted": whitelisted,
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

		baseQuery := db.Model(&models.Email{})
		if status != "" {
			baseQuery = baseQuery.Where("status = ?", status)
		}
		if search != "" {
			baseQuery = baseQuery.Where("address LIKE ?", "%"+search+"%")
		}

		// Count-Query (ohne Limit/Offset/Order)
		countQuery := baseQuery.Session(&gorm.Session{})
		countQuery.Count(&total)

		// Items-Query (mit Limit/Offset/Order)
		itemQuery := baseQuery.Session(&gorm.Session{})
		if orderBy != "ID" && orderBy != "Address" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}
		itemQuery = itemQuery.Order(orderBy + " " + order)
		itemQuery = itemQuery.Offset((pageNum - 1) * limitNum).Limit(limitNum)
		itemQuery.Find(&emails)

		c.JSON(http.StatusOK, gin.H{
			"items": emails,
			"total": total,
		})
	}
}

// Count-Stats für Emails
func GetEmailStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var total, allowed, denied, whitelisted int64
		db.Model(&models.Email{}).Count(&total)
		db.Model(&models.Email{}).Where("status = ?", "allowed").Count(&allowed)
		db.Model(&models.Email{}).Where("status = ?", "denied").Count(&denied)
		db.Model(&models.Email{}).Where("status = ?", "whitelisted").Count(&whitelisted)
		c.JSON(http.StatusOK, gin.H{
			"total":       total,
			"allowed":     allowed,
			"denied":      denied,
			"whitelisted": whitelisted,
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

		baseQuery := db.Model(&models.UserAgent{})
		if status != "" {
			baseQuery = baseQuery.Where("status = ?", status)
		}
		if search != "" {
			baseQuery = baseQuery.Where("user_agent LIKE ?", "%"+search+"%")
		}

		// Count-Query (ohne Limit/Offset/Order)
		countQuery := baseQuery.Session(&gorm.Session{})
		countQuery.Count(&total)

		// Items-Query (mit Limit/Offset/Order)
		itemQuery := baseQuery.Session(&gorm.Session{})
		if orderBy != "ID" && orderBy != "UserAgent" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}
		itemQuery = itemQuery.Order(orderBy + " " + order)
		itemQuery = itemQuery.Offset((pageNum - 1) * limitNum).Limit(limitNum)
		itemQuery.Find(&userAgents)

		c.JSON(http.StatusOK, gin.H{
			"items": userAgents,
			"total": total,
		})
	}
}

// Count-Stats für UserAgents
func GetUserAgentStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var total, allowed, denied, whitelisted int64
		db.Model(&models.UserAgent{}).Count(&total)
		db.Model(&models.UserAgent{}).Where("status = ?", "allowed").Count(&allowed)
		db.Model(&models.UserAgent{}).Where("status = ?", "denied").Count(&denied)
		db.Model(&models.UserAgent{}).Where("status = ?", "whitelisted").Count(&whitelisted)
		c.JSON(http.StatusOK, gin.H{
			"total":       total,
			"allowed":     allowed,
			"denied":      denied,
			"whitelisted": whitelisted,
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

		baseQuery := db.Model(&models.Country{})
		if status != "" {
			baseQuery = baseQuery.Where("status = ?", status)
		}
		if search != "" {
			baseQuery = baseQuery.Where("code LIKE ?", "%"+search+"%")
		}

		// Count-Query (ohne Limit/Offset/Order)
		countQuery := baseQuery.Session(&gorm.Session{})
		countQuery.Count(&total)

		// Items-Query (mit Limit/Offset/Order)
		itemQuery := baseQuery.Session(&gorm.Session{})
		if orderBy != "ID" && orderBy != "Code" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}
		itemQuery = itemQuery.Order(orderBy + " " + order)
		itemQuery = itemQuery.Offset((pageNum - 1) * limitNum).Limit(limitNum)
		itemQuery.Find(&countries)

		c.JSON(http.StatusOK, gin.H{
			"items": countries,
			"total": total,
		})
	}
}

// Count-Stats für Countries
func GetCountryStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var total, allowed, denied, whitelisted int64
		db.Model(&models.Country{}).Count(&total)
		db.Model(&models.Country{}).Where("status = ?", "allowed").Count(&allowed)
		db.Model(&models.Country{}).Where("status = ?", "denied").Count(&denied)
		db.Model(&models.Country{}).Where("status = ?", "whitelisted").Count(&whitelisted)
		c.JSON(http.StatusOK, gin.H{
			"total":       total,
			"allowed":     allowed,
			"denied":      denied,
			"whitelisted": whitelisted,
		})
	}
}

// CreateCharsetRule fügt eine neue Charset-Regel hinzu
func CreateCharsetRule(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rule models.CharsetRule
		if err := c.ShouldBindJSON(&rule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&rule).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save charset rule"})
			return
		}
		_ = services.SyncCharsetToES(rule)
		services.PublishEvent("charset", "created", rule)
		c.JSON(http.StatusOK, rule)
	}
}

// GetCharsetRules listet alle Charset-Regeln
func GetCharsetRules(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rules []models.CharsetRule
		db.Find(&rules)
		c.JSON(http.StatusOK, rules)
	}
}

// UpdateCharsetRule aktualisiert eine Charset-Regel
func UpdateCharsetRule(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rule models.CharsetRule
		id := c.Param("id")
		if err := db.First(&rule, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
			return
		}
		var input models.CharsetRule
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rule.Charset = input.Charset
		rule.Status = input.Status
		if err := db.Save(&rule).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update charset rule"})
			return
		}
		_ = services.SyncCharsetToES(rule)
		services.PublishEvent("charset", "updated", rule)
		c.JSON(http.StatusOK, rule)
	}
}

// DeleteCharsetRule löscht eine Charset-Regel
func DeleteCharsetRule(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.CharsetRule{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete charset rule"})
			return
		}
		// Versuche auch aus ES zu löschen
		_ = services.DeleteCharsetFromES(parseUint(id))
		services.PublishEvent("charset", "deleted", models.CharsetRule{ID: parseUint(id)})
		c.JSON(http.StatusOK, gin.H{"message": "Charset rule deleted"})
	}
}

// Hilfsfunktion für DeleteCharsetRule
func parseUint(s string) uint {
	u, _ := strconv.ParseUint(s, 10, 64)
	return uint(u)
}

// Endpoint: POST /sync/charsets
func SyncCharsetsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := services.SyncAllCharsetsToES(db); err != nil {
			c.JSON(500, gin.H{"error": "Failed to sync charsets to Elasticsearch"})
			return
		}
		c.JSON(200, gin.H{"message": "All charsets synced to Elasticsearch"})
	}
}

// CreateUsernameRule fügt eine neue Username-Regel hinzu
func CreateUsernameRule(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rule models.UsernameRule
		if err := c.ShouldBindJSON(&rule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&rule).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save username rule"})
			return
		}
		_ = services.SyncUsernameToES(rule)
		services.PublishEvent("username", "created", rule)
		c.JSON(http.StatusOK, rule)
	}
}

// GetUsernameRules listet alle Username-Regeln
func GetUsernameRules(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rules []models.UsernameRule
		db.Find(&rules)
		c.JSON(http.StatusOK, rules)
	}
}

// UpdateUsernameRule aktualisiert eine Username-Regel
func UpdateUsernameRule(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rule models.UsernameRule
		id := c.Param("id")
		if err := db.First(&rule, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
			return
		}
		var input models.UsernameRule
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rule.Username = input.Username
		rule.Status = input.Status
		if err := db.Save(&rule).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username rule"})
			return
		}
		_ = services.SyncUsernameToES(rule)
		services.PublishEvent("username", "updated", rule)
		c.JSON(http.StatusOK, rule)
	}
}

// DeleteUsernameRule löscht eine Username-Regel
func DeleteUsernameRule(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.UsernameRule{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete username rule"})
			return
		}
		_ = services.DeleteUsernameFromES(parseUint(id))
		services.PublishEvent("username", "deleted", models.UsernameRule{ID: parseUint(id)})
		c.JSON(http.StatusOK, gin.H{"message": "Username rule deleted"})
	}
}
