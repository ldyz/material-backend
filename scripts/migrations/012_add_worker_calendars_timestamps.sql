-- 012_add_worker_calendars_timestamps.sql
-- 为 worker_calendars 表添加 created_at 和 updated_at 列

-- 添加 created_at 列
ALTER TABLE worker_calendars
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- 添加 updated_at 列
ALTER TABLE worker_calendars
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- 为已有的记录设置初始时间戳
UPDATE worker_calendars
SET created_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE created_at IS NULL OR updated_at IS NULL;

-- 添加注释
COMMENT ON COLUMN worker_calendars.created_at IS '创建时间';
COMMENT ON COLUMN worker_calendars.updated_at IS '更新时间';

-- 创建自动更新 updated_at 的触发器
DROP TRIGGER IF EXISTS update_worker_calendars_updated_at ON worker_calendars;

CREATE TRIGGER update_worker_calendars_updated_at
    BEFORE UPDATE ON worker_calendars
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

SELECT '012_add_worker_calendars_timestamps.sql: added created_at and updated_at columns to worker_calendars' AS status;
