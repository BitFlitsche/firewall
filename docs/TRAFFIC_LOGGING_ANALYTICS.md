# Traffic Logging and Analytics System

## Overview

This document outlines the design for a comprehensive traffic logging and analytics system that captures filter requests and provides insights into data relationships and patterns.

## Goals

- **Log all filter requests** with complete context
- **Track data relationships** (IP ↔ Email ↔ UserAgent ↔ Username ↔ Country ↔ Charset)
- **Provide analytics dashboard** for pattern visualization
- **Enable correlation analysis** for security insights
- **Support real-time monitoring** of traffic patterns

## Architecture Design

### 1. Traffic Logging System

#### Database Schema
```sql
-- Traffic logs table
CREATE TABLE traffic_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    request_id VARCHAR(36) NOT NULL,  -- UUID for request tracking
    
    -- Request data
    ip_address VARCHAR(45),
    email VARCHAR(255),
    user_agent TEXT,
    username VARCHAR(255),
    country VARCHAR(10),
    charset VARCHAR(50),
    content TEXT,
    
    -- Filter results
    final_result ENUM('allowed', 'denied', 'whitelisted') NOT NULL,
    filter_results JSON,  -- Detailed results from each filter
    
    -- Performance metrics
    response_time_ms INT,
    cache_hit BOOLEAN DEFAULT FALSE,
    
    -- Metadata
    user_id VARCHAR(255),  -- If authenticated
    session_id VARCHAR(255),
    client_ip VARCHAR(45),
    user_agent_raw TEXT,
    
    -- Indexes for performance
    INDEX idx_timestamp (timestamp),
    INDEX idx_ip_address (ip_address),
    INDEX idx_email (email),
    INDEX idx_final_result (final_result),
    INDEX idx_request_id (request_id)
);

-- Data relationships table for analytics
CREATE TABLE data_relationships (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Relationship data
    ip_address VARCHAR(45),
    email VARCHAR(255),
    user_agent TEXT,
    username VARCHAR(255),
    country VARCHAR(10),
    charset VARCHAR(50),
    
    -- Relationship metadata
    relationship_type ENUM('ip_email', 'ip_useragent', 'ip_username', 'ip_country', 'ip_charset', 'email_useragent', 'email_username', 'email_country', 'email_charset', 'useragent_username', 'useragent_country', 'useragent_charset', 'username_country', 'username_charset', 'country_charset') NOT NULL,
    frequency INT DEFAULT 1,
    first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_relationship_type (relationship_type),
    INDEX idx_ip_address (ip_address),
    INDEX idx_email (email),
    INDEX idx_timestamp (timestamp)
);

-- Analytics aggregations table
CREATE TABLE analytics_aggregations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    aggregation_date DATE NOT NULL,
    aggregation_type ENUM('daily', 'hourly', 'minute') NOT NULL,
    
    -- Metrics
    total_requests BIGINT DEFAULT 0,
    allowed_requests BIGINT DEFAULT 0,
    denied_requests BIGINT DEFAULT 0,
    whitelisted_requests BIGINT DEFAULT 0,
    
    -- Top data
    top_ips JSON,  -- [{"ip": "192.168.1.1", "count": 150}, ...]
    top_emails JSON,
    top_useragents JSON,
    top_usernames JSON,
    top_countries JSON,
    top_charsets JSON,
    
    -- Relationship insights
    top_relationships JSON,  -- [{"type": "ip_email", "data": {...}, "count": 50}, ...]
    
    -- Performance metrics
    avg_response_time_ms DECIMAL(10,2),
    cache_hit_rate DECIMAL(5,2),
    
    -- Indexes
    UNIQUE KEY unique_aggregation (aggregation_date, aggregation_type)
);
```

