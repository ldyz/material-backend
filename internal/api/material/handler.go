package material

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("material")
	r.Use(jwtpkg.TokenMiddleware())

	service := NewService(db)

	// Material CRUD
	r.GET("/materials", auth.PermissionMiddleware(db, "material_view"), listMaterials(service))
	r.POST("/materials", auth.PermissionMiddleware(db, "material_create"), createMaterial(service))
	r.GET("/materials/:id", auth.PermissionMiddleware(db, "material_view"), getMaterial(service))
	r.PUT("/materials/:id", auth.PermissionMiddleware(db, "material_edit"), updateMaterial(service))
	r.DELETE("/materials/:id", auth.PermissionMiddleware(db, "material_delete"), deleteMaterial(service))

	// Material operations
	r.GET("/materials/export", auth.PermissionMiddleware(db, "material_view"), exportMaterials(service))
	r.POST("/material/materials/import", auth.PermissionMiddleware(db, "material_create"), importMaterials(service))
	r.GET("/materials/:id/logs", auth.PermissionMiddleware(db, "material_view"), getMaterialLogs(service))
	r.GET("/materials/unstored", auth.PermissionMiddleware(db, "material_view"), listUnstoredMaterials(service))
	r.GET("/materials/unstored/export", auth.PermissionMiddleware(db, "material_view"), exportUnstoredMaterials(service))

	// Batch operations
	r.POST("/materials/batch", auth.PermissionMiddleware(db, "material_import"), batchMaterials(service))
	r.POST("/materials/batch-create", auth.PermissionMiddleware(db, "material_import"), batchCreateMaterials(service))
}

// listMaterials lists materials with filters and pagination
func listMaterials(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)
		includeChildren := c.Query("include_children") == "true"

		params := ListQueryParams{
			Page:          page,
			PageSize:      pageSize,
			Search:        c.Query("search"),
			Name:          c.Query("name"),
			Material:      c.Query("material"),
			Spec:          c.Query("spec"),
			Specification: c.Query("specification"),
			Category:      c.Query("category"),
			Filter:        c.Query("filter"),
			IncludeChildren: includeChildren,
		}

		// Handle project filtering
		projectIDsStr := c.Query("project_ids")
		if projectIDsStr != "" {
			for _, idStr := range strings.Split(projectIDsStr, ",") {
				if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64); err == nil {
					params.ProjectIDs = append(params.ProjectIDs, uint(id))
				}
			}
		} else if projectID > 0 && includeChildren {
			params.ProjectIDs = append(params.ProjectIDs, uint(projectID))
			childIDs := service.GetChildProjectIDs(uint(projectID))
			params.ProjectIDs = append(params.ProjectIDs, childIDs...)
		} else if projectID > 0 {
			params.ProjectID = uint(projectID)
		}

		// Apply project access control
		accessibleProjectIDs, err := auth.GetAccessibleProjectIDs(c, service.db)
		if err != nil {
			response.InternalError(c, "获取用户项目权限失败")
			return
		}

		if accessibleProjectIDs != nil {
			if len(accessibleProjectIDs) == 0 {
				response.SuccessWithPagination(c, []map[string]any{}, int64(page), int64(pageSize), 0)
				return
			}

			if len(params.ProjectIDs) > 0 {
				// Verify all requested project IDs are accessible
				for _, pid := range params.ProjectIDs {
					hasAccess := false
					for _, accessibleID := range accessibleProjectIDs {
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
			} else if params.ProjectID > 0 {
				hasAccess := false
				for _, pid := range accessibleProjectIDs {
					if params.ProjectID == pid {
						hasAccess = true
						break
					}
				}
				if !hasAccess {
					response.Forbidden(c, "无权访问该项目")
					return
				}
			} else {
				params.ProjectIDs = accessibleProjectIDs
			}
		}

		results, total, err := service.ListMaterials(params)
		if err != nil {
			response.InternalError(c, "获取物资列表失败")
			return
		}

		// Enrich with plan info
		planMap := service.EnrichWithPlanInfo(results)

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

			planInfo := planMap[m.ID]
			plannedQty := 0
			arrivedQty := 0
			if planInfo.MaterialID > 0 {
				plannedQty = planInfo.PlannedQuantity
				arrivedQty = planInfo.ArrivedQuantity
			}

			remainingQty := plannedQty - arrivedQty
			arrivalPercentage := 0.0
			if plannedQty > 0 {
				arrivalPercentage = float64(arrivedQty) / float64(plannedQty) * 100
			}

			out = append(out, map[string]any{
				"id":                 m.ID,
				"code":               m.Code,
				"name":               m.Name,
				"specification":      spec,
				"unit":               m.Unit,
				"price":              m.Price,
				"description":        m.Description,
				"category":           m.Category,
				"quantity":           m.StockQuantity,
				"planned_quantity":   plannedQty,
				"arrived_quantity":   arrivedQty,
				"remaining_quantity": remainingQty,
				"arrival_percentage": arrivalPercentage,
				"is_fully_arrived":   arrivedQty >= plannedQty && plannedQty > 0,
				"project_id":         projectID,
				"project_name":       m.ProjectName,
				"material":           m.Material,
				"spec":               m.Spec,
			})
		}

		response.SuccessWithPagination(c, out, int64(page), int64(pageSize), total)
	}
}

