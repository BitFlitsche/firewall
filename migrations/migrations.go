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
	// Indizes für Filter-/Sortierspalten
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ip_status ON i_ps (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ip_address ON i_ps (address)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_email_status ON emails (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_email_address ON emails (address)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_useragent_status ON user_agents (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_useragent_useragent ON user_agents (user_agent)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_country_status ON countries (status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_country_code ON countries (code)")
	return nil
}
