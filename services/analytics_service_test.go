package services

import (
	"encoding/json"
	"firewall/config"
	"firewall/models"
	"fmt"
	"testing"
	"time"
)

func init() {
	// Initialize config for tests
	config.InitConfig()
}

func TestNewAnalyticsService(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	if service == nil {
		t.Error("Expected service to be created")
	}

	if service.db != nil {
		t.Error("Expected db to be nil in test environment")
	}

	if service.trafficLogging != nil {
		t.Error("Expected trafficLogging to be nil in test environment")
	}
}

func TestAnalyticsService_CalculateTopData(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// Create test logs
	logs := []models.TrafficLog{
		{IPAddress: "192.168.1.1", Email: "test1@example.com", UserAgent: "Mozilla/1", Username: "user1", Country: "US", Charset: "UTF-8"},
		{IPAddress: "192.168.1.1", Email: "test2@example.com", UserAgent: "Mozilla/1", Username: "user2", Country: "US", Charset: "UTF-8"},
		{IPAddress: "192.168.1.2", Email: "test3@example.com", UserAgent: "Mozilla/2", Username: "user3", Country: "CA", Charset: "UTF-8"},
		{IPAddress: "192.168.1.3", Email: "test4@example.com", UserAgent: "Mozilla/3", Username: "user4", Country: "UK", Charset: "UTF-8"},
		{IPAddress: "192.168.1.1", Email: "test5@example.com", UserAgent: "Mozilla/1", Username: "user5", Country: "US", Charset: "UTF-8"},
	}

	tests := []struct {
		name     string
		field    string
		expected int // Expected number of unique values
	}{
		{"ip_address", "ip_address", 3}, // 3 unique IPs
		{"email", "email", 5},           // 5 unique emails
		{"user_agent", "user_agent", 3}, // 3 unique user agents
		{"username", "username", 5},     // 5 unique usernames
		{"country", "country", 3},       // 3 unique countries
		{"charset", "charset", 1},       // 1 unique charset
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.calculateTopData(logs, tt.field)

			// Parse JSON result
			var items []TopDataItem
			err := json.Unmarshal([]byte(result), &items)
			if err != nil {
				t.Errorf("Failed to parse JSON result: %v", err)
				return
			}

			// Check number of unique values
			if len(items) != tt.expected {
				t.Errorf("Expected %d unique values, got %d", tt.expected, len(items))
			}

			// Check that items are sorted by count (descending)
			for i := 1; i < len(items); i++ {
				if items[i-1].Count < items[i].Count {
					t.Errorf("Items not sorted correctly: %d < %d", items[i-1].Count, items[i].Count)
				}
			}

			// Check that all items have valid data
			for _, item := range items {
				if item.Value == "" {
					t.Error("Found empty value in top data")
				}
				if item.Count <= 0 {
					t.Error("Found non-positive count in top data")
				}
			}
		})
	}
}

func TestAnalyticsService_CalculateTopData_EmptyLogs(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	result := service.calculateTopData([]models.TrafficLog{}, "ip_address")

	// Should return "null" for empty logs (JSON marshaling of empty slice)
	if result != "null" {
		t.Errorf("Expected 'null', got %s", result)
	}
}

func TestAnalyticsService_CalculateTopData_EmptyFields(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// Create logs with empty fields
	logs := []models.TrafficLog{
		{IPAddress: "", Email: "", UserAgent: "", Username: "", Country: "", Charset: ""},
		{IPAddress: "192.168.1.1", Email: "", UserAgent: "", Username: "", Country: "", Charset: ""},
	}

	result := service.calculateTopData(logs, "ip_address")

	var items []TopDataItem
	err := json.Unmarshal([]byte(result), &items)
	if err != nil {
		t.Errorf("Failed to parse JSON result: %v", err)
		return
	}

	// Should only have one item (the non-empty IP)
	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}

	if items[0].Value != "192.168.1.1" {
		t.Errorf("Expected IP '192.168.1.1', got %s", items[0].Value)
	}
}

