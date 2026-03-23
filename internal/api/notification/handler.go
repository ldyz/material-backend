package notification

import (
	"encoding/json"
	"log"
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

	// WebSocket endpoint for real-time notifications
	r.GET("/ws", func(c *gin.Context) {
		// Set user_id in context for WebSocket handler
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("user_id", currentUser.ID)

		// Get the hub and serve WebSocket
		hub := GetHub()
		if hub == nil {
			c.JSON(500, gin.H{"error": "WebSocket hub not initialized"})
			return
		}
		ServeWS(hub)(c)
	})

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

	// Register push notification token
	r.POST("/register-token", func(c *gin.Context) {
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil {
			response.Unauthorized(c, "未授权")
			return
		}

		var req struct {
			Token    string `json:"token" binding:"required"`
			Platform string `json:"platform" binding:"required"` // ios, android, web
			DeviceID string `json:"device_id,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "无效的请求参数")
			return
		}

		// Validate platform
		if req.Platform != "ios" && req.Platform != "android" && req.Platform != "web" {
			response.BadRequest(c, "无效的平台类型")
			return
		}

		// Check if token already exists for this user
		var existingToken DeviceToken
		err = db.Where("user_id = ? AND token = ?", currentUser.ID, req.Token).First(&existingToken).Error

		now := time.Now()
		if err == nil {
			// Token exists, update it
			existingToken.Platform = req.Platform
			existingToken.DeviceID = req.DeviceID
			existingToken.IsActive = true
			existingToken.UpdatedAt = now
			if err := db.Save(&existingToken).Error; err != nil {
				log.Printf("Failed to update device token: %v", err)
			}
		} else {
			// New token, create it
			deviceToken := DeviceToken{
				UserID:    currentUser.ID,
				Token:     req.Token,
				Platform:  req.Platform,
				DeviceID:  req.DeviceID,
				IsActive:  true,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := db.Create(&deviceToken).Error; err != nil {
				log.Printf("Failed to create device token: %v", err)
			}
		}

		response.SuccessWithMessage(c, nil, "推送令牌已注册")
	})

	// Unregister push notification token
	r.DELETE("/unregister-token", func(c *gin.Context) {
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil {
			response.Unauthorized(c, "未授权")
			return
		}

		var req struct {
			Token string `json:"token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "无效的请求参数")
			return
		}

		// Deactivate the token
		result := db.Model(&DeviceToken{}).
			Where("user_id = ? AND token = ?", currentUser.ID, req.Token).
			Update("is_active", false)

		if result.Error != nil {
			response.InternalError(c, "注销令牌失败")
			return
		}

		response.SuccessWithMessage(c, nil, "推送令牌已注销")
	})

	// Get WebSocket hub statistics (admin only)
	r.GET("/hub-stats", func(c *gin.Context) {
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil {
			response.Unauthorized(c, "未授权")
			return
		}

		// Check if user is admin
		if currentUser.Username != "admin" {
			response.Forbidden(c, "无权限访问")
			return
		}

		hub := GetHub()
		if hub == nil {
			response.Success(c, map[string]interface{}{
				"total_users":  0,
				"total_clients": 0,
				"online_users": []uint{},
			})
			return
		}

		stats := hub.GetStats()
		response.Success(c, stats)
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

	if err := db.Create(&notification).Error; err != nil {
		return err
	}

	// Broadcast via WebSocket if hub is available
	if hub := GetHub(); hub != nil {
		hub.BroadcastNotification(userID, notification)
	}

	return nil
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
