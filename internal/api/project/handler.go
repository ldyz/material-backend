package project

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("/project")
	// require token for all project routes
	r.Use(jwtpkg.TokenMiddleware())

	// list projects (supports search, filters, pagination, sorting)
	r.GET("/projects", auth.PermissionMiddleware(db, "project_view"), func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		search := c.Query("search")
		status := c.Query("status")
		manager := c.Query("manager")
		startFrom := c.Query("start_from")
		startTo := c.Query("start_to")
		parentID := c.Query("parent_id")
		sort := c.DefaultQuery("sort", "-id")
		showAll := c.Query("show_all")

		var total int64
		query := db.Model(&Project{})
		if search != "" {
			query = query.Where("name LIKE ?", "%"+search+"%")
		}
		if status != "" {
			query = query.Where("status = ?", status)
		}
		if manager != "" {
			query = query.Where("manager = ?", manager)
		}
		if startFrom != "" {
			if t, err := time.Parse("2006-01-02", startFrom); err == nil { query = query.Where("start_date >= ?", t) }
		}
		if startTo != "" {
			if t, err := time.Parse("2006-01-02", startTo); err == nil { query = query.Where("start_date <= ?", t) }
		}
		if parentID != "" {
			query = query.Where("parent_id = ?", parentID)
		}

		// 获取用户可访问的项目ID列表（数据权限过滤）
		// 如果 show_all 为 true，则跳过权限过滤（允许查看所有项目）
		currentUser, err := auth.GetCurrentUser(c, db)
		if showAll != "true" && err == nil && currentUser != nil && !currentUser.IsAdmin() {
			// 非管理员用户，只能看到自己关联的项目
			query = query.Where("id IN (SELECT project_id FROM user_projects WHERE user_id = ?)", currentUser.ID)
		}

		// sort handling with whitelist
		orderField := "id"
		dir := "DESC"
		if strings.HasPrefix(sort, "-") { orderField = strings.TrimPrefix(sort, "-"); dir = "DESC" } else { orderField = sort; dir = "ASC" }
		allowed := map[string]bool{"name": true, "id": true, "start_date": true, "code": true}
		if !allowed[orderField] { orderField = "id" }
		query = query.Order(orderField + " " + dir)

		query.Count(&total)
		var projects []Project
		query.Offset((page-1)*pageSize).Limit(pageSize).Find(&projects)
		out := make([]map[string]any, 0, len(projects))
		for _, p := range projects { out = append(out, p.ToDTO()) }
		response.SuccessWithPagination(c, out, int64(page), int64(pageSize), total)
	})

	// create project
	r.POST("/projects", auth.PermissionMiddleware(db, "project_create"), func(c *gin.Context) {
		var req struct{
			Name string `json:"name"`
			Code string `json:"code"`
			Description string `json:"description"`
			Location string `json:"location"`
			Manager string `json:"manager"`
			Contact string `json:"contact"`
			Budget string `json:"budget"`
			Status string `json:"status"`
			StartDate string `json:"start_date"`
			EndDate string `json:"end_date"`
			MemberIDs []uint `json:"member_ids"`
			ParentID *uint `json:"parent_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
		if req.Name == "" {
			response.BadRequest(c, "项目名称不能为空")
			return
		}
		// uniqueness checks
		var existing Project
		if err := db.Where("name = ?", req.Name).First(&existing).Error; err == nil { response.BadRequest(c, "项目名已存在"); return }
		if req.Code != "" {
			if err := db.Where("code = ?", req.Code).First(&existing).Error; err == nil { response.BadRequest(c, "项目编号已存在"); return }
		}
		// validate parent exists
		if req.ParentID != nil {
			var parent Project
			if err := db.First(&parent, *req.ParentID).Error; err != nil { response.BadRequest(c, "父项目不存在"); return }
			if parent.Level >= 3 { response.BadRequest(c, "父项目层级已达上限，无法创建子项目"); return }
		}
		// parse dates
		var sd *time.Time
		var ed *time.Time
		if req.StartDate != "" {
			if t, err := time.Parse("2006-01-02", req.StartDate); err == nil { sd = &t } else { response.BadRequest(c, "开始日期格式不正确"); return }
		}
		if req.EndDate != "" {
			if t, err := time.Parse("2006-01-02", req.EndDate); err == nil { ed = &t } else { response.BadRequest(c, "结束日期格式不正确"); return }
		}
		code := req.Code
		if code == "" {
			code = generateProjectNo(db)
		}
		status := req.Status
		if status == "" {
			status = "planning"
		}
		p := Project{
			Name: req.Name,
			Code: code,
			Description: req.Description,
			Location: req.Location,
			Manager: req.Manager,
			Contact: req.Contact,
			Budget: req.Budget,
			Status: status,
			StartDate: sd,
			EndDate: ed,
			ParentID: req.ParentID,
		}
		if err := db.Create(&p).Error; err != nil { response.InternalError(c, err.Error()); return }

		// 处理成员关联（按优先级）
		// 1. 如果明确提供了member_ids（包括空数组），使用提供的成员
		// 2. 如果没有提供member_ids但有parent_id，从父项目继承成员
		// 3. 如果两者都没有，添加创建者为成员
		var memberIDs []uint
		memberIDsProvided := false

		if req.MemberIDs != nil {
			// 明确提供了member_ids（可能是空数组）
			memberIDsProvided = true
			memberIDs = req.MemberIDs
		}

		if memberIDsProvided || len(memberIDs) > 0 {
			// 使用提供的成员或继承的成员
			var users []auth.User
			if err := db.Where("id IN ?", memberIDs).Find(&users).Error; err != nil {
				response.InternalError(c, "查询用户失败: "+err.Error())
				return
			}
			// 清除现有关联并添加新成员
			if err := db.Model(&p).Association("Users").Clear(); err != nil {
				response.InternalError(c, "清除成员失败: "+err.Error())
				return
			}
			if len(users) > 0 {
				if err := db.Model(&p).Association("Users").Append(users); err != nil {
					response.InternalError(c, "添加成员失败: "+err.Error())
					return
				}
			}
		} else if req.ParentID != nil {
			// 从父项目继承成员
			var parent Project
			if err := db.Preload("Users").First(&parent, *req.ParentID).Error; err == nil {
				if len(parent.Users) > 0 {
					// 清除现有关联并添加父项目的成员
					if err := db.Model(&p).Association("Users").Clear(); err != nil {
						response.InternalError(c, "清除成员失败: "+err.Error())
						return
					}
					if err := db.Model(&p).Association("Users").Append(parent.Users); err != nil {
						response.InternalError(c, "继承成员失败: "+err.Error())
						return
					}
				}
			}
		} else {
			// 没有提供member_ids也没有父项目，添加创建者为成员
			currentUser, err := auth.GetCurrentUser(c, db)
			if err == nil && currentUser != nil {
				if err := db.Model(&p).Association("Users").Append([]auth.User{*currentUser}); err != nil {
					response.InternalError(c, "添加创建者失败: "+err.Error())
					return
				}
			}
		}

		// reload project with users
		if err := db.Preload("Users").First(&p, p.ID).Error; err != nil {
			response.InternalError(c, "重新加载项目失败: "+err.Error())
			return
		}
		response.Created(c, p.ToDTO(), "项目创建成功")
	})

	// get project
	r.GET("/projects/:id", auth.PermissionMiddleware(db, "project_view"), func(c *gin.Context) {
		id := c.Param("id")
		var p Project
		if err := db.First(&p, id).Error; err != nil { response.NotFound(c, "项目不存在"); return }
		response.Success(c, p.ToDTO())
	})

	// update
	r.PUT("/projects/:id", auth.PermissionMiddleware(db, "project_edit"), func(c *gin.Context) {
		id := c.Param("id")
		var p Project
		if err := db.First(&p, id).Error; err != nil { response.NotFound(c, "项目不存在"); return }
		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, err.Error()); return }
		
		// 打印接收到的请求参数
		log.Printf("[PROJECT UPDATE] 项目ID: %s", id)
		log.Printf("[PROJECT UPDATE] 请求参数: %+v", req)
		for k, v := range req {
			log.Printf("[PROJECT UPDATE]   %s: %v (类型: %T)", k, v, v)
		}
		
		updates := map[string]any{}
		
		if v, ok := req["name"].(string); ok && v != "" {
			if v != p.Name {
				var ex Project
				if err := db.Where("name = ?", v).First(&ex).Error; err == nil && ex.ID != 0 {
					response.BadRequest(c, "项目名已存在")
					return
				}
			}
			updates["name"] = v
		}
		if v, ok := req["code"].(string); ok {
			if v != p.Code {
				var ex Project
				if err := db.Where("code = ?", v).First(&ex).Error; err == nil && ex.ID != 0 {
					response.BadRequest(c, "项目编号已存在")
					return
				}
			}
			updates["code"] = v
		}
		if v, ok := req["description"].(string); ok { updates["description"] = v }
		if v, ok := req["location"].(string); ok { updates["location"] = v }
		if v, ok := req["manager"].(string); ok { updates["manager"] = v }
		if v, ok := req["contact"].(string); ok { updates["contact"] = v }

		// 处理 parent_id
		if parentIDVal, ok := req["parent_id"]; ok {
			if parentIDVal == nil {
				// 设置为 null（移除父项目）
				updates["parent_id"] = nil
			} else {
				// 尝试解析为数字
				switch v := parentIDVal.(type) {
				case float64:
					parentID := uint(v)
					// 验证父项目存在
					if parentID > 0 {
						var parent Project
						if err := db.First(&parent, parentID).Error; err != nil {
							response.BadRequest(c, "父项目不存在")
							return
						}
						// 检查层级限制
						if parent.Level >= 3 {
							response.BadRequest(c, "父项目层级已达上限，无法添加子项目")
							return
						}
						// 防止循环引用（项目不能是自己的子孙）
						if parentID == p.ID {
							response.BadRequest(c, "项目不能设置自己为父项目")
							return
						}
						updates["parent_id"] = parentID
					} else {
						updates["parent_id"] = nil
					}
				case int:
					parentID := uint(v)
					if parentID > 0 {
						var parent Project
						if err := db.First(&parent, parentID).Error; err != nil {
							response.BadRequest(c, "父项目不存在")
							return
						}
						if parent.Level >= 3 {
							response.BadRequest(c, "父项目层级已达上限，无法添加子项目")
							return
						}
						if parentID == p.ID {
							response.BadRequest(c, "项目不能设置自己为父项目")
							return
						}
						updates["parent_id"] = parentID
					} else {
						updates["parent_id"] = nil
					}
				}
			}
		}
		
		// 处理budget - 可以是string或float64
		if budgetVal, ok := req["budget"]; ok {
			budgetStr := ""
			switch v := budgetVal.(type) {
			case string:
				budgetStr = v
			case float64:
				budgetStr = fmt.Sprintf("%.2f", v)
			}
			if budgetStr != "" {
				updates["budget"] = budgetStr
			}
		}
		
		if v, ok := req["status"].(string); ok { updates["status"] = v }
		
		// 打印要更新的字段
		log.Printf("[PROJECT UPDATE] 待更新字段: %+v", updates)
		for k, v := range updates {
			log.Printf("[PROJECT UPDATE]   %s: %v", k, v)
		}
		
		// apply updates to database
		if len(updates) > 0 {
			log.Printf("[PROJECT UPDATE] 执行数据库UPDATE操作...")
			if err := db.Model(&p).Updates(updates).Error; err != nil {
				log.Printf("[PROJECT UPDATE] 更新失败: %v", err)
				response.InternalError(c, "项目更新失败: "+err.Error())
				return
			}
			log.Printf("[PROJECT UPDATE] UPDATE成功")
		} else {
			log.Printf("[PROJECT UPDATE] 没有字段需要更新")
		}

		// handle member_ids if provided
		if memberIDs, ok := req["member_ids"].([]interface{}); ok && len(memberIDs) > 0 {
			var userIDs []uint
			for _, id := range memberIDs {
				switch v := id.(type) {
				case float64:
					userIDs = append(userIDs, uint(v))
				case int:
					userIDs = append(userIDs, uint(v))
				}
			}
			if len(userIDs) > 0 {
				var users []auth.User
				if err := db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
					response.InternalError(c, "查询用户失败: "+err.Error())
					return
				}
				// clear existing association first
				if err := db.Model(&p).Association("Users").Clear(); err != nil {
					response.InternalError(c, "清除成员失败: "+err.Error())
					return
				}
				// add new users
				if len(users) > 0 {
					if err := db.Model(&p).Association("Users").Append(users); err != nil {
						response.InternalError(c, "添加成员失败: "+err.Error())
						return
					}
				}
			}
		}

		// reload project with users
		if err := db.Preload("Users").First(&p, id).Error; err != nil {
			response.InternalError(c, "重新加载项目失败: "+err.Error())
			return
		}
		response.SuccessWithMessage(c, p.ToDTO(), "项目更新成功")
	})

	// delete (check for linked materials first)
	r.DELETE("/projects/:id", auth.PermissionMiddleware(db, "project_delete"), func(c *gin.Context) {
		id := c.Param("id")
		var p Project
		if err := db.First(&p, id).Error; err != nil { response.NotFound(c, "项目不存在"); return }
		// TODO: check linked materials (Material model implementation will do this)
		// For now, just proceed with delete (update this when Material module is added)
		if err := db.Delete(&p).Error; err != nil {
			response.InternalError(c, "项目删除失败: "+err.Error())
			return
		}
		response.SuccessWithMessage(c, nil, "项目删除成功")
	})

	// members list
	r.GET("/projects/:id/members", auth.PermissionMiddleware(db, "project_view"), func(c *gin.Context) {
		id := c.Param("id")
		var p Project
		if err := db.Preload("Users").First(&p, id).Error; err != nil { response.NotFound(c, "项目不存在"); return }
		out := make([]map[string]any, 0, len(p.Users))
		for _, u := range p.Users { out = append(out, u.ToDTO()) }
		response.Success(c, out)
	})

	// add members (clear existing and replace)
	// 注意：允许空数组，表示清空所有成员
	r.POST("/projects/:id/members", auth.PermissionMiddleware(db, "project_edit"), func(c *gin.Context) {
		id := c.Param("id")
		var body struct{ UserIDs []uint `json:"user_ids"` }
		if err := c.ShouldBindJSON(&body); err != nil { response.BadRequest(c, err.Error()); return }
		var p Project
		if err := db.First(&p, id).Error; err != nil { response.NotFound(c, "项目不存在"); return }

		// clear existing association first
		if err := db.Model(&p).Association("Users").Clear(); err != nil {
			response.InternalError(c, "清除成员失败: "+err.Error())
			return
		}

		// add new users (if any)
		if len(body.UserIDs) > 0 {
			var users []auth.User
			if err := db.Where("id IN ?", body.UserIDs).Find(&users).Error; err != nil {
				response.InternalError(c, "查询用户失败: "+err.Error())
				return
			}
			if len(users) > 0 {
				if err := db.Model(&p).Association("Users").Append(users); err != nil {
					response.InternalError(c, "添加成员失败: "+err.Error())
					return
				}
			}
		}
		response.SuccessWithMessage(c, nil, "成员分配成功")
	})

	// remove member
	r.DELETE("/projects/:id/members/:user_id", auth.PermissionMiddleware(db, "project_delete"), func(c *gin.Context) {
		id := c.Param("id")
		uid := c.Param("user_id")
		var p Project
		if err := db.Preload("Users").First(&p, id).Error; err != nil { response.NotFound(c, "项目不存在"); return }
		var user auth.User
		if err := db.First(&user, uid).Error; err != nil { response.NotFound(c, "用户不存在"); return }
		// remove association
		if err := db.Model(&p).Association("Users").Delete(&user); err != nil { response.InternalError(c, err.Error()); return }
		response.SuccessWithMessage(c, nil, "成员已移除")
	})

	// get project tree (recursive hierarchy)
	r.GET("/projects/:id/tree", auth.PermissionMiddleware(db, "project_view"), func(c *gin.Context) {
		id := c.Param("id")
		var root Project
		if err := db.First(&root, id).Error; err != nil { response.NotFound(c, "项目不存在"); return }

		// Build the tree recursively
		tree := buildProjectTree(db, root.ID)
		response.Success(c, tree)
	})

	// get direct children of a project
	r.GET("/projects/:id/children", auth.PermissionMiddleware(db, "project_view"), func(c *gin.Context) {
		id := c.Param("id")
		var children []Project
		if err := db.Where("parent_id = ?", id).Order("level ASC, id ASC").Find(&children).Error; err != nil {
			response.InternalError(c, "查询子项目失败")
			return
		}

		out := make([]map[string]any, 0, len(children))
		for _, child := range children {
			out = append(out, child.ToDTO())
		}
		response.Success(c, out)
	})

	// aggregate progress from child projects
	r.POST("/projects/:id/aggregate-progress", auth.PermissionMiddleware(db, "project_edit"), func(c *gin.Context) {
		id := c.Param("id")
		var project Project
		if err := db.First(&project, id).Error; err != nil { response.NotFound(c, "项目不存在"); return }

		// Recursively aggregate progress from all descendants
		progress := aggregateProgress(db, project.ID)

		// Update the project's progress
		if err := db.Model(&project).Update("progress_percentage", progress).Error; err != nil {
			response.InternalError(c, "更新进度失败")
			return
		}

		response.SuccessWithMessage(c, map[string]any{"progress_percentage": progress}, "进度聚合成功")
	})
}

// buildProjectTree recursively builds the project tree
func buildProjectTree(db *gorm.DB, projectID uint) map[string]any {
	var project Project
	db.First(&project, projectID)

	// Get children
	var children []Project
	db.Where("parent_id = ?", projectID).Order("level ASC, id ASC").Find(&children)

	// Build result
	result := project.ToDTO()
	childTrees := make([]map[string]any, 0, len(children))
	for _, child := range children {
		childTrees = append(childTrees, buildProjectTree(db, child.ID))
	}
	result["children"] = childTrees

	return result
}

// aggregateProgress recursively calculates aggregated progress from all descendants
func aggregateProgress(db *gorm.DB, projectID uint) float64 {
	// Get direct children
	var children []Project
	db.Where("parent_id = ?", projectID).Find(&children)

	if len(children) == 0 {
		// No children, return project's own progress
		var project Project
		db.First(&project, projectID)
		return project.ProgressPercentage
	}

	// Calculate weighted average based on child project progress
	totalProgress := 0.0
	for _, child := range children {
		childProgress := aggregateProgress(db, child.ID)
		totalProgress += childProgress
	}

	return totalProgress / float64(len(children))
}

// generateProjectNo makes a code like PJ-YYYYMMDDNNN
func generateProjectNo(db *gorm.DB) string {
	// generates unique code like PJ-YYYYMMDDNNN; tries several times to avoid collisions
	today := time.Now().Format("20060102")
	prefix := "PJ-" + today
	// start with count+1
	var baseCount int64
	db.Model(&Project{}).Where("code LIKE ?", prefix+"%").Count(&baseCount)
	try := baseCount + 1
	for i := 0; i < 20; i++ {
		seq := strconv.FormatInt(try+int64(i), 10)
		if len(seq) < 3 { seq = strings.Repeat("0", 3-len(seq)) + seq }
		code := prefix + seq
		var c int64
		db.Model(&Project{}).Where("code = ?", code).Count(&c)
		if c == 0 {
			return code
		}
	}
	// fallback to timestamp-based suffix
	return prefix + strconv.FormatInt(time.Now().UnixNano()%1000000, 10)
}
