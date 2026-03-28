package agent

import (
	"time"

	"gorm.io/gorm"
)

// AIConversation AI对话历史记录
type AIConversation struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	UserID    int64          `gorm:"index;not null" json:"user_id"`
	Role      string         `gorm:"size:20;not null" json:"role"`      // 'user' 或 'assistant'
	Content   string         `gorm:"type:text;not null" json:"content"` // 消息内容
	ToolCalls []byte         `gorm:"type:jsonb" json:"tool_calls"`      // AI工具调用（可选）
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (AIConversation) TableName() string {
	return "ai_conversations"
}

// ConversationRepository 对话历史仓库
type ConversationRepository struct {
	db *gorm.DB
}

// NewConversationRepository 创建对话历史仓库
func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

// SaveMessage 保存一条对话消息
func (r *ConversationRepository) SaveMessage(userID int64, role, content string, toolCalls []byte) error {
	msg := &AIConversation{
		UserID:    userID,
		Role:      role,
		Content:   content,
		ToolCalls: toolCalls,
		CreatedAt: time.Now(),
	}
	return r.db.Create(msg).Error
}

// GetRecentHistory 获取用户最近的对话历史
func (r *ConversationRepository) GetRecentHistory(userID int64, limit int) ([]AIConversation, error) {
	var messages []AIConversation
	if limit <= 0 {
		limit = 20 // 默认获取最近20条
	}
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	// 反转顺序，使旧消息在前
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// ClearHistory 清除用户的对话历史
func (r *ConversationRepository) ClearHistory(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&AIConversation{}).Error
}

// DeleteOldHistory 删除超过指定天数的旧历史记录
func (r *ConversationRepository) DeleteOldHistory(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", cutoff).Delete(&AIConversation{}).Error
}
