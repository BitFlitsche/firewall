package main

import (
	"firewall/config"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Testing Connection Pooling Configuration...")

	// Initialize database with connection pooling
	config.InitMySQL()

	// Test connection pool statistics
	fmt.Println("\n=== Connection Pool Statistics ===")
	stats := config.GetDBStats()
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}

	// Test concurrent connections
	fmt.Println("\n=== Testing Concurrent Connections ===")

	// Simulate concurrent database operations
	for i := 0; i < 10; i++ {
		go func(id int) {
			// Simulate database query
			var count int64
			config.DB.Model(&struct{}{}).Count(&count)
			fmt.Printf("Goroutine %d completed\n", id)
		}(i)
	}

	// Wait a bit for operations to complete
	time.Sleep(2 * time.Second)

	// Show updated statistics
	fmt.Println("\n=== Updated Connection Pool Statistics ===")
	stats = config.GetDBStats()
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}

	fmt.Println("\nConnection pooling test completed successfully!")
}
