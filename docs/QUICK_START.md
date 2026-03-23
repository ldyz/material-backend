# 移动端自动更新 - 快速部署指南

## ✅ 已完成的步骤

### 1. 数据库初始化 ✓

```bash
# 已创建 app_versions 表
# 已插入初始版本数据（1.0.0）
```

### 2. 后端 API ✓

```bash
# 已实现版本检查接口
GET /api/app/version?platform=android&current_version=1.0.0

# 测试结果：
# - 当前版本 0.9.0 → 检测到更新 (1.0.0)
# - 当前版本 1.0.0 → 无更新
```

### 3. 前端集成 ✓

- 自动检测更新功能已集成到 `TabbarLayout.vue`
- 更新对话框组件已创建
- 每 24 小时自动检查一次

## 📱 下一步操作

### 方式 1：手动测试（推荐先做）

#### 测试流程：

1. **添加测试版本**
```bash
PGPASSWORD=julei1984 psql -h localhost -U materials -d materials -c "
INSERT INTO app_versions (platform, version, download_url, force_update, update_message)
VALUES ('android', '9.9.9', 'https://example.com/test.apk', false, '测试更新');
"
```

2. **重新构建移动端**
```bash
cd /home/julei/backend/mobile-app
npm run build
```

3. **在浏览器中打开应用**
```
http://localhost:8088/mobile/
```

4. **检查控制台**
   - 3秒后应该看到更新提示
   - 点击"立即更新"会跳转到下载链接

5. **清理测试数据**
```bash
PGPASSWORD=julei1984 psql -h localhost -U materials -d materials -c "
DELETE FROM app_versions WHERE version = '9.9.9';
"
```

### 方式 2：发布到 Android 设备

#### 第一步：准备环境

```bash
# 1. 安装 Android Studio
# 2. 安装 Capacitor CLI
cd /home/julei/backend/mobile-app
npm install @capacitor/cli
npm install @capacitor/android
```

#### 第二步：初始化 Android 项目

```bash
# 同步 Capacitor 配置
npx cap sync android

# 打开 Android Studio
npx cap open android
```

#### 第三步：构建 APK

在 Android Studio 中：
1. Build → Generate Signed Bundle / APK
2. 选择 APK
3. 选择签名密钥（或创建新的）
4. 构建 release 版本

#### 第四步：上传 APK

```bash
# 创建下载目录
sudo mkdir -p /var/www/html/downloads
sudo chown www-data:www-data /var/www/html/downloads

# 上传 APK
cp app-release.apk /var/www/html/downloads/material-management-1.0.0.apk

# 设置权限
sudo chmod 644 /var/www/html/downloads/material-management-1.0.0.apk
```

#### 第五步：更新数据库

```bash
PGPASSWORD=julei1984 psql -h localhost -U materials -d materials -c "
UPDATE app_versions
SET download_url = 'https://your-domain.com/downloads/material-management-1.0.0.apk'
WHERE platform = 'android' AND version = '1.0.0';
"
```

#### 第六步：安装到设备

1. 将 APK 下载到手机
2. 允许安装未知来源应用
3. 点击 APK 安装

## 🔄 发布新版本的完整流程

### 1. 修改版本号

```bash
cd /home/julei/backend/mobile-app
# 编辑 package.json
vim package.json
# 将 "version": "1.0.0" 改为 "1.0.1"
```

### 2. 构建新版本

```bash
npm run build
npx cap sync android
npx cap open android
# 在 Android Studio 中构建新的 APK
```

### 3. 上传并发布

```bash
# 使用发布脚本
./scripts/release_mobile_app.sh 1.0.1 android /path/to/app.apk

# 或手动执行
APK_FILE="material-management-1.0.1.apk"
cp $APK_FILE /var/www/html/downloads/

PGPASSWORD=julei1984 psql -h localhost -U materials -d materials -c "
INSERT INTO app_versions (platform, version, download_url, update_message, release_notes)
VALUES ('android', '1.0.1', 'https://your-domain.com/downloads/$APK_FILE',
  '修复已知问题', '修复了日期选择器的显示问题');
"
```

### 4. 用户收到更新

用户打开应用后会看到：
- 新版本提示（1.0.1）
- "立即更新"按钮
- 点击后下载新 APK
- 安装后覆盖旧版本

## 🧪 测试检查清单

- [ ] 数据库表已创建
- [ ] API 接口正常工作
- [ ] 前端构建成功
- [ ] 更新对话框正常显示
- [ ] 版本比较逻辑正确
- [ ] 下载链接可访问
- [ ] APK 可正常安装
- [ ] 强制更新功能正常
- [ ] 版本跳过功能正常

## 📊 当前状态

### 已配置的版本

| 平台 | 当前版本 | 下载地址 | 强制更新 |
|------|---------|---------|---------|
| Android | 1.0.0 | example.com/... | 否 |

### API 状态

```bash
# 测试命令
curl --noproxy "*" "http://localhost:8088/api/app/version?platform=android&current_version=0.9.0"

# 返回：有更新（1.0.0）
curl --noproxy "*" "http://localhost:8088/api/app/version?platform=android&current_version=1.0.0"

# 返回：无更新
```

### 前端状态

- ✅ 更新检测逻辑已实现
- ✅ 对话框组件已创建
- ✅ 已集成到主布局
- ✅ 自动检查功能已启用

## 🔧 配置修改

### 修改检查频率

编辑 `src/layouts/TabbarLayout.vue`:

```javascript
useAutoUpdate({
  autoCheck: true,
  checkOnMount: true,
  checkInterval: 12 * 60 * 60 * 1000  // 改为12小时
})
```

### 修改下载地址

在数据库中更新：
```sql
UPDATE app_versions
SET download_url = '新的下载地址'
WHERE platform = 'android' AND version = '1.0.0';
```

### 设置强制更新

```sql
UPDATE app_versions
SET force_update = true
WHERE version = '1.0.1';
```

## 📝 注意事项

1. **APK 签名**
   - 必须使用相同的签名密钥
   - 不同签名的 APK 无法覆盖安装

2. **版本号规则**
   - 格式：主版本.次版本.修订号（如 1.0.0）
   - 只能升级，不能降级
   - 数据库按版本号大小比较

3. **下载地址**
   - 必须是设备可访问的 URL
   - 建议使用 HTTPS
   - 确保 APK 文件存在

4. **强制更新**
   - 谨慎使用，会影响用户体验
   - 仅用于关键安全更新

## 🎯 快速测试

现在你可以：

1. **在浏览器测试**
   ```
   http://localhost:8088/mobile/
   ```
   打开浏览器控制台，3秒后应该看到更新提示

2. **测试 API**
   ```bash
   curl --noproxy "*" "http://localhost:8088/api/app/version?platform=android&current_version=0.0.1"
   ```

3. **查看数据库**
   ```bash
   PGPASSWORD=julei1984 psql -h localhost -U materials -d materials -c "SELECT * FROM app_versions;"
   ```

需要我帮你执行哪个步骤？
