package workflow

import (
	"errors"
	"fmt"
	"time"

	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/notification"
	"gorm.io/gorm"
)

// Engine 工作流引擎
type Engine struct {
	db *gorm.DB
}

// NewEngine 创建工作流引擎
func NewEngine(db *gorm.DB) *Engine {
	return &Engine{db: db}
}

// StartWorkflow 启动工作流
func (e *Engine) StartWorkflow(workflowID uint, businessType string, businessID uint, businessNo string, initiatorID uint, initiatorName string) (*WorkflowInstance, error) {
	// 1. 获取工作流定义
	var workflow WorkflowDefinition
	if err := e.db.Preload("Nodes").Preload("Nodes.Approvers").Preload("Edges").First(&workflow, workflowID).Error; err != nil {
		return nil, fmt.Errorf("工作流不存在: %w", err)
	}

	if !workflow.IsActive {
		return nil, errors.New("工作流未激活")
	}

	// 2. 创建工作流实例
	instance := &WorkflowInstance{
		WorkflowID:    workflowID,
		BusinessType:  businessType,
		BusinessID:    businessID,
		BusinessNo:    businessNo,
		CurrentNode:   NodeTypeStart,
		Status:        InstanceStatusPending,
		InitiatorID:   initiatorID,
		InitiatorName: initiatorName,
		StartedAt:     time.Now(),
	}

	if err := e.db.Create(instance).Error; err != nil {
		return nil, fmt.Errorf("创建工作流实例失败: %w", err)
	}

	// 3. 记录日志
	e.log(instance.ID, NodeTypeStart, ActionStart, initiatorID, initiatorName, nil)

	// 4. 获取开始节点的下一个节点
	nextNode, err := e.getNextNode(&workflow, NodeTypeStart)
	if err != nil {
		return nil, err
	}

	// 5. 移动到下一个节点
	if err := e.moveToNode(instance, &workflow, nextNode, initiatorID, initiatorName); err != nil {
		return nil, err
	}

	return instance, nil
}

// ProcessApproval 处理审批
func (e *Engine) ProcessApproval(instanceID uint, approverID uint, approverName string, action string, remark string, extraData map[string]any) error {
	// 1. 获取工作流实例
	var instance WorkflowInstance
	if err := e.db.Preload("Workflow").Preload("Workflow.Nodes").Preload("Workflow.Edges").First(&instance, instanceID).Error; err != nil {
		return fmt.Errorf("工作流实例不存在: %w", err)
	}

	if instance.Status != InstanceStatusPending {
		return errors.New("工作流已结束")
	}

	// 2. 获取当前待办任务
	var task WorkflowPendingTask
	if err := e.db.Where("instance_id = ? AND approver_id = ? AND status = ?", instanceID, approverID, TaskStatusPending).First(&task).Error; err != nil {
		// 检查是否有其他待办任务（说明该用户没有审批权限）
		var taskCount int64
		e.db.Model(&WorkflowPendingTask{}).Where("instance_id = ? AND status = ?", instanceID, TaskStatusPending).Count(&taskCount)
		if taskCount > 0 {
			return errors.New("您没有审批权限，该任务的审批人是其他角色")
		}
		// 没有待办任务，可能已处理
		return errors.New("该任务已处理或不存在")
	}

	// 3. 获取节点信息
	var node WorkflowNode
	if err := e.db.Preload("Approvers").First(&node, task.NodeID).Error; err != nil {
		return fmt.Errorf("节点不存在: %w", err)
	}

	// 4. 根据操作类型处理
	switch action {
	case ActionApprove:
		return e.processApprove(&instance, &node, &task, instance.Workflow, approverID, approverName, remark, extraData)
	case ActionReject:
		return e.processReject(&instance, &node, approverID, approverName, remark)
	case ActionReturn:
		return e.processReturn(&instance, &node, approverID, approverName, remark)
	case ActionComment:
		return e.processComment(&instance, &node, approverID, approverName, remark)
	default:
		return errors.New("无效的操作类型")
	}
}

