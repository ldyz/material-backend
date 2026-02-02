-- 修复 project_schedules 表的外键约束
-- 允许 created_by 和 updated_by 为 NULL

-- 删除旧的外键约束
ALTER TABLE project_schedules DROP CONSTRAINT IF EXISTS project_schedules_created_by_fkey;
ALTER TABLE project_schedules DROP CONSTRAINT IF EXISTS project_schedules_updated_by_fkey;

-- 修改字段为可空
ALTER TABLE project_schedules ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE project_schedules ALTER COLUMN updated_by DROP NOT NULL;

-- 重新添加外键约束（允许 NULL）
ALTER TABLE project_schedules
ADD CONSTRAINT project_schedules_created_by_fkey
FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE project_schedules
ADD CONSTRAINT project_schedules_updated_by_fkey
FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;

-- 添加注释
COMMENT ON COLUMN project_schedules.created_by IS '创建者ID（可为空）';
COMMENT ON COLUMN project_schedules.updated_by IS '更新者ID（可为空）';
