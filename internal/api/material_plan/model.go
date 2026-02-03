package material_plan

import (
	"time"

	"gorm.io/gorm"
)

// MaterialPlan model maps to 'material_plans' table (unchanged)
type MaterialPlan struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	PlanNo             string     `gorm:"size:50;uniqueIndex;not null" json:"plan_no"`
	PlanName           string     `gorm:"size:200;not null" json:"plan_name"`
	ProjectID          uint       `gorm:"index;not null" json:"project_id"`
	PlanType           string     `gorm:"size:20;default:'procurement'" json:"plan_type"`
	Status             string     `gorm:"size:20;default:'draft';index" json:"status"`
	Priority           string     `gorm:"size:20;default:'normal'" json:"priority"`
	PlannedStartDate   *time.Time `json:"planned_start_date"`
	PlannedEndDate     *time.Time `json:"planned_end_date"`
	TotalBudget        int        `gorm:"default:0" json:"total_budget"`
	ActualCost         int        `gorm:"default:0" json:"actual_cost"`
	Description        string     `gorm:"type:text" json:"description"`
	Remark             string     `gorm:"type:text" json:"remark"`
	WorkflowInstanceID *uint      `gorm:"index" json:"workflow_instance_id"`
	CreatorID          uint       `gorm:"not null;index" json:"creator_id"`
	CreatorName        string     `gorm:"size:100;not null" json:"creator_name"`
	ApproverID         *uint      `json:"approver_id"`
	ApproverName       string     `gorm:"size:100" json:"approver_name"`
	ApprovedAt         *time.Time `json:"approved_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Items              []MaterialPlanItem `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// MaterialPlanItem model maps to 'material_plan_items' table (v2 - simplified)
