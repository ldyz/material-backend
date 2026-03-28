package agent

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	appconfig "github.com/yourorg/material-backend/backend/internal/config"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// RegisterRoutes registers agent API routes
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *appconfig.Config, aiHandler *AIHandler) {
	r := rg.Group("/agent")
	r.Use(jwtpkg.TokenMiddleware())

	h := NewHandler(db)

	// Use shared AI handler if provided, otherwise create one from config
	if aiHandler != nil {
		h.service.SetAIHandler(aiHandler)
	} else {
		asrURL := ""
		if cfg.AI.ASREnabled {
			asrURL = cfg.AI.ASRServiceURL
		}

		// 优先使用百度千帆配置 (Anthropic 兼容 API，使用 OpenAI 客户端)
		if cfg.AI.BaiduAPIKey != "" && cfg.AI.BaiduBaseURL != "" {
			newAIHandler := NewAIHandler(
				db,
				cfg.AI.BaiduAPIKey,      // Auth Token 作为 API Key
				cfg.AI.BaiduModel,       // 模型名称
				cfg.AI.BaiduBaseURL,     // 百度千帆 Base URL
				cfg.AI.OpenAIAPIKey,
				asrURL,
			)
			h.service.SetAIHandler(newAIHandler)
		} else if cfg.AI.DeepSeekAPIKey != "" {
			// 回退到 DeepSeek 配置
			newAIHandler := NewAIHandler(
				db,
				cfg.AI.DeepSeekAPIKey,
				cfg.AI.DeepSeekModel,
				cfg.AI.DeepSeekBaseURL,
				cfg.AI.OpenAIAPIKey,
				asrURL,
			)
			h.service.SetAIHandler(newAIHandler)
		}
	}

	// Core operation endpoints - each requires specific permission
	r.POST("/operate", auth.PermissionMiddleware(db, "ai_agent_operate"), h.HandleOperation)
	r.POST("/query", auth.PermissionMiddleware(db, "ai_agent_query"), h.HandleQuery)
	r.POST("/workflow", auth.PermissionMiddleware(db, "ai_agent_workflow"), h.HandleWorkflow)

	// Chat endpoints for mobile voice/text interaction
	r.POST("/chat", auth.PermissionMiddleware(db, "ai_agent_query"), h.HandleChat)
	r.POST("/voice-chat", auth.PermissionMiddleware(db, "ai_agent_query"), h.HandleVoiceChat)

	// Capability and validation endpoints - view permission is sufficient
	r.GET("/capabilities", auth.PermissionMiddleware(db, "ai_agent_view"), h.HandleCapabilities)
	r.POST("/validate", auth.PermissionMiddleware(db, "ai_agent_query"), h.HandleValidate)

	// Logging endpoints
	r.GET("/logs", auth.PermissionMiddleware(db, "ai_agent_logs"), h.HandleLogs)

	// Model provider endpoints
	r.GET("/providers", auth.PermissionMiddleware(db, "ai_agent_view"), h.HandleGetProviders)
	r.POST("/providers/switch", auth.PermissionMiddleware(db, "ai_agent_operate"), h.HandleSwitchProvider)

	// Conversation history endpoints
	r.GET("/conversation-history", auth.PermissionMiddleware(db, "ai_agent_view"), h.HandleGetConversationHistory)
	r.DELETE("/conversation-history", auth.PermissionMiddleware(db, "ai_agent_operate"), h.HandleClearConversationHistory)
}
