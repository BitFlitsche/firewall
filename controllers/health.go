package controllers

import (
	"context"
	"firewall/config"
	"firewall/services"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthStatus represents the overall health status
// @Description Overall system health status with detailed service information
type HealthStatus struct {
	Status    string                   `json:"status" example:"healthy"`                 // Overall system status: "healthy" or "unhealthy"
	Timestamp string                   `json:"timestamp" example:"2025-07-19T08:34:17Z"` // UTC timestamp of the health check
	Version   string                   `json:"version" example:"1.0.0"`                  // API version
	Services  map[string]ServiceHealth `json:"services"`                                 // Detailed health status of each service
}

// ServiceHealth represents the health of individual services
// @Description Individual service health information
type ServiceHealth struct {
	Status       string `json:"status" example:"healthy"`                      // Service status: "healthy" or "unhealthy"
	Message      string `json:"message,omitempty" example:"Connection failed"` // Error message if service is unhealthy
	ResponseTime int64  `json:"response_time_ms,omitempty" example:"17"`       // Response time in milliseconds
}

// HealthCheckHandler provides a comprehensive health check endpoint
// @Summary      Comprehensive health check
// @Description  Performs detailed health checks on all system components including database, Elasticsearch, cache, event processor, and distributed lock service
// @Tags         health
// @Produce      json
// @Success      200 {object} HealthStatus "System is healthy"
// @Failure      503 {object} HealthStatus "System is unhealthy"
// @Router       /health [get]
func HealthCheckHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Initialize health status
		health := HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Version:   "1.0.0", // You can make this configurable
			Services:  make(map[string]ServiceHealth),
		}

		// Check Database
		dbStart := time.Now()
		var dbHealth ServiceHealth
		if err := db.Raw("SELECT 1").Error; err != nil {
			dbHealth = ServiceHealth{
				Status:  "unhealthy",
				Message: "Database connection failed: " + err.Error(),
			}
			health.Status = "unhealthy"
		} else {
			dbHealth = ServiceHealth{
				Status:       "healthy",
				ResponseTime: time.Since(dbStart).Milliseconds(),
			}
		}
		health.Services["database"] = dbHealth

		// Check Elasticsearch
		esStart := time.Now()
		var esHealth ServiceHealth
		es := config.ESClient
		if es == nil {
			esHealth = ServiceHealth{
				Status:  "unhealthy",
				Message: "Elasticsearch client not initialized",
			}
			health.Status = "unhealthy"
		} else {
			// Ping Elasticsearch
			req := esapi.PingRequest{}
			res, err := req.Do(context.Background(), es)
			if err != nil {
				esHealth = ServiceHealth{
					Status:  "unhealthy",
					Message: "Elasticsearch ping failed: " + err.Error(),
				}
				health.Status = "unhealthy"
			} else {
				defer res.Body.Close()
				if res.IsError() {
					esHealth = ServiceHealth{
						Status:  "unhealthy",
						Message: "Elasticsearch returned error: " + res.String(),
					}
					health.Status = "unhealthy"
				} else {
					esHealth = ServiceHealth{
						Status:       "healthy",
						ResponseTime: time.Since(esStart).Milliseconds(),
					}
				}
			}
		}
		health.Services["elasticsearch"] = esHealth

		// Check Cache
		cacheStart := time.Now()
		var cacheHealth ServiceHealth
		cache := services.GetCacheFactory()
		if cache == nil {
			cacheHealth = ServiceHealth{
				Status:  "unhealthy",
				Message: "Cache not initialized",
			}
			health.Status = "unhealthy"
		} else {
			// Test cache operations
			testKey := "health_check_test"
			testValue := "test_value"

			// Set a test value
			err := cache.Set(testKey, testValue, 10*time.Second)
			if err != nil {
				cacheHealth = ServiceHealth{
					Status:  "unhealthy",
					Message: "Cache set operation failed: " + err.Error(),
				}
				health.Status = "unhealthy"
			} else {
				// Get the test value
				_, exists, err := cache.Get(testKey)
				if err != nil {
					cacheHealth = ServiceHealth{
						Status:  "unhealthy",
						Message: "Cache get operation failed: " + err.Error(),
					}
					health.Status = "unhealthy"
				} else if !exists {
					cacheHealth = ServiceHealth{
						Status:  "unhealthy",
						Message: "Cache get operation returned no value",
					}
					health.Status = "unhealthy"
				} else {
					cacheHealth = ServiceHealth{
						Status:       "healthy",
						ResponseTime: time.Since(cacheStart).Milliseconds(),
					}
				}
			}
		}
		health.Services["cache"] = cacheHealth

		// Check Event Processor
		eventStart := time.Now()
		var eventHealth ServiceHealth
		eventProcessor := services.GetEventProcessor()
		if eventProcessor == nil {
			eventHealth = ServiceHealth{
				Status:  "unhealthy",
				Message: "Event processor not initialized",
			}
			health.Status = "unhealthy"
		} else {
			eventHealth = ServiceHealth{
				Status:       "healthy",
				ResponseTime: time.Since(eventStart).Milliseconds(),
			}
		}
		health.Services["event_processor"] = eventHealth

		// Check Distributed Lock Service
		lockStart := time.Now()
		var lockHealth ServiceHealth
		distributedLock := services.GetDistributedLock()
		if distributedLock == nil {
			lockHealth = ServiceHealth{
				Status:  "unhealthy",
				Message: "Distributed lock service not initialized",
			}
			health.Status = "unhealthy"
		} else {
			lockHealth = ServiceHealth{
				Status:       "healthy",
				ResponseTime: time.Since(lockStart).Milliseconds(),
			}
		}
		health.Services["distributed_lock"] = lockHealth

		// Set appropriate HTTP status code
		if health.Status == "healthy" {
			c.JSON(http.StatusOK, health)
		} else {
			c.JSON(http.StatusServiceUnavailable, health)
		}
	}
}

// SimpleHealthCheckHandler provides a lightweight health check
// @Summary      Simple health check
// @Description  Provides a lightweight health check for load balancers and basic monitoring
// @Tags         health
// @Produce      json
// @Success      200 {object} map[string]interface{} "Service is running"
// @Router       /health/simple [get]
func SimpleHealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
}
