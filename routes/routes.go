// Package routes routes/routes.go
package routes

import (
	"firewall/config"
	"firewall/controllers"
	"firewall/models"
	"firewall/services"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	db := config.DB

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// API routes group
	api := r.Group("/api")

	// CRUD routes for IPs
	api.POST("/ip", controllers.CreateIPAddress(db))
	api.GET("/ips", controllers.GetIPAddresses(db))
	api.PUT("/ip/:id", controllers.UpdateIPAddress(db))
	api.DELETE("/ip/:id", controllers.DeleteIPAddress(db))
	api.GET("/ips/stats", controllers.GetIPStats(db))
	api.POST("/ip/recreate-index", controllers.RecreateIPIndex(db))

	// CRUD routes for Emails
	api.POST("/email", controllers.CreateEmail(db))
	api.GET("/emails", controllers.GetEmails(db))
	api.PUT("/email/:id", controllers.UpdateEmail(db))
	api.DELETE("/email/:id", controllers.DeleteEmail(db))
	api.GET("/emails/stats", controllers.GetEmailStats(db))
	api.POST("/emails/recreate-index", controllers.RecreateEmailIndex(db))

	// CRUD routes for User Agents
	api.POST("/user-agent", controllers.CreateUserAgent(db))
	api.GET("/user-agents", controllers.GetUserAgents(db))
	api.PUT("/user-agent/:id", controllers.UpdateUserAgent(db))
	api.DELETE("/user-agent/:id", controllers.DeleteUserAgent(db))
	api.GET("/user-agents/stats", controllers.GetUserAgentStats(db))
	api.POST("/user-agents/recreate-index", controllers.RecreateUserAgentIndex(db))

	// CRUD routes for Countries
	api.POST("/country", controllers.CreateCountry(db))
	api.GET("/countries", controllers.GetCountries(db))
	api.PUT("/country/:id", controllers.UpdateCountry(db))
	api.DELETE("/country/:id", controllers.DeleteCountry(db))
	api.GET("/countries/stats", controllers.GetCountryStats(db))
	api.POST("/countries/recreate-index", controllers.RecreateCountryIndex(db))

	// CharsetRule CRUD
	api.POST("/charset", controllers.CreateCharsetRule(db))
	api.GET("/charsets", controllers.GetCharsetRules(db))
	api.PUT("/charset/:id", controllers.UpdateCharsetRule(db))
	api.DELETE("/charset/:id", controllers.DeleteCharsetRule(db))
	api.GET("/charsets/stats", controllers.GetCharsetStats(db))
	api.POST("/charsets/recreate-index", controllers.RecreateCharsetIndex(db))

	// UsernameRule CRUD
	api.POST("/username", controllers.CreateUsernameRule(db))
	api.GET("/usernames", controllers.GetUsernameRules(db))
	api.PUT("/username/:id", controllers.UpdateUsernameRule(db))
	api.DELETE("/username/:id", controllers.DeleteUsernameRule(db))
	api.GET("/usernames/stats", controllers.GetUsernameStats(db))
	api.POST("/usernames/recreate-index", controllers.RecreateUsernameIndex(db))

	// Filtering route
	api.POST("/filter", controllers.FilterRequestHandler(db))

	// Manual sync routes
	api.POST("/sync", func(c *gin.Context) {
		if err := services.SyncAllData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync data"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Data synced successfully"})
	})

	api.POST("/sync/full", controllers.ManualFullSync(db))

	// Force sync route
	api.POST("/sync/force", func(c *gin.Context) {
		scheduledSync := services.GetScheduledSync()
		if err := scheduledSync.ForceSync(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to force sync"})
			return
		}

		// Count records that were synced (this is a full sync, so count all records)
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
			"message":        "Force sync completed successfully",
			"records_synced": totalRecords,
		})
	})

	// Health check routes
	api.GET("/health", controllers.HealthCheckHandler(db))
	api.GET("/health/simple", controllers.SimpleHealthCheckHandler())

	// Service status route
	api.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"event_processor": "running",
			"retry_queue":     "running",
			"scheduled_sync":  "running",
			"last_sync":       time.Now().Format(time.RFC3339),
		})
	})

	api.GET("/system-stats", controllers.SystemStatsHandler(db))
	api.POST("/sync/charsets", controllers.SyncCharsetsHandler(db))
	api.POST("/sync/usernames", controllers.SyncUsernamesHandler(db))

	// Analytics routes
	api.GET("/analytics/traffic", controllers.GetTrafficAnalytics(db))
	api.GET("/analytics/relationships", controllers.GetDataRelationships(db))
	api.GET("/analytics/aggregations", controllers.GetAnalyticsAggregations(db))
	api.GET("/analytics/logs", controllers.GetTrafficLogs(db))
	api.GET("/analytics/top-data/:type", controllers.GetTopData(db))
	api.GET("/analytics/insights", controllers.GetRelationshipInsights(db))
	api.GET("/analytics/stats", controllers.GetTrafficStats(db))
	api.POST("/analytics/cleanup", controllers.CleanupOldLogs(db))

	// Cache management route
	api.POST("/cache/flush", func(c *gin.Context) {
		cache := services.GetCacheFactory()
		stats, err := cache.Stats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cache stats"})
			return
		}
		itemsCleared := stats["items"].(int)
		cache.Clear()
		c.JSON(http.StatusOK, gin.H{
			"message":       "Cache flushed successfully",
			"items_cleared": itemsCleared,
		})
	})

	// Sync status route
	api.GET("/sync/status", func(c *gin.Context) {
		distributedLock := services.GetDistributedLock()

		// Get lock information
		fullSyncLock, _ := distributedLock.GetLockInfo("full_sync")
		incrementalSyncLock, _ := distributedLock.GetLockInfo("incremental_sync")

		// Get all active locks
		activeLocks, _ := distributedLock.GetActiveLocks()

		c.JSON(http.StatusOK, gin.H{
			"full_sync_running":           services.IsFullSyncRunning(),
			"distributed_locking_enabled": config.AppConfig.Locking.Enabled,
			"locks": gin.H{
				"full_sync":        fullSyncLock,
				"incremental_sync": incrementalSyncLock,
				"active_locks":     activeLocks,
			},
		})
	})
}
