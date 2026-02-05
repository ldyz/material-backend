package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"gorm.io/gorm"
)

// RegisterReportRoutes 注册报表相关路由
func RegisterReportRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// 报表仪表盘
	r.GET("/reports/dashboard", auth.PermissionMiddleware(db, "system_statistics"), func(c *gin.Context) {
		now := time.Now()
		thirtyDaysAgo := now.AddDate(0, 0, -30)
		var materialCount, stockCount, inboundCount, inboundApproved, reqCount, reqApproved int64
		var totalStockValue float64
		db.Model(&struct{}{}).Table("material_master").Count(&materialCount)
		db.Model(&struct{}{}).Table("stocks").Count(&stockCount)
		db.Model(&struct{}{}).Table("stocks").Select("SUM(quantity * price)").Row().Scan(&totalStockValue)
		db.Model(&struct{}{}).Table("inbound_orders").Where("created_at >= ?", thirtyDaysAgo).Count(&inboundCount)
		db.Model(&struct{}{}).Table("inbound_orders").Where("status = ? AND created_at >= ?", "approved", thirtyDaysAgo).Count(&inboundApproved)
		db.Model(&struct{}{}).Table("requisitions").Where("created_at >= ?", thirtyDaysAgo).Count(&reqCount)
		db.Model(&struct{}{}).Table("requisitions").Where("status = ? AND created_at >= ?", "approved", thirtyDaysAgo).Count(&reqApproved)

		response.Success(c, map[string]interface{}{
			"material_count":       materialCount,
			"stock_count":          stockCount,
			"total_stock_value":    totalStockValue,
			"inbound_count":        inboundCount,
			"inbound_approved":     inboundApproved,
			"requisition_count":    reqCount,
			"requisition_approved": reqApproved,
			"period":               "最近30天",
		})
	})

	// 获取报表列表
	r.GET("/reports", auth.PermissionMiddleware(db, "system_statistics"), func(c *gin.Context) {
		reportDir := "./reports"
		if _, err := os.Stat(reportDir); os.IsNotExist(err) {
			response.SuccessWithMeta(c, []map[string]any{}, map[string]interface{}{
				"reports": []map[string]any{},
				"total":   0,
			})
			return
		}
		files, err := os.ReadDir(reportDir)
		if err != nil {
			response.InternalError(c, "读取报表目录失败")
			return
		}
		reports := []map[string]any{}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			info, _ := file.Info()
			reports = append(reports, map[string]any{
				"name":       file.Name(),
				"size":       info.Size(),
				"created_at": info.ModTime().Format("2006-01-02 15:04:05"),
			})
		}
		response.SuccessWithMeta(c, reports, map[string]interface{}{
			"total": len(reports),
		})
	})

	// 生成报表
	r.POST("/reports/generate", auth.PermissionMiddleware(db, "system_statistics"), func(c *gin.Context) {
		var req struct {
			Type      string `json:"type" binding:"required"`
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
			ProjectID uint   `json:"project_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误")
			return
		}
		reportDir := "./reports"
		if err := os.MkdirAll(reportDir, 0755); err != nil {
			response.InternalError(c, "创建报表目录失败")
			return
		}
		filename := fmt.Sprintf("report_%s_%s.json", req.Type, time.Now().Format("20060102150405"))
		filePath := filepath.Join(reportDir, filename)

		// Generate report data based on type
		var data any
		switch req.Type {
		case "material":
			var materials []map[string]any
			query := db.Model(&struct{}{}).Table("material_master")
			if req.ProjectID > 0 {
				query = query.Where("project_id = ?", req.ProjectID)
			}
			query.Find(&materials)
			data = materials
		case "stock":
			var stocks []map[string]any
			query := db.Model(&struct{}{}).Table("stocks")
			if req.ProjectID > 0 {
				query = query.Where("project_id = ?", req.ProjectID)
			}
			query.Find(&stocks)
			data = stocks
		default:
			data = map[string]any{"message": "Report generated"}
		}

		// Write to file
		file, err := os.Create(filePath)
		if err != nil {
			response.InternalError(c, "创建报表文件失败")
			return
		}
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.Encode(data)

		response.SuccessWithMessage(c, map[string]string{
			"filename": filename,
		}, "报表生成成功")
	})

	// 下载报表(query参数)
	r.GET("/reports/download", auth.PermissionMiddleware(db, "system_statistics"), func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			response.BadRequest(c, "文件名不能为空")
			return
		}
		filePath := filepath.Join("./reports", filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			response.NotFound(c, "报表文件不存在")
			return
		}
		c.FileAttachment(filePath, filename)
	})

	// 下载指定报表(路径参数)
	r.GET("/reports/:report_name/download", auth.PermissionMiddleware(db, "system_statistics"), func(c *gin.Context) {
		filename := c.Param("report_name")
		filePath := filepath.Join("./reports", filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			response.NotFound(c, "报表文件不存在")
			return
		}
		c.FileAttachment(filePath, filename)
	})

	// 删除报表(body参数)
	r.POST("/reports/delete", auth.PermissionMiddleware(db, "system_backup"), func(c *gin.Context) {
		var req struct {
			Filename string `json:"filename" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "文件名不能为空")
			return
		}
		filePath := filepath.Join("./reports", req.Filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			response.NotFound(c, "报表文件不存在")
			return
		}
		if err := os.Remove(filePath); err != nil {
			response.InternalError(c, "删除报表失败")
			return
		}
		response.SuccessOnlyMessage(c, "报表删除成功")
	})

	// 删除指定报表(路径参数)
	r.DELETE("/reports/:report_name", auth.PermissionMiddleware(db, "system_backup"), func(c *gin.Context) {
		filename := c.Param("report_name")
		filePath := filepath.Join("./reports", filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			response.NotFound(c, "报表文件不存在")
			return
		}
		if err := os.Remove(filePath); err != nil {
			response.InternalError(c, "删除报表失败")
			return
		}
		response.SuccessOnlyMessage(c, "报表删除成功")
	})
}
