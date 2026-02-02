package progress

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/project"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"gorm.io/gorm"
)

// ProgressHandler 进度管理处理器
type ProgressHandler struct {
	db *gorm.DB
}

// NewProgressHandler 创建处理器
func NewProgressHandler(db *gorm.DB) *ProgressHandler {
	return &ProgressHandler{db: db}
}
// GetProgressList 获取进度任务列表（支持分页和筛选）
func (h *ProgressHandler) GetProgressList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	projectID := c.Query("project_id")
	status := c.Query("status")

	var tasks []Task
	var total int64

	query := h.db.Model(&Task{})

	if projectID != "" {
		// 检查项目是否有子项目
		var childProjects []struct {
			ID uint
		}
		h.db.Table("projects").Where("parent_id = ?", projectID).Find(&childProjects)

		if len(childProjects) > 0 {
			// 有子项目，查询所有子项目的任务
			pid, _ := strconv.ParseUint(projectID, 10, 32)
			allProjectIDs := make([]uint, len(childProjects)+1)
			allProjectIDs[0] = uint(pid)
			for i, child := range childProjects {
				allProjectIDs[i+1] = child.ID
			}
			query = query.Where("project_id IN ?", allProjectIDs)
		} else {
			// 没有子项目，只查询当前项目
			query = query.Where("project_id = ?", projectID)
		}
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tasks).Error; err != nil {
		response.InternalError(c, "查询任务列表失败")
		return
	}

	// 填充项目名称
	type TaskWithProject struct {
		Task
		TaskName    string `json:"task_name"`
		ProjectName string `json:"project_name"`
	}

	var result []TaskWithProject
	for _, task := range tasks {
		var p project.Project
		h.db.First(&p, task.ProjectID)

		result = append(result, TaskWithProject{
			Task:        task,
			TaskName:    task.Name,
			ProjectName: p.Name,
		})
	}

	response.SuccessWithPagination(c, result, int64(page), int64(pageSize), total)
}

// GetTasks 获取项目任务列表
func (h *ProgressHandler) GetTasks(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空")
		return
	}

	var tasks []Task
	if err := h.db.Where("project_id = ?", projectID).Order("sort_order ASC, id ASC").Find(&tasks).Error; err != nil {
		response.InternalError(c, "查询任务失败")
		return
	}

	response.Success(c, tasks)
}

// CreateTask 创建任务
func (h *ProgressHandler) CreateTask(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空")
		return
	}

	var req struct {
		Name        string    `json:"name"`
		Duration    *float64  `json:"duration"`
		StartDate   *string   `json:"start_date"`
		EndDate     *string   `json:"end_date"`
		Progress    float64   `json:"progress"`
		IsMilestone bool      `json:"is_milestone"`
		ParentID    *uint     `json:"parent_id"`
		SortOrder   int       `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "数据格式错误")
		return
	}

	if req.Name == "" {
		response.BadRequest(c, "任务名称不能为空")
		return
	}

	pid, _ := strconv.ParseUint(projectID, 10, 32)

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			startDate = &t
		}
	}
	if req.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			endDate = &t
		}
	}

	task := Task{
		ProjectID:   uint(pid),
		Name:        req.Name,
		Duration:    req.Duration,
		StartDate:   startDate,
		EndDate:     endDate,
		Progress:    req.Progress,
		IsMilestone: req.IsMilestone,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
	}

	if err := h.db.Create(&task).Error; err != nil {
		response.InternalError(c, "创建任务失败")
		return
	}

	// 更新项目进度
	go h.updateProjectProgress(uint(pid))

	response.Created(c, task, "任务创建成功")
}

// CreateProgress 创建进度任务（兼容旧版本API）
func (h *ProgressHandler) CreateProgress(c *gin.Context) {
	var req struct {
		ProjectID   uint    `json:"project_id" binding:"required"`
		ParentID    *uint   `json:"parent_id"`
		TaskName    string  `json:"task_name" binding:"required"`
		StartDate   string  `json:"start_date" binding:"required"`
		EndDate     string  `json:"end_date" binding:"required"`
		Priority    string  `json:"priority"`
		Status      string  `json:"status"`
		Progress    float64 `json:"progress"`
		Responsible string  `json:"responsible"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "数据格式错误: "+err.Error())
		return
	}

	// 验证项目是否存在
	var project struct {
		ID uint
	}
	if err := h.db.Table("projects").Where("id = ?", req.ProjectID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "项目不存在")
			return
		}
		response.InternalError(c, "查询项目失败")
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response.BadRequest(c, "开始日期格式错误")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response.BadRequest(c, "结束日期格式错误")
		return
	}

	// 计算工期（天数）
	duration := endDate.Sub(startDate).Hours() / 24

	// 创建任务
	task := Task{
		ProjectID:   req.ProjectID,
		ParentID:    req.ParentID,
		Name:        req.TaskName,
		StartDate:   &startDate,
		EndDate:     &endDate,
		Duration:    &duration,
		Progress:    req.Progress,
		Priority:    req.Priority,
		Status:      req.Status,
		Responsible: req.Responsible,
		Description: req.Description,
		SortOrder:   0,
	}

	if err := h.db.Create(&task).Error; err != nil {
		response.InternalError(c, "创建任务失败")
		return
	}

	// 更新项目进度
	go h.updateProjectProgress(req.ProjectID)

	response.Created(c, task, "任务创建成功")
}

