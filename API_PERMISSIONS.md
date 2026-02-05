# 后端API权限完整列表

> 基于后端代码自动生成
> 生成时间: 2026-02-05

## 权限模块概览

| 序号 | 模块 | 权限数量 | 说明 |
|------|------|----------|------|
| 1 | 用户管理 | 4 | 用户增删改查 |
| 2 | 角色管理 | 5 | 角色和权限管理 |
| 3 | 项目管理 | 4 | 项目管理 |
| 4 | 物资管理 | 5 | 物资基础管理 |
| 5 | 物资主数据 | 5 | 物资主数据管理 |
| 6 | 物资计划 | 5 | 物资计划及审核 |
| 7 | 库存管理 | 8 | 库存操作 |
| 8 | 库存日志 | 2 | 库存日志管理 |
| 9 | 入库管理 | 6 | 入库单管理 |
| 10 | 出库管理 | 7 | 出库单管理 |
| 11 | 施工日志 | 4 | 施工日志管理 |
| 12 | 进度管理 | 5 | 进度任务管理 |
| 13 | 审计日志 | 1 | 操作审计 |
| 14 | AI智能体 | 5 | AI功能权限 |
| 15 | 系统管理 | 5 | 系统配置和日志 |
| 16 | 工作流 | 1 | 工作流管理(仅admin) |

**总计：16 个模块，77 个权限点**

---

## 完整权限列表

### 1. 用户管理模块 (user)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `user_view` | 查看用户 | `GET /auth/users`, `GET /auth/users/:id` | 查看用户列表和详情 |
| `user_create` | 创建用户 | `POST /auth/users` | 创建新用户 |
| `user_edit` | 编辑用户 | `PUT /auth/users/:id`, `POST /auth/users/:id/reset-password` | 修改用户信息和重置密码 |
| `user_delete` | 删除用户 | `DELETE /auth/users/:id` | 删除用户 |

### 2. 角色管理模块 (role)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `role_view` | 查看角色 | `GET /auth/roles`, `GET /auth/roles/:id`, `GET /auth/permissions` | 查看角色列表、详情和权限列表 |
| `role_create` | 创建角色 | `POST /auth/roles` | 创建新角色 |
| `role_edit` | 编辑角色 | `PUT /auth/roles/:id` | 修改角色信息 |
| `role_delete` | 删除角色 | `DELETE /auth/roles/:id` | 删除角色 |
| `role_assign_permissions` | 分配权限 | `POST /auth/roles/:id/permissions` | 为角色分配权限 |

### 3. 项目管理模块 (project)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `project_view` | 查看项目 | `GET /projects/projects`, `GET /projects/projects/:id`, `GET /projects/projects/:id/members`, `GET /projects/projects/:id/tree`, `GET /projects/projects/:id/children` | 查看项目列表、详情、成员、树形结构和子项目 |
| `project_create` | 创建项目 | `POST /projects/projects` | 创建新项目 |
| `project_edit` | 编辑项目 | `PUT /projects/projects/:id`, `POST /projects/projects/:id/members`, `POST /projects/projects/:id/aggregate-progress` | 修改项目信息、添加成员、汇总进度 |
| `project_delete` | 删除项目 | `DELETE /projects/projects/:id`, `DELETE /projects/projects/:id/members/:user_id` | 删除项目和移除成员 |

### 4. 物资管理模块 (material)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `material_view` | 查看物资 | `GET /materials/materials`, `GET /materials/materials/:id`, `GET /materials/materials/:id/logs`, `GET /materials/materials/export`, `GET /materials/categories`, `GET /materials/categories/:id` | 查看物资列表、详情、日志、分类和导出 |
| `material_create` | 创建物资 | `POST /materials/materials`, `POST /materials/categories`, `POST /materials/materials/batch-create` | 创建物资、分类和批量创建 |
| `material_edit` | 编辑物资 | `PUT /materials/materials/:id`, `PUT /materials/categories/:id`, `POST /materials/categories/sort` | 修改物资和分类信息、排序 |
| `material_delete` | 删除物资 | `DELETE /materials/materials/:id`, `DELETE /materials/categories/:id` | 删除物资和分类 |
| `material_import` | 导入物资 | `POST /materials/material/materials/import`, `POST /materials/materials/batch` | 批量导入物资 |

