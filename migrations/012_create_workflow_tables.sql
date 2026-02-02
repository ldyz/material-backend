-- 工作流定义表
CREATE TABLE IF NOT EXISTS workflow_definitions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    module VARCHAR(50) NOT NULL, -- 'inbound', 'requisition' 等
    version INT DEFAULT 1,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 工作流节点表
CREATE TABLE IF NOT EXISTS workflow_nodes (
    id SERIAL PRIMARY KEY,
    workflow_id INT NOT NULL REFERENCES workflow_definitions(id) ON DELETE CASCADE,
    node_key VARCHAR(50) NOT NULL, -- 节点标识，如 'start', 'approval_1', 'end'
    node_type VARCHAR(20) NOT NULL, -- 'start', 'approval', 'end', 'parallel', 'merge'
    node_name VARCHAR(100) NOT NULL,
    description TEXT,
    approval_type VARCHAR(20), -- 'sequential' (顺序), 'parallel' (并行), 'any' (任一)
    timeout_hours INT, -- 超时时间（小时）
    auto_approve BOOLEAN DEFAULT false, -- 是否自动通过
    is_required BOOLEAN DEFAULT true, -- 是否必需节点
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(workflow_id, node_key)
);

-- 工作流边（连接线）表
CREATE TABLE IF NOT EXISTS workflow_edges (
    id SERIAL PRIMARY KEY,
    workflow_id INT NOT NULL REFERENCES workflow_definitions(id) ON DELETE CASCADE,
    from_node VARCHAR(50) NOT NULL, -- 源节点 node_key
    to_node VARCHAR(50) NOT NULL, -- 目标节点 node_key
    condition_expression TEXT, -- 条件表达式（可选）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(workflow_id, from_node, to_node)
);

