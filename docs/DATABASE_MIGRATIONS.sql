-- ====================================================================
-- Gantt Chart Editor - Database Migrations
-- Version: 1.0.0
-- Last Updated: 2026-02-19
--
-- This migration adds all required tables, indexes, and constraints
-- for the Gantt Chart Editor functionality.
--
-- Run this migration:
-- psql -U username -d database_name -f DATABASE_MIGRATIONS.sql
-- ====================================================================

-- ====================================================================
-- 1. GANTT TASKS TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_tasks (
    id VARCHAR(50) PRIMARY KEY,
    project_id VARCHAR(50) NOT NULL,
    parent_id VARCHAR(50),

    -- Basic fields
    name VARCHAR(255) NOT NULL,
    description TEXT,

    -- Dates and duration
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    duration INTEGER NOT NULL DEFAULT 1, -- in days

    -- Progress and status
    progress INTEGER DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
    status VARCHAR(20) DEFAULT 'not_started'
        CHECK (status IN ('not_started', 'in_progress', 'completed', 'delayed', 'blocked')),

    -- Priority
    priority VARCHAR(10) DEFAULT 'medium'
        CHECK (priority IN ('low', 'medium', 'high', 'critical')),

    -- Assignment
    assignee_id VARCHAR(50),
    assignee_name VARCHAR(100),

    -- Position/ordering
    position INTEGER NOT NULL DEFAULT 0,

    -- Visual
    color VARCHAR(7), -- hex color code
    milestone BOOLEAN DEFAULT false,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),

    -- Foreign keys
    CONSTRAINT fk_gantt_tasks_project
        FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_gantt_tasks_parent
        FOREIGN KEY (parent_id) REFERENCES gantt_tasks(id) ON DELETE SET NULL,
    CONSTRAINT fk_gantt_tasks_assignee
        FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL,

    -- Check constraints
    CONSTRAINT chk_gantt_tasks_dates
        CHECK (end_date >= start_date),
    CONSTRAINT chk_gantt_tasks_duration
        CHECK (duration > 0)
);

-- Indexes for gantt_tasks
CREATE INDEX idx_gantt_tasks_project_id ON gantt_tasks(project_id);
CREATE INDEX idx_gantt_tasks_parent_id ON gantt_tasks(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_gantt_tasks_status ON gantt_tasks(status);
CREATE INDEX idx_gantt_tasks_priority ON gantt_tasks(priority);
CREATE INDEX idx_gantt_tasks_assignee ON gantt_tasks(assignee_id) WHERE assignee_id IS NOT NULL;
CREATE INDEX idx_gantt_tasks_dates ON gantt_tasks(start_date, end_date);
CREATE INDEX idx_gantt_tasks_position ON gantt_tasks(project_id, position);

-- Full-text search on task name and description
CREATE INDEX idx_gantt_tasks_search ON gantt_tasks
    USING gin(to_tsvector('english', coalesce(name, '') || ' ' || coalesce(description, '')));

-- ====================================================================
-- 2. GANTT DEPENDENCIES TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_dependencies (
    id VARCHAR(50) PRIMARY KEY,
    project_id VARCHAR(50) NOT NULL,
    from_task_id VARCHAR(50) NOT NULL,
    to_task_id VARCHAR(50) NOT NULL,

    -- Dependency type
    type VARCHAR(20) DEFAULT 'finish-to-start'
        CHECK (type IN ('finish-to-start', 'start-to-start', 'finish-to-finish', 'start-to-finish')),

    -- Lag (delay) in days
    lag INTEGER DEFAULT 0,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),

    -- Foreign keys
    CONSTRAINT fk_gantt_deps_project
        FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_gantt_deps_from_task
        FOREIGN KEY (from_task_id) REFERENCES gantt_tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_gantt_deps_to_task
        FOREIGN KEY (to_task_id) REFERENCES gantt_tasks(id) ON DELETE CASCADE,

    -- Prevent circular dependencies (application-level check)
    CONSTRAINT chk_gantt_deps_not_same
        CHECK (from_task_id != to_task_id)
);

