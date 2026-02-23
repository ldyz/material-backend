package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Role struct {
	ID          uint
	Name        string
	Permissions string
}

type UserRole struct {
	UserID uint
	RoleID uint
}

func main() {
	dsn := "host=127.0.0.1 port=5432 user=materials password=julei1984 dbname=materials sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Checking lfh (user ID: 11) roles ===")

	// Get user_roles for lfh
	var userRoles []UserRole
	if err := db.Table("user_roles").Where("user_id = ?", 11).Find(&userRoles).Error; err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Found %d role associations:\n", len(userRoles))
	for _, ur := range userRoles {
		var role Role
		if err := db.First(&role, ur.RoleID).Error; err != nil {
			log.Printf("  RoleID %d: Error loading - %v", ur.RoleID, err)
		} else {
			hasSystemConfig := contains(role.Permissions, "system_config")
			fmt.Printf("  RoleID: %d, Name: %s\n", ur.RoleID, role.Name)
			fmt.Printf("    Has system_config: %v\n", hasSystemConfig)
			if !hasSystemConfig {
				fmt.Printf("    ⚠️  Missing system_config permission!\n")
			}
		}
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