// UpdateProgress 更新进度任务（兼容旧版本API）
func (h *ProgressHandler) UpdateProgress(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	var task Task
	if err := h.db.First(&task, taskID).Error; err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	var req struct {
		ProjectID   *uint    `json:"project_id"`
		ParentID    *uint    `json:"parent_id"`
		TaskName    *string  `json:"task_name"`
		StartDate   *string  `json:"start_date"`
		EndDate     *string  `json:"end_date"`
		Priority    *string  `json:"priority"`
		Status      *string  `json:"status"`
		Progress    *float64 `json:"progress"`
		Responsible *string  `json:"responsible"`
		Description *string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "数据格式错误")
		return
	}

	updates := make(map[string]any)

	if req.ProjectID != nil {
		updates["project_id"] = *req.ProjectID
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.TaskName != nil {
		updates["name"] = *req.TaskName
	}
	if req.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			updates["start_date"] = t
		}
	}
	if req.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			updates["end_date"] = t
			// 重新计算工期
			if task.StartDate != nil {
				duration := t.Sub(*task.StartDate).Hours() / 24
				updates["duration"] = duration
			}
		}
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Progress != nil {
		updates["progress"] = *req.Progress
	}
	if req.Responsible != nil {
		updates["responsible"] = *req.Responsible
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if err := h.db.Model(&task).Updates(updates).Error; err != nil {
		response.InternalError(c, "更新任务失败")
		return
	}

	// 更新项目进度
	go h.updateProjectProgress(task.ProjectID)

	// Reload to get updated data
	h.db.First(&task, taskID)
	response.SuccessWithMessage(c, task, "任务更新成功")
}

// DeleteProgress 删除进度任务（兼容旧版本API）
func (h *ProgressHandler) DeleteProgress(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	// 检查任务是否存在
	var task Task
	if err := h.db.First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "任务不存在")
			return
		}
		response.InternalError(c, "查询任务失败")
		return
	}

	// 删除任务
	if err := h.db.Delete(&task).Error; err != nil {
		response.InternalError(c, "删除任务失败")
		return
	}

	// 更新项目进度
	go h.updateProjectProgress(task.ProjectID)

	response.SuccessWithMessage(c, nil, "任务删除成功")
}

// ExportProgress 导出进度数据为Excel
func (h *ProgressHandler) ExportProgress(c *gin.Context) {
	projectID := c.Query("project_id")
	status := c.Query("status")

	query := h.db.Model(&Task{})

	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var tasks []Task
	if err := query.Order("id DESC").Find(&tasks).Error; err != nil {
		response.InternalError(c, "查询任务失败")
		return
	}

	// 获取项目名称
	projectMap := make(map[uint]string)
	for _, task := range tasks {
		if _, exists := projectMap[task.ProjectID]; !exists {
			var p project.Project
			h.db.First(&p, task.ProjectID)
			projectMap[task.ProjectID] = p.Name
		}
	}

	// 创建CSV数据
	var csvData [][]string
	csvData = append(csvData, []string{"任务名称", "项目名称", "开始日期", "结束日期", "进度(%)", "状态", "优先级", "负责人", "描述"})

	for _, task := range tasks {
		statusText := "未开始"
		switch task.Status {
		case "in_progress":
			statusText = "进行中"
		case "completed":
			statusText = "已完成"
		case "delayed":
			statusText = "已延期"
		}

		priorityText := "中"
		switch task.Priority {
		case "low":
			priorityText = "低"
		case "high":
			priorityText = "高"
		case "urgent":
			priorityText = "紧急"
		}

		startDate := ""
		if task.StartDate != nil {
			startDate = task.StartDate.Format("2006-01-02")
		}

		endDate := ""
		if task.EndDate != nil {
			endDate = task.EndDate.Format("2006-01-02")
		}

		csvData = append(csvData, []string{
			task.Name,
			projectMap[task.ProjectID],
			startDate,
			endDate,
			fmt.Sprintf("%.0f", task.Progress),
			statusText,
			priorityText,
			task.Responsible,
			task.Description,
		})
	}

	// 生成CSV内容
	var csvContent string
	for _, row := range csvData {
		csvContent += fmt.Sprintf("\"%s\"\n", strings.Join(row, "\",\""))
	}

	// 设置响应头
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=progress_%d.csv", time.Now().Unix()))
	c.String(200, csvContent)
}

// UpdateTask 更新任务
func (h *ProgressHandler) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	var task Task
	if err := h.db.First(&task, taskID).Error; err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "数据格式错误")
		return
	}

	updates := make(map[string]any)

	if v, ok := req["name"].(string); ok {
		updates["name"] = v
	}
	if v, ok := req["duration"].(float64); ok {
		updates["duration"] = v
	}
	if v, ok := req["progress"].(float64); ok {
		updates["progress"] = v
	}
	if v, ok := req["is_milestone"].(bool); ok {
		updates["is_milestone"] = v
	}
	if v, ok := req["sort_order"].(float64); ok {
		updates["sort_order"] = int(v)
	}
	// 处理 parent_id：可以是 float64（有值）或 null（根级别）
	if v, ok := req["parent_id"]; ok {
		if v == nil {
			// null 表示根级别，没有父任务
			updates["parent_id"] = nil
		} else if floatVal, ok := v.(float64); ok {
			// 有父任务
			updates["parent_id"] = uint(floatVal)
		}
	}
	if v, ok := req["start_date"].(string); ok {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			updates["start_date"] = t
		}
	}
	if v, ok := req["end_date"].(string); ok {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			updates["end_date"] = t
		}
	}

	if err := h.db.Model(&task).Updates(updates).Error; err != nil {
		response.InternalError(c, "更新任务失败")
		return
	}

	// 获取更新前的 parent_id（如果存在）
	oldParentID := task.ParentID

	// Reload to get updated data
	h.db.First(&task, taskID)

	// 如果 parent_id 发生变化，更新旧父任务和新父任务的日期
	newParentID := task.ParentID
	if oldParentID != newParentID {
		fmt.Printf("parent_id 变化: %d -> %d\n", getUintPtr(oldParentID), getUintPtr(newParentID))

		// 更新旧父任务的日期（移除这个子任务的影响）
		if oldParentID != nil {
			go h.updateParentTaskDates(*oldParentID)
		}

		// 更新新父任务的日期（包含这个子任务）
		if newParentID != nil {
			go h.updateParentTaskDates(*newParentID)
		}
	} else {
		// 检查日期是否发生变化
		_, startOk := updates["start_date"]
		_, endOk := updates["end_date"]
		if startOk || endOk {
			// 如果日期发生变化，更新父任务
			if newParentID != nil {
				go h.updateParentTaskDates(*newParentID)
			}
		}
	}

	// 更新项目进度
	go h.updateProjectProgress(task.ProjectID)

	response.SuccessWithMessage(c, task, "任务更新成功")
}

