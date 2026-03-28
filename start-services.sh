#!/bin/bash
# 启动所有服务脚本

cd /home/julei/backend

echo "========================================"
echo "启动服务..."
echo "========================================"
echo ""
echo "注意: 所有配置已迁移到 config.yaml 文件"
echo "请确保 config.yaml 中已配置:"
echo "  - ai.deepseek_api_key (AI 聊天)"
echo "  - ai.asr_enabled (本地语音识别)"
echo ""

# 1. 启动本地 ASR 服务（如果启用）
if grep -q "asr_enabled: true" config.yaml 2>/dev/null; then
    echo "[1/2] 启动本地语音识别服务..."
    if pgrep -f "asr_service.py" > /dev/null; then
        echo "ASR 服务已在运行中"
    else
        nohup python3 scripts/asr_service.py > /tmp/asr_service.log 2>&1 &
        echo "ASR 服务启动中，PID: $!"
        sleep 3
    fi

    # 检查 ASR 服务状态
    if curl --noproxy '*' -s http://localhost:8089/health > /dev/null 2>&1; then
        echo "✓ ASR 服务运行正常 (http://localhost:8089)"
    else
        echo "⚠ ASR 服务可能需要更长时间加载模型..."
    fi
else
    echo "[1/2] 本地 ASR 服务未启用，跳过"
fi

# 2. 启动后端服务
echo ""
echo "[2/2] 启动后端服务..."
make restart

echo ""
echo "========================================"
echo "服务状态"
echo "========================================"
echo "后端 API: http://localhost:8088"
echo "ASR 服务: http://localhost:8089 (如果启用)"
echo ""
echo "配置文件: config.yaml"
echo ""
echo "日志位置:"
echo "  后端服务:  /tmp/backend.log"
echo "  ASR 服务:  /tmp/asr_service.log"
echo ""
echo "查看日志命令:"
echo "  tail -f /tmp/backend.log"
echo "  tail -f /tmp/asr_service.log"
echo ""
echo "配置 AI 功能:"
echo "  编辑 config.yaml 文件中的 ai 配置项"
