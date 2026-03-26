package appointment

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/yourorg/material-backend/backend/internal/api/notification"
	"gorm.io/gorm"
)

// AppointmentService 预约服务
type AppointmentService struct {
	db              *gorm.DB
	calendarService *CalendarService
}

// NewAppointmentService 创建预约服务
func NewAppointmentService(db *gorm.DB) *AppointmentService {
	calendarService := NewCalendarService(db)
	return &AppointmentService{
		db:              db,
		calendarService: calendarService,
	}
}

// Create 创建预约单
func (s *AppointmentService) Create(req CreateAppointmentRequest, applicantID uint, applicantName string) (*ConstructionAppointment, error) {
	// 解析日期
	workDate, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		return nil, fmt.Errorf("invalid work_date: %w", err)
	}

	// 验证日期不能是过去的日期
	if workDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("作业日期不能早于今天")
	}

	// 验证时间段
	if err := s.calendarService.ValidateTimeSlot(req.TimeSlot); err != nil {
		return nil, err
	}

	// 验证加急原因
	if req.IsUrgent && req.Priority >= 7 && req.UrgentReason == "" {
		return nil, errors.New("高优先级加急预约必须提供加急原因")
	}

	// 创建预约单
	appointment := &ConstructionAppointment{
		ProjectID:           req.ProjectID, // 直接使用指针，可能为 nil
		ApplicantID:         applicantID,
		ApplicantName:       applicantName,
		ContactPhone:        req.ContactPhone,
		ContactPerson:       req.ContactPerson,
		WorkDate:            workDate,
		TimeSlot:            req.TimeSlot,
		WorkLocation:        req.WorkLocation,
		WorkContent:         req.WorkContent,
		WorkType:            req.WorkType,
		IsUrgent:            req.IsUrgent,
		Priority:            req.Priority,
		UrgentReason:        req.UrgentReason,
		AssignedWorkerID:    req.AssignedWorkerID,
		AssignedWorkerIDs:   req.AssignedWorkerIDs,
		AssignedWorkerNames: req.AssignedWorkerNames,
		Status:              StatusDraft,
	}

	// 设置主作业人员姓名（兼容性）
	if appointment.AssignedWorkerID != nil && req.AssignedWorkerNames != "" {
		names := parseCommaSeparatedNames(req.AssignedWorkerNames)
		if len(names) > 0 {
			appointment.AssignedWorkerName = names[0]
		}
	}

	// 如果指定了作业人员，检查可用性
	if req.AssignedWorkerID != nil {
		available, reason, err := s.calendarService.CheckAvailability(*req.AssignedWorkerID, workDate, req.TimeSlot)
		if err != nil {
			return nil, fmt.Errorf("failed to check availability: %w", err)
		}
		if !available {
			return nil, fmt.Errorf("作业人员在指定时间段不可用: %s", reason)
		}
	}

	// 保存到数据库
	if err := s.db.Create(appointment).Error; err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return appointment, nil
}

