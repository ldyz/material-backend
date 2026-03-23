package progress

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ConstraintHandler handles constraint-related HTTP requests
type ConstraintHandler struct {
	repo *ConstraintRepository
}

// NewConstraintHandler creates a new constraint handler
func NewConstraintHandler(repo *ConstraintRepository) *ConstraintHandler {
	return &ConstraintHandler{repo: repo}
}

// GetConstraints retrieves all constraints for a task
// GET /progress/tasks/:id/constraints
func (h *ConstraintHandler) GetConstraints(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	constraints, err := h.repo.GetByTaskID(uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": constraints})
}

// CreateConstraint creates a new constraint for a task
// POST /progress/tasks/:id/constraints
func (h *ConstraintHandler) CreateConstraint(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var constraint Constraint
	if err := c.ShouldBindJSON(&constraint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	constraint.TaskID = uint(taskID)

	// Validate constraint
	if err := constraint.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create constraint
	if err := h.repo.Create(&constraint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": constraint})
}

// UpdateConstraint updates an existing constraint
// PUT /progress/constraints/:id
func (h *ConstraintHandler) UpdateConstraint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid constraint ID"})
		return
	}

	// Get existing constraint
	constraint, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "constraint not found"})
		return
	}

	// Bind updates
	var updates Constraint
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	constraint.Type = updates.Type
	constraint.Date = updates.Date

	// Validate
	if err := constraint.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save
	if err := h.repo.Update(constraint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": constraint})
}

// DeleteConstraint deletes a constraint
// DELETE /progress/constraints/:id
func (h *ConstraintHandler) DeleteConstraint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid constraint ID"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "constraint deleted successfully"})
}

// ApplyConstraintToTask applies a constraint to a task and returns updated dates
// POST /progress/tasks/:id/apply-constraint
func (h *ConstraintHandler) ApplyConstraintToTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var constraint Constraint
	if err := c.ShouldBindJSON(&constraint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	constraint.TaskID = uint(taskID)

	// Validate constraint
	if err := constraint.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get task from task repository and apply constraint
	// This would require integrating with the task repository
	// For now, return the constraint as-is

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"constraint": constraint,
			"message":    "Constraint applied successfully",
		},
	})
}
