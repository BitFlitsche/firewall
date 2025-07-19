package services

import (
	"fmt"
	"log"

	"firewall/models"

	"gorm.io/gorm"
)

// SeedCharsets seeds the charset_rules table with predefined charset rules
func SeedCharsets(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.CharsetRule{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check charsets count: %w", err)
	}
	if count > 0 {
		log.Println("Charset rules table already populated, skipping seed")
		return nil
	}

	// Define charset rules with their associated languages
	charsetRules := []models.CharsetRule{
		// Basic ASCII and Latin scripts
		{Charset: "ASCII", Status: "allowed"},
		{Charset: "Latin", Status: "allowed"},
		{Charset: "Vietnamese", Status: "allowed"},

		// Cyrillic scripts (Russian, Ukrainian, Bulgarian, Serbian, etc.)
		{Charset: "Cyrillic", Status: "allowed"},

		// Arabic scripts (Arabic, Persian, Urdu, etc.)
		{Charset: "Arabic", Status: "allowed"},

		// Hebrew script
		{Charset: "Hebrew", Status: "allowed"},

		// Greek script
		{Charset: "Greek", Status: "allowed"},

		// South Asian scripts
		{Charset: "Devanagari", Status: "allowed"}, // Hindi, Sanskrit, Marathi, etc.
		{Charset: "Bengali", Status: "allowed"},    // Bengali, Assamese
		{Charset: "Tamil", Status: "allowed"},      // Tamil
		{Charset: "Telugu", Status: "allowed"},     // Telugu
		{Charset: "Kannada", Status: "allowed"},    // Kannada
		{Charset: "Malayalam", Status: "allowed"},  // Malayalam
		{Charset: "Gujarati", Status: "allowed"},   // Gujarati
		{Charset: "Gurmukhi", Status: "allowed"},   // Punjabi
		{Charset: "Oriya", Status: "allowed"},      // Odia
		{Charset: "Sinhala", Status: "allowed"},    // Sinhala

		// Southeast Asian scripts
		{Charset: "Thai", Status: "allowed"},    // Thai
		{Charset: "Lao", Status: "allowed"},     // Lao
		{Charset: "Khmer", Status: "allowed"},   // Khmer
		{Charset: "Myanmar", Status: "allowed"}, // Burmese

		// East Asian scripts
		{Charset: "Chinese", Status: "allowed"},  // Chinese (Simplified & Traditional)
		{Charset: "Japanese", Status: "allowed"}, // Japanese (Hiragana, Katakana, Kanji)
		{Charset: "Korean", Status: "allowed"},   // Korean (Hangul)

		// Other scripts
		{Charset: "Armenian", Status: "allowed"},  // Armenian
		{Charset: "Georgian", Status: "allowed"},  // Georgian
		{Charset: "Ethiopic", Status: "allowed"},  // Amharic, Tigrinya, etc.
		{Charset: "Mongolian", Status: "allowed"}, // Mongolian
		{Charset: "Tibetan", Status: "allowed"},   // Tibetan

		// Special categories
		{Charset: "Mixed", Status: "allowed"}, // Mixed scripts
		{Charset: "UTF-8", Status: "allowed"}, // UTF-8 encoded text
		{Charset: "Other", Status: "allowed"}, // Other unrecognized scripts
	}

	// Insert charset rules in batches
	batchSize := 10
	for i := 0; i < len(charsetRules); i += batchSize {
		end := i + batchSize
		if end > len(charsetRules) {
			end = len(charsetRules)
		}

		batch := charsetRules[i:end]
		if err := db.Create(&batch).Error; err != nil {
			return fmt.Errorf("failed to insert charset rules batch %d: %w", i/batchSize+1, err)
		}
	}

	log.Printf("Successfully seeded %d charset rules", len(charsetRules))
	return nil
}
