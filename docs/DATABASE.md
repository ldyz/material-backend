# 数据库文档

本文档说明了数据库的表结构、字段定义和迁移脚本。

## 数据库概述

- **数据库类型**: PostgreSQL 14+
- **ORM框架**: GORM
- **字符编码**: UTF-8
- **时区**: Asia/Shanghai (UTC+8)

## 迁移文件列表

| 文件 | 说明 |
|------|------|
| 001_create_auth_tables.sql | 用户认证相关表 |
| 002_create_projects_tables.sql | 项目管理表 |
| 003_create_user_projects.sql | 用户项目关联表 |
| 004_create_materials_table.sql | 物资主数据表 |
| 005_create_inbound_tables.sql | 入库管理表 |
| 006_create_stock_tables.sql | 库存表 |
| 007_create_stock_logs_table.sql | 库存日志表 |
| 008_create_system_backup_table.sql | 系统备份表 |
| 009_create_project_schedules_table.sql | 项目进度计划表 |
| 011_create_notifications_table.sql | 通知表 |
| 012_create_workflow_tables.sql | 工作流表 |
| 020_create_attendance_tables.sql | 考勤管理表 |
| 025_create_ai_conversations.sql | AI对话历史表 |

---

## 核心数据表

### users 用户表

存储系统用户信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| username | VARCHAR(80) | UNIQUE, NOT NULL | 用户名 |
| password | VARCHAR(255) | NOT NULL | 密码（bcrypt加密） |
| email | VARCHAR(255) | | 邮箱 |
| full_name | VARCHAR(100) | | 全名 |
| avatar | VARCHAR(500) | | 头像URL |
| role | VARCHAR(50) | | 角色代码（兼容字段） |
| group | VARCHAR(100) | | 用户组 |
| is_active | BOOLEAN | DEFAULT true | 是否激活 |
| last_login | TIMESTAMP | | 最后登录时间 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

**索引**:
- `idx_users_username` ON (username)

---

### roles 角色表

存储角色信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| name | VARCHAR(100) | UNIQUE, NOT NULL | 角色名称 |
| description | TEXT | | 角色描述 |
| permissions | TEXT | | 权限列表（逗号分隔） |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

---

### user_roles 用户角色关联表

多对多关系表。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| user_id | INTEGER | REFERENCES users(id) | 用户ID |
| role_id | INTEGER | REFERENCES roles(id) | 角色ID |

**主键**: (user_id, role_id)

---

### projects 项目表

存储项目信息，支持层级结构。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| name | TEXT | NOT NULL | 项目名称 |
| code | TEXT | UNIQUE | 项目编号 |
| location | TEXT | | 项目地点 |
| start_date | DATE | | 开始日期 |
| end_date | DATE | | 结束日期 |
| description | TEXT | | 项目描述 |
| manager | TEXT | | 项目经理 |
| contact | TEXT | | 联系方式 |
| budget | TEXT | | 预算 |
| status | TEXT | DEFAULT 'planning' | 状态 |
| parent_id | INTEGER | REFERENCES projects(id) | 父项目ID |
| level | INTEGER | DEFAULT 0 | 层级深度 |
| path | VARCHAR(500) | | 层级路径 |
| progress_percentage | DECIMAL(5,2) | DEFAULT 0 | 进度百分比 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**索引**:
- `idx_projects_parent_id` ON (parent_id)
- `idx_projects_status` ON (status)

**状态值**:
- `planning`: 规划中
- `active`: 进行中
- `closed`: 已关闭
- `on_hold`: 暂停

---

### user_projects 用户项目关联表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| user_id | INTEGER | REFERENCES users(id) | 用户ID |
| project_id | INTEGER | REFERENCES projects(id) | 项目ID |

**主键**: (user_id, project_id)

---

### materials 物资主数据表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| name | TEXT | NOT NULL | 物资名称 |
| code | TEXT | UNIQUE | 物资编码 |
| category_id | INTEGER | REFERENCES material_categories(id) | 分类ID |
| unit | VARCHAR(50) | | 单位 |
| specification | TEXT | | 规格型号 |
| brand | VARCHAR(100) | | 品牌 |
| price | DECIMAL(10,2) | | 单价 |
| description | TEXT | | 描述 |
| is_active | BOOLEAN | DEFAULT true | 是否启用 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

---

### material_categories 物资分类表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| name | TEXT | NOT NULL | 分类名称 |
| parent_id | INTEGER | REFERENCES material_categories(id) | 父分类ID |
| level | INTEGER | DEFAULT 0 | 层级 |
| sort_order | INTEGER | DEFAULT 0 | 排序 |

---

### material_plans 物资计划表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| plan_no | VARCHAR(50) | UNIQUE | 计划编号 |
| project_id | INTEGER | REFERENCES projects(id) | 项目ID |
| title | TEXT | NOT NULL | 计划标题 |
| description | TEXT | | 描述 |
| plan_date | DATE | | 计划日期 |
| required_date | DATE | | 需求日期 |
| status | VARCHAR(20) | DEFAULT 'draft' | 状态 |
| applicant_id | INTEGER | REFERENCES users(id) | 申请人ID |
| approved_by | INTEGER | REFERENCES users(id) | 审批人ID |
| approved_at | TIMESTAMP | | 审批时间 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

