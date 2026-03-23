# 角色权限配置说明

## 系统角色及权限配置

本文档说明系统中各角色的权限配置。

---

## 1. 管理员 (admin)

**权限标识**: `admin`

**权限**: 拥有所有权限，可以访问系统所有功能

**主要职责**:
- 系统管理
- 用户管理
- 角色和权限管理
- 工作流配置
- 数据备份和恢复
- 系统配置
- 审计日志查看

---

## 2. 保管员

**权限列表**:
```
material_view          - 查看物资
stock_view             - 查看库存
stock_in               - 入库操作
stock_out              - 出库操作
stock_edit             - 编辑库存
stock_export           - 导出库存
inbound_view           - 查看入库单
inbound_approve        - 审批入库单
requisition_view       - 查看领用申请
requisition_approve    - 审批领用申请
requisition_issue      - 发放领用申请
```

**主要职责**:
- 物资入库和出库操作
- 库存管理
- 审批入库单
- 审批并发放领用申请
- 查看物资信息

---

## 3. 施工员

**权限列表**:
```
material_view              - 查看物资
stock_view                 - 查看库存
stocklog_view              - 查看库存日志
requisition_view           - 查看领用申请
requisition_create         - 创建领用申请
requisition_edit           - 编辑领用申请
inbound_view               - 查看入库单
inbound_create             - 创建入库单
project_view               - 查看项目
construction_log_view      - 查看施工日志
construction_log_create    - 创建施工日志
construction_log_edit      - 编辑施工日志
progress_view              - 查看进度
progress_create            - 创建进度
progress_edit              - 编辑进度
```

**主要职责**:
- 创建和管理领用申请
- 创建入库单
- 填写施工日志
- 更新项目进度
- 查看物资和库存信息

---

## 4. 分包材料员

**权限列表**:
```
material_view      - 查看物资
stock_view         - 查看库存
requisition_view   - 查看领用申请
requisition_create - 创建领用申请
```

**主要职责**:
- 查看物资和库存信息
- 创建领用申请（由项目经理或保管员审批）

---

## 5. 项目经理

**权限列表**:
```
project_view              - 查看项目
project_edit              - 编辑项目
material_view             - 查看物资
material_create           - 创建物资
material_edit             - 编辑物资
stock_view                - 查看库存
stock_in                  - 入库操作
stock_out                 - 出库操作
stock_edit                - 编辑库存
stocklog_view             - 查看库存日志
stock_export              - 导出库存
requisition_view          - 查看领用申请
requisition_create        - 创建领用申请
requisition_edit          - 编辑领用申请
requisition_delete        - 删除领用申请
requisition_approve       - 审批领用申请
inbound_view              - 查看入库单
inbound_create            - 创建入库单
inbound_edit              - 编辑入库单
inbound_approve           - 审批入库单
construction_log_view     - 查看施工日志
construction_log_create   - 创建施工日志
construction_log_edit     - 编辑施工日志
progress_view             - 查看进度
progress_create           - 创建进度
progress_edit             - 编辑进度
progress_export           - 导出进度
material_plan_view        - 查看物资计划
material_plan_create      - 创建物资计划
material_plan_edit        - 编辑物资计划
material_plan_approve     - 审批物资计划
workflow_view             - 查看工作流
workflow_task_view        - 查看工作流任务
workflow_task_approve     - 审批工作流任务
workflow_instance_view    - 查看工作流实例
ai_agent_view             - 查看AI能力
ai_agent_query            - AI查询
ai_agent_logs             - 查看AI日志
audit_view                - 查看审计日志
```

**主要职责**:
- 项目全面管理
- 物资计划制定和审批
- 领用申请和入库单审批
- 施工日志管理
- 进度管理
- 工作流审批
- 使用AI辅助功能
- 查看审计日志

---

## 6. 材料员

**权限列表**:
```
material_view          - 查看物资
material_create        - 创建物资
material_edit          - 编辑物资
stock_view             - 查看库存
stock_in               - 入库操作
stock_out              - 出库操作
stock_edit             - 编辑库存
stock_export           - 导出库存
stocklog_view          - 查看库存日志
stock_alerts           - 库存预警
requisition_view       - 查看领用申请
requisition_create     - 创建领用申请
requisition_edit       - 编辑领用申请
requisition_delete     - 删除领用申请
inbound_view           - 查看入库单
inbound_create         - 创建入库单
inbound_edit           - 编辑入库单
inbound_approve        - 审批入库单
material_plan_view     - 查看物资计划
material_plan_create   - 创建物资计划
material_plan_edit     - 编辑物资计划
progress_view          - 查看进度
```

**主要职责**:
- 物资信息维护
- 物资计划制定
- 入库单管理
- 领用申请管理
- 库存操作
- 库存预警管理

---

## 权限命名规范

系统采用 `{模块}_{操作}` 的命名规范：

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

---

## 权限配置SQL脚本

如需重新配置角色权限，可执行以下脚本：

```bash
psql -h 127.0.0.1 -U materials -d materials -f scripts/setup_role_permissions.sql
```

---

**文档版本**: 1.0
**最后更新**: 2026-02-07
