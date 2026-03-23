# 系统权限配置文档（基于实际API）

## 更新时间
2025-02-05

## 权限模块概览

系统共包含 **15 个模块**，**77 个唯一权限点**，覆盖 **167 个 API 端点**。

| 序号 | 模块 | 权限数量 | API端点数量 |
|------|------|----------|-------------|
| 1 | 用户管理 | 9 | 10 |
| 2 | 项目管理 | 4 | 11 |
| 3 | 物资管理 | 5 | 13 |
| 4 | 物资计划 | 5 | 18 |
| 5 | 库存管理 | 8 | 9 |
| 6 | 库存日志 | 2 | 2 |
| 7 | 入库管理 | 6 | 9 |
| 8 | 出库管理 | 7 | 10 |
| 9 | 施工日志 | 5 | 5 |
| 10 | 进度管理 | 5 | 35 |
| 11 | 审计日志 | 2 | 7 |
| 12 | AI智能体 | 5 | 6 |
| 13 | 系统管理 | 5 | 18 |
| 14 | 工作流管理 | 13 | 17 |
| 15 | 通知/上传 | 0（仅登录） | 8 |
| **总计** | **15** | **77** | **167** |

---

## 完整权限列表

### 1. 用户管理模块 (9个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `user_view` | 查看用户 | `GET /api/auth/users` |
| `user_create` | 创建用户 | `POST /api/auth/users` |
| `user_edit` | 编辑用户 | `PUT /api/auth/users/:id` |
| `user_delete` | 删除用户 | `DELETE /api/auth/users/:id` |
| `role_view` | 查看角色 | `GET /api/auth/roles` |
| `role_create` | 创建角色 | `POST /api/auth/roles` |
| `role_edit` | 编辑角色 | `PUT /api/auth/roles/:id` |
| `role_delete` | 删除角色 | `DELETE /api/auth/roles/:id` |
| `role_assign_permissions` | 分配权限 | `POST /api/auth/roles/:id/permissions` |

**说明**：用户和角色管理合并为一个模块。

---

### 2. 项目管理模块 (4个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `project_view` | 查看项目 | `GET /api/project/projects` |
| `project_create` | 创建项目 | `POST /api/project/projects` |
| `project_edit` | 编辑项目 | `PUT /api/project/projects/:id` |
| `project_delete` | 删除项目 | `DELETE /api/project/projects/:id` |

**说明**：项目管理包括项目本身和项目成员管理。

---

### 3. 物资管理模块 (5个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `material_view` | 查看物资 | `GET /api/material/materials` |
| `material_create` | 创建物资 | `POST /api/material/materials` |
| `material_edit` | 编辑物资 | `PUT /api/material/materials/:id` |
| `material_delete` | 删除物资 | `DELETE /api/material/materials/:id` |
| `material_import` | 导入物资 | `POST /api/material/materials/batch-create` |

---

### 4. 物资计划模块 (5个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `material_plan_view` | 查看物资计划 | `GET /api/material-plan/plans` |
| `material_plan_create` | 创建物资计划 | `POST /api/material-plan/plans` |
| `material_plan_edit` | 编辑物资计划 | `PUT /api/material-plan/plans/:id` |
| `material_plan_delete` | 删除物资计划 | `DELETE /api/material-plan/plans/:id` |
| `material_plan_approve` | 审核物资计划 | `POST /api/material-plan/plans/:id/approve` |

---

### 5. 库存管理模块 (8个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `stock_view` | 查看库存 | `GET /api/stock/stocks` |
| `stock_create` | 创建库存 | `POST /api/stock/stocks` |
| `stock_edit` | 编辑库存 | `PUT /api/stock/stocks/:id` |
| `stock_delete` | 删除库存 | `DELETE /api/stock/stocks/:id` |
| `stock_in` | 库存入库 | `POST /api/stock/stocks/:id/in` |
| `stock_out` | 库存出库 | `POST /api/stock/stocks/:id/out` |
| `stock_export` | 导出库存 | `GET /api/stock/stocks/export` |
| `stock_alerts` | 库存预警 | `GET /api/stock/stocks/alerts` |

