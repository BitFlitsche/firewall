package services

import (
	"encoding/json"
	"firewall/config"
	"firewall/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TrafficLoggingService handles traffic logging and analytics
type TrafficLoggingService struct {
	db      *gorm.DB
	enabled bool
}

// FilterRequest represents a filter request
type FilterRequest struct {
	IPAddress string `json:"ip_address"`
	Email     string `json:"email"`
	UserAgent string `json:"user_agent"`
	Username  string `json:"username"`
	Country   string `json:"country"`
	ASN       string `json:"asn"`
	Charset   string `json:"charset"`
	Content   string `json:"content"`
}

// TrafficFilterResult represents the result of a filter operation for traffic logging
type TrafficFilterResult struct {
	FinalResult   string                 `json:"final_result"`
	FilterResults map[string]interface{} `json:"filter_results"`
	ResponseTime  time.Duration          `json:"response_time"`
	CacheHit      bool                   `json:"cache_hit"`
}

// NewTrafficLoggingService creates a new traffic logging service
func NewTrafficLoggingService(db *gorm.DB) *TrafficLoggingService {
	enabled := false
	if config.AppConfig != nil && config.AppConfig.Logging.TrafficLogging {
		enabled = true
	}

	return &TrafficLoggingService{
		db:      db,
		enabled: enabled,
	}
}

// LogFilterRequest logs a filter request with its results
func (tls *TrafficLoggingService) LogFilterRequest(req FilterRequest, result TrafficFilterResult, metadata map[string]string) error {
	if !tls.enabled {
		return nil
	}

	// Create traffic log
	trafficLog := &models.TrafficLog{
		RequestID:      uuid.New().String(),
		IPAddress:      req.IPAddress,
		Email:          req.Email,
		UserAgent:      req.UserAgent,
		Username:       req.Username,
		Country:        req.Country,
		ASN:            req.ASN,
		Charset:        req.Charset,
		Content:        req.Content,
		FinalResult:    result.FinalResult,
		ResponseTimeMs: int(result.ResponseTime.Milliseconds()),
		CacheHit:       result.CacheHit,
		UserID:         metadata["user_id"],
		SessionID:      metadata["session_id"],
		ClientIP:       metadata["client_ip"],
		UserAgentRaw:   metadata["user_agent_raw"],
	}

	// Convert filter results to JSON
	if filterResultsJSON, err := json.Marshal(result.FilterResults); err == nil {
		trafficLog.FilterResults = string(filterResultsJSON)
	}

	// Save to database
	if err := tls.db.Create(trafficLog).Error; err != nil {
		log.Printf("Error logging traffic: %v", err)
		return err
	}

	// Update data relationships asynchronously
	go tls.updateDataRelationships(req, trafficLog.ID)

	return nil
}

// updateDataRelationships updates data relationships for analytics
func (tls *TrafficLoggingService) updateDataRelationships(req FilterRequest, logID uint) {
	relationships := tls.generateRelationships(req)

	for _, rel := range relationships {
		// Try to find existing relationship
		var existing models.DataRelationship
		err := tls.db.Where("relationship_type = ? AND ip_address = ? AND email = ? AND user_agent = ? AND username = ? AND country = ? AND charset = ?",
			rel.RelationshipType, rel.IPAddress, rel.Email, rel.UserAgent, rel.Username, rel.Country, rel.Charset).First(&existing).Error

		if err != nil {
			// Create new relationship
			tls.db.Create(&rel)
		} else {
			// Update existing relationship
			existing.Frequency++
			existing.LastSeen = time.Now()
			tls.db.Save(&existing)
		}
	}
}

// generateRelationships generates all possible relationships from a request
func (tls *TrafficLoggingService) generateRelationships(req FilterRequest) []models.DataRelationship {
	var relationships []models.DataRelationship

	// Generate all possible relationships
	relationshipTypes := []string{
		"ip_email", "ip_useragent", "ip_username", "ip_country", "ip_charset",
		"email_useragent", "email_username", "email_country", "email_charset",
		"useragent_username", "useragent_country", "useragent_charset",
		"username_country", "username_charset", "country_charset",
	}

	for _, relType := range relationshipTypes {
		if tls.isValidRelationship(relType, req) {
			relationships = append(relationships, models.DataRelationship{
				IPAddress:        req.IPAddress,
				Email:            req.Email,
				UserAgent:        req.UserAgent,
				Username:         req.Username,
				Country:          req.Country,
				Charset:          req.Charset,
				RelationshipType: relType,
				Frequency:        1,
				FirstSeen:        time.Now(),
				LastSeen:         time.Now(),
			})
		}
	}

	return relationships
}

// isValidRelationship checks if a relationship is valid based on the request data
func (tls *TrafficLoggingService) isValidRelationship(relType string, req FilterRequest) bool {
	switch relType {
	case "ip_email":
		return req.IPAddress != "" && req.Email != ""
	case "ip_useragent":
		return req.IPAddress != "" && req.UserAgent != ""
	case "ip_username":
		return req.IPAddress != "" && req.Username != ""
	case "ip_country":
		return req.IPAddress != "" && req.Country != ""
	case "ip_charset":
		return req.IPAddress != "" && req.Charset != ""
	case "email_useragent":
		return req.Email != "" && req.UserAgent != ""
	case "email_username":
		return req.Email != "" && req.Username != ""
	case "email_country":
		return req.Email != "" && req.Country != ""
	case "email_charset":
		return req.Email != "" && req.Charset != ""
	case "useragent_username":
		return req.UserAgent != "" && req.Username != ""
	case "useragent_country":
		return req.UserAgent != "" && req.Country != ""
	case "useragent_charset":
		return req.UserAgent != "" && req.Charset != ""
	case "username_country":
		return req.Username != "" && req.Country != ""
	case "username_charset":
		return req.Username != "" && req.Charset != ""
	case "country_charset":
		return req.Country != "" && req.Charset != ""
	default:
		return false
	}
}

// GetTrafficLogs retrieves traffic logs with filtering
func (tls *TrafficLoggingService) GetTrafficLogs(limit int, offset int, filters map[string]string) ([]models.TrafficLog, int64, error) {
	query := tls.db.Model(&models.TrafficLog{})

	// Apply filters
	if filters["ip_address"] != "" {
		query = query.Where("ip_address LIKE ?", "%"+filters["ip_address"]+"%")
	}
	if filters["email"] != "" {
		query = query.Where("email LIKE ?", "%"+filters["email"]+"%")
	}
	if filters["user_agent"] != "" {
		query = query.Where("user_agent LIKE ?", "%"+filters["user_agent"]+"%")
	}
	if filters["username"] != "" {
		query = query.Where("username LIKE ?", "%"+filters["username"]+"%")
	}
	if filters["country"] != "" {
		query = query.Where("country LIKE ?", "%"+filters["country"]+"%")
	}
	if filters["asn"] != "" {
		query = query.Where("asn LIKE ?", "%"+filters["asn"]+"%")
	}
	if filters["final_result"] != "" {
		query = query.Where("final_result = ?", filters["final_result"])
	}
	if filters["start_date"] != "" {
		query = query.Where("timestamp >= ?", filters["start_date"])
	}
	if filters["end_date"] != "" {
		query = query.Where("timestamp <= ?", filters["end_date"])
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderBy := filters["orderBy"]
	order := filters["order"]
	if orderBy == "" {
		orderBy = "timestamp"
	}
	if order == "" {
		order = "desc"
	}

	// Map frontend field names to database column names
	orderByMap := map[string]string{
		"timestamp":        "timestamp",
		"final_result":     "final_result",
		"ip_address":       "ip_address",
		"email":            "email",
		"user_agent":       "user_agent",
		"username":         "username",
		"country":          "country",
		"asn":              "asn",
		"response_time_ms": "response_time_ms",
		"cache_hit":        "cache_hit",
	}

	dbOrderBy := orderByMap[orderBy]
	if dbOrderBy == "" {
		dbOrderBy = "timestamp"
	}

	// Get paginated results with sorting
	var logs []models.TrafficLog
	if err := query.Order(dbOrderBy + " " + order).Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetDataRelationships retrieves data relationships with filtering
func (tls *TrafficLoggingService) GetDataRelationships(limit int, offset int, filters map[string]string) ([]models.DataRelationship, int64, error) {
	query := tls.db.Model(&models.DataRelationship{})

	// Apply filters
	if filters["relationship_type"] != "" {
		query = query.Where("relationship_type = ?", filters["relationship_type"])
	}
	if filters["ip_address"] != "" {
		query = query.Where("ip_address LIKE ?", "%"+filters["ip_address"]+"%")
	}
	if filters["email"] != "" {
		query = query.Where("email LIKE ?", "%"+filters["email"]+"%")
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	var relationships []models.DataRelationship
	if err := query.Order("frequency DESC").Limit(limit).Offset(offset).Find(&relationships).Error; err != nil {
		return nil, 0, err
	}

	return relationships, total, nil
}

// GetTrafficStats returns traffic statistics for a given period
func (tls *TrafficLoggingService) GetTrafficStats(startTime, endTime time.Time) (map[string]interface{}, error) {
	var stats struct {
		TotalRequests       int64   `json:"total_requests"`
		AllowedRequests     int64   `json:"allowed_requests"`
		DeniedRequests      int64   `json:"denied_requests"`
		WhitelistedRequests int64   `json:"whitelisted_requests"`
		AvgResponseTime     float64 `json:"avg_response_time"`
		CacheHitRate        float64 `json:"cache_hit_rate"`
	}

	// Get basic counts
	if err := tls.db.Model(&models.TrafficLog{}).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Select("COUNT(*) as total_requests, " +
			"SUM(CASE WHEN final_result = 'allowed' THEN 1 ELSE 0 END) as allowed_requests, " +
			"SUM(CASE WHEN final_result = 'denied' THEN 1 ELSE 0 END) as denied_requests, " +
			"SUM(CASE WHEN final_result = 'whitelisted' THEN 1 ELSE 0 END) as whitelisted_requests, " +
			"AVG(response_time_ms) as avg_response_time, " +
			"AVG(CASE WHEN cache_hit = 1 THEN 100 ELSE 0 END) as cache_hit_rate").
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_requests":       stats.TotalRequests,
		"allowed_requests":     stats.AllowedRequests,
		"denied_requests":      stats.DeniedRequests,
		"whitelisted_requests": stats.WhitelistedRequests,
		"avg_response_time_ms": stats.AvgResponseTime,
		"cache_hit_rate":       stats.CacheHitRate,
	}, nil
}

// CleanupOldLogs removes old traffic logs based on retention period
func (tls *TrafficLoggingService) CleanupOldLogs(retentionDays int) error {
	if retentionDays <= 0 {
		return nil
	}

	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	// Delete old traffic logs
	if err := tls.db.Where("timestamp < ?", cutoffDate).Delete(&models.TrafficLog{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup traffic logs: %v", err)
	}

	// Delete old data relationships
	if err := tls.db.Where("last_seen < ?", cutoffDate).Delete(&models.DataRelationship{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup data relationships: %v", err)
	}

	log.Printf("Cleaned up traffic logs older than %d days", retentionDays)
	return nil
}
