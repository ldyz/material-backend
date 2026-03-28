#!/bin/bash
# 编译后端并重启服务

set -e

cd /home/julei/backend

BINARY_NAME="server"
LOG_FILE="/tmp/backend.log"

echo "=========================================="
echo "  后端编译重启脚本"
echo "=========================================="

# 1. 停止现有进程
echo ""
echo "[1/3] 停止现有服务..."
pkill -f "./$BINARY_NAME" 2>/dev/null || true
killall -9 "$BINARY_NAME" 2>/dev/null || true
sleep 1
echo "      服务已停止"

# 2. 编译
echo ""
echo "[2/3] 编译后端..."
go build -o "$BINARY_NAME" ./cmd/server
echo "      编译完成: ./$BINARY_NAME"

# 3. 启动服务
echo ""
echo "[3/3] 启动服务..."
nohup ./$BINARY_NAME > "$LOG_FILE" 2>&1 &
sleep 2

# 检查是否启动成功
PID=$(pgrep -f "./$BINARY_NAME" | head -1)
if [ -n "$PID" ]; then
    echo "      服务已启动 (PID: $PID)"
else
    echo "      启动失败，请检查日志: $LOG_FILE"
    exit 1
fi

echo ""
echo "=========================================="
echo "  完成!"
echo "  日志文件: $LOG_FILE"
echo "  查看日志: tail -f $LOG_FILE"
echo "=========================================="
