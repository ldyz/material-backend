# 前端权限配置说明

## 更新时间
2026-02-01

## 前后端权限对照表

### 1. 物资计划模块

| 前端功能 | 后端权限 | 说明 |
|---------|---------|------|
| 查看物资计划列表 | `material_plan_view` | 基础查看权限 |
| 创建物资计划 | `material_plan_create` | 创建新计划 |
| 编辑物资计划 | `material_plan_edit` | 修改计划内容 |
| 删除物资计划 | `material_plan_delete` | 删除计划 |
| 提交审核 | `material_plan_edit` | 编辑权限即可提交 |
| 审核计划 | `material_plan_approve` | 审核通过/拒绝 |
| 添加计划项 | `material_plan_edit` | 编辑权限即可添加 |

**路由配置：**
- 路径：`/material-plans`
- 组件：`MaterialPlans.vue`
- 所需权限：`['material_plan_view']`

**菜单配置：**
```javascript
{
  path: '/material-plans',
  title: '物资计划',
  permissions: ['material_plan_view']
}
```

---

### 2. 物资管理模块

| 前端功能 | 后端权限 | 说明 |
|---------|---------|------|
| 查看物资列表 | `material_view` | 基础查看权限 |
| 创建物资 | `material_create` | 添加新物资 |
| 编辑物资 | `material_edit` | 修改物资信息 |
| 删除物资 | `material_delete` | 删除物资 |
| 导入物资 | `material_import` | 批量导入 |
| 导出物资 | `material_export` | 导出数据 |
| 物资入库 | `material_in` | 入库操作 |

**路由配置：**
- 路径：`/materials`
- 组件：`Materials.vue`
- 所需权限：`['material_view', 'material_create', 'material_edit', 'material_delete', 'material_import', 'material_export', 'material_in']`

**菜单配置：**
```javascript
{
  path: '/materials',
  title: '物资管理',
  permissions: ['material_view']
}
```

---

### 3. 库存管理模块

| 前端功能 | 后端权限 | 说明 |
|---------|---------|------|
| 查看库存 | `stock_view` | 查看库存列表 |
| 库存入库 | `stock_in` | 入库操作 |
| 库存出库 | `stock_out` | 出库操作 |
| 编辑库存 | `stock_edit` | 修改库存信息 |
| 删除库存 | `stock_delete` | 删除库存记录 |
| 查看库存日志 | `stocklog_view` | 查看操作日志 |
| 删除库存日志 | `stocklog_delete` | 删除日志 |
| 导出库存 | `stock_export` | 导出数据 |

**路由配置：**
- 路径：`/stock`
- 组件：`Stock.vue`
- 所需权限：`['stock_view', 'stock_in', 'stock_out', 'stock_edit', 'stock_delete', 'stocklog_view', 'stocklog_delete', 'stock_export']`

**菜单配置：**
```javascript
{
  path: '/stock',
  title: '库存管理',
  permissions: ['stock_view']
}
```

---

### 4. 出库管理模块

| 前端功能 | 后端权限 | 说明 |
|---------|---------|------|
| 查看出库单 | `requisition_view` | 查看出库单列表 |
| 创建出库单 | `requisition_create` | 创建新出库单 |
| 编辑出库单 | `requisition_edit` | 修改出库单 |
| 删除出库单 | `requisition_delete` | 删除出库单 |
| 审核出库单 | `requisition_approve` | 审核通过/拒绝 |
| 发货 | `requisition_issue` | 标记为已发货 |
| 导出出库单 | `requisition_export` | 导出数据 |

**路由配置：**
- 路径：`/requisitions`
- 组件：`Requisitions.vue`
- 所需权限：`['requisition_view']`

**菜单配置：**
```javascript
{
  path: '/requisitions',
  title: '出库管理',
  permissions: ['requisition_view']
}
```

---

### 5. 入库管理模块

