package material

import (
	"time"
)

// MaterialCategory 物资分类模型
type MaterialCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ParentID  uint      `gorm:"default:0;index" json:"parent_id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Code      string    `gorm:"size:20;index" json:"code"`
	Level     int       `gorm:"default:1" json:"level"`
	Path      string    `gorm:"size:255" json:"path"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Remark    string    `gorm:"type:text" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Children  []MaterialCategory `gorm:"-" json:"children,omitempty"`
}

// TableName 指定表名
func (MaterialCategory) TableName() string {
	return "material_categories"
}

// ToDTO 转换为DTO格式
func (mc *MaterialCategory) ToDTO() map[string]any {
	return map[string]any{
		"id":         mc.ID,
		"parent_id":  mc.ParentID,
		"name":       mc.Name,
		"code":       mc.Code,
		"level":      mc.Level,
		"path":       mc.Path,
		"sort":       mc.Sort,
		"remark":     mc.Remark,
		"created_at": mc.CreatedAt,
		"updated_at": mc.UpdatedAt,
	}
}

// 包导入用于gorm标签和time.Time类型
