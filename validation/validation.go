package validation

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode/utf8"

	"firewall/utils"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}
}

// AddError adds a validation error to the result
func (vr *ValidationResult) AddError(field, message, value string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// ValidateIP validates an IP address or CIDR block
func ValidateIP(ip string) *ValidationResult {
	result := NewValidationResult()

	if ip == "" {
		result.AddError("ip", "IP address cannot be empty", "")
		return result
	}

	// Check if it's a CIDR block
	if strings.Contains(ip, "/") {
		if !utils.ValidateCIDR(ip) {
			result.AddError("ip", "Invalid CIDR notation", ip)
			return result
		}
	} else {
		// Check if it's a valid IP address
		if net.ParseIP(ip) == nil {
			result.AddError("ip", "Invalid IP address format", ip)
			return result
		}
	}

	// Check length (IPv6 can be up to 45 characters)
	if len(ip) > 45 {
		result.AddError("ip", "IP address too long (max 45 characters)", ip)
		return result
	}

	return result
}

// ValidateEmail validates an email address
func ValidateEmail(email string) *ValidationResult {
	result := NewValidationResult()

	if email == "" {
		result.AddError("email", "Email address cannot be empty", "")
		return result
	}

	// RFC 5321 max length
	if len(email) > 254 {
		result.AddError("email", "Email address too long (max 254 characters)", email)
		return result
	}

	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		result.AddError("email", "Invalid email address format", email)
		return result
	}

	// Check for common invalid patterns
	if strings.Contains(email, "..") {
		result.AddError("email", "Email address contains consecutive dots", email)
		return result
	}

	if strings.HasPrefix(email, ".") || strings.HasSuffix(email, ".") {
		result.AddError("email", "Email address cannot start or end with a dot", email)
		return result
	}

	return result
}

// ValidateUserAgent validates a user agent string
func ValidateUserAgent(userAgent string) *ValidationResult {
	result := NewValidationResult()

	if userAgent == "" {
		result.AddError("user_agent", "User agent cannot be empty", "")
		return result
	}

	// Check length (reasonable max length)
	if len(userAgent) > 500 {
		result.AddError("user_agent", "User agent too long (max 500 characters)", userAgent)
		return result
	}

	// Check for null bytes or control characters
	for _, r := range userAgent {
		if r < 32 && r != 9 && r != 10 && r != 13 { // Allow tab, newline, carriage return
			result.AddError("user_agent", "User agent contains invalid control characters", userAgent)
			return result
		}
	}

	return result
}

// ValidateCountry validates a country code
func ValidateCountry(country string) *ValidationResult {
	result := NewValidationResult()

	if country == "" {
		result.AddError("country", "Country code cannot be empty", "")
		return result
	}

	// Check length (ISO 3166-1 alpha-2 is exactly 2 characters)
	if len(country) != 2 {
		result.AddError("country", "Country code must be exactly 2 characters", country)
		return result
	}

	// Check if it's alphabetic
	if !regexp.MustCompile(`^[A-Za-z]{2}$`).MatchString(country) {
		result.AddError("country", "Country code must be alphabetic", country)
		return result
	}

	return result
}

// ValidateUsername validates a username
func ValidateUsername(username string) *ValidationResult {
	result := NewValidationResult()

	if username == "" {
		result.AddError("username", "Username cannot be empty", "")
		return result
	}

	// Check length
	if len(username) > 100 {
		result.AddError("username", "Username too long (max 100 characters)", username)
		return result
	}

	// Check for null bytes or control characters
	for _, r := range username {
		if r < 32 && r != 9 { // Allow tab
			result.AddError("username", "Username contains invalid control characters", username)
			return result
		}
	}

	return result
}

// ValidateContent validates content for charset detection
func ValidateContent(content string) *ValidationResult {
	result := NewValidationResult()

	if content == "" {
		result.AddError("content", "Content cannot be empty", "")
		return result
	}

	// Check length (reasonable max length)
	if len(content) > 10000 {
		result.AddError("content", "Content too long (max 10000 characters)", "")
		return result
	}

	// Check if it's valid UTF-8
	if !utf8.ValidString(content) {
		result.AddError("content", "Content contains invalid UTF-8 sequences", "")
		return result
	}

	return result
}

// ValidateCharset validates a charset name
func ValidateCharset(charset string) *ValidationResult {
	result := NewValidationResult()

	if charset == "" {
		result.AddError("charset", "Charset cannot be empty", "")
		return result
	}

	// Check length
	if len(charset) > 100 {
		result.AddError("charset", "Charset name too long (max 100 characters)", charset)
		return result
	}

	// Check for valid charset pattern
	if !regexp.MustCompile(`^[A-Za-z0-9\-_]+$`).MatchString(charset) {
		result.AddError("charset", "Charset name contains invalid characters", charset)
		return result
	}

	return result
}

