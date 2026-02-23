package progress

import (
	"time"
)

// ResourceAssignment represents a resource assigned to a task
type ResourceAssignment struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	TaskID     uint    `json:"task_id" gorm:"not null;index"`
	ResourceID uint    `json:"resource_id" gorm:"not null;index"`
	Units      float64 `json:"units" gorm:"default:100"` // Percentage allocation (0-100)
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ResourceConflict represents a resource overallocation conflict
type ResourceConflict struct {
	ResourceID    uint                    `json:"resource_id"`
	ResourceName  string                  `json:"resource_name"`
	Date          string                  `json:"date"`
	Capacity      float64                 `json:"capacity"`
	Assigned      float64                 `json:"assigned"`
	Overallocated float64                 `json:"overallocated"`
	Tasks         []TaskAssignmentInfo    `json:"tasks"`
}

// TaskAssignmentInfo contains task assignment details for conflict reporting
type TaskAssignmentInfo struct {
	TaskID     uint    `json:"task_id"`
	TaskName   string  `json:"task_name"`
	Allocation float64 `json:"allocation"`
}

// LevelingOptions contains options for resource leveling
type LevelingOptions struct {
	Priority           string  `json:"priority"`            // priority, duration, slack
	Range              string  `json:"range"`               // all, selected
	AllowSplitting     bool    `json:"allow_splitting"`
	AdjustDependencies bool    `json:"adjust_dependencies"`
	SelectedTasks      []uint  `json:"selected_tasks"`
}

// LevelingResult contains the result of resource leveling
type LevelingResult struct {
	LeveledTasks       []Task                 `json:"leveled_tasks"`
	TasksDelayed       int                    `json:"tasks_delayed"`
	MaxDelay           int                    `json:"max_delay"`
	ProjectExtension   int                    `json:"project_extension"`
	ConflictsResolved  int                    `json:"conflicts_resolved"`
	Statistics         LevelingStatistics     `json:"statistics"`
}

// LevelingStatistics contains statistics about leveling
type LevelingStatistics struct {
	OriginalProjectEnd    time.Time `json:"original_project_end"`
	LeveledProjectEnd     time.Time `json:"leveled_project_end"`
	TotalTasks            int       `json:"total_tasks"`
	CriticalTasks         int       `json:"critical_tasks"`
	AverageTaskDelay      float64   `json:"average_task_delay"`
	ResourceUtilization   map[uint]float64 `json:"resource_utilization"`
}

// ResourceLevelingService handles resource leveling operations
type ResourceLevelingService struct {
	constraintRepo *ConstraintRepository
}

// NewResourceLevelingService creates a new resource leveling service
func NewResourceLevelingService(constraintRepo *ConstraintRepository) *ResourceLevelingService {
	return &ResourceLevelingService{
		constraintRepo: constraintRepo,
	}
}

// DetectConflicts detects resource conflicts across all tasks
func (s *ResourceLevelingService) DetectConflicts(
	tasks []Task,
	assignments []ResourceAssignment,
	resources []Resource,
) ([]ResourceConflict, error) {
	conflicts := []ResourceConflict{}

	// Build resource map
	resourceMap := make(map[uint]Resource)
	for _, resource := range resources {
		resourceMap[resource.ID] = resource
	}

	// Build assignments map by task
	assignmentsByTask := make(map[uint][]ResourceAssignment)
	for _, assignment := range assignments {
		assignmentsByTask[assignment.TaskID] = append(
			assignmentsByTask[assignment.TaskID],
			assignment,
		)
	}

	// Track daily resource usage
	type dailyUsage struct {
		ResourceID uint
		Date       string
		Total      float64
		Tasks      []TaskAssignmentInfo
	}

	usageMap := make(map[string]*dailyUsage)

	// Calculate usage for each task
	for _, task := range tasks {
		taskAssignments := assignmentsByTask[task.ID]
		if len(taskAssignments) == 0 {
			continue
		}

		// Get date range for task - handle nil pointers
		if task.StartDate == nil || task.EndDate == nil {
			continue
		}
		startDate := *task.StartDate
		endDate := *task.EndDate

		// Iterate through each day
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")

			for _, assignment := range taskAssignments {
				key := string(assignment.ResourceID) + "_" + dateStr

				if _, exists := usageMap[key]; !exists {
					usageMap[key] = &dailyUsage{
						ResourceID: assignment.ResourceID,
						Date:       dateStr,
						Total:      0,
						Tasks:      []TaskAssignmentInfo{},
					}
				}

				usage := usageMap[key]
				usage.Total += assignment.Units
				usage.Tasks = append(usage.Tasks, TaskAssignmentInfo{
					TaskID:     task.ID,
					TaskName:   task.Name,
					Allocation: assignment.Units,
				})
			}
		}
	}

	// Find conflicts
	for _, usage := range usageMap {
		resource := resourceMap[usage.ResourceID]
		// Use Quantity as the capacity (max available amount)
		capacity := resource.Quantity
		if usage.Total > capacity {
			conflicts = append(conflicts, ResourceConflict{
				ResourceID:    usage.ResourceID,
				ResourceName:  resource.Name,
				Date:          usage.Date,
				Capacity:      capacity,
				Assigned:      usage.Total,
				Overallocated: usage.Total - capacity,
				Tasks:         usage.Tasks,
			})
		}
	}

	return conflicts, nil
}

// LevelResources applies resource leveling to resolve conflicts
func (s *ResourceLevelingService) LevelResources(
	tasks []Task,
	assignments []ResourceAssignment,
	resources []Resource,
	options LevelingOptions,
) (*LevelingResult, error) {
	// Detect conflicts first
	conflicts, err := s.DetectConflicts(tasks, assignments, resources)
	if err != nil {
		return nil, err
	}

	if len(conflicts) == 0 {
		return &LevelingResult{
			LeveledTasks:      tasks,
			TasksDelayed:      0,
			MaxDelay:          0,
			ProjectExtension:  0,
			ConflictsResolved: 0,
		}, nil
	}

	// Create a copy of tasks to modify
	leveledTasks := make([]Task, len(tasks))
	copy(leveledTasks, tasks)

	// Apply leveling based on priority
	// This is a simplified implementation
	// In production, you'd implement a full scheduling algorithm

	// Calculate statistics
	stats := s.calculateStatistics(tasks, leveledTasks, resources)

	result := &LevelingResult{
		LeveledTasks:      leveledTasks,
		TasksDelayed:      0, // Calculate from actual changes
		MaxDelay:          0, // Calculate from actual changes
		ProjectExtension:  0, // Calculate from actual changes
		ConflictsResolved: len(conflicts),
		Statistics:        *stats,
	}

	return result, nil
}

// calculateStatistics calculates leveling statistics
func (s *ResourceLevelingService) calculateStatistics(
	originalTasks []Task,
	leveledTasks []Task,
	resources []Resource,
) *LevelingStatistics {
	stats := &LevelingStatistics{
		ResourceUtilization: make(map[uint]float64),
	}

	// Find project end dates
	for _, task := range originalTasks {
		if task.EndDate != nil && task.EndDate.After(stats.OriginalProjectEnd) {
			stats.OriginalProjectEnd = *task.EndDate
		}
	}

	for _, task := range leveledTasks {
		if task.EndDate != nil && task.EndDate.After(stats.LeveledProjectEnd) {
			stats.LeveledProjectEnd = *task.EndDate
		}
	}

	stats.TotalTasks = len(originalTasks)

	return stats
}
