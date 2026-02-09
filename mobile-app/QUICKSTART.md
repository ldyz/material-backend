# 快速开始 - 解决白屏问题

## ✅ 问题已修复

已修复的原生应用白屏问题：
- ✓ 自动检测 Capacitor 环境
- ✓ 路由 base 路径自适应
- ✓ Vite 构建路径自适应
- ✓ API 地址自动配置
- ✓ 全局错误处理

## 🚀 编译前测试（推荐）

### 方式 1: 本地浏览器测试

```bash
# 启动开发服务器
npm run dev
# 或
npm test
# 或
./test.sh
```

然后访问: **http://localhost:5173/mobile/**

**优势：**
- 快速迭代，即时看到修改
- 浏览器 DevTools 调试
- Vue DevTools 支持
- 无需构建 APK

### 方式 2: 移动设备浏览器测试（局域网）

```bash
# 启动并监听所有网络接口
npm run dev:host
```

1. 查看你的 IP 地址：
   ```bash
   ip addr show    # Linux
   ifconfig        # macOS
   ```

2. 在手机浏览器访问：`http://YOUR_IP:5173/mobile/`

**优势：**
- 真实触摸交互
- 实际设备尺寸
- 无需安装 APK

### 方式 3: Chrome 设备模拟

1. `npm run dev` 启动开发服务器
2. 浏览器按 `F12` 打开 DevTools
3. 按 `Ctrl+Shift+M` 切换到设备模拟
4. 选择设备型号（iPhone 14 Pro, Samsung Galaxy 等）

## 📱 安装测试（APK）

### 快速安装

**新的 APK 文件** (已修复白屏):
```
/home/julei/backend/mobile-app/app-debug.apk
```

**方式 1: USB 安装**
```bash
# 连接设备后
adb install app-debug.apk

# 如果已安装，强制覆盖
adb install -r app-debug.apk
```

**方式 2: 直接传输**
将 `app-debug.apk` 发送到手机，直接打开安装

## 🔍 调试技巧

### 1. Chrome 远程调试（推荐）

```bash
# 1. 连接设备
adb devices

# 2. 安装并启动应用
adb install -r app-debug.apk

# 3. 在 Chrome 浏览器打开
chrome://inspect

# 4. 找到你的应用点击 "inspect"
```

### 2. 查看实时日志

```bash
# 查看 Console 日志
adb logcat | grep "Console"

# 查看所有日志
adb logcat

# 清除日志
adb logcat -c
```

### 3. 浏览器控制台

开发时按 `F12` 查看：
- **Console**: JavaScript 错误和日志
- **Network**: API 请求状态
- **Application**: LocalStorage 等存储

## 📋 修改代码后重新打包

### 只修改 Web 代码（Vue/CSS）

```bash
# 快速测试（浏览器）
npm run dev

# 重新构建 APK
npm run android:build
```

### 修改了原生配置

```bash
# 同步配置到 Android
npx cap sync android

# 如果需要重新构建
cd android && ./gradlew assembleDebug
```

## 🛠️ 常用命令

```bash
# 开发环境测试
npm run dev              # 本地测试 (localhost:5173/mobile)
npm run dev:host         # 局域网测试 (可通过手机访问)

# 构建
npm run build            # 构建生产版本
npm run preview          # 预览生产构建

# Android
npm run android:build    # 构建并生成 APK
npm run cap:android      # 在 Android Studio 中打开
npx cap sync android     # 同步到 Android 项目

# 测试
npm test                 # 启动开发服务器
./test.sh                # 使用测试脚本
```

## ⚠️ 常见问题

### Q1: 浏览器测试正常，但 APK 白屏

**A:** 已修复此问题！如果还遇到：
1. 确认使用最新构建的 APK
2. 使用 `adb logcat` 查看错误日志
3. 使用 Chrome 远程调试

### Q2: API 请求失败

**检查：**
- 后端服务是否运行在 `https://home.mbed.org.cn:9090`
- 网络连接是否正常
- 查看浏览器 Console 和 Network 标签

### Q3: 登录后跳转错误

**检查：**
- Token 是否正确保存
- API 返回的用户信息是否正确
- 查看 router/index.js:100 的 base 配置

## 📚 更多信息

- **详细测试指南**: 查看 `TESTING.md`
- **Android 构建指南**: 查看 `BUILD_ANDROID.md`
- **Capacitor 文档**: https://capacitorjs.com/

## ✨ 下一步

1. **先在浏览器测试**: `npm run dev`
2. **确保所有功能正常**: 登录、列表、详情等
3. **构建新 APK**: `npm run android:build`
4. **真机测试**: `adb install -r app-debug.apk`
5. **使用 Chrome 调试**: `chrome://inspect`

有问题？查看 `TESTING.md` 获取更多调试技巧！
