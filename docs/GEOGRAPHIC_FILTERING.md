# Geographic Filtering

## Overview

The firewall system now supports automatic geographic filtering using MaxMind's GeoLite2 database. This feature allows the system to automatically determine the country of origin for IP addresses and apply country-based rules without requiring manual country code input.

## How It Works

### Automatic Geolocation

When a filter request is made with an IP address but no country code, the system automatically:

1. **Geolocates the IP** using MaxMind's GeoLite2-Country database
2. **Applies country rules** using the existing country filtering system
3. **Returns the result** based on the country's status (allowed/denied/whitelisted)

### Manual Override

Users can still provide a country code manually, which will be used instead of geolocation:

- **When `country` field is empty**: System auto-geolocates the IP
- **When `country` field is provided**: System uses the provided country code

## Implementation Details

### File Structure

```
services/
├── geolocation.go          # MaxMind GeoIP integration
└── filter_service.go       # Enhanced filter logic

main.go                     # GeoIP initialization
```

### Key Components

#### 1. GeoIP Service (`services/geolocation.go`)

**Functions:**
- `InitGeoIP()` - Initializes MaxMind database reader
- `GetCountryFromIP(ipStr string)` - Resolves IP to country code
- `GetCountryFromIPWithFallback(ipStr string)` - Safe geolocation with error handling
- `IsPrivateIP(ip net.IP)` - Detects private/local IP addresses

**Features:**
- **Private IP Detection**: Automatically skips geolocation for private IPs
- **Error Handling**: Graceful fallback when geolocation fails
- **Performance**: Fast local database lookups
- **Privacy**: No external API calls

#### 2. Enhanced Filter Service (`services/filter_service.go`)

**Modified Function:**
- `EvaluateFilters()` - Now includes automatic geolocation logic

**Logic Flow:**
```go
if country == "" && ip != "" {
    country = GetCountryFromIPWithFallback(ip)
    // Log geolocation result for debugging
}
```

#### 3. Application Initialization (`main.go`)

**Startup Process:**
```go
// Initialize GeoIP service
if err := services.InitGeoIP(); err != nil {
    log.Printf("Warning: GeoIP initialization failed: %v", err)
    log.Println("Geographic filtering will be disabled")
} else {
    log.Println("GeoIP service initialized successfully")
}
defer services.CloseGeoIP()
```

## Usage Examples

### API Requests

#### 1. Automatic Geolocation
```bash
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "91.67.0.1",
    "email": "",
    "userAgent": "",
    "country": "",
    "username": ""
  }'
```

**Response:**
```json
{
  "result": "denied",
  "reason": "country denied",
  "field": "country",
  "value": "DE"
}
```

#### 2. Manual Country Override
```bash
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "8.8.8.8",
    "email": "",
    "userAgent": "",
    "country": "DE",
    "username": ""
  }'
```

**Response:**
```json
{
  "result": "denied",
  "reason": "country denied",
  "field": "country",
  "value": "DE"
}
```

#### 3. Private IP Handling
```bash
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "192.168.1.1",
    "email": "",
    "userAgent": "",
    "country": "",
    "username": ""
  }'
```

**Response:**
```json
{
  "result": "denied",
  "reason": "ip denied",
  "field": "ip",
  "value": "192.168.1.1"
}
```

## Configuration

### Required Files

#### GeoLite2-Country.mmdb
- **Location**: Root directory of the application
- **Size**: ~9MB
- **License**: Free for basic use (MaxMind GeoLite2)
- **Updates**: Monthly database updates available

### Dependencies

```go
require (
    github.com/oschwald/geoip2-golang v1.13.0
)
```

## Private IP Ranges

The system automatically skips geolocation for these IP ranges:

- **10.0.0.0/8** - Class A private
- **172.16.0.0/12** - Class B private  
- **192.168.0.0/16** - Class C private
- **127.0.0.0/8** - Loopback
- **169.254.0.0/16** - Link-local
- **::1/128** - IPv6 loopback
- **fe80::/10** - IPv6 link-local
- **fc00::/7** - IPv6 unique local

## Error Handling

### Geolocation Failures

When geolocation fails, the system:

1. **Logs the error** for debugging
2. **Continues processing** with IP-only rules
3. **Returns "allowed"** if no IP rules match
4. **Does not fail the request**

### Database Issues

If the GeoLite2 database is missing or corrupted:

1. **Logs a warning** during startup
2. **Continues operation** without geographic filtering
3. **All other functionality** remains intact

## Performance Considerations

### Optimization Features

- **Singleton Reader**: One GeoIP reader instance shared across requests
- **Local Database**: No network calls required
- **Fast Lookups**: Optimized binary database format
- **Memory Efficient**: Streams data as needed

### Expected Performance

- **Geolocation Speed**: ~1-5ms per IP lookup
- **Memory Usage**: ~50MB for database in memory
- **Concurrent Requests**: Thread-safe implementation

## Performance

