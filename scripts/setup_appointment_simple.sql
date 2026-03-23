-- ====================================================================
-- 施工预约功能 - 角色、权限和工作流配置（适配现有数据库结构）
-- ====================================================================

-- 1. 创建预约相关角色
INSERT INTO roles (id, name, description, permissions, created_at) VALUES
(10, '预约管理员', '负责创建和管理施工预约单，包括客户等预约申请人员',
  'appointment_view,appointment_create,appointment_edit,appointment_delete,appointment_submit,project_view',
  CURRENT_TIMESTAMP),
(11, '预约施工员', '负责第一级审批，确认可以承接预约作业',
  'appointment_view,appointment_approve',
  CURRENT_TIMESTAMP),
(12, '作业人员', '负责执行具体施工作业的人员',
  'appointment_view,appointment_execute',
  CURRENT_TIMESTAMP)
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  permissions = EXCLUDED.permissions;

-- 2. 为现有角色添加预约相关权限
-- 为项目经理(ID=4)添加审批权限
UPDATE roles
SET permissions = permissions || ',appointment_view,appointment_approve,appointment_assign'
WHERE id = 4 AND permissions NOT LIKE '%appointment_approve%';

-- 为经理(ID=2)添加审批权限
UPDATE roles
SET permissions = permissions || ',appointment_view,appointment_approve'
WHERE id = 2 AND permissions NOT LIKE '%appointment_approve%';

-- 3. 查看创建的角色
SELECT '=== 创建的角色 ===' AS info;
SELECT id, name, description FROM roles WHERE id IN (10, 11, 12) ORDER BY id;

-- 4. 创建工作流定义

-- 普通预约工作流
INSERT INTO workflow_definitions (id, name, description, module, is_active, created_at, updated_at) VALUES
(10, '普通施工预约审批流程', '普通施工预约的两级审批流程：施工员 -> 项目经理', 'appointment', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (name) DO UPDATE SET
  is_active = true,
  updated_at = CURRENT_TIMESTAMP;

-- 加急预约工作流
INSERT INTO workflow_definitions (id, name, description, module, is_active, created_at, updated_at) VALUES
(11, '加急施工预约审批流程', '加急施工预约的三级审批流程：施工员 -> 经理 -> 项目经理', 'appointment', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (name) DO UPDATE SET
  is_active = true,
  updated_at = CURRENT_TIMESTAMP;

SELECT '=== 工作流定义 ===' AS info;
SELECT id, name, description, module, is_active FROM workflow_definitions WHERE id IN (10, 11) ORDER BY id;

-- 5. 创建工作流节点和边
DO $$
DECLARE
  wf_id INTEGER;
  start_node_id INTEGER;
  foreman_node_id INTEGER;
  manager_node_id INTEGER;
  pm_node_id INTEGER;
  end_node_id INTEGER;
BEGIN
  -- ============ 普通预约工作流 ============

  -- 获取普通工作流ID
  SELECT id INTO wf_id FROM workflow_definitions WHERE name = '普通施工预约审批流程';

  -- 删除旧的节点和边（如果存在）
  DELETE FROM workflow_node_approvers WHERE node_id IN (
    SELECT id FROM workflow_nodes WHERE workflow_id = wf_id
  );
  DELETE FROM workflow_edges WHERE workflow_id = wf_id;
  DELETE FROM workflow_nodes WHERE workflow_id = wf_id;

  -- 创建开始节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, description, created_at)
  VALUES (wf_id, 'start', '开始', 'start', '预约申请开始', CURRENT_TIMESTAMP)
  RETURNING id INTO start_node_id;

  -- 创建施工员审批节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, approval_type, description, created_at)
  VALUES (wf_id, 'foreman_approve', '施工员审批', 'approval', 'any', '施工员确认可以承接作业', CURRENT_TIMESTAMP)
  RETURNING id INTO foreman_node_id;

  -- 创建项目经理审批节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, approval_type, description, created_at)
  VALUES (wf_id, 'pm_approve', '项目经理审批', 'approval', 'any', '项目经理最终确认', CURRENT_TIMESTAMP)
  RETURNING id INTO pm_node_id;

  -- 创建结束节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, description, created_at)
  VALUES (wf_id, 'end', '结束', 'end', '流程结束', CURRENT_TIMESTAMP)
  RETURNING id INTO end_node_id;

  -- 创建工作流边（连接线）
  INSERT INTO workflow_edges (workflow_id, from_node, to_node, created_at) VALUES
    (wf_id, 'start', 'foreman_approve', CURRENT_TIMESTAMP),
    (wf_id, 'foreman_approve', 'pm_approve', CURRENT_TIMESTAMP),
    (wf_id, 'pm_approve', 'end', CURRENT_TIMESTAMP);

  -- 为节点添加审批人
  -- 施工员角色ID=11
  INSERT INTO workflow_node_approvers (node_id, approver_type, approver_value, approver_name, created_at)
  VALUES (foreman_node_id, 'role', '11', '施工员角色', CURRENT_TIMESTAMP);

  -- 项目经理角色ID=4
  INSERT INTO workflow_node_approvers (node_id, approver_type, approver_value, approver_name, created_at)
  VALUES (pm_node_id, 'role', '4', '项目经理角色', CURRENT_TIMESTAMP);

  RAISE NOTICE '普通预约工作流创建完成，工作流ID: %', wf_id;

