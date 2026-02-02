package material_plan

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

var (
	ErrUserNotAuthenticated = errors.New("User not authenticated")
	ErrInvalidUserID       = errors.New("Invalid user ID type")
)

// getUserInfo extracts and converts user info from gin context
// Returns userID (uint) and username (string)
func getUserInfo(c *gin.Context) (uint, string, error) {
	userIDInt64, exists := c.Get("current_user_id")
	if !exists {
		return 0, "", ErrUserNotAuthenticated
	}

	var userID uint
	switch v := userIDInt64.(type) {
	case int64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	default:
		return 0, "", ErrInvalidUserID
	}

	userName, _ := c.Get("current_username")
	userNameStr, _ := userName.(string)

	return userID, userNameStr, nil
}

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	g := rg.Group("/material-plan")
	g.Use(jwtpkg.TokenMiddleware())

	service := NewService(db)
	workflow := NewWorkflowIntegration(db)

	// Plan CRUD endpoints
	g.GET("/plans", auth.PermissionMiddleware(db, "material_plan_view"), listPlans(db))
	g.POST("/plans", auth.PermissionMiddleware(db, "material_plan_create"), createPlan(service))
	g.GET("/plans/:id", auth.PermissionMiddleware(db, "material_plan_view"), getPlanDetail(db))
	g.PUT("/plans/:id", auth.PermissionMiddleware(db, "material_plan_edit"), updatePlan(service))
	g.DELETE("/plans/:id", auth.PermissionMiddleware(db, "material_plan_delete"), deletePlan(service))

	// Plan workflow endpoints
	g.POST("/plans/:id/submit", auth.PermissionMiddleware(db, "material_plan_edit"), submitPlan(workflow))
	g.POST("/plans/:id/approve", auth.PermissionMiddleware(db, "material_plan_approve"), approvePlan(workflow))
	g.POST("/plans/:id/reject", auth.PermissionMiddleware(db, "material_plan_approve"), rejectPlan(workflow))
	g.POST("/plans/:id/activate", auth.PermissionMiddleware(db, "material_plan_approve"), activatePlan(service))
	g.POST("/plans/:id/resubmit", auth.PermissionMiddleware(db, "material_plan_edit"), resubmitPlan(workflow))
	g.POST("/plans/:id/cancel", auth.PermissionMiddleware(db, "material_plan_approve"), cancelPlan(workflow))

	// Plan items endpoints
	g.GET("/plans/:id/items", auth.PermissionMiddleware(db, "material_plan_view"), getPlanItems(db))
	g.POST("/plans/:id/items", auth.PermissionMiddleware(db, "material_plan_edit"), addPlanItem(service))
	g.PUT("/items/:id", auth.PermissionMiddleware(db, "material_plan_edit"), updatePlanItem(service))
	g.DELETE("/items/:id", auth.PermissionMiddleware(db, "material_plan_edit"), deletePlanItem(service))

	// Workflow status endpoints
	g.GET("/plans/:id/workflow", auth.PermissionMiddleware(db, "material_plan_view"), getPlanWorkflowStatus(workflow))
	g.GET("/plans/:id/approvals", auth.PermissionMiddleware(db, "material_plan_view"), getPlanApprovals(workflow))
	g.GET("/workflow/pending", auth.PermissionMiddleware(db, "material_plan_approve"), getPendingTasks(workflow))
}

