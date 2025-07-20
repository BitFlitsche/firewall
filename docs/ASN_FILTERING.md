# ASN Filtering

## Overview

The firewall system now supports ASN (Autonomous System Number) filtering using MaxMind's GeoLite2-ASN database. This feature allows the system to automatically determine the ASN of an IP address and apply ASN-based rules without requiring manual ASN input.

## How It Works

### Automatic ASN Lookup

When a filter request is made with an IP address, the system automatically:

1. **Looks up the ASN** using MaxMind's GeoLite2-ASN database
2. **Applies ASN rules** using the existing ASN filtering system
3. **Returns the result** based on the ASN's status (allowed/denied/whitelisted)

### Manual ASN Override

Users can provide an ASN manually, which will be used instead of automatic lookup:

- **When `asn` field is empty**: System auto-looks up the ASN from the IP
- **When `asn` field is provided**: System uses the provided ASN

### ASN Database Requirements

The system requires the `GeoLite2-ASN.mmdb` file to be present in the root directory. This file can be downloaded from MaxMind's website.

### Spamhaus ASN-DROP Integration

The system supports importing ASN data from Spamhaus ASN-DROP list, which provides a comprehensive list of ASNs associated with spam and malicious activities.

**Features:**
- **Automatic Import**: Fetches data from Spamhaus ASN-DROP JSON endpoint
- **Data Mapping**: Maps Spamhaus fields to internal database structure
- **Source Tracking**: Tracks imported data with "spamhaus" source identifier
- **Conflict Resolution**: Replaces existing Spamhaus records on re-import
- **Manual Entry Protection**: Preserves manually created/edited entries (source = "manual")
- **Elasticsearch Sync**: Automatically syncs imported data to search index
- **Statistics**: Provides import statistics and sync tracking

**Spamhaus Data Format:**
The Spamhaus ASN-DROP endpoint returns data in JSONL (JSON Lines) format, where each line contains a separate JSON object. The system automatically parses this format and extracts individual ASN records.

**Spamhaus Data Fields:**
- **ASN**: Autonomous System Number (e.g., 7922)
- **RIR**: Regional Internet Registry (e.g., "arin", "ripencc")
- **Domain**: Associated domain name
- **CC**: Country code (ISO 3166-1 alpha-2)
- **ASName**: ASN name/description

**Data Processing:**
- Each line is parsed as a separate JSON object
- Invalid lines are skipped with warnings
- Metadata lines (starting with `{"type":`) are automatically filtered
- Only records with valid ASN numbers are imported
- Manual entries (source = "manual") are preserved and not overwritten

### Automatic Scheduling

The system automatically runs Spamhaus imports daily at midnight with the following features:

**Scheduling:**
- **Daily Import**: Runs automatically at 00:00 (midnight) every day
- **Distributed Locking**: Ensures only one instance runs the import across multiple application instances
- **Lock TTL**: 30-minute timeout to prevent stuck locks
- **Graceful Handling**: Skips import if another instance is already running

**Manual Triggers:**
- **Force Import**: Immediate import with distributed locking
- **Status Monitoring**: Check if import is currently running
- **Next Schedule**: View next scheduled import time

**Configuration:**
- **Lock Name**: `spamhaus_import`
- **Lock TTL**: Configurable via `spamhaus.import_lock_ttl` (default: 30 minutes)
- **Schedule**: Configurable via `spamhaus.import_schedule` (default: daily at midnight)
- **Auto Import**: Configurable via `spamhaus.auto_import_enabled` (default: false)
- **Import URL**: Configurable via `spamhaus.import_url`
- **Retry Logic**: Automatic retry on next schedule if import fails

**Configuration Example:**
```yaml
spamhaus:
  auto_import_enabled: true  # Enable automatic imports
  import_schedule: "0 0 * * *"  # Daily at midnight (cron format)
  import_lock_ttl: "30m"  # 30-minute lock timeout
  import_url: "https://www.spamhaus.org/drop/asndrop.json"
```

## Implementation Details

### File Structure

```
services/
├── geolocation.go          # Enhanced with ASN lookup functionality
└── filter_service.go       # Enhanced with ASN filtering

models/
└── models.go              # ASN model definition

controllers/
└── rules.go               # ASN CRUD controllers

routes/
└── routes.go              # ASN API routes

firewall-app/src/components/
└── ASNForm.js             # ASN management UI
```

### Key Components

#### 1. ASN Model (`models/models.go`)

