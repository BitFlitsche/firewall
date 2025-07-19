# IP/CIDR Conflict Detection System

The firewall application implements a comprehensive conflict detection system to prevent overlapping IP addresses and CIDR ranges, ensuring data integrity and preventing rule conflicts.

## Overview

The conflict detection system provides:
- **Automatic conflict detection** when adding IPs or CIDR ranges
- **Comprehensive conflict reporting** with detailed information
- **Conflict severity levels** (error, warning, info)
- **Pre-creation conflict checking** via dedicated endpoint
- **Update conflict detection** when modifying existing entries

## Conflict Types

### 1. IP Address Conflicts

#### **IP in CIDR Range**
When adding an individual IP address that is already covered by an existing CIDR range.

**Example:**
```
Existing: 192.168.1.0/24 (covers 192.168.1.0 - 192.168.1.255)
New IP: 192.168.1.100
Result: Conflict detected - IP is already covered by CIDR range
```

#### **Exact IP Match**
When adding an IP address that already exists as an individual entry.

**Example:**
```
Existing: 192.168.1.100
New IP: 192.168.1.100
Result: Conflict detected - IP already exists
```

### 2. CIDR Range Conflicts

#### **CIDR Covers Existing IPs**
When adding a CIDR range that would cover existing individual IP addresses.

**Example:**
```
Existing: 192.168.1.100, 192.168.1.200
New CIDR: 192.168.1.0/24
Result: Warning - CIDR would cover existing IPs
```

#### **CIDR Overlap**
When adding a CIDR range that overlaps with existing CIDR ranges.

**Example:**
```
Existing: 192.168.1.0/24
New CIDR: 192.168.1.128/25
Result: Conflict detected - CIDR ranges overlap
```

#### **Exact CIDR Match**
When adding a CIDR range that already exists.

**Example:**
```
Existing: 192.168.1.0/24
New CIDR: 192.168.1.0/24
Result: Conflict detected - CIDR already exists
```

## API Endpoints

### 1. Create IP/CIDR with Conflict Detection

**Endpoint**: `POST /api/ip`

**Description**: Creates a new IP address or CIDR range with automatic conflict detection.

**Request Body**:
```json
{
  "address": "192.168.1.100",
  "status": "denied",
  "is_cidr": false
}
```

**Response (No Conflicts)**:
```json
{
  "id": 1,
  "address": "192.168.1.100",
  "status": "denied",
  "is_cidr": false,
  "created_at": "2025-07-19T08:34:17Z",
  "updated_at": "2025-07-19T08:34:17Z"
}
```

**Response (With Conflicts)**:
```json
{
  "error": "IP/CIDR conflicts detected",
  "conflicts": [
    {
      "type": "ip_in_cidr",
      "message": "IP 192.168.1.100 is already covered by CIDR range 192.168.1.0/24",
      "conflicting": ["192.168.1.0/24"],
      "severity": "error"
    }
  ],
  "message": "Please review conflicts before proceeding"
}
```

### 2. Check Conflicts (Pre-Creation)

**Endpoint**: `POST /api/ip/check-conflicts`

**Description**: Checks for conflicts without creating the entry.

**Request Body**:
```json
{
  "address": "192.168.1.0/24",
  "status": "denied",
  "is_cidr": true
}
```

**Response (No Conflicts)**:
```json
{
  "status": "clean",
  "conflicts": [],
  "conflict_count": 0,
  "can_proceed": true,
  "message": "No conflicts detected"
}
```

**Response (With Warnings)**:
```json
{
  "status": "warning",
  "conflicts": [
    {
      "type": "cidr_covers_ip",
      "message": "CIDR range 192.168.1.0/24 would cover existing IP 192.168.1.100",
      "conflicting": ["192.168.1.100"],
      "severity": "warning"
    }
  ],
  "conflict_count": 1,
  "can_proceed": true,
  "message": "Warnings detected - review before proceeding"
}
```

**Response (With Errors)**:
```json
{
  "status": "error",
  "conflicts": [
    {
      "type": "exact_match",
      "message": "CIDR range 192.168.1.0/24 already exists",
      "conflicting": ["192.168.1.0/24"],
      "severity": "error"
    }
  ],
  "conflict_count": 1,
  "can_proceed": false,
  "message": "Conflicts detected - cannot proceed"
}
```

### 3. Update IP/CIDR with Conflict Detection

**Endpoint**: `PUT /api/ip/:id`

**Description**: Updates an existing IP address or CIDR range with conflict detection.

**Request Body**:
```json
{
  "address": "192.168.1.0/24",
  "status": "denied",
  "is_cidr": true
}
```

**Response**: Same format as create endpoint, with conflict detection excluding the current record.

## Conflict Severity Levels

### Error (Cannot Proceed)
- **Exact matches**: IP or CIDR already exists
- **IP in CIDR**: Individual IP covered by existing CIDR range
- **CIDR overlap**: New CIDR overlaps with existing CIDR ranges