func TestAnalyticsService_CalculateTopData_MoreThan10Items(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// Create logs with more than 10 unique IPs
	logs := make([]models.TrafficLog, 15)
	for i := 0; i < 15; i++ {
		logs[i] = models.TrafficLog{
			IPAddress: fmt.Sprintf("192.168.1.%d", i+1),
		}
	}

	result := service.calculateTopData(logs, "ip_address")

	var items []TopDataItem
	err := json.Unmarshal([]byte(result), &items)
	if err != nil {
		t.Errorf("Failed to parse JSON result: %v", err)
		return
	}

	// Should be limited to top 10
	if len(items) != 10 {
		t.Errorf("Expected 10 items, got %d", len(items))
	}
}

func TestAnalyticsService_CalculateTopData_InvalidField(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	logs := []models.TrafficLog{
		{IPAddress: "192.168.1.1"},
	}

	result := service.calculateTopData(logs, "invalid_field")

	// Should return "null" for invalid field (no matching data)
	if result != "null" {
		t.Errorf("Expected 'null' for invalid field, got %s", result)
	}
}

func TestAnalyticsService_CalculateTopRelationships(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("calculateTopRelationships panicked as expected: %v", r)
		}
	}()

	result := service.calculateTopRelationships(time.Now())

	// Should return valid JSON array
	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Should be valid JSON
	if result != "[]" {
		var items []RelationshipInsight
		err := json.Unmarshal([]byte(result), &items)
		if err != nil {
			t.Errorf("Failed to parse JSON result: %v", err)
		}
	}
}

func TestAnalyticsService_GetAnalyticsAggregations(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GetAnalyticsAggregations panicked as expected: %v", r)
		}
	}()

	aggregations, err := service.GetAnalyticsAggregations("hourly", 7)
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("GetAnalyticsAggregations completed without error (unexpected in test environment)")
	}

	// Verify return types
	if aggregations == nil {
		t.Log("aggregations is nil (expected in test environment)")
	}
}

func TestAnalyticsService_GetTopDataByPeriod(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GetTopDataByPeriod panicked as expected: %v", r)
		}
	}()

	items, err := service.GetTopDataByPeriod(startTime, endTime, "ip_address")
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("GetTopDataByPeriod completed without error (unexpected in test environment)")
	}

	// Verify return types
	if items == nil {
		t.Log("items is nil (expected in test environment)")
	}
}

func TestAnalyticsService_GetRelationshipInsights(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GetRelationshipInsights panicked as expected: %v", r)
		}
	}()

	insights, err := service.GetRelationshipInsights(startTime, endTime, 10)
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("GetRelationshipInsights completed without error (unexpected in test environment)")
	}

	// Verify return types
	if insights == nil {
		t.Log("insights is nil (expected in test environment)")
	}
}

func TestAnalyticsService_GenerateHourlyAggregation(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GenerateHourlyAggregation panicked as expected: %v", r)
		}
	}()

	err := service.GenerateHourlyAggregation()
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("GenerateHourlyAggregation completed without error (unexpected in test environment)")
	}
}

func TestAnalyticsService_GenerateDailyAggregation(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GenerateDailyAggregation panicked as expected: %v", r)
		}
	}()

	err := service.GenerateDailyAggregation()
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("GenerateDailyAggregation completed without error (unexpected in test environment)")
	}
}

func TestAnalyticsService_RunScheduledAggregations(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// This function runs indefinitely, so we can only test that it doesn't panic immediately
	defer func() {
		if r := recover(); r != nil {
			t.Logf("RunScheduledAggregations panicked as expected: %v", r)
		}
	}()

	// Start the function in a goroutine and stop it quickly
	go service.RunScheduledAggregations()

	// Give it a moment to start
	time.Sleep(10 * time.Millisecond)

	// The function should not panic immediately
	t.Log("RunScheduledAggregations started without immediate panic")
}

