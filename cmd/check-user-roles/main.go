package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID     uint
	Username string
	Role   string
}

type UserRole struct {
	UserID  uint
	RoleID  uint
	Role    Role `gorm:"foreignKey:RoleID"`
}

type Role struct {
	ID          uint
	Name        string
	Permissions string
}

func main() {
	dsn := "host=127.0.0.1 port=5432 user=materials password=julei1984 dbname=materials sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Checking user_roles for user lfh (ID: 11) ===")
	var userRoles []UserRole
	if err := db.Where("user_id = ?", 11).Preload("Role").Find(&userRoles).Error; err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Found %d role associations:\n", len(userRoles))
		for _, ur := range userRoles {
			fmt.Printf("  RoleID: %d, RoleName: %s\n", ur.RoleID, ur.Role.Name)
			fmt.Printf("    Permissions: %s\n", ur.Role.Permissions)
		}
	}

	fmt.Println("\n=== All 施工员 roles ===")
	var roles []Role
	if err := db.Where("name = ?", "施工员").Find(&roles).Error; err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Found %d roles:\n", len(roles))
		for _, r := range roles {
			hasRoleView := r.Permissions != "" && contains(r.Permissions, "role_view")
			fmt.Printf("  ID: %d, Permissions: %s\n", r.ID, r.Permissions)
			fmt.Printf("    Has role_view: %v\n", hasRoleView)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
