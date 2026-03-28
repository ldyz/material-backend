-- AI意图缓存表（用于减少LLM请求）
CREATE TABLE IF NOT EXISTS ai_intent_cache (
    id BIGSERIAL PRIMARY KEY,
    intent_pattern VARCHAR(255) NOT NULL UNIQUE,  -- 意图模式（如"今日任务"、"库存预警"）
    tool_name VARCHAR(100) NOT NULL,              -- 工具名称
    tool_args JSONB NOT NULL,                     -- 工具参数（JSON格式）
    description TEXT,                              -- 描述
    hit_count INTEGER DEFAULT 0,                   -- 命中次数
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_ai_intent_cache_pattern ON ai_intent_cache(intent_pattern);
CREATE INDEX idx_ai_intent_cache_tool ON ai_intent_cache(tool_name);

-- 添加注释
COMMENT ON TABLE ai_intent_cache IS 'AI意图缓存表，存储常见查询对应的工具调用';
COMMENT ON COLUMN ai_intent_cache.intent_pattern IS '意图模式关键词';
COMMENT ON COLUMN ai_intent_cache.tool_name IS '对应的工具名称';
COMMENT ON COLUMN ai_intent_cache.tool_args IS '工具参数（可能包含动态值如today）';
COMMENT ON COLUMN ai_intent_cache.hit_count IS '缓存命中次数统计';

-- 插入常见意图缓存数据
INSERT INTO ai_intent_cache (intent_pattern, tool_name, tool_args, description) VALUES
-- 任务/预约相关
('今日任务', 'query_appointments', '{"date": "today"}', '查询今日任务'),
('今天任务', 'query_appointments', '{"date": "today"}', '查询今日任务'),
('今日安排', 'query_appointments', '{"date": "today"}', '查询今日安排'),
('今天安排', 'query_appointments', '{"date": "today"}', '查询今日安排'),
('今天的任务', 'query_appointments', '{"date": "today"}', '查询今日任务'),
('明天任务', 'query_appointments', '{"date": "tomorrow"}', '查询明日任务'),
('明日任务', 'query_appointments', '{"date": "tomorrow"}', '查询明日任务'),
('明天安排', 'query_appointments', '{"date": "tomorrow"}', '查询明日安排'),
('昨天的任务', 'query_appointments', '{"date": "yesterday"}', '查询昨日任务'),

-- 库存相关
('库存预警', 'query_stock_alerts', '{}', '查询库存预警'),
('低库存', 'query_stock_alerts', '{}', '查询库存预警'),
('库存不足', 'query_stock_alerts', '{}', '查询库存预警'),
('库存查询', 'query_stock', '{}', '查询库存列表'),
('查看库存', 'query_stock', '{}', '查询库存列表'),

-- 审批相关
('待审批', 'query_pending_approvals', '{}', '查询待审批事项'),
('待办任务', 'query_pending_approvals', '{}', '查询待审批事项'),
('待办事项', 'query_pending_approvals', '{}', '查询待审批事项'),
('需要审批', 'query_pending_approvals', '{}', '查询待审批事项'),
('我的审批', 'query_pending_approvals', '{}', '查询待审批事项'),

-- 考勤相关
('打卡记录', 'query_attendance', '{"date": "today"}', '查询今日打卡记录'),
('考勤记录', 'query_attendance', '{"date": "today"}', '查询今日考勤记录'),
('今天考勤', 'query_attendance', '{"date": "today"}', '查询今日考勤'),
('今日考勤', 'query_attendance', '{"date": "today"}', '查询今日考勤'),

-- 施工日志相关
('施工日志', 'query_construction_logs', '{"log_date": "today"}', '查询今日施工日志'),
('今天日志', 'query_construction_logs', '{"log_date": "today"}', '查询今日施工日志'),
('今日日志', 'query_construction_logs', '{"log_date": "today"}', '查询今日施工日志'),

-- 项目相关
('项目列表', 'query_projects', '{}', '查询项目列表'),
('所有项目', 'query_projects', '{}', '查询项目列表')
ON CONFLICT (intent_pattern) DO NOTHING;
