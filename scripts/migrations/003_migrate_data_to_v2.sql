-- ============================================
-- 数据库重构迁移脚本 - 第2步: 数据迁移
-- ============================================
-- 说明：此脚本将旧表数据迁移到新表结构
-- 前置条件：002_create_new_tables_v2.sql 已执行
-- 注意：执行前建议备份数据库
-- ============================================

-- ============================================
-- 步骤 1: 创建物资主数据并迁移
-- ============================================
-- 从旧 materials 表中提取唯一物资作为主数据
INSERT INTO material_master (id, code, name, specification, unit, category, safety_stock, description, created_at, updated_at)
SELECT DISTINCT
    NEXTVAL('material_master_id_seq') as id,
    COALESCE(m.code, 'M-' || REPLACE(UPPER(md5(m.name || COALESCE(m.specification, ''))), '-', '')) as code,
    m.name as name,
    COALESCE(m.specification, '') as specification,
    COALESCE(m.unit, '') as unit,
    COALESCE(m.category, '') as category,
    0 as safety_stock,
    COALESCE(m.description, '') as description,
    NOW() as created_at,
    NOW() as updated_at
FROM materials m
WHERE m.name IS NOT NULL AND m.name != ''
-- 去重：相同的 name + specification 视为同一物资
GROUP BY m.name, m.specification, m.unit, m.category, m.description, m.code;

-- ============================================
-- 步骤 2: 创建映射表（旧物资ID -> 新主数据ID）
-- ============================================
-- 创建临时映射表
DROP TABLE IF EXISTS material_id_mapping;
CREATE TEMP TABLE material_id_mapping (
    old_material_id BIGINT,
    new_master_id BIGINT
);

-- 填充映射关系
INSERT INTO material_id_mapping (old_material_id, new_master_id)
SELECT
    m.id as old_material_id,
    mm.id as new_master_id
FROM materials m
JOIN material_master mm ON
    mm.name = m.name AND
    COALESCE(mm.specification, '') = COALESCE(m.specification, '');

-- 为映射表创建索引以提高性能
CREATE INDEX idx_material_id_mapping_old ON material_id_mapping(old_material_id);
CREATE INDEX idx_material_id_mapping_new ON material_id_mapping(new_master_id);

-- ============================================
-- 步骤 3: 迁移库存数据到 stocks_v2
-- ============================================
INSERT INTO stocks_v2 (id, project_id, material_id, quantity, safety_stock, location, unit_cost, created_at, updated_at)
SELECT DISTINCT
    NEXTVAL('stocks_v2_id_seq') as id,
    -- 尝试从旧 materials 表获取 project_id，如果为空则使用默认值
    COALESCE(m.project_id, 0) as project_id,
    map.new_master_id as material_id,
    COALESCE(s.quantity, 0) as quantity,
    COALESCE(s.safety_stock, 0) as safety_stock,
    COALESCE(s.location, '') as location,
    0 as unit_cost,
    s.created_at as created_at,
    s.updated_at as updated_at
FROM stocks s
JOIN materials m ON s.material_id = m.id
JOIN material_id_mapping map ON map.old_material_id = m.id;

-- ============================================
-- 步骤 4: 创建库存ID映射表（旧stock_id -> 新stock_id）
-- ============================================
DROP TABLE IF EXISTS stock_id_mapping;
CREATE TEMP TABLE stock_id_mapping (
    old_stock_id BIGINT,
    new_stock_id BIGINT
);

-- 填充库存映射关系
-- 注意：这里需要根据 material_id 和 project_id 来匹配
INSERT INTO stock_id_mapping (old_stock_id, new_stock_id)
SELECT
    s_old.id as old_stock_id,
    s_new.id as new_stock_id
FROM stocks s_old
JOIN materials m ON s_old.material_id = m.id
JOIN material_id_mapping mat_map ON mat_map.old_material_id = m.id
JOIN stocks_v2 s_new ON
    s_new.material_id = mat_map.new_master_id AND
    s_new.project_id = COALESCE(m.project_id, 0);