**Structure:**
```go
type ASN struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    ASN       string    `gorm:"unique;not null;type:varchar(20)" json:"asn"`
    RIR       string    `gorm:"type:varchar(20)" json:"rir"`                    // Optional
    Domain    string    `gorm:"type:varchar(255)" json:"domain"`                // Optional
    Country   string    `gorm:"type:varchar(2)" json:"cc"`                      // Optional
    Name      string    `gorm:"not null;type:varchar(255)" json:"asname"`
    Status    string    `gorm:"not null;type:varchar(20)" json:"status"`
    Source    string    `gorm:"type:varchar(50)" json:"source"`                 // Optional
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

**Field Details:**
- **ASN** (required): ASN number in format "AS12345"
- **RIR** (optional): Regional Internet Registry (e.g., "arin", "ripencc", "apnic")
- **Domain** (optional): Domain name associated with the ASN
- **Country** (optional): ISO 3166-1 alpha-2 country code (e.g., "US", "GB", "DE")
- **Name** (required): ASN name/description
- **Status** (required): "denied", "allowed", or "whitelisted"
- **Source** (optional): Source of the ASN data (e.g., "spamhaus", "manual")

**Features:**
- **ASN Format**: Must start with "AS" followed by numbers (e.g., "AS12345")
- **Name Field**: Description of the ASN (e.g., "Comcast Cable Communications")
- **Status**: "denied", "allowed", or "whitelisted"

#### 2. Enhanced Geolocation Service (`services/geolocation.go`)

**New Functions:**
- `InitASN()` - Initializes MaxMind ASN database reader
- `GetASNFromIP(ipStr string)` - Resolves IP to ASN
- `GetASNFromIPWithFallback(ipStr string)` - Safe ASN lookup with error handling

**Features:**
- **Private IP Detection**: Automatically skips ASN lookup for private IPs
- **Error Handling**: Graceful fallback when ASN lookup fails
- **Performance**: Fast local database lookups
- **Privacy**: No external API calls

#### 3. Enhanced Filter Service (`services/filter_service.go`)

**Modified Function:**
- `EvaluateFilters()` - Now includes ASN filtering
- `filterASN()` - New ASN filter function

**Logic Flow:**
```go
// When IP is provided
1. Look up ASN using MaxMind database
2. Query ASN rules in Elasticsearch
3. Apply status (denied/allowed/whitelisted)
4. Return result
```

#### 4. ASN CRUD Controllers (`controllers/rules.go`)

**Endpoints:**
- `POST /api/asn` - Create ASN rule
- `GET /api/asns` - List ASN rules with pagination/filtering
- `PUT /api/asn/:id` - Update ASN rule
- `DELETE /api/asn/:id` - Delete ASN rule
- `GET /api/asns/stats` - Get ASN statistics
- `POST /api/asns/recreate-index` - Recreate Elasticsearch index

#### 5. ASN Management UI (`firewall-app/src/components/ASNForm.js`)

**Features:**
- **Form Validation**: Ensures ASN format (AS + numbers)
- **Real-time Filtering**: Search by ASN number or name
- **Status Filtering**: Filter by allowed/denied/whitelisted
- **Sorting**: Sort by ID, ASN, name, or status
- **Pagination**: Server-side pagination for large datasets
- **Edit/Delete**: Full CRUD operations

## API Reference

### ASN Management Endpoints

#### Import Spamhaus ASN-DROP Data
```bash
POST /api/asns/import-spamhaus
```

**Description:** Imports ASN data from Spamhaus ASN-DROP list. This endpoint:
- Fetches data from https://www.spamhaus.org/drop/asndrop.json
- Replaces existing Spamhaus records in the database
- Syncs imported data to Elasticsearch
- Updates sync tracking information

**Response:**
```json
{
  "message": "Spamhaus ASN-DROP data imported successfully"
}
```

#### Get Spamhaus Import Statistics
```bash
GET /api/asns/spamhaus-stats
```

**Description:** Returns statistics about Spamhaus import including:
- Total number of Spamhaus ASNs in database
- Last sync timestamp

**Response:**
```json
{
  "total_spamhaus_asns": 401,
  "last_sync": "2025-01-20T10:30:00Z"
}
```

#### Force Spamhaus Import
```bash

```

**Description:** Triggers an immediate Spamhaus import with distributed locking to ensure only one instance runs the import.

**Response:**
```json
{
  "message": "Spamhaus import triggered successfully"
}
```

#### Get Spamhaus Import Status
```bash
GET /api/asns/spamhaus-status
```

**Description:** Returns the current status of Spamhaus import including:
- Whether import is currently running
- Next scheduled import time
- Time until next scheduled import

**Response:**
```json
{
  "is_running": false,
  "next_scheduled": "2025-01-21 00:00:00",
  "next_scheduled_relative": "8h30m15s"
}
```

#### Create ASN Rule
```bash
POST /api/asn
Content-Type: application/json

{
  "asn": "AS7922",
  "name": "Comcast Cable Communications",
  "status": "denied"
}
```

#### List ASN Rules
```bash
GET /api/asns?page=1&limit=10&status=denied&search=comcast&orderBy=asn&order=asc
```

#### Update ASN Rule
```bash
PUT /api/asn/1
Content-Type: application/json

{
  "asn": "AS7922",
  "name": "Comcast Cable Communications LLC",
  "status": "whitelisted"
}
```

#### Delete ASN Rule
```bash
DELETE /api/asn/1
```

#### Get ASN Statistics
```bash
GET /api/asns/stats
```

### Filter Endpoint (Enhanced)

The existing filter endpoint now includes ASN filtering:

```bash
POST /api/filter
Content-Type: application/json

