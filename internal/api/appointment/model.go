package appointment

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ConstructionAppointment 施工预约单
type ConstructionAppointment struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	AppointmentNo        string    `gorm:"size:50;uniqueIndex" json:"appointment_no"`
	ProjectID            *uint     `gorm:"index" json:"project_id"`
	ApplicantID          uint      `gorm:"index" json:"applicant_id"`
	ApplicantName        string    `gorm:"size:100" json:"applicant_name"`
	ContactPhone         string    `gorm:"size:20" json:"contact_phone"`
	ContactPerson        string    `gorm:"size:100" json:"contact_person"`
	WorkDate             time.Time `gorm:"index;not null" json:"work_date"`
	TimeSlot             string    `gorm:"size:50;not null" json:"time_slot"`
	WorkLocation         string    `gorm:"size:500;not null" json:"work_location"`
	WorkContent          string    `gorm:"type:text;not null" json:"work_content"`
	WorkType             string    `gorm:"size:50" json:"work_type"`
	IsUrgent             bool      `gorm:"default:false;index" json:"is_urgent"`
	Priority             int       `gorm:"default:0" json:"priority"`
	UrgentReason         string    `gorm:"type:text" json:"urgent_reason"`
	AssignedWorkerID     *uint     `gorm:"index" json:"assigned_worker_id"`        // 保留用于兼容性
	AssignedWorkerName   string    `gorm:"size:100" json:"assigned_worker_name"`    // 保留用于兼容性
	AssignedWorkerIDs    string    `gorm:"type:text" json:"assigned_worker_ids"`    // JSON数组格式: [1,2,3]
	AssignedWorkerNames  string    `gorm:"type:text" json:"assigned_worker_names"`  // 逗号分隔: 张三,李四,王五
	SupervisorID         *uint     `gorm:"index" json:"supervisor_id"`           // 监护人ID
	SupervisorName       string    `gorm:"size:100" json:"supervisor_name"`       // 监护人姓名
	Status               string    `gorm:"size:20;default:draft;index" json:"status"`
	WorkflowInstanceID   *uint     `gorm:"index" json:"workflow_instance_id"`
	SubmittedAt          *time.Time `json:"submitted_at"`
	ApprovedAt           *time.Time `json:"approved_at"`
	CompletedAt          *time.Time `json:"completed_at"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// WorkerCalendar 作业人员日历
type WorkerCalendar struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	WorkerID       uint           `gorm:"index;not null" json:"worker_id"`
	CalendarDate   time.Time      `gorm:"index;not null" json:"calendar_date"`
	TimeSlot       string         `gorm:"size:20;not null" json:"time_slot"`
	IsAvailable    bool           `gorm:"default:true" json:"is_available"`
	Status         string         `gorm:"size:20;default:available" json:"status"`
	AppointmentID  *uint          `gorm:"index" json:"appointment_id"`
	BlockedReason  string         `gorm:"type:text" json:"blocked_reason"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// TableName 指定表名
func (ConstructionAppointment) TableName() string {
	return "construction_appointments"
}

// TableName 指定表名
func (WorkerCalendar) TableName() string {
	return "worker_calendars"
}

// ToDTO 转换为 DTO
func (a *ConstructionAppointment) ToDTO() map[string]any {
	var submittedAt, approvedAt, completedAt string
	if a.SubmittedAt != nil {
		submittedAt = a.SubmittedAt.Format(time.RFC3339)
	}
	if a.ApprovedAt != nil {
		approvedAt = a.ApprovedAt.Format(time.RFC3339)
	}
	if a.CompletedAt != nil {
		completedAt = a.CompletedAt.Format(time.RFC3339)
	}

	// 解析作业人员ID列表
	var assignedWorkerIDs []uint
	if a.AssignedWorkerIDs != "" {
		json.Unmarshal([]byte(a.AssignedWorkerIDs), &assignedWorkerIDs)
	}

	return map[string]any{
		"id":                    a.ID,
		"appointment_no":        a.AppointmentNo,
		"project_id":            a.ProjectID,
		"applicant_id":          a.ApplicantID,
		"applicant_name":        a.ApplicantName,
		"contact_phone":         a.ContactPhone,
		"contact_person":        a.ContactPerson,
		"work_date":             a.WorkDate.Format("2006-01-02"),
		"time_slot":             a.TimeSlot,
		"work_location":         a.WorkLocation,
		"work_content":          a.WorkContent,
		"work_type":             a.WorkType,
		"is_urgent":             a.IsUrgent,
		"priority":              a.Priority,
		"urgent_reason":         a.UrgentReason,
		"assigned_worker_id":    a.AssignedWorkerID,
		"assigned_worker_name":  a.AssignedWorkerName,
		"assigned_worker_ids":   assignedWorkerIDs,
		"assigned_worker_names": a.AssignedWorkerNames,
		"status":                a.Status,
		"workflow_instance_id":  a.WorkflowInstanceID,
		"submitted_at":          submittedAt,
		"approved_at":           approvedAt,
		"completed_at":          completedAt,
		"created_at":            a.CreatedAt.Format(time.RFC3339),
		"updated_at":            a.UpdatedAt.Format(time.RFC3339),
	}
}