// Update 更新预约单
func (s *AppointmentService) Update(id uint, req UpdateAppointmentRequest, currentUserID uint) (*ConstructionAppointment, error) {
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("预约单不存在")
		}
		return nil, err
	}

	// 检查是否可编辑（包含权限验证）
	if !appointment.IsEditableBy(currentUserID) {
		// 根据具体情况返回不同的错误信息
		if appointment.Status != StatusDraft && appointment.Status != StatusPending {
			return nil, errors.New("只有草稿或待审批状态的预约单可以编辑")
		}
		return nil, errors.New("只有申请人可以编辑预约单")
	}

	// 解析日期
	if req.WorkDate != "" {
		workDate, err := time.Parse("2006-01-02", req.WorkDate)
		if err != nil {
			return nil, fmt.Errorf("invalid work_date: %w", err)
		}
		appointment.WorkDate = workDate
	}

	if req.TimeSlot != "" {
		if err := s.calendarService.ValidateTimeSlot(req.TimeSlot); err != nil {
			return nil, err
		}
		appointment.TimeSlot = req.TimeSlot
	}

	// 验证加急原因
	if req.IsUrgent && req.Priority >= 7 && req.UrgentReason == "" {
		return nil, errors.New("高优先级加急预约必须提供加急原因")
	}

	// 如果指定了作业人员，检查可用性
	if req.AssignedWorkerID != nil {
		available, reason, err := s.calendarService.CheckAvailability(*req.AssignedWorkerID, appointment.WorkDate, appointment.TimeSlot)
		if err != nil {
			return nil, fmt.Errorf("failed to check availability: %w", err)
		}
		if !available {
			return nil, fmt.Errorf("作业人员在指定时间段不可用: %s", reason)
		}
		appointment.AssignedWorkerID = req.AssignedWorkerID
	}

	// 更新字段
	if req.ProjectID != nil {
		appointment.ProjectID = req.ProjectID
	}
	if req.ContactPhone != "" {
		appointment.ContactPhone = req.ContactPhone
	}
	if req.ContactPerson != "" {
		appointment.ContactPerson = req.ContactPerson
	}
	if req.WorkLocation != "" {
		appointment.WorkLocation = req.WorkLocation
	}
	if req.WorkContent != "" {
		appointment.WorkContent = req.WorkContent
	}
	if req.WorkType != "" {
		appointment.WorkType = req.WorkType
	}
	appointment.IsUrgent = req.IsUrgent
	appointment.Priority = req.Priority
	appointment.UrgentReason = req.UrgentReason

	// 保存
	if err := s.db.Save(&appointment).Error; err != nil {
		return nil, fmt.Errorf("failed to update appointment: %w", err)
	}

	return &appointment, nil
}

// GetByID 根据ID获取预约单
func (s *AppointmentService) GetByID(id uint) (*ConstructionAppointment, error) {
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("预约单不存在")
		}
		return nil, err
	}
	return &appointment, nil
}

// GetByNo 根据预约单号获取预约单
func (s *AppointmentService) GetByNo(appointmentNo string) (*ConstructionAppointment, error) {
	var appointment ConstructionAppointment
	if err := s.db.Where("appointment_no = ?", appointmentNo).First(&appointment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("预约单不存在")
		}
		return nil, err
	}
	return &appointment, nil
}

// List 查询预约单列表
func (s *AppointmentService) List(req AppointmentListRequest) ([]ConstructionAppointment, int64, error) {
	var appointments []ConstructionAppointment
	var total int64

	query := s.db.Model(&ConstructionAppointment{})

	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 加急过滤
	if req.IsUrgent != nil {
		query = query.Where("is_urgent = ?", *req.IsUrgent)
	}

	// 日期范围过滤
	if req.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			query = query.Where("work_date >= ?", startDate.Format("2006-01-02"))
		}
	}
	if req.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			query = query.Where("work_date <= ?", endDate.Format("2006-01-02"))
		}
	}

	// 申请人过滤
	if req.ApplicantID != nil {
		query = query.Where("applicant_id = ?", *req.ApplicantID)
	}

	// 作业人员过滤
	if req.WorkerID != nil {
		query = query.Where("assigned_worker_id = ?", *req.WorkerID)
	}

	// 作业类型过滤
	if req.WorkType != "" {
		query = query.Where("work_type = ?", req.WorkType)
	}

	// 权限过滤：草稿状态的预约单只对创建者可见
	// 如果当前用户ID存在，且没有指定状态过滤，则排除其他人的草稿
	if req.CurrentUserID != 0 && req.Status == "" {
		// 查看所有人的非草稿预约单，或者自己创建的预约单（包括草稿）
		query = query.Where("(status != ? OR applicant_id = ?)", StatusDraft, req.CurrentUserID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&appointments).Error; err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

// Delete 删除预约单（仅草稿状态）
func (s *AppointmentService) Delete(id uint) error {
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("预约单不存在")
		}
		return err
	}

	if !appointment.IsEditable() {
		return errors.New("只有草稿状态的预约单可以删除")
	}

	return s.db.Delete(&appointment).Error
}

