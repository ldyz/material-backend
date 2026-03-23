-- 清空日志类数据表
-- 执行前请确认已备份重要数据

-- 清空库存日志
TRUNCATE TABLE stock_logs CASCADE;

-- 清空库存操作日志
TRUNCATE TABLE stock_op_logs CASCADE;

-- 清空工作流审批记录
TRUNCATE TABLE workflow_approvals CASCADE;

-- 清空工作流日志
TRUNCATE TABLE workflow_logs CASCADE;
