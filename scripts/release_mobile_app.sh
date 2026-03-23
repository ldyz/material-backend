#!/bin/bash

# 移动端快速发布脚本
# 用法: ./scripts/release_mobile_app.sh <version> <platform> <apk_file>

set -e

VERSION=$1
PLATFORM=${2:-android}  # 默认为 android
APK_FILE=$3

# 配置
DB_NAME="material_management"
DB_USER="your_user"
DOWNLOAD_BASE_URL="https://your-domain.com/downloads"
DOWNLOAD_DIR="/path/to/downloads"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

function log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

function log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

function log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查参数
if [ -z "$VERSION" ]; then
    log_error "请提供版本号"
    echo "用法: $0 <version> [platform] [apk_file]"
    echo "示例: $0 1.0.1 android material-management-1.0.1.apk"
    exit 1
fi

log_info "准备发布版本 $VERSION for $PLATFORM"

# 检查 APK 文件
if [ "$PLATFORM" = "android" ] && [ -n "$APK_FILE" ]; then
    if [ ! -f "$APK_FILE" ]; then
        log_error "APK 文件不存在: $APK_FILE"
        exit 1
    fi

    # 上传 APK 文件
    log_info "上传 APK 文件..."
    cp "$APK_FILE" "$DOWNLOAD_DIR/"
    APK_FILENAME=$(basename "$APK_FILE")
    DOWNLOAD_URL="$DOWNLOAD_BASE_URL/$APK_FILENAME"
    log_info "下载地址: $DOWNLOAD_URL"
fi

# 构建 SQL 语句
if [ "$PLATFORM" = "android" ]; then
    SQL="
    INSERT INTO app_versions (platform, version, download_url, update_message, release_notes)
    VALUES ('android', '$VERSION', '$DOWNLOAD_URL',
        '版本 $VERSION 更新',
        '新版本发布：
- 性能优化
- bug修复
- 用户体验改进')
    ON CONFLICT (platform, version) DO UPDATE
    SET download_url = '$DOWNLOAD_URL',
        published_at = NOW();
    "
elif [ "$PLATFORM" = "ios" ]; then
    SQL="
    INSERT INTO app_versions (platform, version, download_url, update_message, release_notes)
    VALUES ('ios', '$VERSION', 'https://apps.apple.com/app/xxxxx',
        '版本 $VERSION 更新',
        '新版本发布，请前往 App Store 更新')
    ON CONFLICT (platform, version) DO UPDATE
    SET published_at = NOW();
    "
fi

# 显示 SQL
log_info "将执行以下 SQL:"
echo "$SQL"

# 询问是否继续
read -p "是否继续？(y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log_warn "已取消"
    exit 0
fi

# 执行 SQL
log_info "执行数据库更新..."
psql -U "$DB_USER" -d "$DB_NAME" -c "$SQL"

log_info "版本 $VERSION 发布成功！"
log_info "请确认："
echo "1. APK 文件已上传到服务器"
echo "2. download_url 配置正确: $DOWNLOAD_URL"
echo "3. 通知用户更新应用"

# 可选：发送通知给所有用户
read -p "是否发送更新通知给所有用户？(y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    log_info "发送更新通知..."
    # 这里可以调用通知 API
fi
