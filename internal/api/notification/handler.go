package notification

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("notification")
	r.Use(jwtpkg.TokenMiddleware())

	// Get user notifications
	r.GET("/notifications", func(c *gin.Context) {
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil {
			response.Unauthorized(c, "未授权")
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}
		unreadOnly := c.DefaultQuery("unread_only", "false") == "true"

		var total int64
		query := db.Model(&Notification{}).Where("user_id = ?", currentUser.ID)

		if unreadOnly {
			query = query.Where("is_read = ?", false)
		}

		query.Count(&total)

		var notifications []Notification
		query.Order("created_at DESC").
			Offset((page-1)*pageSize).
			Limit(pageSize).
			Find(&notifications)

		// Convert to DTO
		items := make([]map[string]any, len(notifications))
		for i, n := range notifications {
			items[i] = n.ToDTO()
		}

		// Get unread count
		var unreadCount int64
		db.Model(&Notification{}).
			Where("user_id = ? AND is_read = ?", currentUser.ID, false).
			Count(&unreadCount)

		// Build meta with unread count
		meta := map[string]interface{}{
			"unread_count": unreadCount,
		}

		response.SuccessWithPaginationAndMeta(c, items, int64(page), int64(pageSize), total, meta)
	})

	// Get notification count
	r.GET("/notifications/count", func(c *gin.Context) {
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil {
			response.Unauthorized(c, "未授权")
			return
		}

		var unreadCount int64
		db.Model(&Notification{}).
			Where("user_id = ? AND is_read = ?", currentUser.ID, false).
			Count(&unreadCount)

		data := map[string]interface{}{
			"unread_count": unreadCount,
		}
		response.Success(c, data)
	})

	// Mark notification as read
	r.PUT("/notifications/:id/read", func(c *gin.Context) {
		id := c.Param("id")
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil {
			response.Unauthorized(c, "未授权")
			return
		}

		var notification Notification
		if err := db.Where("id = ? AND user_id = ?", id, currentUser.ID).First(&notification).Error; err != nil {
			response.NotFound(c, "通知不存在")
			return
		}

		if !notification.IsRead {
			now := time.Now()
			notification.IsRead = true
			notification.ReadAt = &now
			db.Save(&notification)
		}

		response.SuccessWithMessage(c, nil, "已标记为已读")
	})

	// Mark all notifications as read
	r.PUT("/notifications/read-all", func(c *gin.Context) {
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil {
			response.Unauthorized(c, "未授权")
			return
		}

		now := time.Now()
		db.Model(&Notification{}).
			Where("user_id = ? AND is_read = ?", currentUser.ID, false).
			Updates(map[string]interface{}{
				"is_read": true,
				"read_at": now,
			})

		response.SuccessWithMessage(c, nil, "已全部标记为已读")
	})

	// Delete notification
	r.DELETE("/notifications/:id", func(c *gin.Context) {
		id := c.Param("id")
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil {
			response.Unauthorized(c, "未授权")
			return
		}

		result := db.Where("id = ? AND user_id = ?", id, currentUser.ID).Delete(&Notification{})
		if result.RowsAffected == 0 {
			response.NotFound(c, "通知不存在")
			return
		}

		response.SuccessWithMessage(c, nil, "通知已删除")
	})

	// Clear all notifications
	r.DELETE("/notifications", func(c *gin.Context) {
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil {
			response.Unauthorized(c, "未授权")
			return
		}

		db.Where("user_id = ?", currentUser.ID).Delete(&Notification{})

		response.SuccessWithMessage(c, nil, "通知已清空")
	})
}

// CreateNotification creates a new notification for specified users
func CreateNotification(db *gorm.DB, userID uint, notificationType, title, content string, data map[string]interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	notification := Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Content: content,
		Data:    string(dataJSON),
		IsRead:  false,
	}

	return db.Create(&notification).Error
}

// NotifyUsersWithPermission sends notification to all users with specific permission
func NotifyUsersWithPermission(db *gorm.DB, permission, notificationType, title, content string, data map[string]interface{}) error {
	// Get all users with the specified permission
	var users []auth.User
	if err := db.Raw(`
		SELECT DISTINCT u.* FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		JOIN roles r ON ur.role_id = r.id
		WHERE r.permissions LIKE ?
	`, "%"+permission+"%").Find(&users).Error; err != nil {
		return err
	}

	// Create notification for each user
	for _, user := range users {
		if err := CreateNotification(db, user.ID, notificationType, title, content, data); err != nil {
			return err
		}
	}

	return nil
}