// Submit 提交预约单审批
func (s *AppointmentService) Submit(id uint) (*ConstructionAppointment, error) {
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, id).Error; err != nil {
		return nil, err
	}

	if appointment.Status != StatusDraft {
		return nil, errors.New("只有草稿状态的预约单可以提交")
	}

	now := time.Now()
	appointment.Status = StatusPending
	appointment.SubmittedAt = &now

	if err := s.db.Save(&appointment).Error; err != nil {
		return nil, err
	}

	return &appointment, nil
}

// AssignWorker 分配作业人员
func (s *AppointmentService) AssignWorker(id uint, workerID uint, workerName string) (*ConstructionAppointment, error) {
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, id).Error; err != nil {
		return nil, err
	}

	// 检查作业人员可用性
	available, reason, err := s.calendarService.CheckAvailability(workerID, appointment.WorkDate, appointment.TimeSlot)
	if err != nil {
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}
	if !available {
		return nil, fmt.Errorf("作业人员在指定时间段不可用: %s", reason)
	}

	// 如果已有作业人员，先释放其日历
	if appointment.AssignedWorkerID != nil {
		if err := s.calendarService.CancelAppointmentForWorker(&appointment); err != nil {
			return nil, fmt.Errorf("failed to release previous worker: %w", err)
		}
	}

	// 分配新作业人员
	appointment.AssignedWorkerID = &workerID
	appointment.AssignedWorkerName = workerName

	if err := s.db.Save(&appointment).Error; err != nil {
		return nil, err
	}

	// 预约日历
	if err := s.calendarService.BookAppointmentForWorker(&appointment); err != nil {
		return nil, fmt.Errorf("failed to book calendar: %w", err)
	}

	// 通知作业人员
	s.notifyWorkerAssigned(&appointment)

	return &appointment, nil
}

// AssignWorkers 分配多个作业人员
func (s *AppointmentService) AssignWorkers(id uint, workerIDs []uint, workerNames []string, supervisorID *uint, supervisorName string) (*ConstructionAppointment, error) {
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, id).Error; err != nil {
		return nil, err
	}

	// 检查所有作业人员可用性
	for _, workerID := range workerIDs {
		available, reason, err := s.calendarService.CheckAvailability(workerID, appointment.WorkDate, appointment.TimeSlot)
		if err != nil {
			return nil, fmt.Errorf("failed to check availability: %w", err)
		}
		if !available {
			return nil, fmt.Errorf("作业人员在指定时间段不可用: %s", reason)
		}
	}

	// 如果已有作业人员，先释放其日历
	if appointment.AssignedWorkerID != nil {
		if err := s.calendarService.CancelAppointmentForWorker(&appointment); err != nil {
			return nil, fmt.Errorf("failed to release previous worker: %w", err)
		}
	}

	// 将workerID列表序列化为JSON
	workerIDsJSON, err := json.Marshal(workerIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal worker IDs: %w", err)
	}

	// 将workerNames用逗号连接
	workerNamesStr := strings.Join(workerNames, ",")

	// 分配新作业人员
	appointment.AssignedWorkerIDs = string(workerIDsJSON)
	appointment.AssignedWorkerNames = workerNamesStr

	// 为了兼容性，如果有作业人员，设置第一个为默认值
	if len(workerIDs) > 0 {
		appointment.AssignedWorkerID = &workerIDs[0]
		if len(workerNames) > 0 {
			appointment.AssignedWorkerName = workerNames[0]
		}
	}

	// 保存监护人信息
	if supervisorID != nil && *supervisorID > 0 {
		appointment.SupervisorID = supervisorID
		appointment.SupervisorName = supervisorName
	} else {
		appointment.SupervisorID = nil
		appointment.SupervisorName = ""
	}

	if err := s.db.Save(&appointment).Error; err != nil {
		return nil, err
	}

	// 为每个作业人员预约日历
	for _, workerID := range workerIDs {
		tempAppointment := appointment
		tempAppointment.AssignedWorkerID = &workerID
		if err := s.calendarService.BookAppointmentForWorker(&tempAppointment); err != nil {
			// 记录错误但继续
			log.Printf("Warning: failed to book calendar for worker %d: %v", workerID, err)
		}
	}

	// 通知所有作业人员
	s.notifyWorkersAssigned(&appointment, workerIDs)

	return &appointment, nil
}