### 5. 物资主数据模块 (material_master)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `material_view` | 查看物资主数据 | `GET /api/materials/master`, `GET /api/materials/master/:id`, `GET /api/materials/master/project`, `GET /api/materials/master/categories` | 查看物资主数据 |
| `material_create` | 创建物资主数据 | `POST /api/materials/master` | 创建物资主数据 |
| `material_edit` | 编辑物资主数据 | `PUT /api/materials/master/:id` | 修改物资主数据 |
| `material_delete` | 删除物资主数据 | `DELETE /api/materials/master/:id` | 删除物资主数据 |

**注意**：物资主数据模块复用物资管理权限。

### 6. 物资计划模块 (material_plan)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `material_plan_view` | 查看物资计划 | `GET /material-plan/plans`, `GET /material-plan/plans/:id`, `GET /material-plan/plans/:id/items`, `GET /material-plan/plans/:id/workflow`, `GET /material-plan/plans/:id/approvals` | 查看计划列表、详情、项目、工作流状态和审批记录 |
| `material_plan_create` | 创建物资计划 | `POST /material-plan/plans` | 创建新计划 |
| `material_plan_edit` | 编辑物资计划 | `PUT /material-plan/plans/:id`, `POST /material-plan/plans/:id/submit`, `POST /material-plan/plans/:id/resubmit`, `POST /material-plan/plans/:id/items`, `PUT /material-plan/items/:id`, `DELETE /material-plan/items/:id`, `POST /material-plan/plans/:id/sync-materials` | 修改计划、提交、重新提交、管理计划项、同步物资 |
| `material_plan_delete` | 删除物资计划 | `DELETE /material-plan/plans/:id` | 删除计划 |
| `material_plan_approve` | 审核物资计划 | `POST /material-plan/plans/:id/approve`, `POST /material-plan/plans/:id/reject`, `POST /material-plan/plans/:id/activate`, `POST /material-plan/plans/:id/cancel`, `GET /material-plan/workflow/pending` | 审核通过、拒绝、激活、取消、查看待办 |

### 7. 库存管理模块 (stock)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `stock_view` | 查看库存 | `GET /stocks/stocks`, `GET /stocks/stocks/:id`, `GET /stocks/stocks/:id/logs`, `GET /stocks/stock-logs`, `GET /stocks/stocks/alerts` | 查看库存列表、详情、日志和预警 |
| `stock_create` | 创建库存 | `POST /stocks/stocks` | 创建库存记录 |
| `stock_edit` | 编辑库存 | `PUT /stocks/stocks/:id`, `POST /stocks/stocks/:id/adjust` | 修改库存信息和调整 |
| `stock_delete` | 删除库存 | `DELETE /stocks/stocks/:id`, `DELETE /stocks/stock-logs/:id` | 删除库存和日志 |
| `stock_in` | 库存入库 | `POST /stocks/stocks/:id/in` | 入库操作 |
| `stock_out` | 库存出库 | `POST /stocks/stocks/:id/out` | 出库操作 |
| `stock_export` | 导出库存 | `GET /stocks/stocks/export` | 导出库存数据 |
| `stock_alerts` | 库存预警 | `GET /stocks/stocks/alerts` | 查看库存预警 |

### 8. 库存日志模块 (stocklog)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `stocklog_view` | 查看库存日志 | `GET /stocks/stock-logs`, `GET /stocks/stocks/:id/logs` | 查看库存日志 |
| `stocklog_delete` | 删除库存日志 | `DELETE /stocks/stock-logs/:id` | 删除库存日志 |

