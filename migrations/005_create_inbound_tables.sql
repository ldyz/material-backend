-- Create inbound_orders table
CREATE TABLE IF NOT EXISTS inbound_orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(50) NOT NULL UNIQUE,
    supplier VARCHAR(100),
    contact VARCHAR(50),
    project_id BIGINT,
    creator_id BIGINT NOT NULL,
    creator_name VARCHAR(80) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    notes TEXT,
    remark TEXT,
    total_amount BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_inbound_project_id ON inbound_orders(project_id);
CREATE INDEX IF NOT EXISTS idx_inbound_creator_id ON inbound_orders(creator_id);
CREATE INDEX IF NOT EXISTS idx_inbound_status ON inbound_orders(status);
CREATE INDEX IF NOT EXISTS idx_inbound_created_at ON inbound_orders(created_at);

-- Create inbound_order_items table
CREATE TABLE IF NOT EXISTS inbound_order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    material_id BIGINT NOT NULL,
    quantity BIGINT NOT NULL DEFAULT 0,
    unit_price BIGINT NOT NULL DEFAULT 0,
    remark TEXT,
    FOREIGN KEY (order_id) REFERENCES inbound_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_inbound_items_order_id ON inbound_order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_inbound_items_material_id ON inbound_order_items(material_id);
