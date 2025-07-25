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
)

// IndexIPAddress indexes an IP address to Elasticsearch
func IndexIPAddress(ip models.IP) error {
	es := config.ESClient

	// Create the document to index
	doc := map[string]interface{}{
		"address": ip.Address,
		"status":  ip.Status,
		"is_cidr": ip.IsCIDR,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	// Index the document
	// Use database ID as document ID to avoid issues with special characters in CIDR notation
	docID := fmt.Sprintf("%d", ip.ID)
	req := esapi.IndexRequest{
		Index:      "ip-addresses",
		DocumentID: docID,
		Body:       strings.NewReader(string(docJSON)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing IP: %s", res.String())
		return err
	}

	log.Printf("Successfully indexed IP: %s", ip.Address)
	return nil
}

// IndexEmail indexes an email to Elasticsearch
func IndexEmail(email models.Email) error {
	es := config.ESClient

	doc := map[string]interface{}{
		"email":    email.Address,
		"status":   email.Status,
		"is_regex": email.IsRegex,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	// Use database ID as document ID to avoid issues with special characters in patterns
	docID := fmt.Sprintf("%d", email.ID)
	req := esapi.IndexRequest{
		Index:      "emails",
		DocumentID: docID,
		Body:       strings.NewReader(string(docJSON)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing email: %s", res.String())
		return err
	}

	log.Printf("Successfully indexed email: %s", email.Address)
	return nil
}

// IndexUserAgent indexes a user agent to Elasticsearch
func IndexUserAgent(userAgent models.UserAgent) error {
	es := config.ESClient

	doc := map[string]interface{}{
		"user_agent": userAgent.UserAgent,
		"status":     userAgent.Status,
		"is_regex":   userAgent.IsRegex,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	// Use database ID as document ID to avoid issues with special characters in patterns
	docID := fmt.Sprintf("%d", userAgent.ID)
	req := esapi.IndexRequest{
		Index:      "user-agents",
		DocumentID: docID,
		Body:       strings.NewReader(string(docJSON)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing user agent: %s", res.String())
		return err
	}

	log.Printf("Successfully indexed user agent: %s", userAgent.UserAgent)
	return nil
}

// IndexCountry indexes a country to Elasticsearch
func IndexCountry(country models.Country) error {
	es := config.ESClient

	doc := map[string]interface{}{
		"country": country.Code,
		"status":  country.Status,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "countries",
		DocumentID: country.Code,
		Body:       strings.NewReader(string(docJSON)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing country: %s", res.String())
		return err
	}

	log.Printf("Successfully indexed country: %s", country.Code)
	return nil
}

// IndexCharsetRule indexes a charset rule to Elasticsearch
func IndexCharsetRule(charset models.CharsetRule) error {
	es := config.ESClient

	doc := map[string]interface{}{
		"charset": charset.Charset,
		"status":  charset.Status,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "charsets",
		DocumentID: charset.Charset,
		Body:       strings.NewReader(string(docJSON)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing charset rule: %s", res.String())
		return err
	}

	log.Printf("Successfully indexed charset rule: %s", charset.Charset)
	return nil
}

// IndexUsernameRule indexes a username rule to Elasticsearch
func IndexUsernameRule(username models.UsernameRule) error {
	es := config.ESClient

	doc := map[string]interface{}{
		"username": username.Username,
		"status":   username.Status,
		"is_regex": username.IsRegex,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	// Use database ID as document ID to avoid issues with special characters in patterns
	docID := fmt.Sprintf("%d", username.ID)
	req := esapi.IndexRequest{
		Index:      "usernames",
		DocumentID: docID,
		Body:       strings.NewReader(string(docJSON)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing username rule: %s", res.String())
		return err
	}

	log.Printf("Successfully indexed username rule: %s", username.Username)
	return nil
}

// IndexASN indexes an ASN to Elasticsearch
func IndexASN(asn models.ASN) error {
	es := config.ESClient

	doc := map[string]interface{}{
		"asn":    asn.ASN,
		"rir":    asn.RIR,
		"domain": asn.Domain,
		"cc":     asn.Country,
		"asname": asn.Name,
		"status": asn.Status,
		"source": asn.Source,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	// Use database ID as document ID
	docID := fmt.Sprintf("%d", asn.ID)
	req := esapi.IndexRequest{
		Index:      "asns",
		DocumentID: docID,
		Body:       strings.NewReader(string(docJSON)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing ASN: %s", res.String())
		return err
	}

	log.Printf("Successfully indexed ASN: %s", asn.ASN)
	return nil
}

// SyncASNToES synchronizes an ASN rule to Elasticsearch
func SyncASNToES(asn models.ASN) error {
	return IndexASN(asn)
}

// DeleteASNFromES removes an ASN rule from Elasticsearch
func DeleteASNFromES(asnID uint) error {
	es := config.ESClient

	docID := fmt.Sprintf("%d", asnID)
	req := esapi.DeleteRequest{
		Index:      "asns",
		DocumentID: docID,
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error deleting ASN from Elasticsearch: %s", res.String())
		return err
	}

	log.Printf("Successfully deleted ASN from Elasticsearch: %d", asnID)
	return nil
}

// SyncAllIPs syncs all IP addresses from MySQL to Elasticsearch
func SyncAllIPs() error {
	var ips []models.IP
	if err := config.DB.Find(&ips).Error; err != nil {
		return err
	}

	for _, ip := range ips {
		if err := IndexIPAddress(ip); err != nil {
			log.Printf("Error syncing IP %s: %v", ip.Address, err)
		}
	}

	log.Printf("Synced %d IP addresses to Elasticsearch", len(ips))
	return nil
}

// SyncAllEmails syncs all emails from MySQL to Elasticsearch
func SyncAllEmails() error {
	var emails []models.Email
	if err := config.DB.Find(&emails).Error; err != nil {
		return err
	}

	for _, email := range emails {
		if err := IndexEmail(email); err != nil {
			log.Printf("Error syncing email %s: %v", email.Address, err)
		}
	}

	log.Printf("Synced %d emails to Elasticsearch", len(emails))
	return nil
}

// SyncAllUserAgents syncs all user agents from MySQL to Elasticsearch
func SyncAllUserAgents() error {
	var userAgents []models.UserAgent
	if err := config.DB.Find(&userAgents).Error; err != nil {
		return err
	}

	for _, userAgent := range userAgents {
		if err := IndexUserAgent(userAgent); err != nil {
			log.Printf("Error syncing user agent %s: %v", userAgent.UserAgent, err)
		}
	}

	log.Printf("Synced %d user agents to Elasticsearch", len(userAgents))
	return nil
}

// SyncAllCountries syncs all countries from MySQL to Elasticsearch
func SyncAllCountries() error {
	var countries []models.Country
	if err := config.DB.Find(&countries).Error; err != nil {
		return err
	}

	for _, country := range countries {
		if err := IndexCountry(country); err != nil {
			log.Printf("Error syncing country %s: %v", country.Code, err)
		}
	}

	log.Printf("Synced %d countries to Elasticsearch", len(countries))
	return nil
}

// SyncAllCharsetRules syncs all charset rules from MySQL to Elasticsearch
func SyncAllCharsetRules() error {
	var charsets []models.CharsetRule
	if err := config.DB.Find(&charsets).Error; err != nil {
		return err
	}

	for _, charset := range charsets {
		if err := IndexCharsetRule(charset); err != nil {
			log.Printf("Error syncing charset rule %s: %v", charset.Charset, err)
		}
	}

	log.Printf("Synced %d charset rules to Elasticsearch", len(charsets))
	return nil
}

// SyncAllUsernameRules syncs all username rules from MySQL to Elasticsearch
func SyncAllUsernameRules() error {
	var usernames []models.UsernameRule
	if err := config.DB.Find(&usernames).Error; err != nil {
		return err
	}

	for _, username := range usernames {
		if err := IndexUsernameRule(username); err != nil {
			log.Printf("Error syncing username rule %s: %v", username.Username, err)
		}
	}

	log.Printf("Synced %d username rules to Elasticsearch", len(usernames))
	return nil
}

// SyncAllASNs syncs all ASNs from MySQL to Elasticsearch
func SyncAllASNs() error {
	var asns []models.ASN
	if err := config.DB.Find(&asns).Error; err != nil {
		return err
	}

	for _, asn := range asns {
		if err := IndexASN(asn); err != nil {
			log.Printf("Error syncing ASN %s: %v", asn.ASN, err)
		}
	}

	log.Printf("Synced %d ASNs to Elasticsearch", len(asns))
	return nil
}

// SyncAllData syncs all data from MySQL to Elasticsearch
func SyncAllData() error {
	log.Println("Starting full data sync to Elasticsearch...")

	if err := SyncAllIPs(); err != nil {
		log.Printf("Error syncing IPs: %v", err)
	}

	if err := SyncAllEmails(); err != nil {
		log.Printf("Error syncing emails: %v", err)
	}

	if err := SyncAllUserAgents(); err != nil {
		log.Printf("Error syncing user agents: %v", err)
	}

	if err := SyncAllCountries(); err != nil {
		log.Printf("Error syncing countries: %v", err)
	}

	if err := SyncAllCharsetRules(); err != nil {
		log.Printf("Error syncing charset rules: %v", err)
	}

	if err := SyncAllUsernameRules(); err != nil {
		log.Printf("Error syncing username rules: %v", err)
	}

	if err := SyncAllASNs(); err != nil {
		log.Printf("Error syncing ASNs: %v", err)
	}

	log.Println("Full data sync completed")
	return nil
}

// DeleteIPIndex deletes the IP index from Elasticsearch
func DeleteIPIndex() error {
	es := config.ESClient
	req := esapi.IndicesDeleteRequest{
		Index: []string{"ips"},
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting IP index: %s", res.String())
	}

	log.Println("IP index deleted successfully")
	return nil
}

// DeleteEmailIndex deletes the email index from Elasticsearch
func DeleteEmailIndex() error {
	es := config.ESClient
	req := esapi.IndicesDeleteRequest{
		Index: []string{"emails"},
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting email index: %s", res.String())
	}

	log.Println("Email index deleted successfully")
	return nil
}

// DeleteUserAgentIndex deletes the user agent index from Elasticsearch
func DeleteUserAgentIndex() error {
	es := config.ESClient
	req := esapi.IndicesDeleteRequest{
		Index: []string{"user-agents"},
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting user agent index: %s", res.String())
	}

	log.Println("User agent index deleted successfully")
	return nil
}

// DeleteCountryIndex deletes the country index from Elasticsearch
func DeleteCountryIndex() error {
	es := config.ESClient
	req := esapi.IndicesDeleteRequest{
		Index: []string{"countries"},
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting country index: %s", res.String())
	}

	log.Println("Country index deleted successfully")
	return nil
}

// DeleteCharsetIndex deletes the charset index from Elasticsearch
func DeleteCharsetIndex() error {
	es := config.ESClient
	req := esapi.IndicesDeleteRequest{
		Index: []string{"charsets"},
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting charset index: %s", res.String())
	}

	log.Println("Charset index deleted successfully")
	return nil
}

// DeleteUsernameIndex deletes the username index from Elasticsearch
func DeleteUsernameIndex() error {
	es := config.ESClient
	req := esapi.IndicesDeleteRequest{
		Index: []string{"usernames"},
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting username index: %s", res.String())
	}

	log.Println("Username index deleted successfully")
	return nil
}

// DeleteASNIndex deletes the ASN index from Elasticsearch
func DeleteASNIndex() error {
	es := config.ESClient
	req := esapi.IndicesDeleteRequest{
		Index: []string{"asns"},
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting ASN index: %s", res.String())
	}

	log.Println("ASN index deleted successfully")
	return nil
}
