package model

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User represents a system user
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Email     string    `json:"email"`
	IsAdmin   bool      `gorm:"default:false" json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Settings represents system settings
type Settings struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Key       string    `gorm:"unique;not null" json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// InitDB initializes the database
func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("serverpanel.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto migrate tables
	if err := DB.AutoMigrate(&User{}, &Settings{}); err != nil {
		return err
	}

	// Create default admin user if not exists
	var count int64
	DB.Model(&User{}).Count(&count)
	if count == 0 {
		// Default password: admin123 (you should hash this properly)
		admin := User{
			Username: "admin",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // bcrypt hash of "admin123"
			IsAdmin:  true,
		}
		DB.Create(&admin)
	}

	return nil
}
