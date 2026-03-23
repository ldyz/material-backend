-- Migration: Fix inbound_items table
-- Description: Rename inbound_items_v2 to inbound_items and set up auto-increment

-- Rename table from inbound_items_v2 to inbound_items if it doesn't exist
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'inbound_items_v2') THEN
        ALTER TABLE IF EXISTS inbound_items_v2 RENAME TO inbound_items_old;
    END IF;
END $$;

-- Create inbound_items table with proper auto-increment if it doesn't exist
CREATE TABLE IF NOT EXISTS inbound_items (
    id BIGSERIAL PRIMARY KEY,
    inbound_order_id BIGINT NOT NULL,
    stock_id BIGINT,
    material_id BIGINT NOT NULL,
    quantity DECIMAL(15,3) NOT NULL,
    unit_price DECIMAL(15,2),
    status VARCHAR(20) DEFAULT 'pending',
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_inbound_items_order ON inbound_items(inbound_order_id);
CREATE INDEX IF NOT EXISTS idx_inbound_items_stock ON inbound_items(stock_id);
CREATE INDEX IF NOT EXISTS idx_inbound_items_material ON inbound_items(material_id);

-- Migrate data from old table if exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'inbound_items_old') THEN
        INSERT INTO inbound_items (inbound_order_id, stock_id, material_id, quantity, unit_price, status, remark, created_at)
        SELECT inbound_order_id, stock_id, material_id, quantity, unit_price, status, remark, created_at
        FROM inbound_items_old
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

-- Drop old table after migration
DROP TABLE IF EXISTS inbound_items_old;
