package material

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("material")
	r.Use(jwtpkg.TokenMiddleware())

	// list materials (supports search, filters, pagination, sorting)
	r.GET("/materials", auth.PermissionMiddleware(db, "material_view"), func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 { pageSize = 100 }
		search := c.Query("search")
		name := c.Query("name")
		material := c.Query("material")
		spec := c.Query("spec")
		specification := c.Query("specification")
		category := c.Query("category")
		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)
		filter := c.Query("filter") // "unstored" to filter fully stored materials
		includeChildren := c.Query("include_children") == "true" // 是否包含子项目

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
		} else if projectID > 0 && includeChildren {
			// 如果指定了project_id且要求包含子项目，递归获取所有子项目ID
			projectIDsFilter = append(projectIDsFilter, uint(projectID))

			// 递归查询子项目
			var getChildProjectIDs func(uint) []uint
			getChildProjectIDs = func(parentID uint) []uint {
				var children []struct {
					ID uint
				}
				db.Table("projects").Where("parent_id = ?", parentID).Find(&children)

				var ids []uint
				for _, child := range children {
					ids = append(ids, child.ID)
					// 递归获取子项目的子项目
					childIDs := getChildProjectIDs(child.ID)
					ids = append(ids, childIDs...)
				}
				return ids
			}

			childIDs := getChildProjectIDs(uint(projectID))
			projectIDsFilter = append(projectIDsFilter, childIDs...)
		}

		var total int64
		query := db.Model(&Material{})

		// generic search
		if search != "" {
			query = query.Where(
				"materials.name LIKE ? OR materials.material LIKE ? OR materials.spec LIKE ? OR materials.specification LIKE ? OR materials.category LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%",
			)
		}
		// field filters
		if name != "" { query = query.Where("materials.name LIKE ?", "%"+name+"%") }
		if material != "" { query = query.Where("materials.material LIKE ?", "%"+material+"%") }
		if spec != "" { query = query.Where("materials.spec LIKE ? OR materials.specification LIKE ?", "%"+spec+"%", "%"+spec+"%") }
		if specification != "" { query = query.Where("materials.specification LIKE ? OR materials.spec LIKE ?", "%"+specification+"%", "%"+specification+"%") }
		if category != "" { query = query.Where("materials.category LIKE ?", "%"+category+"%") }

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
				query = query.Where("materials.project_id IN ?", projectIDsFilter)
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
				query = query.Where("materials.project_id = ?", projectID)
			} else {
				// 如果没有指定 project_id 参数，自动过滤用户可访问的项目
				query = query.Where("materials.project_id IN ?", projectIDs)
			}
		} else {
			// 管理员
			if len(projectIDsFilter) > 0 {
				query = query.Where("materials.project_id IN ?", projectIDsFilter)
			} else if projectID > 0 {
				query = query.Where("materials.project_id = ?", projectID)
			}
		}

		// 过滤已完全入库的物资
		if filter == "unstored" {
			query = query.Where(`
				NOT EXISTS (
					SELECT 1
					FROM inbound_order_items ioi
					INNER JOIN inbound_orders io ON ioi.order_id = io.id
					WHERE ioi.material_id = materials.id
					AND io.status IN ('approved', 'completed')
					GROUP BY ioi.material_id
					HAVING SUM(ioi.quantity) >= materials.quantity
				)
			`)
		}

		query.Count(&total)

		// Query with JOIN to get project name and stock quantity
		type MaterialWithProject struct {
			ID            uint     `gorm:"column:id"`
			Code          *string  `gorm:"column:code"`
			Name          string   `gorm:"column:name"`
			Specification string   `gorm:"column:specification"`
			Unit          string   `gorm:"column:unit"`
			Price         float64  `gorm:"column:price"`
			Description   string   `gorm:"column:description"`
			Category      string   `gorm:"column:category"`
			Quantity      int      `gorm:"column:quantity"`           // materials表的数量（计划数量）
			StockQuantity float64  `gorm:"column:stock_quantity"`    // stocks表的数量（实际库存）
			ProjectID     *uint    `gorm:"column:project_id"`
			Material      string   `gorm:"column:material"`
			Spec          string   `gorm:"column:spec"`
			ProjectName   *string  `gorm:"column:project_name"`
		}

		var results []MaterialWithProject
		db.Table("materials").
			Select("materials.*, projects.name as project_name, COALESCE(stocks.quantity, 0) as stock_quantity").
			Joins("LEFT JOIN projects ON materials.project_id = projects.id").
			Joins("LEFT JOIN stocks ON stocks.material_id = materials.id").
			Where(query).
			Offset((page-1)*pageSize).Limit(pageSize).Order("materials.id DESC").
			Scan(&results)

		// 批量获取计划信息和到货信息（从material_plan_items表）
		materialIDs := make([]uint, len(results))
		for i, m := range results {
			materialIDs[i] = m.ID
		}

		type PlanInfo struct {
			MaterialID       uint
			PlannedQuantity  int
			ArrivedQuantity  int
		}
		var planInfos []PlanInfo
		db.Raw(`
			SELECT
				material_id,
				COALESCE(SUM(planned_quantity), 0) as planned_quantity,
				COALESCE(SUM(arrived_quantity), 0) as arrived_quantity
			FROM material_plan_items
			WHERE material_id = ANY($1)
			GROUP BY material_id
		`, materialIDs).Scan(&planInfos)

		// 构建计划信息映射
		planMap := make(map[uint]PlanInfo)
		for _, info := range planInfos {
			planMap[info.MaterialID] = info
		}

		out := make([]map[string]any, 0, len(results))
		for _, m := range results {
			var projectID uint
			if m.ProjectID != nil {
				projectID = *m.ProjectID
			}
			spec := m.Specification
			if spec == "" && m.Spec != "" {
				spec = m.Spec
			}

			// 从计划信息中获取计划数量和到货数量
			planInfo := planMap[m.ID]
			plannedQty := 0
			arrivedQty := 0
			if planInfo.MaterialID > 0 {
				plannedQty = planInfo.PlannedQuantity
				arrivedQty = planInfo.ArrivedQuantity
			}

			// 计算剩余数量和到货百分比
			remainingQty := plannedQty - arrivedQty
			arrivalPercentage := 0.0
			if plannedQty > 0 {
				arrivalPercentage = float64(arrivedQty) / float64(plannedQty) * 100
			}

			out = append(out, map[string]any{
				"id":                  m.ID,
				"code":                m.Code,
				"name":                m.Name,
				"specification":       spec,
				"unit":                m.Unit,
				"price":               m.Price,
				"description":         m.Description,
				"category":            m.Category,
				"quantity":            m.StockQuantity,      // 实际库存数量（从stocks表）
				"planned_quantity":    plannedQty,           // 计划数量（从material_plan_items统计）
				"arrived_quantity":    arrivedQty,           // 到货数量（从material_plan_items统计）
				"remaining_quantity":  remainingQty,          // 剩余数量 = 计划 - 到货
				"arrival_percentage":  arrivalPercentage,     // 到货百分比
				"is_fully_arrived":    arrivedQty >= plannedQty && plannedQty > 0, // 是否已完全到货
				"project_id":          projectID,
				"project_name":        m.ProjectName,
				"material":            m.Material,
				"spec":                m.Spec,
			})
		}

		response.SuccessWithPagination(c, out, int64(page), int64(pageSize), total)
	})

	// create material
	r.POST("/materials", auth.PermissionMiddleware(db, "material_create"), func(c *gin.Context) {
		var req struct {
			Code           string  `json:"code"`
			Name           string  `json:"name"`
			Specification  string  `json:"specification"`
			Unit           string  `json:"unit"`
			Price          float64 `json:"price"`
			Description    string  `json:"description"`
			Category       string  `json:"category"`
			Quantity       int     `json:"quantity"`
			ProjectID      string  `json:"project_id"`
			Material       string  `json:"material"`
			Spec           string  `json:"spec"`
		}
		if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, err.Error()); return }
		if req.Name == "" { response.BadRequest(c, "物资名称不能为空"); return }

		// Validate project_id - must be provided and parse to uint
		if req.ProjectID == "" {
			response.BadRequest(c, "请选择所属项目"); return
		}
		projectIDUint, err := strconv.ParseUint(req.ProjectID, 10, 64)
		if err != nil || projectIDUint == 0 {
			response.BadRequest(c, "项目ID格式无效"); return
		}
		projectID := uint(projectIDUint)

		var projectExists int64
		if db.Table("projects").Where("id = ?", projectID).Count(&projectExists); projectExists == 0 {
			response.BadRequest(c, "指定的项目不存在"); return
		}

		// check for duplicate code
		var code *string
		if req.Code != "" {
			code = &req.Code
			var existing Material
			if db.Where("code = ?", req.Code).First(&existing).Error == nil {
				response.BadRequest(c, "物资编码已存在"); return
			}
		}

		projectIDPtr := &projectID

		m := Material{
			Code: code,
			Name: req.Name,
			Specification: req.Specification,
			Unit: req.Unit,
			Price: req.Price,
			Description: req.Description,
			Category: req.Category,
			Quantity: req.Quantity,
			ProjectID: projectIDPtr,
			Material: req.Material,
			Spec: req.Spec,
		}
		if err := db.Create(&m).Error; err != nil { response.InternalError(c, err.Error()); return }
		response.Created(c, m.ToDTO(), "物资创建成功")
	})

	// get single material
	r.GET("/materials/:id", auth.PermissionMiddleware(db, "material_view"), func(c *gin.Context) {
		id := c.Param("id")
		var m Material
		if err := db.First(&m, id).Error; err != nil { response.NotFound(c, "物资不存在"); return }
		response.Success(c, m.ToDTO())
	})

	// update material
	r.PUT("/materials/:id", auth.PermissionMiddleware(db, "material_edit"), func(c *gin.Context) {
		id := c.Param("id")
		var m Material
		if err := db.First(&m, id).Error; err != nil { response.NotFound(c, "物资不存在"); return }

		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, err.Error()); return }

		// Update fields
		if v, ok := req["code"].(string); ok {
			var codePtr *string
			if v != "" {
				codePtr = &v
			}
			// Check if code is different
			if (m.Code == nil && codePtr != nil) || (m.Code != nil && codePtr == nil) || (m.Code != nil && codePtr != nil && *m.Code != *codePtr) {
				var existing Material
				if db.Where("code = ?", v).First(&existing).Error == nil && existing.ID != m.ID {
					response.BadRequest(c, "物资编码已存在"); return
				}
				m.Code = codePtr
			}
		}
		if v, ok := req["name"].(string); ok { m.Name = v }
		if v, ok := req["specification"].(string); ok { m.Specification = v }
		if v, ok := req["unit"].(string); ok { m.Unit = v }
		if v, ok := req["price"].(float64); ok { m.Price = v }
		if v, ok := req["description"].(string); ok { m.Description = v }
		if v, ok := req["category"].(string); ok { m.Category = v }
		if v, ok := req["quantity"].(float64); ok { m.Quantity = int(v) }

		// Handle project_id - accept both string and float64/float/integer
		if v, ok := req["project_id"].(string); ok {
			// String type (e.g., "123")
			if v != "" {
				projectIDUint, err := strconv.ParseUint(v, 10, 64)
				if err == nil && projectIDUint > 0 {
					pid := uint(projectIDUint)
					// Validate project exists
					var projectExists int64
					if db.Table("projects").Where("id = ?", pid).Count(&projectExists); projectExists == 0 {
						response.BadRequest(c, "指定的项目不存在"); return
					}
					m.ProjectID = &pid
				} else if v == "" {
					m.ProjectID = nil
				}
			} else {
				m.ProjectID = nil
			}
		} else if v, ok := req["project_id"].(float64); ok {
			// Float64/Number type (e.g., 123.0)
			if v > 0 {
				pid := uint(v)
				// Validate project exists
				var projectExists int64
				if db.Table("projects").Where("id = ?", pid).Count(&projectExists); projectExists == 0 {
					response.BadRequest(c, "指定的项目不存在"); return
				}
				m.ProjectID = &pid
			} else {
				m.ProjectID = nil
			}
		}
		if v, ok := req["material"].(string); ok { m.Material = v }
		if v, ok := req["spec"].(string); ok { m.Spec = v }

		db.Save(&m)
		response.SuccessWithMessage(c, m.ToDTO(), "物资更新成功")
	})

	// delete material
	r.DELETE("/materials/:id", auth.PermissionMiddleware(db, "material_delete"), func(c *gin.Context) {
		id := c.Param("id")
		var m Material
		if err := db.First(&m, id).Error; err != nil { response.NotFound(c, "物资不存在"); return }
		db.Delete(&m)
		response.SuccessWithMessage(c, nil, "物资删除成功")
	})

	// export materials to Excel
	r.GET("/materials/export", auth.PermissionMiddleware(db, "material_view"), func(c *gin.Context) {
		search := c.Query("search")
		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)

		query := db.Model(&Material{})
		if search != "" {
			query = query.Where(
				"name LIKE ? OR material LIKE ? OR spec LIKE ? OR specification LIKE ? OR category LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%",
			)
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
			}
		} else if projectID > 0 {
			// 管理员，应用指定的 project_id 过滤
			query = query.Where("project_id = ?", projectID)
		}

		var materials []Material
		query.Find(&materials)

		// Create Excel file
		f := excelize.NewFile()
		defer f.Close()
		
		sheet := "物资列表"
		f.NewSheet(sheet)
		f.DeleteSheet("Sheet1") // Remove default sheet
		
		// Set column headers
		headers := []string{"ID", "编码", "物资名称", "规格", "材质", "单位", "价格", "分类", "数量", "描述", "项目ID"}
		for i, header := range headers {
			f.SetCellValue(sheet, fmt.Sprintf("%s1", string(rune('A'+i))), header)
		}
		
		// Set column widths
		colWidths := []float64{8, 12, 15, 15, 12, 8, 10, 10, 8, 20, 10}
		for i, width := range colWidths {
			f.SetColWidth(sheet, string(rune('A'+i)), string(rune('A'+i)), width)
		}
		
		// Fill data
		for idx, m := range materials {
			row := idx + 2
			spec := m.Specification
			if spec == "" && m.Spec != "" {
				spec = m.Spec
			}
			
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), m.ID)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), m.Code)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), m.Name)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), spec)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), m.Material)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), m.Unit)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), m.Price)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), m.Category)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), m.Quantity)
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), m.Description)
			f.SetCellValue(sheet, fmt.Sprintf("K%d", row), m.ProjectID)
		}
		
		// Generate Excel file in memory
		buffer, err := f.WriteToBuffer()
		if err != nil {
			response.InternalError(c, "生成Excel文件失败")
			return
		}
		
		c.Header("Content-Disposition", "attachment; filename=物资导出.xlsx")
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
	})

	// import materials from JSON (with column mapping)
	r.POST("/material/materials/import", auth.PermissionMiddleware(db, "material_create"), func(c *gin.Context) {
		var req struct {
			Materials []map[string]any `json:"materials" binding:"required"`
			ProjectID uint            `json:"project_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误: "+err.Error())
			return
		}

		// Validate project_id
		if req.ProjectID == 0 {
			response.BadRequest(c, "请指定所属项目")
			return
		}
		var projectExists int64
		if db.Table("projects").Where("id = ?", req.ProjectID).Count(&projectExists); projectExists == 0 {
			response.BadRequest(c, "指定的项目不存在")
			return
		}

		successCount := 0
		errorCount := 0
		errors := []string{}

		for i, item := range req.Materials {
			// 提取字段值
			name, _ := item["name"].(string)
			if name == "" {
				// 尝试从其他可能的字段获取
				if val, ok := item["Name"]; ok {
					name, _ = val.(string)
				}
			}

			if name == "" {
				errorCount++
				errors = append(errors, fmt.Sprintf("第%d行: 物资名称不能为空", i+1))
				continue
			}

			// 构建Material对象
			m := Material{
				Name: name,
			}

			// 使用请求中的project_id（项目级别，不允许每行单独指定）
			pid := req.ProjectID
			m.ProjectID = &pid

			// 提取其他字段
			if code, ok := item["code"].(string); ok && code != "" {
				m.Code = &code
				// Check duplicate
				var existing Material
				if db.Where("code = ?", code).First(&existing).Error == nil {
					errorCount++
					errors = append(errors, fmt.Sprintf("第%d行: 物资编码已存在 - %s", i+1, code))
					continue
				}
			}

			if spec, ok := item["specification"].(string); ok {
				m.Specification = spec
			}
			if unit, ok := item["unit"].(string); ok {
				m.Unit = unit
			}
			if price, ok := item["price"].(float64); ok {
				m.Price = price
			} else if priceStr, ok := item["price"].(string); ok {
				if p, err := strconv.ParseFloat(priceStr, 64); err == nil {
					m.Price = p
				}
			}
			if desc, ok := item["description"].(string); ok {
				m.Description = desc
			}
			if cat, ok := item["category"].(string); ok {
				m.Category = cat
			}
			if qty, ok := item["quantity"].(float64); ok {
				m.Quantity = int(qty)
			} else if qtyStr, ok := item["quantity"].(string); ok {
				if q, err := strconv.Atoi(qtyStr); err == nil {
					m.Quantity = q
				}
			}
			if mat, ok := item["material"].(string); ok {
				m.Material = mat
			}
			if spec, ok := item["spec"].(string); ok {
				m.Spec = spec
			}

			if err := db.Create(&m).Error; err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("第%d行: %s", i+1, err.Error()))
			} else {
				successCount++
			}
		}

		response.SuccessWithMessage(c, map[string]any{
			"success_count": successCount,
			"error_count":   errorCount,
			"errors":        errors,
		}, "导入完成，成功"+strconv.Itoa(successCount)+"条，失败"+strconv.Itoa(errorCount)+"条")
	})

	// get material operation logs
	r.GET("/materials/:id/logs", auth.PermissionMiddleware(db, "material_view"), func(c *gin.Context) {
		id := c.Param("id")
		var m Material
		if err := db.First(&m, id).Error; err != nil {
			response.NotFound(c, "物资不存在")
			return
		}

		// Query logs from system logs table (if exists) or stock_logs
		type Log struct {
			ID          uint   `json:"id"`
			Action      string `json:"action"`
			Description string `json:"description"`
			CreatedAt   string `json:"created_at"`
			CreatedBy   string `json:"created_by"`
		}

		var logs []Log
		// Try to get from stock_logs related to this material
		db.Raw(`
			SELECT l.id, l.operation_type as action, l.remark as description, 
			       l.created_at, u.username as created_by
			FROM stock_logs l
			LEFT JOIN stocks s ON l.stock_id = s.id
			LEFT JOIN users u ON l.created_by = u.id
			WHERE s.material_id = ?
			ORDER BY l.created_at DESC
			LIMIT 100
		`, id).Scan(&logs)

		response.Success(c, map[string]any{
			"logs":  logs,
			"total": len(logs),
		})
	})

	// list unstored materials (materials NOT in any inbound order items - not yet received in warehouse)
	r.GET("/materials/unstored", auth.PermissionMiddleware(db, "material_view"), func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}
		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)

		var total int64
		query := db.Model(&Material{})

		// Filter: materials where arrived quantity < planned quantity
		query = query.Where(`
			NOT EXISTS (
				SELECT 1
				FROM inbound_order_items ioi
				INNER JOIN inbound_orders io ON ioi.order_id = io.id
				WHERE ioi.material_id = materials.id
				AND io.status IN ('approved', 'completed')
				GROUP BY ioi.material_id
				HAVING SUM(ioi.quantity) >= materials.quantity
			)
		`)

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
		var materials []Material
		query.Offset((page-1)*pageSize).Limit(pageSize).Order("id DESC").Find(&materials)

		out := make([]map[string]any, 0, len(materials))
		for _, m := range materials {
			out = append(out, m.ToDTO())
		}

		response.SuccessWithPagination(c, out, int64(page), int64(pageSize), total)
	})

	// export unstored materials to Excel
	r.GET("/materials/unstored/export", auth.PermissionMiddleware(db, "material_view"), func(c *gin.Context) {
		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)

		// Query materials where arrived quantity < planned quantity
		query := db.Model(&Material{}).
			Where(`
				NOT EXISTS (
					SELECT 1
					FROM inbound_order_items ioi
					INNER JOIN inbound_orders io ON ioi.order_id = io.id
					WHERE ioi.material_id = materials.id
					AND io.status IN ('approved', 'completed')
					GROUP BY ioi.material_id
					HAVING SUM(ioi.quantity) >= materials.quantity
				)
			`)

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
			}
		} else if projectID > 0 {
			// 管理员，应用指定的 project_id 过滤
			query = query.Where("project_id = ?", projectID)
		}
		query = query.Order("id DESC")

		var materials []Material
		query.Find(&materials)

		// Create Excel file
		f := excelize.NewFile()
		defer f.Close()
		
		sheet := "未入库物资"
		f.NewSheet(sheet)
		f.DeleteSheet("Sheet1") // Remove default sheet
		
		// Set column headers
		headers := []string{"ID", "编码", "物资名称", "规格", "材质", "单位", "价格", "分类", "数量", "描述", "项目ID"}
		for i, header := range headers {
			f.SetCellValue(sheet, fmt.Sprintf("%s1", string(rune('A'+i))), header)
		}
		
		// Set column widths
		colWidths := []float64{8, 12, 15, 15, 12, 8, 10, 10, 8, 20, 10}
		for i, width := range colWidths {
			f.SetColWidth(sheet, string(rune('A'+i)), string(rune('A'+i)), width)
		}
		
		// Fill data
		for idx, m := range materials {
			row := idx + 2
			spec := m.Specification
			if spec == "" && m.Spec != "" {
				spec = m.Spec
			}
			
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), m.ID)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), m.Code)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), m.Name)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), spec)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), m.Material)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), m.Unit)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), m.Price)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), m.Category)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), m.Quantity)
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), m.Description)
			f.SetCellValue(sheet, fmt.Sprintf("K%d", row), m.ProjectID)
		}
		
		// Generate Excel file in memory
		buffer, err := f.WriteToBuffer()
		if err != nil {
			response.InternalError(c, "生成Excel文件失败")
			return
		}
		
		c.Header("Content-Disposition", "attachment; filename=未入库物资导出.xlsx")
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
	})

	// ================== 批量导入物资 ==================
	r.POST("/materials/batch", auth.PermissionMiddleware(db, "material_import"), func(c *gin.Context) {
		var req struct {
			Materials []struct {
				Code             string   `json:"code"`
				Name             string   `json:"name" binding:"required"`
				Category         string   `json:"category"`
				Specification    string   `json:"specification"`
				Unit             string   `json:"unit" binding:"required"`
				Price            *float64 `json:"price"`
				Quantity         *float64 `json:"quantity"`
				QualityStandard  string   `json:"quality_standard"`
				Remark           string   `json:"remark"`
				ProjectID        uint      `json:"project_id" binding:"required"`
			} `json:"materials" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// 获取当前用户
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil || currentUser == nil {
			response.Unauthorized(c, "未授权")
			return
		}

		// 开启事务
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// 批量插入物资
		successCount := 0
		var failedItems []map[string]any

		for idx, mat := range req.Materials {
			// 组合描述字段
			description := mat.QualityStandard
			if mat.Remark != "" {
				if description != "" {
					description += "; " + mat.Remark
				} else {
					description = mat.Remark
				}
			}

			// 处理 Code 指针
			var codePtr *string
			if mat.Code != "" {
				codePtr = &mat.Code
			}

			// 处理 Price
			price := 0.0
			if mat.Price != nil {
				price = *mat.Price
			}

			// 处理 Quantity
			quantity := 0
			if mat.Quantity != nil {
				quantity = int(*mat.Quantity)
			}

			material := Material{
				Code:          codePtr,
				Name:          mat.Name,
				Category:      mat.Category,
				Specification: mat.Specification,
				Unit:          mat.Unit,
				Price:         price,
				Quantity:      quantity,
				Description:   description,
				ProjectID:     &mat.ProjectID,
			}

			if err := tx.Create(&material).Error; err != nil {
				failedItems = append(failedItems, map[string]any{
					"index": idx + 1,
					"name":  mat.Name,
					"error": err.Error(),
				})
				continue
			}

			successCount++
		}

		if len(failedItems) > 0 {
			tx.Rollback()
			response.Error(c, 400, fmt.Sprintf("部分导入失败：成功 %d 条，失败 %d 条", successCount, len(failedItems)))
			c.JSON(http.StatusOK, gin.H{
				"success":      false,
				"successCount": successCount,
				"failCount":    len(failedItems),
				"failedItems":  failedItems,
			})
			return
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			response.InternalError(c, "保存失败")
			return
		}

		response.SuccessWithMeta(c, map[string]any{
			"total":   successCount,
			"success": successCount,
			"failed":  0,
		}, nil)
	})

	// ================== 批量创建物资（用于计划导入） ==================
	r.POST("/materials/batch-create", auth.PermissionMiddleware(db, "material_import"), func(c *gin.Context) {
		var req struct {
			Materials []struct {
				Name          string   `json:"name" binding:"required"`
				Code          string   `json:"code"`
				Specification string   `json:"specification"`
				Category      string   `json:"category"`
				Unit          string   `json:"unit" binding:"required"`
				Price         *float64 `json:"price"`
				Quantity      *float64 `json:"quantity"`
				ProjectID     *uint    `json:"project_id"`
			} `json:"materials" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// 获取当前用户
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil || currentUser == nil {
			response.Unauthorized(c, "未授权")
			return
		}

		// 批量插入物资，如果已存在则使用现有物资
		var createdMaterials []map[string]any

		for _, mat := range req.Materials {
			// 首先检查物资是否已存在（按名称和规格查找）
			var existingMaterial Material
			err := db.Where("name = ? AND (specification = '' OR specification IS NULL OR specification = ?)",
				mat.Name, mat.Specification).
				First(&existingMaterial).Error

			var materialID uint
			var isNew bool

			if err == nil {
				// 物资已存在，使用现有ID
				materialID = existingMaterial.ID
				isNew = false
			} else {
				// 物资不存在，创建新物资
				// 处理 Code
				var codePtr *string
				if mat.Code != "" {
					codePtr = &mat.Code
				}

				// 处理 Price
				price := 0.0
				if mat.Price != nil {
					price = *mat.Price
				}

				// 处理 Quantity
				quantity := 0
				if mat.Quantity != nil {
					quantity = int(*mat.Quantity)
				}

				material := Material{
					Code:          codePtr,
					Name:          mat.Name,
					Category:      mat.Category,
					Specification: mat.Specification,
					Unit:          mat.Unit,
					Price:         price,
					Quantity:      quantity,
					ProjectID:     mat.ProjectID,
				}

				if err := db.Create(&material).Error; err != nil {
					response.InternalError(c, fmt.Sprintf("创建物资失败: %s - %v", mat.Name, err))
					return
				}

				materialID = material.ID
				isNew = true
			}

			// 返回物资信息（用于前端匹配）
			createdMaterials = append(createdMaterials, map[string]any{
				"id":           materialID,
				"name":         mat.Name,
				"specification": mat.Specification,
				"is_new":       isNew,
			})
		}

		response.Success(c, map[string]any{
			"success":  true,
			"materials": createdMaterials,
		})
	})
}