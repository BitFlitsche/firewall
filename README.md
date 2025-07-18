# Firewall Application

A comprehensive firewall management system with Go backend and React frontend.

## Features

- **Real-time filtering** with support for IPs, emails, user agents, countries, charsets, and usernames
- **Server-side filtering, sorting, and pagination** for optimal performance
- **Elasticsearch integration** for advanced search capabilities
- **Configurable distributed locking** for single and multi-instance deployments
- **RESTful API** with comprehensive documentation
- **Modern React frontend** with responsive design
- **Database connection pooling** with monitoring
- **Rate limiting** and security features
- **Comprehensive logging** and metrics

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

- `GET /health` - Application health
- `GET /sync/status` - Sync and locking status

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