// StartWork 开始作业
func (s *AppointmentService) StartWork(id uint) (*ConstructionAppointment, error) {
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, id).Error; err != nil {
		return nil, err
	}

	if appointment.Status != StatusScheduled {
		return nil, errors.New("只有已排期状态的预约单可以开始作业")
	}

	appointment.Status = StatusInProgress
	if err := s.db.Save(&appointment).Error; err != nil {
		return nil, err
	}

	return &appointment, nil
}

// Complete 完成预约单
func (s *AppointmentService) Complete(id uint, req CompleteAppointmentRequest) (*ConstructionAppointment, error) {
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, id).Error; err != nil {
		return nil, err
	}

	if !appointment.CanComplete() {
		return nil, errors.New("当前状态不允许完成")
	}

	now := time.Now()
	appointment.Status = StatusCompleted
	appointment.CompletedAt = &now

	if err := s.db.Save(&appointment).Error; err != nil {
		return nil, err
	}

	// 释放日历
	if err := s.calendarService.CancelAppointmentForWorker(&appointment); err != nil {
		// 日历释放失败不影响完成操作
		fmt.Printf("Warning: failed to release calendar: %v\n", err)
	}

	return &appointment, nil
}

// Cancel 取消预约单
func (s *AppointmentService) Cancel(id uint, reason string) (*ConstructionAppointment, error) {
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, id).Error; err != nil {
		return nil, err
	}

	if !appointment.IsCancellable() && !appointment.IsEditable() {
		return nil, errors.New("当前状态不允许取消")
	}

	appointment.Status = StatusCancelled

	if err := s.db.Save(&appointment).Error; err != nil {
		return nil, err
	}

	// 释放日历
	if err := s.calendarService.CancelAppointmentForWorker(&appointment); err != nil {
		fmt.Printf("Warning: failed to release calendar: %v\n", err)
	}

	return &appointment, nil
}

// GetStats 获取统计数据
func (s *AppointmentService) GetStats(filterDate *time.Time, applicantID *uint) (*StatsResponse, error) {
	var stats StatsResponse

	query := s.db.Model(&ConstructionAppointment{})

	// 日期过滤
	if filterDate != nil {
		startOfMonth := time.Date(filterDate.Year(), filterDate.Month(), 1, 0, 0, 0, 0, filterDate.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, -1)
		query = query.Where("work_date >= ? AND work_date <= ?", startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02"))
	}

	// 申请人过滤
	if applicantID != nil {
		query = query.Where("applicant_id = ?", *applicantID)
	}

	// 总数
	query.Count(&stats.Total)

	// 各状态计数
	s.db.Model(&ConstructionAppointment{}).Where("status = ?", StatusDraft).Count(&stats.Draft)
	s.db.Model(&ConstructionAppointment{}).Where("status = ?", StatusPending).Count(&stats.Pending)
	s.db.Model(&ConstructionAppointment{}).Where("status = ?", StatusScheduled).Count(&stats.Scheduled)
	s.db.Model(&ConstructionAppointment{}).Where("status = ?", StatusInProgress).Count(&stats.InProgress)
	s.db.Model(&ConstructionAppointment{}).Where("status = ?", StatusCompleted).Count(&stats.Completed)
	s.db.Model(&ConstructionAppointment{}).Where("status = ?", StatusCancelled).Count(&stats.Cancelled)
	s.db.Model(&ConstructionAppointment{}).Where("status = ?", StatusRejected).Count(&stats.Rejected)

	// 加急计数
	s.db.Model(&ConstructionAppointment{}).Where("is_urgent = ?", true).Count(&stats.Urgent)

	// 今日计数
	today := time.Now().Format("2006-01-02")
	s.db.Model(&ConstructionAppointment{}).Where("work_date = ?", today).Count(&stats.TodayCount)

	// 本周计数
	startOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 6)
	s.db.Model(&ConstructionAppointment{}).
		Where("work_date >= ? AND work_date <= ?", startOfWeek.Format("2006-01-02"), endOfWeek.Format("2006-01-02")).
		Count(&stats.WeekCount)

	// 本月计数
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	s.db.Model(&ConstructionAppointment{}).
		Where("work_date >= ? AND work_date <= ?", startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02")).
		Count(&stats.MonthCount)

	return &stats, nil
}

// BatchCreate 批量创建预约单
func (s *AppointmentService) BatchCreate(req BatchCreateAppointmentRequest, applicantID uint, applicantName string) ([]ConstructionAppointment, []error) {
	appointments := make([]ConstructionAppointment, 0, len(req.Appointments))
	errs := make([]error, len(req.Appointments))

	for i, apptReq := range req.Appointments {
		appt, err := s.Create(apptReq, applicantID, applicantName)
		if err != nil {
			errs[i] = err
			continue
		}
		appointments = append(appointments, *appt)
	}

	// 过滤掉 nil 错误
	cleanErrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			cleanErrs = append(cleanErrs, err)
		}
	}

	return appointments, cleanErrs
}

