#!/bin/bash
# 快速测试脚本 - 在编译前测试应用

echo "=== 移动应用快速测试 ==="
echo ""
echo "正在启动开发服务器..."
echo ""

# 检查端口是否被占用
if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo "⚠️  端口 5173 已被占用"
    echo "是否要关闭旧进程并重新启动? (y/n)"
    read -r answer
    if [ "$answer" = "y" ]; then
        lsof -ti:5173 | xargs kill -9
        echo "✓ 已关闭旧进程"
    else
        echo "取消启动"
        exit 1
    fi
fi

# 启动开发服务器
npm run dev

# 如果需要局域网访问，使用：
# npm run dev -- --host
