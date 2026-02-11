package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes 注册应用版本相关路由
func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	handler := NewHandler(db)

	// 应用版本路由组 - 不需要认证
	app := r.Group("/app")
	{
		app.GET("/version", handler.CheckVersion)
	}
}
