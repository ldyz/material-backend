-- ====================================================================
-- 施工预约功能 - 角色、权限和工作流配置
-- ====================================================================

-- 1. 插入角色
INSERT INTO roles (id, name, description, created_at, updated_at) VALUES
(10, '预约管理员', '负责创建和管理施工预约单，包括客户等预约申请人员', NOW(), NOW()),
(11, '施工员', '负责第一级审批，确认可以承接预约作业', NOW(), NOW()),
(12, '作业人员', '负责执行具体施工作业的人员', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  updated_at = NOW();

-- 2. 插入权限
INSERT INTO permissions (id, name, description, module, created_at, updated_at) VALUES
(50, 'appointment_view', '查看预约单', 'appointment', NOW(), NOW()),
(51, 'appointment_create', '创建预约单', 'appointment', NOW(), NOW()),
(52, 'appointment_edit', '编辑预约单', 'appointment', NOW(), NOW()),
(53, 'appointment_delete', '删除预约单', 'appointment', NOW(), NOW()),
(54, 'appointment_submit', '提交预约单审批', 'appointment', NOW(), NOW()),
(55, 'appointment_approve', '审批预约单', 'appointment', NOW(), NOW()),
(56, 'appointment_assign', '分配作业人员', 'appointment', NOW(), NOW()),
(57, 'appointment_execute', '执行作业（开始、完成）', 'appointment', NOW(), NOW()),
(58, 'appointment_cancel', '取消预约单', 'appointment', NOW(), NOW()),
(59, 'appointment_manage', '管理预约单日历', 'appointment', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  updated_at = NOW();

-- 3. 分配角色权限

-- 预约管理员权限（客户）
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT 10, id, NOW() FROM permissions WHERE name IN (
  'appointment_view',
  'appointment_create',
  'appointment_edit',
  'appointment_delete',
  'appointment_submit'
)
ON CONFLICT DO NOTHING;

-- 施工员权限（第一级审批）
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT 11, id, NOW() FROM permissions WHERE name IN (
  'appointment_view',
  'appointment_approve'
)
ON CONFLICT DO NOTHING;

-- 作业人员权限
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT 12, id, NOW() FROM permissions WHERE name IN (
  'appointment_view',
  'appointment_execute'
)
ON CONFLICT DO NOTHING;

-- 项目经理已有审批权限（通过 system_config）
-- 经理已有审批权限（通过 system_config）

-- 4. 创建默认工作流定义

-- 普通预约工作流
INSERT INTO workflow_definitions (id, name, description, module, business_type, is_active, created_at, updated_at) VALUES
(2, '普通施工预约审批流程', '普通施工预约的两级审批流程：施工员 -> 项目经理', 'appointment', 'normal', true, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
  is_active = true,
  updated_at = NOW();

-- 获取工作流ID（使用变量）
DO $$
DECLARE
  wf_id INTEGER;
  start_node_id INTEGER;
  foreman_node_id INTEGER;
  pm_node_id INTEGER;
  end_node_id INTEGER;
BEGIN
  -- 获取工作流ID
  SELECT id INTO wf_id FROM workflow_definitions WHERE id = 2;

  -- 创建开始节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, description, created_at, updated_at)
  VALUES (wf_id, 'start', '开始', 'start', '预约申请开始', NOW(), NOW())
  RETURNING id INTO start_node_id;

  -- 创建施工员审批节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, approval_type, description, created_at, updated_at)
  VALUES (wf_id, 'foreman_approve', '施工员审批', 'approval', 'any', '施工员确认可以承接作业', NOW(), NOW())
  RETURNING id INTO foreman_node_id;

  -- 创建项目经理审批节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, approval_type, description, created_at, updated_at)
  VALUES (wf_id, 'pm_approve', '项目经理审批', 'approval', 'any', '项目经理最终确认', NOW(), NOW())
  RETURNING id INTO pm_node_id;

  -- 创建结束节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, description, created_at, updated_at)
  VALUES (wf_id, 'end', '结束', 'end', '流程结束', NOW(), NOW())
  RETURNING id INTO end_node_id;

  -- 创建工作流边（连接线）
  -- 开始 -> 施工员
  INSERT INTO workflow_edges (workflow_id, from_node, to_node, condition_expression, created_at, updated_at)
  VALUES (wf_id, 'start', 'foreman_approve', '', NOW(), NOW());

  -- 施工员 -> 项目经理
  INSERT INTO workflow_edges (workflow_id, from_node, to_node, condition_expression, created_at, updated_at)
  VALUES (wf_id, 'foreman_approve', 'pm_approve', '', NOW(), NOW());

  -- 项目经理 -> 结束
  INSERT INTO workflow_edges (workflow_id, from_node, to_node, condition_expression, created_at, updated_at)
  VALUES (wf_id, 'pm_approve', 'end', '', NOW(), NOW());

  -- 为施工员节点添加审批人（具有施工员角色的用户）
  INSERT INTO workflow_node_approvers (node_id, approver_type, approver_value, approver_name, created_at)
  SELECT foreman_node_id, 'role', '11', '施工员角色', NOW();

  -- 为项目经理节点添加审批人（具有项目经理角色的用户）
  INSERT INTO workflow_node_approvers (node_id, approver_type, approver_value, approver_name, created_at)
  SELECT pm_node_id, 'role', '4', '项目经理角色', NOW();

  RAISE NOTICE '普通预约工作流创建完成，工作流ID: %', wf_id;
