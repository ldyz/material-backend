-- ============================================
-- 清空所有物资、计划和库存相关数据
-- ============================================
-- 警告：此操作将删除以下所有数据，执行前请确认已备份！
--
-- 将被清空的数据：
-- 1. 物资主数据 (material_master)
-- 2. 物资计划 (material_plans, material_plan_items)
-- 3. 入库单 (inbound_orders, inbound_items)
-- 4. 领料单 (requisitions, requisition_items)
-- 5. 库存 (stocks)
-- 6. 库存日志 (stock_logs)
-- 7. 库存操作日志 (stock_op_logs)
-- 8. 相关工作流实例 (workflow_instances, workflow_tasks, workflow_approvals)
-- ============================================

-- 开始事务（可以回滚）
BEGIN;

-- 1. 删除物资计划明细
DELETE FROM material_plan_items;

-- 2. 删除物资计划主表
DELETE FROM material_plans;

-- 3. 删除入库单明细
DELETE FROM inbound_items;

-- 4. 删除入库单主表
DELETE FROM inbound_orders;

-- 5. 删除领料单明细
DELETE FROM requisition_items;

-- 6. 删除领料单主表
DELETE FROM requisitions;

-- 7. 删除库存日志
DELETE FROM stock_logs;

-- 8. 删除库存操作日志
DELETE FROM stock_op_logs;

-- 9. 删除库存数据
DELETE FROM stocks;

-- 10. 删除物资主数据
DELETE FROM material_master;

-- 11. 删除相关的工作流实例（物资计划、入库单、领料单）
DELETE FROM workflow_pending_tasks
WHERE instance_id IN (
    SELECT id FROM workflow_instances
    WHERE business_type IN ('material_plan', 'inbound_order', 'requisition')
);

DELETE FROM workflow_logs
WHERE instance_id IN (
    SELECT id FROM workflow_instances
    WHERE business_type IN ('material_plan', 'inbound_order', 'requisition')
);

DELETE FROM workflow_approvals
WHERE instance_id IN (
    SELECT id FROM workflow_instances
    WHERE business_type IN ('material_plan', 'inbound_order', 'requisition')
);

DELETE FROM workflow_instances
WHERE business_type IN ('material_plan', 'inbound_order', 'requisition');

-- 重置序列
SELECT setval('material_plans_id_seq', 1, false);
SELECT setval('inbound_orders_id_seq', 1, false);
SELECT setval('requisitions_id_seq', 1, false);
SELECT setval('material_master_id_seq', 1, false);
SELECT setval('stocks_id_seq', 1, false);
SELECT setval('stock_logs_id_seq', 1, false);
SELECT setval('workflow_instances_id_seq', 1, false);
SELECT setval('workflow_approvals_id_seq', 1, false);

-- 提交事务
COMMIT;

-- ============================================
-- 执行完成后，可以运行以下查询确认数据已清空
-- ============================================
-- SELECT COUNT(*) FROM material_master;         -- 应该返回 0
-- SELECT COUNT(*) FROM material_plans;          -- 应该返回 0
-- SELECT COUNT(*) FROM material_plan_items;     -- 应该返回 0
-- SELECT COUNT(*) FROM inbound_orders;          -- 应该返回 0
-- SELECT COUNT(*) FROM inbound_items;           -- 应该返回 0
-- SELECT COUNT(*) FROM requisitions;            -- 应该返回 0
-- SELECT COUNT(*) FROM requisition_items;       -- 应该返回 0
-- SELECT COUNT(*) FROM stocks;                  -- 应该返回 0
-- SELECT COUNT(*) FROM stock_logs;              -- 应该返回 0
-- ============================================
