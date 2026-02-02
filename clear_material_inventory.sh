#!/bin/bash

# =============================================
# 一键清空物资管理和库存管理数据
# =============================================

# 数据库配置
DB_HOST="127.0.0.1"
DB_PORT="5432"
DB_USER="materials"
DB_PASS="julei1984"
DB_NAME="materials"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}============================================${NC}"
echo -e "${YELLOW}  物资管理和库存管理数据清空工具${NC}"
echo -e "${YELLOW}============================================${NC}"
echo ""

# 确认操作
echo -e "${RED}警告：此操作将清空以下数据：${NC}"
echo "  - 物资计划 (material_plans, material_plan_items)"
echo "  - 物资主数据 (materials, material_categories)"
echo "  - 库存数据 (stocks, stock_logs, stock_op_logs)"
echo "  - 领料单 (requisitions, requisition_items)"
echo "  - 入库单 (inbound_orders, inbound_order_items)"
echo "  - 备份入库单 (inbound_orders2, inbound_order_items2)"
echo ""

read -p "确认要清空所有数据吗？(输入 YES 确认): " confirm

if [ "$confirm" != "YES" ]; then
    echo -e "${YELLOW}操作已取消${NC}"
    exit 0
fi

echo ""
echo -e "${GREEN}开始清空数据...${NC}"
echo ""

# 执行清空 SQL
PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME << 'EOF'
BEGIN;

-- 1. 库存管理相关表
DELETE FROM stock_op_logs;
DELETE FROM stock_logs;
DELETE FROM stocks;

-- 2. 领料单相关表
DELETE FROM requisition_items;
DELETE FROM requisitions;

-- 3. 入库单相关表
DELETE FROM inbound_order_items;
DELETE FROM inbound_orders;
DELETE FROM inbound_order_items2;
DELETE FROM inbound_orders2;

-- 4. 物资计划相关表
DELETE FROM material_plan_items;
DELETE FROM material_plans;

-- 5. 物资主数据表
DELETE FROM material_categories;
DELETE FROM materials;

COMMIT;

-- 显示清空结果统计
SELECT 'inbound_orders' AS table_name, COUNT(*) AS remaining_count FROM inbound_orders
UNION ALL
SELECT 'inbound_order_items', COUNT(*) FROM inbound_order_items
UNION ALL
SELECT 'inbound_orders2', COUNT(*) FROM inbound_orders2
UNION ALL
SELECT 'inbound_order_items2', COUNT(*) FROM inbound_order_items2
UNION ALL
SELECT 'material_plans', COUNT(*) FROM material_plans
UNION ALL
SELECT 'material_plan_items', COUNT(*) FROM material_plan_items
UNION ALL
SELECT 'materials', COUNT(*) FROM materials
UNION ALL
SELECT 'material_categories', COUNT(*) FROM material_categories
UNION ALL
SELECT 'stocks', COUNT(*) FROM stocks
UNION ALL
SELECT 'stock_logs', COUNT(*) FROM stock_logs
UNION ALL
SELECT 'stock_op_logs', COUNT(*) FROM stock_op_logs
UNION ALL
SELECT 'requisitions', COUNT(*) FROM requisitions
UNION ALL
SELECT 'requisition_items', COUNT(*) FROM requisition_items
ORDER BY table_name;

SELECT '============================================' AS info;
SELECT '物资管理和库存管理数据清空完成！' AS message;
SELECT '============================================' AS info;
EOF

# 检查执行结果
if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}============================================${NC}"
    echo -e "${GREEN}  数据清空完成！${NC}"
    echo -e "${GREEN}============================================${NC}"
else
    echo ""
    echo -e "${RED}数据清空失败，请检查错误信息${NC}"
    exit 1
fi
