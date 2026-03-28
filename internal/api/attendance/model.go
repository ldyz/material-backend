package attendance

import (
	"time"

	"gorm.io/gorm"
)

// 本地时区（中国时区）
var localLocation *time.Location

func init() {
	var err error
	localLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 如果加载失败，使用固定偏移量（UTC+8）
		localLocation = time.FixedZone("CST", 8*3600)
	}
}

// formatTime 格式化时间为本地时区字符串（RFC3339格式）
func formatTime(t time.Time) string {
	return t.In(localLocation).Format(time.RFC3339)
}

// AttendanceRecord 打卡记录
type AttendanceRecord struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	UserID           uint            `gorm:"index;not null" json:"user_id"`
	AppointmentID    *uint           `gorm:"index" json:"appointment_id"`
	WorkContent      string          `gorm:"type:text" json:"work_content"`  // 手动填写的工作内容
	AttendanceType   string          `gorm:"size:20;not null" json:"attendance_type"`
	ClockInTime      time.Time       `gorm:"index;not null" json:"clock_in_time"`
	ClockInLocation  string          `gorm:"size:200" json:"clock_in_location"`
	ClockInLatitude  float64         `gorm:"type:decimal(10,7)" json:"clock_in_latitude"`
	ClockInLongitude float64         `gorm:"type:decimal(10,7)" json:"clock_in_longitude"`
	OvertimeHours    float64         `gorm:"type:decimal(4,1)" json:"overtime_hours"`
	Remark           string          `gorm:"type:text" json:"remark"`
	PhotoURL         string          `gorm:"size:500" json:"photo_url"`       // 打卡照片URL（单张，兼容旧数据）
	PhotoURLs        string          `gorm:"type:text" json:"photo_urls"`     // 打卡照片URLs（多张，JSON数组）
	Status           string          `gorm:"size:20;default:pending;index" json:"status"`
	ConfirmedBy      *uint           `json:"confirmed_by"`
	ConfirmedAt      *time.Time      `json:"confirmed_at"`
	ConfirmedRemark  string          `gorm:"type:text" json:"confirmed_remark"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`

	// 关联数据（不存储在数据库）
	UserName         string          `gorm:"-" json:"user_name,omitempty"`
	AppointmentNo    string          `gorm:"-" json:"appointment_no,omitempty"`
	ConfirmedByName  string          `gorm:"-" json:"confirmed_by_name,omitempty"`
}

