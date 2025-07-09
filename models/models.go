package models

// IP represents the structure for the IP addresses table
type IP struct {
	ID      uint   `gorm:"primaryKey"`
	Address string `gorm:"unique;not null"`
	Status  string `gorm:"not null"` // "denied", "allowed", "whitelisted"
}

// Email represents the structure for the Emails table
type Email struct {
	ID      uint   `gorm:"primaryKey"`
	Address string `gorm:"unique;not null"`
	Status  string `gorm:"not null"` // "denied", "allowed", "whitelisted"
}

// UserAgent represents the structure for the User Agents table
type UserAgent struct {
	ID        uint   `gorm:"primaryKey"`
	UserAgent string `gorm:"unique;not null"`
	Status    string `gorm:"not null"` // "denied", "allowed", "whitelisted"
}

// Country represents the structure for the Countries table
type Country struct {
	ID     uint   `gorm:"primaryKey"`
	Code   string `gorm:"unique;not null"` // ISO 3166-1 alpha-2 code
	Status string `gorm:"not null"`        // "denied", "allowed", "whitelisted"
}

type CharsetRule struct {
	ID      uint   `gorm:"primaryKey" json:"ID"`
	Charset string `gorm:"unique;not null" json:"Charset"`
	Status  string `gorm:"not null" json:"Status"` // denied, allowed, whitelisted
}

type UsernameRule struct {
	ID       uint   `gorm:"primaryKey" json:"ID"`
	Username string `gorm:"unique;not null" json:"Username"`
	Status   string `gorm:"not null" json:"Status"` // denied, allowed, whitelisted
}