// getUintPtr 安全地获取 uint 指针的值
func getUintPtr(ptr *uint) uint {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// updateParentTaskDates 根据子任务更新父任务的开始日期、结束日期和工期
func (h *ProgressHandler) updateParentTaskDates(parentTaskID uint) error {
	// 获取父任务
	var parentTask Task
	if err := h.db.First(&parentTask, parentTaskID).Error; err != nil {
		return err
	}

	// 获取所有子任务
	var childTasks []Task
	if err := h.db.Where("parent_id = ?", parentTaskID).Find(&childTasks).Error; err != nil {
		return err
	}

	// 如果没有子任务，使用默认值
	if len(childTasks) == 0 {
		return nil
	}

	// 计算最早的开始日期和最晚的结束日期
	var earliestStart *time.Time
	var latestEnd *time.Time

	for _, child := range childTasks {
		if child.StartDate != nil {
			if earliestStart == nil || child.StartDate.Before(*earliestStart) {
				earliestStart = child.StartDate
			}
		}
		if child.EndDate != nil {
			if latestEnd == nil || child.EndDate.After(*latestEnd) {
				latestEnd = child.EndDate
			}
		}
	}

	// 计算工期（天数）
	var duration float64
	if earliestStart != nil && latestEnd != nil {
		duration = latestEnd.Sub(*earliestStart).Seconds() / 86400
	}

	// 更新父任务的日期和工期
	updates := make(map[string]interface{})
	if earliestStart != nil {
		updates["start_date"] = earliestStart
	}
	if latestEnd != nil {
		updates["end_date"] = latestEnd
	}
	if duration > 0 {
		updates["duration"] = &duration
	}

	if err := h.db.Model(&parentTask).Updates(updates).Error; err != nil {
		return err
	}

	fmt.Printf("已更新父任务 %d 的日期: start=%v, end=%v, duration=%f\n",
		parentTaskID, earliestStart, latestEnd, duration)

	return nil
}

// DeleteTask 删除任务
func (h *ProgressHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	// 先查询任务以获取project_id
	var task Task
	if err := h.db.First(&task, taskID).Error; err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	if err := h.db.Delete(&Task{}, taskID).Error; err != nil {
		response.InternalError(c, "删除任务失败")
		return
	}

	// 更新项目进度
	go h.updateProjectProgress(task.ProjectID)

	response.SuccessWithMessage(c, nil, "任务删除成功")
}

// GetDependencies 获取任务依赖
func (h *ProgressHandler) GetDependencies(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	var dependencies []TaskDependency
	if err := h.db.Where("task_id = ?", taskID).Find(&dependencies).Error; err != nil {
		response.InternalError(c, "查询依赖失败")
		return
	}

	response.Success(c, dependencies)
}

// AddDependency 添加任务依赖
func (h *ProgressHandler) AddDependency(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	var req struct {
		DependsOn uint   `json:"depends_on"`
		Type      string `json:"type"`
		Lag       int    `json:"lag"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "数据格式错误")
		return
	}

	if req.DependsOn == 0 {
		response.BadRequest(c, "依赖任务ID不能为空")
		return
	}

	tid, _ := strconv.ParseUint(taskID, 10, 32)

	// Check for circular dependency
	if hasCircularDependency(h.db, uint(tid), req.DependsOn) {
		response.BadRequest(c, "无法添加依赖：会产生循环依赖")
		return
	}

	dep := TaskDependency{
		TaskID:    uint(tid),
		DependsOn: req.DependsOn,
		Type:      req.Type,
		Lag:       req.Lag,
	}

	if dep.Type == "" {
		dep.Type = "FS"
	}

	if err := h.db.Create(&dep).Error; err != nil {
		response.InternalError(c, "创建依赖失败")
		return
	}

	response.Created(c, dep, "依赖添加成功")
}

// RemoveDependency 删除任务依赖
func (h *ProgressHandler) RemoveDependency(c *gin.Context) {
	depID := c.Param("id")
	if depID == "" {
		response.BadRequest(c, "依赖ID不能为空")
		return
	}

	if err := h.db.Delete(&TaskDependency{}, depID).Error; err != nil {
		response.InternalError(c, "删除依赖失败")
		return
	}

	response.SuccessWithMessage(c, nil, "依赖删除成功")
}

// UpdateTaskPosition 更新任务位置（网络图可视化）
func (h *ProgressHandler) UpdateTaskPosition(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	var req struct {
		PositionX float64 `json:"position_x"`
		PositionY float64 `json:"position_y"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "数据格式错误")
		return
	}

	updates := map[string]any{
		"position_x": req.PositionX,
		"position_y": req.PositionY,
	}

	if err := h.db.Model(&Task{}).Where("id = ?", taskID).Updates(updates).Error; err != nil {
		response.InternalError(c, "更新位置失败")
		return
	}

	response.SuccessWithMessage(c, nil, "位置更新成功")
}

