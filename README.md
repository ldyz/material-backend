# 材料管理系统 - 后端服务

> 基于 Go + Gin + GORM 的材料管理系统后端 API 服务

## 项目简介

这是材料管理系统的后端服务，提供 RESTful API 接口，支持项目管理、物资管理、库存管理、进度管理、施工预约等核心功能。

## 技术栈

| 类别 | 技术 |
|------|------|
| 语言 | Go 1.18+ |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | PostgreSQL / MySQL / SQLite |
| 认证 | JWT |
| 实时通信 | WebSocket |
| 密码加密 | bcrypt |

## 项目结构

```
backend/
├── cmd/                        # 命令行工具
│   ├── server/                 # 主服务入口
│   ├── check-admin/            # 检查管理员
│   ├── check-data/             # 数据检查工具
│   ├── check-user-roles/       # 角色检查
│   ├── fix-permissions/        # 权限修复工具
│   ├── fix-task-tree/          # 任务树修复
│   └── reset-admin/            # 重置管理员
├── internal/                   # 内部代码
│   ├── api/                    # API 模块
│   │   ├── agent/              # AI 代理服务
│   │   ├── app/                # 应用版本管理
│   │   ├── appointment/        # 施工预约
│   │   ├── audit/              # 操作日志/审计
│   │   ├── auth/               # 认证授权
│   │   ├── construction_log/   # 施工日志
│   │   ├── inbound/            # 入库管理
│   │   ├── material/           # 物资管理
│   │   ├── material_master/    # 物资主数据
│   │   ├── material_plan/      # 物资计划
│   │   ├── notification/       # 通知系统 + WebSocket
│   │   ├── progress/           # 进度管理/甘特图
│   │   ├── project/            # 项目管理
│   │   ├── requisition/        # 出库申请
│   │   ├── response/           # 统一响应格式
│   │   ├── stock/              # 库存管理
│   │   ├── system/             # 系统管理
│   │   ├── upload/             # 文件上传
│   │   └── workflow/           # 工作流引擎
│   ├── config/                 # 配置管理
│   ├── db/                     # 数据库连接
│   └── middleware/             # 中间件 (CORS, 安全等)
├── migrations/                 # 数据库迁移脚本 (19个版本)
├── scripts/                    # 工具脚本
│   └── migrations/             # 增量迁移脚本
├── static/                     # 静态文件
│   └── uploads/                # 上传文件存储
├── mobile-app/                 # 移动端应用 (Capacitor + Vue 3)
├── newstatic/                  # 前端项目 (Vue 3)
├── docs/                       # 文档
├── config.yaml                 # 配置文件
├── config.example.yaml         # 配置模板
├── Makefile                    # 构建脚本
└── Dockerfile                  # Docker 配置
```

## 快速开始

### 环境要求

- Go 1.18+
- PostgreSQL 14+ / MySQL 8+ / SQLite 3
- Node.js 18+ (前端构建)

### 安装依赖

```bash
go mod download
```

### 配置

复制配置模板并修改：

```bash
cp config.example.yaml config.yaml
```

主要配置项：

```yaml
server:
  port: 8088                # 服务端口
  mode: debug               # 运行模式: debug/release

database:
  type: postgresql          # 数据库类型
  host: 127.0.0.1
  port: 5432
  user: materials
  password: your-password
  database: materials

jwt:
  secret: "your-secret-key" # JWT密钥(生产环境必须修改)
  expire_time: 24h          # Token过期时间

upload:
  max_file_size: 5242880    # 最大文件大小 5MB
  upload_dir: "static/uploads"
```

### 数据库初始化

```bash
# 按顺序运行迁移脚本
for f in migrations/*.sql; do
  psql -U materials -d materials -f "$f"
done
```

### 运行

```bash
# 开发模式
go run cmd/server/main.go

# 编译后运行
go build -o server cmd/server/main.go
./server -c config.yaml
```

## API 模块

### 认证模块 `/api/auth`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| POST | `/auth/login` | - | 用户登录 |
| POST | `/auth/logout` | 需登录 | 用户登出 |
| GET | `/auth/me` | 需登录 | 获取当前用户信息 |
| POST | `/auth/change-password` | 需登录 | 修改密码 |
| POST | `/auth/avatar` | 需登录 | 上传头像 |
| GET | `/auth/users` | user_view | 用户列表 |
| POST | `/auth/users` | user_create | 创建用户 |
| GET | `/auth/users/:id` | user_view | 用户详情 |
| PUT | `/auth/users/:id` | user_edit | 更新用户 |
| DELETE | `/auth/users/:id` | user_delete | 删除用户 |
| POST | `/auth/users/:id/reset-password` | user_edit | 重置密码 |
| GET | `/auth/roles` | role_view | 角色列表 |
| POST | `/auth/roles` | role_create | 创建角色 |
| PUT | `/auth/roles/:id` | role_edit | 更新角色 |
| DELETE | `/auth/roles/:id` | role_delete | 删除角色 |
| GET | `/auth/permissions` | role_view | 权限列表 |

