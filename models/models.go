package models

import "time"

// IP represents the structure for the IP addresses table
type IP struct {
	ID        uint      `gorm:"primaryKey"`
	Address   string    `gorm:"unique;not null;type:varchar(45)"` // IPv6 max length or CIDR notation
	Status    string    `gorm:"not null;type:varchar(20)"`        // "denied", "allowed", "whitelisted"
	IsCIDR    bool      `gorm:"default:false;type:boolean"`       // Whether this is a CIDR block
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Email represents the structure for the Emails table
type Email struct {
	ID        uint      `gorm:"primaryKey"`
	Address   string    `gorm:"unique;not null;type:varchar(254)"` // RFC 5321 max length
	Status    string    `gorm:"not null;type:varchar(20)"`         // "denied", "allowed", "whitelisted"
	IsRegex   bool      `gorm:"default:false;type:boolean"`        // Whether this is a regex pattern
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// UserAgent represents the structure for the User Agents table
type UserAgent struct {
	ID        uint      `gorm:"primaryKey"`
	UserAgent string    `gorm:"unique;not null;type:varchar(500)"` // Reasonable max length
	Status    string    `gorm:"not null;type:varchar(20)"`         // "denied", "allowed", "whitelisted"
	IsRegex   bool      `gorm:"default:false;type:boolean"`        // Whether this is a regex pattern
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Country represents the structure for the Countries table
type Country struct {
	ID        uint      `gorm:"primaryKey"`
	Code      string    `gorm:"unique;not null;type:varchar(2)"` // ISO 3166-1 alpha-2 code
	Status    string    `gorm:"not null;type:varchar(20)"`       // "denied", "allowed", "whitelisted"
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type CharsetRule struct {
	ID        uint      `gorm:"primaryKey" json:"ID"`
	Charset   string    `gorm:"unique;not null;type:varchar(100)" json:"Charset"`
	Status    string    `gorm:"not null;type:varchar(20)" json:"Status"` // denied, allowed, whitelisted
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type UsernameRule struct {
	ID        uint      `gorm:"primaryKey" json:"ID"`
	Username  string    `gorm:"unique;not null;type:varchar(100)" json:"Username"`
	Status    string    `gorm:"not null;type:varchar(20)" json:"Status"`   // denied, allowed, whitelisted
	IsRegex   bool      `gorm:"default:false;type:boolean" json:"IsRegex"` // Whether this is a regex pattern
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// SyncTracker tracks the last sync timestamp for each data type
type SyncTracker struct {
	ID        uint      `gorm:"primaryKey"`
	DataType  string    `gorm:"unique;not null;type:varchar(50)"` // "ips", "emails", "user_agents", "countries", "charsets", "usernames"
	LastSync  time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
