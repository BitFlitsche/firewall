package main

import (
	"firewall/config"
	"firewall/migrations"
	"firewall/services"
	"log"
	"time"
)

func main() {
	// Initialize configuration
	config.InitConfig()

	// Initialize MySQL
	config.InitMySQL()

	// Run migrations
	if err := migrations.Migrate(config.DB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize traffic logging service
	trafficLogging := services.NewTrafficLoggingService(config.DB)

	// Test traffic logging
	log.Println("Testing traffic logging...")

	// Create multiple test requests to populate the database
	testRequests := []services.FilterRequest{
		{
			IPAddress: "192.168.1.1",
			Email:     "test1@example.com",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Username:  "user1",
			Country:   "US",
			Charset:   "UTF-8",
			Content:   "Test content 1",
		},
		{
			IPAddress: "192.168.1.2",
			Email:     "test2@example.com",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
			Username:  "user2",
			Country:   "DE",
			Charset:   "UTF-8",
			Content:   "Test content 2",
		},
		{
			IPAddress: "192.168.1.3",
			Email:     "test3@example.com",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
			Username:  "user3",
			Country:   "FR",
			Charset:   "UTF-8",
			Content:   "Test content 3",
		},
	}

	// Create metadata
	metadata := map[string]string{
		"client_ip":      "127.0.0.1",
		"user_agent_raw": "Test User Agent",
		"session_id":     "test-session-123",
	}

	// Log multiple requests
	for i, req := range testRequests {
		// Create test result
		result := services.TrafficFilterResult{
			FinalResult: "allowed",
			FilterResults: map[string]interface{}{
				"result": "allowed",
				"reason": "no matches found",
			},
			ResponseTime: time.Duration(50+i*10) * time.Millisecond,
			CacheHit:     i%2 == 0, // Alternate cache hits
		}

		// Log the request
		if err := trafficLogging.LogFilterRequest(req, result, metadata); err != nil {
			log.Printf("Error logging traffic %d: %v", i+1, err)
		} else {
			log.Printf("Traffic %d logged successfully", i+1)
		}

		// Small delay between requests
		time.Sleep(100 * time.Millisecond)
	}

	// Test analytics service
	log.Println("Testing analytics service...")
	analytics := services.NewAnalyticsService(config.DB, trafficLogging)

	// Generate test aggregation
	if err := analytics.GenerateHourlyAggregation(); err != nil {
		log.Printf("Error generating aggregation: %v", err)
	} else {
		log.Println("Analytics aggregation generated successfully")
	}

	// Get traffic stats
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()
	stats, err := trafficLogging.GetTrafficStats(startTime, endTime)
	if err != nil {
		log.Printf("Error getting traffic stats: %v", err)
	} else {
		log.Printf("Traffic stats: %+v", stats)
	}

	// Check if data was actually written to database
	var count int64
	config.DB.Table("traffic_logs").Count(&count)
	log.Printf("Total traffic logs in database: %d", count)

	var relationshipCount int64
	config.DB.Table("data_relationships").Count(&relationshipCount)
	log.Printf("Total data relationships in database: %d", relationshipCount)

	log.Println("Traffic logging test completed successfully!")
}