// ValidateStatus validates a status value
func ValidateStatus(status string) *ValidationResult {
	result := NewValidationResult()

	if status == "" {
		result.AddError("status", "Status cannot be empty", "")
		return result
	}

	validStatuses := map[string]bool{
		"allowed":     true,
		"denied":      true,
		"whitelisted": true,
	}

	if !validStatuses[status] {
		result.AddError("status", "Invalid status (must be 'allowed', 'denied', or 'whitelisted')", status)
		return result
	}

	return result
}

// ValidateRegex validates a regex pattern
func ValidateRegex(pattern string) *ValidationResult {
	result := NewValidationResult()

	if pattern == "" {
		result.AddError("regex", "Regex pattern cannot be empty", "")
		return result
	}

	// Check length
	if len(pattern) > 500 {
		result.AddError("regex", "Regex pattern too long (max 500 characters)", pattern)
		return result
	}

	// Try to compile the regex
	_, err := regexp.Compile(pattern)
	if err != nil {
		result.AddError("regex", fmt.Sprintf("Invalid regex pattern: %s", err.Error()), pattern)
		return result
	}

	return result
}

// ValidatePagination validates pagination parameters
func ValidatePagination(page, limit string) *ValidationResult {
	result := NewValidationResult()

	// Validate page
	if page != "" {
		var pageNum int
		if _, err := fmt.Sscanf(page, "%d", &pageNum); err != nil {
			result.AddError("page", "Page must be a valid integer", page)
		} else if pageNum < 1 {
			result.AddError("page", "Page must be greater than 0", page)
		}
	}

	// Validate limit
	if limit != "" {
		var limitNum int
		if _, err := fmt.Sscanf(limit, "%d", &limitNum); err != nil {
			result.AddError("limit", "Limit must be a valid integer", limit)
		} else if limitNum < 1 {
			result.AddError("limit", "Limit must be greater than 0", limit)
		} else if limitNum > 1000 {
			result.AddError("limit", "Limit cannot exceed 1000", limit)
		}
	}

	return result
}

// ValidateOrderBy validates order by parameters
func ValidateOrderBy(orderBy, order string, validFields []string) *ValidationResult {
	result := NewValidationResult()

	// Validate orderBy field
	if orderBy != "" {
		valid := false
		for _, field := range validFields {
			if orderBy == field {
				valid = true
				break
			}
		}
		if !valid {
			result.AddError("orderBy", fmt.Sprintf("Invalid order by field (valid: %s)", strings.Join(validFields, ", ")), orderBy)
		}
	}

	// Validate order direction
	if order != "" && order != "asc" && order != "desc" {
		result.AddError("order", "Order must be 'asc' or 'desc'", order)
	}

	return result
}

// ValidateSearch validates search parameters
func ValidateSearch(search string) *ValidationResult {
	result := NewValidationResult()

	if search == "" {
		return result // Empty search is valid
	}

	// Check length
	if len(search) > 255 {
		result.AddError("search", "Search term too long (max 255 characters)", search)
		return result
	}

	// Check for SQL injection patterns (basic)
	sqlInjectionPatterns := []string{
		"';", "--", "/*", "*/", "union", "select", "insert", "update", "delete", "drop", "create",
	}

	searchLower := strings.ToLower(search)
	for _, pattern := range sqlInjectionPatterns {
		if strings.Contains(searchLower, pattern) {
			result.AddError("search", "Search term contains invalid characters", search)
			return result
		}
	}

	return result
}

// ValidateID validates an ID parameter
func ValidateID(id string) *ValidationResult {
	result := NewValidationResult()

	if id == "" {
		result.AddError("id", "ID cannot be empty", "")
		return result
	}

	var idNum uint
	if _, err := fmt.Sscanf(id, "%d", &idNum); err != nil {
		result.AddError("id", "ID must be a valid integer", id)
		return result
	}

	if idNum == 0 {
		result.AddError("id", "ID must be greater than 0", id)
		return result
	}

	return result
}

