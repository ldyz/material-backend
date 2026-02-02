package workflow

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"gorm.io/gorm"
)

// Handler 工作流处理器
type Handler struct {
	db     *gorm.DB
	engine *Engine
}

// NewHandler 创建工作流处理器
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db:     db,
		engine: NewEngine(db),
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	workflows := r.Group("/workflows")
	{
		workflows.GET("", h.listWorkflows)           // 获取工作流列表
		workflows.POST("", h.createWorkflow)         // 创建工作流
		workflows.GET("/:id", h.getWorkflow)         // 获取工作流详情
		workflows.PUT("/:id", h.updateWorkflow)      // 更新工作流
		workflows.DELETE("/:id", h.deleteWorkflow)   // 删除工作流
		workflows.PUT("/:id/activate", h.activateWorkflow) // 激活工作流
		workflows.PUT("/:id/deactivate", h.deactivateWorkflow) // 停用工作流
	}

	instances := r.Group("/workflow-instances")
	{
		instances.GET("", h.listInstances)                    // 获取实例列表
		instances.GET("/:id", h.getInstance)                  // 获取实例详情
		instances.GET("/:id/approvals", h.getInstanceApprovals) // 获取审批记录
		instances.GET("/:id/logs", h.getInstanceLogs)         // 获取操作日志
		instances.POST("/:id/resubmit", h.resubmitInstance)   // 重新提交
	}

	tasks := r.Group("/workflow-tasks")
	{
		tasks.GET("/pending", h.getPendingTasks)              // 获取我的待办任务
		tasks.GET("/pending/:businessType", h.getPendingTasksByType) // 按类型获取待办
		tasks.POST("/:id/approve", h.approveTask)             // 审批通过
		tasks.POST("/:id/reject", h.rejectTask)               // 审批拒绝
		tasks.POST("/:id/return", h.returnTask)               // 退回
		tasks.POST("/:id/comment", h.commentTask)             // 评论
	}
}

// listWorkflows 获取工作流列表
func (h *Handler) listWorkflows(c *gin.Context) {
	var workflows []WorkflowDefinition

	// 支持按模块筛选
	module := c.Query("module")
	query := h.db.Model(&WorkflowDefinition{})
	if module != "" {
		query = query.Where("module = ?", module)
	}

	if err := query.Order("id DESC").Find(&workflows).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "获取工作流列表失败")
		return
	}

	data := make([]map[string]any, 0)
	for _, w := range workflows {
		data = append(data, w.ToDTO())
	}

	response.Success(c, data)
}

