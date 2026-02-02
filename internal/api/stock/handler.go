package stock

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("stock")
	r.Use(jwtpkg.TokenMiddleware())

	// list stocks (supports search, filters, pagination, sorting)
	r.GET("/stocks", auth.PermissionMiddleware(db, "stock_view"), func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}
		search := c.Query("search")
		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)
		quantityMin, _ := strconv.ParseFloat(c.Query("quantity_min"), 64)
		status := c.Query("status")
		category := c.Query("category")

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
		query := db.Model(&Stock{})

		// join material_master table to get material details
		query = query.Joins("LEFT JOIN material_master ON material_master.id = stocks.material_id")

		// generic search - search in material_master table
		if search != "" {
			query = query.Where("material_master.name LIKE ? OR material_master.code LIKE ? OR material_master.specification LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
		}

		// category filter
		if category != "" {
			query = query.Where("material_master.category LIKE ?", "%"+category+"%")
		}

		if quantityMin > 0 {
			query = query.Where("stocks.quantity >= ?", quantityMin)
		}

		// 获取用户可访问的项目ID列表（数据权限过滤）
		projectIDs, err := auth.GetAccessibleProjectIDs(c, db)
		if err != nil {
			response.InternalError(c, "获取用户项目权限失败")
			return
		}

		// 应用项目过滤（使用 stocks.project_id）
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
				query = query.Where("stocks.project_id IN ?", projectIDsFilter)
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
				query = query.Where("stocks.project_id = ?", projectID)
			} else {
				// 如果没有指定 project_id 参数，自动过滤用户可访问的项目
				query = query.Where("stocks.project_id IN ?", projectIDs)
			}
		} else {
			// 管理员
			if len(projectIDsFilter) > 0 {
				query = query.Where("stocks.project_id IN ?", projectIDsFilter)
			} else if projectID > 0 {
				query = query.Where("stocks.project_id = ?", projectID)
			}
		}

		// status filter: normal, low, shortage
		if status != "" {
			switch status {
			case "normal":
				// 库存正常：quantity > safety_stock
				query = query.Where("stocks.quantity > stocks.safety_stock")
			case "low":
				// 库存偏低：quantity <= safety_stock AND quantity > 0
				query = query.Where("stocks.quantity <= stocks.safety_stock AND stocks.quantity > 0")
			case "shortage":
				// 库存不足：quantity = 0
				query = query.Where("stocks.quantity = 0")
			}
		}

		query.Count(&total)

		// Select both stocks and material_master columns
		var results []map[string]any
		db.Model(&Stock{}).
			Select("stocks.*, material_master.code as material_code, material_master.name as material_name, material_master.specification, material_master.unit, material_master.category as category_name, projects.name as project_name, stocks.safety_stock as min_stock, stocks.safety_stock * 2 as max_stock").
			Joins("LEFT JOIN material_master ON material_master.id = stocks.material_id").
			Joins("LEFT JOIN projects ON stocks.project_id = projects.id").
			Where(query).
			Offset((page-1)*pageSize).Limit(pageSize).Order("stocks.updated_at DESC").
			Scan(&results)

		response.SuccessWithPagination(c, results, int64(page), int64(pageSize), total)
	})

	// get stock alerts (stocks with quantity <= safety_stock)
	r.GET("/stocks/alerts", auth.PermissionMiddleware(db, "stock_view"), func(c *gin.Context) {
		query := db.Model(&Stock{}).
			Select("stocks.id, stocks.material_id, stocks.project_id, stocks.quantity, stocks.safety_stock as min_stock, stocks.safety_stock * 2 as max_stock, stocks.location, stocks.unit_cost, stocks.created_at, stocks.updated_at, material_master.code as material_code, material_master.name as material_name, material_master.specification, material_master.unit, material_master.category as category_name, projects.name as project_name").
			Joins("LEFT JOIN material_master ON stocks.material_id = material_master.id").
			Joins("LEFT JOIN projects ON stocks.project_id = projects.id").
			Where("stocks.quantity <= stocks.safety_stock")

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
				response.Success(c, []map[string]any{})
				return
			}
			query = query.Where("stocks.project_id IN ?", projectIDs)
		}

		var results []map[string]any
		query.Order("stocks.updated_at DESC").Scan(&results)

		response.Success(c, results)
	})

	// create stock
	r.POST("/stocks", auth.PermissionMiddleware(db, "stock_create"), func(c *gin.Context) {
		var req struct {
			ProjectID   uint    `json:"project_id" binding:"required"`
			MaterialID  uint    `json:"material_id" binding:"required"`
			Quantity    float64 `json:"quantity"`
			SafetyStock float64 `json:"safety_stock"`
			Location    string  `json:"location"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// Check if material_master exists
		var materialMaster map[string]any
		if err := db.Table("material_master").Where("id = ?", req.MaterialID).First(&materialMaster).Error; err != nil {
			response.BadRequest(c, "物资不存在")
			return
		}

		// Check if stock already exists for this project and material
		var existing Stock
		result := db.Where("project_id = ? AND material_id = ?", req.ProjectID, req.MaterialID).First(&existing)
		if result.RowsAffected > 0 {
			response.BadRequest(c, "该项目该物资的库存记录已存在")
			return
		}

		// Create stock
		stock := Stock{
			ProjectID:   req.ProjectID,
			MaterialID:  req.MaterialID,
			Quantity:    req.Quantity,
			SafetyStock: req.SafetyStock,
			Location:    req.Location,
		}

		if err := db.Create(&stock).Error; err != nil {
			response.InternalError(c, "创建库存失败")
			return
		}

		// Log stock operation - treat initial stock as inbound
		logOp(db, stock.ID, "in", fmt.Sprintf("创建库存记录，初始库存 %.2f", req.Quantity), c, 0, req.Quantity, 0, req.Quantity, "adjust", 0, "")

		response.Created(c, getStockWithMaterial(db, stock.ID), "库存记录创建成功")
	})

	// get single stock
	r.GET("/stocks/:id", auth.PermissionMiddleware(db, "stock_view"), func(c *gin.Context) {
		id := c.Param("id")
		var stock Stock
		if err := db.First(&stock, id).Error; err != nil {
			response.NotFound(c, "库存不存在")
			return
		}
		response.Success(c, getStockWithMaterial(db, stock.ID))
	})

	// get stock logs for a specific stock
	r.GET("/stocks/:id/logs", auth.PermissionMiddleware(db, "stock_view"), func(c *gin.Context) {
		id := c.Param("id")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		var total int64
		db.Model(&StockLog{}).Where("stock_id = ?", id).Count(&total)

		var logs []StockLog
		db.Where("stock_id = ?", id).
			Offset((page-1)*pageSize).Limit(pageSize).
			Order("created_at DESC").
			Find(&logs)

		// Collect unique user IDs
		userIDs := make([]uint, 0)
		userIDSet := make(map[uint]bool)
		for _, log := range logs {
			if log.UserID != nil && *log.UserID > 0 {
				if !userIDSet[*log.UserID] {
					userIDs = append(userIDs, *log.UserID)
					userIDSet[*log.UserID] = true
				}
			}
		}

		// Query usernames in batch
		type UserName struct {
			ID       uint
			Username string
		}
		userNames := make(map[uint]string)
		if len(userIDs) > 0 {
			var users []UserName
			db.Table("users").Select("id, username").Where("id IN ?", userIDs).Find(&users)
			for _, u := range users {
				userNames[u.ID] = u.Username
			}
		}

		out := make([]map[string]any, 0, len(logs))
		for _, log := range logs {
			userName := "系统"
			if log.UserID != nil && *log.UserID > 0 {
				if name, ok := userNames[*log.UserID]; ok && name != "" {
					userName = name
				}
			}
			logData := log.ToDTO(userName)
			out = append(out, logData)
		}

		response.SuccessWithMeta(c, out, map[string]any{"total": total})
	})

	// update stock
	r.PUT("/stocks/:id", auth.PermissionMiddleware(db, "stock_edit"), func(c *gin.Context) {
		id := c.Param("id")
		var stock Stock
		if err := db.First(&stock, id).Error; err != nil {
			response.NotFound(c, "库存不存在")
			return
		}

		var req struct {
			Quantity    *float64 `json:"quantity"`
			SafetyStock *float64 `json:"safety_stock"`
			Location    *string  `json:"location"`
			UnitCost    *float64 `json:"unit_cost"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// Update fields
		if req.Quantity != nil {
			if *req.Quantity < 0 {
				response.BadRequest(c, "数量必须为0或正数")
				return
			}
			stock.Quantity = *req.Quantity
		}
		if req.SafetyStock != nil {
			stock.SafetyStock = *req.SafetyStock
		}
		if req.Location != nil {
			stock.Location = *req.Location
		}
		if req.UnitCost != nil {
			stock.UnitCost = *req.UnitCost
		}

		// Save stock changes
		if err := db.Save(&stock).Error; err != nil {
			response.InternalError(c, "更新库存失败")
			return
		}

		// Log stock operation
		logOp(db, stock.ID, "update", "更新库存记录", c, 0, 0, 0, 0, "adjust", 0, "")

		response.SuccessWithMessage(c, getStockWithMaterial(db, stock.ID), "库存记录更新成功")
	})

	// delete stock
	r.DELETE("/stocks/:id", auth.PermissionMiddleware(db, "stock_delete"), func(c *gin.Context) {
		id := c.Param("id")
		var stock Stock
		if err := db.First(&stock, id).Error; err != nil {
			response.NotFound(c, "库存不存在")
			return
		}

		if err := db.Delete(&stock).Error; err != nil {
			response.InternalError(c, "删除库存失败")
			return
		}

		// Log stock operation
		logOp(db, stock.ID, "delete", "删除库存记录", c, 0, 0, 0, 0, "adjust", 0, "")

		response.SuccessWithMessage(c, nil, "库存记录删除成功")
	})

	// stock logs
	r.GET("/stock-logs", auth.PermissionMiddleware(db, "stock_view"), func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}
		stockID, _ := strconv.ParseUint(c.Query("stock_id"), 10, 64)
		logType := c.Query("type")
		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)

		var total int64
		query := db.Model(&StockLog{})

		if stockID > 0 {
			query = query.Where("stock_id = ?", stockID)
		}
		if logType != "" {
			query = query.Where("type = ?", logType)
		}

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
			// 如果没有指定 project_id 参数，自动过滤用户可访问的项目
			if projectID == 0 {
				query = query.Where("project_id IN ?", projectIDs)
			} else {
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
			}
		} else if projectID > 0 {
			// 管理员，应用指定的 project_id 过滤
			query = query.Where("project_id = ?", projectID)
		}

		query.Count(&total)
		var logs []StockLog
		query.Offset((page-1)*pageSize).Limit(pageSize).Order("created_at DESC").Find(&logs)

		// Collect unique user IDs
		userIDs := make([]uint, 0)
		userIDSet := make(map[uint]bool)
		for _, log := range logs {
			if log.UserID != nil && *log.UserID > 0 {
				if !userIDSet[*log.UserID] {
					userIDs = append(userIDs, *log.UserID)
					userIDSet[*log.UserID] = true
				}
			}
		}

		// Query usernames in batch
		type UserName struct {
			ID       uint
			Username string
		}
		userNames := make(map[uint]string)
		if len(userIDs) > 0 {
			var users []UserName
			db.Table("users").Select("id, username").Where("id IN ?", userIDs).Find(&users)
			for _, u := range users {
				userNames[u.ID] = u.Username
			}
		}

		out := make([]map[string]any, 0, len(logs))
		for _, log := range logs {
			userName := "系统"
			if log.UserID != nil && *log.UserID > 0 {
				if name, ok := userNames[*log.UserID]; ok && name != "" {
					userName = name
				}
			}
			out = append(out, log.ToDTO(userName))
		}
		response.SuccessWithPagination(c, out, int64(page), int64(pageSize), total)
	})

	// stock in (入库操作)
	r.POST("/stocks/:id/in", auth.PermissionMiddleware(db, "stock_in"), func(c *gin.Context) {
		id := c.Param("id")
		var stock Stock
		if err := db.First(&stock, id).Error; err != nil {
			response.NotFound(c, "库存不存在")
			return
		}

		var req struct {
			Quantity float64 `json:"quantity" binding:"required,gt=0"`
			Remark   string  `json:"remark"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误")
			return
		}

		// 记录操作前的数量
		quantityBefore := stock.Quantity
		stock.Quantity += req.Quantity
		quantityAfter := stock.Quantity

		if err := db.Save(&stock).Error; err != nil {
			response.InternalError(c, err.Error())
			return
		}

		logOp(db, stock.ID, "in", fmt.Sprintf("入库 %.2f，备注：%s", req.Quantity, req.Remark), c, 0, req.Quantity, quantityBefore, quantityAfter, "adjust", 0, "")

		response.SuccessWithMessage(c, getStockWithMaterial(db, stock.ID), "入库成功")
	})

	// stock out (出库操作)
	r.POST("/stocks/:id/out", auth.PermissionMiddleware(db, "stock_out"), func(c *gin.Context) {
		id := c.Param("id")
		var stock Stock
		if err := db.First(&stock, id).Error; err != nil {
			response.NotFound(c, "库存不存在")
			return
		}

		var req struct {
			Quantity float64 `json:"quantity" binding:"required,gt=0"`
			Remark   string  `json:"remark"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误")
			return
		}

		if stock.Quantity < req.Quantity {
			response.BadRequest(c, "库存不足")
			return
		}

		// 记录操作前的数量
		quantityBefore := stock.Quantity
		stock.Quantity -= req.Quantity
		quantityAfter := stock.Quantity

		if err := db.Save(&stock).Error; err != nil {
			response.InternalError(c, err.Error())
			return
		}

		logOp(db, stock.ID, "out", fmt.Sprintf("出库 %.2f，备注：%s", req.Quantity, req.Remark), c, 0, req.Quantity, quantityBefore, quantityAfter, "adjust", 0, "")

		response.SuccessWithMessage(c, getStockWithMaterial(db, stock.ID), "出库成功")
	})

	// stock adjust (库存调整)
	r.POST("/stocks/:id/adjust", auth.PermissionMiddleware(db, "stock_edit"), func(c *gin.Context) {
		id := c.Param("id")
		var stock Stock
		if err := db.First(&stock, id).Error; err != nil {
			response.NotFound(c, "库存不存在")
			return
		}

		var req struct {
			Quantity float64 `json:"quantity" binding:"required"`
			Remark   string  `json:"remark" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误或备注不能为空")
			return
		}

		oldQty := stock.Quantity
		stock.Quantity = req.Quantity
		if err := db.Save(&stock).Error; err != nil {
			response.InternalError(c, err.Error())
			return
		}

		// Calculate adjustment delta
		delta := req.Quantity - oldQty
		// Log as in/out based on delta to ensure it shows in stock logs
		var logType string
		var logQty float64
		if delta > 0 {
			logType = "in"
			logQty = delta
		} else if delta < 0 {
			logType = "out"
			logQty = -delta
		} else {
			logType = "adjust"
			logQty = 0
		}
		logOp(db, stock.ID, logType, fmt.Sprintf("库存调整 %.2f -> %.2f，原因：%s", oldQty, req.Quantity, req.Remark), c, 0, logQty, oldQty, req.Quantity, "adjust", 0, "")

		response.SuccessWithMessage(c, getStockWithMaterial(db, stock.ID), "库存调整成功")
	})

	// delete stock log
	r.DELETE("/stock-logs/:id", auth.PermissionMiddleware(db, "stock_delete"), func(c *gin.Context) {
		id := c.Param("id")
		var log StockLog
		if err := db.First(&log, id).Error; err != nil {
			response.NotFound(c, "日志不存在")
			return
		}

		if err := db.Delete(&log).Error; err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, nil, "日志删除成功")
	})

	// export stocks
	r.GET("/stocks/export", auth.PermissionMiddleware(db, "stock_export"), func(c *gin.Context) {
		search := c.Query("search")
		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)
		status := c.Query("status")
		category := c.Query("category")

		// Query stocks with material details via JOIN
		type StockWithMaterial struct {
			ID           uint
			ProjectID    uint
			MaterialID   uint
			Quantity     float64
			SafetyStock  float64
			Location     string
			UnitCost     float64
			CreatedAt    time.Time
			UpdatedAt    time.Time
			Code         string // from material_master
			Name         string // from material_master
			Specification string // from material_master
			Unit         string // from material_master
			Category     string // from material_master
			ProjectName  string // from projects
		}

		var stocks []StockWithMaterial
		query := db.Table("stocks").
			Select("stocks.id, stocks.project_id, stocks.material_id, stocks.quantity, stocks.safety_stock, stocks.location, stocks.unit_cost, stocks.created_at, stocks.updated_at, material_master.code, material_master.name, material_master.specification, material_master.unit, material_master.category, projects.name as project_name").
			Joins("LEFT JOIN material_master ON stocks.material_id = material_master.id").
			Joins("LEFT JOIN projects ON stocks.project_id = projects.id")

		if search != "" {
			query = query.Where("material_master.name LIKE ? OR material_master.specification LIKE ? OR material_master.code LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
		}

		// category filter
		if category != "" {
			query = query.Where("material_master.category LIKE ?", "%"+category+"%")
		}

		// 获取用户可访问的项目ID列表（数据权限过滤）
		projectIDs, err := auth.GetAccessibleProjectIDs(c, db)
		if err != nil {
			response.InternalError(c, "获取用户项目权限失败")
			return
		}

		// 应用项目过滤
		if projectIDs != nil {
			if len(projectIDs) == 0 {
				// 用户无任何项目权限，返回空 Excel
				query = query.Where("1 = 0")
			} else {
				// 如果没有指定 project_id 参数，自动过滤用户可访问的项目
				if projectID == 0 {
					query = query.Where("stocks.project_id IN ?", projectIDs)
				} else {
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
					query = query.Where("stocks.project_id = ?", projectID)
				}
			}
		} else if projectID > 0 {
			// 管理员，应用指定的 project_id 过滤
			query = query.Where("stocks.project_id = ?", projectID)
		}

		// status filter: normal, low, shortage
		if status != "" {
			switch status {
			case "normal":
				query = query.Where("stocks.quantity > stocks.safety_stock")
			case "low":
				query = query.Where("stocks.quantity <= stocks.safety_stock AND stocks.quantity > 0")
			case "shortage":
				query = query.Where("stocks.quantity = 0")
			}
		}

		query.Find(&stocks)

		// Create Excel file
		f := excelize.NewFile()
		defer f.Close()

		sheet := "库存列表"
		f.NewSheet(sheet)
		f.DeleteSheet("Sheet1") // Remove default sheet

		// Set column headers
		headers := []string{"ID", "项目名称", "物资编码", "物资名称", "规格", "单位", "库存数量", "安全库存", "存储位置", "单位成本", "创建时间", "更新时间"}
		for i, header := range headers {
			f.SetCellValue(sheet, fmt.Sprintf("%s1", string(rune('A'+i))), header)
		}

		// Set column widths
		colWidths := []float64{8, 12, 12, 15, 15, 8, 10, 10, 15, 10, 18, 18}
		for i, width := range colWidths {
			f.SetColWidth(sheet, string(rune('A'+i)), string(rune('A'+i)), width)
		}

		// Fill data
		for idx, s := range stocks {
			row := idx + 2

			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), s.ID)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), s.ProjectName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), s.Code)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), s.Name)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), s.Specification)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), s.Unit)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), s.Quantity)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), s.SafetyStock)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), s.Location)
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), s.UnitCost)
			f.SetCellValue(sheet, fmt.Sprintf("K%d", row), s.CreatedAt.Format("2006-01-02 15:04:05"))
			f.SetCellValue(sheet, fmt.Sprintf("L%d", row), s.UpdatedAt.Format("2006-01-02 15:04:05"))
		}

		// Generate Excel file in memory
		buffer, err := f.WriteToBuffer()
		if err != nil {
			response.InternalError(c, "生成Excel文件失败")
			return
		}

		c.Header("Content-Disposition", "attachment; filename=库存导出.xlsx")
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
	})
}

// getStockWithMaterial retrieves a stock with its material details
func getStockWithMaterial(db *gorm.DB, stockID uint) map[string]any {
	var result map[string]any
	db.Table("stocks").
		Select("stocks.id, stocks.project_id, stocks.material_id, stocks.quantity, stocks.safety_stock as min_stock, stocks.safety_stock * 2 as max_stock, stocks.location, stocks.unit_cost, stocks.created_at, stocks.updated_at, material_master.code as material_code, material_master.name as material_name, material_master.specification, material_master.unit, material_master.category as category_name, projects.name as project_name").
		Joins("LEFT JOIN material_master ON stocks.material_id = material_master.id").
		Joins("LEFT JOIN projects ON stocks.project_id = projects.id").
		Where("stocks.id = ?", stockID).
		Scan(&result)
	return result
}

// logOp logs stock operations
// sourceType: inbound/requisition/adjust/transfer
// sourceID: 来源单据ID
// sourceNo: 来源单据号
func logOp(db *gorm.DB, stockID uint, opType string, detail string, c *gin.Context, requisitionID uint, quantity float64, quantityBefore, quantityAfter float64, sourceType string, sourceID uint, sourceNo string) {
	// Get current user info from context
	userID, _ := c.Get("current_user_id")
	var uid *uint
	if userID != nil {
		if id, ok := userID.(int64); ok {
			uidVal := uint(id)
			uid = &uidVal
		}
	}

	// Get stock details to get project_id and material_id
	var stock Stock
	db.First(&stock, stockID)

	// For in, out, and create operations, also create a StockLog entry
	if opType == "in" || opType == "out" {
		// Create StockLog entry with details
		var sourceIDPtr *uint
		if sourceID > 0 {
			sourceIDPtr = &sourceID
		}

		stockLog := StockLog{
			StockID:        stockID,
			Type:           opType,
			Quantity:       quantity,
			QuantityBefore: quantityBefore,
			QuantityAfter:  quantityAfter,
			SourceType:     sourceType,
			SourceID:       sourceIDPtr,
			SourceNo:       sourceNo,
			ProjectID:      stock.ProjectID,
			MaterialID:     stock.MaterialID,
			UserID:         uid,
			Remark:         detail,
			CreatedAt:      time.Now(),
		}
		if err := db.Create(&stockLog).Error; err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to create stock log: %v\n", err)
		}
	}

	// Always create StockOpLog entry
	var logID uint
	if requisitionID > 0 {
		logID = requisitionID
	}
	opLog := StockOpLog{
		StockID: stockID,
		OpType:  opType,
		LogID:   logID,
		Detail:  detail,
		UserID:  0,
		Time:    time.Now(),
	}

	if uid != nil {
		opLog.UserID = *uid
	}

	db.Create(&opLog)
}
