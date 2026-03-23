package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// roleMapping maps English role values to Chinese role names in the database
var roleMapping = map[string]string{
	"foreman":                 "施工员",
	"keeper":                  "保管员",
	"project_manager":         "项目经理",
	"worker":                  "作业人员",
	"material_staff":          "材料员",
	"subcontractor_material_staff": "分包材料员",
	"appointment_admin":       "预约管理员",
}

func main() {
	// Database connection
	dsn := "host=127.0.0.1 port=5432 user=materials password=julei1984 dbname=materials sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// List all users
	fmt.Println("=== All users in database ===")
	var users []auth.User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("Error finding users: %v", err)
	} else {
		for _, u := range users {
			fmt.Printf("  User: %s, Role: %s, Roles count: %d\n", u.Username, u.Role, len(u.Roles))
		}
	}

	fmt.Println()
	fmt.Println("=== All roles with role_view permission ===")
	var roles []auth.Role
	if err := db.Find(&roles).Error; err != nil {
		log.Printf("Error finding roles: %v", err)
	} else {
		for _, r := range roles {
			if strings.Contains(r.Permissions, "role_view") {
				fmt.Printf("  Role: %s (ID: %d)\n", r.Name, r.ID)
				fmt.Printf("    Permissions include: role_view\n")
			}
		}
	}

	fmt.Println()
	fmt.Println("=== All roles in database ===")
	if err := db.Find(&roles).Error; err != nil {
		log.Printf("Error finding roles: %v", err)
	} else {
		for _, r := range roles {
			fmt.Printf("  Role: %s (ID: %d), Permissions count: %d\n", r.Name, r.ID, strings.Count(r.Permissions, ",")+1)
		}
	}

	fmt.Println()

	// Test: Get user "julei" with role="foreman"
	fmt.Println("=== Testing julei user (role=foreman) ===")
	var user auth.User
	if err := db.Where("username = ?", "julei").First(&user).Error; err != nil {
		log.Printf("Error finding julei: %v", err)
	} else {
		fmt.Printf("User: %s (ID: %d)\n", user.Username, user.ID)
		fmt.Printf("  Role field: %s\n", user.Role)
		fmt.Printf("  Roles array: %d items\n", len(user.Roles))

		// Reload with the middleware logic
		if err := db.Preload("Roles").First(&user, user.ID).Error; err != nil {
			log.Printf("Error reloading: %v", err)
		}
		fmt.Printf("  After Preload(Roles): %d items\n", len(user.Roles))
		if len(user.Roles) > 0 {
			fmt.Printf("    Loaded role: %s\n", user.Roles[0].Name)
		}

		// Check permissions before fix
		fmt.Printf("  Before fix - HasPermission(project_view): %v\n", user.HasPermission("project_view"))

		// Simulate the middleware fix with role mapping
		if len(user.Roles) == 0 && user.Role != "" {
			var role auth.Role
			roleName := user.Role

			// Check if there's a mapping for this role
			if mappedName, ok := roleMapping[user.Role]; ok {
				roleName = mappedName
				fmt.Printf("  Mapped role '%s' to '%s'\n", user.Role, roleName)
			}

			// Try with exact match first, then case-insensitive
			if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
				if err := db.Where("LOWER(name) = ?", strings.ToLower(roleName)).First(&role).Error; err != nil {
					fmt.Printf("  Failed to load role: %v\n", err)
				}
			}

			// Only set Roles if we found a match
			if role.ID != 0 {
				user.Roles = []auth.Role{role}
				fmt.Printf("  After middleware fix: Loaded role '%s' (ID: %d)\n", role.Name, role.ID)
			}
		}

		fmt.Printf("  IsAdmin: %v\n", user.IsAdmin())
		fmt.Printf("  HasPermission(project_view): %v\n", user.HasPermission("project_view"))
		fmt.Printf("  HasPermission(project): %v\n", user.HasPermission("project"))
	}

	fmt.Println()
	fmt.Println("=== Testing wqs user (role=keeper) ===")
	var keeper auth.User
	if err := db.Where("username = ?", "wqs").First(&keeper).Error; err != nil {
		log.Printf("Error finding wqs: %v", err)
	} else {
		fmt.Printf("User: %s (ID: %d)\n", keeper.Username, keeper.ID)
		fmt.Printf("  Role field: %s\n", keeper.Role)

		// Simulate the middleware fix with role mapping
		if len(keeper.Roles) == 0 && keeper.Role != "" {
			var role auth.Role
			roleName := keeper.Role

			// Check if there's a mapping for this role
			if mappedName, ok := roleMapping[keeper.Role]; ok {
				roleName = mappedName
				fmt.Printf("  Mapped role '%s' to '%s'\n", keeper.Role, roleName)
			}

			if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
				if err := db.Where("LOWER(name) = ?", strings.ToLower(roleName)).First(&role).Error; err != nil {
					fmt.Printf("  Failed to load role: %v\n", err)
				}
			}

			if role.ID != 0 {
				keeper.Roles = []auth.Role{role}
				fmt.Printf("  After middleware fix: Loaded role '%s' (ID: %d)\n", role.Name, role.ID)
			}
		}

		fmt.Printf("  IsAdmin: %v\n", keeper.IsAdmin())
		fmt.Printf("  HasPermission(project_view): %v\n", keeper.HasPermission("project_view"))
		fmt.Printf("  HasPermission(project): %v\n", keeper.HasPermission("project"))
	}

	fmt.Println()
	fmt.Println("=== Testing admin user ===")
	var admin auth.User
	if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
		log.Printf("Error finding admin: %v", err)
	} else {
		fmt.Printf("User: %s (ID: %d)\n", admin.Username, admin.ID)
		fmt.Printf("  Role field: %s\n", admin.Role)
		fmt.Printf("  Roles array: %d items\n", len(admin.Roles))

		// Reload with the middleware logic
		if err := db.Preload("Roles").First(&admin, admin.ID).Error; err != nil {
			log.Printf("Error reloading: %v", err)
		}
		if len(admin.Roles) > 0 {
			fmt.Printf("  Loaded role: %s (ID: %d)\n", admin.Roles[0].Name, admin.Roles[0].ID)
			fmt.Printf("  Role permissions: %s\n", admin.Roles[0].Permissions)
		}

		fmt.Printf("  IsAdmin: %v\n", admin.IsAdmin())
		fmt.Printf("  HasPermission(role_view): %v\n", admin.HasPermission("role_view"))
		fmt.Printf("  HasPermission(user_view): %v\n", admin.HasPermission("user_view"))
	}
}
