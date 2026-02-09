package appointment

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// RegisterRoutes 注册路由
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	g := rg.Group("/appointments")
	g.Use(jwtpkg.TokenMiddleware())

	service := NewAppointmentService(db)
	workflowService := NewWorkflowService(db)
	calendarService := NewCalendarService(db)

	// 统计
	g.GET("/stats", auth.PermissionMiddleware(db, "appointment_view"), getStats(service))

	// 预约单列表
	g.GET("", auth.PermissionMiddleware(db, "appointment_view"), listAppointments(service))

	// 我的预约
	g.GET("/my", getMyAppointments(service))

	// 待审批
	g.GET("/pending", auth.PermissionMiddleware(db, "appointment_approve"), getPendingApprovals(workflowService))

	// 作业人员的预约列表
	g.GET("/worker/:workerId", auth.PermissionMiddleware(db, "appointment_view"), getWorkerAppointments(service))

	// 搜索
	g.GET("/search", auth.PermissionMiddleware(db, "appointment_view"), searchAppointments(service))

	// 导出
	g.GET("/export", auth.PermissionMiddleware(db, "appointment_view"), exportAppointments(service))

	// 创建预约单
	g.POST("", auth.PermissionMiddleware(db, "appointment_create"), createAppointment(service))

	// 批量创建
	g.POST("/batch", auth.PermissionMiddleware(db, "appointment_create"), batchCreateAppointments(service))

	// 单个预约单详情
	g.GET("/:id", auth.PermissionMiddleware(db, "appointment_view"), getAppointment(service))

	// 更新预约单
	g.PUT("/:id", auth.PermissionMiddleware(db, "appointment_edit"), updateAppointment(service))

	// 删除预约单
	g.DELETE("/:id", auth.PermissionMiddleware(db, "appointment_delete"), deleteAppointment(service))

	// 提交审批
	g.POST("/:id/submit", auth.PermissionMiddleware(db, "appointment_submit"), submitAppointment(service))

	// 启动工作流
	g.POST("/:id/workflow/start", auth.PermissionMiddleware(db, "appointment_approve"), startWorkflow(workflowService))

	// 审批
	g.POST("/:id/approve", auth.PermissionMiddleware(db, "appointment_approve"), approveAppointment(workflowService))

	// 撤回
	g.POST("/:id/recall", auth.PermissionMiddleware(db, "appointment_submit"), recallWorkflow(workflowService))

	// 分配作业人员
	g.POST("/:id/assign", auth.PermissionMiddleware(db, "appointment_assign"), assignWorker(service))

	// 开始作业
	g.POST("/:id/start", auth.PermissionMiddleware(db, "appointment_execute"), startWork(service))

	// 完成作业
	g.POST("/:id/complete", auth.PermissionMiddleware(db, "appointment_execute"), completeAppointment(service))

	// 取消预约
	g.POST("/:id/cancel", auth.PermissionMiddleware(db, "appointment_cancel"), cancelAppointment(service))

	// 审批历史
	g.GET("/:id/approval-history", auth.PermissionMiddleware(db, "appointment_view"), getApprovalHistory(workflowService))

	// 工作流进度
	g.GET("/:id/workflow-progress", auth.PermissionMiddleware(db, "appointment_view"), getWorkflowProgress(workflowService))

	// 当前审批节点
	g.GET("/:id/current-approval", auth.PermissionMiddleware(db, "appointment_view"), getCurrentApproval(workflowService))

	// 批量审批
	g.POST("/batch-approve", auth.PermissionMiddleware(db, "appointment_approve"), batchApprove(workflowService))

	// 日历相关路由
	calendar := g.Group("/calendar")
	{
		// 获取作业人员日历
		calendar.GET("/worker/:workerId", auth.PermissionMiddleware(db, "appointment_view"), getWorkerCalendar(calendarService))

		// 检查可用性
		calendar.POST("/check-availability", auth.PermissionMiddleware(db, "appointment_view"), checkAvailability(calendarService))

		// 批量锁定日历
		calendar.POST("/batch-block", auth.PermissionMiddleware(db, "appointment_manage"), batchBlockCalendar(calendarService))

		// 获取可用作业人员
		calendar.GET("/available-workers", auth.PermissionMiddleware(db, "appointment_view"), getAvailableWorkers(calendarService))

		// 获取日历视图数据
		calendar.GET("/view", auth.PermissionMiddleware(db, "appointment_view"), getCalendarView(calendarService))
	}
}

