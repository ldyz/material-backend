# 项目目录结构文档

## 项目根目录
**路径**: `/home/julei/backend`

## 目录结构

```
/home/julei/backend/
├── cmd/                          # Go 后端命令入口
│   ├── server/                   # 后端服务器主程序
│   │   └── main.go              # 后端入口文件
│   └── check-data/              # 数据检查工具
│
├── internal/                     # Go 后端内部代码
│   ├── api/                     # API 处理器
│   │   ├── auth/               # 认证模块
│   │   ├── construction_log/   # 施工日志模块
│   │   ├── upload/             # 文件上传模块
│   │   ├── material/           # 物资管理模块
│   │   ├── project/            # 项目管理模块
│   │   ├── inbound/            # 入库管理模块
│   │   ├── stock/              # 库存管理模块
│   │   ├── requisition/        # 领用管理模块
│   │   ├── system/             # 系统管理模块
│   │   ├── workflow/           # 工作流模块
│   │   ├── notification/       # 通知模块
│   │   └── progress/           # 进度管理模块
│   ├── db/                     # 数据库连接
│   └── middleware/             # 中间件
│
├── pkg/                         # Go 公共包
│   └── jwt/                    # JWT 认证包
│
├── migrations/                  # 数据库迁移文件
│   └── 014_add_construction_log_fields.sql
│
├── static/                      # 静态文件目录
│   ├── uploads/                # 上传文件存储目录
│   ├── assets/                 # 前端编译输出
│   └── index.html              # 前端入口
│
├── newstatic/                   # Vue 3 前端项目（新Web端）
│   ├── src/                    # 源代码
│   │   ├── views/             # 页面组件
│   │   ├── components/        # 公共组件
│   │   ├── api/               # API 接口
│   │   ├── stores/            # 状态管理
│   │   ├── router/            # 路由配置
│   │   └── utils/             # 工具函数
│   ├── package.json           # npm 依赖配置
│   └── vite.config.js         # Vite 配置
│
├── mobile-app/                  # 移动端项目（Vue 3 + Capacitor）
│   ├── src/                   # 源代码
│   │   ├── views/             # 页面组件
│   │   ├── components/        # 公共组件
│   │   ├── api/               # API 接口
│   │   ├── stores/            # 状态管理
│   │   ├── router/            # 路由配置
│   │   └── utils/             # 工具函数
│   ├── android/               # Android 原生代码
│   ├── package.json           # npm 依赖配置
│   ├── capacitor.config.json  # Capacitor 配置
│   └── .env.development       # 开发环境变量
│
├── server                       # 后端编译输出（可执行文件）
├── server.log                   # 后端运行日志
└── go.mod                       # Go 模块配置
```

## 编译命令速查表

### 后端编译

```bash
# 切换到后端根目录
cd /home/julei/backend

# 编译后端服务器
go build -o server ./cmd/server/main.go

# 重启后端服务器
pkill -f "./server$" && nohup ./server > server.log 2>&1 &

# 查看后端日志
tail -f server.log
```

### 前端（newstatic）编译

```bash
# 切换到前端目录
cd /home/julei/backend/newstatic

# 安装依赖（首次）
npm install

# 开发模式运行
npm run dev

# 生产环境编译
npm run build

# 编译输出目录
# /home/julei/backend/static/
```

### 移动端（mobile-app）编译

```bash
# 切换到移动端目录
cd /home/julei/backend/mobile-app

# 安装依赖（首次）
npm install

# 开发模式运行
npm run dev

# 生产环境编译
npm run build

# 同步到 Android
npm run build
npx cap sync android
```

## 重要路径总结

| 项目 | 源代码路径 | 编译路径 | 工作目录 |
|------|-----------|---------|---------|
| 后端 | `/home/julei/backend/cmd/server/main.go` | `/home/julei/backend/server` | `/home/julei/backend` |
| 前端（新Web） | `/home/julei/backend/newstatic/src/` | `/home/julei/backend/static/` | `/home/julei/backend/newstatic` |
| 移动端 | `/home/julei/backend/mobile-app/src/` | `/home/julei/backend/mobile-app/dist/` | `/home/julei/backend/mobile-app` |

## 服务器配置

- **后端端口**: 8088
- **后端日志**: `/home/julei/backend/server.log`
- **上传目录**: `/home/julei/backend/static/uploads/`
- **进程查看**: `ps aux \| grep "[s]erver$"`

## API 地址

- **开发环境**: `http://localhost:8088/api`
- **生产环境**: `https://home.mbed.org.cn:8088/api`
- **上传接口**: `/api/upload/image`
- **健康检查**: `/api/health`

## 注意事项

1. **编译后端时必须先切换到 `/home/julei/backend` 目录**
2. **编译前端时必须先切换到 `/home/julei/backend/newstatic` 目录**
3. **编译移动端时必须先切换到 `/home/julei/backend/mobile-app` 目录**
4. 前端编译后会自动同步到 `/home/julei/backend/static/` 目录
5. 后端服务器运行在 8088 端口，使用 `nohup` 后台运行
6. 上传目录权限必须为 755，否则会导致 500 错误

## 快速命令别名

可以在 `~/.bashrc` 中添加以下别名：

```bash
# 后端编译
alias backend-build='cd /home/julei/backend && go build -o server ./cmd/server/main.go'

# 后端重启
alias backend-restart='cd /home/julei/backend && pkill -f "./server$" && nohup ./server > server.log 2>&1 &'

# 前端编译
alias frontend-build='cd /home/julei/backend/newstatic && npm run build'

# 移动端编译
alias mobile-build='cd /home/julei/backend/mobile-app && npm run build'

# 查看日志
alias logs='tail -f /home/julei/backend/server.log'
```