END $$;

-- 加急预约工作流
INSERT INTO workflow_definitions (id, name, description, module, business_type, is_active, created_at, updated_at) VALUES
(3, '加急施工预约审批流程', '加急施工预约的三级审批流程：施工员 -> 经理 -> 项目经理', 'appointment', 'urgent', true, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
  is_active = true,
  updated_at = NOW();

DO $$
DECLARE
  wf_id INTEGER;
  start_node_id INTEGER;
  foreman_node_id INTEGER;
  manager_node_id INTEGER;
  pm_node_id INTEGER;
  end_node_id INTEGER;
BEGIN
  -- 获取工作流ID
  SELECT id INTO wf_id FROM workflow_definitions WHERE id = 3;

  -- 创建开始节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, description, created_at, updated_at)
  VALUES (wf_id, 'start', '开始', 'start', '加急预约申请开始', NOW(), NOW())
  RETURNING id INTO start_node_id;

  -- 创建施工员审批节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, approval_type, description, created_at, updated_at)
  VALUES (wf_id, 'foreman_approve', '施工员审批', 'approval', 'any', '施工员确认可以承接作业', NOW(), NOW())
  RETURNING id INTO foreman_node_id;

  -- 创建经理审批节点（加急需要经理批准）
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, approval_type, description, created_at, updated_at)
  VALUES (wf_id, 'manager_approve', '经理审批', 'approval', 'any', '加急作业需要经理批准', NOW(), NOW())
  RETURNING id INTO manager_node_id;

  -- 创建项目经理审批节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, approval_type, description, created_at, updated_at)
  VALUES (wf_id, 'pm_approve', '项目经理审批', 'approval', 'any', '项目经理最终确认', NOW(), NOW())
  RETURNING id INTO pm_node_id;

  -- 创建结束节点
  INSERT INTO workflow_nodes (workflow_id, node_key, node_name, node_type, description, created_at, updated_at)
  VALUES (wf_id, 'end', '结束', 'end', '流程结束', NOW(), NOW())
  RETURNING id INTO end_node_id;

  -- 创建工作流边（连接线）
  INSERT INTO workflow_edges (workflow_id, from_node, to_node, condition_expression, created_at, updated_at)
  VALUES
    (wf_id, 'start', 'foreman_approve', '', NOW(), NOW()),
    (wf_id, 'foreman_approve', 'manager_approve', '', NOW(), NOW()),
    (wf_id, 'manager_approve', 'pm_approve', '', NOW(), NOW()),
    (wf_id, 'pm_approve', 'end', '', NOW(), NOW());

  -- 为节点添加审批人
  INSERT INTO workflow_node_approvers (node_id, approver_type, approver_value, approver_name, created_at)
  VALUES
    (foreman_node_id, 'role', '11', '施工员角色', NOW()),
    (manager_node_id, 'role', '2', '经理角色', NOW()),
    (pm_node_id, 'role', '4', '项目经理角色', NOW());

  RAISE NOTICE '加急预约工作流创建完成，工作流ID: %', wf_id;
END $$;

-- 输出配置信息
SELECT '角色、权限和工作流配置完成！' AS status;

-- 显示创建的角色
SELECT id, name, description FROM roles WHERE id IN (10, 11, 12) ORDER BY id;

-- 显示创建的权限
SELECT id, name, description FROM permissions WHERE id BETWEEN 50 AND 59 ORDER BY id;

-- 显示创建的工作流
SELECT id, name, description, module, business_type FROM workflow_definitions WHERE id IN (10, 11) ORDER BY id;