-- Unique constraint to prevent duplicate dependencies
CREATE UNIQUE INDEX idx_gantt_deps_unique ON gantt_dependencies(from_task_id, to_task_id);

-- Indexes
CREATE INDEX idx_gantt_deps_project ON gantt_dependencies(project_id);
CREATE INDEX idx_gantt_deps_from_task ON gantt_dependencies(from_task_id);
CREATE INDEX idx_gantt_deps_to_task ON gantt_dependencies(to_task_id);

-- ====================================================================
-- 3. GANTT CONSTRAINTS TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_constraints (
    id VARCHAR(50) PRIMARY KEY,
    task_id VARCHAR(50) NOT NULL UNIQUE,

    -- Constraint type
    type VARCHAR(30) NOT NULL
        CHECK (type IN ('start-no-earlier-than', 'finish-no-later-than',
                       'must-start-on', 'must-finish-on')),

    -- Constraint date
    constraint_date TIMESTAMP NOT NULL,

    -- Additional notes
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),

    -- Foreign key
    CONSTRAINT fk_gantt_constraints_task
        FOREIGN KEY (task_id) REFERENCES gantt_tasks(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_gantt_constraints_task ON gantt_constraints(task_id);

-- ====================================================================
-- 4. GANTT COMMENTS TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_comments (
    id VARCHAR(50) PRIMARY KEY,
    task_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    user_name VARCHAR(100) NOT NULL,
    user_avatar VARCHAR(255),

    -- Comment content
    content TEXT NOT NULL,

    -- Mentions (array of user IDs)
    mentions TEXT[],

    -- Attachments (JSON array of file references)
    attachments JSONB,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign keys
    CONSTRAINT fk_gantt_comments_task
        FOREIGN KEY (task_id) REFERENCES gantt_tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_gantt_comments_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_gantt_comments_task ON gantt_comments(task_id);
CREATE INDEX idx_gantt_comments_user ON gantt_comments(user_id);
CREATE INDEX idx_gantt_comments_created ON gantt_comments(created_at DESC);

-- Full-text search on comment content
CREATE INDEX idx_gantt_comments_search ON gantt_comments
    USING gin(to_tsvector('english', content));

-- ====================================================================
-- 5. GANTT CHANGE LOG TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_change_log (
    id VARCHAR(50) PRIMARY KEY,
    task_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    user_name VARCHAR(100) NOT NULL,

    -- Change details (JSON array of changed fields)
    changes JSONB NOT NULL,

    -- Old and new values (JSONB)
    old_values JSONB,
    new_values JSONB,

    -- Metadata
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT,

    -- Foreign key
    CONSTRAINT fk_gantt_changelog_task
        FOREIGN KEY (task_id) REFERENCES gantt_tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_gantt_changelog_user
        FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Indexes
CREATE INDEX idx_gantt_changelog_task ON gantt_change_log(task_id);
CREATE INDEX idx_gantt_changelog_user ON gantt_change_log(user_id);
CREATE INDEX idx_gantt_changelog_timestamp ON gantt_change_log(timestamp DESC);

-- JSONB indexes for querying changes
CREATE INDEX idx_gantt_changelog_changes ON gantt_change_log USING gin(changes);
CREATE INDEX idx_gantt_changelog_old_values ON gantt_change_log USING gin(old_values);
CREATE INDEX idx_gantt_changelog_new_values ON gantt_change_log USING gin(new_values);

-- ====================================================================
-- 6. GANTT REPORTS TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_reports (
    id VARCHAR(50) PRIMARY KEY,
    project_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL,

    -- Report details
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    format VARCHAR(10) NOT NULL,

    -- Report configuration (JSONB)
    config JSONB NOT NULL,

    -- Generated report file reference
    file_path VARCHAR(500),
    file_size BIGINT,

    -- Expiration
    expires_at TIMESTAMP NOT NULL,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),

    -- Foreign keys
    CONSTRAINT fk_gantt_reports_project
        FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_gantt_reports_user
        FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Indexes
CREATE INDEX idx_gantt_reports_project ON gantt_reports(project_id);
CREATE INDEX idx_gantt_reports_user ON gantt_reports(user_id);
CREATE INDEX idx_gantt_reports_expires ON gantt_reports(expires_at);

-- Clean up expired reports
CREATE INDEX idx_gantt_reports_expired ON gantt_reports(expires_at)
    WHERE expires_at < CURRENT_TIMESTAMP;

-- ====================================================================
-- 7. GANTT AI SUGGESTIONS TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_ai_suggestions (
    id VARCHAR(50) PRIMARY KEY,
    project_id VARCHAR(50) NOT NULL,
    task_id VARCHAR(50),

    -- Suggestion details
    type VARCHAR(50) NOT NULL,
    priority VARCHAR(10) DEFAULT 'medium'
        CHECK (priority IN ('low', 'medium', 'high')),

    title VARCHAR(255) NOT NULL,
    description TEXT,

    -- Suggestion data (JSONB)
    impact TEXT,
    effort TEXT,
    actions JSONB,

    -- Status
    status VARCHAR(20) DEFAULT 'pending'
        CHECK (status IN ('pending', 'accepted', 'dismissed', 'expired')),

    -- AI confidence (0-1)
    confidence NUMERIC(3,2) CHECK (confidence >= 0 AND confidence <= 1),

    -- Response
    response_note TEXT,
    responded_at TIMESTAMP,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign keys
    CONSTRAINT fk_gantt_suggestions_project
        FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_gantt_suggestions_task
        FOREIGN KEY (task_id) REFERENCES gantt_tasks(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_gantt_suggestions_project ON gantt_ai_suggestions(project_id);
CREATE INDEX idx_gantt_suggestions_task ON gantt_ai_suggestions(task_id);
CREATE INDEX idx_gantt_suggestions_status ON gantt_ai_suggestions(status);
CREATE INDEX idx_gantt_suggestions_type ON gantt_ai_suggestions(type);
CREATE INDEX idx_gantt_suggestions_created ON gantt_ai_suggestions(created_at DESC);

-- ====================================================================
-- 8. GANTN TEMPLATES TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_templates (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),

    -- Template data (JSONB)
    tasks JSONB NOT NULL,
    dependencies JSONB,

    -- Metadata
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),

    -- Usage tracking
    usage_count INTEGER DEFAULT 0
);

-- Indexes
CREATE INDEX idx_gantt_templates_category ON gantt_templates(category);
CREATE INDEX idx_gantt_templates_public ON gantt_templates(is_public) WHERE is_public = true;
CREATE INDEX idx_gantt_templates_usage ON gantt_templates(usage_count DESC);

-- Full-text search on template name and description
CREATE INDEX idx_gantt_templates_search ON gantt_templates
    USING gin(to_tsvector('english', coalesce(name, '') || ' ' || coalesce(description, '')));

-- ====================================================================
-- 9. GANTT CALENDAR EXCEPTIONS TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_calendar_exceptions (
    id VARCHAR(50) PRIMARY KEY,
    project_id VARCHAR(50) NOT NULL,
    calendar_id VARCHAR(50) NOT NULL,

    -- Exception date or date range
    exception_date DATE,
    start_date DATE,
    end_date DATE,

    -- Exception type
    type VARCHAR(20) NOT NULL
        CHECK (type IN ('holiday', 'non-working', 'special')),

    -- Description
    name VARCHAR(255),
    description TEXT,

    -- Working hours (for special days)
    working_hours JSONB,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),

    -- Foreign keys
    CONSTRAINT fk_gantt_calendar_exceptions_project
        FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,

    -- Check constraints
    CONSTRAINT chk_gantt_calendar_exceptions_dates
        CHECK (
            (exception_date IS NOT NULL AND start_date IS NULL AND end_date IS NULL) OR
            (exception_date IS NULL AND start_date IS NOT NULL AND end_date IS NOT NULL)
        )
);

-- Indexes
CREATE INDEX idx_gantt_calendar_exceptions_project ON gantt_calendar_exceptions(project_id);
CREATE INDEX idx_gantt_calendar_exceptions_calendar ON gantt_calendar_exceptions(calendar_id);
CREATE INDEX idx_gantt_calendar_exceptions_date ON gantt_calendar_exceptions(exception_date);
CREATE INDEX idx_gantt_calendar_exceptions_dates ON gantt_calendar_exceptions(start_date, end_date);

-- ====================================================================
-- 10. GANTT RESOURCES TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_resources (
    id VARCHAR(50) PRIMARY KEY,
    project_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50),

    -- Resource details
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) DEFAULT 'human'
        CHECK (type IN ('human', 'material', 'equipment')),

    -- Capacity
    max_hours_per_day NUMERIC(5,2) DEFAULT 8.00,
    cost_per_hour NUMERIC(10,2),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign keys
    CONSTRAINT fk_gantt_resources_project
        FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_gantt_resources_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes
CREATE INDEX idx_gantt_resources_project ON gantt_resources(project_id);
CREATE INDEX idx_gantt_resources_user ON gantt_resources(user_id);
CREATE INDEX idx_gantt_resources_type ON gantt_resources(type);

-- ====================================================================
-- 11. GANTN TASK ASSIGNMENTS TABLE
-- ====================================================================

CREATE TABLE IF NOT EXISTS gantt_task_assignments (
    id VARCHAR(50) PRIMARY KEY,
    task_id VARCHAR(50) NOT NULL,
    resource_id VARCHAR(50) NOT NULL,

    -- Assignment details
    assigned_units NUMERIC(5,2) DEFAULT 1.00, -- Percentage or count
    role VARCHAR(100),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),

    -- Foreign keys
    CONSTRAINT fk_gantt_assignments_task
        FOREIGN KEY (task_id) REFERENCES gantt_tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_gantt_assignments_resource
        FOREIGN KEY (resource_id) REFERENCES gantt_resources(id) ON DELETE CASCADE,

    -- Unique constraint
    CONSTRAINT uniq_gantt_assignments_task_resource
        UNIQUE (task_id, resource_id)
);

