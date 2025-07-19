# Firewall Application

A comprehensive firewall management system with Go backend and React frontend.

## Features

- **Real-time filtering** with support for IPs, emails, user agents, countries, charsets, and usernames
- **Geographic filtering** with automatic IP geolocation using MaxMind GeoLite2 database
- **Server-side filtering, sorting, and pagination** for optimal performance
- **Elasticsearch integration** for advanced search capabilities
- **Configurable distributed locking** for single and multi-instance deployments
- **RESTful API** with comprehensive documentation
- **Modern React frontend** with responsive design
- **Database connection pooling** with monitoring
- **Rate limiting** and security features
- **Comprehensive logging** and metrics
- **Input validation system** with field-level validation and security checks

## Quick Start

### Prerequisites

- Go 1.24+
- Node.js 18+
- MySQL 8.0+
- Elasticsearch 8.0+ (optional)
- Redis (optional, for distributed locking)

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd firewall
   ```

2. **Configure the application:**
   ```bash
   cp config/config.yaml.example config/config.yaml
   # Edit config/config.yaml with your settings
   ```

3. **Build and run:**
   ```bash
   # Build the backend
   go build -o firewall .
   
   # Install frontend dependencies
   cd firewall-app
   npm install
   npm run build
   cd ..
   
   # Run the application
   ./firewall
   ```

## Configuration

### Distributed Locking

The application supports both single-instance and multi-instance deployments through configurable distributed locking:

#### **Single Instance (Default)**
```yaml
locking:
  enabled: false  # No Redis required
```

#### **Multi-Instance**
```yaml
locking:
  enabled: true   # Requires Redis
redis:
  host: "your-redis-server"
  port: 6379
```

### Environment Variables

You can configure via environment variables:

```bash
# Database
export FIREWALL_DATABASE_HOST=localhost
export FIREWALL_DATABASE_PORT=3306
export FIREWALL_DATABASE_USER=user
export FIREWALL_DATABASE_PASSWORD=password
export FIREWALL_DATABASE_NAME=firewall

# Elasticsearch
export FIREWALL_ELASTIC_HOSTS=http://localhost:9200

# Redis (for distributed locking)
export FIREWALL_REDIS_HOST=localhost
export FIREWALL_REDIS_PORT=6379

# Distributed Locking
export FIREWALL_LOCKING_ENABLED=false
export FIREWALL_LOCKING_INCREMENTAL_TTL=5m
export FIREWALL_LOCKING_FULL_SYNC_TTL=30m
```

## API Documentation

- **Swagger UI**: `http://localhost:8081/swagger/index.html`
- **API Base URL**: `http://localhost:8081/api/v1`

### Key Endpoints

- `GET /api/v1/ips` - List IP addresses
- `GET /api/v1/emails` - List email addresses
- `GET /api/v1/useragents` - List user agents
- `GET /api/v1/countries` - List countries
- `GET /api/v1/charsets` - List charsets
- `GET /api/v1/usernames` - List usernames
- `GET /sync/status` - Check sync status and distributed locking

### Input Validation

All endpoints include comprehensive input validation:

- **IP Addresses**: Valid IPv4, IPv6, and CIDR notation
- **Email Addresses**: RFC 5321 compliant format
- **User Agents**: Length and character validation
- **Country Codes**: ISO 3166-1 alpha-2 format
- **Status Values**: `allowed`, `denied`, `whitelisted`
- **Regex Patterns**: Valid Go regex compilation
- **Pagination**: Page and limit validation
- **Search Parameters**: SQL injection prevention

For detailed validation documentation, see [docs/VALIDATION.md](docs/VALIDATION.md).

### IP/CIDR Conflict Detection

The system automatically detects conflicts when adding IP addresses or CIDR ranges:

- **IP in CIDR**: Detects when adding an IP already covered by a CIDR range
- **CIDR Overlap**: Detects when adding a CIDR that overlaps with existing ranges
- **CIDR Covers IPs**: Warns when a CIDR would cover existing individual IPs
- **Exact Matches**: Detects duplicate IPs or CIDR ranges
- **Pre-creation Checks**: Check for conflicts without creating entries

