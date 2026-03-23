package material_master

import (
	"time"
)

// MaterialMaster 物资主数据模型（全局唯一）
type MaterialMaster struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Code         string    `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Name         string    `gorm:"size:200;not null;index" json:"name"`
	Specification string   `gorm:"size:200" json:"specification"`
	Unit         string    `gorm:"size:20" json:"unit"`
	Material     string    `gorm:"size:100" json:"material"`
	Category     string    `gorm:"size:100;index" json:"category"`
	SafetyStock  float64   `gorm:"type:decimal(15,3);default:0" json:"safety_stock"`
	Description  string    `gorm:"type:text" json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (MaterialMaster) TableName() string {
	return "material_master"
}

// ToDTO 转换为 DTO
func (m *MaterialMaster) ToDTO() map[string]any {
	return map[string]any{
		"id":            m.ID,
		"code":          m.Code,
		"name":          m.Name,
		"specification": m.Specification,
		"unit":          m.Unit,
		"material":      m.Material,
		"category":      m.Category,
		"safety_stock":  m.SafetyStock,
		"description":   m.Description,
		"created_at":    m.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":    m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// CreateMaterialMasterRequest 创建物资主数据请求
type CreateMaterialMasterRequest struct {
	Code         string  `json:"code" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Specification string `json:"specification"`
	Unit         string  `json:"unit"`
	Material     string  `json:"material"`
	Category     string  `json:"category"`
	SafetyStock  float64 `json:"safety_stock"`
	Description  string  `json:"description"`
}

// UpdateMaterialMasterRequest 更新物资主数据请求
type UpdateMaterialMasterRequest struct {
	Code         string  `json:"code" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Specification string `json:"specification"`
	Unit         string  `json:"unit"`
	Material     string  `json:"material"`
	Category     string  `json:"category"`
	SafetyStock  float64 `json:"safety_stock"`
	Description  string  `json:"description"`
}

// MaterialMasterQueryDTO 物资主数据查询响应（带库存信息）
type MaterialMasterQueryDTO struct {
	ID           uint    `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Specification string `json:"specification"`
	Unit         string  `json:"unit"`
	Material     string  `json:"material"`
	Category     string  `json:"category"`
	SafetyStock  float64 `json:"safety_stock"`
	Description  string  `json:"description"`

	// 项目库存信息（如果指定了项目ID）
	ProjectID    *uint   `json:"project_id,omitempty"`
	StockID      *uint   `json:"stock_id,omitempty"`
	Quantity     float64 `json:"quantity,omitempty"`
	Location     string  `json:"location,omitempty"`
	UnitCost     float64 `json:"unit_cost,omitempty"`

	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// ToQueryDTO 转换为查询 DTO（带库存信息）
func (m *MaterialMaster) ToQueryDTO() MaterialMasterQueryDTO {
	return MaterialMasterQueryDTO{
		ID:           m.ID,
		Code:         m.Code,
		Name:         m.Name,
		Specification: m.Specification,
		Unit:         m.Unit,
		Material:     m.Material,
		Category:     m.Category,
		SafetyStock:  m.SafetyStock,
		Description:  m.Description,
		CreatedAt:    m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// MaterialMasterWithStock 物资主数据与库存关联结构
type MaterialMasterWithStock struct {
	MaterialMaster
	ProjectID *uint   `json:"project_id,omitempty"`
	StockID   *uint   `json:"stock_id,omitempty"`
	Quantity  float64 `json:"quantity,omitempty"`
	Location  string  `json:"location,omitempty"`
	UnitCost  float64 `json:"unit_cost,omitempty"`
}
