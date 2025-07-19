// controllers/rules.go
package controllers

import (
	"firewall/models"
	"firewall/services"
	"firewall/validation"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"context"
	"encoding/json"
	"firewall/config"

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

		// Get database connection statistics using the new config function
		dbStats := config.GetDBStats()

		// DB Health
		dbHealth := "ok"
		if err := db.Exec("SELECT 1").Error; err != nil {
			dbHealth = "error"
		}

		// ES Health
		esHealth := "unknown"
		if config.ESClient != nil {
			res, err := config.ESClient.Cluster.Health(
				config.ESClient.Cluster.Health.WithContext(context.Background()),
			)
			if err == nil && res != nil {
				defer res.Body.Close()
				var health map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&health); err == nil {
					if status, ok := health["status"].(string); ok {
						esHealth = status // "green", "yellow", "red"
					} else {
						esHealth = "error"
					}
				} else {
					esHealth = "error"
				}
			} else {
				esHealth = "error"
			}
		}

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
			"db_connections": dbStats,
			"es_health":      esHealth,
			"request_count":  requestCount,
			"error_count":    errorCount,
			"go_routines":    runtime.NumGoroutine(),
			"pid":            pid,
			"cache_stats": func() map[string]interface{} {
				stats, err := services.GetCacheFactory().Stats()
				if err != nil {
					return map[string]interface{}{"error": err.Error()}
				}
				return stats
			}(),
			"cache_type": services.GetCacheFactory().GetCacheType(),
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		// Comprehensive validation
		ipValidation := validation.ValidateIP(ip.Address)
		statusValidation := validation.ValidateStatus(ip.Status)

		if !ipValidation.IsValid || !statusValidation.IsValid {
			errors := []validation.ValidationError{}
			errors = append(errors, ipValidation.Errors...)
			errors = append(errors, statusValidation.Errors...)

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			return
		}

		// Check if IP already exists
		var existingIP models.IP
		if err := db.Where("address = ?", ip.Address).First(&existingIP).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "IP address already exists", "address": ip.Address})
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
		// Query-Parameter
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")
		status := c.Query("status")
		typeFilter := c.Query("type")
		search := c.Query("search")
		orderBy := c.DefaultQuery("orderBy", "ID")
		order := c.DefaultQuery("order", "desc")

		// Validate query parameters
		paginationValidation := validation.ValidatePagination(page, limit)
		if !paginationValidation.IsValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid pagination parameters",
				"details": paginationValidation.Errors,
			})
			return
		}

		// Validate status if provided
		if status != "" {
			statusValidation := validation.ValidateStatus(status)
			if !statusValidation.IsValid {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid status parameter",
					"details": statusValidation.Errors,
				})
				return
			}
		}

		// Validate search if provided
		if search != "" {
			searchValidation := validation.ValidateSearch(search)
			if !searchValidation.IsValid {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid search parameter",
					"details": searchValidation.Errors,
				})
				return
			}
		}

		// Validate orderBy and order
		validFields := []string{"ID", "Address", "Status"}
		orderValidation := validation.ValidateOrderBy(orderBy, order, validFields)
		if !orderValidation.IsValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ordering parameters",
				"details": orderValidation.Errors,
			})
			return
		}

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

		// Build WHERE conditions
		var conditions []string
		var args []interface{}

		if status != "" {
			conditions = append(conditions, "status = ?")
			args = append(args, status)
		}
		if typeFilter != "" {
			if typeFilter == "single" {
				conditions = append(conditions, "is_cidr = ?")
				args = append(args, false)
			} else if typeFilter == "cidr" {
				conditions = append(conditions, "is_cidr = ?")
				args = append(args, true)
			}
		}
		if search != "" {
			conditions = append(conditions, "address LIKE ?")
			args = append(args, "%"+search+"%")
		}

		// Build WHERE clause
		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		// Validate orderBy and order
		if orderBy != "ID" && orderBy != "Address" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM ips 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, orderBy, order)

		// Add pagination parameters
		args = append(args, limitNum, (pageNum-1)*limitNum)

		// Execute query
		var results []struct {
			models.IP
			TotalCount int64 `json:"total_count"`
		}

		if err := db.Raw(query, args...).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch IP addresses"})
			return
		}

		// Extract data
		var ips []models.IP
		var total int64
		if len(results) > 0 {
			total = results[0].TotalCount
			for _, result := range results {
				ips = append(ips, result.IP)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"items": ips,
			"total": total,
		})
	}
}