---

### 6. 库存日志模块 (2个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `stocklog_view` | 查看库存日志 | `GET /api/stock/stock-logs` |
| `stocklog_delete` | 删除库存日志 | `DELETE /api/stock/stock-logs/:id` |

---

### 7. 入库管理模块 (6个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `inbound_view` | 查看入库单 | `GET /api/inbound/inbound-orders` |
| `inbound_create` | 创建入库单 | `POST /api/inbound/inbound-orders` |
| `inbound_edit` | 编辑入库单 | `PUT /api/inbound/inbound-orders/:id` |
| `inbound_delete` | 删除入库单 | `DELETE /api/inbound/inbound-orders/:id` |
| `inbound_approve` | 审核入库单 | `POST /api/inbound/inbound-orders/:id/approve` |
| `inbound_export` | 导出入库单 | `GET /api/inbound/inbound-orders/export` |

---

### 8. 出库管理模块 (7个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `requisition_view` | 查看出库单 | `GET /api/requisition/requisitions` |
| `requisition_create` | 创建出库单 | `POST /api/requisition/requisitions` |
| `requisition_edit` | 编辑出库单 | `PUT /api/requisition/requisitions/:id` |
| `requisition_delete` | 删除出库单 | `DELETE /api/requisition/requisitions/:id` |
| `requisition_approve` | 审核出库单 | `POST /api/requisition/requisitions/:id/approve` |
| `requisition_issue` | 发货 | `POST /api/requisition/requisitions/:id/issue` |
| `requisition_export` | 导出出库单 | `GET /api/requisition/requisitions/export` |

---

### 9. 施工日志模块 (5个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `construction_log_view` | 查看日志 | `GET /api/construction_log/logs` |
| `construction_log_create` | 创建日志 | `POST /api/construction_log/logs` |
| `construction_log_edit` | 编辑日志 | `PUT /api/construction_log/logs/:id` |
| `construction_log_delete` | 删除日志 | `DELETE /api/construction_log/logs/:id` |
| `construction_log_export` | 导出日志 | `GET /api/construction_log/logs/export` |

---

### 10. 进度管理模块 (5个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `progress_view` | 查看进度 | `GET /api/progress` |
| `progress_create` | 创建进度 | `POST /api/progress` |
| `progress_edit` | 编辑进度 | `PUT /api/progress/:id` |
| `progress_delete` | 删除进度 | `DELETE /api/progress/:id` |
| `progress_export` | 导出进度 | `GET /api/progress/export` |

**说明**：包含35个子功能（任务、依赖、资源、AI生成等）。

---

### 11. 审计日志模块 (2个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `audit_view` | 查看审计日志 | `GET /api/audit/operation-logs` |
| `audit_admin` | 审计日志管理 | `DELETE /api/audit/operation-logs/cleanup` |

---

### 12. AI智能体模块 (5个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `ai_agent_view` | 查看AI能力 | `GET /api/agent/capabilities` |
| `ai_agent_query` | AI查询 | `POST /api/agent/query` |
| `ai_agent_operate` | AI操作 | `POST /api/agent/operate` |
| `ai_agent_workflow` | AI工作流 | `POST /api/agent/workflow` |
| `ai_agent_logs` | AI日志 | `GET /api/agent/logs` |

---

### 13. 系统管理模块 (5个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `system_log` | 系统日志 | `GET /api/system/logs` |
| `system_backup` | 数据备份 | `POST /api/system/backup` |
| `system_config` | 系统配置 | `PUT /api/system/settings` |
| `system_statistics` | 系统统计 | `GET /api/system/statistics` |
| `system_activities` | 系统动态 | `GET /api/system/activities` |

**说明**：包含18个子功能（备份、配置、报表、AI分析等）。

---

### 14. 工作流管理模块 (13个权限)

