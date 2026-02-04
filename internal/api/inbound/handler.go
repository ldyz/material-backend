package inbound

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/audit"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"github.com/yourorg/material-backend/backend/internal/api/workflow"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	g := rg.Group("/inbound")
	// 使用JWT中间件进行身份验证
	g.Use(jwtpkg.TokenMiddleware())

	// Get pending inbound orders count
	g.GET("/inbound-orders/pending/count", auth.PermissionMiddleware(db, "inbound_view"), func(c *gin.Context) {
		var count int64
		query := db.Model(&InboundOrder{}).Where("status = ?", StatusPending)

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

	// List inbound orders
	g.GET("/inbound-orders", auth.PermissionMiddleware(db, "inbound_view"), func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		// Filter parameters
		status := strings.TrimSpace(c.Query("status"))
		creatorName := strings.TrimSpace(c.Query("creator_name"))
		supplier := strings.TrimSpace(c.Query("supplier"))
		orderNo := strings.TrimSpace(c.Query("order_no"))
		search := strings.TrimSpace(c.Query("search"))
		startDate := strings.TrimSpace(c.Query("start_date"))
		endDate := strings.TrimSpace(c.Query("end_date"))
		projectID := strings.TrimSpace(c.Query("project_id"))

		query := db.Preload("Items")

		// Status filter
		if status != "" {
			query = query.Where("status = ?", status)
		}

		// Creator filter
		if creatorName != "" {
			query = query.Where("creator_name LIKE ?", "%"+creatorName+"%")
		}

		// Supplier filter
		if supplier != "" {
			query = query.Where("supplier LIKE ?", "%"+supplier+"%")
		}

		// Order number filter
		if orderNo != "" {
			query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
		}

		// General search
		if search != "" {
			query = query.Where("order_no LIKE ? OR supplier LIKE ? OR creator_name LIKE ? OR contact LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
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
			// 解析逗号分隔的项目ID列表
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
				// 用户无任何项目权限，返回空结果
				response.SuccessWithPagination(c, []map[string]any{}, int64(page), int64(pageSize), 0)
				return
			}

			// 如果指定了 project_ids（包含子项目），使用数组过滤
			if len(projectIDsFilter) > 0 {
				// 检查所有请求的项目ID是否在用户可访问列表中
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
				// 如果指定了 project_id，检查是否在用户可访问列表中
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
				// 如果没有指定 project_id 参数，自动过滤用户可访问的项目
				query = query.Where("project_id IN ?", userProjectIDs)
			}
		} else {
			// 管理员
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
		query.Model(&InboundOrder{}).Count(&total)

		// Sort by creation date descending
		query = query.Order("created_at DESC")

		// Pagination
		var orders []InboundOrder
		query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)

		// Enrich with project names and material info
		enrichedOrders := make([]map[string]any, len(orders))
		for i, order := range orders {
			dto := order.ToDTOWithEnrichment(db)
			if order.ProjectID > 0 {
				var project struct {
					ID   uint
					Name string
				}
				if err := db.Model(&struct{}{}).Table("projects").Where("id = ?", order.ProjectID).
					Select("id, name").Scan(&project).Error; err == nil && project.ID > 0 {
					dto["project_name"] = project.Name
				} else {
					dto["project_name"] = "未知项目"
				}
			} else {
				dto["project_name"] = nil
			}
			enrichedOrders[i] = dto
		}

		response.SuccessWithPagination(c, enrichedOrders, int64(page), int64(pageSize), total)
	})

	// Create inbound order
	g.POST("/inbound-orders", auth.PermissionMiddleware(db, "inbound_create"), func(c *gin.Context) {
		// Get user info from JWT middleware
		userIDInt64, _ := c.Get("current_user_id")
		username, _ := c.Get("current_username")

		var creatorID uint
		if userIDInt64 != nil {
			if id, ok := userIDInt64.(int64); ok {
				creatorID = uint(id)
			}
		}

		var creatorName string
		if username != nil {
			if name, ok := username.(string); ok {
				creatorName = name
			}
		}

		var req struct {
			Supplier    string `json:"supplier"`
			Contact     string `json:"contact"`
			ProjectID   string `json:"project_id"`
			PlanID      *uint   `json:"plan_id"` // 关联物资计划ID
			Status      string `json:"status"`
			Notes       string `json:"notes"`
			Remark      string `json:"remark"`
			TotalAmount float64 `json:"total_amount"`
			Items       []struct {
				MaterialID uint    `json:"material_id"`
				Quantity   float64 `json:"quantity"`
				UnitPrice  float64 `json:"unit_price"`
				Remark     string  `json:"remark"`
			} `json:"items"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误")
			return
		}

		// Parse and validate project ID
		var projectIDUint uint
		if req.ProjectID != "" {
			pid, err := strconv.ParseUint(req.ProjectID, 10, 64)
			if err != nil || pid == 0 {
				response.BadRequest(c, "项目ID格式无效")
				return
			}
			projectIDUint = uint(pid)

			// Verify project exists
			var count int64
			if err := db.Table("projects").Where("id = ?", projectIDUint).Count(&count).Error; err != nil || count == 0 {
				response.BadRequest(c, "指定的项目不存在")
				return
			}
		} else {
			response.BadRequest(c, "请选择所属项目")
			return
		}

		// Validate plan if provided
		if req.PlanID != nil {
			// Import material_plan package to check plan status
			var plan struct {
				ID       uint
				Status   string
				ProjectID uint
			}
			if err := db.Table("material_plans").Where("id = ?", req.PlanID).
				Select("id, status, project_id").First(&plan).Error; err != nil {
				response.BadRequest(c, "物资计划不存在")
				return
			}

			// Only approved or active plans can have inbound orders
			if plan.Status != "approved" && plan.Status != "active" {
				response.BadRequest(c, "只有审批通过或激活的计划才能创建入库单")
				return
			}

			// Verify plan belongs to the same project
			if plan.ProjectID != projectIDUint {
				response.BadRequest(c, "计划所属项目与入库单项目不一致")
				return
			}
		}

		if len(req.Items) == 0 {
			response.BadRequest(c, "入库单物资列表不能为空")
			return
		}

		// Validate items
		for _, item := range req.Items {
			if item.MaterialID == 0 {
				response.BadRequest(c, "物资ID无效")
				return
			}
			if item.Quantity <= 0 {
				response.BadRequest(c, "入库数量必须大于0")
				return
			}

			// 检查物资是否存在
			var material struct {
				ID uint
			}
			if err := db.Table("material_master").Where("id = ?", item.MaterialID).
				Select("id").First(&material).Error; err != nil {
				response.BadRequest(c, fmt.Sprintf("物资ID %d 不存在", item.MaterialID))
				return
			}
		}

		// Auto-generate order number
		orderNo := fmt.Sprintf("RK%s", time.Now().Format("20060102150405"))

		order := InboundOrder{
			OrderNo:     orderNo,
			Supplier:    strings.TrimSpace(req.Supplier),
			Contact:     strings.TrimSpace(req.Contact),
			ProjectID:   projectIDUint,
			PlanID:      req.PlanID, // 保存关联的计划ID
			CreatorID:   creatorID,
			CreatorName: creatorName,
			Status:      StatusPending,
			Notes:       strings.TrimSpace(req.Notes),
			Remark:      strings.TrimSpace(req.Remark),
			TotalAmount: int(req.TotalAmount * 100),
		}

		// Create items
		totalAmount := 0
		for _, item := range req.Items {
			orderItem := InboundOrderItem{
				StockID:    nil, // Will be set when approved
				MaterialID:  item.MaterialID,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				Remark:      item.Remark,
			}
			order.Items = append(order.Items, orderItem)
			totalAmount += int(item.Quantity * item.UnitPrice * 100)
		}
		order.TotalAmount = totalAmount

		// 使用事务确保工作流启动成功才创建入库单
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// 创建入库单
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			response.InternalError(c, err.Error())
			return
		}

		// 启动工作流
		wfIntegration := NewWorkflowIntegration(db)
		if err := wfIntegration.StartInboundWorkflow(&order, creatorID, creatorName); err != nil {
			tx.Rollback()
			response.InternalError(c, fmt.Sprintf("启动工作流失败: %v", err))
			return
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			response.InternalError(c, "保存入库单失败")
			return
		}

		// Reload with items for enrichment
		db.Preload("Items").First(&order, order.ID)

		// 记录操作日志
		if err := audit.LogCreate(&creatorID, creatorName, audit.ModuleInbound, audit.ResourceInboundOrder,
			order.ID, order.OrderNo, req); err != nil {
			// 记录日志失败不应影响业务流程，只记录错误
			fmt.Printf("记录操作日志失败: %v\n", err)
		}

		response.Created(c, order.ToDTOWithEnrichment(db), "入库单创建成功")
	})

	// Get single inbound order
	g.GET("/inbound-orders/:id", auth.PermissionMiddleware(db, "inbound_view"), func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
		var order InboundOrder
		if err := db.Preload("Items").First(&order, id).Error; err != nil {
			response.NotFound(c, "入库单不存在")
			return
		}

		dto := order.ToDTOWithEnrichment(db)

		// Enrich with project name
		if order.ProjectID > 0 {
			var project struct {
				ID   uint
				Name string
			}
			if err := db.Model(&struct{}{}).Table("projects").Where("id = ?", order.ProjectID).
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

	// Update inbound order
	g.PUT("/inbound-orders/:id", auth.PermissionMiddleware(db, "inbound_edit"), func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
		var order InboundOrder
		if err := db.Preload("Items").First(&order, id).Error; err != nil {
			response.NotFound(c, "入库单不存在")
			return
		}

		// Only pending status can be edited
		if order.Status != StatusPending {
			response.BadRequest(c, "只有待处理状态的入库单才能修改")
			return
		}

		var req struct {
			Supplier    string `json:"supplier"`
			Contact     string `json:"contact"`
			ProjectID   string `json:"project_id"`
			OrderNo     string `json:"order_no"`
			Notes       string `json:"notes"`
			Remark      string `json:"remark"`
			TotalAmount float64 `json:"total_amount"`
			Items       []struct {
				MaterialID uint    `json:"material_id"`
				Quantity   float64 `json:"quantity"`
				UnitPrice  float64 `json:"unit_price"`
				Remark     string  `json:"remark"`
			} `json:"items"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误")
			return
		}

		// Parse and validate project ID if provided
		if req.ProjectID != "" {
			pid, err := strconv.ParseUint(req.ProjectID, 10, 64)
			if err != nil || pid == 0 {
				response.BadRequest(c, "项目ID格式无效")
				return
			}

			// Verify project exists
			var count int64
			if err := db.Table("projects").Where("id = ?", pid).Count(&count).Error; err != nil || count == 0 {
				response.BadRequest(c, "指定的项目不存在")
				return
			}

			order.ProjectID = uint(pid)
		}

		// Update basic info
		if req.Supplier != "" {
			order.Supplier = strings.TrimSpace(req.Supplier)
		}
		if req.Contact != "" {
			order.Contact = strings.TrimSpace(req.Contact)
		}
		if req.OrderNo != "" {
			order.OrderNo = strings.TrimSpace(req.OrderNo)
		}
		if req.Notes != "" {
			order.Notes = strings.TrimSpace(req.Notes)
		}

		// Update items if provided
		if len(req.Items) > 0 {
			// Delete existing items
			db.Where("inbound_order_id = ?", order.ID).Delete(&InboundOrderItem{})

			// Add new items
			totalAmount := 0
			order.Items = make([]InboundOrderItem, 0)
			for _, item := range req.Items {
				if item.MaterialID == 0 || item.Quantity <= 0 {
					continue
				}
				orderItem := InboundOrderItem{
					InboundOrderID: order.ID,
					MaterialID:     item.MaterialID,
					Quantity:       item.Quantity,
					UnitPrice:      item.UnitPrice,
					Remark:         item.Remark,
				}
				order.Items = append(order.Items, orderItem)
				totalAmount += int(item.Quantity * item.UnitPrice * 100)
			}
			order.TotalAmount = totalAmount
		}

		order.UpdatedAt = time.Now()

		if err := db.Save(&order).Error; err != nil {
			response.InternalError(c, err.Error())
			return
		}

		// Reload with items for enrichment
		db.Preload("Items").First(&order, order.ID)
		response.SuccessWithMessage(c, order.ToDTOWithEnrichment(db), "入库单更新成功")
	})

	// Delete inbound order
	g.DELETE("/inbound-orders/:id", auth.PermissionMiddleware(db, "inbound_delete"), func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
		var order InboundOrder
		if err := db.First(&order, id).Error; err != nil {
			response.NotFound(c, "入库单不存在")
			return
		}

		// Only pending status can be deleted
		if order.Status != StatusPending {
			response.BadRequest(c, "只有待处理状态的入库单才能删除")
			return
		}

		if err := db.Delete(&order).Error; err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, nil, "入库单删除成功")
	})

	// Approve inbound order
	g.POST("/inbound-orders/:id/approve", auth.PermissionMiddleware(db, "inbound_approve"), func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
		var order InboundOrder
		if err := db.Preload("Items").First(&order, id).Error; err != nil {
			response.NotFound(c, "入库单不存在")
			return
		}

		// 检查是否有工作流实例
		wfIntegration := NewWorkflowIntegration(db)
		instance, err := wfIntegration.GetInboundWorkflowStatus(order.ID)

		var req struct {
			Items []struct {
				ID               uint    `json:"id"`
				ApprovedQuantity float64 `json:"approved_quantity"`
			} `json:"items"`
			Remark string `json:"remark"`
			Notes  string `json:"notes"`
		}

		c.ShouldBindJSON(&req)

		// Combine remark and notes
		remark := req.Remark
		if remark == "" {
			remark = req.Notes
		}

		// 如果存在工作流实例，使用工作流引擎处理
		if err == nil && instance != nil && instance.Status == workflow.InstanceStatusPending {
			// 获取用户信息
			userID, _ := c.Get("current_user_id")
			var uid uint
			if userID != nil {
				if id, ok := userID.(int64); ok {
					uid = uint(id)
				}
			}
			userName, _ := c.Get("username")
			var name string
			if userName != nil {
				name = userName.(string)
			} else {
				name = "未知用户"
			}

			// 转换items格式
			items := make([]InboundApprovalItem, 0)
			if len(req.Items) > 0 {
				for _, item := range req.Items {
					items = append(items, InboundApprovalItem{
						ID:               item.ID,
						ApprovedQuantity: int(item.ApprovedQuantity),
					})
				}
			}

			// 使用工作流引擎处理审批
			if err := wfIntegration.ProcessInboundApproval(order.ID, uid, name, workflow.ActionApprove, remark, items); err != nil {
				response.BadRequest(c, err.Error())
				return
			}

			// 重新加载订单
			db.Preload("Items").First(&order, id)

			// 记录操作日志
			if err := audit.LogApprove(&uid, name, audit.ModuleInbound, audit.ResourceInboundOrder,
				order.ID, order.OrderNo, remark); err != nil {
				fmt.Printf("记录操作日志失败: %v\n", err)
			}

			response.SuccessWithMessage(c, order.ToDTOWithEnrichment(db), "审批通过")
			return
		}

		// 如果没有工作流，使用原有逻辑（向后兼容）
		if order.Status != StatusPending {
			response.BadRequest(c, "只有待处理状态的入库单才能审批")
			return
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			// If request body is empty or has no items, approve all quantities as requested
			for i := range order.Items {
				req.Items = append(req.Items, struct {
					ID               uint    `json:"id"`
					ApprovedQuantity float64 `json:"approved_quantity"`
				}{
					ID:               order.Items[i].ID,
					ApprovedQuantity: order.Items[i].Quantity,
				})
			}
		}

		// If items is empty, approve all quantities
		if len(req.Items) == 0 {
			for i := range order.Items {
				req.Items = append(req.Items, struct {
					ID               uint    `json:"id"`
					ApprovedQuantity float64 `json:"approved_quantity"`
				}{
					ID:               order.Items[i].ID,
					ApprovedQuantity: order.Items[i].Quantity,
				})
			}
		}

		// Add remark if provided
		if remark != "" {
			if order.Remark != "" {
				order.Remark = order.Remark + "\n审批备注：" + remark
			} else {
				order.Remark = "审批备注：" + remark
			}
		}

		order.UpdatedAt = time.Now()

		// Create a map for approved quantities
		approvedQtyMap := make(map[uint]float64)
		for _, itemReq := range req.Items {
			approvedQtyMap[itemReq.ID] = itemReq.ApprovedQuantity
		}

		// Sync to stock for each item
		for _, item := range order.Items {
			// Get approved quantity, default to original quantity
			approvedQty := item.Quantity
			if qty, ok := approvedQtyMap[item.ID]; ok {
				approvedQty = qty
			}

			// Find or create stock record for the material
			type StockRecord struct {
				ID          uint
				ProjectID   uint
				MaterialID  uint
				WarehouseID *uint
				Quantity    float64
				SafetyStock float64
				Location    string
				UnitCost    float64
			}

			// Try to find existing stock for this project and material
			var quantityBefore, quantityAfter float64
			var stockRecord StockRecord
			var detail string

			if err := db.Table("stocks").
				Where("project_id = ? AND material_id = ?", order.ProjectID, item.MaterialID).
				First(&stockRecord).Error; err != nil {
				// Get material details to get unit
				var material struct {
					Unit string
				}
				db.Table("material_master").Where("id = ?", item.MaterialID).Select("unit").First(&material)

				// Create new stock record
				quantityBefore = 0
				quantityAfter = approvedQty
				stockRecord = StockRecord{
					ProjectID:   order.ProjectID,
					MaterialID:  item.MaterialID,
					WarehouseID: nil,
					Quantity:    quantityAfter,
					SafetyStock: 0,
					Location:    "",
					UnitCost:    0,
				}
				if err := db.Table("stocks").Create(&stockRecord).Error; err != nil {
					response.InternalError(c, fmt.Sprintf("创建库存记录失败(material_id=%d): %v", item.MaterialID, err))
					return
				}

				// Log the stock operation with material unit
				detail = fmt.Sprintf("入库 %.2f %s，备注：入库单 %s", approvedQty, material.Unit, order.OrderNo)
			} else {
				// Update existing stock quantity
				quantityBefore = stockRecord.Quantity
				stockRecord.Quantity += approvedQty
				quantityAfter = stockRecord.Quantity
				if err := db.Table("stocks").Save(&stockRecord).Error; err != nil {
					response.InternalError(c, fmt.Sprintf("更新库存记录失败(stock_id=%d): %v", stockRecord.ID, err))
					return
				}

				// Get material unit for logging
				var material struct {
					Unit string
				}
				db.Table("material_master").Where("id = ?", item.MaterialID).Select("unit").First(&material)

				// Log the stock operation
				detail = fmt.Sprintf("入库 %.2f %s，备注：入库单 %s", approvedQty, material.Unit, order.OrderNo)
			}

			// Get user ID for logging
			userID, _ := c.Get("current_user_id")
			var uid uint
			if userID != nil {
				if id, ok := userID.(int64); ok {
					uid = uint(id)
				}
			}

			// Get project ID from order
			projectID := order.ProjectID

			// Create StockLog entry with full details
			stockLog := map[string]interface{}{
				"stock_id":        stockRecord.ID,
				"type":            "in",
				"quantity":        approvedQty,
				"quantity_before": quantityBefore,
				"quantity_after":  quantityAfter,
				"source_type":     "inbound",
				"source_id":       order.ID,
				"source_no":       order.OrderNo,
				"remark":          detail,
				"project_id":      projectID,
				"material_id":     item.MaterialID,
				"user_id":         uid,
			}
			if err := db.Table("stock_logs").Create(&stockLog).Error; err != nil {
				response.InternalError(c, fmt.Sprintf("创建库存日志失败: %v", err))
				return
			}

			// Also create StockOpLog entry
			opLog := map[string]interface{}{
				"stock_id": stockRecord.ID,
				"op_type":  "in",
				"detail":   detail,
				"user_id":  uid,
				"time":     time.Now(),
			}
			if err := db.Table("stock_op_logs").Create(&opLog).Error; err != nil {
				response.InternalError(c, fmt.Sprintf("创建库存操作日志失败: %v", err))
				return
			}
		}


		// Update order status to completed
		order.Status = StatusCompleted

		if err := db.Save(&order).Error; err != nil {
			response.InternalError(c, err.Error())
			return
		}

		// Reload with items for enrichment
		db.Preload("Items").First(&order, order.ID)

		// Get user info for logging
		userID, _ := c.Get("current_user_id")
		var uid uint
		if userID != nil {
			if id, ok := userID.(int64); ok {
				uid = uint(id)
			}
		}
		userName, _ := c.Get("username")
		var name string
		if userName != nil {
			name = userName.(string)
		} else {
			name = "未知用户"
		}

		// 记录操作日志
		if err := audit.LogApprove(&uid, name, audit.ModuleInbound, audit.ResourceInboundOrder,
			order.ID, order.OrderNo, remark); err != nil {
			fmt.Printf("记录操作日志失败: %v\n", err)
		}

		response.SuccessWithMessage(c, order.ToDTOWithEnrichment(db), "入库单已批准")
	})

	// Reject inbound order
	g.POST("/inbound-orders/:id/reject", auth.PermissionMiddleware(db, "inbound_approve"), func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
		var order InboundOrder
		if err := db.Preload("Items").First(&order, id).Error; err != nil {
			response.NotFound(c, "入库单不存在")
			return
		}

		// 检查是否有工作流实例
		wfIntegration := NewWorkflowIntegration(db)
		instance, err := wfIntegration.GetInboundWorkflowStatus(order.ID)

		var req struct {
			Remark string `json:"remark"`
		}

		c.ShouldBindJSON(&req)

		// 如果存在工作流实例，使用工作流引擎处理
		if err == nil && instance != nil && instance.Status == workflow.InstanceStatusPending {
			// 获取用户信息
			userID, _ := c.Get("current_user_id")
			var uid uint
			if userID != nil {
				if id, ok := userID.(int64); ok {
					uid = uint(id)
				}
			}
			userName, _ := c.Get("username")
			var name string
			if userName != nil {
				name = userName.(string)
			} else {
				name = "未知用户"
			}

			// 使用工作流引擎处理拒绝
			if err := wfIntegration.ProcessInboundApproval(order.ID, uid, name, workflow.ActionReject, req.Remark, nil); err != nil {
				response.BadRequest(c, err.Error())
				return
			}

			// 重新加载订单
			db.Preload("Items").First(&order, id)

			// 记录操作日志
			if err := audit.LogReject(&uid, name, audit.ModuleInbound, audit.ResourceInboundOrder,
				order.ID, order.OrderNo, req.Remark); err != nil {
				fmt.Printf("记录操作日志失败: %v\n", err)
			}

			response.SuccessWithMessage(c, order.ToDTOWithEnrichment(db), "已拒绝")
			return
		}

		// 如果没有工作流，使用原有逻辑（向后兼容）
		if order.Status != StatusPending {
			response.BadRequest(c, "只有待处理状态的入库单才能拒绝")
			return
		}

		// Update status to rejected
		order.Status = StatusRejected

		// Add remark if provided
		if req.Remark != "" {
			if order.Remark != "" {
				order.Remark = order.Remark + "\n拒绝原因：" + req.Remark
			} else {
				order.Remark = "拒绝原因：" + req.Remark
			}
		}

		order.UpdatedAt = time.Now()

		if err := db.Save(&order).Error; err != nil {
			response.InternalError(c, err.Error())
			return
		}

		// Reload with items for enrichment
		db.Preload("Items").First(&order, order.ID)

		// Get user info for logging
		userID, _ := c.Get("current_user_id")
		var uid uint
		if userID != nil {
			if id, ok := userID.(int64); ok {
				uid = uint(id)
			}
		}
		userName, _ := c.Get("username")
		var name string
		if userName != nil {
			name = userName.(string)
		} else {
			name = "未知用户"
		}

		// 记录操作日志
		if err := audit.LogReject(&uid, name, audit.ModuleInbound, audit.ResourceInboundOrder,
			order.ID, order.OrderNo, req.Remark); err != nil {
			fmt.Printf("记录操作日志失败: %v\n", err)
		}

		response.SuccessWithMessage(c, order.ToDTOWithEnrichment(db), "入库单已拒绝")
	})

	// Get workflow history for inbound order
	g.GET("/inbound-orders/:id/workflow-history", auth.PermissionMiddleware(db, "inbound_view"), func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

		// Get workflow integration
		wfIntegration := NewWorkflowIntegration(db)
		approvals, err := wfIntegration.GetInboundWorkflowApprovals(uint(id))
		if err != nil {
			// If no workflow instance found, return empty array
			response.Success(c, []any{})
			return
		}

		// Transform approval records to match frontend format
		data := make([]map[string]any, 0)
		for _, approval := range approvals {
			item := map[string]any{
				"id":            approval.ID,
				"instance_id":   approval.InstanceID,
				"node_id":       approval.NodeID,
				"node_key":      approval.NodeKey,
				"approver_id":   approval.ApproverID,
				"approver_name": approval.ApproverName,
				"action":        approval.Action,
				"remark":        approval.Remark,
				"approved_at":   approval.ApprovedAt,
				"created_at":    approval.CreatedAt,
			}
			data = append(data, item)
		}

		response.Success(c, data)
	})

	// Submit inbound order (same as create, for compatibility)
	g.POST("/inbound/submit", auth.PermissionMiddleware(db, "inbound_create"), func(c *gin.Context) {
		// Get user info from JWT middleware
		userIDInt64, _ := c.Get("current_user_id")
		username, _ := c.Get("current_username")

		var creatorID uint
		if userIDInt64 != nil {
			if id, ok := userIDInt64.(int64); ok {
				creatorID = uint(id)
			}
		}

		var creatorName string
		if username != nil {
			if name, ok := username.(string); ok {
				creatorName = name
			}
		}

		var req struct {
			Supplier    string `json:"supplier"`
			Contact     string `json:"contact"`
			ProjectID   string `json:"project_id"`
			Notes       string `json:"notes"`
			Remark      string `json:"remark"`
			TotalAmount float64 `json:"total_amount"`
			Items       []struct {
				MaterialID uint    `json:"material_id"`
				Quantity   float64 `json:"quantity"`
				UnitPrice  float64 `json:"unit_price"`
				Remark     string  `json:"remark"`
			} `json:"items"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误")
			return
		}

		// Parse and validate project ID
		var projectIDUint uint
		if req.ProjectID != "" {
			pid, err := strconv.ParseUint(req.ProjectID, 10, 64)
			if err != nil || pid == 0 {
				response.BadRequest(c, "项目ID格式无效")
				return
			}
			projectIDUint = uint(pid)

			// Verify project exists
			var count int64
			if err := db.Table("projects").Where("id = ?", projectIDUint).Count(&count).Error; err != nil || count == 0 {
				response.BadRequest(c, "指定的项目不存在")
				return
			}
		} else {
			response.BadRequest(c, "请选择所属项目")
			return
		}

		if len(req.Items) == 0 {
			response.BadRequest(c, "入库单物资列表不能为空")
			return
		}

		orderNo := fmt.Sprintf("RK%s", time.Now().Format("20060102150405"))
		order := InboundOrder{
			OrderNo:     orderNo,
			Supplier:    strings.TrimSpace(req.Supplier),
			Contact:     strings.TrimSpace(req.Contact),
			ProjectID:   projectIDUint,
			CreatorID:   creatorID,
			CreatorName: creatorName,
			Status:      StatusPending,
			Notes:       strings.TrimSpace(req.Notes),
			Remark:      strings.TrimSpace(req.Remark),
		}

		totalAmount := 0
		for _, item := range req.Items {
			orderItem := InboundOrderItem{
				MaterialID: item.MaterialID,
				Quantity:   item.Quantity,
				UnitPrice:  item.UnitPrice,
				Remark:     item.Remark,
			}
			order.Items = append(order.Items, orderItem)
			totalAmount += int(item.Quantity * item.UnitPrice * 100)
		}
		order.TotalAmount = totalAmount

		if err := db.Create(&order).Error; err != nil {
			response.InternalError(c, err.Error())
			return
		}

		// Reload with items for enrichment
		db.Preload("Items").First(&order, order.ID)
		response.Created(c, order.ToDTOWithEnrichment(db), "入库单提交成功")
	})

	g.GET("/inbound/template", auth.PermissionMiddleware(db, "inbound_view"), func(c *gin.Context) {
		template := map[string]any{"supplier": "示例供应商", "contact": "联系人-13800138000", "project_id": 1, "notes": "备注信息", "items": []map[string]any{{"material_id": 1, "quantity": 10, "unit_price": 100.50, "remark": "第一批"}}}
		c.Header("Content-Disposition", "attachment; filename=inbound_template.json")
		response.SuccessWithMeta(c, map[string]any{
			"template": template,
			"fields":   []string{"supplier", "contact", "project_id", "notes", "items"},
		}, nil)
	})

	g.POST("/inbound/import", auth.PermissionMiddleware(db, "inbound_create"), func(c *gin.Context) {
		// Get user info from JWT middleware
		userIDInt64, _ := c.Get("current_user_id")
		username, _ := c.Get("current_username")

		var creatorID uint
		if userIDInt64 != nil {
			if id, ok := userIDInt64.(int64); ok {
				creatorID = uint(id)
			}
		}

		var creatorName string
		if username != nil {
			if name, ok := username.(string); ok {
				creatorName = name
			}
		}

		var req struct {
			Orders []struct {
				Supplier  string `json:"supplier"`
				Contact   string `json:"contact"`
				ProjectID string `json:"project_id"`
				Notes     string `json:"notes"`
				Items     []struct {
					MaterialID uint    `json:"material_id" binding:"required"`
					Quantity   float64 `json:"quantity" binding:"required"`
					UnitPrice  float64 `json:"unit_price"`
					Remark     string  `json:"remark"`
				} `json:"items" binding:"required"`
			} `json:"orders" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误")
			return
		}

		successCount := 0
		errorCount := 0
		errors := []string{}
		for i, orderReq := range req.Orders {
			if len(orderReq.Items) == 0 {
				errorCount++
				errors = append(errors, strconv.Itoa(i+1)+"行: 入库单物资列表不能为空")
				continue
			}

			// Parse and validate project ID
			var projectIDUint uint
			if orderReq.ProjectID != "" {
				pid, err := strconv.ParseUint(orderReq.ProjectID, 10, 64)
				if err != nil || pid == 0 {
					errorCount++
					errors = append(errors, strconv.Itoa(i+1)+"行: 项目ID格式无效")
					continue
				}
				projectIDUint = uint(pid)

				// Verify project exists
				var count int64
				if err := db.Table("projects").Where("id = ?", projectIDUint).Count(&count).Error; err != nil || count == 0 {
					errorCount++
					errors = append(errors, strconv.Itoa(i+1)+"行: 指定的项目不存在")
					continue
				}
			} else {
				errorCount++
				errors = append(errors, strconv.Itoa(i+1)+"行: 请选择所属项目")
				continue
			}

			orderNo := fmt.Sprintf("RK%s%03d", time.Now().Format("20060102150405"), i)
			order := InboundOrder{OrderNo: orderNo, Supplier: strings.TrimSpace(orderReq.Supplier), Contact: strings.TrimSpace(orderReq.Contact), ProjectID: projectIDUint, CreatorID: creatorID, CreatorName: creatorName, Status: StatusPending, Notes: strings.TrimSpace(orderReq.Notes)}
			totalAmount := 0
			for _, item := range orderReq.Items {
				if item.MaterialID == 0 || item.Quantity <= 0 {
					continue
				}
				orderItem := InboundOrderItem{
					MaterialID: item.MaterialID,
					Quantity:   item.Quantity,
					UnitPrice:  item.UnitPrice,
					Remark:     item.Remark,
				}
				order.Items = append(order.Items, orderItem)
				totalAmount += int(item.Quantity * item.UnitPrice * 100)
			}
			order.TotalAmount = totalAmount
			if err := db.Create(&order).Error; err != nil {
				errorCount++
				errors = append(errors, strconv.Itoa(i+1)+"行: "+err.Error())
			} else {
				successCount++
			}
		}
		response.SuccessWithMessage(c, map[string]interface{}{
			"success_count": successCount,
			"error_count":   errorCount,
			"errors":        errors,
		}, "导入完成，成功"+strconv.Itoa(successCount)+"条，失败"+strconv.Itoa(errorCount)+"条")
	})
}