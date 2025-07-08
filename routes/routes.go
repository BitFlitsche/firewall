// Package routes routes/routes.go
package routes

import (
	"firewall/config"
	"firewall/controllers"
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

	// CRUD routes for IPs
	r.POST("/ip", controllers.CreateIPAddress(db))
	r.GET("/ips", controllers.GetIPAddresses(db))

	// CRUD routes for Emails
	r.POST("/email", controllers.CreateEmail(db))
	r.GET("/emails", controllers.GetEmails(db))

	// CRUD routes for User Agents
	r.POST("/user-agent", controllers.CreateUserAgent(db))
	r.GET("/user-agents", controllers.GetUserAgents(db))

	// CRUD routes for Countries
	r.POST("/country", controllers.CreateCountry(db))
	r.GET("/countries", controllers.GetCountries(db))

	// Filtering route
	r.POST("/filter", controllers.FilterRequestHandler())

	// Manual sync route
	r.POST("/sync", func(c *gin.Context) {
		if err := services.SyncAllData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync data"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Data synced successfully"})
	})

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		// Check MySQL connection
		sqlDB, err := config.DB.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "mysql": "disconnected"})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "mysql": "disconnected"})
			return
		}

		// Check Elasticsearch connection
		es := config.ESClient
		res, err := es.Info()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "elasticsearch": "disconnected"})
			return
		}
		defer res.Body.Close()

		c.JSON(http.StatusOK, gin.H{
			"status":        "healthy",
			"mysql":         "connected",
			"elasticsearch": "connected",
			"services": gin.H{
				"event_processor": "running",
				"retry_queue":     "running",
				"scheduled_sync":  "running",
			},
		})
	})

	// Force sync route
	r.POST("/sync/force", func(c *gin.Context) {
		scheduledSync := services.GetScheduledSync()
		if err := scheduledSync.ForceSync(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to force sync"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Force sync completed successfully"})
	})

	// Service status route
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"event_processor": "running",
			"retry_queue":     "running",
			"scheduled_sync":  "running",
			"last_sync":       time.Now().Format(time.RFC3339),
		})
	})
}
