# 施工预约功能 - 角色、权限和工作流说明

## 一、角色定义

### 1. 预约管理员 (角色ID: 10)
- **职责**: 创建和管理施工预约单
- **适用人员**: 客户、预约申请人员
- **权限**:
  - `appointment_view` - 查看预约单
  - `appointment_create` - 创建预约单
  - `appointment_edit` - 编辑预约单
  - `appointment_delete` - 删除预约单
  - `appointment_submit` - 提交预约单审批

### 2. 施工员 (角色ID: 11)
- **职责**: 第一级审批，确认可以承接作业
- **适用人员**: 施工现场负责人
- **权限**:
  - `appointment_view` - 查看预约单
  - `appointment_approve` - 审批预约单

### 3. 作业人员 (角色ID: 12)
- **职责**: 执行具体施工作业
- **适用人员**: 工人、作业执行人员
- **权限**:
  - `appointment_view` - 查看预约单
  - `appointment_execute` - 执行作业（开始、完成）

### 4. 项目经理 (角色ID: 4)
- **职责**: 最终审批确认
- **适用人员**: 项目负责人
- **权限**: 通过 `system_config` 权限审批

### 5. 经理 (角色ID: 2)
- **职责**: 加急作业审批
- **适用人员**: 部门经理
- **权限**: 通过 `system_config` 权限审批

## 二、权限列表

| 权限ID | 权限代码 | 权限名称 | 说明 |
|--------|----------|----------|------|
| 50 | `appointment_view` | 查看预约单 | 查看预约单列表和详情 |
| 51 | `appointment_create` | 创建预约单 | 创建新的施工预约单 |
| 52 | `appointment_edit` | 编辑预约单 | 编辑草稿状态的预约单 |
| 53 | `appointment_delete` | 删除预约单 | 删除预约单 |
| 54 | `appointment_submit` | 提交预约单审批 | 将预约单提交审批 |
| 55 | `appointment_approve` | 审批预约单 | 审批通过/拒绝预约单 |
| 56 | `appointment_assign` | 分配作业人员 | 为预约单分配作业人员 |
| 57 | `appointment_execute` | 执行作业 | 开始作业、完成作业 |
| 58 | `appointment_cancel` | 取消预约单 | 取消已提交的预约单 |
| 59 | `appointment_manage` | 管理预约单日历 | 管理作业人员日历 |

## 三、工作流定义

### 1. 普通预约工作流 (工作流ID: 10)
- **名称**: 普通施工预约审批流程
- **适用**: 非加急的常规预约
- **审批流程**:
  ```
  开始 → 施工员审批 → 项目经理审批 → 结束
  ```
- **节点说明**:
  - **施工员审批**: 确认可以承接该作业
  - **项目经理审批**: 最终确认，批准执行

### 2. 加急预约工作流 (工作流ID: 11)
- **名称**: 加急施工预约审批流程
- **适用**: 优先级>=7的加急预约
- **审批流程**:
  ```
  开始 → 施工员审批 → 经理审批 → 项目经理审批 → 结束
  ```
- **节点说明**:
  - **施工员审批**: 确认可以承接该作业
  - **经理审批**: 加急作业需要经理特别批准
  - **项目经理审批**: 最终确认，批准执行

## 四、状态流转

```
草稿 (draft)
  ↓ [提交]
待审批 (pending)
  ↓ [施工员通过]
  ↓ [经理通过 - 仅加急]
  ↓ [项目经理通过]
已排期 (scheduled)
  ↓ [开始作业]
进行中 (in_progress)
  ↓ [完成作业]
已完成 (completed)

其他状态:
- 已取消 (cancelled) - 申请人主动取消
- 已拒绝 (rejected) - 任一审批节点拒绝
```

## 五、使用步骤

### 1. 配置数据库
```bash
# 在PostgreSQL中执行
psql -U your_user -d your_database -f scripts/setup_appointment_permissions.sql
```

### 2. 分配用户角色
```sql
-- 查看用户
SELECT id, username, name FROM users;

-- 为用户分配角色（示例）
INSERT INTO user_roles (user_id, role_id, created_at) VALUES (2, 10, NOW());
```

### 3. 创建和提交预约
1. 以"预约管理员"角色登录
2. 创建预约单，填写作业信息
3. 选择作业类型（可多选）
4. 如需加急，勾选"加急"并设置优先级
5. 提交审批

### 4. 审批流程
1. **施工员**登录，审批待办任务
2. **项目经理**登录，进行最终确认
3. 如为加急，需**经理**先审批

### 5. 分配和执行
1. 审批通过后，系统状态变为"已排期"
2. 可为预约单分配作业人员
3. **作业人员**登录，查看分配的任务
4. 开始作业 → 完成作业

## 六、API端点

### 预约单管理
- `GET /api/appointments` - 获取预约单列表
- `GET /api/appointments/my` - 获取我的预约
- `GET /api/appointments/pending` - 获取待审批列表
- `POST /api/appointments` - 创建预约单
- `GET /api/appointments/:id` - 获取预约单详情
- `PUT /api/appointments/:id` - 更新预约单
- `DELETE /api/appointments/:id` - 删除预约单
- `POST /api/appointments/:id/submit` - 提交审批
- `POST /api/appointments/:id/approve` - 审批预约单
- `POST /api/appointments/:id/assign` - 分配作业人员
- `POST /api/appointments/:id/start` - 开始作业
- `POST /api/appointments/:id/complete` - 完成作业
- `POST /api/appointments/:id/cancel` - 取消预约

### 日历管理
- `GET /api/appointments/calendar/worker/:workerId` - 获取作业人员日历
- `POST /api/appointments/calendar/check-availability` - 检查可用性
- `GET /api/appointments/calendar/available-workers` - 获取可用作业人员

### 审批历史
- `GET /api/appointments/:id/approval-history` - 获取审批历史
- `GET /api/appointments/:id/workflow-progress` - 获取工作流进度
- `GET /api/appointments/:id/current-approval` - 获取当前审批节点

## 七、作业类型

支持以下作业类型（可多选）:
- 一般作业
- 动火作业
- 高处作业
- 动土作业
- 受限空间
- 临时用电
- 吊装作业
- 盲板抽堵

## 八、时间段

- morning: 上午 (08:00-12:00)
- afternoon: 下午 (14:00-18:00)
- evening: 晚上 (19:00-22:00)
- full_day: 全天

## 九、优先级

- 0-4: 普通优先级
- 5-6: 重要
- 7-10: 紧急（需要经理审批）
