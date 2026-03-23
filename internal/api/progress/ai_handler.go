package progress

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIHandler handles AI-related requests
type AIHandler struct {
	db *gorm.DB
}

// NewAIHandler creates a new AI handler
func NewAIHandler(db *gorm.DB) *AIHandler {
	return &AIHandler{db: db}
}

// AnalyzeScheduleRequest represents a request to analyze a project schedule
type AnalyzeScheduleRequest struct {
	IncludeRisks        bool   `json:"include_risks"`
	IncludeOptimizations bool  `json:"include_optimizations"`
	CheckResources      bool   `json:"check_resources"`
	HistoricalData      bool   `json:"historical_data"`
	OptimizationLevel   string `json:"optimization_level"`
}

// AnalyzeSchedule analyzes a project schedule and generates suggestions
// @Summary Analyze project schedule
// @Description Run AI analysis on a project schedule to identify risks and optimization opportunities
// @Tags progress
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body AnalyzeScheduleRequest true "Analysis options"
// @Success 200 {object} ScheduleAnalysisResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /progress/project/{id}/analyze [post]
func (h *AIHandler) AnalyzeSchedule(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req AnalyzeScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get project data
	var project struct {
		Name string `json:"name"`
	}
	if err := h.db.Table("projects").Select("name").Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Get tasks
	var tasks []struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Start     int    `json:"start"`
		Duration  int    `json:"duration"`
		Progress  int    `json:"progress"`
		IsMilestone bool `json:"is_milestone"`
		AssignedTo string `json:"assigned_to"`
		Remaining  int    `json:"remaining"`
	}
	if err := h.db.Table("activities").
		Where("project_id = ? AND deleted_at IS NULL", projectID).
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	// Get dependencies
	var dependencies []struct {
		ID           uint `json:"id"`
		PredecessorID uint `json:"predecessor_id"`
		SuccessorID   uint `json:"successor_id"`
		Type         string `json:"type"`
	}
	if err := h.db.Table("dependencies").
		Where("project_id = ?", projectID).
		Find(&dependencies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dependencies"})
		return
	}

	// Get resources if checking resources
	var resources []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	if req.CheckResources {
		if err := h.db.Table("resources").
			Where("project_id = ? OR is_global = ?", projectID, true).
			Find(&resources).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
			return
		}
	}

	// Run analysis
	analysis := h.runAnalysis(uint(projectID), tasks, dependencies, resources, req)

	// Save analysis to database
	analysisRecord := ScheduleAnalysis{
		ProjectID:       uint(projectID),
		OverallScore:    analysis.OverallScore,
		RiskCount:       analysis.RiskCount,
		SuggestionCount: analysis.SuggestionCount,
		AnalyzedAt:      time.Now(),
	}

	// Marshal data to JSON
	criticalPathJSON, _ := json.Marshal(analysis.CriticalPath)
	analysisRecord.CriticalPath = criticalPathJSON

	if len(analysis.OverAllocations) > 0 {
		overAllocationsJSON, _ := json.Marshal(analysis.OverAllocations)
		analysisRecord.OverAllocations = overAllocationsJSON
	}

	if len(analysis.Optimizations) > 0 {
		optimizationsJSON, _ := json.Marshal(analysis.Optimizations)
		analysisRecord.Optimizations = optimizationsJSON
	}

	if err := h.db.Create(&analysisRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save analysis"})
		return
	}

	c.JSON(http.StatusOK, analysisRecord.ToResponse())
}

