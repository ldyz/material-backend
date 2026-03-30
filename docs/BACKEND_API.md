# 后端API文档

本文档详细说明了所有后端API接口的定义、请求参数和响应格式。

## 目录

- [认证模块 (auth)](#认证模块-auth)
- [项目管理模块 (project)](#项目管理模块-project)
- [施工预约模块 (appointment)](#施工预约模块-appointment)
- [考勤管理模块 (attendance)](#考勤管理模块-attendance)
- [AI助手模块 (agent)](#ai助手模块-agent)
- [物资计划模块 (material_plan)](#物资计划模块-material_plan)
- [入库管理模块 (inbound)](#入库管理模块-inbound)
- [出库管理模块 (requisition)](#出库管理模块-requisition)
- [库存管理模块 (stock)](#库存管理模块-stock)
- [通知模块 (notification)](#通知模块-notification)
- [工作流模块 (workflow)](#工作流模块-workflow)
- [进度管理模块 (progress)](#进度管理模块-progress)
- [施工日志模块 (construction_log)](#施工日志模块-construction_log)
- [审计日志模块 (audit)](#审计日志模块-audit)
- [系统管理模块 (system)](#系统管理模块-system)
- [文件上传模块 (upload)](#文件上传模块-upload)
- [应用版本模块 (app)](#应用版本模块-app)

---

## 认证模块 (auth)

### 数据模型

#### User 用户表

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | uint | 主键 |
| username | string | 用户名（唯一） |
| password | string | 密码（bcrypt加密） |
| email | string | 邮箱 |
| full_name | string | 全名 |
| avatar | string | 头像URL |
| role | string | 角色代码（兼容字段） |
| group | string | 用户组 |
| is_active | bool | 是否激活 |
| last_login | time | 最后登录时间 |
| created_at | time | 创建时间 |
| roles | []Role | 关联角色（many2many） |

#### Role 角色表

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | uint | 主键 |
| name | string | 角色名称 |
| description | string | 角色描述 |
| permissions | string | 权限列表（逗号分隔） |
| created_at | time | 创建时间 |

### API端点

#### POST /api/auth/login
用户登录

**权限要求**: 无

**请求参数**:
```json
{
  "username": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "full_name": "管理员",
    "avatar": "/uploads/avatars/admin.jpg",
    "role": "admin",
    "is_active": true,
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "permissions": ["*"]
      }
    ]
  },
  "meta": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### GET /api/auth/me
获取当前用户信息

**权限要求**: 需登录

**响应**: 返回当前用户信息

#### POST /api/auth/logout
用户登出

**权限要求**: 需登录

#### POST /api/auth/change-password
修改密码

**权限要求**: 需登录

**请求参数**:
```json
{
  "old_password": "string",
  "new_password": "string"
}
```

#### POST /api/auth/avatar
上传头像

**权限要求**: 需登录

**请求**: multipart/form-data
- avatar: 图片文件（支持jpg/png/gif/webp，最大2MB）

#### GET /api/auth/users
获取用户列表

**权限要求**: user_view

**查询参数**:
- page: 页码（默认1）
- page_size: 每页数量（默认20）

#### POST /api/auth/users
创建用户

**权限要求**: user_create

**请求参数**:
```json
{
  "username": "string",
  "password": "string",
  "email": "string",
  "full_name": "string",
  "role": "string",
  "role_ids": [1, 2],
  "is_active": true
}
```

#### GET /api/auth/users/:id
获取用户详情

**权限要求**: user_view

#### PUT /api/auth/users/:id
更新用户

**权限要求**: user_edit

#### DELETE /api/auth/users/:id
删除用户

**权限要求**: user_delete

#### POST /api/auth/users/:id/reset-password
重置用户密码

**权限要求**: user_edit

**请求参数**:
```json
{
  "password": "new_password"
}
```

#### GET /api/auth/roles
获取角色列表

**权限要求**: role_view

#### POST /api/auth/roles
创建角色

**权限要求**: role_create

**请求参数**:
```json
{
  "name": "string",
  "description": "string",
  "permissions": ["user_view", "user_create"]
}
```

#### PUT /api/auth/roles/:id
更新角色

**权限要求**: role_edit

#### DELETE /api/auth/roles/:id
删除角色

**权限要求**: role_delete

#### POST /api/auth/roles/:id/permissions
分配权限给角色

**权限要求**: role_assign_permissions

#### GET /api/auth/permissions
获取权限列表

**权限要求**: role_view

---

## 项目管理模块 (project)

### 数据模型

#### Project 项目表

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | uint | 主键 |
| name | string | 项目名称 |
| code | string | 项目编号 |
| location | string | 项目地点 |
| start_date | date | 开始日期 |
| end_date | date | 结束日期 |
| description | text | 项目描述 |
| manager | string | 项目经理 |
| contact | string | 联系方式 |
| budget | string | 预算 |
| status | string | 状态（planning/active/closed/on_hold） |
| parent_id | uint | 父项目ID（支持层级） |
| level | int | 层级深度 |
| path | string | 层级路径 |
| progress_percentage | float | 进度百分比 |
| users | []User | 关联用户 |

### API端点

#### GET /api/project/projects
获取项目列表

**权限要求**: project_view

**查询参数**:
- page: 页码
- page_size: 每页数量
- search: 搜索关键词
- status: 状态筛选
- manager: 项目经理筛选
- start_from: 开始日期范围
- start_to: 开始日期范围
- parent_id: 父项目ID
- sort: 排序字段（-id降序）
- show_all: 是否显示所有项目

#### POST /api/project/projects
创建项目

**权限要求**: project_create

**请求参数**:
```json
{
  "name": "string",
  "code": "string",
  "description": "string",
  "location": "string",
  "manager": "string",
  "contact": "string",
  "budget": "string",
  "status": "planning",
  "start_date": "2024-01-01",
  "end_date": "2024-12-31",
  "member_ids": [1, 2],
  "parent_id": null
}
```

#### GET /api/project/projects/:id
获取项目详情

**权限要求**: project_view

#### PUT /api/project/projects/:id
更新项目

**权限要求**: project_edit

#### DELETE /api/project/projects/:id
删除项目

**权限要求**: project_delete

#### GET /api/project/projects/:id/members
获取项目成员

**权限要求**: project_view

#### POST /api/project/projects/:id/members
分配项目成员

**权限要求**: project_edit

**请求参数**:
```json
{
  "user_ids": [1, 2, 3]
}
```

#### DELETE /api/project/projects/:id/members/:user_id
移除项目成员

**权限要求**: project_delete

#### GET /api/project/projects/:id/tree
获取项目树

**权限要求**: project_view

#### GET /api/project/projects/:id/children
获取子项目列表

**权限要求**: project_view

#### POST /api/project/projects/:id/aggregate-progress
聚合进度

**权限要求**: project_edit

---

## 施工预约模块 (appointment)

### 数据模型

#### ConstructionAppointment 施工预约单

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | uint | 主键 |
| appointment_no | string | 预约单号 |
| project_id | uint | 项目ID |
| applicant_id | uint | 申请人ID |
| applicant_name | string | 申请人姓名 |
| contact_phone | string | 联系电话 |
| contact_person | string | 联系人 |
| work_date | date | 作业日期 |
| time_slot | string | 时间段（morning/noon/afternoon/full_day） |
| work_location | string | 作业地点 |
| work_content | text | 作业内容 |
| work_type | string | 作业类型 |
| is_urgent | bool | 是否加急 |
| priority | int | 优先级（0-10） |
| urgent_reason | text | 加急原因 |
| assigned_worker_ids | text | 分配作业人员ID列表（JSON数组） |
| assigned_worker_names | text | 分配作业人员姓名（逗号分隔） |
| supervisor_id | uint | 监护人ID |
| supervisor_name | string | 监护人姓名 |
| status | string | 状态 |
| workflow_instance_id | uint | 工作流实例ID |
| submitted_at | time | 提交时间 |
| approved_at | time | 审批时间 |
| completed_at | time | 完成时间 |

#### WorkerCalendar 作业人员日历

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | uint | 主键 |
| worker_id | uint | 作业人员ID |
| calendar_date | date | 日期 |
| time_slot | string | 时间段 |
| is_available | bool | 是否可用 |
| status | string | 状态（available/busy/blocked/off） |
| appointment_id | uint | 关联预约ID |
| blocked_reason | text | 锁定原因 |

### 状态说明

| 状态 | 说明 |
|------|------|
| draft | 草稿 |
| pending | 待审批 |
| scheduled | 已排期 |
| in_progress | 进行中 |
| completed | 已完成 |
| cancelled | 已取消 |
| rejected | 已拒绝 |

### API端点

#### GET /api/appointments
获取预约单列表

**权限要求**: appointment_view

**查询参数**:
- page, page_size: 分页
- status: 状态筛选
- is_urgent: 是否加急
- start_date, end_date: 日期范围
- applicant_id: 申请人ID
- worker_id: 作业人员ID
- work_type: 作业类型

#### GET /api/appointments/my
获取我的预约

**权限要求**: 需登录

#### GET /api/appointments/pending
获取待审批列表

**权限要求**: appointment_approve

#### GET /api/appointments/workers
获取作业人员列表

**权限要求**: appointment_view

#### GET /api/appointments/stats
获取统计数据

**权限要求**: appointment_view

#### POST /api/appointments
创建预约单

**权限要求**: appointment_create

**请求参数**:
```json
{
  "project_id": 1,
  "contact_phone": "13800138000",
  "contact_person": "张三",
  "work_date": "2024-01-15",
  "time_slot": "morning",
  "work_location": "A区",
  "work_content": "设备安装",
  "work_type": "安装",
  "is_urgent": false,
  "priority": 5,
  "assigned_worker_ids": "[1,2]"
}
```

#### GET /api/appointments/:id
获取预约单详情

**权限要求**: appointment_view

#### PUT /api/appointments/:id
更新预约单

**权限要求**: appointment_edit

#### DELETE /api/appointments/:id
删除预约单

**权限要求**: appointment_delete

#### POST /api/appointments/:id/submit
提交审批

**权限要求**: appointment_submit

#### POST /api/appointments/:id/approve
审批预约

**权限要求**: appointment_approve

**请求参数**:
```json
{
  "action": "approve",
  "comment": "同意",
  "assign_now": true,
  "worker_id": 1,
  "reschedule": false,
  "new_work_date": null,
  "new_time_slot": null
}
```

#### POST /api/appointments/:id/assign
分配作业人员

**权限要求**: appointment_assign

**请求参数**:
```json
{
  "worker_id": 1,
  "worker_ids": [1, 2],
  "supervisor_id": 3
}
```

#### POST /api/appointments/:id/complete
完成预约

**权限要求**: appointment_execute

#### POST /api/appointments/:id/cancel
取消预约

**权限要求**: appointment_cancel

#### GET /api/appointments/:id/approval-history
获取审批历史

**权限要求**: appointment_view

#### GET /api/appointments/calendar/worker/:workerId
获取作业人员日历

**权限要求**: appointment_view

#### POST /api/appointments/calendar/check-availability
检查可用性

**权限要求**: appointment_view

#### GET /api/appointments/daily-statistics
获取每日统计

**权限要求**: appointment_view

---

## 考勤管理模块 (attendance)

### 数据模型

#### AttendanceRecord 打卡记录

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | uint | 主键 |
| user_id | uint | 用户ID |
| appointment_id | uint | 关联预约ID |
| work_content | text | 工作内容 |
| attendance_type | string | 打卡类型 |
| clock_in_time | time | 打卡时间 |
| clock_in_location | string | 打卡地点 |
| clock_in_latitude | float | 纬度 |
| clock_in_longitude | float | 经度 |
| overtime_hours | float | 加班小时数 |
| remark | text | 备注 |
| photo_url | string | 照片URL（单张） |
| photo_urls | text | 照片URL列表（JSON数组） |
| status | string | 状态 |

#### MonthlyAttendanceSummary 月度汇总

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | uint | 主键 |
| user_id | uint | 用户ID |
| year | int | 年份 |
| month | int | 月份 |
| morning_count | int | 上午打卡次数 |
| afternoon_count | int | 下午打卡次数 |
| noon_overtime_hours | float | 中午加班小时 |
| night_overtime_hours | float | 晚上加班小时 |
| total_work_days | int | 总工作天数 |
| total_overtime_hours | float | 总加班小时 |

### 打卡类型

| 类型 | 说明 |
|------|------|
| morning | 上午打卡 |
| afternoon | 下午打卡 |
| noon_overtime | 中午加班 |
| night_overtime | 晚上加班 |

### API端点

#### POST /api/attendance/clock-in
打卡

**权限要求**: 需登录

**请求参数**:
```json
{
  "appointment_id": 1,
  "work_content": "设备安装",
  "attendance_type": "morning",
  "clock_in_location": "项目现场",
  "clock_in_latitude": 39.9042,
  "clock_in_longitude": 116.4074,
  "overtime_hours": 0,
  "remark": "",
  "photo_urls": ["url1", "url2"],
  "clock_in_time": "2024-01-15 08:30:00"
}
```

#### GET /api/attendance/records
获取打卡记录列表

**权限要求**: attendance_view

#### GET /api/attendance/today-appointments
获取今日待打卡任务

**权限要求**: 需登录

#### GET /api/attendance/statistics
获取考勤统计

**权限要求**: attendance_view

#### POST /api/attendance/records/:id/confirm
确认打卡记录

**权限要求**: attendance_confirm

#### POST /api/attendance/records/:id/reject
驳回打卡记录

**权限要求**: attendance_confirm

#### GET /api/attendance/monthly-summary
获取月度汇总

**权限要求**: attendance_view

---

## AI助手模块 (agent)

### 数据模型

#### AgentOperationLog 操作日志

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | int64 | 主键 |
| operation | string | 操作类型 |
| resource | string | 资源类型 |
| parameters | jsonb | 操作参数 |
| reasoning | text | AI推理过程 |
| result | jsonb | 操作结果 |
| user_id | int | 用户ID |
| agent_id | string | 代理ID |
| status | string | 状态 |
| error | text | 错误信息 |

### 操作类型

| 类型 | 说明 |
|------|------|
| query | 查询操作 |
| analyze | 分析操作 |
| create_material_plan | 创建物资计划 |
| update_stock | 更新库存 |
| approve_workflow | 审批工作流 |
| generate_report | 生成报告 |

### API端点

#### POST /api/agent/chat
文本对话

**权限要求**: ai_agent_query

**请求参数**:
```json
{
  "message": "查询今天的入库单",
  "conversation_history": []
}
```

**响应**:
```json
{
  "message": "今天共有5条入库单...",
  "data": {},
  "intent": "query_inbound"
}
```

#### POST /api/agent/voice-chat
语音对话

**权限要求**: ai_agent_query

**请求**: multipart/form-data
- audio: 音频文件

**响应**:
```json
{
  "transcript": "查询今天的入库单",
  "message": "今天共有5条入库单...",
  "response": "今天共有5条入库单..."
}
```

#### POST /api/agent/operation
执行操作

**权限要求**: ai_agent_operate

**请求参数**:
```json
{
  "operation": "create_material_plan",
  "resource": "material_plan",
  "parameters": {...},
  "context": {},
  "reasoning": "用户请求创建物资计划"
}
```

#### GET /api/agent/capabilities
获取能力列表

**权限要求**: 需登录

#### POST /api/agent/validate
验证操作

**权限要求**: 需登录

#### GET /api/agent/logs
获取操作日志

**权限要求**: ai_agent_logs

#### GET /api/agent/providers
获取AI提供者列表

**权限要求**: 需登录

#### POST /api/agent/switch-provider
切换AI提供者

**权限要求**: 需登录

#### GET /api/agent/conversation-history
获取对话历史

**权限要求**: 需登录

#### DELETE /api/agent/conversation-history
清除对话历史

**权限要求**: 需登录

---

## 物资计划模块 (material_plan)

### API端点

#### GET /api/material-plans
获取计划列表

**权限要求**: material_plan_view

#### POST /api/material-plans
创建计划

**权限要求**: material_plan_create

#### GET /api/material-plans/:id
获取计划详情

**权限要求**: material_plan_view

#### PUT /api/material-plans/:id
更新计划

**权限要求**: material_plan_edit

#### DELETE /api/material-plans/:id
删除计划

**权限要求**: material_plan_delete

#### POST /api/material-plans/:id/submit
提交审批

**权限要求**: material_plan_submit

#### POST /api/material-plans/:id/approve
审批计划

**权限要求**: material_plan_approve

---

## 入库管理模块 (inbound)

### API端点

#### GET /api/inbound/orders
获取入库单列表

**权限要求**: inbound_view

#### POST /api/inbound/orders
创建入库单

**权限要求**: inbound_create

#### GET /api/inbound/orders/:id
获取入库单详情

**权限要求**: inbound_view

#### PUT /api/inbound/orders/:id
更新入库单

**权限要求**: inbound_edit

#### DELETE /api/inbound/orders/:id
删除入库单

**权限要求**: inbound_delete

#### POST /api/inbound/orders/:id/approve
审批入库单

**权限要求**: inbound_approve

---

## 出库管理模块 (requisition)

### API端点

#### GET /api/requisitions
获取出库单列表

**权限要求**: requisition_view

#### POST /api/requisitions
创建出库单

**权限要求**: requisition_create

#### GET /api/requisitions/:id
获取出库单详情

**权限要求**: requisition_view

#### PUT /api/requisitions/:id
更新出库单

**权限要求**: requisition_edit

#### DELETE /api/requisitions/:id
删除出库单

**权限要求**: requisition_delete

#### POST /api/requisitions/:id/approve
审批出库单

**权限要求**: requisition_approve

#### POST /api/requisitions/:id/issue
发货

**权限要求**: requisition_issue

---

## 库存管理模块 (stock)

### API端点

#### GET /api/stock/stocks
获取库存列表

**权限要求**: stock_view

#### GET /api/stock/stocks/:id
获取库存详情

**权限要求**: stock_view

#### PUT /api/stock/stocks/:id
更新库存

**权限要求**: stock_edit

#### POST /api/stock/stocks/:id/in
入库操作

**权限要求**: stock_in

#### POST /api/stock/stocks/:id/out
出库操作

**权限要求**: stock_out

#### GET /api/stock/logs
获取库存日志

**权限要求**: stocklog_view

---

## 通知模块 (notification)

### API端点

#### GET /api/notification/notifications
获取通知列表

**权限要求**: 需登录

#### PUT /api/notification/notifications/:id/read
标记已读

**权限要求**: 需登录

#### GET /api/notification/ws
WebSocket连接

**权限要求**: 需登录

**连接方式**:
```javascript
const ws = new WebSocket('ws://host/api/notification/ws?token=JWT_TOKEN')
```

---

## 工作流模块 (workflow)

### API端点

#### GET /api/workflows/definitions
获取工作流定义列表

**权限要求**: workflow_view

#### POST /api/workflows/definitions
创建工作流定义

**权限要求**: workflow_create

#### GET /api/workflows/definitions/:id
获取工作流定义详情

**权限要求**: workflow_view

#### PUT /api/workflows/definitions/:id
更新工作流定义

**权限要求**: workflow_edit

#### POST /api/workflows/definitions/:id/activate
激活工作流

**权限要求**: workflow_activate

#### GET /api/workflows/pending-tasks
获取待办任务

**权限要求**: workflow_task_view

#### POST /api/workflows/tasks/:id/approve
审批任务

**权限要求**: workflow_task_approve

---

## 进度管理模块 (progress)

### API端点

#### GET /api/progress/project/:id
获取项目进度计划

**权限要求**: progress_view

#### PUT /api/progress/project/:id
更新进度计划

**权限要求**: progress_edit

#### GET /api/progress/project/:id/tasks
获取任务列表

**权限要求**: progress_view

#### POST /api/progress/project/:id/tasks
创建任务

**权限要求**: progress_create

#### PUT /api/progress/tasks/:id
更新任务

**权限要求**: progress_edit

#### DELETE /api/progress/tasks/:id
删除任务

**权限要求**: progress_delete

---

## 施工日志模块 (construction_log)

### API端点

#### GET /api/construction-logs
获取日志列表

**权限要求**: construction_log_view

#### POST /api/construction-logs
创建日志

**权限要求**: construction_log_create

#### GET /api/construction-logs/:id
获取日志详情

**权限要求**: construction_log_view

#### PUT /api/construction-logs/:id
更新日志

**权限要求**: construction_log_edit

#### DELETE /api/construction-logs/:id
删除日志

**权限要求**: construction_log_delete

---

## 审计日志模块 (audit)

### API端点

#### GET /api/audit/logs
获取操作日志列表

**权限要求**: audit_view

#### GET /api/audit/logs/:id
获取日志详情

**权限要求**: audit_view

---

## 系统管理模块 (system)

### API端点

#### GET /api/system/configs
获取系统配置

**权限要求**: system_config

#### PUT /api/system/configs
更新系统配置

**权限要求**: system_config

#### GET /api/system/logs
获取系统日志

**权限要求**: system_log

#### GET /api/system/backups
获取备份列表

**权限要求**: system_backup

#### POST /api/system/backups
创建备份

**权限要求**: system_backup

#### GET /api/system/statistics
获取系统统计

**权限要求**: system_statistics

---

## 文件上传模块 (upload)

### API端点

#### POST /api/upload/upload
单文件上传

**权限要求**: 需登录

**请求**: multipart/form-data
- file: 文件

**响应**:
```json
{
  "url": "/uploads/xxx.jpg",
  "filename": "xxx.jpg"
}
```

#### POST /api/upload/upload-multiple
多文件上传

**权限要求**: 需登录

---

## 应用版本模块 (app)

### API端点

#### GET /api/app/version
获取最新版本信息

**权限要求**: 无

**响应**:
```json
{
  "platform": "android",
  "version": "1.0.100",
  "build_number": 100,
  "download_url": "https://xxx/app.apk",
  "force_update": false,
  "update_message": "修复若干问题",
  "release_notes": "1. 修复登录问题\n2. 优化性能"
}
```

#### GET /api/app/versions
获取版本历史

**权限要求**: 无

---

## 通用响应格式

### 成功响应
```json
{
  "success": true,
  "data": {...},
  "message": "操作成功"
}
```

### 分页响应
```json
{
  "success": true,
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 错误响应
```json
{
  "success": false,
  "error": "错误信息"
}
```

## 认证方式

所有需要认证的接口都需要在请求头中携带JWT Token：

```
Authorization: Bearer <token>
```

Token通过登录接口获取。
