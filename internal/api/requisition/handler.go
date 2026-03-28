package requisition

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/audit"
	"github.com/yourorg/material-backend/backend/internal/api/notification"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"github.com/yourorg/material-backend/backend/internal/api/stock"
	"github.com/yourorg/material-backend/backend/internal/api/workflow"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("requisition")
	r.Use(jwtpkg.TokenMiddleware())

	// Get pending requisitions count
	r.GET("/requisitions/pending/count", auth.PermissionMiddleware(db, "requisition_view"), func(c *gin.Context) {
		var count int64
		query := db.Model(&Requisition{}).Where("status = ?", "pending")

		// 获取用户可访问的项目ID列表（数据权限过滤）
		projectIDs, err := auth.GetAccessibleProjectIDs(c, db)
		if err != nil {
			response.InternalError(c, "获取用户项目权限失败")
			return
		}

		// 应用项目过滤
		if projectIDs != nil {
			if len(projectIDs) == 0 {
				// 用户无任何项目权限，返回 0
				response.SuccessWithMeta(c, map[string]int64{"count": 0}, nil)
				return
			}
			query = query.Where("project_id IN ?", projectIDs)
		}

		query.Count(&count)
		response.SuccessWithMeta(c, map[string]int64{"count": count}, nil)
	})

	// list requisitions (supports filters, pagination, sorting)
	r.GET("/requisitions", auth.PermissionMiddleware(db, "requisition_view"), func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 { pageSize = 100 }
		status := c.Query("status")
		applicant := c.Query("applicant")
		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)
		requisitionNo := c.Query("requisition_no")

		// 支持多个项目ID（用于包含子项目）
		var projectIDsFilter []uint
		projectIDsStr := c.Query("project_ids")
		if projectIDsStr != "" {
			// 解析逗号分隔的项目ID列表
			for _, idStr := range strings.Split(projectIDsStr, ",") {
				if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64); err == nil {
					projectIDsFilter = append(projectIDsFilter, uint(id))
				}
			}
		}

		var total int64
		query := db.Model(&Requisition{})

		// field filters
		if status != "" { query = query.Where("status = ?", status) }
		if applicant != "" { query = query.Where("applicant LIKE ?", "%"+applicant+"%") }
		if requisitionNo != "" { query = query.Where("requisition_no LIKE ?", "%"+requisitionNo+"%") }

		// 获取用户可访问的项目ID列表（数据权限过滤）
		projectIDs, err := auth.GetAccessibleProjectIDs(c, db)
		if err != nil {
			response.InternalError(c, "获取用户项目权限失败")
			return
		}

		// 应用项目过滤
		if projectIDs != nil {
			if len(projectIDs) == 0 {
				// 用户无任何项目权限，返回空结果
				response.SuccessWithPagination(c, []map[string]any{}, int64(page), int64(pageSize), 0)
				return
			}

			// 如果指定了 project_ids（包含子项目），使用数组过滤
			if len(projectIDsFilter) > 0 {
				// 检查所有请求的项目ID是否在用户可访问列表中
				for _, pid := range projectIDsFilter {
					hasAccess := false
					for _, accessibleID := range projectIDs {
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
			} else if projectID > 0 {
				// 如果指定了 project_id，检查是否在用户可访问列表中
				hasAccess := false
				for _, pid := range projectIDs {
					if uint(projectID) == pid {
						hasAccess = true
						break
					}
				}
				if !hasAccess {
					response.Forbidden(c, "无权访问该项目")
					return
				}
				query = query.Where("project_id = ?", projectID)
			} else {
				// 如果没有指定 project_id 参数，自动过滤用户可访问的项目
				query = query.Where("project_id IN ?", projectIDs)
			}
		} else {
			// 管理员
			if len(projectIDsFilter) > 0 {
				query = query.Where("project_id IN ?", projectIDsFilter)
			} else if projectID > 0 {
				query = query.Where("project_id = ?", projectID)
			}
		}

		query.Count(&total)
		var requisitions []Requisition
		query.Offset((page-1)*pageSize).Limit(pageSize).Order("created_at DESC").Preload("Items").Find(&requisitions)

		// Enrich with project names and material details
		out := make([]map[string]any, 0, len(requisitions))
		for _, req := range requisitions {
			dto := req.ToDTOWithEnrichment(db)
			if req.ProjectID > 0 {
				var project struct {
					ID   uint
					Name string
				}
				if err := db.Model(&struct{}{}).Table("projects").Where("id = ?", req.ProjectID).
					Select("id, name").Scan(&project).Error; err == nil && project.ID > 0 {
					dto["project_name"] = project.Name
				} else {
					dto["project_name"] = "未知项目"
				}
			} else {
				dto["project_name"] = nil
			}
			out = append(out, dto)
		}
		response.SuccessWithPagination(c, out, int64(page), int64(pageSize), total)
	})

	// create requisition
	r.POST("/requisitions", auth.PermissionMiddleware(db, "requisition_create"), func(c *gin.Context) {
		var req struct {
			ProjectID  uint                 `json:"project_id" binding:"required"`
			PlanID     *uint                `json:"plan_id"`
			Applicant  string               `json:"applicant"`
			Department string               `json:"department"`
			Purpose    string               `json:"purpose"`
			Urgent     bool                 `json:"urgent"`
			Items      []RequisitionItemReq `json:"items" binding:"required,dive"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Printf("Create requisition validation error: %v\n", err)
			response.BadRequest(c, err.Error())
			return
		}
		if len(req.Items) == 0 {
			response.BadRequest(c, "申请物资不能为空")
			return
		}

		// Get current user from context
		username := ""
		if user, exists := c.Get("current_username"); exists {
			if name, ok := user.(string); ok {
				username = name
			}
		}
		// Fallback: try to get from token_payload
		if username == "" {
			if payload, exists := c.Get("token_payload"); exists {
				if claims, ok := payload.(map[string]any); ok {
					if name, ok := claims["username"].(string); ok {
						username = name
					}
				}
			}
		}

		if req.Applicant == "" {
			req.Applicant = username
		}

		// Generate requisition number
		requisitionNo := generateRequisitionNo(db)

		// Begin transaction
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Create requisition
		requisition := Requisition{
			RequisitionNo: requisitionNo,
			ProjectID:     req.ProjectID,
			PlanID:        req.PlanID,
			Applicant:     req.Applicant,
			Department:    req.Department,
			Purpose:       req.Purpose,
			Status:        "pending",
			Urgent:        boolToInt(req.Urgent),
		}

		if err := tx.Create(&requisition).Error; err != nil {
			tx.Rollback()
			response.InternalError(c, "创建申请单失败")
			return
		}

		// Create requisition items
		for _, itemReq := range req.Items {
			if itemReq.MaterialID == 0 || itemReq.RequestedQuantity <= 0 {
				continue // Skip invalid items
			}

			item := RequisitionItem{
				RequisitionID:     requisition.ID,
				StockID:           itemReq.StockID,
				MaterialID:        itemReq.MaterialID,
				PlanItemID:        itemReq.PlanItemID,
				RequestedQuantity: itemReq.RequestedQuantity,
				ApprovedQuantity:  0,
				ActualQuantity:    0,
				Remark:            itemReq.Remark,
				Status:            "pending",
			}

			if err := tx.Create(&item).Error; err != nil {
				tx.Rollback()
				response.InternalError(c, "创建申请单项失败")
				return
			}
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			response.InternalError(c, "提交申请单失败")
			return
		}

		// Load the created requisition with items
		var createdReq Requisition
		db.Preload("Items").First(&createdReq, requisition.ID)

		// 启动工作流
		wfIntegration := NewWorkflowIntegration(db)
		userID, _ := c.Get("current_user_id")
		var uid uint
		if userID != nil {
			if id, ok := userID.(int64); ok {
				uid = uint(id)
			}
		}

		// 尝试启动工作流（如果配置了工作流）
		if err := wfIntegration.StartRequisitionWorkflow(&createdReq, uid, createdReq.Applicant); err != nil {
			// 工作流启动失败不影响创建，记录日志即可
			fmt.Printf("启动领料单工作流失败: %v\n", err)
		}

		// Send notification to users with approve permission
		go func() {
			// 获取项目名称
			var projectName string
			if createdReq.ProjectID > 0 {
				var project struct {
					Name string
				}
				db.Model(&struct{}{}).Table("projects").Where("id = ?", createdReq.ProjectID).
					Select("name").Scan(&project)
				projectName = project.Name
			}

			// 发送通知给所有有审批权限的用户
			_ = notification.NotifyUsersWithPermission(db,
				"requisition_approve",
				notification.TypeRequisitionApprove,
				"新的出库单待审批",
				fmt.Sprintf("申请人：%s，项目：%s，单号：%s", createdReq.Applicant, projectName, createdReq.RequisitionNo),
				map[string]interface{}{
					"requisition_id":  createdReq.ID,
					"requisition_no":  createdReq.RequisitionNo,
					"project_id":      createdReq.ProjectID,
					"project_name":    projectName,
					"applicant":       createdReq.Applicant,
					"items_count":     len(createdReq.Items),
				},
			)
		}()

		dto := createdReq.ToDTOWithEnrichment(db)

		// Enrich with project name
		if createdReq.ProjectID > 0 {
			var project struct {
				ID   uint
				Name string
			}
			if err := db.Model(&struct{}{}).Table("projects").Where("id = ?", createdReq.ProjectID).
				Select("id, name").Scan(&project).Error; err == nil && project.ID > 0 {
				dto["project_name"] = project.Name
			} else {
				dto["project_name"] = "未知项目"
			}
		} else {
			dto["project_name"] = nil
		}

		// 记录操作日志
		audit.LogCreate(&uid, username, audit.ModuleRequisition, audit.ResourceRequisition,
			createdReq.ID, createdReq.RequisitionNo, req)

		response.Created(c, dto, "申请单创建成功")
	})

	// get single requisition
	r.GET("/requisitions/:id", auth.PermissionMiddleware(db, "requisition_view"), func(c *gin.Context) {
		id := c.Param("id")
		var requisition Requisition
		if err := db.Preload("Items").First(&requisition, id).Error; err != nil {
			response.NotFound(c, "申请单不存在")
			return
		}

		dto := requisition.ToDTOWithEnrichment(db)

		// Enrich with project name
		if requisition.ProjectID > 0 {
			var project struct {
				ID   uint
				Name string
			}
			if err := db.Model(&struct{}{}).Table("projects").Where("id = ?", requisition.ProjectID).
				Select("id, name").Scan(&project).Error; err == nil && project.ID > 0 {
				dto["project_name"] = project.Name
			} else {
				dto["project_name"] = "未知项目"
			}
		} else {
			dto["project_name"] = nil
		}

		response.Success(c, dto)
	})

	// update requisition
	r.PUT("/requisitions/:id", auth.PermissionMiddleware(db, "requisition_edit"), func(c *gin.Context) {
		id := c.Param("id")
		var requisition Requisition
		if err := db.First(&requisition, id).Error; err != nil {
			response.NotFound(c, "申请单不存在")
			return
		}

		// Only pending requisitions can be updated
		if requisition.Status != "pending" {
			response.BadRequest(c, "只有待审核状态的申请单可以修改")
			return
		}

		var req struct {
			ProjectID  *uint `json:"project_id"`
			Department string `json:"department"`
			Purpose    string `json:"purpose"`
			Urgent     *bool  `json:"urgent"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// Update fields
		if req.ProjectID != nil {
			requisition.ProjectID = *req.ProjectID
		}
		if req.Department != "" {
			requisition.Department = req.Department
		}
		if req.Purpose != "" {
			requisition.Purpose = req.Purpose
		}
		if req.Urgent != nil {
			requisition.Urgent = boolToInt(*req.Urgent)
		}

		db.Save(&requisition)

		// Load with items
		var updatedReq Requisition
		db.Preload("Items").First(&updatedReq, requisition.ID)

		dto := updatedReq.ToDTOWithEnrichment(db)

		// Enrich with project name
		if updatedReq.ProjectID > 0 {
			var project struct {
				ID   uint
				Name string
			}
			if err := db.Model(&struct{}{}).Table("projects").Where("id = ?", updatedReq.ProjectID).
				Select("id, name").Scan(&project).Error; err == nil && project.ID > 0 {
				dto["project_name"] = project.Name
			} else {
				dto["project_name"] = "未知项目"
			}
		} else {
			dto["project_name"] = nil
		}

		response.SuccessWithMessage(c, dto, "申请单更新成功")
	})

	// approve requisition
	r.POST("/requisitions/:id/approve", auth.PermissionMiddleware(db, "requisition_approve"), func(c *gin.Context) {
		id := c.Param("id")
		var requisition Requisition
		if err := db.Preload("Items").First(&requisition, id).Error; err != nil {
			response.NotFound(c, "申请单不存在")
			return
		}

		// 检查是否有工作流实例
		wfIntegration := NewWorkflowIntegration(db)
		instance, err := wfIntegration.GetRequisitionWorkflowStatus(requisition.ID)

		// Parse request body for approved quantities
		var req struct {
			Items []struct {
				ID               uint    `json:"id"`
				ApprovedQuantity float64 `json:"approved_quantity"`
			} `json:"items"`
			Notes  string `json:"notes"`
			Remark string `json:"remark"` // 兼容入库单风格
		}
		c.ShouldBindJSON(&req)

		// 兼容 notes 和 remark 参数
		notes := req.Remark
		if notes == "" {
			notes = req.Notes
		}

		// 获取用户信息
		userID, _ := c.Get("current_user_id")
		var uid uint
		if userID != nil {
			if id, ok := userID.(int64); ok {
				uid = uint(id)
			}
		}
		userName, _ := c.Get("current_username")
		var name string
		if userName != nil {
			name = userName.(string)
		} else {
			name = "未知用户"
		}

		// 如果存在工作流实例，使用工作流引擎处理
		if err == nil && instance != nil && instance.Status == workflow.InstanceStatusPending {
			// 转换items格式
			items := make([]RequisitionApprovalItem, 0)
			if len(req.Items) > 0 {
				for _, item := range req.Items {
					items = append(items, RequisitionApprovalItem{
						ID:               item.ID,
						ApprovedQuantity: int(item.ApprovedQuantity),
					})
				}
			} else {
				// 如果没有指定审批数量，默认使用申请数量
				for _, reqItem := range requisition.Items {
					items = append(items, RequisitionApprovalItem{
						ID:               reqItem.ID,
						ApprovedQuantity: int(reqItem.RequestedQuantity),
					})
				}
			}

			// 使用工作流引擎处理审批
			// 注意：ProcessRequisitionApproval 内部会在工作流完全通过时自动执行发放操作
			wfErr := wfIntegration.ProcessRequisitionApproval(requisition.ID, uid, name, workflow.ActionApprove, notes, items)
			if wfErr == nil {
				// 工作流审批成功，重新加载出库单
				db.Preload("Items").First(&requisition, id)

				// 记录操作日志
				audit.LogApprove(&uid, name, audit.ModuleRequisition, audit.ResourceRequisition,
					requisition.ID, requisition.RequisitionNo, notes)

				response.SuccessWithMessage(c, requisition.ToDTO(), "审批通过并已自动发放")
				return
			}
			// 工作流审批失败（如无待办任务），回退到简单审批逻辑
			fmt.Printf("工作流审批失败，回退到简单审批: %v\n", wfErr)
		}

		// 如果没有工作流，使用原有逻辑（向后兼容）
		// Only pending requisitions can be approved
		if requisition.Status != "pending" {
			response.BadRequest(c, "只有待审核状态的申请单可以审核")
			return
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			// If request body is empty, approve all quantities as requested
			for i := range requisition.Items {
				req.Items = append(req.Items, struct {
					ID               uint    `json:"id"`
					ApprovedQuantity float64 `json:"approved_quantity"`
				}{
					ID:               requisition.Items[i].ID,
					ApprovedQuantity: requisition.Items[i].RequestedQuantity,
				})
			}
		}

		// 转换为审批项格式
		items := make([]RequisitionApprovalItem, 0)
		for _, item := range req.Items {
			items = append(items, RequisitionApprovalItem{
				ID:               item.ID,
				ApprovedQuantity: int(item.ApprovedQuantity),
			})
		}

		// 自动执行发放操作（审核通过后直接发放）
		if err := executeRequisitionIssue(db, &requisition, uid, name, items); err != nil {
			response.InternalError(c, fmt.Sprintf("自动发放失败: %v", err))
			return
		}

		// 重新加载出库单以获取最新状态
		db.Preload("Items").First(&requisition, id)
		dto := requisition.ToDTO()

		// Enrich with project name
		if requisition.ProjectID > 0 {
			var project struct {
				ID   uint
				Name string
			}
			if err := db.Model(&struct{}{}).Table("projects").Where("id = ?", requisition.ProjectID).
				Select("id, name").Scan(&project).Error; err == nil && project.ID > 0 {
				dto["project_name"] = project.Name
			} else {
				dto["project_name"] = "未知项目"
			}
		} else {
			dto["project_name"] = nil
		}

		// 记录操作日志
		audit.LogApprove(&uid, name, audit.ModuleRequisition, audit.ResourceRequisition,
			requisition.ID, requisition.RequisitionNo, notes)

		response.SuccessWithMessage(c, dto, "申请单审核通过并已自动发放")
	})

	// issue requisition
	r.POST("/requisitions/:id/issue", auth.PermissionMiddleware(db, "requisition_issue"), func(c *gin.Context) {
		id := c.Param("id")
		var requisition Requisition
		if err := db.Preload("Items").First(&requisition, id).Error; err != nil {
			response.NotFound(c, "申请单不存在")
			return
		}

		// Only approved requisitions can be issued
		if requisition.Status != "approved" {
			response.BadRequest(c, "只有已审核状态的申请单可以发放")
			return
		}

		// Parse request body for actual quantities
		var req struct {
			Items []struct {
				ID             uint    `json:"id"`
				ActualQuantity float64 `json:"actual_quantity"`
			} `json:"items"`
			Notes  string `json:"notes"`
			Remark string `json:"remark"` // 兼容入库单风格
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "无效的请求参数")
			return
		}

		// 兼容 notes 和 remark 参数
		notes := req.Notes
		if notes == "" {
			notes = req.Remark
		}

		// Get current user from context
		username := ""
		if user, exists := c.Get("current_username"); exists {
			if name, ok := user.(string); ok {
				username = name
			}
		}
		// Fallback: try to get from token_payload
		if username == "" {
			if payload, exists := c.Get("token_payload"); exists {
				if claims, ok := payload.(map[string]any); ok {
					if name, ok := claims["username"].(string); ok {
						username = name
					}
				}
			}
		}

		// Get user ID for logging
		userID, _ := c.Get("current_user_id")
		var uid uint
		if userID != nil {
			if id, ok := userID.(int64); ok {
				uid = uint(id)
			}
		}

		// Build a map of item ID to actual quantity
		actualQtyMap := make(map[uint]float64)
		for _, item := range req.Items {
			actualQtyMap[item.ID] = item.ActualQuantity
		}

		// Begin transaction
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Update requisition status
		now := time.Now()
		requisition.Status = "issued"
		requisition.IssuedAt = &now
		requisition.IssuedBy = username
		if notes != "" {
			requisition.Remark = notes
		}

		if err := tx.Save(&requisition).Error; err != nil {
			tx.Rollback()
			response.InternalError(c, "发放申请单失败")
			return
		}

		// Update items status and reduce stock
		for i := range requisition.Items {
			// Get actual quantity from request or use approved quantity
			actualQty := actualQtyMap[requisition.Items[i].ID]
			if actualQty == 0 {
				actualQty = requisition.Items[i].ApprovedQuantity
			}

			// Find stock record
			var stockRecord stock.Stock
			var err error

			// Try to find by StockID first
			if requisition.Items[i].StockID > 0 {
				err = tx.Where("id = ?", requisition.Items[i].StockID).First(&stockRecord).Error
			} else {
				// If no StockID, find by material_id and project
				err = tx.Where("material_id = ? AND project_id = ?",
					requisition.Items[i].MaterialID,
					requisition.ProjectID).
					First(&stockRecord).Error
			}

			if err != nil {
				tx.Rollback()
				response.InternalError(c, fmt.Sprintf("找不到材料ID %d 的库存记录", requisition.Items[i].MaterialID))
				return
			}

			// Check if stock is sufficient
			if stockRecord.Quantity < actualQty {
				tx.Rollback()
				response.BadRequest(c, fmt.Sprintf("材料ID %d 库存不足，当前库存：%.2f，需要：%.2f",
					requisition.Items[i].MaterialID, stockRecord.Quantity, actualQty))
				return
			}

			// 记录操作前后的数量
			quantityBefore := stockRecord.Quantity
			stockRecord.Quantity -= actualQty
			quantityAfter := stockRecord.Quantity

			// Reduce stock quantity
			if err := tx.Save(&stockRecord).Error; err != nil {
				tx.Rollback()
				response.InternalError(c, "更新库存失败")
				return
			}

			// Update item status and set actual quantity
			requisition.Items[i].Status = "issued"
			// Store actual quantity in ActualQuantity field
			requisition.Items[i].ActualQuantity = actualQty
			if err := tx.Save(&requisition.Items[i]).Error; err != nil {
				tx.Rollback()
				response.InternalError(c, "更新申请单项状态失败")
				return
			}

			// Log stock operation
			detail := fmt.Sprintf("出库单发放：%s，出库 %.2f", requisition.RequisitionNo, actualQty)
			if notes != "" {
				detail += fmt.Sprintf("，备注：%s", notes)
			}

			// Use requisition's project ID
			projectID := requisition.ProjectID

			// Create StockLog entry with full details
			stockLog := stock.StockLog{
				StockID:        stockRecord.ID,
				Type:           "out",
				Quantity:       actualQty,
				QuantityBefore: quantityBefore,
				QuantityAfter:  quantityAfter,
				SourceType:     "requisition",
				SourceID:       &requisition.ID,
				SourceNo:       requisition.RequisitionNo,
				ProjectID:      projectID,
				MaterialID:     requisition.Items[i].MaterialID,
				UserID:         &uid,
				Remark:         detail,
				CreatedAt:      time.Now(),
			}
			if err := tx.Create(&stockLog).Error; err != nil {
				tx.Rollback()
				response.InternalError(c, "创建库存日志失败")
				return
			}

			// Create StockOpLog entry
			opLog := stock.StockOpLog{
				StockID: stockRecord.ID,
				OpType:  "out",
				LogID:   requisition.ID, // 关联到领料单ID
				Detail:  detail,
				UserID:  uid,
				Time:    time.Now(),
			}
			if err := tx.Create(&opLog).Error; err != nil {
				tx.Rollback()
				response.InternalError(c, "创建库存操作日志失败")
				return
			}
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			response.InternalError(c, "提交发放失败")
			return
		}

		// Update material plan issued quantities if linked to a plan
		if requisition.PlanID != nil {
			for _, item := range requisition.Items {
				// Get actual quantity
				actualQty := item.ActualQuantity
				if actualQty == 0 {
					actualQty = item.ApprovedQuantity
				}

				// Update issued_quantity in material_plan_items
				result := db.Table("material_plan_items").
					Where("plan_id = ? AND material_id = ?", *requisition.PlanID, item.MaterialID).
					Update("issued_quantity", gorm.Expr("issued_quantity + ?", actualQty))

				if result.Error != nil {
					// Log error but don't fail the response
					fmt.Printf("Warning: Failed to update issued_quantity for plan_id=%d, material_id=%d: %v\n",
						*requisition.PlanID, item.MaterialID, result.Error)
				}
			}
		}

		// Reload the requisition with updated items
		db.Preload("Items").First(&requisition, id)
		dto := requisition.ToDTO()

		// Enrich with project name
		if requisition.ProjectID > 0 {
			var project struct {
				ID   uint
				Name string
			}
			if err := db.Model(&struct{}{}).Table("projects").Where("id = ?", requisition.ProjectID).
				Select("id, name").Scan(&project).Error; err == nil && project.ID > 0 {
				dto["project_name"] = project.Name
			} else {
				dto["project_name"] = "未知项目"
			}
		} else {
			dto["project_name"] = nil
		}

		// 记录操作日志
		audit.LogCreate(&uid, username, audit.ModuleRequisition, "RequisitionIssue",
			requisition.ID, requisition.RequisitionNo, req)

		response.SuccessWithMessage(c, dto, "申请单发放成功")
	})

	// reject requisition
	r.POST("/requisitions/:id/reject", auth.PermissionMiddleware(db, "requisition_approve"), func(c *gin.Context) {
		id := c.Param("id")
		var requisition Requisition
		if err := db.Preload("Items").First(&requisition, id).Error; err != nil {
			response.NotFound(c, "申请单不存在")
			return
		}

		// 检查是否有工作流实例
		wfIntegration := NewWorkflowIntegration(db)
		instance, err := wfIntegration.GetRequisitionWorkflowStatus(requisition.ID)

		var req struct {
			Remark string `json:"remark"`
			Reason string `json:"reason"` // 兼容旧版本
		}
		c.ShouldBindJSON(&req)

		// 兼容 reason 和 remark 参数
		remark := req.Remark
		if remark == "" {
			remark = req.Reason
		}

		// 获取用户信息
		userID, _ := c.Get("current_user_id")
		var uid uint
		if userID != nil {
			if id, ok := userID.(int64); ok {
				uid = uint(id)
			}
		}
		userName, _ := c.Get("current_username")
		var name string
		if userName != nil {
			name = userName.(string)
		} else {
			// Fallback: try to get from token_payload
			if payload, exists := c.Get("token_payload"); exists {
				if claims, ok := payload.(map[string]any); ok {
					if username, ok := claims["username"].(string); ok {
						name = username
					}
				}
			}
		}

		// 如果存在工作流实例，使用工作流引擎处理
		if err == nil && instance != nil && instance.Status == workflow.InstanceStatusPending {
			// 使用工作流引擎处理拒绝
			wfErr := wfIntegration.ProcessRequisitionApproval(requisition.ID, uid, name, workflow.ActionReject, remark, nil)
			if wfErr == nil {
				// 工作流拒绝成功
				db.Preload("Items").First(&requisition, id)

				// 记录操作日志
				audit.LogReject(&uid, name, audit.ModuleRequisition, audit.ResourceRequisition,
					requisition.ID, requisition.RequisitionNo, remark)

				response.SuccessWithMessage(c, requisition.ToDTO(), "已拒绝")
				return
			}
			// 工作流拒绝失败（如无待办任务），回退到简单拒绝逻辑
			fmt.Printf("工作流拒绝失败，回退到简单拒绝: %v\n", wfErr)
		}

		// 如果没有工作流，使用原有逻辑（向后兼容）
		if requisition.Status != "pending" {
			response.BadRequest(c, "只有待审核状态的申请单可以拒绝")
			return
		}

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		now := time.Now()
		requisition.Status = "rejected"
		requisition.ApprovedAt = &now
		requisition.ApprovedBy = name
		if remark != "" {
			requisition.Remark = remark
		}

		if err := tx.Save(&requisition).Error; err != nil {
			tx.Rollback()
			response.InternalError(c, "拒绝申请单失败")
			return
		}

		for i := range requisition.Items {
			requisition.Items[i].Status = "rejected"
			if err := tx.Save(&requisition.Items[i]).Error; err != nil {
				tx.Rollback()
				response.InternalError(c, "更新申请单项状态失败")
				return
			}
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			response.InternalError(c, "提交拒绝失败")
			return
		}

		dto := requisition.ToDTO()

		// Enrich with project name
		if requisition.ProjectID > 0 {
			var project struct {
				ID   uint
				Name string
			}
			if err := db.Model(&struct{}{}).Table("projects").Where("id = ?", requisition.ProjectID).
				Select("id, name").Scan(&project).Error; err == nil && project.ID > 0 {
				dto["project_name"] = project.Name
			} else {
				dto["project_name"] = "未知项目"
			}
		} else {
			dto["project_name"] = nil
		}

		// 记录操作日志
		audit.LogReject(&uid, name, audit.ModuleRequisition, audit.ResourceRequisition,
			requisition.ID, requisition.RequisitionNo, remark)

		response.SuccessWithMessage(c, dto, "申请单已拒绝")
	})

	// Resubmit requisition
	r.POST("/requisitions/:id/resubmit", auth.PermissionMiddleware(db, "requisition_create"), func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
		var requisition Requisition
		if err := db.First(&requisition, id).Error; err != nil {
			response.NotFound(c, "申请单不存在")
			return
		}

		// Get user info from JWT middleware
		userIDInt64, _ := c.Get("current_user_id")
		var submitterID uint
		if userIDInt64 != nil {
			if id, ok := userIDInt64.(int64); ok {
				submitterID = uint(id)
			}
		}
		username, _ := c.Get("current_username")
		var submitterName string
		if username != nil {
			if name, ok := username.(string); ok {
				submitterName = name
			}
		}

		// Use workflow integration to resubmit
		wfIntegration := NewWorkflowIntegration(db)
		if err := wfIntegration.ResubmitRequisition(requisition.ID, submitterID, submitterName); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// Reload requisition
		db.Preload("Items").First(&requisition, id)
		dto := requisition.ToDTO()

		// Enrich with project name
		if requisition.ProjectID > 0 {
			var project struct {
				ID   uint
				Name string
			}
			if err := db.Model(&struct{}{}).Table("projects").Where("id = ?", requisition.ProjectID).
				Select("id, name").Scan(&project).Error; err == nil && project.ID > 0 {
				dto["project_name"] = project.Name
			} else {
				dto["project_name"] = "未知项目"
			}
		} else {
			dto["project_name"] = nil
		}

		// 记录操作日志
		if err := audit.LogCreate(&submitterID, submitterName, audit.ModuleRequisition, audit.ResourceRequisition,
			requisition.ID, requisition.RequisitionNo, map[string]any{"action": "resubmit"}); err != nil {
			fmt.Printf("记录操作日志失败: %v\n", err)
		}

		response.SuccessWithMessage(c, dto, "申请单已重新提交")
	})

	// get requisition items
	r.GET("/requisition-items", auth.PermissionMiddleware(db, "requisition_view"), func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}
		reqID, _ := strconv.ParseUint(c.Query("requisition_id"), 10, 64)
		status := c.Query("status")

		var total int64
		query := db.Model(&RequisitionItem{})

		if reqID > 0 {
			query = query.Where("requisition_id = ?", reqID)
		}
		if status != "" {
			query = query.Where("status = ?", status)
		}

		query.Count(&total)
		var items []RequisitionItem
		query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)

		out := make([]map[string]any, 0, len(items))
		for _, item := range items {
			out = append(out, item.ToDTO())
		}

		response.SuccessWithPagination(c, out, int64(page), int64(pageSize), total)
	})

	// delete requisition
	r.DELETE("/requisitions/:id", auth.PermissionMiddleware(db, "requisition_delete"), func(c *gin.Context) {
		id := c.Param("id")
		var requisition Requisition
		if err := db.First(&requisition, id).Error; err != nil {
			response.NotFound(c, "申请单不存在")
			return
		}

		if requisition.Status != "pending" {
			response.BadRequest(c, "只有待审核状态的申请单可以删除")
			return
		}

		if err := db.Delete(&requisition).Error; err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, nil, "申请单删除成功")
	})
}