### Caching
The system implements intelligent caching for geolocation lookups:

- **Country Lookups**: Cached for 1 hour with automatic cleanup
- **ASN Lookups**: Cached for 1 hour with automatic cleanup
- **Cache Keys**: Prefixed with "country:" or "asn:" to avoid conflicts
- **Cache Misses**: Empty results are also cached to avoid repeated lookups for invalid IPs
- **Memory Management**: Automatic cleanup every 10 minutes removes expired entries

### Cache Statistics
Cache performance can be monitored via the health check endpoint:

```bash
curl http://localhost:8081/api/health
```

Response includes:
```json
{
  "geo_cache": {
    "enabled": true,
    "items_count": 1250
  }
}
```

### Cache Management
- **Automatic Expiration**: 1 hour TTL for all geolocation data
- **LRU Eviction**: Least recently used items are removed when memory is low
- **Manual Clear**: Cache can be cleared programmatically if needed

## Integration with Existing Features

### Country Rules

The geographic filtering integrates seamlessly with existing country rules:

- **Same Status Logic**: Uses existing allowed/denied/whitelisted
- **Existing UI**: Country management interface unchanged
- **Event System**: Country rule changes trigger re-evaluation
- **Conflict Detection**: Works with existing conflict detection

### Filter Pipeline

Geographic filtering works alongside all other filters:

1. **IP Filter** - Direct IP/CIDR matching
2. **Email Filter** - Email pattern matching
3. **User Agent Filter** - User agent pattern matching
4. **Country Filter** - Geographic filtering (auto or manual)
5. **Username Filter** - Username pattern matching

## Monitoring and Logging

### Debug Information

The system logs geolocation results for debugging:

```
Auto-geolocated IP 91.67.0.1 to country: DE
Geolocation failed for IP 192.168.1.1: private IP address: 192.168.1.1
```

### Application Logs

Startup logs indicate GeoIP status:

```
GeoIP service initialized successfully
Warning: GeoIP initialization failed: GeoIP database not found
```

## Troubleshooting

### Common Issues

#### 1. Database Not Found
**Error**: `GeoIP database not found at GeoLite2-Country.mmdb`
**Solution**: Download GeoLite2-Country.mmdb to the root directory

#### 2. Invalid IP Address
**Error**: `invalid IP address: 256.256.256.256`
**Solution**: Ensure valid IP format in requests

#### 3. Private IP Geolocation
**Error**: `private IP address: 192.168.1.1`
**Expected**: Private IPs are skipped for geolocation

### Verification Commands

#### Test Geolocation Service
```bash
go run scripts/test_geolocation.go
```

#### Test Filter Endpoint
```bash
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{"ip": "8.8.8.8", "country": ""}'
```

## Security Considerations

### Privacy
- **No External Calls**: All lookups are local
- **No Data Transmission**: IP addresses stay within the system
- **Compliance**: Meets GDPR and privacy requirements

### Access Control
- **Database Location**: GeoLite2-Country.mmdb in application root
- **File Permissions**: Ensure appropriate read permissions
- **Backup Strategy**: Include database in backup procedures

## Future Enhancements

### Potential Improvements

1. **Caching Layer**: Redis cache for frequent lookups
2. **Database Updates**: Automatic monthly database updates
3. **Regional Filtering**: City/region-level filtering
4. **ISP Filtering**: ISP-based filtering capabilities
5. **Threat Intelligence**: Integration with threat feeds

### Configuration Options

Future versions could include:

- **Database Path**: Configurable database location
- **Update Frequency**: Configurable update intervals
- **Cache TTL**: Configurable cache time-to-live
- **Logging Level**: Configurable geolocation logging

## API Reference

### Filter Endpoint

**URL**: `POST /api/filter`

**Request Body**:
```json
{
  "ip": "string",           // IP address (optional)
  "email": "string",        // Email address (optional)
  "userAgent": "string",    // User agent (optional)
  "country": "string",      // Country code (optional, auto-geolocated if empty)
  "username": "string"      // Username (optional)
}
```

**Response**:
```json
{
  "result": "string",       // "allowed", "denied", "whitelisted", "error"
  "reason": "string",       // Reason for the result
  "field": "string",        // Field that triggered the result
  "value": "string"         // Value that triggered the result
}
```

### Country Codes

The system uses ISO 3166-1 alpha-2 country codes (e.g., "US", "DE", "CN").

## License Information

- **MaxMind GeoLite2**: Free for basic use, commercial license for advanced features
- **Database License**: Creative Commons Attribution-ShareAlike 4.0 International License
- **Usage Terms**: See MaxMind's GeoLite2 license for full terms

## Support

For issues with geographic filtering:

1. **Check Logs**: Review application logs for GeoIP errors
2. **Verify Database**: Ensure GeoLite2-Country.mmdb is present
3. **Test Geolocation**: Use the test script to verify functionality
4. **Check Permissions**: Ensure database file is readable 