package services

import (
	"fmt"
	"log"
	"time"

	"firewall/models"

	"gorm.io/gorm"
)

// CountryData represents the seed data for countries
type CountryData struct {
	Code   string
	Name   string
	Status string
}

// SeedCountries populates the countries table if it's empty
func SeedCountries(db *gorm.DB) error {
	// Check if countries table is empty
	var count int64
	if err := db.Model(&models.Country{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check countries count: %w", err)
	}

	if count > 0 {
		log.Println("Countries table already populated, skipping seed")
		return nil
	}

	log.Println("Countries table is empty, seeding with MaxMind country data...")

	// Country data from MaxMind GeoLite2 database
	countries := []CountryData{
		{"AC", "Ascension Island", "allowed"},
		{"AD", "Andorra", "allowed"},
		{"AE", "United Arab Emirates", "allowed"},
		{"AF", "Afghanistan", "allowed"},
		{"AI", "Anguilla", "allowed"},
		{"AL", "Albania", "allowed"},
		{"AM", "Armenia", "allowed"},
		{"AQ", "Antarctica", "allowed"},
		{"AR", "Argentina", "allowed"},
		{"AS", "American Samoa", "allowed"},
		{"AT", "Austria", "allowed"},
		{"AU", "Australia", "allowed"},
		{"AZ", "Azerbaijan", "allowed"},
		{"BA", "Bosnia and Herzegovina", "allowed"},
		{"BD", "Bangladesh", "allowed"},
		{"BE", "Belgium", "allowed"},
		{"BG", "Bulgaria", "allowed"},
		{"BL", "Saint Barthélemy", "allowed"},
		{"BN", "Brunei", "allowed"},
		{"BO", "Bolivia", "allowed"},
		{"BR", "Brazil", "allowed"},
		{"BT", "Bhutan", "allowed"},
		{"BV", "Bouvet Island", "allowed"},
		{"BY", "Belarus", "allowed"},
		{"CA", "Canada", "allowed"},
		{"CC", "Cocos (Keeling) Islands", "allowed"},
		{"CH", "Switzerland", "allowed"},
		{"CK", "Cook Islands", "allowed"},
		{"CL", "Chile", "allowed"},
		{"CN", "China", "allowed"},
		{"CO", "Colombia", "allowed"},
		{"CR", "Costa Rica", "allowed"},
		{"CU", "Cuba", "allowed"},
		{"CX", "Christmas Island", "allowed"},
		{"CY", "Cyprus", "allowed"},
		{"CZ", "Czech Republic", "allowed"},
		{"DE", "Germany", "allowed"},
		{"DK", "Denmark", "allowed"},
		{"DO", "Dominican Republic", "allowed"},
		{"EC", "Ecuador", "allowed"},
		{"EE", "Estonia", "allowed"},
		{"EG", "Egypt", "allowed"},
		{"ES", "Spain", "allowed"},
		{"FI", "Finland", "allowed"},
		{"FJ", "Fiji", "allowed"},
		{"FK", "Falkland Islands", "allowed"},
		{"FM", "Micronesia", "allowed"},
		{"FR", "France", "allowed"},
		{"GB", "United Kingdom", "allowed"},
		{"GE", "Georgia", "allowed"},
		{"GF", "French Guiana", "allowed"},
		{"GP", "Guadeloupe", "allowed"},
		{"GR", "Greece", "allowed"},
		{"GS", "South Georgia and the South Sandwich Islands", "allowed"},
		{"GT", "Guatemala", "allowed"},
		{"GU", "Guam", "allowed"},
		{"HK", "Hong Kong", "allowed"},
		{"HM", "Heard Island and McDonald Islands", "allowed"},
		{"HN", "Honduras", "allowed"},
		{"HR", "Croatia", "allowed"},
		{"HT", "Haiti", "allowed"},
		{"HU", "Hungary", "allowed"},
		{"ID", "Indonesia", "allowed"},
		{"IE", "Ireland", "allowed"},
		{"IL", "Israel", "allowed"},
		{"IN", "India", "allowed"},
		{"IO", "British Indian Ocean Territory", "allowed"},
		{"IS", "Iceland", "allowed"},
		{"IT", "Italy", "allowed"},
		{"JM", "Jamaica", "allowed"},
		{"JP", "Japan", "allowed"},
		{"KE", "Kenya", "allowed"},
		{"KG", "Kyrgyzstan", "allowed"},
		{"KH", "Cambodia", "allowed"},
		{"KI", "Kiribati", "allowed"},
		{"KR", "South Korea", "allowed"},
		{"KZ", "Kazakhstan", "allowed"},
		{"LA", "Laos", "allowed"},
		{"LI", "Liechtenstein", "allowed"},
		{"LK", "Sri Lanka", "allowed"},
		{"LT", "Lithuania", "allowed"},
		{"LU", "Luxembourg", "allowed"},
		{"LV", "Latvia", "allowed"},
		{"MA", "Morocco", "allowed"},
		{"MC", "Monaco", "allowed"},
		{"MD", "Moldova", "allowed"},
		{"ME", "Montenegro", "allowed"},
		{"MF", "Saint Martin", "allowed"},
		{"MH", "Marshall Islands", "allowed"},
		{"MK", "North Macedonia", "allowed"},
		{"MM", "Myanmar", "allowed"},
		{"MP", "Northern Mariana Islands", "allowed"},
		{"MQ", "Martinique", "allowed"},
		{"MS", "Montserrat", "allowed"},
		{"MT", "Malta", "allowed"},
		{"MV", "Maldives", "allowed"},
		{"MX", "Mexico", "allowed"},
		{"MY", "Malaysia", "allowed"},
		{"NC", "New Caledonia", "allowed"},
		{"NF", "Norfolk Island", "allowed"},
		{"NG", "Nigeria", "allowed"},
		{"NI", "Nicaragua", "allowed"},
		{"NL", "Netherlands", "allowed"},
		{"NO", "Norway", "allowed"},
		{"NP", "Nepal", "allowed"},
		{"NR", "Nauru", "allowed"},
		{"NU", "Niue", "allowed"},
		{"NZ", "New Zealand", "allowed"},
		{"PA", "Panama", "allowed"},
		{"PE", "Peru", "allowed"},
		{"PF", "French Polynesia", "allowed"},
		{"PG", "Papua New Guinea", "allowed"},
		{"PH", "Philippines", "allowed"},
		{"PK", "Pakistan", "allowed"},
		{"PL", "Poland", "allowed"},
		{"PM", "Saint Pierre and Miquelon", "allowed"},
		{"PN", "Pitcairn Islands", "allowed"},
		{"PR", "Puerto Rico", "allowed"},
		{"PT", "Portugal", "allowed"},
		{"PW", "Palau", "allowed"},
		{"PY", "Paraguay", "allowed"},
		{"RE", "Réunion", "allowed"},
		{"RO", "Romania", "allowed"},
		{"RS", "Serbia", "allowed"},
		{"RU", "Russia", "allowed"},
		{"SA", "Saudi Arabia", "allowed"},
		{"SB", "Solomon Islands", "allowed"},
		{"SE", "Sweden", "allowed"},
		{"SG", "Singapore", "allowed"},
		{"SH", "Saint Helena", "allowed"},
		{"SI", "Slovenia", "allowed"},
		{"SK", "Slovakia", "allowed"},
		{"SM", "San Marino", "allowed"},
		{"SV", "El Salvador", "allowed"},
		{"TA", "Tristan da Cunha", "allowed"},
		{"TC", "Turks and Caicos Islands", "allowed"},
		{"TF", "French Southern Territories", "allowed"},
		{"TH", "Thailand", "allowed"},
		{"TJ", "Tajikistan", "allowed"},
		{"TK", "Tokelau", "allowed"},
		{"TL", "Timor-Leste", "allowed"},
		{"TM", "Turkmenistan", "allowed"},
		{"TN", "Tunisia", "allowed"},
		{"TO", "Tonga", "allowed"},
		{"TR", "Turkey", "allowed"},
		{"TV", "Tuvalu", "allowed"},
		{"TW", "Taiwan", "allowed"},
		{"UA", "Ukraine", "allowed"},
		{"US", "United States", "allowed"},
		{"UY", "Uruguay", "allowed"},
		{"UZ", "Uzbekistan", "allowed"},
		{"VA", "Vatican City", "allowed"},
		{"VE", "Venezuela", "allowed"},
		{"VG", "British Virgin Islands", "allowed"},
		{"VI", "U.S. Virgin Islands", "allowed"},
		{"VN", "Vietnam", "allowed"},
		{"VU", "Vanuatu", "allowed"},
		{"WF", "Wallis and Futuna", "allowed"},
		{"WS", "Samoa", "allowed"},
		{"XK", "Kosovo", "allowed"},
		{"YT", "Mayotte", "allowed"},
		{"ZA", "South Africa", "allowed"},
	}

	// Insert countries in batches for better performance
	batchSize := 50
	for i := 0; i < len(countries); i += batchSize {
		end := i + batchSize
		if end > len(countries) {
			end = len(countries)
		}

		batch := countries[i:end]
		var countryModels []models.Country

		for _, country := range batch {
			countryModels = append(countryModels, models.Country{
				Code:      country.Code,
				Name:      country.Name,
				Status:    country.Status,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}

		if err := db.Create(&countryModels).Error; err != nil {
			return fmt.Errorf("failed to insert country batch: %w", err)
		}

		log.Printf("Inserted %d countries (batch %d/%d)", len(batch), (i/batchSize)+1, (len(countries)+batchSize-1)/batchSize)
	}

	log.Printf("Successfully seeded %d countries", len(countries))
	return nil
}
