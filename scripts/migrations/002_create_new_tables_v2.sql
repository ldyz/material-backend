-- ============================================
-- 数据库重构迁移脚本 - 第1步: 创建新表结构
-- ============================================
-- 说明：此脚本创建新的表结构，保留旧表不变
-- 执行时机：业务低峰期
-- ============================================

-- ============================================
-- 1. 物资主数据表（全局唯一）
-- ============================================
CREATE TABLE IF NOT EXISTS material_master (
    id BIGINT PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,          -- 物资编码（全局唯一）
    name VARCHAR(200) NOT NULL,                -- 物资名称
    specification VARCHAR(200),                -- 规格型号
    unit VARCHAR(20),                          -- 基本单位
    category VARCHAR(100),                      -- 分类
    safety_stock DECIMAL(15,3) DEFAULT 0,      -- 默认安全库存
    description TEXT,                           -- 描述
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_material_master_code ON material_master(code);
CREATE INDEX IF NOT EXISTS idx_material_master_category ON material_master(category);
CREATE INDEX IF NOT EXISTS idx_material_master_name ON material_master(name);

-- ============================================
-- 2. 新库存表结构
-- ============================================
CREATE TABLE IF NOT EXISTS stocks_v2 (
    id BIGINT PRIMARY KEY,
    project_id BIGINT NOT NULL,                -- 项目ID
    material_id BIGINT NOT NULL,               -- 关联 material_master.id
    warehouse_id BIGINT,                       -- 仓库ID（可选，支持多仓库）
    quantity DECIMAL(15,3) DEFAULT 0,          -- 当前库存
    safety_stock DECIMAL(15,3) DEFAULT 0,      -- 项目级安全库存
    location VARCHAR(100),                     -- 货位
    unit_cost DECIMAL(15,2) DEFAULT 0,         -- 单位成本
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(project_id, material_id, warehouse_id)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_stocks_v2_project ON stocks_v2(project_id);
CREATE INDEX IF NOT EXISTS idx_stocks_v2_material ON stocks_v2(material_id);
CREATE INDEX IF NOT EXISTS idx_stocks_v2_project_material ON stocks_v2(project_id, material_id);

-- ============================================
-- 3. 新库存变动日志表结构
-- ============================================
CREATE TABLE IF NOT EXISTS stock_logs_v2 (
    id BIGINT PRIMARY KEY,
    stock_id BIGINT NOT NULL,                  -- 关联 stocks_v2.id
    type VARCHAR(10) NOT NULL,                 -- 'in' 入库, 'out' 出库
    quantity DECIMAL(15,3) NOT NULL,           -- 变动数量
    quantity_before DECIMAL(15,3) DEFAULT 0,   -- 变动前
    quantity_after DECIMAL(15,3) DEFAULT 0,    -- 变动后
    source_type VARCHAR(20) NOT NULL,          -- 来源类型：inbound/requisition/adjust/transfer
    source_id BIGINT,                          -- 来源单据ID
    source_no VARCHAR(50),                     -- 来源单据号
    project_id BIGINT NOT NULL,                -- 项目ID（冗余，便于查询）
    material_id BIGINT NOT NULL,               -- 物资ID（冗余，便于查询）
    user_id BIGINT,                            -- 操作人
    remark VARCHAR(500),                       -- 备注
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_stock_logs_v2_stock ON stock_logs_v2(stock_id);
CREATE INDEX IF NOT EXISTS idx_stock_logs_v2_source ON stock_logs_v2(source_type, source_id);
CREATE INDEX IF NOT EXISTS idx_stock_logs_v2_project ON stock_logs_v2(project_id);
CREATE INDEX IF NOT EXISTS idx_stock_logs_v2_material ON stock_logs_v2(material_id);
CREATE INDEX IF NOT EXISTS idx_stock_logs_v2_time ON stock_logs_v2(created_at DESC);

-- ============================================
-- 4. 物资计划明细表（新结构）
-- ============================================
CREATE TABLE IF NOT EXISTS material_plan_items_v2 (
    id BIGINT PRIMARY KEY,
    plan_id BIGINT NOT NULL,                   -- 关联 material_plans.id
    material_id BIGINT NOT NULL,               -- 关联 material_master.id
    planned_quantity DECIMAL(15,3) NOT NULL,   -- 计划数量
    unit_price DECIMAL(15,2) DEFAULT 0,        -- 计划单价
    required_date DATE,                        -- 需求日期
    priority VARCHAR(20) DEFAULT 'normal',
    status VARCHAR(20) DEFAULT 'pending',       -- pending/partial/completed
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(plan_id, material_id)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_plan_items_v2_plan ON material_plan_items_v2(plan_id);
CREATE INDEX IF NOT EXISTS idx_plan_items_v2_material ON material_plan_items_v2(material_id);

-- ============================================
-- 5. 领料单明细表（新结构）
-- ============================================
CREATE TABLE IF NOT EXISTS requisition_items_v2 (
    id BIGINT PRIMARY KEY,
    requisition_id BIGINT NOT NULL,            -- 关联 requisitions.id
    stock_id BIGINT NOT NULL,                  -- 关联 stocks_v2.id
    material_id BIGINT NOT NULL,               -- 关联 material_master.id（冗余）
    requested_quantity DECIMAL(15,3) NOT NULL, -- 申请数量
    approved_quantity DECIMAL(15,3),           -- 批准数量
    actual_quantity DECIMAL(15,3),             -- 实际发放数量
    status VARCHAR(20) DEFAULT 'pending',      -- pending/approved/issued/rejected
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_req_items_v2_requisition ON requisition_items_v2(requisition_id);
CREATE INDEX IF NOT EXISTS idx_req_items_v2_stock ON requisition_items_v2(stock_id);
CREATE INDEX IF NOT EXISTS idx_req_items_v2_material ON requisition_items_v2(material_id);

-- ============================================
-- 6. 入库单明细表（新结构）
-- ============================================
CREATE TABLE IF NOT EXISTS inbound_items_v2 (
    id BIGINT PRIMARY KEY,
    inbound_order_id BIGINT NOT NULL,          -- 关联 inbound_orders.id
    stock_id BIGINT NOT NULL,                  -- 关联 stocks_v2.id
    material_id BIGINT NOT NULL,               -- 关联 material_master.id（冗余）
    quantity DECIMAL(15,3) NOT NULL,           -- 入库数量
    unit_price DECIMAL(15,2),                  -- 入库单价
    status VARCHAR(20) DEFAULT 'pending',      -- pending/received
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_inbound_items_v2_order ON inbound_items_v2(inbound_order_id);
CREATE INDEX IF NOT EXISTS idx_inbound_items_v2_stock ON inbound_items_v2(stock_id);
CREATE INDEX IF NOT EXISTS idx_inbound_items_v2_material ON inbound_items_v2(material_id);

-- ============================================
-- 创建序列用于 ID 生成
-- ============================================
CREATE SEQUENCE IF NOT EXISTS material_master_id_seq START 1;
CREATE SEQUENCE IF NOT EXISTS stocks_v2_id_seq START 1;
CREATE SEQUENCE IF NOT EXISTS stock_logs_v2_id_seq START 1;
CREATE SEQUENCE IF NOT EXISTS material_plan_items_v2_id_seq START 1;
CREATE SEQUENCE IF NOT EXISTS requisition_items_v2_id_seq START 1;
CREATE SEQUENCE IF NOT EXISTS inbound_items_v2_id_seq START 1;

COMMENT ON TABLE material_master IS '物资主数据表（全局唯一）';
COMMENT ON TABLE stocks_v2 IS '库存表（项目维度）';
COMMENT ON TABLE stock_logs_v2 IS '库存变动日志表（统一记录）';
COMMENT ON TABLE material_plan_items_v2 IS '物资计划明细表（简化结构）';
COMMENT ON TABLE requisition_items_v2 IS '领料单明细表（简化结构）';
COMMENT ON TABLE inbound_items_v2 IS '入库单明细表（简化结构）';
