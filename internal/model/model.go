package model

import (
	"fmt"
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
	RoleID    *uint     `json:"roleId"`
	Role      *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Role represents a user role
type Role struct {
	ID          uint         `gorm:"primarykey" json:"id"`
	Name        string       `gorm:"unique;not null" json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

// Permission represents a permission
type Permission struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Resource    string    `gorm:"not null" json:"resource"` // docker, file, database, etc.
	Action      string    `gorm:"not null" json:"action"`   // read, write, delete, execute
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Auto migrate tables
	if err := DB.AutoMigrate(&User{}, &Role{}, &Permission{}, &Settings{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Create default permissions
	if err := createDefaultPermissions(); err != nil {
		return err
	}

	// Create default roles
	if err := createDefaultRoles(); err != nil {
		return err
	}

	// Create default admin user if not exists
	var count int64
	if err := DB.Model(&User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}
	if count == 0 {
		// Default password: admin123 (bcrypt hash)
		admin := User{
			Username: "admin",
			Password: "$2a$10$2URxF2u7riYpkieZii9to.rPNlWKXNmBsKXkxdKzIBA3rPJ9yKoB2", // bcrypt hash of "admin123"
			IsAdmin:  true,
			Enabled:  true,
		}
		if err := DB.Create(&admin).Error; err != nil {
			return fmt.Errorf("failed to create default admin user: %w", err)
		}
	}

	return nil
}

func createDefaultPermissions() error {
	permissions := []Permission{
		// System permissions
		{Name: "system.read", Resource: "system", Action: "read", Description: "View system information"},
		{Name: "system.write", Resource: "system", Action: "write", Description: "Modify system settings"},
		
		// Docker permissions
		{Name: "docker.read", Resource: "docker", Action: "read", Description: "View Docker containers and images"},
		{Name: "docker.write", Resource: "docker", Action: "write", Description: "Manage Docker containers and images"},
		
		// File permissions
		{Name: "file.read", Resource: "file", Action: "read", Description: "View files"},
		{Name: "file.write", Resource: "file", Action: "write", Description: "Modify files"},
		{Name: "file.delete", Resource: "file", Action: "delete", Description: "Delete files"},
		
		// Database permissions
		{Name: "database.read", Resource: "database", Action: "read", Description: "View databases"},
		{Name: "database.write", Resource: "database", Action: "write", Description: "Modify databases"},
		
		// Log permissions
		{Name: "log.read", Resource: "log", Action: "read", Description: "View logs"},
		{Name: "log.write", Resource: "log", Action: "write", Description: "Clear logs"},
		
		// Nginx permissions
		{Name: "nginx.read", Resource: "nginx", Action: "read", Description: "View Nginx configuration"},
		{Name: "nginx.write", Resource: "nginx", Action: "write", Description: "Modify Nginx configuration"},
		
		// Cron permissions
		{Name: "cron.read", Resource: "cron", Action: "read", Description: "View cron jobs"},
		{Name: "cron.write", Resource: "cron", Action: "write", Description: "Modify cron jobs"},
		
		// Backup permissions
		{Name: "backup.read", Resource: "backup", Action: "read", Description: "View backups"},
		{Name: "backup.write", Resource: "backup", Action: "write", Description: "Create and restore backups"},
		
		// User permissions
		{Name: "user.read", Resource: "user", Action: "read", Description: "View users"},
		{Name: "user.write", Resource: "user", Action: "write", Description: "Manage users"},
	}

	for _, perm := range permissions {
		var existing Permission
		if err := DB.Where("name = ?", perm.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := DB.Create(&perm).Error; err != nil {
				return fmt.Errorf("failed to create permission %s: %w", perm.Name, err)
			}
		}
	}

	return nil
}

func createDefaultRoles() error {
	// Admin role with all permissions
	var adminRole Role
	if err := DB.Where("name = ?", "admin").First(&adminRole).Error; err == gorm.ErrRecordNotFound {
		var allPermissions []Permission
		DB.Find(&allPermissions)

		adminRole = Role{
			Name:        "admin",
			Description: "Administrator with full access",
			Permissions: allPermissions,
		}
		if err := DB.Create(&adminRole).Error; err != nil {
			return fmt.Errorf("failed to create admin role: %w", err)
		}
	}

	// Viewer role with read-only permissions
	var viewerRole Role
	if err := DB.Where("name = ?", "viewer").First(&viewerRole).Error; err == gorm.ErrRecordNotFound {
		var readPermissions []Permission
		DB.Where("action = ?", "read").Find(&readPermissions)

		viewerRole = Role{
			Name:        "viewer",
			Description: "Read-only access",
			Permissions: readPermissions,
		}
		if err := DB.Create(&viewerRole).Error; err != nil {
			return fmt.Errorf("failed to create viewer role: %w", err)
		}
	}

	return nil
}
