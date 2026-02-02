# 物资管理系统 API 端点文档

## 基础信息

**基础URL**: `http://your-domain:8088/api`

**认证方式**: JWT Token (Bearer Token)

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {token}
```

**响应格式**:
```json
{
  "success": true,
  "data": {},
  "message": "操作成功",
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "pages": 5
  }
}
```

---

## 一、认证模块 (`/api/auth`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| POST | `/auth/login` | - | 用户登录 |
| POST | `/auth/register` | - | 用户注册 |
| GET | `/auth/debug-token` | - | 调试 Token |

---

## 二、项目管理 (`/api/projects`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/projects` | `project_view` | 获取项目列表 |
| POST | `/projects` | `project_create` | 创建项目 |
| GET | `/projects/:id` | `project_view` | 获取项目详情 |
| PUT | `/projects/:id` | `project_edit` | 更新项目 |
| DELETE | `/projects/:id` | `project_delete` | 删除项目 |
| GET | `/projects/:id/members` | `project_view` | 获取项目成员 |
| POST | `/projects/:id/members` | `project_edit` | 添加项目成员 |
| DELETE | `/projects/:id/members/:user_id` | `project_delete` | 移除项目成员 |
| GET | `/projects/:id/tree` | `project_view` | 获取项目树 |
| GET | `/projects/:id/children` | `project_view` | 获取子项目 |
| POST | `/projects/:id/aggregate-progress` | `project_edit` | 聚合项目进度 |

---

## 三、物资管理 (`/api/material`)

### 3.1 物资主数据 (`/api/materials/master`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/materials/master` | - | 获取物资主数据列表 |
| GET | `/materials/master/:id` | - | 获取物资主数据详情 |
| POST | `/materials/master` | - | 创建物资主数据 |
| PUT | `/materials/master/:id` | - | 更新物资主数据 |
| DELETE | `/materials/master/:id` | - | 删除物资主数据 |
| GET | `/materials/master/project` | - | 获取项目物资列表（带库存） |
| GET | `/materials/master/categories` | - | 获取物资分类列表 |

### 3.2 物资管理 (`/api/materials`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/materials` | `material_view` | 获取物资列表 |
| POST | `/materials` | `material_create` | 创建物资 |
| GET | `/materials/:id` | `material_view` | 获取物资详情 |
| PUT | `/materials/:id` | `material_edit` | 更新物资 |
| DELETE | `/materials/:id` | `material_delete` | 删除物资 |
| GET | `/materials/export` | `material_view` | 导出物资列表 |
| GET | `/materials/:id/logs` | `material_view` | 获取物资日志 |
| GET | `/materials/unstored` | `material_view` | 获取未入库物资 |
| GET | `/materials/unstored/export` | `material_view` | 导出未入库物资 |
| POST | `/materials/batch` | `material_import` | 批量操作物资 |
| POST | `/materials/batch-create` | `material_import` | 批量创建物资 |
| POST | `/material/materials/import` | `material_create` | 导入物资 |

### 3.3 物资分类 (`/api/material/categories`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/material/categories` | `material_view` | 获取分类列表（树形） |
| GET | `/material/categories/:id` | `material_view` | 获取分类详情 |
| POST | `/material/categories` | `material_create` | 创建分类 |
| PUT | `/material/categories/:id` | `material_edit` | 更新分类 |
| DELETE | `/material/categories/:id` | `material_delete` | 删除分类 |
| POST | `/material/categories/sort` | `material_edit` | 批量更新排序 |

---

## 四、库存管理 (`/api/stock`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/stocks` | `stock_view` | 获取库存列表 |
| GET | `/stocks/alerts` | `stock_view` | 获取库存预警 |
| POST | `/stocks` | `stock_create` | 创建库存记录 |
| GET | `/stocks/:id` | `stock_view` | 获取库存详情 |
| GET | `/stocks/:id/logs` | `stock_view` | 获取库存操作日志 |
| PUT | `/stocks/:id` | `stock_edit` | 更新库存 |
| DELETE | `/stocks/:id` | `stock_delete` | 删除库存 |
| GET | `/stock-logs` | `stock_view` | 获取库存日志列表 |
| POST | `/stocks/:id/in` | `stock_in` | 入库操作 |
| POST | `/stocks/:id/out` | `stock_out` | 出库操作 |
| POST | `/stocks/:id/adjust` | `stock_edit` | 库存调整 |
| DELETE | `/stock-logs/:id` | `stock_delete` | 删除库存日志 |
| GET | `/stocks/export` | `stock_export` | 导出库存列表 |