// Count-Stats für IPs
func GetIPStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var total, allowed, denied, whitelisted, single, cidr int64
		db.Model(&models.IP{}).Count(&total)
		db.Model(&models.IP{}).Where("status = ?", "allowed").Count(&allowed)
		db.Model(&models.IP{}).Where("status = ?", "denied").Count(&denied)
		db.Model(&models.IP{}).Where("status = ?", "whitelisted").Count(&whitelisted)

		// Use GORM model field names instead of database column names
		db.Model(&models.IP{}).Where("is_cidr = ?", false).Count(&single)
		db.Model(&models.IP{}).Where("is_cidr = ?", true).Count(&cidr)

		// Debug: Let's also check what the actual column name is
		var testIP models.IP
		if err := db.First(&testIP, 26).Error; err == nil {
			log.Printf("Debug: IP 26 IsCIDR = %v", testIP.IsCIDR)
		}

		c.JSON(http.StatusOK, gin.H{
			"total":       total,
			"allowed":     allowed,
			"denied":      denied,
			"whitelisted": whitelisted,
			"single":      single,
			"cidr":        cidr,
		})
	}
}

// UpdateIPAddress aktualisiert eine IP-Adresse
func UpdateIPAddress(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate ID parameter
		idValidation := validation.ValidateID(c.Param("id"))
		if !idValidation.IsValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID parameter",
				"details": idValidation.Errors,
			})
			return
		}

		var ip models.IP
		id := c.Param("id")
		if err := db.First(&ip, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "IP address not found"})
			return
		}

		var input models.IP
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		// Comprehensive validation
		ipValidation := validation.ValidateIP(input.Address)
		statusValidation := validation.ValidateStatus(input.Status)

		if !ipValidation.IsValid || !statusValidation.IsValid {
			errors := []validation.ValidationError{}
			errors = append(errors, ipValidation.Errors...)
			errors = append(errors, statusValidation.Errors...)

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			return
		}

		// Check if the new address conflicts with existing records (excluding current record)
		var existingIP models.IP
		if err := db.Where("address = ? AND id != ?", input.Address, id).First(&existingIP).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "IP address already exists", "address": input.Address})
			return
		}

		ip.Address = input.Address
		ip.Status = input.Status
		ip.IsCIDR = input.IsCIDR

		if err := db.Save(&ip).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update IP address"})
			return
		}

		services.PublishEvent("ip", "updated", ip)
		c.JSON(http.StatusOK, ip)
	}
}

// DeleteIPAddress löscht eine IP-Adresse
func DeleteIPAddress(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.IP{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IP address"})
			return
		}
		services.PublishEvent("ip", "deleted", models.IP{ID: parseUint(id)})
		c.JSON(http.StatusOK, gin.H{"message": "IP address deleted"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		// Comprehensive validation
		emailValidation := validation.ValidateEmail(email.Address)
		statusValidation := validation.ValidateStatus(email.Status)

		if !emailValidation.IsValid || !statusValidation.IsValid {
			errors := []validation.ValidationError{}
			errors = append(errors, emailValidation.Errors...)
			errors = append(errors, statusValidation.Errors...)

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			return
		}

		// Check if email already exists
		var existingEmail models.Email
		if err := db.Where("address = ?", email.Address).First(&existingEmail).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email address already exists", "address": email.Address})
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

		// Build WHERE conditions
		var conditions []string
		var args []interface{}

		if status != "" {
			conditions = append(conditions, "status = ?")
			args = append(args, status)
		}
		if search != "" {
			conditions = append(conditions, "address LIKE ?")
			args = append(args, "%"+search+"%")
		}

		// Build WHERE clause
		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		// Validate orderBy and order
		if orderBy != "ID" && orderBy != "Address" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM emails 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, orderBy, order)

		// Add pagination parameters
		args = append(args, limitNum, (pageNum-1)*limitNum)

		// Execute query
		var results []struct {
			models.Email
			TotalCount int64 `json:"total_count"`
		}

		if err := db.Raw(query, args...).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch emails"})
			return
		}

		// Extract data
		var emails []models.Email
		var total int64
		if len(results) > 0 {
			total = results[0].TotalCount
			for _, result := range results {
				emails = append(emails, result.Email)
			}
		}

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

// UpdateEmail aktualisiert eine E-Mail-Adresse
func UpdateEmail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var email models.Email
		id := c.Param("id")
		if err := db.First(&email, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			return
		}
		var input models.Email
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email.Address = input.Address
		email.Status = input.Status
		email.IsRegex = input.IsRegex
		if err := db.Save(&email).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email"})
			return
		}
		services.PublishEvent("email", "updated", email)
		c.JSON(http.StatusOK, email)
	}
}

