package appointment

import (
	"errors"
	"fmt"
	"time"

	"github.com/yourorg/material-backend/backend/internal/api/workflow"
	"gorm.io/gorm"
)

// WorkflowService 工作流服务
type WorkflowService struct {
	db                *gorm.DB
	appointmentService *AppointmentService
	workflowEngine    *workflow.Engine
}

// NewWorkflowService 创建工作流服务
func NewWorkflowService(db *gorm.DB) *WorkflowService {
	return &WorkflowService{
		db:                 db,
		appointmentService: NewAppointmentService(db),
		workflowEngine:     workflow.NewEngine(db),
	}
}

// StartApprovalWorkflow 启动审批工作流
func (s *WorkflowService) StartApprovalWorkflow(appointmentID uint, workflowID uint) (*ConstructionAppointment, *workflow.WorkflowInstance, error) {
	// 获取预约单
	appointment, err := s.appointmentService.GetByID(appointmentID)
	if err != nil {
		return nil, nil, err
	}

	// 检查状态
	if appointment.Status != StatusPending {
		return nil, nil, errors.New("只有待审批状态的预约单可以启动工作流")
	}

	// 检查是否已有工作流实例
	if appointment.WorkflowInstanceID != nil {
		return nil, nil, errors.New("该预约单已启动工作流")
	}

	// 确定工作流类型（普通或加急）
	actualWorkflowID := workflowID
	if actualWorkflowID == 0 {
		var err error
		actualWorkflowID, err = s.getDefaultWorkflowID(appointment)
		if err != nil {
			return nil, nil, err
		}
	}

	// 启动工作流
	instance, err := s.workflowEngine.StartWorkflow(
		actualWorkflowID,
		"construction_appointment",
		appointment.ID,
		appointment.AppointmentNo,
		appointment.ApplicantID,
		appointment.ApplicantName,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("启动工作流失败: %w", err)
	}

	// 更新预约单的工作流实例ID
	appointment.WorkflowInstanceID = &instance.ID
	if err := s.db.Save(appointment).Error; err != nil {
		return nil, nil, fmt.Errorf("更新预约单工作流失败: %w", err)
	}

	return appointment, instance, nil
}

// ProcessApproval 处理审批
func (s *WorkflowService) ProcessApproval(instanceID uint, approverID uint, approverName string, req ApproveAppointmentRequest) error {
	// 构建额外数据
	extraData := make(map[string]any)
	if req.AssignNow && req.WorkerID != nil {
		extraData["assign_worker"] = true
		extraData["worker_id"] = *req.WorkerID
	}

	// 处理审批
	action := req.Action
	if action == "approve" {
		action = workflow.ActionApprove
	} else if action == "reject" {
		action = workflow.ActionReject
	}

	err := s.workflowEngine.ProcessApproval(instanceID, approverID, approverName, action, req.Comment, extraData)
	if err != nil {
		return err
	}

	// 获取工作流实例
	var instance workflow.WorkflowInstance
	if err := s.db.First(&instance, instanceID).Error; err != nil {
		return err
	}

	// 根据工作流状态更新预约单状态
	return s.updateAppointmentStatusByWorkflow(&instance, extraData)
}

// updateAppointmentStatusByWorkflow 根据工作流状态更新预约单状态
func (s *WorkflowService) updateAppointmentStatusByWorkflow(instance *workflow.WorkflowInstance, extraData map[string]any) error {
	var appointment ConstructionAppointment
	if err := s.db.Where("id = ?", instance.BusinessID).First(&appointment).Error; err != nil {
		return err
	}

	switch instance.Status {
	case workflow.InstanceStatusApproved:
		// 工作流审批通过
		appointment.Status = StatusScheduled
		now := time.Now()
		appointment.ApprovedAt = &now

		// 如果指定了作业人员，立即分配
		if workerID, ok := extraData["worker_id"].(uint); ok && workerID != 0 {
			appointment.AssignedWorkerID = &workerID
			// 获取作业人员姓名
			var workerName string
			s.db.Table("users").Where("id = ?", workerID).Pluck("name", &workerName)
			appointment.AssignedWorkerName = workerName
		}

		// 保存并预约日历
		if err := s.db.Save(&appointment).Error; err != nil {
			return err
		}

		// 预约日历
		if appointment.AssignedWorkerID != nil {
			if err := s.appointmentService.AssignWorker(appointment.ID, *appointment.AssignedWorkerID, appointment.AssignedWorkerName); err != nil {
				// 日历预约失败不影响状态更新
				fmt.Printf("Warning: failed to book calendar: %v\n", err)
			}
		}

	case workflow.InstanceStatusRejected:
		// 工作流拒绝
		appointment.Status = StatusRejected

	case workflow.InstanceStatusCancelled:
		// 工作流取消
		appointment.Status = StatusCancelled
	}

	return s.db.Save(&appointment).Error
}