// ToDTO 转换为 DTO
func (w *WorkerCalendar) ToDTO() map[string]any {
	return map[string]any{
		"id":             w.ID,
		"worker_id":      w.WorkerID,
		"calendar_date":  w.CalendarDate.Format("2006-01-02"),
		"time_slot":      w.TimeSlot,
		"is_available":   w.IsAvailable,
		"status":         w.Status,
		"appointment_id": w.AppointmentID,
		"blocked_reason": w.BlockedReason,
		"created_at":     w.CreatedAt.Format(time.RFC3339),
		"updated_at":     w.UpdatedAt.Format(time.RFC3339),
	}
}

// CreateAppointmentRequest 创建预约单请求
type CreateAppointmentRequest struct {
	ProjectID          *uint  `json:"project_id"`
	ContactPhone       string `json:"contact_phone"`
	ContactPerson      string `json:"contact_person"`
	WorkDate           string `json:"work_date" binding:"required"`
	TimeSlot           string `json:"time_slot" binding:"required,oneof=morning noon afternoon full_day"`
	WorkLocation       string `json:"work_location" binding:"required"`
	WorkContent        string `json:"work_content" binding:"required"`
	WorkType           string `json:"work_type"`
	IsUrgent           bool   `json:"is_urgent"`
	Priority           int    `json:"priority" binding:"min=0,max=10"`
	UrgentReason       string `json:"urgent_reason"`
	AssignedWorkerID   *uint  `json:"assigned_worker_id"`
	AssignedWorkerIDs  string `json:"assigned_worker_ids"`  // JSON数组格式
	AssignedWorkerNames string `json:"assigned_worker_names"` // 逗号分隔
}

// UpdateAppointmentRequest 更新预约单请求
type UpdateAppointmentRequest struct {
	ProjectID           *uint   `json:"project_id"`
	ContactPhone        string  `json:"contact_phone"`
	ContactPerson       string  `json:"contact_person"`
	WorkDate            string  `json:"work_date"`
	TimeSlot            string  `json:"time_slot" binding:"omitempty,oneof=morning noon afternoon full_day"`
	WorkLocation        string  `json:"work_location"`
	WorkContent         string  `json:"work_content"`
	WorkType            string  `json:"work_type"`
	IsUrgent            bool    `json:"is_urgent"`
	Priority            int     `json:"priority" binding:"omitempty,min=0,max=10"`
	UrgentReason        string  `json:"urgent_reason"`
	AssignedWorkerID    *uint   `json:"assigned_worker_id"`     // 兼容单选
	AssignedWorkerIDs   string  `json:"assigned_worker_ids"`    // JSON数组格式，支持多选
	AssignedWorkerNames string  `json:"assigned_worker_names"`  // 逗号分隔
	Status              string  `json:"status" binding:"omitempty,oneof=draft pending scheduled in_progress completed cancelled rejected"`
}