// GetSuggestions retrieves all suggestions for a project
// @Summary Get project suggestions
// @Description Get all AI-generated suggestions for a project
// @Tags progress
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param status query string false "Filter by status (pending, accepted, rejected, dismissed)"
// @Param type query string false "Filter by type (optimization, risk, resource, dependency)"
// @Success 200 {array} SuggestionResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /progress/project/{id}/suggestions [get]
func (h *AIHandler) GetSuggestions(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	status := c.Query("status")
	suggestionType := c.Query("type")

	query := h.db.Model(&Suggestion{}).Where("project_id = ?", projectID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if suggestionType != "" {
		query = query.Where("type = ?", suggestionType)
	}

	var suggestions []Suggestion
	if err := query.Order("created_at DESC").Find(&suggestions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch suggestions"})
		return
	}

	responses := make([]SuggestionResponse, len(suggestions))
	for i, suggestion := range suggestions {
		responses[i] = suggestion.ToResponse()
	}

	c.JSON(http.StatusOK, responses)
}

// AcceptSuggestion accepts a suggestion and applies it
// @Summary Accept suggestion
// @Description Accept a suggestion and mark it as accepted
// @Tags progress
// @Accept json
// @Produce json
// @Param id path int true "Suggestion ID"
// @Success 200 {object} SuggestionResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /progress/suggestions/{id}/accept [post]
func (h *AIHandler) AcceptSuggestion(c *gin.Context) {
	suggestionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid suggestion ID"})
		return
	}

	var suggestion Suggestion
	if err := h.db.First(&suggestion, suggestionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Suggestion not found"})
		return
	}

	// Update status
	suggestion.Status = "accepted"
	if err := h.db.Save(&suggestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update suggestion"})
		return
	}

	// TODO: Apply the suggestion actions to the project
	// This would involve updating tasks, dependencies, resources, etc.

	c.JSON(http.StatusOK, suggestion.ToResponse())
}

// RejectSuggestion rejects a suggestion
// @Summary Reject suggestion
// @Description Reject a suggestion and mark it as rejected
// @Tags progress
// @Accept json
// @Produce json
// @Param id path int true "Suggestion ID"
// @Success 200 {object} SuggestionResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /progress/suggestions/{id}/reject [post]
func (h *AIHandler) RejectSuggestion(c *gin.Context) {
	suggestionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid suggestion ID"})
		return
	}

	var suggestion Suggestion
	if err := h.db.First(&suggestion, suggestionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Suggestion not found"})
		return
	}

	// Update status
	suggestion.Status = "rejected"
	if err := h.db.Save(&suggestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update suggestion"})
		return
	}

	c.JSON(http.StatusOK, suggestion.ToResponse())
}

// DismissSuggestion dismisses a suggestion
// @Summary Dismiss suggestion
// @Description Dismiss a suggestion and mark it as dismissed
// @Tags progress
// @Accept json
// @Produce json
// @Param id path int true "Suggestion ID"
// @Success 200 {object} SuggestionResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /progress/suggestions/{id}/dismiss [post]
func (h *AIHandler) DismissSuggestion(c *gin.Context) {
	suggestionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid suggestion ID"})
		return
	}

	var suggestion Suggestion
	if err := h.db.First(&suggestion, suggestionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Suggestion not found"})
		return
	}

	// Update status
	suggestion.Status = "dismissed"
	if err := h.db.Save(&suggestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update suggestion"})
		return
	}

	c.JSON(http.StatusOK, suggestion.ToResponse())
}

