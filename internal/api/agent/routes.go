package agent

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/middleware"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// RegisterRoutes registers agent API routes
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("/agent")
	r.Use(jwtpkg.TokenMiddleware())
	// Use agent permission middleware for all agent routes
	r.Use(middleware.AgentPermissionMiddleware())

	h := NewHandler(db)

	// Core operation endpoints
	r.POST("/operate", h.HandleOperation)
	r.POST("/query", h.HandleQuery)
	r.POST("/workflow", h.HandleWorkflow)

	// Capability and validation endpoints
	r.GET("/capabilities", h.HandleCapabilities)
	r.POST("/validate", h.HandleValidate)

	// Logging endpoints
	r.GET("/logs", auth.PermissionMiddleware(db, "ai_agent_logs"), h.HandleLogs)
}
