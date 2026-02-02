-- Fix system_backup table to use proper serial type for auto-increment
-- First, drop the existing table if it exists (it will be recreated)
DROP TABLE IF EXISTS system_backup CASCADE;

-- Recreate with proper SERIAL type for ID
CREATE TABLE system_backup (
    id BIGSERIAL PRIMARY KEY,
    filename VARCHAR(200) NOT NULL,
    filepath VARCHAR(500) NOT NULL,
    size BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'completed',
    created_by VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);
