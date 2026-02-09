#!/bin/bash
# Android 构建环境配置脚本

echo "=== Android 构建环境配置 ==="

# 检查 ANDROID_HOME
if [ -z "$ANDROID_HOME" ]; then
    echo "⚠️  ANDROID_HOME 环境变量未设置"
    echo ""
    echo "请设置 Android SDK 路径："
    echo "  export ANDROID_HOME=/path/to/android/sdk"
    echo "  export PATH=\$PATH:\$ANDROID_HOME/platform-tools:\$ANDROID_HOME/tools"
    echo ""
    echo "常见位置："
    echo "  - Linux: ~/Android/sdk"
    echo "  - macOS: ~/Library/Android/sdk"
    echo "  - Windows: C:\\Users\\YourName\\AppData\\Local\\Android\\Sdk"
    echo ""
    read -p "请输入 Android SDK 路径: " SDK_PATH

    if [ ! -d "$SDK_PATH" ]; then
        echo "❌ 路径不存在: $SDK_PATH"
        exit 1
    fi

    export ANDROID_HOME="$SDK_PATH"
    echo "export ANDROID_HOME=$SDK_PATH" >> ~/.bashrc
    echo "export PATH=\$PATH:\$ANDROID_HOME/platform-tools:\$ANDROID_HOME/tools" >> ~/.bashrc
    echo "✓ 已添加到 ~/.bashrc"
fi

# 创建 local.properties
echo "sdk.dir=$ANDROID_HOME" > android/local.properties
echo "✓ 已创建 android/local.properties"

# 检查必要工具
echo ""
echo "=== 检查工具 ==="

if [ -f "$ANDROID_HOME/platform-tools/adb" ]; then
    echo "✓ adb 已安装"
else
    echo "❌ adb 未找到，请安装 Android SDK Platform-Tools"
fi

if [ -d "$ANDROID_HOME/build-tools" ]; then
    LATEST_BUILD_TOOLS=$(ls "$ANDROID_HOME/build-tools" | sort -V | tail -n 1)
    echo "✓ Build Tools: $LATEST_BUILD_TOOLS"
else
    echo "⚠️  Build Tools 未找到"
fi

echo ""
echo "=== 配置完成 ==="
echo "现在可以运行: npm run android:build"
