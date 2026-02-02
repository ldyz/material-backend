#!/bin/bash
# 使用Go后端的数据库配置查询
echo "检查库存数据完整性..."
echo ""

# 检查数据库连接
if [ -f .env ]; then
    source .env
fi

echo "请手动运行以下SQL查询来诊断问题："
echo ""
echo "1. 查看没有关联material的库存："
echo "   SELECT s.id, s.material_id, s.quantity, s.unit"
echo "   FROM stocks s"
echo "   LEFT JOIN materials m ON s.material_id = m.id"
echo "   WHERE m.id IS NULL"
echo "   LIMIT 10;"
echo ""
echo "2. 查看material_id为107的记录是否存在："
echo "   SELECT * FROM materials WHERE id = 107;"
echo ""
echo "3. 统计数据完整性："
echo "   SELECT"
echo "     COUNT(*) as total_stocks,"
echo "     COUNT(m.name) as stocks_with_material_name"
echo "   FROM stocks s"
echo "   LEFT JOIN materials m ON s.material_id = m.id;"
