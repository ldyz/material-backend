-- 修复 roles 表 permissions 字段长度不足的问题
-- 运行此脚本来将 permissions 字段从 VARCHAR(500) 改为 TEXT (PostgreSQL)

ALTER TABLE roles ALTER COLUMN permissions TYPE TEXT;

-- 验证修改
SELECT column_name, data_type, character_maximum_length
FROM information_schema.columns
WHERE table_name = 'roles' AND column_name = 'permissions';