### Warning (Can Proceed with Review)
- **CIDR covers IPs**: New CIDR would cover existing individual IPs
- **Potential conflicts**: Situations that might cause confusion

### Info (Informational)
- **General information**: Non-critical conflicts or suggestions

## Conflict Detection Logic

### For Individual IP Addresses
1. **Parse the IP**: Validate IP address format
2. **Check existing CIDR ranges**: See if IP is covered by any CIDR
3. **Check exact matches**: See if IP already exists as individual entry
4. **Report conflicts**: Return detailed conflict information

### For CIDR Ranges
1. **Parse the CIDR**: Validate CIDR notation
2. **Check existing IPs**: See if CIDR would cover any individual IPs
3. **Check existing CIDRs**: See if CIDR overlaps with existing ranges
4. **Check exact matches**: See if CIDR already exists
5. **Report conflicts**: Return detailed conflict information

## Conflict Resolution Strategies

### 1. Remove Conflicting Entries
```bash
# Remove individual IP that conflicts with new CIDR
DELETE /api/ip/123

# Remove overlapping CIDR range
DELETE /api/ip/456
```

### 2. Modify Existing Entries
```bash
# Change individual IP to different address
PUT /api/ip/123
{
  "address": "192.168.2.100",
  "status": "denied",
  "is_cidr": false
}
```

### 3. Use Smaller CIDR Ranges
```bash
# Instead of 192.168.1.0/24, use smaller ranges
POST /api/ip
{
  "address": "192.168.1.0/25",
  "status": "denied",
  "is_cidr": true
}
```

## Best Practices

### For Administrators
1. **Check conflicts before creating**: Use the check-conflicts endpoint
2. **Review warnings carefully**: Warnings indicate potential issues
3. **Plan CIDR ranges**: Design network ranges to minimize conflicts
4. **Document decisions**: Keep track of conflict resolution decisions

### For Developers
1. **Handle conflict responses**: Implement proper error handling
2. **Provide user feedback**: Show conflict details to users
3. **Allow conflict resolution**: Provide options to resolve conflicts
4. **Validate input**: Always validate before checking conflicts

### For API Consumers
1. **Check conflicts first**: Use pre-creation endpoint
2. **Handle different status codes**: 200 (OK), 409 (Conflict)
3. **Parse conflict details**: Extract specific conflict information
4. **Provide resolution options**: Allow users to resolve conflicts

## Testing Conflict Detection

### Test Cases

#### **IP in CIDR Conflict**
```bash
# Create CIDR range
curl -X POST http://localhost:8081/api/ip \
  -H "Content-Type: application/json" \
  -d '{"address": "192.168.1.0/24", "status": "denied", "is_cidr": true}'

# Try to add IP in range
curl -X POST http://localhost:8081/api/ip \
  -H "Content-Type: application/json" \
  -d '{"address": "192.168.1.100", "status": "denied", "is_cidr": false}'
```

#### **CIDR Overlap Conflict**
```bash
# Create first CIDR
curl -X POST http://localhost:8081/api/ip \
  -H "Content-Type: application/json" \
  -d '{"address": "192.168.1.0/24", "status": "denied", "is_cidr": true}'

# Try overlapping CIDR
curl -X POST http://localhost:8081/api/ip \
  -H "Content-Type: application/json" \
  -d '{"address": "192.168.1.128/25", "status": "denied", "is_cidr": true}'
```

#### **Pre-Creation Check**
```bash
# Check for conflicts without creating
curl -X POST http://localhost:8081/api/ip/check-conflicts \
  -H "Content-Type: application/json" \
  -d '{"address": "192.168.1.100", "status": "denied", "is_cidr": false}'
```

## Monitoring and Logging

### Conflict Metrics
- **Conflict frequency**: Track how often conflicts occur
- **Conflict types**: Monitor which types of conflicts are most common
- **Resolution time**: Track how long conflicts take to resolve
- **User behavior**: Monitor how users handle conflicts

### Logging
```go
// Log conflict detection
if len(conflicts) > 0 {
    log.Printf("Conflicts detected for %s: %+v", ip.Address, conflicts)
}
```

## Future Enhancements

### Planned Features
1. **Conflict suggestions**: Automatic suggestions for conflict resolution
2. **Bulk conflict checking**: Check multiple entries at once
3. **Conflict history**: Track conflict resolution history
4. **Conflict analytics**: Analyze conflict patterns
5. **Auto-resolution**: Automatic conflict resolution options

### Integration Opportunities
1. **Frontend validation**: Real-time conflict checking in UI
2. **API documentation**: Enhanced Swagger documentation
3. **Testing framework**: Automated conflict testing
4. **Monitoring dashboards**: Conflict monitoring and alerting
5. **Audit trails**: Track all conflict-related actions 