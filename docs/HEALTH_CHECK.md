# Health Check Endpoints

The firewall provides comprehensive health check endpoints to monitor the status of all system components.

## Overview

The health check system monitors the following services:
- **Database**: MySQL connection and query performance
- **Elasticsearch**: Connection and ping response time
- **Cache**: Set/get operations and performance
- **Event Processor**: Service availability
- **Distributed Lock Service**: Service availability

## Endpoints

### 1. Comprehensive Health Check

**Endpoint**: `GET /api/health`

**Description**: Performs detailed health checks on all system components and returns comprehensive status information.

**Response Codes**:
- `200 OK`: All services are healthy
- `503 Service Unavailable`: One or more services are unhealthy

**Response Format**:
```json
{
  "status": "healthy",
  "timestamp": "2025-07-19T08:34:17Z",
  "version": "1.0.0",
  "services": {
    "database": {
      "status": "healthy"
    },
    "elasticsearch": {
      "status": "healthy",
      "response_time_ms": 17
    },
    "cache": {
      "status": "healthy"
    },
    "event_processor": {
      "status": "healthy"
    },
    "distributed_lock": {
      "status": "healthy"
    }
  }
}
```

**Use Cases**:
- Detailed system monitoring
- Alerting systems
- Performance monitoring
- Troubleshooting

### 2. Simple Health Check

**Endpoint**: `GET /api/health/simple`

**Description**: Provides a lightweight health check suitable for load balancers and basic monitoring.

**Response Codes**:
- `200 OK`: Service is running

**Response Format**:
```json
{
  "status": "ok",
  "timestamp": "2025-07-19T08:34:01Z"
}
```

**Use Cases**:
- Load balancer health checks
- Basic service availability monitoring
- Quick status verification

## Service Monitoring Details

### Database Health Check
- **Test**: Executes `SELECT 1` query
- **Metrics**: Response time in milliseconds
- **Failure**: Connection errors or query timeouts

### Elasticsearch Health Check
- **Test**: Sends ping request to Elasticsearch cluster
- **Metrics**: Response time in milliseconds
- **Failure**: Connection errors, cluster health issues

### Cache Health Check
- **Test**: Performs set/get operations with test data
- **Metrics**: Response time in milliseconds
- **Failure**: Cache connection errors, operation failures

### Event Processor Health Check
- **Test**: Verifies service initialization
- **Metrics**: Service availability
- **Failure**: Service not initialized

### Distributed Lock Service Health Check
- **Test**: Verifies service initialization
- **Metrics**: Service availability
- **Failure**: Service not initialized

## Monitoring Integration

### Prometheus Metrics
The health check endpoints can be integrated with Prometheus for metrics collection:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'firewall-health'
    static_configs:
      - targets: ['localhost:8081']
    metrics_path: '/api/health'
    scrape_interval: 30s
```

### Alerting Rules
Example alerting rules for monitoring systems:

```yaml
# Alert when any service is unhealthy
- alert: FirewallServiceUnhealthy
  expr: firewall_health_status == 0
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "Firewall service is unhealthy"
    description: "One or more firewall services are reporting unhealthy status"

# Alert when response times are high
- alert: FirewallHighResponseTime
  expr: firewall_service_response_time_ms > 1000
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "Firewall service response time is high"
    description: "Service response time exceeds 1 second"
```

### Load Balancer Configuration
Example configuration for load balancers:

```nginx
# nginx.conf
upstream firewall_backend {
    server 127.0.0.1:8081;
}

server {
    listen 80;
    server_name firewall.example.com;

    location / {
        proxy_pass http://firewall_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://firewall_backend/api/health/simple;
        access_log off;
    }
}
```

## Error Handling

### Unhealthy Service Response
When a service is unhealthy, the response includes detailed error information:

```json
{
  "status": "unhealthy",
  "timestamp": "2025-07-19T08:34:17Z",
  "version": "1.0.0",
  "services": {
    "database": {
      "status": "healthy"
    },
    "elasticsearch": {
      "status": "unhealthy",
      "message": "Elasticsearch ping failed: connection refused"
    },
    "cache": {
      "status": "healthy"
    },
    "event_processor": {
      "status": "healthy"
    },
    "distributed_lock": {
      "status": "healthy"
    }
  }
}
```

### Common Error Messages
- `"Database connection failed: ..."` - MySQL connection issues
- `"Elasticsearch ping failed: ..."` - Elasticsearch connectivity issues
- `"Cache set operation failed: ..."` - Cache service issues
- `"Event processor not initialized"` - Service initialization issues

## Performance Considerations

### Response Time Monitoring
- **Target**: < 100ms for simple health check
- **Target**: < 500ms for comprehensive health check
- **Alert**: > 1000ms response time

### Frequency Recommendations
- **Simple Health Check**: Every 30 seconds for load balancers
- **Comprehensive Health Check**: Every 5 minutes for monitoring systems
- **High-Frequency Monitoring**: Every 10 seconds for critical systems

## Security Considerations

### Access Control
- Health check endpoints are public by default
- Consider implementing authentication for production environments
- Use network-level access controls for sensitive deployments

### Information Disclosure
- Health check responses may contain system information
- Consider what information is exposed in error messages
- Implement proper logging and monitoring

## Troubleshooting

### Common Issues

1. **Database Connection Failures**
   - Check MySQL service status
   - Verify connection pool configuration
   - Check network connectivity

2. **Elasticsearch Connection Failures**
   - Verify Elasticsearch service is running
   - Check cluster health status
   - Verify network connectivity

3. **Cache Operation Failures**
   - Check cache service status
   - Verify memory availability
   - Check cache configuration

4. **High Response Times**
   - Monitor system resources (CPU, memory, disk)
   - Check network latency
   - Review service configurations

### Debugging Commands

```bash
# Test simple health check
curl -s http://localhost:8081/api/health/simple

# Test comprehensive health check
curl -s http://localhost:8081/api/health | jq .

# Test individual service health
curl -s http://localhost:8081/api/health | jq '.services.database'
curl -s http://localhost:8081/api/health | jq '.services.elasticsearch'

# Monitor response times
curl -w "Response time: %{time_total}s\n" -s http://localhost:8081/api/health
```

## API Documentation

The health check endpoints are documented in the Swagger UI at:
`http://localhost:8081/swagger/index.html`

Navigate to the "health" tag to view detailed API documentation including:
- Request/response schemas
- Example responses
- Error codes
- Testing interface 