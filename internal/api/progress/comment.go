package progress

import (
	"time"

	"gorm.io/gorm"
)

// Comment represents a task comment
type Comment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `gorm:"not null;index:idx_task_comments" json:"task_id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	ParentID   *uint     `gorm:"index" json:"parent_id,omitempty"` // For threading
	IsResolved bool      `gorm:"default:false" json:"is_resolved"`
	Mentions   []string  `gorm:"type:text[]" json:"mentions,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	User     *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Task     *Task    `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Parent   *Comment `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies  []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// TableName specifies the table name for comments
func (Comment) TableName() string {
	return "task_comments"
}

// CommentRepository handles comment data operations
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create creates a new comment
func (r *CommentRepository) Create(comment *Comment) error {
	return r.db.Create(comment).Error
}

// GetByID retrieves a comment by ID
func (r *CommentRepository) GetByID(id uint) (*Comment, error) {
	var comment Comment
	err := r.db.Preload("User").
		Preload("Parent").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByTaskID retrieves all comments for a task
func (r *CommentRepository) GetByTaskID(taskID uint) ([]Comment, error) {
	var comments []Comment
	err := r.db.Where("task_id = ? AND parent_id IS NULL", taskID).
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at ASC")
		}).
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}

// GetByProjectID retrieves all comments for a project
func (r *CommentRepository) GetByProjectID(projectID uint) ([]Comment, error) {
	var comments []Comment
	err := r.db.Joins("JOIN tasks ON tasks.id = task_comments.task_id").
		Where("tasks.project_id = ?", projectID).
		Preload("User").
		Preload("Task").
		Order("task_comments.created_at DESC").
		Find(&comments).Error
	return comments, err
}

// Update updates a comment
func (r *CommentRepository) Update(comment *Comment) error {
	return r.db.Save(comment).Error
}

// Delete deletes a comment
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&Comment{}, id).Error
}

// Resolve marks a comment as resolved
func (r *CommentRepository) Resolve(id uint) error {
	return r.db.Model(&Comment{}).Where("id = ?", id).Update("is_resolved", true).Error
}

// Unresolve marks a comment as unresolved
func (r *CommentRepository) Unresolve(id uint) error {
	return r.db.Model(&Comment{}).Where("id = ?", id).Update("is_resolved", false).Error
}

// GetUnresolvedCount returns the count of unresolved comments for a task
func (r *CommentRepository) GetUnresolvedCount(taskID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Comment{}).
		Where("task_id = ? AND is_resolved = false AND parent_id IS NULL", taskID).
		Count(&count).Error
	return count, err
}

// GetByUserID retrieves all comments by a user
func (r *CommentRepository) GetByUserID(userID uint, limit int) ([]Comment, error) {
	var comments []Comment
	err := r.db.Where("user_id = ?", userID).
		Preload("Task").
		Order("created_at DESC").
		Limit(limit).
		Find(&comments).Error
	return comments, err
}

// GetMentionsForUser retrieves all comments mentioning a user
func (r *CommentRepository) GetMentionsForUser(userID uint, limit int) ([]Comment, error) {
	var comments []Comment
	err := r.db.Where("? = ANY(mentions)", userID).
		Preload("Task").
		Order("created_at DESC").
		Limit(limit).
		Find(&comments).Error
	return comments, err
}

// User represents a user (simplified for foreign key)
type User struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// TableName specifies the table name for users
func (User) TableName() string {
	return "users"
}
