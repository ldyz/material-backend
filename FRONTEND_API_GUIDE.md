# 前端/移动端 API 使用指南

## 更新日期
2025-01-28

## 概述
本文档说明前端和移动端如何适配后端统一的响应格式。

---

## 一、后端统一响应格式

### 1.1 成功响应

```json
{
  "success": true,
  "data": {
    // 业务数据
  },
  "message": "操作成功（可选）",
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "pages": 5
  }
}
```

### 1.2 带分页的成功响应

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

### 1.3 带额外元数据的响应

```json
{
  "success": true,
  "data": {...},
  "meta": {
    // 额外元数据
  }
}
```

### 1.4 错误响应

```json
{
  "success": false,
  "error": "错误描述信息"
}
```

### 1.5 带错误码的错误响应

```json
{
  "success": false,
  "error": "错误描述信息",
  "code": "ERROR_CODE"
}
```

---

## 二、前端 (Vue 3 + Element Plus)

### 2.1 请求拦截器配置

**文件位置：** `/newstatic/src/api/request.js`

**核心功能：**
- ✅ 自动检查 `success` 字段
- ✅ 自动处理业务错误（`success: false`）
- ✅ 自动显示错误提示
- ✅ 自动处理 401 未授权（自动登出）
- ✅ 支持可选链操作符（`?.`）访问 `error` 字段

### 2.2 组件中使用 API

#### 示例 1：获取列表数据（带分页）

```javascript
import request from '@/api/request'

export default {
  data() {
    return {
      materials: [],
      pagination: {
        page: 1,
        perPage: 20,
        total: 0
      }
    }
  },
  methods: {
    async fetchMaterials() {
      try {
        // 直接解构响应中的 data 和 pagination
        const { data, pagination } = await request({
          url: '/materials',
          method: 'GET',
          params: {
            page: this.pagination.page,
            page_size: this.pagination.perPage
          }
        })

        // data 包含业务数据
        this.materials = data

        // pagination 包含分页信息
        this.pagination = pagination
      } catch (error) {
        // 错误已被 request.js 自动处理（显示错误提示）
        // 这里可以添加额外的错误处理逻辑
        console.error('获取物资列表失败:', error)
      }
    }
  }
}
```

#### 示例 2：创建资源

```javascript
async createMaterial(materialData) {
  try {
    // 成功响应：{ success: true, data: {...}, message: "创建成功" }
    const { data, message } = await request({
      url: '/materials',
      method: 'POST',
      data: materialData
    })

    // 显示成功消息（如果有）
    if (message) {
      this.$message.success(message)
    }

    // 使用返回的数据
    this.materials.push(data)
  } catch (error) {
    // 错误已被自动处理
    console.error('创建物资失败:', error)
  }
}
```

#### 示例 3：更新资源

```javascript
async updateMaterial(id, updateData) {
  try {
    const { data, message } = await request({
      url: `/materials/${id}`,
      method: 'PUT',
      data: updateData
    })

    if (message) {
      this.$message.success(message)
    }

    // 更新本地数据
    const index = this.materials.findIndex(m => m.id === id)
    if (index !== -1) {
      this.materials[index] = data
    }
  } catch (error) {
    console.error('更新物资失败:', error)
  }
}
```

#### 示例 4：删除资源

```javascript
async deleteMaterial(id) {
  try {
    const { message } = await request({
      url: `/materials/${id}`,
      method: 'DELETE'
    })

    if (message) {
      this.$message.success(message)
    }

    // 从本地列表中移除
    this.materials = this.materials.filter(m => m.id !== id)
  } catch (error) {
    console.error('删除物资失败:', error)
  }
}
```

#### 示例 5：处理带 meta 的响应

```javascript
async fetchStats() {
  try {
    // 响应：{ success: true, data: {...}, meta: {...} }
    const { data, meta } = await request({
      url: '/stats',
      method: 'GET'
    })

    // 使用业务数据
    this.stats = data

    // 使用元数据
    this.meta = meta
  } catch (error) {
    console.error('获取统计失败:', error)
  }
}
```

### 2.3 错误处理

**request.js 已经自动处理以下情况：**

1. **业务错误**（`success: false`）
   - 自动显示错误提示
   - 返回被拒绝的 Promise

2. **HTTP 401 错误**
   - 自动清除用户信息
   - 自动跳转到登录页

3. **其他 HTTP 错误**（403, 404, 500 等）
   - 自动显示错误提示

**组件中的额外错误处理：**

```javascript
try {
  const { data } = await request({ url: '/materials', method: 'POST', data: {...} })
  // 成功处理
} catch (error) {
  // 错误已被 request.js 自动处理
  // 这里可以添加特定场景的额外逻辑
  if (error.message.includes('特定错误')) {
    // 特殊处理
  }
}
```

---

## 三、移动端 (Vue 3 + Vant)

### 3.1 请求拦截器配置

**文件位置：** `/mobile-app/src/utils/request.js`

**核心功能：**
- ✅ 自动检查 `success` 字段
- ✅ 自动处理业务错误（`success: false`）
- ✅ 使用 Vant Toast 显示错误提示
- ✅ 自动处理 401/403 未授权（自动登出）

### 3.2 组件中使用 API

#### 示例 1：获取列表数据

```javascript
import request from '@/utils/request'

export default {
  data() {
    return {
      materials: [],
      loading: false,
      finished: false
    }
  },
  methods: {
    async fetchMaterials() {
      this.loading = true
      try {
        const { data, pagination } = await request({
          url: '/materials',
          method: 'GET'
        })

        this.materials = data

        // 检查是否还有更多数据
        this.finished = this.materials.length >= pagination.total
      } catch (error) {
        console.error('获取物资列表失败:', error)
      } finally {
        this.loading = false
      }
    }
  }
}
```

