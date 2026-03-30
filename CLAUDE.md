# 项目开发指南

## 重要规则

### 后端编译规则
**编译后端时必须使用脚本，不要直接运行 go build 命令。**

- 后端编译重启脚本: `./rebuild.sh`
- 如果脚本不存在，先创建脚本再执行
- 脚本功能: 停止旧服务 -> 编译 -> 启动新服务

### 前端构建规则
**构建前端时必须使用脚本，不要直接运行 npm run build 命令。**

- 移动端构建脚本: `./mobile-app/build.sh`
- 管理后台构建脚本: `./newstatic/build.sh`
- 如果脚本不存在，先创建脚本再执行

### 移动端发布规则
**修改移动端文件后，如果需要发布新版本，必须使用发布脚本。**

- 移动端发布脚本: `./mobile-app/build-release.sh "更新说明"`
- 脚本功能: 更新版本号 -> 构建 -> 打包APK -> 发布到服务器
- 示例: `./mobile-app/build-release.sh "修复登录问题"`

---

## 常用脚本命令

### 后端编译重启
```bash
./rebuild.sh
```
或使用 Makefile:
```bash
make restart
```

### 其他常用命令
```bash
# 编译后端
go build -o server ./cmd/server

# 后台启动服务
nohup ./server > /tmp/backend.log 2>&1 &

# 查看日志
tail -f /tmp/backend.log

# 停止服务
pkill -f "./server"
```

## Makefile 命令

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

## 项目结构

```
/home/julei/backend/
├── cmd/server/main.go    # 入口文件
├── internal/api/         # API 模块
├── internal/config/      # 配置
├── migrations/           # 数据库迁移
├── mobile-app/           # 移动端
├── newstatic/            # 管理后台前端
├── rebuild.sh            # 编译重启脚本
└── Makefile              # Make 命令
```

## 日志文件

- 后端日志: `/tmp/backend.log`
