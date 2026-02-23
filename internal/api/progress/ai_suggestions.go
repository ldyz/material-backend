package progress

import (
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
)

// JSON is a type alias for json.RawMessage
type JSON = json.RawMessage

// Suggestion represents an AI-powered suggestion for schedule optimization
type Suggestion struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ProjectID   uint           `json:"project_id" gorm:"index"`
	Type        string         `json:"type" gorm:"index"` // optimization, risk, resource, dependency
	Title       string         `json:"title" gorm:"size:255"`
	Description string         `json:"description" gorm:"type:text"`
	Impact      string         `json:"impact" gorm:"size:20"` // high, medium, low
	Actions     JSON           `json:"actions" gorm:"type:json"`
	Status      string         `json:"status" gorm:"index;size:20"` // pending, accepted, rejected, dismissed
	Metadata    JSON           `json:"metadata" gorm:"type:json"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
}

// ScheduleAnalysis represents the results of AI schedule analysis
type ScheduleAnalysis struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	ProjectID       uint      `json:"project_id" gorm:"index"`
	OverallScore    float64   `json:"overall_score"`
	RiskCount       int       `json:"risk_count"`
	SuggestionCount int       `json:"suggestion_count"`
	CriticalPath    JSON      `json:"critical_path" gorm:"type:json"`
	OverAllocations JSON      `json:"over_allocations" gorm:"type:json"`
	Optimizations   JSON      `json:"optimizations" gorm:"type:json"`
	AnalyzedAt      time.Time `json:"analyzed_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RiskAnalysis represents a schedule risk
type RiskAnalysis struct {
	TaskID     uint    `json:"task_id"`
	TaskName   string  `json:"task_name"`
	Type       string  `json:"type"` // long_duration, high_complexity, no_slack, resource_overload
	Level      string  `json:"level"` // high, medium, low
	Message    string  `json:"message"`
	Score      float64 `json:"score"`
	Factors    []string `json:"factors"`
}

// OptimizationSuggestion represents an optimization opportunity
type OptimizationSuggestion struct {
	Type          string                 `json:"type"` // parallelization, task_splitting, resource_leveling
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	PotentialSavings int                 `json:"potential_savings"` // in days or percentage
	Tasks         []OptimizationTask     `json:"tasks"`
	Actions       []SuggestedAction      `json:"actions"`
}

// OptimizationTask represents a task in an optimization
type OptimizationTask struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Duration int     `json:"duration"`
	Start    int     `json:"start"`
}

// SuggestedAction represents a specific action to take
type SuggestedAction struct {
	Type        string                 `json:"type"` // reschedule, reassign, split, remove, etc.
	Description string                 `json:"description"`
	Changes     map[string]interface{} `json:"changes"`
	Impact      string                 `json:"impact"` // expected impact
}

// ResourceOverload represents a resource allocation issue
type ResourceOverload struct {
	ResourceID   uint     `json:"resource_id"`
	ResourceName string   `json:"resource_name"`
	OverloadedDays []int  `json:"overloaded_days"`
	Severity     string   `json:"severity"` // high, medium, low
	Suggestions  []string `json:"suggestions"`
}

// AnalysisOptions represents options for schedule analysis
type AnalysisOptions struct {
	IncludeRisks        bool   `json:"include_risks"`
	IncludeOptimizations bool  `json:"include_optimizations"`
	CheckResources      bool   `json:"check_resources"`
	HistoricalData      bool   `json:"historical_data"`
	OptimizationLevel   string `json:"optimization_level"` // basic, standard, aggressive
}