// getStats 获取统计数据
func getStats(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filterDate *time.Time
		var applicantID *uint

		// 解析日期参数
		if dateStr := c.Query("date"); dateStr != "" {
			if d, err := time.Parse("2006-01-02", dateStr); err == nil {
				filterDate = &d
			}
		}

		// 解析申请人ID
		if applicantStr := c.Query("applicant_id"); applicantStr != "" {
			if id, err := strconv.ParseUint(applicantStr, 10, 32); err == nil {
				aid := uint(id)
				applicantID = &aid
			}
		}

		stats, err := service.GetStats(filterDate, applicantID)
		if err != nil {
			response.InternalError(c, "获取统计数据失败")
			return
		}

		response.Success(c, stats)
	}
}

// listAppointments 获取预约单列表
func listAppointments(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		var isUrgent *bool
		if urgentStr := c.Query("is_urgent"); urgentStr != "" {
			if urgent, err := strconv.ParseBool(urgentStr); err == nil {
				isUrgent = &urgent
			}
		}

		var applicantID, workerID *uint
		if applicantStr := c.Query("applicant_id"); applicantStr != "" {
			if id, err := strconv.ParseUint(applicantStr, 10, 32); err == nil {
				aid := uint(id)
				applicantID = &aid
			}
		}
		if workerStr := c.Query("worker_id"); workerStr != "" {
			if id, err := strconv.ParseUint(workerStr, 10, 32); err == nil {
				wid := uint(id)
				workerID = &wid
			}
		}

		req := AppointmentListRequest{
			Page:        page,
			PageSize:    pageSize,
			Status:      strings.TrimSpace(c.Query("status")),
			IsUrgent:    isUrgent,
			StartDate:   strings.TrimSpace(c.Query("start_date")),
			EndDate:     strings.TrimSpace(c.Query("end_date")),
			ApplicantID: applicantID,
			WorkerID:    workerID,
			WorkType:    strings.TrimSpace(c.Query("work_type")),
		}

		appointments, total, err := service.List(req)
		if err != nil {
			response.InternalError(c, "获取预约单列表失败")
			return
		}

		// 转换为DTO
		items := make([]map[string]any, len(appointments))
		for i, a := range appointments {
			items[i] = a.ToDTO()
		}

		meta := map[string]any{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		}

		response.SuccessWithMeta(c, items, meta)
	}
}

// getMyAppointments 获取我的预约
func getMyAppointments(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := jwtpkg.GetUserID(c)
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		req := AppointmentListRequest{
			Page:        page,
			PageSize:    pageSize,
			Status:      strings.TrimSpace(c.Query("status")),
			StartDate:   strings.TrimSpace(c.Query("start_date")),
			EndDate:     strings.TrimSpace(c.Query("end_date")),
			ApplicantID: &userID,
		}

		appointments, total, err := service.List(req)
		if err != nil {
			response.InternalError(c, "获取预约列表失败")
			return
		}

		items := make([]map[string]any, len(appointments))
		for i, a := range appointments {
			items[i] = a.ToDTO()
		}

		meta := map[string]any{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		}

		response.SuccessWithMeta(c, items, meta)
	}
}