// createMaterial creates a new material
func createMaterial(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Code          string  `json:"code"`
			Name          string  `json:"name"`
			Specification string  `json:"specification"`
			Unit          string  `json:"unit"`
			Price         float64 `json:"price"`
			Description   string  `json:"description"`
			Category      string  `json:"category"`
			Quantity      int     `json:"quantity"`
			ProjectID     string  `json:"project_id"`
			Material      string  `json:"material"`
			Spec          string  `json:"spec"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
		if req.Name == "" {
			response.BadRequest(c, "物资名称不能为空")
			return
		}

		if req.ProjectID == "" {
			response.BadRequest(c, "请选择所属项目")
			return
		}

		projectIDUint, err := strconv.ParseUint(req.ProjectID, 10, 64)
		if err != nil || projectIDUint == 0 {
			response.BadRequest(c, "项目ID格式无效")
			return
		}

		createReq := &CreateMaterialRequest{
			Code:          req.Code,
			Name:          req.Name,
			Specification: req.Specification,
			Unit:          req.Unit,
			Price:         req.Price,
			Description:   req.Description,
			Category:      req.Category,
			Quantity:      req.Quantity,
			ProjectID:     uint(projectIDUint),
			Material:      req.Material,
			Spec:          req.Spec,
		}

		m, err := service.CreateMaterial(createReq)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Created(c, m.ToDTO(), "物资创建成功")
	}
}

// getMaterial gets a single material
func getMaterial(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的物资ID")
			return
		}

		m, err := service.GetMaterial(uint(idUint))
		if err != nil {
			response.NotFound(c, err.Error())
			return
		}

		response.Success(c, m.ToDTO())
	}
}

// updateMaterial updates a material
func updateMaterial(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的物资ID")
			return
		}

		var req map[string]any
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		m, err := service.UpdateMaterial(uint(idUint), req)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, m.ToDTO(), "物资更新成功")
	}
}

// deleteMaterial deletes a material
func deleteMaterial(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的物资ID")
			return
		}

		if err := service.DeleteMaterial(uint(idUint)); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, nil, "物资删除成功")
	}
}

// exportMaterials exports materials to Excel
func exportMaterials(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")
		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)

		var projectIDs []uint
		accessibleProjectIDs, err := auth.GetAccessibleProjectIDs(c, service.db)
		if err != nil {
			response.InternalError(c, "获取用户项目权限失败")
			return
		}

		if accessibleProjectIDs != nil {
			if len(accessibleProjectIDs) == 0 {
				projectIDs = []uint{}
			} else {
				if projectID == 0 {
					projectIDs = accessibleProjectIDs
				} else {
					hasAccess := false
					for _, pid := range accessibleProjectIDs {
						if uint(projectID) == pid {
							hasAccess = true
							break
						}
					}
					if !hasAccess {
						response.Forbidden(c, "无权访问该项目")
						return
					}
					projectIDs = []uint{uint(projectID)}
				}
			}
		} else if projectID > 0 {
			projectIDs = []uint{uint(projectID)}
		}

		materials, err := service.ExportMaterials(search, projectIDs)
		if err != nil {
			response.InternalError(c, "获取物资列表失败")
			return
		}

		// Create Excel file
		f := excelize.NewFile()
		defer f.Close()

		sheet := "物资列表"
		f.NewSheet(sheet)
		f.DeleteSheet("Sheet1")

		// Set headers
		headers := []string{"ID", "编码", "物资名称", "规格", "材质", "单位", "价格", "分类", "数量", "描述", "项目ID"}
		for i, header := range headers {
			cell := fmt.Sprintf("%s1", string(rune('A'+i)))
			f.SetCellValue(sheet, cell, header)
		}

		// Set column widths
		colWidths := []float64{8, 12, 15, 15, 12, 8, 10, 10, 8, 20, 10}
		for i, width := range colWidths {
			col := string(rune('A' + i))
			f.SetColWidth(sheet, col, col, width)
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

		buffer, err := f.WriteToBuffer()
		if err != nil {
			response.InternalError(c, "生成Excel文件失败")
			return
		}

		c.Header("Content-Disposition", "attachment; filename=物资导出.xlsx")
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
	}
}

// importMaterials imports materials from JSON
func importMaterials(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Materials []map[string]any `json:"materials" binding:"required"`
			ProjectID uint            `json:"project_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误: "+err.Error())
			return
		}

		importItems := make([]ImportMaterial, 0, len(req.Materials))
		for _, item := range req.Materials {
			name, _ := item["name"].(string)
			if name == "" {
				if val, ok := item["Name"]; ok {
					name, _ = val.(string)
				}
			}

			if name == "" {
				continue
			}

			importItem := ImportMaterial{Name: name}

			if code, ok := item["code"].(string); ok {
				importItem.Code = code
			}
			if spec, ok := item["specification"].(string); ok {
				importItem.Specification = spec
			}
			if unit, ok := item["unit"].(string); ok {
				importItem.Unit = unit
			}
			if price, ok := item["price"].(float64); ok {
				importItem.Price = price
			} else if priceStr, ok := item["price"].(string); ok {
				if p, err := strconv.ParseFloat(priceStr, 64); err == nil {
					importItem.Price = p
				}
			}
			if desc, ok := item["description"].(string); ok {
				importItem.Description = desc
			}
			if cat, ok := item["category"].(string); ok {
				importItem.Category = cat
			}
			if qty, ok := item["quantity"].(float64); ok {
				importItem.Quantity = int(qty)
			} else if qtyStr, ok := item["quantity"].(string); ok {
				if q, err := strconv.Atoi(qtyStr); err == nil {
					importItem.Quantity = q
				}
			}
			if mat, ok := item["material"].(string); ok {
				importItem.Material = mat
			}
			if spec, ok := item["spec"].(string); ok {
				importItem.Spec = spec
			}

			importItems = append(importItems, importItem)
		}

		result, err := service.ImportMaterials(req.ProjectID, importItems)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.SuccessWithMessage(c, map[string]any{
			"success_count": result.SuccessCount,
			"error_count":   result.ErrorCount,
			"errors":        result.Errors,
		}, fmt.Sprintf("导入完成，成功%d条，失败%d条", result.SuccessCount, result.ErrorCount))
	}
}