// DeleteEmail löscht eine E-Mail-Adresse
func DeleteEmail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.Email{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete email"})
			return
		}
		services.PublishEvent("email", "deleted", models.Email{ID: parseUint(id)})
		c.JSON(http.StatusOK, gin.H{"message": "Email deleted"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		// Comprehensive validation
		userAgentValidation := validation.ValidateUserAgent(userAgent.UserAgent)
		statusValidation := validation.ValidateStatus(userAgent.Status)

		if !userAgentValidation.IsValid || !statusValidation.IsValid {
			errors := []validation.ValidationError{}
			errors = append(errors, userAgentValidation.Errors...)
			errors = append(errors, statusValidation.Errors...)

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			return
		}

		// Check if user agent already exists
		var existingUserAgent models.UserAgent
		if err := db.Where("user_agent = ?", userAgent.UserAgent).First(&existingUserAgent).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User agent already exists", "user_agent": userAgent.UserAgent})
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

		// Build WHERE conditions
		var conditions []string
		var args []interface{}

		if status != "" {
			conditions = append(conditions, "status = ?")
			args = append(args, status)
		}
		if search != "" {
			conditions = append(conditions, "user_agent LIKE ?")
			args = append(args, "%"+search+"%")
		}

		// Build WHERE clause
		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		// Validate orderBy and order
		if orderBy != "ID" && orderBy != "UserAgent" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM user_agents 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, orderBy, order)

		// Add pagination parameters
		args = append(args, limitNum, (pageNum-1)*limitNum)

		// Execute query
		var results []struct {
			models.UserAgent
			TotalCount int64 `json:"total_count"`
		}

		if err := db.Raw(query, args...).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user agents"})
			return
		}

		// Extract data
		var userAgents []models.UserAgent
		var total int64
		if len(results) > 0 {
			total = results[0].TotalCount
			for _, result := range results {
				userAgents = append(userAgents, result.UserAgent)
			}
		}

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

// UpdateUserAgent aktualisiert einen User-Agent
func UpdateUserAgent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userAgent models.UserAgent
		id := c.Param("id")
		if err := db.First(&userAgent, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User agent not found"})
			return
		}
		var input models.UserAgent
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userAgent.UserAgent = input.UserAgent
		userAgent.Status = input.Status
		userAgent.IsRegex = input.IsRegex
		if err := db.Save(&userAgent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user agent"})
			return
		}
		services.PublishEvent("user_agent", "updated", userAgent)
		c.JSON(http.StatusOK, userAgent)
	}
}

