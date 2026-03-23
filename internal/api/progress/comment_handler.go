package progress

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CommentHandler handles comment-related HTTP requests
type CommentHandler struct {
	repo *CommentRepository
}

// NewCommentHandler creates a new comment handler
func NewCommentHandler(repo *CommentRepository) *CommentHandler {
	return &CommentHandler{repo: repo}
}

// CreateCommentRequest represents the request to create a comment
type CreateCommentRequest struct {
	Content  string   `json:"content" binding:"required"`
	Mentions []string `json:"mentions"`
	ParentID *uint    `json:"parent_id"`
}

// UpdateCommentRequest represents the request to update a comment
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// GetComments retrieves all comments for a task
// @Summary Get task comments
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {array} Comment
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/tasks/{id}/comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	comments, err := h.repo.GetByTaskID(uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// CreateComment creates a new comment
// @Summary Create comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body CreateCommentRequest true "Comment data"
// @Success 201 {object} Comment
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/tasks/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	comment := &Comment{
		TaskID:   uint(taskID),
		UserID:   userID.(uint),
		Content:  strings.TrimSpace(req.Content),
		Mentions: req.Mentions,
		ParentID: req.ParentID,
	}

	if err := h.repo.Create(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Load relations
	comment, _ = h.repo.GetByID(comment.ID)

	// Broadcast via WebSocket if available
	if GlobalHub != nil {
		BroadcastTaskChange(uint(taskID), "comment:create", comment.ID, comment)
	}

	c.JSON(http.StatusCreated, gin.H{"data": comment})
}

// UpdateComment updates a comment
// @Summary Update comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param request body UpdateCommentRequest true "Updated comment data"
// @Success 200 {object} Comment
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check ownership
	userID, exists := c.Get("userID")
	if !exists || comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this comment"})
		return
	}

	comment.Content = strings.TrimSpace(req.Content)

	if err := h.repo.Update(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comment})
}

// DeleteComment deletes a comment
// @Summary Delete comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	comment, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check ownership
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Allow admin or comment author
	isAdmin, _ := c.Get("isAdmin")
	if !isAdmin.(bool) && comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this comment"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// ResolveComment marks a comment as resolved
// @Summary Resolve comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} Comment
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/comments/{id}/resolve [post]
func (h *CommentHandler) ResolveComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if err := h.repo.Resolve(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve comment"})
		return
	}

	comment, _ := h.repo.GetByID(uint(id))
	c.JSON(http.StatusOK, gin.H{"data": comment})
}

// UnresolveComment marks a comment as unresolved
// @Summary Unresolve comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} Comment
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/comments/{id}/unresolve [post]
func (h *CommentHandler) UnresolveComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if err := h.repo.Unresolve(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unresolve comment"})
		return
	}

	comment, _ := h.repo.GetByID(uint(id))
	c.JSON(http.StatusOK, gin.H{"data": comment})
}

// GetProjectComments retrieves all comments for a project
// @Summary Get project comments
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} Comment
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/project/{id}/comments [get]
func (h *CommentHandler) GetProjectComments(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	comments, err := h.repo.GetByProjectID(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// GetTaskHistory retrieves comment history for a task
// @Summary Get task comment history
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {array} Comment
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/tasks/{id}/history [get]
func (h *CommentHandler) GetTaskHistory(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	comments, err := h.repo.GetByTaskID(uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// GetUserComments retrieves all comments by current user
// @Summary Get user comments
// @Tags comments
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Success 200 {array} Comment
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/comments/my [get]
func (h *CommentHandler) GetUserComments(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	comments, err := h.repo.GetByUserID(userID.(uint), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// GetMentions retrieves all comments mentioning current user
// @Summary Get mentions
// @Tags comments
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Success 200 {array} Comment
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/comments/mentions [get]
func (h *CommentHandler) GetMentions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	comments, err := h.repo.GetMentionsForUser(userID.(uint), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mentions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// GetUnresolvedCount retrieves unresolved comment count for a task
// @Summary Get unresolved count
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/tasks/{id}/comments/count [get]
func (h *CommentHandler) GetUnresolvedCount(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	count, err := h.repo.GetUnresolvedCount(uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"count": count}})
}
