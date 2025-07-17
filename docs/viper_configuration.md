# Viper Configuration System

This document describes the comprehensive configuration system using Viper for the firewall application.

## Overview

The application now uses [Viper](https://github.com/spf13/viper) for configuration management, which provides:

- **Multiple formats**: YAML, JSON, TOML, HCL, INI, envfile, Java properties
- **Environment variables**: Automatic binding with prefix support
- **Configuration files**: Hierarchical configuration file support
- **Remote configuration**: Etcd, Consul, AWS Parameter Store
- **Live watching**: Hot reloading of configuration files
- **Default values**: Sensible defaults for all settings

## Configuration Sources

Viper loads configuration in the following order (later sources override earlier ones):

1. **Default values** (hardcoded in the application)
2. **Configuration file** (`config.yaml`, `config.json`, etc.)
3. **Environment variables** (with `FIREWALL_` prefix)
4. **Command line flags** (if implemented)

## Configuration Structure

### Server Configuration
```yaml
server:
  port: 8081                    # Server port
  host: "0.0.0.0"              # Server host
  read_timeout: "30s"           # Read timeout
  write_timeout: "30s"          # Write timeout
  idle_timeout: "60s"           # Idle timeout
  mode: "debug"                 # Gin mode (debug, release)
```

### Database Configuration
```yaml
database:
  host: "127.0.0.1"            # Database host
  port: 3306                    # Database port
  user: "user"                  # Database user
  password: "password"          # Database password
  name: "firewall"              # Database name
  max_open_conns: 25           # Max open connections
  max_idle_conns: 5            # Max idle connections
  conn_max_lifetime: "5m"      # Connection max lifetime
  conn_max_idle_time: "1m"     # Connection max idle time
  ssl_mode: "false"            # SSL mode
  charset: "utf8mb4"           # Character set
  parse_time: true             # Parse time
  loc: "Local"                 # Location
```

### Elasticsearch Configuration
```yaml
elastic:
  hosts:                        # Elasticsearch hosts
    - "http://localhost:9200"
  username: ""                  # Username (if authentication enabled)
  password: ""                  # Password (if authentication enabled)
  timeout: "30s"               # Connection timeout
  index: "firewall"            # Default index name
```

### Redis Configuration
```yaml
redis:
  host: "localhost"             # Redis host
  port: 6379                   # Redis port
  password: ""                 # Redis password
  db: 0                        # Redis database number
  timeout: "5s"                # Connection timeout
```

### Logging Configuration
```yaml
logging:
  level: "info"                # Log level (debug, info, warn, error)
  format: "json"               # Log format (json, text)
  output: "stdout"             # Output destination
  max_size: 100                # Max file size in MB
  max_backups: 3               # Max number of backup files
  max_age: 28                  # Max age in days
  compress: true               # Compress rotated files
```

### Security Configuration
```yaml
security:
  cors_enabled: true           # Enable CORS
  cors_origins:                # Allowed origins
    - "http://localhost:3000"
    - "http://localhost:8081"
  cors_methods:                # Allowed HTTP methods
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  cors_headers:                # Allowed headers
    - "Content-Type"
    - "Authorization"
    - "X-Requested-With"
  rate_limit: 1000             # Rate limit per window
  rate_limit_window: "1m"      # Rate limit window
```

## Environment Variables

All configuration options can be overridden using environment variables with the `FIREWALL_` prefix:

### Server Environment Variables
```bash
export FIREWALL_SERVER_PORT=8081
export FIREWALL_SERVER_HOST="0.0.0.0"
export FIREWALL_SERVER_MODE="release"
```

### Database Environment Variables
```bash
export FIREWALL_DATABASE_HOST="production-db.example.com"
export FIREWALL_DATABASE_PORT=3306
export FIREWALL_DATABASE_USER="firewall_user"
export FIREWALL_DATABASE_PASSWORD="secure_password"
export FIREWALL_DATABASE_NAME="firewall_prod"
export FIREWALL_DATABASE_MAX_OPEN_CONNS=50
export FIREWALL_DATABASE_MAX_IDLE_CONNS=10
export FIREWALL_DATABASE_CONN_MAX_LIFETIME="10m"
export FIREWALL_DATABASE_CONN_MAX_IDLE_TIME="2m"
```

### Elasticsearch Environment Variables
```bash
export FIREWALL_ELASTIC_HOSTS="http://es1:9200,http://es2:9200"
export FIREWALL_ELASTIC_USERNAME="elastic"
export FIREWALL_ELASTIC_PASSWORD="secure_password"
export FIREWALL_ELASTIC_TIMEOUT="30s"
```

### Security Environment Variables
```bash
export FIREWALL_SECURITY_CORS_ENABLED=true
export FIREWALL_SECURITY_CORS_ORIGINS="http://app.example.com,https://app.example.com"
export FIREWALL_SECURITY_RATE_LIMIT=2000
export FIREWALL_SECURITY_RATE_LIMIT_WINDOW="1m"
```

## Configuration File Locations

Viper searches for configuration files in the following locations:

1. Current working directory (`.`)
2. `./config/` directory
3. `/etc/firewall/` directory

Supported file formats:
- `config.yaml`
- `config.yml`
- `config.json`
- `config.toml`
- `config.hcl`
- `config.ini`

## Configuration Validation

The application validates configuration on startup:

### Server Validation
- Port must be between 1 and 65535
- Host must be a valid IP address or hostname

### Database Validation
- Port must be between 1 and 65535
- Max open connections must be > 0
- Max idle connections must be >= 0
- Max idle connections cannot exceed max open connections

### Logging Validation
- Log level must be one of: debug, info, warn, error
- Output must be a valid destination

### Security Validation
- Rate limit must be > 0
- CORS origins must be valid URLs (if specified)

## Usage Examples

### Development Environment
```yaml
# config.yaml
server:
  port: 8081
  mode: "debug"

database:
  host: "localhost"
  port: 3306
  user: "user"
  password: "password"
  name: "firewall_dev"
  max_open_conns: 10
  max_idle_conns: 3

logging:
  level: "debug"
  format: "text"
  output: "stdout"
```

### Production Environment
```yaml
# config.yaml
server:
  port: 8081
  mode: "release"
  host: "0.0.0.0"

database:
  host: "production-db.example.com"
  port: 3306
  user: "firewall_user"
  password: "secure_password"
  name: "firewall_prod"
  max_open_conns: 50
  max_idle_conns: 10
  conn_max_lifetime: "10m"
  conn_max_idle_time: "2m"

logging:
  level: "info"
  format: "json"
  output: "/var/log/firewall/app.log"
  max_size: 100
  max_backups: 5
  max_age: 30
  compress: true

security:
  cors_enabled: true
  cors_origins:
    - "https://app.example.com"
  rate_limit: 2000
  rate_limit_window: "1m"
```

### Docker Environment
```yaml
# docker-compose.yml
services:
  app:
    environment:
      - FIREWALL_SERVER_PORT=8081
      - FIREWALL_SERVER_MODE=release
      - FIREWALL_DATABASE_HOST=mariadb
      - FIREWALL_DATABASE_PORT=3306
      - FIREWALL_DATABASE_USER=user
      - FIREWALL_DATABASE_PASSWORD=password
      - FIREWALL_DATABASE_NAME=firewall
      - FIREWALL_ELASTIC_HOSTS=http://elasticsearch:9200
      - FIREWALL_REDIS_HOST=redis
      - FIREWALL_REDIS_PORT=6379
```

## Hot Reloading

Viper supports hot reloading of configuration files. To enable this feature:

```go
// In your application code
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    log.Println("Config file changed:", e.Name)
    // Reload configuration
    config.InitConfig()
})
```

## Configuration Access

Access configuration values in your application:

```go
import "firewall/config"

// Access server configuration
port := config.AppConfig.Server.Port
host := config.AppConfig.Server.Host

// Access database configuration
dbHost := config.AppConfig.Database.Host
dbPort := config.AppConfig.Database.Port

// Access security configuration
corsEnabled := config.AppConfig.Security.CORSEnabled
rateLimit := config.AppConfig.Security.RateLimit

// Check environment
isDev := config.AppConfig.Server.IsDevelopment()
isProd := config.AppConfig.Server.IsProduction()
```

## Best Practices

1. **Use environment variables for secrets**: Never put passwords in configuration files
2. **Use configuration files for defaults**: Set sensible defaults in configuration files
3. **Validate configuration**: Always validate configuration on startup
4. **Use hierarchical configuration**: Organize configuration logically
5. **Document configuration**: Document all configuration options
6. **Use type-safe access**: Use the structured configuration objects
7. **Test configuration**: Test configuration loading in your tests

## Migration from Environment Variables

If you were using the old environment variable system, here's how to migrate:

### Old Environment Variables
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=user
export DB_PASSWORD=password
export DB_NAME=firewall
export DB_MAX_OPEN_CONNS=25
export DB_MAX_IDLE_CONNS=5
export DB_CONN_MAX_LIFETIME=5m
export DB_CONN_MAX_IDLE_TIME=1m
```

### New Environment Variables
```bash
export FIREWALL_DATABASE_HOST=localhost
export FIREWALL_DATABASE_PORT=3306
export FIREWALL_DATABASE_USER=user
export FIREWALL_DATABASE_PASSWORD=password
export FIREWALL_DATABASE_NAME=firewall
export FIREWALL_DATABASE_MAX_OPEN_CONNS=25
export FIREWALL_DATABASE_MAX_IDLE_CONNS=5
export FIREWALL_DATABASE_CONN_MAX_LIFETIME=5m
export FIREWALL_DATABASE_CONN_MAX_IDLE_TIME=1m
```

The old environment variables are still supported for backward compatibility, but the new `FIREWALL_` prefixed variables are recommended. 