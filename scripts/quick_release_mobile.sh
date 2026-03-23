#!/bin/bash

# 移动端快速发布脚本（仅更新Web，不构建APK）
# 用法: ./scripts/quick_release_mobile.sh [version] [notes]

set -e

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
MOBILE_DIR="$PROJECT_DIR/mobile-app"
DB_NAME="materials"
DB_USER="materials"
DB_PASS="julei1984"
DB_HOST="localhost"
DOWNLOAD_BASE_URL="https://home.mbed.org.cn:9090/mobile-updates"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

function log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

function log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

function log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

function log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 获取当前版本号
CURRENT_VERSION=$(grep '"version"' "$MOBILE_DIR/package.json" | head -1 | sed 's/.*"version": "\(.*\)".*/\1/')

# 处理版本号参数
if [ -n "$1" ]; then
    NEW_VERSION="$1"
else
    # 自动递增版本号
    IFS='.' read -r major minor patch <<< "$CURRENT_VERSION"
    patch=$((patch + 1))
    NEW_VERSION="$major.$minor.$patch"
fi

log_info "当前版本: $CURRENT_VERSION"
log_info "新版本: $NEW_VERSION"

# 更新 package.json 中的版本号
log_step "更新 package.json 版本号..."
sed -i "s/\"version\": \"$CURRENT_VERSION\"/\"version\": \"$NEW_VERSION\"/" "$MOBILE_DIR/package.json"

# 获取更新说明
if [ -n "$2" ]; then
    RELEASE_NOTES="$2"
else
    RELEASE_NOTES="版本 $NEW_VERSION 更新：
- 功能优化
- bug修复
- 用户体验改进"
fi

log_info "更新说明: $RELEASE_NOTES"

# 进入移动端目录
cd "$MOBILE_DIR"

# 构建移动端
log_step "构建移动端 Web 静态资源..."
npm run build

log_info "Web 静态资源构建完成！"
log_info "移动端访问地址: https://home.mbed.org.cn:9090/mobile/"
log_info "用户刷新页面即可看到最新版本"

# 返回项目根目录
cd "$PROJECT_DIR"

# 构建 SQL 插入版本信息
SQL="INSERT INTO app_versions (platform, version, download_url, update_message, release_notes, force_update, published_at)
VALUES ('android', '$NEW_VERSION', '$DOWNLOAD_BASE_URL/latest.apk',
        '发现新版本 $NEW_VERSION，刷新页面即可更新',
        $(echo "$RELEASE_NOTES" | sed "s/'/''/g"),
        false,
        NOW())
ON CONFLICT (platform, version) DO UPDATE
SET download_url = '$DOWNLOAD_BASE_URL/latest.apk',
    update_message = '发现新版本 $NEW_VERSION，刷新页面即可更新',
    release_notes = $(echo "$RELEASE_NOTES" | sed "s/'/''/g"),
    published_at = NOW();"

# 显示 SQL
log_step "将执行以下 SQL:"
echo "$SQL"

# 询问是否继续
read -p "是否继续执行数据库更新？(y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log_warn "已取消数据库更新，但 Web 资源已构建完成"
    exit 0
fi

# 执行 SQL
log_step "执行数据库更新..."
PGPASSWORD="$DB_PASS" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "$SQL" 2>/dev/null || {
    log_warn "数据库更新失败（可能需要手动执行），但 Web 资源已构建完成"
    echo "$SQL" > /tmp/mobile_update_$NEW_VERSION.sql
    log_info "SQL 已保存到: /tmp/mobile_update_$NEW_VERSION.sql"
}

# 完成
echo ""
log_info "=========================================="
log_info "快速发布完成！"
log_info "=========================================="
log_info "版本: $NEW_VERSION"
log_info "类型: Web 更新（无需下载APK）"
log_info "访问地址: https://home.mbed.org.cn:9090/mobile/"
log_info "=========================================="

# 提示重启服务器
read -p "是否重启服务器以使更新生效？(y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    log_step "重启服务器..."
    pkill -f "./server" || true
    sleep 1
    nohup ./server > server.log 2>&1 &
    log_info "服务器已重启，PID: $!"
fi
