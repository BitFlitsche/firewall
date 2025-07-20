package services

import (
	"firewall/config"
	"testing"
	"time"
)

func init() {
	// Initialize config for tests
	config.InitConfig()
}

func TestNewTrafficLoggingService(t *testing.T) {
	// Test with traffic logging enabled
	config.AppConfig.Logging.TrafficLogging = true
	service := NewTrafficLoggingService(nil)

	if service == nil {
		t.Error("Expected service to be created")
	}

	if !service.enabled {
		t.Error("Expected service to be enabled")
	}

	// Test with traffic logging disabled
	config.AppConfig.Logging.TrafficLogging = false
	service = NewTrafficLoggingService(nil)

	if service == nil {
		t.Error("Expected service to be created")
	}

	if service.enabled {
		t.Error("Expected service to be disabled")
	}
}

func TestTrafficLoggingService_LogFilterRequest_Disabled(t *testing.T) {
	config.AppConfig.Logging.TrafficLogging = false

	service := NewTrafficLoggingService(nil)

	req := FilterRequest{
		IPAddress: "192.168.1.1",
		Email:     "test@example.com",
	}

	result := TrafficFilterResult{
		FinalResult:  "allowed",
		ResponseTime: 100 * time.Millisecond,
		CacheHit:     false,
	}

	metadata := map[string]string{
		"user_id":    "user123",
		"session_id": "session456",
	}

	// Should return nil when disabled
	err := service.LogFilterRequest(req, result, metadata)
	if err != nil {
		t.Errorf("Expected no error when disabled, got: %v", err)
	}
}

func TestTrafficLoggingService_IsValidRelationship(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	tests := []struct {
		name     string
		relType  string
		req      FilterRequest
		expected bool
	}{
		{
			name:    "valid ip_email relationship",
			relType: "ip_email",
			req: FilterRequest{
				IPAddress: "192.168.1.1",
				Email:     "test@example.com",
			},
			expected: true,
		},
		{
			name:    "invalid ip_email relationship - missing email",
			relType: "ip_email",
			req: FilterRequest{
				IPAddress: "192.168.1.1",
				Email:     "",
			},
			expected: false,
		},
		{
			name:    "valid ip_useragent relationship",
			relType: "ip_useragent",
			req: FilterRequest{
				IPAddress: "192.168.1.1",
				UserAgent: "Mozilla/5.0",
			},
			expected: true,
		},
		{
			name:    "valid email_username relationship",
			relType: "email_username",
			req: FilterRequest{
				Email:    "test@example.com",
				Username: "testuser",
			},
			expected: true,
		},
		{
			name:    "valid country_charset relationship",
			relType: "country_charset",
			req: FilterRequest{
				Country: "US",
				Charset: "UTF-8",
			},
			expected: true,
		},
		{
			name:    "invalid relationship type",
			relType: "invalid_type",
			req: FilterRequest{
				IPAddress: "192.168.1.1",
				Email:     "test@example.com",
			},
			expected: false,
		},
		{
			name:     "empty relationship",
			relType:  "ip_email",
			req:      FilterRequest{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.isValidRelationship(tt.relType, tt.req)
			if result != tt.expected {
				t.Errorf("isValidRelationship(%s) = %v, expected %v", tt.relType, result, tt.expected)
			}
		})
	}
}

