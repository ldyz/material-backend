-- 打卡记录表
CREATE TABLE IF NOT EXISTS attendance_records (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    appointment_id INTEGER,
    attendance_type VARCHAR(20) NOT NULL, -- morning, afternoon, noon_overtime, night_overtime
    clock_in_time TIMESTAMP NOT NULL,
    clock_in_location VARCHAR(200),
    clock_in_latitude DECIMAL(10, 7),
    clock_in_longitude DECIMAL(10, 7),
    overtime_hours DECIMAL(4, 1),
    remark TEXT,
    status VARCHAR(20) DEFAULT 'pending', -- pending, confirmed, rejected
    confirmed_by INTEGER,
    confirmed_at TIMESTAMP,
    confirmed_remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_attendance_records_user_id ON attendance_records(user_id);
CREATE INDEX IF NOT EXISTS idx_attendance_records_appointment_id ON attendance_records(appointment_id);
CREATE INDEX IF NOT EXISTS idx_attendance_records_clock_in_time ON attendance_records(clock_in_time);
CREATE INDEX IF NOT EXISTS idx_attendance_records_status ON attendance_records(status);

-- 月度考勤汇总表
CREATE TABLE IF NOT EXISTS monthly_attendance_summary (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    morning_count INTEGER DEFAULT 0,
    afternoon_count INTEGER DEFAULT 0,
    noon_overtime_hours DECIMAL(10, 1) DEFAULT 0,
    night_overtime_hours DECIMAL(10, 1) DEFAULT 0,
    total_work_days INTEGER DEFAULT 0,
    total_overtime_hours DECIMAL(10, 1) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'draft', -- draft, confirmed
    confirmed_by INTEGER,
    confirmed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建唯一索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_monthly_attendance_user_year_month ON monthly_attendance_summary(user_id, year, month);

-- 添加注释
COMMENT ON TABLE attendance_records IS '打卡记录表';
COMMENT ON COLUMN attendance_records.attendance_type IS '打卡类型：morning(上午), afternoon(下午), noon_overtime(中午加班), night_overtime(晚上加班)';
COMMENT ON COLUMN attendance_records.status IS '状态：pending(待确认), confirmed(已确认), rejected(已驳回)';

COMMENT ON TABLE monthly_attendance_summary IS '月度考勤汇总表';
COMMENT ON COLUMN monthly_attendance_summary.status IS '状态：draft(草稿), confirmed(已确认)';