#### 示例 2：创建资源（带表单验证）

```javascript
import { showToast } from 'vant'

async onSubmit(formData) {
  try {
    const { data, message } = await request({
      url: '/materials',
      method: 'POST',
      data: formData
    })

    // 显示成功提示
    if (message) {
      showToast({
        type: 'success',
        message: message
      })
    }

    // 返回创建的数据
    return data
  } catch (error) {
    // 错误已被 request.js 自动处理
    console.error('创建失败:', error)
  }
}
```

#### 示例 3：下拉刷新和上拉加载

```javascript
import { List } from 'vant'

export default {
  data() {
    return {
      list: [],
      page: 1,
      loading: false,
      finished: false
    }
  },
  methods: {
    async onLoad() {
      try {
        const { data, pagination } = await request({
          url: '/materials',
          method: 'GET',
          params: {
            page: this.page,
            page_size: 20
          }
        })

        // 追加数据到列表
        this.list.push(...data)

        // 更新页码
        this.page++

        // 检查是否加载完成
        this.finished = this.list.length >= pagination.total
      } catch (error) {
        console.error('加载失败:', error)
      } finally {
        this.loading = false
      }
    },
    async onRefresh() {
      this.page = 1
      this.finished = false
      await this.onLoad()
    }
  }
}
```

### 3.3 错误处理

**request.js 已经自动处理：**

1. **业务错误**（`success: false`）
   - 使用 Vant Toast 显示错误提示
   - 返回被拒绝的 Promise

2. **HTTP 401/403 错误**
   - 清除本地存储
   - 跳转到登录页

3. **其他错误**
   - 显示相应的错误提示

---

## 四、响应拦截器对比

| 功能 | 前端 (Element Plus) | 移动端 (Vant) |
|------|---------------------|---------------|
| success 字段检查 | ✅ | ✅ |
| 业务错误处理 | ✅ | ✅ |
| 错误提示组件 | ElMessage | showToast |
| 401 自动登出 | ✅ | ✅ |
| 403 权限错误 | ✅ | ✅ |
| 可选链支持 | ✅ (data?.error) | ✅ (data?.error) |
| 返回完整响应 | ✅ | ✅ |

---

## 五、迁移指南

### 5.1 从旧代码迁移

**迁移前（旧代码）：**

```javascript
const response = await request({ url: '/materials', method: 'GET' })
if (response.data) {
  this.materials = response.data
}
```

**迁移后（新代码）：**

```javascript
const { data } = await request({ url: '/materials', method: 'GET' })
this.materials = data
```

### 5.2 错误处理迁移

**迁移前：**

```javascript
try {
  const response = await request({ url: '/materials', method: 'POST', data: {...} })
  if (response.success === false) {
    this.$message.error(response.error)
  }
} catch (error) {
  this.$message.error(error.message)
}
```

**迁移后：**

```javascript
try {
  const { data, message } = await request({
    url: '/materials',
    method: 'POST',
    data: {...}
  })
  // 业务错误已被 request.js 自动处理
  if (message) {
    this.$message.success(message)
  }
} catch (error) {
  // HTTP 错误已被 request.js 自动处理
  console.error('操作失败:', error)
}
```

---

## 六、最佳实践

### 6.1 统一的 API 调用模式

```javascript
// ✅ 推荐：使用解构赋值
const { data, pagination, message } = await request({ url: '/materials' })

// ❌ 不推荐：直接访问 response.data
const response = await request({ url: '/materials' })
const data = response.data
```

### 6.2 错误处理

```javascript
// ✅ 推荐：利用自动错误处理
try {
  const { data } = await request({ url: '/materials', method: 'POST', data: {...} })
  // 成功处理
} catch (error) {
  // 只处理特定场景的额外逻辑
}

// ❌ 不推荐：重复的错误处理
try {
  const response = await request({ url: '/materials', method: 'POST', data: {...} })
  if (response.success === false) {
    this.$message.error(response.error)
  }
} catch (error) {
  this.$message.error(error.message)
}
```

### 6.3 分页处理

```javascript
// ✅ 推荐：使用标准的分页结构
const { data, pagination } = await request({
  url: '/materials',
  params: { page: 1, page_size: 20 }
})

this.materials = data
this.total = pagination.total
```

---

## 七、常见问题

### Q1: 为什么 request.js 返回的不是直接的 data？

**A:** request.js 返回完整的响应对象，这样你可以访问：
- `data` - 业务数据
- `pagination` - 分页信息
- `message` - 成功消息
- `meta` - 其他元数据

如果只需要 data，可以使用解构赋值：
```javascript
const { data } = await request({ url: '/materials' })
```

### Q2: 如何禁用某个请求的错误提示？

**A:** 可以在 catch 块中阻止默认行为：
```javascript
try {
  const { data } = await request({ url: '/materials' })
} catch (error) {
  // 错误已自动显示，这里可以添加额外逻辑
  // 不需要再次显示错误提示
}
```

### Q3: 前端和移动端的响应处理有什么区别？

**A:**
- **功能相同**：都检查 `success` 字段，自动处理错误
- **UI 组件不同**：前端使用 Element Plus，移动端使用 Vant
- **配置略有不同**：超时时间、baseURL 等

---

## 八、相关文档

- [API_FIXES_REPORT.md](./API_FIXES_REPORT.md) - 后端 API 统一规范报告
- [API_STANDARD.md](./API_STANDARD.md) - API 标准规范文档
- [后端 response/builder.go](./internal/api/response/builder.go) - 后端响应包实现

---

**文档维护：** 本文档随 API 规范更新而更新，最后更新于 2025-01-28
