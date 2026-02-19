package appointment

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CalendarService 日历服务
type CalendarService struct {
	db *gorm.DB
}

// NewCalendarService 创建日历服务
func NewCalendarService(db *gorm.DB) *CalendarService {
	return &CalendarService{db: db}
}

// CheckAvailability 检查作业人员在指定时间段的可用性
func (s *CalendarService) CheckAvailability(workerID uint, workDate time.Time, timeSlot string) (bool, string, error) {
	var calendar WorkerCalendar
	err := s.db.Where("worker_id = ? AND calendar_date = ? AND time_slot = ?",
		workerID, workDate.Format("2006-01-02"), timeSlot).
		First(&calendar).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有记录，表示可用
		return true, "", nil
	}

	if err != nil {
		return false, "", err
	}

	if !calendar.IsAvailable || calendar.Status != CalendarStatusAvailable {
		reason := calendar.BlockedReason
		if reason == "" {
			reason = "该时间段已被占用"
		}
		return false, reason, nil
	}

	return true, "", nil
}

// CheckMultipleAvailability 批量检查多个作业人员的可用性
func (s *CalendarService) CheckMultipleAvailability(workerIDs []uint, workDate time.Time, timeSlot string) (map[uint]bool, error) {
	result := make(map[uint]bool)

	// 初始化所有人都为可用
	for _, id := range workerIDs {
		result[id] = true
	}

	if len(workerIDs) == 0 {
		return result, nil
	}

	var calendars []WorkerCalendar
	err := s.db.Where("worker_id IN ? AND calendar_date = ? AND time_slot = ?",
		workerIDs, workDate.Format("2006-01-02"), timeSlot).
		Find(&calendars).Error

	if err != nil {
		return nil, err
	}

	for _, cal := range calendars {
		if !cal.IsAvailable || cal.Status != CalendarStatusAvailable {
			result[cal.WorkerID] = false
		}
	}

	return result, nil
}

// BlockCalendar 锁定日历时间段
func (s *CalendarService) BlockCalendar(workerID uint, workDate time.Time, timeSlot, reason string, appointmentID *uint) error {
	// 检查是否已存在记录
	var calendar WorkerCalendar
	err := s.db.Where("worker_id = ? AND calendar_date = ? AND time_slot = ?",
		workerID, workDate.Format("2006-01-02"), timeSlot).
		First(&calendar).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新记录
		calendar = WorkerCalendar{
			WorkerID:      workerID,
			CalendarDate:  workDate,
			TimeSlot:      timeSlot,
			IsAvailable:   false,
			Status:        CalendarStatusBlocked,
			AppointmentID: appointmentID,
			BlockedReason: reason,
		}
		return s.db.Create(&calendar).Error
	}

	if err != nil {
		return err
	}

	// 更新现有记录
	calendar.IsAvailable = false
	calendar.Status = CalendarStatusBlocked
	calendar.AppointmentID = appointmentID
	calendar.BlockedReason = reason
	return s.db.Save(&calendar).Error
}

// ReleaseCalendar 释放日历时间段
func (s *CalendarService) ReleaseCalendar(workerID uint, workDate time.Time, timeSlot string) error {
	return s.db.Model(&WorkerCalendar{}).
		Where("worker_id = ? AND calendar_date = ? AND time_slot = ?",
			workerID, workDate.Format("2006-01-02"), timeSlot).
		Updates(map[string]interface{}{
			"is_available":   true,
			"status":         CalendarStatusAvailable,
			"appointment_id": nil,
			"blocked_reason": "",
		}).Error
}

// BookAppointmentForWorker 为作业人员预约时间段
func (s *CalendarService) BookAppointmentForWorker(appointment *ConstructionAppointment) error {
	if appointment.AssignedWorkerID == nil {
		return nil // 没有分配作业人员，不需要预约
	}

	return s.BlockCalendar(
		*appointment.AssignedWorkerID,
		appointment.WorkDate,
		appointment.TimeSlot,
		fmt.Sprintf("预约单: %s - %s", appointment.AppointmentNo, appointment.WorkContent),
		&appointment.ID,
	)
}

// CancelAppointmentForWorker 取消作业人员的预约时间段
func (s *CalendarService) CancelAppointmentForWorker(appointment *ConstructionAppointment) error {
	if appointment.AssignedWorkerID == nil {
		return nil
	}

	return s.ReleaseCalendar(
		*appointment.AssignedWorkerID,
		appointment.WorkDate,
		appointment.TimeSlot,
	)
}

