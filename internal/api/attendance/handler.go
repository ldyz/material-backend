package attendance

import (
	"log"
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
	g := rg.Group("/attendance")
	g.Use(jwtpkg.TokenMiddleware())

	service := NewService(db)

	// 移动端接口
	{
		// 获取今日待打卡任务
		g.GET("/today-appointments", getTodayAppointments(service))

		// 打卡
		g.POST("/clock-in", clockIn(service))

		// 获取我的打卡记录
		g.GET("/my-records", getMyRecords(service))

		// 获取我的月度汇总
		g.GET("/my-summary", getMySummary(service))

		// 获取打卡日历统计
		g.GET("/calendar-statistics", getCalendarStatistics(service))
	}

	// Web端接口（需要权限）
	{
		// 获取打卡记录列表
		g.GET("/records", auth.PermissionMiddleware(db, "attendance_view"), getRecords(service))

		// 获取打卡记录详情
		g.GET("/records/:id", auth.PermissionMiddleware(db, "attendance_view"), getRecordByID(service))

		// 确认打卡记录
		g.POST("/records/:id/confirm", auth.PermissionMiddleware(db, "attendance_confirm"), confirmRecord(service))

		// 驳回打卡记录
		g.POST("/records/:id/reject", auth.PermissionMiddleware(db, "attendance_confirm"), rejectRecord(service))

		// 获取月度汇总
		g.GET("/monthly-summary", auth.PermissionMiddleware(db, "attendance_view"), getMonthlySummary(service))

		// 生成月度汇总
		g.POST("/generate-monthly", auth.PermissionMiddleware(db, "attendance_manage"), generateMonthly(service))

		// 确认月度汇总
		g.POST("/monthly-summary/:id/confirm", auth.PermissionMiddleware(db, "attendance_manage"), confirmMonthlySummary(service))

		// 获取统计数据
		g.GET("/statistics", auth.PermissionMiddleware(db, "attendance_view"), getStatistics(service))
	}
}

// getTodayAppointments 获取今日待打卡任务
func getTodayAppointments(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("current_user_id")
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		appointments, err := service.GetTodayAppointments(uint(userID))
		if err != nil {
			log.Printf("获取今日待打卡任务失败: %v", err)
			response.InternalError(c, "获取今日待打卡任务失败")
			return
		}

		response.Success(c, appointments)
	}
}

// clockIn 打卡
func clockIn(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("current_user_id")
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		var req ClockInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		record, err := service.ClockIn(uint(userID), req)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, record.ToDTO(), "打卡成功")
	}
}

// getMyRecords 获取我的打卡记录
func getMyRecords(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("current_user_id")
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		req := RecordListRequest{
			Page:           page,
			PageSize:       pageSize,
			AttendanceType: strings.TrimSpace(c.Query("attendance_type")),
			Status:         strings.TrimSpace(c.Query("status")),
			StartDate:      strings.TrimSpace(c.Query("start_date")),
			EndDate:        strings.TrimSpace(c.Query("end_date")),
		}

		records, total, err := service.GetMyRecords(uint(userID), req)
		if err != nil {
			response.InternalError(c, "获取打卡记录失败")
			return
		}

		items := make([]map[string]any, len(records))
		for i, r := range records {
			items[i] = r.ToDTO()
		}

		meta := map[string]any{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		}

		response.SuccessWithMeta(c, items, meta)
	}
}

// getMySummary 获取我的月度汇总
func getMySummary(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("current_user_id")
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
		month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))

		// 默认当前年月
		if year == 0 || month == 0 {
			now := time.Now()
			year = now.Year()
			month = int(now.Month())
		}

		summary, err := service.GetMyMonthlySummary(uint(userID), year, month)
		if err != nil {
			response.InternalError(c, "获取月度汇总失败")
			return
		}

		response.Success(c, summary.ToDTO())
	}
}

