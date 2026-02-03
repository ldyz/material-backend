#!/bin/bash
# 启动后端服务器

cd /home/julei/backend

# 停止现有进程
pkill -9 server 2>/dev/null

# 等待进程结束
sleep 1

# 启动服务器
./bin/server