// MonthlyAttendanceSummary 月度考勤汇总
type MonthlyAttendanceSummary struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	UserID             uint       `gorm:"index;not null" json:"user_id"`
	Year               int        `gorm:"not null" json:"year"`
	Month              int        `gorm:"not null" json:"month"`
	MorningCount       int        `gorm:"default:0" json:"morning_count"`
	AfternoonCount     int        `gorm:"default:0" json:"afternoon_count"`
	NoonOvertimeHours  float64    `gorm:"type:decimal(10,1);default:0" json:"noon_overtime_hours"`
	NightOvertimeHours float64    `gorm:"type:decimal(10,1);default:0" json:"night_overtime_hours"`
	TotalWorkDays      int        `gorm:"default:0" json:"total_work_days"`
	TotalOvertimeHours float64    `gorm:"type:decimal(10,1);default:0" json:"total_overtime_hours"`
	Status             string     `gorm:"size:20;default:draft" json:"status"`
	ConfirmedBy        *uint      `json:"confirmed_by"`
	ConfirmedAt        *time.Time `json:"confirmed_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// 关联数据
	UserName          string      `gorm:"-" json:"user_name,omitempty"`
	ConfirmedByName   string      `gorm:"-" json:"confirmed_by_name,omitempty"`
}

// TableName 指定表名
func (AttendanceRecord) TableName() string {
	return "attendance_records"
}

// TableName 指定表名
func (MonthlyAttendanceSummary) TableName() string {
	return "monthly_attendance_summary"
}

// 打卡类型常量
const (
	AttendanceTypeMorning       = "morning"
	AttendanceTypeAfternoon     = "afternoon"
	AttendanceTypeNoonOvertime  = "noon_overtime"
	AttendanceTypeNightOvertime = "night_overtime"
)

// 状态常量
const (
	RecordStatusPending   = "pending"
	RecordStatusConfirmed = "confirmed"
	RecordStatusRejected  = "rejected"
)

// 汇总状态常量
const (
	SummaryStatusDraft     = "draft"
	SummaryStatusConfirmed = "confirmed"
)

// ClockInRequest 打卡请求
type ClockInRequest struct {
	AppointmentID    *uint    `json:"appointment_id"`  // 可选，关联预约任务
	WorkContent      string   `json:"work_content"`    // 手动填写的工作内容（无任务时必填）
	AttendanceType   string   `json:"attendance_type" binding:"required,oneof=morning afternoon noon_overtime night_overtime"`
	ClockInLocation  string   `json:"clock_in_location"`
	ClockInLatitude  float64  `json:"clock_in_latitude"`
	ClockInLongitude float64  `json:"clock_in_longitude"`
	OvertimeHours    float64  `json:"overtime_hours"`  // 加班小时数
	Remark           string   `json:"remark"`
	PhotoURL         string   `json:"photo_url"`        // 打卡照片URL（单张，兼容旧版）
	PhotoURLs        []string `json:"photo_urls"`       // 打卡照片URLs（多张）
	ClockInTime      string   `json:"clock_in_time"`    // 补卡时间（格式：2006-01-02 15:04:05），为空则使用当前时间
}

// ConfirmRecordRequest 确认打卡记录请求
type ConfirmRecordRequest struct {
	Remark string `json:"remark"`
}

// RejectRecordRequest 驳回打卡记录请求
type RejectRecordRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// RecordListRequest 打卡记录列表请求
type RecordListRequest struct {
	Page           int    `form:"page,default=1" binding:"min=1"`
	PageSize       int    `form:"page_size,default=20" binding:"min=1,max=100"`
	UserID         *uint  `form:"user_id"`
	AttendanceType string `form:"attendance_type"`
	Status         string `form:"status"`
	StartDate      string `form:"start_date"`
	EndDate        string `form:"end_date"`
}

// MonthlySummaryRequest 月度汇总请求
type MonthlySummaryRequest struct {
	Year  int `form:"year" binding:"required"`
	Month int `form:"month" binding:"required,min=1,max=12"`
}

// TodayAppointment 今日待打卡任务
type TodayAppointment struct {
	ID             uint      `json:"id"`
	AppointmentNo  string    `json:"appointment_no"`
	WorkDate       string    `json:"work_date"`
	TimeSlot       string    `json:"time_slot"`
	WorkLocation   string    `json:"work_location"`
	WorkContent    string    `json:"work_content"`
	WorkType       string    `json:"work_type"`
	IsUrgent       bool      `json:"is_urgent"`
	Status         string    `json:"status"`
	// 已打卡信息
	ClockedTypes   []string  `json:"clocked_types"`  // 已打卡的类型
}

// AttendanceStatistics 考勤统计
type AttendanceStatistics struct {
	TotalRecords      int64   `json:"total_records"`
	PendingRecords    int64   `json:"pending_records"`
	ConfirmedRecords  int64   `json:"confirmed_records"`
	RejectedRecords   int64   `json:"rejected_records"`
	TodayClockIns     int64   `json:"today_clock_ins"`
	MonthClockIns     int64   `json:"month_clock_ins"`
	MonthOvertimeHours float64 `json:"month_overtime_hours"`
}

// ToDTO 转换为 DTO
func (r *AttendanceRecord) ToDTO() map[string]any {
	var confirmedAt string
	if r.ConfirmedAt != nil {
		confirmedAt = formatTime(*r.ConfirmedAt)
	}

	dto := map[string]any{
		"id":                 r.ID,
		"user_id":            r.UserID,
		"appointment_id":     r.AppointmentID,
		"attendance_type":    r.AttendanceType,
		"attendance_type_label": GetAttendanceTypeLabel(r.AttendanceType),
		"clock_in_time":      formatTime(r.ClockInTime),
		"clock_in_location":  r.ClockInLocation,
		"clock_in_latitude":  r.ClockInLatitude,
		"clock_in_longitude": r.ClockInLongitude,
		"overtime_hours":     r.OvertimeHours,
		"remark":             r.Remark,
		"status":             r.Status,
		"status_label":       GetStatusLabel(r.Status),
		"confirmed_by":       r.ConfirmedBy,
		"confirmed_at":       confirmedAt,
		"confirmed_remark":   r.ConfirmedRemark,
		"created_at":         formatTime(r.CreatedAt),
		"updated_at":         formatTime(r.UpdatedAt),
	}

	if r.UserName != "" {
		dto["user_name"] = r.UserName
	}
	if r.AppointmentNo != "" {
		dto["appointment_no"] = r.AppointmentNo
	}
	if r.WorkContent != "" {
		dto["work_content"] = r.WorkContent
	}
	if r.ConfirmedByName != "" {
		dto["confirmed_by_name"] = r.ConfirmedByName
	}
	if r.PhotoURL != "" {
		dto["photo_url"] = r.PhotoURL
	}
	if r.PhotoURLs != "" {
		dto["photo_urls"] = r.PhotoURLs
	}

	return dto
}

// ToDTO 转换为 DTO
func (s *MonthlyAttendanceSummary) ToDTO() map[string]any {
	var confirmedAt string
	if s.ConfirmedAt != nil {
		confirmedAt = formatTime(*s.ConfirmedAt)
	}

	dto := map[string]any{
		"id":                   s.ID,
		"user_id":              s.UserID,
		"year":                 s.Year,
		"month":                s.Month,
		"morning_count":        s.MorningCount,
		"afternoon_count":      s.AfternoonCount,
		"noon_overtime_hours":  s.NoonOvertimeHours,
		"night_overtime_hours": s.NightOvertimeHours,
		"total_work_days":      s.TotalWorkDays,
		"total_overtime_hours": s.TotalOvertimeHours,
		"status":               s.Status,
		"status_label":         getSummaryStatusLabel(s.Status),
		"confirmed_by":         s.ConfirmedBy,
		"confirmed_at":         confirmedAt,
		"created_at":           formatTime(s.CreatedAt),
		"updated_at":           formatTime(s.UpdatedAt),
	}

	if s.UserName != "" {
		dto["user_name"] = s.UserName
	}
	if s.ConfirmedByName != "" {
		dto["confirmed_by_name"] = s.ConfirmedByName
	}

	return dto
}

// GetAttendanceTypeLabel 获取打卡类型标签
func GetAttendanceTypeLabel(attendanceType string) string {
	labels := map[string]string{
		AttendanceTypeMorning:       "上午打卡",
		AttendanceTypeAfternoon:     "下午打卡",
		AttendanceTypeNoonOvertime:  "中午加班",
		AttendanceTypeNightOvertime: "晚上加班",
	}
	if label, ok := labels[attendanceType]; ok {
		return label
	}
	return attendanceType
}

// GetStatusLabel 获取状态标签
func GetStatusLabel(status string) string {
	labels := map[string]string{
		RecordStatusPending:   "待确认",
		RecordStatusConfirmed: "已确认",
		RecordStatusRejected:  "已驳回",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}

// getSummaryStatusLabel 获取汇总状态标签
func getSummaryStatusLabel(status string) string {
	labels := map[string]string{
		SummaryStatusDraft:     "草稿",
		SummaryStatusConfirmed: "已确认",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}

// IsOvertimeType 是否是加班类型
func IsOvertimeType(attendanceType string) bool {
	return attendanceType == AttendanceTypeNoonOvertime || attendanceType == AttendanceTypeNightOvertime
}

// BeforeCreate GORM hook
func (r *AttendanceRecord) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM hook
func (r *AttendanceRecord) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate GORM hook
func (s *MonthlyAttendanceSummary) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM hook
func (s *MonthlyAttendanceSummary) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}
