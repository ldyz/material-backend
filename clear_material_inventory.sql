-- =============================================
-- 清空物资管理和库存管理相关数据
-- 执行顺序：先删除子表，再删除主表
-- =============================================

BEGIN;

-- 1. 库存管理相关表
-- 1.1 清空库存操作日志
DELETE FROM stock_op_logs;
-- 1.2 清空库存日志
DELETE FROM stock_logs;
-- 1.3 清空库存
DELETE FROM stocks;

-- 2. 领料单相关表
-- 2.1 清空领料单项目
DELETE FROM requisition_items;
-- 2.2 清空领料单
DELETE FROM requisitions;

-- 3. 入库单相关表
-- 3.1 清空入库单项目
DELETE FROM inbound_items;
-- 3.2 清空入库单
DELETE FROM inbounds;

-- 4. 请购单相关表
-- 4.1 清空请购单项目
DELETE FROM request_items;
-- 4.2 清空请购单
DELETE FROM requests;

-- 5. 物资计划相关表
-- 5.1 清空物资计划项目
DELETE FROM material_plan_items;
-- 5.2 清空物资计划
DELETE FROM material_plans;

-- 6. 物资主数据表
DELETE FROM materials;

-- 7. 重置序列（如果有使用自增ID）
-- SELECT setval('material_plans_id_seq', 1, false);
-- SELECT setval('material_plan_items_id_seq', 1, false);
-- SELECT setval('stocks_id_seq', 1, false);
-- SELECT setval('stock_logs_id_seq', 1, false);
-- SELECT setval('stock_op_logs_id_seq', 1, false);
-- SELECT setval('requisitions_id_seq', 1, false);
-- SELECT setval('requisition_items_id_seq', 1, false);
-- SELECT setval('inbounds_id_seq', 1, false);
-- SELECT setval('inbound_items_id_seq', 1, false);
-- SELECT setval('requests_id_seq', 1, false);
-- SELECT setval('request_items_id_seq', 1, false);
-- SELECT setval('materials_id_seq', 1, false);

COMMIT;

-- =============================================
-- 显示清空结果统计
-- =============================================
SELECT 'material_plans' AS table_name, COUNT(*) AS remaining_count FROM material_plans
UNION ALL
SELECT 'material_plan_items', COUNT(*) FROM material_plan_items
UNION ALL
SELECT 'materials', COUNT(*) FROM materials
UNION ALL
SELECT 'stocks', COUNT(*) FROM stocks
UNION ALL
SELECT 'stock_logs', COUNT(*) FROM stock_logs
UNION ALL
SELECT 'stock_op_logs', COUNT(*) FROM stock_op_logs
UNION ALL
SELECT 'requisitions', COUNT(*) FROM requisitions
UNION ALL
SELECT 'requisition_items', COUNT(*) FROM requisition_items
UNION ALL
SELECT 'inbounds', COUNT(*) FROM inbounds
UNION ALL
SELECT 'inbound_items', COUNT(*) FROM inbound_items
UNION ALL
SELECT 'requests', COUNT(*) FROM requests
UNION ALL
SELECT 'request_items', COUNT(*) FROM request_items
ORDER BY table_name;

-- 显示完成消息
SELECT '============================================' AS info;
SELECT '物资管理和库存管理数据清空完成！' AS message;
SELECT '============================================' AS info;
