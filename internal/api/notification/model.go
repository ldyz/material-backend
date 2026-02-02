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
	TypeRequisitionApprove = "requisition_approve"
	TypeInboundApprove     = "inbound_approve"
	TypeStockAlert         = "stock_alert"
	TypeSystem             = "system"
)

func (n *Notification) ToDTO() map[string]any {
	var readAtStr string
	if n.ReadAt != nil {
		readAtStr = n.ReadAt.Format("2006-01-02 15:04:05")
	}

	return map[string]any{
		"id":         n.ID,
		"user_id":    n.UserID,
		"type":       n.Type,
		"title":      n.Title,
		"content":    n.Content,
		"data":       n.Data,
		"is_read":    n.IsRead,
		"created_at": n.CreatedAt.Format("2006-01-02 15:04:05"),
		"read_at":    readAtStr,
	}
}
