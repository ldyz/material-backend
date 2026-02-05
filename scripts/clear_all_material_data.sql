-- ============================================
-- 清空所有物资相关数据
-- 日期: 2025-02-05
-- 注意: 此操作不可逆，请先备份数据库！
-- ============================================

-- 备份数据库（请在执行前手动备份）
-- pg_dump -U materials -d materials > backup_$(date +%Y%m%d_%H%M%S).sql

BEGIN;

-- ============================================
-- 第一步：清空物资计划相关数据
-- ============================================

-- 删除物资计划明细
DELETE FROM material_plan_items;

-- 重置自增序列
-- ALTER SEQUENCE material_plan_items_id_seq RESTART WITH 1;

-- 删除物资计划主表
DELETE FROM material_plans;

-- 重置自增序列
-- ALTER SEQUENCE material_plans_id_seq RESTART WITH 1;


-- ============================================
-- 第二步：清空库存相关数据
-- ============================================

-- 删除库存日志
DELETE FROM stock_logs;

-- 删除库存操作日志
DELETE FROM stock_op_logs;

-- 删除库存记录
DELETE FROM stocks;

-- 重置自增序列
-- ALTER SEQUENCE stocks_id_seq RESTART WITH 1;
-- ALTER SEQUENCE stock_logs_id_seq RESTART WITH 1;


-- ============================================
-- 第三步：清空入库单相关数据
-- ============================================

-- 删除入库单明细
DELETE FROM inbound_items;

-- 删除入库单主表
DELETE FROM inbound_orders;

-- 重置自增序列
-- ALTER SEQUENCE inbound_orders_id_seq RESTART WITH 1;
-- ALTER SEQUENCE inbound_items_id_seq RESTART WITH 1;


-- ============================================
-- 第四步：清空出库单相关数据
-- ============================================

-- 删除出库单明细
DELETE FROM requisition_items;

-- 删除出库单主表
DELETE FROM requisitions;

-- 重置自增序列
-- ALTER SEQUENCE requisitions_id_seq RESTART WITH 1;
-- ALTER SEQUENCE requisition_items_id_seq RESTART WITH 1;


-- ============================================
-- 第五步：清空物资主数据
-- ============================================

DELETE FROM material_master;

-- 重置自增序列
-- ALTER SEQUENCE material_master_id_seq RESTART WITH 1;


-- ============================================
-- 第六步：清空物资分类（可选）
-- ============================================

-- DELETE FROM material_categories;
-- ALTER SEQUENCE material_categories_id_seq RESTART WITH 1;


COMMIT;

-- ============================================
-- 验证清空结果
-- ============================================

-- 查看各表记录数
SELECT
    'material_master' AS table_name, COUNT(*) AS count FROM material_master
UNION ALL
SELECT 'material_plans', COUNT(*) FROM material_plans
UNION ALL
SELECT 'material_plan_items', COUNT(*) FROM material_plan_items
UNION ALL
SELECT 'stocks', COUNT(*) FROM stocks
UNION ALL
SELECT 'stock_logs', COUNT(*) FROM stock_logs
UNION ALL
SELECT 'inbound_orders', COUNT(*) FROM inbound_orders
UNIONION ALL
SELECT 'inbound_items', COUNT(*) FROM inbound_items
UNION ALL
SELECT 'requisitions', COUNT(*) FROM requisitions
UNION ALL
SELECT 'requisition_items', COUNT(*) FROM requisition_items;

-- ============================================
-- 完成提示
-- ============================================

DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE '数据清空完成！';
    RAISE NOTICE '已清空以下表：';
    RAISE NOTICE '1. material_master (物资主数据)';
    RAISE NOTICE '2. material_plans + material_plan_items (物资计划)';
    RAISE NOTICE '3. stocks + stock_logs (库存数据)';
    RAISE NOTICE '4. inbound_orders + inbound_items (入库单)';
    RAISE NOTICE '5. requisitions + requisition_items (出库单)';
    RAISE NOTICE '============================================';
END $$;
