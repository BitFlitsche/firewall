# Connection Pooling Configuration

This document describes the connection pooling configuration for the firewall application.

## Environment Variables

The following environment variables can be used to configure database connection pooling:

### Database Connection
- `DB_HOST`: Database host (default: 127.0.0.1)
- `DB_PORT`: Database port (default: 3306)
- `DB_USER`: Database username (default: user)
- `DB_PASSWORD`: Database password (default: password)
- `DB_NAME`: Database name (default: firewall)

### Connection Pooling
- `DB_MAX_OPEN_CONNS`: Maximum number of open connections (default: 25)
- `DB_MAX_IDLE_CONNS`: Maximum number of idle connections (default: 5)
- `DB_CONN_MAX_LIFETIME`: Maximum lifetime of a connection (default: 5m)
- `DB_CONN_MAX_IDLE_TIME`: Maximum idle time of a connection (default: 1m)

## Configuration Examples

### Development Environment
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=user
export DB_PASSWORD=password
export DB_NAME=firewall
export DB_MAX_OPEN_CONNS=10
export DB_MAX_IDLE_CONNS=3
export DB_CONN_MAX_LIFETIME=5m
export DB_CONN_MAX_IDLE_TIME=1m
```

### Production Environment
```bash
export DB_HOST=production-db.example.com
export DB_PORT=3306
export DB_USER=firewall_user
export DB_PASSWORD=secure_password
export DB_NAME=firewall_prod
export DB_MAX_OPEN_CONNS=50
export DB_MAX_IDLE_CONNS=10
export DB_CONN_MAX_LIFETIME=10m
export DB_CONN_MAX_IDLE_TIME=2m
```

### High Load Environment
```bash
export DB_MAX_OPEN_CONNS=100
export DB_MAX_IDLE_CONNS=20
export DB_CONN_MAX_LIFETIME=15m
export DB_CONN_MAX_IDLE_TIME=5m
```

## Connection Pooling Benefits

1. **Performance**: Reuses connections instead of creating new ones for each request
2. **Resource Management**: Limits the number of concurrent connections
3. **Connection Lifecycle**: Automatically closes idle connections and recycles old ones
4. **Monitoring**: Provides detailed connection statistics

## Monitoring Connection Pool

The application provides connection pool statistics through the `/api/stats` endpoint:

```json
{
  "db_connections": {
    "max_open_connections": 25,
    "open_connections": 8,
    "in_use": 3,
    "idle": 5,
    "wait_count": 0,
    "wait_duration": "0s",
    "max_idle_closed": 12,
    "max_lifetime_closed": 5
  }
}
```

## Best Practices

1. **Max Open Connections**: Set to 2-3x the number of CPU cores for most applications
2. **Max Idle Connections**: Set to 25-50% of max open connections
3. **Connection Lifetime**: Set to 5-15 minutes depending on database stability
4. **Idle Time**: Set to 1-5 minutes to balance resource usage and responsiveness

## Troubleshooting

### Connection Exhaustion
If you see high wait counts, increase `DB_MAX_OPEN_CONNS`:
```bash
export DB_MAX_OPEN_CONNS=50
```

### Memory Usage
If memory usage is high, reduce `DB_MAX_IDLE_CONNS`:
```bash
export DB_MAX_IDLE_CONNS=3
```

### Connection Timeouts
If connections are timing out, increase `DB_CONN_MAX_LIFETIME`:
```bash
export DB_CONN_MAX_LIFETIME=10m
```

## Docker Compose Configuration

For Docker environments, add environment variables to your `docker-compose.yml`:

```yaml
services:
  app:
    environment:
      - DB_HOST=mariadb
      - DB_PORT=3306
      - DB_USER=user
      - DB_PASSWORD=password
      - DB_NAME=firewall
      - DB_MAX_OPEN_CONNS=25
      - DB_MAX_IDLE_CONNS=5
      - DB_CONN_MAX_LIFETIME=5m
      - DB_CONN_MAX_IDLE_TIME=1m
``` 