// GetWorkerCalendar 获取作业人员的日历
func (s *CalendarService) GetWorkerCalendar(workerID uint, startDate, endDate time.Time) ([]WorkerCalendar, error) {
	var calendars []WorkerCalendar
	err := s.db.Where("worker_id = ? AND calendar_date >= ? AND calendar_date <= ?",
		workerID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("calendar_date ASC, time_slot ASC").
		Find(&calendars).Error

	if err != nil {
		return nil, err
	}

	// 如果某些日期没有记录，返回完整的日历（包含可用的时间段）
	return s.fillMissingDates(calendars, workerID, startDate, endDate), nil
}

// fillMissingDates 填充缺失的日期
func (s *CalendarService) fillMissingDates(existing []WorkerCalendar, workerID uint, startDate, endDate time.Time) []WorkerCalendar {
	// 创建已存在日历的映射
	calendarMap := make(map[string]*WorkerCalendar)
	for i := range existing {
		key := fmt.Sprintf("%s_%s", existing[i].CalendarDate.Format("2006-01-02"), existing[i].TimeSlot)
		calendarMap[key] = &existing[i]
	}

	// 生成所有日期
	var result []WorkerCalendar
	timeSlots := []string{TimeSlotMorning, TimeSlotNoon, TimeSlotAfternoon}

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		for _, slot := range timeSlots {
			key := fmt.Sprintf("%s_%s", dateStr, slot)
			if cal, ok := calendarMap[key]; ok {
				result = append(result, *cal)
			} else {
				// 创建默认的可用日历项
				result = append(result, WorkerCalendar{
					WorkerID:     workerID,
					CalendarDate: d,
					TimeSlot:     slot,
					IsAvailable:  true,
					Status:       CalendarStatusAvailable,
				})
			}
		}
	}

	return result
}

// BatchBlockCalendar 批量锁定日历
func (s *CalendarService) BatchBlockCalendar(req BatchBlockCalendarRequest) error {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end_date: %w", err)
	}

	// 验证时间段
	validSlots := map[string]bool{
		TimeSlotMorning: true, TimeSlotNoon: true, TimeSlotAfternoon: true, TimeSlotFullDay: true,
	}
	for _, slot := range req.TimeSlots {
		if !validSlots[slot] {
			return fmt.Errorf("invalid time_slot: %s", slot)
		}
	}

	// 遍历日期和时间段
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		for _, slot := range req.TimeSlots {
			err := s.BlockCalendar(req.WorkerID, d, slot, req.BlockedReason, nil)
			if err != nil {
				return fmt.Errorf("failed to block %s %s: %w", d.Format("2006-01-02"), slot, err)
			}
		}
	}

	return nil
}

// GetAvailableWorkers 获取指定时间段可用的作业人员列表
func (s *CalendarService) GetAvailableWorkers(workDate time.Time, timeSlot string) ([]uint, error) {
	// 获取所有作业人员角色的用户
	type Worker struct {
		ID uint
	}
	var workers []Worker
	err := s.db.Table("users").
		Select("id").
		Where("role = ?", "worker"). // 假设作业人员角色为 "worker"
		Find(&workers).Error

	if err != nil {
		return nil, err
	}

	workerIDs := make([]uint, len(workers))
	for i, w := range workers {
		workerIDs[i] = w.ID
	}

	// 检查可用性
	availability, err := s.CheckMultipleAvailability(workerIDs, workDate, timeSlot)
	if err != nil {
		return nil, err
	}

	// 返回可用的作业人员
	var available []uint
	for _, id := range workerIDs {
		if availability[id] {
			available = append(available, id)
		}
	}

	return available, nil
}

// WorkerWithInfo 带信息的作业人员
type WorkerWithInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	IsAvailable bool   `json:"is_available"`
}