// createWorkflow 创建工作流
func (h *Handler) createWorkflow(c *gin.Context) {
	var req struct {
		Name        string                   `json:"name" binding:"required"`
		Description string                   `json:"description"`
		Module      string                   `json:"module" binding:"required"`
		Nodes       []map[string]any         `json:"nodes" binding:"required"`
		Edges       []map[string]any         `json:"edges" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 开始事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建工作流定义
	workflow := &WorkflowDefinition{
		Name:        req.Name,
		Description: req.Description,
		Module:      req.Module,
		Version:     1,
		IsActive:    false, // 默认不激活
	}

	if err := tx.Create(workflow).Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "创建工作流失败")
		return
	}

	// 创建节点
	for _, nodeData := range req.Nodes {
		node := &WorkflowNode{
			WorkflowID:    workflow.ID,
			NodeKey:       nodeData["node_key"].(string),
			NodeType:      nodeData["node_type"].(string),
			NodeName:      nodeData["node_name"].(string),
			Description:   getStringValue(nodeData, "description"),
			ApprovalType:  getStringValue(nodeData, "approval_type"),
			TimeoutHours:  getIntValue(nodeData, "timeout_hours"),
			AutoApprove:   getBoolValue(nodeData, "auto_approve"),
			IsRequired:    getBoolValue(nodeData, "is_required"),
			X:             getIntValue(nodeData, "x"),
			Y:             getIntValue(nodeData, "y"),
		}

		if err := tx.Create(node).Error; err != nil {
			tx.Rollback()
			response.Error(c, http.StatusInternalServerError, "创建节点失败")
			return
		}

		// 创建审批人配置
		if approversData, ok := nodeData["approvers"].([]any); ok {
			for _, approverData := range approversData {
				approverMap := approverData.(map[string]any)
				approver := &WorkflowNodeApprover{
					NodeID:       node.ID,
					ApproverType: approverMap["approver_type"].(string),
					ApproverID:   int(approverMap["approver_id"].(float64)),
					ApproverName: approverMap["approver_name"].(string),
					Sequence:     int(approverMap["sequence"].(float64)),
				}

				if err := tx.Create(approver).Error; err != nil {
					tx.Rollback()
					response.Error(c, http.StatusInternalServerError, "创建审批人配置失败")
					return
				}
			}
		}
	}

	// 创建边
	for _, edgeData := range req.Edges {
		edge := &WorkflowEdge{
			WorkflowID:          workflow.ID,
			FromNode:            edgeData["from_node"].(string),
			ToNode:              edgeData["to_node"].(string),
			ConditionExpression: getStringValue(edgeData, "condition_expression"),
		}

		if err := tx.Create(edge).Error; err != nil {
			tx.Rollback()
			response.Error(c, http.StatusInternalServerError, "创建连接线失败")
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "创建工作流失败")
		return
	}

	response.SuccessWithMessage(c, workflow.ToDTO(), "工作流创建成功")
}

// getWorkflow 获取工作流详情
func (h *Handler) getWorkflow(c *gin.Context) {
	id := c.Param("id")

	var workflow WorkflowDefinition
	if err := h.db.Preload("Nodes").Preload("Nodes.Approvers").Preload("Edges").First(&workflow, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "工作流不存在")
		return
	}

	response.Success(c, workflow.ToDTO())
}

// updateWorkflow 更新工作流
func (h *Handler) updateWorkflow(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string           `json:"name"`
		Description string           `json:"description"`
		Nodes       []map[string]any `json:"nodes"`
		Edges       []map[string]any `json:"edges"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 获取现有工作流
	var workflow WorkflowDefinition
	if err := h.db.First(&workflow, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "工作流不存在")
		return
	}

	// 如果要修改节点或边，检查是否有正在运行的实例
	if len(req.Nodes) > 0 || len(req.Edges) > 0 {
		var count int64
		h.db.Model(&WorkflowInstance{}).Where("workflow_id = ? AND status = ?", id, InstanceStatusPending).Count(&count)
		if count > 0 {
			response.Error(c, http.StatusBadRequest, "该工作流有正在运行的实例，无法修改流程结构")
			return
		}
	}

	// 开始事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新基本信息
	if req.Name != "" {
		workflow.Name = req.Name
	}
	if req.Description != "" {
		workflow.Description = req.Description
	}
	if err := tx.Save(&workflow).Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "更新工作流失败")
		return
	}

	// 如果提供了节点和边，则重新创建
	if len(req.Nodes) > 0 {
		// 获取现有节点的ID
		var existingNodeIDs []uint
		h.db.Model(&WorkflowNode{}).Where("workflow_id = ?", workflow.ID).Pluck("id", &existingNodeIDs)

		// 删除旧的节点和审批人配置
		// 注意：必须按照外键依赖顺序删除：审批记录 -> 待办任务 -> 审批人配置 -> 节点
		if len(existingNodeIDs) > 0 {
			// 删除工作流实例的审批记录（通过实例关联）
			var instanceIDs []uint
			tx.Model(&WorkflowInstance{}).Where("workflow_id = ?", workflow.ID).Pluck("id", &instanceIDs)
			if len(instanceIDs) > 0 {
				tx.Where("instance_id IN ?", instanceIDs).Delete(&WorkflowApproval{})
				tx.Where("instance_id IN ?", instanceIDs).Delete(&WorkflowPendingTask{})
			}
			// 删除节点审批人配置
			tx.Where("node_id IN ?", existingNodeIDs).Delete(&WorkflowNodeApprover{})
		}
		tx.Where("workflow_id = ?", workflow.ID).Delete(&WorkflowNode{})

		// 创建新节点
		for _, nodeData := range req.Nodes {
			node := &WorkflowNode{
				WorkflowID:    workflow.ID,
				NodeKey:       nodeData["node_key"].(string),
				NodeType:      nodeData["node_type"].(string),
				NodeName:      nodeData["node_name"].(string),
				Description:   getStringValue(nodeData, "description"),
				ApprovalType:  getStringValue(nodeData, "approval_type"),
				TimeoutHours:  getIntValue(nodeData, "timeout_hours"),
				AutoApprove:   getBoolValue(nodeData, "auto_approve"),
				IsRequired:    getBoolValue(nodeData, "is_required"),
				X:             getIntValue(nodeData, "x"),
				Y:             getIntValue(nodeData, "y"),
			}

			if err := tx.Create(node).Error; err != nil {
				tx.Rollback()
				response.Error(c, http.StatusInternalServerError, "创建节点失败")
				return
			}

			// 创建审批人配置
			if approversData, ok := nodeData["approvers"].([]any); ok {
				for _, approverData := range approversData {
					approverMap := approverData.(map[string]any)
					approver := &WorkflowNodeApprover{
						NodeID:       node.ID,
						ApproverType: approverMap["approver_type"].(string),
						ApproverID:   int(approverMap["approver_id"].(float64)),
						ApproverName: approverMap["approver_name"].(string),
						Sequence:     int(approverMap["sequence"].(float64)),
					}

					if err := tx.Create(approver).Error; err != nil {
						tx.Rollback()
						response.Error(c, http.StatusInternalServerError, "创建审批人配置失败")
						return
					}
				}
			}
		}
	}

	// 删除旧的边并创建新边（不管 req.Edges 是否为空，都需要删除旧边）
	tx.Where("workflow_id = ?", workflow.ID).Delete(&WorkflowEdge{})

	// 创建新边
	for _, edgeData := range req.Edges {
		edge := &WorkflowEdge{
			WorkflowID:          workflow.ID,
			FromNode:            edgeData["from_node"].(string),
			ToNode:              edgeData["to_node"].(string),
			ConditionExpression: getStringValue(edgeData, "condition_expression"),
		}

		if err := tx.Create(edge).Error; err != nil {
			tx.Rollback()
			response.Error(c, http.StatusInternalServerError, "创建连接线失败")
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "更新工作流失败")
		return
	}

	// 重新加载并返回
	h.db.Preload("Nodes").Preload("Nodes.Approvers").Preload("Edges").First(&workflow, id)
	response.SuccessWithMessage(c, workflow.ToDTO(), "工作流更新成功")
}

