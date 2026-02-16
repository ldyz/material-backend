# 审批流程测试指南

## 准备工作

### 1. 启动后端服务
```bash
cd /home/julei/backend
./bin/server > server.log 2>&1 &
```

### 2. 启动前端服务（可选，用于Web测试）
```bash
cd /home/julei/backend/newstatic
npm run dev
```

### 3. 准备测试数据

测试预约单ID: **12**
- 状态: draft（草稿）
- 单号: TEST20260216001

## 测试步骤

### 方式一：使用移动端 APK 测试

1. **下载并安装最新版本**
   - 版本: v1.0.22
   - 下载地址: https://home.mbed.org.cn:9090/downloads/material-management-1.0.22.apk

2. **登录应用**
   - 用户名: admin
   - 密码: (需要确认密码，可能需要重置)

3. **导航到预约单详情**
   - 点击 "预约管理"
   - 找到测试预约单 (TEST20260216001)
   - 点击进入详情页

4. **测试审批流程**

   **步骤 A: 提交审批**
   - 确认状态为 "草稿"
   - 点击 "提交审批" 按钮
   - 观察以下变化：
     - ✓ Toast 提示 "提交成功"
     - ✓ 状态变为 "待审批"
     - ✓ 审批历史显示 "提交申请" 记录

   **步骤 B: 查看审批历史**
   - 滚动到 "审批历史" 部分
   - 应该看到时间线显示：
     ```
     ● 提交申请 ✓ 已通过
       admin - 申请人
       [时间戳]

     ⏰ 项目经理审批
       等待项目经理审批
     ```

   **步骤 C: 审批通过**
   - 点击 "审批通过" 按钮
   - 输入审批意见（可选）: "测试审批通过"
   - 点击 "确定"
   - 观察以下变化：
     - ✓ Toast 提示 "审批通过"
     - ✓ 对话框关闭
     - ✓ 审批历史立即更新
     - ✓ 时间线显示新的审批记录

   **步骤 D: 验证时间线更新**
   - 审批历史应显示：
     ```
     ● 提交申请 ✓ 已通过
       admin - 申请人
       [提交时间]

     ● 项目经理审批 ✓ 已通过
       admin - 项目经理
       [审批时间]
       备注：测试审批通过
     ```

### 方式二：使用浏览器控制台测试

1. **打开移动端应用**
   - 在浏览器中访问: http://192.168.18.1:5173 或 localhost
   - 登录应用

2. **打开开发者工具**
   - 按 F12 打开控制台
   - 切换到 Console 标签

3. **执行测试命令**

```javascript
// 1. 提交审批
fetch('/api/appointment/12/submit', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token')
  }
})
.then(r => r.json())
.then(d => console.log('提交结果:', d))

// 2. 查看审批历史
fetch('/api/appointment/12/approval-history', {
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token')
  }
})
.then(r => r.json())
.then(d => console.log('审批历史:', d))

// 3. 审批通过
fetch('/api/appointment/12/approve', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token'),
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    action: 'approve',
    comment: '测试审批通过'
  })
})
.then(r => r.json())
.then(d => console.log('审批结果:', d))

// 4. 再次查看审批历史（验证更新）
fetch('/api/appointment/12/approval-history', {
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token')
  }
})
.then(r => r.json())
.then(d => console.log('更新后的审批历史:', d))
```

### 方式三：使用 Postman/curl 测试

1. **登录获取 Token**
```bash
curl -X POST http://127.0.0.1:8088/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YOUR_PASSWORD"}'
```

2. **提交审批**
```bash
TOKEN="YOUR_TOKEN_HERE"
curl -X POST http://127.0.0.1:8088/api/appointment/12/submit \
  -H "Authorization: Bearer $TOKEN"
```

3. **查询审批历史**
```bash
curl -X GET http://127.0.0.1:8088/api/appointment/12/approval-history \
  -H "Authorization: Bearer $TOKEN"
```

4. **审批通过**
```bash
curl -X POST http://127.0.0.1:8088/api/appointment/12/approve \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action":"approve","comment":"测试审批通过"}'
```

## 预期结果

### 成功的标志

✅ **提交审批后:**
- 状态从 `draft` 变为 `pending`
- 创建 workflow_instance
- 添加一条 workflow_log (action: submit)
- 审批历史显示一条记录

✅ **审批通过后:**
- 状态从 `pending` 变为 `approved`
- 添加一条 workflow_approval (action: approve)
- 添加一条 workflow_log (action: approve)
- **审批历史立即更新**，显示新的审批记录
- 时间线显示正确的状态图标和颜色

### 失败排查

如果审批历史没有更新，检查：

1. **查看浏览器控制台日志**
   - 应该看到 `[审批历史]` 开头的日志
   - 记录原始响应、更新后的记录数等

2. **查看 Network 标签**
   - 检查 `/api/appointment/12/approval-history` 请求
   - 查看响应数据格式
   - 确认 `data` 字段是数组

3. **检查组件 key**
   - 组件的 key 应该包含 `approvalLogs.length`
   - 当记录数变化时，key 应该改变
   - Vue 会销毁旧组件并创建新组件

4. **检查数据库**
```sql
-- 查看 workflow_logs
SELECT * FROM workflow_logs WHERE instance_id = (
  SELECT workflow_instance_id FROM construction_appointments WHERE id = 12
) ORDER BY created_at;

-- 查看 workflow_approvals
SELECT * FROM workflow_approvals WHERE instance_id = (
  SELECT workflow_instance_id FROM construction_appointments WHERE id = 12
) ORDER BY created_at;
```

## 关键改进点

本次更新 (v1.0.22) 包含以下关键改进：

1. **组件 Key 强制重新渲染**
   ```vue
   :key="`timeline-${appointment.id}-${approvalLogs.length}`"
   ```

2. **增强的错误处理和日志**
   - 详细的控制台日志
   - Toast 提示错误信息

3. **响应式数据更新**
   - 使用 `computed` 构建时间线节点
   - 使用 `watch` 监听数据变化

## 联系方式

如有问题，请检查：
- 后端日志: `/home/julei/backend/server.log`
- 浏览器控制台日志
- API 响应数据
