-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL, -- requisition_approve, inbound_approve, etc.
    title VARCHAR(200) NOT NULL,
    content TEXT,
    data JSONB, -- 存储相关数据，如 {requisition_id: 22, requisition_no: "CK20260126001"}
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);

-- 创建序列
CREATE SEQUENCE IF NOT EXISTS notifications_id_seq;
ALTER TABLE notifications ALTER COLUMN id SET DEFAULT nextval('notifications_id_seq'::regclass);
SELECT setval('notifications_id_seq', COALESCE((SELECT MAX(id) FROM notifications), 0) + 1, false);

-- 添加注释
COMMENT ON TABLE notifications IS '用户通知表';
COMMENT ON COLUMN notifications.type IS '通知类型：requisition_approve(出库审批), inbound_approve(入库审批)等';
COMMENT ON COLUMN notifications.data IS '相关数据JSON，如申请单ID、单号等';
