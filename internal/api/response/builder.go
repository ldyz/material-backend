package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Pagination 分页信息
type Pagination struct {
	Page    int64 `json:"page"`
	PerPage int64 `json:"per_page"`
	Total   int64 `json:"total"`
	Pages   int64 `json:"pages"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Success bool         `json:"success"`
	Error   string       `json:"error,omitempty"`
	Code    string       `json:"code,omitempty"`
	Details []ErrorDetail `json:"details,omitempty"`
	Context interface{}  `json:"context,omitempty"`
}

// ==================== 成功响应 ====================

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（带消息）
func SuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// SuccessOnlyMessage 成功响应（仅消息，无data）
func SuccessOnlyMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
	})
}

// SuccessWithPagination 成功响应（带分页）
func SuccessWithPagination(c *gin.Context, data interface{}, page, perPage, total int64) {
	pages := int64(0)
	if perPage > 0 {
		pages = (total + perPage - 1) / perPage
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
		Pagination: &Pagination{
			Page:    page,
			PerPage: perPage,
			Total:   total,
			Pages:   pages,
		},
	})
}

// SuccessWithMeta 成功响应（带元数据）
func SuccessWithMeta(c *gin.Context, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// SuccessWithPaginationAndMeta 成功响应（带分页和元数据）
func SuccessWithPaginationAndMeta(c *gin.Context, data interface{}, page, perPage, total int64, meta interface{}) {
	pages := int64(0)
	if perPage > 0 {
		pages = (total + perPage - 1) / perPage
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
		Pagination: &Pagination{
			Page:    page,
			PerPage: perPage,
			Total:   total,
			Pages:   pages,
		},
		Meta: meta,
	})
}

// Created 创建成功响应 (201)
func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// Accepted 已接受响应 (202)
func Accepted(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusAccepted, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// NoContent 无内容响应 (204)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ==================== 错误响应 ====================

// Error 错误响应（基础方法）
func Error(c *gin.Context, statusCode int, err string, code ...string) {
	resp := ErrorResponse{
		Success: false,
		Error:   err,
	}
	if len(code) > 0 {
		resp.Code = code[0]
	}
	c.JSON(statusCode, resp)
}

// ErrorWithDetails 错误响应（带详情）
func ErrorWithDetails(c *gin.Context, statusCode int, err string, code string, details []ErrorDetail) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   err,
		Code:    code,
		Details: details,
	})
}

// ErrorWithContext 错误响应（带上下文）
func ErrorWithContext(c *gin.Context, statusCode int, err string, code string, context interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   err,
		Code:    code,
		Context: context,
	})
}

// ErrorWithCode 带错误码的错误响应
func ErrorWithCode(c *gin.Context, statusCode int, err string, code string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   err,
		Code:    code,
	})
}

// BadRequest 400错误 - 请求参数错误
func BadRequest(c *gin.Context, err string) {
	Error(c, http.StatusBadRequest, err)
}

// BadRequestWithCode 400错误 - 带错误码
func BadRequestWithCode(c *gin.Context, err string, code string) {
	ErrorWithCode(c, http.StatusBadRequest, err, code)
}

// Unauthorized 401错误 - 未认证
func Unauthorized(c *gin.Context, err string) {
	Error(c, http.StatusUnauthorized, err)
}

// UnauthorizedWithCode 401错误 - 带错误码
func UnauthorizedWithCode(c *gin.Context, err string, code string) {
	ErrorWithCode(c, http.StatusUnauthorized, err, code)
}

// Forbidden 403错误 - 无权限
func Forbidden(c *gin.Context, err string) {
	Error(c, http.StatusForbidden, err)
}

// ForbiddenWithCode 403错误 - 带错误码
func ForbiddenWithCode(c *gin.Context, err string, code string) {
	ErrorWithCode(c, http.StatusForbidden, err, code)
}

// NotFound 404错误 - 资源不存在
func NotFound(c *gin.Context, err string) {
	Error(c, http.StatusNotFound, err)
}

// NotFoundWithCode 404错误 - 带错误码
func NotFoundWithCode(c *gin.Context, err string, code string) {
	ErrorWithCode(c, http.StatusNotFound, err, code)
}

// Conflict 409错误 - 资源冲突
func Conflict(c *gin.Context, err string) {
	Error(c, http.StatusConflict, err)
}

// ConflictWithCode 409错误 - 带错误码
func ConflictWithCode(c *gin.Context, err string, code string) {
	ErrorWithCode(c, http.StatusConflict, err, code)
}

// UnprocessableEntity 422错误 - 无法处理的实体
func UnprocessableEntity(c *gin.Context, err string) {
	Error(c, http.StatusUnprocessableEntity, err)
}

// UnprocessableEntityWithCode 422错误 - 带错误码
func UnprocessableEntityWithCode(c *gin.Context, err string, code string) {
	ErrorWithCode(c, http.StatusUnprocessableEntity, err, code)
}

// InternalError 500错误 - 服务器内部错误
func InternalError(c *gin.Context, err string) {
	Error(c, http.StatusInternalServerError, err)
}

// InternalErrorWithCode 500错误 - 带错误码
func InternalErrorWithCode(c *gin.Context, err string, code string) {
	ErrorWithCode(c, http.StatusInternalServerError, err, code)
}

// NotImplemented 501错误 - 未实现
func NotImplemented(c *gin.Context, err string) {
	Error(c, http.StatusNotImplemented, err)
}

// ServiceUnavailable 503错误 - 服务不可用
func ServiceUnavailable(c *gin.Context, err string) {
	Error(c, http.StatusServiceUnavailable, err)
}

// ==================== 验证辅助函数 ====================

// IsSuccess 检查响应是否成功
func IsSuccess(resp map[string]interface{}) bool {
	if success, ok := resp["success"].(bool); ok {
		return success
	}
	return false
}

// GetErrorMessage 从响应中获取错误消息
func GetErrorMessage(resp map[string]interface{}) string {
	if err, ok := resp["error"].(string); ok {
		return err
	}
	return "未知错误"
}