---

## 五、入库管理 (`/api/inbound`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/inbound/inbound-orders/pending/count` | `inbound_view` | 获取待处理入库单数量 |
| GET | `/inbound/inbound-orders` | `inbound_view` | 获取入库单列表 |
| POST | `/inbound/inbound-orders` | `inbound_create` | 创建入库单 |
| GET | `/inbound/inbound-orders/:id` | `inbound_view` | 获取入库单详情 |
| PUT | `/inbound/inbound-orders/:id` | `inbound_edit` | 更新入库单 |
| DELETE | `/inbound/inbound-orders/:id` | `inbound_delete` | 删除入库单 |
| POST | `/inbound/inbound-orders/:id/approve` | `inbound_approve` | 审批入库单 |
| POST | `/inbound/inbound-orders/:id/reject` | `inbound_approve` | 拒绝入库单 |
| POST | `/inbound/submit` | `inbound_create` | 提交入库单 |
| GET | `/inbound/template` | `inbound_view` | 获取入库单模板 |
| POST | `/inbound/import` | `inbound_create` | 导入入库单 |

---

## 六、领料管理 (`/api/requisitions`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/requisitions/pending/count` | `requisition_view` | 获取待处理领料单数量 |
| GET | `/requisitions` | `requisition_view` | 获取领料单列表 |
| POST | `/requisitions` | `requisition_create` | 创建领料单 |
| GET | `/requisitions/:id` | `requisition_view` | 获取领料单详情 |
| PUT | `/requisitions/:id` | `requisition_edit` | 更新领料单 |
| DELETE | `/requisitions/:id` | `requisition_delete` | 删除领料单 |
| POST | `/requisitions/:id/approve` | `requisition_approve` | 审批领料单 |
| POST | `/requisitions/:id/reject` | `requisition_approve` | 拒绝领料单 |
| POST | `/requisitions/:id/issue` | `requisition_issue` | 发放领料单 |
| GET | `/requisition-items` | `requisition_view` | 获取领料单明细列表 |

---

## 七、物资计划 (`/api/material-plan`)

### 7.1 计划管理

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/material-plan/plans` | `material_plan_view` | 获取计划列表 |
| POST | `/material-plan/plans` | `material_plan_create` | 创建计划 |
| GET | `/material-plan/plans/:id` | `material_plan_view` | 获取计划详情 |
| PUT | `/material-plan/plans/:id` | `material_plan_edit` | 更新计划 |
| DELETE | `/material-plan/plans/:id` | `material_plan_delete` | 删除计划 |
| POST | `/material-plan/plans/:id/submit` | `material_plan_edit` | 提交计划 |
| POST | `/material-plan/plans/:id/approve` | `material_plan_approve` | 审批计划 |
| POST | `/material-plan/plans/:id/reject` | `material_plan_approve` | 拒绝计划 |
| POST | `/material-plan/plans/:id/activate` | `material_plan_approve` | 激活计划 |
| POST | `/material-plan/plans/:id/resubmit` | `material_plan_edit` | 重新提交计划 |
| POST | `/material-plan/plans/:id/cancel` | `material_plan_approve` | 取消计划 |

### 7.2 计划明细

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/material-plan/plans/:id/items` | `material_plan_view` | 获取计划明细列表 |
| POST | `/material-plan/plans/:id/items` | `material_plan_edit` | 添加计划明细 |
| PUT | `/material-plan/items/:id` | `material_plan_edit` | 更新计划明细 |
| DELETE | `/material-plan/items/:id` | `material_plan_edit` | 删除计划明细 |

### 7.3 工作流

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/material-plan/plans/:id/workflow` | `material_plan_view` | 获取计划工作流状态 |
| GET | `/material-plan/plans/:id/approvals` | `material_plan_view` | 获取计划审批记录 |
| GET | `/material-plan/workflow/pending` | `material_plan_approve` | 获取待办任务 |

---

## 八、施工日志 (`/api/construction-logs`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/construction-logs/logs` | `constructionlog_view` | 获取施工日志列表 |
| GET | `/construction-logs/:log_id` | `constructionlog_view` | 获取施工日志详情 |
| POST | `/construction-logs` | `constructionlog_create` | 创建施工日志 |
| PUT | `/construction-logs/:log_id` | `constructionlog_edit` | 更新施工日志 |
| DELETE | `/construction-logs/:log_id` | `constructionlog_delete` | 删除施工日志 |
| POST | `/construction-logs/upload_image` | `constructionlog_create` | 上传日志图片 |

