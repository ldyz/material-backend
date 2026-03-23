-- Add arrived_quantity column to material_plan_items if it doesn't exist
ALTER TABLE material_plan_items
ADD COLUMN IF NOT EXISTS arrived_quantity DECIMAL(15,3) DEFAULT 0;

COMMENT ON COLUMN material_plan_items.arrived_quantity IS '已到货数量';