### 9. 入库管理模块 (inbound)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `inbound_view` | 查看入库单 | `GET /inbound/inbound-orders`, `GET /inbound/inbound-orders/:id`, `GET /inbound/inbound-orders/:id/workflow-history`, `GET /inbound/inbound/template`, `GET /inbound/inbound-orders/pending/count` | 查看入库单列表、详情、工作流历史、模板和待办数量 |
| `inbound_create` | 创建入库单 | `POST /inbound/inbound-orders`, `POST /inbound/inbound/submit`, `POST /inbound/inbound/import` | 创建、提交和导入入库单 |
| `inbound_edit` | 编辑入库单 | `PUT /inbound/inbound-orders/:id` | 修改入库单 |
| `inbound_delete` | 删除入库单 | `DELETE /inbound/inbound-orders/:id` | 删除入库单 |
| `inbound_approve` | 审核入库单 | `POST /inbound/inbound-orders/:id/approve`, `POST /inbound/inbound-orders/:id/reject` | 审核通过或拒绝入库单 |
| `inbound_export` | 导出入库单 | (未在代码中找到对应端点) | 导出入库单数据 |

### 10. 出库管理模块 (requisition)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `requisition_view` | 查看出库单 | `GET /requisitions/requisitions`, `GET /requisitions/requisitions/:id`, `GET /requisitions/requisition-items`, `GET /requisitions/requisitions/pending/count` | 查看出库单列表、详情、项目和待办数量 |
| `requisition_create` | 创建出库单 | `POST /requisitions/requisitions` | 创建新出库单 |
| `requisition_edit` | 编辑出库单 | `PUT /requisitions/requisitions/:id` | 修改出库单 |
| `requisition_delete` | 删除出库单 | `DELETE /requisitions/requisitions/:id` | 删除出库单 |
| `requisition_approve` | 审核出库单 | `POST /requisitions/requisitions/:id/approve`, `POST /requisitions/requisitions/:id/reject` | 审核通过或拒绝出库单 |
| `requisition_issue` | 发货 | `POST /requisitions/requisitions/:id/issue` | 标记为已发货 |
| `requisition_export` | 导出出库单 | (未在代码中找到对应端点) | 导出出库单数据 |

### 11. 施工日志模块 (construction_log)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `constructionlog_view` | 查看日志 | `GET /construction/logs`, `GET /construction/:log_id` | 查看施工日志列表和详情 |
| `constructionlog_create` | 创建日志 | `POST /construction/`, `POST /construction/upload_image` | 创建日志和上传图片 |
| `constructionlog_edit` | 编辑日志 | `PUT /construction/:log_id` | 修改日志内容 |
| `constructionlog_delete` | 删除日志 | `DELETE /construction/:log_id` | 删除日志 |
| `constructionlog_export` | 导出日志 | (未在代码中找到对应端点) | 导出施工日志数据 |

**注意**：权限代码使用 `constructionlog` 前缀而非 `construction_log`。

### 12. 进度管理模块 (progress)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `progress_view` | 查看进度 | `GET /progress/`, `GET /progress/project/:id`, `GET /progress/project/:id/tasks`, `GET /progress/tasks/:id/dependencies` | 查看进度列表、项目进度、任务和依赖关系 |
| `progress_create` | 创建进度 | `POST /progress/`, `POST /progress/project/:id/tasks`, `POST /progress/project/:id/generate-plan` | 创建进度任务、任务和生成计划 |
| `progress_edit` | 编辑进度 | `PUT /progress/:id`, `PUT /progress/tasks/:id`, `PUT /progress/tasks/:id/position`, `PUT /progress/project/:id` | 修改进度、任务、位置和项目 |
| `progress_delete` | 删除进度 | `DELETE /progress/:id`, `DELETE /progress/tasks/:id`, `DELETE /progress/dependencies/:id`, `DELETE /progress/project/:id/schedule` | 删除进度、任务、依赖和计划 |
| `progress_export` | 导出进度 | `GET /progress/export` | 导出进度数据 |

