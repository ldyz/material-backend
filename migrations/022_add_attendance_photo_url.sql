-- 添加打卡照片URL字段
ALTER TABLE attendance_records ADD COLUMN IF NOT EXISTS photo_url VARCHAR(500);
COMMENT ON COLUMN attendance_records.photo_url IS '打卡照片URL';
