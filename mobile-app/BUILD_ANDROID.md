# Android App 打包指南

## 环境要求

1. **Java JDK** - 已安装 OpenJDK 21 ✓
2. **Android SDK** - 需要安装

## 安装 Android SDK

### 方式 1: 使用 Android Studio (推荐)

1. 下载 Android Studio: https://developer.android.com/studio
2. 安装后打开 Android Studio
3. 进入 Settings -> Appearance & Behavior -> System Settings -> Android SDK
4. 安装 Android SDK (建议 API Level 33+)
5. 安装 Android SDK Build-Tools
6. 记录 SDK 路径，设置环境变量:
   ```bash
   export ANDROID_HOME=/path/to/android/sdk
   export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/tools
   ```

### 方式 2: 命令行安装 (Linux)

```bash
# 安装 Android SDK
wget https://dl.google.com/android/repository/commandlinetools-linux-9477386_latest.zip
unzip commandlinetools-linux-9477386_latest.zip
mkdir -p ~/Android/sdk/cmdline-tools/latest
mv cmdline-tools/* ~/Android/sdk/cmdline-tools/latest/

# 设置环境变量
export ANDROID_HOME=~/Android/sdk
export PATH=$PATH:$ANDROID_HOME/cmdline-tools/latest/bin:$ANDROID_HOME/platform-tools

# 接受许可
yes | sdkmanager --licenses

# 安装必要组件
sdkmanager "platform-tools" "platforms;android-34" "build-tools;34.0.0"
```

## 构建 APK

### 快速构建

安装完 Android SDK 后，运行:

```bash
# 方式 1: 使用脚本
npm run build:android

# 方式 2: 手动执行
npm run build
npx cap sync android
cd android
./gradlew assembleDebug
```

### APK 文件位置

构建完成后，APK 文件位于:
```
android/app/build/outputs/apk/debug/app-debug.apk
```

### 安装到设备

```bash
# 连接 Android 设备或启动模拟器
adb install android/app/build/outputs/apk/debug/app-debug.apk

# 或者使用 Capacitor 同步直接运行
npx cap run android
```

## 常见问题

### 1. Gradle 无法找到 Android SDK
检查 `local.properties` 文件:
```bash
cd android
echo "sdk.dir=$ANDROID_HOME" > local.properties
```

### 2. HTTP 网络请求失败
应用已配置使用 HTTPS API URL `https://home.mbed.org.cn:9090/api`

### 3. 修改 API 地址
编辑 `src/utils/request.js`:
```javascript
const baseURL = isCapacitor ? 'YOUR_API_URL' : '/api'
```

## 发布版本 (Release APK)

1. 生成签名密钥:
```bash
keytool -genkey -v -keystore release-key.keystore -alias release -keyalg RSA -keysize 2048 -validity 10000
```

2. 在 `capacitor.config.json` 中配置密钥路径

3. 构建发布版本:
```bash
cd android
./gradlew assembleRelease
```

4. 发布 APK 位置:
```
android/app/build/outputs/apk/release/app-release.apk
```
