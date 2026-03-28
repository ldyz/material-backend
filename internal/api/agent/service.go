package agent

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// Service handles agent business logic
type Service struct {
	db        *gorm.DB
	aiHandler *AIHandler
}

// NewService creates a new agent service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// SetAIHandler sets the AI handler
func (s *Service) SetAIHandler(handler *AIHandler) {
	s.aiHandler = handler
}

// GetAIHandler returns the AI handler
func (s *Service) GetAIHandler() *AIHandler {
	return s.aiHandler
}

// HandleOperation executes an AI Agent operation
func (s *Service) HandleOperation(req *AgentOperation, userID int, agentID string) (*AgentOperationResponse, error) {
	// Create log entry
	logEntry := &AgentOperationLog{
		Operation:  req.Operation,
		Resource:   req.Resource,
		Parameters: req.Parameters,
		Reasoning:  req.Reasoning,
		UserID:     &userID,
		AgentID:    agentID,
		Status:     "pending",
	}

	// Execute operation based on type
	var result any
	var err error
	var affectedRows int

	switch req.Operation {
	case OpQuery:
		result, affectedRows, err = s.handleQuery(req)
	case OpAnalyze:
		result, affectedRows, err = s.handleAnalyze(req)
	case OpCreatePlan:
		result, affectedRows, err = s.handleCreatePlan(req)
	case OpUpdateStock:
		result, affectedRows, err = s.handleUpdateStock(req)
	case OpApproveWorkflow:
		result, affectedRows, err = s.handleApproveWorkflow(req)
	case OpGenerateReport:
		result, affectedRows, err = s.handleGenerateReport(req)
	default:
		logEntry.Status = "failed"
		logEntry.Error = fmt.Sprintf("unknown operation: %s", req.Operation)
		if err := s.db.Create(logEntry).Error; err != nil {
			log.Printf("[Agent] Failed to create log entry: %v", err)
		}
		return nil, fmt.Errorf("unknown operation: %s", req.Operation)
	}

	// Update log entry with result
	if err != nil {
		logEntry.Status = "failed"
		logEntry.Error = err.Error()
	} else {
		logEntry.Status = "completed"
		logEntry.Result = result
	}

	if saveErr := s.db.Create(logEntry).Error; saveErr != nil {
		log.Printf("[Agent] Failed to save log entry: %v", saveErr)
	}

	if err != nil {
		return nil, err
	}

	return &AgentOperationResponse{
		Success:      true,
		Operation:    req.Operation,
		Result:       result.(map[string]any),
		Message:      fmt.Sprintf("Operation %s completed successfully", req.Operation),
		AffectedRows: affectedRows,
	}, nil
}

// handleQuery handles query operations
func (s *Service) handleQuery(req *AgentOperation) (map[string]any, int, error) {
	// Extract query parameters
	question, _ := req.Parameters["question"].(string)
	limit := 10
	if l, ok := req.Parameters["limit"].(float64); ok {
		limit = int(l)
	}

	// Build query based on resource
	var results []map[string]any
	var count int64

	switch req.Resource {
	case "material":
		query := s.db.Table("material_master").
			Select("id, name, specification, unit, category, description").
			Limit(limit)

		if search, ok := req.Parameters["search"].(string); ok && search != "" {
			query = query.Where("name ILIKE ? OR specification ILIKE ?",
				"%"+search+"%", "%"+search+"%")
		}

		if err := query.Find(&results).Error; err != nil {
			return nil, 0, err
		}
		query.Count(&count)

	case "stock":
		query := s.db.Table("stocks").
			Select("id, material_id, quantity, safety_stock").
			Limit(limit)

		if alert, ok := req.Parameters["low_stock_alert"].(bool); ok && alert {
			query = query.Where("quantity < safety_stock")
		}

		if err := query.Find(&results).Error; err != nil {
			return nil, 0, err
		}
		query.Count(&count)

	case "material_plan":
		query := s.db.Table("material_plans").
			Select("id, plan_no, project_id, status, total_budget").
			Limit(limit)

		if status, ok := req.Parameters["status"].(string); ok && status != "" {
			query = query.Where("status = ?", status)
		}

		if err := query.Find(&results).Error; err != nil {
			return nil, 0, err
		}
		query.Count(&count)

	default:
		return nil, 0, fmt.Errorf("unsupported resource for query: %s", req.Resource)
	}

	return map[string]any{
		"data":  results,
		"total": count,
		"question": question,
	}, len(results), nil
}

