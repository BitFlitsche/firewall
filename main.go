// main.go
package main

import (
	"firewall/config"
	"firewall/migrations"
	"firewall/routes"
	"firewall/services"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "firewall/docs" // Swagger-Dokumentation

	"firewall/controllers"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Initialize configuration
	config.InitConfig()

	// Initialize MySQL and Elasticsearch
	config.InitMySQL()
	// Run migrations
	if err := migrations.Migrate(config.DB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	config.InitElasticsearch()

	// Initialize all services
	log.Println("Initializing services...")

	// Initialize cache factory (switches between in-memory and distributed based on config)
	_ = services.GetCacheFactory()

	// Initialize event processor
	eventProcessor := services.GetEventProcessor()

	// Initialize retry queue
	retryQueue := services.GetRetryQueue()

	// Initialize distributed lock service
	distributedLock := services.GetDistributedLock()

	// Initialize scheduled sync
	scheduledSync := services.GetScheduledSync()

	// Initial sync of existing data
	log.Println("Performing initial sync of existing data...")
	if err := services.SyncAllData(); err != nil {
		log.Printf("Warning: Initial sync failed: %v", err)
	}
	if err := services.SyncAllCharsetsToES(config.DB); err != nil {
		log.Printf("Warning: Initial charset sync failed: %v", err)
	}

	// Set up Gin and routes
	r := gin.Default()
	r.Use(controllers.MetricsMiddleware())
	routes.SetupRoutes(r)

	// Serve React build static files
	r.Static("/static", "./firewall-app/build/static")
	r.StaticFile("/favicon.ico", "./firewall-app/build/favicon.ico")
	r.StaticFile("/manifest.json", "./firewall-app/build/manifest.json")
	r.StaticFile("/logo192.png", "./firewall-app/build/logo192.png")
	r.StaticFile("/logo512.png", "./firewall-app/build/logo512.png")

	// Fallback: serve index.html for all other routes (client-side routing)
	r.NoRoute(func(c *gin.Context) {
		c.File("./firewall-app/build/index.html")
	})

	// Swagger-UI Route
	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)))
	// Die Swagger-UI ist jetzt erreichbar unter: http://localhost:8081/swagger/index.html

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		serverAddr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
		log.Printf("Starting server on %s", serverAddr)
		if err := r.Run(serverAddr); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down server...")

	// Stop all services gracefully
	log.Println("Stopping services...")

	// Stop distributed lock service
	distributedLock.Stop()

	// Stop scheduled sync
	scheduledSync.Stop()

	// Stop retry queue
	retryQueue.Stop()

	// Stop event processor
	eventProcessor.Stop()

	log.Println("Server stopped gracefully")
}
