-- 角色权限配置脚本
-- 为系统中的角色配置正确的权限

-- 1. 更新管理员角色 (admin) - 拥有所有权限
UPDATE roles SET permissions = 'admin' WHERE name = 'admin';

-- 2. 更新保管员角色 (保管员)
-- 保管员负责物资的入库、出库、库存管理，以及领用申请的审批和发放
UPDATE roles SET permissions = 'material_view,stock_view,stock_in,stock_out,stock_edit,stock_export,inbound_view,inbound_approve,requisition_view,requisition_approve,requisition_issue' WHERE name = '保管员';

-- 3. 更新施工员角色 (施工员)
-- 施工员负责现场物资管理、领用申请、施工日志等
UPDATE roles SET permissions = 'material_view,stock_view,stocklog_view,requisition_view,requisition_create,requisition_edit,inbound_view,inbound_create,project_view,construction_log_view,construction_log_create,construction_log_edit,progress_view,progress_create,progress_edit' WHERE name = '施工员';

-- 4. 更新分包材料员角色 (分包材料员)
-- 分包材料员只能查看和创建领用申请
UPDATE roles SET permissions = 'material_view,stock_view,requisition_view,requisition_create' WHERE name = '分包材料员';

-- 5. 更新项目经理角色 (项目经理)
-- 项目经理拥有项目内的大部分权限，包括审批、查看报告等
UPDATE roles SET permissions = 'project_view,project_edit,material_view,material_create,material_edit,stock_view,stock_in,stock_out,stock_edit,stocklog_view,stock_export,requisition_view,requisition_create,requisition_edit,requisition_delete,requisition_approve,inbound_view,inbound_create,inbound_edit,inbound_approve,construction_log_view,construction_log_create,construction_log_edit,progress_view,progress_create,progress_edit,progress_export,material_plan_view,material_plan_create,material_plan_edit,material_plan_approve,workflow_view,workflow_task_view,workflow_task_approve,workflow_instance_view,ai_agent_view,ai_agent_query,ai_agent_logs,audit_view' WHERE name = '项目经理';

-- 6. 更新材料员角色 (材料员)
-- 材料员负责物资计划和库存管理
UPDATE roles SET permissions = 'material_view,material_create,material_edit,stock_view,stock_in,stock_out,stock_edit,stock_export,stocklog_view,stock_alerts,requisition_view,requisition_create,requisition_edit,requisition_delete,inbound_view,inbound_create,inbound_edit,inbound_approve,material_plan_view,material_plan_create,material_plan_edit,material_plan_view,progress_view' WHERE name = '材料员';

-- 查看更新结果
SELECT id, name, description,
       array_length(string_to_array(permissions, ','), 1) as permission_count,
       permissions
FROM roles
ORDER BY id;