// handleAnalyze handles analyze operations (enhanced query)
func (s *Service) handleAnalyze(req *AgentOperation) (map[string]any, int, error) {
	// Analyze is similar to query but returns aggregated/processed data
	question, _ := req.Parameters["question"].(string)

	var result map[string]any
	var count int

	switch req.Resource {
	case "inventory":
		// Get inventory analysis
		var lowStockCount int64
		var totalStock float64
		var materialCount int64

		s.db.Table("stocks").Where("quantity < safety_stock").Count(&lowStockCount)
		s.db.Table("stocks").Select("COALESCE(SUM(quantity), 0)").Scan(&totalStock)
		s.db.Table("material_master").Count(&materialCount)

		result = map[string]any{
			"low_stock_count": lowStockCount,
			"total_stock":     totalStock,
			"material_count":  materialCount,
			"analysis":        "Inventory overview generated",
		}
		count = int(materialCount)

	case "workflow":
		// Get workflow analysis
		var pendingCount int64
		var approvedCount int64
		var rejectedCount int64

		s.db.Table("workflow_pending_tasks").Count(&pendingCount)
		s.db.Table("workflow_instances").Where("status = ?", "completed").Count(&approvedCount)
		s.db.Table("workflow_instances").Where("status = ?", "cancelled").Count(&rejectedCount)

		result = map[string]any{
			"pending_count":  pendingCount,
			"approved_count": approvedCount,
			"rejected_count": rejectedCount,
			"analysis":       "Workflow overview generated",
		}
		count = int(pendingCount)

	default:
		return nil, 0, fmt.Errorf("unsupported resource for analyze: %s", req.Resource)
	}

	result["question"] = question
	return result, count, nil
}

// handleCreatePlan handles material plan creation
func (s *Service) handleCreatePlan(req *AgentOperation) (map[string]any, int, error) {
	projectIDFloat, ok := req.Parameters["project_id"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("invalid project_id")
	}
	projectID := int(projectIDFloat)

	items, ok := req.Parameters["items"].([]any)
	if !ok {
		return nil, 0, fmt.Errorf("invalid items")
	}

	// Generate plan number
	var planNo string
	if err := s.db.Raw("SELECT generate_plan_no()").Scan(&planNo).Error; err != nil {
		// Fallback if function doesn't exist
		planNo = fmt.Sprintf("MP%d", time.Now().Unix())
	}

	// Create plan
	planData := map[string]any{
		"plan_no":   planNo,
		"project_id": projectID,
		"status":    "draft",
		"remark":    req.Parameters["remark"],
	}

	if err := s.db.Table("material_plans").Create(&planData).Error; err != nil {
		return nil, 0, err
	}

	planID := planData["id"]
	affectedRows := 1

	// Create plan items
	if len(items) > 0 {
		for _, item := range items {
			itemMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			itemMap["plan_id"] = planID
			s.db.Table("material_plan_items").Create(&itemMap)
			affectedRows++
		}
	}

	return map[string]any{
		"plan_id": planID,
		"plan_no": planNo,
		"message": "Material plan created successfully",
	}, affectedRows, nil
}

// handleUpdateStock handles stock updates
func (s *Service) handleUpdateStock(req *AgentOperation) (map[string]any, int, error) {
	stockIDFloat, ok := req.Parameters["stock_id"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("invalid stock_id")
	}
	stockID := int(stockIDFloat)

	quantityFloat, ok := req.Parameters["quantity"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("invalid quantity")
	}
	quantity := int(quantityFloat)

	remark, _ := req.Parameters["remark"].(string)

	// Update stock
	var result map[string]any
	if err := s.db.Table("stocks").
		Where("id = ?", stockID).
		Updates(map[string]any{
			"quantity": quantity,
		}).Error; err != nil {
		return nil, 0, err
	}

	// Create stock log
	logData := map[string]any{
		"stock_id":  stockID,
		"quantity":  quantity,
		"remark":    remark,
		"operation": "adjust",
	}
	s.db.Table("stock_op_logs").Create(&logData)

	result = map[string]any{
		"stock_id":  stockID,
		"quantity":  quantity,
		"message":   "Stock updated successfully",
	}

	return result, 1, nil
}