// listPlans lists all material plans with pagination and filters
func listPlans(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		// Filter parameters
		status := strings.TrimSpace(c.Query("status"))
		planType := strings.TrimSpace(c.Query("plan_type"))
		priority := strings.TrimSpace(c.Query("priority"))
		search := strings.TrimSpace(c.Query("search"))
		startDate := strings.TrimSpace(c.Query("start_date"))
		endDate := strings.TrimSpace(c.Query("end_date"))
		projectID := strings.TrimSpace(c.Query("project_id"))

		query := db.Preload("Items")

		// Status filter
		if status != "" {
			query = query.Where("status = ?", status)
		}

		// Type filter
		if planType != "" {
			query = query.Where("plan_type = ?", planType)
		}

		// Priority filter
		if priority != "" {
			query = query.Where("priority = ?", priority)
		}

		// General search
		if search != "" {
			query = query.Where("plan_no LIKE ? OR plan_name LIKE ? OR creator_name LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
		}

		// Project filter
		var projectIDUint uint
		if projectID != "" {
			if pid, err := strconv.ParseUint(projectID, 10, 64); err == nil {
				projectIDUint = uint(pid)
			}
		}

		// 支持多个项目ID（用于包含子项目）
		var projectIDsFilter []uint
		projectIDsStr := c.Query("project_ids")
		if projectIDsStr != "" {
			for _, idStr := range strings.Split(projectIDsStr, ",") {
				if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64); err == nil {
					projectIDsFilter = append(projectIDsFilter, uint(id))
				}
			}
		}

		// 获取用户可访问的项目ID列表（数据权限过滤）
		userProjectIDs, err := auth.GetAccessibleProjectIDs(c, db)
		if err != nil {
			response.InternalError(c, "获取用户项目权限失败")
			return
		}

		// 应用项目过滤
		if userProjectIDs != nil {
			if len(userProjectIDs) == 0 {
				response.SuccessWithPagination(c, []map[string]any{}, int64(page), int64(pageSize), 0)
				return
			}

			if len(projectIDsFilter) > 0 {
				for _, pid := range projectIDsFilter {
					hasAccess := false
					for _, accessibleID := range userProjectIDs {
						if pid == accessibleID {
							hasAccess = true
							break
						}
					}
					if !hasAccess {
						response.Forbidden(c, "无权访问该项目")
						return
					}
				}
				query = query.Where("project_id IN ?", projectIDsFilter)
			} else if projectIDUint > 0 {
				hasAccess := false
				for _, pid := range userProjectIDs {
					if projectIDUint == pid {
						hasAccess = true
						break
					}
				}
				if !hasAccess {
					response.Forbidden(c, "无权访问该项目")
					return
				}
				query = query.Where("project_id = ?", projectIDUint)
			} else {
				query = query.Where("project_id IN ?", userProjectIDs)
			}
		} else {
			if len(projectIDsFilter) > 0 {
				query = query.Where("project_id IN ?", projectIDsFilter)
			} else if projectIDUint > 0 {
				query = query.Where("project_id = ?", projectIDUint)
			}
		}

		// Date range filters
		if startDate != "" {
			if t, err := time.Parse("2006-01-02", startDate); err == nil {
				query = query.Where("created_at >= ?", t)
			}
		}

		if endDate != "" {
			if t, err := time.Parse("2006-01-02", endDate); err == nil {
				t = t.Add(24*time.Hour - time.Second)
				query = query.Where("created_at <= ?", t)
			}
		}

		// Get total count
		var total int64
		query.Model(&MaterialPlan{}).Count(&total)

		// Sort by creation date descending
		query = query.Order("created_at DESC")

		// Pagination
		var plans []MaterialPlan
		query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&plans)

		// Enrich with project names
		enrichedPlans := make([]map[string]any, len(plans))
		for i, plan := range plans {
			dto := plan.ToDTO()
			dto["progress"] = plan.CalculateProgress()

			// Get project name
			var project struct {
				ID   uint
				Name string
			}
			if err := db.Table("projects").Where("id = ?", plan.ProjectID).
				Select("id, name").First(&project).Error; err == nil {
				dto["project_name"] = project.Name
			}

			enrichedPlans[i] = dto
		}

		response.SuccessWithPagination(c, enrichedPlans, int64(page), int64(pageSize), total)
	}
}

// createPlan creates a new material plan
func createPlan(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateMaterialPlanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "Invalid request format: "+err.Error())
			return
		}

		// Get user info from JWT
		userID, userName, err := getUserInfo(c)
		if err != nil {
			response.Unauthorized(c, "User not authenticated")
			return
		}

		// Create plan
		plan, err := service.CreatePlan(&req, userID, userName)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// Load items for response
		service.db.Preload("Items").First(plan, plan.ID)

		response.Created(c, plan.ToDTOWithEnrichment(service.db), "计划创建成功")
	}
}

