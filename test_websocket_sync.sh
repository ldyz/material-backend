#!/bin/bash

# 禁用代理
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY
export NO_PROXY="127.0.0.1,localhost"
export no_proxy="127.0.0.1,localhost"

echo "=================================="
echo "WebSocket 实时同步测试"
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
    "work_date": "2026-02-20",
    "time_slot": "morning",
    "work_type": "施工作业",
    "work_location": "WebSocket测试",
    "work_content": "测试WebSocket实时同步功能",
    "contact_person": "测试",
    "contact_phone": "13800138000"
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

# 4. 获取 julei token
echo "4. 获取审批人 token..."
JULEI_LOGIN=$(curl -s -X POST "http://127.0.0.1:8088/api/auth/login" -H "Content-Type: application/json" -d '{"username":"julei","password":"julei1984"}')
JULEI_TOKEN=$(echo $JULEI_LOGIN | jq -r '.meta.token // empty')

if [ -z "$JULEI_TOKEN" ]; then
  echo "错误: 无法获取 julei token"
  exit 1
fi
echo "julei token: ${JULEI_TOKEN:0:20}..."
echo ""

# 5. 执行审批（这应该触发WebSocket广播）
echo "5. 执行审批操作（应该触发WebSocket广播）..."
echo "=================================="

APPROVE_RESPONSE=$(curl -s -X POST "http://127.0.0.1:8088/api/appointments/$APPOINTMENT_ID/approve" \
  -H "Authorization: Bearer $JULEI_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action":"approve","comment":"WebSocket实时同步测试"}')

echo "审批响应:"
echo $APPROVE_RESPONSE | jq '.'
echo ""

# 6. 检查服务器日志
echo "6. 检查服务器日志中的WebSocket广播记录..."
sleep 1
echo "=================================="
tail -30 server.log | grep -i "broadcast\|websocket\|广播" || echo "未找到广播日志"

echo ""
echo "7. 检查审批历史..."
HISTORY_RESPONSE=$(curl -s -X GET "http://127.0.0.1:8088/api/appointments/$APPOINTMENT_ID/approval-history" \
  -H "Authorization: Bearer $WQS_TOKEN")

echo "审批历史记录:"
echo $HISTORY_RESPONSE | jq '.data | length'
echo $HISTORY_RESPONSE | jq '.data[] | {node_name, action, approver_name, created_at}'

echo ""
echo "=================================="
echo "测试完成"
echo "=================================="