// GeneratePlanWithAI 使用AI生成进度计划
func (h *ProgressHandler) GeneratePlanWithAI(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空")
		return
	}

	var req struct {
		Mode         string `json:"mode"`         // auto, manual
		TaskCount    int    `json:"task_count"`   // for auto mode
		Requirements string `json:"requirements"` // user requirements
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "数据格式错误")
		return
	}

	// For now, return a simple template plan
	// TODO: Integrate with DeepSeek API in ai_helper.go
	scheduleData := GenerateTemplateSchedule(req.TaskCount)

	response.SuccessWithMessage(c, scheduleData, "计划生成成功")
}

// AggregateChildPlans 聚合子项目计划到主项目
func (h *ProgressHandler) AggregateChildPlans(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空")
		return
	}

	pid, _ := strconv.ParseUint(projectID, 10, 32)

	// Get all child projects
	var childProjects []struct {
		ID   uint
		Name string
	}
	if err := h.db.Table("projects").Where("parent_id = ?", pid).Find(&childProjects).Error; err != nil {
		response.InternalError(c, "查询子项目失败")
		return
	}

	if len(childProjects) == 0 {
		response.BadRequest(c, "该项目没有子项目")
		return
	}

	// Aggregate schedules from all children
	aggregatedData := ScheduleData{
		Nodes:      make(map[string]Node),
		Activities: make(map[string]Activity),
	}

	nodeOffset := 0
	activityOffset := 0

	for _, child := range childProjects {
		var schedule ProjectSchedule
		if err := h.db.Where("project_id = ?", child.ID).First(&schedule).Error; err == nil {
			// Merge nodes and activities with offset to avoid conflicts
			for id, node := range schedule.Data.Nodes {
				newID := fmt.Sprintf("%d_%s", child.ID, id)
				node.ID = newID
				node.Label = fmt.Sprintf("[%s] %s", child.Name, node.Label)
				// Offset positions for visualization
				node.X += float64(nodeOffset * 200)
				node.Y += float64(nodeOffset * 100)
				aggregatedData.Nodes[newID] = node
			}

			for id, activity := range schedule.Data.Activities {
				newID := fmt.Sprintf("%d_%s", child.ID, id)
				activity.ID = newID
				if activity.FromNode != "" {
					activity.FromNode = fmt.Sprintf("%d_%s", child.ID, activity.FromNode)
				}
				if activity.ToNode != "" {
					activity.ToNode = fmt.Sprintf("%d_%s", child.ID, activity.ToNode)
				}
				aggregatedData.Activities[newID] = activity
			}

			nodeOffset += len(schedule.Data.Nodes)
			activityOffset += len(schedule.Data.Activities)
		}
	}

	// Save or update the aggregated schedule
	var existingSchedule ProjectSchedule
	err := h.db.Where("project_id = ?", pid).First(&existingSchedule).Error

	if err == gorm.ErrRecordNotFound {
		existingSchedule = ProjectSchedule{
			ProjectID: uint(pid),
			Data:      aggregatedData,
		}
		if err := h.db.Create(&existingSchedule).Error; err != nil {
			response.InternalError(c, "保存计划失败")
			return
		}
	} else if err != nil {
		response.InternalError(c, "查询计划失败")
		return
	} else {
		existingSchedule.Data = aggregatedData
		if err := h.db.Save(&existingSchedule).Error; err != nil {
			response.InternalError(c, "更新计划失败")
			return
		}
	}

	response.SuccessWithMessage(c, aggregatedData, "子项目计划聚合成功")
}

// GenerateTemplateSchedule 生成模板计划（临时方案，将被AI替代）
func GenerateTemplateSchedule(taskCount int) ScheduleData {
	if taskCount < 3 {
		taskCount = 5
	}

	data := ScheduleData{
		Nodes:      make(map[string]Node),
		Activities: make(map[string]Activity),
	}

	// Create start and end nodes
	data.Nodes["start"] = Node{
		ID:     "start",
		Label:  "开始",
		X:      100,
		Y:      300,
		Type:   "start",
		Number: intPtr(1),
	}

	data.Nodes["end"] = Node{
		ID:     "end",
		Label:  "结束",
		X:      800,
		Y:      300,
		Type:   "end",
		Number: intPtr(taskCount + 2),
	}

	// Create intermediate nodes and activities
	for i := 1; i <= taskCount; i++ {
		nodeID := fmt.Sprintf("node_%d", i)
		data.Nodes[nodeID] = Node{
			ID:     nodeID,
			Label:  fmt.Sprintf("节点%d", i+1),
			X:      100 + float64(i)*150,
			Y:      200 + float64(i%3)*100,
			Type:   "event",
			Number: intPtr(i + 1),
		}

		activityID := fmt.Sprintf("activity_%d", i)
		var fromNode, toNode string
		if i == 1 {
			fromNode = "start"
		} else {
			fromNode = fmt.Sprintf("node_%d", i-1)
		}
		toNode = nodeID

		data.Activities[activityID] = Activity{
			ID:            activityID,
			Name:          fmt.Sprintf("活动%d", i),
			Duration:      5.0,
			FromNode:      fromNode,
			ToNode:        toNode,
			EarliestStart: float64(i - 1) * 5,
			EarliestFinish: float64(i) * 5,
			LatestStart:   float64(i - 1) * 5,
			LatestFinish:  float64(i) * 5,
			TotalFloat:    0,
			FreeFloat:     0,
			IsCritical:    true,
			IsDummy:       false,
		}
	}

	// Final activity to end node
	data.Activities[fmt.Sprintf("activity_%d", taskCount+1)] = Activity{
		ID:            fmt.Sprintf("activity_%d", taskCount+1),
		Name:          "结束活动",
		Duration:      0,
		FromNode:      fmt.Sprintf("node_%d", taskCount),
		ToNode:        "end",
		EarliestStart: float64(taskCount) * 5,
		EarliestFinish: float64(taskCount) * 5,
		LatestStart:   float64(taskCount) * 5,
		LatestFinish:  float64(taskCount) * 5,
		TotalFloat:    0,
		FreeFloat:     0,
		IsCritical:    true,
		IsDummy:       true,
	}

	return data
}