// ValidateFilterRequest validates a filter request
func ValidateFilterRequest(ip, email, userAgent, country, username, content string) *ValidationResult {
	result := NewValidationResult()

	// Validate individual fields
	if ip != "" {
		ipResult := ValidateIP(ip)
		if !ipResult.IsValid {
			result.Errors = append(result.Errors, ipResult.Errors...)
			result.IsValid = false
		}
	}

	if email != "" {
		emailResult := ValidateEmail(email)
		if !emailResult.IsValid {
			result.Errors = append(result.Errors, emailResult.Errors...)
			result.IsValid = false
		}
	}

	if userAgent != "" {
		userAgentResult := ValidateUserAgent(userAgent)
		if !userAgentResult.IsValid {
			result.Errors = append(result.Errors, userAgentResult.Errors...)
			result.IsValid = false
		}
	}

	if country != "" {
		countryResult := ValidateCountry(country)
		if !countryResult.IsValid {
			result.Errors = append(result.Errors, countryResult.Errors...)
			result.IsValid = false
		}
	}

	if username != "" {
		usernameResult := ValidateUsername(username)
		if !usernameResult.IsValid {
			result.Errors = append(result.Errors, usernameResult.Errors...)
			result.IsValid = false
		}
	}

	if content != "" {
		contentResult := ValidateContent(content)
		if !contentResult.IsValid {
			result.Errors = append(result.Errors, contentResult.Errors...)
			result.IsValid = false
		}
	}

	// Check that at least one field is provided
	if ip == "" && email == "" && userAgent == "" && country == "" && username == "" && content == "" {
		result.AddError("request", "At least one filter field must be provided (ip, email, user_agent, country, username, or content)", "")
	}

	return result
}

// ValidateCreateRequest validates a create request for any entity
func ValidateCreateRequest(entityType string, data map[string]interface{}) *ValidationResult {
	result := NewValidationResult()

	switch entityType {
	case "ip":
		if address, ok := data["address"].(string); ok {
			ipResult := ValidateIP(address)
			if !ipResult.IsValid {
				result.Errors = append(result.Errors, ipResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("address", "IP address is required", "")
		}

		if status, ok := data["status"].(string); ok {
			statusResult := ValidateStatus(status)
			if !statusResult.IsValid {
				result.Errors = append(result.Errors, statusResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("status", "Status is required", "")
		}

	case "email":
		if address, ok := data["address"].(string); ok {
			emailResult := ValidateEmail(address)
			if !emailResult.IsValid {
				result.Errors = append(result.Errors, emailResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("address", "Email address is required", "")
		}

		if status, ok := data["status"].(string); ok {
			statusResult := ValidateStatus(status)
			if !statusResult.IsValid {
				result.Errors = append(result.Errors, statusResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("status", "Status is required", "")
		}

	case "user_agent":
		if userAgent, ok := data["user_agent"].(string); ok {
			userAgentResult := ValidateUserAgent(userAgent)
			if !userAgentResult.IsValid {
				result.Errors = append(result.Errors, userAgentResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("user_agent", "User agent is required", "")
		}

		if status, ok := data["status"].(string); ok {
			statusResult := ValidateStatus(status)
			if !statusResult.IsValid {
				result.Errors = append(result.Errors, statusResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("status", "Status is required", "")
		}

	case "country":
		if code, ok := data["code"].(string); ok {
			countryResult := ValidateCountry(code)
			if !countryResult.IsValid {
				result.Errors = append(result.Errors, countryResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("code", "Country code is required", "")
		}

		if status, ok := data["status"].(string); ok {
			statusResult := ValidateStatus(status)
			if !statusResult.IsValid {
				result.Errors = append(result.Errors, statusResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("status", "Status is required", "")
		}

	case "charset":
		if charset, ok := data["charset"].(string); ok {
			charsetResult := ValidateCharset(charset)
			if !charsetResult.IsValid {
				result.Errors = append(result.Errors, charsetResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("charset", "Charset is required", "")
		}

		if status, ok := data["status"].(string); ok {
			statusResult := ValidateStatus(status)
			if !statusResult.IsValid {
				result.Errors = append(result.Errors, statusResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("status", "Status is required", "")
		}

	case "username":
		if username, ok := data["username"].(string); ok {
			usernameResult := ValidateUsername(username)
			if !usernameResult.IsValid {
				result.Errors = append(result.Errors, usernameResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("username", "Username is required", "")
		}

		if status, ok := data["status"].(string); ok {
			statusResult := ValidateStatus(status)
			if !statusResult.IsValid {
				result.Errors = append(result.Errors, statusResult.Errors...)
				result.IsValid = false
			}
		} else {
			result.AddError("status", "Status is required", "")
		}

	default:
		result.AddError("entity_type", "Unknown entity type", entityType)
	}

	return result
}