// getPlanDetail gets a single plan detail
func getPlanDetail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		var plan MaterialPlan
		if err := db.Preload("Items").First(&plan, planID).Error; err != nil {
			response.NotFound(c, "Plan not found")
			return
		}

		// Check project access
		userProjectIDs, err := auth.GetAccessibleProjectIDs(c, db)
		if err != nil {
			response.InternalError(c, "获取用户项目权限失败")
			return
		}

		if userProjectIDs != nil {
			hasAccess := false
			for _, pid := range userProjectIDs {
				if plan.ProjectID == pid {
					hasAccess = true
					break
				}
			}
			if !hasAccess {
				response.Forbidden(c, "无权访问该计划")
				return
			}
		}

		dto := plan.ToDTOWithEnrichment(db)
		dto["progress"] = plan.CalculateProgress()

		response.Success(c, dto)
	}
}

// updatePlan updates an existing plan
func updatePlan(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		var req UpdateMaterialPlanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "Invalid request format: "+err.Error())
			return
		}

		plan, err := service.UpdatePlan(uint(planID), &req)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, plan.ToDTOWithEnrichment(service.db), "计划更新成功")
	}
}

// deletePlan deletes a plan
func deletePlan(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		if err := service.DeletePlan(uint(planID)); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessOnlyMessage(c, "计划删除成功")
	}
}

// submitPlan submits a plan for approval - starts workflow
func submitPlan(workflow *WorkflowIntegration) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		// Get user info
		userID, userName, err := getUserInfo(c)
		if err != nil {
			response.Unauthorized(c, "User not authenticated")
			return
		}

		// Get plan
		var plan MaterialPlan
		if err := workflow.db.First(&plan, planID).Error; err != nil {
			response.BadRequest(c, "计划不存在")
			return
		}

		// Validate plan can be submitted
		if plan.Status != PlanStatusDraft {
			response.BadRequest(c, "只有草稿状态的计划可以提交")
			return
		}

		// Check if plan has items
		var itemCount int64
		workflow.db.Model(&MaterialPlanItem{}).Where("plan_id = ?", planID).Count(&itemCount)
		if itemCount == 0 {
			response.BadRequest(c, "计划必须包含至少一个计划项")
			return
		}

		// Start workflow instead of simple status change
		if err := workflow.StartPlanWorkflow(&plan, userID, userName); err != nil {
			response.BadRequest(c, fmt.Sprintf("启动工作流失败: %s", err.Error()))
			return
		}

		response.SuccessOnlyMessage(c, "计划已提交审核")
	}
}

// approvePlan approves a plan - uses workflow
func approvePlan(workflow *WorkflowIntegration) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		var req struct {
			Remark string `json:"remark"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "Invalid request format")
			return
		}

		// Get user info
		userID, userName, err := getUserInfo(c)
		if err != nil {
			response.Unauthorized(c, "User not authenticated")
			return
		}

		// Process approval through workflow
		if err := workflow.ProcessPlanApproval(uint(planID), userID, userName, "approve", req.Remark); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessOnlyMessage(c, "审批成功")
	}
}

// rejectPlan rejects a plan - uses workflow
func rejectPlan(workflow *WorkflowIntegration) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		var req struct {
			Remark string `json:"remark"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "Invalid request format")
			return
		}

		// Get user info
		userID, userName, err := getUserInfo(c)
		if err != nil {
			response.Unauthorized(c, "User not authenticated")
			return
		}

		// Process rejection through workflow
		if err := workflow.ProcessPlanApproval(uint(planID), userID, userName, "reject", req.Remark); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessOnlyMessage(c, "计划已拒绝")
	}
}

// activatePlan activates an approved plan
func activatePlan(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		// Get plan
		var plan MaterialPlan
		if err := service.db.First(&plan, planID).Error; err != nil {
			response.NotFound(c, "Plan not found")
			return
		}

		// Check if plan can be activated
		if plan.Status != PlanStatusApproved {
			response.BadRequest(c, "只有审批通过的计划才能激活")
			return
		}

		// Update status to active
		plan.Status = PlanStatusActive
		if err := service.db.Save(&plan).Error; err != nil {
			response.InternalError(c, "Failed to activate plan")
			return
		}

		response.SuccessOnlyMessage(c, "计划已激活")
	}
}

// resubmitPlan resubmits a rejected plan
func resubmitPlan(workflow *WorkflowIntegration) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		// Get user info
		userID, userName, err := getUserInfo(c)
		if err != nil {
			response.Unauthorized(c, "User not authenticated")
			return
		}

		if err := workflow.ResubmitPlan(uint(planID), userID, userName); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessOnlyMessage(c, "计划已重新提交")
	}
}

