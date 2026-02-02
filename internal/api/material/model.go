package material

// Material model maps to 'materials' table
type Material struct {
	ID             uint     `gorm:"primaryKey" json:"id"`
	Code           *string  `gorm:"type:text" json:"code"`
	Name           string   `gorm:"type:text;not null" json:"name"`
	Specification  string   `gorm:"type:text" json:"specification"`
	Unit           string   `gorm:"type:text" json:"unit"`
	Price          float64  `gorm:"type:real" json:"price"`
	Description    string   `gorm:"type:text" json:"description"`
	Category       string   `gorm:"type:text" json:"category"`
	Quantity       int      `gorm:"type:integer" json:"quantity"`
	ProjectID      *uint    `gorm:"type:integer" json:"project_id"`
	Material       string   `gorm:"type:text" json:"material"`
	Spec           string   `gorm:"type:text" json:"spec"`
}

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
		"id": m.ID,
		"code": m.Code,
		"name": m.Name,
		"specification": spec,
		"unit": m.Unit,
		"price": m.Price,
		"description": m.Description,
		"category": m.Category,
		"quantity": m.Quantity,
		"project_id": projectID,
		"material": m.Material,
		"spec": m.Spec,
	}
}
