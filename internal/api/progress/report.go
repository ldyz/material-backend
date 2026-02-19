package progress

import (
	"time"
)

// ReportConfig defines the configuration for generating reports
type ReportConfig struct {
	ProjectID   uint      `json:"project_id" gorm:"index"`
	Type        string    `json:"type" gorm:"type:varchar(50);index"` // task, resource, milestone, progress
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Columns     []string  `json:"columns" gorm:"type:text"`     // JSON array
	GroupBy     string    `json:"group_by" gorm:"type:varchar(50)"` // status, assignee, priority, phase, milestone
	Filters     []Filter  `json:"filters" gorm:"type:text"`     // JSON array
	SortBy      string    `json:"sort_by" gorm:"type:varchar(50)"`
	SortOrder   string    `json:"sort_order" gorm:"type:varchar(10)"` // asc, desc
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Filter defines a filter condition
type Filter struct {
	Field    string `json:"field"`    // name, status, assignee, priority, etc.
	Operator string `json:"operator"` // eq, ne, contains, gt, lt
	Value    string `json:"value"`
}

// Report represents a saved report
type Report struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	Name      string      `json:"name" gorm:"type:varchar(200);not null" binding:"required"`
	Config    ReportConfig `json:"config" gorm:"embedded;embeddedPrefix:config_"`
	CreatedBy uint        `json:"created_by" gorm:"index"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`

	// Associations
	User     *User   `json:"user,omitempty" gorm:"foreignKey:CreatedBy"`
	Project  *Project `json:"project,omitempty" gorm:"foreignKey:Config.ProjectID"`
}

// ReportData holds the generated report data
type ReportData struct {
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	DateRange []time.Time            `json:"date_range"`
	Columns   []ColumnDef            `json:"columns"`
	Data      []map[string]interface{} `json:"data"`
	Groups    []GroupData            `json:"groups,omitempty"`
	Summary   SummaryStats           `json:"summary"`
	Metadata  ReportMetadata         `json:"metadata"`
}

// ColumnDef defines a column in the report
type ColumnDef struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Width int    `json:"width"`
	Type  string `json:"type"` // string, number, date, currency
}

// GroupData holds grouped data
type GroupData struct {
	Key   string                   `json:"key"`
	Count int                      `json:"count"`
	Items []map[string]interface{} `json:"items"`
}

// SummaryStats holds summary statistics
type SummaryStats struct {
	TotalRecords int                    `json:"total_records"`
	Stats        map[string]interface{} `json:"stats"`
}

// ReportMetadata holds metadata about the report
type ReportMetadata struct {
	GeneratedAt time.Time `json:"generated_at"`
	GeneratedBy uint      `json:"generated_by"`
	ProjectID   uint      `json:"project_id"`
	ProjectName string    `json:"project_name"`
}

// ReportTemplate is a reusable report configuration
type ReportTemplate struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	Name        string      `json:"name" gorm:"type:varchar(200);not null" binding:"required"`
	Description string      `json:"description" gorm:"type:text"`
	Config      ReportConfig `json:"config" gorm:"embedded;embeddedPrefix:config_"`
	IsSystem    bool        `json:"is_system" gorm:"default:false"` // System templates cannot be deleted
	CreatedBy   uint        `json:"created_by" gorm:"index"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// ReportRequest is the request payload for generating reports
type ReportRequest struct {
	ProjectID uint      `json:"project_id" binding:"required"`
	Type      string    `json:"type" binding:"required,oneof=task resource milestone progress"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Columns   []string  `json:"columns"`
	GroupBy   string    `json:"group_by"`
	Filters   []Filter  `json:"filters"`
	SortBy    string    `json:"sort_by"`
	SortOrder string    `json:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// ReportSaveRequest is the request payload for saving reports
type ReportSaveRequest struct {
	Name   string      `json:"name" binding:"required"`
	Config ReportConfig `json:"config" binding:"required"`
}

// ReportExportRequest is the request payload for exporting reports
type ReportExportRequest struct {
	Format string `json:"format" binding:"required,oneof=pdf excel print"`
}
