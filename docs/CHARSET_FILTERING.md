# Charset Filtering Documentation

## Overview

The charset filtering system allows you to detect and filter requests based on the Unicode scripts and character sets used in various request fields. This feature helps identify and control requests containing non-Latin characters, which can be useful for security and content filtering purposes.

## Features

### ✅ Dynamic Field Management
- **Standard Fields**: username, email, user_agent (can be enabled/disabled)
- **Custom Fields**: Any field name can be added for charset detection
- **Real-time Configuration**: Changes take effect immediately
- **Cache Management**: Automatic cache clearing when fields are modified

### ✅ Unicode Script Detection
- **30+ Unicode Scripts**: Supports Latin, Cyrillic, Arabic, Chinese, Japanese, Korean, and more
- **Mixed Content Detection**: Identifies content with multiple scripts
- **"Other" Category**: Catches symbols, emojis, and special characters

### ✅ Flexible Configuration
- **Field Selection**: Choose which request fields to check
- **Status Management**: Allow, deny, or whitelist specific charsets
- **Custom Fields**: Add any field name for charset detection

## Architecture

### Backend Components

#### 1. Charset Detection Engine (`controllers/filter.go`)
```go
func detectCharset(s string) string {
    // Analyzes Unicode scripts in the input string
    // Returns the dominant script or "Mixed" for multiple scripts
}

func getUnicodeScript(r rune) string {
    // Maps individual runes to Unicode script names
    // Supports 30+ scripts including Latin, Cyrillic, Arabic, etc.
}
```

#### 2. Field Configuration Service (`services/charset_fields_service.go`)
```go
type CharsetField struct {
    Name    string `json:"name"`
    Enabled bool   `json:"enabled"`
    Type    string `json:"type"` // "standard" or "custom"
}

type CharsetFieldsConfig struct {
    standardFields []CharsetField
    customFields   []CharsetField
}
```

#### 3. Cache Management
- **Automatic Cache Clearing**: Cache is cleared when field configuration changes
- **5-minute Cache Duration**: Filter results are cached for performance
- **Selective Clearing**: Only charset-related cache items are cleared

### Frontend Components

#### 1. Charset Management Page
- **Charset Rules Table**: View, add, edit, delete charset rules
- **Field Configuration**: Collapsible accordion for field management
- **Real-time Updates**: Changes are saved immediately

#### 2. Custom Fields Manager
- **Standard Fields**: Checkboxes for username, email, user_agent
- **Custom Fields**: Add/remove custom field names
- **Visual Feedback**: Success/error messages for all operations

## Configuration

### Standard Fields

The system comes with three standard fields that can be enabled or disabled:

| Field | Default | Description |
|-------|---------|-------------|
| `username` | ✅ Enabled | User's login name |
| `email` | ✅ Enabled | Email address |
| `user_agent` | ✅ Enabled | Browser/application identifier |

### Adding Custom Fields

You can add any field name for charset detection:

```bash
# Add a custom field
curl -X POST "http://localhost:8081/api/charset-fields/add-custom" \
  -H "Content-Type: application/json" \
  -d '{"field_name": "content"}'
```

Common custom field examples:
- `content` - Request body content
- `description` - Item descriptions
- `notes` - User notes
- `comment` - User comments
- `title` - Page or item titles

### Field Management

#### Enable/Disable Standard Fields
```bash
# Toggle a standard field
curl -X POST "http://localhost:8081/api/charset-fields/toggle-standard" \
  -H "Content-Type: application/json" \
  -d '{"field_name": "email"}'
```

#### Add Custom Fields
```bash
# Add a new custom field
curl -X POST "http://localhost:8081/api/charset-fields/add-custom" \
  -H "Content-Type: application/json" \
  -d '{"field_name": "description"}'
```

#### Delete Custom Fields
```bash
# Delete a custom field
curl -X DELETE "http://localhost:8081/api/charset-fields/custom/description"
```

#### View Current Configuration
```bash
# Get current field configuration
curl "http://localhost:8081/api/charset-fields"
```

Response:
```json
{
  "standard_fields": [
    {"name": "username", "enabled": true, "type": "standard"},
    {"name": "email", "enabled": false, "type": "standard"},
    {"name": "user_agent", "enabled": true, "type": "standard"}
  ],
  "custom_fields": [
    {"name": "content", "enabled": true, "type": "custom"}
  ]
}
```