// getRecords 获取打卡记录列表（管理员）
func getRecords(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		var userID *uint
		if userStr := c.Query("user_id"); userStr != "" {
			if id, err := strconv.ParseUint(userStr, 10, 32); err == nil {
				uid := uint(id)
				userID = &uid
			}
		}

		req := RecordListRequest{
			Page:           page,
			PageSize:       pageSize,
			UserID:         userID,
			AttendanceType: strings.TrimSpace(c.Query("attendance_type")),
			Status:         strings.TrimSpace(c.Query("status")),
			StartDate:      strings.TrimSpace(c.Query("start_date")),
			EndDate:        strings.TrimSpace(c.Query("end_date")),
		}

		records, total, err := service.GetRecords(req)
		if err != nil {
			response.InternalError(c, "获取打卡记录失败")
			return
		}

		items := make([]map[string]any, len(records))
		for i, r := range records {
			items[i] = r.ToDTO()
		}

		meta := map[string]any{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		}

		response.SuccessWithMeta(c, items, meta)
	}
}

// getRecordByID 获取打卡记录详情
func getRecordByID(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的记录ID")
			return
		}

		record, err := service.GetRecordByID(uint(id))
		if err != nil {
			response.NotFound(c, "打卡记录不存在")
			return
		}

		response.Success(c, record.ToDTO())
	}
}

// confirmRecord 确认打卡记录
func confirmRecord(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的记录ID")
			return
		}

		userID := c.GetInt64("current_user_id")
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		var req ConfirmRecordRequest
		c.ShouldBindJSON(&req)

		if err := service.ConfirmRecord(uint(id), uint(userID), req.Remark); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, nil, "确认成功")
	}
}

// rejectRecord 驳回打卡记录
func rejectRecord(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的记录ID")
			return
		}

		userID := c.GetInt64("current_user_id")
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		var req RejectRecordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		if err := service.RejectRecord(uint(id), uint(userID), req.Reason); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, nil, "驳回成功")
	}
}

// getMonthlySummary 获取月度汇总
func getMonthlySummary(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		year, _ := strconv.Atoi(c.Query("year"))
		month, _ := strconv.Atoi(c.Query("month"))

		if year == 0 || month == 0 {
			now := time.Now()
			year = now.Year()
			month = int(now.Month())
		}

		var userID *uint
		if userStr := c.Query("user_id"); userStr != "" {
			if id, err := strconv.ParseUint(userStr, 10, 32); err == nil {
				uid := uint(id)
				userID = &uid
			}
		}

		summaries, err := service.GetMonthlySummary(year, month, userID)
		if err != nil {
			response.InternalError(c, "获取月度汇总失败")
			return
		}

		items := make([]map[string]any, len(summaries))
		for i, s := range summaries {
			items[i] = s.ToDTO()
		}

		response.Success(c, items)
	}
}

// generateMonthly 生成月度汇总
func generateMonthly(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MonthlySummaryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		if err := service.GenerateMonthlySummary(req.Year, req.Month); err != nil {
			response.InternalError(c, "生成月度汇总失败")
			return
		}

		response.SuccessWithMessage(c, nil, "生成成功")
	}
}

// confirmMonthlySummary 确认月度汇总
func confirmMonthlySummary(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的汇总ID")
			return
		}

		userID := c.GetInt64("current_user_id")
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		if err := service.ConfirmMonthlySummary(uint(id), uint(userID)); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, nil, "确认成功")
	}
}

// getStatistics 获取统计数据
func getStatistics(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := service.GetStatistics()
		if err != nil {
			response.InternalError(c, "获取统计数据失败")
			return
		}

		response.Success(c, stats)
	}
}

// getCalendarStatistics 获取打卡日历统计
func getCalendarStatistics(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("current_user_id")
		if userID == 0 {
			response.Unauthorized(c, "未授权")
			return
		}

		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		if startDate == "" || endDate == "" {
			response.BadRequest(c, "请提供开始和结束日期")
			return
		}

		statistics, err := service.GetCalendarStatistics(uint(userID), startDate, endDate)
		if err != nil {
			response.InternalError(c, "获取日历统计失败")
			return
		}

		response.Success(c, statistics)
	}
}
