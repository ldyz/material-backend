package progress

import "time"

// Resource 资源表
type Resource struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"not null;index;comment:项目ID" json:"project_id"`
	Name        string    `gorm:"type:varchar(100);not null;comment:资源名称" json:"name"`
	Type        string    `gorm:"type:varchar(20);not null;comment:资源类型" json:"type"`     // labor, equipment, material
	Unit        string    `gorm:"type:varchar(20);comment:计量单位" json:"unit"`             // 人/d, 台/d, kg
	Quantity    float64   `gorm:"default:0;comment:可用数量" json:"quantity"`
	CostPerUnit float64   `gorm:"default:0;comment:单位成本" json:"cost_per_unit"`
	Color       string    `gorm:"type:varchar(20);comment:显示颜色" json:"color"`
	IsActive    bool      `gorm:"default:true;comment:是否启用" json:"is_active"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (Resource) TableName() string {
	return "resources"
}

// TaskResource 任务资源分配表
type TaskResource struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `gorm:"not null;index;comment:任务ID" json:"task_id"`
	ResourceID uint      `gorm:"not null;index;comment:资源ID" json:"resource_id"`
	Quantity   float64   `gorm:"default:0;comment:分配数量" json:"quantity"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (TaskResource) TableName() string {
	return "task_resources"
}
