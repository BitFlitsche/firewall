# API Reference

## Base URL

All API endpoints are prefixed with `/api`

## Authentication

Currently, the API does not require authentication. In production, consider implementing proper authentication.

## Common Response Format

All API responses follow this format:

```json
{
  "result": "string",       // "allowed", "denied", "whitelisted", "error"
  "reason": "string",       // Reason for the result (optional)
  "field": "string",        // Field that triggered the result (optional)
  "value": "string"         // Value that triggered the result (optional)
}
```

## Filter Endpoint

### POST /api/filter

Evaluates a request against all firewall rules and returns the result.

**Request Body:**
```json
{
  "ip": "string",           // IP address (required) - sufficient for filtering
  "email": "string",        // Email address (optional)
  "userAgent": "string",    // User agent string (optional)
  "country": "string",      // Country code (optional, auto-geolocated if empty)
  "username": "string"      // Username (optional)
}
```

**Response Examples:**

#### Allowed Request
```json
{
  "result": "allowed"
}
```

#### Denied by IP
```json
{
  "result": "denied",
  "reason": "ip denied",
  "field": "ip",
  "value": "192.168.1.1"
}
```

#### Denied by Country (Auto-geolocated)
```json
{
  "result": "denied",
  "reason": "country denied",
  "field": "country",
  "value": "DE"
}
```

#### Whitelisted
```json
{
  "result": "whitelisted",
  "reason": "ip whitelisted",
  "field": "ip",
  "value": "10.0.0.1"
}
```

#### Error
```json
{
  "result": "error",
  "reason": "elasticsearch error",
  "field": "ip",
  "value": "invalid-ip"
}
```

## Geographic Filtering

### Automatic Geolocation

When the `country` field is empty and an `ip` is provided, the system automatically:

1. **Geolocates the IP** using MaxMind's GeoLite2-Country database
2. **Applies country rules** using the existing country filtering system
3. **Returns the result** based on the country's status

### Manual Country Override

When the `country` field is provided, the system uses the provided country code instead of geolocation.

### Private IP Handling

Private/local IP addresses are automatically skipped for geolocation and processed using IP rules only.

**Private IP Ranges:**
- `10.0.0.0/8` - Class A private
- `172.16.0.0/12` - Class B private
- `192.168.0.0/16` - Class C private
- `127.0.0.0/8` - Loopback
- `169.254.0.0/16` - Link-local
- IPv6 equivalents

## Usage Examples

### cURL Examples

#### Automatic Geolocation
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

#### Manual Country Override
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

#### Private IP (No Geolocation)
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

### JavaScript Examples

#### Using fetch()
```javascript
const response = await fetch('/api/filter', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    ip: '91.67.0.1',
    email: '',
    userAgent: '',
    country: '',
    username: ''
  })
});

const result = await response.json();
console.log(result);
// {result: "denied", reason: "country denied", field: "country", value: "DE"}
```

#### Using axios
```javascript
import axios from 'axios';

const response = await axios.post('/api/filter', {
  ip: '91.67.0.1',
  email: '',
  userAgent: '',
  country: '',
  username: ''
});

console.log(response.data);
// {result: "denied", reason: "country denied", field: "country", value: "DE"}
```

## Error Handling

### Common Error Scenarios

#### Geolocation Failed
```json
{
  "result": "allowed"
}
```
*Note: When geolocation fails, the system continues with IP-only rules and returns "allowed" if no rules match.*

#### Invalid IP Address
```json
{
  "result": "error",
  "reason": "elasticsearch error",
  "field": "ip",
  "value": "invalid-ip"
}
```

#### Database Connection Error
```json
{
  "result": "error",
  "reason": "elasticsearch error",
  "field": "ip",
  "value": "91.67.0.1"
}
```

### HTTP Status Codes

- `200 OK` - Request processed successfully
- `400 Bad Request` - Invalid request format
- `500 Internal Server Error` - Server error

## Performance Considerations

### Response Times

- **Typical response time**: 10-50ms
- **Geolocation lookup**: 1-5ms
- **Elasticsearch query**: 5-20ms
- **Database query**: 2-10ms

### Rate Limiting

The API includes rate limiting to prevent abuse:

