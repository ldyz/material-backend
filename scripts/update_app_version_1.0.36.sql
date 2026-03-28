-- 插入新版本 1.0.36
INSERT INTO app_versions (platform, version, build_number, download_url, force_update, update_message, release_notes, published_at)
VALUES ('android', '1.0.36', 36, 'https://home.mbed.org.cn:9090/mobile-updates/android/material-management-1.0.36.apk', false,
        '发现新版本 1.0.36，优化图标颜色和返回逻辑',
        '版本 1.0.36 更新：
- 修复预约管理图标颜色太浅的问题
- 修复施工预约日历返回需要两次的问题
- 修复项目选择器不显示的问题
- 优化首页图标颜色
- 移除系统公告栏
- 修复新建页面点击首页无效的问题
- 优化作业人员名字显示（最多4个）',
        NOW())
ON CONFLICT (platform, version) DO UPDATE SET
    build_number = 36,
    download_url = 'https://home.mbed.org.cn:9090/mobile-updates/android/material-management-1.0.36.apk',
    force_update = false,
    update_message = '发现新版本 1.0.36，优化图标颜色和返回逻辑',
    release_notes = '版本 1.0.36 更新：
- 修复预约管理图标颜色太浅的问题
- 修复施工预约日历返回需要两次的问题
- 修复项目选择器不显示的问题
- 优化首页图标颜色
- 移除系统公告栏
- 修复新建页面点击首页无效的问题
- 优化作业人员名字显示（最多4个）',
    published_at = NOW(),
    updated_at = NOW();
