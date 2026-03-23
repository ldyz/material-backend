package progress

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// 全局进度监听器实例
var globalWatcher *ProgressWatcher

// RegisterRoutes 注册进度管理路由
func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	handler := NewProgressHandler(db)

	// 创建并设置全局进度监听器
	globalWatcher = NewProgressWatcher(db)

	// 创建进度子组，添加JWT认证
	progressGroup := r.Group("progress")
	progressGroup.Use(jwtpkg.TokenMiddleware())

	// 进度列表（必须放在所有路由之前）
	progressGroup.GET("", auth.PermissionMiddleware(db, "progress_view"), handler.GetProgressList)

	// 导出功能
	progressGroup.GET("/export", auth.PermissionMiddleware(db, "progress_export"), handler.ExportProgress)

	// 创建任务（兼容旧版本）
	progressGroup.POST("", auth.PermissionMiddleware(db, "progress_create"), handler.CreateProgress)

	// 更新/删除任务（兼容旧版本，必须放在/project/:id之前）
	progressGroup.PUT("/:id", auth.PermissionMiddleware(db, "progress_edit"), handler.UpdateProgress)
	progressGroup.DELETE("/:id", auth.PermissionMiddleware(db, "progress_delete"), handler.DeleteProgress)

	// 获取所有项目进度计划状态
	progressGroup.GET("/project-schedules", auth.PermissionMiddleware(db, "progress_view"), handler.GetAllProjectSchedules)

	// 项目进度管理路由
	progressGroup.GET("/project/:id", auth.PermissionMiddleware(db, "progress_view"), handler.GetProjectSchedule)
	progressGroup.PUT("/project/:id", auth.PermissionMiddleware(db, "progress_edit"), handler.UpdateProjectSchedule)
	progressGroup.DELETE("/project/:id/schedule", auth.PermissionMiddleware(db, "progress_delete"), handler.DeleteProjectSchedule)
	progressGroup.GET("/project/:id/exists", auth.PermissionMiddleware(db, "progress_view"), handler.CheckScheduleExists)

	// Task management endpoints
	progressGroup.GET("/project/:id/tasks", auth.PermissionMiddleware(db, "progress_view"), handler.GetTasks)
	progressGroup.POST("/project/:id/tasks", auth.PermissionMiddleware(db, "progress_create"), handler.CreateTask)
	progressGroup.PUT("/tasks/:id", auth.PermissionMiddleware(db, "progress_edit"), handler.UpdateTask)
	progressGroup.DELETE("/tasks/:id", auth.PermissionMiddleware(db, "progress_delete"), handler.DeleteTask)

	// Dependency management
	progressGroup.GET("/tasks/:id/dependencies", auth.PermissionMiddleware(db, "progress_view"), handler.GetDependencies)
	progressGroup.POST("/tasks/:id/dependencies", auth.PermissionMiddleware(db, "progress_create"), handler.AddDependency)
	progressGroup.DELETE("/dependencies/:id", auth.PermissionMiddleware(db, "progress_delete"), handler.RemoveDependency)

	// Position persistence
	progressGroup.PUT("/tasks/:id/position", auth.PermissionMiddleware(db, "progress_edit"), handler.UpdateTaskPosition)

	// AI generation
	progressGroup.POST("/project/:id/generate-plan", auth.PermissionMiddleware(db, "progress_create"), handler.GeneratePlanWithAI)

	// Aggregate child plans
	progressGroup.POST("/project/:id/aggregate-plan", auth.PermissionMiddleware(db, "progress_edit"), handler.AggregateChildPlans)

	// ==================== 子任务进度管理 ====================
	// 计算父任务进度
	progressGroup.POST("/tasks/:id/calculate-parent-progress", auth.PermissionMiddleware(db, "progress_edit"), handler.CalculateParentTaskProgress)
	progressGroup.POST("/tasks/:id/update-parent-progress", auth.PermissionMiddleware(db, "progress_edit"), handler.UpdateTaskParentProgress)

	// ==================== 资源管理 ====================
	// 项目资源管理
	progressGroup.GET("/project/:id/resources", auth.PermissionMiddleware(db, "progress_view"), handler.GetProjectResources)
	progressGroup.POST("/project/:id/resources", auth.PermissionMiddleware(db, "progress_create"), handler.CreateResource)
	progressGroup.PUT("/project/:id/resources/:resourceId", auth.PermissionMiddleware(db, "progress_edit"), handler.UpdateResource)
	progressGroup.DELETE("/project/:id/resources/:resourceId", auth.PermissionMiddleware(db, "progress_delete"), handler.DeleteResource)

	// 任务资源分配
	progressGroup.GET("/tasks/:id/resources", auth.PermissionMiddleware(db, "progress_view"), handler.GetTaskResources)
	progressGroup.POST("/tasks/:id/resources", auth.PermissionMiddleware(db, "progress_edit"), handler.AllocateTaskResource)
	progressGroup.DELETE("/tasks/:id/resources/:resourceId", auth.PermissionMiddleware(db, "progress_delete"), handler.RemoveTaskResource)

	// ==================== 可视化创建依赖关系 ====================
	// 可视化创建依赖关系（从前端调用）
	// 使用单独的路由避免与 /tasks/:id/resources 冲突
	progressGroup.POST("/dependencies/visual/:fromId/:toId", auth.PermissionMiddleware(db, "progress_create"), handler.CreateDependencyVisual)

	// 进度同步管理接口（需要admin权限）
	adminGroup := r.Group("admin/progress")
	adminGroup.Use(jwtpkg.TokenMiddleware())
	{
		// 手动触发所有项目的进度同步
		adminGroup.POST("/sync-all", auth.PermissionMiddleware(db, "admin"), handleSyncAllProgress)
		// 手动触发指定项目的进度同步
		adminGroup.POST("/sync/:projectId", auth.PermissionMiddleware(db, "admin"), handleSyncProjectProgress)
		// 获取同步状态
		adminGroup.GET("/sync-status", auth.PermissionMiddleware(db, "admin"), handleGetSyncStatus)
	}
}

// handleSyncAllProgress 处理所有项目的进度同步请求
func handleSyncAllProgress(c *gin.Context) {
	if globalWatcher == nil {
		c.JSON(500, gin.H{"error": "进度监听器未初始化"})
		return
	}

	if err := globalWatcher.UpdateAllProjectsProgress(); err != nil {
		c.JSON(500, gin.H{"error": "同步失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "所有项目进度同步成功"})
}

// handleSyncProjectProgress 处理指定项目的进度同步请求
func handleSyncProjectProgress(c *gin.Context) {
	if globalWatcher == nil {
		c.JSON(500, gin.H{"error": "进度监听器未初始化"})
		return
	}

	var params struct {
		ProjectID uint `uri:"projectId" binding:"required"`
	}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(400, gin.H{"error": "无效的项目ID"})
		return
	}

	if err := globalWatcher.ForceUpdateProjectProgress(params.ProjectID); err != nil {
		c.JSON(500, gin.H{"error": "同步失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "项目进度同步成功"})
}

// handleGetSyncStatus 获取同步状态
func handleGetSyncStatus(c *gin.Context) {
	if globalWatcher == nil {
		c.JSON(500, gin.H{"error": "进度监听器未初始化"})
		return
	}

	status := globalWatcher.GetStatus()
	c.JSON(200, status)
}