// getMaterialLogs gets operation logs for a material
func getMaterialLogs(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的物资ID")
			return
		}

		logs, err := service.GetMaterialLogs(uint(idUint))
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		response.Success(c, map[string]any{
			"logs":  logs,
			"total": len(logs),
		})
	}
}

// listUnstoredMaterials lists materials that are not fully received
func listUnstoredMaterials(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)

		params := ListQueryParams{
			Page:    page,
			PageSize: pageSize,
			Filter:  "unstored",
		}

		// Apply project access control
		accessibleProjectIDs, err := auth.GetAccessibleProjectIDs(c, service.db)
		if err != nil {
			response.InternalError(c, "获取用户项目权限失败")
			return
		}

		if accessibleProjectIDs != nil {
			if len(accessibleProjectIDs) == 0 {
				response.SuccessWithPagination(c, []map[string]any{}, int64(page), int64(pageSize), 0)
				return
			}

			if projectID == 0 {
				params.ProjectIDs = accessibleProjectIDs
			} else {
				hasAccess := false
				for _, pid := range accessibleProjectIDs {
					if uint(projectID) == pid {
						hasAccess = true
						break
					}
				}
				if !hasAccess {
					response.Forbidden(c, "无权访问该项目")
					return
				}
				params.ProjectID = uint(projectID)
			}
		} else if projectID > 0 {
			params.ProjectID = uint(projectID)
		}

		results, total, err := service.ListMaterials(params)
		if err != nil {
			response.InternalError(c, "获取未入库物资列表失败")
			return
		}

		out := make([]map[string]any, 0, len(results))
		for _, m := range results {
			out = append(out, map[string]any{
				"id":            m.ID,
				"code":          m.Code,
				"name":          m.Name,
				"specification": m.Specification,
				"unit":          m.Unit,
				"price":         m.Price,
				"description":   m.Description,
				"category":      m.Category,
				"quantity":      m.Quantity,
				"project_id":    m.ProjectID,
				"material":      m.Material,
				"spec":          m.Spec,
			})
		}

		response.SuccessWithPagination(c, out, int64(page), int64(pageSize), total)
	}
}

