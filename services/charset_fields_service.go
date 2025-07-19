package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// CharsetField represents a field that should be checked for charset detection
type CharsetField struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"` // "standard" or "custom"
}

// CharsetFieldsConfig manages the configuration of which fields to check for charset detection
type CharsetFieldsConfig struct {
	mu             sync.RWMutex
	standardFields []CharsetField
	customFields   []CharsetField
}

var (
	charsetFieldsConfig *CharsetFieldsConfig
	configOnce          sync.Once
)

// GetCharsetFieldsConfig returns the singleton instance of CharsetFieldsConfig
func GetCharsetFieldsConfig() *CharsetFieldsConfig {
	configOnce.Do(func() {
		charsetFieldsConfig = &CharsetFieldsConfig{
			standardFields: []CharsetField{
				{Name: "username", Enabled: true, Type: "standard"},
				{Name: "email", Enabled: true, Type: "standard"},
				{Name: "user_agent", Enabled: true, Type: "standard"},
			},
			customFields: []CharsetField{},
		}
	})
	return charsetFieldsConfig
}

// GetEnabledFields returns all enabled fields (both standard and custom)
func (c *CharsetFieldsConfig) GetEnabledFields() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var enabledFields []string

	// Add enabled standard fields
	for _, field := range c.standardFields {
		if field.Enabled {
			enabledFields = append(enabledFields, field.Name)
		}
	}

	// Add enabled custom fields
	for _, field := range c.customFields {
		if field.Enabled {
			enabledFields = append(enabledFields, field.Name)
		}
	}

	return enabledFields
}

// GetAllFields returns all fields (enabled and disabled)
func (c *CharsetFieldsConfig) GetAllFields() map[string][]CharsetField {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string][]CharsetField{
		"standard": c.standardFields,
		"custom":   c.customFields,
	}
}

// ToggleStandardField enables/disables a standard field
func (c *CharsetFieldsConfig) ToggleStandardField(fieldName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, field := range c.standardFields {
		if field.Name == fieldName {
			c.standardFields[i].Enabled = !field.Enabled
			log.Printf("Standard field '%s' %s", fieldName, map[bool]string{true: "enabled", false: "disabled"}[c.standardFields[i].Enabled])
			c.clearCharsetCache()
			return nil
		}
	}

	return fmt.Errorf("standard field '%s' not found", fieldName)
}

// AddCustomField adds a new custom field
func (c *CharsetFieldsConfig) AddCustomField(fieldName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if field already exists (standard or custom)
	for _, field := range c.standardFields {
		if field.Name == fieldName {
			return fmt.Errorf("field '%s' already exists as a standard field", fieldName)
		}
	}

	for _, field := range c.customFields {
		if field.Name == fieldName {
			return fmt.Errorf("field '%s' already exists as a custom field", fieldName)
		}
	}

	newField := CharsetField{
		Name:    fieldName,
		Enabled: true,
		Type:    "custom",
	}

	c.customFields = append(c.customFields, newField)
	log.Printf("Custom field '%s' added", fieldName)
	c.clearCharsetCache()
	return nil
}

// DeleteCustomField removes a custom field
func (c *CharsetFieldsConfig) DeleteCustomField(fieldName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, field := range c.customFields {
		if field.Name == fieldName {
			c.customFields = append(c.customFields[:i], c.customFields[i+1:]...)
			log.Printf("Custom field '%s' deleted", fieldName)
			c.clearCharsetCache()
			return nil
		}
	}

	return fmt.Errorf("custom field '%s' not found", fieldName)
}

// ToggleCustomField enables/disables a custom field
func (c *CharsetFieldsConfig) ToggleCustomField(fieldName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, field := range c.customFields {
		if field.Name == fieldName {
			c.customFields[i].Enabled = !field.Enabled
			log.Printf("Custom field '%s' %s", fieldName, map[bool]string{true: "enabled", false: "disabled"}[c.customFields[i].Enabled])
			c.clearCharsetCache()
			return nil
		}
	}

	return fmt.Errorf("custom field '%s' not found", fieldName)
}

// GetConfigJSON returns the current configuration as JSON
func (c *CharsetFieldsConfig) GetConfigJSON() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	config := map[string]interface{}{
		"standard_fields": c.standardFields,
		"custom_fields":   c.customFields,
	}

	return json.Marshal(config)
}

// LoadConfigFromJSON loads configuration from JSON
func (c *CharsetFieldsConfig) LoadConfigFromJSON(data []byte) error {
	var config struct {
		StandardFields []CharsetField `json:"standard_fields"`
		CustomFields   []CharsetField `json:"custom_fields"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.standardFields = config.StandardFields
	c.customFields = config.CustomFields

	log.Printf("Charset fields config loaded: %d standard fields, %d custom fields",
		len(c.standardFields), len(c.customFields))
	return nil
}

// clearCharsetCache removes cache items that are affected by charset field changes
func (c *CharsetFieldsConfig) clearCharsetCache() {
	cache := GetCacheFactory()
	if cache == nil {
		log.Printf("Warning: Cache factory not available for charset cache clearing")
		return
	}

	// Get all cache keys and remove those that start with "filter:"
	// This is a simplified approach - in a production environment you might want
	// to use a more sophisticated cache tagging system
	stats, err := cache.Stats()
	if err != nil {
		log.Printf("Warning: Could not get cache stats for charset cache clearing: %v", err)
		return
	}

	// For now, we'll clear all filter cache items since charset changes affect all filter results
	// In a more sophisticated implementation, you could tag cache items with field names
	itemsCleared := 0
	if items, ok := stats["items"].(int); ok && items > 0 {
		// Clear all filter cache items
		// Note: This is a simplified approach. In production, you might want to implement
		// cache tagging to only clear items that contain specific fields
		cache.Clear()
		itemsCleared = items
	}

	log.Printf("Cleared %d cache items due to charset field configuration change", itemsCleared)
}