**注意**：进度模块的路由在代码中未明确添加 `PermissionMiddleware`，建议添加。

### 13. 审计日志模块 (audit)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `audit_view` | 查看审计日志 | `GET /audit/operation-logs`, `GET /audit/operation-logs/statistics`, `GET /audit/operation-logs/export`, `GET /audit/operation-logs/:id`, `GET /audit/operation-logs/resource/:resource_type/:resource_id`, `GET /audit/operation-logs/user/:user_id` | 查看审计日志、统计、导出、详情、资源日志和用户日志 |

### 14. AI智能体模块 (ai_agent)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `ai_agent_view` | 查看AI能力 | `GET /agent/capabilities`, `GET /agent/validate` | 查看AI能力和验证 |
| `ai_agent_query` | AI查询 | `POST /agent/query` | 使用AI查询功能 |
| `ai_agent_operate` | AI操作 | `POST /agent/operate` | 执行AI操作 |
| `ai_agent_workflow` | AI工作流 | `POST /agent/workflow` | AI工作流操作 |
| `ai_agent_logs` | AI日志 | `GET /agent/logs` | 查看AI操作日志 |

**注意**：AI智能体的大部分端点使用 `AgentPermissionMiddleware` 中间件而非 `PermissionMiddleware`。

### 15. 系统管理模块 (system)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `system_log` | 查看系统日志 | `POST /system/logs/clear`, `DELETE /system/logs` | 清除和删除系统日志 |
| `system_backup` | 数据备份 | `POST /system/backup`, `GET /system/backup/history`, `GET /system/backup`, `GET /system/backup/:backup_name/download`, `DELETE /system/backup/:backup_name`, `GET /system/backup/download`, `POST /system/backup/delete`, `POST /system/backup/restore`, `POST /system/backup/create`, `POST /system/reports/delete`, `DELETE /system/reports/:report_name` | 创建、查看、下载、删除和恢复备份，删除报告 |
| `system_config` | 系统配置 | `GET /system/settings`, `PUT /system/settings`, `GET /system/config`, `POST /system/config` | 查看和修改系统设置、AI分析配置 |
| `system_statistics` | 系统统计 | `GET /system/info`, `GET /system/material-category-stats`, `GET /system/project-material-stats`, `GET /system/reports/dashboard`, `GET /system/reports`, `POST /system/reports/generate`, `GET /system/reports/download`, `GET /system/reports/:report_name/download`, `POST /system/analyze`, `GET /system/suggestions`, `GET /system/insights`, `GET /system/history`, `GET /system/history/:id`, `DELETE /system/history/:id`, `GET /system/stats` | 查看系统信息、统计、报告、AI分析等 |
| `system_activities` | 系统动态 | `GET /system/recent-activities` | 查看系统动态和活动 |

### 16. 工作流管理模块 (workflow)

| 权限代码 | 权限名称 | API端点 | 说明 |
|---------|---------|---------|------|
| `admin` | 工作流管理 | `GET /workflows`, `POST /workflows`, `GET /workflows/:id`, `PUT /workflows/:id`, `DELETE /workflows/:id`, `PUT /workflows/:id/activate`, `PUT /workflows/:id/deactivate` | 工作流的全部操作都需要admin权限 |

**注意**：工作流实例和任务相关的端点只要求登录认证，不要求特定权限。

---

## 权限命名规范

系统采用 `{模块}_{操作}` 的命名规范：

### 操作类型

| 操作后缀 | 说明 | 示例 |
|---------|------|------|
| `_view` | 查看权限 | `material_view`, `stock_view` |
| `_create` | 创建权限 | `material_create`, `stock_create` |
| `_edit` | 编辑权限 | `material_edit`, `stock_edit` |
| `_delete` | 删除权限 | `material_delete`, `stock_delete` |
| `_approve` | 审核权限 | `material_plan_approve`, `requisition_approve` |
| `_import` | 导入权限 | `material_import` |
| `_export` | 导出权限 | `stock_export`, `progress_export` |
| `_in` | 入库操作 | `stock_in` |
| `_out` | 出库操作 | `stock_out` |
| `_issue` | 发货操作 | `requisition_issue` |

