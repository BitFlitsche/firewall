package services

import (
	"context"
	"encoding/json"
	"firewall/config"
	"firewall/models"
	"firewall/utils"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
)

// FilterResult defines the structure of the response for filtering
// Vereinheitlicht: result, reason, field, value
type FilterResult struct {
	Result string      `json:"result"`
	Reason string      `json:"reason,omitempty"`
	Field  string      `json:"field,omitempty"`
	Value  interface{} `json:"value,omitempty"`
}

const NumFilters = 6

// FilterResultWithResolvedData includes the resolved country and ASN values
type FilterResultWithResolvedData struct {
	FilterResult
	ResolvedCountry string
	ResolvedASN     string
}

// EvaluateFilters runs only the necessary filters concurrently and returns the final result with resolved data
func EvaluateFilters(ctx context.Context, ip, email, userAgent, country, asn, username string) (FilterResultWithResolvedData, error) {
	// Auto-geolocate IP if country is empty and IP is provided
	resolvedCountry := country
	if country == "" && ip != "" {
		resolvedCountry = GetCountryFromIPWithFallback(ip)
		// Log the geolocation result for debugging
		if resolvedCountry != "" {
			fmt.Printf("Auto-geolocated IP %s to country: %s\n", ip, resolvedCountry)
		}
	}

	// Auto-resolve ASN if not provided but IP is available
	resolvedASN := asn
	if asn == "" && ip != "" {
		resolvedASN = GetASNFromIPWithFallback(ip)
	}

	// Determine which filters to run based on non-empty fields
	var filtersToRun int
	var results chan FilterResult

	// Count how many filters we need to run
	if ip != "" {
		filtersToRun++
	}
	if email != "" {
		filtersToRun++
	}
	if userAgent != "" {
		filtersToRun++
	}
	if country != "" {
		filtersToRun++
	}
	if username != "" {
		filtersToRun++
	}
	if asn != "" || ip != "" { // ASN filter runs if either ASN is provided or IP is provided (for auto-ASN lookup)
		filtersToRun++
	}

	// If no filters to run, return allowed
	if filtersToRun == 0 {
		return FilterResultWithResolvedData{
			FilterResult:    FilterResult{Result: "allowed", Reason: "no filter fields provided"},
			ResolvedCountry: resolvedCountry,
			ResolvedASN:     resolvedASN,
		}, nil
	}

	// Create channel with exact size needed
	results = make(chan FilterResult, filtersToRun)

	// Start only the necessary filters concurrently
	if ip != "" {
		go filterIP(ctx, ip, results)
	}
	if email != "" {
		go filterEmail(ctx, email, results)
	}
	if userAgent != "" {
		go filterUserAgent(ctx, userAgent, results)
	}
	if resolvedCountry != "" {
		go filterCountry(ctx, resolvedCountry, results)
	}
	if username != "" {
		go filterUsername(ctx, username, results)
	}
	if asn != "" || ip != "" {
		go filterASN(ctx, ip, resolvedASN, results)
	}

	// Collect and evaluate the results
	filterResult, err := collectResults(ctx, results, filtersToRun)
	return FilterResultWithResolvedData{
		FilterResult:    filterResult,
		ResolvedCountry: resolvedCountry,
		ResolvedASN:     resolvedASN,
	}, err
}

// collectResults processes filter results with dynamic count
func collectResults(ctx context.Context, result chan FilterResult, filterCount int) (FilterResult, error) {
	output := FilterResult{Result: "allowed"}

	for i := 0; i < filterCount; i++ {
		select {
		case res := <-result:
			if res.Result == "whitelisted" {
				return res, nil
			} else if res.Result == "denied" {
				output = res // Update output to the denied result
			}
		case <-ctx.Done():
			return FilterResult{Result: "timeout", Reason: "timeout"}, ctx.Err()
		}
	}

	return output, nil
}

