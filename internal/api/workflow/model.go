package workflow

import (
	"time"
)

// WorkflowDefinition 工作流定义
type WorkflowDefinition struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;uniqueIndex" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Module      string    `gorm:"size:50;index" json:"module"` // 'inbound', 'requisition' 等
	Version     int       `gorm:"default:1" json:"version"`
	IsActive    bool      `gorm:"default:true;index" json:"is_active"`
	Nodes       []WorkflowNode `gorm:"foreignKey:WorkflowID;constraint:OnDelete:CASCADE" json:"nodes,omitempty"`
	Edges       []WorkflowEdge `gorm:"foreignKey:WorkflowID;constraint:OnDelete:CASCADE" json:"edges,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WorkflowNode 工作流节点
type WorkflowNode struct {
	ID            uint                `gorm:"primaryKey" json:"id"`
	WorkflowID    uint                `gorm:"not null;index" json:"workflow_id"`
	NodeKey       string              `gorm:"size:50" json:"node_key"` // 节点标识
	NodeType      string              `gorm:"size:20" json:"node_type"` // 'start', 'approval', 'end', 'parallel', 'merge'
	NodeName      string              `gorm:"size:100" json:"node_name"`
	Description   string              `gorm:"type:text" json:"description"`
	ApprovalType  string              `gorm:"size:20" json:"approval_type"` // 'sequential', 'parallel', 'any'
	TimeoutHours  int                 `json:"timeout_hours"`
	AutoApprove   bool                `gorm:"default:false" json:"auto_approve"`
	IsRequired    bool                `gorm:"default:true" json:"is_required"`
	X             int                 `gorm:"default:0" json:"x"` // 节点X坐标
	Y             int                 `gorm:"default:0" json:"y"` // 节点Y坐标
	Approvers     []WorkflowNodeApprover `gorm:"foreignKey:NodeID;constraint:OnDelete:CASCADE" json:"approvers,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
}

// WorkflowEdge 工作流边（连接线）
type WorkflowEdge struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	WorkflowID          uint      `gorm:"not null;index" json:"workflow_id"`
	FromNode            string    `gorm:"size:50" json:"from_node"` // 源节点 node_key
	ToNode              string    `gorm:"size:50" json:"to_node"`   // 目标节点 node_key
	ConditionExpression string    `gorm:"type:text" json:"condition_expression"` // 条件表达式
	CreatedAt           time.Time `json:"created_at"`
}

// WorkflowNodeApprover 工作流节点审批人配置
type WorkflowNodeApprover struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NodeID        uint      `gorm:"not null;index" json:"node_id"`
	ApproverType  string    `gorm:"size:20" json:"approver_type"` // 'user', 'role', 'department', 'superior'
	ApproverID    int       `json:"approver_id"`
	ApproverName  string    `gorm:"size:100" json:"approver_name"`
	Sequence      int       `gorm:"default:0" json:"sequence"` // 审批顺序
	CreatedAt     time.Time `json:"created_at"`
}

// WorkflowInstance 工作流实例
type WorkflowInstance struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	WorkflowID    uint      `gorm:"not null;index" json:"workflow_id"`
	BusinessType  string    `gorm:"size:50;index" json:"business_type"` // 'inbound_order', 'requisition' 等
	BusinessID    uint      `gorm:"index" json:"business_id"`
	BusinessNo    string    `gorm:"size:50" json:"business_no"`
	CurrentNode   string    `gorm:"size:50;index" json:"current_node"`
	Status        string    `gorm:"size:20;default:'pending';index" json:"status"` // 'pending', 'approved', 'rejected', 'cancelled'
	InitiatorID   uint      `gorm:"not null" json:"initiator_id"`
	InitiatorName string    `gorm:"size:100;not null" json:"initiator_name"`
	ProjectID     *uint     `gorm:"index" json:"project_id"` // 关联项目ID，用于项目角色审批
	StartedAt     time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联数据
	Workflow      *WorkflowDefinition `gorm:"foreignKey:WorkflowID" json:"-"`
	Approvals     []WorkflowApproval  `gorm:"foreignKey:InstanceID;constraint:OnDelete:CASCADE" json:"approvals,omitempty"`
	PendingTasks  []WorkflowPendingTask `gorm:"foreignKey:InstanceID;constraint:OnDelete:CASCADE" json:"pending_tasks,omitempty"`
}