// Helper function to check for circular dependencies
func hasCircularDependency(db *gorm.DB, taskID, dependsOn uint) bool {
	visited := make(map[uint]bool)
	return checkCircular(db, dependsOn, taskID, visited)
}

func checkCircular(db *gorm.DB, current, target uint, visited map[uint]bool) bool {
	if current == target {
		return true
	}
	if visited[current] {
		return false
	}
	visited[current] = true

	var deps []TaskDependency
	db.Where("task_id = ?", current).Find(&deps)

	for _, dep := range deps {
		if checkCircular(db, dep.DependsOn, target, visited) {
			return true
		}
	}

	return false
}

func intPtr(i int) *int {
	return &i
}

// getTaskIDs 从任务数组中提取所有任务ID
func getTaskIDs(tasks []Task) []uint {
	ids := make([]uint, len(tasks))
	for i, task := range tasks {
		ids[i] = task.ID
	}
	return ids
}

// GetProjectSchedule 获取项目进度计划
func (h *ProgressHandler) GetProjectSchedule(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空")
		return
	}

	// 首先尝试从 project_schedules 表获取
	var schedule ProjectSchedule
	err := h.db.Where("project_id = ?", projectID).First(&schedule).Error

	if err == gorm.ErrRecordNotFound {
		// 没有找到 schedule 数据，从 tasks 表动态生成
		scheduleData, err := h.generateScheduleFromTasks(projectID)
		if err != nil {
			response.InternalError(c, "生成进度计划失败")
			return
		}
		response.Success(c, scheduleData)
		return
	}

	if err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, schedule.Data)
}

// GetAllProjectSchedules 获取所有项目进度计划状态
func (h *ProgressHandler) GetAllProjectSchedules(c *gin.Context) {
	var schedules []ProjectSchedule
	if err := h.db.Find(&schedules).Error; err != nil {
		response.InternalError(c, "查询进度计划失败")
		return
	}

	// 转换为 map 格式，key 为项目ID
	result := make(map[string]interface{})
	for _, schedule := range schedules {
		result[fmt.Sprintf("%d", schedule.ProjectID)] = schedule.Data
	}

	// 获取所有有任务的项目ID
	var projectIDs []uint
	h.db.Model(&Task{}).Distinct("project_id").Pluck("project_id", &projectIDs)

	// 为有任务但没有schedule的项目添加空schedule数据
	for _, projectID := range projectIDs {
		key := fmt.Sprintf("%d", projectID)
		if _, exists := result[key]; !exists {
			// 项目有任务但没有schedule记录，动态生成schedule数据
			scheduleData, err := h.generateScheduleFromTasks(fmt.Sprintf("%d", projectID))
			if err == nil {
				result[key] = scheduleData
			}
		}
	}

	response.Success(c, result)
}

// UpdateProjectSchedule 更新项目进度计划
func (h *ProgressHandler) UpdateProjectSchedule(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空")
		return
	}

	// 验证项目是否存在
	var project struct {
		ID uint
	}
	if err := h.db.Table("projects").Where("id = ?", projectID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "项目不存在")
			return
		}
		response.InternalError(c, "查询项目失败")
		return
	}

	// 解析请求体
	var reqData ScheduleData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		response.BadRequest(c, "数据格式错误")
		return
	}

	// 获取当前用户ID（从JWT中获取）
	var uidPtr *uint
	if userID, exists := c.Get("user_id"); exists && userID != nil {
		if uidVal, ok := userID.(uint); ok {
			uidPtr = &uidVal
		}
	}

	// 查找或创建进度计划
	var schedule ProjectSchedule
	err := h.db.Where("project_id = ?", projectID).First(&schedule).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		pid, _ := strconv.ParseUint(projectID, 10, 32)
		schedule = ProjectSchedule{
			ProjectID: uint(pid),
			Data:      reqData,
			CreatedBy: uidPtr,
		}
		if err := h.db.Create(&schedule).Error; err != nil {
			response.InternalError(c, "创建失败")
			return
		}
	} else if err != nil {
		response.InternalError(c, "查询失败")
		return
	} else {
		// 更新现有记录
		schedule.Data = reqData
		schedule.UpdatedBy = uidPtr
		if err := h.db.Save(&schedule).Error; err != nil {
			response.InternalError(c, "更新失败")
			return
		}
	}

	response.SuccessWithMessage(c, schedule.Data, "保存成功")
}

// DeleteProjectSchedule 删除项目进度计划
func (h *ProgressHandler) DeleteProjectSchedule(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空")
		return
	}

	err := h.db.Where("project_id = ?", projectID).Delete(&ProjectSchedule{}).Error
	if err != nil {
		response.InternalError(c, "删除失败")
		return
	}

	response.SuccessOnlyMessage(c, "删除成功")
}

// CheckScheduleExists 检查项目是否有进度计划
func (h *ProgressHandler) CheckScheduleExists(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空")
		return
	}

	var count int64
	err := h.db.Model(&ProjectSchedule{}).Where("project_id = ?", projectID).Count(&count).Error
	if err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	response.SuccessWithMeta(c, map[string]interface{}{
		"exists":   count > 0,
		"has_data": count > 0,
	}, nil)
}

