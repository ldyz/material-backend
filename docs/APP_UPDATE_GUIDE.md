# 移动端自动升级功能说明

## 功能概述

移动端应用已集成自动检测更新功能，支持：
- 应用启动时自动检查更新
- 定期后台检查更新（每24小时）
- 手动检查更新
- 强制更新支持
- 版本跳过功能

## 使用前准备

### 1. 创建数据库表

执行 SQL 脚本创建版本管理表：

```bash
psql -U your_user -d your_database -f scripts/create_app_versions_table.sql
```

### 2. 配置版本信息

在数据库中插入应用版本信息：

```sql
-- 示例：添加新版本
INSERT INTO app_versions (platform, version, build_number, download_url, force_update, update_message, release_notes)
VALUES
    ('android', '1.0.1', 2, 'https://your-domain.com/downloads/material-management-1.0.1.apk', false,
     '修复已知问题，优化用户体验',
     '修复：
- 修复日期选择器显示问题
- 修复作业人员选择问题

优化：
- 优化应用性能
- 改进用户界面');
```

### 3. 上传 APK 文件

将编译好的 APK 文件上传到服务器可访问的位置：

```bash
# 创建下载目录
mkdir -p /path/to/downloads

# 上传 APK 文件
cp material-management-1.0.1.apk /path/to/downloads/

# 确保 Web 服务器可以访问这些文件
```

## 工作流程

### 1. 发布新版本

当需要发布新版本时：

1. **修改应用版本号**

   编辑 `mobile-app/package.json`：
   ```json
   {
     "version": "1.0.1"
   }
   ```

2. **构建应用**

   ```bash
   cd mobile-app
   npm run build
   npx cap sync android
   npx cap open android
   ```

   在 Android Studio 中生成签名的 APK

3. **上传 APK**

   将生成的 APK 上传到服务器

4. **添加版本记录**

   在数据库中插入新版本信息

### 2. 用户端更新流程

1. **自动检测**
   - 应用启动时自动检查
   - 每24小时后台检查一次
   - 发现新版本时显示更新提示

2. **用户操作**
   - 点击"立即更新"下载并安装新版本
   - 或选择"稍后提醒"（非强制更新）
   - 可以选择跳过该版本

3. **强制更新**
   - 如果设置了 `force_update = true`
   - 用户必须更新才能继续使用应用

## API 接口说明

### 检查版本接口

**请求：**
```
GET /api/app/version?platform=android&current_version=1.0.0
```

**参数：**
- `platform`: 平台类型（android/ios）
- `current_version`: 当前应用版本

**响应：**
```json
{
  "success": true,
  "data": {
    "has_update": true,
    "latest_version": "1.0.1",
    "download_url": "https://example.com/downloads/app-1.0.1.apk",
    "force_update": false,
    "update_message": "修复已知问题",
    "release_notes": "详细更新日志..."
  }
}
```

## 配置选项

### 前端配置

在 `src/layouts/TabbarLayout.vue` 中可以调整：

```javascript
const { performCheck } = useAutoUpdate({
  autoCheck: true,              // 是否自动检查
  checkOnMount: true,           // 启动时检查
  checkInterval: 24 * 60 * 60 * 1000  // 检查间隔（毫秒）
})
```

### 后端配置

在数据库中配置每个版本：

| 字段 | 说明 | 示例 |
|------|------|------|
| platform | 平台类型 | android, ios |
| version | 版本号 | 1.0.0 |
| build_number | 构建号 | 1 |
| download_url | 下载地址 | https://... |
| force_update | 是否强制更新 | true/false |
| update_message | 更新提示 | "修复了..." |
| release_notes | 更新日志 | 详细说明 |

## 手动检查更新

用户可以在"我的"页面手动检查更新（需要添加该功能）：

```javascript
// 在"我的"页面添加按钮
<van-button @click="checkUpdate">检查更新</van-button>

<script setup>
import { useAppUpdate } from '@/composables/useAppUpdate'

const { checkUpdate } = useAppUpdate()

async function checkUpdate() {
  const result = await checkUpdate()
  if (result.hasUpdate) {
    // 显示更新对话框
    showUpdateDialog.value = true
  } else {
    showToast('当前已是最新版本')
  }
}
</script>
```

## 版本号规则

版本号采用语义化版本规范：`主版本.次版本.修订号`

- **主版本**：重大功能变更或架构调整
- **次版本**：新增功能
- **修订号**：bug修复和小改进

示例：
- `1.0.0` → `1.0.1`：bug修复
- `1.0.1` → `1.1.0`：新增功能
- `1.1.0` → `2.0.0`：重大更新

## 安全建议

1. **APK签名**
   - 使用正式签名密钥打包
   - 不要将密钥提交到版本控制

2. **下载地址**
   - 使用 HTTPS 协议
   - 设置合理的访问权限

3. **强制更新**
   - 仅在关键安全更新时使用
   - 给用户足够的更新时间

## 常见问题

### Q: 如何测试更新功能？

A: 可以在数据库中插入一个更高的版本号，然后重新启动应用测试。

### Q: 更新下载失败怎么办？

A: 检查：
- download_url 是否正确
- 网络连接是否正常
- APK 文件是否存在

### Q: 如何回滚版本？

A: 保留旧版本的 APK，用户可以手动下载安装。

## 维护建议

1. **定期清理旧版本**
   ```sql
   -- 删除3个月前的旧版本记录
   DELETE FROM app_versions
   WHERE published_at < NOW() - INTERVAL '3 months';
   ```

2. **保留最近版本**
   ```sql
   -- 每个平台保留最近5个版本
   DELETE FROM app_versions
   WHERE id NOT IN (
     SELECT id FROM (
       SELECT id, ROW_NUMBER() OVER (
         PARTITION BY platform ORDER BY published_at DESC
       ) as rn
       FROM app_versions
     ) t WHERE rn <= 5
   );
   ```

3. **监控更新情况**
   - 记录更新请求日志
   - 统计用户更新率