// WorkflowApproval 工作流审批记录
type WorkflowApproval struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	InstanceID   uint      `gorm:"not null;index" json:"instance_id"`
	NodeID       uint      `gorm:"not null;index" json:"node_id"`
	NodeKey      string    `gorm:"size:50" json:"node_key"`
	ApproverID   uint      `gorm:"not null" json:"approver_id"`
	ApproverName string    `gorm:"size:100;not null" json:"approver_name"`
	Action       string    `gorm:"size:20" json:"action"` // 'approve', 'reject', 'return', 'comment'
	Remark       string    `gorm:"type:text" json:"remark"`
	Attachments  string    `gorm:"type:text" json:"attachments"` // JSON格式
	ApprovedAt   time.Time `json:"approved_at"`
	CreatedAt    time.Time `json:"created_at"`

	Instance     *WorkflowInstance `gorm:"foreignKey:InstanceID" json:"-"`
	Node         *WorkflowNode     `gorm:"foreignKey:NodeID" json:"-"`
}

// WorkflowPendingTask 工作流待办任务
type WorkflowPendingTask struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	InstanceID   uint       `gorm:"not null;index" json:"instance_id"`
	NodeID       uint       `gorm:"not null" json:"node_id"`
	NodeKey      string     `gorm:"size:50" json:"node_key"`
	NodeName     string     `gorm:"size:100" json:"node_name"`
	BusinessType string     `gorm:"size:50;index" json:"business_type"`
	BusinessID   uint       `gorm:"index" json:"business_id"`
	BusinessNo   string     `gorm:"size:50" json:"business_no"`
	ApproverID   uint       `gorm:"not null;index" json:"approver_id"`
	ApproverName string     `gorm:"size:100;not null" json:"approver_name"`
	Status       string     `gorm:"size:20;default:'pending';index" json:"status"` // 'pending', 'approved', 'rejected', 'returned', 'cancelled'
	IsParallel   bool       `gorm:"default:false" json:"is_parallel"`
	ArrivedAt    time.Time  `json:"arrived_at"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Instance     *WorkflowInstance `gorm:"foreignKey:InstanceID" json:"-"`
	Node         *WorkflowNode     `gorm:"foreignKey:NodeID" json:"-"`
}

// WorkflowLog 工作流操作日志
type WorkflowLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	InstanceID uint      `gorm:"not null;index" json:"instance_id"`
	NodeKey    string    `gorm:"size:50" json:"node_key"`
	Action     string    `gorm:"size:50" json:"action"` // 'start', 'approve', 'reject', 'return', 'cancel', 'comment'
	ActorID    uint      `gorm:"not null" json:"actor_id"`
	ActorName  string    `gorm:"size:100;not null" json:"actor_name"`
	ActionData string    `gorm:"type:text" json:"action_data"` // JSON格式
	CreatedAt  time.Time `json:"created_at"`

	Instance   *WorkflowInstance `gorm:"foreignKey:InstanceID" json:"-"`
}

// 审批动作常量
const (
	ActionApprove = "approve" // 通过
	ActionReject  = "reject"  // 拒绝
	ActionReturn  = "return"  // 退回
	ActionComment = "comment" // 评论
	ActionStart   = "start"   // 开始
	ActionCancel  = "cancel"  // 取消
)

// 节点类型常量
const (
	NodeTypeStart    = "start"
	NodeTypeApproval = "approval"
	NodeTypeEnd      = "end"
	NodeTypeParallel = "parallel"
	NodeTypeMerge    = "merge"
)

// 审批类型常量
const (
	ApprovalTypeSequential = "sequential" // 顺序审批（所有人都需要审批）
	ApprovalTypeParallel   = "parallel"   // 并行审批（同时审批，所有人都需要通过）
	ApprovalTypeAny        = "any"        // 任一审批（任意一人通过即可）
)

// 实例状态常量
const (
	InstanceStatusPending   = "pending"   // 进行中
	InstanceStatusApproved  = "approved"  // 已通过
	InstanceStatusRejected  = "rejected"  // 已拒绝
	InstanceStatusCancelled = "cancelled" // 已取消
)

// 待办状态常量
const (
	TaskStatusPending   = "pending"   // 待处理
	TaskStatusApproved  = "approved"  // 已通过
	TaskStatusRejected  = "rejected"  // 已拒绝
	TaskStatusReturned  = "returned"  // 已退回
	TaskStatusCancelled = "cancelled" // 已取消
)

// 审批人类型常量
const (
	ApproverTypeUser       = "user"       // 指定用户
	ApproverTypeRole       = "role"       // 角色
	ApproverTypeDepartment = "department" // 部门
	ApproverTypeSuperior   = "superior"   // 上级
	ApproverTypeProjectRole = "project_role" // 项目+角色组合（从项目关联用户中筛选特定角色）
)

// ToDTO 转换为DTO格式
func (w *WorkflowDefinition) ToDTO() map[string]any {
	nodes := make([]map[string]any, 0)
	if w.Nodes != nil {
		for _, node := range w.Nodes {
			nodes = append(nodes, node.ToDTO())
		}
	}

	edges := make([]map[string]any, 0)
	if w.Edges != nil {
		for _, edge := range w.Edges {
			edges = append(edges, edge.ToDTO())
		}
	}

	return map[string]any{
		"id":          w.ID,
		"name":        w.Name,
		"description": w.Description,
		"module":      w.Module,
		"version":     w.Version,
		"is_active":   w.IsActive,
		"nodes":       nodes,
		"edges":       edges,
		"created_at":  w.CreatedAt,
		"updated_at":  w.UpdatedAt,
	}
}

func (n *WorkflowNode) ToDTO() map[string]any {
	approvers := make([]map[string]any, 0)
	if n.Approvers != nil {
		for _, approver := range n.Approvers {
			approvers = append(approvers, approver.ToDTO())
		}
	}

	return map[string]any{
		"id":            n.ID,
		"workflow_id":   n.WorkflowID,
		"node_key":      n.NodeKey,
		"node_type":     n.NodeType,
		"node_name":     n.NodeName,
		"description":   n.Description,
		"approval_type": n.ApprovalType,
		"timeout_hours": n.TimeoutHours,
		"auto_approve":  n.AutoApprove,
		"is_required":   n.IsRequired,
		"x":             n.X,
		"y":             n.Y,
		"approvers":     approvers,
		"created_at":    n.CreatedAt,
	}
}

func (n *WorkflowNodeApprover) ToDTO() map[string]any {
	return map[string]any{
		"id":            n.ID,
		"node_id":       n.NodeID,
		"approver_type": n.ApproverType,
		"approver_id":   n.ApproverID,
		"approver_name": n.ApproverName,
		"sequence":      n.Sequence,
		"created_at":    n.CreatedAt,
	}
}

func (e *WorkflowEdge) ToDTO() map[string]any {
	return map[string]any{
		"id":                   e.ID,
		"workflow_id":          e.WorkflowID,
		"from_node":            e.FromNode,
		"to_node":              e.ToNode,
		"condition_expression": e.ConditionExpression,
		"created_at":           e.CreatedAt,
	}
}

func (i *WorkflowInstance) ToDTO() map[string]any {
	return map[string]any{
		"id":            i.ID,
		"workflow_id":   i.WorkflowID,
		"business_type": i.BusinessType,
		"business_id":   i.BusinessID,
		"business_no":   i.BusinessNo,
		"current_node":  i.CurrentNode,
		"status":        i.Status,
		"initiator_id":  i.InitiatorID,
		"initiator_name": i.InitiatorName,
		"started_at":    i.StartedAt,
		"finished_at":   i.FinishedAt,
		"created_at":    i.CreatedAt,
		"updated_at":    i.UpdatedAt,
	}
}

func (a *WorkflowApproval) ToDTO() map[string]any {
	return map[string]any{
		"id":            a.ID,
		"instance_id":   a.InstanceID,
		"node_id":       a.NodeID,
		"node_key":      a.NodeKey,
		"approver_id":   a.ApproverID,
		"approver_name": a.ApproverName,
		"action":        a.Action,
		"remark":        a.Remark,
		"attachments":   a.Attachments,
		"approved_at":   a.ApprovedAt,
		"created_at":    a.CreatedAt,
	}
}

func (t *WorkflowPendingTask) ToDTO() map[string]any {
	return map[string]any{
		"id":            t.ID,
		"instance_id":   t.InstanceID,
		"node_id":       t.NodeID,
		"node_key":      t.NodeKey,
		"node_name":     t.NodeName,
		"business_type": t.BusinessType,
		"business_id":   t.BusinessID,
		"business_no":   t.BusinessNo,
		"approver_id":   t.ApproverID,
		"approver_name": t.ApproverName,
		"status":        t.Status,
		"is_parallel":   t.IsParallel,
		"arrived_at":    t.ArrivedAt,
		"processed_at":  t.ProcessedAt,
		"created_at":    t.CreatedAt,
		"updated_at":    t.UpdatedAt,
	}
}

func (l *WorkflowLog) ToDTO() map[string]any {
	return map[string]any{
		"id":          l.ID,
		"instance_id": l.InstanceID,
		"node_key":    l.NodeKey,
		"action":      l.Action,
		"actor_id":    l.ActorID,
		"actor_name":  l.ActorName,
		"action_data": l.ActionData,
		"created_at":  l.CreatedAt,
	}
}