// generateScheduleFromTasks 从 tasks 表生成 ScheduleData 格式的数据
func (h *ProgressHandler) generateScheduleFromTasks(projectID string) (*ScheduleData, error) {
	pid, err := strconv.ParseUint(projectID, 10, 32)
	if err != nil {
		return nil, err
	}

	// 获取项目的所有任务
	var tasks []Task
	if err := h.db.Where("project_id = ?", pid).Order("id ASC").Find(&tasks).Error; err != nil {
		return nil, err
	}

	// 如果没有任务，返回空数据
	if len(tasks) == 0 {
		return &ScheduleData{
			Nodes:      make(map[string]Node),
			Activities: make(map[string]Activity),
		}, nil
		}

	// 创建 nodes 和 activities
	nodes := make(map[string]Node)
	activities := make(map[string]Activity)

	// 为每个任务创建 activity
	nodeNumber := 1
	taskIDMap := make(map[uint]string) // 任务ID到task_X的映射
	for _, task := range tasks {
		taskID := fmt.Sprintf("task_%d", task.ID)
		taskIDMap[task.ID] = taskID

		// 计算时间戳
		var earliestStart, latestFinish float64
		if task.StartDate != nil {
			earliestStart = float64(task.StartDate.Unix())
		}
		if task.EndDate != nil {
			latestFinish = float64(task.EndDate.Unix())
		}

		var duration float64
		if task.Duration != nil {
			duration = *task.Duration * 86400 // 转换为秒
		} else if task.StartDate != nil && task.EndDate != nil {
			duration = float64(task.EndDate.Sub(*task.StartDate).Seconds())
		}

		// 创建 activity
		activities[taskID] = Activity{
			ID:            taskID,
			Name:          task.Name,
			Duration:      duration / 86400, // 转换为天
			FromNode:      "",
			ToNode:        "",
			EarliestStart: earliestStart,
			EarliestFinish: earliestStart + duration,
			LatestStart:   earliestStart,
			LatestFinish:  latestFinish,
			TotalFloat:    0,
			FreeFloat:     0,
			IsCritical:    false,
			IsDummy:       false,
			Breakpoints:   nil,
			Progress:      task.Progress,
			Status:        task.Status,
			Priority:      task.Priority,
			Predecessors:  []int{},   // 初始化前置任务数组
			Successors:    []int{},   // 初始化后置任务数组
			SortOrder:     int(task.SortOrder),
		}
	}

	// 加载任务依赖关系
	var dependencies []TaskDependency
	if err := h.db.Where("task_id IN ?", getTaskIDs(tasks)).Find(&dependencies).Error; err == nil {
		// 构建前置任务映射
		predecessorMap := make(map[string][]int)
		successorMap := make(map[string][]int)

		for _, dep := range dependencies {
			taskIDStr := taskIDMap[dep.TaskID]
			if dep.DependsOn > 0 {
				if depTaskIDStr, ok := taskIDMap[dep.DependsOn]; ok {
					// 将 depends_on 转换为数字ID
					parts := strings.TrimPrefix(depTaskIDStr, "task_")
					if depID, err := strconv.Atoi(parts); err == nil {
						predecessorMap[taskIDStr] = append(predecessorMap[taskIDStr], depID)
						successorMap[depTaskIDStr] = append(successorMap[depTaskIDStr], int(dep.TaskID))
					}
				}
			}
		}

		// 将依赖关系填充到 activities
		for taskIDStr, preds := range predecessorMap {
			if activity, ok := activities[taskIDStr]; ok {
				activity.Predecessors = preds
				activities[taskIDStr] = activity
			}
		}
		for taskIDStr, succs := range successorMap {
			if activity, ok := activities[taskIDStr]; ok {
				activity.Successors = succs
				activities[taskIDStr] = activity
			}
		}
	}

	// 加载任务资源分配
	var taskResources []TaskResource
	if err := h.db.Where("task_id IN ?", getTaskIDs(tasks)).Find(&taskResources).Error; err == nil {
		// 按任务ID分组资源
		taskResourcesMap := make(map[uint][]Resource)
		for _, tr := range taskResources {
			var resource Resource
			if err := h.db.First(&resource, tr.ResourceID).Error; err == nil {
				taskResourcesMap[tr.TaskID] = append(taskResourcesMap[tr.TaskID], Resource{
					ID:          resource.ID,
					Name:        resource.Name,
					Type:        resource.Type,
					Unit:        resource.Unit,
					Quantity:    resource.Quantity,
					CostPerUnit: resource.CostPerUnit,
					Color:       resource.Color,
				})
			}
		}

		// 将资源填充到 activities
		for _, task := range tasks {
			taskIDStr := taskIDMap[task.ID]
			if activity, ok := activities[taskIDStr]; ok {
				activity.Resources = taskResourcesMap[task.ID]
				activities[taskIDStr] = activity
			}
		}
	}

	// 创建起点和终点节点
	nodes["start"] = Node{
		ID:           "start",
		Label:        "开始",
		X:            100,
		Y:            300,
		Type:         "start",
		Number:       intPtr(0),
		EarliestTime: float64Ptr(0),
		LatestTime:   float64Ptr(0),
		Activities:   nil,
	}

	nodes["end"] = Node{
		ID:           "end",
		Label:        "结束",
		X:            800,
		Y:            300,
		Type:         "end",
		Number:       intPtr(nodeNumber),
		EarliestTime: float64Ptr(0),
		LatestTime:   float64Ptr(0),
		Activities:   nil,
	}

	return &ScheduleData{
		Nodes:      nodes,
		Activities: activities,
	}, nil
}

func float64Ptr(f float64) *float64 {
	return &f
}