// DeleteUserAgent löscht einen User-Agent
func DeleteUserAgent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.UserAgent{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user agent"})
			return
		}
		services.PublishEvent("user_agent", "deleted", models.UserAgent{ID: parseUint(id)})
		c.JSON(http.StatusOK, gin.H{"message": "User agent deleted"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		// Comprehensive validation
		countryValidation := validation.ValidateCountry(country.Code)
		statusValidation := validation.ValidateStatus(country.Status)

		if !countryValidation.IsValid || !statusValidation.IsValid {
			errors := []validation.ValidationError{}
			errors = append(errors, countryValidation.Errors...)
			errors = append(errors, statusValidation.Errors...)

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			return
		}

		// Check if country already exists
		var existingCountry models.Country
		if err := db.Where("code = ?", country.Code).First(&existingCountry).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Country code already exists", "code": country.Code})
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

		// Build WHERE conditions
		var conditions []string
		var args []interface{}

		if status != "" {
			conditions = append(conditions, "status = ?")
			args = append(args, status)
		}
		if search != "" {
			conditions = append(conditions, "code LIKE ?")
			args = append(args, "%"+search+"%")
		}

		// Build WHERE clause
		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		// Validate orderBy and order
		if orderBy != "ID" && orderBy != "Code" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM countries 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, orderBy, order)

		// Add pagination parameters
		args = append(args, limitNum, (pageNum-1)*limitNum)

		// Execute query
		var results []struct {
			models.Country
			TotalCount int64 `json:"total_count"`
		}

		if err := db.Raw(query, args...).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch countries"})
			return
		}

		// Extract data
		var countries []models.Country
		var total int64
		if len(results) > 0 {
			total = results[0].TotalCount
			for _, result := range results {
				countries = append(countries, result.Country)
			}
		}

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

// UpdateCountry aktualisiert ein Land
func UpdateCountry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var country models.Country
		id := c.Param("id")
		if err := db.First(&country, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
			return
		}
		var input models.Country
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		country.Code = input.Code
		country.Status = input.Status
		if err := db.Save(&country).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update country"})
			return
		}
		services.PublishEvent("country", "updated", country)
		c.JSON(http.StatusOK, country)
	}
}

// DeleteCountry löscht ein Land
func DeleteCountry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.Country{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete country"})
			return
		}
		services.PublishEvent("country", "deleted", models.Country{ID: parseUint(id)})
		c.JSON(http.StatusOK, gin.H{"message": "Country deleted"})
	}
}

// CreateCharsetRule fügt eine neue Charset-Regel hinzu
// @Summary      Neue Charset-Regel anlegen
// @Description  Legt eine neue Charset-Regel mit Status an
// @Tags         charset
// @Accept       json
// @Produce      json
// @Param        charset  body      models.CharsetRule  true  "Charset-Daten"
// @Success      200 {object}  models.CharsetRule
// @Failure      400 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /charset [post]
func CreateCharsetRule(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var charset models.CharsetRule
		if err := c.ShouldBindJSON(&charset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		// Comprehensive validation
		charsetValidation := validation.ValidateCharset(charset.Charset)
		statusValidation := validation.ValidateStatus(charset.Status)

		if !charsetValidation.IsValid || !statusValidation.IsValid {
			errors := []validation.ValidationError{}
			errors = append(errors, charsetValidation.Errors...)
			errors = append(errors, statusValidation.Errors...)

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			return
		}

		// Check if charset already exists
		var existingCharset models.CharsetRule
		if err := db.Where("charset = ?", charset.Charset).First(&existingCharset).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Charset already exists", "charset": charset.Charset})
			return
		}

		// Save to MySQL first
		if err := db.Create(&charset).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save charset rule"})
			return
		}

		// Publish event for async processing
		services.PublishEvent("charset", "created", charset)

		c.JSON(http.StatusOK, charset)
	}
}

// GetCharsetRules listet alle Charset-Regeln mit Paginierung, Filterung und Sortierung
// @Summary      Charset-Regeln auflisten
// @Description  Gibt paginierte, gefilterte und sortierte Charset-Regeln zurück
// @Tags         charset
// @Produce      json
// @Param        page     query     int     false  "Seite (beginnend bei 1)"
// @Param        limit    query     int     false  "Einträge pro Seite"
// @Param        status   query     string  false  "Status-Filter (allowed, denied, whitelisted)"
// @Param        search   query     string  false  "Suche nach Charset"
// @Param        orderBy  query     string  false  "Sortierfeld (ID, Charset, Status)"
// @Param        order    query     string  false  "asc oder desc"
// @Success      200 {object} map[string]interface{}
// @Router       /charsets [get]
func GetCharsetRules(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Build WHERE conditions
		var conditions []string
		var args []interface{}

		if status != "" {
			conditions = append(conditions, "status = ?")
			args = append(args, status)
		}
		if search != "" {
			conditions = append(conditions, "charset LIKE ?")
			args = append(args, "%"+search+"%")
		}

		// Build WHERE clause
		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		// Validate orderBy and order
		if orderBy != "ID" && orderBy != "Charset" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM charset_rules 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, orderBy, order)

		// Add pagination parameters
		args = append(args, limitNum, (pageNum-1)*limitNum)

		// Execute query
		var results []struct {
			models.CharsetRule
			TotalCount int64 `json:"total_count"`
		}

		if err := db.Raw(query, args...).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch charset rules"})
			return
		}

		// Extract data
		var rules []models.CharsetRule
		var total int64
		if len(results) > 0 {
			total = results[0].TotalCount
			for _, result := range results {
				rules = append(rules, result.CharsetRule)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"items": rules,
			"total": total,
		})
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

