package appointment

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/yourorg/material-backend/backend/internal/api/notification"
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
	if err := s.updateAppointmentStatusByWorkflow(&instance, extraData); err != nil {
		return err
	}

	// 通过 WebSocket 广播审批更新通知
	s.broadcastApprovalUpdate(instance.BusinessID, approverID, approverName, req.Action, req.Comment)

	return nil
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
			s.db.Table("users").Where("id = ?", workerID).Pluck("full_name", &workerName)
			appointment.AssignedWorkerName = workerName
		}

		// 保存并预约日历
		if err := s.db.Save(&appointment).Error; err != nil {
			return err
		}

		// 预约日历
		if appointment.AssignedWorkerID != nil {
			if _, err := s.appointmentService.AssignWorker(appointment.ID, *appointment.AssignedWorkerID, appointment.AssignedWorkerName); err != nil {
				// 日历预约失败不影响状态更新
				fmt.Printf("Warning: failed to book calendar: %v\n", err)
			}
		}

		// 通知申请人和作业人员
		s.notifyAppointmentApproved(&appointment)

	case workflow.InstanceStatusRejected:
		// 工作流拒绝
		appointment.Status = StatusRejected
		s.db.Save(&appointment)
		// 通知申请人被拒绝
		s.notifyAppointmentRejected(&appointment)

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
		// 返回加急工作流ID（数据库中配置的ID为11）
		return 11, nil
	}

	// 返回普通工作流ID（数据库中配置的ID为10）
	return 10, nil
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
			"id":            log.ID,
			"node_name":     log.NodeKey,
			"action":        log.Action,
			"approver_id":   log.ActorID,
			"approver_name": log.ActorName,
			"remark":        log.ActionData,
			"created_at":    log.CreatedAt.Format("2006-01-02 15:04:05"),
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

	// 更新工作流实例状态为已取消
	var instance workflow.WorkflowInstance
	if err := s.db.First(&instance, *appointment.WorkflowInstanceID).Error; err != nil {
		return err
	}
	instance.Status = workflow.InstanceStatusCancelled
	now := time.Now()
	instance.FinishedAt = &now
	if err := s.db.Save(&instance).Error; err != nil {
		return err
	}

	// 取消所有待办任务
	s.db.Model(&workflow.WorkflowPendingTask{}).
		Where("instance_id = ? AND status = ?", *appointment.WorkflowInstanceID, workflow.TaskStatusPending).
		Updates(map[string]any{
			"status":       workflow.TaskStatusCancelled,
			"processed_at": now,
		})

	// 更新预约单状态
	appointment.Status = StatusDraft
	appointment.WorkflowInstanceID = nil
	appointment.SubmittedAt = nil

	return s.db.Save(&appointment).Error
}

// TransferApproval 转交审批
func (s *WorkflowService) TransferApproval(instanceID uint, fromUserID uint, toUserID uint, toUserName string, remark string) error {
	// TODO: 实现转交审批功能
	return errors.New("转交审批功能暂未实现")
}

