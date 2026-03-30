# Web前端文档

本文档详细说明了Web端前端项目的结构、组件和使用方法。

## 项目概述

- **框架**: Vue 3 + Composition API
- **UI组件库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **构建工具**: Vite
- **HTTP客户端**: Axios

## 目录结构

```
newstatic/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API接口定义
│   │   └── index.js       # 统一API导出
│   ├── assets/            # 静态资源
│   ├── components/        # 公共组件
│   │   ├── layout/        # 布局组件
│   │   │   └── MainLayout.vue
│   │   └── common/        # 通用组件
│   │       └── AiChat.vue # AI对话组件
│   ├── router/            # 路由配置
│   │   └── index.js
│   ├── stores/            # Pinia状态管理
│   │   ├── auth.js        # 认证状态
│   │   └── ...
│   ├── utils/             # 工具函数
│   │   └── index.js
│   ├── views/             # 页面组件
│   │   ├── Login.vue
│   │   ├── Dashboard.vue
│   │   ├── Projects.vue
│   │   ├── Materials.vue
│   │   ├── MaterialPlans.vue
│   │   ├── Stock.vue
│   │   ├── Inbound.vue
│   │   ├── Requisitions.vue
│   │   ├── Appointments.vue
│   │   ├── Attendance.vue
│   │   ├── Progress.vue
│   │   ├── ConstructionLog.vue
│   │   ├── WorkflowManagement.vue
│   │   ├── UserManagement.vue
│   │   ├── RoleManagement.vue
│   │   ├── System.vue
│   │   ├── OperationLogs.vue
│   │   ├── Notifications.vue
│   │   └── ...
│   ├── App.vue
│   └── main.js
├── index.html
├── vite.config.js
├── package.json
└── build.sh               # 构建脚本
```

## 页面清单

| 页面 | 路由 | 权限 | 说明 |
|------|------|------|------|
| Login | /login | 无 | 用户登录 |
| Dashboard | /dashboard | 登录用户 | 仪表板 |
| Projects | /projects | project_view | 项目管理 |
| Materials | /materials | material_view | 物资浏览 |
| MaterialCategories | /material-categories | material_view | 物资分类 |
| MaterialPlans | /material-plans | material_plan_view | 物资计划 |
| PlanStatistics | /plan-statistics | material_plan_view | 计划统计 |
| Stock | /stock | stock_view | 库存浏览 |
| Inbound | /inbound | inbound_view | 入库管理 |
| Requisitions | /requisitions | requisition_view | 出库管理 |
| Appointments | /appointments | appointment_view | 施工预约 |
| Attendance | /attendance | attendance_view | 考勤管理 |
| Progress | /progress | progress_view | 进度管理 |
| ConstructionLog | /construction-log | construction_log_view | 施工日志 |
| Workflows | /workflows | system_config | 工作流管理 |
| UserManagement | /system/users | user_view | 用户管理 |
| RoleManagement | /system/roles | role_view | 角色管理 |
| System | /system | system_config | 系统管理 |
| OperationLogs | /operation-logs | audit_view | 操作日志 |
| Notifications | /notifications | 登录用户 | 通知中心 |
| ResetPassword | /reset-password | 登录用户 | 修改密码 |

## 核心组件

### MainLayout 主布局组件

**路径**: `src/components/layout/MainLayout.vue`

**功能**:
- 左侧菜单导航
- 顶部用户信息栏
- 右侧内容区域
- AI助手悬浮按钮

**Props**: 无

**使用示例**:
```vue
<template>
  <MainLayout>
    <router-view />
  </MainLayout>
</template>
```

### AiChat AI对话组件

**路径**: `src/components/common/AiChat.vue`

**功能**:
- 文本对话界面
- 语音输入支持
- 对话历史显示
- 模型切换

**Props**:
| 名称 | 类型 | 默认值 | 说明 |
|------|------|-------|------|
| visible | Boolean | false | 是否显示 |

**Events**:
| 名称 | 参数 | 说明 |
|------|------|------|
| update:visible | Boolean | 显示状态变化 |
| close | - | 关闭对话框 |

**使用示例**:
```vue
<template>
  <AiChat v-model:visible="showChat" />
</template>
```

## 状态管理

### auth.js 认证状态

**状态**:
```javascript
{
  user: null,           // 当前用户信息
  token: null,          // JWT Token
  isAuthenticated: false, // 是否已认证
  isAdmin: false,       // 是否管理员
  permissions: []       // 权限列表
}
```

**Actions**:
- `login(credentials)` - 用户登录
- `logout()` - 用户登出
- `fetchUser()` - 获取用户信息
- `hasPermission(perm)` - 检查权限