// runAnalysis performs the actual schedule analysis
func (h *AIHandler) runAnalysis(
	projectID uint,
	tasks []struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Start     int    `json:"start"`
		Duration  int    `json:"duration"`
		Progress  int    `json:"progress"`
		IsMilestone bool `json:"is_milestone"`
		AssignedTo string `json:"assigned_to"`
		Remaining  int    `json:"remaining"`
	},
	dependencies []struct {
		ID           uint `json:"id"`
		PredecessorID uint `json:"predecessor_id"`
		SuccessorID   uint `json:"successor_id"`
		Type         string `json:"type"`
	},
	resources []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	},
	req AnalyzeScheduleRequest,
) ScheduleAnalysisResponse {
	analysis := ScheduleAnalysisResponse{
		ProjectID: projectID,
	}

	overallScore := 100.0
	riskCount := 0
	suggestionCount := 0

	// 1. Critical Path Analysis
	criticalPath := h.calculateCriticalPath(tasks, dependencies)
	analysis.CriticalPath = criticalPath

	// 2. Risk Analysis
	if req.IncludeRisks {
		risks := h.analyzeRisks(tasks, dependencies)
		riskCount = len(risks)

		// Generate suggestions for risks
		for _, risk := range risks {
			if risk.Score >= 50 {
				suggestionCount++
				// Create suggestion in database
				actions := []SuggestedAction{
					{
						Type:        "mitigate",
						Description: risk.Message,
						Changes: map[string]interface{}{
							"task_id": risk.TaskID,
							"risk_type": risk.Type,
						},
						Impact: risk.Level,
					},
				}
				actionsJSON, _ := json.Marshal(actions)

				suggestion := Suggestion{
					ProjectID:   projectID,
					Type:        "risk",
					Title:       "Risk Detected: " + risk.TaskName,
					Description: risk.Message,
					Impact:      risk.Level,
					Actions:     actionsJSON,
					Status:      "pending",
				}
				h.db.Create(&suggestion)

				// Reduce overall score based on risk level
				if risk.Level == "high" {
					overallScore -= 10
				} else if risk.Level == "medium" {
					overallScore -= 5
				}
			}
		}
	}

	// 3. Optimization Analysis
	if req.IncludeOptimizations {
		optimizations := h.findOptimizations(tasks, dependencies)
		analysis.Optimizations = optimizations

		// Generate suggestions for optimizations
		for _, opt := range optimizations {
			if opt.PotentialSavings >= 2 {
				suggestionCount++
				actions := make([]SuggestedAction, len(opt.Actions))
				for i, action := range opt.Actions {
					actions[i] = action
				}
				actionsJSON, _ := json.Marshal(actions)

				suggestion := Suggestion{
					ProjectID:   projectID,
					Type:        "optimization",
					Title:       opt.Title,
					Description: opt.Description,
					Impact:      "medium",
					Actions:     actionsJSON,
					Status:      "pending",
				}
				h.db.Create(&suggestion)
			}
		}
	}

	// 4. Resource Analysis
	if req.CheckResources && len(resources) > 0 {
		overAllocations := h.analyzeResourceAllocation(tasks, resources)
		analysis.OverAllocations = overAllocations

		for _, overload := range overAllocations {
			if overload.Severity == "high" || overload.Severity == "medium" {
				suggestionCount++
				actions := []SuggestedAction{
					{
						Type:        "reassign",
						Description: "Reassign tasks to reduce overallocation",
						Changes: map[string]interface{}{
							"resource_id": overload.ResourceID,
						},
						Impact: "high",
					},
				}
				actionsJSON, _ := json.Marshal(actions)

				suggestion := Suggestion{
					ProjectID:   projectID,
					Type:        "resource",
					Title:       "Resource Overallocation: " + overload.ResourceName,
					Description: overload.ResourceName + " is overloaded on " + strconv.Itoa(len(overload.OverloadedDays)) + " days",
					Impact:      overload.Severity,
					Actions:     actionsJSON,
					Status:      "pending",
				}
				h.db.Create(&suggestion)

				overallScore -= 10
			}
		}
	}

	// Ensure score is within 0-100
	analysis.OverallScore = math.Max(0, math.Min(100, overallScore))
	analysis.RiskCount = riskCount
	analysis.SuggestionCount = suggestionCount
	analysis.AnalyzedAt = time.Now()

	return analysis
}

