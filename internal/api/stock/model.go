package stock

import (
	"time"
)

// Stock model maps to 'stocks' table (v2)
type Stock struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"index;not null" json:"project_id"`
	MaterialID  uint      `gorm:"index;not null" json:"material_id"`
	WarehouseID *uint     `gorm:"index" json:"warehouse_id"` // 仓库ID（可选）
	Quantity    float64   `gorm:"type:decimal(15,3);default:0" json:"quantity"`
	SafetyStock float64   `gorm:"type:decimal(15,3);default:0" json:"safety_stock"`
	Location    string    `gorm:"size:100" json:"location"`
	UnitCost    float64   `gorm:"type:decimal(15,2);default:0" json:"unit_cost"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for Stock
func (Stock) TableName() string {
	return "stocks"
}

// StockLog model maps to 'stock_logs' table (v2)
type StockLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	StockID        uint      `gorm:"index" json:"stock_id"`
	Type           string    `gorm:"size:10;not null" json:"type"` // 'in' or 'out'
	Quantity       float64   `gorm:"type:decimal(15,3);not null" json:"quantity"`
	QuantityBefore float64   `gorm:"type:decimal(15,3);default:0" json:"quantity_before"`
	QuantityAfter  float64   `gorm:"type:decimal(15,3);default:0" json:"quantity_after"`
	SourceType     string    `gorm:"size:20;not null" json:"source_type"` // inbound/requisition/adjust/transfer
	SourceID       *uint     `gorm:"index" json:"source_id"`
	SourceNo       string    `gorm:"size:50" json:"source_no"`
	ProjectID      uint      `gorm:"index;not null" json:"project_id"`
	MaterialID     uint      `gorm:"index;not null" json:"material_id"`
	UserID         *uint     `gorm:"index" json:"user_id"`
	Remark         string    `gorm:"size:500" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
}

// TableName specifies the table name for StockLog
func (StockLog) TableName() string {
	return "stock_logs"
}

// StockOpLog model maps to 'stock_op_logs' table (unchanged)
type StockOpLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	OpType    string    `gorm:"size:20" json:"op_type"`
	StockID   uint      `gorm:"index" json:"stock_id"`
	LogID     uint      `json:"log_id"`
	Detail    string    `gorm:"type:text" json:"detail"`
	Time      time.Time `json:"time"`
}

// ToDTO converts Stock to a data transfer object
func (s *Stock) ToDTO() map[string]any {
	return map[string]any{
		"id":           s.ID,
		"project_id":   s.ProjectID,
		"material_id":  s.MaterialID,
		"warehouse_id": s.WarehouseID,
		"quantity":     s.Quantity,
		"safety_stock": s.SafetyStock,
		"location":     s.Location,
		"unit_cost":    s.UnitCost,
		"created_at":   s.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":   s.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToDTOWithMaterial converts Stock to DTO with material information
func (s *Stock) ToDTOWithMaterial(materialCode, materialName, specification, unit string) map[string]any {
	dto := s.ToDTO()
	dto["material_code"] = materialCode
	dto["material_name"] = materialName
	dto["specification"] = specification
	dto["unit"] = unit
	return dto
}

func (sl *StockLog) ToDTO(userName string) map[string]any {
	if userName == "" {
		userName = "系统"
	}
	return map[string]any{
		"id":             sl.ID,
		"stock_id":       sl.StockID,
		"type":           sl.Type,
		"quantity":       sl.Quantity,
		"quantity_before": sl.QuantityBefore,
		"quantity_after":  sl.QuantityAfter,
		"source_type":    sl.SourceType,
		"source_id":      sl.SourceID,
		"source_no":      sl.SourceNo,
		"project_id":     sl.ProjectID,
		"material_id":    sl.MaterialID,
		"user_id":        sl.UserID,
		"user":           userName,
		"operator":       userName,
		"remark":         sl.Remark,
		"created_at":     sl.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
