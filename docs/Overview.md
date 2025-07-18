# Firewall System Architecture Overview

## Table of Contents
1. [System Overview](#system-overview)
2. [Architecture Components](#architecture-components)
3. [Data Flow](#data-flow)
4. [Caching Strategy](#caching-strategy)
5. [Scaling Strategy](#scaling-strategy)
6. [Configuration Management](#configuration-management)
7. [Deployment Options](#deployment-options)

## System Overview

The Firewall System is a comprehensive web application that provides real-time filtering and blocking capabilities for network traffic, emails, user agents, and other data types. It's built as a modern web application with a Go backend and React frontend, designed to be scalable from single-instance deployments to full horizontal scaling.

### Key Features
- **Real-time filtering** of IP addresses, emails, user agents, countries, and usernames
- **Regex pattern matching** for flexible rule definitions
- **Character set detection** and filtering
- **Elasticsearch integration** for fast search and indexing
- **Configurable caching** (in-memory or distributed)
- **Distributed locking** for multi-instance deployments
- **Horizontal scaling** support
- **RESTful API** with comprehensive endpoints
- **Modern web interface** with real-time updates

## Architecture Components

### 1. Backend (Go)
```
├── main.go                 # Application entry point
├── config/                 # Configuration management
│   ├── config.go          # Viper-based config with environment support
│   └── elasticsearch.go   # Elasticsearch client setup
├── controllers/           # HTTP request handlers
│   ├── filter.go         # Filter request processing
│   └── rules.go          # CRUD operations for all data types
├── models/               # Database models
│   └── models.go         # GORM models for all entities
├── services/             # Business logic layer
│   ├── cache_factory.go  # Cache abstraction (in-memory/distributed)
│   ├── cache_service.go  # In-memory cache implementation
│   ├── distributed_cache.go # Redis-based distributed cache
│   ├── event_service.go  # Event processing system
│   ├── filter_service.go # Core filtering logic
│   ├── elasticsearch_sync.go # ES indexing and sync
│   ├── scheduled_sync.go # Background sync service
│   └── retry_service.go  # Failed operation retry queue
├── middleware/           # HTTP middleware
│   └── rate_limit.go    # Rate limiting
├── routes/              # Route definitions
│   └── routes.go        # API route setup
└── migrations/          # Database migrations
    └── migrations.go    # GORM auto-migrations
```

### 2. Frontend (React)
```
firewall-app/
├── src/
│   ├── components/       # Reusable UI components
│   │   ├── DashboardStats.js
│   │   ├── FilterForm.js
│   │   ├── ListView.js
│   │   └── SystemStats.js
│   ├── pages/           # Page components
│   │   └── SystemHealthPage.js
│   ├── App.js           # Main application component
│   └── axiosConfig.js   # API client configuration
```

### 3. Data Storage
- **MySQL Database**: Primary data storage for all entities
- **Elasticsearch**: Search and indexing engine for fast filtering
- **Redis** (optional): Distributed caching and locking for horizontal scaling

## Data Flow

### 1. Filter Request Flow
```
Client Request → Gin Router → Rate Limiting Middleware → Filter Controller → Filter Service → Elasticsearch → Response
```

**Detailed Flow:**
1. **Client sends filter request** with IP, email, user agent, country, username
2. **Gin Router** receives request and routes to appropriate handler
3. **Rate Limiting Middleware** applies request throttling if configured
4. **Filter Controller** normalizes data (e.g., Gmail email normalization)
5. **Cache Check**: Look for cached filter result
6. **Filter Service** performs concurrent filtering:
   - IP address lookup in Elasticsearch
   - Email pattern matching (exact + regex)
   - User agent pattern matching (exact + regex)
   - Country code validation
   - Username pattern matching (exact + regex)
   - Character set detection
7. **Result Aggregation**: Combine all filter results
8. **Cache Storage**: Store result for future requests
9. **Response**: Return filter decision (allowed/denied/whitelisted)

### 2. Data Management Flow
```
Admin Action → Controller → Database → Event Service → Elasticsearch Sync → Cache Invalidation
```

**Detailed Flow:**
1. **Admin creates/updates/deletes** rule via web interface
2. **Controller** validates and saves to MySQL database
3. **Event Service** publishes change event
4. **Elasticsearch Sync** indexes data for fast searching
5. **Cache Invalidation** clears related cached data
6. **Real-time updates** reflected in web interface

## Caching Strategy

### Cache Architecture
The system uses a **configurable caching layer** that automatically switches between implementations:

```go
// Cache Factory Pattern
cache := services.GetCacheFactory()
```

### Cache Types
1. **In-Memory Cache** (`caching.distributed: false`)
   - Fast local caching
   - No external dependencies
   - Perfect for single-instance deployments

2. **Distributed Cache** (`caching.distributed: true`)
   - Redis-based shared cache
   - Consistent across all instances
   - Required for horizontal scaling

### Cache Patterns
```go
// Filter Results
cacheKey := cache.FilterKey("ip", "192.168.1.1")
cache.Set(cacheKey, result, 5*time.Minute)

// List Data
listKey := cache.ListKey("ip", "1", "10", "", "all")
cache.Set(listKey, data, 2*time.Minute)

// System Stats
statsKey := cache.StatsKey("system")
cache.Set(statsKey, stats, 30*time.Second)
```

### Cache Invalidation
```go
// After data changes
cache.InvalidateAll("ip")      // Clear all IP-related cache
cache.InvalidateFilter("ip")   // Clear IP filter cache
cache.InvalidateList("ip")     // Clear IP list cache
cache.InvalidateStats("ip")    // Clear IP stats cache
```

## Scaling Strategy

### Single Instance Deployment
```
┌─────────────────┐
│   Web Browser   │
└─────────┬───────┘
          │
┌─────────▼───────┐
│   Firewall App  │
│                 │
│  ┌───────────┐  │
│  │   React   │  │  ← Frontend
│  └───────────┘  │
│                 │
│  ┌───────────┐  │
│  │    Go     │  │  ← Backend
│  └───────────┘  │
│                 │
│  ┌───────────┐  │
│  │ In-Memory │  │  ← Local Cache
│  │   Cache   │  │
│  └───────────┘  │
└─────────┬───────┘
          │
┌─────────▼───────┐
│     MySQL       │  ← Database
└─────────┬───────┘
          │
┌─────────▼───────┐
│ Elasticsearch   │  ← Search Engine
└─────────────────┘
```

### Horizontal Scaling Deployment
```
┌─────────────────┐
│   Load Balancer │
└─────────┬───────┘
          │
    ┌─────┴─────┐
    │           │
┌───▼───┐   ┌───▼───┐
│App #1 │   │App #2 │
│       │   │       │
│ ┌───┐ │   │ ┌───┐ │
│ │Go │ │   │ │Go │ │  ← Multiple Backend Instances
│ └───┘ │   │ └───┘ │
│       │   │       │
│ ┌───┐ │   │ ┌───┐ │
│ │React│ │   │ │React│ │  ← Frontend (static)
│ └───┘ │   │ └───┘ │
└───┬───┘   └───┬───┘
    │           │
    └─────┬─────┘
          │
    ┌─────▼─────┐
    │   Redis   │  ← Distributed Cache & Locking
    └─────┬─────┘
          │
    ┌─────▼─────┐
    │   MySQL   │  ← Database (shared)
    └─────┬─────┘
          │
    ┌─────▼─────┐
    │Elasticsearch│  ← Search Cluster
    └─────────────┘
```

### Scaling Components

#### 1. Distributed Locking
```yaml
locking:
  distributed: true  # Enable for horizontal scaling
  lock_ttl: "5m"
  incremental_ttl: "5m"
  full_sync_ttl: "30m"
```

#### 2. Distributed Caching
```yaml
caching:
  distributed: true  # Enable for horizontal scaling
  default_ttl: "5m"
  filter_ttl: "5m"
  list_ttl: "2m"
  stats_ttl: "30s"
```

#### 3. Database Connection Pooling
```yaml
database:
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: "5m"
  conn_max_idle_time: "1m"
```

## Configuration Management

### Configuration Sources (Priority Order)
1. **Environment Variables** (highest priority)
2. **Configuration Files** (`config.yaml`)
3. **Default Values** (lowest priority)

### Environment Variable Examples
```bash
# Database
FIREWALL_DATABASE_HOST=localhost
FIREWALL_DATABASE_PORT=3306
FIREWALL_DATABASE_USER=user
FIREWALL_DATABASE_PASSWORD=password

# Elasticsearch
FIREWALL_ELASTIC_HOSTS=http://localhost:9200
FIREWALL_ELASTIC_USERNAME=
FIREWALL_ELASTIC_PASSWORD=

# Redis (for horizontal scaling)
FIREWALL_REDIS_HOST=localhost
FIREWALL_REDIS_PORT=6379
FIREWALL_REDIS_PASSWORD=

# Scaling Configuration
FIREWALL_LOCKING_DISTRIBUTED=false
FIREWALL_CACHING_DISTRIBUTED=false

# Server Configuration
FIREWALL_SERVER_PORT=8081
FIREWALL_SERVER_HOST=0.0.0.0
```

### Configuration File Structure
```yaml
server:
  port: 8081
  host: "0.0.0.0"
  mode: "debug"

database:
  host: "127.0.0.1"
  port: 3306
  max_open_conns: 25
  max_idle_conns: 5

elastic:
  hosts: ["http://localhost:9200"]
  index: "firewall"

redis:
  host: "localhost"
  port: 6379

locking:
  distributed: false  # Set to true for horizontal scaling
  lock_ttl: "5m"

caching:
  distributed: false  # Set to true for horizontal scaling
  default_ttl: "5m"
```

## Deployment Options

### 1. Development Environment
```bash
# Single instance with in-memory cache
go run main.go
```

### 2. Production Single Instance
```bash
# Build and run
go build -o firewall main.go
./firewall
```

### 3. Docker Single Instance
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o firewall main.go
EXPOSE 8081
CMD ["./firewall"]
```

### 4. Docker Compose (Horizontal Scaling)
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8081:8081"
    environment:
      - FIREWALL_LOCKING_DISTRIBUTED=true
      - FIREWALL_CACHING_DISTRIBUTED=true
      - FIREWALL_REDIS_HOST=redis
    depends_on:
      - redis
      - mysql
      - elasticsearch

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: firewall
    ports:
      - "3306:3306"

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.8.0
    environment:
      - discovery.type=single-node
    ports:
      - "9200:9200"
```

### 5. Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: firewall
spec:
  replicas: 3
  selector:
    matchLabels:
      app: firewall
  template:
    metadata:
      labels:
        app: firewall
    spec:
      containers:
      - name: firewall
        image: firewall:latest
        ports:
        - containerPort: 8081
        env:
        - name: FIREWALL_LOCKING_DISTRIBUTED
          value: "true"
        - name: FIREWALL_CACHING_DISTRIBUTED
          value: "true"
        - name: FIREWALL_REDIS_HOST
          value: "redis-service"
```

## Performance Characteristics

### Throughput
- **Single Instance**: ~1000-5000 requests/second (depending on filter complexity)
- **Horizontal Scaling**: Linear scaling with number of instances

### Latency
- **Cache Hit**: <1ms
- **Cache Miss + Database**: 5-50ms
- **Elasticsearch Query**: 10-100ms

### Memory Usage
- **In-Memory Cache**: ~50-200MB (depending on cache size)
- **Application**: ~100-500MB
- **Total per instance**: ~200-800MB

### Storage Requirements
- **MySQL**: ~1-10GB (depending on data volume)
- **Elasticsearch**: ~2-20GB (depending on data volume)
- **Redis**: ~100MB-1GB (depending on cache size)

## Monitoring and Observability

### Health Checks
- `/api/system-stats` - System health and performance metrics
- `/api/status` - Service status information
- `/api/sync/status` - Sync operation status

### Metrics Available
- Database connection pool stats
- Cache hit/miss ratios
- Filter request latency
- Sync operation status
- Memory and CPU usage
- Active locks and cache items

### Logging
- Structured JSON logging
- Configurable log levels
- Request/response logging
- Error tracking and debugging

## Security Considerations

### API Security
- Rate limiting on all endpoints
- CORS configuration for web interface
- Input validation and sanitization
- SQL injection protection via GORM

### Data Security
- Database connection encryption
- Elasticsearch security (optional)
- Redis authentication (optional)
- Environment variable configuration

### Network Security
- HTTPS/TLS support (configured in reverse proxy)
- Network isolation in containerized deployments
- Firewall rules for service communication

## Troubleshooting

### Common Issues
1. **Cache Inconsistency**: Check distributed cache connectivity
2. **Sync Failures**: Verify Elasticsearch connectivity and permissions
3. **Performance Issues**: Monitor cache hit rates and database connections
4. **Scaling Problems**: Ensure Redis and distributed locking are properly configured

### Debug Tools
- `/api/system-stats` - System performance metrics
- `/api/cache/flush` - Cache management
- `/api/sync/force` - Manual sync operations
- Log analysis for error tracking

This architecture provides a robust, scalable foundation for real-time filtering and blocking operations, with clear separation of concerns and configurable scaling options. 