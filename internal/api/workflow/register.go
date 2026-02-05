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

	// 工作流配置路由（需要工作流管理权限）
	workflows := r.Group("/workflows")
	workflows.Use(jwtpkg.TokenMiddleware())
	{
		workflows.GET("", auth.PermissionMiddleware(db, "workflow_view"), handler.listWorkflows)
		workflows.POST("", auth.PermissionMiddleware(db, "workflow_create"), handler.createWorkflow)
		workflows.GET("/:id", auth.PermissionMiddleware(db, "workflow_view"), handler.getWorkflow)
		workflows.PUT("/:id", auth.PermissionMiddleware(db, "workflow_edit"), handler.updateWorkflow)
		workflows.DELETE("/:id", auth.PermissionMiddleware(db, "workflow_delete"), handler.deleteWorkflow)
		workflows.PUT("/:id/activate", auth.PermissionMiddleware(db, "workflow_activate"), handler.activateWorkflow)
		workflows.PUT("/:id/deactivate", auth.PermissionMiddleware(db, "workflow_activate"), handler.deactivateWorkflow)
	}

	// 工作流实例路由（需要登录）
	instances := r.Group("/workflow-instances")
	instances.Use(jwtpkg.TokenMiddleware())
	{
		instances.GET("", auth.PermissionMiddleware(db, "workflow_instance_view"), handler.listInstances)
		instances.GET("/:id", auth.PermissionMiddleware(db, "workflow_instance_view"), handler.getInstance)
		instances.GET("/:id/approvals", auth.PermissionMiddleware(db, "workflow_instance_view"), handler.getInstanceApprovals)
		instances.GET("/:id/logs", auth.PermissionMiddleware(db, "workflow_log_view"), handler.getInstanceLogs)
		instances.POST("/:id/resubmit", auth.PermissionMiddleware(db, "workflow_instance_resubmit"), handler.resubmitInstance)
	}

	// 工作流任务路由（需要登录）
	tasks := r.Group("/workflow-tasks")
	tasks.Use(jwtpkg.TokenMiddleware())
	{
		tasks.GET("/pending", auth.PermissionMiddleware(db, "workflow_task_view"), handler.getPendingTasks)
		tasks.GET("/pending/:businessType", auth.PermissionMiddleware(db, "workflow_task_view"), handler.getPendingTasksByType)
		tasks.POST("/:id/approve", auth.PermissionMiddleware(db, "workflow_task_approve"), handler.approveTask)
		tasks.POST("/:id/reject", auth.PermissionMiddleware(db, "workflow_task_reject"), handler.rejectTask)
		tasks.POST("/:id/return", auth.PermissionMiddleware(db, "workflow_task_reject"), handler.returnTask)
		tasks.POST("/:id/comment", auth.PermissionMiddleware(db, "workflow_task_approve"), handler.commentTask)
	}
}
