package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/agent"
	"github.com/yourorg/material-backend/backend/internal/api/app"
	"github.com/yourorg/material-backend/backend/internal/api/appointment"
	"github.com/yourorg/material-backend/backend/internal/api/attendance"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/construction_log"
	"github.com/yourorg/material-backend/backend/internal/api/inbound"
	audit2 "github.com/yourorg/material-backend/backend/internal/api/audit"
	"github.com/yourorg/material-backend/backend/internal/api/material"
	"github.com/yourorg/material-backend/backend/internal/api/material_master"
	"github.com/yourorg/material-backend/backend/internal/api/material_plan"
	"github.com/yourorg/material-backend/backend/internal/api/notification"
	appconfig "github.com/yourorg/material-backend/backend/internal/config"
	"github.com/yourorg/material-backend/backend/internal/api/project"
	"github.com/yourorg/material-backend/backend/internal/api/progress"
	"github.com/yourorg/material-backend/backend/internal/api/requisition"
	"github.com/yourorg/material-backend/backend/internal/api/stock"
	"github.com/yourorg/material-backend/backend/internal/api/system"
	"github.com/yourorg/material-backend/backend/internal/api/upload"
	"github.com/yourorg/material-backend/backend/internal/api/workflow"
	"github.com/yourorg/material-backend/backend/internal/db"
	"github.com/yourorg/material-backend/backend/internal/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// isMobileDevice 检测是否为移动端设备
func isMobileDevice(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)
	mobileKeywords := []string{
		"mobile", "android", "iphone", "ipad", "ipod",
		"windows phone", "blackberry", "webos", "opera mini",
		"opera mobi", "symbian", "meego", "kindle",
		"capacitor", // Capacitor 应用
	}

	for _, keyword := range mobileKeywords {
		if strings.Contains(userAgent, keyword) {
			return true
		}
	}

	// 检查屏幕宽度提示（部分浏览器会发送）
	if strings.Contains(userAgent, "mobi") {
		return true
	}

	return false
}

// initializeUploadConfig 初始化上传配置
func initializeUploadConfig(dbConn *gorm.DB) {
	defaultConfigs := []system.SystemConfig{
		{Key: "upload_directory", Value: "static/uploads", Type: "string", Description: "文件上传目录路径"},
		{Key: "max_file_size", Value: "5", Type: "integer", Description: "最大文件上传大小(MB)"},
		{Key: "allowed_file_types", Value: "jpg,jpeg,png,gif,bmp,webp,svg", Type: "string", Description: "允许上传的文件类型"},
		{Key: "max_upload_count", Value: "10", Type: "integer", Description: "单次最多上传文件数量"},
	}

	for _, config := range defaultConfigs {
		var count int64
		dbConn.Model(&system.SystemConfig{}).Where("key = ?", config.Key).Count(&count)
		if count == 0 {
			if err := dbConn.Create(&config).Error; err != nil {
				log.Printf("初始化配置失败 [%s]: %v", config.Key, err)
			} else {
				log.Printf("初始化配置成功: %s = %s", config.Key, config.Value)
			}
		}
	}
}

// initializeMaterialCategories 初始化默认物资分类
func initializeMaterialCategories(dbConn *gorm.DB) {
	defaultCategories := []struct {
		name  string
		code  string
		sort  int
		remark string
	}{
		{"钢材", "STEEL", 1, "各种钢材材料，如钢筋、钢板、钢管等"},
		{"水泥", "CEMENT", 2, "各类水泥及水泥制品"},
		{"砂石", "SAND_STONE", 3, "砂、石、砂石等骨料"},
		{"电气材料", "ELECTRICAL", 4, "电线、电缆、开关、插座等电气材料"},
		{"管道材料", "PIPE", 5, "各类管道及管件"},
		{"木材", "WOOD", 6, "原木、板材、木方等木材"},
		{"涂料", "PAINT", 7, "各类油漆、涂料"},
		{"保温材料", "INSULATION", 8, "保温棉、泡沫板等保温材料"},
		{"防水材料", "WATERPROOF", 9, "防水卷材、防水涂料等"},
		{"五金配件", "HARDWARE", 10, "螺丝、螺母、钉子等五金配件"},
		{"劳保用品", "SAFETY", 11, "安全帽、手套、工作服等劳保用品"},
		{"工具", "TOOLS", 12, "电动工具、手动工具等"},
		{"其他", "OTHER", 99, "其他物资"},
	}

	for _, cat := range defaultCategories {
		var count int64
		dbConn.Table("material_categories").Where("name = ?", cat.name).Count(&count)
		if count == 0 {
			category := map[string]any{
				"name":   cat.name,
				"code":   cat.code,
				"sort":   cat.sort,
				"remark": cat.remark,
			}
			if err := dbConn.Table("material_categories").Create(&category).Error; err != nil {
				log.Printf("初始化物资分类失败 [%s]: %v", cat.name, err)
			} else {
				log.Printf("初始化物资分类成功: %s", cat.name)
			}
		}
	}
}

