-- 添加 photo_urls 字段用于存储多张照片URL（JSON数组）
ALTER TABLE attendance_records ADD COLUMN IF NOT EXISTS photo_urls TEXT;

-- 添加注释
COMMENT ON COLUMN attendance_records.photo_urls IS '打卡照片URLs（多张，JSON数组格式）';
