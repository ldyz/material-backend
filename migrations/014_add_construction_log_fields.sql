-- 添加施工日志表的缺失字段
-- 执行时间: 2026-01-27

-- 检查并添加温度字段
DO $$
BEGIN
    ALTER TABLE construction_log ADD COLUMN temperature float;
EXCEPTION
    WHEN duplicate_column THEN null;
END $$;

-- 检查并添加施工进度字段
DO $$
BEGIN
    ALTER TABLE construction_log ADD COLUMN progress TEXT;
EXCEPTION
    WHEN duplicate_column THEN null;
END $$;

-- 检查并添加存在问题字段
DO $$
BEGIN
    ALTER TABLE construction_log ADD COLUMN issues TEXT;
EXCEPTION
    WHEN duplicate_column THEN null;
END $$;

-- 检查并添加日志日期字段
DO $$
BEGIN
    ALTER TABLE construction_log ADD COLUMN log_date VARCHAR(20);
EXCEPTION
    WHEN duplicate_column THEN null;
END $$;

-- 检查并添加备注字段
DO $$
BEGIN
    ALTER TABLE construction_log ADD COLUMN remark TEXT;
EXCEPTION
    WHEN duplicate_column THEN null;
END $$;

COMMENT ON COLUMN construction_log.temperature IS '温度';
COMMENT ON COLUMN construction_log.progress IS '施工进度';
COMMENT ON COLUMN construction_log.issues IS '存在问题';
COMMENT ON COLUMN construction_log.log_date IS '日志日期';
COMMENT ON COLUMN construction_log.remark IS '备注';
