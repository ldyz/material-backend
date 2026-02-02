# 后端API统一规范

## 一、请求格式规范

### 1.1 统一请求头

```http
Content-Type: application/json
Authorization: Bearer {token}
```

### 1.2 分页参数

所有列表接口统一使用以下分页参数：

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码，从1开始 |
| page_size | int | 否 | 20 | 每页数量，最大100 |

### 1.3 排序参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| sort_by | string | 否 | id | 排序字段 |
| sort_order | string | 否 | asc | 排序方向：asc/desc |

### 1.4 搜索参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| search | string | 否 | 全文搜索关键词 |

### 1.5 过滤参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | string | 否 | 状态过滤 |
| project_id | int | 否 | 项目ID过滤 |
| start_date | string | 否 | 开始日期（YYYY-MM-DD） |
| end_date | string | 否 | 结束日期（YYYY-MM-DD） |

## 二、响应格式规范

### 2.1 成功响应

```json
{
  "success": true,
  "data": {
    // 业务数据
  },
  "message": "操作成功",
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "pages": 5
  }
}
```

### 2.2 错误响应

```json
{
  "success": false,
  "error": "错误描述信息",
  "code": "ERROR_CODE"
}
```

### 2.3 分页响应

```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "pages": 5
  }
}
```

### 2.4 带额外元数据的响应

```json
{
  "success": true,
  "data": [...],
  "meta": {
    // 额外元数据
  }
}
```

## 三、HTTP状态码规范

| 状态码 | 说明 | 使用场景 |
|--------|------|----------|
| 200 | OK | 请求成功 |
| 201 | Created | 创建成功 |
| 204 | No Content | 删除成功 |
| 400 | Bad Request | 请求参数错误 |
| 401 | Unauthorized | 未认证 |
| 403 | Forbidden | 无权限 |
| 404 | Not Found | 资源不存在 |
| 409 | Conflict | 资源冲突 |
| 422 | Unprocessable Entity | 验证失败 |
| 500 | Internal Server Error | 服务器错误 |

## 四、字段命名规范

### 4.1 JSON字段命名（snake_case）

```json
{
  "user_id": 1,
  "user_name": "admin",
  "created_at": "2024-01-01T00:00:00Z",
  "page_size": 20
}
```

### 4.2 Go结构体命名（PascalCase + json tag）

```go
type User struct {
    ID       uint      `json:"id"`
    UserName string    `json:"user_name"`
    CreatedAt time.Time `json:"created_at"`
}
```

## 五、常用请求体示例

### 5.1 创建资源

```json
{
  "name": "项目名称",
  "code": "PRJ001",
  "description": "项目描述"
}
```

### 5.2 更新资源

```json
{
  "name": "新项目名称",
  "status": "active"
}
```

### 5.3 批量操作

```json
{
  "ids": [1, 2, 3],
  "action": "approve",
  "data": {
    "comment": "审批通过"
  }
}
```

### 5.4 带物资明细的入库单

```json
{
  "supplier": "供应商名称",
  "inbound_date": "2024-01-01",
  "items": [
    {
      "material_id": 1,
      "quantity": 100,
      "price": 10.5,
      "remark": "备注"
    }
  ],
  "remark": "入库单备注"
}
```

## 六、验证规则

### 6.1 字符串验证

- `required`: 必填
- `min=x`: 最小长度
- `max=x`: 最大长度
- `email`: 邮箱格式

### 6.2 数值验证

- `min=x`: 最小值
- `max=x`: 最大值
- `gt=x`: 大于x
- `gte=x`: 大于等于x

### 6.3 日期验证

- 格式：`2006-01-02` 或 `2006-01-02T15:04:05Z`

## 七、错误码规范

| 错误码 | 说明 |
|--------|------|
| INVALID_PARAMS | 参数错误 |
| UNAUTHORIZED | 未认证 |
| FORBIDDEN | 无权限 |
| NOT_FOUND | 资源不存在 |
| ALREADY_EXISTS | 资源已存在 |
| CONFLICT | 资源冲突 |
| INTERNAL_ERROR | 内部错误 |

