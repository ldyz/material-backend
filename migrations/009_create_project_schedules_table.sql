-- 项目进度计划表
CREATE TABLE IF NOT EXISTS project_schedules (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    data JSONB NOT NULL DEFAULT '{"nodes": {}, "activities": {}}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER REFERENCES users(id),
    updated_by INTEGER REFERENCES users(id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_project_schedules_project_id ON project_schedules(project_id);
CREATE INDEX IF NOT EXISTS idx_project_schedules_created_at ON project_schedules(created_at);
CREATE INDEX IF NOT EXISTS idx_project_schedules_created_by ON project_schedules(created_by);

-- 添加注释
COMMENT ON TABLE project_schedules IS '项目进度计划表，存储项目的甘特图和双代号网络图数据';
COMMENT ON COLUMN project_schedules.project_id IS '项目ID，关联projects表';
COMMENT ON COLUMN project_schedules.data IS '进度数据，JSON格式，包含nodes（节点）和activities（活动）';
COMMENT ON COLUMN project_schedules.created_by IS '创建者ID，关联users表';
COMMENT ON COLUMN project_schedules.updated_by IS '更新者ID，关联users表';

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_project_schedules_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_project_schedules_updated_at ON project_schedules;
CREATE TRIGGER trigger_update_project_schedules_updated_at
    BEFORE UPDATE ON project_schedules
    FOR EACH ROW
    EXECUTE FUNCTION update_project_schedules_updated_at();