// cancelPlan cancels a plan
func cancelPlan(workflow *WorkflowIntegration) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		var req struct {
			Reason string `json:"reason"`
		}
		c.ShouldBindJSON(&req)

		// Get user info
		userID, userName, err := getUserInfo(c)
		if err != nil {
			response.Unauthorized(c, "User not authenticated")
			return
		}

		if err := workflow.CancelPlanWorkflow(uint(planID), userID, userName, req.Reason); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessOnlyMessage(c, "计划已取消")
	}
}

// getPlanItems gets items for a plan
func getPlanItems(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		var items []MaterialPlanItem
		if err := db.Where("plan_id = ?", planID).Order("sort_order ASC").Find(&items).Error; err != nil {
			response.InternalError(c, "Failed to fetch items")
			return
		}

		result := make([]map[string]any, len(items))
		for i, item := range items {
			result[i] = item.ToDTO()
		}

		response.Success(c, result)
	}
}

// addPlanItem adds an item to a plan
func addPlanItem(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		var req CreateMaterialPlanItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "Invalid request format: "+err.Error())
			return
		}

		// Use transaction to create item
		tx := service.db.Begin()
		item, err := service.CreatePlanItem(tx, uint(planID), &req)
		if err != nil {
			tx.Rollback()
			response.BadRequest(c, err.Error())
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.InternalError(c, "Failed to commit transaction")
			return
		}

		response.Created(c, item.ToDTO(), "项目添加成功")
	}
}

// updatePlanItem updates a plan item
func updatePlanItem(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		itemID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid item ID")
			return
		}

		var req CreateMaterialPlanItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "Invalid request format: "+err.Error())
			return
		}

		// Get existing item
		var item MaterialPlanItem
		if err := service.db.First(&item, itemID).Error; err != nil {
			response.NotFound(c, "Item not found")
			return
		}

		// Parse required date
		var requiredDate *time.Time
		if req.RequiredDate != "" {
			t, err := time.Parse("2006-01-02", req.RequiredDate)
			if err != nil {
				response.BadRequest(c, "Invalid required_date format, use YYYY-MM-DD")
				return
			}
			requiredDate = &t
		}

		// Update fields
		item.MaterialID = req.MaterialID
		item.PlannedQuantity = req.PlannedQuantity
		item.UnitPrice = req.UnitPrice
		item.RequiredDate = requiredDate
		item.Priority = req.Priority
		item.Remark = req.Remark

		if err := service.db.Save(&item).Error; err != nil {
			response.InternalError(c, "Failed to update item")
			return
		}

		response.SuccessWithMessage(c, item.ToDTO(), "项目更新成功")
	}
}

// deletePlanItem deletes a plan item
func deletePlanItem(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		itemID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid item ID")
			return
		}

		if err := service.db.Delete(&MaterialPlanItem{}, itemID).Error; err != nil {
			response.InternalError(c, "Failed to delete item")
			return
		}

		response.SuccessOnlyMessage(c, "项目删除成功")
	}
}

// getPlanWorkflowStatus gets workflow status for a plan
func getPlanWorkflowStatus(workflow *WorkflowIntegration) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		instance, err := workflow.GetPlanWorkflowStatus(uint(planID))
		if err != nil {
			response.NotFound(c, "Workflow instance not found")
			return
		}

		response.Success(c, instance)
	}
}

// getPlanApprovals gets approval history for a plan
func getPlanApprovals(workflow *WorkflowIntegration) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		planID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid plan ID")
			return
		}

		approvals, err := workflow.GetPlanWorkflowApprovals(uint(planID))
		if err != nil {
			response.NotFound(c, "Approvals not found")
			return
		}

		response.Success(c, approvals)
	}
}

// getPendingTasks gets pending tasks for the current user
func getPendingTasks(workflow *WorkflowIntegration) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, err := getUserInfo(c)
		if err != nil {
			response.Unauthorized(c, "User not authenticated")
			return
		}

		tasks, err := workflow.GetPlanPendingTasks(userID)
		if err != nil {
			response.InternalError(c, "Failed to get pending tasks")
			return
		}

		response.Success(c, tasks)
	}
}