// GetAllWorkersWithAvailability 获取所有作业人员及其可用状态
func (s *CalendarService) GetAllWorkersWithAvailability(workDate time.Time, timeSlot string) ([]WorkerWithInfo, error) {
	// 获取作业人员和班组长角色的用户（支持中文名称）
	type UserInfo struct {
		ID     uint
		Name   string
		Avatar string
	}
	var users []UserInfo
	err := s.db.Table("users").
		Select("id, full_name as name, avatar").
		Where("is_active = ? AND (role = ? OR role = ? OR role = ? OR role = ?)",
			true,
			"worker",           // 英文作业人员
			"作业人员",          // 中文作业人员
			"team_leader",      // 英文班组长
			"班组长",            // 中文班组长
		).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	// 如果没有用户，返回空列表
	if len(users) == 0 {
		return []WorkerWithInfo{}, nil
	}

	// 提取所有用户ID
	userIDs := make([]uint, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	// 检查可用性
	availability, err := s.CheckMultipleAvailability(userIDs, workDate, timeSlot)
	if err != nil {
		return nil, err
	}

	// 构建结果，包含可用状态
	result := make([]WorkerWithInfo, len(users))
	for i, u := range users {
		result[i] = WorkerWithInfo{
			ID:          u.ID,
			Name:        u.Name,
			Avatar:      u.Avatar,
			IsAvailable: availability[u.ID],
		}
	}

	return result, nil
}

// GetWorkerAvailabilitySummary 获取作业人员可用性摘要
func (s *CalendarService) GetWorkerAvailabilitySummary(workerID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	// 获取日历
	calendars, err := s.GetWorkerCalendar(workerID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 统计
	totalSlots := len(calendars)
	availableSlots := 0
	busySlots := 0
	blockedSlots := 0

	for _, cal := range calendars {
		switch cal.Status {
		case CalendarStatusAvailable:
			availableSlots++
		case CalendarStatusBusy:
			busySlots++
		case CalendarStatusBlocked:
			blockedSlots++
		}
	}

	return map[string]interface{}{
		"total_slots":     totalSlots,
		"available_slots": availableSlots,
		"busy_slots":      busySlots,
		"blocked_slots":   blockedSlots,
		"availability_rate": float64(availableSlots) / float64(totalSlots) * 100,
	}, nil
}

// GetAppointmentsByDateRange 获取指定日期范围内的预约
func (s *CalendarService) GetAppointmentsByDateRange(startDate, endDate time.Time, workerID *uint) ([]ConstructionAppointment, error) {
	var appointments []ConstructionAppointment
	query := s.db.Where("work_date >= ? AND work_date <= ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	if workerID != nil {
		query = query.Where("assigned_worker_id = ?", *workerID)
	}

	err := query.Where("status IN ?", []string{StatusPending, StatusScheduled, StatusInProgress}).
		Order("work_date ASC, time_slot ASC, priority DESC").
		Find(&appointments).Error

	return appointments, err
}

// ValidateTimeSlot 验证时间段是否有效
func (s *CalendarService) ValidateTimeSlot(timeSlot string) error {
	validSlots := []string{TimeSlotMorning, TimeSlotNoon, TimeSlotAfternoon, TimeSlotFullDay}
	for _, slot := range validSlots {
		if timeSlot == slot {
			return nil
		}
	}
	return fmt.Errorf("invalid time_slot: %s", timeSlot)
}

// GetConflictingAppointments 获取冲突的预约
func (s *CalendarService) GetConflictingAppointments(workerID uint, workDate time.Time, timeSlot string, excludeID *uint) ([]ConstructionAppointment, error) {
	var appointments []ConstructionAppointment
	query := s.db.Where("assigned_worker_id = ? AND work_date = ? AND time_slot = ?",
		workerID, workDate.Format("2006-01-02"), timeSlot).
		Where("status IN ?", []string{StatusScheduled, StatusInProgress})

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Find(&appointments).Error
	return appointments, err
}

// GetAvailableWorkersCount 获取指定时间段可用作业人员数量
func (s *CalendarService) GetAvailableWorkersCount(workDateStr, timeSlot string) (int, error) {
	workDate, err := time.Parse("2006-01-02", workDateStr)
	if err != nil {
		return 0, fmt.Errorf("invalid work_date format: %w", err)
	}

	// 获取所有作业人员及其可用状态
	workers, err := s.GetAllWorkersWithAvailability(workDate, timeSlot)
	if err != nil {
		return 0, err
	}

	// 计算可用人数
	count := 0
	for _, w := range workers {
		if w.IsAvailable {
			count++
		}
	}

	return count, nil
}

// ClearAppointmentCalendar 清除预约单关联的所有日历记录
func (s *CalendarService) ClearAppointmentCalendar(appointmentID uint) error {
	// 查找该预约单关联的所有日历记录
	var calendars []WorkerCalendar
	if err := s.db.Where("appointment_id = ?", appointmentID).Find(&calendars).Error; err != nil {
		return err
	}

	// 释放所有关联的日历
	for _, cal := range calendars {
		if err := s.ReleaseCalendar(cal.WorkerID, cal.CalendarDate, cal.TimeSlot); err != nil {
			fmt.Printf("Warning: failed to release calendar for worker %d: %v\n", cal.WorkerID, err)
		}
	}

	return nil
}
