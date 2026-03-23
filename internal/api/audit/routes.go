package audit

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// RegisterRoutes 注册审计日志路由
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("/audit")
	r.Use(jwtpkg.TokenMiddleware())

	h := NewHandler(db)

	// 操作日志查询（需要审计权限）
	r.GET("/operation-logs", auth.PermissionMiddleware(db, "audit_view"), h.GetLogs)
	r.GET("/operation-logs/statistics", auth.PermissionMiddleware(db, "audit_view"), h.GetLogStatistics)
	r.GET("/operation-logs/export", auth.PermissionMiddleware(db, "audit_view"), h.ExportLogs)
	r.GET("/operation-logs/:id", auth.PermissionMiddleware(db, "audit_view"), h.GetLogDetail)
	r.GET("/operation-logs/resource/:resource_type/:resource_id", auth.PermissionMiddleware(db, "audit_view"), h.GetResourceLogs)
	r.GET("/operation-logs/user/:user_id", auth.PermissionMiddleware(db, "audit_view"), h.GetUserLogs)

	// 清理旧日志（需要管理员权限）
	r.DELETE("/operation-logs/cleanup", auth.PermissionMiddleware(db, "admin"), h.DeleteOldLogs)
}
