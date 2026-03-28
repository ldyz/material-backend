package system

import (
	"time"
)

// SystemConfig 系统配置表
type SystemConfig struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Key         string    `gorm:"size:100;unique;not null;comment:配置键" json:"key"`
	Value       string    `gorm:"type:text;comment:配置值" json:"value"`
	Type        string    `gorm:"size:20;default:string;comment:配置类型:string,integer,float,boolean" json:"type"`
	Description string    `gorm:"size:200;comment:配置描述" json:"description"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}

// TableName 设置表名
func (SystemConfig) TableName() string {
	return "system_config"
}

// SystemBackup 系统备份记录表
type SystemBackup struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Filename    string    `gorm:"size:200;not null;comment:备份文件名" json:"filename"`
	Filepath    string    `gorm:"size:500;not null;comment:备份文件路径" json:"filepath"`
	Size        int64     `gorm:"default:0;comment:备份文件大小(字节)" json:"size"`
	Status      string    `gorm:"size:20;default:completed;comment:备份状态:pending,completed,failed" json:"status"`
	CreatedBy   string    `gorm:"size:50;comment:创建者" json:"created_by"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	Description string    `gorm:"type:text;comment:备份描述" json:"description"`
}

// TableName 设置表名
func (SystemBackup) TableName() string {
	return "system_backup"
}

// SystemActivity 系统活动记录表
type SystemActivity struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityType  string    `gorm:"size:50;not null;comment:活动类型" json:"activity_type"`
	Title         string    `gorm:"size:200;not null;comment:活动标题" json:"title"`
	Description   string    `gorm:"type:text;comment:活动描述" json:"description"`
	UserID        uint      `gorm:"comment:用户ID" json:"user_id"`
	Username      string    `gorm:"size:50;comment:用户名" json:"username"`
	IPAddress     string    `gorm:"size:45;comment:IP地址" json:"ip_address"`
	UserAgent     string    `gorm:"size:500;comment:用户代理" json:"user_agent"`
	Status        string    `gorm:"size:20;default:success;comment:状态:success,failed,warning" json:"status"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	ExtraData     string    `gorm:"type:text;comment:额外数据(JSON格式)" json:"extra_data"`
}

// TableName 设置表名
func (SystemActivity) TableName() string {
	return "system_activity"
}

// SystemLog 系统日志表
type SystemLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Level     string    `gorm:"size:20;not null" json:"level"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Module    string    `gorm:"size:100" json:"module"`
	UserID    uint      `gorm:"comment:用户ID" json:"user_id"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 设置表名
func (SystemLog) TableName() string {
	return "system_logs"
}

// AIAnalysisLog AI分析日志表
type AIAnalysisLog struct {
	ID             uint      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Question       string    `gorm:"type:text;not null;comment:用户问题" json:"question"`
	Answer         string    `gorm:"type:text;comment:AI回答" json:"answer"`
	QueryUsed      string    `gorm:"type:text;comment:使用的SQL查询" json:"query_used"`
	DataSummary    string    `gorm:"type:text;comment:数据摘要" json:"data_summary"`
	ProcessingTime float64   `gorm:"comment:处理时间(秒)" json:"processing_time"`
	UserID         uint      `gorm:"comment:用户ID" json:"user_id"`
	IPAddress      string    `gorm:"size:45;comment:用户IP地址" json:"ip_address"`
	Status         string    `gorm:"size:20;default:completed;comment:状态:completed,failed,processing" json:"status"`
	ErrorMessage   string    `gorm:"type:text;comment:错误信息" json:"error_message"`
	ModelUsed      string    `gorm:"size:50;default:deepseek-chat;comment:使用的AI模型" json:"model_used"`
	TokensUsed     int       `gorm:"comment:使用的token数量" json:"tokens_used"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
}

// TableName 设置表名
func (AIAnalysisLog) TableName() string {
	return "ai_analysis_logs"
}

// ToDTO 将SystemBackup转换为DTO
func (sb *SystemBackup) ToDTO() map[string]any {
	return map[string]any{
		"id":          sb.ID,
		"filename":    sb.Filename,
		"filepath":    sb.Filepath,
		"size":        sb.Size,
		"status":      sb.Status,
		"created_by":  sb.CreatedBy,
		"created_at":  sb.CreatedAt.Format(time.RFC3339),
		"description": sb.Description,
	}
}

// ToDTO 将SystemActivity转换为DTO
func (sa *SystemActivity) ToDTO() map[string]any {
	return map[string]any{
		"id":          sa.ID,
		"activity_type": sa.ActivityType,
		"title":       sa.Title,
		"description": sa.Description,
		"user_id":     sa.UserID,
		"username":    sa.Username,
		"ip_address":  sa.IPAddress,
		"user_agent":  sa.UserAgent,
		"status":      sa.Status,
		"created_at":  sa.CreatedAt.Format(time.RFC3339),
		"extra_data":  sa.ExtraData,
	}
}

// ToDTO 将SystemConfig转换为DTO
func (sc *SystemConfig) ToDTO() map[string]any {
	return map[string]any{
		"id":          sc.ID,
		"key":         sc.Key,
		"value":       sc.Value,
		"type":        sc.Type,
		"description": sc.Description,
		"created_at":  sc.CreatedAt.Format(time.RFC3339),
		"updated_at":  sc.UpdatedAt.Format(time.RFC3339),
	}
}