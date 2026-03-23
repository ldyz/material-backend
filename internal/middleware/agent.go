package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
)

// AgentPermissionMiddleware validates AI Agent requests
func AgentPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user
		userID, exists := c.Get("current_user_id")
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		log.Printf("[AgentMiddleware] User ID: %v", userID)

		// Get permissions from context (set by TokenMiddleware)
		permissions, exists := c.Get("permissions")
		if !exists {
			response.Forbidden(c, "Unable to verify permissions")
			c.Abort()
			return
		}

		// Check if user has any AI agent permission
		permStr, ok := permissions.(string)
		if !ok || permStr == "" {
			response.Forbidden(c, "No permissions found")
			c.Abort()
			return
		}

		// Check for ai_agent_operate permission (allows all operations)
		if !auth.HasPermissionString(permStr, "ai_agent_operate") &&
			!auth.HasPermissionString(permStr, "ai_agent_query") &&
			!auth.HasPermissionString(permStr, "ai_agent_workflow") {
			response.Forbidden(c, "AI Agent permissions required")
			c.Abort()
			return
		}

		// Validate Agent ID header (should be present for audit)
		agentID := c.GetHeader("X-Agent-ID")
		if agentID == "" {
			// Allow requests without Agent ID for web client
			// But log it for debugging
			log.Printf("[AgentMiddleware] Warning: No X-Agent-ID header")
		} else {
			log.Printf("[AgentMiddleware] Agent ID: %s", agentID)
		}

		// Validate signature if provided (optional security measure)
		signature := c.GetHeader("X-Agent-Signature")
		if signature != "" {
			if !validateAgentSignature(c, signature) {
				response.Forbidden(c, "Invalid agent signature")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// validateAgentSignature validates the request signature
func validateAgentSignature(c *gin.Context, signature string) bool {
	// Get the shared secret from environment
	secret := "your-shared-secret" // In production, load from env

	// Create a signature from the request
	// In a real implementation, you would sign specific request parameters
	data := c.Request.URL.Path + c.Request.Method

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// CheckOperationPermission checks if user has permission for specific operation
func CheckOperationPermission(permissions string, operation string) bool {
	switch operation {
	case "query", "analyze":
		return auth.HasPermissionString(permissions, "ai_agent_query") ||
			auth.HasPermissionString(permissions, "ai_agent_operate")
	case "approve_workflow", "reject_workflow":
		return auth.HasPermissionString(permissions, "ai_agent_workflow") ||
			auth.HasPermissionString(permissions, "ai_agent_operate")
	case "create_material_plan", "update_stock", "generate_report":
		return auth.HasPermissionString(permissions, "ai_agent_operate")
	default:
		return auth.HasPermissionString(permissions, "ai_agent_operate")
	}
}
