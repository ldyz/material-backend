# 移动端自动更新系统使用说明

## 概述

移动端自动更新系统支持两种更新方式：
1. **快速 Web 更新**：仅更新 Web 静态资源，用户刷新页面即可获得新版本
2. **完整 APK 更新**：构建新的 Android APK，用户需要下载并安装

## 文件结构

```
backend/
├── mobile-app/              # 移动端源代码
├── mobile-app-updates/      # 移动端更新文件存储目录
│   └── android/             # Android APK 和版本信息
├── scripts/
│   ├── quick_release_mobile.sh      # 快速 Web 更新脚本
│   ├── build_and_release_mobile.sh  # 完整 APK 构建脚本
│   └── setup_app_versions.sql       # 数据库表创建脚本
└── docs/
    └── MOBILE_AUTO_UPDATE.md        # 本文档
```

## 使用方法

### 方式一：快速 Web 更新（推荐）

适用于大多数更新，用户只需刷新页面即可获得新版本。

```bash
# 基本用法 - 自动递增版本号
./scripts/quick_release_mobile.sh

# 指定版本号
./scripts/quick_release_mobile.sh 1.0.3

# 指定版本号和更新说明
./scripts/quick_release_mobile.sh 1.0.3 "修复了几个bug：
- 头像上传功能优化
- 预约单提交问题修复
- UI优化"
```

脚本会自动：
1. 更新 `package.json` 中的版本号
2. 构建 Web 静态资源
3. 更新数据库中的版本记录
4. 提示是否重启服务器

### 方式二：完整 APK 构建

适用于需要更新原生代码或 Capacitor 插件的更新。

```bash
# 基本用法 - 自动递增版本号
./scripts/build_and_release_mobile.sh

# 指定版本号
./scripts/build_and_release_mobile.sh 1.1.0

# 指定版本号和更新说明
./scripts/build_and_release_mobile.sh 1.1.0 "重大功能更新：
- 新增离线功能
- 性能优化
- 全新UI设计"
```

脚本会自动：
1. 更新 `package.json` 中的版本号
2. 构建 Web 静态资源
3. 同步到 Capacitor Android
4. 构建 Debug APK
5. 复制 APK 到更新目录
6. 创建最新版本符号链接
7. 更新数据库记录
8. 提示是否重启服务器

## API 端点

### 版本检查 API

```
GET /api/app/version?platform=android&current_version=1.0.1
```

响应示例：
```json
{
  "success": true,
  "data": {
    "has_update": true,
    "latest_version": "1.0.2",
    "download_url": "https://home.mbed.org.cn:9090/mobile-updates/latest.apk",
    "force_update": false,
    "update_message": "发现新版本 1.0.2，立即更新",
    "release_notes": "版本 1.0.2 更新：\n- 功能优化\n- bug修复"
  }
}
```

### 文件下载

- Web 资源：`https://home.mbed.org.cn:9090/mobile/`
- APK 下载：`https://home.mbed.org.cn:9090/mobile-updates/latest.apk`
- 版本信息：`https://home.mbed.org.cn:9090/mobile-updates/latest.json`

## 移动端更新检测

移动端应用在启动时会自动检查版本更新。更新提示组件位于：
```
mobile-app/src/components/AppUpdateDialog.vue
```

## 版本管理

### 查看当前版本

```bash
# 查看 package.json 中的版本
grep '"version"' mobile-app/package.json

# 查看数据库中的版本记录
PGPASSWORD=julei1984 psql -h 127.0.0.1 -U materials -d materials \
  -c "SELECT * FROM app_versions ORDER BY published_at DESC LIMIT 5;"
```

### 手动添加版本记录

```sql
INSERT INTO app_versions (platform, version, download_url, update_message, release_notes, force_update)
VALUES ('android', '1.0.3',
    'https://home.mbed.org.cn:9090/mobile-updates/material-management-1.0.3.apk',
    '发现新版本 1.0.3，立即更新',
    '版本 1.0.3 更新：
- 新功能A
- 修复bugB
- 优化C',
    false)
ON CONFLICT (platform, version) DO UPDATE SET
    download_url = EXCLUDED.download_url,
    update_message = EXCLUDED.update_message,
    release_notes = EXCLUDED.release_notes,
    published_at = NOW();
```

## 常见问题

### 1. 快速更新后用户看不到新版本？

确保：
- 服务器已重启（静态文件缓存已刷新）
- 用户清除了浏览器缓存或强制刷新（Ctrl+Shift+R）
- CDN 缓存已刷新（如果使用 CDN）

### 2. 如何强制用户更新？

在数据库中将 `force_update` 设置为 `true`：
```sql
UPDATE app_versions SET force_update = true WHERE platform = 'android' AND version = '1.0.3';
```

### 3. 如何回滚到旧版本？

1. 恢复旧版本的 `mobile-app/dist/` 目录
2. 在数据库中将旧版本设置为最新：
```sql
UPDATE app_versions SET published_at = NOW()
WHERE platform = 'android' AND version = '1.0.1';
```

### 4. APK 构建失败？

确保：
- Android SDK 已正确安装
- `mobile-app/android` 目录存在且已初始化
- `gradlew` 文件有执行权限

## 注意事项

1. **版本号规则**：采用语义化版本号 (Semantic Versioning)，格式为 `主版本.次版本.补丁`
   - 主版本：不兼容的 API 修改
   - 次版本：向下兼容的功能性新增
   - 补丁：向下兼容的问题修正

2. **测试流程**：建议先在测试环境验证，再发布到生产环境

3. **备份**：发布前建议备份当前版本，以便必要时回滚

4. **数据库**：确保 `app_versions` 表已创建（运行 `scripts/setup_app_versions.sql`）

5. **权限**：脚本需要数据库写入权限和文件系统写入权限