{
  "ip": "68.85.108.1",
  "email": "",
  "user_agent": "",
  "country": "",
  "asn": "",
  "username": ""
}
```

**Response Examples:**

#### Denied by ASN
```json
{
  "result": "denied",
  "reason": "asn denied",
  "field": "asn",
  "value": "AS7922"
}
```

#### Allowed (No ASN Rule)
```json
{
  "result": "allowed",
  "field": "asn",
  "value": "AS15169"
}
```

## Usage Examples

### cURL Examples

#### Create ASN Rule
```bash
curl -X POST http://localhost:8081/api/asn \
  -H "Content-Type: application/json" \
  -d '{
    "asn": "AS7922",
    "name": "Comcast Cable Communications",
    "status": "denied"
  }'
```

#### Test ASN Filtering (Automatic Lookup)
```bash
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "68.85.108.1",
    "email": "",
    "user_agent": "",
    "country": "",
    "asn": "",
    "username": ""
  }'
```

#### Test ASN Filtering (Manual Override)
```bash
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "8.8.8.8",
    "email": "",
    "user_agent": "",
    "country": "",
    "asn": "AS7922",
    "username": ""
  }'
```

#### Import Spamhaus ASN-DROP Data
```bash
curl -X POST http://localhost:8081/api/asns/import-spamhaus
```

#### Get Spamhaus Import Statistics
```bash
curl -X GET http://localhost:8081/api/asns/spamhaus-stats
```

### JavaScript Examples

#### Using fetch() (Automatic ASN Lookup)
```javascript
const response = await fetch('/api/filter', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    ip: '68.85.108.1',
    email: '',
    userAgent: '',
    country: '',
    asn: '',
    username: ''
  })
});

const result = await response.json();
console.log(result);
// {result: "denied", reason: "asn denied", field: "asn", value: "AS7922"}
```

#### Using fetch() (Manual ASN Override)
```javascript
const response = await fetch('/api/filter', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    ip: '8.8.8.8',
    email: '',
    userAgent: '',
    country: '',
    asn: 'AS7922',
    username: ''
  })
});

const result = await response.json();
console.log(result);
// {result: "denied", reason: "asn denied", field: "asn", value: "AS7922"}
```

## Database Schema

### ASN Table
```sql
CREATE TABLE asns (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    asn VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_asn_status (status),
    INDEX idx_asn_asn (asn),
    INDEX idx_asn_status_asn (status, asn),
    INDEX idx_asn_asn_status (asn, status),
    INDEX idx_asn_status_id (status, id),
    INDEX idx_asn_status_asn_id (status, asn, id)
);
```

## Elasticsearch Index

### ASN Index Mapping
```json
{
  "mappings": {
    "properties": {
      "id": {"type": "long"},
      "asn": {"type": "keyword"},
      "name": {"type": "text"},
      "status": {"type": "keyword"}
    }
  }
}
```

## Configuration

### Required Files
- `GeoLite2-ASN.mmdb` - MaxMind ASN database (must be in root directory)

### Environment Variables
No additional environment variables are required. The system will automatically detect if the ASN database is available.

## Error Handling

### ASN Database Not Found
If the ASN database is not available:
- System logs a warning but continues to function
- ASN filtering is disabled
- Other filters continue to work normally

### Invalid ASN Format
When creating ASN rules:
- ASN must start with "AS" followed by numbers
- Validation error is returned for invalid format

### ASN Lookup Failures
When looking up ASN for an IP:
- Private IPs are skipped (no error)
- Invalid IPs return error
- Database lookup failures are logged but don't fail the request

## Performance Considerations

### Database Indexes
- Optimized indexes for filtering and sorting
- Composite indexes for common query patterns

### Elasticsearch
- Keyword fields for exact matching
- Text fields for search functionality
- Optimized mappings for fast queries

### Caching
- Filter results are cached for 5 minutes
- Reduces database and Elasticsearch load

## Security Considerations

### Input Validation
- ASN format validation (AS + numbers)
- SQL injection prevention through parameterized queries
- XSS prevention through proper escaping

### Access Control
- Currently no authentication required
- Consider implementing proper authentication for production

## Monitoring and Logging

### Log Messages
- ASN database initialization status
- ASN lookup successes and failures
- Filter application results

### Metrics
- ASN filter response times
- Cache hit rates
- Error rates for ASN lookups

## Troubleshooting

### Common Issues

1. **ASN Database Not Found**
   - Ensure `GeoLite2-ASN.mmdb` is in the root directory
   - Check file permissions

2. **ASN Lookup Fails**
   - Verify database file integrity
   - Check logs for specific error messages

3. **Filter Not Working**
   - Verify ASN rules exist in database
   - Check Elasticsearch index status
   - Review filter logs

### Debug Commands

```bash
# Check ASN database status
curl http://localhost:8081/api/health

# Test ASN lookup
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{"ip": "8.8.8.8"}'

# List ASN rules
curl http://localhost:8081/api/asns
```

## Future Enhancements

### Planned Features
- ASN geolocation (country/region from ASN)
- ASN reputation scoring
- Bulk ASN import/export
- ASN change detection and alerts

### Potential Improvements
- Real-time ASN database updates
- Integration with external ASN reputation services
- Advanced ASN analytics and reporting 