// Endpoint: POST /sync/usernames
func SyncUsernamesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := services.SyncAllUsernamesToES(db); err != nil {
			c.JSON(500, gin.H{"error": "Failed to sync usernames to Elasticsearch"})
			return
		}
		c.JSON(200, gin.H{"message": "All usernames synced to Elasticsearch"})
	}
}

// CreateUsernameRule fügt eine neue Username-Regel hinzu
// @Summary      Neue Username-Regel anlegen
// @Description  Legt eine neue Username-Regel mit Status an
// @Tags         username
// @Accept       json
// @Produce      json
// @Param        username  body      models.UsernameRule  true  "Username-Daten"
// @Success      200 {object}  models.UsernameRule
// @Failure      400 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /username [post]
func CreateUsernameRule(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var username models.UsernameRule
		if err := c.ShouldBindJSON(&username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		// Comprehensive validation
		usernameValidation := validation.ValidateUsername(username.Username)
		statusValidation := validation.ValidateStatus(username.Status)

		if !usernameValidation.IsValid || !statusValidation.IsValid {
			errors := []validation.ValidationError{}
			errors = append(errors, usernameValidation.Errors...)
			errors = append(errors, statusValidation.Errors...)

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			return
		}

		// Validate regex if IsRegex is true
		if username.IsRegex {
			regexValidation := validation.ValidateRegex(username.Username)
			if !regexValidation.IsValid {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid regex pattern",
					"details": regexValidation.Errors,
				})
				return
			}
		}

		// Check if username already exists
		var existingUsername models.UsernameRule
		if err := db.Where("username = ?", username.Username).First(&existingUsername).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username rule already exists", "username": username.Username})
			return
		}

		// Save to MySQL first
		if err := db.Create(&username).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save username rule"})
			return
		}

		// Publish event for async processing
		services.PublishEvent("username", "created", username)

		c.JSON(http.StatusOK, username)
	}
}

// GetUsernameRules listet alle Username-Regeln mit Paginierung, Filterung und Sortierung
// @Summary      Username-Regeln auflisten
// @Description  Gibt paginierte, gefilterte und sortierte Username-Regeln zurück
// @Tags         username
// @Produce      json
// @Param        page     query     int     false  "Seite (beginnend bei 1)"
// @Param        limit    query     int     false  "Einträge pro Seite"
// @Param        status   query     string  false  "Status-Filter (allowed, denied, whitelisted)"
// @Param        search   query     string  false  "Suche nach Username"
// @Param        orderBy  query     string  false  "Sortierfeld (ID, Username, Status)"
// @Param        order    query     string  false  "asc oder desc"
// @Success      200 {object} map[string]interface{}
// @Router       /usernames [get]
func GetUsernameRules(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Build WHERE conditions
		var conditions []string
		var args []interface{}

		if status != "" {
			conditions = append(conditions, "status = ?")
			args = append(args, status)
		}
		if search != "" {
			conditions = append(conditions, "username LIKE ?")
			args = append(args, "%"+search+"%")
		}

		// Build WHERE clause
		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		// Validate orderBy and order
		if orderBy != "ID" && orderBy != "Username" && orderBy != "Status" {
			orderBy = "ID"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM username_rules 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, orderBy, order)

		// Add pagination parameters
		args = append(args, limitNum, (pageNum-1)*limitNum)

		// Execute query
		var results []struct {
			models.UsernameRule
			TotalCount int64 `json:"total_count"`
		}

		if err := db.Raw(query, args...).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch username rules"})
			return
		}

		// Extract data
		var rules []models.UsernameRule
		var total int64
		if len(results) > 0 {
			total = results[0].TotalCount
			for _, result := range results {
				rules = append(rules, result.UsernameRule)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"items": rules,
			"total": total,
		})
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
		rule.IsRegex = input.IsRegex
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

