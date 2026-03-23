package progress

import (
	"time"

	"gorm.io/gorm"
)

// Calendar represents a working time calendar
type Calendar struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	Name        string          `json:"name" gorm:"size:255;not null"`
	Description string          `json:"description" gorm:"size:500"`
	ProjectID   *uint           `json:"project_id" gorm:"index"`
	IsDefault   bool            `json:"is_default" gorm:"default:false"`
	WorkingDays string          `json:"working_days" gorm:"size:20"` // JSON array of day numbers (0-6)
	WorkingHours string         `json:"working_hours" gorm:"size:100"` // JSON object with start/end times
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`

	// Associations
	Holidays    []Holiday       `json:"holidays" gorm:"foreignKey:CalendarID"`
	Exceptions  []CalendarException `json:"exceptions" gorm:"foreignKey:CalendarID"`
}

// Holiday represents a non-working day
type Holiday struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	CalendarID uint       `json:"calendar_id" gorm:"not null;index"`
	Date       time.Time  `json:"date" gorm:"not null"`
	Name       string     `json:"name" gorm:"size:255;not null"`
	Recurring  bool       `json:"recurring" gorm:"default:false"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// CalendarException represents an exception to the standard calendar
type CalendarException struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CalendarID   uint      `json:"calendar_id" gorm:"not null;index"`
	Date         time.Time `json:"date" gorm:"not null"`
	IsWorkingDay bool      `json:"is_working_day" gorm:"default:true"`
	Description  string    `json:"description" gorm:"size:500"`
	WorkingHours string    `json:"working_hours" gorm:"size:100"` // Override working hours for this day
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TaskCalendar represents a task-specific calendar assignment
type TaskCalendar struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	TaskID     uint      `json:"task_id" gorm:"not null;uniqueIndex"`
	CalendarID uint      `json:"calendar_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CalendarRepository handles calendar data operations
type CalendarRepository struct {
	db *gorm.DB
}

// NewCalendarRepository creates a new calendar repository
func NewCalendarRepository(db *gorm.DB) *CalendarRepository {
	return &CalendarRepository{db: db}
}

// Create creates a new calendar
func (r *CalendarRepository) Create(calendar *Calendar) error {
	return r.db.Create(calendar).Error
}

// GetByID retrieves a calendar by ID
func (r *CalendarRepository) GetByID(id uint) (*Calendar, error) {
	var calendar Calendar
	err := r.db.Preload("Holidays").Preload("Exceptions").First(&calendar, id).Error
	if err != nil {
		return nil, err
	}
	return &calendar, nil
}

// GetByProjectID retrieves all calendars for a project
func (r *CalendarRepository) GetByProjectID(projectID uint) ([]Calendar, error) {
	var calendars []Calendar
	err := r.db.Where("project_id = ?", projectID).
		Preload("Holidays").
		Preload("Exceptions").
		Find(&calendars).Error
	return calendars, err
}

// GetDefaultCalendar retrieves the default calendar for a project
func (r *CalendarRepository) GetDefaultCalendar(projectID uint) (*Calendar, error) {
	var calendar Calendar
	err := r.db.Where("project_id = ? AND is_default = ?", projectID, true).
		Preload("Holidays").
		Preload("Exceptions").
		First(&calendar).Error
	if err != nil {
		return nil, err
	}
	return &calendar, nil
}

// Update updates a calendar
func (r *CalendarRepository) Update(calendar *Calendar) error {
	return r.db.Save(calendar).Error
}

// Delete deletes a calendar
func (r *CalendarRepository) Delete(id uint) error {
	return r.db.Delete(&Calendar{}, id).Error
}

// AddHoliday adds a holiday to a calendar
func (r *CalendarRepository) AddHoliday(holiday *Holiday) error {
	return r.db.Create(holiday).Error
}

// GetHolidays retrieves all holidays for a calendar
func (r *CalendarRepository) GetHolidays(calendarID uint) ([]Holiday, error) {
	var holidays []Holiday
	err := r.db.Where("calendar_id = ?", calendarID).Find(&holidays).Error
	return holidays, err
}

// DeleteHoliday deletes a holiday
func (r *CalendarRepository) DeleteHoliday(id uint) error {
	return r.db.Delete(&Holiday{}, id).Error
}

// AddException adds an exception to a calendar
func (r *CalendarRepository) AddException(exception *CalendarException) error {
	return r.db.Create(exception).Error
}

// GetExceptions retrieves all exceptions for a calendar
func (r *CalendarRepository) GetExceptions(calendarID uint) ([]CalendarException, error) {
	var exceptions []CalendarException
	err := r.db.Where("calendar_id = ?", calendarID).Find(&exceptions).Error
	return exceptions, err
}

// DeleteException deletes an exception
func (r *CalendarRepository) DeleteException(id uint) error {
	return r.db.Delete(&CalendarException{}, id).Error
}

// AssignTaskCalendar assigns a calendar to a task
func (r *CalendarRepository) AssignTaskCalendar(taskCalendar *TaskCalendar) error {
	return r.db.Create(taskCalendar).Error
}

// GetTaskCalendar retrieves the calendar for a task
func (r *CalendarRepository) GetTaskCalendar(taskID uint) (*TaskCalendar, error) {
	var taskCalendar TaskCalendar
	err := r.db.Where("task_id = ?", taskID).First(&taskCalendar).Error
	if err != nil {
		return nil, err
	}
	return &taskCalendar, nil
}

// RemoveTaskCalendar removes a task's calendar assignment
func (r *CalendarRepository) RemoveTaskCalendar(taskID uint) error {
	return r.db.Where("task_id = ?", taskID).Delete(&TaskCalendar{}).Error
}

// WorkingTimeCalculator provides working time calculations
type WorkingTimeCalculator struct {
	calendarRepo *CalendarRepository
}

// NewWorkingTimeCalculator creates a new working time calculator
func NewWorkingTimeCalculator(calendarRepo *CalendarRepository) *WorkingTimeCalculator {
	return &WorkingTimeCalculator{
		calendarRepo: calendarRepo,
	}
}

// IsWorkingDay checks if a date is a working day
func (c *WorkingTimeCalculator) IsWorkingDay(date time.Time, calendar *Calendar) bool {
	if calendar == nil {
		// Default: Monday-Friday are working days
		weekday := date.Weekday()
		return weekday >= time.Monday && weekday <= time.Friday
	}

	// Check exceptions first
	for _, exception := range calendar.Exceptions {
		if exception.Date.Year() == date.Year() &&
			exception.Date.Month() == date.Month() &&
			exception.Date.Day() == date.Day() {
			return exception.IsWorkingDay
		}
	}

	// Check holidays
	for _, holiday := range calendar.Holidays {
		isSameDay := holiday.Date.Year() == date.Year() &&
			holiday.Date.Month() == date.Month() &&
			holiday.Date.Day() == date.Day()

		if isSameDay {
			return false
		}

		// Check recurring holidays (same month and day)
		if holiday.Recurring &&
			holiday.Date.Month() == date.Month() &&
			holiday.Date.Day() == date.Day() {
			return false
		}
	}

	// Check working days
	weekday := date.Weekday()
	// TODO: Parse workingDays JSON and check if weekday is included
	// For now, assume Monday-Friday
	return weekday >= time.Monday && weekday <= time.Friday
}

// CalculateWorkingDays calculates the number of working days between two dates
func (c *WorkingTimeCalculator) CalculateWorkingDays(
	start, end time.Time,
	calendar *Calendar,
) int {
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if c.IsWorkingDay(d, calendar) {
			count++
		}
	}
	return count
}

// AddWorkingDays adds working days to a date
func (c *WorkingTimeCalculator) AddWorkingDays(
	date time.Time,
	days int,
	calendar *Calendar,
) time.Time {
	result := date
	added := 0

	for added < days {
		result = result.AddDate(0, 0, 1)
		if c.IsWorkingDay(result, calendar) {
			added++
		}
	}

	return result
}
