-- ============================================
-- 数据迁移修复脚本
-- ============================================

-- 1. 先处理物资计划明细迁移（处理重复键）
-- ============================================
-- 删除失败的迁移记录并重新迁移
TRUNCATE TABLE material_plan_items_v2;

INSERT INTO material_plan_items_v2 (
    id, plan_id, material_id, planned_quantity, unit_price,
    required_date, priority, status, remark, created_at, updated_at
)
SELECT
    NEXTVAL('material_plan_items_v2_id_seq') as id,
    mpi_old.plan_id as plan_id,
    -- 通过旧的 material_id 找到新的 master_id
    COALESCE(
        (SELECT mm.id FROM material_master mm WHERE mm.id = mpi_old.material_id),
        -- 如果没有直接映射，尝试通过名称匹配
        (SELECT mm.id FROM material_master mm WHERE mm.name = mpi_old.material_name LIMIT 1),
        0
    ) as material_id,
    mpi_old.planned_quantity::DECIMAL(15,3) as planned_quantity,
    mpi_old.unit_price / 100.0 as unit_price,
    mpi_old.required_date as required_date,
    mpi_old.priority as priority,
    mpi_old.status as status,
    COALESCE(mpi_old.remark, '') as remark,
    mpi_old.created_at as created_at,
    mpi_old.updated_at as updated_at
FROM material_plan_items mpi_old
WHERE COALESCE(
        (SELECT mm.id FROM material_master mm WHERE mm.id = mpi_old.material_id),
        (SELECT mm.id FROM material_master mm WHERE mm.name = mpi_old.material_name LIMIT 1),
        0
    ) > 0
ON CONFLICT (plan_id, material_id) DO NOTHING;

-- ============================================
-- 2. 为入库单明细创建缺失的 stock 记录
-- ============================================
-- 为入库单明细中涉及的项目+物资组合创建 stock 记录
INSERT INTO stocks_v2 (id, project_id, material_id, quantity, safety_stock, location, unit_cost, created_at, updated_at)
SELECT DISTINCT
    NEXTVAL('stocks_v2_id_seq') as id,
    COALESCE(
        (SELECT p.id FROM projects p LIMIT 1),  -- 获取第一个项目作为默认
        1  -- 如果没有项目则使用1
    ) as project_id,
    ii_old.material_id as material_id,
    0 as quantity,
    0 as safety_stock,
    '' as location,
    0 as unit_cost,
    NOW() as created_at,
    NOW() as updated_at
FROM inbound_order_items ii_old
WHERE NOT EXISTS (
    SELECT 1 FROM stocks_v2 s
    WHERE s.material_id = ii_old.material_id
    LIMIT 1
);

-- 更新 stock_id_mapping 表
DROP TABLE IF EXISTS stock_id_mapping;
CREATE TEMP TABLE stock_id_mapping (
    old_stock_id BIGINT,
    new_stock_id BIGINT
);

-- 填充映射关系
INSERT INTO stock_id_mapping (old_stock_id, new_stock_id)
SELECT
    s_old.id as old_stock_id,
    s_new.id as new_stock_id
FROM stocks s_old
JOIN material_master mm ON mm.id = s_old.material_id
JOIN stocks_v2 s_new ON
    s_new.material_id = mm.id;

CREATE INDEX idx_stock_id_mapping_old ON stock_id_mapping(old_stock_id);

-- ============================================
-- 3. 迁移领料单明细
-- ============================================
TRUNCATE TABLE requisition_items_v2;

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
    0 as actual_quantity,
    ri_old.status as status,
    COALESCE(ri_old.remark, '') as remark,
    ri_old.created_at as created_at,
    ri_old.updated_at as updated_at
FROM requisition_items ri_old
LEFT JOIN stock_id_mapping map_stock ON map_stock.old_stock_id = ri_old.stock_id
LEFT JOIN stocks_v2 s_v2 ON s_v2.id = map_stock.new_stock_id
WHERE map_stock.new_stock_id IS NOT NULL;

-- ============================================
-- 4. 迁移入库单明细
-- ============================================
-- 先通过 material_id 创建 stock 记录
INSERT INTO stocks_v2 (id, project_id, material_id, quantity, safety_stock, location, unit_cost, created_at, updated_at)
SELECT DISTINCT
    NEXTVAL('stocks_v2_id_seq') as id,
    1 as project_id,  -- 默认项目ID
    ii_old.material_id as material_id,
    0 as quantity,
    0 as safety_stock,
    '' as location,
    ii_old.unit_price / 100.0 as unit_cost,
    NOW() as created_at,
    NOW() as updated_at
FROM inbound_order_items ii_old
ON CONFLICT (project_id, material_id, warehouse_id) DO NOTHING;

-- 重新映射 stock_id
TRUNCATE TABLE stock_id_mapping;
INSERT INTO stock_id_mapping (old_stock_id, new_stock_id)
SELECT
    s_old.id as old_stock_id,
    s_new.id as new_stock_id
FROM stocks s_old
JOIN material_master mm ON mm.id = s_old.material_id
JOIN stocks_v2 s_new ON
    s_new.material_id = mm.id;

-- 迁移入库单明细
INSERT INTO inbound_items_v2 (
    id, inbound_order_id, stock_id, material_id,
    quantity, unit_price, status, remark, created_at
)
SELECT
    NEXTVAL('inbound_items_v2_id_seq') as id,
    ii_old.order_id as inbound_order_id,
    s_v2.id as stock_id,
    ii_old.material_id as material_id,
    ii_old.quantity::DECIMAL(15,3) as quantity,
    ii_old.unit_price / 100.0 as unit_price,
    'received' as status,
    COALESCE(ii_old.remark, '') as remark,
    NOW() as created_at
FROM inbound_order_items ii_old
JOIN stocks_v2 s_v2 ON s_v2.material_id = ii_old.material_id;

-- ============================================
-- 验证迁移结果
-- ============================================
DO $$
BEGIN
    RAISE NOTICE '=== 数据迁移验证 ===';
    RAISE NOTICE 'material_master: % 条', (SELECT COUNT(*) FROM material_master);
    RAISE NOTICE 'stocks_v2: % 条', (SELECT COUNT(*) FROM stocks_v2);
    RAISE NOTICE 'stock_logs_v2: % 条', (SELECT COUNT(*) FROM stock_logs_v2);
    RAISE NOTICE 'material_plan_items_v2: % 条', (SELECT COUNT(*) FROM material_plan_items_v2);
    RAISE NOTICE 'requisition_items_v2: % 条', (SELECT COUNT(*) FROM requisition_items_v2);
    RAISE NOTICE 'inbound_items_v2: % 条', (SELECT COUNT(*) FROM inbound_items_v2);
END $$;

COMMENT ON TABLE material_master IS '物资主数据表（全局唯一）';
COMMENT ON TABLE stocks_v2 IS '库存表（项目维度）';
COMMENT ON TABLE stock_logs_v2 IS '库存变动日志表（统一记录）';
COMMENT ON TABLE material_plan_items_v2 IS '物资计划明细表（简化结构）';
COMMENT ON TABLE requisition_items_v2 IS '领料单明细表（简化结构）';
COMMENT ON TABLE inbound_items_v2 IS '入库单明细表（简化结构）';
