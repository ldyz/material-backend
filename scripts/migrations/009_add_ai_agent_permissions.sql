-- Migration: Add AI Agent Permissions
-- Version: 009
-- Date: 2025-02-03
-- Description: Adds AI Agent related permissions to the roles table

-- Note: In this system, permissions are stored as comma-separated strings in the roles.permissions column.
-- This script adds new permissions that can be assigned to roles.

-- AI Agent Permissions:
-- These permissions control access to AI Agent functionality

-- To add these permissions to the admin role:
UPDATE roles
SET permissions = CONCAT(
    COALESCE(permissions, ''),
    CASE
        WHEN COALESCE(permissions, '') = '' THEN ''
        ELSE ','
    END,
    'ai_agent_view,ai_agent_query,ai_agent_operate,ai_agent_workflow,ai_agent_logs'
)
WHERE name = 'admin' AND permissions NOT LIKE '%ai_agent_%';

-- Create additional role for AI operators (optional)
INSERT INTO roles (name, description, permissions, created_at)
VALUES (
    'ai_operator',
    'AI Agent Operator - Can use AI features for queries and basic operations',
    'ai_agent_view,ai_agent_query',
    CURRENT_TIMESTAMP
)
ON CONFLICT (name) DO NOTHING;

-- Verify the permissions were added
SELECT name, permissions FROM roles WHERE name = 'admin' OR name = 'ai_operator';

-- Documentation of new permissions:
/*
 * ai_agent_view    - View AI Agent capabilities and configuration
 * ai_agent_query   - Use AI query and analyze features
 * ai_agent_operate - Execute AI operations (create plans, update stock, etc.)
 * ai_agent_workflow - AI workflow operations (approve/reject tasks)
 * ai_agent_logs    - View AI Agent operation logs
 */

-- Migration complete
SELECT '009_add_ai_agent_permissions.sql: AI agent permissions added' AS status;