// deleteWorkflow 删除工作流
func (h *Handler) deleteWorkflow(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有正在运行的实例
	var count int64
	h.db.Model(&WorkflowInstance{}).Where("workflow_id = ? AND status = ?", id, InstanceStatusPending).Count(&count)
	if count > 0 {
		response.Error(c, http.StatusBadRequest, "该工作流有正在运行的实例，无法删除")
		return
	}

	if err := h.db.Delete(&WorkflowDefinition{}, id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除工作流失败")
		return
	}

	response.SuccessWithMessage(c, nil, "工作流删除成功")
}

// activateWorkflow 激活工作流
func (h *Handler) activateWorkflow(c *gin.Context) {
	id := c.Param("id")

	// 先停用该模块的其他工作流
	var workflow WorkflowDefinition
	if err := h.db.First(&workflow, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "工作流不存在")
		return
	}

	h.db.Model(&WorkflowDefinition{}).Where("module = ? AND id != ?", workflow.Module, id).Update("is_active", false)

	// 激活当前工作流
	workflow.IsActive = true
	if err := h.db.Save(&workflow).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "激活工作流失败")
		return
	}

	response.SuccessWithMessage(c, workflow.ToDTO(), "工作流已激活")
}

// deactivateWorkflow 停用工作流
func (h *Handler) deactivateWorkflow(c *gin.Context) {
	id := c.Param("id")

	var workflow WorkflowDefinition
	if err := h.db.First(&workflow, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "工作流不存在")
		return
	}

	workflow.IsActive = false
	if err := h.db.Save(&workflow).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "停用工作流失败")
		return
	}

	response.SuccessWithMessage(c, workflow.ToDTO(), "工作流已停用")
}

// listInstances 获取工作流实例列表
func (h *Handler) listInstances(c *gin.Context) {
	var instances []WorkflowInstance

	// 支持按业务类型、状态筛选
	businessType := c.Query("business_type")
	status := c.Query("status")
	query := h.db.Model(&WorkflowInstance{}).Preload("Workflow")

	if businessType != "" {
		query = query.Where("business_type = ?", businessType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("id DESC").Find(&instances).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "获取实例列表失败")
		return
	}

	data := make([]map[string]any, 0)
	for _, i := range instances {
		data = append(data, i.ToDTO())
	}

	response.Success(c, data)
}