// processApprove 处理通过
func (e *Engine) processApprove(instance *WorkflowInstance, node *WorkflowNode, task *WorkflowPendingTask, workflow *WorkflowDefinition, approverID uint, approverName string, remark string, extraData map[string]any) error {
	// 1. 更新待办任务状态
	now := time.Now()
	task.Status = TaskStatusApproved
	task.ProcessedAt = &now
	if err := e.db.Save(task).Error; err != nil {
		return err
	}

	// 2. 创建审批记录
	approval := &WorkflowApproval{
		InstanceID:   instance.ID,
		NodeID:       node.ID,
		NodeKey:      node.NodeKey,
		ApproverID:   approverID,
		ApproverName: approverName,
		Action:       ActionApprove,
		Remark:       remark,
		ApprovedAt:   now,
	}
	if err := e.db.Create(approval).Error; err != nil {
		return err
	}

	// 3. 记录日志
	e.log(instance.ID, node.NodeKey, ActionApprove, approverID, approverName, map[string]any{
		"remark": remark,
	})

	// 4. 检查是否所有人都已审批（根据审批类型）
	allApproved, err := e.checkAllApproved(instance, node)
	if err != nil {
		return err
	}

	if !allApproved {
		// 还需要其他人审批，不推进流程
		return nil
	}

	// 5. 所有人都已审批，移动到下一个节点
	nextNode, err := e.getNextNode(workflow, node.NodeKey)
	if err != nil {
		return err
	}

	return e.moveToNode(instance, workflow, nextNode, approverID, approverName)
}

// processReject 处理拒绝
func (e *Engine) processReject(instance *WorkflowInstance, node *WorkflowNode, approverID uint, approverName string, remark string) error {
	// 1. 更新实例状态为已拒绝
	instance.Status = InstanceStatusRejected
	now := time.Now()
	instance.FinishedAt = &now
	if err := e.db.Save(instance).Error; err != nil {
		return err
	}

	// 2. 取消所有待办任务
	e.db.Model(&WorkflowPendingTask{}).Where("instance_id = ? AND status = ?", instance.ID, TaskStatusPending).Updates(map[string]any{
		"status":        TaskStatusCancelled,
		"processed_at":  now,
	})

	// 3. 创建审批记录
	approval := &WorkflowApproval{
		InstanceID:   instance.ID,
		NodeID:       node.ID,
		NodeKey:      node.NodeKey,
		ApproverID:   approverID,
		ApproverName: approverName,
		Action:       ActionReject,
		Remark:       remark,
		ApprovedAt:   now,
	}
	if err := e.db.Create(approval).Error; err != nil {
		return err
	}

	// 4. 记录日志
	e.log(instance.ID, node.NodeKey, ActionReject, approverID, approverName, map[string]any{
		"remark": remark,
	})

	return nil
}

// processReturn 处理退回
func (e *Engine) processReturn(instance *WorkflowInstance, node *WorkflowNode, approverID uint, approverName string, remark string) error {
	// 1. 取消当前节点的所有待办任务
	now := time.Now()
	e.db.Model(&WorkflowPendingTask{}).Where("instance_id = ? AND node_id = ? AND status = ?", instance.ID, node.ID, TaskStatusPending).Updates(map[string]any{
		"status":       TaskStatusReturned,
		"processed_at": now,
	})

	// 2. 创建审批记录
	approval := &WorkflowApproval{
		InstanceID:   instance.ID,
		NodeID:       node.ID,
		NodeKey:      node.NodeKey,
		ApproverID:   approverID,
		ApproverName: approverName,
		Action:       ActionReturn,
		Remark:       remark,
		ApprovedAt:   now,
	}
	if err := e.db.Create(approval).Error; err != nil {
		return err
	}

	// 3. 记录日志
	e.log(instance.ID, node.NodeKey, ActionReturn, approverID, approverName, map[string]any{
		"remark": remark,
	})

	// 4. 回退到开始节点（申请人重新提交）
	instance.CurrentNode = NodeTypeStart
	if err := e.db.Save(instance).Error; err != nil {
		return err
	}

	// 5. 为申请人创建待办任务
	task := &WorkflowPendingTask{
		InstanceID:   instance.ID,
		NodeID:       node.ID, // 使用当前节点ID记录
		NodeKey:      NodeTypeStart,
		NodeName:     "重新提交",
		BusinessType: instance.BusinessType,
		BusinessID:   instance.BusinessID,
		BusinessNo:   instance.BusinessNo,
		ApproverID:   instance.InitiatorID,
		ApproverName: instance.InitiatorName,
		Status:       TaskStatusPending,
		ArrivedAt:    now,
	}
	if err := e.db.Create(task).Error; err != nil {
		return err
	}

	return nil
}