// filterIP runs the IP filter
func filterIP(ctx context.Context, ip string, result chan FilterResult) {
	// Handle empty IP addresses - treat as invalid input
	if ip == "" {
		result <- FilterResult{Result: "allowed", Reason: "empty ip address", Field: "ip", Value: ip}
		return
	}

	es := config.ESClient

	// First check for exact IP matches (including existing data without is_cidr field)
	exactQuery := `{
		"query": {
			"bool": {
				"should": [
					{"bool": {"must": [{"term": {"address.keyword": "` + ip + `"}}, {"term": {"is_cidr": false}}]}},
					{"bool": {"must": [{"term": {"address.keyword": "` + ip + `"}}, {"bool": {"must_not": {"exists": {"field": "is_cidr"}}}}]}}
				],
				"minimum_should_match": 1
			}
		}
	}`

	req := esapi.SearchRequest{
		Index: []string{"ip-addresses"},
		Body:  strings.NewReader(exactQuery),
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Result: "error", Reason: "elasticsearch error", Field: "ip", Value: ip}
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Result: "error", Reason: "decode error", Field: "ip", Value: ip}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		totalHits := hits["total"].(map[string]interface{})["value"].(float64)
		if totalHits > 0 {
			firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
			source := firstHit["_source"].(map[string]interface{})
			status := source["status"].(string)

			if status == "denied" {
				result <- FilterResult{Result: "denied", Reason: "ip denied", Field: "ip", Value: ip}
			} else if status == "whitelisted" {
				result <- FilterResult{Result: "whitelisted", Reason: "ip whitelisted", Field: "ip", Value: ip}
			} else {
				result <- FilterResult{Result: "allowed", Field: "ip", Value: ip}
			}
			return
		}
	}

	// If no exact match, check CIDR blocks
	cidrQuery := `{
		"query": {
			"bool": {
				"must": [
					{"term": {"is_cidr": true}}
				]
			}
		}
	}`

	req = esapi.SearchRequest{
		Index: []string{"ip-addresses"},
		Body:  strings.NewReader(cidrQuery),
	}

	res, err = req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Result: "error", Reason: "elasticsearch error", Field: "ip", Value: ip}
		return
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Result: "error", Reason: "decode error", Field: "ip", Value: ip}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		hitsList := hits["hits"].([]interface{})
		for _, hit := range hitsList {
			hitMap := hit.(map[string]interface{})
			source := hitMap["_source"].(map[string]interface{})
			cidrBlock := source["address"].(string)
			status := source["status"].(string)

			// Check if IP falls within this CIDR block
			inRange, err := utils.IsIPInCIDR(ip, cidrBlock)
			if err != nil {
				continue // Skip invalid CIDR blocks
			}
			if inRange {
				if status == "denied" {
					result <- FilterResult{Result: "denied", Reason: "ip cidr denied", Field: "ip", Value: ip}
				} else if status == "whitelisted" {
					result <- FilterResult{Result: "whitelisted", Reason: "ip cidr whitelisted", Field: "ip", Value: ip}
				} else {
					result <- FilterResult{Result: "allowed", Field: "ip", Value: ip}
				}
				return
			}
		}
	}

	// No matches found
	result <- FilterResult{Result: "allowed", Field: "ip", Value: ip}
}

// filterEmail runs the email filter
func filterEmail(ctx context.Context, email string, result chan FilterResult) {
	// Handle empty email addresses - treat as allowed
	if email == "" {
		result <- FilterResult{Result: "allowed", Reason: "empty email address", Field: "email", Value: email}
		return
	}

	es := config.ESClient

	// First try exact match
	exactQuery := `{
		"query": {
			"bool": {
				"must": [
					{"term": {"email.keyword": "` + email + `"}},
					{"term": {"is_regex": false}}
				]
			}
		}
	}`

	req := esapi.SearchRequest{
		Index: []string{"emails"},
		Body:  strings.NewReader(exactQuery),
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Result: "error", Reason: "elasticsearch error", Field: "email", Value: email}
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Result: "error", Reason: "decode error", Field: "email", Value: email}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		totalHits := hits["total"].(map[string]interface{})["value"].(float64)
		if totalHits > 0 {
			firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
			source := firstHit["_source"].(map[string]interface{})
			status := source["status"].(string)

			if status == "denied" {
				result <- FilterResult{Result: "denied", Reason: "email denied", Field: "email", Value: email}
			} else if status == "whitelisted" {
				result <- FilterResult{Result: "whitelisted", Reason: "email whitelisted", Field: "email", Value: email}
			} else {
				result <- FilterResult{Result: "allowed", Field: "email", Value: email}
			}
			return
		}
	}

	// If no exact match, try regex patterns
	regexQuery := `{
		"query": {
			"bool": {
				"must": [
					{"term": {"is_regex": true}}
				]
			}
		}
	}`

	req = esapi.SearchRequest{
		Index: []string{"emails"},
		Body:  strings.NewReader(regexQuery),
	}

	res, err = req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Result: "error", Reason: "elasticsearch error", Field: "email", Value: email}
		return
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Result: "error", Reason: "decode error", Field: "email", Value: email}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		hitsList := hits["hits"].([]interface{})
		for _, hit := range hitsList {
			hitMap := hit.(map[string]interface{})
			source := hitMap["_source"].(map[string]interface{})
			pattern := source["email"].(string)
			status := source["status"].(string)

			// Check if email matches the regex pattern
			matched, err := regexp.MatchString(pattern, email)
			if err != nil {
				continue // Skip invalid regex patterns
			}
			if matched {
				if status == "denied" {
					result <- FilterResult{Result: "denied", Reason: "email regex denied", Field: "email", Value: email}
				} else if status == "whitelisted" {
					result <- FilterResult{Result: "whitelisted", Reason: "email regex whitelisted", Field: "email", Value: email}
				} else {
					result <- FilterResult{Result: "allowed", Field: "email", Value: email}
				}
				return
			}
		}
	}

	// No matches found
	result <- FilterResult{Result: "allowed", Field: "email", Value: email}
}

