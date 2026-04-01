package attendance

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Service 打卡服务
type Service struct {
	db *gorm.DB
}

// NewService 创建服务
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// ClockIn 打卡
func (s *Service) ClockIn(userID uint, req ClockInRequest) (*AttendanceRecord, error) {
	// 解析打卡时间
	var clockInTime time.Time
	if req.ClockInTime != "" {
		// 补卡模式：解析指定时间（使用本地时区）
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.ClockInTime, localLocation)
		if err != nil {
			return nil, errors.New("打卡时间格式错误，应为：YYYY-MM-DD HH:mm:ss")
		}
		clockInTime = parsedTime
	} else {
		// 正常打卡：使用当前时间（本地时区）
		clockInTime = time.Now().In(localLocation)
	}

	// 如果选择了预约任务，验证任务是否存在且属于该用户
	if req.AppointmentID != nil && *req.AppointmentID > 0 {
		var appointmentCount int64
		s.db.Table("construction_appointments").
			Where("id = ? AND status IN ?", *req.AppointmentID, []string{"scheduled", "in_progress"}).
			Count(&appointmentCount)
		if appointmentCount == 0 {
			return nil, errors.New("预约任务不存在或状态不允许打卡")
		}

		// 检查用户是否被分配到该任务
		var isAssigned int64
		s.db.Table("construction_appointments").
			Where("id = ? AND (assigned_worker_id = ? OR assigned_worker_ids::jsonb @> ?::jsonb)",
				*req.AppointmentID, userID, "["+string(rune('0'+userID))+"]").
			Count(&isAssigned)
		if isAssigned == 0 {
			// 尝试另一种方式检查 JSON 数组
			var assignedWorkerIDs string
			s.db.Table("construction_appointments").
				Where("id = ?", *req.AppointmentID).
				Pluck("assigned_worker_ids", &assignedWorkerIDs)
			// 如果用户未被分配，返回错误
			if assignedWorkerIDs == "" {
				return nil, errors.New("您未被分配到此任务")
			}
		}
	} else if req.WorkContent == "" {
		// 没有选择任务时，必须填写工作内容
		return nil, errors.New("请选择任务或填写工作内容")
	}

	// 如果是加班打卡，检查是否已完成当天正常打卡
	if IsOvertimeType(req.AttendanceType) {
		hasNormalClockIn, err := s.hasNormalClockInOnDate(userID, clockInTime.Format("2006-01-02"))
		if err != nil {
			return nil, err
		}
		if !hasNormalClockIn {
			return nil, errors.New("请先完成上午或下午打卡后再打加班卡")
		}
	}

	// 检查该日期该类型是否已打卡
	targetDate := clockInTime.Format("2006-01-02")
	var existingRecord AttendanceRecord
	err := s.db.Where("user_id = ? AND attendance_type = ? AND DATE(clock_in_time) = ?",
		userID, req.AttendanceType, targetDate).First(&existingRecord).Error
	if err == nil {
		return nil, errors.New("该日期已完成该类型打卡")
	}

	// 验证加班小时数
	if IsOvertimeType(req.AttendanceType) && req.OvertimeHours <= 0 {
		return nil, errors.New("加班打卡需要输入加班小时数")
	}

	// 创建打卡记录
	record := &AttendanceRecord{
		UserID:           userID,
		AppointmentID:    req.AppointmentID,
		WorkContent:      req.WorkContent,
		AttendanceType:   req.AttendanceType,
		ClockInTime:      clockInTime,
		ClockInLocation:  req.ClockInLocation,
		ClockInLatitude:  req.ClockInLatitude,
		ClockInLongitude: req.ClockInLongitude,
		OvertimeHours:    req.OvertimeHours,
		Remark:           req.Remark,
		PhotoURL:         req.PhotoURL,
		Status:           RecordStatusPending,
	}

	// 处理多张照片URL
	if len(req.PhotoURLs) > 0 {
		photoURLsJSON, err := json.Marshal(req.PhotoURLs)
		if err == nil {
			record.PhotoURLs = string(photoURLsJSON)
		}
		// 兼容旧版：第一张照片也存入 photo_url
		if req.PhotoURL == "" && len(req.PhotoURLs) > 0 {
			record.PhotoURL = req.PhotoURLs[0]
		}
	}

	if err := s.db.Create(record).Error; err != nil {
		return nil, err
	}

	return record, nil
}