## Charset Rules

### Supported Charsets

The system detects and supports the following Unicode scripts:

#### Latin Scripts
- **ASCII** - Basic Latin characters (A-Z, a-z, 0-9)
- **Latin** - Extended Latin with diacritics
- **Latin Extended** - Additional Latin characters

#### European Scripts
- **Cyrillic** - Russian, Bulgarian, Serbian, etc.
- **Greek** - Greek alphabet
- **Armenian** - Armenian alphabet
- **Georgian** - Georgian alphabet

#### Asian Scripts
- **Chinese** - Simplified and Traditional Chinese
- **Japanese** - Hiragana, Katakana, Kanji
- **Korean** - Hangul characters
- **Thai** - Thai alphabet
- **Vietnamese** - Latin with Vietnamese diacritics

#### Middle Eastern Scripts
- **Arabic** - Arabic alphabet
- **Hebrew** - Hebrew alphabet

#### South Asian Scripts
- **Devanagari** - Hindi, Marathi, etc.
- **Bengali** - Bengali alphabet
- **Tamil** - Tamil alphabet
- **Telugu** - Telugu alphabet
- **Kannada** - Kannada alphabet
- **Malayalam** - Malayalam alphabet
- **Gujarati** - Gujarati alphabet
- **Gurmukhi** - Punjabi alphabet
- **Oriya** - Odia alphabet

#### Southeast Asian Scripts
- **Khmer** - Cambodian alphabet
- **Lao** - Lao alphabet
- **Myanmar** - Burmese alphabet
- **Thai** - Thai alphabet

#### Other Scripts
- **Ethiopic** - Amharic, Tigrinya, etc.
- **Mongolian** - Mongolian alphabet
- **Tibetan** - Tibetan alphabet
- **Sinhala** - Sinhalese alphabet

#### Special Categories
- **Mixed** - Content with multiple scripts
- **Other** - Symbols, emojis, special characters

### Managing Charset Rules

#### View All Charset Rules
```bash
curl "http://localhost:8081/api/charsets"
```

#### Add a Charset Rule
```bash
curl -X POST "http://localhost:8081/api/charset" \
  -H "Content-Type: application/json" \
  -d '{"charset": "Cyrillic", "status": "denied"}'
```

#### Update a Charset Rule
```bash
curl -X PUT "http://localhost:8081/api/charset/1" \
  -H "Content-Type: application/json" \
  -d '{"charset": "Cyrillic", "status": "allowed"}'
```

#### Delete a Charset Rule
```bash
curl -X DELETE "http://localhost:8081/api/charset/1"
```

## Usage Examples

### Basic Filtering

#### Test with Latin Content
```bash
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"username": "john_doe", "email": "john@example.com"}'
```

#### Test with Cyrillic Content
```bash
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"username": "Привет", "content": "Hello world"}'
```

#### Test with Mixed Content
```bash
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"content": "Hello 世界", "user_agent": "Mozilla/5.0"}'
```

### Field-Specific Filtering

#### Disable Username Charset Checking
1. Go to Charset Management page
2. Expand "Charset Filter Fields" section
3. Uncheck "username" checkbox
4. Test with Cyrillic username:
```bash
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"username": "ПриветМир"}'
# Result: allowed (charset checking disabled for username)
```

#### Enable Custom Field Checking
1. Add custom field "content":
```bash
curl -X POST "http://localhost:8081/api/charset-fields/add-custom" \
  -H "Content-Type: application/json" \
  -d '{"field_name": "content"}'
```

2. Test with Cyrillic content:
```bash
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"content": "Привет мир"}'
# Result: denied (if Cyrillic is set to denied)
```

## API Reference

### Charset Fields Management

#### GET /api/charset-fields
Get current field configuration.

**Response:**
```json
{
  "standard_fields": [
    {"name": "username", "enabled": true, "type": "standard"},
    {"name": "email", "enabled": true, "type": "standard"},
    {"name": "user_agent", "enabled": true, "type": "standard"}
  ],
  "custom_fields": [
    {"name": "content", "enabled": true, "type": "custom"}
  ]
}
```

#### POST /api/charset-fields/toggle-standard
Toggle a standard field (enable/disable).

**Request:**
```json
{"field_name": "email"}
```

**Response:**
```json
{"message": "Field toggled successfully"}
```

#### POST /api/charset-fields/add-custom
Add a new custom field.