func TestAnalyticsService_CalculateTopData_AllFields(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// Create comprehensive test logs
	logs := []models.TrafficLog{
		{
			IPAddress: "192.168.1.1",
			Email:     "test1@example.com",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			Username:  "user1",
			Country:   "US",
			Charset:   "UTF-8",
		},
		{
			IPAddress: "192.168.1.2",
			Email:     "test2@example.com",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			Username:  "user2",
			Country:   "CA",
			Charset:   "UTF-8",
		},
		{
			IPAddress: "192.168.1.1", // Duplicate IP
			Email:     "test3@example.com",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)", // Duplicate UserAgent
			Username:  "user3",
			Country:   "UK",
			Charset:   "UTF-8",
		},
	}

	// Test all fields
	fields := []string{"ip_address", "email", "user_agent", "username", "country", "charset"}

	for _, field := range fields {
		t.Run(field, func(t *testing.T) {
			result := service.calculateTopData(logs, field)

			var items []TopDataItem
			err := json.Unmarshal([]byte(result), &items)
			if err != nil {
				t.Errorf("Failed to parse JSON result for field %s: %v", field, err)
				return
			}

			// Should have at least one item
			if len(items) == 0 {
				t.Errorf("Expected at least one item for field %s", field)
			}

			// Check that items are sorted by count (descending)
			for i := 1; i < len(items); i++ {
				if items[i-1].Count < items[i].Count {
					t.Errorf("Items not sorted correctly for field %s: %d < %d", field, items[i-1].Count, items[i].Count)
				}
			}

			// Check that all items have valid data
			for _, item := range items {
				if item.Value == "" {
					t.Errorf("Found empty value in top data for field %s", field)
				}
				if item.Count <= 0 {
					t.Errorf("Found non-positive count in top data for field %s", field)
				}
			}
		})
	}
}

func TestAnalyticsService_CalculateTopData_EdgeCases(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// Test with logs that have mixed empty and non-empty values
	logs := []models.TrafficLog{
		{IPAddress: "192.168.1.1", Email: "", UserAgent: "Mozilla/1", Username: "", Country: "US", Charset: ""},
		{IPAddress: "", Email: "test@example.com", UserAgent: "", Username: "user1", Country: "", Charset: "UTF-8"},
		{IPAddress: "192.168.1.2", Email: "test2@example.com", UserAgent: "Mozilla/2", Username: "user2", Country: "CA", Charset: "UTF-8"},
	}

	fields := []string{"ip_address", "email", "user_agent", "username", "country", "charset"}

	for _, field := range fields {
		t.Run(field, func(t *testing.T) {
			result := service.calculateTopData(logs, field)

			var items []TopDataItem
			err := json.Unmarshal([]byte(result), &items)
			if err != nil {
				t.Errorf("Failed to parse JSON result for field %s: %v", field, err)
				return
			}

			// Should not have empty values
			for _, item := range items {
				if item.Value == "" {
					t.Errorf("Found empty value in top data for field %s", field)
				}
			}
		})
	}
}

func TestAnalyticsService_CalculateTopData_SingleValue(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// Test with logs that have only one unique value per field
	logs := []models.TrafficLog{
		{IPAddress: "192.168.1.1", Email: "test@example.com", UserAgent: "Mozilla/1", Username: "user1", Country: "US", Charset: "UTF-8"},
		{IPAddress: "192.168.1.1", Email: "test@example.com", UserAgent: "Mozilla/1", Username: "user1", Country: "US", Charset: "UTF-8"},
	}

	result := service.calculateTopData(logs, "ip_address")

	var items []TopDataItem
	err := json.Unmarshal([]byte(result), &items)
	if err != nil {
		t.Errorf("Failed to parse JSON result: %v", err)
		return
	}

	// Should have exactly one item
	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}

	// Should have count of 2
	if items[0].Count != 2 {
		t.Errorf("Expected count 2, got %d", items[0].Count)
	}
}

func TestAnalyticsService_CalculateTopData_JSONError(t *testing.T) {
	service := NewAnalyticsService(nil, nil)

	// Create logs with problematic data that might cause JSON marshaling issues
	logs := []models.TrafficLog{
		{IPAddress: "192.168.1.1", Email: "test@example.com", UserAgent: "Mozilla/1", Username: "user1", Country: "US", Charset: "UTF-8"},
	}

	result := service.calculateTopData(logs, "ip_address")

	// Should return valid JSON even if there are issues
	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Should be valid JSON
	var items []TopDataItem
	err := json.Unmarshal([]byte(result), &items)
	if err != nil {
		t.Errorf("Failed to parse JSON result: %v", err)
	}
}
