package controllers

import (
	"net/http"
	"strconv"
	"time"

	"firewall/models"
	"firewall/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTrafficAnalytics returns traffic analytics for a given period
func GetTrafficAnalytics(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get query parameters
		period := c.DefaultQuery("period", "24h")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

		var startTime time.Time
		switch period {
		case "1h":
			startTime = time.Now().Add(-1 * time.Hour)
		case "24h":
			startTime = time.Now().Add(-24 * time.Hour)
		case "7d":
			startTime = time.Now().Add(-7 * 24 * time.Hour)
		case "30d":
			startTime = time.Now().Add(-30 * 24 * time.Hour)
		default:
			startTime = time.Now().Add(-24 * time.Hour)
		}

		// Get traffic logs
		var logs []models.TrafficLog
		if err := db.Where("timestamp >= ?", startTime).Limit(limit).Order("timestamp DESC").Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traffic logs"})
			return
		}

		// Calculate summary statistics
		var total, allowed, denied, whitelisted int64
		var totalResponseTime int64
		var cacheHits int64

		for _, log := range logs {
			total++
			switch log.FinalResult {
			case "allowed":
				allowed++
			case "denied":
				denied++
			case "whitelisted":
				whitelisted++
			}
			totalResponseTime += int64(log.ResponseTimeMs)
			if log.CacheHit {
				cacheHits++
			}
		}

		avgResponseTime := float64(0)
		cacheHitRate := float64(0)
		if total > 0 {
			avgResponseTime = float64(totalResponseTime) / float64(total)
			cacheHitRate = float64(cacheHits) / float64(total) * 100
		}

		c.JSON(http.StatusOK, gin.H{
			"period":               period,
			"total_requests":       total,
			"allowed_requests":     allowed,
			"denied_requests":      denied,
			"whitelisted_requests": whitelisted,
			"avg_response_time_ms": avgResponseTime,
			"cache_hit_rate":       cacheHitRate,
			"logs":                 logs,
		})
	}
}

// GetDataRelationships returns data relationships with filtering
func GetDataRelationships(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		relationshipType := c.Query("type")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

		query := db.Model(&models.DataRelationship{})
		if relationshipType != "" {
			query = query.Where("relationship_type = ?", relationshipType)
		}

		var relationships []models.DataRelationship
		if err := query.Order("frequency DESC").Limit(limit).Find(&relationships).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch relationships"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"relationships": relationships,
			"total":         len(relationships),
		})
	}
}

// GetAnalyticsAggregations returns analytics aggregations
func GetAnalyticsAggregations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		aggregationType := c.DefaultQuery("type", "hourly")
		days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

		startDate := time.Now().AddDate(0, 0, -days)

		var aggregations []models.AnalyticsAggregation
		if err := db.Where("aggregation_type = ? AND aggregation_date >= ?", aggregationType, startDate).
			Order("aggregation_date DESC").Find(&aggregations).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch aggregations"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"aggregations": aggregations,
			"type":         aggregationType,
			"days":         days,
		})
	}
}

// GetTrafficLogs returns paginated traffic logs with filtering
func GetTrafficLogs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
		offset := (page - 1) * limit

		// Build filters
		filters := make(map[string]string)
		if ip := c.Query("ip_address"); ip != "" {
			filters["ip_address"] = ip
		}
		if email := c.Query("email"); email != "" {
			filters["email"] = email
		}
		if userAgent := c.Query("user_agent"); userAgent != "" {
			filters["user_agent"] = userAgent
		}
		if username := c.Query("username"); username != "" {
			filters["username"] = username
		}
		if country := c.Query("country"); country != "" {
			filters["country"] = country
		}
		if asn := c.Query("asn"); asn != "" {
			filters["asn"] = asn
		}
		if result := c.Query("final_result"); result != "" {
			filters["final_result"] = result
		}
		if startDate := c.Query("start_date"); startDate != "" {
			filters["start_date"] = startDate
		}
		if endDate := c.Query("end_date"); endDate != "" {
			filters["end_date"] = endDate
		}

		// Get traffic logs using service
		trafficLogging := services.NewTrafficLoggingService(db)
		logs, total, err := trafficLogging.GetTrafficLogs(limit, offset, filters)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traffic logs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"logs":        logs,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (int(total) + limit - 1) / limit,
		})
	}
}

