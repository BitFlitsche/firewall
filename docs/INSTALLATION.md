# Installation Guide

## Prerequisites

- **Go 1.24+** - Backend runtime
- **Node.js 18+** - Frontend build tools
- **MySQL 8.0+** - Database
- **Elasticsearch 8.0+** - Search engine (optional)
- **Redis** - Distributed locking (optional)

## Quick Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd firewall
```

### 2. Install Dependencies

#### Backend Dependencies
```bash
go mod download
```

#### Frontend Dependencies
```bash
cd firewall-app
npm install
cd ..
```

### 3. Configure the Application

```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your settings
```

### 4. Set Up the Database

```bash
# Create MySQL database
mysql -u root -p
CREATE DATABASE firewall;
CREATE USER 'firewall_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON firewall.* TO 'firewall_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 5. Download GeoLite2 Database (Required for Geographic Filtering)

```bash
# Download from MaxMind (requires free account)
wget https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=YOUR_LICENSE_KEY&suffix=tar.gz

# Extract to root directory
tar -xzf GeoLite2-Country_*.tar.gz
cp GeoLite2-Country_*/GeoLite2-Country.mmdb .
rm -rf GeoLite2-Country_*
```

**Alternative: Manual Download**
1. Visit [MaxMind GeoLite2](https://dev.maxmind.com/geoip/geoip2/geolite2/)
2. Create a free account
3. Download `GeoLite2-Country.mmdb`
4. Place in the root directory of the application

### 6. Build the Application

```bash
# Build backend
go build -o firewall .

# Build frontend
cd firewall-app
npm run build
cd ..
```

### 7. Run the Application

```bash
./firewall
```

The application will be available at `http://localhost:8081`

## Configuration

### Basic Configuration

Edit `config.yaml`:

```yaml
server:
  port: 8081
  host: "0.0.0.0"

database:
  host: "127.0.0.1"
  port: 3306
  user: "firewall_user"
  password: "your_password"
  name: "firewall"

elastic:
  hosts:
    - "http://localhost:9200"
```

### Geographic Filtering Configuration

The geographic filtering feature requires:

1. **GeoLite2-Country.mmdb** in the root directory
2. **Country rules** in the database (can be added via UI)
3. **No additional configuration** - works automatically

### Distributed Locking (Optional)

For multi-instance deployments:

```yaml
locking:
  distributed: true
redis:
  host: "localhost"
  port: 6379
```

## Verification

### 1. Check Application Health

```bash
curl http://localhost:8081/api/health
```

Expected response:
```json
{
  "status": "healthy",
  "services": {
    "database": {"status": "healthy"},
    "elasticsearch": {"status": "healthy"},
    "geoip": {"status": "healthy"}
  }
}
```

### 2. Test Geographic Filtering

```bash
# Test automatic geolocation
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{"ip": "8.8.8.8", "country": ""}'
```

Expected response:
```json
{
  "result": "allowed"
}
```

### 3. Check GeoIP Status

Look for these log messages during startup:

```
GeoIP service initialized successfully
```

If you see:
```
Warning: GeoIP initialization failed: GeoIP database not found
```

Then the `GeoLite2-Country.mmdb` file is missing or in the wrong location.

## Troubleshooting

### Common Issues

#### 1. GeoIP Database Not Found

**Error**: `GeoIP database not found at GeoLite2-Country.mmdb`

**Solution**:
```bash
# Check if file exists
ls -la GeoLite2-Country.mmdb

# Download if missing
wget https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=YOUR_LICENSE_KEY&suffix=tar.gz
```

#### 2. Database Connection Failed

**Error**: `Failed to connect to database`

**Solution**:
```bash
# Check MySQL status
sudo systemctl status mysql

# Verify credentials
mysql -u firewall_user -p firewall
```

#### 3. Elasticsearch Connection Failed

**Error**: `Failed to connect to Elasticsearch`

**Solution**:
```bash
# Check Elasticsearch status
curl http://localhost:9200

# Start if stopped
sudo systemctl start elasticsearch
```

#### 4. Port Already in Use

**Error**: `Address already in use`

**Solution**:
```bash
# Check what's using the port
lsof -i :8081

# Kill the process or change port in config.yaml
```

### Performance Issues

#### 1. Slow Geolocation

**Symptoms**: High response times for filter requests

**Solutions**:
- Ensure `GeoLite2-Country.mmdb` is on fast storage (SSD)
- Check system memory usage
- Consider upgrading hardware

#### 2. Database Slow

**Symptoms**: Slow API responses

**Solutions**:
- Optimize MySQL configuration
- Add database indexes
- Consider connection pooling settings

## Production Deployment

### 1. Environment Variables

Set production environment variables:

```bash
export FIREWALL_DATABASE_HOST=your-db-host
export FIREWALL_DATABASE_PASSWORD=your-secure-password
export FIREWALL_ELASTIC_HOSTS=http://your-es-host:9200
export FIREWALL_REDIS_HOST=your-redis-host
```

### 2. Systemd Service

Create `/etc/systemd/system/firewall.service`:

```ini
[Unit]
Description=Firewall Application
After=network.target

[Service]
Type=simple
User=firewall
WorkingDirectory=/opt/firewall
ExecStart=/opt/firewall/firewall
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable firewall
sudo systemctl start firewall
```

### 3. Reverse Proxy (Nginx)

Create `/etc/nginx/sites-available/firewall`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable:
```bash
sudo ln -s /etc/nginx/sites-available/firewall /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 4. SSL/TLS

Use Let's Encrypt for free SSL certificates:

```bash
sudo certbot --nginx -d your-domain.com
```

### 5. Monitoring

Set up monitoring with:

- **Prometheus** for metrics
- **Grafana** for dashboards
- **AlertManager** for alerts

## Security Considerations

### 1. Database Security

- Use strong passwords
- Limit database user privileges
- Enable SSL connections
- Regular backups

### 2. Network Security

- Use firewall rules
- Enable HTTPS
- Rate limiting
- IP whitelisting

### 3. Application Security

- Regular updates
- Security patches
- Input validation
- Error handling

## Backup Strategy

### 1. Database Backup

```bash
# Create backup script
#!/bin/bash
mysqldump -u firewall_user -p firewall > backup_$(date +%Y%m%d_%H%M%S).sql
```

### 2. Configuration Backup

```bash
# Backup configuration
cp config.yaml backup/config_$(date +%Y%m%d_%H%M%S).yaml
```

### 3. GeoIP Database Backup

```bash
# Backup GeoIP database
cp GeoLite2-Country.mmdb backup/GeoLite2-Country_$(date +%Y%m%d_%H%M%S).mmdb
```

## Updates

### 1. Application Updates

```bash
git pull origin main
go build -o firewall .
sudo systemctl restart firewall
```

### 2. GeoIP Database Updates

```bash
# Download new database
wget https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=YOUR_LICENSE_KEY&suffix=tar.gz

# Replace old database
tar -xzf GeoLite2-Country_*.tar.gz
cp GeoLite2-Country_*/GeoLite2-Country.mmdb .
rm -rf GeoLite2-Country_*

# Restart application
sudo systemctl restart firewall
```

## Support

For installation issues:

1. Check the logs: `sudo journalctl -u firewall -f`
2. Verify prerequisites are installed
3. Test each component individually
4. Review the troubleshooting section
5. Open an issue on GitHub

## Related Documentation

- [API Reference](API_REFERENCE.md) - Complete API documentation
- [Geographic Filtering](GEOGRAPHIC_FILTERING.md) - Geographic filtering guide
- [Configuration](CONFIGURATION.md) - Detailed configuration options
- [Monitoring](MONITORING.md) - Monitoring and alerting setup 