// getPendingApprovals 获取待审批列表
func getPendingApprovals(service *WorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := jwtpkg.GetUserID(c)
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		items, total, err := service.GetPendingApprovals(userID, page, pageSize)
		if err != nil {
			response.InternalError(c, "获取待审批列表失败")
			return
		}

		meta := map[string]any{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		}

		response.SuccessWithMeta(c, items, meta)
	}
}

// getWorkerAppointments 获取作业人员的预约列表
func getWorkerAppointments(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workerIDStr := c.Param("workerId")
		workerID, err := strconv.ParseUint(workerIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的作业人员ID")
			return
		}

		startDateStr := c.DefaultQuery("start_date", time.Now().Format("2006-01-02"))
		endDateStr := c.DefaultQuery("end_date", time.Now().AddDate(0, 0, 30).Format("2006-01-02"))

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.BadRequest(c, "无效的开始日期")
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.BadRequest(c, "无效的结束日期")
			return
		}

		appointments, err := service.GetWorkerAppointments(uint(workerID), startDate, endDate)
		if err != nil {
			response.InternalError(c, "获取预约列表失败")
			return
		}

		items := make([]map[string]any, len(appointments))
		for i, a := range appointments {
			items[i] = a.ToDTO()
		}

		response.Success(c, items)
	}
}

// searchAppointments 搜索预约单
func searchAppointments(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := strings.TrimSpace(c.Query("keyword"))
		if keyword == "" {
			response.BadRequest(c, "搜索关键词不能为空")
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		appointments, total, err := service.SearchByKeyword(keyword, page, pageSize)
		if err != nil {
			response.InternalError(c, "搜索失败")
			return
		}

		items := make([]map[string]any, len(appointments))
		for i, a := range appointments {
			items[i] = a.ToDTO()
		}

		meta := map[string]any{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		}

		response.SuccessWithMeta(c, items, meta)
	}
}

// exportAppointments 导出预约单
func exportAppointments(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idsStr := c.Query("ids")
		if idsStr == "" {
			response.BadRequest(c, "请指定要导出的预约单ID")
			return
		}

		idStrs := strings.Split(idsStr, ",")
		ids := make([]uint, 0, len(idStrs))
		for _, s := range idStrs {
			if id, err := strconv.ParseUint(strings.TrimSpace(s), 10, 32); err == nil {
				ids = append(ids, uint(id))
			}
		}

		if len(ids) == 0 {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		data, err := service.ExportToJSON(ids)
		if err != nil {
			response.InternalError(c, "导出失败")
			return
		}

		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=appointments.json")
		c.Data(200, "application/json", data)
	}
}

// createAppointment 创建预约单
func createAppointment(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := jwtpkg.GetUserID(c)
		userName := jwtpkg.GetUserName(c)
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		var req CreateAppointmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		appointment, err := service.Create(req, userID, userName)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Success(c, appointment.ToDTO())
	}
}

// batchCreateAppointments 批量创建预约单
func batchCreateAppointments(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := jwtpkg.GetUserID(c)
		userName := jwtpkg.GetUserName(c)
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		var req BatchCreateAppointmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		appointments, errs := service.BatchCreate(req, userID, userName)
		if len(errs) > 0 {
			errorMsgs := make([]string, len(errs))
			for i, e := range errs {
				errorMsgs[i] = e.Error()
			}
			response.SuccessWithMessage(c, map[string]any{
				"created": len(appointments),
				"total":   len(req.Appointments),
				"errors":  errorMsgs,
			}, "部分创建失败")
			return
		}

		items := make([]map[string]any, len(appointments))
		for i, a := range appointments {
			items[i] = a.ToDTO()
		}

		response.Success(c, items)
	}
}

// getAppointment 获取预约单详情
func getAppointment(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		appointment, err := service.GetByID(uint(id))
		if err != nil {
			response.NotFound(c, "预约单不存在")
			return
		}

		response.Success(c, appointment.ToDTO())
	}
}

// updateAppointment 更新预约单
func updateAppointment(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		var req UpdateAppointmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		appointment, err := service.Update(uint(id), req)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Success(c, appointment.ToDTO())
	}
}

