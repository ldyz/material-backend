-- 修复 inbound_orders 表的主键序列
-- 将序列重置为当前最大ID + 1

-- 重置主键序列
SELECT setval(
    pg_get_serial_sequence('inbound_orders', 'id'),
    COALESCE((SELECT MAX(id) FROM inbound_orders), 0) + 1,
    false
);

-- 重置 inbound_order_items 表的主键序列
SELECT setval(
    pg_get_serial_sequence('inbound_order_items', 'id'),
    COALESCE((SELECT MAX(id) FROM inbound_order_items), 0) + 1,
    false
);

-- 查看当前序列值（用于验证）
SELECT
    'inbound_orders' as table_name,
    pg_get_serial_sequence('inbound_orders', 'id') as sequence_name,
    nextval(pg_get_serial_sequence('inbound_orders', 'id')) - 1 as current_value
UNION ALL
SELECT
    'inbound_order_items' as table_name,
    pg_get_serial_sequence('inbound_order_items', 'id') as sequence_name,
    nextval(pg_get_serial_sequence('inbound_order_items', 'id')) - 1 as current_value;
