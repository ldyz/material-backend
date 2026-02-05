package audit

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"gorm.io/gorm"
)

// Handler 操作日志处理器
type Handler struct {
	db *gorm.DB
}

// NewHandler 创建操作日志处理器
func NewHandler(database *gorm.DB) *Handler {
	return &Handler{db: database}
}

// GetLogs 获取操作日志列表
// GET /audit/operation-logs
func (h *Handler) GetLogs(c *gin.Context) {
	// 解析查询参数
	var filter OperationLogFilter
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid32 := uint(uid)
			filter.UserID = &uid32
		}
	}
	filter.Operation = c.Query("operation")
	filter.Module = c.Query("module")
	filter.ResourceType = c.Query("resource_type")
	filter.ResourceNo = c.Query("resource_no")
	filter.Status = c.Query("status")
	filter.Keyword = c.Query("keyword")

	// 日期范围
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}

	// 分页
	filter.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	filter.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 获取日志
	service := NewService(h.db)
	result, err := service.List(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取操作日志失败")
		return
	}

	fmt.Printf("[审计API] 返回日志数据: total=%d, page=%d, data_len=%d\n", result.Total, result.Page, len(result.Data))

	response.Success(c, gin.H{
		"data":         result.Data,
		"total":        result.Total,
		"page":         result.Page,
		"page_size":    result.PageSize,
		"total_pages":  result.TotalPages,
	})
}

// GetLogDetail 获取操作日志详情
// GET /audit/operation-logs/:id
func (h *Handler) GetLogDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的日志ID")
		return
	}

	service := NewService(h.db)
	log, err := service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusNotFound, "操作日志不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "获取操作日志详情失败")
		return
	}

	response.Success(c, log.ToDetailDTO())
}

// GetLogStatistics 获取操作统计信息
// GET /audit/operation-logs/statistics
func (h *Handler) GetLogStatistics(c *gin.Context) {
	days := 7 // 默认7天
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	service := NewService(h.db)
	stats, err := service.GetStatistics(days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// GetResourceLogs 获取指定资源的操作日志
// GET /audit/operation-logs/resource/:resource_type/:resource_id
func (h *Handler) GetResourceLogs(c *gin.Context) {
	resourceType := c.Param("resource_type")
	resourceIDStr := c.Param("resource_id")
	resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的资源ID")
		return
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	service := NewService(h.db)
	logs, err := service.GetByResource(resourceType, uint(resourceID), limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取资源日志失败")
		return
	}

	// 转换为DTO列表
	dtos := make([]map[string]any, len(logs))
	for i, log := range logs {
		dtos[i] = log.ToDTO()
	}

	response.Success(c, dtos)
}

// GetUserLogs 获取指定用户的操作日志
// GET /audit/operation-logs/user/:user_id
func (h *Handler) GetUserLogs(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	service := NewService(h.db)
	logs, err := service.GetByUser(uint(userID), limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户日志失败")
		return
	}

	// 转换为DTO列表
	dtos := make([]map[string]any, len(logs))
	for i, log := range logs {
		dtos[i] = log.ToDTO()
	}

	response.Success(c, dtos)
}

// ExportLogs 导出操作日志
// GET /audit/operation-logs/export
func (h *Handler) ExportLogs(c *gin.Context) {
	// 解析查询参数（与GetLogs相同）
	var filter OperationLogFilter
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid32 := uint(uid)
			filter.UserID = &uid32
		}
	}
	filter.Operation = c.Query("operation")
	filter.Module = c.Query("module")
	filter.ResourceType = c.Query("resource_type")
	filter.Status = c.Query("status")

	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}

	service := NewService(h.db)
	logs, err := service.ExportLogs(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "导出日志失败")
		return
	}

	// 转换为详细的DTO列表
	dtos := make([]map[string]any, len(logs))
	for i, log := range logs {
		dtos[i] = log.ToDetailDTO()
	}

	// 设置下载文件名
	filename := time.Now().Format("operation-logs-20060102-150405.json")

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.JSON(http.StatusOK, dtos)
}

// DeleteOldLogs 清理旧日志（需要管理员权限）
// DELETE /audit/operation-logs/cleanup
func (h *Handler) DeleteOldLogs(c *gin.Context) {
	// 检查是否是管理员
	permissions := c.GetString("permissions")
	if !auth.HasPermissionString(permissions, "admin") {
		response.Error(c, http.StatusForbidden, "需要管理员权限")
		return
	}

	days := 90 // 默认保留90天
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	service := NewService(h.db)
	count, err := service.DeleteOldLogs(days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "清理旧日志失败")
		return
	}

	response.Success(c, gin.H{
		"deleted_count": count,
		"message":      fmt.Sprintf("已清理 %d 天前的操作日志", days),
	})
}
