package material

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// RegisterCategoryRoutes 注册物资分类路由
func RegisterCategoryRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("material/categories")
	r.Use(jwtpkg.TokenMiddleware())

	// 获取分类列表
	r.GET("", auth.PermissionMiddleware(db, "material_view"), func(c *gin.Context) {
		var categories []MaterialCategory
		db.Order("sort ASC, id ASC").Find(&categories)

		// 构建树形结构
		tree := buildCategoryTree(categories, 0)

		response.Success(c, tree)
	})

	// 获取单个分类
	r.GET("/:id", auth.PermissionMiddleware(db, "material_view"), func(c *gin.Context) {
		id := c.Param("id")
		var category MaterialCategory
		if err := db.First(&category, id).Error; err != nil {
			response.NotFound(c, "分类不存在")
			return
		}
		response.Success(c, category.ToDTO())
	})

	// 创建分类
	r.POST("", auth.PermissionMiddleware(db, "material_create"), func(c *gin.Context) {
		var req struct {
			Name     string `json:"name" binding:"required"`
			Code     string `json:"code"`
			ParentID uint   `json:"parent_id"`
			Sort     int    `json:"sort"`
			Remark   string `json:"remark"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// 验证层级（最高4级）
		if req.ParentID > 0 {
			var parent MaterialCategory
			if err := db.First(&parent, req.ParentID).Error; err != nil {
				response.BadRequest(c, "父分类不存在")
				return
			}
			if parent.Level >= 4 {
				response.BadRequest(c, "分类层级不能超过4级")
				return
			}
		}

		// 检查名称是否在同级中重复
		var count int64
		query := db.Model(&MaterialCategory{}).Where("name = ?", req.Name)
		if req.ParentID > 0 {
			query = query.Where("parent_id = ?", req.ParentID)
		} else {
			query = query.Where("parent_id = 0")
		}
		query.Count(&count)
		if count > 0 {
			response.BadRequest(c, "同级分类名称已存在")
			return
		}

		// 计算层级和路径
		level := 1
		path := ""
		if req.ParentID > 0 {
			var parent MaterialCategory
			db.First(&parent, req.ParentID)
			level = parent.Level + 1
			path = fmt.Sprintf("%s/%d", parent.Path, 0) // 临时path，创建后更新
		}

		category := MaterialCategory{
			ParentID: req.ParentID,
			Name:     req.Name,
			Code:     req.Code,
			Level:    level,
			Path:     path,
			Sort:     req.Sort,
			Remark:   req.Remark,
		}

		if err := db.Create(&category).Error; err != nil {
			response.InternalError(c, "创建分类失败")
			return
		}

		// 更新path
		if req.ParentID > 0 {
			var parent MaterialCategory
			db.First(&parent, req.ParentID)
			category.Path = fmt.Sprintf("%s/%d", parent.Path, category.ID)
			db.Save(&category)
		} else {
			category.Path = fmt.Sprintf("%d", category.ID)
			db.Save(&category)
		}

		response.Created(c, category.ToDTO(), "分类创建成功")
	})

	// 更新分类
	r.PUT("/:id", auth.PermissionMiddleware(db, "material_edit"), func(c *gin.Context) {
		id := c.Param("id")
		var category MaterialCategory
		if err := db.First(&category, id).Error; err != nil {
			response.NotFound(c, "分类不存在")
			return
		}

		var req struct {
			Name     string `json:"name"`
			Code     string `json:"code"`
			ParentID *uint  `json:"parent_id"`
			Sort     int    `json:"sort"`
			Remark   string `json:"remark"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// 检查名称是否与其他分类重复
		if req.Name != "" && req.Name != category.Name {
			var count int64
			parentID := category.ParentID
			if req.ParentID != nil {
				parentID = *req.ParentID
			}
			db.Model(&MaterialCategory{}).Where("name = ? AND parent_id = ? AND id != ?", req.Name, parentID, id).Count(&count)
			if count > 0 {
				response.BadRequest(c, "同级分类名称已存在")
				return
			}
			category.Name = req.Name
		}

		// 处理父分类变更
		if req.ParentID != nil {
			newParentID := *req.ParentID

			// 不能设置自己为父分类
			if newParentID == category.ID {
				response.BadRequest(c, "不能设置自己为父分类")
				return
			}

			// 验证新父分类
			if newParentID > 0 {
				var newParent MaterialCategory
				if err := db.First(&newParent, newParentID).Error; err != nil {
					response.BadRequest(c, "父分类不存在")
					return
				}

				// 检查是否会造成循环引用
				if isDescendant(db, newParentID, category.ID) {
					response.BadRequest(c, "不能将分类移动到其子分类下")
					return
				}

				// 验证层级
				if newParent.Level >= 4 {
					response.BadRequest(c, "分类层级不能超过4级")
					return
				}

				category.ParentID = newParentID
				category.Level = newParent.Level + 1
			} else {
				category.ParentID = 0
				category.Level = 1
			}
		}

		if req.Code != "" {
			category.Code = req.Code
		}
		category.Sort = req.Sort
		category.Remark = req.Remark

		// 更新path
		if category.ParentID > 0 {
			var parent MaterialCategory
			db.First(&parent, category.ParentID)
			category.Path = fmt.Sprintf("%s/%d", parent.Path, category.ID)
		} else {
			category.Path = fmt.Sprintf("%d", category.ID)
		}

		if err := db.Save(&category).Error; err != nil {
			response.InternalError(c, "更新分类失败")
			return
		}

		// 更新所有子分类的path和level
		updateChildrenPath(db, category)

		response.SuccessWithMessage(c, category.ToDTO(), "分类更新成功")
	})

	// 删除分类
	r.DELETE("/:id", auth.PermissionMiddleware(db, "material_delete"), func(c *gin.Context) {
		id := c.Param("id")
		var category MaterialCategory
		if err := db.First(&category, id).Error; err != nil {
			response.NotFound(c, "分类不存在")
			return
		}

		// 检查是否有子分类
		var childCount int64
		db.Model(&MaterialCategory{}).Where("parent_id = ?", category.ID).Count(&childCount)
		if childCount > 0 {
			response.BadRequest(c, "该分类下还有子分类，无法删除")
			return
		}

		// 检查是否有物资使用此分类或其子分类
		var count int64
		db.Table("material_master").Where("category = ?", category.Name).Count(&count)
		if count > 0 {
			response.BadRequest(c, "该分类下还有物资，无法删除")
			return
		}

		if err := db.Delete(&category).Error; err != nil {
			response.InternalError(c, "删除分类失败")
			return
		}

		response.SuccessWithMessage(c, nil, "分类删除成功")
	})

	// 批量更新排序
	r.POST("/sort", auth.PermissionMiddleware(db, "material_edit"), func(c *gin.Context) {
		var req struct {
			Sorts []struct {
				ID    uint `json:"id"`
				Sort  int  `json:"sort"`
			} `json:"sorts" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// 使用事务更新排序
		tx := db.Begin()
		for _, item := range req.Sorts {
			if err := tx.Model(&MaterialCategory{}).Where("id = ?", item.ID).Update("sort", item.Sort).Error; err != nil {
				tx.Rollback()
				response.InternalError(c, "更新排序失败")
				return
			}
		}
		tx.Commit()

		response.SuccessWithMessage(c, nil, "排序更新成功")
	})
}

// buildCategoryTree 构建分类树
func buildCategoryTree(categories []MaterialCategory, parentID uint) []map[string]any {
	tree := make([]map[string]any, 0)

	for _, cat := range categories {
		if cat.ParentID == parentID {
			node := cat.ToDTO()
			// 递归查找子分类
			children := buildCategoryTree(categories, cat.ID)
			if len(children) > 0 {
				node["children"] = children
			}
			tree = append(tree, node)
		}
	}

	return tree
}

// isDescendant 检查是否是子孙节点
func isDescendant(db *gorm.DB, ancestorID, descendantID uint) bool {
	var current MaterialCategory
	if err := db.First(&current, descendantID).Error; err != nil {
		return false
	}

	for current.ParentID > 0 {
		if current.ParentID == ancestorID {
			return true
		}
		if err := db.First(&current, current.ParentID).Error; err != nil {
			return false
		}
	}

	return false
}

// updateChildrenPath 递归更新子分类的path和level
func updateChildrenPath(db *gorm.DB, parent MaterialCategory) {
	var children []MaterialCategory
	db.Where("parent_id = ?", parent.ID).Find(&children)

	for _, child := range children {
		child.Level = parent.Level + 1
		child.Path = fmt.Sprintf("%s/%d", parent.Path, child.ID)
		db.Save(&child)

		// 递归更新子分类
		updateChildrenPath(db, child)
	}
}
