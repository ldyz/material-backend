#!/usr/bin/env python3
"""
从入库单和出库单重建库存操作日志
"""

import psycopg2
from datetime import datetime

# PostgreSQL 连接配置
PG_CONFIG = {
    'host': '127.0.0.1',
    'port': 5432,
    'database': 'materials',
    'user': 'materials',
    'password': 'julei1984'
}

def rebuild_stock_logs():
    """从入库单和出库单重建库存操作日志"""
    print("========== 开始重建库存操作日志 ==========")

    pg_conn = psycopg2.connect(**PG_CONFIG)
    pg_cur = pg_conn.cursor()

    # 清空现有的 stock_logs
    pg_cur.execute("DELETE FROM stock_logs")
    pg_conn.commit()
    print("✓ 清空现有日志")

    count = 0

    # 1. 从已完成的入库单生成日志
    print("\n========== 处理入库单 ==========")
    pg_cur.execute("""
        SELECT
            io.id,
            ioi.material_id,
            ioi.quantity,
            io.order_no,
            io.supplier,
            io.status,
            io.created_at,
            s.id as stock_id,
            m.name as material_name
        FROM inbound_orders io
        JOIN inbound_order_items ioi ON io.id = ioi.order_id
        JOIN stocks s ON s.material_id = ioi.material_id
        JOIN materials m ON m.id = ioi.material_id
        WHERE io.status IN ('completed', 'approved')
        ORDER BY io.created_at ASC
    """)

    inbound_rows = pg_cur.fetchall()
    print(f"找到 {len(inbound_rows)} 条入库记录")

    for row in inbound_rows:
        (order_id, material_id, quantity, order_no, supplier,
         status, created_at, stock_id, material_name) = row
        supplier = supplier or '-'

        try:
            # stock_logs 表的字段: stock_id, type, quantity, time, remark
            pg_cur.execute('''
                INSERT INTO stock_logs
                (stock_id, type, quantity, time, remark, created_at)
                VALUES (%s, %s, %s, %s, %s, %s)
            ''', (
                stock_id,
                'in',
                quantity,
                created_at,
                f'入库单 {order_no} - {material_name}',
                created_at
            ))
            count += 1
        except Exception as e:
            print(f"  ✗ 入库单 {order_no} 失败: {e}")

    pg_conn.commit()
    print(f"✓ 导入 {count} 条入库日志")

    # 2. 从已发料的出库单生成日志
    print("\n========== 处理出库单 ==========")
    pg_cur.execute("""
        SELECT
            r.id,
            ri.stock_id,
            ri.quantity as quantity,
            r.requisition_no,
            r.purpose,
            r.status,
            r.created_at,
            m.name as material_name
        FROM requisitions r
        JOIN requisition_items ri ON r.id = ri.requisition_id
        JOIN stocks s ON s.id = ri.stock_id
        JOIN materials m ON m.id = s.material_id
        WHERE r.status IN ('issued', 'approved')
        ORDER BY r.created_at ASC
    """)

    outbound_rows = pg_cur.fetchall()
    print(f"找到 {len(outbound_rows)} 条出库记录")

    outbound_count = 0
    for row in outbound_rows:
        (requisition_id, stock_id, quantity, requisition_no,
         purpose, status, created_at, material_name) = row

        if quantity is None or quantity <= 0:
            continue

        try:
            pg_cur.execute('''
                INSERT INTO stock_logs
                (stock_id, type, quantity, time, remark, created_at, requisition_id)
                VALUES (%s, %s, %s, %s, %s, %s, %s)
            ''', (
                stock_id,
                'out',
                quantity,
                created_at,
                f'出库单 {requisition_no} - {material_name}',
                created_at,
                requisition_id
            ))
            outbound_count += 1
        except Exception as e:
            print(f"  ✗ 出库单 {requisition_no} 失败: {e}")

    pg_conn.commit()
    print(f"✓ 导入 {outbound_count} 条出库日志")

    # 验证结果
    pg_cur.execute("SELECT COUNT(*) FROM stock_logs")
    total = pg_cur.fetchone()[0]

    print(f"\n========== 重建完成 ==========")
    print(f"总共生成 {total} 条库存操作日志")

    pg_cur.close()
    pg_conn.close()

if __name__ == '__main__':
    rebuild_stock_logs()
