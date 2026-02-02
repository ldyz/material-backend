package construction_log

import (
	"time"
)

// ConstructionLog 施工日志模型
type ConstructionLog struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"size:128;not null;comment:标题" json:"title"`
	Content     string     `gorm:"type:text;comment:图文内容" json:"content"`
	Images      string     `gorm:"type:text;comment:图片URL，逗号分隔或JSON" json:"images"`
	Weather     string     `gorm:"size:64;comment:天气" json:"weather"`
	Temperature float64    `gorm:"comment:温度" json:"temperature"`
	Progress    string     `gorm:"type:text;comment:施工进度" json:"progress"`
	Issues      string     `gorm:"type:text;comment:存在问题" json:"issues"`
	LogDate     string     `gorm:"comment:日志日期" json:"log_date"`
	Remark      string     `gorm:"type:text;comment:备注" json:"remark"`
	ProjectID   uint       `gorm:"not null;comment:项目ID" json:"project_id"`
	CreatorID   uint       `gorm:"not null;comment:创建者ID" json:"creator_id"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}

// TableName 设置表名
func (ConstructionLog) TableName() string {
	return "construction_log"
}

// ToDict 将施工日志转换为字典格式
func (log *ConstructionLog) ToDict() map[string]any {
	return map[string]any{
		"id":         log.ID,
		"title":      log.Title,
		"content":    log.Content,
		"images":     log.Images,
		"weather":    log.Weather,
		"project_id": log.ProjectID,
		"creator_id": log.CreatorID,
		"created_at": log.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": log.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}