// handleApproveWorkflow handles workflow approval
func (s *Service) handleApproveWorkflow(req *AgentOperation) (map[string]any, int, error) {
	taskIDFloat, ok := req.Parameters["task_id"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("invalid task_id")
	}
	taskID := int64(taskIDFloat)

	remark, _ := req.Parameters["remark"].(string)

	// Get workflow task info
	var task map[string]any
	if err := s.db.Table("workflow_pending_tasks").
		Where("id = ?", taskID).
		First(&task).Error; err != nil {
		return nil, 0, fmt.Errorf("workflow task not found: %w", err)
	}

	// Get user ID from context
	userID, ok := req.Context["user_id"].(int)
	if !ok {
		return nil, 0, fmt.Errorf("invalid user_id in context")
	}

	// Create approval record
	approvalData := map[string]any{
		"instance_id": task["instance_id"],
		"node_id":     task["node_id"],
		"user_id":     userID,
		"status":      "approved",
		"comment":     remark,
	}
	if err := s.db.Table("workflow_approvals").Create(&approvalData).Error; err != nil {
		return nil, 0, err
	}

	// Update task status
	s.db.Table("workflow_pending_tasks").
		Where("id = ?", taskID).
		Update("status", "completed")

	// Move workflow to next node
	// This would typically involve complex workflow logic
	// For now, we'll mark the instance as completed
	s.db.Table("workflow_instances").
		Where("id = ?", task["instance_id"]).
		Update("status", "completed")

	return map[string]any{
		"task_id":     taskID,
		"instance_id": task["instance_id"],
		"status":      "approved",
		"message":     "Workflow approved successfully",
	}, 1, nil
}

// handleGenerateReport handles report generation
func (s *Service) handleGenerateReport(req *AgentOperation) (map[string]any, int, error) {
	reportType, _ := req.Parameters["report_type"].(string)

	switch reportType {
	case "inventory_summary":
		var summary map[string]any
		if err := s.db.Raw(`
			SELECT
				COUNT(*) as total_items,
				COALESCE(SUM(s.quantity), 0) as total_quantity,
				COUNT(CASE WHEN s.quantity < m.safety_stock THEN 1 END) as low_stock_items
			FROM stocks s
			JOIN material_master m ON s.material_id = m.id
		`).Scan(&summary).Error; err != nil {
			return nil, 0, err
		}
		return summary, 1, nil

	case "material_plan_summary":
		var summary map[string]any
		if err := s.db.Raw(`
			SELECT
				COUNT(*) as total_plans,
				COUNT(CASE WHEN status = 'draft' THEN 1 END) as draft_plans,
				COUNT(CASE WHEN status = 'active' THEN 1 END) as active_plans,
				COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_plans,
				COALESCE(SUM(total_budget), 0) as total_budget
			FROM material_plans
		`).Scan(&summary).Error; err != nil {
			return nil, 0, err
		}
		return summary, 1, nil

	default:
		return nil, 0, fmt.Errorf("unsupported report type: %s", reportType)
	}
}