### 项目管理 `/api/project`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/project/projects` | project_view | 项目列表 (支持分页/搜索/筛选) |
| POST | `/project/projects` | project_create | 创建项目 |
| GET | `/project/projects/:id` | project_view | 项目详情 |
| PUT | `/project/projects/:id` | project_edit | 更新项目 |
| DELETE | `/project/projects/:id` | project_delete | 删除项目 |
| GET | `/project/projects/:id/members` | project_view | 项目成员 |
| POST | `/project/projects/:id/members` | project_edit | 分配成员 |
| GET | `/project/projects/:id/tree` | project_view | 项目树 |
| GET | `/project/projects/:id/children` | project_view | 子项目列表 |
| POST | `/project/projects/:id/aggregate-progress` | project_edit | 聚合进度 |

### 物资管理 `/api/material`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/material/materials` | material_view | 物资列表 |
| POST | `/material/materials` | material_create | 创建物资 |
| GET | `/material/materials/:id` | material_view | 物资详情 |
| PUT | `/material/materials/:id` | material_edit | 更新物资 |
| DELETE | `/material/materials/:id` | material_delete | 删除物资 |
| GET | `/material/categories` | material_view | 物资分类列表 |
| POST | `/material/categories` | material_create | 创建分类 |

### 库存管理 `/api/stock`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/stock/stocks` | stock_view | 库存列表 |
| GET | `/stock/stocks/:id` | stock_view | 库存详情 |
| PUT | `/stock/stocks/:id` | stock_edit | 更新库存 |
| POST | `/stock/stocks/:id/in` | stock_in | 入库操作 |
| POST | `/stock/stocks/:id/out` | stock_out | 出库操作 |
| GET | `/stock/logs` | stocklog_view | 库存日志 |

### 入库管理 `/api/inbound`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/inbound/orders` | inbound_view | 入库单列表 |
| POST | `/inbound/orders` | inbound_create | 创建入库单 |
| GET | `/inbound/orders/:id` | inbound_view | 入库单详情 |
| PUT | `/inbound/orders/:id` | inbound_edit | 更新入库单 |
| DELETE | `/inbound/orders/:id` | inbound_delete | 删除入库单 |
| POST | `/inbound/orders/:id/approve` | inbound_approve | 审核入库单 |

### 出库管理 `/api/requisition`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/requisition/requisitions` | requisition_view | 出库单列表 |
| POST | `/requisition/requisitions` | requisition_create | 创建出库单 |
| GET | `/requisition/requisitions/:id` | requisition_view | 出库单详情 |
| PUT | `/requisition/requisitions/:id` | requisition_edit | 更新出库单 |
| DELETE | `/requisition/requisitions/:id` | requisition_delete | 删除出库单 |
| POST | `/requisition/requisitions/:id/approve` | requisition_approve | 审核出库单 |
| POST | `/requisition/requisitions/:id/issue` | requisition_issue | 发货 |

### 进度管理 `/api/progress`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/progress` | progress_view | 进度列表 |
| GET | `/progress/project/:id` | progress_view | 项目进度计划 |
| PUT | `/progress/project/:id` | progress_edit | 更新进度计划 |
| GET | `/progress/project/:id/tasks` | progress_view | 任务列表 |
| POST | `/progress/project/:id/tasks` | progress_create | 创建任务 |
| PUT | `/progress/tasks/:id` | progress_edit | 更新任务 |
| DELETE | `/progress/tasks/:id` | progress_delete | 删除任务 |
| GET | `/progress/tasks/:id/dependencies` | progress_view | 任务依赖 |
| POST | `/progress/project/:id/resources` | progress_create | 创建资源 |
| POST | `/progress/project/:id/generate-plan` | progress_create | AI生成计划 |

### 施工预约 `/api/appointments`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/appointments` | appointment_view | 预约列表 |
| GET | `/appointments/my` | 需登录 | 我的预约 |
| POST | `/appointments` | appointment_create | 创建预约 |
| GET | `/appointments/:id` | appointment_view | 预约详情 |
| PUT | `/appointments/:id` | appointment_edit | 更新预约 |
| DELETE | `/appointments/:id` | appointment_delete | 删除预约 |
| POST | `/appointments/:id/submit` | appointment_submit | 提交审批 |
| POST | `/appointments/:id/approve` | appointment_approve | 审批预约 |
| GET | `/appointments/workers` | appointment_view | 作业人员列表 |