// processComment 处理评论
func (e *Engine) processComment(instance *WorkflowInstance, node *WorkflowNode, approverID uint, approverName string, remark string) error {
	// 1. 创建审批记录（仅评论，不推进流程）
	approval := &WorkflowApproval{
		InstanceID:   instance.ID,
		NodeID:       node.ID,
		NodeKey:      node.NodeKey,
		ApproverID:   approverID,
		ApproverName: approverName,
		Action:       ActionComment,
		Remark:       remark,
		ApprovedAt:   time.Now(),
	}
	if err := e.db.Create(approval).Error; err != nil {
		return err
	}

	// 2. 记录日志
	e.log(instance.ID, node.NodeKey, ActionComment, approverID, approverName, map[string]any{
		"remark": remark,
	})

	return nil
}

// Resubmit 重新提交（退回后）
func (e *Engine) Resubmit(instanceID uint, submitterID uint, submitterName string) error {
	// 1. 获取工作流实例
	var instance WorkflowInstance
	if err := e.db.Preload("Workflow").Preload("Workflow.Nodes").Preload("Workflow.Edges").First(&instance, instanceID).Error; err != nil {
		return fmt.Errorf("工作流实例不存在: %w", err)
	}

	if instance.CurrentNode != NodeTypeStart {
		return errors.New("当前状态不允许重新提交")
	}

	// 2. 清除申请人的待办任务
	e.db.Model(&WorkflowPendingTask{}).Where("instance_id = ? AND node_key = ? AND approver_id = ?", instanceID, NodeTypeStart, submitterID).Updates(map[string]any{
		"status":       TaskStatusCancelled,
		"processed_at": time.Now(),
	})

	// 3. 记录日志
	e.log(instance.ID, NodeTypeStart, "resubmit", submitterID, submitterName, nil)

	// 4. 获取开始节点的下一个节点
	nextNode, err := e.getNextNode(instance.Workflow, NodeTypeStart)
	if err != nil {
		return err
	}

	// 5. 移动到下一个节点
	return e.moveToNode(&instance, instance.Workflow, nextNode, submitterID, submitterName)
}

// moveToNode 移动到指定节点
func (e *Engine) moveToNode(instance *WorkflowInstance, workflow *WorkflowDefinition, node *WorkflowNode, actorID uint, actorName string) error {
	// 更新当前节点
	instance.CurrentNode = node.NodeKey
	if err := e.db.Save(instance).Error; err != nil {
		return err
	}

	// 如果是结束节点
	if node.NodeType == NodeTypeEnd {
		return e.finishWorkflow(instance, actorID, actorName)
	}

	// 如果是审批节点，创建待办任务
	if node.NodeType == NodeTypeApproval {
		return e.createApprovalTasks(instance, node)
	}

	return nil
}

// finishWorkflow 完成工作流
func (e *Engine) finishWorkflow(instance *WorkflowInstance, actorID uint, actorName string) error {
	now := time.Now()
	instance.Status = InstanceStatusApproved
	instance.FinishedAt = &now
	if err := e.db.Save(instance).Error; err != nil {
		return err
	}

	// 记录日志
	e.log(instance.ID, NodeTypeEnd, "finish", actorID, actorName, nil)

	return nil
}

