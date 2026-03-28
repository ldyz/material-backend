package agent

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	openai "github.com/yourorg/material-backend/backend/pkg/openai"
	"gorm.io/gorm"
)

// Handler handles HTTP requests for agent operations
type Handler struct {
	service *Service
	db      *gorm.DB
}

// NewHandler creates a new agent handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		service: NewService(db),
		db:      db,
	}
}

// HandleOperation handles AI Agent operations
func (h *Handler) HandleOperation(c *gin.Context) {
	var req AgentOperation
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get user ID from context
	userIDVal, exists := c.Get("current_user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userID := int(userIDVal.(int64))

	// Get agent ID from header
	agentID := c.GetHeader("X-Agent-ID")
	if agentID == "" {
		agentID = "unknown"
	}

	// Add context info
	if req.Context == nil {
		req.Context = make(map[string]any)
	}
	req.Context["user_id"] = userID

	// Check for sensitive operations
	if IsSensitiveOperation(req.Operation) {
		// Additional validation for sensitive operations
		if req.Reasoning == "" {
			response.BadRequest(c, "Sensitive operations require reasoning explanation")
			return
		}
	}

	// Execute operation
	result, err := h.service.HandleOperation(&req, userID, agentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, result, result.Message)
}

// HandleQuery handles AI queries
func (h *Handler) HandleQuery(c *gin.Context) {
	var req AgentQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Set default limit
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Get user ID
	userIDVal, exists := c.Get("current_user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	userID := int(userIDVal.(int64))

	// Build operation request
	opReq := &AgentOperation{
		Operation: OpQuery,
		Resource:  "material", // Default resource
		Parameters: map[string]any{
			"question":           req.Question,
			"limit":              req.Limit,
			"fields":             req.Fields,
			"filters":            req.Filters,
			"order_by":           req.OrderBy,
		},
		Context: map[string]any{
			"user_id": userID,
		},
		Reasoning: "AI query from user",
	}

	// Execute
	agentID := c.GetHeader("X-Agent-ID")
	if agentID == "" {
		agentID = "web-client"
	}

	result, err := h.service.HandleOperation(opReq, userID, agentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// HandleWorkflow handles workflow operations
func (h *Handler) HandleWorkflow(c *gin.Context) {
	var req AgentWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get user ID
	userIDVal, exists := c.Get("current_user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	userID := int(userIDVal.(int64))

	// Build operation request
	opReq := &AgentOperation{
		Operation: OpApproveWorkflow,
		Resource:  "workflow",
		Parameters: map[string]any{
			"task_id":   req.TaskID,
			"action":    req.Action,
			"remark":    req.Remark,
			"to_node_id": req.ToNodeID,
		},
		Context: map[string]any{
			"user_id": userID,
		},
		Reasoning: "AI workflow operation",
	}

	agentID := c.GetHeader("X-Agent-ID")
	if agentID == "" {
		agentID = "web-client"
	}

	result, err := h.service.HandleOperation(opReq, userID, agentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, result, "Workflow operation completed successfully")
}

// HandleCapabilities returns available capabilities
func (h *Handler) HandleCapabilities(c *gin.Context) {
	caps, err := h.service.GetCapabilities()
	if err != nil {
		response.InternalError(c, "Failed to get capabilities")
		return
	}

	response.Success(c, caps)
}

// HandleValidate validates an operation
func (h *Handler) HandleValidate(c *gin.Context) {
	var req ValidateOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get user ID
	userIDVal, exists := c.Get("current_user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	userID := int(userIDVal.(int64))

	result, err := h.service.ValidateOperation(&req, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// HandleLogs retrieves agent operation logs
func (h *Handler) HandleLogs(c *gin.Context) {
	// Parse query parameters
	params := &AgentLogsQueryParams{
		Page:     1,
		PageSize: 20,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params.Page = p
		}
	}

	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			params.PageSize = ps
		}
	}

	params.Operation = c.Query("operation")
	params.Resource = c.Query("resource")
	params.Status = c.Query("status")
	params.AgentID = c.Query("agent_id")
	params.StartDate = c.Query("start_date")
	params.EndDate = c.Query("end_date")

	// Get current user
	user, err := auth.GetCurrentUser(c, h.db)
	if err != nil || user == nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Non-admin users can only see their own logs
	if !user.IsAdmin() {
		uid := int(user.ID)
		params.UserID = &uid
	}

	logs, total, err := h.service.GetLogs(params)
	if err != nil {
		response.InternalError(c, "Failed to retrieve logs")
		return
	}

	response.SuccessWithPagination(c, logs, int64(params.Page), int64(params.PageSize), total)
}

// GetCurrentUser gets the current authenticated user
func (h *Handler) GetCurrentUser(c *gin.Context) (int, error) {
	userIDVal, exists := c.Get("current_user_id")
	if !exists {
		return 0, nil
	}

	switch v := userIDVal.(type) {
	case int64:
		return int(v), nil
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, nil
	}
}

// HandleChat handles text chat requests
func (h *Handler) HandleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get user ID
	userID, err := h.GetCurrentUser(c)
	if err != nil || userID == 0 {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get AI handler from service
	aiHandler := h.service.GetAIHandler()
	if aiHandler == nil {
		response.InternalError(c, "AI service not configured")
		return
	}

	// Process chat
	ctx := c.Request.Context()
	resp, err := aiHandler.HandleAIChat(ctx, &AIChatRequest{
		Message:             req.Message,
		ConversationHistory: req.ConversationHistory,
		UserID:              userID,
		Context: map[string]interface{}{
			"user_id": userID,
		},
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// HandleVoiceChat handles voice chat requests
func (h *Handler) HandleVoiceChat(c *gin.Context) {
	// Get user ID
	userID, err := h.GetCurrentUser(c)
	if err != nil || userID == 0 {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get audio file from form
	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		response.BadRequest(c, "Audio file is required")
		return
	}
	defer file.Close()

	// Get AI handler from service
	aiHandler := h.service.GetAIHandler()
	if aiHandler == nil {
		response.InternalError(c, "AI service not configured")
		return
	}

	// Transcribe audio
	ctx := c.Request.Context()
	transcript, err := aiHandler.TranscribeAudio(ctx, file, header.Filename)
	if err != nil {
		response.InternalError(c, "Failed to transcribe audio: "+err.Error())
		return
	}

	// Process chat with transcribed text
	resp, err := aiHandler.HandleAIChat(ctx, &AIChatRequest{
		Message:             transcript,
		ConversationHistory: []openai.Message{},
		UserID:              userID,
		Context: map[string]interface{}{
			"user_id": userID,
		},
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// Return response with transcript
	response.Success(c, gin.H{
		"transcript": transcript,
		"message":    resp.Message,
		"response":   resp.Message,
	})
}

// ChatRequest represents a chat request
type ChatRequest struct {
	Message             string          `json:"message"`
	ConversationHistory []openai.Message `json:"conversation_history,omitempty"`
}

// HandleGetProviders 获取所有可用的模型提供者
func (h *Handler) HandleGetProviders(c *gin.Context) {
	aiHandler := h.service.GetAIHandler()
	if aiHandler == nil {
		response.InternalError(c, "AI service not configured")
		return
	}

	providers := aiHandler.GetProviders()
	currentProvider := aiHandler.GetCurrentProvider()

	// 构建响应
	type ProviderResponse struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Model   string `json:"model"`
		BaseURL string `json:"base_url"`
		Current bool   `json:"current"`
	}

	result := make([]ProviderResponse, 0, len(providers))
	for id, config := range providers {
		result = append(result, ProviderResponse{
			ID:      id,
			Name:    config.Name,
			Model:   config.Model,
			BaseURL: config.BaseURL,
			Current: id == currentProvider,
		})
	}

	response.Success(c, gin.H{
		"providers":        result,
		"current_provider": currentProvider,
	})
}

// HandleSwitchProvider 切换模型提供者
func (h *Handler) HandleSwitchProvider(c *gin.Context) {
	var req struct {
		Provider string `json:"provider" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	aiHandler := h.service.GetAIHandler()
	if aiHandler == nil {
		response.InternalError(c, "AI service not configured")
		return
	}

	if err := aiHandler.SwitchProvider(req.Provider); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":           "模型切换成功",
		"current_provider":  aiHandler.GetCurrentProvider(),
	})
}

// HandleClearConversationHistory 清除用户的对话历史
func (h *Handler) HandleClearConversationHistory(c *gin.Context) {
	// Get user ID
	userID, err := h.GetCurrentUser(c)
	if err != nil || userID == 0 {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Clear conversation history
	repo := NewConversationRepository(h.db)
	if err := repo.ClearHistory(int64(userID)); err != nil {
		response.InternalError(c, "Failed to clear conversation history")
		return
	}

	response.Success(c, gin.H{
		"message": "对话历史已清除",
	})
}

// HandleGetConversationHistory 获取用户的对话历史
func (h *Handler) HandleGetConversationHistory(c *gin.Context) {
	// Get user ID
	userID, err := h.GetCurrentUser(c)
	if err != nil || userID == 0 {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get limit from query
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	// Get conversation history
	repo := NewConversationRepository(h.db)
	messages, err := repo.GetRecentHistory(int64(userID), limit)
	if err != nil {
		response.InternalError(c, "Failed to get conversation history")
		return
	}

	response.Success(c, gin.H{
		"messages": messages,
	})
}
