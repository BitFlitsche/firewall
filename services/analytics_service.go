package services

import (
	"encoding/json"
	"firewall/models"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// AnalyticsService handles analytics aggregation and processing
type AnalyticsService struct {
	db             *gorm.DB
	trafficLogging *TrafficLoggingService
}

// TopDataItem represents a top data item with count
type TopDataItem struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// RelationshipInsight represents a relationship insight
type RelationshipInsight struct {
	Type  string            `json:"type"`
	Data  map[string]string `json:"data"`
	Count int64             `json:"count"`
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(db *gorm.DB, trafficLogging *TrafficLoggingService) *AnalyticsService {
	return &AnalyticsService{
		db:             db,
		trafficLogging: trafficLogging,
	}
}

// GenerateHourlyAggregation generates hourly analytics aggregation
func (as *AnalyticsService) GenerateHourlyAggregation() error {
	now := time.Now()
	hourStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// Check if aggregation already exists for this hour
	var existing models.AnalyticsAggregation
	if err := as.db.Where("aggregation_date = ? AND aggregation_type = ?", hourStart, "hourly").First(&existing).Error; err == nil {
		// Aggregation already exists, skip
		return nil
	}

	// Get traffic logs for the hour
	var logs []models.TrafficLog
	if err := as.db.Where("timestamp >= ? AND timestamp < ?", hourStart, hourStart.Add(time.Hour)).Find(&logs).Error; err != nil {
		return err
	}

	if len(logs) == 0 {
		return nil // No data to aggregate
	}

	// Calculate metrics
	aggregation := &models.AnalyticsAggregation{
		AggregationDate: hourStart,
		AggregationType: "hourly",
	}

	// Count results
	for _, log := range logs {
		aggregation.TotalRequests++
		switch log.FinalResult {
		case "allowed":
			aggregation.AllowedRequests++
		case "denied":
			aggregation.DeniedRequests++
		case "whitelisted":
			aggregation.WhitelistedRequests++
		}
	}

	// Calculate top data
	aggregation.TopIPs = as.calculateTopData(logs, "ip_address")
	aggregation.TopEmails = as.calculateTopData(logs, "email")
	aggregation.TopUserAgents = as.calculateTopData(logs, "user_agent")
	aggregation.TopUsernames = as.calculateTopData(logs, "username")
	aggregation.TopCountries = as.calculateTopData(logs, "country")
	aggregation.TopCharsets = as.calculateTopData(logs, "charset")

	// Calculate top relationships
	aggregation.TopRelationships = as.calculateTopRelationships(hourStart)

	// Calculate performance metrics
	var totalResponseTime int64
	var cacheHits int64
	for _, log := range logs {
		totalResponseTime += int64(log.ResponseTimeMs)
		if log.CacheHit {
			cacheHits++
		}
	}

	if aggregation.TotalRequests > 0 {
		aggregation.AvgResponseTimeMs = float64(totalResponseTime) / float64(aggregation.TotalRequests)
		aggregation.CacheHitRate = float64(cacheHits) / float64(aggregation.TotalRequests) * 100
	}

	// Save aggregation
	return as.db.Save(aggregation).Error
}

// calculateTopData calculates top data items for a specific field
func (as *AnalyticsService) calculateTopData(logs []models.TrafficLog, field string) string {
	counts := make(map[string]int64)

	for _, log := range logs {
		var value string
		switch field {
		case "ip_address":
			value = log.IPAddress
		case "email":
			value = log.Email
		case "user_agent":
			value = log.UserAgent
		case "username":
			value = log.Username
		case "country":
			value = log.Country
		case "charset":
			value = log.Charset
		}

		if value != "" {
			counts[value]++
		}
	}

	// Convert to sorted slice
	var items []TopDataItem
	for value, count := range counts {
		items = append(items, TopDataItem{Value: value, Count: count})
	}

	// Sort by count (descending) and take top 10
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	if len(items) > 10 {
		items = items[:10]
	}

	if jsonData, err := json.Marshal(items); err == nil {
		return string(jsonData)
	}

	return "[]"
}

// calculateTopRelationships calculates top relationships for a time period
func (as *AnalyticsService) calculateTopRelationships(hourStart time.Time) string {
	var relationships []models.DataRelationship
	if err := as.db.Where("timestamp >= ? AND timestamp < ?", hourStart, hourStart.Add(time.Hour)).Find(&relationships).Error; err != nil {
		return "[]"
	}

	// Group by relationship type and data
	relationshipCounts := make(map[string]int64)

	for _, rel := range relationships {
		key := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", rel.RelationshipType, rel.IPAddress, rel.Email, rel.UserAgent, rel.Username, rel.Country, rel.Charset)
		relationshipCounts[key] += int64(rel.Frequency)
	}

	// Convert to insights
	var insights []RelationshipInsight
	for key, count := range relationshipCounts {
		// Parse key to extract relationship type and data
		parts := strings.Split(key, ":")
		if len(parts) >= 7 {
			insight := RelationshipInsight{
				Type: parts[0],
				Data: map[string]string{
					"ip_address": parts[1],
					"email":      parts[2],
					"user_agent": parts[3],
					"username":   parts[4],
					"country":    parts[5],
					"charset":    parts[6],
				},
				Count: count,
			}
			insights = append(insights, insight)
		}
	}

	// Sort by count and take top 10
	sort.Slice(insights, func(i, j int) bool {
		return insights[i].Count > insights[j].Count
	})

	if len(insights) > 10 {
		insights = insights[:10]
	}

	if jsonData, err := json.Marshal(insights); err == nil {
		return string(jsonData)
	}

	return "[]"
}

// GetAnalyticsAggregations retrieves analytics aggregations
func (as *AnalyticsService) GetAnalyticsAggregations(aggregationType string, days int) ([]models.AnalyticsAggregation, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	var aggregations []models.AnalyticsAggregation
	if err := as.db.Where("aggregation_type = ? AND aggregation_date >= ?", aggregationType, startDate).
		Order("aggregation_date DESC").Find(&aggregations).Error; err != nil {
		return nil, err
	}

	return aggregations, nil
}

// GetTopDataByPeriod retrieves top data for a specific period
func (as *AnalyticsService) GetTopDataByPeriod(startTime, endTime time.Time, dataType string) ([]TopDataItem, error) {
	var logs []models.TrafficLog
	if err := as.db.Where("timestamp BETWEEN ? AND ?", startTime, endTime).Find(&logs).Error; err != nil {
		return nil, err
	}

	// Calculate top data
	topDataJSON := as.calculateTopData(logs, dataType)

	var items []TopDataItem
	if err := json.Unmarshal([]byte(topDataJSON), &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetRelationshipInsights retrieves relationship insights for a period
func (as *AnalyticsService) GetRelationshipInsights(startTime, endTime time.Time, limit int) ([]RelationshipInsight, error) {
	var relationships []models.DataRelationship
	if err := as.db.Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Order("frequency DESC").Limit(limit).Find(&relationships).Error; err != nil {
		return nil, err
	}

	// Group by relationship type and data
	relationshipCounts := make(map[string]int64)

	for _, rel := range relationships {
		key := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", rel.RelationshipType, rel.IPAddress, rel.Email, rel.UserAgent, rel.Username, rel.Country, rel.Charset)
		relationshipCounts[key] += int64(rel.Frequency)
	}

	// Convert to insights
	var insights []RelationshipInsight
	for key, count := range relationshipCounts {
		parts := strings.Split(key, ":")
		if len(parts) >= 7 {
			insight := RelationshipInsight{
				Type: parts[0],
				Data: map[string]string{
					"ip_address": parts[1],
					"email":      parts[2],
					"user_agent": parts[3],
					"username":   parts[4],
					"country":    parts[5],
					"charset":    parts[6],
				},
				Count: count,
			}
			insights = append(insights, insight)
		}
	}

	// Sort by count
	sort.Slice(insights, func(i, j int) bool {
		return insights[i].Count > insights[j].Count
	})

	return insights, nil
}

// GenerateDailyAggregation generates daily analytics aggregation
func (as *AnalyticsService) GenerateDailyAggregation() error {
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Check if aggregation already exists for this day
	var existing models.AnalyticsAggregation
	if err := as.db.Where("aggregation_date = ? AND aggregation_type = ?", dayStart, "daily").First(&existing).Error; err == nil {
		// Aggregation already exists, skip
		return nil
	}

	// Get traffic logs for the day
	var logs []models.TrafficLog
	if err := as.db.Where("timestamp >= ? AND timestamp < ?", dayStart, dayStart.Add(24*time.Hour)).Find(&logs).Error; err != nil {
		return err
	}

	if len(logs) == 0 {
		return nil // No data to aggregate
	}

	// Calculate metrics (similar to hourly but for daily)
	aggregation := &models.AnalyticsAggregation{
		AggregationDate: dayStart,
		AggregationType: "daily",
	}

	// Count results
	for _, log := range logs {
		aggregation.TotalRequests++
		switch log.FinalResult {
		case "allowed":
			aggregation.AllowedRequests++
		case "denied":
			aggregation.DeniedRequests++
		case "whitelisted":
			aggregation.WhitelistedRequests++
		}
	}

	// Calculate top data
	aggregation.TopIPs = as.calculateTopData(logs, "ip_address")
	aggregation.TopEmails = as.calculateTopData(logs, "email")
	aggregation.TopUserAgents = as.calculateTopData(logs, "user_agent")
	aggregation.TopUsernames = as.calculateTopData(logs, "username")
	aggregation.TopCountries = as.calculateTopData(logs, "country")
	aggregation.TopCharsets = as.calculateTopData(logs, "charset")

	// Calculate top relationships
	aggregation.TopRelationships = as.calculateTopRelationships(dayStart)

	// Calculate performance metrics
	var totalResponseTime int64
	var cacheHits int64
	for _, log := range logs {
		totalResponseTime += int64(log.ResponseTimeMs)
		if log.CacheHit {
			cacheHits++
		}
	}

	if aggregation.TotalRequests > 0 {
		aggregation.AvgResponseTimeMs = float64(totalResponseTime) / float64(aggregation.TotalRequests)
		aggregation.CacheHitRate = float64(cacheHits) / float64(aggregation.TotalRequests) * 100
	}

	// Save aggregation
	return as.db.Save(aggregation).Error
}

// RunScheduledAggregations runs scheduled analytics aggregations
func (as *AnalyticsService) RunScheduledAggregations() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Generate hourly aggregation
			if err := as.GenerateHourlyAggregation(); err != nil {
				log.Printf("Error generating hourly aggregation: %v", err)
			}

			// Generate daily aggregation at midnight
			now := time.Now()
			if now.Hour() == 0 && now.Minute() < 5 {
				if err := as.GenerateDailyAggregation(); err != nil {
					log.Printf("Error generating daily aggregation: %v", err)
				}
			}
		}
	}
}
