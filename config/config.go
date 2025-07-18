// Package config config/config.go
package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Elastic  ElasticConfig  `mapstructure:"elastic"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Security SecurityConfig `mapstructure:"security"`
	Locking  LockingConfig  `mapstructure:"locking"`
	Caching  CachingConfig  `mapstructure:"caching"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	Mode         string        `mapstructure:"mode"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parse_time"`
	Loc             string        `mapstructure:"loc"`
}

// ElasticConfig holds Elasticsearch-related configuration
type ElasticConfig struct {
	Hosts    []string      `mapstructure:"hosts"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Index    string        `mapstructure:"index"`
}

// RedisConfig holds Redis-related configuration
type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	CORSEnabled     bool          `mapstructure:"cors_enabled"`
	CORSOrigins     []string      `mapstructure:"cors_origins"`
	CORSMethods     []string      `mapstructure:"cors_methods"`
	CORSHeaders     []string      `mapstructure:"cors_headers"`
	RateLimit       int           `mapstructure:"rate_limit"`
	RateLimitWindow time.Duration `mapstructure:"rate_limit_window"`
}

// LockingConfig holds distributed locking configuration
type LockingConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	LockTTL         time.Duration `mapstructure:"lock_ttl"`
	IncrementalTTL  time.Duration `mapstructure:"incremental_ttl"`
	FullSyncTTL     time.Duration `mapstructure:"full_sync_ttl"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

// CachingConfig holds distributed caching configuration
type CachingConfig struct {
	Distributed bool          `mapstructure:"distributed"`
	DefaultTTL  time.Duration `mapstructure:"default_ttl"`
	FilterTTL   time.Duration `mapstructure:"filter_ttl"`
	ListTTL     time.Duration `mapstructure:"list_ttl"`
	StatsTTL    time.Duration `mapstructure:"stats_ttl"`
}

// Global config instance
var AppConfig *Config

// InitConfig initializes the configuration using Viper
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/firewall")

	// Set environment variable prefix
	viper.SetEnvPrefix("FIREWALL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Enable environment variable binding
	viper.AutomaticEnv()

	// Set sensible defaults
	setDefaults()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: Error reading config file: %v", err)
		}
	}

	// Unmarshal configuration
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Validate configuration
	if err := validateConfig(AppConfig); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	log.Println("Configuration loaded successfully")
}

// setDefaults sets sensible defaults for all configuration options
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8081)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "60s")
	viper.SetDefault("server.mode", "debug")

	// Database defaults
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "user")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "firewall")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	viper.SetDefault("database.conn_max_idle_time", "1m")
	viper.SetDefault("database.ssl_mode", "false")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parse_time", true)
	viper.SetDefault("database.loc", "Local")

	// Elasticsearch defaults
	viper.SetDefault("elastic.hosts", []string{"http://localhost:9200"})
	viper.SetDefault("elastic.username", "")
	viper.SetDefault("elastic.password", "")
	viper.SetDefault("elastic.timeout", "30s")
	viper.SetDefault("elastic.index", "firewall")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.timeout", "5s")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.max_size", 100)
	viper.SetDefault("logging.max_backups", 3)
	viper.SetDefault("logging.max_age", 28)
	viper.SetDefault("logging.compress", true)

	// Security defaults
	viper.SetDefault("security.cors_enabled", true)
	viper.SetDefault("security.cors_origins", []string{"http://localhost:3000"})
	viper.SetDefault("security.cors_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("security.cors_headers", []string{"Content-Type", "Authorization"})
	viper.SetDefault("security.rate_limit", 1000)
	viper.SetDefault("security.rate_limit_window", "1m")

	// Distributed locking defaults
	viper.SetDefault("locking.enabled", false) // Disabled by default for single instances
	viper.SetDefault("locking.lock_ttl", "5m")
	viper.SetDefault("locking.incremental_ttl", "5m")
	viper.SetDefault("locking.full_sync_ttl", "30m")
	viper.SetDefault("locking.cleanup_interval", "10m")

	// Distributed caching defaults
	viper.SetDefault("caching.distributed", false) // In-memory by default for single instances
	viper.SetDefault("caching.default_ttl", "5m")
	viper.SetDefault("caching.filter_ttl", "5m")
	viper.SetDefault("caching.list_ttl", "2m")
	viper.SetDefault("caching.stats_ttl", "30s")
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate server configuration
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// Validate database configuration
	if config.Database.Port < 1 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Database.Port)
	}
	if config.Database.MaxOpenConns < 1 {
		return fmt.Errorf("invalid max open connections: %d", config.Database.MaxOpenConns)
	}
	if config.Database.MaxIdleConns < 0 {
		return fmt.Errorf("invalid max idle connections: %d", config.Database.MaxIdleConns)
	}
	if config.Database.MaxIdleConns > config.Database.MaxOpenConns {
		return fmt.Errorf("max idle connections cannot be greater than max open connections")
	}

	// Validate logging configuration
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[config.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", config.Logging.Level)
	}

	// Validate security configuration
	if config.Security.RateLimit < 1 {
		return fmt.Errorf("invalid rate limit: %d", config.Security.RateLimit)
	}

	return nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s&timeout=10s&readTimeout=30s&writeTimeout=30s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.Charset, c.ParseTime, c.Loc)
}

// GetElasticsearchURLs returns the Elasticsearch URLs
func (c *ElasticConfig) GetElasticsearchURLs() []string {
	return c.Hosts
}

// GetRedisAddr returns the Redis address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsDevelopment returns true if the application is in development mode
func (c *ServerConfig) IsDevelopment() bool {
	return c.Mode == "debug" || c.Mode == "development"
}

// IsProduction returns true if the application is in production mode
func (c *ServerConfig) IsProduction() bool {
	return c.Mode == "release" || c.Mode == "production"
}

// InitMySQL initializes and connects to the MySQL database with connection pooling
func InitMySQL() {
	// Initialize config if not already done
	if AppConfig == nil {
		InitConfig()
	}

	// Build DSN with connection pooling parameters
	dsn := AppConfig.Database.GetDSN()

	// Configure GORM with connection pooling
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Get underlying sql.DB to configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	// Configure connection pooling
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(AppConfig.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(AppConfig.Database.ConnMaxIdleTime)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	DB = db

	log.Printf("MySQL connected successfully with connection pooling:")
	log.Printf("  Max Open Connections: %d", AppConfig.Database.MaxOpenConns)
	log.Printf("  Max Idle Connections: %d", AppConfig.Database.MaxIdleConns)
	log.Printf("  Connection Max Lifetime: %v", AppConfig.Database.ConnMaxLifetime)
	log.Printf("  Connection Max Idle Time: %v", AppConfig.Database.ConnMaxIdleTime)
}

// GetDBStats returns current database connection statistics
func GetDBStats() map[string]interface{} {
	if DB == nil {
		return map[string]interface{}{
			"error": "Database not initialized",
		}
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to get sql.DB: %v", err),
		}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}