// updateProjectProgress 根据任务进度更新项目进度
func (h *ProgressHandler) updateProjectProgress(projectID uint) error {
	// 获取项目的所有任务
	var tasks []Task
	if err := h.db.Where("project_id = ?", projectID).Find(&tasks).Error; err != nil {
		return err
	}

	// 计算平均进度
	var totalProgress float64
	var taskCount int

	for _, task := range tasks {
		totalProgress += task.Progress
		taskCount++
	}

	var progressPercentage float64
	if taskCount > 0 {
		progressPercentage = totalProgress / float64(taskCount)
	}

	// 更新项目的进度百分比
	if err := h.db.Table("projects").
		Where("id = ?", projectID).
		Update("progress_percentage", progressPercentage).Error; err != nil {
		return err
	}

	return nil
}

// ==================== 子任务进度计算 ====================

// CalculateParentTaskProgress 根据子任务计算父任务进度
// 按工期加权平均计算子任务进度
func (h *ProgressHandler) CalculateParentTaskProgress(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	tid, _ := strconv.ParseUint(taskID, 10, 32)

	if err := h.CalculateParentTaskProgressInternal(uint(tid)); err != nil {
		response.InternalError(c, "计算父任务进度失败")
		return
	}

	response.SuccessWithMessage(c, nil, "父任务进度已更新")
}

// CalculateParentTaskProgressInternal 计算父任务进度的内部方法
func (h *ProgressHandler) CalculateParentTaskProgressInternal(parentTaskID uint) error {
	// 1. 查询所有直接子任务
	var childTasks []Task
	if err := h.db.Where("parent_id = ?", parentTaskID).Find(&childTasks).Error; err != nil {
		return err
	}

	// 如果没有子任务，返回
	if len(childTasks) == 0 {
		return nil
	}

	// 2. 计算加权平均进度（按工期加权）
	var totalWeightedProgress float64
	var totalWeight float64

	for _, child := range childTasks {
		var duration float64
		if child.Duration != nil {
			duration = *child.Duration
		} else if child.StartDate != nil && child.EndDate != nil {
			duration = child.EndDate.Sub(*child.StartDate).Hours() / 24
		}

		// 最小权重为1天
		if duration < 1 {
			duration = 1
		}

		totalWeightedProgress += child.Progress * duration
		totalWeight += duration
	}

	var parentProgress float64
	if totalWeight > 0 {
		parentProgress = totalWeightedProgress / totalWeight
	}

	// 3. 更新父任务进度
	if err := h.db.Model(&Task{}).Where("id = ?", parentTaskID).Update("progress", parentProgress).Error; err != nil {
		return err
	}

	// 4. 递归向上计算（如果父任务也有父任务）
	var parentTask Task
	if err := h.db.First(&parentTask, parentTaskID).Error; err != nil {
		return err
	}

	if parentTask.ParentID != nil {
		return h.CalculateParentTaskProgressInternal(*parentTask.ParentID)
	}

	return nil
}

// UpdateTaskParentProgress 更新任务及其所有父任务的进度
func (h *ProgressHandler) UpdateTaskParentProgress(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	tid, _ := strconv.ParseUint(taskID, 10, 32)

	if err := h.CalculateParentTaskProgressInternal(uint(tid)); err != nil {
		response.InternalError(c, "更新父任务进度失败")
		return
	}

	response.SuccessWithMessage(c, nil, "父任务进度已更新")
}

// ==================== 资源管理 ====================

// GetProjectResources 获取项目资源列表
func (h *ProgressHandler) GetProjectResources(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空")
		return
	}

	var resources []Resource
	if err := h.db.Where("project_id = ? AND is_active = ?", projectID, true).Order("type, name").Find(&resources).Error; err != nil {
		response.InternalError(c, "查询资源失败")
		return
	}

	response.Success(c, resources)
}

