package audit

import (
	"encoding/json"
	"time"
)

// OperationLog 通用操作日志模型
type OperationLog struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	UserID        *uint           `gorm:"index" json:"user_id"`
	Username      string          `gorm:"size:100" json:"username"`
	Operation     string          `gorm:"size:100;not null;index" json:"operation"`     // 操作类型
	Module        string          `gorm:"size:50;not null;index" json:"module"`          // 模块
	ResourceType  string          `gorm:"size:100" json:"resource_type"`                  // 资源类型
	ResourceID    *uint           `json:"resource_id"`                                    // 资源ID
	ResourceNo    string          `gorm:"size:100;index" json:"resource_no"`              // 资源编号
	Changes       json.RawMessage `gorm:"type:jsonb" json:"changes,omitempty"`          // 变更内容
	RequestMethod string         `gorm:"size:10" json:"request_method"`                  // 请求方法
	RequestPath   string          `gorm:"size:500" json:"request_path"`                   // 请求路径
	RequestParams json.RawMessage `gorm:"type:jsonb" json:"request_params,omitempty"`    // 请求参数
	IPAddress     string          `gorm:"size:45" json:"ip_address,omitempty"`            // IP地址
	UserAgent     string          `gorm:"type:text" json:"user_agent,omitempty"`        // 用户代理
	Status        string          `gorm:"size:20;not null;index" json:"status"`           // 状态
	ErrorMessage string          `gorm:"type:text" json:"error_message,omitempty"`      // 错误信息
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// OperationType 操作类型常量
const (
	OpCreate   = "create"
	OpUpdate   = "update"
	OpDelete   = "delete"
	OpApprove  = "approve"
	OpReject   = "reject"
	OpSubmit   = "submit"
	OpCancel   = "cancel"
	OpActivate = "activate"
	OpComplete = "complete"
	OpIssue    = "issue"    // 发放
	OpReceive  = "receive"  // 接收
	OpAdjust   = "adjust"   // 调整
	OpTransfer = "transfer"
	OpAssign   = "assign"   // 分配
	OpStart    = "start"    // 开始
	OpLogin    = "login"
	OpLogout   = "logout"
)

// Module 模块常量
const (
	ModuleMaterialPlan = "material_plan"
	ModuleInbound      = "inbound"
	ModuleOutbound     = "outbound"
	ModuleRequisition  = "requisition"
	ModuleStock        = "stock"
	ModuleMaterial     = "material"
	ModuleWorkflow     = "workflow"
	ModuleProject      = "project"
	ModuleConstruction = "construction"
	ModuleAppointment  = "appointment"
	ModuleSystem       = "system"
	ModuleAuth         = "auth"
)

// ResourceType 资源类型常量
const (
	ResourceMaterialPlan     = "MaterialPlan"
	ResourceInboundOrder     = "InboundOrder"
	ResourceRequisition      = "Requisition"
	ResourceStock            = "Stock"
	ResourceMaterial         = "Material"
	ResourceWorkflowTask     = "WorkflowTask"
	ResourceWorkflowInstance = "WorkflowInstance"
	ResourceAppointment      = "ConstructionAppointment"
)

// LogStatus 日志状态常量
const (
	LogStatusSuccess = "success"
	LogStatusError   = "error"
)

// LogChanges 记录变更内容
type LogChanges struct {
	Before interface{} `json:"before,omitempty"`
	After  interface{} `json:"after,omitempty"`
	Diff   []FieldDiff `json:"diff,omitempty"`
}

// FieldDiff 字段差异
type FieldDiff struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"old_value,omitempty"`
	NewValue interface{} `json:"new_value,omitempty"`
}

// ToDTO 转换为DTO
func (ol *OperationLog) ToDTO() map[string]any {
	return map[string]any{
		"id":             ol.ID,
		"user_id":        ol.UserID,
		"username":       ol.Username,
		"operation":      ol.Operation,
		"module":         ol.Module,
		"resource_type":  ol.ResourceType,
		"resource_id":    ol.ResourceID,
		"resource_no":    ol.ResourceNo,
		"request_method": ol.RequestMethod,
		"request_path":   ol.RequestPath,
		"ip_address":     ol.IPAddress,
		"status":         ol.Status,
		"error_message":  ol.ErrorMessage,
		"created_at":     ol.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToDetailDTO 转换为详细DTO（包含变更内容）
func (ol *OperationLog) ToDetailDTO() map[string]any {
	dto := ol.ToDTO()
	if ol.Changes != nil {
		var changes LogChanges
		_ = json.Unmarshal(ol.Changes, &changes)
		dto["changes"] = changes
	}
	if ol.RequestParams != nil {
		var params map[string]any
		_ = json.Unmarshal(ol.RequestParams, &params)
		dto["request_params"] = params
	}
	return dto
}

// OperationLogFilter 查询过滤器
type OperationLogFilter struct {
	UserID      *uint      `json:"user_id"`
	Operation   string     `json:"operation"`
	Module      string     `json:"module"`
	ResourceType string    `json:"resource_type"`
	ResourceNo  string     `json:"resource_no"`
	Status      string     `json:"status"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Keyword     string     `json:"keyword"`     // 搜索关键词（resource_no, username等）
	Page        int        `json:"page"`
	PageSize    int        `json:"page_size"`
}

// OperationLogListResponse 操作日志列表响应
type OperationLogListResponse struct {
	Data       []OperationLog     `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}
