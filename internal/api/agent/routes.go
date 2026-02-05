package agent

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// RegisterRoutes registers agent API routes
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("/agent")
	r.Use(jwtpkg.TokenMiddleware())

	h := NewHandler(db)

	// Core operation endpoints - each requires specific permission
	r.POST("/operate", auth.PermissionMiddleware(db, "ai_agent_operate"), h.HandleOperation)
	r.POST("/query", auth.PermissionMiddleware(db, "ai_agent_query"), h.HandleQuery)
	r.POST("/workflow", auth.PermissionMiddleware(db, "ai_agent_workflow"), h.HandleWorkflow)

	// Capability and validation endpoints - view permission is sufficient
	r.GET("/capabilities", auth.PermissionMiddleware(db, "ai_agent_view"), h.HandleCapabilities)
	r.POST("/validate", auth.PermissionMiddleware(db, "ai_agent_query"), h.HandleValidate)

	// Logging endpoints
	r.GET("/logs", auth.PermissionMiddleware(db, "ai_agent_logs"), h.HandleLogs)
}