---

## 九、进度管理 (`/api/progress`)

### 9.1 进度计划

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/progress` | - | 获取进度列表 |
| POST | `/progress` | - | 创建进度 |
| GET | `/progress/export` | - | 导出进度 |
| PUT | `/progress/:id` | - | 更新进度 |
| DELETE | `/progress/:id` | - | 删除进度 |
| GET | `/progress/project-schedules` | - | 获取所有项目进度计划 |
| GET | `/progress/project/:id` | - | 获取项目进度计划 |
| PUT | `/progress/project/:id` | - | 更新项目进度计划 |
| DELETE | `/progress/project/:id/schedule` | - | 删除项目进度计划 |
| GET | `/progress/project/:id/exists` | - | 检查进度计划是否存在 |

### 9.2 任务管理

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/progress/project/:id/tasks` | - | 获取项目任务列表 |
| POST | `/progress/project/:id/tasks` | - | 创建任务 |
| PUT | `/progress/tasks/:id` | - | 更新任务 |
| DELETE | `/progress/tasks/:id` | - | 删除任务 |

### 9.3 依赖管理

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/progress/tasks/:id/dependencies` | - | 获取任务依赖 |
| POST | `/progress/tasks/:id/dependencies` | - | 添加任务依赖 |
| DELETE | `/progress/dependencies/:id` | - | 删除任务依赖 |
| POST | `/progress/dependencies/visual/:fromId/:toId` | - | 可视化创建依赖关系 |

### 9.4 资源管理

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/progress/project/:id/resources` | - | 获取项目资源 |
| POST | `/progress/project/:id/resources` | - | 创建资源 |
| PUT | `/progress/project/:id/resources/:resourceId` | - | 更新资源 |
| DELETE | `/progress/project/:id/resources/:resourceId` | - | 删除资源 |
| GET | `/progress/tasks/:id/resources` | - | 获取任务资源 |
| POST | `/progress/tasks/:id/resources` | - | 分配任务资源 |
| DELETE | `/progress/tasks/:id/resources/:resourceId` | - | 移除任务资源 |

### 9.5 AI 辅助

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| POST | `/progress/project/:id/generate-plan` | - | AI 生成进度计划 |
| POST | `/progress/project/:id/aggregate-plan` | - | 聚合子项目计划 |

### 9.6 管理接口

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| POST | `/admin/progress/sync-all` | - | 同步所有项目进度 |
| POST | `/admin/progress/sync/:projectId` | - | 同步指定项目进度 |
| GET | `/admin/progress/sync-status` | - | 获取同步状态 |

---

## 十、工作流管理 (`/api/workflows`)

### 10.1 工作流定义

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/workflows` | - | 获取工作流列表 |
| POST | `/workflows` | - | 创建工作流 |
| GET | `/workflows/:id` | - | 获取工作流详情 |
| PUT | `/workflows/:id` | - | 更新工作流 |
| DELETE | `/workflows/:id` | - | 删除工作流 |
| PUT | `/workflows/:id/activate` | - | 激活工作流 |
| PUT | `/workflows/:id/deactivate` | - | 停用工作流 |

### 10.2 工作流实例

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/workflow-instances` | - | 获取实例列表 |
| GET | `/workflow-instances/:id` | - | 获取实例详情 |
| GET | `/workflow-instances/:id/approvals` | - | 获取审批记录 |
| GET | `/workflow-instances/:id/logs` | - | 获取操作日志 |
| POST | `/workflow-instances/:id/resubmit` | - | 重新提交 |

### 10.3 待办任务

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/workflow-tasks/pending` | - | 获取我的待办任务 |
| GET | `/workflow-tasks/pending/:businessType` | - | 按类型获取待办 |
| POST | `/workflow-tasks/:id/approve` | - | 审批通过 |
| POST | `/workflow-tasks/:id/reject` | - | 审批拒绝 |
| POST | `/workflow-tasks/:id/return` | - | 退回 |
| POST | `/workflow-tasks/:id/comment` | - | 评论 |

---

## 十一、通知管理 (`/api/notifications`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/notifications` | - | 获取通知列表 |
| GET | `/notifications/count` | - | 获取未读通知数量 |
| PUT | `/notifications/:id/read` | - | 标记通知为已读 |
| PUT | `/notifications/read-all` | - | 全部标记为已读 |
| DELETE | `/notifications/:id` | - | 删除通知 |
| DELETE | `/notifications` | - | 清空所有通知 |

