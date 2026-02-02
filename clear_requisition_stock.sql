-- 清空领料单和库存相关数据
-- 执行顺序：先删除子表，再删除主表

-- 1. 清空库存操作日志 (stock_op_logs)
DELETE FROM stock_op_logs;

-- 2. 清空库存日志 (stock_logs)
DELETE FROM stock_logs;

-- 3. 清空库存 (stocks)
DELETE FROM stocks;

-- 4. 清空领料单项目 (requisition_items)
DELETE FROM requisition_items;

-- 5. 清空领料单 (requisitions)
DELETE FROM requisitions;

-- 显示清空结果
SELECT 'requisitions 清空完成' AS status;
SELECT 'requisition_items 清空完成' AS status;
SELECT 'stocks 清空完成' AS status;
SELECT 'stock_logs 清空完成' AS status;
SELECT 'stock_op_logs 清空完成' AS status;
