-- 011_create_appointment_tables.sql
-- 创建施工预约相关表

-- 施工预约单表
CREATE TABLE IF NOT EXISTS construction_appointments (
    id SERIAL PRIMARY KEY,
    appointment_no VARCHAR(50) UNIQUE NOT NULL,
    project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,

    -- 预约信息
    applicant_id INTEGER NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    applicant_name VARCHAR(100) NOT NULL,
    contact_phone VARCHAR(20),
    contact_person VARCHAR(100),

    -- 作业信息
    work_date DATE NOT NULL,
    time_slot VARCHAR(50) NOT NULL,
    work_location VARCHAR(500) NOT NULL,
    work_content TEXT NOT NULL,
    work_type VARCHAR(50),

    -- 优先级
    is_urgent BOOLEAN DEFAULT FALSE,
    priority INTEGER DEFAULT 0,
    urgent_reason TEXT,

    -- 分配信息
    assigned_worker_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    assigned_worker_name VARCHAR(100),

    -- 状态
    status VARCHAR(20) DEFAULT 'draft',
    workflow_instance_id INTEGER REFERENCES workflow_instances(id) ON DELETE SET NULL,

    -- 时间戳
    submitted_at TIMESTAMP,
    approved_at TIMESTAMP,
    completed_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 作业人员日历表
CREATE TABLE IF NOT EXISTS worker_calendars (
    id SERIAL PRIMARY KEY,
    worker_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    calendar_date DATE NOT NULL,
    time_slot VARCHAR(20) NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    status VARCHAR(20) DEFAULT 'available',
    appointment_id INTEGER REFERENCES construction_appointments(id) ON DELETE SET NULL,
    blocked_reason TEXT,

    UNIQUE(worker_id, calendar_date, time_slot)
);

-- 创建索引 - construction_appointments
CREATE INDEX IF NOT EXISTS idx_appointments_no ON construction_appointments(appointment_no);
CREATE INDEX IF NOT EXISTS idx_appointments_project ON construction_appointments(project_id);
CREATE INDEX IF NOT EXISTS idx_appointments_applicant ON construction_appointments(applicant_id);
CREATE INDEX IF NOT EXISTS idx_appointments_worker ON construction_appointments(assigned_worker_id);
CREATE INDEX IF NOT EXISTS idx_appointments_work_date ON construction_appointments(work_date);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON construction_appointments(status);
CREATE INDEX IF NOT EXISTS idx_appointments_workflow ON construction_appointments(workflow_instance_id);
CREATE INDEX IF NOT EXISTS idx_appointments_urgent ON construction_appointments(is_urgent);
CREATE INDEX IF NOT EXISTS idx_appointments_date_status ON construction_appointments(work_date, status);

-- 创建索引 - worker_calendars
CREATE INDEX IF NOT EXISTS idx_worker_calendars_worker ON worker_calendars(worker_id);
CREATE INDEX IF NOT EXISTS idx_worker_calendars_date ON worker_calendars(calendar_date);
CREATE INDEX IF NOT EXISTS idx_worker_calendars_slot ON worker_calendars(time_slot);
CREATE INDEX IF NOT EXISTS idx_worker_calendars_available ON worker_calendars(is_available);
CREATE INDEX IF NOT EXISTS idx_worker_calendars_status ON worker_calendars(status);
CREATE INDEX IF NOT EXISTS idx_worker_calendars_appointment ON worker_calendars(appointment_id);
CREATE INDEX IF NOT EXISTS idx_worker_calendars_worker_date ON worker_calendars(worker_id, calendar_date);

-- 添加注释 - construction_appointments
COMMENT ON TABLE construction_appointments IS '施工预约单表';
COMMENT ON COLUMN construction_appointments.id IS '主键ID';
COMMENT ON COLUMN construction_appointments.appointment_no IS '预约单号';
COMMENT ON COLUMN construction_appointments.project_id IS '关联项目ID';
COMMENT ON COLUMN construction_appointments.applicant_id IS '预约人ID';
COMMENT ON COLUMN construction_appointments.applicant_name IS '预约人姓名';
COMMENT ON COLUMN construction_appointments.contact_phone IS '联系电话';
COMMENT ON COLUMN construction_appointments.contact_person IS '联系人';
COMMENT ON COLUMN construction_appointments.work_date IS '作业日期';
COMMENT ON COLUMN construction_appointments.time_slot IS '时间段 (morning/afternoon/evening)';
COMMENT ON COLUMN construction_appointments.work_location IS '作业地点';
COMMENT ON COLUMN construction_appointments.work_content IS '作业内容';
COMMENT ON COLUMN construction_appointments.work_type IS '作业类型';
COMMENT ON COLUMN construction_appointments.is_urgent IS '是否加急';
COMMENT ON COLUMN construction_appointments.priority IS '优先级 (0-10)';
COMMENT ON COLUMN construction_appointments.urgent_reason IS '加急原因';
COMMENT ON COLUMN construction_appointments.assigned_worker_id IS '指派作业人员ID';
COMMENT ON COLUMN construction_appointments.assigned_worker_name IS '指派作业人员姓名';
COMMENT ON COLUMN construction_appointments.status IS '状态 (draft/pending/scheduled/in_progress/completed/cancelled/rejected)';
COMMENT ON COLUMN construction_appointments.workflow_instance_id IS '关联工作流实例ID';
COMMENT ON COLUMN construction_appointments.submitted_at IS '提交时间';
COMMENT ON COLUMN construction_appointments.approved_at IS '审批通过时间';
COMMENT ON COLUMN construction_appointments.completed_at IS '完成时间';

-- 添加注释 - worker_calendars
COMMENT ON TABLE worker_calendars IS '作业人员日历表';
COMMENT ON COLUMN worker_calendars.id IS '主键ID';
COMMENT ON COLUMN worker_calendars.worker_id IS '作业人员ID';
COMMENT ON COLUMN worker_calendars.calendar_date IS '日历日期';
COMMENT ON COLUMN worker_calendars.time_slot IS '时间段';
COMMENT ON COLUMN worker_calendars.is_available IS '是否可用';
COMMENT ON COLUMN worker_calendars.status IS '状态';
COMMENT ON COLUMN worker_calendars.appointment_id IS '关联预约单ID';
COMMENT ON COLUMN worker_calendars.blocked_reason IS '不可用原因';

-- 添加约束
ALTER TABLE construction_appointments ADD CONSTRAINT chk_appointments_status
    CHECK (status IN ('draft', 'pending', 'scheduled', 'in_progress', 'completed', 'cancelled', 'rejected'));

ALTER TABLE construction_appointments ADD CONSTRAINT chk_appointments_time_slot
    CHECK (time_slot IN ('morning', 'afternoon', 'evening', 'full_day'));

ALTER TABLE construction_appointments ADD CONSTRAINT chk_appointments_priority
    CHECK (priority >= 0 AND priority <= 10);

ALTER TABLE worker_calendars ADD CONSTRAINT chk_worker_calendars_time_slot
    CHECK (time_slot IN ('morning', 'afternoon', 'evening', 'full_day'));

ALTER TABLE worker_calendars ADD CONSTRAINT chk_worker_calendars_status
    CHECK (status IN ('available', 'busy', 'blocked', 'off'));

-- 创建预约单号序列
CREATE SEQUENCE IF NOT EXISTS appointment_no_seq;
CREATE OR REPLACE FUNCTION generate_appointment_no() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.appointment_no IS NULL OR NEW.appointment_no = '' THEN
        NEW.appointment_no := 'APT' || TO_CHAR(CURRENT_DATE, 'YYYYMMDD') || LPAD(nextval('appointment_no_seq')::TEXT, 4, '0');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_generate_appointment_no
    BEFORE INSERT ON construction_appointments
    FOR EACH ROW
    EXECUTE FUNCTION generate_appointment_no();

-- 自动更新 updated_at 触发器
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_construction_appointments_updated_at
    BEFORE UPDATE ON construction_appointments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

SELECT '011_create_appointment_tables.sql: construction_appointments and worker_calendars tables created' AS status;