// getDefaultWorkflowID 获取默认工作流ID
func (s *WorkflowService) getDefaultWorkflowID(appointment *ConstructionAppointment) (uint, error) {
	// 检查是否是加急预约
	if appointment.IsUrgent && appointment.NeedsUrgentApproval() {
		// 返回加急工作流ID（需要在数据库中配置）
		// 这里假设加急工作流ID为3
		return 3, nil
	}

	// 返回普通工作流ID
	// 这里假设普通工作流ID为2
	return 2, nil
}

// GetPendingApprovals 获取待审批列表
func (s *WorkflowService) GetPendingApprovals(userID uint, page, pageSize int) ([]map[string]any, int64, error) {
	// 获取待办任务
	var tasks []workflow.WorkflowPendingTask
	var total int64

	query := s.db.Where("approver_id = ? AND status = ?", userID, workflow.TaskStatusPending).
		Where("business_type = ?", "construction_appointment")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	// 获取预约单详情
	result := make([]map[string]any, 0, len(tasks))
	for _, task := range tasks {
		var appointment ConstructionAppointment
		if err := s.db.Where("id = ?", task.BusinessID).First(&appointment).Error; err != nil {
			continue
		}

		item := appointment.ToDTO()
		item["workflow_instance_id"] = task.InstanceID
		item["node_id"] = task.NodeID
		item["node_name"] = task.NodeName
		item["task_id"] = task.ID
		result = append(result, item)
	}

	return result, total, nil
}

// GetApprovalHistory 获取审批历史
func (s *WorkflowService) GetApprovalHistory(appointmentID uint) ([]map[string]any, error) {
	// 获取预约单
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, appointmentID).Error; err != nil {
		return nil, err
	}

	if appointment.WorkflowInstanceID == nil {
		return []map[string]any{}, nil
	}

	// 获取工作流日志
	var logs []workflow.WorkflowLog
	if err := s.db.Where("instance_id = ?", *appointment.WorkflowInstanceID).
		Order("created_at ASC").
		Find(&logs).Error; err != nil {
		return nil, err
	}

	// 转换为DTO
	result := make([]map[string]any, len(logs))
	for i, log := range logs {
		result[i] = map[string]any{
			"id":           log.ID,
			"node_name":    log.NodeName,
			"action":       log.Action,
			"approver_id":  log.ApproverID,
			"approver_name": log.ApproverName,
			"remark":       log.Remark,
			"created_at":   log.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return result, nil
}

// GetCurrentApprovalNode 获取当前审批节点信息
func (s *WorkflowService) GetCurrentApprovalNode(appointmentID uint) (map[string]any, error) {
	// 获取预约单
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, appointmentID).Error; err != nil {
		return nil, err
	}

	if appointment.WorkflowInstanceID == nil {
		return nil, errors.New("该预约单未启动工作流")
	}

	// 获取待办任务
	var tasks []workflow.WorkflowPendingTask
	if err := s.db.Where("instance_id = ? AND status = ?", *appointment.WorkflowInstanceID, workflow.TaskStatusPending).
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, errors.New("没有待审批任务")
	}

	// 构建审批节点信息
	approvers := make([]map[string]any, len(tasks))
	for i, task := range tasks {
		approvers[i] = map[string]any{
			"id":   task.ApproverID,
			"name": task.ApproverName,
		}
	}

	return map[string]any{
		"node_id":   tasks[0].NodeID,
		"node_name": tasks[0].NodeName,
		"approvers": approvers,
	}, nil
}

// RecallWorkflow 撤回工作流
func (s *WorkflowService) RecallWorkflow(appointmentID uint, userID uint) error {
	// 获取预约单
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, appointmentID).Error; err != nil {
		return err
	}

	// 检查权限（只有发起人可以撤回）
	if appointment.ApplicantID != userID {
		return errors.New("只有发起人可以撤回")
	}

	if appointment.WorkflowInstanceID == nil {
		return errors.New("该预约单未启动工作流")
	}

	// 取消工作流实例
	if err := s.workflowEngine.CancelInstance(*appointment.WorkflowInstanceID, userID, appointment.ApplicantName, "申请人撤回"); err != nil {
		return err
	}

	// 更新预约单状态
	appointment.Status = StatusDraft
	appointment.WorkflowInstanceID = nil
	appointment.SubmittedAt = nil

	return s.db.Save(&appointment).Error
}

