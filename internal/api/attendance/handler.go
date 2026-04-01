package attendance

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"github.com/yourorg/material-backend/backend/internal/pkg/logger"
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

		// 导出打卡记录 Excel
		g.GET("/records/export", auth.PermissionMiddleware(db, "attendance_view"), exportRecords(service))

		// 导出打卡照片 ZIP
		g.GET("/records/export-photos", auth.PermissionMiddleware(db, "attendance_view"), exportPhotos(service))

		// 获取月度汇总
		g.GET("/monthly-summary", auth.PermissionMiddleware(db, "attendance_view"), getMonthlySummary(service))

		// 生成月度汇总
		g.POST("/generate-monthly", auth.PermissionMiddleware(db, "attendance_manage"), generateMonthly(service))

		// 确认月度汇总
		g.POST("/monthly-summary/:id/confirm", auth.PermissionMiddleware(db, "attendance_manage"), confirmMonthlySummary(service))

		// 获取统计数据
		g.GET("/statistics", auth.PermissionMiddleware(db, "attendance_view"), getStatistics(service))

		// 统计接口
		g.GET("/statistics/daily", auth.PermissionMiddleware(db, "attendance_view"), getDailyStatistics(service))
		g.GET("/statistics/by-user", auth.PermissionMiddleware(db, "attendance_view"), getUserStatistics(service))
		g.GET("/statistics/by-task", auth.PermissionMiddleware(db, "attendance_view"), getTaskStatistics(service))
		g.GET("/statistics/by-project", auth.PermissionMiddleware(db, "attendance_view"), getProjectStatistics(service))
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
			logger.Warnf("获取今日待打卡任务失败: %v", err)
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

		var appointmentID *uint
		if aptStr := c.Query("appointment_id"); aptStr != "" {
			if id, err := strconv.ParseUint(aptStr, 10, 32); err == nil {
				aid := uint(id)
				appointmentID = &aid
			}
		}

		var projectID *uint
		if projStr := c.Query("project_id"); projStr != "" {
			if id, err := strconv.ParseUint(projStr, 10, 32); err == nil {
				pid := uint(id)
				projectID = &pid
			}
		}

		req := RecordListRequest{
			Page:           page,
			PageSize:       pageSize,
			UserID:         userID,
			AppointmentID:  appointmentID,
			ProjectID:      projectID,
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

// getDailyStatistics 按日期统计
func getDailyStatistics(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		if startDate == "" || endDate == "" {
			// 默认最近30天
			now := time.Now()
			endDate = now.Format("2006-01-02")
			startDate = now.AddDate(0, 0, -30).Format("2006-01-02")
		}

		statistics, err := service.GetDailyStatistics(startDate, endDate)
		if err != nil {
			response.InternalError(c, "获取统计数据失败")
			return
		}

		response.Success(c, statistics)
	}
}

// getUserStatistics 按人员统计
func getUserStatistics(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		if startDate == "" || endDate == "" {
			// 默认当月
			now := time.Now()
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
			endDate = now.Format("2006-01-02")
		}

		statistics, err := service.GetUserStatistics(startDate, endDate)
		if err != nil {
			response.InternalError(c, "获取统计数据失败")
			return
		}

		response.Success(c, statistics)
	}
}

// getTaskStatistics 按任务统计
func getTaskStatistics(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		if startDate == "" || endDate == "" {
			// 默认最近30天
			now := time.Now()
			endDate = now.Format("2006-01-02")
			startDate = now.AddDate(0, 0, -30).Format("2006-01-02")
		}

		statistics, err := service.GetTaskStatistics(startDate, endDate)
		if err != nil {
			response.InternalError(c, "获取统计数据失败")
			return
		}

		response.Success(c, statistics)
	}
}

// getProjectStatistics 按项目统计
func getProjectStatistics(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		if startDate == "" || endDate == "" {
			// 默认当月
			now := time.Now()
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
			endDate = now.Format("2006-01-02")
		}

		statistics, err := service.GetProjectStatistics(startDate, endDate)
		if err != nil {
			response.InternalError(c, "获取统计数据失败")
			return
		}

		response.Success(c, statistics)
	}
}