-- 工作流节点审批人配置表
CREATE TABLE IF NOT EXISTS workflow_node_approvers (
    id SERIAL PRIMARY KEY,
    node_id INT NOT NULL REFERENCES workflow_nodes(id) ON DELETE CASCADE,
    approver_type VARCHAR(20) NOT NULL, -- 'user', 'role', 'department', 'superior'
    approver_id INT, -- 审批人ID（user）或 角色ID（role）
    approver_name VARCHAR(100), -- 审批人或角色名称
    sequence INT DEFAULT 0, -- 审批顺序（用于同一节点多人审批）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 工作流实例表
CREATE TABLE IF NOT EXISTS workflow_instances (
    id SERIAL PRIMARY KEY,
    workflow_id INT NOT NULL REFERENCES workflow_definitions(id),
    business_type VARCHAR(50) NOT NULL, -- 'inbound_order', 'requisition' 等
    business_id INT NOT NULL, -- 业务单据ID
    business_no VARCHAR(50), -- 业务单据编号
    current_node VARCHAR(50), -- 当前节点
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected', 'cancelled'
    initiator_id INT NOT NULL, -- 发起人ID
    initiator_name VARCHAR(100) NOT NULL, -- 发起人姓名
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 工作流审批记录表
CREATE TABLE IF NOT EXISTS workflow_approvals (
    id SERIAL PRIMARY KEY,
    instance_id INT NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    node_id INT NOT NULL REFERENCES workflow_nodes(id),
    node_key VARCHAR(50) NOT NULL,
    approver_id INT NOT NULL, -- 审批人ID
    approver_name VARCHAR(100) NOT NULL, -- 审批人姓名
    action VARCHAR(20) NOT NULL, -- 'approve', 'reject', 'return', 'comment'
    remark TEXT, -- 审批意见
    attachments TEXT, -- 附件信息（JSON格式）
    approved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 工作流待办表
CREATE TABLE IF NOT EXISTS workflow_pending_tasks (
    id SERIAL PRIMARY KEY,
    instance_id INT NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    node_id INT NOT NULL REFERENCES workflow_nodes(id),
    node_key VARCHAR(50) NOT NULL,
    node_name VARCHAR(100) NOT NULL,
    business_type VARCHAR(50) NOT NULL,
    business_id INT NOT NULL,
    business_no VARCHAR(50),
    approver_id INT NOT NULL,
    approver_name VARCHAR(100) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'approved', 'rejected', 'returned', 'cancelled'
    is_parallel BOOLEAN DEFAULT false, -- 是否为并行审批
    arrived_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 工作流操作日志表
CREATE TABLE IF NOT EXISTS workflow_logs (
    id SERIAL PRIMARY KEY,
    instance_id INT NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    node_key VARCHAR(50),
    action VARCHAR(50) NOT NULL, -- 'start', 'approve', 'reject', 'return', 'cancel', 'comment'
    actor_id INT NOT NULL,
    actor_name VARCHAR(100) NOT NULL,
    action_data TEXT, -- 操作数据（JSON格式）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_workflow_definitions_module ON workflow_definitions(module);
CREATE INDEX IF NOT EXISTS idx_workflow_definitions_active ON workflow_definitions(is_active);
CREATE INDEX IF NOT EXISTS idx_workflow_nodes_workflow ON workflow_nodes(workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_edges_workflow ON workflow_edges(workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_node_approvers_node ON workflow_node_approvers(node_id);
CREATE INDEX IF NOT EXISTS idx_workflow_instances_business ON workflow_instances(business_type, business_id);
CREATE INDEX IF NOT EXISTS idx_workflow_instances_status ON workflow_instances(status);
CREATE INDEX IF NOT EXISTS idx_workflow_instances_current ON workflow_instances(current_node);
CREATE INDEX IF NOT EXISTS idx_workflow_approvals_instance ON workflow_approvals(instance_id);
CREATE INDEX IF NOT EXISTS idx_workflow_approvals_node ON workflow_approvals(node_id);
CREATE INDEX IF NOT EXISTS idx_workflow_pending_tasks_approver ON workflow_pending_tasks(approver_id, status);
CREATE INDEX IF NOT EXISTS idx_workflow_pending_tasks_business ON workflow_pending_tasks(business_type, business_id);
CREATE INDEX IF NOT EXISTS idx_workflow_logs_instance ON workflow_logs(instance_id);

-- 创建默认的简单审批流程（入库单）
INSERT INTO workflow_definitions (name, description, module, version, is_active)
VALUES ('入库单默认审批流程', '入库单的两级审批流程', 'inbound', 1, true)
ON CONFLICT (name) DO NOTHING;

-- 为入库单创建默认节点
INSERT INTO workflow_nodes (workflow_id, node_key, node_type, node_name, description, approval_type, is_required)
SELECT
    id,
    'start',
    'start',
    '开始',
    '流程开始节点',
    'sequential',
    true
FROM workflow_definitions
WHERE module = 'inbound'
ON CONFLICT (workflow_id, node_key) DO NOTHING;

INSERT INTO workflow_nodes (workflow_id, node_key, node_type, node_name, description, approval_type, is_required)
SELECT
    id,
    'approval_1',
    'approval',
    '一级审批',
    '部门主管审批',
    'sequential',
    true
FROM workflow_definitions
WHERE module = 'inbound'
ON CONFLICT (workflow_id, node_key) DO NOTHING;

INSERT INTO workflow_nodes (workflow_id, node_key, node_type, node_name, description, approval_type, is_required)
SELECT
    id,
    'approval_2',
    'approval',
    '二级审批',
    '仓库管理员审批',
    'sequential',
    true
FROM workflow_definitions
WHERE module = 'inbound'
ON CONFLICT (workflow_id, node_key) DO NOTHING;

INSERT INTO workflow_nodes (workflow_id, node_key, node_type, node_name, description, approval_type, is_required)
SELECT
    id,
    'end',
    'end',
    '结束',
    '流程结束节点',
    'sequential',
    true
FROM workflow_definitions
WHERE module = 'inbound'
ON CONFLICT (workflow_id, node_key) DO NOTHING;

-- 创建连接线
INSERT INTO workflow_edges (workflow_id, from_node, to_node)
SELECT id, 'start', 'approval_1' FROM workflow_definitions WHERE module = 'inbound'
ON CONFLICT (workflow_id, from_node, to_node) DO NOTHING;

INSERT INTO workflow_edges (workflow_id, from_node, to_node)
SELECT id, 'approval_1', 'approval_2' FROM workflow_definitions WHERE module = 'inbound'
ON CONFLICT (workflow_id, from_node, to_node) DO NOTHING;

INSERT INTO workflow_edges (workflow_id, from_node, to_node)
SELECT id, 'approval_2', 'end' FROM workflow_definitions WHERE module = 'inbound'
ON CONFLICT (workflow_id, from_node, to_node) DO NOTHING;

-- 创建默认的简单审批流程（领料单）
INSERT INTO workflow_definitions (name, description, module, version, is_active)
VALUES ('领料单默认审批流程', '领料单的两级审批流程', 'requisition', 1, true)
ON CONFLICT (name) DO NOTHING;

-- 为领料单创建默认节点
INSERT INTO workflow_nodes (workflow_id, node_key, node_type, node_name, description, approval_type, is_required)
SELECT
    id,
    'start',
    'start',
    '开始',
    '流程开始节点',
    'sequential',
    true
FROM workflow_definitions
WHERE module = 'requisition'
ON CONFLICT (workflow_id, node_key) DO NOTHING;

INSERT INTO workflow_nodes (workflow_id, node_key, node_type, node_name, description, approval_type, is_required)
SELECT
    id,
    'approval_1',
    'approval',
    '一级审批',
    '部门主管审批',
    'sequential',
    true
FROM workflow_definitions
WHERE module = 'requisition'
ON CONFLICT (workflow_id, node_key) DO NOTHING;

INSERT INTO workflow_nodes (workflow_id, node_key, node_type, node_name, description, approval_type, is_required)
SELECT
    id,
    'approval_2',
    'approval',
    '二级审批',
    '仓库管理员审批',
    'sequential',
    true
FROM workflow_definitions
WHERE module = 'requisition'
ON CONFLICT (workflow_id, node_key) DO NOTHING;

INSERT INTO workflow_nodes (workflow_id, node_key, node_type, node_name, description, approval_type, is_required)
SELECT
    id,
    'end',
    'end',
    '结束',
    '流程结束节点',
    'sequential',
    true
FROM workflow_definitions
WHERE module = 'requisition'
ON CONFLICT (workflow_id, node_key) DO NOTHING;

-- 创建连接线
INSERT INTO workflow_edges (workflow_id, from_node, to_node)
SELECT id, 'start', 'approval_1' FROM workflow_definitions WHERE module = 'requisition'
ON CONFLICT (workflow_id, from_node, to_node) DO NOTHING;

INSERT INTO workflow_edges (workflow_id, from_node, to_node)
SELECT id, 'approval_1', 'approval_2' FROM workflow_definitions WHERE module = 'requisition'
ON CONFLICT (workflow_id, from_node, to_node) DO NOTHING;

INSERT INTO workflow_edges (workflow_id, from_node, to_node)
SELECT id, 'approval_2', 'end' FROM workflow_definitions WHERE module = 'requisition'
ON CONFLICT (workflow_id, from_node, to_node) DO NOTHING;