-- Indexes
CREATE INDEX idx_gantt_assignments_task ON gantt_task_assignments(task_id);
CREATE INDEX idx_gantt_assignments_resource ON gantt_task_assignments(resource_id);

-- ====================================================================
-- 12. FUNCTIONS AND TRIGGERS
-- ====================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_gantt_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to tables with updated_at
CREATE TRIGGER trigger_gantt_tasks_updated_at
    BEFORE UPDATE ON gantt_tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_gantt_updated_at();

CREATE TRIGGER trigger_gantt_constraints_updated_at
    BEFORE UPDATE ON gantt_constraints
    FOR EACH ROW
    EXECUTE FUNCTION update_gantt_updated_at();

CREATE TRIGGER trigger_gantt_templates_updated_at
    BEFORE UPDATE ON gantt_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_gantt_updated_at();

CREATE TRIGGER trigger_gantt_resources_updated_at
    BEFORE UPDATE ON gantt_resources
    FOR EACH ROW
    EXECUTE FUNCTION update_gantt_updated_at();

-- ====================================================================
-- 13. VIEWS FOR COMMON QUERIES
-- ====================================================================

-- View for task with dependencies
CREATE OR REPLACE VIEW v_gantt_tasks_with_deps AS
SELECT
    t.*,
    COUNT(de.id) as dependency_count,
    ARRAY_AGG(de.to_task_id) FILTER (WHERE de.id IS NOT NULL) as dependent_task_ids