// exportUnstoredMaterials exports unstored materials to Excel
func exportUnstoredMaterials(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 64)

		// Apply project access control
		accessibleProjectIDs, err := auth.GetAccessibleProjectIDs(c, service.db)
		if err != nil {
			response.InternalError(c, "获取用户项目权限失败")
			return
		}

		var projectIDs []uint
		if accessibleProjectIDs != nil {
			if len(accessibleProjectIDs) == 0 {
				projectIDs = []uint{}
			} else {
				if projectID == 0 {
					projectIDs = accessibleProjectIDs
				} else {
					hasAccess := false
					for _, pid := range accessibleProjectIDs {
						if uint(projectID) == pid {
							hasAccess = true
							break
						}
					}
					if !hasAccess {
						response.Forbidden(c, "无权访问该项目")
						return
					}
					projectIDs = []uint{uint(projectID)}
				}
			}
		} else if projectID > 0 {
			projectIDs = []uint{uint(projectID)}
		}

		params := ListQueryParams{
			Filter:     "unstored",
			ProjectIDs: projectIDs,
		}

		results, _, err := service.ListMaterials(params)
		if err != nil {
			response.InternalError(c, "获取未入库物资列表失败")
			return
		}

		// Create Excel file
		f := excelize.NewFile()
		defer f.Close()

		sheet := "未入库物资"
		f.NewSheet(sheet)
		f.DeleteSheet("Sheet1")

		// Set headers
		headers := []string{"ID", "编码", "物资名称", "规格", "材质", "单位", "价格", "分类", "数量", "描述", "项目ID"}
		for i, header := range headers {
			cell := fmt.Sprintf("%s1", string(rune('A'+i)))
			f.SetCellValue(sheet, cell, header)
		}

		// Set column widths
		colWidths := []float64{8, 12, 15, 15, 12, 8, 10, 10, 8, 20, 10}
		for i, width := range colWidths {
			col := string(rune('A' + i))
			f.SetColWidth(sheet, col, col, width)
		}

		// Fill data
		for idx, m := range results {
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

		buffer, err := f.WriteToBuffer()
		if err != nil {
			response.InternalError(c, "生成Excel文件失败")
			return
		}

		c.Header("Content-Disposition", "attachment; filename=未入库物资导出.xlsx")
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
	}
}

// batchMaterials creates materials in batch
func batchMaterials(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Materials []struct {
				Code            string   `json:"code"`
				Name            string   `json:"name" binding:"required"`
				Category        string   `json:"category"`
				Specification   string   `json:"specification"`
				Unit            string   `json:"unit" binding:"required"`
				Price           *float64 `json:"price"`
				Quantity        *float64 `json:"quantity"`
				QualityStandard string   `json:"quality_standard"`
				Remark          string   `json:"remark"`
				ProjectID       uint     `json:"project_id" binding:"required"`
			} `json:"materials" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		tx := service.db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		var failedItems []map[string]any
		successCount := 0

		for _, mat := range req.Materials {
			description := mat.QualityStandard
			if mat.Remark != "" {
				if description != "" {
					description += "; " + mat.Remark
				} else {
					description = mat.Remark
				}
			}

			price := 0.0
			if mat.Price != nil {
				price = *mat.Price
			}

			quantity := 0
			if mat.Quantity != nil {
				quantity = int(*mat.Quantity)
			}

			var codePtr *string
			if mat.Code != "" {
				codePtr = &mat.Code
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
					"name":  mat.Name,
					"error": err.Error(),
				})
				continue
			}

			successCount++
		}

		if len(failedItems) > 0 {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{
				"success":      false,
				"successCount": successCount,
				"failCount":    len(failedItems),
				"failedItems":  failedItems,
			})
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.InternalError(c, "保存失败")
			return
		}

		response.SuccessWithMeta(c, map[string]any{
			"total":   successCount,
			"success": successCount,
			"failed":  0,
		}, nil)
	}
}

// batchCreateMaterials creates or finds materials in batch
func batchCreateMaterials(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		materials := make([]BatchCreateMaterial, 0, len(req.Materials))
		for _, mat := range req.Materials {
			materials = append(materials, BatchCreateMaterial{
				Name:          mat.Name,
				Code:          mat.Code,
				Specification: mat.Specification,
				Category:      mat.Category,
				Unit:          mat.Unit,
				Price:         mat.Price,
				Quantity:      mat.Quantity,
				ProjectID:     mat.ProjectID,
			})
		}

		result, err := service.BatchCreateMaterials(materials)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, map[string]any{
			"success":   true,
			"materials": result.Materials,
		})
	}
}
