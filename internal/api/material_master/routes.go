package material_master

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册物资主数据路由
func RegisterRoutes(r *gin.Engine, handler *Handler) {
	masterGroup := r.Group("/api/materials/master")
	{
		// 物资主数据 CRUD
		masterGroup.POST("", handler.CreateMaterialMaster)           // 创建物资主数据
		masterGroup.PUT("/:id", handler.UpdateMaterialMaster)        // 更新物资主数据
		masterGroup.DELETE("/:id", handler.DeleteMaterialMaster)     // 删除物资主数据
		masterGroup.GET("/:id", handler.GetMaterialMaster)           // 获取物资主数据详情
		masterGroup.GET("", handler.ListMaterialsMaster)             // 获取物资主数据列表
		masterGroup.GET("/project", handler.ListProjectMaterials)    // 获取项目物资列表（带库存）
		masterGroup.GET("/categories", handler.GetCategories)        // 获取物资分类列表
	}
}
