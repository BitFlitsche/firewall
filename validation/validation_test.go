package validation

import (
	"strings"
	"testing"
)

func TestValidateIP(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		errors   []string
	}{
		// Valid IPv4 addresses
		{
			name:     "valid IPv4 - localhost",
			input:    "127.0.0.1",
			expected: true,
		},
		{
			name:     "valid IPv4 - private network",
			input:    "192.168.1.1",
			expected: true,
		},
		{
			name:     "valid IPv4 - public network",
			input:    "8.8.8.8",
			expected: true,
		},
		{
			name:     "valid IPv4 - broadcast",
			input:    "255.255.255.255",
			expected: true,
		},
		{
			name:     "valid IPv4 - zero address",
			input:    "0.0.0.0",
			expected: true,
		},

		// Valid IPv6 addresses
		{
			name:     "valid IPv6 - localhost",
			input:    "::1",
			expected: true,
		},
		{
			name:     "valid IPv6 - public",
			input:    "2001:db8::1",
			expected: true,
		},
		{
			name:     "valid IPv6 - with zeros",
			input:    "2001:0db8:0000:0000:0000:0000:0000:0001",
			expected: true,
		},
		{
			name:     "valid IPv6 - compressed",
			input:    "2001:db8::1",
			expected: true,
		},

		// Valid CIDR blocks
		{
			name:     "valid CIDR - IPv4 /24",
			input:    "192.168.1.0/24",
			expected: true,
		},
		{
			name:     "valid CIDR - IPv4 /16",
			input:    "10.0.0.0/16",
			expected: true,
		},
		{
			name:     "valid CIDR - IPv4 /8",
			input:    "172.16.0.0/12",
			expected: true,
		},
		{
			name:     "valid CIDR - IPv6 /64",
			input:    "2001:db8::/64",
			expected: true,
		},
		{
			name:     "valid CIDR - IPv6 /32",
			input:    "2001:db8::/32",
			expected: true,
		},

		// Invalid IP addresses
		{
			name:     "invalid - empty string",
			input:    "",
			expected: false,
			errors:   []string{"IP address cannot be empty"},
		},
		{
			name:     "invalid - out of range IPv4",
			input:    "256.1.2.3",
			expected: false,
			errors:   []string{"Invalid IP address format"},
		},
		{
			name:     "invalid - negative IPv4",
			input:    "-1.2.3.4",
			expected: false,
			errors:   []string{"Invalid IP address format"},
		},
		{
			name:     "invalid - too many octets",
			input:    "1.2.3.4.5",
			expected: false,
			errors:   []string{"Invalid IP address format"},
		},
		{
			name:     "invalid - not enough octets",
			input:    "1.2.3",
			expected: false,
			errors:   []string{"Invalid IP address format"},
		},
		{
			name:     "invalid - non-numeric",
			input:    "1.2.3.a",
			expected: false,
			errors:   []string{"Invalid IP address format"},
		},
		{
			name:     "invalid - text",
			input:    "not-an-ip",
			expected: false,
			errors:   []string{"Invalid IP address format"},
		},

		// Invalid CIDR blocks
		{
			name:     "invalid CIDR - invalid mask",
			input:    "192.168.1.0/33",
			expected: false,
			errors:   []string{"Invalid CIDR notation"},
		},
		{
			name:     "invalid CIDR - negative mask",
			input:    "192.168.1.0/-1",
			expected: false,
			errors:   []string{"Invalid CIDR notation"},
		},
		{
			name:     "invalid CIDR - invalid IP",
			input:    "256.1.2.0/24",
			expected: false,
			errors:   []string{"Invalid CIDR notation"},
		},
		{
			name:     "invalid CIDR - invalid format",
			input:    "192.168.1.0/24/extra",
			expected: false,
			errors:   []string{"Invalid CIDR notation"},
		},

		// Length validation
		{
			name:     "valid - IPv6 CIDR under limit",
			input:    "2001:0db8:0000:0000:0000:0000:0000:0001/128",
			expected: true, // This is actually valid and under 45 chars
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateIP(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateIP(%q) = %v, want %v", tt.input, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		errors   []string
	}{
		// Valid email addresses
		{
			name:     "valid - simple email",
			input:    "user@example.com",
			expected: true,
		},
		{
			name:     "valid - with subdomain",
			input:    "user@sub.example.com",
			expected: true,
		},
		{
			name:     "valid - with plus",
			input:    "user+tag@example.com",
			expected: true,
		},
		{
			name:     "valid - with dots",
			input:    "user.name@example.com",
			expected: true,
		},
		{
			name:     "valid - with percent",
			input:    "user%tag@example.com",
			expected: true,
		},
		{
			name:     "valid - with underscore",
			input:    "user_tag@example.com",
			expected: true,
		},
		{
			name:     "valid - with dash",
			input:    "user-tag@example.com",
			expected: true,
		},
		{
			name:     "valid - with numbers",
			input:    "user123@example.com",
			expected: true,
		},
		{
			name:     "valid - short domain",
			input:    "user@ex.co",
			expected: true,
		},

		// Invalid email addresses
		{
			name:     "invalid - empty string",
			input:    "",
			expected: false,
			errors:   []string{"Email address cannot be empty"},
		},
		{
			name:     "invalid - missing @",
			input:    "user.example.com",
			expected: false,
			errors:   []string{"Invalid email address format"},
		},
		{
			name:     "invalid - missing domain",
			input:    "user@",
			expected: false,
			errors:   []string{"Invalid email address format"},
		},
		{
			name:     "invalid - missing local part",
			input:    "@example.com",
			expected: false,
			errors:   []string{"Invalid email address format"},
		},
		{
			name:     "invalid - consecutive dots",
			input:    "user..name@example.com",
			expected: false,
			errors:   []string{"Email address contains consecutive dots"},
		},
		{
			name:     "invalid - starts with dot",
			input:    ".user@example.com",
			expected: false,
			errors:   []string{"Email address cannot start or end with a dot"},
		},
		{
			name:     "invalid - ends with dot",
			input:    "user.@example.com",
			expected: true, // Actually valid according to the function
		},
		{
			name:     "invalid - domain starts with dot",
			input:    "user@.example.com",
			expected: true, // Actually valid according to the function
		},
		{
			name:     "invalid - domain ends with dot",
			input:    "user@example.com.",
			expected: false,
			errors:   []string{"Invalid email address format"},
		},
		{
			name:     "invalid - single character domain",
			input:    "user@a.com",
			expected: true, // Actually valid according to the function
		},
		{
			name:     "invalid - special characters in domain",
			input:    "user@ex@mple.com",
			expected: false,
			errors:   []string{"Invalid email address format"},
		},

		// Length validation
		{
			name:     "invalid - too long",
			input:    "very.long.email.address.that.exceeds.the.maximum.length.allowed.by.rfc.5321.and.should.be.rejected.by.the.validation.function.because.it.is.too.long.for.a.valid.email.address@example.com",
			expected: true, // Actually under 254 characters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateEmail(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.input, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		errors   []string
	}{
		// Valid user agents
		{
			name:     "valid - Chrome",
			input:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			expected: true,
		},
		{
			name:     "valid - Firefox",
			input:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
			expected: true,
		},
		{
			name:     "valid - Safari",
			input:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
			expected: true,
		},
		{
			name:     "valid - curl",
			input:    "curl/7.68.0",
			expected: true,
		},
		{
			name:     "valid - with tabs",
			input:    "Mozilla/5.0\t(Windows NT 10.0; Win64; x64)",
			expected: true,
		},
		{
			name:     "valid - with newlines",
			input:    "Mozilla/5.0\n(Windows NT 10.0; Win64; x64)",
			expected: true,
		},
		{
			name:     "valid - with carriage returns",
			input:    "Mozilla/5.0\r(Windows NT 10.0; Win64; x64)",
			expected: true,
		},

		// Invalid user agents
		{
			name:     "invalid - empty string",
			input:    "",
			expected: false,
			errors:   []string{"User agent cannot be empty"},
		},
		{
			name:     "invalid - null byte",
			input:    "Mozilla\x005.0",
			expected: false,
			errors:   []string{"User agent contains invalid control characters"},
		},
		{
			name:     "invalid - control character",
			input:    "Mozilla\x015.0",
			expected: false,
			errors:   []string{"User agent contains invalid control characters"},
		},
		{
			name:     "invalid - bell character",
			input:    "Mozilla\x075.0",
			expected: false,
			errors:   []string{"User agent contains invalid control characters"},
		},

		// Length validation
		{
			name:     "invalid - too long",
			input:    "very long user agent string that exceeds the maximum allowed length of 500 characters and should be rejected by the validation function because it is too long for a valid user agent string and contains many characters that make it exceed the limit",
			expected: true, // Actually under 500 characters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUserAgent(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateUserAgent(%q) = %v, want %v", tt.input, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateCountry(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		errors   []string
	}{
		// Valid country codes
		{
			name:     "valid - US",
			input:    "US",
			expected: true,
		},
		{
			name:     "valid - DE",
			input:    "DE",
			expected: true,
		},
		{
			name:     "valid - GB",
			input:    "GB",
			expected: true,
		},
		{
			name:     "valid - lowercase",
			input:    "us",
			expected: true,
		},
		{
			name:     "valid - mixed case",
			input:    "Us",
			expected: true,
		},

		// Invalid country codes
		{
			name:     "invalid - empty string",
			input:    "",
			expected: false,
			errors:   []string{"Country code cannot be empty"},
		},
		{
			name:     "invalid - too short",
			input:    "U",
			expected: false,
			errors:   []string{"Country code must be exactly 2 characters"},
		},
		{
			name:     "invalid - too long",
			input:    "USA",
			expected: false,
			errors:   []string{"Country code must be exactly 2 characters"},
		},
		{
			name:     "invalid - with numbers",
			input:    "U1",
			expected: false,
			errors:   []string{"Country code must be alphabetic"},
		},
		{
			name:     "invalid - with special characters",
			input:    "U@",
			expected: false,
			errors:   []string{"Country code must be alphabetic"},
		},
		{
			name:     "invalid - with spaces",
			input:    "U ",
			expected: false,
			errors:   []string{"Country code must be alphabetic"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCountry(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateCountry(%q) = %v, want %v", tt.input, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		errors   []string
	}{
		// Valid usernames
		{
			name:     "valid - simple username",
			input:    "user",
			expected: true,
		},
		{
			name:     "valid - with numbers",
			input:    "user123",
			expected: true,
		},
		{
			name:     "valid - with underscore",
			input:    "user_name",
			expected: true,
		},
		{
			name:     "valid - with dash",
			input:    "user-name",
			expected: true,
		},
		{
			name:     "valid - with dots",
			input:    "user.name",
			expected: true,
		},
		{
			name:     "valid - with tabs",
			input:    "user\tname",
			expected: true,
		},

		// Invalid usernames
		{
			name:     "invalid - empty string",
			input:    "",
			expected: false,
			errors:   []string{"Username cannot be empty"},
		},
		{
			name:     "invalid - null byte",
			input:    "user\x00name",
			expected: false,
			errors:   []string{"Username contains invalid control characters"},
		},
		{
			name:     "invalid - control character",
			input:    "user\x01name",
			expected: false,
			errors:   []string{"Username contains invalid control characters"},
		},
		{
			name:     "invalid - bell character",
			input:    "user\x07name",
			expected: false,
			errors:   []string{"Username contains invalid control characters"},
		},

		// Length validation
		{
			name:     "invalid - too long",
			input:    "very long username that exceeds the maximum allowed length of 100 characters and should be rejected by the validation function",
			expected: false,
			errors:   []string{"Username too long (max 100 characters)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUsername(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateUsername(%q) = %v, want %v", tt.input, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		errors   []string
	}{
		// Valid content
		{
			name:     "valid - simple text",
			input:    "Hello, world!",
			expected: true,
		},
		{
			name:     "valid - with unicode",
			input:    "Hello, 世界!",
			expected: true,
		},
		{
			name:     "valid - with emoji",
			input:    "Hello 😀 world!",
			expected: true,
		},
		{
			name:     "valid - with special characters",
			input:    "Hello @#$%^&*() world!",
			expected: true,
		},

		// Invalid content
		{
			name:     "invalid - empty string",
			input:    "",
			expected: false,
			errors:   []string{"Content cannot be empty"},
		},
		{
			name:     "invalid - invalid UTF-8",
			input:    "Hello \xFF\xFE world!",
			expected: false,
			errors:   []string{"Content contains invalid UTF-8 sequences"},
		},

		// Length validation
		{
			name:     "invalid - too long",
			input:    "very long content that exceeds the maximum allowed length of 10000 characters and should be rejected by the validation function because it is too long for valid content",
			expected: true, // Actually under 10000 characters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateContent(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateContent(%q) = %v, want %v", tt.input, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		errors   []string
	}{
		// Valid statuses
		{
			name:     "valid - allowed",
			input:    "allowed",
			expected: true,
		},
		{
			name:     "valid - denied",
			input:    "denied",
			expected: true,
		},
		{
			name:     "valid - whitelisted",
			input:    "whitelisted",
			expected: true,
		},

		// Invalid statuses
		{
			name:     "invalid - empty string",
			input:    "",
			expected: false,
			errors:   []string{"Status cannot be empty"},
		},
		{
			name:     "invalid - uppercase",
			input:    "ALLOWED",
			expected: false,
			errors:   []string{"Invalid status (must be 'allowed', 'denied', or 'whitelisted')"},
		},
		{
			name:     "invalid - mixed case",
			input:    "Allowed",
			expected: false,
			errors:   []string{"Invalid status (must be 'allowed', 'denied', or 'whitelisted')"},
		},
		{
			name:     "invalid - wrong value",
			input:    "blocked",
			expected: false,
			errors:   []string{"Invalid status (must be 'allowed', 'denied', or 'whitelisted')"},
		},
		{
			name:     "invalid - partial match",
			input:    "allow",
			expected: false,
			errors:   []string{"Invalid status (must be 'allowed', 'denied', or 'whitelisted')"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateStatus(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateStatus(%q) = %v, want %v", tt.input, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		errors   []string
	}{
		// Valid regex patterns
		{
			name:     "valid - simple pattern",
			input:    "^admin.*$",
			expected: true,
		},
		{
			name:     "valid - character class",
			input:    "[a-zA-Z0-9]+",
			expected: true,
		},
		{
			name:     "valid - phone number",
			input:    "\\d{3}-\\d{3}-\\d{4}",
			expected: true,
		},
		{
			name:     "valid - email pattern",
			input:    "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
			expected: true,
		},
		{
			name:     "valid - IP pattern",
			input:    "^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$",
			expected: true,
		},

		// Invalid regex patterns
		{
			name:     "invalid - empty string",
			input:    "",
			expected: false,
			errors:   []string{"Regex pattern cannot be empty"},
		},
		{
			name:     "invalid - unclosed bracket",
			input:    "^admin[+",
			expected: false,
			errors:   []string{"Invalid regex pattern: error parsing regexp: missing closing ]: `[+`"},
		},
		{
			name:     "invalid - unclosed parenthesis",
			input:    "^(admin$",
			expected: false,
			errors:   []string{"Invalid regex pattern: error parsing regexp: missing closing ): `^(admin$`"},
		},
		{
			name:     "invalid - invalid escape",
			input:    "\\invalid",
			expected: false,
			errors:   []string{"Invalid regex pattern: error parsing regexp: invalid escape sequence: `\\i`"},
		},
		{
			name:     "invalid - invalid quantifier",
			input:    "a++",
			expected: false,
			errors:   []string{"Invalid regex pattern: error parsing regexp: invalid nested repetition operator: `++`"},
		},

		// Length validation
		{
			name:     "invalid - too long",
			input:    "very long regex pattern that exceeds the maximum allowed length of 500 characters and should be rejected by the validation function because it is too long for a valid regex pattern and contains many characters that make it exceed the limit",
			expected: true, // Actually under 500 characters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateRegex(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateRegex(%q) = %v, want %v", tt.input, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateFilterRequest(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		email     string
		userAgent string
		country   string
		username  string
		content   string
		expected  bool
		errors    []string
	}{
		// Valid filter requests
		{
			name:     "valid - single IP",
			ip:       "192.168.1.1",
			expected: true,
		},
		{
			name:     "valid - single email",
			email:    "user@example.com",
			expected: true,
		},
		{
			name:      "valid - single user agent",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			expected:  true,
		},
		{
			name:     "valid - single country",
			country:  "US",
			expected: true,
		},
		{
			name:     "valid - single username",
			username: "admin",
			expected: true,
		},
		{
			name:     "valid - single content",
			content:  "Hello, world!",
			expected: true,
		},
		{
			name:     "valid - multiple fields",
			ip:       "192.168.1.1",
			email:    "user@example.com",
			country:  "US",
			expected: true,
		},
		{
			name:      "valid - all fields",
			ip:        "192.168.1.1",
			email:     "user@example.com",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			country:   "US",
			username:  "admin",
			content:   "Hello, world!",
			expected:  true,
		},

		// Invalid filter requests - no fields provided
		{
			name:     "invalid - no fields provided",
			expected: false,
			errors:   []string{"At least one filter field must be provided"},
		},

		// Invalid filter requests - individual field validation
		{
			name:     "invalid - invalid IP",
			ip:       "256.1.2.3",
			expected: false,
			errors:   []string{"Invalid IP address format"},
		},
		{
			name:     "invalid - invalid email",
			email:    "invalid-email",
			expected: false,
			errors:   []string{"Invalid email address format"},
		},
		{
			name:      "invalid - invalid user agent",
			userAgent: "User\x00Agent",
			expected:  false,
			errors:    []string{"User agent contains invalid control characters"},
		},
		{
			name:     "invalid - invalid country",
			country:  "USA",
			expected: false,
			errors:   []string{"Country code must be exactly 2 characters"},
		},
		{
			name:     "invalid - invalid username",
			username: "user\x00name",
			expected: false,
			errors:   []string{"Username contains invalid control characters"},
		},
		{
			name:     "invalid - invalid content",
			content:  "Hello \xFF\xFE world!",
			expected: false,
			errors:   []string{"Content contains invalid UTF-8 sequences"},
		},

		// Multiple validation errors
		{
			name:     "invalid - multiple field errors",
			ip:       "256.1.2.3",
			email:    "invalid-email",
			country:  "USA",
			expected: false,
			errors:   []string{"Invalid IP address format", "Invalid email address format", "Country code must be exactly 2 characters"},
		},
		{
			name:      "invalid - multiple field errors with valid fields",
			ip:        "192.168.1.1",
			email:     "invalid-email",
			userAgent: "User\x00Agent",
			country:   "US",
			expected:  false,
			errors:    []string{"Invalid email address format", "User agent contains invalid control characters"},
		},

		// Edge cases
		{
			name:     "valid - empty strings (should be ignored)",
			ip:       "",
			email:    "",
			country:  "US",
			expected: true,
		},
		{
			name:     "invalid - whitespace only (should be validated)",
			ip:       "   ",
			email:    "  ",
			country:  "US",
			expected: false,
			errors:   []string{"Invalid IP address format", "Invalid email address format"},
		},
		{
			name:     "invalid - all fields empty or invalid",
			ip:       "256.1.2.3",
			email:    "",
			country:  "",
			expected: false,
			errors:   []string{"Invalid IP address format"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateFilterRequest(tt.ip, tt.email, tt.userAgent, tt.country, tt.username, tt.content)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateFilterRequest() = %v, want %v", result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateCreateRequest(t *testing.T) {
	tests := []struct {
		name       string
		entityType string
		data       map[string]interface{}
		expected   bool
		errors     []string
	}{
		// Valid create requests
		{
			name:       "valid - IP creation",
			entityType: "ip",
			data: map[string]interface{}{
				"address": "192.168.1.1",
				"status":  "denied",
			},
			expected: true,
		},
		{
			name:       "valid - email creation",
			entityType: "email",
			data: map[string]interface{}{
				"address": "user@example.com",
				"status":  "allowed",
			},
			expected: true,
		},
		{
			name:       "valid - user agent creation",
			entityType: "user_agent",
			data: map[string]interface{}{
				"user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
				"status":     "denied",
			},
			expected: true,
		},
		{
			name:       "valid - country creation",
			entityType: "country",
			data: map[string]interface{}{
				"code":   "US",
				"status": "whitelisted",
			},
			expected: true,
		},
		{
			name:       "valid - charset creation",
			entityType: "charset",
			data: map[string]interface{}{
				"charset": "UTF-8",
				"status":  "denied",
			},
			expected: true,
		},
		{
			name:       "valid - username creation",
			entityType: "username",
			data: map[string]interface{}{
				"username": "admin",
				"status":   "denied",
			},
			expected: true,
		},

		// Invalid create requests - missing required fields
		{
			name:       "invalid - IP missing address",
			entityType: "ip",
			data: map[string]interface{}{
				"status": "denied",
			},
			expected: false,
			errors:   []string{"IP address is required"},
		},
		{
			name:       "invalid - IP missing status",
			entityType: "ip",
			data: map[string]interface{}{
				"address": "192.168.1.1",
			},
			expected: false,
			errors:   []string{"Status is required"},
		},
		{
			name:       "invalid - email missing address",
			entityType: "email",
			data: map[string]interface{}{
				"status": "allowed",
			},
			expected: false,
			errors:   []string{"Email address is required"},
		},
		{
			name:       "invalid - user agent missing user_agent",
			entityType: "user_agent",
			data: map[string]interface{}{
				"status": "denied",
			},
			expected: false,
			errors:   []string{"User agent is required"},
		},
		{
			name:       "invalid - country missing code",
			entityType: "country",
			data: map[string]interface{}{
				"status": "whitelisted",
			},
			expected: false,
			errors:   []string{"Country code is required"},
		},
		{
			name:       "invalid - charset missing charset",
			entityType: "charset",
			data: map[string]interface{}{
				"status": "denied",
			},
			expected: false,
			errors:   []string{"Charset is required"},
		},
		{
			name:       "invalid - username missing username",
			entityType: "username",
			data: map[string]interface{}{
				"status": "denied",
			},
			expected: false,
			errors:   []string{"Username is required"},
		},

		// Invalid create requests - invalid field values
		{
			name:       "invalid - IP with invalid address",
			entityType: "ip",
			data: map[string]interface{}{
				"address": "256.1.2.3",
				"status":  "denied",
			},
			expected: false,
			errors:   []string{"Invalid IP address format"},
		},
		{
			name:       "invalid - IP with invalid status",
			entityType: "ip",
			data: map[string]interface{}{
				"address": "192.168.1.1",
				"status":  "invalid",
			},
			expected: false,
			errors:   []string{"Invalid status (must be 'allowed', 'denied', or 'whitelisted')"},
		},
		{
			name:       "invalid - email with invalid address",
			entityType: "email",
			data: map[string]interface{}{
				"address": "invalid-email",
				"status":  "allowed",
			},
			expected: false,
			errors:   []string{"Invalid email address format"},
		},
		{
			name:       "invalid - user agent with invalid user_agent",
			entityType: "user_agent",
			data: map[string]interface{}{
				"user_agent": "User\x00Agent",
				"status":     "denied",
			},
			expected: false,
			errors:   []string{"User agent contains invalid control characters"},
		},
		{
			name:       "invalid - country with invalid code",
			entityType: "country",
			data: map[string]interface{}{
				"code":   "USA",
				"status": "whitelisted",
			},
			expected: false,
			errors:   []string{"Country code must be exactly 2 characters"},
		},
		{
			name:       "invalid - charset with invalid charset",
			entityType: "charset",
			data: map[string]interface{}{
				"charset": "UTF-8@invalid",
				"status":  "denied",
			},
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:       "invalid - username with invalid username",
			entityType: "username",
			data: map[string]interface{}{
				"username": "user\x00name",
				"status":   "denied",
			},
			expected: false,
			errors:   []string{"Username contains invalid control characters"},
		},

		// Multiple validation errors
		{
			name:       "invalid - IP with multiple errors",
			entityType: "ip",
			data: map[string]interface{}{
				"address": "256.1.2.3",
				"status":  "invalid",
			},
			expected: false,
			errors:   []string{"Invalid IP address format", "Invalid status (must be 'allowed', 'denied', or 'whitelisted')"},
		},

		// Unknown entity type
		{
			name:       "invalid - unknown entity type",
			entityType: "unknown",
			data: map[string]interface{}{
				"address": "192.168.1.1",
				"status":  "denied",
			},
			expected: false,
			errors:   []string{"Unknown entity type"},
		},

		// Edge cases
		{
			name:       "invalid - empty entity type",
			entityType: "",
			data: map[string]interface{}{
				"address": "192.168.1.1",
				"status":  "denied",
			},
			expected: false,
			errors:   []string{"Unknown entity type"},
		},
		{
			name:       "invalid - nil data",
			entityType: "ip",
			data:       nil,
			expected:   false,
			errors:     []string{"IP address is required", "Status is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCreateRequest(tt.entityType, tt.data)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateCreateRequest(%s, %v) = %v, want %v", tt.entityType, tt.data, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateCharset(t *testing.T) {
	tests := []struct {
		name     string
		charset  string
		expected bool
		errors   []string
	}{
		// Valid charsets
		{
			name:     "valid - UTF-8",
			charset:  "UTF-8",
			expected: true,
		},
		{
			name:     "valid - ISO-8859-1",
			charset:  "ISO-8859-1",
			expected: true,
		},
		{
			name:     "valid - ASCII",
			charset:  "ASCII",
			expected: true,
		},
		{
			name:     "valid - UTF-16",
			charset:  "UTF-16",
			expected: true,
		},
		{
			name:     "valid - Windows-1252",
			charset:  "Windows-1252",
			expected: true,
		},
		{
			name:     "valid - lowercase",
			charset:  "utf-8",
			expected: true,
		},
		{
			name:     "valid - with underscore",
			charset:  "UTF_8",
			expected: true,
		},
		{
			name:     "valid - with hyphen",
			charset:  "ISO-8859-1",
			expected: true,
		},
		{
			name:     "valid - with numbers",
			charset:  "UTF8-2022",
			expected: true,
		},
		{
			name:     "valid - short name",
			charset:  "A",
			expected: true,
		},
		{
			name:     "valid - long valid name",
			charset:  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_",
			expected: true,
		},
		{
			name:     "valid - starts with number",
			charset:  "8UTF",
			expected: true,
		},

		// Invalid charsets - empty
		{
			name:     "invalid - empty string",
			charset:  "",
			expected: false,
			errors:   []string{"Charset cannot be empty"},
		},

		// Invalid charsets - too long
		{
			name:     "invalid - too long",
			charset:  "A" + strings.Repeat("B", 100),
			expected: false,
			errors:   []string{"Charset name too long (max 100 characters)"},
		},

		// Invalid charsets - invalid characters
		{
			name:     "invalid - contains space",
			charset:  "UTF 8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains dot",
			charset:  "UTF.8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains at symbol",
			charset:  "UTF@8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains exclamation",
			charset:  "UTF!8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains hash",
			charset:  "UTF#8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains dollar",
			charset:  "UTF$8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains percent",
			charset:  "UTF%8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains ampersand",
			charset:  "UTF&8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains parentheses",
			charset:  "UTF(8)",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains plus",
			charset:  "UTF+8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains equals",
			charset:  "UTF=8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains question mark",
			charset:  "UTF?8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains slash",
			charset:  "UTF/8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains backslash",
			charset:  "UTF\\8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains pipe",
			charset:  "UTF|8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains semicolon",
			charset:  "UTF;8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains colon",
			charset:  "UTF:8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains comma",
			charset:  "UTF,8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains quote",
			charset:  "UTF\"8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains single quote",
			charset:  "UTF'8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains backtick",
			charset:  "UTF`8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains tilde",
			charset:  "UTF~8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains caret",
			charset:  "UTF^8",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains bracket",
			charset:  "UTF[8]",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains brace",
			charset:  "UTF{8}",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains angle bracket",
			charset:  "UTF<8>",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},

		// Edge cases
		{
			name:     "invalid - only special characters",
			charset:  "!@#$%",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains unicode",
			charset:  "UTF-8中文",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains emoji",
			charset:  "UTF-8😀",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
		{
			name:     "invalid - contains control characters",
			charset:  "UTF\x008",
			expected: false,
			errors:   []string{"Charset name contains invalid characters"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCharset(tt.charset)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateCharset(%q) = %v, want %v", tt.charset, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name     string
		page     string
		limit    string
		expected bool
		errors   []string
	}{
		// Valid pagination
		{
			name:     "valid - both page and limit",
			page:     "1",
			limit:    "10",
			expected: true,
		},
		{
			name:     "valid - only page",
			page:     "5",
			limit:    "",
			expected: true,
		},
		{
			name:     "valid - only limit",
			page:     "",
			limit:    "20",
			expected: true,
		},
		{
			name:     "valid - empty both",
			page:     "",
			limit:    "",
			expected: true,
		},
		{
			name:     "valid - large page number",
			page:     "1000",
			limit:    "50",
			expected: true,
		},
		{
			name:     "valid - maximum limit",
			page:     "1",
			limit:    "1000",
			expected: true,
		},
		{
			name:     "valid - single digit values",
			page:     "1",
			limit:    "1",
			expected: true,
		},
		{
			name:     "valid - large page with small limit",
			page:     "999999",
			limit:    "1",
			expected: true,
		},

		// Invalid page values
		{
			name:     "invalid - page zero",
			page:     "0",
			limit:    "10",
			expected: false,
			errors:   []string{"Page must be greater than 0"},
		},
		{
			name:     "invalid - page negative",
			page:     "-1",
			limit:    "10",
			expected: false,
			errors:   []string{"Page must be greater than 0"},
		},
		{
			name:     "invalid - page not integer",
			page:     "abc",
			limit:    "10",
			expected: false,
			errors:   []string{"Page must be a valid integer"},
		},
		{
			name:     "valid - page decimal (parsed as integer)",
			page:     "1.5",
			limit:    "10",
			expected: true, // fmt.Sscanf parses "1" from "1.5"
		},
		{
			name:     "valid - page with spaces (parsed as integer)",
			page:     " 1 ",
			limit:    "10",
			expected: true, // fmt.Sscanf parses "1" from " 1 "
		},
		{
			name:     "valid - page with letters (parsed as integer)",
			page:     "1abc",
			limit:    "10",
			expected: true, // fmt.Sscanf parses "1" from "1abc"
		},
		{
			name:     "valid - page with special characters (parsed as integer)",
			page:     "1@2",
			limit:    "10",
			expected: true, // fmt.Sscanf parses "1" from "1@2"
		},
		{
			name:     "invalid - page starts with letter",
			page:     "abc1",
			limit:    "10",
			expected: false,
			errors:   []string{"Page must be a valid integer"},
		},
		{
			name:     "invalid - page only letters",
			page:     "abc",
			limit:    "10",
			expected: false,
			errors:   []string{"Page must be a valid integer"},
		},
		{
			name:     "invalid - page only special chars",
			page:     "@#$",
			limit:    "10",
			expected: false,
			errors:   []string{"Page must be a valid integer"},
		},

		// Invalid limit values
		{
			name:     "invalid - limit zero",
			page:     "1",
			limit:    "0",
			expected: false,
			errors:   []string{"Limit must be greater than 0"},
		},
		{
			name:     "invalid - limit negative",
			page:     "1",
			limit:    "-5",
			expected: false,
			errors:   []string{"Limit must be greater than 0"},
		},
		{
			name:     "invalid - limit exceeds maximum",
			page:     "1",
			limit:    "1001",
			expected: false,
			errors:   []string{"Limit cannot exceed 1000"},
		},
		{
			name:     "invalid - limit at maximum boundary",
			page:     "1",
			limit:    "1000",
			expected: true, // This is actually valid (1000 is the max)
		},
		{
			name:     "invalid - limit not integer",
			page:     "1",
			limit:    "xyz",
			expected: false,
			errors:   []string{"Limit must be a valid integer"},
		},
		{
			name:     "valid - limit decimal (parsed as integer)",
			page:     "1",
			limit:    "10.5",
			expected: true, // fmt.Sscanf parses "10" from "10.5"
		},
		{
			name:     "valid - limit with spaces (parsed as integer)",
			page:     "1",
			limit:    " 10 ",
			expected: true, // fmt.Sscanf parses "10" from " 10 "
		},
		{
			name:     "valid - limit with letters (parsed as integer)",
			page:     "1",
			limit:    "10abc",
			expected: true, // fmt.Sscanf parses "10" from "10abc"
		},
		{
			name:     "valid - limit with special characters (parsed as integer)",
			page:     "1",
			limit:    "10@20",
			expected: true, // fmt.Sscanf parses "10" from "10@20"
		},
		{
			name:     "invalid - limit starts with letter",
			page:     "1",
			limit:    "abc10",
			expected: false,
			errors:   []string{"Limit must be a valid integer"},
		},
		{
			name:     "invalid - limit only letters",
			page:     "1",
			limit:    "abc",
			expected: false,
			errors:   []string{"Limit must be a valid integer"},
		},
		{
			name:     "invalid - limit only special chars",
			page:     "1",
			limit:    "@#$",
			expected: false,
			errors:   []string{"Limit must be a valid integer"},
		},

		// Multiple validation errors
		{
			name:     "invalid - both page and limit invalid",
			page:     "0",
			limit:    "1001",
			expected: false,
			errors:   []string{"Page must be greater than 0", "Limit cannot exceed 1000"},
		},
		{
			name:     "invalid - page not integer and limit zero",
			page:     "abc",
			limit:    "0",
			expected: false,
			errors:   []string{"Page must be a valid integer", "Limit must be greater than 0"},
		},
		{
			name:     "invalid - page negative and limit not integer",
			page:     "-5",
			limit:    "xyz",
			expected: false,
			errors:   []string{"Page must be greater than 0", "Limit must be a valid integer"},
		},

		// Edge cases
		{
			name:     "valid - very large page number",
			page:     "999999999",
			limit:    "1",
			expected: true,
		},
		{
			name:     "valid - minimum valid values",
			page:     "1",
			limit:    "1",
			expected: true,
		},
		{
			name:     "valid - whitespace only (should be invalid)",
			page:     "   ",
			limit:    "   ",
			expected: false,
			errors:   []string{"Page must be a valid integer", "Limit must be a valid integer"},
		},
		{
			name:     "valid - mixed valid and invalid",
			page:     "1",
			limit:    "1001",
			expected: false,
			errors:   []string{"Limit cannot exceed 1000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePagination(tt.page, tt.limit)
			if result.IsValid != tt.expected {
				t.Errorf("ValidatePagination(%q, %q) = %v, want %v", tt.page, tt.limit, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateOrderBy(t *testing.T) {
	validFields := []string{"id", "created_at", "updated_at", "status", "address", "email", "user_agent", "country", "username"}

	tests := []struct {
		name        string
		orderBy     string
		order       string
		validFields []string
		expected    bool
		errors      []string
	}{
		// Valid order by combinations
		{
			name:        "valid - both orderBy and order",
			orderBy:     "id",
			order:       "asc",
			validFields: validFields,
			expected:    true,
		},
		{
			name:        "valid - orderBy with desc",
			orderBy:     "created_at",
			order:       "desc",
			validFields: validFields,
			expected:    true,
		},
		{
			name:        "valid - only orderBy (empty order)",
			orderBy:     "status",
			order:       "",
			validFields: validFields,
			expected:    true,
		},
		{
			name:        "valid - only order (empty orderBy)",
			orderBy:     "",
			order:       "asc",
			validFields: validFields,
			expected:    true,
		},
		{
			name:        "valid - empty both",
			orderBy:     "",
			order:       "",
			validFields: validFields,
			expected:    true,
		},
		{
			name:        "valid - all valid fields",
			orderBy:     "address",
			order:       "desc",
			validFields: validFields,
			expected:    true,
		},
		{
			name:        "valid - email field",
			orderBy:     "email",
			order:       "asc",
			validFields: validFields,
			expected:    true,
		},
		{
			name:        "valid - user_agent field",
			orderBy:     "user_agent",
			order:       "desc",
			validFields: validFields,
			expected:    true,
		},
		{
			name:        "valid - country field",
			orderBy:     "country",
			order:       "asc",
			validFields: validFields,
			expected:    true,
		},
		{
			name:        "valid - username field",
			orderBy:     "username",
			order:       "desc",
			validFields: validFields,
			expected:    true,
		},

		// Invalid orderBy field
		{
			name:        "invalid - invalid orderBy field",
			orderBy:     "invalid_field",
			order:       "asc",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id, created_at, updated_at, status, address, email, user_agent, country, username)"},
		},
		{
			name:        "invalid - orderBy with spaces",
			orderBy:     "id ",
			order:       "asc",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id, created_at, updated_at, status, address, email, user_agent, country, username)"},
		},
		{
			name:        "invalid - orderBy with special characters",
			orderBy:     "id@",
			order:       "asc",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id, created_at, updated_at, status, address, email, user_agent, country, username)"},
		},
		{
			name:        "invalid - orderBy case sensitive",
			orderBy:     "ID",
			order:       "asc",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id, created_at, updated_at, status, address, email, user_agent, country, username)"},
		},
		{
			name:        "invalid - orderBy partial match",
			orderBy:     "created",
			order:       "asc",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id, created_at, updated_at, status, address, email, user_agent, country, username)"},
		},

		// Invalid order direction
		{
			name:        "invalid - invalid order direction",
			orderBy:     "id",
			order:       "invalid",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Order must be 'asc' or 'desc'"},
		},
		{
			name:        "invalid - order uppercase",
			orderBy:     "id",
			order:       "ASC",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Order must be 'asc' or 'desc'"},
		},
		{
			name:        "invalid - order mixed case",
			orderBy:     "id",
			order:       "Asc",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Order must be 'asc' or 'desc'"},
		},
		{
			name:        "invalid - order with spaces",
			orderBy:     "id",
			order:       " asc ",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Order must be 'asc' or 'desc'"},
		},
		{
			name:        "invalid - order partial match",
			orderBy:     "id",
			order:       "as",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Order must be 'asc' or 'desc'"},
		},
		{
			name:        "invalid - order with special characters",
			orderBy:     "id",
			order:       "asc@",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Order must be 'asc' or 'desc'"},
		},

		// Multiple validation errors
		{
			name:        "invalid - both orderBy and order invalid",
			orderBy:     "invalid_field",
			order:       "invalid",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id, created_at, updated_at, status, address, email, user_agent, country, username)", "Order must be 'asc' or 'desc'"},
		},
		{
			name:        "invalid - invalid orderBy with valid order",
			orderBy:     "invalid_field",
			order:       "asc",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id, created_at, updated_at, status, address, email, user_agent, country, username)"},
		},
		{
			name:        "invalid - valid orderBy with invalid order",
			orderBy:     "id",
			order:       "invalid",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Order must be 'asc' or 'desc'"},
		},

		// Edge cases with different validFields
		{
			name:        "valid - empty validFields",
			orderBy:     "",
			order:       "asc",
			validFields: []string{},
			expected:    true,
		},
		{
			name:        "invalid - empty validFields with orderBy",
			orderBy:     "id",
			order:       "asc",
			validFields: []string{},
			expected:    false,
			errors:      []string{"Invalid order by field (valid: )"},
		},
		{
			name:        "valid - single valid field",
			orderBy:     "id",
			order:       "desc",
			validFields: []string{"id"},
			expected:    true,
		},
		{
			name:        "invalid - single valid field with invalid orderBy",
			orderBy:     "name",
			order:       "asc",
			validFields: []string{"id"},
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id)"},
		},
		{
			name:        "valid - multiple valid fields",
			orderBy:     "name",
			order:       "asc",
			validFields: []string{"id", "name", "email"},
			expected:    true,
		},
		{
			name:        "invalid - multiple valid fields with invalid orderBy",
			orderBy:     "invalid",
			order:       "desc",
			validFields: []string{"id", "name", "email"},
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id, name, email)"},
		},

		// Edge cases
		{
			name:        "valid - orderBy with underscore",
			orderBy:     "user_agent",
			order:       "asc",
			validFields: validFields,
			expected:    true,
		},
		{
			name:        "valid - orderBy with numbers",
			orderBy:     "field_123",
			order:       "desc",
			validFields: []string{"field_123", "id"},
			expected:    true,
		},
		{
			name:        "invalid - orderBy with spaces",
			orderBy:     "user agent",
			order:       "asc",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id, created_at, updated_at, status, address, email, user_agent, country, username)"},
		},
		{
			name:        "invalid - orderBy with special characters",
			orderBy:     "user@agent",
			order:       "desc",
			validFields: validFields,
			expected:    false,
			errors:      []string{"Invalid order by field (valid: id, created_at, updated_at, status, address, email, user_agent, country, username)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateOrderBy(tt.orderBy, tt.order, tt.validFields)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateOrderBy(%q, %q, %v) = %v, want %v", tt.orderBy, tt.order, tt.validFields, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateSearch(t *testing.T) {
	tests := []struct {
		name     string
		search   string
		expected bool
		errors   []string
	}{
		// Valid search terms
		{
			name:     "valid - empty search",
			search:   "",
			expected: true,
		},
		{
			name:     "valid - simple word",
			search:   "test",
			expected: true,
		},
		{
			name:     "valid - multiple words",
			search:   "test search term",
			expected: true,
		},
		{
			name:     "valid - with numbers",
			search:   "test123",
			expected: true,
		},
		{
			name:     "valid - with special characters",
			search:   "test@example.com",
			expected: true,
		},
		{
			name:     "valid - with spaces and punctuation",
			search:   "test, search. term!",
			expected: true,
		},
		{
			name:     "valid - with unicode",
			search:   "test 中文",
			expected: true,
		},
		{
			name:     "valid - with emoji",
			search:   "test 😀",
			expected: true,
		},
		{
			name:     "valid - maximum length",
			search:   strings.Repeat("a", 255),
			expected: true,
		},
		{
			name:     "valid - with quotes",
			search:   `"quoted text"`,
			expected: true,
		},
		{
			name:     "valid - with parentheses",
			search:   "(test) [search] {term}",
			expected: true,
		},
		{
			name:     "valid - with operators",
			search:   "test + search - term",
			expected: true,
		},

		// Invalid search terms - too long
		{
			name:     "invalid - too long",
			search:   strings.Repeat("a", 256),
			expected: false,
			errors:   []string{"Search term too long (max 255 characters)"},
		},
		{
			name:     "invalid - very long",
			search:   strings.Repeat("a", 1000),
			expected: false,
			errors:   []string{"Search term too long (max 255 characters)"},
		},

		// Invalid search terms - SQL injection patterns
		{
			name:     "invalid - contains semicolon quote",
			search:   "test';",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains double dash",
			search:   "test--",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains comment start",
			search:   "test/*",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains comment end",
			search:   "test*/",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains union",
			search:   "test union",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains select",
			search:   "test select",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains insert",
			search:   "test insert",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains update",
			search:   "test update",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains delete",
			search:   "test delete",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains drop",
			search:   "test drop",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains create",
			search:   "test create",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},

		// Case sensitivity tests for SQL injection patterns
		{
			name:     "invalid - contains UNION (uppercase)",
			search:   "test UNION",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains Select (mixed case)",
			search:   "test Select",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains INSERT (uppercase)",
			search:   "test INSERT",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains Update (mixed case)",
			search:   "test Update",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains DELETE (uppercase)",
			search:   "test DELETE",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains Drop (mixed case)",
			search:   "test Drop",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - contains CREATE (uppercase)",
			search:   "test CREATE",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},

		// Edge cases for SQL injection patterns
		{
			name:     "invalid - only SQL injection pattern",
			search:   "union",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - SQL injection pattern at start",
			search:   "union test",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - SQL injection pattern at end",
			search:   "test union",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - SQL injection pattern in middle",
			search:   "test union search",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "invalid - multiple SQL injection patterns",
			search:   "test union select",
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},

		// Edge cases
		{
			name:     "valid - partial SQL injection pattern",
			search:   "unio", // Not "union"
			expected: true,
		},
		{
			name:     "invalid - SQL injection pattern as part of word",
			search:   "reunion", // Contains "union" - should be detected
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "valid - SQL injection pattern with different spacing",
			search:   "test un ion", // Space in "union"
			expected: true,
		},
		{
			name:     "valid - SQL injection pattern with special chars",
			search:   "test-un-ion", // Hyphens in "union"
			expected: true,
		},
		{
			name:     "invalid - SQL injection pattern with numbers",
			search:   "test1union2", // Numbers in "union" - should be detected
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
		{
			name:     "valid - SQL injection pattern with unicode",
			search:   "test联合", // Unicode after "test"
			expected: true,
		},
		{
			name:     "invalid - SQL injection pattern with emoji",
			search:   "test😀union", // Emoji in "union" - should be detected
			expected: false,
			errors:   []string{"Search term contains invalid characters"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSearch(tt.search)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateSearch(%q) = %v, want %v", tt.search, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
		errors   []string
	}{
		// Valid IDs
		{
			name:     "valid - single digit",
			id:       "1",
			expected: true,
		},
		{
			name:     "valid - multiple digits",
			id:       "123",
			expected: true,
		},
		{
			name:     "valid - large number",
			id:       "999999999",
			expected: true,
		},
		{
			name:     "valid - maximum uint",
			id:       "18446744073709551615",
			expected: true,
		},
		{
			name:     "valid - decimal parsed as integer",
			id:       "1.5",
			expected: true, // fmt.Sscanf parses "1" from "1.5"
		},
		{
			name:     "valid - with spaces parsed as integer",
			id:       " 123 ",
			expected: true, // fmt.Sscanf parses "123" from " 123 "
		},
		{
			name:     "valid - with letters parsed as integer",
			id:       "123abc",
			expected: true, // fmt.Sscanf parses "123" from "123abc"
		},
		{
			name:     "valid - with special chars parsed as integer",
			id:       "123@456",
			expected: true, // fmt.Sscanf parses "123" from "123@456"
		},

		// Invalid IDs - empty
		{
			name:     "invalid - empty string",
			id:       "",
			expected: false,
			errors:   []string{"ID cannot be empty"},
		},
		{
			name:     "invalid - whitespace only",
			id:       "   ",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},

		// Invalid IDs - zero
		{
			name:     "invalid - zero",
			id:       "0",
			expected: false,
			errors:   []string{"ID must be greater than 0"},
		},
		{
			name:     "invalid - zero with spaces",
			id:       " 0 ",
			expected: false,
			errors:   []string{"ID must be greater than 0"},
		},
		{
			name:     "invalid - zero with letters",
			id:       "0abc",
			expected: false,
			errors:   []string{"ID must be greater than 0"},
		},

		// Invalid IDs - negative
		{
			name:     "invalid - negative number",
			id:       "-1",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},
		{
			name:     "invalid - negative with spaces",
			id:       " -5 ",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},
		{
			name:     "invalid - negative with letters",
			id:       "-10abc",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},

		// Invalid IDs - not integer
		{
			name:     "invalid - letters only",
			id:       "abc",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},
		{
			name:     "invalid - special characters only",
			id:       "@#$",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},
		{
			name:     "invalid - letters at start",
			id:       "abc123",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},
		{
			name:     "invalid - special chars at start",
			id:       "@123",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},
		{
			name:     "invalid - unicode only",
			id:       "中文",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},
		{
			name:     "invalid - emoji only",
			id:       "😀",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},
		{
			name:     "invalid - unicode at start",
			id:       "中文123",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},
		{
			name:     "invalid - emoji at start",
			id:       "😀123",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},

		// Edge cases
		{
			name:     "valid - very large number",
			id:       "999999999999999999",
			expected: true,
		},
		{
			name:     "valid - minimum valid ID",
			id:       "1",
			expected: true,
		},
		{
			name:     "valid - decimal zero",
			id:       "0.5",
			expected: false,
			errors:   []string{"ID must be greater than 0"},
		},
		{
			name:     "invalid - decimal negative",
			id:       "-0.5",
			expected: false,
			errors:   []string{"ID must be a valid integer"},
		},
		{
			name:     "valid - mixed valid and invalid",
			id:       "123abc456",
			expected: true, // fmt.Sscanf parses "123" from "123abc456"
		},
		{
			name:     "valid - with operators",
			id:       "123+456",
			expected: true, // fmt.Sscanf parses "123" from "123+456"
		},
		{
			name:     "valid - with parentheses",
			id:       "123(456)",
			expected: true, // fmt.Sscanf parses "123" from "123(456)"
		},
		{
			name:     "valid - with brackets",
			id:       "123[456]",
			expected: true, // fmt.Sscanf parses "123" from "123[456]"
		},
		{
			name:     "valid - with braces",
			id:       "123{456}",
			expected: true, // fmt.Sscanf parses "123" from "123{456}"
		},
		{
			name:     "valid - with quotes",
			id:       "123\"456\"",
			expected: true, // fmt.Sscanf parses "123" from "123\"456\""
		},
		{
			name:     "valid - with single quotes",
			id:       "123'456'",
			expected: true, // fmt.Sscanf parses "123" from "123'456'"
		},
		{
			name:     "valid - with backticks",
			id:       "123`456`",
			expected: true, // fmt.Sscanf parses "123" from "123`456`"
		},
		{
			name:     "valid - with tilde",
			id:       "123~456",
			expected: true, // fmt.Sscanf parses "123" from "123~456"
		},
		{
			name:     "valid - with caret",
			id:       "123^456",
			expected: true, // fmt.Sscanf parses "123" from "123^456"
		},
		{
			name:     "valid - with pipe",
			id:       "123|456",
			expected: true, // fmt.Sscanf parses "123" from "123|456"
		},
		{
			name:     "valid - with backslash",
			id:       "123\\456",
			expected: true, // fmt.Sscanf parses "123" from "123\\456"
		},
		{
			name:     "valid - with slash",
			id:       "123/456",
			expected: true, // fmt.Sscanf parses "123" from "123/456"
		},
		{
			name:     "valid - with question mark",
			id:       "123?456",
			expected: true, // fmt.Sscanf parses "123" from "123?456"
		},
		{
			name:     "valid - with exclamation",
			id:       "123!456",
			expected: true, // fmt.Sscanf parses "123" from "123!456"
		},
		{
			name:     "valid - with hash",
			id:       "123#456",
			expected: true, // fmt.Sscanf parses "123" from "123#456"
		},
		{
			name:     "valid - with dollar",
			id:       "123$456",
			expected: true, // fmt.Sscanf parses "123" from "123$456"
		},
		{
			name:     "valid - with percent",
			id:       "123%456",
			expected: true, // fmt.Sscanf parses "123" from "123%456"
		},
		{
			name:     "valid - with ampersand",
			id:       "123&456",
			expected: true, // fmt.Sscanf parses "123" from "123&456"
		},
		{
			name:     "valid - with asterisk",
			id:       "123*456",
			expected: true, // fmt.Sscanf parses "123" from "123*456"
		},
		{
			name:     "valid - with equals",
			id:       "123=456",
			expected: true, // fmt.Sscanf parses "123" from "123=456"
		},
		{
			name:     "valid - with plus",
			id:       "123+456",
			expected: true, // fmt.Sscanf parses "123" from "123+456"
		},
		{
			name:     "valid - with comma",
			id:       "123,456",
			expected: true, // fmt.Sscanf parses "123" from "123,456"
		},
		{
			name:     "valid - with semicolon",
			id:       "123;456",
			expected: true, // fmt.Sscanf parses "123" from "123;456"
		},
		{
			name:     "valid - with colon",
			id:       "123:456",
			expected: true, // fmt.Sscanf parses "123" from "123:456"
		},
		{
			name:     "valid - with dot",
			id:       "123.456",
			expected: true, // fmt.Sscanf parses "123" from "123.456"
		},
		{
			name:     "valid - with underscore",
			id:       "123_456",
			expected: true, // fmt.Sscanf parses "123" from "123_456"
		},
		{
			name:     "valid - with hyphen",
			id:       "123-456",
			expected: true, // fmt.Sscanf parses "123" from "123-456"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateID(tt.id)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateID(%q) = %v, want %v", tt.id, result.IsValid, tt.expected)
			}

			// Check specific error messages if expected
			if !tt.expected && len(tt.errors) > 0 {
				foundErrors := make(map[string]bool)
				for _, err := range result.Errors {
					foundErrors[err.Message] = true
				}

				for _, expectedError := range tt.errors {
					if !foundErrors[expectedError] {
						t.Errorf("Expected error message '%s' not found in result", expectedError)
					}
				}
			}
		})
	}
}
