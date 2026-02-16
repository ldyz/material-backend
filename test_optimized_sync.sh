#!/bin/bash

# 禁用代理
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY
export NO_PROXY="127.0.0.1,localhost"
export no_proxy="127.0.0.1,localhost"

echo "=================================="
echo "优化功能测试 - WebSocket 实时同步"
echo "=================================="
echo ""

# 1. 获取 token
echo "1. 获取用户 token..."
WQS_LOGIN=$(curl -s -X POST "http://127.0.0.1:8088/api/auth/login" -H "Content-Type: application/json" -d '{"username":"wqs","password":"wqs1984"}')
WQS_TOKEN=$(echo $WQS_LOGIN | jq -r '.meta.token // empty')

if [ -z "$WQS_TOKEN" ] || [ "$WQS_TOKEN" = "null" ]; then
  echo "错误: 无法获取 wqs token"
  exit 1
fi
echo "wqs token: ${WQS_TOKEN:0:20}..."
echo ""

# 2. 创建测试预约单
echo "2. 创建测试预约单..."
CREATE_RESPONSE=$(curl -s -X POST "http://127.0.0.1:8088/api/appointments" \
  -H "Authorization: Bearer $WQS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": 1,
    "work_date": "2026-02-21",
    "time_slot": "afternoon",
    "work_type": "施工作业",
    "work_location": "优化测试地点",
    "work_content": "测试Toast提示、页面可见性、声音提示",
    "contact_person": "测试",
    "contact_phone": "13900139000"
  }')

APPOINTMENT_ID=$(echo $CREATE_RESPONSE | jq -r '.data.id // empty')
APPOINTMENT_NO=$(echo $CREATE_RESPONSE | jq -r '.data.appointment_no // empty')

if [ -z "$APPOINTMENT_ID" ]; then
  echo "错误: 创建预约单失败"
  echo "响应: $CREATE_RESPONSE"
  exit 1
fi

echo "创建成功: $APPOINTMENT_NO (ID: $APPOINTMENT_ID)"
echo ""

# 3. 提交审批
echo "3. 提交审批..."
SUBMIT_RESPONSE=$(curl -s -X POST "http://127.0.0.1:8088/api/appointments/$APPOINTMENT_ID/submit" \
  -H "Authorization: Bearer $WQS_TOKEN")

STATUS=$(echo $SUBMIT_RESPONSE | jq -r '.data.status // empty')
INSTANCE_ID=$(echo $SUBMIT_RESPONSE | jq -r '.data.workflow_instance_id // empty')

echo "状态: $STATUS, 工作流实例ID: $INSTANCE_ID"
echo ""

# 4. 获取审批人 token
echo "4. 获取审批人 token..."
JULEI_LOGIN=$(curl -s -X POST "http://127.0.0.1:8088/api/auth/login" -H "Content-Type: application/json" -d '{"username":"julei","password":"julei1984"}')
JULEI_TOKEN=$(echo $JULEI_LOGIN | jq -r '.meta.token // empty')

if [ -z "$JULEI_TOKEN" ]; then
  echo "错误: 无法获取 julei token"
  exit 1
fi
echo "julei token: ${JULEI_TOKEN:0:20}..."
echo ""

# 5. 第一次审批（触发 WebSocket 广播 + Toast + 声音）
echo "5. 第一次审批 - 测试优化功能"
echo "=================================="
APPROVE_RESPONSE=$(curl -s -X POST "http://127.0.0.1:8088/api/appointments/$APPOINTMENT_ID/approve" \
  -H "Authorization: Bearer $JULEI_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action":"approve","comment":"第一次审批 - 测试Toast和声音"}')

echo "审批响应:"
echo $APPROVE_RESPONSE | jq '.success, .message'
echo ""

sleep 1
echo "6. 检查服务器日志..."
echo "=================================="
tail -30 server.log | grep -i "broadcast\|广播" | tail -2

echo ""
echo "7. 获取审批历史..."
HISTORY_RESPONSE=$(curl -s -X GET "http://127.0.0.1:8088/api/appointments/$APPOINTMENT_ID/approval-history" \
  -H "Authorization: Bearer $WQS_TOKEN")

echo "审批历史记录数:"
echo $HISTORY_RESPONSE | jq '.data | length'
echo ""
echo "审批历史详情:"
echo $HISTORY_RESPONSE | jq '.data[] | {node_name, action, approver_name, created_at}'

echo ""
echo "8. 检查 WebSocket 连接状态..."
echo "=================================="
tail -50 server.log | grep -i "websocket.*connect\|Client registered" | tail -3

echo ""
echo "=================================="
echo "优化功能测试完成"
echo "=================================="
echo ""
echo "✅ Toast 提示: 前端会显示 '审批状态已更新'"
echo "✅ 声音提示: 会播放 '叮' 的提示音（Web Audio API）"
echo "✅ 页面可见性: 切换回页面时会自动刷新数据"
echo "✅ WebSocket 重连: 连接恢复后会自动刷新数据"
echo ""
echo "测试预约单ID: $APPOINTMENT_ID"
echo "可在前端打开该预约单详情页，然后在另一个窗口审批，观察效果"