| 权限代码 | 权限名称 | API端点示例 |
|---------|---------|-------------|
| `workflow_view` | 查看工作流 | `GET /api/workflows` |
| `workflow_create` | 创建工作流 | `POST /api/workflows` |
| `workflow_edit` | 编辑工作流 | `PUT /api/workflows/:id` |
| `workflow_delete` | 删除工作流 | `DELETE /api/workflows/:id` |
| `workflow_activate` | 激活/停用工作流 | `PUT /api/workflows/:id/activate` |
| `workflow_instance_view` | 查看工作流实例 | `GET /api/workflow-instances` |
| `workflow_instance_resubmit` | 重新提交实例 | `POST /api/workflow-instances/:id/resubmit` |
| `workflow_task_view` | 查看工作流任务 | `GET /api/workflow-tasks/pending` |
| `workflow_task_approve` | 审批任务 | `POST /api/workflow-tasks/:id/approve` |
| `workflow_task_reject` | 拒绝任务 | `POST /api/workflow-tasks/:id/reject` |
| `workflow_task_delegate` | 委派任务 | `POST /api/workflow-tasks/:id/return` |
| `workflow_task_comment` | 评论任务 | `POST /api/workflow-tasks/:id/comment` |
| `workflow_log_view` | 查看流程日志 | `GET /api/workflow-instances/:id/logs` |

---

### 15. 通知/文件上传模块 (无需特定权限)

| API端点 | 权限要求 |
|---------|---------|
| `GET /api/notification/notifications` | 仅需登录 |
| `POST /api/upload/files` | 仅需登录 |

**说明**：通知和文件上传仅需要用户登录即可访问。

---

## 特殊权限

### admin 权限

| API端点 | 说明 |
|---------|------|
| `POST /api/admin/progress/sync-all` | 同步所有项目进度 |
| `POST /api/admin/progress/sync/:projectId` | 同步指定项目进度 |
| `GET /api/admin/progress/sync-status` | 获取同步状态 |
| `DELETE /api/audit/operation-logs/cleanup` | 清理审计日志 |

**说明**：admin权限用于系统管理和维护操作。

---

## 权限命名规范

权限代码遵循以下命名规范：

```
{模块}_{操作}
```

### 模块前缀

| 前缀 | 模块名称 |
|------|---------|
| `user_` | 用户管理 |
| `role_` | 角色管理 |
| `project_` | 项目管理 |
| `material_` | 物资管理 |
| `material_plan_` | 物资计划 |
| `stock_` | 库存管理 |
| `inbound_` | 入库管理 |
| `requisition_` | 出库管理 |
| `construction_log_` | 施工日志 |
| `progress_` | 进度管理 |
| `audit_` | 审计日志 |
| `ai_agent_` | AI智能体 |
| `system_` | 系统管理 |
| `workflow_` | 工作流管理 |

### 操作后缀

| 后缀 | 含义 |
|------|------|
| `_view` | 查看/列表 |
| `_create` | 创建/新建 |
| `_edit` | 编辑/修改 |
| `_delete` | 删除 |
| `_approve` | 审批/审核 |
| `_export` | 导出 |
| `_import` | 导入 |
| `_in` | 入库 |
| `_out` | 出库 |

---

## 权限检查实现

### 后端权限中间件

后端使用 `PermissionMiddleware` 进行权限检查：

```go
// 示例：需要物资查看权限
r.GET("/materials", auth.PermissionMiddleware(db, "material_view"), handler)

// 示例：需要物资创建权限
r.POST("/materials", auth.PermissionMiddleware(db, "material_create"), handler)

// 示例：需要管理员权限
r.POST("/admin/progress/sync-all", auth.PermissionMiddleware(db, "admin"), handler)
```

### 前端权限检查

前端使用 `hasPermission` 方法进行权限检查：

```javascript
// 在组件中使用
<el-button
  v-if="authStore.hasPermission('user_create')"
  @click="handleAddUser">
  添加用户
</el-button>

// 在路由守卫中使用
meta: {
  permissions: ['user_view', 'user_edit', 'user_delete']
}
```

---

## 权限分配建议