## 八、使用示例

### 8.1 导入请求包

```go
import (
    "github.com/yourorg/material-backend/backend/internal/api/request"
    "github.com/yourorg/material-backend/backend/internal/api/response"
)
```

### 8.2 分页列表

```go
var req request.ListRequest
if err := request.BindQuery(c, &req); err != nil {
    response.BadRequest(c, err.Error())
    return
}

page := req.GetPage()
pageSize := req.GetPageSize()
offset := req.GetOffset()

// 查询数据...
response.SuccessWithPagination(c, data, int64(page), int64(pageSize), total)
```

### 8.3 创建资源

```go
var req struct {
    Name string `json:"name" binding:"required"`
    Code string `json:"code" binding:"required"`
}
if err := request.BindJSON(c, &req); err != nil {
    response.BadRequest(c, err.Error())
    return
}

// 创建数据...
response.Created(c, result, "创建成功")
```

### 8.4 获取单个资源

```go
var req request.IDRequest
if err := request.BindURI(c, &req); err != nil {
    response.BadRequest(c, err.Error())
    return
}

// 查询数据...
response.Success(c, data)
```

### 8.5 更新资源

```go
var uriReq request.IDRequest
if err := request.BindURI(c, &uriReq); err != nil {
    response.BadRequest(c, err.Error())
    return
}

var bodyReq struct {
    Name string `json:"name"`
}
if err := request.BindJSON(c, &bodyReq); err != nil {
    response.BadRequest(c, err.Error())
    return
}

// 更新数据...
response.SuccessWithMessage(c, data, "更新成功")
```

### 8.6 删除资源

```go
var req request.IDRequest
if err := request.BindURI(c, &req); err != nil {
    response.BadRequest(c, err.Error())
    return
}

// 删除数据...
response.SuccessOnlyMessage(c, "删除成功")
```

## 九、响应包使用方法

```go
// 成功响应
response.Success(c, data)
response.SuccessWithMessage(c, data, "成功消息")
response.SuccessWithPagination(c, data, page, pageSize, total)
response.SuccessWithMeta(c, data, meta)
response.SuccessOnlyMessage(c, "仅返回消息")
response.Created(c, data, "创建成功")

// 错误响应
response.BadRequest(c, "参数错误")
response.Unauthorized(c, "未认证")
response.Forbidden(c, "无权限")
response.NotFound(c, "资源不存在")
response.Conflict(c, "资源冲突")
response.InternalError(c, "服务器错误")
```

## 十、实施状态

### 10.1 完成情况 ✅

| 项目 | 状态 | 完成度 |
|------|------|--------|
| 响应包实现 | ✅ 完成 | 100% |
| 请求包实现 | ✅ 完成 | 100% |
| Handler更新 | ✅ 完成 | 100% (14/14) |
| 编译验证 | ✅ 通过 | - |

### 10.2 已更新的Handler列表

- ✅ `progress/handler.go` - 进度管理
- ✅ `auth/handler.go` - 认证管理
- ✅ `auth/middleware.go` - 认证中间件
- ✅ `construction_log/handler.go` - 施工日志
- ✅ `upload/handler.go` - 文件上传
- ✅ `inbound/handler.go` - 入库管理
- ✅ `requisition/handler.go` - 领用管理
- ✅ `system/handler.go` - 系统管理
- ✅ `system/ai_handler.go` - AI分析 (23处替换)
- ✅ `system/report_handler.go` - 报表管理 (21处替换)
- ✅ `material/handler.go` - 物资管理
- ✅ `stock/handler.go` - 库存管理
- ✅ `project/handler.go` - 项目管理
- ✅ `workflow/handler.go` - 工作流

### 10.3 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| v3.0 | 2025-01-28 | 完成所有Handler响应格式统一 |
| v2.0 | 2025-01-28 | 创建统一响应包和请求包 |
| v1.0 | 2025-01-27 | 初始版本 |

---

**文档维护：** 本文档随API规范更新而更新，最后更新于 2025-01-28