#### Go Models
```go
// models/traffic_logs.go
type TrafficLog struct {
    ID            uint      `json:"id" gorm:"primaryKey"`
    Timestamp     time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`
    RequestID     string    `json:"request_id" gorm:"size:36;not null"`
    
    // Request data
    IPAddress     string `json:"ip_address" gorm:"size:45"`
    Email         string `json:"email" gorm:"size:255"`
    UserAgent     string `json:"user_agent" gorm:"type:text"`
    Username      string `json:"username" gorm:"size:255"`
    Country       string `json:"country" gorm:"size:10"`
    Charset       string `json:"charset" gorm:"size:50"`
    Content       string `json:"content" gorm:"type:text"`
    
    // Filter results
    FinalResult   string `json:"final_result" gorm:"type:enum('allowed','denied','whitelisted');not null"`
    FilterResults string `json:"filter_results" gorm:"type:json"`
    
    // Performance metrics
    ResponseTimeMs int  `json:"response_time_ms"`
    CacheHit       bool `json:"cache_hit" gorm:"default:false"`
    
    // Metadata
    UserID        string `json:"user_id" gorm:"size:255"`
    SessionID     string `json:"session_id" gorm:"size:255"`
    ClientIP      string `json:"client_ip" gorm:"size:45"`
    UserAgentRaw  string `json:"user_agent_raw" gorm:"type:text"`
}

