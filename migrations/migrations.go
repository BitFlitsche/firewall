package migrations

import (
	"firewall/models" // Import your models package

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
	)
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

	return nil
}