// CreateResource 创建资源
func (h *ProgressHandler) CreateResource(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空")
		return
	}

	var req struct {
		Name        string  `json:"name" binding:"required"`
		Type        string  `json:"type" binding:"required"`
		Unit        string  `json:"unit"`
		Quantity    float64 `json:"quantity"`
		CostPerUnit float64 `json:"cost_per_unit"`
		Color       string  `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "数据格式错误")
		return
	}

	pid, _ := strconv.ParseUint(projectID, 10, 32)

	resource := Resource{
		ProjectID:   uint(pid),
		Name:        req.Name,
		Type:        req.Type,
		Unit:        req.Unit,
		Quantity:    req.Quantity,
		CostPerUnit: req.CostPerUnit,
		Color:       req.Color,
		IsActive:    true,
	}

	if err := h.db.Create(&resource).Error; err != nil {
		response.InternalError(c, "创建资源失败")
		return
	}

	response.Created(c, resource, "资源创建成功")
}

// UpdateResource 更新资源
func (h *ProgressHandler) UpdateResource(c *gin.Context) {
	resourceID := c.Param("resourceId")
	if resourceID == "" {
		response.BadRequest(c, "资源ID不能为空")
		return
	}

	var req struct {
		Name        *string  `json:"name"`
		Type        *string  `json:"type"`
		Unit        *string  `json:"unit"`
		Quantity    *float64 `json:"quantity"`
		CostPerUnit *float64 `json:"cost_per_unit"`
		Color       *string  `json:"color"`
		IsActive    *bool    `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "数据格式错误")
		return
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Unit != nil {
		updates["unit"] = *req.Unit
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.CostPerUnit != nil {
		updates["cost_per_unit"] = *req.CostPerUnit
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.db.Model(&Resource{}).Where("id = ?", resourceID).Updates(updates).Error; err != nil {
		response.InternalError(c, "更新资源失败")
		return
	}

	// 重新查询获取更新后的数据
	var resource Resource
	h.db.First(&resource, resourceID)
	response.SuccessWithMessage(c, resource, "资源更新成功")
}

// DeleteResource 删除资源（软删除）
func (h *ProgressHandler) DeleteResource(c *gin.Context) {
	resourceID := c.Param("resourceId")
	if resourceID == "" {
		response.BadRequest(c, "资源ID不能为空")
		return
	}

	if err := h.db.Model(&Resource{}).Where("id = ?", resourceID).Update("is_active", false).Error; err != nil {
		response.InternalError(c, "删除资源失败")
		return
	}

	response.SuccessWithMessage(c, nil, "资源已删除")
}

// ==================== 任务资源分配 ====================

// GetTaskResources 获取任务资源分配
func (h *ProgressHandler) GetTaskResources(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	var taskResources []struct {
		TaskResource
		ResourceName string `json:"resource_name"`
		ResourceType string `json:"type"`
		Unit         string `json:"unit"`
		Color        string `json:"color"`
	}

	if err := h.db.Table("task_resources tr").
		Select("tr.*, r.name as resource_name, r.type, r.unit, r.color").
		Joins("INNER JOIN resources r ON tr.resource_id = r.id").
		Where("tr.task_id = ?", taskID).
		Find(&taskResources).Error; err != nil {
		response.InternalError(c, "查询任务资源失败")
		return
	}

	response.Success(c, taskResources)
}

// AllocateTaskResource 分配资源给任务
func (h *ProgressHandler) AllocateTaskResource(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	var req struct {
		ResourceID uint    `json:"resource_id" binding:"required"`
		Quantity   float64 `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "数据格式错误")
		return
	}

	// 检查资源是否存在且激活
	var resource Resource
	if err := h.db.Where("id = ? AND is_active = ?", req.ResourceID, true).First(&resource).Error; err != nil {
		response.NotFound(c, "资源不存在或已停用")
		return
	}

	// 检查是否已存在相同的分配
	var existing TaskResource
	err := h.db.Where("task_id = ? AND resource_id = ?", taskID, req.ResourceID).First(&existing).Error

	if err == nil {
		// 更新现有分配
		if err := h.db.Model(&existing).Update("quantity", req.Quantity).Error; err != nil {
			response.InternalError(c, "更新资源分配失败")
			return
		}
		existing.Quantity = req.Quantity
		response.SuccessWithMessage(c, existing, "资源分配已更新")
		return
	}

	if err != gorm.ErrRecordNotFound {
		response.InternalError(c, "查询资源分配失败")
		return
	}

	// 创建新分配
	tid, _ := strconv.ParseUint(taskID, 10, 32)
	taskResource := TaskResource{
		TaskID:     uint(tid),
		ResourceID: req.ResourceID,
		Quantity:   req.Quantity,
	}

	if err := h.db.Create(&taskResource).Error; err != nil {
		response.InternalError(c, "分配资源失败")
		return
	}

	response.Created(c, taskResource, "资源分配成功")
}

// RemoveTaskResource 移除任务资源
func (h *ProgressHandler) RemoveTaskResource(c *gin.Context) {
	taskID := c.Param("id")
	resourceID := c.Param("resourceId")

	if taskID == "" || resourceID == "" {
		response.BadRequest(c, "任务ID和资源ID不能为空")
		return
	}

	if err := h.db.Where("task_id = ? AND resource_id = ?", taskID, resourceID).Delete(&TaskResource{}).Error; err != nil {
		response.InternalError(c, "移除资源失败")
		return
	}

	response.SuccessWithMessage(c, nil, "资源已移除")
}

// ==================== 依赖关系可视化创建 ====================

// CreateDependencyVisual 可视化创建依赖关系（从前端API调用）
func (h *ProgressHandler) CreateDependencyVisual(c *gin.Context) {
	fromTaskID := c.Param("fromId")
	toTaskID := c.Param("toId")

	if fromTaskID == "" || toTaskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	fromID, err1 := strconv.ParseUint(fromTaskID, 10, 32)
	toID, err2 := strconv.ParseUint(toTaskID, 10, 32)

	if err1 != nil || err2 != nil {
		response.BadRequest(c, "任务ID格式错误")
		return
	}

	var req struct {
		Type string `json:"type"`
		Lag  int    `json:"lag"`
	}

	// 绑定可选的请求体
	c.ShouldBindJSON(&req)

	fmt.Printf("创建依赖关系: fromID=%d, toID=%d, type=%s, lag=%d\n", fromID, toID, req.Type, req.Lag)

	// 检查循环依赖
	if hasCircularDependency(h.db, uint(fromID), uint(toID)) {
		fmt.Println("检测到循环依赖")
		response.BadRequest(c, "无法添加依赖：会产生循环依赖")
		return
	}

	// 检查是否已存在
	var existing TaskDependency
	err := h.db.Where("task_id = ? AND depends_on = ?", uint(toID), uint(fromID)).First(&existing).Error
	if err == nil {
		// 已存在，更新类型和延迟
		existing.Type = req.Type
		if existing.Type == "" {
			existing.Type = "FS"
		}
		existing.Lag = req.Lag
		if err := h.db.Save(&existing).Error; err != nil {
			response.InternalError(c, "更新依赖关系失败")
			return
		}
		fmt.Println("依赖关系已更新")
		response.SuccessWithMessage(c, existing, "依赖关系已更新")
		return
	}
	if err != gorm.ErrRecordNotFound {
		response.InternalError(c, "检查依赖关系失败")
		return
	}

	// 创建依赖
	dep := TaskDependency{
		TaskID:    uint(toID),
		DependsOn: uint(fromID),
		Type:      req.Type,
		Lag:       req.Lag,
	}

	if dep.Type == "" {
		dep.Type = "FS"
	}

	if err := h.db.Create(&dep).Error; err != nil {
		response.InternalError(c, "创建依赖失败")
		return
	}

	response.Created(c, dep, "依赖创建成功")
}
