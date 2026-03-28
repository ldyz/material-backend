package notification

import (
	"time"
)

// Notification model maps to 'notifications' table
type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Type      string    `gorm:"size:50;index" json:"type"`
	Title     string    `gorm:"size:200" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Data      string    `gorm:"type:jsonb" json:"data"` // JSON string
	IsRead    bool      `gorm:"index;default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt    *time.Time `json:"read_at"`
}

// Notification types
const (
	TypeRequisitionApprove  = "requisition_approve"
	TypeInboundApprove      = "inbound_approve"
	TypeMaterialPlanApprove = "material_plan_approve"
	TypeAppointmentApprove  = "appointment_approve"
	TypeStockAlert          = "stock_alert"
	TypeSystem              = "system"
)

func (n *Notification) ToDTO() map[string]any {
	var readAtStr string
	if n.ReadAt != nil {
		readAtStr = n.ReadAt.Format(time.RFC3339)
	}

	return map[string]any{
		"id":         n.ID,
		"user_id":    n.UserID,
		"type":       n.Type,
		"title":      n.Title,
		"content":    n.Content,
		"data":       n.Data,
		"is_read":    n.IsRead,
		"created_at": n.CreatedAt.Format(time.RFC3339),
		"read_at":    readAtStr,
	}
}

// DeviceToken model for storing push notification tokens
type DeviceToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Token     string    `gorm:"size:500;index" json:"token"`
	Platform  string    `gorm:"size:20" json:"platform"` // ios, android, web
	DeviceID  string    `gorm:"size:200" json:"device_id,omitempty"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for DeviceToken
func (DeviceToken) TableName() string {
	return "device_tokens"
}