// calculateCriticalPath calculates the critical path
func (h *AIHandler) calculateCriticalPath(
	tasks []struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Start     int    `json:"start"`
		Duration  int    `json:"duration"`
		Progress  int    `json:"progress"`
		IsMilestone bool `json:"is_milestone"`
		AssignedTo string `json:"assigned_to"`
		Remaining  int    `json:"remaining"`
	},
	dependencies []struct {
		ID           uint `json:"id"`
		PredecessorID uint `json:"predecessor_id"`
		SuccessorID   uint `json:"successor_id"`
		Type         string `json:"type"`
	},
) []uint {
	// Build adjacency list
	graph := make(map[uint][]uint)
	inDegree := make(map[uint]int)

	for _, task := range tasks {
		graph[task.ID] = []uint{}
		inDegree[task.ID] = 0
	}

	for _, dep := range dependencies {
		graph[dep.PredecessorID] = append(graph[dep.PredecessorID], dep.SuccessorID)
		inDegree[dep.SuccessorID]++
	}

	// Find start tasks (no incoming edges)
	queue := []uint{}
	for _, task := range tasks {
		if inDegree[task.ID] == 0 {
			queue = append(queue, task.ID)
		}
	}

	// Topological sort and calculate longest path
	dist := make(map[uint]int)
	for _, task := range tasks {
		dist[task.ID] = task.Start + task.Duration
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range graph[current] {
			currentTask := h.findTask(tasks, current)
			neighborTask := h.findTask(tasks, neighbor)

			if currentTask != nil && neighborTask != nil {
				if dist[current] > neighborTask.Start {
					dist[neighbor] = dist[current] + neighborTask.Duration
				}
			}

			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Find the task with maximum distance
	maxDist := 0
	var lastTask uint
	for taskID, distance := range dist {
		if distance > maxDist {
			maxDist = distance
			lastTask = taskID
		}
	}

	// Reconstruct path (simplified)
	path := []uint{lastTask}
	// TODO: Implement proper path reconstruction

	return path
}

// analyzeRisks analyzes schedule risks
func (h *AIHandler) analyzeRisks(
	tasks []struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Start     int    `json:"start"`
		Duration  int    `json:"duration"`
		Progress  int    `json:"progress"`
		IsMilestone bool `json:"is_milestone"`
		AssignedTo string `json:"assigned_to"`
		Remaining  int    `json:"remaining"`
	},
	dependencies []struct {
		ID           uint `json:"id"`
		PredecessorID uint `json:"predecessor_id"`
		SuccessorID   uint `json:"successor_id"`
		Type         string `json:"type"`
	},
) []RiskAnalysis {
	risks := []RiskAnalysis{}

	for _, task := range tasks {
		score := 0.0
		factors := []string{}

		// Duration risk
		if task.Duration > 20 {
			score += 30
			factors = append(factors, "Long duration")
		}

		// Dependency complexity
		depCount := 0
		for _, dep := range dependencies {
			if dep.PredecessorID == task.ID || dep.SuccessorID == task.ID {
				depCount++
			}
		}
		if depCount > 5 {
			score += 30
			factors = append(factors, "High complexity")
		}

		// Progress risk
		if task.Remaining > 0 && task.Duration > 0 {
			elapsed := task.Duration - task.Remaining
			progress := float64(elapsed) / float64(task.Duration)
			if progress > 0.8 && task.Progress < 50 {
				score += 35
				factors = append(factors, "Behind schedule")
			}
		}

		// Determine risk level
		level := "low"
		if score >= 70 {
			level = "high"
		} else if score >= 40 {
			level = "medium"
		}

		if score > 0 {
			risks = append(risks, RiskAnalysis{
				TaskID:   task.ID,
				TaskName: task.Name,
				Type:     "schedule_risk",
				Level:    level,
				Message:  "Task " + task.Name + " has risk score: " + strconv.Itoa(int(score)),
				Score:    score,
				Factors:  factors,
			})
		}
	}

	return risks
}

// findOptimizations finds optimization opportunities
func (h *AIHandler) findOptimizations(
	tasks []struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Start     int    `json:"start"`
		Duration  int    `json:"duration"`
		Progress  int    `json:"progress"`
		IsMilestone bool `json:"is_milestone"`
		AssignedTo string `json:"assigned_to"`
		Remaining  int    `json:"remaining"`
	},
	dependencies []struct {
		ID           uint `json:"id"`
		PredecessorID uint `json:"predecessor_id"`
		SuccessorID   uint `json:"successor_id"`
		Type         string `json:"type"`
	},
) []OptimizationSuggestion {
	optimizations := []OptimizationSuggestion{}

	// Find long tasks that can be split
	for _, task := range tasks {
		if task.Duration > 10 && !task.IsMilestone {
			splits := int(math.Ceil(float64(task.Duration) / 7))
			if splits > 1 {
				optimizations = append(optimizations, OptimizationSuggestion{
					Type:              "task_splitting",
					Title:             "Split Long Task: " + task.Name,
					Description:       "This task can be split into " + strconv.Itoa(splits) + " smaller tasks",
					PotentialSavings:   0,
					Tasks:             []OptimizationTask{{ID: task.ID, Name: task.Name, Duration: task.Duration, Start: task.Start}},
					Actions: []SuggestedAction{
						{
							Type:        "split",
							Description: "Split task into " + strconv.Itoa(splits) + " parts",
							Changes:     map[string]interface{}{"task_id": task.ID, "splits": splits},
							Impact:      "medium",
						},
					},
				})
			}
		}
	}

	return optimizations
}

