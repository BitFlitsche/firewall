// controllers/rules.go
package controllers

import (
	"firewall/models"
	"firewall/services"
	"firewall/utils"
	"firewall/validation"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
var isImportRunning sync.Map

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
			"geo_cache":  services.GetGeoCacheStats(),
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
// @Failure      409 {object}  map[string]string
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

		// Check for conflicts with existing entries
		var existingIPs []models.IP
		if err := db.Find(&existingIPs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for conflicts"})
			return
		}

		// Extract existing IPs and CIDR ranges with their statuses
		var existingIPAddresses []string
		var existingCIDRs []string
		existingStatuses := make(map[string]string)

		for _, existing := range existingIPs {
			if existing.IsCIDR {
				existingCIDRs = append(existingCIDRs, existing.Address)
				existingStatuses[existing.Address] = existing.Status
			} else {
				existingIPAddresses = append(existingIPAddresses, existing.Address)
				existingStatuses[existing.Address] = existing.Status
			}
		}

		// Check conflicts based on whether new entry is IP or CIDR
		var conflicts []utils.ConflictInfo
		var err error

		if ip.IsCIDR {
			// New entry is a CIDR range - check for conflicts
			conflicts, err = utils.CheckCIDRConflicts(ip.Address, existingIPAddresses, existingCIDRs, existingStatuses, ip.Status)
		} else {
			// New entry is an IP address - check if it's covered by existing CIDR ranges
			conflicts, err = utils.CheckIPConflicts(ip.Address, existingCIDRs, existingStatuses, ip.Status)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check conflicts", "details": err.Error()})
			return
		}

		// If there are conflicts, return them
		if len(conflicts) > 0 {
			// Check if any conflicts are errors (not just warnings)
			hasErrors := false
			for _, conflict := range conflicts {
				if conflict.Severity == "error" {
					hasErrors = true
					break
				}
			}

			if hasErrors {
				// Build detailed error message with conflicting records
				var conflictDetails []string
				for _, conflict := range conflicts {
					conflictDetails = append(conflictDetails, conflict.Message)
				}

				errorMessage := "IP/CIDR conflicts detected"
				if len(conflictDetails) > 0 {
					errorMessage = fmt.Sprintf("IP/CIDR conflicts detected: %s", strings.Join(conflictDetails, "; "))
				}

				c.JSON(http.StatusConflict, gin.H{
					"error":     errorMessage,
					"conflicts": conflicts,
					"message":   "Please review conflicts before proceeding",
				})
				return
			} else {
				// Only warnings - allow creation but inform user
				var conflictDetails []string
				for _, conflict := range conflicts {
					conflictDetails = append(conflictDetails, conflict.Message)
				}

				warningMessage := "IP/CIDR overlaps detected"
				if len(conflictDetails) > 0 {
					warningMessage = fmt.Sprintf("IP/CIDR overlaps detected: %s", strings.Join(conflictDetails, "; "))
				}

				// Log warning but continue with creation
				log.Printf("Warning: %s", warningMessage)
			}
		}

		// Check if IP already exists (exact match)
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
		orderBy := c.DefaultQuery("orderBy", "id")
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
		if orderBy != "id" && orderBy != "address" && orderBy != "status" {
			orderBy = "id"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Map frontend field names to database column names
		orderByMap := map[string]string{
			"id":      "ID",
			"address": "Address",
			"status":  "Status",
		}
		dbOrderBy := orderByMap[orderBy]

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
				conditions = append(conditions, "is_c_id_r = ?")
				args = append(args, false)
			} else if typeFilter == "cidr" {
				conditions = append(conditions, "is_c_id_r = ?")
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

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM ips 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, dbOrderBy, order)

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

		// Use the correct column name for CIDR
		db.Raw("SELECT COUNT(*) FROM ips WHERE is_c_id_r = 0").Scan(&single)
		db.Raw("SELECT COUNT(*) FROM ips WHERE is_c_id_r = 1").Scan(&cidr)

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

		// Check for conflicts with existing entries (excluding current record)
		var existingIPs []models.IP
		if err := db.Where("id != ?", id).Find(&existingIPs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for conflicts"})
			return
		}

		// Extract existing IPs and CIDR ranges with their statuses
		var existingIPAddresses []string
		var existingCIDRs []string
		existingStatuses := make(map[string]string)

		for _, existing := range existingIPs {
			if existing.IsCIDR {
				existingCIDRs = append(existingCIDRs, existing.Address)
				existingStatuses[existing.Address] = existing.Status
			} else {
				existingIPAddresses = append(existingIPAddresses, existing.Address)
				existingStatuses[existing.Address] = existing.Status
			}
		}

		// Check conflicts based on whether new entry is IP or CIDR
		var conflicts []utils.ConflictInfo
		var err error

		if input.IsCIDR {
			// New entry is a CIDR range - check for conflicts
			conflicts, err = utils.CheckCIDRConflicts(input.Address, existingIPAddresses, existingCIDRs, existingStatuses, input.Status)
		} else {
			// New entry is an IP address - check if it's covered by existing CIDR ranges
			conflicts, err = utils.CheckIPConflicts(input.Address, existingCIDRs, existingStatuses, input.Status)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check conflicts", "details": err.Error()})
			return
		}

		// If there are conflicts, return them
		if len(conflicts) > 0 {
			// Check if any conflicts are errors (not just warnings)
			hasErrors := false
			for _, conflict := range conflicts {
				if conflict.Severity == "error" {
					hasErrors = true
					break
				}
			}

			statusCode := http.StatusConflict
			if !hasErrors {
				statusCode = http.StatusOK // If only warnings, still allow update
			}

			// Build detailed error message with conflicting records
			var conflictDetails []string
			for _, conflict := range conflicts {
				conflictDetails = append(conflictDetails, conflict.Message)
			}

			errorMessage := "IP/CIDR conflicts detected"
			if len(conflictDetails) > 0 {
				errorMessage = fmt.Sprintf("IP/CIDR conflicts detected: %s", strings.Join(conflictDetails, "; "))
			}

			c.JSON(statusCode, gin.H{
				"error":     errorMessage,
				"conflicts": conflicts,
				"message":   "Please review conflicts before proceeding",
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
		orderBy := c.DefaultQuery("orderBy", "id")
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
		if orderBy != "id" && orderBy != "address" && orderBy != "status" {
			orderBy = "id"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Map frontend field names to database column names
		orderByMap := map[string]string{
			"id":      "ID",
			"address": "Address",
			"status":  "Status",
		}
		dbOrderBy := orderByMap[orderBy]

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM emails 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, dbOrderBy, order)

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
		orderBy := c.DefaultQuery("orderBy", "id")
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
		if orderBy != "id" && orderBy != "user_agent" && orderBy != "status" {
			orderBy = "id"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Map frontend field names to database column names
		orderByMap := map[string]string{
			"id":         "ID",
			"user_agent": "UserAgent",
			"status":     "Status",
		}
		dbOrderBy := orderByMap[orderBy]

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM user_agents 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, dbOrderBy, order)

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
		_ = services.DeleteUsernameFromES(parseUint(id))
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
		orderBy := c.DefaultQuery("orderBy", "name")
		order := c.DefaultQuery("order", "asc")

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
			// Search in both code and name fields
			conditions = append(conditions, "(code LIKE ? OR name LIKE ?)")
			args = append(args, "%"+search+"%", "%"+search+"%")
		}

		// Build WHERE clause
		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		// Validate orderBy and order
		if orderBy != "id" && orderBy != "code" && orderBy != "name" && orderBy != "status" {
			orderBy = "name"
		}
		if order != "asc" && order != "desc" {
			order = "asc"
		}

		// Map frontend field names to database column names
		orderByMap := map[string]string{
			"id":     "ID",
			"code":   "Code",
			"name":   "Name",
			"status": "Status",
		}
		dbOrderBy := orderByMap[orderBy]

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM countries 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, dbOrderBy, order)

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
		orderBy := c.DefaultQuery("orderBy", "id")
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
		if orderBy != "id" && orderBy != "charset" && orderBy != "status" {
			orderBy = "id"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Map frontend field names to database column names
		orderByMap := map[string]string{
			"id":      "ID",
			"charset": "Charset",
			"status":  "Status",
		}
		dbOrderBy := orderByMap[orderBy]

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM charset_rules 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, dbOrderBy, order)

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
		orderBy := c.DefaultQuery("orderBy", "id")
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
		if orderBy != "id" && orderBy != "username" && orderBy != "status" {
			orderBy = "id"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Map frontend field names to database column names
		orderByMap := map[string]string{
			"id":       "ID",
			"username": "Username",
			"status":   "Status",
		}
		dbOrderBy := orderByMap[orderBy]

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM username_rules 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, dbOrderBy, order)

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

// CheckIPConflicts checks for conflicts when adding a new IP or CIDR
// @Summary      Check IP/CIDR conflicts
// @Description  Checks if a new IP or CIDR would conflict with existing entries
// @Tags         ip
// @Accept       json
// @Produce      json
// @Param        ip  body      models.IP  true  "IP/CIDR to check"
// @Success      200 {object}  map[string]interface{}
// @Failure      400 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /ip/check-conflicts [post]
func CheckIPConflicts(db *gorm.DB) gin.HandlerFunc {
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

		// Get all existing IPs and CIDR ranges
		var existingIPs []models.IP
		if err := db.Find(&existingIPs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for conflicts"})
			return
		}

		// Extract existing IPs and CIDR ranges
		var existingIPAddresses []string
		var existingCIDRs []string
		existingStatuses := make(map[string]string)

		for _, existing := range existingIPs {
			if existing.IsCIDR {
				existingCIDRs = append(existingCIDRs, existing.Address)
				existingStatuses[existing.Address] = existing.Status
			} else {
				existingIPAddresses = append(existingIPAddresses, existing.Address)
				existingStatuses[existing.Address] = existing.Status
			}
		}

		// Check conflicts based on whether new entry is IP or CIDR
		var conflicts []utils.ConflictInfo
		var err error

		if ip.IsCIDR {
			// New entry is a CIDR range - check for conflicts
			conflicts, err = utils.CheckCIDRConflicts(ip.Address, existingIPAddresses, existingCIDRs, existingStatuses, ip.Status)
		} else {
			// New entry is an IP address - check if it's covered by existing CIDR ranges
			conflicts, err = utils.CheckIPConflicts(ip.Address, existingCIDRs, existingStatuses, ip.Status)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check conflicts", "details": err.Error()})
			return
		}

		// Determine overall status
		hasErrors := false
		hasWarnings := false
		for _, conflict := range conflicts {
			if conflict.Severity == "error" {
				hasErrors = true
			} else if conflict.Severity == "warning" {
				hasWarnings = true
			}
		}

		status := "clean"
		if hasErrors {
			status = "error"
		} else if hasWarnings {
			status = "warning"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":         status,
			"conflicts":      conflicts,
			"conflict_count": len(conflicts),
			"can_proceed":    !hasErrors,
			"message": func() string {
				if hasErrors {
					return "Conflicts detected - cannot proceed"
				} else if hasWarnings {
					return "Warnings detected - review before proceeding"
				}
				return "No conflicts detected"
			}(),
		})
	}
}