// Count-Stats für CharsetRules
func GetCharsetStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var total, allowed, denied, whitelisted int64
		db.Model(&models.CharsetRule{}).Count(&total)
		db.Model(&models.CharsetRule{}).Where("status = ?", "allowed").Count(&allowed)
		db.Model(&models.CharsetRule{}).Where("status = ?", "denied").Count(&denied)
		db.Model(&models.CharsetRule{}).Where("status = ?", "whitelisted").Count(&whitelisted)
		c.JSON(http.StatusOK, gin.H{
			"total":       total,
			"allowed":     allowed,
			"denied":      denied,
			"whitelisted": whitelisted,
		})
	}
}

// Count-Stats für UsernameRules
func GetUsernameStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var total, allowed, denied, whitelisted int64
		db.Model(&models.UsernameRule{}).Count(&total)
		db.Model(&models.UsernameRule{}).Where("status = ?", "allowed").Count(&allowed)
		db.Model(&models.UsernameRule{}).Where("status = ?", "denied").Count(&denied)
		db.Model(&models.UsernameRule{}).Where("status = ?", "whitelisted").Count(&whitelisted)
		c.JSON(http.StatusOK, gin.H{
			"total":       total,
			"allowed":     allowed,
			"denied":      denied,
			"whitelisted": whitelisted,
		})
	}
}

// RecreateIPIndex löscht und erstellt den IP-Index neu
// @Summary      IP-Index neu erstellen
// @Description  Löscht den IP-Index und erstellt ihn mit allen Daten aus der Datenbank neu
// @Tags         ip
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /ip/recreate-index [post]
func RecreateIPIndex(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Delete existing index
		if err := services.DeleteIPIndex(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IP index: " + err.Error()})
			return
		}

		// Recreate index with all data
		if err := services.SyncAllIPs(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recreate IP index: " + err.Error()})
			return
		}

		// Count records indexed
		var recordCount int64
		db.Model(&models.IP{}).Count(&recordCount)

		c.JSON(http.StatusOK, gin.H{
			"message":         "IP index recreated successfully",
			"records_indexed": recordCount,
		})
	}
}

// RecreateEmailIndex löscht und erstellt den Email-Index neu
// @Summary      Email-Index neu erstellen
// @Description  Löscht den Email-Index und erstellt ihn mit allen Daten aus der Datenbank neu
// @Tags         emails
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /emails/recreate-index [post]
func RecreateEmailIndex(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Delete existing index
		if err := services.DeleteEmailIndex(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete email index: " + err.Error()})
			return
		}

		// Recreate index with all data
		if err := services.SyncAllEmails(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recreate email index: " + err.Error()})
			return
		}

		// Count records indexed
		var recordCount int64
		db.Model(&models.Email{}).Count(&recordCount)

		c.JSON(http.StatusOK, gin.H{
			"message":         "Email index recreated successfully",
			"records_indexed": recordCount,
		})
	}
}

// RecreateUserAgentIndex löscht und erstellt den User-Agent-Index neu
// @Summary      User-Agent-Index neu erstellen
// @Description  Löscht den User-Agent-Index und erstellt ihn mit allen Daten aus der Datenbank neu
// @Tags         user-agents
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /user-agents/recreate-index [post]
func RecreateUserAgentIndex(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Delete existing index
		if err := services.DeleteUserAgentIndex(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user agent index: " + err.Error()})
			return
		}

		// Recreate index with all data
		if err := services.SyncAllUserAgents(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recreate user agent index: " + err.Error()})
			return
		}

		// Count records indexed
		var recordCount int64
		db.Model(&models.UserAgent{}).Count(&recordCount)

		c.JSON(http.StatusOK, gin.H{
			"message":         "User agent index recreated successfully",
			"records_indexed": recordCount,
		})
	}
}

