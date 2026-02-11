-- 应用版本管理表
CREATE TABLE IF NOT EXISTS app_versions (
    id SERIAL PRIMARY KEY,
    platform VARCHAR(20) NOT NULL,           -- 平台: android, ios
    version VARCHAR(20) NOT NULL,            -- 版本号: 1.0.0
    build_number INTEGER DEFAULT 0,          -- 构建号
    download_url TEXT,                       -- 下载链接
    force_update BOOLEAN DEFAULT false,      -- 是否强制更新
    update_message TEXT,                     -- 更新提示消息
    release_notes TEXT,                      -- 更新日志
    published_at TIMESTAMP DEFAULT NOW(),    -- 发布时间
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(platform, version)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_app_versions_platform ON app_versions(platform);
CREATE INDEX IF NOT EXISTS idx_app_versions_published_at ON app_versions(published_at DESC);

-- 插入初始版本记录
INSERT INTO app_versions (platform, version, download_url, update_message, release_notes)
VALUES ('android', '1.0.1', 'https://home.mbed.org.cn:9090/mobile-updates/material-management-1.0.1.apk',
        '发现新版本 1.0.1，立即更新',
        '版本 1.0.1 更新：
- 施工预约管理功能
- 头像上传裁剪功能
- 多选作业人员功能
- UI优化和bug修复')
ON CONFLICT (platform, version) DO NOTHING;
