# 移动端自动更新功能

## 功能说明

已为移动端添加完整的自动更新功能，支持 Android 平台（iOS 可通过跳转 App Store 实现）。

## 主要特性

✅ **自动检测更新** - 应用启动时和定期后台检查
✅ **版本管理** - 数据库管理所有版本信息
✅ **强制更新** - 支持强制用户更新重要版本
✅ **版本跳过** - 用户可跳过非强制更新的版本
✅ **手动检查** - 支持用户手动检查更新
✅ **更新提示** - 友好的更新提示对话框

## 快速开始

### 1. 初始化数据库

```bash
cd /home/julei/backend
psql -U your_user -d your_database -f scripts/create_app_versions_table.sql
```

### 2. 构建移动端

```bash
cd mobile-app
npm run build
```

### 3. 发布新版本

```bash
# 方式1: 使用发布脚本
./scripts/release_mobile_app.sh 1.0.1 android /path/to/app.apk

# 方式2: 手动执行 SQL
psql -U your_user -d your_database
```

然后执行：
```sql
INSERT INTO app_versions (platform, version, download_url, update_message)
VALUES ('android', '1.0.1', 'https://example.com/downloads/app-1.0.1.apk', '新版本发布');
```

## 文件说明

### 前端文件

- `src/composables/useAppUpdate.js` - 更新检测逻辑
- `src/components/AppUpdateDialog.vue` - 更新提示对话框
- `src/layouts/TabbarLayout.vue` - 集成更新功能

### 后端文件

- `internal/api/app/model.go` - 版本数据模型
- `internal/api/app/handler.go` - 版本检查接口
- `internal/api/app/routes.go` - 路由注册

### 脚本文件

- `scripts/create_app_versions_table.sql` - 数据库表创建
- `scripts/release_mobile_app.sh` - 快速发布脚本

### 文档文件

- `docs/APP_UPDATE_GUIDE.md` - 详细使用指南

## API 接口

```
GET /api/app/version?platform=android&current_version=1.0.0
```

响应：
```json
{
  "success": true,
  "data": {
    "has_update": true,
    "latest_version": "1.0.1",
    "download_url": "https://example.com/app.apk",
    "force_update": false,
    "update_message": "新版本可用"
  }
}
```

## 配置说明

### 检查频率

在 `TabbarLayout.vue` 中配置：

```javascript
useAutoUpdate({
  autoCheck: true,                   // 启用自动检查
  checkOnMount: true,                // 启动时检查
  checkInterval: 24 * 60 * 60 * 1000 // 24小时
})
```

### 版本配置

数据库字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| platform | string | android 或 ios |
| version | string | 版本号，如 1.0.0 |
| download_url | string | APK 下载地址 |
| force_update | boolean | 是否强制更新 |
| update_message | string | 更新提示文本 |
| release_notes | text | 详细更新日志 |

## 测试方法

1. **添加测试版本**
   ```sql
   INSERT INTO app_versions (platform, version, download_url, force_update)
   VALUES ('android', '9.9.9', 'https://example.com/test.apk', false);
   ```

2. **重新构建应用**
   ```bash
   cd mobile-app
   npm run build
   ```

3. **启动应用**
   - 应该会自动检测到更新
   - 显示更新提示对话框

4. **清理测试数据**
   ```sql
   DELETE FROM app_versions WHERE version = '9.9.9';
   ```

## 发布流程

1. **修改版本号**
   编辑 `mobile-app/package.json` 中的 `version` 字段

2. **构建 APK**
   ```bash
   cd mobile-app
   npm run build
   npx cap sync android
   npx cap open android
   ```
   在 Android Studio 中生成签名 APK

3. **上传 APK**
   将 APK 上传到服务器下载目录

4. **发布版本**
   ```bash
   ./scripts/release_mobile_app.sh <version> android <apk_file>
   ```

5. **验证**
   - 用户打开应用会收到更新提示
   - 可以下载并安装新版本

## 注意事项

⚠️ **APK 签名**
- 必须使用相同的签名密钥
- 否则无法覆盖安装

⚠️ **版本号规则**
- 遵循语义化版本：主版本.次版本.修订号
- 只能升级，不能降级

⚠️ **强制更新**
- 仅在必要时使用
- 给用户足够的准备时间

⚠️ **网络配置**
- 确保 download_url 可访问
- 建议使用 HTTPS

## 故障排查

### 问题：检测不到更新

**检查：**
1. 数据库中是否有新版本记录
2. 版本号比较是否正确
3. API 接口是否正常

**调试：**
```javascript
// 在浏览器控制台
window.checkAppUpdate()
```

### 问题：下载失败

**检查：**
1. download_url 是否正确
2. APK 文件是否存在
3. 网络连接是否正常

### 问题：强制更新不生效

**检查：**
1. force_update 字段是否为 true
2. 对话框是否正确显示

## 后续扩展

可以考虑的功能：

1. **增量更新** - 仅下载变更部分
2. **后台下载** - 不阻塞用户操作
3. **WiFi 提示** - 提醒在 WiFi 下更新
4. **更新统计** - 记录更新率
5. **灰度发布** - 逐步推送更新
6. **版本回滚** - 支持回退到旧版本

## 支持

如有问题，请参考：
- [详细使用指南](./APP_UPDATE_GUIDE.md)
- Capacitor 官方文档
- Android 版本管理最佳实践