- **Default limit**: 1000 requests per minute
- **Configurable**: Via `config.yaml`
- **Response headers**: Include rate limit information

### Caching

- **Geolocation results**: Not cached (fast local database)
- **Filter results**: Not cached (real-time evaluation)
- **Country rules**: Cached in Elasticsearch

## Country Codes

The system uses ISO 3166-1 alpha-2 country codes:

- `US` - United States
- `DE` - Germany
- `CN` - China
- `JP` - Japan
- `UK` - United Kingdom
- `CA` - Canada
- And many more...

For a complete list, see [ISO 3166-1](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2).

## Testing

### Test Endpoints

#### Health Check
```bash
curl http://localhost:8081/api/health
```

#### System Stats
```bash
curl http://localhost:8081/api/system-stats
```

### Test Data

You can test with these sample IPs:

- `8.8.8.8` - Google DNS (US)
- `1.1.1.1` - Cloudflare DNS (US)
- `91.67.0.1` - German IP (DE)
- `192.168.1.1` - Private IP (skipped for geolocation)

## Monitoring

### Logs

The application logs geolocation results for debugging:

```
Auto-geolocated IP 91.67.0.1 to country: DE
Geolocation failed for IP 192.168.1.1: private IP address: 192.168.1.1
```

### Metrics

Monitor these metrics for performance:

- Response times
- Geolocation success rate
- Error rates
- Rate limiting hits

## Security Considerations

### Input Validation

All inputs are validated:

- **IP addresses**: Valid IPv4/IPv6 format
- **Email addresses**: RFC 5321 compliant
- **Country codes**: ISO 3166-1 alpha-2 format
- **User agents**: Length and character limits

### Privacy

- **No external calls**: All geolocation is local
- **No data transmission**: IPs stay within your system
- **Compliance**: Meets GDPR and privacy requirements

## Support

For API issues:

1. Check the application logs
2. Verify the GeoLite2 database is present
3. Test with the health check endpoint
4. Review the detailed documentation

## Related Documentation

- [Geographic Filtering](GEOGRAPHIC_FILTERING.md) - Detailed geographic filtering guide
- [Validation](VALIDATION.md) - Input validation rules
- [Conflict Detection](CONFLICT_DETECTION.md) - IP/CIDR conflict detection
- [Health Check](HEALTH_CHECK.md) - Monitoring and health checks 

## Analytics Endpoints

### GET /api/analytics/logs

Returns paginated traffic logs with filtering and sorting capabilities.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Number of items per page (default: 50)
- `orderBy` (optional): Field to sort by (default: timestamp)
- `order` (optional): Sort direction - "asc" or "desc" (default: desc)
- `ip_address` (optional): Filter by IP address
- `email` (optional): Filter by email address
- `user_agent` (optional): Filter by user agent
- `username` (optional): Filter by username
- `country` (optional): Filter by country
- `asn` (optional): Filter by ASN
- `final_result` (optional): Filter by result (allowed, denied, whitelisted)

**Sortable Fields:**
- `timestamp` - Request timestamp
- `final_result` - Filter result
- `ip_address` - IP address
- `email` - Email address
- `user_agent` - User agent string
- `username` - Username
- `country` - Country code
- `asn` - ASN number
- `response_time_ms` - Response time in milliseconds
- `cache_hit` - Cache hit status

**Response:**
```json
{
  "logs": [
    {
      "id": 1,
      "timestamp": "2024-01-01T12:00:00Z",
      "final_result": "denied",
      "ip_address": "192.168.1.1",
      "email": "test@example.com",
      "user_agent": "Mozilla/5.0...",
      "username": "user123",
      "country": "US",
      "asn": "AS7922",
      "response_time_ms": 150,
      "cache_hit": false
    }
  ],
  "total": 1000,
  "page": 1,
  "limit": 50,
  "total_pages": 20
}
```

**Example Usage:**
```bash
# Sort by IP address ascending
curl "http://localhost:8081/api/analytics/logs?orderBy=ip_address&order=asc"

# Sort by response time descending
curl "http://localhost:8081/api/analytics/logs?orderBy=response_time_ms&order=desc"

# Filter and sort
curl "http://localhost:8081/api/analytics/logs?final_result=denied&orderBy=timestamp&order=desc"
``` 