package progress

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ChangeLog represents a change log entry
type ChangeLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProjectID  uint      `gorm:"not null;index:idx_project_changes" json:"project_id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	EntityType string    `gorm:"type:varchar(50);not null;index" json:"entity_type"` // task, dependency, resource
	EntityID   uint      `gorm:"not null;index" json:"entity_id"`
	ActionType string    `gorm:"type:varchar(20);not null" json:"action_type"` // create, update, delete
	Changes    ChangeData `gorm:"type:jsonb;not null" json:"changes"`
	CreatedAt  time.Time `json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for change logs
func (ChangeLog) TableName() string {
	return "change_logs"
}

// ChangeData represents before/after data
type ChangeData struct {
	Before map[string]interface{} `json:"before"`
	After  map[string]interface{} `json:"after"`
}

// Scan implements sql.Scanner for ChangeData
func (cd *ChangeData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, cd)
}

// Value implements driver.Valuer for ChangeData
func (cd ChangeData) Value() (driver.Value, error) {
	if cd.Before == nil {
		cd.Before = make(map[string]interface{})
	}
	if cd.After == nil {
		cd.After = make(map[string]interface{})
	}
	return json.Marshal(cd)
}

// ChangeLogRepository handles change log data operations
type ChangeLogRepository struct {
	db *gorm.DB
}

// NewChangeLogRepository creates a new change log repository
func NewChangeLogRepository(db *gorm.DB) *ChangeLogRepository {
	return &ChangeLogRepository{db: db}
}

// Create creates a new change log entry
func (r *ChangeLogRepository) Create(log *ChangeLog) error {
	return r.db.Create(log).Error
}

// GetByID retrieves a change log entry by ID
func (r *ChangeLogRepository) GetByID(id uint) (*ChangeLog, error) {
	var log ChangeLog
	err := r.db.Preload("User").First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetByProjectID retrieves all change log entries for a project
func (r *ChangeLogRepository) GetByProjectID(projectID uint, opts QueryOptions) ([]ChangeLog, int64, error) {
	var logs []ChangeLog
	var total int64

	query := r.db.Model(&ChangeLog{}).Where("project_id = ?", projectID)

	// Apply filters
	if opts.EntityType != "" {
		query = query.Where("entity_type = ?", opts.EntityType)
	}
	if opts.ActionType != "" {
		query = query.Where("action_type = ?", opts.ActionType)
	}
	if opts.UserID != 0 {
		query = query.Where("user_id = ?", opts.UserID)
	}
	if !opts.StartDate.IsZero() {
		query = query.Where("created_at >= ?", opts.StartDate)
	}
	if !opts.EndDate.IsZero() {
		query = query.Where("created_at <= ?", opts.EndDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		query = query.Offset(opts.Offset)
	}

	err := query.Preload("User").
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

// GetByEntityID retrieves all change log entries for an entity
func (r *ChangeLogRepository) GetByEntityID(entityType string, entityID uint, limit int) ([]ChangeLog, error) {
	var logs []ChangeLog
	query := r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID)

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Preload("User").
		Order("created_at DESC").
		Find(&logs).Error

	return logs, err
}

// GetByUserID retrieves all change log entries by a user
func (r *ChangeLogRepository) GetByUserID(userID uint, limit int) ([]ChangeLog, error) {
	var logs []ChangeLog
	query := r.db.Where("user_id = ?", userID)

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Preload("User").
		Order("created_at DESC").
		Find(&logs).Error

	return logs, err
}

// Delete deletes a change log entry
func (r *ChangeLogRepository) Delete(id uint) error {
	return r.db.Delete(&ChangeLog{}, id).Error
}

// DeleteOldEntries deletes entries older than specified days
func (r *ChangeLogRepository) DeleteOldEntries(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", cutoff).Delete(&ChangeLog{}).Error
}

// GetStatistics returns statistics about changes for a project
func (r *ChangeLogRepository) GetStatistics(projectID uint) (*ChangeStatistics, error) {
	var stats ChangeStatistics

	stats.ProjectID = projectID

	// Total changes
	if err := r.db.Model(&ChangeLog{}).
		Where("project_id = ?", projectID).
		Count(&stats.TotalChanges).Error; err != nil {
		return nil, err
	}

	// By action type
	rows, err := r.db.Model(&ChangeLog{}).
		Select("action_type, count(*) as count").
		Where("project_id = ?", projectID).
		Group("action_type").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats.ByActionType = make(map[string]int64)
	for rows.Next() {
		var actionType string
		var count int64
		if err := rows.Scan(&actionType, &count); err != nil {
			continue
		}
		stats.ByActionType[actionType] = count
	}

	// By entity type
	rows, err = r.db.Model(&ChangeLog{}).
		Select("entity_type, count(*) as count").
		Where("project_id = ?", projectID).
		Group("entity_type").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats.ByEntityType = make(map[string]int64)
	for rows.Next() {
		var entityType string
		var count int64
		if err := rows.Scan(&entityType, &count); err != nil {
			continue
		}
		stats.ByEntityType[entityType] = count
	}

	// Date range
	var earliest, latest time.Time
	if err := r.db.Model(&ChangeLog{}).
		Where("project_id = ?", projectID).
		Select("min(created_at)").Scan(&earliest).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&ChangeLog{}).
		Where("project_id = ?", projectID).
		Select("max(created_at)").Scan(&latest).Error; err != nil {
		return nil, err
	}

	stats.DateRange = DateRange{
		Earliest: earliest,
		Latest:   latest,
	}

	return &stats, nil
}

// ChangeStatistics represents change statistics
type ChangeStatistics struct {
	ProjectID     uint              `json:"project_id"`
	TotalChanges  int64             `json:"total_changes"`
	ByActionType  map[string]int64  `json:"by_action_type"`
	ByEntityType  map[string]int64  `json:"by_entity_type"`
	DateRange     DateRange         `json:"date_range"`
}

// DateRange represents a date range
type DateRange struct {
	Earliest time.Time `json:"earliest"`
	Latest   time.Time `json:"latest"`
}

// QueryOptions represents query options
type QueryOptions struct {
	EntityType string
	ActionType string
	UserID     uint
	StartDate  time.Time
	EndDate    time.Time
	Limit      int
	Offset     int
}

// RollbackData represents data needed for rollback
type RollbackData struct {
	EntityID   uint                   `json:"entity_id"`
	EntityType string                 `json:"entity_type"`
	Data       map[string]interface{} `json:"data"`
}

// GetRollbackData retrieves data for rolling back a change
func (r *ChangeLogRepository) GetRollbackData(changeID uint) (*RollbackData, error) {
	log, err := r.GetByID(changeID)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if log.ActionType == "delete" {
		// For delete, restore from before state
		data = log.Changes.Before
	} else {
		// For create/update, use after state
		data = log.Changes.After
	}

	return &RollbackData{
		EntityID:   log.EntityID,
		EntityType: log.EntityType,
		Data:       data,
	}, nil
}
