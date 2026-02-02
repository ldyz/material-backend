-- 添加系统配置表
-- 执行时间: 2026-01-28

-- 创建系统配置表
CREATE TABLE IF NOT EXISTS system_settings (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    description TEXT,
    category VARCHAR(50) DEFAULT 'general',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入默认配置
INSERT INTO system_settings (key, value, description, category) VALUES
    ('upload_directory', 'static/uploads', '文件上传目录路径', 'upload'),
    ('max_file_size', '5', '最大文件上传大小(MB)', 'upload'),
    ('allowed_file_types', 'jpg,jpeg,png,gif,bmp,webp,svg', '允许上传的文件类型', 'upload'),
    ('max_upload_count', '10', '单次最多上传文件数量', 'upload')
ON CONFLICT (key) DO NOTHING;

-- 添加注释
COMMENT ON TABLE system_settings IS '系统配置表';
COMMENT ON COLUMN system_settings.key IS '配置键';
COMMENT ON COLUMN system_settings.value IS '配置值';
COMMENT ON COLUMN system_settings.description IS '配置说明';
COMMENT ON COLUMN system_settings.category IS '配置分类';
