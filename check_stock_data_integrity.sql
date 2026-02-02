-- ========== 库存数据完整性检查 ==========

-- 1. 统计总库存数
SELECT '总库存记录数:' as description, COUNT(*) as count FROM stocks;

-- 2. 统计有材料名称的库存数
SELECT '有材料名称的库存记录数:' as description, COUNT(*) as count
FROM stocks s
LEFT JOIN materials m ON s.material_id = m.id
WHERE m.id IS NOT NULL AND m.name IS NOT NULL AND m.name != '';

-- 3. 统计孤立的库存记录（material_id指向不存在的记录）
SELECT '孤立的库存记录数:' as description, COUNT(*) as count
FROM stocks s
LEFT JOIN materials m ON s.material_id = m.id
WHERE s.material_id IS NOT NULL AND m.id IS NULL;

-- 4. 查看孤立的库存详情（前20条）
SELECT '孤立的库存记录详情:' as info,
       s.id as stock_id,
       s.material_id,
       s.quantity,
       s.unit
FROM stocks s
LEFT JOIN materials m ON s.material_id = m.id
WHERE s.material_id IS NOT NULL AND m.id IS NULL
ORDER BY s.id
LIMIT 20;

-- 5. 查看没有material_id的库存
SELECT 'material_id为NULL的库存记录数:' as description, COUNT(*) as count
FROM stocks
WHERE material_id IS NULL;

-- 6. 统计materials总数
SELECT '材料记录总数:' as description, COUNT(*) as count FROM materials;

-- 7. 检查ID范围
SELECT '库存中material_id的最大值:' as description, MAX(material_id) as max_value FROM stocks WHERE material_id IS NOT NULL;
SELECT '库存中material_id的最小值:' as description, MIN(material_id) as min_value FROM stocks WHERE material_id IS NOT NULL;
SELECT 'materials表中最大ID:' as description, MAX(id) as max_value FROM materials;

-- ========== 清理方案 ==========

-- 方案1: 删除孤立的库存记录（小心使用！）
-- DELETE FROM stocks
-- WHERE material_id IS NOT NULL
--   AND material_id NOT IN (SELECT id FROM materials);

-- 方案2: 将孤立的material_id设为NULL（推荐）
-- UPDATE stocks
-- SET material_id = NULL
-- WHERE material_id IS NOT NULL
--   AND material_id NOT IN (SELECT id FROM materials);

-- 执行后查看结果
-- SELECT '清理后剩余的孤立记录:' as description, COUNT(*) as count
-- FROM stocks s
-- LEFT JOIN materials m ON s.material_id = m.id
-- WHERE s.material_id IS NOT NULL AND m.id IS NULL;
