-- 创建应用版本信息表
CREATE TABLE IF NOT EXISTS app_versions (
    id SERIAL PRIMARY KEY,
    platform VARCHAR(10) NOT NULL DEFAULT 'android', -- android, ios
    version VARCHAR(20) NOT NULL, -- 版本号，如 1.0.0
    build_number INTEGER DEFAULT 0, -- 构建号
    download_url TEXT, -- 下载链接
    force_update BOOLEAN DEFAULT FALSE, -- 是否强制更新
    update_message TEXT, -- 更新提示
    release_notes TEXT, -- 更新日志
    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 发布时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_platform_version UNIQUE (platform, version)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_app_versions_platform ON app_versions(platform);
CREATE INDEX IF NOT EXISTS idx_app_versions_published_at ON app_versions(published_at DESC);

-- 插入最新版本记录
INSERT INTO app_versions (platform, version, build_number, force_update, update_message, release_notes, published_at)
VALUES (
    'android',
    '1.0.26',
    26,
    FALSE,
    '优化作业类型显示，提升用户体验',
    '1. 预约单详情中作业类型现在显示中文标签
2. 优化了多项用户界面细节
3. 修复了已知问题并提升性能',
    NOW()
)
ON CONFLICT (platform, version) DO NOTHING;

-- 添加注释
COMMENT ON TABLE app_versions IS '移动应用版本信息表';
COMMENT ON COLUMN app_versions.platform IS '平台类型：android, ios';
COMMENT ON COLUMN app_versions.version IS '版本号，格式如 1.0.0';
COMMENT ON COLUMN app_versions.build_number IS '构建号';
COMMENT ON COLUMN app_versions.download_url IS 'APK下载链接';
COMMENT ON COLUMN app_versions.force_update IS '是否强制更新';
COMMENT ON COLUMN app_versions.update_message IS '简短的更新提示';
COMMENT ON COLUMN app_versions.release_notes IS '详细的更新日志';
