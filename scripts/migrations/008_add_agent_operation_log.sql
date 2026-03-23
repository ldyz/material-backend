-- Migration: Add Agent Operation Log Table
-- Version: 008
-- Date: 2025-02-03
-- Description: Creates table for logging AI Agent operations for audit trail

-- AI Agent Operation Logs Table
CREATE TABLE IF NOT EXISTS agent_operation_logs (
    id BIGSERIAL PRIMARY KEY,
    operation VARCHAR(100) NOT NULL,           -- Operation type (query, analyze, create_material_plan, etc.)
    resource VARCHAR(100) NOT NULL,            -- Resource type (material, stock, workflow, etc.)
    parameters JSONB,                          -- Operation parameters (flexible schema)
    reasoning TEXT,                            -- AI reasoning process/explanation
    result JSONB,                              -- Operation result
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    agent_id VARCHAR(255),                     -- Agent identifier (e.g., claude-desktop-v1)
    status VARCHAR(50) NOT NULL,               -- pending/completed/failed
    error TEXT,                                -- Error message if failed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_agent_logs_user ON agent_operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_agent_logs_status ON agent_operation_logs(status);
CREATE INDEX IF NOT EXISTS idx_agent_logs_operation ON agent_operation_logs(operation);
CREATE INDEX IF NOT EXISTS idx_agent_logs_resource ON agent_operation_logs(resource);
CREATE INDEX IF NOT EXISTS idx_agent_logs_agent_id ON agent_operation_logs(agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_logs_created_at ON agent_operation_logs(created_at DESC);

-- Create composite index for common queries
CREATE INDEX IF NOT EXISTS idx_agent_logs_user_status ON agent_operation_logs(user_id, status);
CREATE INDEX IF NOT EXISTS idx_agent_logs_operation_created ON agent_operation_logs(operation, created_at DESC);

-- Add comment for documentation
COMMENT ON TABLE agent_operation_logs IS 'Logs all AI Agent operations for audit and debugging';
COMMENT ON COLUMN agent_operation_logs.operation IS 'Type of operation performed (query, analyze, create_material_plan, etc.)';
COMMENT ON COLUMN agent_operation_logs.resource IS 'Target resource type (material, stock, workflow, etc.)';
COMMENT ON COLUMN agent_operation_logs.parameters IS 'Operation parameters stored as JSONB';
COMMENT ON COLUMN agent_operation_logs.reasoning IS 'AI explanation of why this operation was chosen';
COMMENT ON COLUMN agent_operation_logs.result IS 'Operation result stored as JSONB';
COMMENT ON COLUMN agent_operation_logs.agent_id IS 'Identifier for the AI agent (e.g., claude-desktop-v1)';
COMMENT ON COLUMN agent_operation_logs.status IS 'Operation status: pending, completed, or failed';

-- Migration complete
SELECT '008_add_agent_operation_log.sql: agent_operation_logs table created' AS status;
