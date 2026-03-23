package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/db"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres dbname=material sslmode=disable"
	}

	database, err := db.New(dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// Check if admin exists
	var admin auth.User
	result := database.Where("username = ?", "admin").First(&admin)
	
	if result.Error != nil {
		// Create admin user
		admin = auth.User{
			Username: "admin",
			Email:    "admin@example.com",
			Role:     "admin",
			FullName: "系统管理员",
			IsActive: true,
		}
		admin.SetPassword("admin123")
		if err := database.Create(&admin).Error; err != nil {
			log.Fatalf("创建admin用户失败: %v", err)
		}
		fmt.Println("✅ 已创建admin用户，密码: admin123")
	} else {
		// Update password
		admin.SetPassword("admin123")
		if err := database.Save(&admin).Error; err != nil {
			log.Fatalf("更新密码失败: %v", err)
		}
		fmt.Println("✅ 已重置admin密码为: admin123")
	}
}
