package auth

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"gorm.io/gorm"
)

// simple DTOs
type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type createUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

// RegisterRoutes registers auth routes
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("")
	// public
	r.POST("/auth/login", func(c *gin.Context) {
		var req loginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
		// find user
		var user User
		if err := db.Preload("Roles").Where("username = ?", req.Username).First(&user).Error; err != nil {
			response.Unauthorized(c, "用户名或密码错误")
			return
		}
		if !user.CheckPassword(req.Password) {
			response.Unauthorized(c, "用户名或密码错误")
			return
		}
		if !user.IsActive {
			response.Unauthorized(c, "账户已被禁用")
			return
		}
		// update last login (只更新 LastLogin 字段，避免触发 Roles 关联保存)
		now := time.Now().UTC()
		db.Model(&User{}).Where("id = ?", user.ID).Update("last_login", now)
			// create token
		token, err := jwtpkg.GenerateToken(user.ID, user.Username)
		if err != nil {
			response.InternalError(c, "token 创建失败")
			return
		}

		response.SuccessWithMeta(c, user.ToDTO(), map[string]interface{}{
			"token": token,
		})
	})

	// protected routes
	auth := r.Group("/auth")
	auth.Use(jwtpkg.TokenMiddleware())
	{
		// current user info
		auth.GET("/me", func(c *gin.Context) {
			user, err := GetCurrentUser(c, db)
			if err != nil || user == nil {
				response.NotFound(c, "用户不存在")
				return
			}
			response.Success(c, user.ToDTO())
		})

		// logout (logs to stdout for now)
		auth.POST("/logout", func(c *gin.Context) {
			user, _ := GetCurrentUser(c, db)
			if user != nil {
				// placeholder for system log
				c.Writer.WriteString("user logout: " + user.Username)
			}
			response.SuccessWithMessage(c, nil, "退出成功")
		})

		// change password (token-only required: do not require DB lookup to validate token format)
		auth.POST("/change-password", jwtpkg.TokenOnlyMiddleware(), func(c *gin.Context) {
			var req struct{
				OldPassword string `json:"old_password"`
				NewPassword string `json:"new_password"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, err.Error())
				return
			}
			user, err := GetCurrentUser(c, db)
			if err != nil || user == nil {
				response.NotFound(c, "用户不存在")
				return
			}
			if !user.CheckPassword(req.OldPassword) {
				response.BadRequest(c, "原密码错误")
				return
			}
			if len(req.NewPassword) < 6 {
				response.BadRequest(c, "新密码长度不能少于6位")
				return
			}
			if err := user.SetPassword(req.NewPassword); err != nil {
				response.InternalError(c, "设置密码失败")
				return
			}
			if err := db.Save(user).Error; err != nil {
				response.InternalError(c, err.Error())
				return
			}
			response.SuccessWithMessage(c, nil, "密码修改成功")
		})

		// update avatar
		auth.POST("/avatar", func(c *gin.Context) {
			user, err := GetCurrentUser(c, db)
			if err != nil || user == nil {
				response.NotFound(c, "用户不存在")
				return
			}

			// 获取上传的文件
			file, err := c.FormFile("avatar")
			if err != nil {
				response.BadRequest(c, "请选择要上传的头像文件")
				return
			}

			// 检查文件大小（最大2MB）
			if file.Size > 2*1024*1024 {
				response.BadRequest(c, "头像文件大小不能超过2MB")
				return
			}

			// 检查文件类型
			ext := filepath.Ext(file.Filename)
			allowedExts := map[string]bool{
				".jpg":  true,
				".jpeg": true,
				".png":  true,
				".gif":  true,
				".webp": true,
			}
			if !allowedExts[strings.ToLower(ext)] {
				response.BadRequest(c, "只支持 JPG、PNG、GIF、WEBP 格式的图片")
				return
			}

			// 创建上传目录
			uploadDir := "static/uploads/avatars"
			if err := c.SaveUploadedFile(file, filepath.Join(uploadDir, file.Filename)); err != nil {
				response.InternalError(c, "文件保存失败")
				return
			}

			// 更新用户头像路径
			user.Avatar = "/uploads/avatars/" + file.Filename
			if err := db.Save(user).Error; err != nil {
				response.InternalError(c, "头像更新失败")
				return
			}

			response.SuccessWithMessage(c, map[string]string{"avatar": user.Avatar}, "头像更新成功")
		})

		// USERS CRUD
		// list users
		auth.GET("/users", PermissionMiddleware(db, "user_view"), func(c *gin.Context) {
			page := 1
			perPage := 20
			// naive query params parsing (improve later)
			if p := c.Query("page"); p != "" {
				// ignore errors for brevity
			}
			var total int64
			db.Model(&User{}).Count(&total)
			var users []User
			db.Preload("Roles").Find(&users)
			out := make([]map[string]any, 0, len(users))
			for _, u := range users { out = append(out, u.ToDTO()) }
			response.SuccessWithPagination(c, out, int64(page), int64(perPage), total)
		})

		// create user
		auth.POST("/users", PermissionMiddleware(db, "user_create"), func(c *gin.Context) {
			var req struct {
				Username string `json:"username"`
				Password string `json:"password"`
				Email    string `json:"email"`
				FullName string `json:"full_name"`
				Role     string `json:"role"`
				RoleIDs  []uint `json:"role_ids"`
				IsActive bool   `json:"is_active"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, err.Error())
				return
			}
			u := User{Username: req.Username, Email: req.Email, FullName: req.FullName, Role: req.Role, IsActive: true}
			// 默认启用，只有在请求明确指定false时才禁用
			if req.IsActive == false {
				u.IsActive = false
			}
			if err := u.SetPassword(req.Password); err != nil {
				response.InternalError(c, "设置密码失败")
				return
			}
			if err := db.Create(&u).Error; err != nil {
				response.InternalError(c, err.Error())
				return
			}

			 // 如果提供了 role_ids，使用 many2many 关联
		 if len(req.RoleIDs) > 0 {
			 var roles []Role
			 if err := db.Find(&roles, req.RoleIDs).Error; err != nil {
				 response.InternalError(c, "角色查找失败")
				 return
			 }
			 if err := db.Model(&u).Association("Roles").Append(roles); err != nil {
				 response.InternalError(c, "角色关联失败")
				 return
			 }
			 // 重新加载用户数据
			 db.Preload("Roles").First(&u, u.ID)
		 }

			response.Created(c, u.ToDTO(), "用户创建成功")
		})

		// get single user
		auth.GET("/users/:id", PermissionMiddleware(db, "user_view"), func(c *gin.Context) {
			var u User
			if err := db.Preload("Roles").First(&u, c.Param("id")).Error; err != nil {
				response.NotFound(c, "用户不存在")
				return
			}
			response.Success(c, u.ToDTO())
		})

		// reset password
		auth.POST("/users/:id/reset-password", PermissionMiddleware(db, "user_edit"), func(c *gin.Context) {
			var body struct{ Password string `json:"password"` }
			if err := c.ShouldBindJSON(&body); err != nil || body.Password == "" {
				response.BadRequest(c, "新密码不能为空")
				return
			}
			var u User
			if err := db.First(&u, c.Param("id")).Error; err != nil { response.NotFound(c, "用户不存在"); return }
			if err := u.SetPassword(body.Password); err != nil { response.InternalError(c, "设置密码失败"); return }
			db.Save(&u)
			response.SuccessWithMessage(c, nil, "密码重置成功")
		})

		// update user
		auth.PUT("/users/:id", PermissionMiddleware(db, "user_edit"), func(c *gin.Context) {
			var req map[string]any
			if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, err.Error()); return }
			var u User
			if err := db.Preload("Roles").First(&u, c.Param("id")).Error; err != nil { response.NotFound(c, "用户不存在"); return }
			if v, ok := req["username"].(string); ok && v != "" { u.Username = v }
			if v, ok := req["email"].(string); ok { u.Email = v }
			if v, ok := req["full_name"].(string); ok { u.FullName = v }
			if v, ok := req["role"].(string); ok { u.Role = v }
			if v, ok := req["is_active"].(bool); ok { u.IsActive = v }
			if v, ok := req["password"].(string); ok && v != "" { u.SetPassword(v) }

			// 处理 role_ids many2many 关联
			if roleIDs, ok := req["role_ids"].([]interface{}); ok {
				// 转换为 []uint
				var ids []uint
				for _, id := range roleIDs {
					switch v := id.(type) {
					case float64:
						ids = append(ids, uint(v))
					case int:
						ids = append(ids, uint(v))
					case int64:
						ids = append(ids, uint(v))
					}
				}

				// 查找角色
				var roles []Role
				if len(ids) > 0 {
					if err := db.Find(&roles, ids).Error; err != nil {
						response.InternalError(c, "角色查找失败")
						return
					}
				}

				// 替换关联
				if err := db.Model(&u).Association("Roles").Replace(roles); err != nil {
					response.InternalError(c, "角色更新失败")
					return
				}
			}

			db.Save(&u)
			// 重新加载用户数据以包含更新后的角色
			db.Preload("Roles").First(&u, u.ID)
			response.SuccessWithMessage(c, u.ToDTO(), "用户更新成功")
		})

		// delete user
		auth.DELETE("/users/:id", PermissionMiddleware(db, "user_delete"), func(c *gin.Context) {
			// prevent deleting self
			cur, _ := GetCurrentUser(c, db)
			if cur != nil && cur.ID == 0 {
				// noop
			}
			var u User
			if err := db.First(&u, c.Param("id")).Error; err != nil { response.NotFound(c, "用户不存在"); return }
			if cur != nil && cur.ID == u.ID { response.BadRequest(c, "不能删除自己"); return }
			db.Delete(&u)
			response.SuccessWithMessage(c, nil, "用户删除成功")
		})

		// ROLES CRUD
		auth.GET("/roles", PermissionMiddleware(db, "role_view"), func(c *gin.Context) {
			var roles []Role
			db.Find(&roles)
			out := make([]map[string]any, 0, len(roles))
			for _, r := range roles {
				perms := []string{}
				if r.Permissions != "" {
					for _, it := range strings.Split(r.Permissions, ",") { if it = strings.TrimSpace(it); it != "" { perms = append(perms, it) } }
				}
				out = append(out, map[string]any{"id": r.ID, "name": r.Name, "description": r.Description, "permissions": perms, "created_at": r.CreatedAt.Format("2006-01-02 15:04:05")})
			}
			response.Success(c, out)
		})

		auth.POST("/roles", PermissionMiddleware(db, "role_create"), func(c *gin.Context) {
			var body struct{ Name, Description string; Permissions []string }
			if err := c.ShouldBindJSON(&body); err != nil || body.Name == "" { response.BadRequest(c, "角色名称不能为空"); return }
			r := Role{Name: body.Name, Description: body.Description, Permissions: strings.Join(body.Permissions, ",")}
			if err := db.Create(&r).Error; err != nil { response.InternalError(c, err.Error()); return }
			roleData := map[string]any{"id": r.ID, "name": r.Name, "description": r.Description, "permissions": body.Permissions, "created_at": r.CreatedAt.Format("2006-01-02 15:04:05")}
			response.Created(c, roleData, "角色创建成功")
		})

		auth.GET("/roles/:id", PermissionMiddleware(db, "role_view"), func(c *gin.Context) {
			var r Role
			if err := db.First(&r, c.Param("id")).Error; err != nil { response.NotFound(c, "角色不存在"); return }
			perms := []string{}
			if r.Permissions != "" { for _, it := range strings.Split(r.Permissions, ",") { if it = strings.TrimSpace(it); it != "" { perms = append(perms, it) } } }
			data := map[string]any{"id": r.ID, "name": r.Name, "description": r.Description, "permissions": perms, "created_at": r.CreatedAt.Format("2006-01-02 15:04:05")}
			response.Success(c, data)
		})

		auth.PUT("/roles/:id", PermissionMiddleware(db, "role_edit"), func(c *gin.Context) {
			var body struct {
				Name          string   `json:"name"`
				Description   string   `json:"description"`
				PermissionIDs []string `json:"permission_ids"`
			}
			if err := c.ShouldBindJSON(&body); err != nil { response.BadRequest(c, err.Error()); return }
			var r Role
			if err := db.First(&r, c.Param("id")).Error; err != nil { response.NotFound(c, "角色不存在"); return }
			if body.Name != "" && body.Name != r.Name {
				// check unique
				var exists Role
				if db.Where("name = ?", body.Name).First(&exists); exists.ID != 0 { response.BadRequest(c, "角色名已存在"); return }
				r.Name = body.Name
			}
			if body.Description != "" { r.Description = body.Description }
			if body.PermissionIDs != nil { r.Permissions = strings.Join(body.PermissionIDs, ",") }
			db.Save(&r)
			roleData := map[string]any{"id": r.ID, "name": r.Name, "description": r.Description, "permissions": strings.Split(r.Permissions, ",")}
			response.SuccessWithMessage(c, roleData, "角色更新成功")
		})

		// Assign permissions to a role
		auth.POST("/roles/:id/permissions", PermissionMiddleware(db, "role_assign_permissions"), func(c *gin.Context) {
			var body struct {
				Permissions []string `json:"permissions"`
			}
			if err := c.ShouldBindJSON(&body); err != nil { response.BadRequest(c, err.Error()); return }

			var r Role
			if err := db.First(&r, c.Param("id")).Error; err != nil { response.NotFound(c, "角色不存在"); return }

			// Update role permissions
			r.Permissions = strings.Join(body.Permissions, ",")
			if err := db.Save(&r).Error; err != nil { response.InternalError(c, "权限配置失败"); return }

			roleData := map[string]any{
				"id": r.ID,
				"name": r.Name,
				"description": r.Description,
				"permissions": body.Permissions,
			}
			response.SuccessWithMessage(c, roleData, "权限配置成功")
		})

		auth.DELETE("/roles/:id", PermissionMiddleware(db, "role_delete"), func(c *gin.Context) {
			var r Role
			if err := db.First(&r, c.Param("id")).Error; err != nil { response.NotFound(c, "角色不存在"); return }
			// check users assigned to role
			var count int64
			db.Table("user_roles").Where("role_id = ?", r.ID).Count(&count)
			if count > 0 { response.BadRequest(c, "该角色仍被用户使用，无法删除"); return }
			db.Delete(&r)
			response.SuccessWithMessage(c, nil, "角色删除成功")
		})

		// permissions list
		auth.GET("/permissions", PermissionMiddleware(db, "role_view"), func(c *gin.Context) {
			// Define flat permission list first
			permissionList := []struct {
				Key  string
				Name string
			}{
				// 用户管理 (4)
				{"user_view", "查看用户"},
				{"user_create", "创建用户"},
				{"user_edit", "编辑用户"},
				{"user_delete", "删除用户"},
				// 角色管理 (5)
				{"role_view", "查看角色"},
				{"role_create", "创建角色"},
				{"role_edit", "编辑角色"},
				{"role_delete", "删除角色"},
				{"role_assign_permissions", "分配权限"},
				// 项目管理 (4)
				{"project_view", "查看项目"},
				{"project_create", "创建项目"},
				{"project_edit", "编辑项目"},
				{"project_delete", "删除项目"},
				// 物资管理 (5)
				{"material_view", "查看物资"},
				{"material_create", "创建物资"},
				{"material_edit", "编辑物资"},
				{"material_delete", "删除物资"},
				{"material_import", "导入物资"},
				// 物资计划 (5)
				{"material_plan_view", "查看物资计划"},
				{"material_plan_create", "创建物资计划"},
				{"material_plan_edit", "编辑物资计划"},
				{"material_plan_delete", "删除物资计划"},
				{"material_plan_approve", "审核物资计划"},
				// 库存管理 (8)
				{"stock_view", "查看库存"},
				{"stock_create", "创建库存"},
				{"stock_edit", "编辑库存"},
				{"stock_delete", "删除库存"},
				{"stock_in", "库存入库"},
				{"stock_out", "库存出库"},
				{"stock_export", "导出库存"},
				{"stock_alerts", "库存预警"},
				// 库存日志 (2)
				{"stocklog_view", "查看库存日志"},
				{"stocklog_delete", "删除库存日志"},
				// 入库管理 (6)
				{"inbound_view", "查看入库单"},
				{"inbound_create", "创建入库单"},
				{"inbound_edit", "编辑入库单"},
				{"inbound_delete", "删除入库单"},
				{"inbound_approve", "审核入库单"},
				{"inbound_export", "导出入库单"},
				// 出库管理 (7)
				{"requisition_view", "查看出库单"},
				{"requisition_create", "创建出库单"},
				{"requisition_edit", "编辑出库单"},
				{"requisition_delete", "删除出库单"},
				{"requisition_approve", "审核出库单"},
				{"requisition_issue", "发货"},
				{"requisition_export", "导出出库单"},
				// 施工日志 (5)
				{"construction_log_view", "查看日志"},
				{"construction_log_create", "创建日志"},
				{"construction_log_edit", "编辑日志"},
				{"construction_log_delete", "删除日志"},
				{"construction_log_export", "导出日志"},
				// 进度管理 (5)
				{"progress_view", "查看进度"},
				{"progress_create", "创建进度"},
				{"progress_edit", "编辑进度"},
				{"progress_delete", "删除进度"},
				{"progress_export", "导出进度"},
				// 审计日志 (1)
				{"audit_view", "查看审计日志"},
				// AI 智能体 (5)
				{"ai_agent_view", "查看 AI"},
				{"ai_agent_query", "AI 查询"},
				{"ai_agent_operate", "AI 操作"},
				{"ai_agent_workflow", "AI 工作流"},
				{"ai_agent_logs", "AI 日志"},
				// 系统管理 (5)
				{"system_log", "查看系统日志"},
				{"system_backup", "数据备份"},
				{"system_config", "系统配置"},
				{"system_statistics", "系统统计"},
				{"system_activities", "系统动态"},
				// 工作流管理 (12)
				{"workflow_view", "查看工作流"},
				{"workflow_create", "创建工作流"},
				{"workflow_edit", "编辑工作流"},
				{"workflow_delete", "删除工作流"},
				{"workflow_activate", "激活工作流"},
				{"workflow_instance_view", "查看实例"},
				{"workflow_instance_resubmit", "重新提交"},
				{"workflow_task_view", "查看任务"},
				{"workflow_task_approve", "审批任务"},
				{"workflow_task_reject", "拒绝任务"},
				{"workflow_task_delegate", "委派任务"},
				{"workflow_log_view", "查看流程日志"},
			}

			// Group permissions by module
			groups := make(map[string]map[string]string)
			for _, p := range permissionList {
				// Extract module name from permission key
				var module string
				if len(p.Key) > 0 {
					// Find the first underscore to separate module from action
					idx := strings.Index(p.Key, "_")
					if idx > 0 {
						module = p.Key[:idx]
					} else {
						module = "other"
					}
				} else {
					module = "other"
				}

				// Map module keys to Chinese names
				var moduleName string
				switch module {
				case "user":
					moduleName = "用户管理"
				case "role":
					moduleName = "角色管理"
				case "project":
					moduleName = "项目管理"
				case "material":
					moduleName = "物资管理"
				case "stock", "stocklog":
					moduleName = "库存管理"
				case "inbound":
					moduleName = "入库管理"
				case "requisition":
					moduleName = "出库管理"
				case "progress":
					moduleName = "进度管理"
				case "construction_log":
					moduleName = "施工日志"
				case "audit":
					moduleName = "审计日志"
				case "ai_agent":
					moduleName = "AI 智能体"
				case "workflow":
					moduleName = "工作流管理"
				case "system":
					moduleName = "系统管理"
				default:
					moduleName = "其他"
				}

				if groups[moduleName] == nil {
					groups[moduleName] = make(map[string]string)
				}
				groups[moduleName][p.Key] = p.Name
			}

			// Convert to response format
			data := make([]map[string]any, 0, len(groups))
			for moduleName, perms := range groups {
				data = append(data, map[string]any{
					"key":   moduleName,
					"label": perms,
				})
			}

			response.Success(c, data)
		})
	}
}