// hasNormalClockInToday 检查今天是否有正常打卡
func (s *Service) hasNormalClockInToday(userID uint) (bool, error) {
	today := time.Now().Format("2006-01-02")
	return s.hasNormalClockInOnDate(userID, today)
}

// hasNormalClockInOnDate 检查指定日期是否有正常打卡
func (s *Service) hasNormalClockInOnDate(userID uint, date string) (bool, error) {
	var count int64
	err := s.db.Model(&AttendanceRecord{}).
		Where("user_id = ? AND attendance_type IN ? AND DATE(clock_in_time) = ?",
			userID, []string{AttendanceTypeMorning, AttendanceTypeAfternoon}, date).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetTodayAppointments 获取今日待打卡任务
func (s *Service) GetTodayAppointments(userID uint) ([]TodayAppointment, error) {
	today := time.Now().Format("2006-01-02")

	// 查询今天分配给该用户的预约任务
	query := `
		SELECT
			ca.id, ca.appointment_no, ca.work_date, ca.time_slot,
			ca.work_location, ca.work_content, ca.work_type,
			ca.is_urgent, ca.status
		FROM construction_appointments ca
		WHERE ca.work_date = ?
		AND ca.status IN ('scheduled', 'in_progress')
		AND (
			ca.assigned_worker_id = ?
			OR ca.assigned_worker_ids::jsonb @> ?::jsonb
		)
	`

	var appointments []TodayAppointment
	rows, err := s.db.Raw(query, today, userID, "["+string(rune('0'+userID))+"]").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var apt TodayAppointment
		var workDate time.Time
		err := rows.Scan(
			&apt.ID, &apt.AppointmentNo, &workDate, &apt.TimeSlot,
			&apt.WorkLocation, &apt.WorkContent, &apt.WorkType,
			&apt.IsUrgent, &apt.Status,
		)
		if err != nil {
			continue
		}
		apt.WorkDate = workDate.Format("2006-01-02")

		// 获取该任务的已打卡类型
		var records []AttendanceRecord
		s.db.Where("appointment_id = ? AND user_id = ? AND DATE(clock_in_time) = ?",
			apt.ID, userID, today).
			Find(&records)
		apt.ClockedTypes = make([]string, 0)
		for _, r := range records {
			apt.ClockedTypes = append(apt.ClockedTypes, r.AttendanceType)
		}

		appointments = append(appointments, apt)
	}

	// 如果上面的查询没有结果，尝试另一种方式
	if len(appointments) == 0 {
		appointments = s.getTodayAppointmentsAlternative(userID, today)
	}

	return appointments, nil
}