// createApprovalTasks 创建审批待办任务
func (e *Engine) createApprovalTasks(instance *WorkflowInstance, node *WorkflowNode) error {
	now := time.Now()
	tasks := make([]WorkflowPendingTask, 0)

	fmt.Printf("[DEBUG] 节点 %s (ID=%d) 有 %d 个审批人配置\n", node.NodeKey, node.ID, len(node.Approvers))

	for _, approver := range node.Approvers {
		fmt.Printf("[DEBUG] 处理审批人: type=%s, id=%d, name=%s\n", approver.ApproverType, approver.ApproverID, approver.ApproverName)
		// 根据审批人类型处理
		switch approver.ApproverType {
		case ApproverTypeUser:
			// 直接给指定用户创建任务
			task := WorkflowPendingTask{
				InstanceID:   instance.ID,
				NodeID:       node.ID,
				NodeKey:      node.NodeKey,
				NodeName:     node.NodeName,
				BusinessType: instance.BusinessType,
				BusinessID:   instance.BusinessID,
				BusinessNo:   instance.BusinessNo,
				ApproverID:   uint(approver.ApproverID),
				ApproverName: approver.ApproverName,
				Status:       TaskStatusPending,
				IsParallel:   node.ApprovalType == ApprovalTypeParallel,
				ArrivedAt:    now,
			}
			tasks = append(tasks, task)

		case ApproverTypeRole:
			// 获取拥有该角色的所有用户
			users, err := e.getUsersByRole(approver.ApproverID)
			if err != nil {
				fmt.Printf("获取角色 %d 的用户失败: %v\n", approver.ApproverID, err)
				continue
			}
			for _, user := range users {
				task := WorkflowPendingTask{
					InstanceID:   instance.ID,
					NodeID:       node.ID,
					NodeKey:      node.NodeKey,
					NodeName:     node.NodeName,
					BusinessType: instance.BusinessType,
					BusinessID:   instance.BusinessID,
					BusinessNo:   instance.BusinessNo,
					ApproverID:   user.ID,
					ApproverName: user.FullName,
					Status:       TaskStatusPending,
					IsParallel:   node.ApprovalType == ApprovalTypeParallel,
					ArrivedAt:    now,
				}
				tasks = append(tasks, task)
			}

		case ApproverTypeDepartment:
			// 获取该部门的所有用户
			users, err := e.getUsersByDepartment(approver.ApproverID)
			if err != nil {
				fmt.Printf("获取部门 %d 的用户失败: %v\n", approver.ApproverID, err)
				continue
			}
			for _, user := range users {
				task := WorkflowPendingTask{
					InstanceID:   instance.ID,
					NodeID:       node.ID,
					NodeKey:      node.NodeKey,
					NodeName:     node.NodeName,
					BusinessType: instance.BusinessType,
					BusinessID:   instance.BusinessID,
					BusinessNo:   instance.BusinessNo,
					ApproverID:   user.ID,
					ApproverName: user.FullName,
					Status:       TaskStatusPending,
					IsParallel:   node.ApprovalType == ApprovalTypeParallel,
					ArrivedAt:    now,
				}
				tasks = append(tasks, task)
			}

		case ApproverTypeSuperior:
			// 获取申请人的上级
			superior, err := e.getUserSuperior(instance.InitiatorID)
			if err != nil {
				fmt.Printf("获取用户 %d 的上级失败: %v\n", instance.InitiatorID, err)
				continue
			}
			if superior != nil {
				task := WorkflowPendingTask{
					InstanceID:   instance.ID,
					NodeID:       node.ID,
					NodeKey:      node.NodeKey,
					NodeName:     node.NodeName,
					BusinessType: instance.BusinessType,
					BusinessID:   instance.BusinessID,
					BusinessNo:   instance.BusinessNo,
					ApproverID:   superior.ID,
					ApproverName: superior.FullName,
					Status:       TaskStatusPending,
					IsParallel:   node.ApprovalType == ApprovalTypeParallel,
					ArrivedAt:    now,
				}
				tasks = append(tasks, task)
			}
		}
	}

	if len(tasks) > 0 {
		if err := e.db.Create(&tasks).Error; err != nil {
			return err
		}
		fmt.Printf("为工作流实例 %d 的节点 %s 创建了 %d 个待办任务\n", instance.ID, node.NodeKey, len(tasks))
	} else {
		fmt.Printf("警告：工作流实例 %d 的节点 %s 没有创建任何待办任务\n", instance.ID, node.NodeKey)
	}

	// 发送通知给审批人
	if err := e.sendApprovalNotifications(instance, node); err != nil {
		// 通知发送失败不影响任务创建，仅记录日志
		fmt.Printf("发送审批通知失败: %v\n", err)
	}

	return nil
}

