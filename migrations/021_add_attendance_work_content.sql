-- 添加 work_content 字段到打卡记录表
ALTER TABLE attendance_records ADD COLUMN IF NOT EXISTS work_content TEXT;

COMMENT ON COLUMN attendance_records.work_content IS '手动填写的工作内容（无关联任务时）';
