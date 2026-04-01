#!/bin/bash

# 小程序构建脚本

echo "=== 开始构建微信小程序 ==="

# 进入项目目录
cd "$(dirname "$0")"

# 安装依赖
echo ">>> 安装依赖..."
npm install

# 构建
echo ">>> 构建小程序..."
npm run build:mp-weixin

echo "=== 构建完成 ==="
echo "构建产物在 dist/build/mp-weixin 目录下"
echo "请使用微信开发者工具打开该目录进行预览和上传"