// checkAllApproved 检查是否所有审批人都已审批
func (e *Engine) checkAllApproved(instance *WorkflowInstance, node *WorkflowNode) (bool, error) {
	var totalApprovers int64
	if err := e.db.Model(&WorkflowNodeApprover{}).Where("node_id = ?", node.ID).Count(&totalApprovers).Error; err != nil {
		return false, err
	}

	var approvedCount int64
	if err := e.db.Model(&WorkflowPendingTask{}).Where("instance_id = ? AND node_id = ? AND status = ?", instance.ID, node.ID, TaskStatusApproved).Count(&approvedCount).Error; err != nil {
		return false, err
	}

	// 根据审批类型判断
	switch node.ApprovalType {
	case ApprovalTypeAny:
		// 任一审批通过即可
		return approvedCount > 0, nil
	case ApprovalTypeSequential, ApprovalTypeParallel:
		// 所有人都需要审批
		return approvedCount >= totalApprovers, nil
	default:
		return false, errors.New("未知的审批类型")
	}
}

// getNextNode 获取下一个节点
func (e *Engine) getNextNode(workflow *WorkflowDefinition, currentNodeKey string) (*WorkflowNode, error) {
	// 查找从当前节点出发的边
	var edge WorkflowEdge
	if err := e.db.Where("workflow_id = ? AND from_node = ?", workflow.ID, currentNodeKey).First(&edge).Error; err != nil {
		return nil, fmt.Errorf("未找到从节点 %s 出发的边", currentNodeKey)
	}

	// 查找目标节点（包含审批人）
	var node WorkflowNode
	if err := e.db.Preload("Approvers").Where("workflow_id = ? AND node_key = ?", workflow.ID, edge.ToNode).First(&node).Error; err != nil {
		return nil, fmt.Errorf("未找到节点 %s", edge.ToNode)
	}

	return &node, nil
}

// log 记录工作流日志
func (e *Engine) log(instanceID uint, nodeKey string, action string, actorID uint, actorName string, data map[string]any) error {
	// 简化版日志，不存储详细数据
	log := &WorkflowLog{
		InstanceID: instanceID,
		NodeKey:    nodeKey,
		Action:     action,
		ActorID:    actorID,
		ActorName:  actorName,
		// ActionData 可以序列化 data，这里简化处理
	}
	return e.db.Create(log).Error
}

// GetPendingTasks 获取用户的待办任务
func (e *Engine) GetPendingTasks(approverID uint) ([]WorkflowPendingTask, error) {
	var tasks []WorkflowPendingTask
	err := e.db.Preload("Instance").Where("approver_id = ? AND status = ?", approverID, TaskStatusPending).Find(&tasks).Error
	return tasks, err
}

// GetPendingTasksByBusiness 根据业务类型获取待办任务
func (e *Engine) GetPendingTasksByBusiness(approverID uint, businessType string) ([]WorkflowPendingTask, error) {
	var tasks []WorkflowPendingTask
	err := e.db.Preload("Instance").Where("approver_id = ? AND business_type = ? AND status = ?", approverID, businessType, TaskStatusPending).Find(&tasks).Error
	return tasks, err
}

// GetInstanceApprovals 获取实例的所有审批记录
func (e *Engine) GetInstanceApprovals(instanceID uint) ([]WorkflowApproval, error) {
	var approvals []WorkflowApproval
	err := e.db.Where("instance_id = ?", instanceID).Order("created_at ASC").Find(&approvals).Error
	return approvals, err
}

// GetInstanceLogs 获取实例的所有日志
func (e *Engine) GetInstanceLogs(instanceID uint) ([]WorkflowLog, error) {
	var logs []WorkflowLog
	err := e.db.Where("instance_id = ?", instanceID).Order("created_at ASC").Find(&logs).Error
	return logs, err
}

// GetWorkflowByModule 根据模块获取激活的工作流
func (e *Engine) GetWorkflowByModule(module string) (*WorkflowDefinition, error) {
	var workflow WorkflowDefinition
	err := e.db.Preload("Nodes").Preload("Nodes.Approvers").Preload("Edges").Where("module = ? AND is_active = ?", module, true).First(&workflow).Error
	return &workflow, err
}

// GetInstanceByBusiness 根据业务信息获取工作流实例
func (e *Engine) GetInstanceByBusiness(businessType string, businessID uint) (*WorkflowInstance, error) {
	var instance WorkflowInstance
	err := e.db.Preload("Workflow").Preload("Workflow.Nodes").Where("business_type = ? AND business_id = ?", businessType, businessID).First(&instance).Error
	return &instance, err
}