### 特殊权限

| 权限代码 | 说明 |
|---------|------|
| `admin` | 管理员权限，拥有所有权限 |
| `system_log` | 系统日志管理 |
| `system_backup` | 数据备份管理 |
| `system_config` | 系统配置管理 |
| `system_statistics` | 系统统计分析 |
| `system_activities` | 系统动态查看 |

---

## 需要注意的问题

### 1. 权限不一致问题

- **施工日志权限**：代码中使用 `constructionlog_*` 前缀，建议统一为 `construction_log_*`
  - 当前：`constructionlog_view`, `constructionlog_create`, `constructionlog_edit`, `constructionlog_delete`
  - 建议：`construction_log_view`, `construction_log_create`, `construction_log_edit`, `construction_log_delete`

### 2. 缺少权限中间件

以下模块的路由在代码中未明确添加 `PermissionMiddleware`，建议补充：

- **进度管理模块** (`/progress/*`) - 所有路由
- **通知模块** (`/notification/*`) - 所有路由
- **上传模块** (`/upload/*`) - 所有路由

### 3. 工作流权限过于严格

工作流管理模块的所有端点都要求 `admin` 权限，建议细分：
- `workflow_view` - 查看工作流定义
- `workflow_create` - 创建工作流
- `workflow_edit` - 编辑工作流
- `workflow_delete` - 删除工作流
- `workflow_activate` - 激活/停用工作流

### 4. AI智能体权限

AI智能体的大部分端点使用自定义的 `AgentPermissionMiddleware`，与系统的 `PermissionMiddleware` 不一致，建议统一。

---

## API路由注册文件清单

| 模块 | 路由注册文件 |
|------|-------------|
| 用户与角色 | `/internal/api/auth/handler.go` |
| 项目 | `/internal/api/project/handler.go` |
| 物资 | `/internal/api/material/handler.go`, `/internal/api/material/category_handler.go` |
| 物资主数据 | `/internal/api/material_master/routes.go` |
| 物资计划 | `/internal/api/material_plan/handler.go` |
| 库存 | `/internal/api/stock/handler.go` |
| 入库 | `/internal/api/inbound/handler.go` |
| 出库 | `/internal/api/requisition/handler.go` |
| 施工日志 | `/internal/api/construction_log/handler.go` |
| 进度 | `/internal/api/progress/routes.go` |
| 审计 | `/internal/api/audit/routes.go` |
| AI智能体 | `/internal/api/agent/routes.go` |
| 系统 | `/internal/api/system/handler.go`, `/internal/api/system/ai_handler.go`, `/internal/api/system/report_handler.go` |
| 工作流 | `/internal/api/workflow/register.go` |
| 通知 | `/internal/api/notification/handler.go` |
| 上传 | `/internal/api/upload/handler.go` |

---

## 使用建议

### 为新功能添加权限

1. **后端**：在路由注册时使用 `auth.PermissionMiddleware(db, "permission_key")`
2. **权限列表**：在 `auth/handler.go` 的 `/permissions` 端点中添加权限定义
3. **前端**：在 `System.vue` 的 `permissionTree` 中添加权限节点
4. **文档**：更新本文档和 `PERMISSION_CONFIG.md`

### 示例代码

```go
// 后端路由注册
r.GET("/api/resource", auth.PermissionMiddleware(db, "resource_view"), handler)
r.POST("/api/resource", auth.PermissionMiddleware(db, "resource_create"), handler)
r.PUT("/api/resource/:id", auth.PermissionMiddleware(db, "resource_edit"), handler)
r.DELETE("/api/resource/:id", auth.PermissionMiddleware(db, "resource_delete"), handler)
```

---

**文档版本**: 1.0
**最后更新**: 2026-02-05
