#!/bin/bash
# 移动端前端构建脚本

set -e

cd /home/julei/backend/mobile-app

echo "=========================================="
echo "  移动端前端构建脚本"
echo "=========================================="

echo ""
echo "[1/3] 构建前端 (Capacitor 模式)..."
rm -rf dist-capacitor/*
CAPACITOR_BUILD=true npm run build

echo ""
echo "[2/3] 同步到 dist 目录（用于 Web 访问）..."
rm -rf dist/*
cp -r dist-capacitor/* dist/

echo ""
echo "[3/3] 更新版本号..."
node -e "
const fs = require('fs');
const pkg = JSON.parse(fs.readFileSync('./package.json', 'utf8'));
const parts = pkg.version.split('.');
parts[2] = parseInt(parts[2]) + 1;
const newVersion = parts.join('.');

pkg.version = newVersion;
fs.writeFileSync('./package.json', JSON.stringify(pkg, null, 2));

const versionInfo = { version: newVersion, buildTime: new Date().toISOString() };
fs.writeFileSync('./dist-capacitor/version.json', JSON.stringify(versionInfo, null, 2));
fs.writeFileSync('./dist/version.json', JSON.stringify(versionInfo, null, 2));
fs.writeFileSync('./public/version.json', JSON.stringify(versionInfo, null, 2));

console.log('Version updated to:', newVersion);
"

echo ""
echo "=========================================="
echo "  构建完成!"
echo "=========================================="