// filterUserAgent runs the user agent filter
func filterUserAgent(ctx context.Context, userAgent string, result chan FilterResult) {
	// Handle empty user agent strings - treat as allowed
	if userAgent == "" {
		result <- FilterResult{Result: "allowed", Reason: "empty user agent", Field: "user_agent", Value: userAgent}
		return
	}

	es := config.ESClient

	// First try exact match
	exactQuery := `{
		"query": {
			"bool": {
				"must": [
					{"term": {"user_agent.keyword": "` + userAgent + `"}},
					{"term": {"is_regex": false}}
				]
			}
		}
	}`

	req := esapi.SearchRequest{
		Index: []string{"user-agents"},
		Body:  strings.NewReader(exactQuery),
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Result: "error", Reason: "elasticsearch error", Field: "user_agent", Value: userAgent}
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Result: "error", Reason: "decode error", Field: "user_agent", Value: userAgent}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		totalHits := hits["total"].(map[string]interface{})["value"].(float64)
		if totalHits > 0 {
			firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
			source := firstHit["_source"].(map[string]interface{})
			status := source["status"].(string)

			if status == "denied" {
				result <- FilterResult{Result: "denied", Reason: "user_agent denied", Field: "user_agent", Value: userAgent}
			} else if status == "whitelisted" {
				result <- FilterResult{Result: "whitelisted", Reason: "user_agent whitelisted", Field: "user_agent", Value: userAgent}
			} else {
				result <- FilterResult{Result: "allowed", Field: "user_agent", Value: userAgent}
			}
			return
		}
	}

	// If no exact match, try regex patterns
	regexQuery := `{
		"query": {
			"bool": {
				"must": [
					{"term": {"is_regex": true}}
				]
			}
		}
	}`

	req = esapi.SearchRequest{
		Index: []string{"user-agents"},
		Body:  strings.NewReader(regexQuery),
	}

	res, err = req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Result: "error", Reason: "elasticsearch error", Field: "user_agent", Value: userAgent}
		return
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Result: "error", Reason: "decode error", Field: "user_agent", Value: userAgent}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		hitsList := hits["hits"].([]interface{})
		for _, hit := range hitsList {
			hitMap := hit.(map[string]interface{})
			source := hitMap["_source"].(map[string]interface{})
			pattern := source["user_agent"].(string)
			status := source["status"].(string)

			// Check if user agent matches the regex pattern
			matched, err := regexp.MatchString(pattern, userAgent)
			if err != nil {
				continue // Skip invalid regex patterns
			}
			if matched {
				if status == "denied" {
					result <- FilterResult{Result: "denied", Reason: "user_agent regex denied", Field: "user_agent", Value: userAgent}
				} else if status == "whitelisted" {
					result <- FilterResult{Result: "whitelisted", Reason: "user_agent regex whitelisted", Field: "user_agent", Value: userAgent}
				} else {
					result <- FilterResult{Result: "allowed", Field: "user_agent", Value: userAgent}
				}
				return
			}
		}
	}

	// No matches found
	result <- FilterResult{Result: "allowed", Field: "user_agent", Value: userAgent}
}

// filterCountry runs the country filter
func filterCountry(ctx context.Context, country string, result chan FilterResult) {
	// Handle empty country codes - treat as allowed
	if country == "" {
		result <- FilterResult{Result: "allowed", Reason: "empty country code", Field: "country", Value: country}
		return
	}

	es := config.ESClient
	query := `{
		"query": {
			"match": {
				"country": "` + country + `"
			}
		}
	}`

	req := esapi.SearchRequest{
		Index: []string{"countries"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Result: "error", Reason: "elasticsearch error", Field: "country", Value: country}
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Result: "error", Reason: "decode error", Field: "country", Value: country}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		totalHits := hits["total"].(map[string]interface{})["value"].(float64)
		if totalHits > 0 {
			firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
			source := firstHit["_source"].(map[string]interface{})
			status := source["status"].(string)

			if status == "denied" {
				result <- FilterResult{Result: "denied", Reason: "country denied", Field: "country", Value: country}
			} else if status == "whitelisted" {
				result <- FilterResult{Result: "whitelisted", Reason: "country whitelisted", Field: "country", Value: country}
			} else {
				result <- FilterResult{Result: "allowed", Field: "country", Value: country}
			}
		} else {
			result <- FilterResult{Result: "allowed", Field: "country", Value: country}
		}
	} else {
		result <- FilterResult{Result: "error", Reason: "no hits", Field: "country", Value: country}
	}
}