FROM gantt_tasks t
LEFT JOIN gantt_dependencies de ON de.from_task_id = t.id
GROUP BY t.id;

-- View for task summary statistics
CREATE OR REPLACE VIEW v_gantt_project_stats AS
SELECT
    project_id,
    COUNT(*) as total_tasks,
    COUNT(*) FILTER (WHERE status = 'completed') as completed_tasks,
    COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress_tasks,
    COUNT(*) FILTER (WHERE status = 'not_started') as not_started_tasks,
    COUNT(*) FILTER (WHERE status = 'delayed') as delayed_tasks,
    COUNT(*) FILTER (WHERE status = 'blocked') as blocked_tasks,
    AVG(progress) as avg_progress,
    MIN(start_date) as earliest_start,
    MAX(end_date) as latest_end
FROM gantt_tasks
GROUP BY project_id;

-- View for resource utilization
CREATE OR REPLACE VIEW v_gantt_resource_utilization AS
SELECT
    r.id,
    r.name,
    r.type,
    r.max_hours_per_day,
    COUNT(ta.id) as assignment_count,
    SUM(ta.assigned_units) as total_assigned_units
FROM gantt_resources r
LEFT JOIN gantt_task_assignments ta ON ta.resource_id = r.id
GROUP BY r.id, r.name, r.type, r.max_hours_per_day;

