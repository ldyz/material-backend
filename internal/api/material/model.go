package material

import "time"

// Material model maps to 'materials' table
type Material struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Code          *string   `gorm:"type:text" json:"code"`
	Name          string    `gorm:"type:text;not null" json:"name"`
	Specification string    `gorm:"type:text" json:"specification"`
	Unit          string    `gorm:"type:text" json:"unit"`
	Price         float64   `gorm:"type:real" json:"price"`
	Description   string    `gorm:"type:text" json:"description"`
	Category      string    `gorm:"type:text" json:"category"`
	Quantity      int       `gorm:"type:integer" json:"quantity"`
	ProjectID     *uint     `gorm:"type:integer" json:"project_id"`
	Material      string    `gorm:"type:text" json:"material"`
	Spec          string    `gorm:"type:text" json:"spec"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName specifies the table name for Material
func (Material) TableName() string {
	return "materials"
}

// ToDTO converts Material to DTO format
func (m *Material) ToDTO() map[string]any {
	spec := m.Specification
	if spec == "" && m.Spec != "" {
		spec = m.Spec
	}

	var projectID uint
	if m.ProjectID != nil {
		projectID = *m.ProjectID
	}

	return map[string]any{
		"id":            m.ID,
		"code":          m.Code,
		"name":          m.Name,
		"specification": spec,
		"unit":          m.Unit,
		"price":         m.Price,
		"description":   m.Description,
		"category":      m.Category,
		"quantity":      m.Quantity,
		"project_id":    projectID,
		"material":      m.Material,
		"spec":          m.Spec,
		"created_at":    m.CreatedAt,
		"updated_at":    m.UpdatedAt,
	}
}

// MaterialCategory represents a material category
type MaterialCategory struct {
	ID        uint                `gorm:"primaryKey" json:"id"`
	ParentID  uint                `gorm:"default:0;index" json:"parent_id"`
	Name      string              `gorm:"size:50;not null" json:"name"`
	Code      string              `gorm:"size:20;index" json:"code"`
	Level     int                 `gorm:"default:1" json:"level"`
	Path      string              `gorm:"size:255" json:"path"`
	Sort      int                 `gorm:"default:0" json:"sort"`
	Remark    string              `gorm:"type:text" json:"remark"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Children  []MaterialCategory  `gorm:"-" json:"children,omitempty"`
}

// TableName specifies the table name for MaterialCategory
func (MaterialCategory) TableName() string {
	return "material_categories"
}

// ToDTO converts MaterialCategory to DTO format
func (mc *MaterialCategory) ToDTO() map[string]any {
	return map[string]any{
		"id":         mc.ID,
		"parent_id":  mc.ParentID,
		"name":       mc.Name,
		"code":       mc.Code,
		"level":      mc.Level,
		"path":       mc.Path,
		"sort":       mc.Sort,
		"remark":     mc.Remark,
		"created_at": mc.CreatedAt,
		"updated_at": mc.UpdatedAt,
	}
}