// GetPendingApprovals 获取待审批的预约单
func (s *AppointmentService) GetPendingApprovals(page, pageSize int) ([]ConstructionAppointment, int64, error) {
	var appointments []ConstructionAppointment
	var total int64

	query := s.db.Model(&ConstructionAppointment{}).Where("status = ?", StatusPending)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("is_urgent DESC, priority DESC, created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&appointments).Error; err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

// GetWorkerAppointments 获取作业人员的预约列表
func (s *AppointmentService) GetWorkerAppointments(workerID uint, startDate, endDate time.Time) ([]ConstructionAppointment, error) {
	var appointments []ConstructionAppointment

	err := s.db.Where("assigned_worker_id = ? AND work_date >= ? AND work_date <= ?",
		workerID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Where("status IN ?", []string{StatusScheduled, StatusInProgress}).
		Order("work_date ASC, time_slot ASC").
		Find(&appointments).Error

	return appointments, err
}

// SearchByKeyword 根据关键词搜索预约单
func (s *AppointmentService) SearchByKeyword(keyword string, page, pageSize int) ([]ConstructionAppointment, int64, error) {
	var appointments []ConstructionAppointment
	var total int64

	query := s.db.Model(&ConstructionAppointment{}).
		Where("appointment_no LIKE ? OR work_location LIKE ? OR work_content LIKE ? OR applicant_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&appointments).Error; err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

// UpdateStatus 更新状态
func (s *AppointmentService) UpdateStatus(id uint, status string) (*ConstructionAppointment, error) {
	var appointment ConstructionAppointment
	if err := s.db.First(&appointment, id).Error; err != nil {
		return nil, err
	}

	// 状态流转验证
	validTransitions := map[string][]string{
		StatusDraft:     {StatusPending, StatusCancelled},
		StatusPending:   {StatusScheduled, StatusRejected, StatusCancelled},
		StatusScheduled: {StatusInProgress, StatusCancelled},
		StatusInProgress: {StatusCompleted},
	}

	if allowedStates, ok := validTransitions[appointment.Status]; ok {
		allowed := false
		for _, s := range allowedStates {
			if s == status {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("不能从状态 %s 转换到 %s", appointment.Status, status)
		}
	} else {
		return nil, fmt.Errorf("未知的状态: %s", appointment.Status)
	}

	appointment.Status = status

	// 更新时间戳
	now := time.Now()
	if status == StatusScheduled {
		appointment.ApprovedAt = &now
	}

	if err := s.db.Save(&appointment).Error; err != nil {
		return nil, err
	}

	return &appointment, nil
}

// ExportToJSON 导出为JSON
func (s *AppointmentService) ExportToJSON(ids []uint) ([]byte, error) {
	var appointments []ConstructionAppointment
	if err := s.db.Where("id IN ?", ids).Find(&appointments).Error; err != nil {
		return nil, err
	}

	data := make([]map[string]interface{}, len(appointments))
	for i, a := range appointments {
		data[i] = a.ToDTO()
	}

	return json.MarshalIndent(data, "", "  ")
}

// notifyWorkerAssigned 通知单个作业人员被分配任务
func (s *AppointmentService) notifyWorkerAssigned(appointment *ConstructionAppointment) {
	if appointment.AssignedWorkerID == nil || *appointment.AssignedWorkerID == 0 {
		return
	}

	notificationData := map[string]interface{}{
		"appointment_id": appointment.ID,
		"appointment_no": appointment.AppointmentNo,
		"work_date":      appointment.WorkDate.Format("2006-01-02"),
		"time_slot":      appointment.TimeSlot,
		"work_location":  appointment.WorkLocation,
		"work_content":   appointment.WorkContent,
		"assigned_at":    time.Now(),
	}

	title := "新的作业任务分配"
	content := fmt.Sprintf("您被分配了一项施工任务，单号：%s，作业时间：%s %s，地点：%s，内容：%s",
		appointment.AppointmentNo,
		appointment.WorkDate.Format("2006-01-02"),
		appointment.TimeSlot,
		appointment.WorkLocation,
		appointment.WorkContent)

	if err := notification.CreateNotification(s.db, *appointment.AssignedWorkerID, notification.TypeAppointmentApprove, title, content, notificationData); err != nil {
		log.Printf("通知作业人员失败: %v", err)
	}
}

// notifyWorkersAssigned 通知多个作业人员被分配任务
func (s *AppointmentService) notifyWorkersAssigned(appointment *ConstructionAppointment, workerIDs []uint) {
	now := time.Now()

	for _, workerID := range workerIDs {
		notificationData := map[string]interface{}{
			"appointment_id": appointment.ID,
			"appointment_no": appointment.AppointmentNo,
			"work_date":      appointment.WorkDate.Format("2006-01-02"),
			"time_slot":      appointment.TimeSlot,
			"work_location":  appointment.WorkLocation,
			"work_content":   appointment.WorkContent,
			"assigned_at":    now,
		}

		title := "新的作业任务分配"
		content := fmt.Sprintf("您被分配了一项施工任务，单号：%s，作业时间：%s %s，地点：%s",
			appointment.AppointmentNo,
			appointment.WorkDate.Format("2006-01-02"),
			appointment.TimeSlot,
			appointment.WorkLocation)

		if err := notification.CreateNotification(s.db, workerID, notification.TypeAppointmentApprove, title, content, notificationData); err != nil {
			log.Printf("通知作业人员 %d 失败: %v", workerID, err)
		}
	}
}

// WorkerInfo 作业人员信息
type WorkerInfo struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// GetWorkersList 获取作业人员列表
func (s *AppointmentService) GetWorkersList() ([]WorkerInfo, error) {
	var workers []WorkerInfo

	// 查询作业人员角色的用户（支持中英文角色名）
	err := s.db.Table("users").
		Select("id, full_name, username, email, avatar").
		Where("is_active = ? AND (role = ? OR role = ? OR role = ? OR role = ?)",
			true,
			"worker",           // 英文作业人员
			"作业人员",          // 中文作业人员
			"team_leader",      // 英文班组长
			"班组长",            // 中文班组长
		).
		Find(&workers).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get workers list: %w", err)
	}

	return workers, nil
}

// GetWorkerByID 根据ID获取作业人员信息
func (s *AppointmentService) GetWorkerByID(workerID uint) (*WorkerInfo, error) {
	var worker WorkerInfo

	err := s.db.Table("users").
		Select("id, full_name, username, email, avatar").
		Where("id = ?", workerID).
		Where("is_active = ? AND (role = ? OR role = ? OR role = ? OR role = ?)",
			true,
			"worker",           // 英文作业人员
			"作业人员",          // 中文作业人员
			"team_leader",      // 英文班组长
			"班组长",            // 中文班组长
		).
		First(&worker).Error

	if err != nil {
		return nil, fmt.Errorf("worker not found: %w", err)
	}

	return &worker, nil
}

// GetDailyStatistics 获取每日预约统计数据
func (s *AppointmentService) GetDailyStatistics(startDate, endDate string) (*DailyStatisticsResponse, error) {
	// 获取总作业人员数
	var totalWorkers int64
	err := s.db.Table("users").
		Where("is_active = ? AND (role = ? OR role = ? OR role = ? OR role = ?)",
			true,
			"worker",
			"作业人员",
			"team_leader",
			"班组长",
		).
		Count(&totalWorkers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get total workers count: %w", err)
	}

	// 获取每日预约统计
	type DailyCount struct {
		Date        string
		TotalCount  int64
		UrgentCount int64
	}

	var dailyCounts []DailyCount
	err = s.db.Model(&ConstructionAppointment{}).
		Select("TO_CHAR(work_date, 'YYYY-MM-DD') as date, COUNT(*) as total_count, SUM(CASE WHEN is_urgent = true THEN 1 ELSE 0 END) as urgent_count").
		Where("work_date >= ? AND work_date <= ?", startDate, endDate).
		Where("status IN ?", []string{StatusPending, StatusScheduled, StatusInProgress}).
		Group("TO_CHAR(work_date, 'YYYY-MM-DD')").
		Order("date ASC").
		Scan(&dailyCounts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get daily statistics: %w", err)
	}

	// 转换为响应格式
	statistics := make([]DailyStatistics, len(dailyCounts))
	for i, dc := range dailyCounts {
		statistics[i] = DailyStatistics{
			Date:         dc.Date,
			TotalCount:   dc.TotalCount,
			UrgentCount:  dc.UrgentCount,
			TotalWorkers: totalWorkers,
		}
	}

	return &DailyStatisticsResponse{
		Statistics:   statistics,
		TotalWorkers: totalWorkers,
	}, nil
}

// GetTimeSlotStatistics 获取指定日期的时间段统计数据
func (s *AppointmentService) GetTimeSlotStatistics(date string) (*TimeSlotStatisticsResponse, error) {
	// 获取总作业人员数
	var totalWorkers int64
	err := s.db.Table("users").
		Where("is_active = ? AND (role = ? OR role = ? OR role = ? OR role = ?)",
			true,
			"worker",
			"作业人员",
			"team_leader",
			"班组长",
		).
		Count(&totalWorkers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get total workers count: %w", err)
	}

	// 定义所有时间段
	timeSlots := []string{TimeSlotMorning, TimeSlotNoon, TimeSlotAfternoon, TimeSlotFullDay}

	// 获取每个时间段的统计
	type TimeSlotCount struct {
		TimeSlot   string
		TotalCount int64
	}

	var timeSlotCounts []TimeSlotCount
	err = s.db.Model(&ConstructionAppointment{}).
		Select("time_slot as time_slot, COUNT(*) as total_count").
		Where("work_date = ?", date).
		Where("time_slot IN ?", timeSlots).
		Where("status IN ?", []string{StatusPending, StatusScheduled, StatusInProgress}).
		Group("time_slot").
		Scan(&timeSlotCounts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get time slot statistics: %w", err)
	}

	// 构建结果，确保所有时间段都有数据
	countMap := make(map[string]int64)
	for _, tsc := range timeSlotCounts {
		countMap[tsc.TimeSlot] = tsc.TotalCount
	}

	statistics := make([]TimeSlotStatistics, len(timeSlots))
	for i, slot := range timeSlots {
		statistics[i] = TimeSlotStatistics{
			Date:         date,
			TimeSlot:     slot,
			TotalCount:   countMap[slot],
			TotalWorkers: totalWorkers,
		}
	}

	return &TimeSlotStatisticsResponse{
		Statistics:   statistics,
		TotalWorkers: totalWorkers,
	}, nil
}

// parseCommaSeparatedNames 解析逗号分隔的姓名列表
func parseCommaSeparatedNames(namesStr string) []string {
	if namesStr == "" {
		return []string{}
	}
	names := strings.Split(namesStr, ",")
	result := make([]string, 0, len(names))
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// GetPendingApprovalCount 获取待审批数量
func (s *AppointmentService) GetPendingApprovalCount() (int, error) {
	var count int64
	s.db.Model(&ConstructionAppointment{}).Where("status = ?", StatusPending).Count(&count)
	return int(count), nil
}
