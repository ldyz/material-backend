package auth

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"gorm.io/gorm"
)

// roleMapping maps English role values to Chinese role names in the database
var roleMapping = map[string]string{
	"foreman":                 "施工员",
	"keeper":                  "保管员",
	"project_manager":         "项目经理",
	"worker":                  "作业人员",
	"material_staff":          "材料员",
	"subcontractor_material_staff": "分包材料员",
	"appointment_admin":       "预约管理员",
}

// GetCurrentUser fetches current user from DB using context's current_user_id
func GetCurrentUser(c *gin.Context, db *gorm.DB) (*User, error) {
	uid, exists := c.Get("current_user_id")
	if !exists {
		return nil, nil
	}
	var id int64
	switch v := uid.(type) {
	case int64:
		id = v
	case int:
		id = int64(v)
	case float64:
		id = int64(v)
	default:
		return nil, nil
	}
	var user User
	if err := db.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}

	// Fix: If Roles is empty, load role based on role field for backward compatibility
	if len(user.Roles) == 0 && user.Role != "" {
		// Try to find the role by mapping or directly
		var role Role
		roleName := user.Role

		// Check if there's a mapping for this role
		if mappedName, ok := roleMapping[user.Role]; ok {
			roleName = mappedName
		}

		// Try with exact match first, then case-insensitive
		if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
			if err := db.Where("LOWER(name) = ?", strings.ToLower(roleName)).First(&role).Error; err != nil {
				// If still not found, try the original role value
				db.Where("LOWER(name) = ?", strings.ToLower(user.Role)).First(&role)
			}
		}

		// Only set Roles if we found a match
		if role.ID != 0 {
			user.Roles = []Role{role}
		}
	}

	return &user, nil
}

// PermissionMiddleware ensures the current user has the given permission
func PermissionMiddleware(db *gorm.DB, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetCurrentUser(c, db)
		if err != nil || user == nil {
			log.Printf("[PermissionMiddleware] 用户不存在或获取失败: err=%v, user=%v", err, user)
			response.Unauthorized(c, "用户不存在")
			c.Abort()
			return
		}
		log.Printf("[PermissionMiddleware] 用户信息: id=%d, username=%s, role=%s, isAdmin=%v",
			user.ID, user.Username, user.Role, user.IsAdmin())
		if user.IsAdmin() {
			c.Next()
			return
		}
		if !user.HasPermission(permission) {
			log.Printf("[PermissionMiddleware] 权限不足: user=%s, permission=%s", user.Username, permission)
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetAccessibleProjectIDs 获取用户可访问的项目ID列表
// 管理员返回 nil（表示无限制）
// 普通用户返回其关联的项目ID列表
func GetAccessibleProjectIDs(c *gin.Context, db *gorm.DB) ([]uint, error) {
	user, err := GetCurrentUser(c, db)
	if err != nil || user == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 管理员可以访问所有项目
	if user.IsAdmin() {
		return nil, nil // nil 表示无限制
	}

	// 查询用户关联的项目ID列表
	var projectIDs []uint
	err = db.Table("user_projects").
		Where("user_id = ?", user.ID).
		Pluck("project_id", &projectIDs).
		Error

	if err != nil {
		return nil, err
	}

	return projectIDs, nil
}

// ApplyProjectFilter 将项目过滤应用到 GORM 查询
// 如果 projectIDs 为 nil（管理员），不添加过滤
// 如果 projectIDs 为空列表（无权限项目），返回空结果
// 如果 projectIDs 有值，添加 IN 过滤
func ApplyProjectFilter(query *gorm.DB, projectIDs []uint, fieldName string) *gorm.DB {
	if projectIDs == nil {
		// 管理员，不添加过滤
		return query
	}

	if len(projectIDs) == 0 {
		// 用户无任何项目权限，返回空结果
		return query.Where("1 = 0")
	}

	// 添加项目过滤
	return query.Where(fieldName+" IN ?", projectIDs)
}
