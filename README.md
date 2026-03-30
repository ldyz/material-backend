# 工地物资管理系统

一个完整的工地物资管理系统，包含Web管理后台和移动端应用，支持物资计划、入库、出库、库存管理、施工预约、考勤管理等功能，并集成了AI智能助手。

## 技术栈

| 层级 | 技术 |
|-----|------|
| 后端 | Go 1.21+ / Gin / GORM / PostgreSQL |
| Web前端 | Vue 3 / Element Plus / Pinia / Vite |
| 移动端 | Vue 3 / Vant / Capacitor |
| 实时通信 | WebSocket |
| AI集成 | 百度千帆 / DeepSeek |
| 文件存储 | 本地存储 |

## 项目结构

```
/home/julei/backend/
├── cmd/server/main.go          # 后端入口文件
├── internal/
│   ├── api/                    # API模块
│   │   ├── agent/              # AI助手
│   │   ├── appointment/        # 施工预约
│   │   ├── attendance/         # 考勤管理
│   │   ├── audit/              # 操作日志
│   │   ├── auth/               # 用户认证
│   │   ├── construction_log/   # 施工日志
│   │   ├── inbound/            # 入库管理
│   │   ├── material/           # 物资管理
│   │   ├── material_master/    # 物资主数据
│   │   ├── material_plan/      # 物资计划
│   │   ├── notification/       # 通知服务
│   │   ├── progress/           # 进度管理
│   │   ├── project/            # 项目管理
│   │   ├── requisition/        # 出库/领用
│   │   ├── stock/              # 库存管理
│   │   ├── system/             # 系统配置
│   │   ├── upload/             # 文件上传
│   │   ├── workflow/           # 工作流引擎
│   │   └── app/                # 应用版本管理
│   ├── config/                 # 配置管理
│   ├── db/                     # 数据库连接
│   └── middleware/             # 中间件
├── pkg/
│   ├── jwt/                    # JWT工具
│   └── openai/                 # AI服务封装
├── migrations/                 # 数据库迁移文件
├── newstatic/                  # Web端前端
├── mobile-app/                 # 移动端前端
├── mobile-app-updates/         # 移动端OTA更新包
├── static/uploads/             # 上传文件存储
├── scripts/                    # 脚本文件
├── rebuild.sh                  # 后端编译重启脚本
├── Makefile                    # Make命令
└── docs/                       # 项目文档
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- pnpm / npm

### 后端启动

```bash
# 编译并重启服务
./rebuild.sh

# 或使用 Makefile
make restart

# 查看日志
tail -f /tmp/backend.log
```

### Web前端启动

```bash
cd newstatic

# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
./build.sh
```

### 移动端启动

```bash
cd mobile-app

# 安装依赖
npm install

# 开发模式
npm run dev

# Android构建
./build.sh
```

## 功能模块

### Web端功能

| 模块 | 功能 |
|-----|------|
| 项目管理 | 项目创建、编辑、成员管理、层级结构 |
| 物资主数据 | 物资分类、物资信息管理 |
| 物资计划 | 计划编制、审批、执行 |
| 入库管理 | 入库单创建、审批、查询 |
| 出库管理 | 出库申请、审批、执行 |
| 库存管理 | 库存查询、盘点、调拨 |
| 施工预约 | 预约创建、审批、日历视图 |
| 考勤管理 | 打卡记录、统计报表 |
| 进度管理 | 进度填报、审核、甘特图 |
| 施工日志 | 日志录入、查询 |
| 工作流管理 | 流程定义、节点配置 |
| 系统管理 | 用户、角色、权限配置 |
| AI助手 | 智能对话、数据查询、语音交互 |

### 移动端功能

| 模块 | 功能 |
|-----|------|
| 计划管理 | 查看、创建、审批物资计划 |
| 入库管理 | 查看、创建、审批入库单 |
| 出库管理 | 查看、创建、审批出库单 |
| 预约管理 | 查看、创建、审批施工预约、日历视图 |
| 考勤打卡 | GPS定位打卡、拍照签到、加班记录 |
| 消息通知 | 实时消息推送、WebSocket |
| AI助手 | 语音交互、智能问答 |

## API模块

| 模块 | 路径前缀 | 功能 |
|-----|---------|------|
| auth | /api/auth | 用户认证、用户管理、角色管理 |
| project | /api/project | 项目管理、成员分配 |
| material_master | /api/material-master | 物资主数据管理 |
| material_plan | /api/material-plans | 物资计划管理 |
| inbound | /api/inbound | 入库管理 |
| requisition | /api/requisitions | 出库管理 |
| stock | /api/stock | 库存管理 |
| appointment | /api/appointments | 施工预约、日历管理 |
| attendance | /api/attendance | 考勤管理 |
| notification | /api/notification | 通知服务、WebSocket |
| agent | /api/agent | AI助手、语音交互 |
| workflow | /api/workflows | 工作流管理 |
| progress | /api/progress | 进度管理 |
| construction_log | /api/construction-logs | 施工日志 |
| audit | /api/audit | 操作日志 |
| system | /api/system | 系统配置 |
| upload | /api/upload | 文件上传 |
| app | /api/app | 应用版本管理 |

## 数据库

使用PostgreSQL数据库，通过GORM进行ORM映射。

### 主要数据表

| 表名 | 说明 |
|-----|------|
| users | 用户表 |
| roles | 角色表 |
| user_roles | 用户角色关联表 |
| permissions | 权限表 |
| projects | 项目表（支持层级结构） |
| user_projects | 用户项目关联表 |
| materials | 物资主数据 |
| material_categories | 物资分类 |
| material_plans | 物资计划 |
| material_plan_items | 计划明细 |
| inbound_orders | 入库单 |
| inbound_order_items | 入库明细 |
| requisitions | 出库单 |
| requisition_items | 出库明细 |
| stocks | 库存表 |
| construction_appointments | 施工预约 |
| worker_calendars | 作业人员日历 |
| attendance_records | 考勤记录 |
| monthly_attendance_summary | 月度考勤汇总 |
| notifications | 通知表 |
| workflow_definitions | 工作流定义 |
| workflow_instances | 工作流实例 |
| ai_conversations | AI对话历史 |

## 配置说明

配置文件：`config.yaml`

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

ai:
  provider: baidu           # AI提供者: baidu/deepseek
  baidu_api_key: xxx
  baidu_secret_key: xxx
  deepseek_api_key: xxx
```

