-- 优化到货数量计算的复合索引
CREATE INDEX IF NOT EXISTS idx_inbound_items_material_quantity
ON inbound_order_items(material_id, quantity);

CREATE INDEX IF NOT EXISTS idx_inbound_status_material
ON inbound_orders(status, id)
WHERE status IN ('approved', 'completed');