// analyzeResourceAllocation analyzes resource allocation
func (h *AIHandler) analyzeResourceAllocation(
	tasks []struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Start     int    `json:"start"`
		Duration  int    `json:"duration"`
		Progress  int    `json:"progress"`
		IsMilestone bool `json:"is_milestone"`
		AssignedTo string `json:"assigned_to"`
		Remaining  int    `json:"remaining"`
	},
	resources []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	},
) []ResourceOverload {
	overAllocations := []ResourceOverload{}

	// Count tasks per resource per day
	resourceLoad := make(map[string]map[int]int)
	for _, task := range tasks {
		if task.AssignedTo == "" {
			continue
		}
		if resourceLoad[task.AssignedTo] == nil {
			resourceLoad[task.AssignedTo] = make(map[int]int)
		}
		for day := task.Start; day < task.Start+task.Duration; day++ {
			resourceLoad[task.AssignedTo][day]++
		}
	}

	// Check for overloads
	for _, resource := range resources {
		overloadedDays := []int{}
		for day, load := range resourceLoad[resource.Name] {
			if load > 1 {
				overloadedDays = append(overloadedDays, day)
			}
		}

		if len(overloadedDays) > 0 {
			severity := "low"
			if len(overloadedDays) > 10 {
				severity = "high"
			} else if len(overloadedDays) > 5 {
				severity = "medium"
			}

			overAllocations = append(overAllocations, ResourceOverload{
				ResourceID:      resource.ID,
				ResourceName:    resource.Name,
				OverloadedDays:  overloadedDays,
				Severity:        severity,
				Suggestions:     []string{"Reassign some tasks to other resources", "Adjust task dates to reduce overlap"},
			})
		}
	}

	return overAllocations
}

// findTask finds a task by ID
func (h *AIHandler) findTask(
	tasks []struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Start     int    `json:"start"`
		Duration  int    `json:"duration"`
		Progress  int    `json:"progress"`
		IsMilestone bool `json:"is_milestone"`
		AssignedTo string `json:"assigned_to"`
		Remaining  int    `json:"remaining"`
	},
	id uint,
) *struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Start     int    `json:"start"`
	Duration  int    `json:"duration"`
	Progress  int    `json:"progress"`
	IsMilestone bool `json:"is_milestone"`
	AssignedTo string `json:"assigned_to"`
	Remaining  int    `json:"remaining"`
} {
	for _, task := range tasks {
		if task.ID == id {
			return &task
		}
	}
	return nil
}
