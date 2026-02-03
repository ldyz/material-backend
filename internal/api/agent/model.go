package agent

import "time"

// Operation types supported by AI Agent
const (
	// Query operations
	OpQuery     = "query"
	OpAnalyze   = "analyze"

	// Business operations
	OpCreatePlan      = "create_material_plan"
	OpUpdateStock     = "update_stock"
	OpApproveWorkflow = "approve_workflow"

	// Report operations
	OpGenerateReport = "generate_report"
)

// AgentOperation represents an AI Agent operation request
type AgentOperation struct {
	Operation  string                 `json:"operation" binding:"required"`  // Operation type
	Resource   string                 `json:"resource" binding:"required"`   // Resource type
	Parameters map[string]any         `json:"parameters"`                    // Operation parameters
	Context    map[string]any         `json:"context"`                       // Context information
	Reasoning  string                 `json:"reasoning"`                     // AI reasoning process
}

// AgentOperationResponse represents an AI Agent operation response
type AgentOperationResponse struct {
	Success      bool         `json:"success"`
	Operation    string       `json:"operation"`
	Result       map[string]any `json:"result"`
	Message      string       `json:"message"`
	AffectedRows int          `json:"affected_rows"`
}

// AgentQueryRequest represents an AI query request
type AgentQueryRequest struct {
	Question          string   `json:"question" binding:"required"`
	Limit             int      `json:"limit"`
	Fields            []string `json:"fields"`
	Filters           map[string]any `json:"filters"`
	OrderBy           string   `json:"order_by"`
	ConversationMode  bool     `json:"conversation_mode"`
	ConversationID    string   `json:"conversation_id"`
	ConversationHistory []map[string]any `json:"conversation_history"`
	MaxIterations     int      `json:"max_iterations"`
}

// AgentWorkflowRequest represents an AI workflow operation request
type AgentWorkflowRequest struct {
	TaskID    int64  `json:"task_id" binding:"required"`
	Action    string `json:"action" binding:"required"` // approve/reject/return
	Remark    string `json:"remark"`
	ToNodeID  *int   `json:"to_node_id,omitempty"`     // Used for return action
}

// AgentCapability represents a single capability
type AgentCapability struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  []ParameterDefinition `json:"parameters,omitempty"`
}

// ParameterDefinition defines a parameter for a capability
type ParameterDefinition struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// AgentResource represents an AI resource
type AgentResource struct {
	URI         string `json:"uri"`
	Description string `json:"description"`
	RefreshInterval int `json:"refresh_interval,omitempty"`
}

// CapabilitiesResponse represents the capabilities response
type CapabilitiesResponse struct {
	Operations []AgentCapability `json:"operations"`
	Resources  []AgentResource   `json:"resources"`
}

// ValidateOperationRequest represents a validation request
type ValidateOperationRequest struct {
	Operation  string                 `json:"operation" binding:"required"`
	Resource   string                 `json:"resource" binding:"required"`
	Parameters map[string]any         `json:"parameters"`
}

// ValidateOperationResponse represents a validation response
type ValidateOperationResponse struct {
	Valid       bool   `json:"valid"`
	Message     string `json:"message,omitempty"`
	Warnings    []string `json:"warnings,omitempty"`
	RequiredPerms []string `json:"required_permissions,omitempty"`
}

// AgentOperationLog represents an agent operation log entry in database
type AgentOperationLog struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	Operation  string    `gorm:"type:varchar(100);not null" json:"operation"`
	Resource   string    `gorm:"type:varchar(100);not null" json:"resource"`
	Parameters any       `gorm:"type:jsonb" json:"parameters"`
	Reasoning  string    `gorm:"type:text" json:"reasoning"`
	Result     any       `gorm:"type:jsonb" json:"result"`
	UserID     *int      `gorm:"index" json:"user_id"`
	AgentID    string    `gorm:"type:varchar(255)" json:"agent_id"`
	Status     string    `gorm:"type:varchar(50);not null" json:"status"` // pending/completed/failed
	Error      string    `gorm:"type:text" json:"error,omitempty"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName specifies the table name for AgentOperationLog
func (AgentOperationLog) TableName() string {
	return "agent_operation_logs"
}

// AgentLogsQueryParams represents query parameters for agent logs
type AgentLogsQueryParams struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Operation  string `json:"operation"`
	Resource   string `json:"resource"`
	Status     string `json:"status"`
	UserID     *int   `json:"user_id"`
	AgentID    string `json:"agent_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}

// Sensitive operations that require additional validation
var SensitiveOperations = []string{
	"approve_workflow",
	"delete_material",
	"update_stock",
	"delete_plan",
	"cancel_plan",
}

// IsSensitiveOperation checks if an operation is sensitive
func IsSensitiveOperation(op string) bool {
	for _, s := range SensitiveOperations {
		if op == s {
			return true
		}
	}
	return false
}