// deleteAppointment 删除预约单
func deleteAppointment(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		if err := service.Delete(uint(id)); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, nil, "删除成功")
	}
}

// submitAppointment 提交审批
func submitAppointment(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		appointment, err := service.Submit(uint(id))
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Success(c, appointment.ToDTO())
	}
}

// startWorkflow 启动工作流
func startWorkflow(service *WorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		var req struct {
			WorkflowID uint `json:"workflow_id"`
		}
		c.ShouldBindJSON(&req)

		appointment, instance, err := service.StartApprovalWorkflow(uint(id), req.WorkflowID)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Success(c, map[string]any{
			"appointment": appointment.ToDTO(),
			"instance":    instance,
		})
	}
}

// approveAppointment 审批
func approveAppointment(service *WorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		appointment, err := NewAppointmentService(service.db).GetByID(uint(id))
		if err != nil {
			response.NotFound(c, "预约单不存在")
			return
		}

		if appointment.WorkflowInstanceID == nil {
			response.BadRequest(c, "该预约单未启动工作流")
			return
		}

		userID := jwtpkg.GetUserID(c)
		userName := jwtpkg.GetUserName(c)
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		var req ApproveAppointmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		if err := service.ProcessApproval(*appointment.WorkflowInstanceID, userID, userName, req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// 获取更新后的预约单
		appointment, _ = NewAppointmentService(service.db).GetByID(uint(id))
		response.SuccessWithMessage(c, appointment.ToDTO(), "审批成功")
	}
}

// recallWorkflow 撤回工作流
func recallWorkflow(service *WorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		userID := jwtpkg.GetUserID(c)
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		if err := service.RecallWorkflow(uint(id), userID); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, nil, "撤回成功")
	}
}

// assignWorker 分配作业人员
func assignWorker(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		var req AssignWorkerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		// 获取作业人员姓名
		var workerName string
		db := service.db
		db.Table("users").Where("id = ?", req.WorkerID).Pluck("name", &workerName)

		appointment, err := service.AssignWorker(uint(id), req.WorkerID, workerName)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Success(c, appointment.ToDTO())
	}
}

// startWork 开始作业
func startWork(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		appointment, err := service.StartWork(uint(id))
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Success(c, appointment.ToDTO())
	}
}

// completeAppointment 完成作业
func completeAppointment(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		var req CompleteAppointmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		appointment, err := service.Complete(uint(id), req)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Success(c, appointment.ToDTO())
	}
}

// cancelAppointment 取消预约
func cancelAppointment(service *AppointmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		var req struct {
			Reason string `json:"reason"`
		}
		c.ShouldBindJSON(&req)

		appointment, err := service.Cancel(uint(id), req.Reason)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Success(c, appointment.ToDTO())
	}
}

// getApprovalHistory 获取审批历史
func getApprovalHistory(service *WorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		history, err := service.GetApprovalHistory(uint(id))
		if err != nil {
			response.InternalError(c, "获取审批历史失败")
			return
		}

		response.Success(c, history)
	}
}

// getWorkflowProgress 获取工作流进度
func getWorkflowProgress(service *WorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		progress, err := service.GetWorkflowProgress(uint(id))
		if err != nil {
			response.InternalError(c, "获取工作流进度失败")
			return
		}

		response.Success(c, progress)
	}
}

// getCurrentApproval 获取当前审批节点
func getCurrentApproval(service *WorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的预约单ID")
			return
		}

		approval, err := service.GetCurrentApprovalNode(uint(id))
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Success(c, approval)
	}
}

