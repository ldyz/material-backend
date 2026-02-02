package project

import (
	"time"

	"github.com/yourorg/material-backend/backend/internal/api/auth"
)

// ProjectStatus constants
const (
	StatusPlanning = "planning"
	StatusActive   = "active"
	StatusClosed   = "closed"
	StatusOnHold   = "on_hold"
)

// Project model maps to 'projects' table
type Project struct {
	ID                 uint        `gorm:"primaryKey" json:"id"`
	Name               string      `gorm:"type:text;not null" json:"name"`
	Code               string      `gorm:"type:text" json:"code"`
	Location           string      `gorm:"type:text" json:"location"`
	StartDate          *time.Time  `gorm:"type:date" json:"start_date"`
	EndDate            *time.Time  `gorm:"type:date" json:"end_date"`
	Description        string      `gorm:"type:text" json:"description"`
	Manager            string      `gorm:"type:text" json:"manager"`
	Contact            string      `gorm:"type:text" json:"contact"`
	Budget             string      `gorm:"type:text" json:"budget"` // PostgreSQL中是text类型，不是decimal
	Status             string      `gorm:"type:text" json:"status"`
	Users              []auth.User `gorm:"many2many:user_projects" json:"users,omitempty"`
	// Hierarchy fields
	ParentID           *uint       `gorm:"index" json:"parent_id,omitempty"`
	Level              int         `gorm:"default:0" json:"level"`
	Path               string      `gorm:"type:varchar(500)" json:"path,omitempty"`
	ProgressPercentage float64     `gorm:"default:0" json:"progress_percentage"`
	Children           []Project   `gorm:"-" json:"children,omitempty"`
}

func (p *Project) ToDTO() map[string]any {
	users := make([]map[string]any, 0, len(p.Users))
	for _, u := range p.Users {
		users = append(users, u.ToDTO())
	}
	return map[string]any{
		"id":                    p.ID,
		"name":                  p.Name,
		"code":                  p.Code,
		"location":              p.Location,
		"start_date":            p.StartDate,
		"end_date":              p.EndDate,
		"description":           p.Description,
		"manager":               p.Manager,
		"contact":               p.Contact,
		"budget":                p.Budget,
		"status":                p.Status,
		"users":                 users,
		"parent_id":             p.ParentID,
		"level":                 p.Level,
		"path":                  p.Path,
		"progress_percentage":   p.ProgressPercentage,
		"children":              p.Children,
	}
}