// sendApprovalNotifications 发送审批通知给审批人
func (e *Engine) sendApprovalNotifications(instance *WorkflowInstance, node *WorkflowNode) error {
	// 获取业务类型显示名称
	businessTypeName := getBusinessTypeName(instance.BusinessType)

	// 构建通知数据
	notificationData := map[string]interface{}{
		"instance_id":   instance.ID,
		"business_type": instance.BusinessType,
		"business_id":   instance.BusinessID,
		"business_no":   instance.BusinessNo,
		"node_key":      node.NodeKey,
		"node_name":     node.NodeName,
		"initiator":     instance.InitiatorName,
	}

	// 根据审批人类型发送通知
	for _, approver := range node.Approvers {
		// 确定通知类型
		var notificationType string
		if instance.BusinessType == "inbound_order" {
			notificationType = notification.TypeInboundApprove
		} else if instance.BusinessType == "requisition" {
			notificationType = notification.TypeRequisitionApprove
		} else {
			notificationType = notification.TypeSystem
		}

		// 生成通知标题和内容
		title := fmt.Sprintf("待审批：%s", businessTypeName)
		content := fmt.Sprintf("%s %s 需要您审批，单号：%s", businessTypeName, node.NodeName, instance.BusinessNo)

		// 根据审批人类型处理
		switch approver.ApproverType {
		case ApproverTypeUser:
			// 直接给指定用户发送通知
			if err := notification.CreateNotification(e.db, uint(approver.ApproverID), notificationType, title, content, notificationData); err != nil {
				return fmt.Errorf("发送通知给用户 %d 失败: %w", approver.ApproverID, err)
			}

		case ApproverTypeRole:
			// 获取拥有该角色的所有用户
			users, err := e.getUsersByRole(approver.ApproverID)
			if err != nil {
				return fmt.Errorf("获取角色 %d 的用户失败: %w", approver.ApproverID, err)
			}
			for _, user := range users {
				if err := notification.CreateNotification(e.db, user.ID, notificationType, title, content, notificationData); err != nil {
					return fmt.Errorf("发送通知给用户 %d 失败: %w", user.ID, err)
				}
			}

		case ApproverTypeDepartment:
			// 获取该部门的所有用户
			users, err := e.getUsersByDepartment(approver.ApproverID)
			if err != nil {
				return fmt.Errorf("获取部门 %d 的用户失败: %w", approver.ApproverID, err)
			}
			for _, user := range users {
				if err := notification.CreateNotification(e.db, user.ID, notificationType, title, content, notificationData); err != nil {
					return fmt.Errorf("发送通知给用户 %d 失败: %w", user.ID, err)
				}
			}

		case ApproverTypeSuperior:
			// 获取申请人的上级
			superior, err := e.getUserSuperior(instance.InitiatorID)
			if err != nil {
				return fmt.Errorf("获取用户 %d 的上级失败: %w", instance.InitiatorID, err)
			}
			if superior != nil {
				if err := notification.CreateNotification(e.db, superior.ID, notificationType, title, content, notificationData); err != nil {
					return fmt.Errorf("发送通知给上级 %d 失败: %w", superior.ID, err)
				}
			}
		}
	}

	return nil
}

// getUsersByRole 获取拥有指定角色的所有用户
func (e *Engine) getUsersByRole(roleID int) ([]auth.User, error) {
	var users []auth.User
	err := e.db.Raw(`
		SELECT DISTINCT u.* FROM users u
		INNER JOIN user_roles ur ON u.id = ur.user_id
		WHERE ur.role_id = ?
	`, roleID).Find(&users).Error
	return users, err
}

// getUsersByDepartment 获取指定部门的所有用户
func (e *Engine) getUsersByDepartment(departmentID int) ([]auth.User, error) {
	var users []auth.User
	// 假设 departmentID 存储在 group 字段中
	err := e.db.Where("group = ?", fmt.Sprintf("%d", departmentID)).Find(&users).Error
	return users, err
}

// getUserSuperior 获取用户的上级
func (e *Engine) getUserSuperior(userID uint) (*auth.User, error) {
	// 这里简化处理，实际可能需要根据组织架构查询
	// 暂时返回空，可以根据实际业务逻辑实现
	return nil, nil
}

// getBusinessTypeName 获取业务类型的显示名称
func getBusinessTypeName(businessType string) string {
	names := map[string]string{
		"inbound_order": "入库单",
		"requisition":   "领料单",
	}
	if name, ok := names[businessType]; ok {
		return name
	}
	return "业务单据"
}