**Request:**
```json
{"field_name": "description"}
```

**Response:**
```json
{"message": "Custom field added successfully"}
```

#### DELETE /api/charset-fields/custom/:field
Delete a custom field.

**Response:**
```json
{"message": "Custom field deleted successfully"}
```

#### POST /api/charset-fields/toggle-custom
Toggle a custom field (enable/disable).

**Request:**
```json
{"field_name": "content"}
```

**Response:**
```json
{"message": "Custom field toggled successfully"}
```

### Charset Rules Management

#### GET /api/charsets
Get all charset rules with pagination and filtering.

**Query Parameters:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)
- `status` - Filter by status (allowed, denied, whitelisted)
- `search` - Search in charset names
- `orderBy` - Sort field (charset, status, created_at)
- `order` - Sort direction (asc, desc)

#### POST /api/charset
Create a new charset rule.

**Request:**
```json
{"charset": "Cyrillic", "status": "denied"}
```

#### PUT /api/charset/:id
Update a charset rule.

**Request:**
```json
{"charset": "Cyrillic", "status": "allowed"}
```

#### DELETE /api/charset/:id
Delete a charset rule.

#### GET /api/charsets/stats
Get charset statistics.

**Response:**
```json
{
  "total": 35,
  "allowed": 20,
  "denied": 10,
  "whitelisted": 5
}
```

## Cache Management

### Automatic Cache Clearing

The system automatically clears relevant cache items when field configuration changes:

- **Field Toggle**: Cache cleared when standard/custom fields are enabled/disabled
- **Field Addition**: Cache cleared when new custom fields are added
- **Field Deletion**: Cache cleared when custom fields are removed

### Manual Cache Clearing

```bash
curl -X POST "http://localhost:8081/api/cache/flush"
```

**Response:**
```json
{
  "message": "Cache flushed successfully",
  "items_cleared": 5
}
```

## Best Practices

### 1. Field Selection
- **Start with Standard Fields**: Use username, email, user_agent for basic filtering
- **Add Custom Fields**: Add specific fields based on your application needs
- **Disable Unused Fields**: Disable fields you don't need to improve performance

### 2. Charset Rules
- **Whitelist Approach**: Start with all charsets allowed, then deny specific ones
- **Deny Approach**: Start with all charsets denied, then allow specific ones
- **Mixed Content**: Consider how to handle content with multiple scripts

### 3. Performance
- **Cache Duration**: 5-minute cache provides good balance of performance and freshness
- **Field Count**: Limit custom fields to only those you need
- **Regular Monitoring**: Check charset statistics regularly

### 4. Security
- **Logging**: Monitor charset detection results for security insights
- **Rate Limiting**: Consider rate limiting for filter requests
- **Validation**: Always validate input before charset detection

## Troubleshooting

### Common Issues

#### 1. Field Not Being Checked
**Problem**: Added custom field but it's not being checked for charset detection.

**Solution**: 
- Verify the field is enabled in the configuration
- Check if the field name matches exactly
- Clear cache to ensure new configuration is applied

#### 2. Cache Not Updating
**Problem**: Changes to field configuration not taking effect.

**Solution**:
- Cache is automatically cleared when fields are modified
- Manual cache clearing: `POST /api/cache/flush`
- Check server logs for cache clearing messages

#### 3. Unexpected Denials
**Problem**: Requests being denied unexpectedly.

**Solution**:
- Check charset rules for the detected script
- Verify field configuration (enabled/disabled)
- Test with different content to isolate the issue

#### 4. Performance Issues
**Problem**: Slow filter response times.

**Solution**:
- Reduce number of enabled fields
- Check cache hit rates
- Monitor server resources

### Debug Information

#### Check Field Configuration
```bash
curl "http://localhost:8081/api/charset-fields"
```

#### Check Charset Rules
```bash
curl "http://localhost:8081/api/charsets"
```

#### Test Specific Content
```bash
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"content": "test content"}'
```

## Related Documentation

- [Charset Detection System](./CHARSET_DETECTION.md) - Detailed charset detection algorithm
- [Filter API](./FILTER_API.md) - Complete filter endpoint documentation
- [Cache Management](./CACHE_MANAGEMENT.md) - Cache system documentation
- [Field Configuration](./FIELD_CONFIGURATION.md) - Field management guide

---

*Last updated: July 2024* 