**Endpoints:**
- `POST /api/ip` - Create with automatic conflict detection
- `POST /api/ip/check-conflicts` - Check conflicts before creating
- `PUT /api/ip/:id` - Update with conflict detection

For detailed conflict detection documentation, see [docs/CONFLICT_DETECTION.md](docs/CONFLICT_DETECTION.md).

### Geographic Filtering

The system supports automatic geographic filtering using MaxMind's GeoLite2 database:

- **Automatic Geolocation**: IP addresses are automatically geolocated to countries
- **Manual Override**: Users can provide country codes manually
- **Private IP Handling**: Private/local IPs are skipped for geolocation
- **Country Rules**: Uses existing country filtering system
- **Performance**: Fast local database lookups with no external API calls

**Usage Examples:**
```bash
# Automatic geolocation
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{"ip": "91.67.0.1", "country": ""}'
# Response: {"result":"denied","reason":"country denied","field":"country","value":"DE"}

# Manual country override
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{"ip": "8.8.8.8", "country": "DE"}'
# Response: {"result":"denied","reason":"country denied","field":"country","value":"DE"}
```

**Required Files:**
- `GeoLite2-Country.mmdb` in the root directory (~9MB)

For detailed geographic filtering documentation, see [docs/GEOGRAPHIC_FILTERING.md](docs/GEOGRAPHIC_FILTERING.md).

## Development

### Backend Development

```bash
# Run with hot reload
go run main.go

# Run tests
go test ./...

# Build for production
go build -ldflags="-s -w" -o firewall .
```

### Frontend Development

```bash
cd firewall-app

# Install dependencies
npm install

# Start development server
npm start

# Build for production
npm run build
```

## Deployment

### Single Instance

For single-instance deployments, no additional infrastructure is required:

```yaml
# config/config.yaml
locking:
  enabled: false  # No Redis needed
```

### Multi-Instance

For high-availability deployments:

```yaml
# config/config.yaml
locking:
  enabled: true
redis:
  host: "redis-cluster.example.com"
  port: 6379
```

### Docker

```bash
# Build the application
docker build -t firewall .

# Run with configuration
docker run -p 8081:8081 \
  -v $(pwd)/config:/app/config \
  firewall
```

## Monitoring

### Health Checks

The firewall provides comprehensive health check endpoints for monitoring system status:

#### **Comprehensive Health Check**
- **Endpoint**: `GET /api/health`
- **Description**: Detailed health checks for all system components
- **Response**: `200 OK` (healthy) or `503 Service Unavailable` (unhealthy)
- **Use Cases**: Detailed monitoring, alerting systems, troubleshooting

#### **Simple Health Check**
- **Endpoint**: `GET /api/health/simple`
- **Description**: Lightweight health check for load balancers
- **Response**: `200 OK` with basic status
- **Use Cases**: Load balancer health checks, basic availability monitoring

#### **Monitored Services**
- **Database**: MySQL connection and query performance
- **Elasticsearch**: Connection and ping response time
- **Cache**: Set/get operations and performance
- **Event Processor**: Service availability
- **Distributed Lock Service**: Service availability

#### **Example Response**
```json
{
  "status": "healthy",
  "timestamp": "2025-07-19T08:34:17Z",
  "version": "1.0.0",
  "services": {
    "database": {"status": "healthy"},
    "elasticsearch": {"status": "healthy", "response_time_ms": 17},
    "cache": {"status": "healthy"},
    "event_processor": {"status": "healthy"},
    "distributed_lock": {"status": "healthy"}
  }
}
```

For detailed documentation, see [docs/HEALTH_CHECK.md](docs/HEALTH_CHECK.md).

### Additional Monitoring

- `GET /api/sync/status` - Sync and locking status
- `GET /api/system-stats` - System performance metrics

### Metrics

The application provides metrics for:
- Database connection pool status
- Sync operation status
- Distributed lock status (when enabled)
- Request/response statistics

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Support

For issues and questions:
- Check the [documentation](docs/)
- Review the [API documentation](http://localhost:8081/swagger/index.html)
- Open an issue on GitHub 