package progress

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ChangeLogHandler handles change log-related HTTP requests
type ChangeLogHandler struct {
	repo *ChangeLogRepository
	db   *gorm.DB
}

// NewChangeLogHandler creates a new change log handler
func NewChangeLogHandler(repo *ChangeLogRepository, db *gorm.DB) *ChangeLogHandler {
	return &ChangeLogHandler{
		repo: repo,
		db:   db,
	}
}

// GetProjectChangeLog retrieves change log for a project
// @Summary Get project change log
// @Tags change-log
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param entity_type query string false "Entity type filter"
// @Param action_type query string false "Action type filter"
// @Param user_id query int false "User ID filter"
// @Param start_date query string false "Start date filter (ISO 8601)"
// @Param end_date query string false "End date filter (ISO 8601)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/project/{id}/change-log [get]
func (h *ChangeLogHandler) GetProjectChangeLog(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Build query options
	opts := QueryOptions{
		EntityType: c.Query("entity_type"),
		ActionType: c.Query("action_type"),
	}

	if userID := c.Query("user_id"); userID != "" {
		if id, err := strconv.ParseUint(userID, 10, 32); err == nil {
			opts.UserID = uint(id)
		}
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			opts.StartDate = t
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			opts.EndDate = t
		}
	}

	// Pagination
	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		if num, err := strconv.Atoi(p); err == nil && num > 0 {
			page = num
		}
	}
	if l := c.Query("limit"); l != "" {
		if num, err := strconv.Atoi(l); err == nil && num > 0 {
			limit = num
		}
	}
	opts.Offset = (page - 1) * limit
	opts.Limit = limit

	// Fetch logs
	logs, total, err := h.repo.GetByProjectID(uint(projectID), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve change log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetTaskHistory retrieves change history for a specific task
// @Summary Get task history
// @Tags change-log
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param limit query int false "Limit" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/tasks/{id}/history [get]
func (h *ChangeLogHandler) GetTaskHistory(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		if num, err := strconv.Atoi(l); err == nil && num > 0 {
			limit = num
		}
	}

	logs, err := h.repo.GetByEntityID("task", uint(taskID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// GetEntityHistory retrieves change history for any entity
// @Summary Get entity history
// @Tags change-log
// @Accept json
// @Produce json
// @Param entity_type path string true "Entity type"
// @Param entity_id path int true "Entity ID"
// @Param limit query int false "Limit" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/history/{entity_type}/{entity_id} [get]
func (h *ChangeLogHandler) GetEntityHistory(c *gin.Context) {
	entityType := c.Param("entity_type")
	entityID, err := strconv.ParseUint(c.Param("entity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		if num, err := strconv.Atoi(l); err == nil && num > 0 {
			limit = num
		}
	}

	logs, err := h.repo.GetByEntityID(entityType, uint(entityID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// RollbackChange rolls back a specific change
// @Summary Rollback change
// @Tags change-log
// @Accept json
// @Produce json
// @Param id path int true "Change ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/changes/{id}/rollback [post]
func (h *ChangeLogHandler) RollbackChange(c *gin.Context) {
	changeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid change ID"})
		return
	}

	// Get change log entry
	changeLog, err := h.repo.GetByID(uint(changeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Change not found"})
		return
	}

	// Get rollback data
	rollbackData, err := h.repo.GetRollbackData(uint(changeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rollback data"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Perform rollback based on entity type
	switch rollbackData.EntityType {
	case "task":
		err = h.rollbackTask(changeLog.ProjectID, rollbackData, userID.(uint))
	case "dependency":
		err = h.rollbackDependency(changeLog.ProjectID, rollbackData, userID.(uint))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported entity type for rollback"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback change"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Change rolled back successfully",
		"data":    rollbackData,
	})
}

// rollbackTask rolls back a task change
func (h *ChangeLogHandler) rollbackTask(projectID uint, data *RollbackData, userID uint) error {
	// Check if task exists
	var task Task
	err := h.db.Where("id = ? AND project_id = ?", data.EntityID, projectID).First(&task).Error
	if err != nil {
		// Task doesn't exist, create it
		task = Task{
			ProjectID: projectID,
		}
	}

	// Apply rollback data
	if name, ok := data.Data["name"].(string); ok {
		task.Name = name
	}
	if duration, ok := data.Data["duration"].(float64); ok {
		d := float64(duration)
		task.Duration = &d
	}
	if progress, ok := data.Data["progress"].(float64); ok {
		task.Progress = progress
	}
	if priority, ok := data.Data["priority"].(string); ok {
		task.Priority = priority
	}
	if status, ok := data.Data["status"].(string); ok {
		task.Status = status
	}
	if responsible, ok := data.Data["responsible"].(string); ok {
		task.Responsible = responsible
	}
	if description, ok := data.Data["description"].(string); ok {
		task.Description = description
	}

	// Save task
	if task.ID == 0 {
		return h.db.Create(&task).Error
	}
	return h.db.Save(&task).Error
}

// rollbackDependency rolls back a dependency change
func (h *ChangeLogHandler) rollbackDependency(projectID uint, data *RollbackData, userID uint) error {
	// Implementation would depend on dependency repository
	// For now, return not implemented error
	return nil
}

// GetStatistics retrieves change statistics for a project
// @Summary Get change statistics
// @Tags change-log
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/project/{id}/statistics [get]
func (h *ChangeLogHandler) GetStatistics(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	stats, err := h.repo.GetStatistics(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// DeleteOldChangeLogs deletes old change log entries
// @Summary Delete old change logs
// @Tags change-log
// @Accept json
// @Produce json
// @Param days query int true "Days to keep" default(90)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/change-log/cleanup [post]
func (h *ChangeLogHandler) DeleteOldChangeLogs(c *gin.Context) {
	days := 90
	if d := c.Query("days"); d != "" {
		if num, err := strconv.Atoi(d); err == nil && num > 0 {
			days = num
		}
	}

	if err := h.repo.DeleteOldEntries(days); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old entries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Old change logs deleted successfully",
		"data": gin.H{
			"days_kept": days,
		},
	})
}

// GetUserChanges retrieves all changes by current user
// @Summary Get user changes
// @Tags change-log
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/changes/my [get]
func (h *ChangeLogHandler) GetUserChanges(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		if num, err := strconv.Atoi(l); err == nil && num > 0 {
			limit = num
		}
	}

	logs, err := h.repo.GetByUserID(userID.(uint), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve changes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}
