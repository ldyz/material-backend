-- 清空所有业务数据
-- 警告：此操作将删除所有入库单、出库单、库存和计划数据
-- 执行前请确认已备份重要数据

-- 1. 删除计划明细（先删除子表）
DELETE FROM material_plan_items;

-- 2. 删除计划主表
DELETE FROM material_plans;

-- 3. 删除入库单明细
DELETE FROM inbound_items;

-- 4. 删除入库单主表
DELETE FROM inbound_orders;

-- 5. 删除出库单明细
DELETE FROM requisition_items;

-- 6. 删除出库单主表
DELETE FROM requisitions;

-- 7. 删除库存数据
DELETE FROM stocks;

-- 重置序列
SELECT setval('material_plans_id_seq', 1, false);
SELECT setval('inbound_orders_id_seq', 1, false);
SELECT setval('requisitions_id_seq', 1, false);