// RecreateCountryIndex löscht und erstellt den Country-Index neu
// @Summary      Country-Index neu erstellen
// @Description  Löscht den Country-Index und erstellt ihn mit allen Daten aus der Datenbank neu
// @Tags         countries
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /countries/recreate-index [post]
func RecreateCountryIndex(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Delete existing index
		if err := services.DeleteCountryIndex(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete country index: " + err.Error()})
			return
		}

		// Recreate index with all data
		if err := services.SyncAllCountries(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recreate country index: " + err.Error()})
			return
		}

		// Count records indexed
		var recordCount int64
		db.Model(&models.Country{}).Count(&recordCount)

		c.JSON(http.StatusOK, gin.H{
			"message":         "Country index recreated successfully",
			"records_indexed": recordCount,
		})
	}
}

// RecreateCharsetIndex löscht und erstellt den Charset-Index neu
// @Summary      Charset-Index neu erstellen
// @Description  Löscht den Charset-Index und erstellt ihn mit allen Daten aus der Datenbank neu
// @Tags         charsets
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /charsets/recreate-index [post]
func RecreateCharsetIndex(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Delete existing index
		if err := services.DeleteCharsetIndex(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete charset index: " + err.Error()})
			return
		}

		// Recreate index with all data
		if err := services.SyncAllCharsetRules(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recreate charset index: " + err.Error()})
			return
		}

		// Count records indexed
		var recordCount int64
		db.Model(&models.CharsetRule{}).Count(&recordCount)

		c.JSON(http.StatusOK, gin.H{
			"message":         "Charset index recreated successfully",
			"records_indexed": recordCount,
		})
	}
}

// RecreateUsernameIndex löscht und erstellt den Username-Index neu
// @Summary      Username-Index neu erstellen
// @Description  Löscht den Username-Index und erstellt ihn mit allen Daten aus der Datenbank neu
// @Tags         usernames
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /usernames/recreate-index [post]
func RecreateUsernameIndex(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Delete existing index
		if err := services.DeleteUsernameIndex(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete username index: " + err.Error()})
			return
		}

		// Recreate index with all data
		if err := services.SyncAllUsernameRules(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recreate username index: " + err.Error()})
			return
		}

		// Count records indexed
		var recordCount int64
		db.Model(&models.UsernameRule{}).Count(&recordCount)

		c.JSON(http.StatusOK, gin.H{
			"message":         "Username index recreated successfully",
			"records_indexed": recordCount,
		})
	}
}

// ManualFullSync performs a manual full sync of all data to Elasticsearch
// @Summary      Manual full sync
// @Description  Performs a full sync of all data from MySQL to Elasticsearch
// @Tags         sync
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /sync/full [post]
func ManualFullSync(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Manual full sync requested...")

		// Create incremental sync instance to update timestamps
		incrementalSync := services.NewIncrementalSync()

		// Perform full sync and update timestamps
		if err := incrementalSync.ForceFullSync(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Full sync failed: " + err.Error()})
			return
		}

		// Count total records synced
		var totalRecords int64
		db.Model(&models.IP{}).Count(&totalRecords)
		var emailCount int64
		db.Model(&models.Email{}).Count(&emailCount)
		totalRecords += emailCount
		var userAgentCount int64
		db.Model(&models.UserAgent{}).Count(&userAgentCount)
		totalRecords += userAgentCount
		var countryCount int64
		db.Model(&models.Country{}).Count(&countryCount)
		totalRecords += countryCount
		var charsetCount int64
		db.Model(&models.CharsetRule{}).Count(&charsetCount)
		totalRecords += charsetCount
		var usernameCount int64
		db.Model(&models.UsernameRule{}).Count(&usernameCount)
		totalRecords += usernameCount

		c.JSON(http.StatusOK, gin.H{
			"message":        "Full sync completed successfully",
			"records_synced": totalRecords,
		})
	}
}