### 施工日志 `/api/construction-logs`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/construction-logs` | construction_log_view | 日志列表 |
| POST | `/construction-logs` | construction_log_create | 创建日志 |
| GET | `/construction-logs/:id` | construction_log_view | 日志详情 |
| PUT | `/construction-logs/:id` | construction_log_edit | 更新日志 |
| DELETE | `/construction-logs/:id` | construction_log_delete | 删除日志 |

### 工作流 `/api/workflow`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/workflow/definitions` | workflow_view | 工作流定义列表 |
| POST | `/workflow/definitions` | workflow_create | 创建工作流 |
| GET | `/workflow/definitions/:id` | workflow_view | 工作流详情 |
| PUT | `/workflow/definitions/:id` | workflow_edit | 更新工作流 |
| POST | `/workflow/definitions/:id/activate` | workflow_activate | 激活工作流 |
| GET | `/workflow/pending-tasks` | workflow_task_view | 待办任务 |
| POST | `/workflow/tasks/:id/approve` | workflow_task_approve | 审批任务 |

### 通知系统 `/api/notification`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/notification/notifications` | 需登录 | 通知列表 |
| PUT | `/notification/notifications/:id/read` | 需登录 | 标记已读 |
| GET | `/notification/ws` | 需登录 | WebSocket 连接 |

### 系统管理 `/api/system`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/system/configs` | system_config | 系统配置 |
| PUT | `/system/configs` | system_config | 更新配置 |
| GET | `/system/logs` | system_log | 系统日志 |
| GET | `/system/backups` | system_backup | 备份列表 |
| POST | `/system/backups` | system_backup | 创建备份 |
| GET | `/system/statistics` | system_statistics | 系统统计 |
| GET | `/system/activities` | system_activities | 系统动态 |

### 操作审计 `/api/audit`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/audit/logs` | audit_view | 操作日志列表 |
| GET | `/audit/logs/:id` | audit_view | 日志详情 |

### 文件上传 `/api/upload`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| POST | `/upload/upload` | 需登录 | 单文件上传 |
| POST | `/upload/upload-multiple` | 需登录 | 多文件上传 |

### AI 代理 `/api/agent`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| POST | `/agent/query` | ai_agent_query | AI 查询 |
| GET | `/agent/logs` | ai_agent_logs | AI 操作日志 |

### 应用版本 `/api/app`

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/app/version` | - | 获取最新版本信息 |
| GET | `/app/versions` | - | 版本历史 |

## 权限系统

系统采用 RBAC (基于角色的访问控制) 模型，支持多角色关联。

### 预置角色

| 角色 | 说明 |
|------|------|
| admin | 系统管理员，拥有所有权限 |
| 项目经理 | 管理项目和进度 |
| 材料员 | 管理物资和计划 |
| 保管员 | 管理库存和出入库 |
| 施工员 | 审批预约和任务 |
| 作业人员 | 执行施工任务 |

### 权限模块

| 模块 | 权限数量 | 示例 |
|------|----------|------|
| 用户管理 | 4 | user_view, user_create, user_edit, user_delete |
| 角色管理 | 5 | role_view, role_create, role_edit, role_delete, role_assign_permissions |
| 项目管理 | 4 | project_view, project_create, project_edit, project_delete |
| 物资管理 | 6 | material_view, material_create, material_edit, material_delete, material_import, material_export |
| 库存管理 | 8 | stock_view, stock_in, stock_out, stock_edit, stock_delete, stock_export, stock_alerts, stocklog_view |
| 入库管理 | 6 | inbound_view, inbound_create, inbound_edit, inbound_delete, inbound_approve, inbound_export |
| 出库管理 | 7 | requisition_view, requisition_create, requisition_edit, requisition_delete, requisition_approve, requisition_issue, requisition_export |
| 进度管理 | 5 | progress_view, progress_create, progress_edit, progress_delete, progress_export |
| 施工日志 | 5 | construction_log_view, construction_log_create, construction_log_edit, construction_log_delete, construction_log_export |
| 工作流管理 | 12 | workflow_view, workflow_create, workflow_task_approve 等 |
| 系统管理 | 5 | system_log, system_backup, system_config, system_statistics, system_activities |
| 审计日志 | 1 | audit_view |
| AI 智能体 | 5 | ai_agent_view, ai_agent_query, ai_agent_operate, ai_agent_workflow, ai_agent_logs |

## 数据库模型

### 核心表

| 表名 | 说明 |
|------|------|
| users | 用户表 |
| roles | 角色表 |
| user_roles | 用户角色关联表 |
| projects | 项目表 (支持层级结构) |
| user_projects | 用户项目关联表 |
| materials | 物资表 |
| material_categories | 物资分类表 |
| material_plans | 物资计划表 |
| material_plan_items | 计划明细表 |
| stocks | 库存表 |
| stock_logs | 库存日志表 |
| inbound_orders | 入库单表 |
| inbound_order_items | 入库明细表 |
| requisitions | 出库单表 |
| requisition_items | 出库明细表 |
| project_schedules | 项目进度计划表 |
| schedule_tasks | 任务表 |
| task_dependencies | 任务依赖表 |
| resources | 资源表 |
| task_resources | 任务资源分配表 |
| construction_appointments | 施工预约单表 |
| worker_calendars | 作业人员日历表 |
| construction_logs | 施工日志表 |
| workflow_definitions | 工作流定义表 |
| workflow_nodes | 工作流节点表 |
| workflow_instances | 工作流实例表 |
| workflow_pending_tasks | 待办任务表 |
| notifications | 通知表 |
| device_tokens | 设备推送令牌表 |
| system_configs | 系统配置表 |
| system_logs | 系统日志表 |
| operation_logs | 操作审计日志表 |

## 构建和部署

### 后端构建

```bash
# 编译
go build -o server cmd/server/main.go

