-- 添加任务表的新字段
-- 添加优先级字段
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS priority VARCHAR(20) DEFAULT 'medium';

-- 添加状态字段
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'not_started';

-- 添加负责人字段
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS responsible VARCHAR(100);

-- 添加描述字段
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS description TEXT;

-- 添加索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
