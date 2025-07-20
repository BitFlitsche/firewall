package models

import "time"

// IP represents the structure for the IP addresses table
type IP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Address   string    `gorm:"unique;not null;type:varchar(45)" json:"address" binding:"required"`                          // IPv6 max length or CIDR notation
	Status    string    `gorm:"not null;type:varchar(20)" json:"status" binding:"required,oneof=allowed denied whitelisted"` // "denied", "allowed", "whitelisted"
	IsCIDR    bool      `gorm:"column:is_c_id_r;default:false;type:boolean" json:"is_cidr"`                                  // Correct column for CIDR flag
	Source    string    `gorm:"type:varchar(50)" json:"source"`                                                              // Source of the IP data (e.g., "stopforumspam_toxic_cidr", "manual")
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Email represents the structure for the Emails table
type Email struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Address   string    `gorm:"unique;not null;type:varchar(254)" json:"address" binding:"required,email"`                   // RFC 5321 max length
	Status    string    `gorm:"not null;type:varchar(20)" json:"status" binding:"required,oneof=allowed denied whitelisted"` // "denied", "allowed", "whitelisted"
	IsRegex   bool      `gorm:"default:false;type:boolean" json:"is_regex"`                                                  // Whether this is a regex pattern
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// UserAgent represents the structure for the User Agents table
type UserAgent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserAgent string    `gorm:"unique;not null;type:varchar(500)" json:"user_agent" binding:"required,max=500"`              // Reasonable max length
	Status    string    `gorm:"not null;type:varchar(20)" json:"status" binding:"required,oneof=allowed denied whitelisted"` // "denied", "allowed", "whitelisted"
	IsRegex   bool      `gorm:"default:false;type:boolean" json:"is_regex"`                                                  // Whether this is a regex pattern
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Country represents the structure for the Countries table
type Country struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"unique;not null;type:varchar(2)" json:"code" binding:"required,len=2,alpha"`                  // ISO 3166-1 alpha-2 code
	Name      string    `gorm:"not null;type:varchar(100)" json:"name" binding:"required,max=100"`                           // Country name
	Status    string    `gorm:"not null;type:varchar(20)" json:"status" binding:"required,oneof=allowed denied whitelisted"` // "denied", "allowed", "whitelisted"
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CharsetRule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Charset   string    `gorm:"unique;not null;type:varchar(100)" json:"charset" binding:"required,max=100,alphanum"`
	Status    string    `gorm:"not null;type:varchar(20)" json:"status" binding:"required,oneof=allowed denied whitelisted"` // denied, allowed, whitelisted
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UsernameRule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null;type:varchar(100)" json:"username" binding:"required,max=100"`
	Status    string    `gorm:"not null;type:varchar(20)" json:"status" binding:"required,oneof=allowed denied whitelisted"` // denied, allowed, whitelisted
	IsRegex   bool      `gorm:"default:false;type:boolean" json:"is_regex"`                                                  // Whether this is a regex pattern
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// SyncTracker tracks the last sync timestamp for each data type
type SyncTracker struct {
	ID        uint      `gorm:"primaryKey"`
	DataType  string    `gorm:"unique;not null;type:varchar(50)"` // "ips", "emails", "user_agents", "countries", "charsets", "usernames", "asns"
	LastSync  time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// ASN represents the structure for the ASNs table
type ASN struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ASN       string    `gorm:"unique;not null;type:varchar(20)" json:"asn" binding:"required,max=20"`                       // ASN number (e.g., "AS12345")
	RIR       string    `gorm:"type:varchar(20)" json:"rir" binding:"omitempty,max=20"`                                      // Regional Internet Registry (e.g., "arin", "ripencc") - optional
	Domain    string    `gorm:"type:varchar(255)" json:"domain" binding:"omitempty,max=255"`                                 // Domain name - optional
	Country   string    `gorm:"type:varchar(2)" json:"cc" binding:"omitempty,len=2,alpha"`                                   // Country code (ISO 3166-1 alpha-2) - optional
	Name      string    `gorm:"not null;type:varchar(255)" json:"asname" binding:"required,max=255"`                         // ASN name/description
	Status    string    `gorm:"not null;type:varchar(20)" json:"status" binding:"required,oneof=allowed denied whitelisted"` // "denied", "allowed", "whitelisted"
	Source    string    `gorm:"type:varchar(50)" json:"source"`                                                              // Source of the ASN data (e.g., "spamhaus", "manual")
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