// getTodayAppointmentsAlternative 备用查询方法
func (s *Service) getTodayAppointmentsAlternative(userID uint, today string) []TodayAppointment {
	var appointments []TodayAppointment

	// 查询所有今天的预约任务
	rows, err := s.db.Table("construction_appointments").
		Select("id, appointment_no, work_date, time_slot, work_location, work_content, work_type, is_urgent, status, assigned_worker_ids, assigned_worker_id").
		Where("work_date = ? AND status IN ?", today, []string{"scheduled", "in_progress"}).
		Rows()
	if err != nil {
		return appointments
	}
	defer rows.Close()

	for rows.Next() {
		var id, assignedWorkerID uint
		var appointmentNo, timeSlot, workLocation, workContent, workType, status string
		var workDate time.Time
		var isUrgent bool
		var assignedWorkerIDsStr string

		err := rows.Scan(&id, &appointmentNo, &workDate, &timeSlot, &workLocation, &workContent, &workType, &isUrgent, &status, &assignedWorkerIDsStr, &assignedWorkerID)
		if err != nil {
			continue
		}

		// 检查用户是否被分配到该任务
		isAssigned := false
		if assignedWorkerID == userID {
			isAssigned = true
		} else if assignedWorkerIDsStr != "" {
			// 检查 JSON 数组中是否包含该用户ID
			var ids []uint
			if err := parseWorkerIDs(assignedWorkerIDsStr, &ids); err == nil {
				for _, uid := range ids {
					if uid == userID {
						isAssigned = true
						break
					}
				}
			}
		}

		if !isAssigned {
			continue
		}

		apt := TodayAppointment{
			ID:            id,
			AppointmentNo: appointmentNo,
			WorkDate:      workDate.Format("2006-01-02"),
			TimeSlot:      timeSlot,
			WorkLocation:  workLocation,
			WorkContent:   workContent,
			WorkType:      workType,
			IsUrgent:      isUrgent,
			Status:        status,
			ClockedTypes:  make([]string, 0),
		}

		// 获取该任务的已打卡类型
		var records []AttendanceRecord
		s.db.Where("appointment_id = ? AND user_id = ? AND DATE(clock_in_time) = ?",
			id, userID, today).
			Find(&records)
		for _, r := range records {
			apt.ClockedTypes = append(apt.ClockedTypes, r.AttendanceType)
		}

		appointments = append(appointments, apt)
	}

	return appointments
}

// parseWorkerIDs 解析作业人员ID JSON数组
func parseWorkerIDs(jsonStr string, ids *[]uint) error {
	// 简单解析 JSON 数组
	if jsonStr == "" || jsonStr == "[]" {
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), ids)
}