### 管理员角色 (admin)

拥有所有权限，包括特殊的 `admin` 权限。

### 项目经理角色

| 模块 | 推荐权限 |
|------|---------|
| 项目管理 | `project_view`, `project_create`, `project_edit` |
| 进度管理 | `progress_view`, `progress_create`, `progress_edit` |
| 施工日志 | `construction_log_view`, `construction_log_create`, `construction_log_edit` |
| 物资计划 | `material_plan_view`, `material_plan_create`, `material_plan_edit`, `material_plan_approve` |

### 仓库管理员角色

| 模块 | 推荐权限 |
|------|---------|
| 库存管理 | `stock_view`, `stock_create`, `stock_edit`, `stock_in`, `stock_out`, `stock_alerts` |
| 入库管理 | `inbound_view`, `inbound_create`, `inbound_edit`, `inbound_approve` |
| 出库管理 | `requisition_view`, `requisition_create`, `requisition_edit`, `requisition_approve`, `requisition_issue` |
| 物资管理 | `material_view`, `material_create`, `material_edit`, `material_import` |

### 普通员工角色

| 模块 | 推荐权限 |
|------|---------|
| 项目管理 | `project_view` |
| 进度管理 | `progress_view` |
| 施工日志 | `construction_log_view`, `construction_log_create` |
| 出库管理 | `requisition_view`, `requisition_create` |

---

## 权限统计汇总

### 按模块统计

| 模块 | 权限数 | 占比 |
|------|--------|------|
| 工作流管理 | 13 | 19.1% |
| 库存管理 | 8 | 11.8% |
| 出库管理 | 7 | 10.3% |
| 进度管理 | 5 | 7.4% |
| 物资计划 | 5 | 7.4% |
| 物资管理 | 5 | 7.4% |
| 施工日志 | 5 | 7.4% |
| 系统管理 | 5 | 7.4% |
| AI智能体 | 5 | 7.4% |
| 用户管理 | 5 | 7.4% |
| 入库管理 | 5 | 7.4% |
| 项目管理 | 4 | 5.9% |
| 物资主数据 | 4 | 5.9% |
| 审计日志 | 2 | 2.9% |
| **总计** | **68** | **100%** |

### 按操作类型统计

| 操作 | 权限数 | 占比 |
|------|--------|------|
| 查看 (view) | 14 | 20.6% |
| 创建 (create) | 13 | 19.1% |
| 编辑 (edit) | 13 | 19.1% |
| 删除 (delete) | 11 | 16.2% |
| 审批 (approve) | 6 | 8.8% |
| 导出 (export) | 6 | 8.8% |
| 其他 | 5 | 7.4% |
| **总计** | **68** | **100%** |

---

## 更新日志

### 2025-02-05
- 基于实际API端点重新生成权限列表
- 删除了已废弃的API权限（debug-token, system/logs等）
- 更新了权限数量统计
- 添加了权限分配建议

---

## 注意事项

1. **权限名称唯一性**：每个权限代码在系统中是唯一的
2. **权限继承性**：admin角色自动拥有所有权限
3. **权限粒度**：权限按功能模块划分，便于管理和分配
4. **权限扩展性**：新增功能时按规范添加新权限
5. **前后端一致性**：前端权限配置需与后端保持一致

---

## 开发指南

### 新增权限步骤

1. **后端添加权限**
   ```go
   // 在 handler 中添加权限检查
   r.POST("/api/new-endpoint",
       auth.PermissionMiddleware(db, "new_permission"),
       handler)
   ```

2. **更新权限配置**
   - 在本文档中添加权限定义
   - 更新权限统计信息

3. **前端配置权限**
   ```javascript
   // 在路由配置中添加权限
   {
     path: '/new-feature',
     meta: { permissions: ['new_permission'] }
   }
   ```

4. **分配权限给角色**
   - 在角色管理页面为相应角色分配新权限
   - 测试权限是否生效

---

**文档版本**: 2.0
**最后更新**: 2025-02-05
**维护者**: System Development Team
