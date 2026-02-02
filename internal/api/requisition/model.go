package requisition

import (
	"time"
)

// Requisition model maps to 'requisitions' table (unchanged)
type Requisition struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	RequisitionNo string            `gorm:"size:50;index" json:"requisition_no"`
	ProjectID     uint              `gorm:"index" json:"project_id"`
	PlanID        *uint             `gorm:"index" json:"plan_id"`
	Applicant     string            `gorm:"size:80" json:"applicant"`
	Department    string            `gorm:"size:80" json:"department"`
	Status        string            `gorm:"size:20" json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	Remark        string            `gorm:"size:200" json:"remark"`
	ApprovedAt    *time.Time        `json:"approved_at"`
	ApprovedBy    string            `gorm:"size:80" json:"approved_by"`
	Urgent        int               `gorm:"default:0" json:"urgent"`
	Purpose       string            `gorm:"size:200" json:"purpose"`
	IssuedBy      string            `gorm:"size:80" json:"issued_by"`
	IssuedAt      *time.Time        `json:"issued_at"`
	Items         []RequisitionItem `gorm:"foreignKey:RequisitionID" json:"items"`
}

// RequisitionItem model maps to 'requisition_items' table (v2 - simplified)
type RequisitionItem struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	RequisitionID     uint      `gorm:"index" json:"requisition_id"`
	StockID           uint      `gorm:"index" json:"stock_id"`
	MaterialID        uint      `gorm:"index" json:"material_id"`
	RequestedQuantity float64   `gorm:"type:decimal(15,3)" json:"requested_quantity"`
	ApprovedQuantity  float64   `gorm:"type:decimal(15,3)" json:"approved_quantity"`
	ActualQuantity    float64   `gorm:"type:decimal(15,3)" json:"actual_quantity"`
	Status            string    `gorm:"size:20" json:"status"`
	Remark            string    `gorm:"type:text" json:"remark"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (r *Requisition) ToDTO() map[string]any {
	itemsCount := len(r.Items)
	items := make([]map[string]any, 0, itemsCount)
	for _, item := range r.Items {
		items = append(items, item.ToDTO())
	}

	// Handle nullable datetime fields
	var approvedAtStr, issuedAtStr string
	if r.ApprovedAt != nil {
		approvedAtStr = r.ApprovedAt.Format("2006-01-02 15:04:05")
	}
	if r.IssuedAt != nil {
		issuedAtStr = r.IssuedAt.Format("2006-01-02 15:04:05")
	}

	return map[string]any{
		"id":                 r.ID,
		"requisition_no":     r.RequisitionNo,
		"requisition_number": r.RequisitionNo,
		"project_id":         r.ProjectID,
		"plan_id":            r.PlanID,
		"applicant":          r.Applicant,
		"applicant_name":     r.Applicant,
		"department":         r.Department,
		"status":             r.Status,
		"created_at":         r.CreatedAt.Format("2006-01-02 15:04:05"),
		"requisition_date":   r.CreatedAt.Format("2006-01-02"),
		"remark":             r.Remark,
		"approved_at":        approvedAtStr,
		"approved_by":        r.ApprovedBy,
		"urgent":             r.Urgent == 1,
		"purpose":            r.Purpose,
		"issued_by":          r.IssuedBy,
		"issued_at":          issuedAtStr,
		"items_count":        itemsCount,
		"items":              items,
	}
}

func (ri *RequisitionItem) ToDTO() map[string]any {
	return map[string]any{
		"id":                 ri.ID,
		"requisition_id":     ri.RequisitionID,
		"stock_id":           ri.StockID,
		"material_id":        ri.MaterialID,
		"requested_quantity": ri.RequestedQuantity,
		"approved_quantity":  ri.ApprovedQuantity,
		"actual_quantity":    ri.ActualQuantity,
		"remark":             ri.Remark,
		"status":             ri.Status,
		"created_at":         ri.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":         ri.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToDTOWithMaterial 转换为 DTO，包含物资信息
func (ri *RequisitionItem) ToDTOWithMaterial(materialCode, materialName, specification, unit string) map[string]any {
	dto := ri.ToDTO()
	dto["material_code"] = materialCode
	dto["material_name"] = materialName
	dto["specification"] = specification
	dto["unit"] = unit
	return dto
}

// CreateRequisitionRequest 创建领料单请求
type CreateRequisitionRequest struct {
	ProjectID  uint                         `json:"project_id" binding:"required"`
	PlanID     *uint                        `json:"plan_id"`
	Applicant  string                       `json:"applicant" binding:"required"`
	Department string                       `json:"department"`
	Purpose    string                       `json:"purpose"`
	Urgent     bool                         `json:"urgent"`
	Remark     string                       `json:"remark"`
	Items      []CreateRequisitionItemRequest `json:"items" binding:"required"`
}

// CreateRequisitionItemRequest 创建领料单明细请求
type CreateRequisitionItemRequest struct {
	StockID           uint    `json:"stock_id" binding:"required"`
	MaterialID        uint    `json:"material_id" binding:"required"`
	RequestedQuantity float64 `json:"requested_quantity" binding:"required,gt=0"`
	Remark            string  `json:"remark"`
}

// UpdateRequisitionRequest 更新领料单请求
type UpdateRequisitionRequest struct {
	Status    string                          `json:"status" binding:"required"`
	Remark    string                          `json:"remark"`
	Items     []UpdateRequisitionItemRequest  `json:"items"`
}

// UpdateRequisitionItemRequest 更新领料单明细请求
type UpdateRequisitionItemRequest struct {
	ID                uint    `json:"id" binding:"required"`
	ApprovedQuantity  float64 `json:"approved_quantity"`
	ActualQuantity    float64 `json:"actual_quantity"`
	Status            string  `json:"status"`
	Remark            string  `json:"remark"`
}
