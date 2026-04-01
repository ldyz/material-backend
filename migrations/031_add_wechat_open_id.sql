-- 添加微信 OpenID 字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS wechat_open_id VARCHAR(128);

-- 创建唯一索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_wechat_open_id ON users(wechat_open_id) WHERE wechat_open_id IS NOT NULL;