// exportRecords 导出打卡记录 Excel
func exportRecords(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := parseExportRequest(c)

		records, err := service.GetRecordsForExport(req)
		if err != nil {
			response.InternalError(c, "获取打卡记录失败")
			return
		}

		if len(records) == 0 {
			response.BadRequest(c, "没有可导出的数据")
			return
		}

		// 创建 Excel 文件
		f := excelize.NewFile()
		defer f.Close()

		sheet := "打卡记录"
		f.SetSheetName("Sheet1", sheet)

		// 设置表头
		headers := []string{"序号", "打卡人", "打卡类型", "打卡时间", "关联任务", "工作内容", "打卡位置", "加班小时", "状态", "确认人", "确认时间"}
		for i, header := range headers {
			cell := fmt.Sprintf("%s1", string(rune('A'+i)))
			f.SetCellValue(sheet, cell, header)
		}

		// 设置列宽
		colWidths := []float64{6, 10, 10, 18, 14, 30, 20, 10, 8, 10, 18}
		for i, width := range colWidths {
			col := string(rune('A' + i))
			f.SetColWidth(sheet, col, col, width)
		}

		// 填充数据
		for idx, record := range records {
			row := idx + 2
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), idx+1)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), record.UserName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), GetAttendanceTypeLabel(record.AttendanceType))
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), record.ClockInTime.Format("2006-01-02 15:04:05"))
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), record.AppointmentNo)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), record.WorkContent)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), record.ClockInLocation)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), record.OvertimeHours)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), GetStatusLabel(record.Status))
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), record.ConfirmedByName)
			if record.ConfirmedAt != nil {
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), record.ConfirmedAt.Format("2006-01-02 15:04:05"))
			}
		}

		// 设置响应头
		filename := fmt.Sprintf("打卡记录_%s.xlsx", time.Now().Format("20060102150405"))
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

		// 写入响应
		buf := new(bytes.Buffer)
		if err := f.Write(buf); err != nil {
			response.InternalError(c, "生成Excel失败")
			return
		}

		c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	}
}

// exportPhotos 导出打卡照片 ZIP
func exportPhotos(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := parseExportRequest(c)

		photos, records, err := service.GetPhotosForExport(req)
		if err != nil {
			response.InternalError(c, "获取照片信息失败")
			return
		}

		if len(photos) == 0 {
			response.BadRequest(c, "没有可导出的照片")
			return
		}

		// 创建 ZIP 缓冲区
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		successCount := 0
		// 读取本地照片文件并添加到 ZIP
		for _, photo := range photos {
			photoURL := photo["url"]

			// 处理相对路径，转换为本地文件路径
			localPath := photoURL
			if strings.HasPrefix(photoURL, "/uploads/") {
				localPath = "static" + photoURL
			} else if strings.HasPrefix(photoURL, "/static/") {
				localPath = strings.TrimPrefix(photoURL, "/")
			}

			// 读取本地文件
			fileData, err := os.ReadFile(localPath)
			if err != nil {
				logger.Warnf("读取照片文件失败: %s, 错误: %v", localPath, err)
				continue
			}

			// 创建 ZIP 文件条目
			writer, err := zipWriter.Create(photo["filename"])
			if err != nil {
				continue
			}

			// 写入照片数据
			_, err = writer.Write(fileData)
			if err != nil {
				continue
			}
			successCount++
		}

		// 添加说明文件
		readmeWriter, err := zipWriter.Create("导出说明.txt")
		if err == nil {
			readmeContent := fmt.Sprintf("打卡照片导出\n导出时间: %s\n共 %d 条记录, %d 张照片成功导出\n",
				time.Now().Format("2006-01-02 15:04:05"),
				len(records),
				successCount)
			readmeWriter.Write([]byte(readmeContent))
		}

		zipWriter.Close()

		if successCount == 0 {
			response.BadRequest(c, "照片文件读取失败，请检查文件是否存在")
			return
		}

		// 设置响应头
		filename := fmt.Sprintf("打卡照片_%s.zip", time.Now().Format("20060102150405"))
		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

		c.Data(200, "application/zip", buf.Bytes())
	}
}

// parseExportRequest 解析导出请求参数
func parseExportRequest(c *gin.Context) ExportRequest {
	var userID, appointmentID, projectID *uint

	if userStr := c.Query("user_id"); userStr != "" {
		if id, err := strconv.ParseUint(userStr, 10, 32); err == nil {
			uid := uint(id)
			userID = &uid
		}
	}

	if aptStr := c.Query("appointment_id"); aptStr != "" {
		if id, err := strconv.ParseUint(aptStr, 10, 32); err == nil {
			aid := uint(id)
			appointmentID = &aid
		}
	}

	if projStr := c.Query("project_id"); projStr != "" {
		if id, err := strconv.ParseUint(projStr, 10, 32); err == nil {
			pid := uint(id)
			projectID = &pid
		}
	}

	return ExportRequest{
		UserID:        userID,
		AppointmentID: appointmentID,
		ProjectID:     projectID,
		StartDate:     strings.TrimSpace(c.Query("start_date")),
		EndDate:       strings.TrimSpace(c.Query("end_date")),
		Status:        strings.TrimSpace(c.Query("status")),
	}
}