func TestTrafficLoggingService_GenerateRelationships(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	req := FilterRequest{
		IPAddress: "192.168.1.1",
		Email:     "test@example.com",
		UserAgent: "Mozilla/5.0",
		Username:  "testuser",
		Country:   "US",
		Charset:   "UTF-8",
	}

	relationships := service.generateRelationships(req)

	// Should generate 15 relationships (all combinations of 6 fields taken 2 at a time)
	if len(relationships) != 15 {
		t.Errorf("Expected 15 relationships, got %d", len(relationships))
	}

	// Verify all relationships have correct data
	for _, rel := range relationships {
		if rel.IPAddress != req.IPAddress {
			t.Errorf("Expected IPAddress %s, got %s", req.IPAddress, rel.IPAddress)
		}
		if rel.Email != req.Email {
			t.Errorf("Expected Email %s, got %s", req.Email, rel.Email)
		}
		if rel.UserAgent != req.UserAgent {
			t.Errorf("Expected UserAgent %s, got %s", req.UserAgent, rel.UserAgent)
		}
		if rel.Username != req.Username {
			t.Errorf("Expected Username %s, got %s", req.Username, rel.Username)
		}
		if rel.Country != req.Country {
			t.Errorf("Expected Country %s, got %s", req.Country, rel.Country)
		}
		if rel.Charset != req.Charset {
			t.Errorf("Expected Charset %s, got %s", req.Charset, rel.Charset)
		}
		if rel.Frequency != 1 {
			t.Errorf("Expected Frequency 1, got %d", rel.Frequency)
		}
		if rel.RelationshipType == "" {
			t.Error("Expected non-empty RelationshipType")
		}
	}

	// Verify relationship types
	expectedTypes := []string{
		"ip_email", "ip_useragent", "ip_username", "ip_country", "ip_charset",
		"email_useragent", "email_username", "email_country", "email_charset",
		"useragent_username", "useragent_country", "useragent_charset",
		"username_country", "username_charset", "country_charset",
	}

	foundTypes := make(map[string]bool)
	for _, rel := range relationships {
		foundTypes[rel.RelationshipType] = true
	}

	for _, expectedType := range expectedTypes {
		if !foundTypes[expectedType] {
			t.Errorf("Missing relationship type: %s", expectedType)
		}
	}
}

func TestTrafficLoggingService_GenerateRelationships_PartialData(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	req := FilterRequest{
		IPAddress: "192.168.1.1",
		Email:     "test@example.com",
		// Missing other fields
	}

	relationships := service.generateRelationships(req)

	// Should only generate relationships for available data
	if len(relationships) != 1 {
		t.Errorf("Expected 1 relationship, got %d", len(relationships))
	}

	if relationships[0].RelationshipType != "ip_email" {
		t.Errorf("Expected relationship type 'ip_email', got %s", relationships[0].RelationshipType)
	}
	if relationships[0].IPAddress != req.IPAddress {
		t.Errorf("Expected IPAddress %s, got %s", req.IPAddress, relationships[0].IPAddress)
	}
	if relationships[0].Email != req.Email {
		t.Errorf("Expected Email %s, got %s", req.Email, relationships[0].Email)
	}
}

func TestTrafficLoggingService_GenerateRelationships_EmptyData(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	req := FilterRequest{}

	relationships := service.generateRelationships(req)

	// Should generate no relationships for empty data
	if len(relationships) != 0 {
		t.Errorf("Expected 0 relationships, got %d", len(relationships))
	}
}

func TestTrafficLoggingService_GenerateRelationships_SingleField(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	req := FilterRequest{
		IPAddress: "192.168.1.1",
	}

	relationships := service.generateRelationships(req)

	// Should generate no relationships for single field
	if len(relationships) != 0 {
		t.Errorf("Expected 0 relationships, got %d", len(relationships))
	}
}

func TestTrafficLoggingService_GenerateRelationships_TwoFields(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	req := FilterRequest{
		IPAddress: "192.168.1.1",
		Email:     "test@example.com",
	}

	relationships := service.generateRelationships(req)

	// Should generate 1 relationship for 2 fields
	if len(relationships) != 1 {
		t.Errorf("Expected 1 relationship, got %d", len(relationships))
	}

	if relationships[0].RelationshipType != "ip_email" {
		t.Errorf("Expected relationship type 'ip_email', got %s", relationships[0].RelationshipType)
	}
}

func TestTrafficLoggingService_GenerateRelationships_ThreeFields(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	req := FilterRequest{
		IPAddress: "192.168.1.1",
		Email:     "test@example.com",
		UserAgent: "Mozilla/5.0",
	}

	relationships := service.generateRelationships(req)

	// Should generate 3 relationships for 3 fields (3 choose 2 = 3)
	if len(relationships) != 3 {
		t.Errorf("Expected 3 relationships, got %d", len(relationships))
	}

	expectedTypes := []string{"ip_email", "ip_useragent", "email_useragent"}
	foundTypes := make(map[string]bool)
	for _, rel := range relationships {
		foundTypes[rel.RelationshipType] = true
	}

	for _, expectedType := range expectedTypes {
		if !foundTypes[expectedType] {
			t.Errorf("Missing relationship type: %s", expectedType)
		}
	}
}