// GetTopData returns top data for a specific type and period
func GetTopData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		dataType := c.Param("type") // ip_address, email, user_agent, username, country, charset
		period := c.DefaultQuery("period", "24h")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		var startTime, endTime time.Time
		endTime = time.Now()

		switch period {
		case "1h":
			startTime = endTime.Add(-1 * time.Hour)
		case "24h":
			startTime = endTime.Add(-24 * time.Hour)
		case "7d":
			startTime = endTime.Add(-7 * 24 * time.Hour)
		case "30d":
			startTime = endTime.Add(-30 * 24 * time.Hour)
		default:
			startTime = endTime.Add(-24 * time.Hour)
		}

		// Get top data using service
		trafficLogging := services.NewTrafficLoggingService(db)
		analytics := services.NewAnalyticsService(db, trafficLogging)

		topData, err := analytics.GetTopDataByPeriod(startTime, endTime, dataType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top data"})
			return
		}

		// Limit results
		if len(topData) > limit {
			topData = topData[:limit]
		}

		c.JSON(http.StatusOK, gin.H{
			"data_type": dataType,
			"period":    period,
			"top_data":  topData,
		})
	}
}

// GetRelationshipInsights returns relationship insights for a period
func GetRelationshipInsights(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		period := c.DefaultQuery("period", "24h")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		var startTime, endTime time.Time
		endTime = time.Now()

		switch period {
		case "1h":
			startTime = endTime.Add(-1 * time.Hour)
		case "24h":
			startTime = endTime.Add(-24 * time.Hour)
		case "7d":
			startTime = endTime.Add(-7 * 24 * time.Hour)
		case "30d":
			startTime = endTime.Add(-30 * 24 * time.Hour)
		default:
			startTime = endTime.Add(-24 * time.Hour)
		}

		// Get relationship insights using service
		trafficLogging := services.NewTrafficLoggingService(db)
		analytics := services.NewAnalyticsService(db, trafficLogging)

		insights, err := analytics.GetRelationshipInsights(startTime, endTime, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch relationship insights"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"period":   period,
			"insights": insights,
		})
	}
}

// GetTrafficStats returns traffic statistics for a period
func GetTrafficStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		period := c.DefaultQuery("period", "24h")

		var startTime, endTime time.Time
		endTime = time.Now()

		switch period {
		case "1h":
			startTime = endTime.Add(-1 * time.Hour)
		case "24h":
			startTime = endTime.Add(-24 * time.Hour)
		case "7d":
			startTime = endTime.Add(-7 * 24 * time.Hour)
		case "30d":
			startTime = endTime.Add(-30 * 24 * time.Hour)
		default:
			startTime = endTime.Add(-24 * time.Hour)
		}

		// Get traffic stats using service
		trafficLogging := services.NewTrafficLoggingService(db)
		stats, err := trafficLogging.GetTrafficStats(startTime, endTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traffic stats"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"period": period,
			"stats":  stats,
		})
	}
}

// GetTrafficLogStats returns statistics for traffic logs (for dropdown counts)
func GetTrafficLogStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stats struct {
			Total       int64 `json:"total"`
			Allowed     int64 `json:"allowed"`
			Denied      int64 `json:"denied"`
			Whitelisted int64 `json:"whitelisted"`
		}

		// Get total count
		db.Model(&models.TrafficLog{}).Count(&stats.Total)

		// Get counts by final_result
		db.Model(&models.TrafficLog{}).Where("final_result = ?", "allowed").Count(&stats.Allowed)
		db.Model(&models.TrafficLog{}).Where("final_result = ?", "denied").Count(&stats.Denied)
		db.Model(&models.TrafficLog{}).Where("final_result = ?", "whitelisted").Count(&stats.Whitelisted)

		c.JSON(http.StatusOK, stats)
	}
}

// CleanupOldLogs cleans up old traffic logs
func CleanupOldLogs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		retentionDays, _ := strconv.Atoi(c.DefaultQuery("days", "90"))

		trafficLogging := services.NewTrafficLoggingService(db)
		if err := trafficLogging.CleanupOldLogs(retentionDays); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup old logs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "Old logs cleaned up successfully",
			"retention_days": retentionDays,
		})
	}
}
