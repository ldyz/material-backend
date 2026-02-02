-- Migration: Add project hierarchy and task management tables
-- This migration adds support for 4-level project hierarchy and task management

-- Add project hierarchy fields to projects table
ALTER TABLE projects
  ADD COLUMN IF NOT EXISTS parent_id BIGINT REFERENCES projects(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS level INTEGER DEFAULT 0 CHECK (level >= 0 AND level <= 3),
  ADD COLUMN IF NOT EXISTS path VARCHAR(500),
  ADD COLUMN IF NOT EXISTS progress_percentage DECIMAL(5,2) DEFAULT 0;

-- Create indexes for hierarchy queries
CREATE INDEX IF NOT EXISTS idx_projects_parent_id ON projects(parent_id);
CREATE INDEX IF NOT EXISTS idx_projects_level ON projects(level);
CREATE INDEX IF NOT EXISTS idx_projects_path ON projects(path);

-- Function to update project hierarchy (level and path)
CREATE OR REPLACE FUNCTION update_project_hierarchy()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.parent_id IS NULL THEN
    -- Root level project
    NEW.level := 0;
    IF NEW.id IS NULL THEN
      -- For INSERT, id will be assigned, path will be updated by another trigger
      NEW.path := '/temp/';
    ELSE
      NEW.path := '/' || NEW.id || '/';
    END IF;
  ELSE
    -- Child project
    SELECT level INTO NEW.level FROM projects WHERE id = NEW.parent_id;
    IF NEW.level IS NULL THEN
      RAISE EXCEPTION 'Parent project does not exist';
    END IF;
    IF NEW.level >= 3 THEN
      RAISE EXCEPTION 'Maximum hierarchy depth (3) reached. Cannot create more than 4 levels.';
    END IF;
    NEW.level := NEW.level + 1;

    SELECT path INTO NEW.path FROM projects WHERE id = NEW.parent_id;
    IF NEW.path IS NULL THEN
      RAISE EXCEPTION 'Parent project path is invalid';
    END IF;
    IF NEW.id IS NULL THEN
      NEW.path := NEW.path || 'temp/';
    ELSE
      NEW.path := NEW.path || NEW.id || '/';
    END IF;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update hierarchy on insert/update
DROP TRIGGER IF EXISTS trigger_update_project_hierarchy ON projects;
CREATE TRIGGER trigger_update_project_hierarchy
  BEFORE INSERT OR UPDATE OF parent_id ON projects
  FOR EACH ROW
  EXECUTE FUNCTION update_project_hierarchy();

-- Function to update path after id is assigned
CREATE OR REPLACE FUNCTION update_project_path()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.path IS NOT NULL AND (NEW.path LIKE '%/temp/%' OR NEW.path = '/temp/') THEN
    IF NEW.parent_id IS NULL THEN
      NEW.path := '/' || NEW.id || '/';
    ELSE
      SELECT path INTO NEW.path FROM projects WHERE id = NEW.parent_id;
      NEW.path := NEW.path || NEW.id || '/';
    END IF;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_project_path ON projects;
CREATE TRIGGER trigger_update_project_path
  AFTER INSERT ON projects
  FOR EACH ROW
  EXECUTE FUNCTION update_project_path();

-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    schedule_id INTEGER REFERENCES project_schedules(id) ON DELETE CASCADE,
    parent_id INTEGER REFERENCES tasks(id) ON DELETE SET NULL,
    name VARCHAR(200) NOT NULL,
    duration DECIMAL(10,2),
    start_date DATE,
    end_date DATE,
    progress DECIMAL(5,2) DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
    is_milestone BOOLEAN DEFAULT FALSE,
    sort_order INTEGER DEFAULT 0,
    position_x DECIMAL(10,2),
    position_y DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for tasks
CREATE INDEX IF NOT EXISTS idx_tasks_project ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_schedule ON tasks(schedule_id);
CREATE INDEX IF NOT EXISTS idx_tasks_parent ON tasks(parent_id);
CREATE INDEX IF NOT EXISTS idx_tasks_sort_order ON tasks(sort_order);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_tasks_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_tasks_updated_at ON tasks;
CREATE TRIGGER trigger_update_tasks_updated_at
  BEFORE UPDATE ON tasks
  FOR EACH ROW
  EXECUTE FUNCTION update_tasks_updated_at();

-- Create task dependencies table
CREATE TABLE IF NOT EXISTS task_dependencies (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    depends_on INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    type VARCHAR(2) DEFAULT 'FS' CHECK (type IN ('FS', 'FF', 'SS', 'SF')),
    lag INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(task_id, depends_on)
);

-- Create indexes for task dependencies
CREATE INDEX IF NOT EXISTS idx_task_dependencies_task ON task_dependencies(task_id);
CREATE INDEX IF NOT EXISTS idx_task_dependencies_depends_on ON task_dependencies(depends_on);

-- Comment on tables
COMMENT ON TABLE tasks IS 'Stores project tasks/activities for progress management';
COMMENT ON TABLE task_dependencies IS 'Stores dependencies between tasks (FS=Finish-to-Start, FF=Finish-to-Finish, SS=Start-to-Start, SF=Start-to-Finish)';
COMMENT ON COLUMN tasks.is_milestone IS 'Indicates if this task is a milestone (zero duration event)';
COMMENT ON COLUMN tasks.position_x IS 'X coordinate for network diagram visualization';
COMMENT ON COLUMN tasks.position_y IS 'Y coordinate for network diagram visualization';
COMMENT ON COLUMN task_dependencies.type IS 'Dependency type: FS, FF, SS, or SF';
COMMENT ON COLUMN task_dependencies.lag IS 'Lag time in days between tasks';
