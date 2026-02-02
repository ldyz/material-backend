-- 为 stock_logs 表添加详细字段
-- 用于记录操作前后的数量、单价，以及关联的入库单/出库单号

-- 添加操作前数量
ALTER TABLE stock_logs ADD COLUMN IF NOT EXISTS quantity_before DECIMAL(15,3) DEFAULT 0;

-- 添加操作后数量
ALTER TABLE stock_logs ADD COLUMN IF NOT EXISTS quantity_after DECIMAL(15,3) DEFAULT 0;

-- 添加单价
ALTER TABLE stock_logs ADD COLUMN IF NOT EXISTS price DECIMAL(15,2) DEFAULT 0;

-- 添加入库单号
ALTER TABLE stock_logs ADD COLUMN IF NOT EXISTS inbound_code VARCHAR(50);

-- 添加出库单号（领料单号）
ALTER TABLE stock_logs ADD COLUMN IF NOT EXISTS requisition_code VARCHAR(50);

-- 添加注释
COMMENT ON COLUMN stock_logs.quantity_before IS '操作前数量';
COMMENT ON COLUMN stock_logs.quantity_after IS '操作后数量';
COMMENT ON COLUMN stock_logs.price IS '单价';
COMMENT ON COLUMN stock_logs.inbound_code IS '入库单号';
COMMENT ON COLUMN stock_logs.requisition_code IS '出库单号（领料单号）';