func main() {
	// 定义命令行参数
	configFile := flag.String("c", "config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置文件
	cfg, err := appconfig.Load(*configFile)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	log.Printf("配置文件加载成功: %s", *configFile)
	log.Printf("数据库类型: %s, 主机: %s:%d, 数据库: %s",
		cfg.Database.Type, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)

	// 使用配置文件中的数据库连接
	dbConn, err := db.New(cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	// Sync PostgreSQL sequences (fixes out-of-sync sequences)
	if err := db.SyncPostgreSQLSequences(dbConn); err != nil {
		log.Printf("Warning: Failed to sync PostgreSQL sequences: %v", err)
	}

	// Auto migrate auth, project, material and inbound models
	dbConn.AutoMigrate(
		&auth.User{}, &auth.Role{},
		&project.Project{},
		&material.Material{}, &material.MaterialCategory{},
		&material_master.MaterialMaster{}, // 物资主数据表 (v2)
		&inbound.InboundOrder{}, &inbound.InboundOrderItem{},
		&requisition.Requisition{}, &requisition.RequisitionItem{},
		&stock.Stock{}, &stock.StockLog{}, &stock.StockOpLog{},
		&system.SystemConfig{}, &system.SystemBackup{}, &system.SystemActivity{}, &system.SystemLog{}, &system.AIAnalysisLog{},
		&construction_log.ConstructionLog{},
		&progress.ProjectSchedule{},
		&progress.Resource{}, &progress.TaskResource{}, // 资源管理表
		&notification.Notification{},
		&material_plan.MaterialPlan{}, &material_plan.MaterialPlanItem{},
		&workflow.WorkflowDefinition{}, &workflow.WorkflowNode{}, &workflow.WorkflowEdge{},
		&workflow.WorkflowNodeApprover{}, &workflow.WorkflowInstance{}, &workflow.WorkflowApproval{},
		&workflow.WorkflowPendingTask{}, &workflow.WorkflowLog{},
		&agent.AgentOperationLog{},
		&audit2.OperationLog{}, // 操作日志表
		&notification.DeviceToken{}, // 设备推送令牌表
		&appointment.ConstructionAppointment{}, // 施工预约单表
		&appointment.WorkerCalendar{}, // 作业人员日历表
		&attendance.AttendanceRecord{}, // 打卡记录表
		&attendance.MonthlyAttendanceSummary{}, // 月度考勤汇总表
	)

	// create default admin role/user if not exists
	var count int64
	dbConn.Model(&auth.Role{}).Where("name = ?", "admin").Count(&count)
	if count == 0 {
		adminRole := auth.Role{Name: "admin", Description: "System administrator", Permissions: "admin"}
		dbConn.Create(&adminRole)
	}
	// create admin user 'admin' with password 'admin' in dev only if no users
	dbConn.Model(&auth.User{}).Count(&count)
	if count == 0 {
		passwd := "admin"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
		u := auth.User{Username: "admin", Password: string(hashed), Email: "admin@example.com", Role: "admin", IsActive: true}
		dbConn.Create(&u)
	}

	// 初始化上传配置
	initializeUploadConfig(dbConn)

	// 初始化默认物资分类
	initializeMaterialCategories(dbConn)

	// 初始化操作日志服务
	audit2.InitAuditService(dbConn)

	// 初始化 WebSocket Hub
	notification.InitHub()
	log.Println("WebSocket Hub initialized")

	// 初始化 AI Handler 和 Voice Handler
	var aiHandler *agent.AIHandler
	asrURL := ""
	if cfg.AI.ASREnabled {
		asrURL = cfg.AI.ASRServiceURL
	}

	// 初始化多模型 AI Handler
	aiHandler = agent.NewMultiModelAIHandler(
		dbConn,
		cfg.AI.BaiduAPIKey,
		cfg.AI.BaiduModel,
		cfg.AI.BaiduBaseURL,
		cfg.AI.DeepSeekAPIKey,
		cfg.AI.DeepSeekModel,
		cfg.AI.DeepSeekBaseURL,
		cfg.AI.OpenAIAPIKey,
		asrURL,
		cfg.AI.DefaultProvider,
	)
	log.Println("AI Handler initialized with multi-model support")

	if aiHandler != nil {
		// 创建并设置 VoiceHandler
		voiceHandler := notification.NewVoiceHandler(dbConn, aiHandler, asrURL)
		notification.GetHub().SetVoiceHandler(voiceHandler)
		log.Println("Voice Handler initialized and set on WebSocket Hub")
	}

	r := gin.Default()

	// 配置受信任的代理（安全设置）
	// 如果使用反向代理（如 Nginx），需要配置受信任的代理IP
	// 如果直接暴露给互联网，不要信任任何代理
	r.SetTrustedProxies([]string{"127.0.0.1", "::1","home.mbed.org.cn"}) // 只信任本地代理
	// 如果在反向代理后面，可以添加反向代理的IP：
	// r.SetTrustedProxies([]string{"10.0.0.1", "172.17.0.1"})

	// 全局中间件（需要在静态文件路由之前注册）
	// CORS 中间件（允许跨域请求）
	r.Use(middleware.CORS())

	// 全局安全中间件（应用到所有路由）
	r.Use(middleware.SecurityMiddleware())

	// JSON 大小限制中间件（限制请求体大小为 10MB）
	r.Use(middleware.ValidateJSON(10<<20)) // 10MB

	// 设置静态文件服务
	r.Static("/static", "./newstatic/dist")
	// /assets 路径根据设备类型动态选择（移动端或PC端）
	r.Use(func(c *gin.Context) {
		// 只处理 /assets 路径
		if !strings.HasPrefix(c.Request.URL.Path, "/assets/") {
			c.Next()
			return
		}

		// 检测设备类型
		userAgent := c.GetHeader("User-Agent")
		if isMobileDevice(userAgent) {
			// 移动端设备：从 dist-capacitor/assets 提供资源
			filePath := "./mobile-app/dist-capacitor" + c.Request.URL.Path
			c.File(filePath)
			c.Abort()
			return
		}

		// PC端：继续使用默认的静态文件路由
		c.Next()
	})
	r.Static("/assets", "./newstatic/dist/assets")
	r.Static("/uploads", "./static/uploads")  // 上传文件访问路径（保留在原位置）
	// 移动端应用（生产模式使用构建后的 dist，开发模式可配置源码目录）
	// 添加缓存控制中间件，确保 index.html 和 JS 文件不被缓存
	r.Use(func(c *gin.Context) {
		// 只处理 /mobile 路径
		if !strings.HasPrefix(c.Request.URL.Path, "/mobile") {
			c.Next()
			return
		}

		// 设置缓存控制头
		// index.html 和 version.json 不缓存
		// JS/CSS 文件使用短缓存（1小时）
		if c.Request.URL.Path == "/mobile/index.html" ||
		   c.Request.URL.Path == "/mobile/version.json" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		} else if strings.HasSuffix(c.Request.URL.Path, ".js") ||
		          strings.HasSuffix(c.Request.URL.Path, ".css") {
			c.Header("Cache-Control", "public, max-age=3600") // 1小时
		} else {
			c.Header("Cache-Control", "public, max-age=86400") // 24小时
		}

		c.Next()
	})
	r.Static("/mobile", "./mobile-app/dist")
	// 移动端更新文件（APK下载）
	r.Static("/mobile-updates", "./mobile-app-updates")
	// 版本文件（用于移动端版本检测）
	r.GET("/version.json", func(c *gin.Context) {
		c.File("./mobile-app/dist/version.json")
	})

	// API路由组
	api := r.Group("/api")
	{
		// 健康检查端点
		api.GET("/health", func(c *gin.Context) {
			// 检查数据库连接
			sqlDB, err := dbConn.DB()
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status":  "error",
					"message": "Failed to get database connection",
					"error":   err.Error(),
				})
				return
			}

			err = sqlDB.Ping()
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status":  "error",
					"message": "Database connection failed",
					"error":   err.Error(),
				})
				return
			}

			// 返回健康状态
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Service is healthy",
				"database": "PostgreSQL",
			})
		})
		auth.RegisterRoutes(api, dbConn)
		project.RegisterRoutes(api, dbConn)
		material.RegisterRoutes(api, dbConn)
		material.RegisterCategoryRoutes(api, dbConn)
		material_master.RegisterRoutes(r, material_master.NewHandler(material_master.NewService(dbConn)))
		inbound.RegisterRoutes(api, dbConn)
		requisition.RegisterRoutes(api, dbConn)
		stock.RegisterRoutes(api, dbConn)
		system.RegisterRoutes(api, dbConn)
		upload.RegisterRoutes(api, dbConn)
		construction_log.RegisterRoutes(api, dbConn)
		progress.RegisterRoutes(api, dbConn)
		notification.RegisterRoutes(api, dbConn)
		material_plan.RegisterRoutes(api, dbConn)
		workflow.RegisterRoutes(api, dbConn)
		agent.RegisterRoutes(api, dbConn, cfg, aiHandler)
		audit2.RegisterRoutes(api, dbConn) // 操作日志路由
		appointment.RegisterRoutes(api, dbConn) // 施工预约路由
		attendance.RegisterRoutes(api, dbConn) // 打卡考勤路由
		app.RegisterRoutes(api, dbConn) // 应用版本路由
	}

	// SPA路由回退 - 所有非API和非静态文件请求都返回对应的前端入口
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 排除API路径
		if filepath.HasPrefix(path, "/api") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// 排除明确指定PC端的静态资源路径
		if filepath.HasPrefix(path, "/static") ||
		   filepath.HasPrefix(path, "/uploads") ||
		   filepath.HasPrefix(path, "/mobile-updates") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// 移动端应用路由回退（排除 mobile-updates 和静态资源）
		// 只对移动端的主路径和没有扩展名的路径返回 index.html
		if filepath.HasPrefix(path, "/mobile") {
			// 检查是否请求静态资源（有文件扩展名）
			ext := filepath.Ext(path)
			if ext != "" {
				// 有扩展名，让静态文件路由处理
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			// 没有扩展名，返回 index.html（SPA 路由）
			c.File("./mobile-app/dist/index.html")
			return
		}

		// 根路径和其他路径：根据设备类型自动选择前端
		userAgent := c.GetHeader("User-Agent")
		if isMobileDevice(userAgent) {
			// 移动端设备：返回移动端应用（使用相对路径版本）
			c.File("./mobile-app/dist-capacitor/index.html")
			return
		}

		// PC端 SPA 路由回退 - 使用新构建的 Vue 3 应用
		c.File("./newstatic/dist/index.html")
	})

	// 设置Gin运行模式
	gin.SetMode(cfg.Server.Mode)

	// 创建并启动进度监听器
	progressWatcher := progress.NewProgressWatcher(dbConn)
	progressWatcher.Start()
	log.Println("进度监听器已启动，将自动同步进度计划变化到项目进度")

	// 创建HTTP服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("服务器启动在 %s (模式: %s)", addr, cfg.Server.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 停止进度监听器
	progressWatcher.Stop()

	// 优雅关闭HTTP服务器
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭失败: %v", err)
	}

	log.Println("服务器已关闭")
}