**使用示例**:
```javascript
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 登录
await authStore.login({ username: 'admin', password: '123456' })

// 检查权限
if (authStore.hasPermission('user_create')) {
  // 有创建用户权限
}
```

## API调用层

API定义位于 `src/api/index.js`，使用Axios封装。

**配置**:
```javascript
const api = axios.create({
  baseURL: '/api',
  timeout: 30000
})

// 请求拦截器 - 添加Token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器 - 处理错误
api.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response?.status === 401) {
      // Token过期，跳转登录
      router.push('/login')
    }
    return Promise.reject(error)
  }
)
```

**API方法**:
```javascript
// 认证
login(data)
logout()
getMe()
changePassword(data)

// 用户
getUsers(params)
createUser(data)
updateUser(id, data)
deleteUser(id)
resetPassword(id, data)

// 项目
getProjects(params)
createProject(data)
getProject(id)
updateProject(id, data)
deleteProject(id)
getProjectMembers(id)
addProjectMembers(id, data)

// 物资计划
getMaterialPlans(params)
createMaterialPlan(data)
getMaterialPlan(id)
updateMaterialPlan(id, data)
deleteMaterialPlan(id)
submitMaterialPlan(id)
approveMaterialPlan(id, data)

// 入库
getInboundOrders(params)
createInboundOrder(data)
getInboundOrder(id)
updateInboundOrder(id, data)
deleteInboundOrder(id)
approveInboundOrder(id, data)

// 出库
getRequisitions(params)
createRequisition(data)
getRequisition(id)
updateRequisition(id, data)
deleteRequisition(id)
approveRequisition(id, data)
issueRequisition(id, data)

// 施工预约
getAppointments(params)
createAppointment(data)
getAppointment(id)
updateAppointment(id, data)
deleteAppointment(id)
submitAppointment(id)
approveAppointment(id, data)
assignWorkers(id, data)
completeAppointment(id, data)

// 考勤
clockIn(data)
getAttendanceRecords(params)
getTodayAppointments()
confirmRecord(id, data)
rejectRecord(id, data)

// AI助手
chat(data)
voiceChat(formData)
getProviders()
switchProvider(data)
getConversationHistory()
clearConversationHistory()

// 通知
getNotifications(params)
markAsRead(id)
connectWebSocket()
```

## 路由配置

路由配置位于 `src/router/index.js`。

**路由守卫**:
```javascript
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - 材料管理系统` : '材料管理系统'

  // 检查认证
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  // 检查权限
  if (to.meta.permissions && !hasAnyPermission(user, to.meta.permissions)) {
    next({ name: 'Dashboard' })
    return
  }

  next()
})
```

## 工具函数

### permissions.js 权限工具

```javascript
// 检查是否有任一权限
export function hasAnyPermission(user, permissions) {
  if (user.isAdmin) return true
  return permissions.some(p => user.permissions.includes(p))
}

// 检查是否有所有权限
export function hasAllPermissions(user, permissions) {
  if (user.isAdmin) return true
  return permissions.every(p => user.permissions.includes(p))
}
```

### 日期格式化

```javascript
// 格式化日期
export function formatDate(date, format = 'YYYY-MM-DD') {
  return dayjs(date).format(format)
}

// 格式化日期时间
export function formatDateTime(date) {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}
```

## 开发命令

```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build

# 使用构建脚本（推荐）
./build.sh

# 预览构建结果
npm run preview
```

## 构建脚本

`build.sh` 脚本内容：

```bash
#!/bin/bash
echo "Building web frontend..."
npm run build
echo "Build completed!"
echo "Output directory: dist/"
```

## 环境配置

### 开发环境 `.env.development`

```
VITE_API_BASE_URL=/api
VITE_WS_URL=ws://localhost:8088/api/notification/ws
```

### 生产环境 `.env.production`

```
VITE_API_BASE_URL=/api
VITE_WS_URL=wss://your-domain.com/api/notification/ws
```

## 样式规范

- 使用Element Plus默认主题
- 自定义变量在 `src/assets/styles/variables.scss`
- 响应式断点：
  - xs: < 576px
  - sm: >= 576px
  - md: >= 768px
  - lg: >= 992px
  - xl: >= 1200px

## 常见问题

### 1. Token过期处理

系统会自动检测401响应并跳转到登录页。

### 2. 权限不足

无权限访问的页面会自动重定向到仪表板。

### 3. WebSocket断线重连

通知模块实现了自动重连机制。

### 4. 文件上传

使用 `el-upload` 组件，支持拖拽上传和点击上传。
