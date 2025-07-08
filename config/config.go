// Package config config/config.go
package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// InitMySQL initializes and connects to the MySQL database
func InitMySQL() {
	dsn := "user:password@tcp(127.0.0.1:3306)/firewall?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	DB = db
	log.Println("MySQL connected successfully")
}
