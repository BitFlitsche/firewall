# Charset Filtering - Quick Reference

## 🚀 Quick Start

### 1. View Current Configuration
```bash
curl "http://localhost:8081/api/charset-fields"
```

### 2. Add Custom Field
```bash
curl -X POST "http://localhost:8081/api/charset-fields/add-custom" \
  -H "Content-Type: application/json" \
  -d '{"field_name": "content"}'
```

### 3. Toggle Standard Field
```bash
curl -X POST "http://localhost:8081/api/charset-fields/toggle-standard" \
  -H "Content-Type: application/json" \
  -d '{"field_name": "username"}'
```

### 4. Test Filter
```bash
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"content": "Привет мир"}'
```

## 📋 Common Operations

### Field Management
| Operation | Command |
|-----------|---------|
| **View Fields** | `GET /api/charset-fields` |
| **Add Custom** | `POST /api/charset-fields/add-custom` |
| **Toggle Standard** | `POST /api/charset-fields/toggle-standard` |
| **Delete Custom** | `DELETE /api/charset-fields/custom/:field` |

### Charset Rules
| Operation | Command |
|-----------|---------|
| **View Rules** | `GET /api/charsets` |
| **Add Rule** | `POST /api/charset` |
| **Update Rule** | `PUT /api/charset/:id` |
| **Delete Rule** | `DELETE /api/charset/:id` |
| **View Stats** | `GET /api/charsets/stats` |

### Cache Management
| Operation | Command |
|-----------|---------|
| **Clear Cache** | `POST /api/cache/flush` |

## 🎯 Common Use Cases

### Allow Only Latin Characters
```bash
# Deny all non-Latin charsets
curl -X POST "http://localhost:8081/api/charset" \
  -H "Content-Type: application/json" \
  -d '{"charset": "Cyrillic", "status": "denied"}'

curl -X POST "http://localhost:8081/api/charset" \
  -H "Content-Type: application/json" \
  -d '{"charset": "Chinese", "status": "denied"}'
```

### Check Custom Field
```bash
# Add content field
curl -X POST "http://localhost:8081/api/charset-fields/add-custom" \
  -H "Content-Type: application/json" \
  -d '{"field_name": "content"}'

# Test with Cyrillic content
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"content": "Привет мир"}'
```

### Disable Username Checking
```bash
# Toggle username field off
curl -X POST "http://localhost:8081/api/charset-fields/toggle-standard" \
  -H "Content-Type: application/json" \
  -d '{"field_name": "username"}'

# Test with Cyrillic username
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"username": "ПриветМир"}'
```

## 📊 Supported Charsets

### Major Scripts
- **Latin** - English, French, Spanish, etc.
- **Cyrillic** - Russian, Bulgarian, Serbian
- **Chinese** - Simplified & Traditional
- **Japanese** - Hiragana, Katakana, Kanji
- **Korean** - Hangul
- **Arabic** - Arabic alphabet
- **Hebrew** - Hebrew alphabet
- **Thai** - Thai alphabet
- **Greek** - Greek alphabet

### Special Categories
- **Mixed** - Multiple scripts in one field
- **Other** - Symbols, emojis, special characters

## 🔧 Frontend Usage

### 1. Access Charset Management
- Navigate to **Charset Management** page
- Expand **"Charset Filter Fields"** accordion

### 2. Manage Standard Fields
- Use checkboxes for: username, email, user_agent
- Can be enabled/disabled (cannot be deleted)

### 3. Manage Custom Fields
- Click **"Add Custom Field"** button
- Enter field name (e.g., content, description)
- Delete with trash icon

### 4. View Results
- Success/error messages appear immediately
- Changes take effect instantly
- Cache is automatically cleared

## 🚨 Troubleshooting

### Field Not Working?
```bash
# Check configuration
curl "http://localhost:8081/api/charset-fields"

# Clear cache
curl -X POST "http://localhost:8081/api/cache/flush"
```

### Unexpected Denials?
```bash
# Check charset rules
curl "http://localhost:8081/api/charsets"

# Test specific content
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"content": "test"}'
```

### Performance Issues?
- Reduce number of enabled fields
- Check cache hit rates
- Monitor server resources

## 📝 Examples

### Basic Latin Content
```bash
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"username": "john_doe", "email": "john@example.com"}'
# Result: allowed
```

### Cyrillic Content
```bash
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"content": "Привет мир"}'
# Result: denied (if Cyrillic is denied)
```

### Mixed Content
```bash
curl -X POST "http://localhost:8081/api/filter" \
  -H "Content-Type: application/json" \
  -d '{"content": "Hello 世界"}'
# Result: denied (if Chinese is denied)
```

---

*For detailed documentation, see [CHARSET_FILTERING.md](./CHARSET_FILTERING.md)* 