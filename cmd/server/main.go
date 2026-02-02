package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/project"
	"github.com/yourorg/material-backend/backend/internal/api/material"
	"github.com/yourorg/material-backend/backend/internal/api/material_master"
	"github.com/yourorg/material-backend/backend/internal/api/inbound"
	"github.com/yourorg/material-backend/backend/internal/api/requisition"
	"github.com/yourorg/material-backend/backend/internal/api/stock"
	"github.com/yourorg/material-backend/backend/internal/api/system"
	"github.com/yourorg/material-backend/backend/internal/api/construction_log"
	"github.com/yourorg/material-backend/backend/internal/api/upload"
	"github.com/yourorg/material-backend/backend/internal/api/progress"
	"github.com/yourorg/material-backend/backend/internal/api/notification"
	"github.com/yourorg/material-backend/backend/internal/api/workflow"
	"github.com/yourorg/material-backend/backend/internal/api/material_plan"
	"github.com/yourorg/material-backend/backend/internal/db"
	"github.com/yourorg/material-backend/backend/internal/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
	// Load .env file if exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or error loading")
	}
	// Load config from env
	// Only use PostgreSQL database as requested
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN environment variable is required. No fallback database will be used.")
	}
	log.Println("Using PostgreSQL database from POSTGRES_DSN")

	dbConn, err := db.New(dsn)
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
	r.Static("/static", "./static")
	r.Static("/assets", "./static/assets")
	// 移动端应用（生产模式使用构建后的 dist，开发模式可配置源码目录）
	r.Static("/mobile", "./mobile-app/dist")

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
	}

	// SPA路由回退 - 所有非API和非静态文件请求都返回对应的前端入口
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 移动端应用路由回退
		if filepath.HasPrefix(path, "/mobile") {
			c.File("./mobile-app/dist/index.html")
			return
		}

		// PC端 SPA 路由回退 - 使用新构建的 Vue 3 应用
		// 排除 API 路径和静态资源路径
		if !filepath.HasPrefix(path, "/api") &&
		   !filepath.HasPrefix(path, "/static") &&
		   !filepath.HasPrefix(path, "/assets") &&
		   !filepath.HasPrefix(path, "/mobile") {
			c.File("./static/index.html")
			return
		}

		// 其他情况返回404
		c.AbortWithStatus(http.StatusNotFound)
	})

	// 创建并启动进度监听器
	progressWatcher := progress.NewProgressWatcher(dbConn)
	progressWatcher.Start()
	log.Println("进度监听器已启动，将自动同步进度计划变化到项目进度")

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":8088",
		Handler: r,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Println("服务器启动在 :8088")
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
	log.Println("服务器已关闭")
}