# Input Validation System

The firewall application implements a comprehensive input validation system to ensure data integrity, security, and prevent common attack vectors.

## Overview

The validation system provides:
- **Field-level validation** for all input types
- **Comprehensive error reporting** with detailed messages
- **Security validation** to prevent injection attacks
- **Format validation** for emails, IPs, countries, etc.
- **Length and range validation** for all fields
- **Duplicate detection** for unique fields

## Validation Components

### 1. Validation Package (`validation/validation.go`)

The core validation package provides:

#### **Validation Functions**
- `ValidateIP(ip string)` - Validates IP addresses and CIDR blocks
- `ValidateEmail(email string)` - Validates email addresses
- `ValidateUserAgent(userAgent string)` - Validates user agent strings
- `ValidateCountry(country string)` - Validates country codes
- `ValidateUsername(username string)` - Validates usernames
- `ValidateContent(content string)` - Validates content for charset detection
- `ValidateCharset(charset string)` - Validates charset names
- `ValidateStatus(status string)` - Validates status values
- `ValidateRegex(pattern string)` - Validates regex patterns
- `ValidatePagination(page, limit string)` - Validates pagination parameters
- `ValidateOrderBy(orderBy, order string, validFields []string)` - Validates ordering
- `ValidateSearch(search string)` - Validates search parameters
- `ValidateID(id string)` - Validates ID parameters
- `ValidateFilterRequest(...)` - Validates filter requests
- `ValidateCreateRequest(...)` - Validates create requests

#### **Validation Results**
```go
type ValidationResult struct {
    IsValid bool              `json:"is_valid"`
    Errors  []ValidationError `json:"errors,omitempty"`
}

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   string `json:"value,omitempty"`
}
```

### 2. Model Validation Tags

All models include validation tags using Gin's binding system:

#### **IP Model**
```go
type IP struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Address   string    `gorm:"unique;not null;type:varchar(45)" json:"address" binding:"required"`
    Status    string    `gorm:"not null;type:varchar(20)" json:"status" binding:"required,oneof=allowed denied whitelisted"`
    IsCIDR    bool      `gorm:"default:false;type:boolean" json:"is_cidr"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

#### **Email Model**
```go
type Email struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Address   string    `gorm:"unique;not null;type:varchar(254)" json:"address" binding:"required,email"`
    Status    string    `gorm:"not null;type:varchar(20)" json:"status" binding:"required,oneof=allowed denied whitelisted"`
    IsRegex   bool      `gorm:"default:false;type:boolean" json:"is_regex"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

#### **Filter Request**
```go
type FilterRequest struct {
    IP        string `json:"ip" binding:"omitempty,max=45"`
    Email     string `json:"email" binding:"omitempty,email,max=254"`
    UserAgent string `json:"user_agent" binding:"omitempty,max=500"`
    Country   string `json:"country" binding:"omitempty,len=2,alpha"`
    Content   string `json:"content" binding:"omitempty,max=10000"`
    Username  string `json:"username" binding:"omitempty,max=100"`
}
```

## Validation Rules

### IP Address Validation
- **Format**: Valid IPv4, IPv6, or CIDR notation
- **Length**: Maximum 45 characters (IPv6 + CIDR)
- **CIDR**: Valid network/mask format
- **Single IP**: Valid IP address format

**Examples:**
```json
// Valid
"192.168.1.1"
"2001:db8::1"
"192.168.1.0/24"
"2001:db8::/32"

// Invalid
"256.1.2.3"
"192.168.1.1/33"
"not-an-ip"
```

### Email Validation
- **Format**: RFC 5321 compliant email format
- **Length**: Maximum 254 characters
- **Local part**: Valid characters, no consecutive dots
- **Domain**: Valid domain format

**Examples:**
```json
// Valid
"user@example.com"
"user.name@domain.co.uk"
"user+tag@example.com"

// Invalid
"user..name@example.com"
"user@.example.com"
"user@example."
"@example.com"
```

### User Agent Validation
- **Length**: Maximum 500 characters
- **Characters**: No null bytes or invalid control characters
- **Format**: Standard user agent string format

**Examples:**
```json
// Valid
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
"curl/7.68.0"

// Invalid
"User\x00Agent"  // Contains null byte
"Very long user agent..." // > 500 characters
```

### Country Code Validation
- **Length**: Exactly 2 characters
- **Format**: ISO 3166-1 alpha-2 format
- **Characters**: Alphabetic only

**Examples:**
```json
// Valid
"US"
"DE"
"GB"

// Invalid
"USA"     // Too long
"1"       // Too short
"12"      // Not alphabetic
```

### Status Validation
- **Values**: `allowed`, `denied`, `whitelisted`
- **Case**: Case-sensitive

**Examples:**
```json
// Valid
"allowed"
"denied"
"whitelisted"

// Invalid
"ALLOWED"     // Wrong case
"blocked"     // Invalid value
"permitted"   // Invalid value
```

### Regex Pattern Validation
- **Format**: Valid Go regex pattern
- **Length**: Maximum 500 characters
- **Compilation**: Must compile successfully

**Examples:**
```json
// Valid
"^admin.*$"
"[a-zA-Z0-9]+"
"\\d{3}-\\d{3}-\\d{4}"

// Invalid
"^admin[+"     // Unclosed bracket
"[a-z"         // Unclosed bracket
"\\invalid"    // Invalid escape
```

## API Validation

### Create Endpoints

All create endpoints (`POST /api/ip`, `POST /api/email`, etc.) include:

