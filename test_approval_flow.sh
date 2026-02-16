#!/bin/bash

# 审批流程测试脚本
BASE_URL="http://127.0.0.1:8088/api"
TOKEN=""  # 将在登录后填充

echo "========== 审批流程测试 =========="
echo ""

# 1. 登录获取 token
echo "1. 登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

echo "登录响应: $LOGIN_RESPONSE"
echo ""

# 尝试从 meta.token 或直接从 data.token 获取
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.meta.token // .token // empty' 2>/dev/null)

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ] || [ "$TOKEN" = "empty" ]; then
  # 尝试其他可能的字段
  TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token // .token // .access_token // empty' 2>/dev/null)
fi

echo "Token: $TOKEN"
echo ""

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ] || [ "$TOKEN" = "empty" ]; then
  echo "❌ 登录失败，无法获取 token"
  echo "请检查："
  echo "1. 后端服务是否在运行: ps aux | grep bin/server"
  echo "2. 尝试手动登录: curl -X POST http://127.0.0.1:8088/api/auth/login -H 'Content-Type: application/json' -d '{\"username\":\"admin\",\"password\":\"admin123\"}'"
  exit 1
fi

# 2. 获取预约单详情
echo "2. 获取测试预约单详情 (ID=12)..."
APPOINTMENT=$(curl -s -X GET "$BASE_URL/appointment/12" \
  -H "Authorization: Bearer $TOKEN")

echo "预约单信息:"
echo $APPOINTMENT | jq '.data | {id, appointment_no, status, workflow_instance_id}'
echo ""

# 3. 提交审批
echo "3. 提交审批..."
SUBMIT_RESPONSE=$(curl -s -X POST "$BASE_URL/appointment/12/submit" \
  -H "Authorization: Bearer $TOKEN")

echo "提交响应: $SUBMIT_RESPONSE"
echo ""

# 4. 查询审批历史（提交后）
echo "4. 查询审批历史（提交后）..."
sleep 1
HISTORY_AFTER_SUBMIT=$(curl -s -X GET "$BASE_URL/appointment/12/approval-history" \
  -H "Authorization: Bearer $TOKEN")

echo "审批历史记录数:"
echo $HISTORY_AFTER_SUBMIT | jq '.data | length'
echo "审批历史详情:"
echo $HISTORY_AFTER_SUBMIT | jq '.data'
echo ""

# 5. 获取当前审批节点
echo "5. 获取当前审批节点..."
CURRENT_APPROVAL=$(curl -s -X GET "$BASE_URL/appointment/12/current-approval" \
  -H "Authorization: Bearer $TOKEN")

echo "当前审批节点:"
echo $CURRENT_APPROVAL | jq '.data'
echo ""

# 6. 审批通过
echo "6. 审批通过..."
APPROVE_RESPONSE=$(curl -s -X POST "$BASE_URL/appointment/12/approve" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action":"approve","comment":"测试审批通过"}')

echo "审批响应: $APPROVE_RESPONSE"
echo ""

# 7. 查询审批历史（审批后）
echo "7. 查询审批历史（审批后）..."
sleep 1
HISTORY_AFTER_APPROVE=$(curl -s -X GET "$BASE_URL/appointment/12/approval-history" \
  -H "Authorization: Bearer $TOKEN")

echo "审批历史记录数:"
echo $HISTORY_AFTER_APPROVE | jq '.data | length'
echo "审批历史详情:"
echo $HISTORY_AFTER_APPROVE | jq '.data'
echo ""

# 8. 获取更新后的预约单状态
echo "8. 获取更新后的预约单状态..."
UPDATED_APPOINTMENT=$(curl -s -X GET "$BASE_URL/appointment/12" \
  -H "Authorization: Bearer $TOKEN")

echo "预约单当前状态:"
echo $UPDATED_APPOINTMENT | jq '.data | {id, appointment_no, status, approved_at}'
echo ""

echo "========== 测试完成 =========="
