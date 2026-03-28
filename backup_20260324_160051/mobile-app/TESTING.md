# 移动端应用测试指南

## 编译前测试（推荐先在浏览器测试）

### 1. 启动开发服务器

```bash
npm run dev
```

这将在 `http://localhost:5173/mobile/` 启动开发服务器

### 2. 浏览器测试

**Chrome 桌面浏览器测试：**
```bash
# 启动开发服务器
npm run dev

# 打开浏览器访问
http://localhost:5173/mobile/
```

**Chrome 移动设备模拟：**
1. 按 `F12` 打开开发者工具
2. 点击工具栏上的"切换设备工具栏"图标（或按 `Ctrl+Shift+M`）
3. 选择设备类型（如 iPhone 14 Pro, Samsung Galaxy 等）
4. 测试触摸交互

**真机测试（通过局域网）：**
```bash
# 1. 启动开发服务器并监听所有网络接口
npm run dev -- --host

# 2. 查看你的局域网 IP
ip addr show  # Linux
ifconfig      # macOS

# 3. 在手机浏览器访问
http://YOUR_IP:5173/mobile/
```

### 3. Vue DevTools 调试

安装 Vue DevTools 浏览器扩展：
- Chrome: https://chrome.google.com/webstore/detail/vuejs-devtools
- Firefox: https://addons.mozilla.org/firefox/addon/vue-js-devtools/

### 4. 常见问题排查

**白屏问题：**
1. 打开浏览器控制台查看错误
2. 检查 Network 标签，确认 API 请求是否成功
3. 查看 Console 标签是否有 JavaScript 错误

**API 请求失败：**
- 确认后端服务运行在 `https://home.mbed.org.cn:9090`
- 开发环境会自动代理 `/api` 到后端
- 检查 CORS 设置

**路由问题：**
- 确认 URL 以 `/mobile/` 开头（开发环境）
- 确认已登录（未登录会跳转到 `/login`）

## Capacitor 本地测试（无需构建 APK）

### 1. 安装依赖

```bash
npm install
```

### 2. 同步到 Capacitor

```bash
npm run build
npx cap sync android
```

### 3. 在 Android Studio 中运行

```bash
npm run cap:android
```

这会打开 Android Studio，然后：
1. 连接 Android 设备或启动模拟器
2. 点击 Android Studio 中的"运行"按钮
3. 应用会实时同步，无需重新构建 APK

### 4. 使用 Capacitor 实时重载

```bash
# 终端 1: 启动 Vite 开发服务器
npm run dev

# 终端 2: 启动 Capacitor 实时同步
npx cap sync android && npx cap run android
```

代码修改后自动刷新应用！

## Chrome 远程调试（真机调试）

### 1. 启用 USB 调试

在 Android 设备上：
1. 设置 → 关于手机 → 连续点击"版本号" 7 次
2. 返回设置 → 系统 → 开发者选项
3. 启用"USB 调试"

### 2. 连接设备

```bash
# 检查设备连接
adb devices

# 安装 APK
adb install app-debug.apk

# 启动 Chrome 远程调试
chrome://inspect
```

### 3. Chrome 调试

1. 在 Chrome 浏览器访问 `chrome://inspect`
2. 在"目标"部分找到你的应用
3. 点击"inspect"打开 DevTools
4. 可以像在浏览器一样调试

## 日志调试

### 使用 adb 查看日志

```bash
# 查看实时日志
adb logcat | grep "Console"

# 查看所有日志
adb logcat

# 清除日志
adb logcat -c
```

### 在代码中添加日志

```javascript
console.log('调试信息:', data)
console.error('错误信息:', error)
console.warn('警告信息:', warning)
```

## 测试清单

### 功能测试
- [ ] 登录功能
- [ ] 物料计划列表和详情
- [ ] 入库管理
- [ ] 领料管理
- [ ] 个人中心
- [ ] 表单提交
- [ ] 图片上传

### UI/UX 测试
- [ ] 页面切换流畅
- [ ] 加载状态显示
- [ ] 错误提示清晰
- [ ] 表单验证
- [ ] 网络错误处理

### 兼容性测试
- [ ] 不同屏幕尺寸
- [ ] Android 不同版本
- [ ] 横竖屏切换
- [ ] 网络状况（4G/5G/WiFi）

## 性能优化建议

1. **减少初始加载体积**
   ```bash
   npm run build
   # 检查 dist/ 目录大小
   ```

2. **启用生产模式测试**
   ```bash
   npm run preview
   ```

3. **分析包大小**
   ```bash
   npx vite-bundle-visualizer
   ```

## 快速命令参考

```bash
# 开发环境测试
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview

# 同步到 Android
npx cap sync android

# 构建 APK
npm run android:build

# 实时运行（开发模式）
npm run dev & npx cap run android
```
