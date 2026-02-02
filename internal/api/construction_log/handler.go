package construction_log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/request"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// RegisterRoutes 注册施工日志模块路由
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// 创建路由组
	r := rg.Group("/construction_log")
	// 使用JWT中间件进行身份验证
	r.Use(jwtpkg.TokenMiddleware())

	// 配置文件上传
	uploadFolder := os.Getenv("CONSTRUCTION_LOG_UPLOAD")
	if uploadFolder == "" {
		uploadFolder = "static/uploads/construction_log"
	}
	// 确保上传目录存在
	os.MkdirAll(uploadFolder, 0755)

	// 允许的文件扩展名
	allowedExtensions := map[string]bool{
		"png":  true,
		"jpg":  true,
		"jpeg": true,
		"gif":  true,
		"bmp":  true,
	}

	// 检查文件是否允许上传
	allowedFile := func(filename string) bool {
		ext := strings.ToLower(filepath.Ext(filename))
		if ext == "" || !strings.HasPrefix(ext, ".") {
			return false
		}
		ext = ext[1:] // 去掉点号
		_, ok := allowedExtensions[ext]
		return ok
	}

	// ================== 日志列表接口 ==================
	r.GET("/logs", auth.PermissionMiddleware(db, "constructionlog_view"), func(c *gin.Context) {
		var req request.PaginationRequest
		if err := request.BindQuery(c, &req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		page := req.GetPage()
		perPage := req.GetPageSize()
		offset := req.GetOffset()

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

		// 兼容旧的单项目ID参数
		projectID := c.Query("project_id")
		if projectID != "" && len(projectIDsFilter) == 0 {
			if pid, err := strconv.ParseUint(projectID, 10, 64); err == nil {
				projectIDsFilter = append(projectIDsFilter, uint(pid))
			}
		}

		startDate := c.Query("start_date")
		endDate := c.Query("end_date")
		search := strings.TrimSpace(c.Query("search"))

		// 构建查询
		query := db.Model(&ConstructionLog{})

		// 权限过滤：只查用户有权限的项目
		// 暂时移除用户项目权限过滤，仅保留基本权限检查
		// TODO: 实现用户项目权限过滤

		// 项目ID过滤（支持多个）
		if len(projectIDsFilter) > 0 {
			query = query.Where("project_id IN ?", projectIDsFilter)
		}

		// 日期范围过滤
		if startDate != "" {
			if dt, err := time.Parse("2006-01-02", startDate); err == nil {
				query = query.Where("log_date >= ?", dt)
			}
		}
		if endDate != "" {
			if dt, err := time.Parse("2006-01-02", endDate); err == nil {
				query = query.Where("log_date <= ?", dt)
			}
		}

		// 搜索过滤（标题和内容）
		if search != "" {
			query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
		}

		// 获取总数
		var total int64
		query.Count(&total)

		// 查询日志
		var logs []ConstructionLog
		query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&logs)

		// 补充项目名称
		logDicts := make([]map[string]any, len(logs))
		for i, log := range logs {
			d := log.ToDict()
			// 获取项目名称
			var projectName string
			db.Raw(`SELECT name FROM projects WHERE id = ?`, log.ProjectID).Scan(&projectName)
			if projectName == "" {
				projectName = "未知项目"
			}
			d["project_name"] = projectName
			logDicts[i] = d
		}

		response.SuccessWithPagination(c, logDicts, int64(page), int64(perPage), total)
	})

	// ================== 日志详情接口 ==================
	r.GET("/:log_id", auth.PermissionMiddleware(db, "constructionlog_view"), func(c *gin.Context) {
		var uriReq struct {
			LogID int `uri:"log_id" binding:"required,min=1"`
		}
		if err := request.BindURI(c, &uriReq); err != nil {
			response.BadRequest(c, "无效的日志ID")
			return
		}

		// 构建查询
		query := db.Model(&ConstructionLog{})

		// 权限过滤
		// 暂时移除用户项目权限过滤，仅保留基本权限检查
		// TODO: 实现用户项目权限过滤

		// 查询日志
		var log ConstructionLog
		if err := query.Where("id = ?", uriReq.LogID).First(&log).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.NotFound(c, "日志不存在")
			} else {
				response.InternalError(c, "查询失败")
			}
			return
		}

		response.Success(c, log.ToDict())
	})

	// ================== 创建日志接口 ==================
	r.POST("/", auth.PermissionMiddleware(db, "constructionlog_create"), func(c *gin.Context) {
		// 获取当前用户
		currentUser, err := auth.GetCurrentUser(c, db)
		if err != nil || currentUser == nil {
			response.Unauthorized(c, "未授权")
			return
		}

		// 绑定请求数据
		var req struct {
			ProjectID uint   `json:"project_id" binding:"required"`
			Title     string `json:"title"`
			Content   string `json:"content"`
			Images    string `json:"images"`
			Weather   string `json:"weather"`
		}

		if err := request.BindJSON(c, &req); err != nil {
			response.BadRequest(c, "无效的请求数据")
			return
		}

		// 检查项目ID是否有效
		if req.ProjectID <= 0 {
			response.BadRequest(c, "项目ID无效")
			return
		}

		// 权限检查：确保用户有权限访问该项目
		// 暂时移除用户项目权限过滤，仅保留基本权限检查
		// TODO: 实现用户项目权限过滤

		// 设置默认标题
		if req.Title == "" {
			req.Title = "无标题"
		}

		// 创建日志
		log := ConstructionLog{
			Title:     req.Title,
			Content:   req.Content,
			Images:    req.Images,
			Weather:   req.Weather,
			ProjectID: req.ProjectID,
			CreatorID: currentUser.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// 保存到数据库
		if err := db.Create(&log).Error; err != nil {
			response.InternalError(c, "创建日志失败")
			return
		}

		response.SuccessWithMeta(c, map[string]uint{"id": log.ID}, nil)
	})

	// ================== 编辑日志接口 ==================
	r.PUT("/:log_id", auth.PermissionMiddleware(db, "constructionlog_edit"), func(c *gin.Context) {
		var uriReq struct {
			LogID int `uri:"log_id" binding:"required,min=1"`
		}
		if err := request.BindURI(c, &uriReq); err != nil {
			response.BadRequest(c, "无效的日志ID")
			return
		}

		// 绑定请求数据
		var req struct {
			ProjectID   *int     `json:"project_id"`
			Title       *string  `json:"title"`
			LogDate     *string  `json:"log_date"`
			Weather     *string  `json:"weather"`
			Temperature *float64 `json:"temperature"`
			Content     *string  `json:"content"`
			Progress    *string  `json:"progress"`
			Issues      *string  `json:"issues"`
			Images      *string  `json:"images"`
			Remark      *string  `json:"remark"`
		}

		if err := request.BindJSON(c, &req); err != nil {
			response.BadRequest(c, "无效的请求数据")
			return
		}

		// 构建查询
		query := db.Model(&ConstructionLog{})

		// 查询日志
		var log ConstructionLog
		if err := query.Where("id = ?", uriReq.LogID).First(&log).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.NotFound(c, "日志不存在")
			} else {
				response.InternalError(c, "查询失败")
			}
			return
		}

		// 更新日志
		updates := map[string]any{}
		if req.ProjectID != nil {
			updates["project_id"] = *req.ProjectID
		}
		if req.Title != nil {
			updates["title"] = *req.Title
		}
		if req.LogDate != nil {
			updates["log_date"] = *req.LogDate
		}
		if req.Weather != nil {
			updates["weather"] = *req.Weather
		}
		if req.Temperature != nil {
			updates["temperature"] = *req.Temperature
		}
		if req.Content != nil {
			updates["content"] = *req.Content
		}
		if req.Progress != nil {
			updates["progress"] = *req.Progress
		}
		if req.Issues != nil {
			updates["issues"] = *req.Issues
		}
		if req.Images != nil {
			updates["images"] = *req.Images
		}
		if req.Remark != nil {
			updates["remark"] = *req.Remark
		}
		if len(updates) > 0 {
			updates["updated_at"] = time.Now()
			if err := db.Model(&log).Updates(updates).Error; err != nil {
				response.InternalError(c, "更新失败")
				return
			}
		}

		response.SuccessOnlyMessage(c, "更新成功")
	})

	// ================== 删除日志接口 ==================
	r.DELETE("/:log_id", auth.PermissionMiddleware(db, "constructionlog_delete"), func(c *gin.Context) {
		var uriReq struct {
			LogID int `uri:"log_id" binding:"required,min=1"`
		}
		if err := request.BindURI(c, &uriReq); err != nil {
			response.BadRequest(c, "无效的日志ID")
			return
		}

		// 构建查询
		query := db.Model(&ConstructionLog{})

		// 权限过滤
		// 暂时移除用户项目权限过滤，仅保留基本权限检查
		// TODO: 实现用户项目权限过滤

		// 查询日志
		var log ConstructionLog
		if err := query.Where("id = ?", uriReq.LogID).First(&log).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.NotFound(c, "日志不存在")
			} else {
				response.InternalError(c, "查询失败")
			}
			return
		}

		// 删除日志
		if err := db.Delete(&log).Error; err != nil {
			response.InternalError(c, "删除日志失败")
			return
		}

		response.SuccessOnlyMessage(c, "删除成功")
	})

	// ================== 图片上传接口 ==================
	r.POST("/upload_image", auth.PermissionMiddleware(db, "constructionlog_create"), func(c *gin.Context) {
		// 获取上传文件
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			response.BadRequest(c, "未选择文件")
			return
		}
		defer file.Close()

		// 检查文件名
		if header.Filename == "" {
			response.BadRequest(c, "未选择文件")
			return
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if !allowedFile(header.Filename) {
			response.BadRequest(c, "文件类型不支持")
			return
		}

		// 生成唯一文件名（使用当前时间戳+文件扩展名）
		timestamp := time.Now().UnixNano() / int64(time.Millisecond)
		filename := fmt.Sprintf("%d%s", timestamp, ext)
		filePath := filepath.Join(uploadFolder, filename)

		// 创建目标文件
		dst, err := os.Create(filePath)
		if err != nil {
			response.InternalError(c, "创建文件失败")
			return
		}
		defer dst.Close()

		// 复制文件内容
		if _, err := io.Copy(dst, file); err != nil {
			response.InternalError(c, "保存文件失败")
			return
		}

		// 生成文件URL
		url := fmt.Sprintf("/static/uploads/construction_log/%s", filename)

		response.SuccessWithMeta(c, map[string]string{"url": url}, nil)
	})
}
