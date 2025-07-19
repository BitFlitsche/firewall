# Enhanced Charset Detection System

## Overview

The firewall now includes an enhanced charset detection system that can identify and filter content based on Unicode scripts and character sets. This system supports a wide range of languages and scripts from around the world.

## Supported Scripts and Languages

### Basic Scripts
- **ASCII**: Basic English characters (0-127)
- **Latin**: Extended Latin characters (European languages)
- **Vietnamese**: Vietnamese with diacritics

### Cyrillic Scripts
- **Cyrillic**: Russian, Ukrainian, Bulgarian, Serbian, Belarusian, Macedonian, etc.

### Middle Eastern Scripts
- **Arabic**: Arabic, Persian (Farsi), Urdu, Pashto, etc.
- **Hebrew**: Hebrew, Yiddish

### European Scripts
- **Greek**: Greek language

### South Asian Scripts
- **Devanagari**: Hindi, Sanskrit, Marathi, Nepali, etc.
- **Bengali**: Bengali, Assamese
- **Tamil**: Tamil
- **Telugu**: Telugu
- **Kannada**: Kannada
- **Malayalam**: Malayalam
- **Gujarati**: Gujarati
- **Gurmukhi**: Punjabi
- **Oriya**: Odia
- **Sinhala**: Sinhala

### Southeast Asian Scripts
- **Thai**: Thai
- **Lao**: Lao
- **Khmer**: Khmer (Cambodian)
- **Myanmar**: Burmese

### East Asian Scripts
- **Chinese**: Chinese (Simplified & Traditional)
- **Japanese**: Japanese (Hiragana, Katakana, Kanji)
- **Korean**: Korean (Hangul)

### Other Scripts
- **Armenian**: Armenian
- **Georgian**: Georgian
- **Ethiopic**: Amharic, Tigrinya, etc.
- **Mongolian**: Mongolian
- **Tibetan**: Tibetan

### Special Categories
- **Mixed**: Content with multiple scripts
- **UTF-8**: Valid UTF-8 encoded text
- **Other**: Unrecognized characters

## How It Works

### Detection Algorithm

1. **Character Analysis**: Each character in the input is analyzed for its Unicode script
2. **Script Counting**: Characters are counted by their script type
3. **Dominant Script Detection**: The script with the highest count is identified
4. **Threshold Check**: If the dominant script has >50% of characters, it's used
5. **Mixed Content**: If no single script dominates, marked as "Mixed"

### Detection Logic

```go
func detectCharset(s string) string {
    // Count characters by script
    scriptCounts := make(map[string]int)
    totalChars := 0
    
    for _, r := range s {
        totalChars++
        script := getUnicodeScript(r)
        scriptCounts[script]++
    }
    
    // Find dominant script
    var dominantScript string
    maxCount := 0
    for script, count := range scriptCounts {
        if count > maxCount {
            maxCount = count
            dominantScript = script
        }
    }
    
    // Apply threshold
    if float64(maxCount)/float64(totalChars) > 0.5 {
        return dominantScript
    }
    
    return "Mixed"
}
```

### Unicode Range Support

The system supports Unicode ranges for:

- **Basic Latin**: 0x0020-0x007F
- **Latin Extended**: 0x0080-0x024F
- **Cyrillic**: 0x0400-0x04FF, 0x0500-0x052F
- **Arabic**: 0x0600-0x06FF, 0x0750-0x077F, 0xFB50-0xFDFF, 0xFE70-0xFEFF
- **Hebrew**: 0x0590-0x05FF
- **Greek**: 0x0370-0x03FF, 0x1F00-0x1FFF
- **Thai**: 0x0E00-0x0E7F
- **Devanagari**: 0x0900-0x097F
- **Bengali**: 0x0980-0x09FF
- **Tamil**: 0x0B80-0x0BFF
- **Telugu**: 0x0C00-0x0C7F
- **Kannada**: 0x0C80-0x0CFF
- **Malayalam**: 0x0D00-0x0D7F
- **Gujarati**: 0x0A80-0x0AFF
- **Gurmukhi**: 0x0A00-0x0A7F
- **Oriya**: 0x0B00-0x0B7F
- **Sinhala**: 0x0D80-0x0DFF
- **Chinese**: 0x4E00-0x9FFF, 0x3400-0x4DBF, 0x20000-0x2A6DF
- **Japanese**: 0x3040-0x309F, 0x30A0-0x30FF, 0x31F0-0x31FF
- **Korean**: 0xAC00-0xD7AF, 0x1100-0x11FF, 0x3130-0x318F, 0xA960-0xA97F, 0xD7B0-0xD7FF
- **Vietnamese**: 0x1EA0-0x1EFF
- **Armenian**: 0x0530-0x058F
- **Georgian**: 0x10A0-0x10FF
- **Ethiopic**: 0x1200-0x137F
- **Mongolian**: 0x1800-0x18AF
- **Tibetan**: 0x0F00-0x0FFF
- **Khmer**: 0x1780-0x17FF
- **Lao**: 0x0E80-0x0EFF
- **Myanmar**: 0x1000-0x109F

## Usage

### API Endpoint

The charset detection is automatically applied to the `/api/filter` endpoint for these fields:
- `email`
- `user_agent`
- `content`
- `username`

### Example Requests

```bash
# Cyrillic content
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"username": "Привет"}'

# Arabic content
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"username": "مرحبا"}'

# Chinese content
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"username": "你好"}'

# Mixed content
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"username": "Hello 你好 Привет"}'
```

### Response Format

```json
{
  "result": "denied|allowed|whitelisted",
  "reason": "charset denied|charset whitelisted",
  "field": "username|email|user_agent|content",
  "value": "detected text"
}
```

## Database Schema

### Charset Rules Table

```sql
CREATE TABLE charset_rules (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    charset VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('allowed', 'denied', 'whitelisted')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Seeded Data

The system automatically seeds the charset_rules table with 30+ predefined charset rules:

- ASCII, Latin, Vietnamese
- Cyrillic, Arabic, Hebrew, Greek
- South Asian scripts (Devanagari, Bengali, Tamil, etc.)
- Southeast Asian scripts (Thai, Lao, Khmer, Myanmar)
- East Asian scripts (Chinese, Japanese, Korean)
- Other scripts (Armenian, Georgian, Ethiopic, etc.)
- Special categories (Mixed, UTF-8, Other)

## Configuration

### Default Status

All charset rules are seeded with `status: "allowed"` by default. Administrators can change the status to:
- `"denied"`: Block content in this script
- `"whitelisted"`: Always allow content in this script
- `"allowed"`: Normal processing (default)

### Management

Charset rules can be managed through the API:

```bash
# List all charset rules
GET /api/charsets

# Create new charset rule
POST /api/charset
{
  "charset": "NewScript",
  "status": "denied"
}

# Update charset rule
PUT /api/charset/{id}
{
  "charset": "UpdatedScript",
  "status": "whitelisted"
}

# Delete charset rule
DELETE /api/charset/{id}
```

## Performance Considerations

- **Caching**: Filter results are cached for 5 minutes
- **Batch Processing**: Charset rules are loaded once per request
- **Efficient Detection**: Uses character-by-character analysis with early termination
- **Memory Efficient**: Minimal memory footprint for script detection

## Security Features

- **Input Validation**: All input is validated before charset detection
- **SQL Injection Protection**: Uses parameterized queries
- **XSS Protection**: Proper content encoding
- **Rate Limiting**: Integrated with existing rate limiting system

## Monitoring and Analytics

Charset detection events are logged for:
- **Traffic Analysis**: Track charset usage patterns
- **Security Monitoring**: Monitor blocked charset attempts
- **Performance Metrics**: Response times and cache hit rates
- **Audit Trail**: Complete request/response logging

## Future Enhancements

- **Language Detection**: Identify specific languages within scripts
- **Machine Learning**: Adaptive charset detection based on patterns
- **Custom Scripts**: User-defined script detection rules
- **Real-time Updates**: Dynamic charset rule updates
- **Advanced Analytics**: Deep charset usage analytics 