package auth

import (
	"github.com/gin-gonic/gin"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// RegisterRoutes 注册所有认证相关路由
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("")

	// 注册认证路由（登录、登出、修改密码、头像）
	RegisterAuthRoutes(r, db)

	// 注册微信登录路由
	RegisterWechatRoutes(r, db)

	// 需要认证的路由组
	auth := r.Group("/auth")
	auth.Use(jwtpkg.TokenMiddleware())

	// 注册用户管理路由
	RegisterUserRoutes(auth, db)

	// 注册角色管理路由
	RegisterRoleRoutes(auth, db)
}