// AddNodeApprover 添加节点审批人
func (s *WorkflowService) AddNodeApprover(instanceID uint, nodeID uint, newApproverID uint, newApproverName string) error {
	// TODO: 实现添加节点审批人功能
	return errors.New("添加节点审批人功能暂未实现")
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
	if err := s.db.Where("workflow_id = ?", instance.WorkflowID).Find(&nodes).Error; err != nil {
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
		// 查找该节点的日志（通过 NodeKey 匹配）
		var nodeLog *workflow.WorkflowLog
		for j := range logs {
			if logs[j].NodeKey == node.NodeKey {
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
		isCurrent := instance.CurrentNode == node.NodeKey

		nodeProgress[i] = map[string]any{
			"node_id":      node.ID,
			"node_name":    node.NodeName,
			"node_type":    node.NodeType,
			"node_key":     node.NodeKey,
			"status":       status,
			"is_current":   isCurrent,
			"processed_at": nil,
		}

		if nodeLog != nil {
			nodeProgress[i]["processed_at"] = nodeLog.CreatedAt.Format("2006-01-02 15:04:05")
			nodeProgress[i]["approver_name"] = nodeLog.ActorName
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

// notifyAppointmentApproved 通知申请人和作业人员审批通过
func (s *WorkflowService) notifyAppointmentApproved(appointment *ConstructionAppointment) {
	// 构建通知数据
	notificationData := map[string]interface{}{
		"appointment_id":   appointment.ID,
		"appointment_no":   appointment.AppointmentNo,
		"work_date":        appointment.WorkDate.Format("2006-01-02"),
		"time_slot":        appointment.TimeSlot,
		"work_location":    appointment.WorkLocation,
		"work_content":     appointment.WorkContent,
		"approved_at":      appointment.ApprovedAt,
	}

	// 通知申请人（创建人）
	applicantTitle := "预约审批通过"
	applicantContent := fmt.Sprintf("您提交的施工预约单 %s 已审批通过，作业时间：%s %s，地点：%s",
		appointment.AppointmentNo,
		appointment.WorkDate.Format("2006-01-02"),
		appointment.TimeSlot,
		appointment.WorkLocation)

	if err := notification.CreateNotification(s.db, appointment.ApplicantID, notification.TypeAppointmentApprove, applicantTitle, applicantContent, notificationData); err != nil {
		fmt.Printf("通知申请人失败: %v\n", err)
	}

	// 通知作业人员
	if appointment.AssignedWorkerID != nil && *appointment.AssignedWorkerID != 0 {
		workerTitle := "新的作业任务分配"
		workerContent := fmt.Sprintf("您被分配了一项施工任务，单号：%s，作业时间：%s %s，地点：%s，内容：%s",
			appointment.AppointmentNo,
			appointment.WorkDate.Format("2006-01-02"),
			appointment.TimeSlot,
			appointment.WorkLocation,
			appointment.WorkContent)

		workerData := map[string]interface{}{
			"appointment_id": appointment.ID,
			"appointment_no": appointment.AppointmentNo,
			"work_date":      appointment.WorkDate.Format("2006-01-02"),
			"time_slot":      appointment.TimeSlot,
			"work_location":  appointment.WorkLocation,
			"work_content":   appointment.WorkContent,
			"assigned_at":    appointment.ApprovedAt,
		}

		if err := notification.CreateNotification(s.db, *appointment.AssignedWorkerID, notification.TypeAppointmentApprove, workerTitle, workerContent, workerData); err != nil {
			fmt.Printf("通知作业人员失败: %v\n", err)
		}
	}

	// 处理多作业人员通知
	if appointment.AssignedWorkerIDs != "" {
		var workerIDs []uint
		if err := json.Unmarshal([]byte(appointment.AssignedWorkerIDs), &workerIDs); err == nil {
			for _, workerID := range workerIDs {
				// 跳过已经通知的主作业人员
				if appointment.AssignedWorkerID != nil && workerID == *appointment.AssignedWorkerID {
					continue
				}
				workerTitle := "新的作业任务分配"
				workerContent := fmt.Sprintf("您被分配了一项施工任务，单号：%s，作业时间：%s %s，地点：%s",
					appointment.AppointmentNo,
					appointment.WorkDate.Format("2006-01-02"),
					appointment.TimeSlot,
					appointment.WorkLocation)

				if err := notification.CreateNotification(s.db, workerID, notification.TypeAppointmentApprove, workerTitle, workerContent, notificationData); err != nil {
					fmt.Printf("通知作业人员 %d 失败: %v\n", workerID, err)
				}
			}
		}
	}
}

// notifyAppointmentRejected 通知申请人审批被拒绝
func (s *WorkflowService) notifyAppointmentRejected(appointment *ConstructionAppointment) {
	// 构建通知数据
	notificationData := map[string]interface{}{
		"appointment_id": appointment.ID,
		"appointment_no": appointment.AppointmentNo,
		"work_date":      appointment.WorkDate.Format("2006-01-02"),
		"time_slot":      appointment.TimeSlot,
		"work_location":  appointment.WorkLocation,
		"work_content":   appointment.WorkContent,
	}

	title := "预约审批未通过"
	content := fmt.Sprintf("您提交的施工预约单 %s 未通过审批，请查看详情并修改后重新提交", appointment.AppointmentNo)

	if err := notification.CreateNotification(s.db, appointment.ApplicantID, notification.TypeAppointmentApprove, title, content, notificationData); err != nil {
		fmt.Printf("通知申请人失败: %v\n", err)
	}
}

// broadcastApprovalUpdate 通过 WebSocket 广播审批更新通知
func (s *WorkflowService) broadcastApprovalUpdate(appointmentID uint, approverID uint, approverName string, action string, comment string) {
	// 获取预约单信息
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, appointmentID).Error; err != nil {
		fmt.Printf("获取预约单失败，无法广播通知: %v\n", err)
		return
	}

	// 构建需要通知的用户列表
	userIDs := make([]uint, 0)
	userIDsMap := make(map[uint]bool)

	// 添加申请人
	userIDsMap[appointment.ApplicantID] = true

	// 添加当前审批人（待办任务的审批人）
	var pendingTasks []struct {
		ApproverID uint
	}
	if err := s.db.Table("workflow_pending_tasks").
		Where("instance_id = ? AND status = ?", appointment.WorkflowInstanceID, "pending").
		Pluck("approver_id", &pendingTasks).Error; err == nil {
		for _, task := range pendingTasks {
			userIDsMap[task.ApproverID] = true
		}
	}

	// 转换为数组
	for userID := range userIDsMap {
		userIDs = append(userIDs, userID)
	}

	// 构建消息数据
	messageData := map[string]interface{}{
		"appointment_id":   appointment.ID,
		"appointment_no":   appointment.AppointmentNo,
		"approver_id":      approverID,
		"approver_name":    approverName,
		"action":           action,
		"comment":          comment,
		"status":           appointment.Status,
		"workflow_instance_id": appointment.WorkflowInstanceID,
	}

	// 通过 WebSocket hub 广播消息
	hub := notification.GetHub()
	if hub != nil {
		for _, userID := range userIDs {
			// 构建完整消息
			message := map[string]interface{}{
				"type": "appointment_approval_update",
				"data": messageData,
			}

			// 序列化消息
			data, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("序列化审批更新消息失败: %v\n", err)
				continue
			}

			// 发送给指定用户
			hub.BroadcastToUser(userID, data)
			fmt.Printf("已向用户 %d 广播审批更新通知: 预约单 %s, 操作: %s\n", userID, appointment.AppointmentNo, action)
		}
	}
}