// GetMyRecords 获取我的打卡记录
func (s *Service) GetMyRecords(userID uint, req RecordListRequest) ([]AttendanceRecord, int64, error) {
	query := s.db.Model(&AttendanceRecord{}).Where("user_id = ?", userID)

	if req.AttendanceType != "" {
		query = query.Where("attendance_type = ?", req.AttendanceType)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StartDate != "" {
		query = query.Where("DATE(clock_in_time) >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("DATE(clock_in_time) <= ?", req.EndDate)
	}

	var total int64
	query.Count(&total)

	var records []AttendanceRecord
	offset := (req.Page - 1) * req.PageSize
	err := query.Order("clock_in_time DESC").
		Offset(offset).Limit(req.PageSize).
		Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	// 填充关联数据
	s.fillRecordDetails(records)

	return records, total, nil
}

// GetRecords 获取打卡记录列表（管理员）
func (s *Service) GetRecords(req RecordListRequest) ([]AttendanceRecord, int64, error) {
	query := s.db.Model(&AttendanceRecord{})

	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.AppointmentID != nil {
		query = query.Where("appointment_id = ?", *req.AppointmentID)
	}
	if req.ProjectID != nil {
		// 通过预约任务关联项目
		query = query.Where("appointment_id IN (SELECT id FROM construction_appointments WHERE project_id = ?)", *req.ProjectID)
	}
	if req.AttendanceType != "" {
		query = query.Where("attendance_type = ?", req.AttendanceType)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StartDate != "" {
		query = query.Where("DATE(clock_in_time) >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("DATE(clock_in_time) <= ?", req.EndDate)
	}

	var total int64
	query.Count(&total)

	var records []AttendanceRecord
	offset := (req.Page - 1) * req.PageSize
	err := query.Order("clock_in_time DESC").
		Offset(offset).Limit(req.PageSize).
		Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	// 填充关联数据
	s.fillRecordDetails(records)

	return records, total, nil
}

// fillRecordDetails 填充打卡记录详情
func (s *Service) fillRecordDetails(records []AttendanceRecord) {
	for i := range records {
		// 获取用户名
		s.db.Table("users").Where("id = ?", records[i].UserID).
			Pluck("full_name", &records[i].UserName)

		// 获取预约任务信息
		if records[i].AppointmentID != nil {
			s.db.Table("construction_appointments").Where("id = ?", *records[i].AppointmentID).
				Pluck("appointment_no", &records[i].AppointmentNo)
			// 如果打卡记录没有工作内容，从预约任务获取
			if records[i].WorkContent == "" {
				s.db.Table("construction_appointments").Where("id = ?", *records[i].AppointmentID).
					Pluck("work_content", &records[i].WorkContent)
			}
		}

		// 获取确认人姓名
		if records[i].ConfirmedBy != nil {
			s.db.Table("users").Where("id = ?", *records[i].ConfirmedBy).
				Pluck("full_name", &records[i].ConfirmedByName)
		}
	}
}

// GetRecordByID 获取打卡记录详情
func (s *Service) GetRecordByID(id uint) (*AttendanceRecord, error) {
	var record AttendanceRecord
	err := s.db.First(&record, id).Error
	if err != nil {
		return nil, err
	}

	// 填充关联数据
	s.fillRecordDetails([]AttendanceRecord{record})

	return &record, nil
}

// ConfirmRecord 确认打卡记录
func (s *Service) ConfirmRecord(id uint, confirmedBy uint, remark string) error {
	var record AttendanceRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return errors.New("打卡记录不存在")
	}

	if record.Status != RecordStatusPending {
		return errors.New("该记录已处理")
	}

	now := time.Now()
	updates := map[string]any{
		"status":           RecordStatusConfirmed,
		"confirmed_by":     confirmedBy,
		"confirmed_at":     now,
		"confirmed_remark": remark,
		"updated_at":       now,
	}

	return s.db.Model(&record).Updates(updates).Error
}

// RejectRecord 驳回打卡记录
func (s *Service) RejectRecord(id uint, confirmedBy uint, reason string) error {
	var record AttendanceRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return errors.New("打卡记录不存在")
	}

	if record.Status != RecordStatusPending {
		return errors.New("该记录已处理")
	}

	now := time.Now()
	updates := map[string]any{
		"status":           RecordStatusRejected,
		"confirmed_by":     confirmedBy,
		"confirmed_at":     now,
		"confirmed_remark": reason,
		"updated_at":       now,
	}

	return s.db.Model(&record).Updates(updates).Error
}

// GetMonthlySummary 获取月度考勤汇总
func (s *Service) GetMonthlySummary(year, month int, userID *uint) ([]MonthlyAttendanceSummary, error) {
	query := s.db.Model(&MonthlyAttendanceSummary{}).
		Where("year = ? AND month = ?", year, month)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	var summaries []MonthlyAttendanceSummary
	err := query.Find(&summaries).Error
	if err != nil {
		return nil, err
	}

	// 填充用户名
	for i := range summaries {
		s.db.Table("users").Where("id = ?", summaries[i].UserID).
			Pluck("full_name", &summaries[i].UserName)
		if summaries[i].ConfirmedBy != nil {
			s.db.Table("users").Where("id = ?", *summaries[i].ConfirmedBy).
				Pluck("full_name", &summaries[i].ConfirmedByName)
		}
	}

	return summaries, nil
}

// GenerateMonthlySummary 生成月度考勤汇总
func (s *Service) GenerateMonthlySummary(year, month int) error {
	// 获取该月所有有打卡记录的用户
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	var userIDs []uint
	s.db.Model(&AttendanceRecord{}).
		Where("clock_in_time >= ? AND clock_in_time <= ? AND status = ?",
			startDate, endDate, RecordStatusConfirmed).
		Distinct("user_id").
		Pluck("user_id", &userIDs)

	for _, userID := range userIDs {
		// 检查是否已存在汇总
		var existing MonthlyAttendanceSummary
		err := s.db.Where("user_id = ? AND year = ? AND month = ?", userID, year, month).
			First(&existing).Error
		if err == nil && existing.Status == SummaryStatusConfirmed {
			continue // 已确认的汇总跳过
		}

		// 统计各类打卡
		var morningCount, afternoonCount int64
		var noonOvertimeHours, nightOvertimeHours float64

		s.db.Model(&AttendanceRecord{}).
			Where("user_id = ? AND attendance_type = ? AND clock_in_time >= ? AND clock_in_time <= ? AND status = ?",
				userID, AttendanceTypeMorning, startDate, endDate, RecordStatusConfirmed).
			Count(&morningCount)

		s.db.Model(&AttendanceRecord{}).
			Where("user_id = ? AND attendance_type = ? AND clock_in_time >= ? AND clock_in_time <= ? AND status = ?",
				userID, AttendanceTypeAfternoon, startDate, endDate, RecordStatusConfirmed).
			Count(&afternoonCount)

		s.db.Model(&AttendanceRecord{}).
			Where("user_id = ? AND attendance_type = ? AND clock_in_time >= ? AND clock_in_time <= ? AND status = ?",
				userID, AttendanceTypeNoonOvertime, startDate, endDate, RecordStatusConfirmed).
			Select("COALESCE(SUM(overtime_hours), 0)").Scan(&noonOvertimeHours)

		s.db.Model(&AttendanceRecord{}).
			Where("user_id = ? AND attendance_type = ? AND clock_in_time >= ? AND clock_in_time <= ? AND status = ?",
				userID, AttendanceTypeNightOvertime, startDate, endDate, RecordStatusConfirmed).
			Select("COALESCE(SUM(overtime_hours), 0)").Scan(&nightOvertimeHours)

		// 计算工作天数（有上午或下午打卡的天数）
		var workDays int64
		s.db.Raw(`
			SELECT COUNT(DISTINCT DATE(clock_in_time))
			FROM attendance_records
			WHERE user_id = ?
			AND attendance_type IN ('morning', 'afternoon')
			AND clock_in_time >= ? AND clock_in_time <= ?
			AND status = 'confirmed'
		`, userID, startDate, endDate).Scan(&workDays)

		summary := MonthlyAttendanceSummary{
			UserID:             userID,
			Year:               year,
			Month:              month,
			MorningCount:       int(morningCount),
			AfternoonCount:     int(afternoonCount),
			NoonOvertimeHours:  noonOvertimeHours,
			NightOvertimeHours: nightOvertimeHours,
			TotalWorkDays:      int(workDays),
			TotalOvertimeHours: noonOvertimeHours + nightOvertimeHours,
			Status:             SummaryStatusDraft,
		}

		if existing.ID > 0 {
			s.db.Model(&existing).Updates(map[string]any{
				"morning_count":        summary.MorningCount,
				"afternoon_count":      summary.AfternoonCount,
				"noon_overtime_hours":  summary.NoonOvertimeHours,
				"night_overtime_hours": summary.NightOvertimeHours,
				"total_work_days":      summary.TotalWorkDays,
				"total_overtime_hours": summary.TotalOvertimeHours,
				"updated_at":           time.Now(),
			})
		} else {
			s.db.Create(&summary)
		}
	}

	return nil
}

// GetMyMonthlySummary 获取我的月度考勤汇总
func (s *Service) GetMyMonthlySummary(userID uint, year, month int) (*MonthlyAttendanceSummary, error) {
	var summary MonthlyAttendanceSummary
	err := s.db.Where("user_id = ? AND year = ? AND month = ?", userID, year, month).
		First(&summary).Error
	if err != nil {
		// 如果不存在，实时计算
		s.GenerateMonthlySummary(year, month)
		err = s.db.Where("user_id = ? AND year = ? AND month = ?", userID, year, month).
			First(&summary).Error
		if err != nil {
			return nil, err
		}
	}

	// 填充用户名
	s.db.Table("users").Where("id = ?", summary.UserID).
		Pluck("full_name", &summary.UserName)

	return &summary, nil
}

// GetStatistics 获取考勤统计
func (s *Service) GetStatistics() (*AttendanceStatistics, error) {
	stats := &AttendanceStatistics{}

	s.db.Model(&AttendanceRecord{}).Count(&stats.TotalRecords)
	s.db.Model(&AttendanceRecord{}).Where("status = ?", RecordStatusPending).Count(&stats.PendingRecords)
	s.db.Model(&AttendanceRecord{}).Where("status = ?", RecordStatusConfirmed).Count(&stats.ConfirmedRecords)
	s.db.Model(&AttendanceRecord{}).Where("status = ?", RecordStatusRejected).Count(&stats.RejectedRecords)

	today := time.Now().Format("2006-01-02")
	s.db.Model(&AttendanceRecord{}).Where("DATE(clock_in_time) = ?", today).Count(&stats.TodayClockIns)

	// 本月统计
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	s.db.Model(&AttendanceRecord{}).
		Where("clock_in_time >= ?", startOfMonth).
		Count(&stats.MonthClockIns)

	s.db.Model(&AttendanceRecord{}).
		Where("clock_in_time >= ? AND attendance_type IN ?", startOfMonth,
			[]string{AttendanceTypeNoonOvertime, AttendanceTypeNightOvertime}).
		Select("COALESCE(SUM(overtime_hours), 0)").Scan(&stats.MonthOvertimeHours)

	return stats, nil
}

// ConfirmMonthlySummary 确认月度汇总
func (s *Service) ConfirmMonthlySummary(id uint, confirmedBy uint) error {
	var summary MonthlyAttendanceSummary
	if err := s.db.First(&summary, id).Error; err != nil {
		return errors.New("汇总记录不存在")
	}

	if summary.Status == SummaryStatusConfirmed {
		return errors.New("该汇总已确认")
	}

	now := time.Now()
	return s.db.Model(&summary).Updates(map[string]any{
		"status":       SummaryStatusConfirmed,
		"confirmed_by": confirmedBy,
		"confirmed_at": now,
		"updated_at":   now,
	}).Error
}

// DailyStatistics 每日统计
type DailyStatistics struct {
	Date           string  `json:"date"`
	MorningCount   int     `json:"morning_count"`
	AfternoonCount int     `json:"afternoon_count"`
	OvertimeHours  float64 `json:"overtime_hours"`
	TotalCount     int     `json:"total_count"`
}

// GetCalendarStatistics 获取打卡日历统计
func (s *Service) GetCalendarStatistics(userID uint, startDate, endDate string) ([]DailyStatistics, error) {
	var records []AttendanceRecord
	err := s.db.Where("user_id = ? AND DATE(clock_in_time) >= ? AND DATE(clock_in_time) <= ?",
		userID, startDate, endDate).
		Order("clock_in_time").
		Find(&records).Error
	if err != nil {
		return nil, err
	}

	// 按日期分组统计
	statsMap := make(map[string]*DailyStatistics)

	for _, record := range records {
		dateStr := record.ClockInTime.Format("2006-01-02")
		if _, exists := statsMap[dateStr]; !exists {
			statsMap[dateStr] = &DailyStatistics{
				Date:           dateStr,
				MorningCount:   0,
				AfternoonCount: 0,
				OvertimeHours:  0,
				TotalCount:     0,
			}
		}

		statsMap[dateStr].TotalCount++
		if record.AttendanceType == AttendanceTypeMorning {
			statsMap[dateStr].MorningCount++
		} else if record.AttendanceType == AttendanceTypeAfternoon {
			statsMap[dateStr].AfternoonCount++
		} else if IsOvertimeType(record.AttendanceType) {
			statsMap[dateStr].OvertimeHours += record.OvertimeHours
		}
	}

	// 转换为数组
	result := make([]DailyStatistics, 0, len(statsMap))
	for _, stat := range statsMap {
		result = append(result, *stat)
	}

	return result, nil
}

// GetDailyStatistics 按日期统计
func (s *Service) GetDailyStatistics(startDate, endDate string) ([]DailyStatisticsResponse, error) {
	query := `
		SELECT
			DATE(clock_in_time) as date,
			COUNT(*) as total_count,
			COUNT(*) FILTER (WHERE attendance_type = 'morning') as morning_count,
			COUNT(*) FILTER (WHERE attendance_type = 'afternoon') as afternoon_count,
			COUNT(*) FILTER (WHERE attendance_type = 'noon_overtime') as noon_overtime,
			COUNT(*) FILTER (WHERE attendance_type = 'night_overtime') as night_overtime,
			COALESCE(SUM(overtime_hours), 0) as overtime_hours,
			COUNT(DISTINCT user_id) as user_count
		FROM attendance_records
		WHERE DATE(clock_in_time) >= ? AND DATE(clock_in_time) <= ?
		GROUP BY DATE(clock_in_time)
		ORDER BY date DESC
	`

	var results []DailyStatisticsResponse
	err := s.db.Raw(query, startDate, endDate).Scan(&results).Error
	return results, err
}

// GetUserStatistics 按人员统计
func (s *Service) GetUserStatistics(startDate, endDate string) ([]UserStatisticsResponse, error) {
	query := `
		SELECT
			ar.user_id,
			u.full_name as user_name,
			COUNT(*) as total_count,
			COUNT(*) FILTER (WHERE ar.attendance_type = 'morning') as morning_count,
			COUNT(*) FILTER (WHERE ar.attendance_type = 'afternoon') as afternoon_count,
			COUNT(*) FILTER (WHERE ar.attendance_type = 'noon_overtime') as noon_overtime_count,
			COUNT(*) FILTER (WHERE ar.attendance_type = 'night_overtime') as night_overtime_count,
			COALESCE(SUM(ar.overtime_hours), 0) as total_overtime_hours,
			COUNT(DISTINCT DATE(ar.clock_in_time)) as work_days
		FROM attendance_records ar
		LEFT JOIN users u ON ar.user_id = u.id
		WHERE DATE(ar.clock_in_time) >= ? AND DATE(ar.clock_in_time) <= ?
		GROUP BY ar.user_id, u.full_name
		ORDER BY total_count DESC
	`

	var results []UserStatisticsResponse
	err := s.db.Raw(query, startDate, endDate).Scan(&results).Error
	return results, err
}

// GetTaskStatistics 按任务统计
func (s *Service) GetTaskStatistics(startDate, endDate string) ([]TaskStatisticsResponse, error) {
	query := `
		SELECT
			ca.id as appointment_id,
			ca.appointment_no,
			COALESCE(ca.work_content, '') as work_content,
			TO_CHAR(ca.work_date, 'YYYY-MM-DD') as work_date,
			COUNT(ar.id) as total_clock_ins,
			COUNT(DISTINCT ar.user_id) as unique_users,
			SUM(CASE
				WHEN ar.photo_urls IS NOT NULL AND ar.photo_urls != '' THEN
					jsonb_array_length(CASE WHEN ar.photo_urls::jsonb IS NOT NULL THEN ar.photo_urls::jsonb ELSE '[]'::jsonb END)
				WHEN ar.photo_url IS NOT NULL AND ar.photo_url != '' THEN 1
				ELSE 0
			END) as photo_count,
			COUNT(*) FILTER (WHERE ar.attendance_type = 'morning') as morning_count,
			COUNT(*) FILTER (WHERE ar.attendance_type = 'afternoon') as afternoon_count,
			COUNT(*) FILTER (WHERE ar.attendance_type IN ('noon_overtime', 'night_overtime')) as overtime_count,
			COALESCE(SUM(ar.overtime_hours), 0) as overtime_hours
		FROM construction_appointments ca
		LEFT JOIN attendance_records ar ON ca.id = ar.appointment_id
			AND DATE(ar.clock_in_time) >= ? AND DATE(ar.clock_in_time) <= ?
		GROUP BY ca.id, ca.appointment_no, ca.work_content, ca.work_date
		HAVING COUNT(ar.id) > 0
		ORDER BY ca.work_date DESC
	`

	var results []TaskStatisticsResponse
	err := s.db.Raw(query, startDate, endDate).Scan(&results).Error
	return results, err
}

// GetProjectStatistics 按项目统计
func (s *Service) GetProjectStatistics(startDate, endDate string) ([]ProjectStatisticsResponse, error) {
	query := `
		SELECT
			p.id as project_id,
			p.name as project_name,
			COUNT(ar.id) as total_clock_ins,
			COUNT(DISTINCT ar.user_id) as unique_users,
			COUNT(*) FILTER (WHERE ar.attendance_type = 'morning') as morning_count,
			COUNT(*) FILTER (WHERE ar.attendance_type = 'afternoon') as afternoon_count,
			COALESCE(SUM(ar.overtime_hours), 0) as total_overtime_hours,
			COUNT(DISTINCT DATE(ar.clock_in_time)) as work_days
		FROM projects p
		LEFT JOIN construction_appointments ca ON p.id = ca.project_id
		LEFT JOIN attendance_records ar ON ca.id = ar.appointment_id
			AND DATE(ar.clock_in_time) >= ? AND DATE(ar.clock_in_time) <= ?
		GROUP BY p.id, p.name
		HAVING COUNT(ar.id) > 0
		ORDER BY total_clock_ins DESC
	`

	var results []ProjectStatisticsResponse
	err := s.db.Raw(query, startDate, endDate).Scan(&results).Error
	return results, err
}

// GetRecordsForExport 获取用于导出的打卡记录
func (s *Service) GetRecordsForExport(req ExportRequest) ([]AttendanceRecord, error) {
	query := s.db.Model(&AttendanceRecord{})

	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.AppointmentID != nil {
		query = query.Where("appointment_id = ?", *req.AppointmentID)
	}
	if req.ProjectID != nil {
		query = query.Where("appointment_id IN (SELECT id FROM construction_appointments WHERE project_id = ?)", *req.ProjectID)
	}
	if req.StartDate != "" {
		query = query.Where("DATE(clock_in_time) >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("DATE(clock_in_time) <= ?", req.EndDate)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	var records []AttendanceRecord
	err := query.Order("clock_in_time DESC").Find(&records).Error
	if err != nil {
		return nil, err
	}

	// 填充关联数据
	s.fillRecordDetails(records)

	return records, nil
}

// GetPhotosForExport 获取用于导出的照片URL列表
func (s *Service) GetPhotosForExport(req ExportRequest) ([]map[string]string, []AttendanceRecord, error) {
	records, err := s.GetRecordsForExport(req)
	if err != nil {
		return nil, nil, err
	}

	var photos []map[string]string
	for _, record := range records {
		// 解析照片URLs
		var urls []string
		if record.PhotoURLs != "" {
			json.Unmarshal([]byte(record.PhotoURLs), &urls)
		}
		if record.PhotoURL != "" && !contains(urls, record.PhotoURL) {
			urls = append(urls, record.PhotoURL)
		}

		for i, url := range urls {
			if url == "" || url == "null" {
				continue
			}
			photos = append(photos, map[string]string{
				"url":       url,
				"filename":  fmt.Sprintf("%s_%s_%s_%d.jpg",
					record.ClockInTime.Format("2006-01-02"),
					record.UserName,
					record.AttendanceType,
					i+1),
				"record_id": fmt.Sprintf("%d", record.ID),
			})
		}
	}

	return photos, records, nil
}

// contains 检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
