package services

import (
	"context"
	"encoding/json"
	"firewall/config"
	"firewall/models"
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
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	// Index the document
	req := esapi.IndexRequest{
		Index:      "ip-addresses",
		DocumentID: ip.Address, // Use IP as document ID
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
		"email":  email.Address,
		"status": email.Status,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "emails",
		DocumentID: email.Address,
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
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "user-agents",
		DocumentID: userAgent.UserAgent,
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
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "usernames",
		DocumentID: username.Username,
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

	log.Println("Full data sync completed")
	return nil
}
