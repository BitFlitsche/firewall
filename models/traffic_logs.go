package models

import (
	"time"
)

// TrafficLog represents a logged filter request
type TrafficLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Timestamp time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`
	RequestID string    `json:"request_id" gorm:"size:36;not null"`

	// Request data
	IPAddress string `json:"ip_address" gorm:"size:45"`
	Email     string `json:"email" gorm:"size:255"`
	UserAgent string `json:"user_agent" gorm:"type:text"`
	Username  string `json:"username" gorm:"size:255"`
	Country   string `json:"country" gorm:"size:10"`
	Charset   string `json:"charset" gorm:"size:50"`
	Content   string `json:"content" gorm:"type:text"`

	// Filter results
	FinalResult   string `json:"final_result" gorm:"type:enum('allowed','denied','whitelisted');not null"`
	FilterResults string `json:"filter_results" gorm:"type:json"`

	// Performance metrics
	ResponseTimeMs int  `json:"response_time_ms"`
	CacheHit       bool `json:"cache_hit" gorm:"default:false"`

	// Metadata
	UserID       string `json:"user_id" gorm:"size:255"`
	SessionID    string `json:"session_id" gorm:"size:255"`
	ClientIP     string `json:"client_ip" gorm:"size:45"`
	UserAgentRaw string `json:"user_agent_raw" gorm:"type:text"`
}

// DataRelationship represents relationships between different data types
type DataRelationship struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Timestamp time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationship data
	IPAddress string `json:"ip_address" gorm:"size:45"`
	Email     string `json:"email" gorm:"size:255"`
	UserAgent string `json:"user_agent" gorm:"type:text"`
	Username  string `json:"username" gorm:"size:255"`
	Country   string `json:"country" gorm:"size:10"`
	Charset   string `json:"charset" gorm:"size:50"`

	// Relationship metadata
	RelationshipType string    `json:"relationship_type" gorm:"type:enum('ip_email','ip_useragent','ip_username','ip_country','ip_charset','email_useragent','email_username','email_country','email_charset','useragent_username','useragent_country','useragent_charset','username_country','username_charset','country_charset');not null"`
	Frequency        int       `json:"frequency" gorm:"default:1"`
	FirstSeen        time.Time `json:"first_seen" gorm:"default:CURRENT_TIMESTAMP"`
	LastSeen         time.Time `json:"last_seen" gorm:"default:CURRENT_TIMESTAMP"`
}

// AnalyticsAggregation represents aggregated analytics data
type AnalyticsAggregation struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	AggregationDate time.Time `json:"aggregation_date" gorm:"not null"`
	AggregationType string    `json:"aggregation_type" gorm:"type:enum('daily','hourly','minute');not null"`

	// Metrics
	TotalRequests       int64 `json:"total_requests" gorm:"default:0"`
	AllowedRequests     int64 `json:"allowed_requests" gorm:"default:0"`
	DeniedRequests      int64 `json:"denied_requests" gorm:"default:0"`
	WhitelistedRequests int64 `json:"whitelisted_requests" gorm:"default:0"`

	// Top data
	TopIPs        string `json:"top_ips" gorm:"type:json"`
	TopEmails     string `json:"top_emails" gorm:"type:json"`
	TopUserAgents string `json:"top_useragents" gorm:"type:json"`
	TopUsernames  string `json:"top_usernames" gorm:"type:json"`
	TopCountries  string `json:"top_countries" gorm:"type:json"`
	TopCharsets   string `json:"top_charsets" gorm:"type:json"`

	// Relationship insights
	TopRelationships string `json:"top_relationships" gorm:"type:json"`

	// Performance metrics
	AvgResponseTimeMs float64 `json:"avg_response_time_ms"`
	CacheHitRate      float64 `json:"cache_hit_rate"`
}