// batchApprove 批量审批
func batchApprove(service *WorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := jwtpkg.GetUserID(c)
		userName := jwtpkg.GetUserName(c)
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		var req struct {
			InstanceIDs []uint `json:"instance_ids" binding:"required"`
			Action      string `json:"action" binding:"required"`
			Comment     string `json:"comment"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		results, errs := service.BatchApprove(req.InstanceIDs, userID, userName, req.Action, req.Comment)

		response.Success(c, map[string]any{
			"results": results,
			"errors":  errs,
		})
	}
}

// 日历相关处理函数

// getWorkerCalendar 获取作业人员日历
func getWorkerCalendar(service *CalendarService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workerIDStr := c.Param("workerId")
		workerID, err := strconv.ParseUint(workerIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的作业人员ID")
			return
		}

		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.BadRequest(c, "无效的开始日期")
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.BadRequest(c, "无效的结束日期")
			return
		}

		calendars, err := service.GetWorkerCalendar(uint(workerID), startDate, endDate)
		if err != nil {
			response.InternalError(c, "获取日历失败")
			return
		}

		items := make([]map[string]any, len(calendars))
		for i, cal := range calendars {
			items[i] = cal.ToDTO()
		}

		response.Success(c, items)
	}
}

// checkAvailability 检查可用性
func checkAvailability(service *CalendarService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AvailabilityCheckRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		workDate, err := time.Parse("2006-01-02", req.WorkDate)
		if err != nil {
			response.BadRequest(c, "无效的作业日期")
			return
		}

		available, reason, err := service.CheckAvailability(req.WorkerID, workDate, req.TimeSlot)
		if err != nil {
			response.InternalError(c, "检查可用性失败")
			return
		}

		response.Success(c, map[string]any{
			"available": available,
			"reason":    reason,
		})
	}
}

// batchBlockCalendar 批量锁定日历
func batchBlockCalendar(service *CalendarService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchBlockCalendarRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		if err := service.BatchBlockCalendar(req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, nil, "日历锁定成功")
	}
}

// getAvailableWorkers 获取可用作业人员
func getAvailableWorkers(service *CalendarService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workDateStr := c.Query("work_date")
		timeSlot := c.Query("time_slot")

		if workDateStr == "" || timeSlot == "" {
			response.BadRequest(c, "请提供作业日期和时间段")
			return
		}

		workDate, err := time.Parse("2006-01-02", workDateStr)
		if err != nil {
			response.BadRequest(c, "无效的作业日期")
			return
		}

		workers, err := service.GetAvailableWorkers(workDate, timeSlot)
		if err != nil {
			response.InternalError(c, "获取可用作业人员失败")
			return
		}

		// 获取作业人员详细信息
		type WorkerInfo struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}
		var workerInfos []WorkerInfo
		db := service.db
		db.Table("users").Where("id IN ?", workers).Find(&workerInfos)

		response.Success(c, workerInfos)
	}
}

// getCalendarView 获取日历视图数据
func getCalendarView(service *CalendarService) gin.HandlerFunc {
	return func(c *gin.Context) {
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")
		workerIDStr := c.Query("worker_id")

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.BadRequest(c, "无效的开始日期")
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.BadRequest(c, "无效的结束日期")
			return
		}

		var workerID *uint
		if workerIDStr != "" {
			if id, err := strconv.ParseUint(workerIDStr, 10, 32); err == nil {
				wid := uint(id)
				workerID = &wid
			}
		}

		appointments, err := service.GetAppointmentsByDateRange(startDate, endDate, workerID)
		if err != nil {
			response.InternalError(c, "获取日历视图数据失败")
			return
		}

		// 构建日历视图响应
		calendarData := make(map[string][]map[string]any)
		for _, apt := range appointments {
			dateKey := apt.WorkDate.Format("2006-01-02")
			if calendarData[dateKey] == nil {
				calendarData[dateKey] = []map[string]any{}
			}
			calendarData[dateKey] = append(calendarData[dateKey], apt.ToDTO())
		}

		response.Success(c, map[string]any{
			"start_date": startDateStr,
			"end_date":   endDateStr,
			"events":     calendarData,
		})
	}
}