type DataRelationship struct {
    ID              uint      `json:"id" gorm:"primaryKey"`
    Timestamp       time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`
    
    // Relationship data
    IPAddress       string `json:"ip_address" gorm:"size:45"`
    Email           string `json:"email" gorm:"size:255"`
    UserAgent       string `json:"user_agent" gorm:"type:text"`
    Username        string `json:"username" gorm:"size:255"`
    Country         string `json:"country" gorm:"size:10"`
    Charset         string `json:"charset" gorm:"size:50"`
    
    // Relationship metadata
    RelationshipType string    `json:"relationship_type" gorm:"type:enum('ip_email','ip_useragent','ip_username','ip_country','ip_charset','email_useragent','email_username','email_country','email_charset','useragent_username','useragent_country','useragent_charset','username_country','username_charset','country_charset');not null"`
    Frequency       int       `json:"frequency" gorm:"default:1"`
    FirstSeen       time.Time `json:"first_seen" gorm:"default:CURRENT_TIMESTAMP"`
    LastSeen        time.Time `json:"last_seen" gorm:"default:CURRENT_TIMESTAMP"`
}

type AnalyticsAggregation struct {
    ID              uint      `json:"id" gorm:"primaryKey"`
    AggregationDate time.Time `json:"aggregation_date" gorm:"not null"`
    AggregationType string    `json:"aggregation_type" gorm:"type:enum('daily','hourly','minute');not null"`
    
    // Metrics
    TotalRequests     int64   `json:"total_requests" gorm:"default:0"`
    AllowedRequests   int64   `json:"allowed_requests" gorm:"default:0"`
    DeniedRequests    int64   `json:"denied_requests" gorm:"default:0"`
    WhitelistedRequests int64 `json:"whitelisted_requests" gorm:"default:0"`
    
    // Top data
    TopIPs           string `json:"top_ips" gorm:"type:json"`
    TopEmails        string `json:"top_emails" gorm:"type:json"`
    TopUserAgents    string `json:"top_useragents" gorm:"type:json"`
    TopUsernames     string `json:"top_usernames" gorm:"type:json"`
    TopCountries     string `json:"top_countries" gorm:"type:json"`
    TopCharsets      string `json:"top_charsets" gorm:"type:json"`
    
    // Relationship insights
    TopRelationships string `json:"top_relationships" gorm:"type:json"`
    
    // Performance metrics
    AvgResponseTimeMs float64 `json:"avg_response_time_ms"`
    CacheHitRate      float64 `json:"cache_hit_rate"`
}
```

### 2. Logging Service

#### Traffic Logging Service
```go
// services/traffic_logging.go
package services

import (
    "encoding/json"
    "firewall/config"
    "firewall/models"
    "fmt"
    "log"
    "time"
    "github.com/google/uuid"
)

type TrafficLoggingService struct {
    db *gorm.DB
    enabled bool
}

type FilterRequest struct {
    IPAddress string `json:"ip_address"`
    Email     string `json:"email"`
    UserAgent string `json:"user_agent"`
    Username  string `json:"username"`
    Country   string `json:"country"`
    Charset   string `json:"charset"`
    Content   string `json:"content"`
}

type FilterResult struct {
    FinalResult   string                 `json:"final_result"`
    FilterResults map[string]interface{} `json:"filter_results"`
    ResponseTime  time.Duration          `json:"response_time"`
    CacheHit      bool                   `json:"cache_hit"`
}

func NewTrafficLoggingService(db *gorm.DB) *TrafficLoggingService {
    return &TrafficLoggingService{
        db: db,
        enabled: config.AppConfig.Logging.TrafficLogging,
    }
}

func (tls *TrafficLoggingService) LogFilterRequest(req FilterRequest, result FilterResult, metadata map[string]string) error {
    if !tls.enabled {
        return nil
    }

    // Create traffic log
    trafficLog := &models.TrafficLog{
        RequestID:     uuid.New().String(),
        IPAddress:     req.IPAddress,
        Email:         req.Email,
        UserAgent:     req.UserAgent,
        Username:      req.Username,
        Country:       req.Country,
        Charset:       req.Charset,
        Content:       req.Content,
        FinalResult:   result.FinalResult,
        ResponseTimeMs: int(result.ResponseTime.Milliseconds()),
        CacheHit:      result.CacheHit,
        UserID:        metadata["user_id"],
        SessionID:     metadata["session_id"],
        ClientIP:      metadata["client_ip"],
        UserAgentRaw:  metadata["user_agent_raw"],
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

    // Update data relationships
    go tls.updateDataRelationships(req, trafficLog.ID)

    return nil
}

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
                IPAddress:       req.IPAddress,
                Email:           req.Email,
                UserAgent:       req.UserAgent,
                Username:        req.Username,
                Country:         req.Country,
                Charset:         req.Charset,
                RelationshipType: relType,
                Frequency:       1,
                FirstSeen:       time.Now(),
                LastSeen:        time.Now(),
            })
        }
    }

    return relationships
}

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
```

### 3. Analytics Service

#### Analytics Aggregation Service
```go
// services/analytics_service.go
package services

import (
    "encoding/json"
    "firewall/config"
    "firewall/models"
    "log"
    "time"
)

type AnalyticsService struct {
    db *gorm.DB
    trafficLogging *TrafficLoggingService
}

type TopDataItem struct {
    Value string `json:"value"`
    Count int64  `json:"count"`
}

type RelationshipInsight struct {
    Type string                 `json:"type"`
    Data map[string]string     `json:"data"`
    Count int64                `json:"count"`
}

func NewAnalyticsService(db *gorm.DB, trafficLogging *TrafficLoggingService) *AnalyticsService {
    return &AnalyticsService{
        db: db,
        trafficLogging: trafficLogging,
    }
}

func (as *AnalyticsService) GenerateHourlyAggregation() error {
    now := time.Now()
    hourStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
    
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

func (as *AnalyticsService) calculateTopRelationships(hourStart time.Time) string {
    var relationships []models.DataRelationship
    if err := as.db.Where("timestamp >= ? AND timestamp < ?", hourStart, hourStart.Add(time.Hour)).Find(&relationships).Error; err != nil {
        return "[]"
    }

    // Group by relationship type and data
    relationshipCounts := make(map[string]int64)
    
    for _, rel := range relationships {
        key := fmt.Sprintf("%s:%s:%s:%s:%s:%s", rel.RelationshipType, rel.IPAddress, rel.Email, rel.UserAgent, rel.Username, rel.Country)
        relationshipCounts[key] += rel.Frequency
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
                    "email": parts[2],
                    "user_agent": parts[3],
                    "username": parts[4],
                    "country": parts[5],
                    "charset": parts[6],
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
```

### 4. API Endpoints

#### Analytics API Controllers
```go
// controllers/analytics.go
package controllers

import (
    "firewall/models"
    "firewall/services"
    "net/http"
    "strconv"
    "time"
    "github.com/gin-gonic/gin"
)

func GetTrafficAnalytics(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get query parameters
        period := c.DefaultQuery("period", "24h")
        limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
        
        var startTime time.Time
        switch period {
        case "1h":
            startTime = time.Now().Add(-1 * time.Hour)
        case "24h":
            startTime = time.Now().Add(-24 * time.Hour)
        case "7d":
            startTime = time.Now().Add(-7 * 24 * time.Hour)
        case "30d":
            startTime = time.Now().Add(-30 * 24 * time.Hour)
        default:
            startTime = time.Now().Add(-24 * time.Hour)
        }

        // Get traffic logs
        var logs []models.TrafficLog
        if err := db.Where("timestamp >= ?", startTime).Limit(limit).Order("timestamp DESC").Find(&logs).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traffic logs"})
            return
        }

        // Calculate summary statistics
        var total, allowed, denied, whitelisted int64
        var totalResponseTime int64
        var cacheHits int64

        for _, log := range logs {
            total++
            switch log.FinalResult {
            case "allowed":
                allowed++
            case "denied":
                denied++
            case "whitelisted":
                whitelisted++
            }
            totalResponseTime += int64(log.ResponseTimeMs)
            if log.CacheHit {
                cacheHits++
            }
        }

        avgResponseTime := float64(0)
        cacheHitRate := float64(0)
        if total > 0 {
            avgResponseTime = float64(totalResponseTime) / float64(total)
            cacheHitRate = float64(cacheHits) / float64(total) * 100
        }

        c.JSON(http.StatusOK, gin.H{
            "period": period,
            "total_requests": total,
            "allowed_requests": allowed,
            "denied_requests": denied,
            "whitelisted_requests": whitelisted,
            "avg_response_time_ms": avgResponseTime,
            "cache_hit_rate": cacheHitRate,
            "logs": logs,
        })
    }
}

func GetDataRelationships(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        relationshipType := c.Query("type")
        limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
        
        query := db.Model(&models.DataRelationship{})
        if relationshipType != "" {
            query = query.Where("relationship_type = ?", relationshipType)
        }
        
        var relationships []models.DataRelationship
        if err := query.Order("frequency DESC").Limit(limit).Find(&relationships).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch relationships"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "relationships": relationships,
            "total": len(relationships),
        })
    }
}

func GetAnalyticsAggregations(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        aggregationType := c.DefaultQuery("type", "hourly")
        days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
        
        startDate := time.Now().AddDate(0, 0, -days)
        
        var aggregations []models.AnalyticsAggregation
        if err := db.Where("aggregation_type = ? AND aggregation_date >= ?", aggregationType, startDate).
            Order("aggregation_date DESC").Find(&aggregations).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch aggregations"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "aggregations": aggregations,
            "type": aggregationType,
            "days": days,
        })
    }
}
```

### 5. Frontend Analytics Dashboard

#### React Analytics Components
```jsx
// firewall-app/src/components/AnalyticsDashboard.js
import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
    LineChart, Line, BarChart, Bar, PieChart, Pie,
    XAxis, YAxis, CartesianGrid, Tooltip, Legend,
    ResponsiveContainer
} from 'recharts';

const AnalyticsDashboard = () => {
    const [analytics, setAnalytics] = useState(null);
    const [relationships, setRelationships] = useState([]);
    const [period, setPeriod] = useState('24h');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchAnalytics();
        fetchRelationships();
    }, [period]);

    const fetchAnalytics = async () => {
        try {
            const response = await axios.get(`/api/analytics/traffic?period=${period}`);
            setAnalytics(response.data);
        } catch (error) {
            console.error('Error fetching analytics:', error);
        } finally {
            setLoading(false);
        }
    };

    const fetchRelationships = async () => {
        try {
            const response = await axios.get('/api/analytics/relationships?limit=20');
            setRelationships(response.data.relationships);
        } catch (error) {
            console.error('Error fetching relationships:', error);
        }
    };

    if (loading) {
        return <div>Loading analytics...</div>;
    }

    return (
        <div className="analytics-dashboard">
            <div className="dashboard-header">
                <h1>Traffic Analytics Dashboard</h1>
                <select value={period} onChange={(e) => setPeriod(e.target.value)}>
                    <option value="1h">Last Hour</option>
                    <option value="24h">Last 24 Hours</option>
                    <option value="7d">Last 7 Days</option>
                    <option value="30d">Last 30 Days</option>
                </select>
            </div>

            {analytics && (
                <div className="metrics-grid">
                    <div className="metric-card">
                        <h3>Total Requests</h3>
                        <div className="metric-value">{analytics.total_requests}</div>
                    </div>
                    <div className="metric-card">
                        <h3>Allowed</h3>
                        <div className="metric-value allowed">{analytics.allowed_requests}</div>
                    </div>
                    <div className="metric-card">
                        <h3>Denied</h3>
                        <div className="metric-value denied">{analytics.denied_requests}</div>
                    </div>
                    <div className="metric-card">
                        <h3>Whitelisted</h3>
                        <div className="metric-value whitelisted">{analytics.whitelisted_requests}</div>
                    </div>
                    <div className="metric-card">
                        <h3>Avg Response Time</h3>
                        <div className="metric-value">{analytics.avg_response_time_ms.toFixed(2)}ms</div>
                    </div>
                    <div className="metric-card">
                        <h3>Cache Hit Rate</h3>
                        <div className="metric-value">{analytics.cache_hit_rate.toFixed(1)}%</div>
                    </div>
                </div>
            )}

            <div className="charts-section">
                <div className="chart-container">
                    <h3>Request Results Distribution</h3>
                    <ResponsiveContainer width="100%" height={300}>
                        <PieChart>
                            <Pie
                                data={[
                                    { name: 'Allowed', value: analytics?.allowed_requests || 0, fill: '#4CAF50' },
                                    { name: 'Denied', value: analytics?.denied_requests || 0, fill: '#F44336' },
                                    { name: 'Whitelisted', value: analytics?.whitelisted_requests || 0, fill: '#2196F3' }
                                ]}
                                dataKey="value"
                                nameKey="name"
                            />
                            <Tooltip />
                            <Legend />
                        </PieChart>
                    </ResponsiveContainer>
                </div>

                <div className="chart-container">
                    <h3>Top Data Relationships</h3>
                    <div className="relationships-list">
                        {relationships.map((rel, index) => (
                            <div key={index} className="relationship-item">
                                <div className="relationship-type">{rel.relationship_type}</div>
                                <div className="relationship-data">
                                    {Object.entries(rel.data).map(([key, value]) => (
                                        value && <span key={key} className="data-item">{key}: {value}</span>
                                    ))}
                                </div>
                                <div className="relationship-frequency">{rel.frequency} occurrences</div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AnalyticsDashboard;
```

### 6. Configuration

#### Add to config.yaml
```yaml
logging:
  traffic_logging: true  # Enable traffic logging
  analytics_enabled: true  # Enable analytics processing
  retention_days: 90  # How long to keep traffic logs
  aggregation_schedule: "hourly"  # How often to run aggregations
```

### 7. Integration with Existing System

#### Update Filter Controller
```go
// In controllers/filter.go, update FilterRequestHandler
func FilterRequestHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()
        
        var input FilterRequest
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            return
        }

        // ... existing filter logic ...

        // Log the traffic
        trafficLogging := services.NewTrafficLoggingService(db)
        metadata := map[string]string{
            "client_ip": c.ClientIP(),
            "user_agent_raw": c.GetHeader("User-Agent"),
            "session_id": c.GetHeader("X-Session-ID"),
        }
        
        result := services.FilterResult{
            FinalResult:   finalResult.Result,
            FilterResults: map[string]interface{}{"details": finalResult},
            ResponseTime:  time.Since(startTime),
            CacheHit:      cacheHit,
        }
        
        go trafficLogging.LogFilterRequest(input, result, metadata)

        c.JSON(http.StatusOK, finalResult)
    }
}
```

## Benefits

✅ **Complete Traffic Visibility** - Log all filter requests with full context  
✅ **Data Relationship Analysis** - Understand correlations between different data types  
✅ **Performance Monitoring** - Track response times and cache effectiveness  
✅ **Security Insights** - Identify patterns in denied/whitelisted requests  
✅ **Real-time Analytics** - Live dashboard for monitoring traffic patterns  
✅ **Historical Analysis** - Long-term trend analysis and reporting  
✅ **Scalable Architecture** - Separate logging from main request flow  

This system provides comprehensive traffic logging and analytics capabilities while maintaining high performance through asynchronous processing and efficient data storage. 