// Charset Fields Management Controllers

// GetCharsetFields returns the current charset fields configuration
func GetCharsetFields(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fieldsConfig := services.GetCharsetFieldsConfig()
		allFields := fieldsConfig.GetAllFields()

		c.JSON(http.StatusOK, gin.H{
			"standard_fields": allFields["standard"],
			"custom_fields":   allFields["custom"],
		})
	}
}

// ToggleStandardField enables/disables a standard field
func ToggleStandardField(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			FieldName string `json:"field_name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		fieldsConfig := services.GetCharsetFieldsConfig()
		if err := fieldsConfig.ToggleStandardField(request.FieldName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Field toggled successfully"})
	}
}

// AddCustomField adds a new custom field
func AddCustomField(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			FieldName string `json:"field_name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		fieldsConfig := services.GetCharsetFieldsConfig()
		if err := fieldsConfig.AddCustomField(request.FieldName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Custom field added successfully"})
	}
}

// DeleteCustomField removes a custom field
func DeleteCustomField(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fieldName := c.Param("field")
		if fieldName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Field name is required"})
			return
		}

		fieldsConfig := services.GetCharsetFieldsConfig()
		if err := fieldsConfig.DeleteCustomField(fieldName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Custom field deleted successfully"})
	}
}

// ToggleCustomField enables/disables a custom field
func ToggleCustomField(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			FieldName string `json:"field_name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		fieldsConfig := services.GetCharsetFieldsConfig()
		if err := fieldsConfig.ToggleCustomField(request.FieldName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Custom field toggled successfully"})
	}
}

// ASN CRUD Controllers

// CreateASN fügt eine neue ASN-Regel hinzu
// @Summary      Neue ASN-Regel anlegen
// @Description  Legt eine neue ASN-Regel mit Status an
// @Tags         asn
// @Accept       json
// @Produce      json
// @Param        asn  body      models.ASN  true  "ASN-Daten"
// @Success      200 {object}  models.ASN
// @Failure      400 {object}  map[string]string
// @Failure      409 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /asn [post]
func CreateASN(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var asn models.ASN

		if err := c.ShouldBindJSON(&asn); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		// Validate ASN format (should start with "AS" followed by numbers)
		if len(asn.ASN) < 3 || !strings.HasPrefix(asn.ASN, "AS") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ASN must start with 'AS' followed by numbers"})
			return
		}

		// Validate status
		statusValidation := validation.ValidateStatus(asn.Status)
		if !statusValidation.IsValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": statusValidation.Errors,
			})
			return
		}

		// Check for existing ASN
		var existingASN models.ASN
		if err := db.Where("asn = ?", asn.ASN).First(&existingASN).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "ASN already exists"})
			return
		}

		// Create the ASN
		if err := db.Create(&asn).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ASN"})
			return
		}

		// Sync to Elasticsearch
		go func() {
			if err := services.SyncASNToES(asn); err != nil {
				log.Printf("Failed to sync ASN to Elasticsearch: %v", err)
			}
		}()

		c.JSON(http.StatusOK, asn)
	}
}

// GetASNs returns all ASN rules with pagination and filtering
func GetASNs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query-Parameter
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")
		status := c.Query("status")
		rir := c.Query("rir")
		country := c.Query("country")
		search := c.Query("search")
		orderBy := c.DefaultQuery("orderBy", "id")
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
		if orderBy != "id" && orderBy != "asn" && orderBy != "name" && orderBy != "status" {
			orderBy = "id"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		// Map frontend field names to database column names
		orderByMap := map[string]string{
			"id":     "ID",
			"asn":    "ASN",
			"name":   "Name",
			"status": "Status",
		}
		dbOrderBy := orderByMap[orderBy]

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
		if rir != "" {
			conditions = append(conditions, "rir = ?")
			args = append(args, rir)
		}
		if country != "" {
			conditions = append(conditions, "country = ?")
			args = append(args, country)
		}
		if search != "" {
			conditions = append(conditions, "(asn LIKE ? OR domain LIKE ? OR name LIKE ?)")
			args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
		}

		// Build WHERE clause
		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		// Use COUNT(*) OVER() for single query optimization
		query := fmt.Sprintf(`
			SELECT *, COUNT(*) OVER() as total_count 
			FROM asns 
			%s 
			ORDER BY %s %s 
			LIMIT ? OFFSET ?
		`, whereClause, dbOrderBy, order)

		// Add pagination parameters
		args = append(args, limitNum, (pageNum-1)*limitNum)

		// Execute query
		var results []struct {
			models.ASN
			TotalCount int64 `json:"total_count"`
		}

		if err := db.Raw(query, args...).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ASNs"})
			return
		}

		// Extract data
		var asns []models.ASN
		var total int64
		if len(results) > 0 {
			total = results[0].TotalCount
			for _, result := range results {
				asns = append(asns, result.ASN)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"items": asns,
			"total": total,
		})
	}
}

// GetASNStats returns statistics for ASN rules
func GetASNStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stats struct {
			Total       int64 `json:"total"`
			Allowed     int64 `json:"allowed"`
			Denied      int64 `json:"denied"`
			Whitelisted int64 `json:"whitelisted"`
		}

		// Get total count
		db.Model(&models.ASN{}).Count(&stats.Total)

		// Get counts by status
		db.Model(&models.ASN{}).Where("status = ?", "allowed").Count(&stats.Allowed)
		db.Model(&models.ASN{}).Where("status = ?", "denied").Count(&stats.Denied)
		db.Model(&models.ASN{}).Where("status = ?", "whitelisted").Count(&stats.Whitelisted)

		c.JSON(http.StatusOK, stats)
	}
}

// GetASNFilterStats returns statistics for ASN filtering (RIR and Country counts)
func GetASNFilterStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get RIR counts
		var rirStats []struct {
			RIR   string `json:"rir"`
			Count int64  `json:"count"`
		}
		db.Model(&models.ASN{}).
			Select("rir, COUNT(*) as count").
			Where("rir IS NOT NULL AND rir != ''").
			Group("rir").
			Order("count DESC").
			Scan(&rirStats)

		// Get Country counts
		var countryStats []struct {
			Country string `json:"country"`
			Count   int64  `json:"count"`
		}
		db.Model(&models.ASN{}).
			Select("country, COUNT(*) as count").
			Where("country IS NOT NULL AND country != ''").
			Group("country").
			Order("count DESC").
			Scan(&countryStats)

		// Convert to map format for frontend
		rirCounts := make(map[string]int64)
		var totalRIR int64
		for _, stat := range rirStats {
			rirCounts[stat.RIR] = stat.Count
			totalRIR += stat.Count
		}
		rirCounts["total"] = totalRIR

		countryCounts := make(map[string]int64)
		var totalCountry int64
		for _, stat := range countryStats {
			countryCounts[stat.Country] = stat.Count
			totalCountry += stat.Count
		}
		countryCounts["total"] = totalCountry

		c.JSON(http.StatusOK, gin.H{
			"rir_counts":     rirCounts,
			"country_counts": countryCounts,
		})
	}
}

// UpdateASN updates an existing ASN rule
func UpdateASN(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ASN ID is required"})
			return
		}

		var asn models.ASN
		if err := c.ShouldBindJSON(&asn); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		// Validate ASN format
		if len(asn.ASN) < 3 || !strings.HasPrefix(asn.ASN, "AS") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ASN must start with 'AS' followed by numbers"})
			return
		}

		// Validate status
		statusValidation := validation.ValidateStatus(asn.Status)
		if !statusValidation.IsValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": statusValidation.Errors,
			})
			return
		}

		// Check if ASN exists
		var existingASN models.ASN
		if err := db.First(&existingASN, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ASN not found"})
			return
		}

		// Update the ASN
		if err := db.Model(&existingASN).Updates(asn).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ASN"})
			return
		}

		// Sync to Elasticsearch
		go func() {
			if err := services.SyncASNToES(existingASN); err != nil {
				log.Printf("Failed to sync ASN to Elasticsearch: %v", err)
			}
		}()

		c.JSON(http.StatusOK, existingASN)
	}
}

// DeleteASN deletes an ASN rule
func DeleteASN(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ASN ID is required"})
			return
		}

		// Check if ASN exists
		var asn models.ASN
		if err := db.First(&asn, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ASN not found"})
			return
		}

		// Delete the ASN
		if err := db.Delete(&asn).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ASN"})
			return
		}

		// Delete from Elasticsearch
		go func() {
			if err := services.DeleteASNFromES(asn.ID); err != nil {
				log.Printf("Failed to delete ASN from Elasticsearch: %v", err)
			}
		}()

		c.JSON(http.StatusOK, gin.H{"message": "ASN deleted successfully"})
	}
}

// RecreateASNIndex recreates the ASN index in Elasticsearch
func RecreateASNIndex(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Delete existing index
		es := config.ESClient
		_, err := es.Indices.Delete([]string{"asns"}, es.Indices.Delete.WithContext(context.Background()))
		if err != nil {
			// Index might not exist, which is fine
			log.Printf("Could not delete existing ASN index: %v", err)
		}

		// Create new index with mapping
		mapping := `{
			"mappings": {
				"properties": {
					"id": {"type": "long"},
					"asn": {"type": "keyword"},
					"name": {"type": "text"},
					"status": {"type": "keyword"}
				}
			}
		}`

		_, err = es.Indices.Create(
			"asns",
			es.Indices.Create.WithBody(strings.NewReader(mapping)),
			es.Indices.Create.WithContext(context.Background()),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ASN index"})
			return
		}

		// Sync all ASNs to Elasticsearch
		var asns []models.ASN
		if err := db.Find(&asns).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ASNs"})
			return
		}

		for _, asn := range asns {
			if err := services.SyncASNToES(asn); err != nil {
				log.Printf("Failed to sync ASN %d to Elasticsearch: %v", asn.ID, err)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ASN index recreated successfully",
			"synced":  len(asns),
		})
	}
}

// ImportSpamhausASNDrop imports ASN data from Spamhaus ASN-DROP list
func ImportSpamhausASNDrop(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Import ASN data from Spamhaus
		if err := services.ImportSpamhausASNDrop(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import Spamhaus ASN-DROP data", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Spamhaus ASN-DROP data imported successfully"})
	}
}

// GetSpamhausImportStats returns statistics about the Spamhaus import
func GetSpamhausImportStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := services.GetSpamhausImportStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Spamhaus import stats", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}

// GetSpamhausImportStatus returns the current status of Spamhaus import
func GetSpamhausImportStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		isRunning := services.IsSpamhausImportRunning()

		// Get next scheduled import time
		now := time.Now()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

		c.JSON(http.StatusOK, gin.H{
			"is_running":              isRunning,
			"auto_import_enabled":     config.AppConfig.Spamhaus.AutoImportEnabled,
			"next_scheduled":          nextMidnight.Format("2006-01-02 15:04:05"),
			"next_scheduled_relative": time.Until(nextMidnight).String(),
		})
	}
}

// ImportStopForumSpamToxicCIDRs imports toxic IP addresses in CIDR format from StopForumSpam
func ImportStopForumSpamToxicCIDRs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if import is already running
		if _, running := isImportRunning.Load("stopforumspam"); running {
			c.JSON(http.StatusConflict, gin.H{"error": "Import already running"})
			return
		}

		// Set import as running
		isImportRunning.Store("stopforumspam", true)
		defer isImportRunning.Delete("stopforumspam")

		// Create StopForumSpam import service
		stopForumSpamService := services.NewStopForumSpamImportService(db)

		// Run import in a goroutine to avoid blocking
		go func() {
			if err := stopForumSpamService.ImportToxicCIDRs(); err != nil {
				log.Printf("StopForumSpam toxic CIDR import failed: %v", err)
			}
		}()

		c.JSON(http.StatusOK, gin.H{"message": "StopForumSpam toxic CIDR import started"})
	}
}

// GetStopForumSpamImportStats returns statistics about StopForumSpam imports
func GetStopForumSpamImportStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stopForumSpamService := services.NewStopForumSpamImportService(db)

		stats, err := stopForumSpamService.GetStopForumSpamImportStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get StopForumSpam import stats"})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}

// GetStopForumSpamImportStatus returns the current status of StopForumSpam imports
func GetStopForumSpamImportStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stopForumSpamService := services.NewStopForumSpamImportService(db)

		status, err := stopForumSpamService.GetStopForumSpamImportStatus()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get StopForumSpam import status"})
			return
		}

		c.JSON(http.StatusOK, status)
	}
}