# 或使用 Makefile
make build

# 运行
./server -c config.yaml
```

### PC 前端构建 (newstatic)

```bash
cd newstatic
npm run build
# 输出到 dist/ 目录
```

### 移动端构建 (mobile-app)

**重要：移动端构建必须设置 `CAPACITOR_BUILD=true` 环境变量，否则资源路径会使用绝对路径导致白屏！**

```bash
cd mobile-app

# 方式1：使用 package.json 预定义命令（推荐）
npm run android:build    # 构建 Debug APK（自动设置环境变量）

# 方式2：手动构建
CAPACITOR_BUILD=true npm run build   # 构建前端（输出到 dist-capacitor/）
npx cap sync android                  # 同步到 Android 项目
cd android && ./gradlew assembleDebug # 构建 APK

# 构建产物位置
# Android APK: mobile-app/android/app/build/outputs/apk/debug/app-debug.apk
```

### 移动端版本更新流程

```bash
# 1. 更新版本号
# 编辑 mobile-app/public/version.json
# 编辑 mobile-app/android/app/build.gradle (versionCode, versionName)

# 2. 构建 APK
cd mobile-app
CAPACITOR_BUILD=true npm run build
npx cap sync android
cd android && ./gradlew assembleDebug

# 3. 复制 APK 到更新目录
cp app/build/outputs/apk/debug/app-debug.apk ../../mobile-app-updates/android/material-management-X.X.X.apk

# 4. 创建版本信息文件
# 创建 mobile-app-updates/android/latest.json
# 创建 mobile-app-updates/android/version-X.X.X.json

# 5. 更新数据库版本记录
psql -U materials -d materials -c "
INSERT INTO app_versions (platform, version, build_number, download_url, force_update, update_message, release_notes, published_at, created_at, updated_at)
VALUES ('android', 'X.X.X', XXX, 'https://home.mbed.org.cn:9090/mobile-updates/android/material-management-X.X.X.apk', false, '更新说明', '详细更新内容', NOW(), NOW(), NOW());
"
```

### Docker 部署

```bash
# 构建镜像
docker build -t material-backend .

# 运行容器
docker run -d -p 8088:8088 -v ./config.yaml:/app/config.yaml material-backend
```

### 生产部署要点

1. 修改 `config.yaml` 中的敏感配置
2. 设置 `server.mode: release`
3. 使用 Nginx 反向代理
4. 配置 HTTPS
5. 定期备份数据库

## 静态资源服务

服务启动后会自动提供：

- `/static/*` → 前端构建文件 (`newstatic/dist`)
- `/mobile/*` → 移动端应用 (`mobile-app/dist`)
- `/uploads/*` → 上传文件 (`static/uploads`)
- `/mobile-updates/*` → 移动端更新包

## WebSocket 支持

通知模块提供 WebSocket 实时通信：

```javascript
// 连接 WebSocket
const ws = new WebSocket('ws://localhost:8088/api/notification/ws?token=YOUR_JWT_TOKEN')

ws.onmessage = (event) => {
  const notification = JSON.parse(event.data)
  console.log('收到通知:', notification)
}
```

## 相关文档

- [配置说明](CONFIG.md)
- [API 接口详情](docs/API_ENDPOINTS.md)
- [数据库迁移](docs/DATABASE_MIGRATIONS.sql)
- [角色权限配置](docs/ROLES_PERMISSIONS.md)
- [前端项目文档](newstatic/README.md)

## 许可证

MIT License
