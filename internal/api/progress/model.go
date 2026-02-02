package progress

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Task 任务表
type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ProjectID   uint       `gorm:"not null;index" json:"project_id"`
	ScheduleID  *uint      `gorm:"index" json:"schedule_id,omitempty"`
	ParentID    *uint      `gorm:"index" json:"parent_id,omitempty"`
	Name        string     `gorm:"type:varchar(200);not null" json:"name"`
	Duration    *float64   `json:"duration,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Progress    float64    `gorm:"default:0" json:"progress"`
	IsMilestone bool       `gorm:"default:false" json:"is_milestone"`
	Priority    string     `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	Status      string     `gorm:"type:varchar(20);default:'not_started'" json:"status"`
	Responsible string     `gorm:"type:varchar(100)" json:"responsible,omitempty"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	PositionX   *float64   `json:"position_x,omitempty"`
	PositionY   *float64   `json:"position_y,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Task) TableName() string {
	return "tasks"
}

// TaskDependency 任务依赖表
type TaskDependency struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    uint      `gorm:"not null;index" json:"task_id"`
	DependsOn uint      `gorm:"not null;index" json:"depends_on"`
	Type      string    `gorm:"default:FS" json:"type"` // FS, FF, SS, SF
	Lag       int       `gorm:"default:0" json:"lag"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (TaskDependency) TableName() string {
	return "task_dependencies"
}

// ProjectSchedule 项目进度计划表
type ProjectSchedule struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID uint      `gorm:"not null;index;comment:项目ID" json:"project_id"`
	Data      ScheduleData `gorm:"type:jsonb;not null;comment:进度数据(JSON格式)" json:"data"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	CreatedBy *uint     `gorm:"comment:创建者ID" json:"created_by,omitempty"`
	UpdatedBy *uint     `gorm:"comment:更新者ID" json:"updated_by,omitempty"`
}

// TableName 指定表名
func (ProjectSchedule) TableName() string {
	return "project_schedules"
}

// ScheduleData 进度数据结构（统一格式）
type ScheduleData struct {
	Nodes      map[string]Node      `json:"nodes"`
	Activities map[string]Activity  `json:"activities"`
}

// Node 节点数据
type Node struct {
	ID           string  `json:"id"`
	Label        string  `json:"label"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Type         string  `json:"type"` // start, end, event
	Number       *int    `json:"number,omitempty"`
	EarliestTime *float64 `json:"earliest_time,omitempty"`
	LatestTime   *float64 `json:"latest_time,omitempty"`
	Activities   []string `json:"activities,omitempty"`
}

// Activity 活动数据
type Activity struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Duration      float64       `json:"duration"`
	FromNode      string        `json:"from_node"`
	ToNode        string        `json:"to_node"`
	EarliestStart float64       `json:"earliest_start"`
	EarliestFinish float64       `json:"earliest_finish"`
	LatestStart   float64       `json:"latest_start"`
	LatestFinish  float64       `json:"latest_finish"`
	TotalFloat    float64       `json:"total_float"`
	FreeFloat     float64       `json:"free_float"`
	IsCritical    bool          `json:"is_critical"`
	IsDummy       bool          `json:"is_dummy"`
	Breakpoints   []Breakpoint  `json:"breakpoints,omitempty"`
	Progress      float64       `json:"progress"`
	Status        string        `json:"status,omitempty"`
	Priority      string        `json:"priority,omitempty"`
	Predecessors  []int         `json:"predecessors,omitempty"`  // 前置任务ID列表
	Successors    []int         `json:"successors,omitempty"`    // 后置任务ID列表
	Resources     []Resource    `json:"resources,omitempty"`
	ParentID      *uint         `json:"parent_id,omitempty"`
	SortOrder     int           `json:"sort_order,omitempty"`
}

// Breakpoint 路径折点
type Breakpoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Scan 实现 sql.Scanner 接口
func (sd *ScheduleData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, sd)
}

// Value 实现 driver.Valuer 接口
func (sd ScheduleData) Value() (driver.Value, error) {
	if sd.Nodes == nil {
		sd.Nodes = make(map[string]Node)
	}
	if sd.Activities == nil {
		sd.Activities = make(map[string]Activity)
	}
	return json.Marshal(sd)
}

// ToDTO 转换为 DTO
func (ps *ProjectSchedule) ToDTO() map[string]interface{} {
	return map[string]interface{}{
		"id":         ps.ID,
		"project_id": ps.ProjectID,
		"data":       ps.Data,
		"created_at": ps.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": ps.UpdatedAt.Format("2006-01-02 15:04:05"),
		"created_by": ps.CreatedBy,
		"updated_by": ps.UpdatedBy,
	}
}