1. **JSON Binding Validation**: Validates required fields and formats
2. **Custom Field Validation**: Validates specific field formats
3. **Duplicate Detection**: Checks for existing records
4. **Status Validation**: Validates status values

**Example Response (Success):**
```json
{
  "id": 1,
  "address": "192.168.1.1",
  "status": "denied",
  "is_cidr": false,
  "created_at": "2025-07-19T08:34:17Z",
  "updated_at": "2025-07-19T08:34:17Z"
}
```

**Example Response (Validation Error):**
```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "address",
      "message": "Invalid IP address format",
      "value": "256.1.2.3"
    },
    {
      "field": "status",
      "message": "Invalid status (must be 'allowed', 'denied', or 'whitelisted')",
      "value": "blocked"
    }
  ]
}
```

### Update Endpoints

All update endpoints (`PUT /api/ip/:id`, etc.) include:

1. **ID Validation**: Validates the ID parameter
2. **Record Existence**: Checks if the record exists
3. **Input Validation**: Same as create endpoints
4. **Conflict Detection**: Checks for conflicts with other records

### List Endpoints

All list endpoints (`GET /api/ips`, etc.) include:

1. **Pagination Validation**: Validates page and limit parameters
2. **Status Validation**: Validates status filter if provided
3. **Search Validation**: Validates search parameters
4. **Ordering Validation**: Validates orderBy and order parameters

**Example Query Parameters:**
```
GET /api/ips?page=1&limit=10&status=denied&search=192.168&orderBy=Address&order=asc
```

### Filter Endpoint

The filter endpoint (`POST /api/filter`) includes:

1. **JSON Binding**: Validates request structure
2. **Field Validation**: Validates each provided field
3. **Minimum Fields**: Ensures at least one field is provided
4. **Content Validation**: Validates content for charset detection

## Security Validation

### SQL Injection Prevention
- **Search Parameters**: Validates against SQL injection patterns
- **Input Sanitization**: Removes dangerous characters
- **Parameter Binding**: Uses prepared statements

### XSS Prevention
- **Content Validation**: Validates content for malicious patterns
- **Character Encoding**: Ensures proper UTF-8 encoding
- **Length Limits**: Prevents oversized inputs

### Input Length Limits
- **IP Address**: 45 characters max
- **Email**: 254 characters max
- **User Agent**: 500 characters max
- **Country Code**: 2 characters exact
- **Username**: 100 characters max
- **Content**: 10,000 characters max
- **Charset**: 100 characters max

## Error Handling

### Validation Error Format
All validation errors follow a consistent format:

```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "field_name",
      "message": "Human-readable error message",
      "value": "The invalid value (if applicable)"
    }
  ]
}
```

### HTTP Status Codes
- **200 OK**: Validation passed, operation successful
- **400 Bad Request**: Validation failed
- **404 Not Found**: Record not found (for updates/deletes)
- **409 Conflict**: Duplicate record exists

### Error Categories
1. **Format Errors**: Invalid data format
2. **Length Errors**: Data too long or too short
3. **Value Errors**: Invalid values (e.g., status)
4. **Conflict Errors**: Duplicate records
5. **Security Errors**: Potentially malicious input

## Testing Validation

### Manual Testing
```bash
# Test valid IP creation
curl -X POST http://localhost:8081/api/ip \
  -H "Content-Type: application/json" \
  -d '{"address": "192.168.1.1", "status": "denied"}'

# Test invalid IP creation
curl -X POST http://localhost:8081/api/ip \
  -H "Content-Type: application/json" \
  -d '{"address": "256.1.2.3", "status": "blocked"}'

# Test filter with validation
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{"ip": "192.168.1.1", "email": "user@example.com"}'
```

### Unit Testing
```go
func TestValidateIP(t *testing.T) {
    tests := []struct {
        input    string
        expected bool
    }{
        {"192.168.1.1", true},
        {"256.1.2.3", false},
        {"192.168.1.0/24", true},
        {"192.168.1.0/33", false},
    }
    
    for _, test := range tests {
        result := validation.ValidateIP(test.input)
        if result.IsValid != test.expected {
            t.Errorf("ValidateIP(%s) = %v, expected %v", test.input, result.IsValid, test.expected)
        }
    }
}
```

## Best Practices

### For Developers
1. **Always validate input** before processing
2. **Use the validation package** for consistent validation
3. **Provide clear error messages** to users
4. **Log validation failures** for monitoring
5. **Test edge cases** thoroughly

### For API Consumers
1. **Validate input client-side** for better UX
2. **Handle validation errors** gracefully
3. **Use appropriate HTTP status codes**
4. **Provide user-friendly error messages**
5. **Retry with corrected input**

## Monitoring and Logging

### Validation Metrics
- **Validation Success Rate**: Track validation pass/fail rates
- **Common Validation Errors**: Identify frequent validation issues
- **Field-specific Errors**: Monitor which fields fail most often
- **Performance Impact**: Monitor validation processing time

### Logging
```go
// Log validation failures
if !validationResult.IsValid {
    log.Printf("Validation failed for %s: %+v", entityType, validationResult.Errors)
}
```

## Future Enhancements

### Planned Features
1. **Custom Validation Rules**: Allow custom validation rules
2. **Validation Caching**: Cache validation results for performance
3. **Async Validation**: Validate large datasets asynchronously
4. **Validation Plugins**: Extensible validation system
5. **Real-time Validation**: Client-side validation feedback

### Integration Opportunities
1. **OpenAPI/Swagger**: Enhanced API documentation
2. **GraphQL**: Schema-based validation
3. **Message Queues**: Async validation processing
4. **Microservices**: Distributed validation
5. **API Gateway**: Gateway-level validation 