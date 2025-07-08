// main.go
package main

import (
	"firewall/config"
	"firewall/migrations"
	"firewall/routes"
	"firewall/services"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "firewall/docs" // Swagger-Dokumentation

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Initialize MySQL and Elasticsearch
	config.InitMySQL()
	// Run migrations
	if err := migrations.Migrate(config.DB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	config.InitElasticsearch()

	// Initialize all services
	log.Println("Initializing services...")

	// Initialize event processor
	eventProcessor := services.GetEventProcessor()

	// Initialize retry queue
	retryQueue := services.GetRetryQueue()

	// Initialize scheduled sync
	scheduledSync := services.GetScheduledSync()

	// Initial sync of existing data
	log.Println("Performing initial sync of existing data...")
	if err := services.SyncAllData(); err != nil {
		log.Printf("Warning: Initial sync failed: %v", err)
	}

	// Set up Gin and routes
	r := gin.Default()
	routes.SetupRoutes(r)

	// Swagger-UI Route
	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler()))
	// Die Swagger-UI ist jetzt erreichbar unter: http://localhost:8081/swagger/index.html

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :8081")
		if err := r.Run(":8081"); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down server...")

	// Stop all services gracefully
	log.Println("Stopping services...")

	// Stop scheduled sync
	scheduledSync.Stop()

	// Stop retry queue
	retryQueue.Stop()

	// Stop event processor
	eventProcessor.Stop()

	log.Println("Server stopped gracefully")
}