**状态值**:
- `draft`: 草稿
- `pending`: 待审批
- `approved`: 已审批
- `rejected`: 已拒绝
- `completed`: 已完成

---

### material_plan_items 计划明细表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| plan_id | INTEGER | REFERENCES material_plans(id) | 计划ID |
| material_id | INTEGER | REFERENCES materials(id) | 物资ID |
| quantity | DECIMAL(10,2) | NOT NULL | 数量 |
| unit | VARCHAR(50) | | 单位 |
| remark | TEXT | | 备注 |

---

### inbound_orders 入库单表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| order_no | VARCHAR(50) | UNIQUE | 入库单号 |
| project_id | INTEGER | REFERENCES projects(id) | 项目ID |
| supplier | VARCHAR(200) | | 供应商 |
| total_amount | DECIMAL(12,2) | | 总金额 |
| status | VARCHAR(20) | DEFAULT 'draft' | 状态 |
| applicant_id | INTEGER | REFERENCES users(id) | 申请人ID |
| approved_by | INTEGER | REFERENCES users(id) | 审批人ID |
| approved_at | TIMESTAMP | | 审批时间 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

**状态值**:
- `draft`: 草稿
- `pending`: 待审批
- `approved`: 已审批
- `rejected`: 已拒绝

---

### inbound_order_items 入库明细表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| order_id | INTEGER | REFERENCES inbound_orders(id) | 入库单ID |
| material_id | INTEGER | REFERENCES materials(id) | 物资ID |
| quantity | DECIMAL(10,2) | NOT NULL | 数量 |
| unit_price | DECIMAL(10,2) | | 单价 |
| amount | DECIMAL(12,2) | | 金额 |

---

### requisitions 出库单表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| requisition_no | VARCHAR(50) | UNIQUE | 出库单号 |
| project_id | INTEGER | REFERENCES projects(id) | 项目ID |
| recipient | VARCHAR(100) | | 领用人 |
| purpose | TEXT | | 用途 |
| status | VARCHAR(20) | DEFAULT 'draft' | 状态 |
| applicant_id | INTEGER | REFERENCES users(id) | 申请人ID |
| approved_by | INTEGER | REFERENCES users(id) | 审批人ID |
| issued_by | INTEGER | REFERENCES users(id) | 发货人ID |
| approved_at | TIMESTAMP | | 审批时间 |
| issued_at | TIMESTAMP | | 发货时间 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

---

### stocks 库存表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| material_id | INTEGER | REFERENCES materials(id) | 物资ID |
| project_id | INTEGER | REFERENCES projects(id) | 项目ID |
| quantity | DECIMAL(10,2) | DEFAULT 0 | 库存数量 |
| unit | VARCHAR(50) | | 单位 |
| location | VARCHAR(200) | | 存放位置 |
| min_quantity | DECIMAL(10,2) | | 最低库存 |
| max_quantity | DECIMAL(10,2) | | 最高库存 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**索引**:
- `idx_stocks_material_project` ON (material_id, project_id)

---

### construction_appointments 施工预约表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| appointment_no | VARCHAR(50) | UNIQUE | 预约单号 |
| project_id | INTEGER | REFERENCES projects(id) | 项目ID |
| applicant_id | INTEGER | REFERENCES users(id) | 申请人ID |
| applicant_name | VARCHAR(100) | | 申请人姓名 |
| contact_phone | VARCHAR(20) | | 联系电话 |
| contact_person | VARCHAR(100) | | 联系人 |
| work_date | DATE | NOT NULL | 作业日期 |
| time_slot | VARCHAR(50) | NOT NULL | 时间段 |
| work_location | VARCHAR(500) | NOT NULL | 作业地点 |
| work_content | TEXT | NOT NULL | 作业内容 |
| work_type | VARCHAR(50) | | 作业类型 |
| is_urgent | BOOLEAN | DEFAULT false | 是否加急 |
| priority | INTEGER | DEFAULT 0 | 优先级 |
| urgent_reason | TEXT | | 加急原因 |
| assigned_worker_ids | TEXT | | 作业人员ID列表（JSON） |
| assigned_worker_names | TEXT | | 作业人员姓名 |
| supervisor_id | INTEGER | REFERENCES users(id) | 监护人ID |
| supervisor_name | VARCHAR(100) | | 监护人姓名 |
| status | VARCHAR(20) | DEFAULT 'draft' | 状态 |
| workflow_instance_id | INTEGER | | 工作流实例ID |
| submitted_at | TIMESTAMP | | 提交时间 |
| approved_at | TIMESTAMP | | 审批时间 |
| completed_at | TIMESTAMP | | 完成时间 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**索引**:
- `idx_appointments_date` ON (work_date)
- `idx_appointments_status` ON (status)
- `idx_appointments_applicant` ON (applicant_id)