// GetCapabilities returns available agent capabilities
func (s *Service) GetCapabilities() (*CapabilitiesResponse, error) {
	return &CapabilitiesResponse{
		Operations: []AgentCapability{
			{
				Name:        OpQuery,
				Description: "Query data from resources",
				Parameters: []ParameterDefinition{
					{Name: "question", Type: "string", Required: true, Description: "Natural language question"},
					{Name: "limit", Type: "integer", Required: false, Description: "Result limit"},
				},
			},
			{
				Name:        OpAnalyze,
				Description: "Analyze data and generate insights",
				Parameters: []ParameterDefinition{
					{Name: "question", Type: "string", Required: true, Description: "Analysis question"},
				},
			},
			{
				Name:        OpCreatePlan,
				Description: "Create material plan",
				Parameters: []ParameterDefinition{
					{Name: "project_id", Type: "integer", Required: true, Description: "Project ID"},
					{Name: "items", Type: "array", Required: true, Description: "Plan items"},
					{Name: "remark", Type: "string", Required: false, Description: "Remarks"},
				},
			},
			{
				Name:        OpUpdateStock,
				Description: "Update stock quantity",
				Parameters: []ParameterDefinition{
					{Name: "stock_id", Type: "integer", Required: true, Description: "Stock ID"},
					{Name: "quantity", Type: "integer", Required: true, Description: "New quantity"},
					{Name: "remark", Type: "string", Required: false, Description: "Remarks"},
				},
			},
			{
				Name:        OpApproveWorkflow,
				Description: "Approve workflow task",
				Parameters: []ParameterDefinition{
					{Name: "task_id", Type: "integer", Required: true, Description: "Task ID"},
					{Name: "remark", Type: "string", Required: false, Description: "Approval remarks"},
				},
			},
			{
				Name:        OpGenerateReport,
				Description: "Generate reports",
				Parameters: []ParameterDefinition{
					{Name: "report_type", Type: "string", Required: true, Description: "Report type"},
				},
			},
		},
		Resources: []AgentResource{
			{URI: "inventory://alerts", Description: "Low stock alerts", RefreshInterval: 300},
			{URI: "inventory://summary", Description: "Inventory summary", RefreshInterval: 60},
			{URI: "workflow://pending-tasks", Description: "Pending workflow tasks", RefreshInterval: 30},
			{URI: "material://plans", Description: "Material plans", RefreshInterval: 60},
		},
	}, nil
}

// ValidateOperation validates if an operation is allowed
func (s *Service) ValidateOperation(req *ValidateOperationRequest, userID int) (*ValidateOperationResponse, error) {
	// Check if operation is valid
	validOps := map[string]bool{
		OpQuery:           true,
		OpAnalyze:         true,
		OpCreatePlan:      true,
		OpUpdateStock:     true,
		OpApproveWorkflow: true,
		OpGenerateReport:  true,
	}

	if !validOps[req.Operation] {
		return &ValidateOperationResponse{
			Valid:   false,
			Message: fmt.Sprintf("Invalid operation: %s", req.Operation),
		}, nil
	}

	// Check required permissions
	var perms []string
	switch req.Operation {
	case OpQuery, OpAnalyze:
		perms = []string{"ai_agent_query"}
	case OpCreatePlan, OpUpdateStock, OpGenerateReport:
		perms = []string{"ai_agent_operate"}
	case OpApproveWorkflow:
		perms = []string{"ai_agent_workflow"}
	}

	// Check sensitive operations
	warnings := []string{}
	if IsSensitiveOperation(req.Operation) {
		warnings = append(warnings, "This is a sensitive operation that will be logged")
	}

	return &ValidateOperationResponse{
		Valid:           true,
		RequiredPerms:   perms,
		Warnings:        warnings,
	}, nil
}

// GetLogs retrieves agent operation logs
func (s *Service) GetLogs(params *AgentLogsQueryParams) ([]map[string]any, int64, error) {
	var logs []AgentOperationLog
	var total int64

	query := s.db.Model(&AgentOperationLog{})

	// Apply filters
	if params.Operation != "" {
		query = query.Where("operation = ?", params.Operation)
	}
	if params.Resource != "" {
		query = query.Where("resource = ?", params.Resource)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}
	if params.AgentID != "" {
		query = query.Where("agent_id = ?", params.AgentID)
	}
	if params.StartDate != "" {
		query = query.Where("created_at >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("created_at <= ?", params.EndDate)
	}

	// Count total
	query.Count(&total)

	// Apply pagination
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize).
		Order("created_at DESC")

	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	// Convert to map format
	result := make([]map[string]any, len(logs))
	for i, log := range logs {
		result[i] = map[string]any{
			"id":         log.ID,
			"operation":  log.Operation,
			"resource":   log.Resource,
			"parameters": log.Parameters,
			"reasoning":  log.Reasoning,
			"result":     log.Result,
			"user_id":    log.UserID,
			"agent_id":   log.AgentID,
			"status":     log.Status,
			"error":      log.Error,
			"created_at": log.CreatedAt,
		}
	}

	return result, total, nil
}
