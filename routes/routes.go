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
	r.PUT("/ip/:id", controllers.UpdateIPAddress(db))
	r.DELETE("/ip/:id", controllers.DeleteIPAddress(db))
	r.GET("/ips/stats", controllers.GetIPStats(db))
	r.POST("/ip/recreate-index", controllers.RecreateIPIndex(db))

	// CRUD routes for Emails
	r.POST("/email", controllers.CreateEmail(db))
	r.GET("/emails", controllers.GetEmails(db))
	r.PUT("/email/:id", controllers.UpdateEmail(db))
	r.DELETE("/email/:id", controllers.DeleteEmail(db))
	r.GET("/emails/stats", controllers.GetEmailStats(db))
	r.POST("/emails/recreate-index", controllers.RecreateEmailIndex(db))

	// CRUD routes for User Agents
	r.POST("/user-agent", controllers.CreateUserAgent(db))
	r.GET("/user-agents", controllers.GetUserAgents(db))
	r.PUT("/user-agent/:id", controllers.UpdateUserAgent(db))
	r.DELETE("/user-agent/:id", controllers.DeleteUserAgent(db))
	r.GET("/user-agents/stats", controllers.GetUserAgentStats(db))
	r.POST("/user-agents/recreate-index", controllers.RecreateUserAgentIndex(db))

	// CRUD routes for Countries
	r.POST("/country", controllers.CreateCountry(db))
	r.GET("/countries", controllers.GetCountries(db))
	r.PUT("/country/:id", controllers.UpdateCountry(db))
	r.DELETE("/country/:id", controllers.DeleteCountry(db))
	r.GET("/countries/stats", controllers.GetCountryStats(db))
	r.POST("/countries/recreate-index", controllers.RecreateCountryIndex(db))

	// CharsetRule CRUD
	r.POST("/charset", controllers.CreateCharsetRule(db))
	r.GET("/charsets", controllers.GetCharsetRules(db))
	r.PUT("/charset/:id", controllers.UpdateCharsetRule(db))
	r.DELETE("/charset/:id", controllers.DeleteCharsetRule(db))
	r.GET("/charsets/stats", controllers.GetCharsetStats(db))
	r.POST("/charsets/recreate-index", controllers.RecreateCharsetIndex(db))

	// UsernameRule CRUD
	r.POST("/username", controllers.CreateUsernameRule(db))
	r.GET("/usernames", controllers.GetUsernameRules(db))
	r.PUT("/username/:id", controllers.UpdateUsernameRule(db))
	r.DELETE("/username/:id", controllers.DeleteUsernameRule(db))
	r.GET("/usernames/stats", controllers.GetUsernameStats(db))
	r.POST("/usernames/recreate-index", controllers.RecreateUsernameIndex(db))

	// Filtering route
	r.POST("/filter", controllers.FilterRequestHandler(db))

	// Manual sync route
	r.POST("/sync", func(c *gin.Context) {
		if err := services.SyncAllData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync data"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Data synced successfully"})
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

	r.GET("/system-stats", controllers.SystemStatsHandler(db))
	r.POST("/sync/charsets", controllers.SyncCharsetsHandler(db))
	r.POST("/sync/usernames", controllers.SyncUsernamesHandler(db))
}