// getInstance 获取实例详情
func (h *Handler) getInstance(c *gin.Context) {
	id := c.Param("id")

	var instance WorkflowInstance
	if err := h.db.Preload("Workflow").Preload("Workflow.Nodes").Preload("Approvals").First(&instance, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "实例不存在")
		return
	}

	response.Success(c, instance.ToDTO())
}

// getInstanceApprovals 获取实例的审批记录
func (h *Handler) getInstanceApprovals(c *gin.Context) {
	id := c.Param("id")

	approvals, err := h.engine.GetInstanceApprovals(str2uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取审批记录失败")
		return
	}

	data := make([]map[string]any, 0)
	for _, a := range approvals {
		data = append(data, a.ToDTO())
	}

	response.Success(c, data)
}

// getInstanceLogs 获取实例的操作日志
func (h *Handler) getInstanceLogs(c *gin.Context) {
	id := c.Param("id")

	logs, err := h.engine.GetInstanceLogs(str2uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取操作日志失败")
		return
	}

	data := make([]map[string]any, 0)
	for _, l := range logs {
		data = append(data, l.ToDTO())
	}

	response.Success(c, data)
}

// resubmitInstance 重新提交实例
func (h *Handler) resubmitInstance(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	if err := h.engine.Resubmit(str2uint(id), userID, userName); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, nil, "重新提交成功")
}

// getPendingTasks 获取我的待办任务
func (h *Handler) getPendingTasks(c *gin.Context) {
	userID := c.GetUint("user_id")

	tasks, err := h.engine.GetPendingTasks(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取待办任务失败")
		return
	}

	data := make([]map[string]any, 0)
	for _, t := range tasks {
		data = append(data, t.ToDTO())
	}

	response.Success(c, data)
}

// getPendingTasksByType 按类型获取待办任务
func (h *Handler) getPendingTasksByType(c *gin.Context) {
	userID := c.GetUint("user_id")
	businessType := c.Param("businessType")

	tasks, err := h.engine.GetPendingTasksByBusiness(userID, businessType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取待办任务失败")
		return
	}

	data := make([]map[string]any, 0)
	for _, t := range tasks {
		data = append(data, t.ToDTO())
	}

	response.Success(c, data)
}

// approveTask 审批通过
func (h *Handler) approveTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	var req struct {
		Remark string                 `json:"remark"`
		Items  []map[string]any       `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 获取待办任务
	var task WorkflowPendingTask
	if err := h.db.First(&task, taskID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "待办任务不存在")
		return
	}

	// 构建额外数据
	extraData := map[string]any{
		"items": req.Items,
	}

	// 处理审批
	if err := h.engine.ProcessApproval(task.InstanceID, userID, userName, ActionApprove, req.Remark, extraData); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, nil, "审批通过")
}

// rejectTask 审批拒绝
func (h *Handler) rejectTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	var req struct {
		Remark string `json:"remark" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 获取待办任务
	var task WorkflowPendingTask
	if err := h.db.First(&task, taskID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "待办任务不存在")
		return
	}

	// 处理拒绝
	if err := h.engine.ProcessApproval(task.InstanceID, userID, userName, ActionReject, req.Remark, nil); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, nil, "已拒绝")
}

// returnTask 退回
func (h *Handler) returnTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	var req struct {
		Remark string `json:"remark" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 获取待办任务
	var task WorkflowPendingTask
	if err := h.db.First(&task, taskID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "待办任务不存在")
		return
	}

	// 处理退回
	if err := h.engine.ProcessApproval(task.InstanceID, userID, userName, ActionReturn, req.Remark, nil); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, nil, "已退回给申请人")
}

// commentTask 评论
func (h *Handler) commentTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	var req struct {
		Remark string `json:"remark" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 获取待办任务
	var task WorkflowPendingTask
	if err := h.db.First(&task, taskID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "待办任务不存在")
		return
	}

	// 处理评论
	if err := h.engine.ProcessApproval(task.InstanceID, userID, userName, ActionComment, req.Remark, nil); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, nil, "评论成功")
}

// 辅助函数
func getStringValue(m map[string]any, key string) string {
	if val, ok := m[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getIntValue(m map[string]any, key string) int {
	if val, ok := m[key]; ok && val != nil {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		}
	}
	return 0
}

func getBoolValue(m map[string]any, key string) bool {
	if val, ok := m[key]; ok && val != nil {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func str2uint(s string) uint {
	id, _ := strconv.Atoi(s)
	return uint(id)
}
