package main

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null;uniqueIndex"`
	Permissions string `gorm:"type:text"`
}

func main() {
	dsn := "host=127.0.0.1 port=5432 user=materials password=julei1984 dbname=materials sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Add system permissions to roles that need system settings access
	rolesToUpdate := []string{
		"项目经理",   // Project managers
		"材料员",     // Material staff
		"分包材料员", // Subcontractor material staff
		"施工员",     // Foreman - to fix the 403 errors
		"保管员",     // Keeper
	}

	permissions := []string{"system_config"}

	for _, roleName := range rolesToUpdate {
		fmt.Printf("\n=== Processing role: %s ===\n", roleName)
		var roles []Role
		if err := db.Where("name = ?", roleName).Find(&roles).Error; err != nil {
			log.Printf("  Error finding roles: %v\n", err)
			continue
		}

		for _, role := range roles {
			fmt.Printf("  Found role: %s (ID: %d)\n", role.Name, role.ID)
			updated := false
			rolePerms := role.Permissions

			for _, perm := range permissions {
				if !containsPermission(rolePerms, perm) {
					if rolePerms == "" {
						rolePerms = perm
					} else {
						rolePerms += "," + perm
					}
					updated = true
					fmt.Printf("  Adding permission: %s\n", perm)
				}
			}

			if updated {
				role.Permissions = rolePerms
				if err := db.Save(&role).Error; err != nil {
					log.Printf("  Error updating: %v\n", err)
				} else {
					fmt.Printf("  Updated!\n")
				}
			} else {
				fmt.Printf("  No changes needed\n")
			}
		}
	}

	fmt.Println("\n=== Done ===")
}

func containsPermission(permissions, permission string) bool {
	if permissions == "" {
		return false
	}
	permList := strings.Split(permissions, ",")
	for _, p := range permList {
		if strings.TrimSpace(p) == permission {
			return true
		}
	}
	return false
}