CREATE INDEX idx_stock_id_mapping_old ON stock_id_mapping(old_stock_id);
CREATE INDEX idx_stock_id_mapping_new ON stock_id_mapping(new_stock_id);

-- ============================================
-- 步骤 5: 迁移库存日志数据到 stock_logs_v2
-- ============================================
INSERT INTO stock_logs_v2 (
    id, stock_id, type, quantity, quantity_before, quantity_after,
    source_type, source_id, source_no, project_id, material_id,
    user_id, remark, created_at
)
SELECT
    NEXTVAL('stock_logs_v2_id_seq') as id,
    map_stock.new_stock_id as stock_id,
    sl.type as type,
    sl.quantity as quantity,
    sl.quantity_before as quantity_before,
    sl.quantity_after as quantity_after,
    -- 确定 source_type
    CASE
        WHEN sl.requisition_id IS NOT NULL THEN 'requisition'
        WHEN sl.inbound_code IS NOT NULL THEN 'inbound'
        ELSE 'adjust'
    END as source_type,
    -- 确定 source_id
    COALESCE(sl.requisition_id, 0) as source_id,
    -- 确定 source_no
    COALESCE(sl.requisition_code, sl.inbound_code, '') as source_no,
    sl.project_id as project_id,
    -- 通过 stock_id 找到对应的 material_id
    s_v2.material_id as material_id,
    sl.user_id as user_id,
    COALESCE(sl.remark, '') as remark,
    sl.time as created_at
FROM stock_logs sl
JOIN stock_id_mapping map_stock ON map_stock.old_stock_id = sl.stock_id
JOIN stocks_v2 s_v2 ON s_v2.id = map_stock.new_stock_id;

-- ============================================
-- 步骤 6: 迁移物资计划明细数据到 material_plan_items_v2
-- ============================================
INSERT INTO material_plan_items_v2 (
    id, plan_id, material_id, planned_quantity, unit_price,
    required_date, priority, status, remark, created_at, updated_at
)
SELECT
    NEXTVAL('material_plan_items_v2_id_seq') as id,
    mpi_old.plan_id as plan_id,
    -- 通过旧的 material_id 找到新的 master_id
    -- 首先尝试直接匹配 material_id
    COALESCE(
        (SELECT map.new_master_id FROM material_id_mapping map WHERE map.old_material_id = mpi_old.material_id),
        -- 如果没有直接映射，尝试通过名称匹配
        (SELECT mm.id FROM material_master mm WHERE mm.name = mpi_old.material_name LIMIT 1),
        0
    ) as material_id,
    mpi_old.planned_quantity as planned_quantity,
    mpi_old.unit_price / 100.0 as unit_price,  -- 从分转换为元
    mpi_old.required_date as required_date,
    mpi_old.priority as priority,
    mpi_old.status as status,
    COALESCE(mpi_old.remark, '') as remark,
    mpi_old.created_at as created_at,
    mpi_old.updated_at as updated_at
FROM material_plan_items mpi_old;

-- ============================================
-- 步骤 7: 迁移领料单明细数据到 requisition_items_v2
-- ============================================
INSERT INTO requisition_items_v2 (
    id, requisition_id, stock_id, material_id,
    requested_quantity, approved_quantity, actual_quantity,
    status, remark, created_at, updated_at
)
SELECT
    NEXTVAL('requisition_items_v2_id_seq') as id,
    ri_old.requisition_id as requisition_id,
    map_stock.new_stock_id as stock_id,
    s_v2.material_id as material_id,
    ri_old.quantity as requested_quantity,
    ri_old.approved_quantity as approved_quantity,
    0 as actual_quantity,  -- 旧数据可能没有这个字段
    ri_old.status as status,
    COALESCE(ri_old.remark, '') as remark,
    ri_old.created_at as created_at,
    ri_old.updated_at as updated_at
