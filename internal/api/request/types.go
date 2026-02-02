package request

import (
	"github.com/gin-gonic/gin"
)

// ==================== 通用请求结构 ====================

// BaseRequest 基础请求结构
type BaseRequest struct {
	// 预留字段，可用于扩展
}

// ==================== 分页请求 ====================

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
}

// GetPage 获取页码，默认为1
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页数量，默认为20
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset 获取偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// ==================== 排序请求 ====================

// SortRequest 排序请求
type SortRequest struct {
	SortBy  string `form:"sort_by" json:"sort_by" binding:"omitempty"`           // 排序字段
	SortOrder string `form:"sort_order" json:"sort_order" binding:"omitempty,oneof=asc desc ASC DESC"` // 排序方向
}

// GetOrderBy 获取排序字段
func (s *SortRequest) GetOrderBy(defaultField string) string {
	if s.SortBy == "" {
		return defaultField
	}
	return s.SortBy
}

// GetSortOrder 获取排序方向
func (s *SortRequest) GetSortOrder() string {
	if s.SortOrder == "" {
		return "asc"
	}
	if s.SortOrder == "ASC" {
		return "ASC"
	}
	if s.SortOrder == "DESC" {
		return "DESC"
	}
	return s.SortOrder
}

// ==================== 搜索请求 ====================

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword string `form:"search" json:"search" binding:"omitempty"` // 搜索关键词
}

// ==================== 基础列表请求 ====================

// ListRequest 基础列表请求（分页+排序+搜索）
type ListRequest struct {
	PaginationRequest
	SortRequest
	SearchRequest
}

// ==================== ID请求 ====================

// IDRequest 单个ID请求
type IDRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// IDsRequest 多个ID请求
type IDsRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1,dive,min=1"`
}

// ==================== 范围过滤请求 ====================

// DateRangeRequest 日期范围请求
type DateRangeRequest struct {
	StartDate string `form:"start_date" json:"start_date" binding:"omitempty" format:"2006-01-02"`
	EndDate   string `form:"end_date" json:"end_date" binding:"omitempty" format:"2006-01-02"`
}

// NumberRangeRequest 数值范围请求
type NumberRangeRequest struct {
	Min float64 `form:"min" json:"min" binding:"omitempty"`
	Max float64 `form:"max" json:"max" binding:"omitempty"`
}

// ==================== 批量操作请求 ====================

// BatchOperationRequest 批量操作请求
type BatchOperationRequest struct {
	IDs    []uint                 `json:"ids" binding:"required,min=1,dive,min=1"`
	Action string                 `json:"action" binding:"required"` // 操作类型：delete, approve, reject等
	Data   map[string]interface{} `json:"data"`                      // 操作所需的数据
}

// ==================== 通用工具函数 ====================

// BindQuery 绑定查询参数
func BindQuery(c *gin.Context, obj interface{}) error {
	return c.ShouldBindQuery(obj)
}

// BindJSON 绑定JSON请求体
func BindJSON(c *gin.Context, obj interface{}) error {
	return c.ShouldBindJSON(obj)
}

// BindURI 绑定URI参数
func BindURI(c *gin.Context, obj interface{}) error {
	return c.ShouldBindUri(obj)
}

// BindAll 绑定所有参数（Query + JSON）
func BindAll(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return err
	}
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
}