**时间段值**:
- `morning`: 上午 (8:00-11:30)
- `noon`: 中午 (12:00-13:30)
- `afternoon`: 下午 (13:30-16:30)
- `full_day`: 全天

---

### worker_calendars 作业人员日历表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| worker_id | INTEGER | REFERENCES users(id) | 作业人员ID |
| calendar_date | DATE | NOT NULL | 日期 |
| time_slot | VARCHAR(20) | NOT NULL | 时间段 |
| is_available | BOOLEAN | DEFAULT true | 是否可用 |
| status | VARCHAR(20) | DEFAULT 'available' | 状态 |
| appointment_id | INTEGER | REFERENCES construction_appointments(id) | 关联预约ID |
| blocked_reason | TEXT | | 锁定原因 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**唯一约束**: (worker_id, calendar_date, time_slot)

---

### attendance_records 考勤记录表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| user_id | INTEGER | REFERENCES users(id) | 用户ID |
| appointment_id | INTEGER | REFERENCES construction_appointments(id) | 关联预约ID |
| work_content | TEXT | | 工作内容 |
| attendance_type | VARCHAR(20) | NOT NULL | 打卡类型 |
| clock_in_time | TIMESTAMP | NOT NULL | 打卡时间 |
| clock_in_location | VARCHAR(200) | | 打卡地点 |
| clock_in_latitude | DECIMAL(10,7) | | 纬度 |
| clock_in_longitude | DECIMAL(10,7) | | 经度 |
| overtime_hours | DECIMAL(4,1) | | 加班小时数 |
| remark | TEXT | | 备注 |
| photo_url | VARCHAR(500) | | 照片URL |
| photo_urls | TEXT | | 照片URL列表（JSON） |
| status | VARCHAR(20) | DEFAULT 'pending' | 状态 |
| confirmed_by | INTEGER | REFERENCES users(id) | 确认人ID |
| confirmed_at | TIMESTAMP | | 确认时间 |
| confirmed_remark | TEXT | | 确认备注 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**打卡类型**:
- `morning`: 上午打卡
- `afternoon`: 下午打卡
- `noon_overtime`: 中午加班
- `night_overtime`: 晚上加班

---

### notifications 通知表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| user_id | INTEGER | REFERENCES users(id) | 用户ID |
| title | VARCHAR(200) | NOT NULL | 标题 |
| content | TEXT | | 内容 |
| type | VARCHAR(50) | | 通知类型 |
| data | JSONB | | 附加数据 |
| is_read | BOOLEAN | DEFAULT false | 是否已读 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

**索引**:
- `idx_notifications_user` ON (user_id)
- `idx_notifications_read` ON (is_read)

---

### workflow_definitions 工作流定义表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| name | VARCHAR(100) | NOT NULL | 工作流名称 |
| description | TEXT | | 描述 |
| type | VARCHAR(50) | | 工作流类型 |
| is_active | BOOLEAN | DEFAULT false | 是否激活 |
| version | INTEGER | DEFAULT 1 | 版本 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

---

### workflow_nodes 工作流节点表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| workflow_id | INTEGER | REFERENCES workflow_definitions(id) | 工作流ID |
| name | VARCHAR(100) | NOT NULL | 节点名称 |
| type | VARCHAR(50) | NOT NULL | 节点类型 |
| order | INTEGER | | 顺序 |
| approver_type | VARCHAR(50) | | 审批人类型 |
| approver_ids | TEXT | | 审批人ID列表 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

---

### ai_conversations AI对话历史表

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| user_id | INTEGER | REFERENCES users(id) | 用户ID |
| role | VARCHAR(20) | NOT NULL | 角色（user/assistant） |
| content | TEXT | NOT NULL | 内容 |
| intent | VARCHAR(100) | | 意图 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

**索引**:
- `idx_ai_conversations_user` ON (user_id)

---

## 数据库维护

### 备份命令

```bash
# 全库备份
pg_dump -U materials -d materials > backup_$(date +%Y%m%d).sql

# 恢复
psql -U materials -d materials < backup_20240115.sql
```

### 常用查询

```sql
-- 查看所有表
SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';

-- 查看表结构
\d table_name

-- 查看索引
SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'users';

-- 查看表大小
SELECT pg_size_pretty(pg_total_relation_size('users'));
```

### 性能优化建议

1. 为经常查询的字段创建索引
2. 定期执行 VACUUM ANALYZE 清理和更新统计信息
3. 大表考虑分区
4. 监控慢查询日志

## 数据迁移

迁移文件存放在 `migrations/` 目录，按顺序执行：

```bash
# 执行单个迁移
psql -U materials -d materials -f migrations/001_create_auth_tables.sql

# 执行所有迁移
for f in migrations/*.sql; do
  psql -U materials -d materials -f "$f"
done
```