-- ====================================================================
-- 14. CLEANUP JOBS
-- ====================================================================

-- Function to clean up expired reports
CREATE OR REPLACE FUNCTION cleanup_expired_reports()
RETURNS void AS $$
BEGIN
    DELETE FROM gantt_reports
    WHERE expires_at < CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;

-- Function to clean up old change logs (older than 90 days)
CREATE OR REPLACE FUNCTION cleanup_old_change_logs()
RETURNS void AS $$
BEGIN
    DELETE FROM gantt_change_log
    WHERE timestamp < CURRENT_TIMESTAMP - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- ====================================================================
-- 15. GRANT PERMISSIONS (Adjust as needed)
-- ====================================================================

-- Grant permissions to application user
-- Replace 'your_app_user' with your actual database user
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_app_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO your_app_user;
-- GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO your_app_user;

-- ====================================================================
-- 16. SAMPLE DATA (Optional - Comment out for production)
-- ====================================================================

-- Insert sample tasks for testing
-- INSERT INTO gantt_tasks (id, project_id, name, start_date, end_date, duration, progress, status, priority, position)
-- VALUES
--     ('task-1', 'proj-1', 'Project Planning', '2024-01-01', '2024-01-05', 5, 100, 'completed', 'high', 0),
--     ('task-2', 'proj-1', 'Design Phase', '2024-01-06', '2024-01-15', 10, 60, 'in_progress', 'high', 1),
--     ('task-3', 'proj-1', 'Development', '2024-01-16', '2024-02-15', 30, 0, 'not_started', 'high', 2);

-- ====================================================================
-- MIGRATION COMPLETE
-- ====================================================================

-- Display success message
DO $$
BEGIN
    RAISE NOTICE '====================================================================';
    RAISE NOTICE 'Gantt Chart Database Migration Complete!';
    RAISE NOTICE '====================================================================';
    RAISE NOTICE 'Created tables:';
    RAISE NOTICE '  - gantt_tasks';
    RAISE NOTICE '  - gantt_dependencies';
    RAISE NOTICE '  - gantt_constraints';
    RAISE NOTICE '  - gantt_comments';
    RAISE NOTICE '  - gantt_change_log';
    RAISE NOTICE '  - gantt_reports';
    RAISE NOTICE '  - gantt_ai_suggestions';
    RAISE NOTICE '  - gantt_templates';
    RAISE NOTICE '  - gantt_calendar_exceptions';
    RAISE NOTICE '  - gantt_resources';
    RAISE NOTICE '  - gantt_task_assignments';
    RAISE NOTICE '';
    RAISE NOTICE 'Created views:';
    RAISE NOTICE '  - v_gantt_tasks_with_deps';
    RAISE NOTICE '  - v_gantt_project_stats';
    RAISE NOTICE '  - v_gantt_resource_utilization';
    RAISE NOTICE '';
    RAISE NOTICE 'Created functions:';
    RAISE NOTICE '  - update_gantt_updated_at()';
    RAISE NOTICE '  - cleanup_expired_reports()';
    RAISE NOTICE '  - cleanup_old_change_logs()';
    RAISE NOTICE '====================================================================';
END $$;