// AppointmentListRequest 预约单列表查询请求
type AppointmentListRequest struct {
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
	Status     string `form:"status"`
	IsUrgent   *bool  `form:"is_urgent"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	ApplicantID *uint `form:"applicant_id"`
	WorkerID   *uint  `form:"worker_id"`
	WorkType     string `form:"work_type"`
	CurrentUserID uint  // 当前用户ID，用于权限过滤
}

// CalendarListRequest 日历列表查询请求
type CalendarListRequest struct {
	WorkerID   *uint  `form:"worker_id"`
	StartDate  string `form:"start_date" binding:"required"`
	EndDate    string `form:"end_date" binding:"required"`
	TimeSlot   string `form:"time_slot" binding:"omitempty,oneof=morning noon afternoon full_day"`
	Status     string `form:"status"`
}

// AvailabilityCheckRequest 可用性检查请求
type AvailabilityCheckRequest struct {
	WorkerID   uint   `json:"worker_id" binding:"required"`
	WorkDate   string `json:"work_date" binding:"required"`
	TimeSlot   string `json:"time_slot" binding:"required,oneof=morning noon afternoon full_day"`
}

// AssignWorkerRequest 分配作业人员请求
type AssignWorkerRequest struct {
	WorkerID     uint   `json:"worker_id,omitempty"`   // 单个作业人员ID（兼容旧版）
	WorkerIDs    []uint `json:"worker_ids,omitempty"`  // 多个作业人员ID（新版）
	SupervisorID *uint  `json:"supervisor_id,omitempty"` // 监护人ID
}

// BatchBlockCalendarRequest 批量锁定日历请求
type BatchBlockCalendarRequest struct {
	WorkerID      uint     `json:"worker_id" binding:"required"`
	StartDate     string   `json:"start_date" binding:"required"`
	EndDate       string   `json:"end_date" binding:"required"`
	TimeSlots     []string `json:"time_slots" binding:"required,min=1,dive,oneof=morning noon afternoon full_day"`
	BlockedReason string   `json:"blocked_reason"`
}

// BatchCreateAppointmentRequest 批量创建预约请求
type BatchCreateAppointmentRequest struct {
	Appointments []CreateAppointmentRequest `json:"appointments" binding:"required,min=1,max=10"`
}

// ApproveAppointmentRequest 审批预约请求
type ApproveAppointmentRequest struct {
	Action       string `json:"action" binding:"required,oneof=approve reject"`
	Comment      string `json:"comment"`
	AssignNow    bool   `json:"assign_now"` // 是否立即分配作业人员
	WorkerID     *uint  `json:"worker_id"`  // 指定的作业人员ID

	// 作业时间修改
	Reschedule   bool   `json:"reschedule"`    // 是否修改作业时间
	NewWorkDate  string `json:"new_work_date"` // 新作业日期 (YYYY-MM-DD)
	NewTimeSlot  string `json:"new_time_slot"` // 新时间段 (morning/noon/afternoon/full_day)
}

// CompleteAppointmentRequest 完成预约请求
type CompleteAppointmentRequest struct {
	CompletionNote string `json:"completion_note"`
	Photos         []string `json:"photos"` // 完成照片URL
}

// CalendarViewResponse 日历视图响应
type CalendarViewResponse struct {
	Date     string                    `json:"date"`
	TimeSlots map[string]TimeSlotInfo `json:"time_slots"`
}

// TimeSlotInfo 时间段信息
type TimeSlotInfo struct {
	IsAvailable bool                   `json:"is_available"`
	Status      string                 `json:"status"`
	Appointment *AppointmentSummary    `json:"appointment,omitempty"`
	Reason      string                 `json:"reason,omitempty"`
}

// AppointmentSummary 预约摘要
type AppointmentSummary struct {
	ID            uint   `json:"id"`
	AppointmentNo string `json:"appointment_no"`
	WorkLocation  string `json:"work_location"`
	WorkContent   string `json:"work_content"`
	IsUrgent      bool   `json:"is_urgent"`
	Status        string `json:"status"`
}

// StatsResponse 统计响应
type StatsResponse struct {
	Total       int64 `json:"total"`
	Draft       int64 `json:"draft"`
	Pending     int64 `json:"pending"`
	Scheduled   int64 `json:"scheduled"`
	InProgress  int64 `json:"in_progress"`
	Completed   int64 `json:"completed"`
	Cancelled   int64 `json:"cancelled"`
	Rejected    int64 `json:"rejected"`
	Urgent      int64 `json:"urgent"`
	TodayCount  int64 `json:"today_count"`
	WeekCount   int64 `json:"week_count"`
	MonthCount  int64 `json:"month_count"`
}

// DailyStatistics 每日统计数据
type DailyStatistics struct {
	Date         string `json:"date"`         // YYYY-MM-DD
	TotalCount   int64  `json:"total_count"`  // 当天任务总数
	UrgentCount  int64  `json:"urgent_count"` // 加急任务数
	TotalWorkers int64  `json:"total_workers"` // 总作业人员数
}

// DailyStatisticsResponse 每日统计响应
type DailyStatisticsResponse struct {
	Statistics []DailyStatistics `json:"statistics"`
	TotalWorkers int64           `json:"total_workers"` // 总作业人员数
}

// TimeSlotStatistics 时间段统计数据
type TimeSlotStatistics struct {
	Date        string `json:"date"`        // YYYY-MM-DD
	TimeSlot    string `json:"time_slot"`   // morning, noon, afternoon, full_day
	TotalCount  int64  `json:"total_count"` // 该时间段任务数
	TotalWorkers int64  `json:"total_workers"` // 总作业人员数
}

// TimeSlotStatisticsResponse 时间段统计响应
type TimeSlotStatisticsResponse struct {
	Statistics   []TimeSlotStatistics `json:"statistics"`
	TotalWorkers int64                `json:"total_workers"` // 总作业人员数
}

// TimeSlotConstants 时间段常量
const (
	TimeSlotMorning   = "morning"
	TimeSlotNoon      = "noon"
	TimeSlotAfternoon = "afternoon"
	TimeSlotFullDay   = "full_day"
)

// StatusConstants 状态常量
const (
	StatusDraft        = "draft"
	StatusPending      = "pending"
	StatusScheduled    = "scheduled"
	StatusInProgress   = "in_progress"
	StatusCompleted    = "completed"
	StatusCancelled    = "cancelled"
	StatusRejected     = "rejected"
)

// CalendarStatusConstants 日历状态常量
const (
	CalendarStatusAvailable = "available"
	CalendarStatusBusy      = "busy"
	CalendarStatusBlocked   = "blocked"
	CalendarStatusOff       = "off"
)

// BeforeCreate GORM hook - 创建前
func (a *ConstructionAppointment) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM hook - 更新前
func (a *ConstructionAppointment) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate GORM hook - 创建前
func (w *WorkerCalendar) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	w.CreatedAt = now
	w.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM hook - 更新前
func (w *WorkerCalendar) BeforeUpdate(tx *gorm.DB) error {
	w.UpdatedAt = time.Now()
	return nil
}

// IsEditable 检查预约单是否可编辑
func (a *ConstructionAppointment) IsEditable() bool {
	return a.Status == StatusDraft
}

// IsEditableBy 检查预约单是否可被指定用户编辑
// 只有申请人可以在 draft 或 pending 状态下编辑
func (a *ConstructionAppointment) IsEditableBy(userID uint) bool {
	// 状态必须是 draft 或 pending
	if a.Status != StatusDraft && a.Status != StatusPending {
		return false
	}
	// 必须是申请人本人
	return a.ApplicantID == userID
}

// IsCancellable 检查预约单是否可取消
func (a *ConstructionAppointment) IsCancellable() bool {
	return a.Status == StatusPending || a.Status == StatusScheduled
}

// CanComplete 检查预约单是否可完成
func (a *ConstructionAppointment) CanComplete() bool {
	return a.Status == StatusInProgress || a.Status == StatusScheduled
}

// NeedsUrgentApproval 检查是否需要加急审批
func (a *ConstructionAppointment) NeedsUrgentApproval() bool {
	return a.IsUrgent && a.Priority >= 7
}

// GetTimeSlotLabel 获取时间段标签
func GetTimeSlotLabel(timeSlot string) string {
	labels := map[string]string{
		TimeSlotMorning:   "上午 (8:00-11:30)",
		TimeSlotNoon:      "中午 (12:00-13:30)",
		TimeSlotAfternoon: "下午 (13:30-16:30)",
		TimeSlotFullDay:   "全天",
	}
	if label, ok := labels[timeSlot]; ok {
		return label
	}
	return timeSlot
}

// ContactInfo 联系人信息
type ContactInfo struct {
	ContactPerson string `json:"contact_person"`
	ContactPhone  string `json:"contact_phone"`
	Count         int    `json:"count"` // 使用次数
}

// GetStatusLabel 获取状态标签
func GetStatusLabel(status string) string {
	labels := map[string]string{
		StatusDraft:        "草稿",
		StatusPending:      "待审批",
		StatusScheduled:    "已排期",
		StatusInProgress:   "进行中",
		StatusCompleted:    "已完成",
		StatusCancelled:    "已取消",
		StatusRejected:     "已拒绝",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}
