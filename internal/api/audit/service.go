package audit

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Service 操作日志服务
type Service struct {
	db *gorm.DB
}

// NewService 创建操作日志服务
func NewService(database *gorm.DB) *Service {
	return &Service{db: database}
}

// Log 记录操作日志
func (s *Service) Log(log *OperationLog) error {
	log.CreatedAt = time.Now()
	log.UpdatedAt = time.Now()

	if log.Status == "" {
		log.Status = LogStatusSuccess
	}

	return s.db.Create(log).Error
}

// LogWithChanges 记录带变更内容的操作日志
func (s *Service) LogWithChanges(log *OperationLog, changes LogChanges) error {
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		return fmt.Errorf("failed to marshal changes: %w", err)
	}
	log.Changes = changesJSON
	return s.Log(log)
}

// LogCreate 记录创建操作
func (s *Service) LogCreate(userID *uint, username, module, resourceType string, resourceID uint, resourceNo string, data interface{}) error {
	var paramsJSON json.RawMessage
	if data != nil {
		var err error
		paramsJSON, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	log := &OperationLog{
		UserID:        userID,
		Username:      username,
		Operation:     OpCreate,
		Module:        module,
		ResourceType:  resourceType,
		ResourceID:    &resourceID,
		ResourceNo:    resourceNo,
		RequestParams: paramsJSON,
		Status:        LogStatusSuccess,
	}
	return s.Log(log)
}

// LogUpdate 记录更新操作
func (s *Service) LogUpdate(userID *uint, username, module, resourceType string, resourceID uint, resourceNo string, before, after interface{}) error {
	changes := LogChanges{
		Before: before,
		After:  after,
	}
	return s.LogWithChanges(&OperationLog{
		UserID:       userID,
		Username:     username,
			Operation:    OpUpdate,
		Module:       module,
		ResourceType:  resourceType,
		ResourceID:   &resourceID,
		ResourceNo:   resourceNo,
		Status:       LogStatusSuccess,
	}, changes)
}

// LogDelete 记录删除操作
func (s *Service) LogDelete(userID *uint, username, module, resourceType string, resourceID uint, resourceNo string, deletedData interface{}) error {
	changes := LogChanges{
		Before: deletedData,
	}
	return s.LogWithChanges(&OperationLog{
		UserID:       userID,
		Username:     username,
		Operation:    OpDelete,
		Module:       module,
		ResourceType:  resourceType,
		ResourceID:   &resourceID,
		ResourceNo:   resourceNo,
		Status:       LogStatusSuccess,
	}, changes)
}

// LogApprove 记录审批操作
func (s *Service) LogApprove(userID *uint, username, module, resourceType string, resourceID uint, resourceNo string, comment string) error {
	params := map[string]any{
		"comment": comment,
	}
	paramsJSON, _ := json.Marshal(params)

	log := &OperationLog{
		UserID:        userID,
		Username:      username,
		Operation:     OpApprove,
		Module:        module,
		ResourceType:  resourceType,
		ResourceID:    &resourceID,
		ResourceNo:    resourceNo,
		RequestParams: paramsJSON,
		Status:        LogStatusSuccess,
	}
	return s.Log(log)
}

// LogReject 记录拒绝操作
func (s *Service) LogReject(userID *uint, username, module, resourceType string, resourceID uint, resourceNo string, reason string) error {
	params := map[string]any{
		"reason": reason,
	}
	paramsJSON, _ := json.Marshal(params)

	log := &OperationLog{
		UserID:        userID,
		Username:      username,
		Operation:     OpReject,
		Module:        module,
		ResourceType:  resourceType,
		ResourceID:    &resourceID,
		ResourceNo:    resourceNo,
		RequestParams: paramsJSON,
		Status:        LogStatusSuccess,
	}
	return s.Log(log)
}

// LogError 记录错误操作
func (s *Service) LogError(userID *uint, username, operation, module string, errorMessage string, requestPath string, params interface{}) error {
	var paramsJSON json.RawMessage
	if params != nil {
		var err error
		paramsJSON, err = json.Marshal(params)
		if err != nil {
			return err
		}
	}

	log := &OperationLog{
		UserID:        userID,
		Username:      username,
		Operation:     operation,
		Module:        module,
		RequestPath:   requestPath,
		RequestParams: paramsJSON,
		Status:        LogStatusError,
		ErrorMessage:  errorMessage,
	}
	return s.Log(log)
}

// GetByID 根据ID获取操作日志
func (s *Service) GetByID(id uint) (*OperationLog, error) {
	var log OperationLog
	err := s.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// List 查询操作日志列表
func (s *Service) List(filter OperationLogFilter) (*OperationLogListResponse, error) {
	var logs []OperationLog
	var total int64

	query := s.db.Model(&OperationLog{})

	// 应用过滤条件
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Operation != "" {
		query = query.Where("operation = ?", filter.Operation)
	}
	if filter.Module != "" {
		query = query.Where("module = ?", filter.Module)
	}
	if filter.ResourceType != "" {
		query = query.Where("resource_type = ?", filter.ResourceType)
	}
	if filter.ResourceNo != "" {
		query = query.Where("resource_no = ?", filter.ResourceNo)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	if filter.Keyword != "" {
		query = query.Where("resource_no LIKE ? OR username LIKE ?",
			"%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	offset := (filter.Page - 1) * filter.PageSize

	// 查询数据
	if err := query.
		Order("created_at DESC").
		Limit(filter.PageSize).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &OperationLogListResponse{
		Data:       logs,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByResource 获取指定资源的操作日志
func (s *Service) GetByResource(resourceType string, resourceID uint, limit int) ([]OperationLog, error) {
	var logs []OperationLog
	query := s.db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetByUser 获取指定用户的操作日志
func (s *Service) GetByUser(userID uint, limit int) ([]OperationLog, error) {
	var logs []OperationLog
	query := s.db.Where("user_id = ?", userID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetStatistics 获取操作统计信息
func (s *Service) GetStatistics(days int) (map[string]interface{}, error) {
	if days <= 0 {
		days = 7
	}

	since := time.Now().AddDate(0, 0, -days)

	type StatResult struct {
		Operation string `json:"operation"`
		Count     int64  `json:"count"`
	}

	var operationStats []StatResult
	err := s.db.Model(&OperationLog{}).
		Select("operation, count(*) as count").
		Where("created_at >= ?", since).
		Group("operation").
		Scan(&operationStats).Error

	if err != nil {
		return nil, err
	}

	// 总操作数
	var totalCount int64
	s.db.Model(&OperationLog{}).Where("created_at >= ?", since).Count(&totalCount)

	// 成功/失败统计
	var successCount, errorCount int64
	s.db.Model(&OperationLog{}).Where("created_at >= ? AND status = ?", since, LogStatusSuccess).Count(&successCount)
	s.db.Model(&OperationLog{}).Where("created_at >= ? AND status = ?", since, LogStatusError).Count(&errorCount)

	// 最活跃用户
	type UserStat struct {
		Username string `json:"username"`
		Count    int64  `json:"count"`
	}
	var topUsers []UserStat
	s.db.Model(&OperationLog{}).
		Select("username, count(*) as count").
		Where("created_at >= ?", since).
		Group("username").
		Order("count DESC").
		Limit(5).
		Scan(&topUsers)

	return map[string]interface{}{
		"total_operations": totalCount,
		"success_count":    successCount,
		"error_count":      errorCount,
		"by_operation":     operationStats,
		"top_users":        topUsers,
		"period_days":      days,
	}, nil
}

// DeleteOldLogs 删除旧日志（定期清理）
func (s *Service) DeleteOldLogs(days int) (int64, error) {
	if days <= 0 {
		days = 90 // 默认保留90天
	}

	cutoffDate := time.Now().AddDate(0, 0, -days)
	result := s.db.Where("created_at < ?", cutoffDate).Delete(&OperationLog{})

	return result.RowsAffected, result.Error
}

// ExportLogs 导出日志（可用于审计导出）
func (s *Service) ExportLogs(filter OperationLogFilter) ([]OperationLog, error) {
	var logs []OperationLog
	query := s.db.Model(&OperationLog{})

	// 应用相同的过滤条件
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Operation != "" {
		query = query.Where("operation = ?", filter.Operation)
	}
	if filter.Module != "" {
		query = query.Where("module = ?", filter.Module)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// 全局操作日志服务实例（用于便捷访问）
var auditService *Service

// InitAuditService 初始化审计服务
func InitAuditService(database *gorm.DB) {
	if auditService == nil {
		auditService = NewService(database)
	}
}

// 便捷函数（需要在应用启动时初始化）
func LogCreate(userID *uint, username, module, resourceType string, resourceID uint, resourceNo string, data interface{}) error {
	if auditService != nil {
		return auditService.LogCreate(userID, username, module, resourceType, resourceID, resourceNo, data)
	}
	return nil
}

func LogUpdate(userID *uint, username, module, resourceType string, resourceID uint, resourceNo string, before, after interface{}) error {
	if auditService != nil {
		return auditService.LogUpdate(userID, username, module, resourceType, resourceID, resourceNo, before, after)
	}
	return nil
}

func LogDelete(userID *uint, username, module, resourceType string, resourceID uint, resourceNo string, deletedData interface{}) error {
	if auditService != nil {
		return auditService.LogDelete(userID, username, module, resourceType, resourceID, resourceNo, deletedData)
	}
	return nil
}

func LogApprove(userID *uint, username, module, resourceType string, resourceID uint, resourceNo string, comment string) error {
	if auditService != nil {
		return auditService.LogApprove(userID, username, module, resourceType, resourceID, resourceNo, comment)
	}
	return nil
}

func LogReject(userID *uint, username, module, resourceType string, resourceID uint, resourceNo string, reason string) error {
	if auditService != nil {
		return auditService.LogReject(userID, username, module, resourceType, resourceID, resourceNo, reason)
	}
	return nil
}