---

## 十二、系统管理 (`/api/system`)

### 12.1 系统信息

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/system/info` | `system_statistics` | 获取系统信息 |
| GET | `/system/stats` | - | 获取系统统计 |
| GET | `/system/recent-activities` | `system_activities` | 获取最近活动 |
| GET | `/system/material-category-stats` | `system_statistics` | 获取物资分类统计 |
| GET | `/system/project-material-stats` | `system_statistics` | 获取项目物资统计 |

### 12.2 系统设置

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/system/settings` | `system_config` | 获取系统设置 |
| PUT | `/system/settings` | `system_config` | 更新系统设置 |
| GET | `/system/public/settings` | - | 获取公开设置 |

### 12.3 系统日志

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/system/logs` | - | 获取系统日志 |
| POST | `/system/logs/clear` | `system_log` | 清空系统日志 |
| DELETE | `/system/logs` | `system_log` | 删除系统日志 |

### 12.4 数据库备份

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| POST | `/system/backup` | `system_backup` | 创建备份 |
| GET | `/system/backup` | `system_backup` | 获取备份列表 |
| GET | `/system/backup/history` | `system_backup` | 获取备份历史 |
| GET | `/system/backup/:backup_name/download` | `system_backup` | 下载备份 |
| GET | `/system/backup/download` | `system_backup` | 下载备份（旧接口） |
| DELETE | `/system/backup/:backup_name` | `system_backup` | 删除备份 |
| POST | `/system/backup/delete` | `system_backup` | 删除备份（旧接口） |
| POST | `/system/backup/create` | `system_backup` | 创建备份（报表接口） |
| POST | `/system/backup/restore` | `system_backup` | 恢复备份 |

### 12.5 报表管理

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/system/reports/dashboard` | `system_statistics` | 获取仪表板数据 |
| GET | `/system/reports` | `system_statistics` | 获取报表列表 |
| POST | `/system/reports/generate` | `system_statistics` | 生成报表 |
| GET | `/system/reports/download` | `system_statistics` | 下载报表 |
| GET | `/system/reports/:report_name/download` | `system_statistics` | 下载指定报表 |
| POST | `/system/reports/delete` | `system_backup` | 删除报表 |
| DELETE | `/system/reports/:report_name` | `system_backup` | 删除指定报表 |

### 12.6 AI 分析

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| POST | `/system/ai/analyze` | `system_statistics` | AI 智能分析 |
| GET | `/system/ai/suggestions` | `system_statistics` | 获取 AI 建议 |
| GET | `/system/ai/insights` | `system_statistics` | 获取数据洞察 |
| GET | `/system/ai/history` | `system_statistics` | 获取分析历史 |
| GET | `/system/ai/history/:id` | `system_statistics` | 获取分析历史详情 |
| DELETE | `/system/ai/history/:id` | `system_statistics` | 删除分析历史 |
| GET | `/system/ai/stats` | `system_statistics` | 获取 AI 统计 |
| GET | `/system/ai/config` | `system_config` | 获取 AI 配置 |
| POST | `/system/ai/config` | `system_config` | 更新 AI 配置 |
| GET | `/system/ai/status` | - | 获取 AI 状态 |

---

## 十三、文件上传 (`/api/upload`)

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| POST | `/upload/image` | - | 上传单张图片 |
| POST | `/upload/images` | - | 上传多张图片 |

---

## 十四、健康检查

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/health` | - | 健康检查端点 |

---

## 分页参数

所有列表接口支持以下分页参数：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码（从1开始） |
| page_size | int | 否 | 20 | 每页数量（最大100） |

## 过滤参数

| 参数 | 类型 | 说明 |
|------|------|------|
| status | string | 状态过滤 |
| project_id | int | 项目ID过滤 |
| search | string | 搜索关键词 |
| start_date | string | 开始日期 (YYYY-MM-DD) |
| end_date | string | 结束日期 (YYYY-MM-DD) |

## 排序参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| sort_by | string | id | 排序字段 |
| sort_order | string | desc | 排序方向 (asc/desc) |

---

**最后更新时间**: 2026-02-02
