#!/bin/bash
# 移动端应用自动构建发布脚本
# 用法: ./build-release.sh ["更新说明"]
# 示例: ./build-release.sh "修复登录问题"

set -e

PROJECT_DIR="/home/julei/backend/mobile-app"
UPDATE_DIR="/home/julei/backend/mobile-app-updates/android"
BASE_URL="https://home.mbed.org.cn:9090/mobile-updates/android"

cd "$PROJECT_DIR"

# 获取更新说明
RELEASE_NOTES="${1:-版本更新}"
echo "更新说明: $RELEASE_NOTES"

# 获取当前版本并计算新版本
CURRENT_VERSION=$(grep '"version"' package.json | cut -d'"' -f4)
MAJOR=$(echo $CURRENT_VERSION | cut -d'.' -f1)
MINOR=$(echo $CURRENT_VERSION | cut -d'.' -f2)
PATCH=$(echo $CURRENT_VERSION | cut -d'.' -f3)
NEW_PATCH=$((PATCH + 1))
NEW_VERSION="$MAJOR.$MINOR.$NEW_PATCH"
NEW_VERSION_CODE=$NEW_PATCH

echo "当前版本: $CURRENT_VERSION -> 新版本: $NEW_VERSION"

# 1. 更新 package.json
node -e "const fs=require('fs');const p=JSON.parse(fs.readFileSync('package.json','utf8'));p.version='$NEW_VERSION';fs.writeFileSync('package.json',JSON.stringify(p,null,2))"
echo "✓ 已更新 package.json"

# 2. 同步版本号
node sync-version.cjs
echo "✓ 已同步版本号"

# 3. 构建 Vue (Capacitor 模式)
CAPACITOR_BUILD=true npm run build 2>&1 | tail -3
echo "✓ Vue 构建完成"

# 4. 同步 Capacitor
npx cap sync android 2>&1 | tail -3
echo "✓ Capacitor 同步完成"

# 5. 构建 APK
cd android && ./gradlew assembleRelease -q 2>&1 | tail -5 && cd ..
echo "✓ Android APK 构建完成"

# 6. 复制 APK
APK_NAME="material-management-$NEW_VERSION.apk"
cp android/app/build/outputs/apk/release/app-release.apk "$UPDATE_DIR/$APK_NAME"
echo "✓ APK 已复制: $UPDATE_DIR/$APK_NAME"

# 7. 更新 version JSON
cat > "$UPDATE_DIR/version-$NEW_VERSION.json" << EOF
{
  "version": "$NEW_VERSION",
  "versionCode": $NEW_VERSION_CODE,
  "url": "/mobile-updates/android/$APK_NAME",
  "releaseNotes": "$RELEASE_NOTES",
  "forceUpdate": false,
  "publishDate": "$(date +%Y-%m-%d)"
}
EOF

rm -f "$UPDATE_DIR/latest.json"
ln -s "version-$NEW_VERSION.json" "$UPDATE_DIR/latest.json"
echo "✓ 版本 JSON 已更新"

# 8. 更新数据库
PGPASSWORD=julei1984 psql -h 127.0.0.1 -U materials -d materials -c "
INSERT INTO app_versions (platform, version, build_number, download_url, force_update, update_message, release_notes, published_at, created_at, updated_at)
VALUES ('android', '$NEW_VERSION', $NEW_VERSION_CODE, '$BASE_URL/$APK_NAME', false, '$RELEASE_NOTES', '$RELEASE_NOTES', NOW(), NOW(), NOW())
" > /dev/null 2>&1
echo "✓ 数据库记录已更新"

echo ""
echo "========================================"
echo "构建发布完成!"
echo "版本: $NEW_VERSION"
echo "下载: $BASE_URL/$APK_NAME"
echo "========================================"
