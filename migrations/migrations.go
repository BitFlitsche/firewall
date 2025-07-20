package migrations

import (
	"firewall/models" // Import your models package
	"fmt"

	"gorm.io/gorm"
)

// Migrate runs the migrations for all tables
func Migrate(db *gorm.DB) error {
	// Auto migrate the models
	err := db.AutoMigrate(
		&models.IP{},
		&models.Email{},
		&models.UserAgent{},
		&models.Country{},
		&models.CharsetRule{},
		&models.UsernameRule{},
		&models.ASN{},
		&models.SyncTracker{},
		&models.TrafficLog{},
		&models.DataRelationship{},
		&models.AnalyticsAggregation{},
	)
	if err != nil {
		return err
	}

	// Run manual migrations for schema changes
	if err := runManualMigrations(db); err != nil {
		return err
	}
	if err != nil {
		return err
	}

	// Single column indexes for basic filtering
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ip_status ON i_ps (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ip_address ON i_ps (address)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_email_status ON emails (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_email_address ON emails (address)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_useragent_status ON user_agents (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_useragent_useragent ON user_agents (user_agent(100))")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_country_status ON countries (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_country_code ON countries (code)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_charset_status ON charset_rules (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_charset_charset ON charset_rules (charset)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_username_status ON username_rules (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_username_username ON username_rules (username)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_asn_status ON asns (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_asn_asn ON asns (asn)")

	// Composite indexes for common filter combinations (status + search field)
	// Using limited key lengths to prevent MySQL key length errors
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ip_status_address ON i_ps (status, address(45))")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ip_address_status ON i_ps (address(45), status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_email_status_address ON emails (status, address(100))")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_email_address_status ON emails (address(100), status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_useragent_status_useragent ON user_agents (status, user_agent(100))")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_useragent_useragent_status ON user_agents (user_agent(100), status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_country_status_code ON countries (status, code)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_country_code_status ON countries (code, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_charset_status_charset ON charset_rules (status, charset)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_charset_charset_status ON charset_rules (charset, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_username_status_username ON username_rules (status, username)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_username_username_status ON username_rules (username, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_asn_status_asn ON asns (status, asn)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_asn_asn_status ON asns (asn, status)")

	// Composite indexes for sorting with filtering
	// Using limited key lengths to prevent MySQL key length errors
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ip_status_id ON i_ps (status, id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ip_status_address_id ON i_ps (status, address(45), id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_email_status_id ON emails (status, id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_email_status_address_id ON emails (status, address(100), id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_useragent_status_id ON user_agents (status, id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_useragent_status_useragent_id ON user_agents (status, user_agent(100), id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_country_status_id ON countries (status, id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_country_status_code_id ON countries (status, code, id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_charset_status_id ON charset_rules (status, id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_charset_status_charset_id ON charset_rules (status, charset, id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_username_status_id ON username_rules (status, id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_username_status_username_id ON username_rules (status, username, id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_asn_status_id ON asns (status, id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_asn_status_asn_id ON asns (status, asn, id)")

	// Traffic logging indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_traffic_logs_timestamp ON traffic_logs (timestamp)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_traffic_logs_ip_address ON traffic_logs (ip_address)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_traffic_logs_email ON traffic_logs (email)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_traffic_logs_final_result ON traffic_logs (final_result)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_traffic_logs_request_id ON traffic_logs (request_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_traffic_logs_timestamp_final_result ON traffic_logs (timestamp, final_result)")

	// Data relationships indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_data_relationships_relationship_type ON data_relationships (relationship_type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_data_relationships_ip_address ON data_relationships (ip_address)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_data_relationships_email ON data_relationships (email)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_data_relationships_timestamp ON data_relationships (timestamp)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_data_relationships_frequency ON data_relationships (frequency)")

	// Analytics aggregations indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_analytics_aggregations_date_type ON analytics_aggregations (aggregation_date, aggregation_type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_analytics_aggregations_type ON analytics_aggregations (aggregation_type)")

	return nil
}

// runManualMigrations runs manual SQL migrations
func runManualMigrations(db *gorm.DB) error {
	// Migration to make ASN fields optional
	migrations := []string{
		"ALTER TABLE asns MODIFY COLUMN rir VARCHAR(20) NULL",
		"ALTER TABLE asns MODIFY COLUMN domain VARCHAR(255) NULL",
		"ALTER TABLE asns MODIFY COLUMN country VARCHAR(2) NULL",
		"ALTER TABLE asns COMMENT = 'ASN table with optional RIR, Domain, and Country fields'",
	}

	// Execute each migration separately
	for _, migration := range migrations {
		if err := db.Exec(migration).Error; err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
