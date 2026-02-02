package workflow

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// RegisterRoutes 注册工作流路由
func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	handler := NewHandler(db)

	// 工作流配置路由（需要管理员权限）
	workflows := r.Group("/workflows")
	workflows.Use(jwtpkg.TokenMiddleware())
	workflows.Use(auth.PermissionMiddleware(db, "admin"))
	{
		workflows.GET("", handler.listWorkflows)
		workflows.POST("", handler.createWorkflow)
		workflows.GET("/:id", handler.getWorkflow)
		workflows.PUT("/:id", handler.updateWorkflow)
		workflows.DELETE("/:id", handler.deleteWorkflow)
		workflows.PUT("/:id/activate", handler.activateWorkflow)
		workflows.PUT("/:id/deactivate", handler.deactivateWorkflow)
	}

	// 工作流实例路由（需要登录）
	instances := r.Group("/workflow-instances")
	instances.Use(jwtpkg.TokenMiddleware())
	{
		instances.GET("", handler.listInstances)
		instances.GET("/:id", handler.getInstance)
		instances.GET("/:id/approvals", handler.getInstanceApprovals)
		instances.GET("/:id/logs", handler.getInstanceLogs)
		instances.POST("/:id/resubmit", handler.resubmitInstance)
	}

	// 工作流任务路由（需要登录）
	tasks := r.Group("/workflow-tasks")
	tasks.Use(jwtpkg.TokenMiddleware())
	{
		tasks.GET("/pending", handler.getPendingTasks)
		tasks.GET("/pending/:businessType", handler.getPendingTasksByType)
		tasks.POST("/:id/approve", handler.approveTask)
		tasks.POST("/:id/reject", handler.rejectTask)
		tasks.POST("/:id/return", handler.returnTask)
		tasks.POST("/:id/comment", handler.commentTask)
	}
}