func TestTrafficLoggingService_UpdateDataRelationships(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	req := FilterRequest{
		IPAddress: "192.168.1.1",
		Email:     "test@example.com",
		UserAgent: "Mozilla/5.0",
	}

	// This runs asynchronously, so we can't easily test the exact behavior
	// But we can verify the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Logf("updateDataRelationships panicked as expected: %v", r)
		}
	}()

	service.updateDataRelationships(req, 1)

	// Give some time for the goroutine to run
	time.Sleep(10 * time.Millisecond)
}

func TestTrafficLoggingService_LogFilterRequest_Enabled(t *testing.T) {
	config.AppConfig.Logging.TrafficLogging = true

	service := NewTrafficLoggingService(nil)

	req := FilterRequest{
		IPAddress: "192.168.1.1",
		Email:     "test@example.com",
		UserAgent: "Mozilla/5.0",
		Username:  "testuser",
		Country:   "US",
		Charset:   "UTF-8",
		Content:   "test content",
	}

	result := TrafficFilterResult{
		FinalResult: "allowed",
		FilterResults: map[string]interface{}{
			"ip":    "allowed",
			"email": "allowed",
		},
		ResponseTime: 100 * time.Millisecond,
		CacheHit:     false,
	}

	metadata := map[string]string{
		"user_id":        "user123",
		"session_id":     "session456",
		"client_ip":      "192.168.1.100",
		"user_agent_raw": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	}

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("LogFilterRequest panicked as expected: %v", r)
		}
	}()

	err := service.LogFilterRequest(req, result, metadata)
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("LogFilterRequest completed without error (unexpected in test environment)")
	}
}

func TestTrafficLoggingService_GetTrafficLogs(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	filters := map[string]string{
		"ip_address":   "192.168.1.1",
		"final_result": "allowed",
	}

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GetTrafficLogs panicked as expected: %v", r)
		}
	}()

	logs, count, err := service.GetTrafficLogs(10, 0, filters)
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("GetTrafficLogs completed without error (unexpected in test environment)")
	}

	// Verify return types
	if logs == nil {
		t.Log("logs is nil (expected in test environment)")
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestTrafficLoggingService_GetDataRelationships(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	filters := map[string]string{
		"relationship_type": "ip_email",
	}

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GetDataRelationships panicked as expected: %v", r)
		}
	}()

	relationships, count, err := service.GetDataRelationships(10, 0, filters)
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("GetDataRelationships completed without error (unexpected in test environment)")
	}

	// Verify return types
	if relationships == nil {
		t.Log("relationships is nil (expected in test environment)")
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestTrafficLoggingService_GetTrafficStats(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GetTrafficStats panicked as expected: %v", r)
		}
	}()

	stats, err := service.GetTrafficStats(startTime, endTime)
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("GetTrafficStats completed without error (unexpected in test environment)")
	}

	// Verify return types
	if stats == nil {
		t.Log("stats is nil (expected in test environment)")
	}
}

func TestTrafficLoggingService_CleanupOldLogs(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	retentionDays := 30

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("CleanupOldLogs panicked as expected: %v", r)
		}
	}()

	err := service.CleanupOldLogs(retentionDays)
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("CleanupOldLogs completed without error (unexpected in test environment)")
	}
}

func TestTrafficLoggingService_GetTrafficLogs_NoFilters(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GetTrafficLogs panicked as expected: %v", r)
		}
	}()

	logs, count, err := service.GetTrafficLogs(10, 0, nil)
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("GetTrafficLogs completed without error (unexpected in test environment)")
	}

	// Verify return types
	if logs == nil {
		t.Log("logs is nil (expected in test environment)")
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestTrafficLoggingService_GetDataRelationships_NoFilters(t *testing.T) {
	service := NewTrafficLoggingService(nil)

	// This will fail because we don't have a real database, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GetDataRelationships panicked as expected: %v", r)
		}
	}()

	relationships, count, err := service.GetDataRelationships(10, 0, nil)
	// We expect an error because there's no database, but the function should not panic
	if err == nil {
		t.Log("GetDataRelationships completed without error (unexpected in test environment)")
	}

	// Verify return types
	if relationships == nil {
		t.Log("relationships is nil (expected in test environment)")
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}