type MaterialPlanItem struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	PlanID          uint       `gorm:"not null;index" json:"plan_id"`
	MaterialID      uint       `gorm:"not null;index" json:"material_id"`
	PlannedQuantity float64    `gorm:"type:decimal(15,3);not null;default:0" json:"planned_quantity"`
	UnitPrice       float64    `gorm:"type:decimal(15,2);default:0" json:"unit_price"`
	RequiredDate    *time.Time `json:"required_date"`
	Priority        string     `gorm:"size:20;default:'normal'" json:"priority"`
	Status          string     `gorm:"size:20;default:'pending'" json:"status"`
	Remark          string     `gorm:"type:text" json:"remark"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Plan            *MaterialPlan `gorm:"foreignKey:PlanID" json:"-"`
}

// Status constants for MaterialPlan
const (
	PlanStatusDraft      = "draft"
	PlanStatusPending    = "pending"
	PlanStatusApproved   = "approved"
	PlanStatusRejected   = "rejected"
	PlanStatusActive     = "active"
	PlanStatusCompleted  = "completed"
	PlanStatusCancelled  = "cancelled"
)

// Plan type constants
const (
	PlanTypeProcurement = "procurement"
	PlanTypeUsage       = "usage"
	PlanTypeMixed       = "mixed"
)

// Priority constants
const (
	PriorityLow      = "low"
	PriorityNormal   = "normal"
	PriorityHigh     = "high"
	PriorityUrgent   = "urgent"
)

// Item status constants
const (
	ItemStatusPending  = "pending"
	ItemStatusPartial  = "partial"
	ItemStatusCompleted = "completed"
	ItemStatusCancelled = "cancelled"
)

// CreateMaterialPlanRequest DTO for creating a plan
type CreateMaterialPlanRequest struct {
	PlanName           string                    `json:"plan_name" binding:"required"`
	ProjectID          uint                      `json:"project_id" binding:"required"`
	PlanType           string                    `json:"plan_type"`
	Priority           string                    `json:"priority"`
	PlannedStartDate   string                    `json:"planned_start_date"`
	PlannedEndDate     string                    `json:"planned_end_date"`
	TotalBudget        float64                   `json:"total_budget"`
	Description        string                    `json:"description"`
	Remark             string                    `json:"remark"`
	Items              []CreateMaterialPlanItemRequest `json:"items"`
}

// CreateMaterialPlanItemRequest DTO for creating plan items
type CreateMaterialPlanItemRequest struct {
	MaterialID      uint    `json:"material_id"`
	MaterialName    string  `json:"material_name"`
	MaterialCode    string  `json:"material_code"`
	Specification   string  `json:"specification"`
	Category        string  `json:"category"`
	Unit            string  `json:"unit"`
	PlannedQuantity float64 `json:"planned_quantity" binding:"required,gt=0"`
	UnitPrice       float64 `json:"unit_price"`
	RequiredDate    string  `json:"required_date"`
	Priority        string  `json:"priority"`
	Remark          string  `json:"remark"`
}

// UpdateMaterialPlanRequest DTO for updating a plan
type UpdateMaterialPlanRequest struct {
	PlanName           string                    `json:"plan_name"`
	PlanType           string                    `json:"plan_type"`
	Priority           string                    `json:"priority"`
	PlannedStartDate   string                    `json:"planned_start_date"`
	PlannedEndDate     string                    `json:"planned_end_date"`
	TotalBudget        float64                   `json:"total_budget"`
	Description        string                    `json:"description"`
	Remark             string                    `json:"remark"`
	Items              []CreateMaterialPlanItemRequest `json:"items"`
}

// ToDTO returns basic plan DTO
func (p *MaterialPlan) ToDTO() map[string]any {
	return map[string]any{
		"id":                 p.ID,
		"plan_no":            p.PlanNo,
		"plan_name":          p.PlanName,
		"project_id":         p.ProjectID,
		"plan_type":          p.PlanType,
		"status":             p.Status,
		"priority":           p.Priority,
		"planned_start_date": p.PlannedStartDate,
		"planned_end_date":   p.PlannedEndDate,
		"total_budget":       float64(p.TotalBudget) / 100.0,
		"actual_cost":        float64(p.ActualCost) / 100.0,
		"description":        p.Description,
		"remark":             p.Remark,
		"workflow_instance_id": p.WorkflowInstanceID,
		"creator_id":         p.CreatorID,
		"creator_name":       p.CreatorName,
		"approver_id":        p.ApproverID,
		"approver_name":      p.ApproverName,
		"approved_at":        p.ApprovedAt,
		"created_at":         p.CreatedAt,
		"updated_at":         p.UpdatedAt,
		"items_count":        len(p.Items),
	}
}

// ToDTOWithEnrichment returns plan DTO with enriched items and project info
func (p *MaterialPlan) ToDTOWithEnrichment(db *gorm.DB) map[string]any {
	dto := p.ToDTO()

	// Get project name
	var project struct {
		ID   uint
		Name string
	}
	if err := db.Table("projects").Where("id = ?", p.ProjectID).
		Select("id, name").First(&project).Error; err == nil {
		dto["project_name"] = project.Name
	}

	// Get all requisition IDs for this plan
	var requisitionIDs []uint
	db.Table("requisitions").Where("plan_id = ?", p.ID).Pluck("id", &requisitionIDs)

	// Enrich items with actual issued quantity from stock_logs
	items := make([]map[string]any, 0, len(p.Items))
	for _, item := range p.Items {
		itemDTO := item.ToDTO()

		// Get material details
		var material struct {
			Code          string
			Name          string
			Specification string
			Unit          string
		}
		if err := db.Table("material_master").Where("id = ?", item.MaterialID).
			Select("code, name, specification, unit").First(&material).Error; err == nil {
			itemDTO["material_code"] = material.Code
			itemDTO["material_name"] = material.Name
			itemDTO["specification"] = material.Specification
			itemDTO["unit"] = material.Unit
		}

		// Calculate actual issued quantity from stock_logs
		if len(requisitionIDs) > 0 {
			var stockID struct {
				ID uint
			}
			// Get stock_id for this material in this project
			if err := db.Table("stocks").Where("project_id = ? AND material_id = ?", p.ProjectID, item.MaterialID).
				Select("id").First(&stockID).Error; err == nil {

				// Calculate total issued quantity from stock_logs
				var totalIssued float64
				db.Table("stock_logs").
					Where("stock_id = ? AND type = ? AND source_type = ? AND source_id IN ?",
						stockID.ID, "out", "requisition", requisitionIDs).
					Select("COALESCE(SUM(quantity), 0)").
					Scan(&totalIssued)

				itemDTO["issued_quantity"] = totalIssued
				itemDTO["issue_progress"] = calculateProgress(totalIssued, item.PlannedQuantity)
				itemDTO["remaining_quantity"] = item.PlannedQuantity - totalIssued
			}
		}

		items = append(items, itemDTO)
	}
	dto["items"] = items
	dto["items_count"] = len(items)

	return dto
}

// ToDTO returns basic item DTO
func (i *MaterialPlanItem) ToDTO() map[string]any {
	return map[string]any{
		"id":               i.ID,
		"plan_id":          i.PlanID,
		"material_id":      i.MaterialID,
		"planned_quantity": i.PlannedQuantity,
		"unit_price":       i.UnitPrice,
		"required_date":    i.RequiredDate,
		"priority":         i.Priority,
		"status":           i.Status,
		"remark":           i.Remark,
		"created_at":       i.CreatedAt,
		"updated_at":       i.UpdatedAt,
	}
}

// calculateProgress calculates progress percentage
func calculateProgress(current, total float64) float64 {
	if total == 0 {
		return 0
	}
	return current / total * 100
}

// CalculateTotalBudget calculates total budget from items
func (p *MaterialPlan) CalculateTotalBudget() int {
	total := 0
	for _, item := range p.Items {
		total += int(item.UnitPrice * item.PlannedQuantity * 100) // 转换为分
	}
	p.TotalBudget = total
	return total
}

// CalculateProgress calculates overall progress
func (p *MaterialPlan) CalculateProgress() map[string]float64 {
	totalPlanned := 0.0
	totalIssued := 0.0

	for _, item := range p.Items {
		totalPlanned += item.PlannedQuantity

		// 从库存日志中获取实际发放数量
		if item.Plan != nil {
			// 这里需要从数据库查询，简化处理
		}
	}

	issueProgress := 0.0
	if totalPlanned > 0 {
		issueProgress = totalIssued / totalPlanned * 100
	}

	return map[string]float64{
		"issue_progress":   issueProgress,
		"overall_progress": issueProgress,
	}
}

// GetCompletionStatus returns completion status based on progress
func (p *MaterialPlan) GetCompletionStatus() string {
	progress := p.CalculateProgress()
	overall := progress["overall_progress"]

	if overall >= 100 {
		return PlanStatusCompleted
	} else if overall > 0 {
		return PlanStatusActive
	}
	return p.Status
}
