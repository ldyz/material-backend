-- 创建应用版本表
CREATE TABLE IF NOT EXISTS app_versions (
    id SERIAL PRIMARY KEY,
    platform VARCHAR(20) NOT NULL,           -- android, ios
    version VARCHAR(20) NOT NULL,            -- 版本号，如 1.0.0
    build_number INTEGER DEFAULT 0,          -- 构建号
    download_url TEXT,                       -- 下载链接
    force_update BOOLEAN DEFAULT false,      -- 是否强制更新
    update_message TEXT,                     -- 更新提示
    release_notes TEXT,                      -- 更新日志
    published_at TIMESTAMP DEFAULT NOW(),    -- 发布时间
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_app_versions_platform ON app_versions(platform);
CREATE INDEX IF NOT EXISTS idx_app_versions_published ON app_versions(published_at DESC);

-- 插入初始版本数据（请根据实际情况修改）
INSERT INTO app_versions (platform, version, build_number, download_url, force_update, update_message, release_notes)
VALUES
    ('android', '1.0.0', 1, 'https://example.com/downloads/material-management-1.0.0.apk', false,
     '初始版本',
     '首次发布材料管理系统移动端
- 完整的物资管理功能
- 预约管理功能
- 入库/出库管理
- 数据统计与报表')
ON CONFLICT DO NOTHING;

-- 可以为 iOS 添加单独的版本记录
-- INSERT INTO app_versions (platform, version, download_url, force_update, update_message, release_notes)
-- VALUES
--     ('ios', '1.0.0', 'https://apps.apple.com/app/xxxxx', false,
--      '初始版本',
--      '首次发布材料管理系统移动端');

COMMENT ON TABLE app_versions IS '应用版本信息表';
COMMENT ON COLUMN app_versions.platform IS '平台类型：android, ios';
COMMENT ON COLUMN app_versions.version IS '版本号，格式：主版本.次版本.修订号';
COMMENT ON COLUMN app_versions.build_number IS '构建号，用于比较同版本的不同构建';
COMMENT ON COLUMN app_versions.download_url IS '应用下载地址';
COMMENT ON COLUMN app_versions.force_update IS '是否强制更新（用户必须更新才能使用）';
COMMENT ON COLUMN app_versions.update_message IS '简短的更新提示，显示给用户';
COMMENT ON COLUMN app_versions.release_notes IS '详细的更新日志，新功能说明';
