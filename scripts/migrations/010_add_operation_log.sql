-- 010_add_operation_log.sql
-- 创建通用操作日志表，用于审计所有用户操作

CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    username VARCHAR(100),
    operation VARCHAR(100) NOT NULL,           -- 操作类型：create, update, delete, approve, reject, etc.
    module VARCHAR(50) NOT NULL,               -- 模块：material_plan, inbound, outbound, stock, etc.
    resource_type VARCHAR(100),                -- 资源类型：MaterialPlan, InboundOrder, Requisition, etc.
    resource_id INTEGER,                        -- 资源ID
    resource_no VARCHAR(100),                   -- 资源编号（如计划号、单号）
    changes JSONB,                             -- 变更内容（记录before/after）
    request_method VARCHAR(10),                -- 请求方法：GET, POST, PUT, DELETE
    request_path VARCHAR(500),                 -- 请求路径
    request_params JSONB,                      -- 请求参数
    ip_address VARCHAR(45),                    -- IP地址
    user_agent TEXT,                           -- 用户代理
    status VARCHAR(20) NOT NULL DEFAULT 'success', -- 状态：success, error
    error_message TEXT,                        -- 错误信息
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_op_logs_user ON operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_op_logs_operation ON operation_logs(operation);
CREATE INDEX IF NOT EXISTS idx_op_logs_module ON operation_logs(module);
CREATE INDEX IF NOT EXISTS idx_op_logs_resource ON operation_logs(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_op_logs_resource_no ON operation_logs(resource_no);
CREATE INDEX IF NOT EXISTS idx_op_logs_created_at ON operation_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_op_logs_user_created ON operation_logs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_op_logs_module_created ON operation_logs(module, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_op_logs_status ON operation_logs(status);

-- 添加注释
COMMENT ON TABLE operation_logs IS '通用操作日志表，记录所有用户操作用于审计';
COMMENT ON COLUMN operation_logs.user_id IS '操作用户ID';
COMMENT ON COLUMN operation_logs.username IS '操作用户名';
COMMENT ON COLUMN operation_logs.operation IS '操作类型：create, update, delete, approve, reject, submit, cancel, etc.';
COMMENT ON COLUMN operation_logs.module IS '模块名称：material_plan, inbound, outbound, stock, material, workflow, etc.';
COMMENT ON COLUMN operation_logs.resource_type IS '资源类型';
COMMENT ON COLUMN operation_logs.resource_id IS '资源ID';
COMMENT ON COLUMN operation_logs.resource_no IS '资源编号';
COMMENT ON COLUMN operation_logs.changes IS '变更内容JSON，记录before和after状态';
COMMENT ON COLUMN operation_logs.request_method IS 'HTTP请求方法';
COMMENT ON COLUMN operation_logs.request_path IS '请求路径';
COMMENT ON COLUMN operation_logs.request_params IS '请求参数JSON';
COMMENT ON COLUMN operation_logs.ip_address IS '客户端IP地址';
COMMENT ON COLUMN operation_logs.user_agent IS '客户端用户代理';
COMMENT ON COLUMN operation_logs.status IS '操作状态：success, error';

-- 为常见操作创建检查约束
ALTER TABLE operation_logs ADD CONSTRAINT chk_op_logs_operation
    CHECK (operation IN ('create', 'update', 'delete', 'approve', 'reject', 'submit', 'cancel', 'activate', 'complete', 'issue', 'receive', 'adjust', 'transfer', 'login', 'logout'));

ALTER TABLE operation_logs ADD CONSTRAINT chk_op_logs_status
    CHECK (status IN ('success', 'error'));

-- 创建分区表（可选，用于大数据量场景）
-- CREATE TABLE operation_logs_y2026m01 PARTITION OF operation_logs
--     FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

SELECT '010_add_operation_log.sql: operation_logs table created' AS status;
