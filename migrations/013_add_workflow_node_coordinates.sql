-- Add x and y coordinates to workflow_nodes table
-- This migration adds position fields for automatic layout in the workflow editor

-- Add x column
ALTER TABLE workflow_nodes ADD COLUMN IF NOT EXISTS x INT DEFAULT 0;

-- Add y column
ALTER TABLE workflow_nodes ADD COLUMN IF NOT EXISTS y INT DEFAULT 0;

-- Add comment
COMMENT ON COLUMN workflow_nodes.x IS 'Node X coordinate for canvas layout';
COMMENT ON COLUMN workflow_nodes.y IS 'Node Y coordinate for canvas layout';
