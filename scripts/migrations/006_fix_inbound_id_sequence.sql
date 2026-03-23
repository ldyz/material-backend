-- Migration: Fix inbound_items table id column auto-increment
-- Description: Set up sequence for inbound_items id column

-- Drop the table if it exists with wrong schema and recreate with correct schema
DROP TABLE IF EXISTS inbound_items CASCADE;

-- Create inbound_items table with proper auto-increment
CREATE TABLE inbound_items (
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
CREATE INDEX idx_inbound_items_order ON inbound_items(inbound_order_id);
CREATE INDEX idx_inbound_items_stock ON inbound_items(stock_id);
CREATE INDEX idx_inbound_items_material ON inbound_items(material_id);
