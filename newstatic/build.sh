#!/bin/bash
# 管理后台前端构建脚本

set -e

cd /home/julei/backend/newstatic

echo "=========================================="
echo "  管理后台前端构建脚本"
echo "=========================================="

echo ""
echo "[1/2] 构建前端..."
npm run build

echo ""
echo "[2/2] 更新版本号..."
node -e "
const fs = require('fs');
const pkgPath = './package.json';
if (fs.existsSync(pkgPath)) {
  const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
  if (pkg.version) {
    const parts = pkg.version.split('.');
    parts[2] = parseInt(parts[2]) + 1;
    pkg.version = parts.join('.');
    fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));
    console.log('Version updated to:', pkg.version);
  }
}
" 2>/dev/null || echo "跳过版本更新"

echo ""
echo "=========================================="
echo "  构建完成!"
echo "=========================================="