// filterUsername runs the username filter
func filterUsername(ctx context.Context, username string, result chan FilterResult) {
	// Handle empty usernames - treat as allowed
	if username == "" {
		result <- FilterResult{Result: "allowed", Reason: "empty username", Field: "username", Value: username}
		return
	}

	es := config.ESClient

	// First try exact match
	exactQuery := `{
		"query": {
			"bool": {
				"must": [
					{"term": {"username.keyword": "` + username + `"}},
					{"term": {"is_regex": false}}
				]
			}
		}
	}`

	req := esapi.SearchRequest{
		Index: []string{"usernames"},
		Body:  strings.NewReader(exactQuery),
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Result: "error", Reason: "elasticsearch error", Field: "username", Value: username}
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Result: "error", Reason: "decode error", Field: "username", Value: username}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		totalHits := hits["total"].(map[string]interface{})["value"].(float64)
		if totalHits > 0 {
			firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
			source := firstHit["_source"].(map[string]interface{})
			status := source["status"].(string)

			if status == "denied" {
				result <- FilterResult{Result: "denied", Reason: "username denied", Field: "username", Value: username}
			} else if status == "whitelisted" {
				result <- FilterResult{Result: "whitelisted", Reason: "username whitelisted", Field: "username", Value: username}
			} else {
				result <- FilterResult{Result: "allowed", Field: "username", Value: username}
			}
			return
		}
	}

	// If no exact match, try regex patterns
	regexQuery := `{
		"query": {
			"bool": {
				"must": [
					{"term": {"is_regex": true}}
				]
			}
		}
	}`

	req = esapi.SearchRequest{
		Index: []string{"usernames"},
		Body:  strings.NewReader(regexQuery),
	}

	res, err = req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Result: "error", Reason: "elasticsearch error", Field: "username", Value: username}
		return
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Result: "error", Reason: "decode error", Field: "username", Value: username}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		hitsList := hits["hits"].([]interface{})
		for _, hit := range hitsList {
			hitMap := hit.(map[string]interface{})
			source := hitMap["_source"].(map[string]interface{})
			pattern := source["username"].(string)
			status := source["status"].(string)

			// Check if username matches the regex pattern
			matched, err := regexp.MatchString(pattern, username)
			if err != nil {
				continue // Skip invalid regex patterns
			}
			if matched {
				if status == "denied" {
					result <- FilterResult{Result: "denied", Reason: "username regex denied", Field: "username", Value: username}
				} else if status == "whitelisted" {
					result <- FilterResult{Result: "whitelisted", Reason: "username regex whitelisted", Field: "username", Value: username}
				} else {
					result <- FilterResult{Result: "allowed", Field: "username", Value: username}
				}
				return
			}
		}
	}

	// No matches found
	result <- FilterResult{Result: "allowed", Field: "username", Value: username}
}

// filterASN runs the ASN filter
func filterASN(ctx context.Context, ip, asn string, result chan FilterResult) {
	// If ASN is provided manually, use it
	if asn != "" {
		// Validate ASN format
		if len(asn) < 3 || !strings.HasPrefix(asn, "AS") {
			result <- FilterResult{Result: "error", Reason: "invalid asn format", Field: "asn", Value: asn}
			return
		}
	} else if ip != "" {
		// ASN should already be resolved in the main function
		// If we still don't have an ASN, just allow
		if asn == "" {
			result <- FilterResult{Result: "allowed", Reason: "no asn found", Field: "asn", Value: ip}
			return
		}
	} else {
		// No IP or ASN provided
		result <- FilterResult{Result: "allowed", Reason: "no ip or asn provided", Field: "asn", Value: ""}
		return
	}

	es := config.ESClient

	// Query ASN rules
	query := `{
		"query": {
			"match": {
				"asn": "` + asn + `"
			}
		}
	}`

	req := esapi.SearchRequest{
		Index: []string{"asns"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Result: "error", Reason: "elasticsearch error", Field: "asn", Value: asn}
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Result: "error", Reason: "decode error", Field: "asn", Value: asn}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		totalHits := hits["total"].(map[string]interface{})["value"].(float64)
		if totalHits > 0 {
			firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
			source := firstHit["_source"].(map[string]interface{})
			status := source["status"].(string)

			if status == "denied" {
				result <- FilterResult{Result: "denied", Reason: "asn denied", Field: "asn", Value: asn}
			} else if status == "whitelisted" {
				result <- FilterResult{Result: "whitelisted", Reason: "asn whitelisted", Field: "asn", Value: asn}
			} else {
				result <- FilterResult{Result: "allowed", Field: "asn", Value: asn}
			}
		} else {
			result <- FilterResult{Result: "allowed", Field: "asn", Value: asn}
		}
	} else {
		result <- FilterResult{Result: "error", Reason: "no hits", Field: "asn", Value: asn}
	}
}