// TransferApproval 转交审批
func (s *WorkflowService) TransferApproval(instanceID uint, fromUserID uint, toUserID uint, toUserName string, remark string) error {
	return s.workflowEngine.TransferApproval(instanceID, fromUserID, toUserID, toUserName, remark)
}

// AddNodeApprover 添加节点审批人
func (s *WorkflowService) AddNodeApprover(instanceID uint, nodeID uint, newApproverID uint, newApproverName string) error {
	return s.workflowEngine.AddNodeApprover(instanceID, nodeID, newApproverID, newApproverName)
}

// BatchApprove 批量审批
func (s *WorkflowService) BatchApprove(instanceIDs []uint, approverID uint, approverName string, action string, comment string) ([]map[string]any, []error) {
	results := make([]map[string]any, 0)
	errs := make([]error, 0)

	for _, instanceID := range instanceIDs {
		req := ApproveAppointmentRequest{
			Action: action,
			Comment: comment,
		}
		err := s.ProcessApproval(instanceID, approverID, approverName, req)
		if err != nil {
			errs = append(errs, err)
			results = append(results, map[string]any{
				"instance_id": instanceID,
				"success":     false,
				"error":       err.Error(),
			})
		} else {
			results = append(results, map[string]any{
				"instance_id": instanceID,
				"success":     true,
			})
		}
	}

	return results, errs
}

// GetWorkflowProgress 获取工作流进度
func (s *WorkflowService) GetWorkflowProgress(appointmentID uint) (map[string]any, error) {
	// 获取预约单
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, appointmentID).Error; err != nil {
		return nil, err
	}

	if appointment.WorkflowInstanceID == nil {
		return map[string]any{
			"has_workflow": false,
			"status":       appointment.Status,
		}, nil
	}

	// 获取工作流实例
	var instance workflow.WorkflowInstance
	if err := s.db.First(&instance, *appointment.WorkflowInstanceID).Error; err != nil {
		return nil, err
	}

	// 获取所有节点
	var nodes []workflow.WorkflowNode
	if err := s.db.Where("workflow_id = ?", instance.WorkflowID).Order("node_order ASC").Find(&nodes).Error; err != nil {
		return nil, err
	}

	// 获取所有日志
	var logs []workflow.WorkflowLog
	if err := s.db.Where("instance_id = ?", *appointment.WorkflowInstanceID).
		Order("created_at ASC").
		Find(&logs).Error; err != nil {
		return nil, err
	}

	// 构建进度信息
	nodeProgress := make([]map[string]any, len(nodes))
	for i, node := range nodes {
		// 查找该节点的日志
		var nodeLog *workflow.WorkflowLog
		for j := range logs {
			if logs[j].NodeID == node.ID {
				nodeLog = &logs[j]
				break
			}
		}

		status := "pending"
		if nodeLog != nil {
			status = "completed"
			if nodeLog.Action == workflow.ActionReject {
				status = "rejected"
			}
		}

		// 检查是否是当前节点
		isCurrent := instance.CurrentNode == node.NodeType

		nodeProgress[i] = map[string]any{
			"node_id":     node.ID,
			"node_name":   node.NodeName,
			"node_type":   node.NodeType,
			"node_order":  node.NodeOrder,
			"status":      status,
			"is_current":  isCurrent,
			"processed_at": nil,
		}

		if nodeLog != nil {
			nodeProgress[i]["processed_at"] = nodeLog.CreatedAt.Format("2006-01-02 15:04:05")
			nodeProgress[i]["approver_name"] = nodeLog.ApproverName
		}
	}

	return map[string]any{
		"has_workflow":     true,
		"workflow_id":      instance.WorkflowID,
		"instance_id":      instance.ID,
		"status":           instance.Status,
		"current_node":     instance.CurrentNode,
		"started_at":       instance.StartedAt.Format("2006-01-02 15:04:05"),
		"finished_at":      instance.FinishedAt,
		"nodes":            nodeProgress,
	}, nil
}