// CreateSuggestionRequest represents a request to create a suggestion
type CreateSuggestionRequest struct {
	ProjectID   uint                   `json:"project_id" binding:"required"`
	Type        string                 `json:"type" binding:"required"`
	Title       string                 `json:"title" binding:"required"`
	Description string                 `json:"description"`
	Impact      string                 `json:"impact"`
	Actions     json.RawMessage        `json:"actions"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateSuggestionRequest represents a request to update a suggestion
type UpdateSuggestionRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Impact      string          `json:"impact"`
	Actions     json.RawMessage `json:"actions"`
	Status      string          `json:"status"`
	Metadata    json.RawMessage `json:"metadata"`
}

// SuggestionResponse represents a suggestion response
type SuggestionResponse struct {
	ID          uint                   `json:"id"`
	ProjectID   uint                   `json:"project_id"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Impact      string                 `json:"impact"`
	Actions     []SuggestedAction      `json:"actions"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
}

// ToResponse converts a Suggestion to SuggestionResponse
func (s *Suggestion) ToResponse() SuggestionResponse {
	var actions []SuggestedAction
	if s.Actions != nil {
		json.Unmarshal(s.Actions, &actions)
	}

	var metadata map[string]interface{}
	if s.Metadata != nil {
		json.Unmarshal(s.Metadata, &metadata)
	}

	return SuggestionResponse{
		ID:          s.ID,
		ProjectID:   s.ProjectID,
		Type:        s.Type,
		Title:       s.Title,
		Description: s.Description,
		Impact:      s.Impact,
		Actions:     actions,
		Status:      s.Status,
		Metadata:    metadata,
		CreatedAt:   s.CreatedAt,
	}
}

// ScheduleAnalysisResponse represents a schedule analysis response
type ScheduleAnalysisResponse struct {
	ID              uint                  `json:"id"`
	ProjectID       uint                  `json:"project_id"`
	OverallScore    float64               `json:"overall_score"`
	RiskCount       int                   `json:"risk_count"`
	SuggestionCount int                   `json:"suggestion_count"`
	CriticalPath    []uint                `json:"critical_path"`
	OverAllocations []ResourceOverload    `json:"over_allocations"`
	Optimizations   []OptimizationSuggestion `json:"optimizations"`
	AnalyzedAt      time.Time             `json:"analyzed_at"`
}

// ToResponse converts a ScheduleAnalysis to ScheduleAnalysisResponse
func (sa *ScheduleAnalysis) ToResponse() ScheduleAnalysisResponse {
	var criticalPath []uint
	if sa.CriticalPath != nil {
		json.Unmarshal(sa.CriticalPath, &criticalPath)
	}

	var overAllocations []ResourceOverload
	if sa.OverAllocations != nil {
		json.Unmarshal(sa.OverAllocations, &overAllocations)
	}

	var optimizations []OptimizationSuggestion
	if sa.Optimizations != nil {
		json.Unmarshal(sa.Optimizations, &optimizations)
	}

	return ScheduleAnalysisResponse{
		ID:              sa.ID,
		ProjectID:       sa.ProjectID,
		OverallScore:    sa.OverallScore,
		RiskCount:       sa.RiskCount,
		SuggestionCount: sa.SuggestionCount,
		CriticalPath:    criticalPath,
		OverAllocations: overAllocations,
		Optimizations:   optimizations,
		AnalyzedAt:      sa.AnalyzedAt,
	}
}

// GenerateAnalysisID generates a unique analysis ID
func GenerateAnalysisID() string {
	return uuid.New().String()
}

// CalculateRiskScore calculates a risk score for a task
func CalculateRiskScore(taskDuration int, dependencyCount int, hasSlack bool, isOverloaded bool) float64 {
	score := 0.0

	// Duration risk
	if taskDuration > 20 {
		score += 30
	} else if taskDuration > 14 {
		score += 20
	} else if taskDuration > 7 {
		score += 10
	}

	// Complexity risk
	if dependencyCount > 5 {
		score += 30
	} else if dependencyCount > 3 {
		score += 20
	} else if dependencyCount > 1 {
		score += 10
	}

	// Slack risk
	if !hasSlack {
		score += 25
	}

	// Resource risk
	if isOverloaded {
		score += 25
	}

	// Normalize to 0-100
	return math.Min(score, 100.0)
}
