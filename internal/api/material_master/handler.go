package material_master

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 物资主数据处理器
type Handler struct {
	service *Service
}

// NewHandler 创建物资主数据处理器
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateMaterialMaster 创建物资主数据
// @Summary 创建物资主数据
// @Tags MaterialMaster
// @Accept json
// @Produce json
// @Param request body CreateMaterialMasterRequest true "创建请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/materials/master [post]
func (h *Handler) CreateMaterialMaster(c *gin.Context) {
	var req CreateMaterialMasterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	material, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"data":    material.ToDTO(),
	})
}

// UpdateMaterialMaster 更新物资主数据
// @Summary 更新物资主数据
// @Tags MaterialMaster
// @Accept json
// @Produce json
// @Param id path int true "物资ID"
// @Param request body UpdateMaterialMasterRequest true "更新请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/materials/master/:id [put]
func (h *Handler) UpdateMaterialMaster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的物资ID"})
		return
	}

	var req UpdateMaterialMasterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	material, err := h.service.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    material.ToDTO(),
	})
}

// DeleteMaterialMaster 删除物资主数据
// @Summary 删除物资主数据
// @Tags MaterialMaster
// @Produce json
// @Param id path int true "物资ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/materials/master/:id [delete]
func (h *Handler) DeleteMaterialMaster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的物资ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

// GetMaterialMaster 获取物资主数据详情
// @Summary 获取物资主数据详情
// @Tags MaterialMaster
// @Produce json
// @Param id path int true "物资ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/materials/master/:id [get]
func (h *Handler) GetMaterialMaster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的物资ID"})
		return
	}

	material, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": material.ToDTO(),
	})
}

// ListMaterialsMaster 获取物资主数据列表
// @Summary 获取物资主数据列表
// @Tags MaterialMaster
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param category query string false "物资分类"
// @Success 200 {object} map[string]interface{}
// @Router /api/materials/master [get]
func (h *Handler) ListMaterialsMaster(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	materials, total, err := h.service.List(page, pageSize, keyword, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为 DTO
	items := make([]map[string]any, 0, len(materials))
	for _, m := range materials {
		items = append(items, m.ToDTO())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": items,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ListProjectMaterials 获取指定项目的物资列表（带库存）
// @Summary 获取指定项目的物资列表（带库存）
// @Tags MaterialMaster
// @Produce json
// @Param project_id query int true "项目ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param category query string false "物资分类"
// @Success 200 {object} map[string]interface{}
// @Router /api/materials/master/project [get]
func (h *Handler) ListProjectMaterials(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Query("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	materials, total, err := h.service.ListWithProjectStock(uint(projectID), page, pageSize, keyword, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": materials,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetCategories 获取物资分类列表
// @Summary 获取物资分类列表
// @Tags MaterialMaster
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/materials/master/categories [get]
func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}