FROM requisition_items ri_old
JOIN stock_id_mapping map_stock ON map_stock.old_stock_id = ri_old.stock_id
JOIN stocks_v2 s_v2 ON s_v2.id = map_stock.new_stock_id;

-- ============================================
-- 步骤 8: 迁移入库单明细数据到 inbound_items_v2
-- ============================================
-- 注意：需要先检查 inbound_order_items 表结构是否与 inbound_items 一致
-- 这里假设表名为 inbound_items 或 inbound_order_items

-- 首先检查表是否存在
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'inbound_items') THEN
        -- 使用 inbound_items 表
        INSERT INTO inbound_items_v2 (
            id, inbound_order_id, stock_id, material_id,
            quantity, unit_price, status, remark, created_at
        )
        SELECT
            NEXTVAL('inbound_items_v2_id_seq') as id,
            ii_old.inbound_order_id as inbound_order_id,
            map_stock.new_stock_id as stock_id,
            s_v2.material_id as material_id,
            ii_old.quantity as quantity,
            ii_old.unit_price / 100.0 as unit_price,
            'received' as status,  -- 假设旧数据都是已入库
            COALESCE(ii_old.remark, '') as remark,
            NOW() as created_at
        FROM inbound_items ii_old
        JOIN stock_id_mapping map_stock ON map_stock.old_stock_id = ii_old.stock_id
        JOIN stocks_v2 s_v2 ON s_v2.id = map_stock.new_stock_id;
    ELSIF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'inbound_order_items') THEN
        -- 使用 inbound_order_items 表
        INSERT INTO inbound_items_v2 (
            id, inbound_order_id, stock_id, material_id,
            quantity, unit_price, status, remark, created_at
        )
        SELECT
            NEXTVAL('inbound_items_v2_id_seq') as id,
            ii_old.order_id as inbound_order_id,
            map_stock.new_stock_id as stock_id,
            s_v2.material_id as material_id,
            ii_old.quantity as quantity,
            ii_old.unit_price / 100.0 as unit_price,
            'received' as status,
            COALESCE(ii_old.remark, '') as remark,
            NOW() as created_at
        FROM inbound_order_items ii_old
        JOIN stock_id_mapping map_stock ON map_stock.old_stock_id = ii_old.stock_id
        JOIN stocks_v2 s_v2 ON s_v2.id = map_stock.new_stock_id;
    END IF;
END $$;

-- ============================================
-- 数据迁移验证查询
-- ============================================

-- 验证物资主数据
DO $$
BEGIN
    RAISE NOTICE '=== 数据迁移验证 ===';
    RAISE NOTICE '物资主数据记录数: %', (SELECT COUNT(*) FROM material_master);
    RAISE NOTICE '旧物资记录数: %', (SELECT COUNT(DISTINCT name, specification) FROM materials);
    RAISE NOTICE '新库存记录数: %', (SELECT COUNT(*) FROM stocks_v2);
    RAISE NOTICE '旧库存记录数: %', (SELECT COUNT(*) FROM stocks);
    RAISE NOTICE '新库存日志记录数: %', (SELECT COUNT(*) FROM stock_logs_v2);
    RAISE NOTICE '旧库存日志记录数: %', (SELECT COUNT(*) FROM stock_logs);
    RAISE NOTICE '新计划明细记录数: %', (SELECT COUNT(*) FROM material_plan_items_v2);
    RAISE NOTICE '旧计划明细记录数: %', (SELECT COUNT(*) FROM material_plan_items);
    RAISE NOTICE '新领料明细记录数: %', (SELECT COUNT(*) FROM requisition_items_v2);
    RAISE NOTICE '旧领料明细记录数: %', (SELECT COUNT(*) FROM requisition_items);
END $$;

-- ============================================
-- 注意事项
-- ============================================
-- 1. 检查映射表数据完整性
-- 2. 验证外键关系是否正确
-- 3. 确认业务数据（数量、金额）迁移正确
-- 4. 测试环境充分测试后再在生产环境执行
-- ============================================