// RequisitionItemReq is the request structure for requisition items
type RequisitionItemReq struct {
	StockID           uint    `json:"stock_id"`
	MaterialID        uint    `json:"material_id" binding:"required_without=StockID"`
	PlanItemID        *uint   `json:"plan_item_id"`
	RequestedQuantity float64 `json:"requested_quantity" binding:"required,gt=0"`
	Remark            string  `json:"remark"`
}

// generateRequisitionNo generates a unique requisition number
func generateRequisitionNo(db *gorm.DB) string {
	// Format: CK + YYYYMMDD + 3-digit sequence
	today := time.Now().Format("20060102")
	prefix := fmt.Sprintf("CK%s", today)

	// Find the latest requisition number with today's prefix
	var latestReq Requisition
	result := db.Where("requisition_no LIKE ?", prefix+"%").Order("requisition_no DESC").First(&latestReq)

	sequence := 1
	if result.RowsAffected > 0 {
		// Extract sequence number from the latest requisition number
		seqStr := latestReq.RequisitionNo[len(prefix):]
		if seq, err := strconv.Atoi(seqStr); err == nil {
			sequence = seq + 1
		}
	}

	// Format sequence with leading zeros
	return fmt.Sprintf("%s%03d", prefix, sequence)
}

// boolToInt converts boolean to integer (0 or 1)
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// executeRequisitionIssue 执行出库单的实际发放逻辑（扣减库存、创建日志）
// 审批通过后自动调用此函数完成发放
func executeRequisitionIssue(db *gorm.DB, requisition *Requisition, approverID uint, approverName string, approverItems []RequisitionApprovalItem) error {
	// 构建批准数量映射
	approvedQtyMap := make(map[uint]float64)
	if len(approverItems) > 0 {
		for _, item := range approverItems {
			approvedQtyMap[item.ID] = float64(item.ApprovedQuantity)
		}
	} else {
		// 默认使用申请数量
		for _, reqItem := range requisition.Items {
			approvedQtyMap[reqItem.ID] = reqItem.RequestedQuantity
		}
	}

	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update requisition status to issued
	now := time.Now()
	requisition.Status = "issued"
	requisition.ApprovedAt = &now
	requisition.ApprovedBy = approverName
	requisition.IssuedAt = &now
	requisition.IssuedBy = approverName

	if err := tx.Save(&requisition).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新出库单状态失败: %w", err)
	}

	// Update items status and reduce stock
	for _, item := range requisition.Items {
		// Get approved quantity
		approvedQty := approvedQtyMap[item.ID]
		if approvedQty == 0 {
			approvedQty = item.RequestedQuantity
		}

		// Find stock record
		var stockRecord stock.Stock
		var err error

		// Try to find by StockID first
		if item.StockID > 0 {
			err = tx.Where("id = ?", item.StockID).First(&stockRecord).Error
		} else {
			// If no StockID, find by material_id and project
			err = tx.Where("material_id = ? AND project_id = ?",
				item.MaterialID,
				requisition.ProjectID).
				First(&stockRecord).Error
		}

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("找不到材料ID %d 的库存记录: %w", item.MaterialID, err)
		}

		// Check if stock is sufficient
		if stockRecord.Quantity < approvedQty {
			tx.Rollback()
			return fmt.Errorf("材料ID %d 库存不足，当前库存：%.2f，需要：%.2f",
				item.MaterialID, stockRecord.Quantity, approvedQty)
		}

		// 记录操作前后的数量
		quantityBefore := stockRecord.Quantity
		stockRecord.Quantity -= approvedQty
		quantityAfter := stockRecord.Quantity

		// Reduce stock quantity
		if err := tx.Save(&stockRecord).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("更新库存失败: %w", err)
		}

		// Update item status
		item.Status = "issued"
		item.ApprovedQuantity = approvedQty
		item.ActualQuantity = approvedQty
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("更新申请单项状态失败: %w", err)
		}

		// Log stock operation
		detail := fmt.Sprintf("出库单发放：%s，出库 %.2f", requisition.RequisitionNo, approvedQty)

		// Use requisition's project ID
		projectID := requisition.ProjectID

		// Create StockLog entry with full details
		stockLog := stock.StockLog{
			StockID:        stockRecord.ID,
			Type:           "out",
			Quantity:       approvedQty,
			QuantityBefore: quantityBefore,
			QuantityAfter:  quantityAfter,
			SourceType:     "requisition",
			SourceID:       &requisition.ID,
			SourceNo:       requisition.RequisitionNo,
			ProjectID:      projectID,
			MaterialID:     item.MaterialID,
			UserID:         &approverID,
			Remark:         detail,
			CreatedAt:      time.Now(),
		}
		if err := tx.Create(&stockLog).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建库存日志失败: %w", err)
		}

		// Create StockOpLog entry
		opLog := stock.StockOpLog{
			StockID: stockRecord.ID,
			OpType:  "out",
			LogID:   requisition.ID,
			Detail:  detail,
			UserID:  approverID,
			Time:    time.Now(),
		}
		if err := tx.Create(&opLog).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建库存操作日志失败: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交发放失败: %w", err)
	}

	// Update material plan issued quantities if linked to a plan
	if requisition.PlanID != nil {
		for _, item := range requisition.Items {
			// Get actual quantity
			actualQty := item.ActualQuantity
			if actualQty == 0 {
				actualQty = item.ApprovedQuantity
			}

			// Update issued_quantity in material_plan_items
			result := db.Table("material_plan_items").
				Where("plan_id = ? AND material_id = ?", *requisition.PlanID, item.MaterialID).
				Update("issued_quantity", gorm.Expr("issued_quantity + ?", actualQty))

			if result.Error != nil {
				// Log error but don't fail
				fmt.Printf("Warning: Failed to update issued_quantity for plan_id=%d, material_id=%d: %v\n",
					*requisition.PlanID, item.MaterialID, result.Error)
			}
		}
	}

	return nil
}