package inbound

import (
	"time"

	"gorm.io/gorm"
)

// InboundOrder model maps to 'inbound_orders' table (unchanged)
type InboundOrder struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderNo      string    `gorm:"size:50;index" json:"order_no"`
	Supplier     string    `gorm:"size:100" json:"supplier"`
	Contact      string    `gorm:"size:50" json:"contact"`
	ProjectID    uint      `gorm:"index" json:"project_id"`
	PlanID       *uint     `gorm:"index" json:"plan_id"`
	CreatorID    uint      `gorm:"not null;index" json:"creator_id"`
	CreatorName  string    `gorm:"size:80;not null" json:"creator_name"`
	Status       string    `gorm:"size:20;default:'pending';index" json:"status"`
	Notes        string    `gorm:"type:text" json:"notes"`
	Remark       string    `gorm:"type:text" json:"remark"`
	TotalAmount  int       `gorm:"default:0" json:"total_amount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Items        []InboundOrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// InboundOrderItem model maps to 'inbound_items' table (v2 - simplified)
type InboundOrderItem struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	InboundOrderID  uint          `gorm:"not null;index" json:"inbound_order_id"`
	StockID         uint          `gorm:"index" json:"stock_id"`
	MaterialID      uint          `gorm:"not null;index" json:"material_id"`
	Quantity        float64       `gorm:"type:decimal(15,3);not null;default:0" json:"quantity"`
	UnitPrice       float64       `gorm:"type:decimal(15,2);default:0" json:"unit_price"`
	Status          string        `gorm:"size:20;default:'pending'" json:"status"`
	Remark          string        `gorm:"type:text" json:"remark"`
	CreatedAt       time.Time     `json:"created_at"`
	Order           *InboundOrder `gorm:"foreignKey:InboundOrderID" json:"-"`
}

// Status constants
const (
	StatusPending   = "pending"
	StatusApproved  = "approved"
	StatusRejected  = "rejected"
	StatusCompleted = "completed"
)

// ToDTOWithEnrichment returns order DTO with enriched items
func (o *InboundOrder) ToDTOWithEnrichment(db *gorm.DB) map[string]any {
	// Collect material IDs
	materialIDs := make([]uint, 0, len(o.Items))
	for _, item := range o.Items {
		materialIDs = append(materialIDs, item.MaterialID)
	}

	// Fetch materials in batch from material_master
	matMap := make(map[uint]map[string]any)
	if len(materialIDs) > 0 {
		type MaterialMaster struct {
			ID            uint
			Code          string
			Name          string
			Specification string
			Unit          string
		}
		var materials []MaterialMaster
		db.Table("material_master").Where("id IN ?", materialIDs).Find(&materials)
		for _, m := range materials {
			matMap[m.ID] = map[string]any{
				"code":          m.Code,
				"name":          m.Name,
				"specification": m.Specification,
				"unit":          m.Unit,
			}
		}
	}

	// Build enriched items
	itemsCount := len(o.Items)
	items := make([]map[string]any, 0, itemsCount)
	for _, item := range o.Items {
		itemDTO := item.ToDTO()
		if mat, ok := matMap[item.MaterialID]; ok {
			itemDTO["material_code"] = mat["code"]
			itemDTO["material_name"] = mat["name"]
			itemDTO["specification"] = mat["specification"]
			itemDTO["unit"] = mat["unit"]
		}
		items = append(items, itemDTO)
	}

	return map[string]any{
		"id":           o.ID,
		"order_no":     o.OrderNo,
		"supplier":     o.Supplier,
		"contact":      o.Contact,
		"project_id":   o.ProjectID,
		"plan_id":      o.PlanID,
		"creator_id":   o.CreatorID,
		"creator_name": o.CreatorName,
		"status":       o.Status,
		"notes":        o.Notes,
		"remark":       o.Remark,
		"total_amount": float64(o.TotalAmount) / 100.0,
		"items_count":  itemsCount,
		"items":        items,
		"created_at":   o.CreatedAt,
		"updated_at":   o.UpdatedAt,
	}
}

func (o *InboundOrder) ToDTO() map[string]any {
	itemsCount := 0
	items := make([]map[string]any, 0)
	if o.Items != nil {
		itemsCount = len(o.Items)
		for _, item := range o.Items {
			items = append(items, item.ToDTO())
		}
	}
	return map[string]any{
		"id":           o.ID,
		"order_no":     o.OrderNo,
		"supplier":     o.Supplier,
		"contact":      o.Contact,
		"project_id":   o.ProjectID,
		"plan_id":      o.PlanID,
		"creator_id":   o.CreatorID,
		"creator_name": o.CreatorName,
		"status":       o.Status,
		"notes":        o.Notes,
		"remark":       o.Remark,
		"total_amount": float64(o.TotalAmount) / 100.0,
		"items_count":  itemsCount,
		"items":        items,
		"created_at":   o.CreatedAt,
		"updated_at":   o.UpdatedAt,
	}
}

func (i *InboundOrderItem) ToDTO() map[string]any {
	return map[string]any{
		"id":              i.ID,
		"inbound_order_id": i.InboundOrderID,
		"stock_id":        i.StockID,
		"material_id":     i.MaterialID,
		"quantity":        i.Quantity,
		"unit_price":      i.UnitPrice,
		"status":          i.Status,
		"remark":          i.Remark,
		"created_at":      i.CreatedAt,
	}
}

// CreateInboundOrderRequest 创建入库单请求
type CreateInboundOrderRequest struct {
	OrderNo    string                       `json:"order_no" binding:"required"`
	ProjectID  uint                         `json:"project_id" binding:"required"`
	PlanID     *uint                        `json:"plan_id"`
	Supplier   string                       `json:"supplier"`
	Contact    string                       `json:"contact"`
	Notes      string                       `json:"notes"`
	Remark     string                       `json:"remark"`
	Items      []CreateInboundItemRequest   `json:"items" binding:"required"`
}

// CreateInboundItemRequest 创建入库单明细请求
type CreateInboundItemRequest struct {
	StockID    uint    `json:"stock_id" binding:"required"`
	MaterialID uint    `json:"material_id" binding:"required"`
	Quantity   float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice  float64 `json:"unit_price"`
	Remark     string  `json:"remark"`
}

// UpdateInboundOrderRequest 更新入库单请求
type UpdateInboundOrderRequest struct {
	Status   string                      `json:"status" binding:"required"`
	Notes    string                      `json:"notes"`
	Remark   string                      `json:"remark"`
	Items    []UpdateInboundItemRequest  `json:"items"`
}

// UpdateInboundItemRequest 更新入库单明细请求
type UpdateInboundItemRequest struct {
	ID        uint    `json:"id" binding:"required"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Status    string  `json:"status"`
	Remark    string  `json:"remark"`
}
