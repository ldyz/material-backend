# Makefile for development conveniences

# Include .env file if it exists
-include .env
export

.PHONY: dev run build test restart install-deps clean

# 颜定义变量
BINARY_NAME=server
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go' -type f)
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.CommitHash=$(COMMIT_HASH) -X main.BuildTime=$(BUILD_TIME)"

# 开发模式：自动编译并运行
dev:
	@echo "🔄 Auto-compiling and starting server..."
	@make restart

# 手动运行（使用go run）
run:
	@echo "Starting dev server with go run..."
	@go run ./cmd/server

# 编译二进制文件
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 快速编译（输出到当前目录）
build-fast:
	@echo "🔨 Fast build..."
	@go build -o $(BINARY_NAME) ./cmd/server
	@echo "✅ Fast build complete: ./$(BINARY_NAME)"

# 停止正在运行的服务
stop:
	@echo "🛑 Stopping $(BINARY_NAME)..."
	@-pkill -f "./$(BINARY_NAME)" || true
	@-killall -9 $(BINARY_NAME) 2>/dev/null || true
	@sleep 1
	@echo "✅ Server stopped"

# 重启服务：停止旧进程 -> 编译 -> 启动新进程
restart: stop build-fast run-bg
	@echo "🚀 Server restarted successfully"
	@echo ""
	@echo "📊 Server info:"
	@ps aux | grep -v grep | grep "./$(BINARY_NAME)" || echo "  No process found"

# 后台运行服务
run-bg:
	@echo "🚀 Starting $(BINARY_NAME) in background..."
	@nohup ./$(BINARY_NAME) > /tmp/backend.log 2>&1 &
	@sleep 2
	@ps aux | grep -v grep | grep "./$(BINARY_NAME)" > /dev/null && \
		echo "✅ Server started (PID: $$(pgrep -f '$(BINARY_NAME)' | head -1))" || \
		echo "❌ Failed to start server"

# 前台运行服务
run-fg:
	@echo "🚀 Starting $(BINARY_NAME) in foreground..."
	@./$(BINARY_NAME)

# 查看日志
logs:
	@echo "📋 Showing recent logs:"
	@tail -f /tmp/backend.log

# 查看最近的日志
log:
	@tail -50 /tmp/backend.log

# 安装依赖
install-deps:
	@echo "📦 Installing Go dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

# 清理构建文件
clean:
	@echo "🧹 Cleaning build files..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# 运行测试
test:
	@echo "🧪 Running tests..."
	@go test ./... -v

# 代码检查
vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

# 格式化代码
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# 注意: 数据库迁移现在使用 AutoMigrate (在 main.go 中)
# 手动 SQL 迁移脚本位于 scripts/migrations/ 目录

# 完整的CI流程
ci: build test
	@echo "✅ CI target: built and tests ran"

# 帮助信息
help:
	@echo "📖 Available targets:"
	@echo ""
	@echo "Development:"
	@echo "  make dev         - Auto-compile and restart server (default)"
	@echo "  make restart     - Stop, compile and start server in background"
	@echo "  make run         - Run server with go run"
	@echo "  make run-bg      - Start server in background"
	@echo "  make run-fg      - Start server in foreground"
	@echo "  make stop        - Stop running server"
	@echo ""
	@echo "Build:"
	@echo "  make build       - Build to bin/$(BINARY_NAME)"
	@echo "  make build-fast  - Quick build to ./$(BINARY_NAME)"
	@echo "  make clean       - Remove build artifacts"
	@echo ""
	@echo "Logs:"
	@echo "  make logs        - Follow logs in real-time"
	@echo "  make log         - Show recent logs"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt         - Format code"
	@echo "  make vet         - Run go vet"
	@echo "  make test        - Run tests"
	@echo "  make install-deps- Install Go dependencies"
	@echo ""
	@echo "Database:"
	@echo "  Database migrations now use AutoMigrate (see main.go)"
	@echo "  Manual SQL scripts are in scripts/migrations/"
