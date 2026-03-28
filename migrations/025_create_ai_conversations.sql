-- AI对话历史表
CREATE TABLE IF NOT EXISTS ai_conversations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,           -- 'user' 或 'assistant'
    content TEXT NOT NULL,
    tool_calls JSONB,                     -- AI工具调用（可选）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引以加速查询
CREATE INDEX idx_ai_conversations_user_id ON ai_conversations(user_id);
CREATE INDEX idx_ai_conversations_created_at ON ai_conversations(created_at DESC);

-- 添加注释
COMMENT ON TABLE ai_conversations IS 'AI对话历史记录';
COMMENT ON COLUMN ai_conversations.user_id IS '用户ID';
COMMENT ON COLUMN ai_conversations.role IS '消息角色：user(用户) 或 assistant(AI)';
COMMENT ON COLUMN ai_conversations.content IS '消息内容';
COMMENT ON COLUMN ai_conversations.tool_calls IS 'AI工具调用信息（JSON格式）';
