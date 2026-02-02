-- Migration: Add plan_id foreign keys to inbound_orders and requisitions
-- Description: Associates inbound orders and requisitions with material plans

-- Add plan_id column to inbound_orders table
ALTER TABLE inbound_orders
ADD COLUMN IF NOT EXISTS plan_id BIGINT REFERENCES material_plans(id) ON DELETE SET NULL;

-- Add plan_id column to requisitions table
ALTER TABLE requisitions
ADD COLUMN IF NOT EXISTS plan_id BIGINT REFERENCES material_plans(id) ON DELETE SET NULL;

-- Create indexes for plan_id columns
CREATE INDEX IF NOT EXISTS idx_inbound_plan ON inbound_orders(plan_id);
CREATE INDEX IF NOT EXISTS idx_requisition_plan ON requisitions(plan_id);

-- Add comment for documentation
COMMENT ON COLUMN inbound_orders.plan_id IS '关联的物资计划ID';
COMMENT ON COLUMN requisitions.plan_id IS '关联的物资计划ID';
