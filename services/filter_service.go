package services

import (
	"context"
	"encoding/json"
	"firewall/config"
	"firewall/models"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
)

// FilterResult defines the structure of the response for filtering
type FilterResult struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

const NumFilters = 4

// EvaluateFilters runs all filters concurrently and returns the final result
func EvaluateFilters(ctx context.Context, ip, email, userAgent, country string) (FilterResult, error) {
	// Channel to collect filter results
	results := make(chan FilterResult, 4)

	// Start filters concurrently
	go filterIP(ctx, ip, results)
	go filterEmail(ctx, email, results)
	go filterUserAgent(ctx, userAgent, results)
	go filterCountry(ctx, country, results)

	// Collect and evaluate the results
	return collectResults(ctx, results)
}

// collectResults processes all filter results
func collectResults(ctx context.Context, result chan FilterResult) (FilterResult, error) {
	output := FilterResult{Status: "allowed"}

	for i := 0; i < NumFilters; i++ {
		select {
		case res := <-result:
			if res.Status == "whitelisted" {
				return res, nil
			} else if res.Status == "denied" {
				output = res // Update output to the denied result
			}
		case <-ctx.Done():
			return FilterResult{Status: "timeout"}, ctx.Err()
		}
	}

	return output, nil
}

// filterIP runs the IP filter
func filterIP(ctx context.Context, ip string, result chan FilterResult) {
	es := config.ESClient
	query := `{
		"query": {
			"match": {
				"address": "` + ip + `"
			}
		}
	}`

	req := esapi.SearchRequest{
		Index: []string{"ip-addresses"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Type: "IP", Status: "error"}
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Type: "IP", Status: "error"}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		totalHits := hits["total"].(map[string]interface{})["value"].(float64)
		if totalHits > 0 {
			firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
			source := firstHit["_source"].(map[string]interface{})
			status := source["status"].(string)

			if status == "denied" {
				result <- FilterResult{Type: "IP", Status: "denied"}
			} else if status == "whitelisted" {
				result <- FilterResult{Type: "IP", Status: "whitelisted"}
			} else {
				result <- FilterResult{Type: "IP", Status: "allowed"}
			}
		} else {
			result <- FilterResult{Type: "IP", Status: "allowed"}
		}
	} else {
		result <- FilterResult{Type: "IP", Status: "error"}
	}
}

// filterEmail runs the email filter
func filterEmail(ctx context.Context, email string, result chan FilterResult) {
	es := config.ESClient
	query := `{
		"query": {
			"match": {
				"email": "` + email + `"
			}
		}
	}`

	req := esapi.SearchRequest{
		Index: []string{"emails"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Type: "Email", Status: "error"}
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Type: "Email", Status: "error"}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		totalHits := hits["total"].(map[string]interface{})["value"].(float64)
		if totalHits > 0 {
			firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
			source := firstHit["_source"].(map[string]interface{})
			status := source["status"].(string)

			if status == "denied" {
				result <- FilterResult{Type: "Email", Status: "denied"}
			} else if status == "whitelisted" {
				result <- FilterResult{Type: "Email", Status: "whitelisted"}
			} else {
				result <- FilterResult{Type: "Email", Status: "allowed"}
			}
		} else {
			result <- FilterResult{Type: "Email", Status: "allowed"}
		}
	} else {
		result <- FilterResult{Type: "Email", Status: "error"}
	}
}

// filterUserAgent runs the user agent filter
func filterUserAgent(ctx context.Context, userAgent string, result chan FilterResult) {
	es := config.ESClient
	query := `{
		"query": {
			"match": {
				"user_agent": "` + userAgent + `"
			}
		}
	}`

	req := esapi.SearchRequest{
		Index: []string{"user-agents"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(ctx, es)
	if err != nil {
		result <- FilterResult{Type: "UserAgent", Status: "error"}
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Type: "UserAgent", Status: "error"}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		totalHits := hits["total"].(map[string]interface{})["value"].(float64)
		if totalHits > 0 {
			firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
			source := firstHit["_source"].(map[string]interface{})
			status := source["status"].(string)

			if status == "denied" {
				result <- FilterResult{Type: "UserAgent", Status: "denied"}
			} else if status == "whitelisted" {
				result <- FilterResult{Type: "UserAgent", Status: "whitelisted"}
			} else {
				result <- FilterResult{Type: "UserAgent", Status: "allowed"}
			}
		} else {
			result <- FilterResult{Type: "UserAgent", Status: "allowed"}
		}
	} else {
		result <- FilterResult{Type: "UserAgent", Status: "error"}
	}
}

// filterCountry runs the country filter
func filterCountry(ctx context.Context, country string, result chan FilterResult) {
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
		result <- FilterResult{Type: "Country", Status: "error"}
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		result <- FilterResult{Type: "Country", Status: "error"}
		return
	}

	if hits, found := r["hits"].(map[string]interface{}); found {
		totalHits := hits["total"].(map[string]interface{})["value"].(float64)
		if totalHits > 0 {
			firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
			source := firstHit["_source"].(map[string]interface{})
			status := source["status"].(string)

			if status == "denied" {
				result <- FilterResult{Type: "Country", Status: "denied"}
			} else if status == "whitelisted" {
				result <- FilterResult{Type: "Country", Status: "whitelisted"}
			} else {
				result <- FilterResult{Type: "Country", Status: "allowed"}
			}
		} else {
			result <- FilterResult{Type: "Country", Status: "allowed"}
		}
	} else {
		result <- FilterResult{Type: "Country", Status: "error"}
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
		strings.NewReader(fmt.Sprintf(`{"id": %d, "username": "%s", "status": "%s"}`, username.ID, username.Username, username.Status)),
		config.ESClient.Index.WithDocumentID(fmt.Sprintf("%d", username.ID)),
		config.ESClient.Index.WithContext(ctx),
	)
	if err == nil {
		log.Printf("Indexed Username: %d %s %s", username.ID, username.Username, username.Status)
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