END $$;

DO $$
DECLARE
  wf_id INTEGER;
  start_node_id INTEGER;
  foreman_node_id INTEGER;
  manager_node_id INTEGER;
  pm_node_id INTEGER;
  end_node_id INTEGER;
BEGIN
  -- ============ 加急预约工作流 ============

  -- 获取加急工作流ID
  SELECT id INTO wf_id FROM workflow_definitions WHERE name = '加急施工预约审批流程';

  -- 删除旧的节点和边（如果存在）
  DELETE FROM workflow_node_approvers WHERE node_id IN (
    SELECT id FROM workflow_nodes WHERE workflow_id = wf_id
  );
  DELETE FROM workflow_edges WHERE workflow_id = wf_id;
  DELETE FROM workflow_nodes WHERE workflow_id = wf_id;

  -- 创建开始节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, description, created_at)
  VALUES (wf_id, 'start', '开始', 'start', '加急预约申请开始', CURRENT_TIMESTAMP)
  RETURNING id INTO start_node_id;

  -- 创建施工员审批节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, approval_type, description, created_at)
  VALUES (wf_id, 'foreman_approve', '施工员审批', 'approval', 'any', '施工员确认可以承接作业', CURRENT_TIMESTAMP)
  RETURNING id INTO foreman_node_id;

  -- 创建经理审批节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, approval_type, description, created_at)
  VALUES (wf_id, 'manager_approve', '经理审批', 'approval', 'any', '加急作业需要经理批准', CURRENT_TIMESTAMP)
  RETURNING id INTO manager_node_id;

  -- 创建项目经理审批节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, approval_type, description, created_at)
  VALUES (wf_id, 'pm_approve', '项目经理审批', 'approval', 'any', '项目经理最终确认', CURRENT_TIMESTAMP)
  RETURNING id INTO pm_node_id;

  -- 创建结束节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, description, created_at)
  VALUES (wf_id, 'end', '结束', 'end', '流程结束', CURRENT_TIMESTAMP)
  RETURNING id INTO end_node_id;

  -- 创建工作流边（连接线）
  INSERT INTO workflow_edges (workflow_id, from_node, to_node, created_at) VALUES
    (wf_id, 'start', 'foreman_approve', CURRENT_TIMESTAMP),
    (wf_id, 'foreman_approve', 'manager_approve', CURRENT_TIMESTAMP),
    (wf_id, 'manager_approve', 'pm_approve', CURRENT_TIMESTAMP),
    (wf_id, 'pm_approve', 'end', CURRENT_TIMESTAMP);

  -- 为节点添加审批人
  INSERT INTO workflow_node_approvers (node_id, approver_type, approver_value, approver_name, created_at) VALUES
    (foreman_node_id, 'role', '11', '施工员角色', CURRENT_TIMESTAMP),
    (manager_node_id, 'role', '2', '经理角色', CURRENT_TIMESTAMP),
    (pm_node_id, 'role', '4', '项目经理角色', CURRENT_TIMESTAMP);

  RAISE NOTICE '加急预约工作流创建完成，工作流ID: %', wf_id;

END $$;

-- 6. 查看创建的工作流节点
SELECT '=== 工作流节点 ===' AS info;
SELECT
  wd.name as workflow_name,
  wn.node_key,
  wn.node_name,
  wn.node_type,
  wn.approval_type
FROM workflow_nodes wn
JOIN workflow_definitions wd ON wn.workflow_id = wd.id
WHERE wd.id IN (10, 11)
ORDER BY wd.id, wn.id;

-- 7. 查看审批人配置
SELECT '=== 审批人配置 ===' AS info;
SELECT
  wn.node_name,
  wna.approver_type,
  wna.approver_value,
  wna.approver_name
FROM workflow_node_approvers wna
JOIN workflow_nodes wn ON wna.node_id = wn.id
JOIN workflow_definitions wd ON wn.workflow_id = wd.id
WHERE wd.id IN (10, 11)
ORDER BY wd.id, wn.id;

SELECT '=== 配置完成！ ===' AS status;

-- 8. 显示下一步操作说明
SELECT '请使用以下命令为用户分配预约角色：
UPDATE roles SET permissions = permissions || '',appointment_view,appointment_create'' WHERE id = 用户角色ID;

或者手动编辑角色的permissions字段，添加以下权限：
- appointment_view: 查看预约单
- appointment_create: 创建预约单
- appointment_edit: 编辑预约单
- appointment_delete: 删除预约单
- appointment_submit: 提交审批
- appointment_approve: 审批预约单
- appointment_assign: 分配作业人员
- appointment_execute: 执行作业
- appointment_cancel: 取消预约单' AS next_steps;
