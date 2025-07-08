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
	)
	if err != nil {
		return err
	}
	return nil
}
