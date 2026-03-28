-- 添加 plan_item_id 字段到 requisition_items 表，用于关联物资计划明细
-- 支持出库单关联物资计划功能

ALTER TABLE requisition_items ADD COLUMN IF NOT EXISTS plan_item_id INTEGER REFERENCES material_plan_items(id);

-- 添加索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_requisition_items_plan_item_id ON requisition_items(plan_item_id);
