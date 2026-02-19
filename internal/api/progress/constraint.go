package progress

import (
	"time"

	"gorm.io/gorm"
)

// Constraint represents a task scheduling constraint
type Constraint struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskID    uint      `json:"task_id" gorm:"not null;index"`
	Type      string    `json:"type" gorm:"size:10;not null"` // MSO, MFO, SNET, SNLT, FNET, FNLT
	Date      time.Time `json:"date" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Constraint types
const (
	ConstraintMustStartOn           = "MSO" // Must Start On
	ConstraintMustFinishOn          = "MFO" // Must Finish On
	ConstraintStartNoEarlierThan    = "SNET" // Start No Earlier Than
	ConstraintStartNoLaterThan      = "SNLT" // Start No Later Than
	ConstraintFinishNoEarlierThan   = "FNET" // Finish No Earlier Than
	ConstraintFinishNoLaterThan     = "FNLT" // Finish No Later Than
)

// Validate checks if the constraint is valid
func (c *Constraint) Validate() error {
	// Validate constraint type
	validTypes := map[string]bool{
		ConstraintMustStartOn:         true,
		ConstraintMustFinishOn:        true,
		ConstraintStartNoEarlierThan:  true,
		ConstraintStartNoLaterThan:    true,
		ConstraintFinishNoEarlierThan: true,
		ConstraintFinishNoLaterThan:   true,
	}

	if !validTypes[c.Type] {
		return &ValidationError{
			Field:   "type",
			Message: "invalid constraint type",
		}
	}

	// Validate date
	if c.Date.IsZero() {
		return &ValidationError{
			Field:   "date",
			Message: "constraint date is required",
		}
	}

	return nil
}

// ApplyToTask applies the constraint to a task and adjusts dates
func (c *Constraint) ApplyToTask(task *Task) error {
	if task == nil {
		return &ValidationError{
			Field:   "task",
			Message: "task cannot be nil",
		}
	}

	// Calculate task duration
	duration := task.EndDate.Sub(task.StartDate)

	switch c.Type {
	case ConstraintMustStartOn:
		// Task must start exactly on constraint date
		task.StartDate = c.Date
		task.EndDate = c.Date.Add(duration)

	case ConstraintMustFinishOn:
		// Task must finish exactly on constraint date
		task.EndDate = c.Date
		task.StartDate = c.Date.Add(-duration)

	case ConstraintStartNoEarlierThan:
		// Task cannot start before constraint date
		if task.StartDate.Before(c.Date) {
			task.StartDate = c.Date
			task.EndDate = c.Date.Add(duration)
		}

	case ConstraintStartNoLaterThan:
		// Task cannot start after constraint date
		if task.StartDate.After(c.Date) {
			task.StartDate = c.Date
			task.EndDate = c.Date.Add(duration)
		}

	case ConstraintFinishNoEarlierThan:
		// Task cannot finish before constraint date
		if task.EndDate.Before(c.Date) {
			task.EndDate = c.Date
			task.StartDate = c.Date.Add(-duration)
		}

	case ConstraintFinishNoLaterThan:
		// Task cannot finish after constraint date
		if task.EndDate.After(c.Date) {
			task.EndDate = c.Date
			task.StartDate = c.Date.Add(-duration)
		}
	}

	return nil
}

// IsSatisfied checks if the task satisfies the constraint
func (c *Constraint) IsSatisfied(task *Task) bool {
	if task == nil {
		return false
	}

	switch c.Type {
	case ConstraintMustStartOn:
		return task.StartDate.Equal(c.Date)

	case ConstraintMustFinishOn:
		return task.EndDate.Equal(c.Date)

	case ConstraintStartNoEarlierThan:
		return !task.StartDate.Before(c.Date)

	case ConstraintStartNoLaterThan:
		return !task.StartDate.After(c.Date)

	case ConstraintFinishNoEarlierThan:
		return !task.EndDate.Before(c.Date)

	case ConstraintFinishNoLaterThan:
		return !task.EndDate.After(c.Date)
	}

	return true
}

// ConstraintRepository handles constraint data operations
type ConstraintRepository struct {
	db *gorm.DB
}

// NewConstraintRepository creates a new constraint repository
func NewConstraintRepository(db *gorm.DB) *ConstraintRepository {
	return &ConstraintRepository{db: db}
}

// Create creates a new constraint
func (r *ConstraintRepository) Create(constraint *Constraint) error {
	return r.db.Create(constraint).Error
}

// GetByID retrieves a constraint by ID
func (r *ConstraintRepository) GetByID(id uint) (*Constraint, error) {
	var constraint Constraint
	err := r.db.First(&constraint, id).Error
	if err != nil {
		return nil, err
	}
	return &constraint, nil
}

// GetByTaskID retrieves all constraints for a task
func (r *ConstraintRepository) GetByTaskID(taskID uint) ([]Constraint, error) {
	var constraints []Constraint
	err := r.db.Where("task_id = ?", taskID).Find(&constraints).Error
	return constraints, err
}

// Update updates a constraint
func (r *ConstraintRepository) Update(constraint *Constraint) error {
	return r.db.Save(constraint).Error
}

// Delete deletes a constraint
func (r *ConstraintRepository) Delete(id uint) error {
	return r.db.Delete(&Constraint{}, id).Error
}

// DeleteByTaskID deletes all constraints for a task
func (r *ConstraintRepository) DeleteByTaskID(taskID uint) error {
	return r.db.Where("task_id = ?", taskID).Delete(&Constraint{}).Error
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
