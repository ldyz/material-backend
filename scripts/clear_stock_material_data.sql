-- 清空库存和物资数据
-- 执行时间: 2026-02-09
-- 警告: 此操作将删除所有库存和物资数据，不可恢复！

BEGIN;

-- 设置客户端编码
SET client_encoding = 'UTF8';

-- 显示将被删除的数据量
DO $$
DECLARE
    stock_count INT;
    stock_log_count INT;
    stock_op_log_count INT;
    category_count INT;
    master_count INT;
    inbound_count INT;
    requisition_count INT;
BEGIN
    SELECT COUNT(*) INTO stock_count FROM stocks;
    SELECT COUNT(*) INTO stock_log_count FROM stock_logs;
    SELECT COUNT(*) INTO stock_op_log_count FROM stock_op_logs;
    SELECT COUNT(*) INTO category_count FROM material_categories;
    SELECT COUNT(*) INTO master_count FROM material_master;
    SELECT COUNT(*) INTO inbound_count FROM inbound_orders;
    SELECT COUNT(*) INTO requisition_count FROM requisitions;

    RAISE NOTICE '即将删除的数据统计:';
    RAISE NOTICE '  库存记录: % 条', stock_count;
    RAISE NOTICE '  库存日志: % 条', stock_log_count;
    RAISE NOTICE '  库存操作日志: % 条', stock_op_log_count;
    RAISE NOTICE '  物资分类: % 条', category_count;
    RAISE NOTICE '  物资主数据: % 条', master_count;
    RAISE NOTICE '  入库单: % 条', inbound_count;
    RAISE NOTICE '  出库单: % 条', requisition_count;
END $$;

-- 1. 清空库存操作日志（无外键依赖）
DELETE FROM stock_op_logs;

-- 2. 清空库存日志（无外键依赖）
DELETE FROM stock_logs;

-- 3. 清空物资计划项
DELETE FROM material_plan_items;

-- 4. 清空物资计划
DELETE FROM material_plans;

-- 5. 清空入库单项
DELETE FROM inbound_items;

-- 6. 清空入库单
DELETE FROM inbound_orders;

-- 7. 清空出库单项
DELETE FROM requisition_items;

-- 8. 清空出库单
DELETE FROM requisitions;

-- 9. 清空库存记录
DELETE FROM stocks;

-- 10. 清空物资分类
DELETE FROM material_categories;

-- 11. 清空物资主数据
DELETE FROM material_master;

COMMIT;

-- 显示执行结果
DO $$
DECLARE
    remaining_stock INT;
    remaining_master INT;
BEGIN
    SELECT COUNT(*) INTO remaining_stock FROM stocks;
    SELECT COUNT(*) INTO remaining_master FROM material_master;

    RAISE NOTICE '数据清空完成！';
    RAISE NOTICE '  剩余库存记录: % 条', remaining_stock;
    RAISE NOTICE '  剩余物资主数据: % 条', remaining_master;
    RAISE NOTICE '  所有库存和物资数据已清空';
END $$;