## 文档目录

- [后端API文档](docs/BACKEND_API.md) - 详细的API接口说明
- [数据库文档](docs/DATABASE.md) - 表结构和迁移说明
- [Web前端文档](docs/WEB_FRONTEND.md) - 组件和页面说明
- [移动端文档](docs/MOBILE_FRONTEND.md) - 移动端开发指南
- [开发指南](docs/DEVELOPMENT.md) - 环境搭建和开发流程
- [部署指南](docs/DEPLOYMENT.md) - 生产部署说明

## Makefile命令

| 命令 | 说明 |
|------|------|
| `make restart` | 停止、编译、重启服务 |
| `make build` | 编译到 bin/server |
| `make build-fast` | 快速编译到 ./server |
| `make run` | 使用 go run 启动 |
| `make run-bg` | 后台运行服务 |
| `make stop` | 停止服务 |
| `make logs` | 实时查看日志 |
| `make log` | 查看最近50行日志 |

## 权限系统

系统采用RBAC（基于角色的访问控制）模型，支持多角色关联。

### 预置角色

| 角色 | 说明 |
|------|------|
| admin | 系统管理员，拥有所有权限 |
| 项目经理 | 管理项目和进度 |
| 材料员 | 管理物资和计划 |
| 保管员 | 管理库存和出入库 |
| 施工员 | 审批预约和任务 |
| 作业人员 | 执行施工任务 |
| 预约管理员 | 管理施工预约 |

### 权限模块

系统包含以下权限模块：
- 用户管理（user_view, user_create, user_edit, user_delete）
- 角色管理（role_view, role_create, role_edit, role_delete, role_assign_permissions）
- 项目管理（project_view, project_create, project_edit, project_delete）
- 物资管理（material_view, material_create, material_edit, material_delete, material_import）
- 物资计划（material_plan_view, material_plan_create, material_plan_edit, material_plan_delete, material_plan_approve）
- 库存管理（stock_view, stock_create, stock_edit, stock_in, stock_out, stock_export, stock_alerts）
- 入库管理（inbound_view, inbound_create, inbound_edit, inbound_delete, inbound_approve, inbound_export）
- 出库管理（requisition_view, requisition_create, requisition_edit, requisition_delete, requisition_approve, requisition_issue）
- 施工预约（appointment_view, appointment_create, appointment_edit, appointment_approve, appointment_assign）
- 考勤管理（attendance_view, attendance_clock_in, attendance_confirm）
- 工作流管理（workflow_view, workflow_create, workflow_task_approve等）
- AI助手（ai_agent_view, ai_agent_query, ai_agent_operate）
- 系统管理（system_log, system_backup, system_config, system_statistics）

## 开发团队

本项目为内部管理系统，如有问题请联系开发团队。

## 许可证

私有项目，仅供内部使用。
