package services

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/patrickmn/go-cache"
)

var geoipReader *geoip2.Reader
var asnReader *geoip2.Reader

// Cache for geolocation lookups
var geoCache *cache.Cache

// InitGeoIP initializes the MaxMind GeoIP database reader
func InitGeoIP() error {
	// Look for the database file in the root directory
	dbPath := "GeoLite2-Country.mmdb"

	// Check if file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("GeoIP database not found at %s. Please download GeoLite2-Country.mmdb to the root directory", dbPath)
	}

	reader, err := geoip2.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open GeoIP database: %v", err)
	}

	geoipReader = reader

	// Initialize cache with 1 hour default expiration and cleanup every 10 minutes
	geoCache = cache.New(1*time.Hour, 10*time.Minute)

	return nil
}

// InitASN initializes the MaxMind ASN database reader
func InitASN() error {
	// Look for the ASN database file in the root directory
	dbPath := "GeoLite2-ASN.mmdb"

	// Check if file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("ASN database not found at %s. Please download GeoLite2-ASN.mmdb to the root directory", dbPath)
	}

	reader, err := geoip2.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open ASN database: %v", err)
	}

	asnReader = reader
	return nil
}

// CloseGeoIP closes the GeoIP database reader
func CloseGeoIP() {
	if geoipReader != nil {
		geoipReader.Close()
	}
	if asnReader != nil {
		asnReader.Close()
	}
}

// IsPrivateIP checks if an IP address is private/local
func IsPrivateIP(ip net.IP) bool {
	if ip == nil {
		return true
	}

	// Check for private IP ranges
	privateRanges := []string{
		"10.0.0.0/8",     // Class A private
		"172.16.0.0/12",  // Class B private
		"192.168.0.0/16", // Class C private
		"127.0.0.0/8",    // Loopback
		"169.254.0.0/16", // Link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local
	}

	for _, cidr := range privateRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}

// GetCountryFromIP resolves an IP address to a country code
func GetCountryFromIP(ipStr string) (string, error) {
	if geoipReader == nil {
		return "", fmt.Errorf("GeoIP database not initialized")
	}

	// Parse the IP address
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address: %s", ipStr)
	}

	// Skip private/local IPs
	if IsPrivateIP(ip) {
		return "", fmt.Errorf("private IP address: %s", ipStr)
	}

	// Look up the country
	record, err := geoipReader.Country(ip)
	if err != nil {
		return "", fmt.Errorf("failed to lookup country for IP %s: %v", ipStr, err)
	}

	// Return the ISO country code
	countryCode := record.Country.IsoCode
	if countryCode == "" {
		return "", fmt.Errorf("no country code found for IP: %s", ipStr)
	}

	return strings.ToUpper(countryCode), nil
}

// GetCountryFromIPWithFallback resolves an IP to country code, returns empty string on error
func GetCountryFromIPWithFallback(ipStr string) string {
	// Check cache first
	if geoCache != nil {
		if cached, found := geoCache.Get("country:" + ipStr); found {
			fmt.Printf("Country cache HIT for IP %s: %s\n", ipStr, cached.(string))
			return cached.(string)
		}
		fmt.Printf("Country cache MISS for IP %s\n", ipStr)
	}

	countryCode, err := GetCountryFromIP(ipStr)
	if err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Geolocation failed for IP %s: %v\n", ipStr, err)
		// Cache the empty result to avoid repeated lookups for invalid IPs
		if geoCache != nil {
			geoCache.Set("country:"+ipStr, "", cache.DefaultExpiration)
		}
		return ""
	}

	// Cache the result
	if geoCache != nil {
		geoCache.Set("country:"+ipStr, countryCode, cache.DefaultExpiration)
		fmt.Printf("Country cached for IP %s: %s\n", ipStr, countryCode)
	}

	return countryCode
}

// GetASNFromIP resolves an IP address to an ASN
func GetASNFromIP(ipStr string) (string, error) {
	if asnReader == nil {
		return "", fmt.Errorf("ASN database not initialized")
	}

	// Parse the IP address
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address: %s", ipStr)
	}

	// Skip private/local IPs
	if IsPrivateIP(ip) {
		return "", fmt.Errorf("private IP address: %s", ipStr)
	}

	// Look up the ASN
	record, err := asnReader.ASN(ip)
	if err != nil {
		return "", fmt.Errorf("failed to lookup ASN for IP %s: %v", ipStr, err)
	}

	// Return the ASN number
	asnNumber := record.AutonomousSystemNumber
	if asnNumber == 0 {
		return "", fmt.Errorf("no ASN found for IP: %s", ipStr)
	}

	return fmt.Sprintf("AS%d", asnNumber), nil
}

// GetASNFromIPWithFallback resolves an IP to ASN, returns empty string on error
func GetASNFromIPWithFallback(ipStr string) string {
	// Check cache first
	if geoCache != nil {
		if cached, found := geoCache.Get("asn:" + ipStr); found {
			fmt.Printf("ASN cache HIT for IP %s: %s\n", ipStr, cached.(string))
			return cached.(string)
		}
		fmt.Printf("ASN cache MISS for IP %s\n", ipStr)
	}

	asn, err := GetASNFromIP(ipStr)
	if err != nil {
		// Log the error but don't fail the request
		fmt.Printf("ASN lookup failed for IP %s: %v\n", ipStr, err)
		// Cache the empty result to avoid repeated lookups for invalid IPs
		if geoCache != nil {
			geoCache.Set("asn:"+ipStr, "", cache.DefaultExpiration)
		}
		return ""
	}

	// Cache the result
	if geoCache != nil {
		geoCache.Set("asn:"+ipStr, asn, cache.DefaultExpiration)
		fmt.Printf("ASN cached for IP %s: %s\n", ipStr, asn)
	}

	return asn
}

// GetGeoCacheStats returns statistics about the geolocation cache
func GetGeoCacheStats() map[string]interface{} {
	if geoCache == nil {
		return map[string]interface{}{
			"enabled": false,
			"error":   "Cache not initialized",
		}
	}

	return map[string]interface{}{
		"enabled":     true,
		"items_count": geoCache.ItemCount(),
	}
}

// ClearGeoCache clears all cached geolocation data
func ClearGeoCache() {
	if geoCache != nil {
		geoCache.Flush()
	}
}

// GetGeoCacheSize returns the number of cached items
func GetGeoCacheSize() int {
	if geoCache == nil {
		return 0
	}
	return geoCache.ItemCount()
}
