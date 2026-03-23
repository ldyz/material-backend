-- Migration: Add sort_order column to material_plan_items
-- Description: Adds the missing sort_order column for ordering plan items

-- Add sort_order column if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'material_plan_items'
        AND column_name = 'sort_order'
    ) THEN
        ALTER TABLE material_plan_items ADD COLUMN sort_order INTEGER DEFAULT 0;
    END IF;
END $$;

-- Create index on sort_order for better query performance
CREATE INDEX IF NOT EXISTS idx_material_plan_items_sort_order ON material_plan_items(sort_order);
