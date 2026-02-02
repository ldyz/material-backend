-- ========================================
-- 清空物资计划数据 SQL 脚本
-- ========================================
-- 警告：此脚本将删除所有物资计划及相关数据，操作不可逆！
-- 建议在执行前先备份数据库
-- ========================================

-- 开始事务
BEGIN;

-- 显示将要删除的数据量
SELECT '物资计划主表记录数:' as info, COUNT(*) as count FROM material_plans;
SELECT '物资计划明细记录数:' as info, COUNT(*) as count FROM material_plan_items;

-- 删除物资计划明细（因为有关联约束，需要先删除子表）
DELETE FROM material_plan_items;

-- 删除物资计划主表
DELETE FROM material_plans;

-- 重置序列（如果使用自增主键）
-- SELECT setval('material_plans_id_seq', 1, false);

-- 验证删除结果
SELECT '删除后物资计划主表记录数:' as info, COUNT(*) as count FROM material_plans;
SELECT '删除后物资计划明细记录数:' as info, COUNT(*) as count FROM material_plan_items;

-- 提交事务
-- COMMIT;
-- 如果确认无误，取消上面的 COMMIT 注释来执行删除
-- 如果需要撤销，执行 ROLLBACK;

-- 回滚事务（默认不执行删除）
ROLLBACK;