| 前端功能 | 后端权限 | 说明 |
|---------|---------|------|
| 查看入库单 | `inbound_view` | 查看入库单列表 |
| 创建入库单 | `inbound_create` | 创建新入库单 |
| 编辑入库单 | `inbound_edit` | 修改入库单 |
| 删除入库单 | `inbound_delete` | 删除入库单 |
| 审核入库单 | `inbound_approve` | 审核通过/拒绝 |
| 导出入库单 | `inbound_export` | 导出数据 |

**路由配置：**
- 路径：`/inbound`
- 组件：`Inbound.vue`
- 所需权限：`['inbound_view']`

**菜单配置：**
```javascript
{
  path: '/inbound',
  title: '入库管理',
  permissions: ['inbound_view']
}
```

---

### 6. 项目管理模块

| 前端功能 | 后端权限 | 说明 |
|---------|---------|------|
| 查看项目 | `project_view` | 查看项目列表 |
| 创建项目 | `project_create` | 创建新项目 |
| 编辑项目 | `project_edit` | 修改项目信息 |
| 删除项目 | `project_delete` | 删除项目 |
| 管理项目成员 | `project_member_manage` | 管理成员权限 |

---

### 7. 进度管理模块

| 前端功能 | 后端权限 | 说明 |
|---------|---------|------|
| 查看进度 | `progress_view` | 查看进度计划 |
| 创建进度 | `progress_create` | 创建新进度 |
| 编辑进度 | `progress_edit` | 修改进度 |
| 删除进度 | `progress_delete` | 删除进度 |

---

### 8. 施工日志模块

| 前端功能 | 后端权限 | 说明 |
|---------|---------|------|
| 查看日志 | `construction_log_view` | 查看施工日志 |
| 添加日志 | `constructionlog_add` | 添加新日志 |
| 编辑日志 | `construction_log_edit` | 修改日志内容 |
| 删除日志 | `construction_log_delete` | 删除日志 |

---

### 9. 系统管理模块

| 前端功能 | 后端权限 | 说明 |
|---------|---------|------|
| 查看日志 | `system_log` | 查看系统日志 |
| 数据备份 | `system_backup` | 备份数据 |
| 系统配置 | `system_config` | 系统配置 |
| 数据报告 | `system_report` | 生成报告 |
| 统计分析 | `system_statistics` | 查看统计 |

---

## 权限检查优先级

1. **管理员（admin）**：拥有所有权限，自动通过所有检查
2. **普通用户**：根据角色拥有的权限进行判断
3. **未登录用户**：只能访问公开页面（登录页）

## 权限检查层次

```
前端三层权限检查：
┌─────────────────────────────────────┐
│  1. 菜单显示权限                      │
│     - 控制 Sidebar 菜单显示           │
│     - 根据权限过滤菜单项              │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│  2. 路由导航权限                      │
│     - 路由守卫检查                   │
│     - 无权限重定向到 Dashboard        │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│  3. 后端API权限                       │
│     - 所有API调用都检查权限          │
│     - 无权限返回403错误              │
└─────────────────────────────────────┘
```

## 注意事项

1. **前端权限只是UI优化**，真正的安全控制在后端
2. **菜单权限过滤**在 MainLayout 组件初始化时执行一次
3. **路由权限检查**在路由守卫中每次导航时执行
4. **后端API权限**在每个请求时验证
5. **权限变更**后需要刷新页面才能看到菜单变化

## 文件清单

| 文件 | 作用 |
|------|------|
| `@/utils/permissions.js` | 权限检查工具函数 |
| `@/composables/usePermission.js` | 权限检查组合式函数 |
| `@/components/layout/MainLayout.vue` | 菜单权限过滤 |
| `@/router/index.js` | 路由权限守卫 |

## 开发建议

1. **新增功能时**：
   - 后端先定义权限
   - 前端在路由配置中添加权限
   - 菜单配置同步更新

2. **测试权限时**：
   - 测试不同角色的用户
   - 验证菜单显示正确
   - 验证路由跳转正确
   - 验证API调用权限正确

3. **调试问题时**：
   - 打开控制台查看权限警告
   - 检查用户角色和权限
   - 验证后端API权限配置
