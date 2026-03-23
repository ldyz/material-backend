package agent

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
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