// SyncCharsetToES synchronisiert eine CharsetRule zu Elasticsearch
func SyncCharsetToES(charset models.CharsetRule) error {
	ctx := context.Background()
	_, err := config.ESClient.Index(
		"charsets",
		strings.NewReader(fmt.Sprintf(`{"id": %d, "charset": "%s", "status": "%s"}`, charset.ID, charset.Charset, charset.Status)),
		config.ESClient.Index.WithDocumentID(fmt.Sprintf("%d", charset.ID)),
		config.ESClient.Index.WithContext(ctx),
	)
	if err == nil {
		log.Printf("Indexed Charset: %d %s %s", charset.ID, charset.Charset, charset.Status)
	}
	return err
}

// DeleteCharsetFromES entfernt eine CharsetRule aus Elasticsearch
func DeleteCharsetFromES(id uint) error {
	ctx := context.Background()
	_, err := config.ESClient.Delete(
		"charsets",
		fmt.Sprintf("%d", id),
		config.ESClient.Delete.WithContext(ctx),
	)
	return err
}

// SyncAllCharsetsToES synchronisiert alle CharsetRules nach Elasticsearch
func SyncAllCharsetsToES(db *gorm.DB) error {
	var charsets []models.CharsetRule
	if err := db.Find(&charsets).Error; err != nil {
		return err
	}
	for _, c := range charsets {
		_ = SyncCharsetToES(c)
	}
	log.Printf("Synced %d charsets to Elasticsearch", len(charsets))
	return nil
}

// SyncUsernameToES synchronisiert eine UsernameRule zu Elasticsearch
func SyncUsernameToES(username models.UsernameRule) error {
	ctx := context.Background()
	_, err := config.ESClient.Index(
		"usernames",
		strings.NewReader(fmt.Sprintf(`{"id": %d, "username": "%s", "status": "%s", "is_regex": %t}`, username.ID, username.Username, username.Status, username.IsRegex)),
		config.ESClient.Index.WithDocumentID(fmt.Sprintf("%d", username.ID)),
		config.ESClient.Index.WithContext(ctx),
	)
	if err == nil {
		log.Printf("Indexed Username: %d %s %s %t", username.ID, username.Username, username.Status, username.IsRegex)
	}
	return err
}

// DeleteUsernameFromES entfernt eine UsernameRule aus Elasticsearch
func DeleteUsernameFromES(id uint) error {
	ctx := context.Background()
	_, err := config.ESClient.Delete(
		"usernames",
		fmt.Sprintf("%d", id),
		config.ESClient.Delete.WithContext(ctx),
	)
	return err
}

// SyncAllUsernamesToES synchronisiert alle UsernameRules nach Elasticsearch
func SyncAllUsernamesToES(db *gorm.DB) error {
	var usernames []models.UsernameRule
	if err := db.Find(&usernames).Error; err != nil {
		return err
	}
	for _, u := range usernames {
		_ = SyncUsernameToES(u)
	}
	log.Printf("Synced %d usernames to Elasticsearch", len(usernames))
	return nil
}

// Event-Handler für Charset-Events
func HandleCharsetEvent(action string, data interface{}) {
	charset, ok := data.(models.CharsetRule)
	if !ok {
		return
	}
	switch action {
	case "created", "updated":
		_ = SyncCharsetToES(charset)
	case "deleted":
		_ = DeleteCharsetFromES(charset.ID)
	}
}

// Event-Handler für Username-Events
func HandleUsernameEvent(action string, data interface{}) {
	username, ok := data.(models.UsernameRule)
	if !ok {
		return
	}
	switch action {
	case "created", "updated":
		_ = SyncUsernameToES(username)
	case "deleted":
		_ = DeleteUsernameFromES(username.ID)
	}
}
